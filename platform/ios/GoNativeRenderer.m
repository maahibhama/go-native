#import "GoNativeRenderer.h"
#import "GNProtocolReader.h"
#import "GNViewRegistry.h"
#import "../abi/GoNativeApp.h"
#include <time.h>
#include <math.h>

typedef NS_ENUM(uint8_t, GNMutation) { GNCreate=1, GNDelete, GNUpdate, GNInsert, GNRemove, GNMove };
typedef NS_ENUM(uint8_t, GNNode) { GNView=1, GNText, GNButton, GNRow, GNColumn, GNSafeArea, GNTextInput, GNSwitch, GNProgressIndicator, GNImage, GNScrollView };

@interface GNSafeAreaView : UIView
@property(nonatomic) uint8_t gnAlignment;
@end
@implementation GNSafeAreaView
@end

@interface GNTextView : UITextView
@property(nonatomic,strong) UILabel *gnPlaceholder;
@end
@implementation GNTextView
- (instancetype)init { if((self=[super init])){_gnPlaceholder=[UILabel new];_gnPlaceholder.textColor=UIColor.placeholderTextColor;_gnPlaceholder.translatesAutoresizingMaskIntoConstraints=NO;[self addSubview:_gnPlaceholder];[NSLayoutConstraint activateConstraints:@[[_gnPlaceholder.leadingAnchor constraintEqualToAnchor:self.leadingAnchor constant:13],[_gnPlaceholder.topAnchor constraintEqualToAnchor:self.topAnchor constant:8],[_gnPlaceholder.trailingAnchor constraintLessThanOrEqualToAnchor:self.trailingAnchor constant:-8]]];}return self; }
@end

@interface GNAction : NSObject <UITextFieldDelegate, UITextViewDelegate>
@property(nonatomic) uint64_t handler;
@property(nonatomic) uint64_t nodeID;
@property(nonatomic) uint64_t submitHandler;
@property(nonatomic) uint64_t selectionHandler;
@property(nonatomic) uint64_t linkHandler;
@property(nonatomic) int32_t maxLength;
@property(nonatomic) BOOL submitOnReturn;
- (void)invoke;
- (void)change:(UITextField *)sender;
- (void)toggle:(UISwitch *)sender;
- (void)focus:(UITextField *)sender;
- (void)blur:(UITextField *)sender;
@end
@implementation GNAction
- (void)invoke { GoNativeDispatchEvent(self.handler); }
- (void)change:(UITextField *)sender { GoNativeDispatchValueEvent(self.handler, (char *)sender.text.UTF8String); }
- (void)toggle:(UISwitch *)sender { GoNativeDispatchBoolEvent(self.handler, sender.isOn ? 1 : 0); }
- (void)focus:(UITextField *)sender { (void)sender; GoNativeDispatchFocus(self.nodeID, 1); }
- (void)blur:(UITextField *)sender { (void)sender; GoNativeDispatchFocus(self.nodeID, 0); }
- (BOOL)textFieldShouldReturn:(UITextField *)sender { if(self.submitHandler)GoNativeDispatchEvent(self.submitHandler);[sender resignFirstResponder];return YES; }
- (BOOL)textField:(UITextField *)sender shouldChangeCharactersInRange:(NSRange)range replacementString:(NSString *)string { (void)sender;if(self.maxLength<=0)return YES;return (int64_t)sender.text.length-(int64_t)range.length+(int64_t)string.length<=self.maxLength; }
- (void)textFieldDidChangeSelection:(UITextField *)sender { if(!self.selectionHandler)return;UITextRange *range=sender.selectedTextRange;if(!range)return;NSInteger start=[sender offsetFromPosition:sender.beginningOfDocument toPosition:range.start],end=[sender offsetFromPosition:sender.beginningOfDocument toPosition:range.end];GoNativeDispatchSelection(self.selectionHandler,(int32_t)start,(int32_t)end); }
- (void)textViewDidChange:(UITextView *)sender { if([sender isKindOfClass:GNTextView.class])((GNTextView *)sender).gnPlaceholder.hidden=sender.text.length>0;if(self.handler)GoNativeDispatchValueEvent(self.handler,(char *)sender.text.UTF8String); }
- (void)textViewDidBeginEditing:(UITextView *)sender { (void)sender;GoNativeDispatchFocus(self.nodeID,1); }
- (void)textViewDidEndEditing:(UITextView *)sender { (void)sender;GoNativeDispatchFocus(self.nodeID,0); }
- (void)textViewDidChangeSelection:(UITextView *)sender { if(!self.selectionHandler)return;NSRange range=sender.selectedRange;GoNativeDispatchSelection(self.selectionHandler,(int32_t)range.location,(int32_t)(range.location+range.length)); }
- (BOOL)textView:(UITextView *)sender shouldChangeTextInRange:(NSRange)range replacementText:(NSString *)text { if(self.submitOnReturn&&[text isEqualToString:@"\n"]){if(self.submitHandler)GoNativeDispatchEvent(self.submitHandler);[sender resignFirstResponder];return NO;}if(self.maxLength<=0)return YES;return (int64_t)sender.text.length-(int64_t)range.length+(int64_t)text.length<=self.maxLength; }
- (BOOL)textView:(UITextView *)textView shouldInteractWithURL:(NSURL *)URL inRange:(NSRange)characterRange interaction:(UITextItemInteraction)interaction API_AVAILABLE(ios(10.0)) { (void)textView;(void)characterRange;(void)interaction;if(self.linkHandler)GoNativeDispatchValueEvent(self.linkHandler,(char *)URL.absoluteString.UTF8String);return NO; }
@end

@interface GNGestureAction : NSObject
@property(nonatomic) uint64_t handler;
@property(nonatomic) uint8_t kind;
- (void)recognize:(UIGestureRecognizer *)recognizer;
@end
@implementation GNGestureAction
- (void)recognize:(UIGestureRecognizer *)recognizer {
    if ((self.kind == 2 && recognizer.state != UIGestureRecognizerStateBegan) ||
        (self.kind == 4 && recognizer.state != UIGestureRecognizerStateEnded)) return;
    if (self.kind != 2 && self.kind != 4 && recognizer.state != UIGestureRecognizerStateRecognized) return;
    CGPoint translation=CGPointZero, velocity=CGPointZero;
    if ([recognizer isKindOfClass:UIPanGestureRecognizer.class]) {
        UIPanGestureRecognizer *pan=(UIPanGestureRecognizer *)recognizer;
        translation=[pan translationInView:pan.view]; velocity=[pan velocityInView:pan.view];
    }
    GoNativeDispatchGestureEvent(self.handler,(float)translation.x,(float)translation.y,(float)velocity.x,(float)velocity.y);
}
@end

static uint64_t GNNowNanos(void){struct timespec t;clock_gettime(CLOCK_MONOTONIC_RAW,&t);return (uint64_t)t.tv_sec*1000000000ull+(uint64_t)t.tv_nsec;}

static UIViewAnimationOptions GNAnimationOptions(uint8_t curve) {
    switch(curve){case 1:return UIViewAnimationOptionCurveEaseIn;case 2:return UIViewAnimationOptionCurveEaseOut;case 3:return UIViewAnimationOptionCurveLinear;default:return UIViewAnimationOptionCurveEaseInOut;}
}

static void GNConfigureInteractions(uint64_t nodeID, UIView *view, NSData *payload, BOOL animate) {
    for (UIGestureRecognizer *recognizer in [view.gestureRecognizers copy]) [view removeGestureRecognizer:recognizer];
    NSMutableArray<GNGestureAction *> *targets=[NSMutableArray array];
    GNReader r={(const uint8_t *)payload.bytes,(const uint8_t *)payload.bytes+payload.length};
    uint32_t gestureCount=u32(&r);
    for(uint32_t i=0;i<gestureCount && r.p<r.end;i++){
        uint8_t kind=u8(&r),direction=u8(&r);uint64_t minimumPress=u64(&r);float minimumTravel=f32(&r);uint64_t handler=u64(&r);
        GNGestureAction *target=[GNGestureAction new];target.handler=handler;target.kind=kind;
        UIGestureRecognizer *recognizer=nil;
        if(kind==1) recognizer=[[UITapGestureRecognizer alloc]initWithTarget:target action:@selector(recognize:)];
        else if(kind==2){UILongPressGestureRecognizer *longPress=[[UILongPressGestureRecognizer alloc]initWithTarget:target action:@selector(recognize:)];longPress.minimumPressDuration=(NSTimeInterval)minimumPress/1e9;recognizer=longPress;}
        else if(kind==3){UISwipeGestureRecognizer *swipe=[[UISwipeGestureRecognizer alloc]initWithTarget:target action:@selector(recognize:)];swipe.direction=direction==1?UISwipeGestureRecognizerDirectionUp:direction==2?UISwipeGestureRecognizerDirectionDown:direction==3?UISwipeGestureRecognizerDirectionLeft:direction==4?UISwipeGestureRecognizerDirectionRight:(UISwipeGestureRecognizerDirectionUp|UISwipeGestureRecognizerDirectionDown|UISwipeGestureRecognizerDirectionLeft|UISwipeGestureRecognizerDirectionRight);recognizer=swipe;}
        else if(kind==4){UIPanGestureRecognizer *pan=[[UIPanGestureRecognizer alloc]initWithTarget:target action:@selector(recognize:)];pan.minimumNumberOfTouches=1;recognizer=pan;}
        (void)minimumTravel;
        if(recognizer){[view addGestureRecognizer:recognizer];[targets addObject:target];}
    }
    GNGestureActions[@(nodeID)]=targets;
    uint32_t animationCount=u32(&r);
    for(uint32_t i=0;i<animationCount && r.p<r.end;i++){
        uint8_t property=u8(&r);int64_t duration=(int64_t)u64(&r),delay=(int64_t)u64(&r);uint8_t curve=u8(&r);float damping=f32(&r),velocity=f32(&r);BOOL reduceMotionOK=u8(&r);
        float from=f32(&r),to=f32(&r),fromX=f32(&r),fromY=f32(&r),toX=f32(&r),toY=f32(&r);
        void (^initial)(void)=^{if(property==1)view.alpha=from;else if(property==2)view.transform=CGAffineTransformMakeScale(from,from);else if(property==3)view.transform=CGAffineTransformMakeTranslation(fromX,fromY);};
        void (^final)(void)=^{if(property==1)view.alpha=to;else if(property==2)view.transform=CGAffineTransformMakeScale(to,to);else if(property==3)view.transform=CGAffineTransformMakeTranslation(toX,toY);else if(property==4)[view.superview layoutIfNeeded];};
        if(!animate) continue;
        initial();
        BOOL reduce=UIAccessibilityIsReduceMotionEnabled() && !reduceMotionOK;
        if(reduce || duration<=0){final();continue;}
        NSTimeInterval seconds=(NSTimeInterval)duration/1e9, wait=(NSTimeInterval)delay/1e9;
        if(curve==4)[UIView animateWithDuration:seconds delay:wait usingSpringWithDamping:MAX(0.01,MIN(1,damping)) initialSpringVelocity:velocity options:UIViewAnimationOptionBeginFromCurrentState animations:final completion:nil];
        else [UIView animateWithDuration:seconds delay:wait options:(GNAnimationOptions(curve)|UIViewAnimationOptionBeginFromCurrentState) animations:final completion:nil];
    }
}

static UIColor *GNColor(const uint8_t *p);
static NSLineBreakMode GNLineBreak(uint8_t wrap,uint8_t overflow){
    if(overflow==1)return NSLineBreakByTruncatingHead;if(overflow==2)return NSLineBreakByTruncatingMiddle;if(overflow==3)return NSLineBreakByTruncatingTail;return wrap==2?NSLineBreakByClipping:NSLineBreakByWordWrapping;
}
static UIKeyboardType GNKeyboardType(uint8_t kind){switch(kind){case 1:return UIKeyboardTypeEmailAddress;case 2:return UIKeyboardTypePhonePad;case 3:return UIKeyboardTypeURL;case 4:return UIKeyboardTypeNumberPad;case 5:return UIKeyboardTypeDecimalPad;case 6:return UIKeyboardTypeWebSearch;default:return UIKeyboardTypeDefault;}}
static UIReturnKeyType GNReturnKey(uint8_t value){switch(value){case 1:return UIReturnKeyDone;case 2:return UIReturnKeyGo;case 3:return UIReturnKeyNext;case 4:return UIReturnKeySearch;case 5:return UIReturnKeySend;default:return UIReturnKeyDefault;}}
static UITextAutocapitalizationType GNCapitalization(uint8_t value){switch(value){case 1:return UITextAutocapitalizationTypeSentences;case 2:return UITextAutocapitalizationTypeWords;case 3:return UITextAutocapitalizationTypeAllCharacters;default:return UITextAutocapitalizationTypeNone;}}
static UITextAutocorrectionType GNAutocorrection(uint8_t value){return value==1?UITextAutocorrectionTypeYes:value==2?UITextAutocorrectionTypeNo:UITextAutocorrectionTypeDefault;}
static void GNSetSelection(id<UITextInput> input,int32_t start,int32_t end){if(start<0||end<start)return;UITextPosition *from=[input positionFromPosition:input.beginningOfDocument offset:start],*to=[input positionFromPosition:input.beginningOfDocument offset:end];if(from&&to)input.selectedTextRange=[input textRangeFromPosition:from toPosition:to];}
static NSAttributedString *GNRichText(NSData *payload,UIFont *fallback,NSLineBreakMode lineBreak){
    if(payload.length<6)return nil;GNReader r={(const uint8_t *)payload.bytes,(const uint8_t *)payload.bytes+payload.length};if(u16(&r)!=1)return nil;uint32_t count=u32(&r);if(count>100000)return nil;NSMutableAttributedString *result=[NSMutableAttributedString new];NSMutableParagraphStyle *paragraph=[NSMutableParagraphStyle new];paragraph.lineBreakMode=lineBreak;
    for(uint32_t i=0;i<count;i++){if(r.p+4>r.end)return nil;NSString *text=str(&r);if(r.p+4>r.end)return nil;NSString *link=str(&r);if(r.p+12>r.end)return nil;float size=f32(&r);uint16_t weight=u16(&r);uint8_t flags=u8(&r),hasColor=u8(&r);const uint8_t *rgba=r.p;r.p+=4;if(r.p>r.end)return nil;UIFont *font=[UIFont systemFontOfSize:size>0?size:fallback.pointSize weight:weight>=600?UIFontWeightBold:UIFontWeightRegular];if(flags&2)font=[UIFont italicSystemFontOfSize:font.pointSize];NSMutableDictionary *attrs=[@{NSFontAttributeName:font,NSParagraphStyleAttributeName:paragraph} mutableCopy];if(flags&1)attrs[NSUnderlineStyleAttributeName]=@(NSUnderlineStyleSingle);if(hasColor)attrs[NSForegroundColorAttributeName]=GNColor(rgba);if(link.length)attrs[NSLinkAttributeName]=link;[result appendAttributedString:[[NSAttributedString alloc]initWithString:text attributes:attrs]];}
    return r.p==r.end?result:nil;
}

static void GNStyle(uint64_t nodeID, UIView *view, GNNode kind, NSString *text, float width, float height, float padding, float gap, uint8_t alignment, float fontSize, BOOL bold, uint64_t handler, uint64_t changeHandler, uint64_t toggleHandler, BOOL checked, float progress, NSString *accessibility, NSString *hint, uint8_t role, BOOL focused, BOOL scalesText, NSString *imageSource, uint8_t imageMode, BOOL horizontal, uint8_t textWrap, uint8_t textOverflow, uint32_t maxLines, BOOL selectable, NSData *richText, uint64_t linkHandler, NSString *placeholder, uint8_t inputMode, uint8_t inputKind, uint8_t returnKey, uint8_t capitalization, uint8_t autocorrect, BOOL secure, BOOL multiline, BOOL readOnly, uint8_t validation, NSString *errorText, int32_t selectionStart, int32_t selectionEnd, int32_t maxLength, uint64_t submitHandler, uint64_t selectionHandler, NSData *interactions, BOOL animate) {
    UIFont*base=bold?[UIFont boldSystemFontOfSize:fontSize>0?fontSize:17]:[UIFont systemFontOfSize:fontSize>0?fontSize:17];base=scalesText?[[UIFontMetrics defaultMetrics] scaledFontForFont:base]:base;NSLineBreakMode lineBreak=GNLineBreak(textWrap,textOverflow);
    if ([view isKindOfClass:UILabel.class]) { UILabel*l=(UILabel*)view;l.text=text;l.font=base;l.adjustsFontForContentSizeCategory=scalesText;l.numberOfLines=maxLines>0?(NSInteger)maxLines:(textWrap==2?1:0);l.lineBreakMode=lineBreak;NSAttributedString *rich=GNRichText(richText,base,lineBreak);if(rich)l.attributedText=rich; }
    if ([view isKindOfClass:UIButton.class]) {
        UIButton*b=(UIButton*)view;
        UIFont*btnFont=bold?[UIFont boldSystemFontOfSize:fontSize>0?fontSize:16]:[UIFont systemFontOfSize:fontSize>0?fontSize:16];
        if (@available(iOS 15.0, *)) {
            UIButtonConfiguration *cfg = b.configuration ?: [UIButtonConfiguration filledButtonConfiguration];
            cfg.attributedTitle = [[NSAttributedString alloc] initWithString:text?:@"" attributes:@{NSFontAttributeName:btnFont}];
            cfg.baseBackgroundColor = [UIColor systemBlueColor];
            cfg.baseForegroundColor = [UIColor whiteColor];
            cfg.cornerStyle = UIButtonConfigurationCornerStyleMedium;
            cfg.contentInsets = NSDirectionalEdgeInsetsMake(12,20,12,20);
            b.configuration = cfg;
        } else {
            [b setTitle:text forState:UIControlStateNormal];
            b.titleLabel.font = btnFont;
            b.backgroundColor = [UIColor systemBlueColor];
            [b setTitleColor:UIColor.whiteColor forState:UIControlStateNormal];
            b.layer.cornerRadius = 8.0;
            b.clipsToBounds = YES;
        }
        b.titleLabel.numberOfLines = 1;
        b.titleLabel.lineBreakMode = NSLineBreakByTruncatingTail;
        NSNumber*actionKey=@(nodeID);GNAction*a=GNActions[actionKey];if(!a&&handler){a=[GNAction new];GNActions[actionKey]=a;[b addTarget:a action:@selector(invoke) forControlEvents:UIControlEventTouchUpInside];}a.handler=handler;
    }
    if ([view isKindOfClass:UITextField.class]) {
        UITextField*f=(UITextField*)view;
        if(inputMode==0&&![f.text isEqualToString:text]) f.text=text;
        f.placeholder=placeholder.length?placeholder:nil;
        f.font = bold?[UIFont boldSystemFontOfSize:fontSize>0?fontSize:16]:[UIFont systemFontOfSize:fontSize>0?fontSize:16];
        f.borderStyle = UITextBorderStyleRoundedRect;
        f.keyboardType=GNKeyboardType(inputKind);f.returnKeyType=GNReturnKey(returnKey);f.autocapitalizationType=GNCapitalization(capitalization);f.autocorrectionType=GNAutocorrection(autocorrect);f.secureTextEntry=secure;f.enabled=!readOnly;f.clearButtonMode=inputKind==6?UITextFieldViewModeWhileEditing:UITextFieldViewModeNever;
        NSNumber*actionKey=@(nodeID);GNAction*a=GNActions[actionKey];if(!a){a=[GNAction new];a.nodeID=nodeID;GNActions[actionKey]=a;[f addTarget:a action:@selector(change:) forControlEvents:UIControlEventEditingChanged];[f addTarget:a action:@selector(focus:) forControlEvents:UIControlEventEditingDidBegin];[f addTarget:a action:@selector(blur:) forControlEvents:UIControlEventEditingDidEnd];}a.handler=changeHandler;
        a.submitHandler=submitHandler;a.selectionHandler=selectionHandler;a.maxLength=maxLength;f.delegate=a;GNSetSelection(f,selectionStart,selectionEnd);
    }
    if ([view isKindOfClass:UITextView.class]) { UITextView*v=(UITextView*)view;GNAction*a=GNActions[@(nodeID)];if(!a){a=[GNAction new];a.nodeID=nodeID;GNActions[@(nodeID)]=a;}a.handler=changeHandler;a.submitHandler=submitHandler;a.selectionHandler=selectionHandler;a.linkHandler=linkHandler;a.maxLength=maxLength;a.submitOnReturn=returnKey!=0;v.delegate=a;v.font=base;v.adjustsFontForContentSizeCategory=scalesText;v.editable=kind==GNTextInput&&!readOnly;v.selectable=selectable||kind==GNTextInput||richText.length>0;v.scrollEnabled=kind==GNTextInput&&multiline;v.dataDetectorTypes=UIDataDetectorTypeNone;v.textContainer.lineBreakMode=lineBreak;v.textContainer.maximumNumberOfLines=maxLines;v.textContainerInset=kind==GNText?UIEdgeInsetsZero:UIEdgeInsetsMake(8,8,8,8);v.textContainer.lineFragmentPadding=kind==GNText?0:5;NSAttributedString *rich=GNRichText(richText,base,lineBreak);if(rich)v.attributedText=rich;else if(inputMode==0||kind==GNText)v.text=text;if(kind==GNTextInput){v.keyboardType=GNKeyboardType(inputKind);v.returnKeyType=GNReturnKey(returnKey);v.autocapitalizationType=GNCapitalization(capitalization);v.autocorrectionType=GNAutocorrection(autocorrect);v.secureTextEntry=secure;GNSetSelection(v,selectionStart,selectionEnd);} }
    if ([view isKindOfClass:GNTextView.class]) { GNTextView*v=(GNTextView*)view;v.gnPlaceholder.text=kind==GNTextInput?placeholder:@"";v.gnPlaceholder.font=base;v.gnPlaceholder.hidden=v.text.length>0||v.gnPlaceholder.text.length==0; }
    if ([view isKindOfClass:UISwitch.class]) { UISwitch*s=(UISwitch*)view;[s setOn:checked animated:NO];NSNumber*actionKey=@(nodeID);GNAction*a=GNActions[actionKey];if(!a&&toggleHandler){a=[GNAction new];GNActions[actionKey]=a;[s addTarget:a action:@selector(toggle:) forControlEvents:UIControlEventValueChanged];}a.handler=toggleHandler; }
    if ([view isKindOfClass:UIProgressView.class]) { [(UIProgressView*)view setProgress:progress animated:NO]; }
    if ([view isKindOfClass:UIImageView.class]) {
        UIImageView*i=(UIImageView*)view;
        UIImage *img = [UIImage imageNamed:imageSource];
        if(!img && [imageSource isEqualToString:@"app_logo"]) {
            img = [UIImage systemImageNamed:@"lock.shield.fill"];
            i.tintColor = UIColor.systemBlueColor;
        } else if(!img && [imageSource isEqualToString:@"avatar"]) {
            img = [UIImage systemImageNamed:@"person.crop.circle.fill"];
            i.tintColor = UIColor.systemBlueColor;
        } else if(!img && imageSource.length) {
            img = [UIImage systemImageNamed:imageSource];
        }
        i.image=img;
        i.contentMode=imageMode==1?UIViewContentModeScaleAspectFill:imageMode==2?UIViewContentModeCenter:UIViewContentModeScaleAspectFit;
        i.clipsToBounds=YES;
    }
    if ([view isKindOfClass:UIStackView.class]) { UIStackView*s=(UIStackView*)view;s.spacing=gap;s.layoutMarginsRelativeArrangement=YES;s.directionalLayoutMargins=NSDirectionalEdgeInsetsMake(padding,padding,padding,padding);s.alignment=alignment==1?UIStackViewAlignmentCenter:alignment==2?UIStackViewAlignmentTrailing:UIStackViewAlignmentLeading; }
    if ([view isKindOfClass:GNSafeAreaView.class]) { ((GNSafeAreaView *)view).gnAlignment=alignment; }
    if(width>0)[view.widthAnchor constraintEqualToConstant:width].active=YES;if(height>0)[view.heightAnchor constraintEqualToConstant:height].active=YES;
    view.isAccessibilityElement=(kind==GNText||kind==GNButton||kind==GNTextInput||kind==GNSwitch||kind==GNProgressIndicator||role!=0);NSString *fallbackLabel=kind==GNTextInput&&placeholder.length?placeholder:text;view.accessibilityLabel=accessibility.length?accessibility:fallbackLabel;NSString *validationHint=validation==2&&errorText.length?errorText:hint;view.accessibilityHint=validationHint;UIAccessibilityTraits traits=UIAccessibilityTraitNone;if(role==2||kind==GNButton)traits|=UIAccessibilityTraitButton;if(role==3)traits|=UIAccessibilityTraitHeader;if(role==4)traits|=UIAccessibilityTraitImage;if(readOnly)traits|=UIAccessibilityTraitNotEnabled;view.accessibilityTraits=traits;if(focused){if(view.canBecomeFirstResponder&&![view isFirstResponder])[view becomeFirstResponder];UIAccessibilityPostNotification(UIAccessibilityScreenChangedNotification,view);}else if([view isFirstResponder]){[view resignFirstResponder];}
    if(validation==2){view.layer.borderColor=UIColor.systemRedColor.CGColor;view.layer.borderWidth=MAX(1,view.layer.borderWidth);}else if(validation==1){view.layer.borderColor=UIColor.systemGreenColor.CGColor;view.layer.borderWidth=MAX(1,view.layer.borderWidth);}else if(kind==GNTextInput){view.layer.borderColor=[view isKindOfClass:UITextView.class]?UIColor.separatorColor.CGColor:nil;view.layer.borderWidth=[view isKindOfClass:UITextView.class]?1:0;}
    (void)horizontal;(void)maxLength;
    GNConfigureInteractions(nodeID,view,interactions,animate);
}

static UIView *GNMake(GNNode kind,BOOL multiline,BOOL richOrSelectable){UIView*v;if(kind==GNText){if(richOrSelectable){GNTextView*t=[GNTextView new];t.editable=NO;t.scrollEnabled=NO;t.backgroundColor=UIColor.clearColor;v=t;}else{UILabel*l=[UILabel new];l.numberOfLines=0;l.textColor=UIColor.labelColor;v=l;}}else if(kind==GNButton){UIButton*b=[UIButton buttonWithType:UIButtonTypeSystem];v=b;}else if(kind==GNTextInput){if(multiline){GNTextView*t=[GNTextView new];t.layer.borderColor=UIColor.separatorColor.CGColor;t.layer.borderWidth=1;t.layer.cornerRadius=6;v=t;}else{UITextField*f=[UITextField new];f.borderStyle=UITextBorderStyleRoundedRect;v=f;}}else if(kind==GNSwitch){v=[UISwitch new];}else if(kind==GNProgressIndicator){v=[[UIProgressView alloc]initWithProgressViewStyle:UIProgressViewStyleDefault];}else if(kind==GNImage){v=[UIImageView new];}else if(kind==GNScrollView){v=[UIScrollView new];}else if(kind==GNRow||kind==GNColumn){UIStackView*s=[UIStackView new];s.axis=kind==GNRow?UILayoutConstraintAxisHorizontal:UILayoutConstraintAxisVertical;v=s;}else if(kind==GNSafeArea){v=[GNSafeAreaView new];v.backgroundColor=UIColor.systemBackgroundColor;}else{v=[UIView new];v.backgroundColor=UIColor.systemBackgroundColor;}v.translatesAutoresizingMaskIntoConstraints=NO;return v;}

static UIColor *GNColor(const uint8_t *p){return [UIColor colorWithRed:p[0]/255.0 green:p[1]/255.0 blue:p[2]/255.0 alpha:p[3]/255.0];}
static NSUInteger GNStyleSize(const uint8_t *style,const uint8_t *end){if(style+185>end)return 0;uint32_t fontLength=0;memcpy(&fontLength,style+181,4);NSUInteger size=220+(NSUInteger)fontLength;return style+size<=end?size:0;}
static BOOL GNHasTypedValues(const uint8_t *p,NSUInteger length,const uint8_t *end){if(p+length>end)return NO;for(NSUInteger i=0;i<length;i++)if(p[i]!=0)return YES;return NO;}
static void GNApplyTypedStyle(UIView *view, NSData *payload){
    if(!view||payload.length<187)return;const uint8_t*record=payload.bytes,*end=record+payload.length;uint16_t version=0;memcpy(&version,record,2);if(version!=1)return;
    const uint8_t*portable=record+2;NSUInteger portableSize=GNStyleSize(portable,end);if(!portableSize)return;const uint8_t*ios=portable+portableSize;NSUInteger iosSize=GNStyleSize(ios,end);if(!iosSize)return;
    const uint8_t*appearance=GNHasTypedValues(ios+112,69,end)?ios+112:portable+112;
    uint32_t iosFontLength=0;memcpy(&iosFontLength,ios+181,4);const uint8_t*text=GNHasTypedValues(ios+181,22+(NSUInteger)iosFontLength,end)?ios:portable;
    const uint8_t*interaction=GNHasTypedValues(ios+203+(NSUInteger)iosFontLength,17,end)?ios:portable;
    float borderWidth=0,cornerRadius=0,opacity=0;memcpy(&borderWidth,appearance+8,4);memcpy(&cornerRadius,appearance+16,4);memcpy(&opacity,appearance+44,4);
    if(appearance[3]>0)view.backgroundColor=GNColor(appearance);
    if([view isKindOfClass:UILabel.class]&&appearance[7]>0)((UILabel*)view).textColor=GNColor(appearance+4);
    if([view isKindOfClass:UIButton.class]&&appearance[7]>0)[((UIButton*)view) setTitleColor:GNColor(appearance+4) forState:UIControlStateNormal];
    if(borderWidth>0){view.layer.borderWidth=borderWidth;view.layer.borderColor=GNColor(appearance+12).CGColor;}
    if(cornerRadius>0){view.layer.cornerRadius=cornerRadius;view.clipsToBounds=YES;}
    if(opacity>0)view.alpha=MIN(1,opacity);view.hidden=appearance[68]!=0;
    float tx=0,ty=0,sx=0,sy=0,rotation=0;memcpy(&tx,appearance+48,4);memcpy(&ty,appearance+52,4);memcpy(&sx,appearance+56,4);memcpy(&sy,appearance+60,4);memcpy(&rotation,appearance+64,4);view.transform=CGAffineTransformRotate(CGAffineTransformScale(CGAffineTransformMakeTranslation(tx,ty),sx==0?1:sx,sy==0?1:sy),rotation*(CGFloat)M_PI/180.0);
    float shadowX=0,shadowY=0,shadowBlur=0,shadowOpacity=0;memcpy(&shadowX,appearance+24,4);memcpy(&shadowY,appearance+28,4);memcpy(&shadowBlur,appearance+32,4);memcpy(&shadowOpacity,appearance+40,4);if(shadowOpacity>0){view.layer.shadowColor=GNColor(appearance+20).CGColor;view.layer.shadowOffset=CGSizeMake(shadowX,shadowY);view.layer.shadowRadius=shadowBlur;view.layer.shadowOpacity=shadowOpacity;}
    uint32_t fontLength=0;memcpy(&fontLength,text+181,4);NSUInteger disabledOffset=203;uint32_t interactionFontLength=0;memcpy(&interactionFontLength,interaction+181,4);disabledOffset+=(NSUInteger)interactionFontLength;if(interaction+disabledOffset>=end)return;NSUInteger fontOffset=185+(NSUInteger)fontLength;float fontSize=0,lineHeight=0,letterSpacing=0;uint16_t fontWeight=0;memcpy(&fontSize,text+fontOffset,4);memcpy(&fontWeight,text+fontOffset+4,2);memcpy(&lineHeight,text+fontOffset+6,4);memcpy(&letterSpacing,text+fontOffset+10,4);NSString*family=[[NSString alloc]initWithBytes:text+185 length:fontLength encoding:NSUTF8StringEncoding]?:@"";UIFont*font=family.length?[UIFont fontWithName:family size:fontSize]:nil;if(!font&&fontSize>0)font=[UIFont systemFontOfSize:fontSize weight:fontWeight>=600?UIFontWeightBold:UIFontWeightRegular];if(font){if([view isKindOfClass:UILabel.class])((UILabel*)view).font=font;else if([view isKindOfClass:UIButton.class])((UIButton*)view).titleLabel.font=font;else if([view isKindOfClass:UITextField.class])((UITextField*)view).font=font;else if([view isKindOfClass:UITextView.class])((UITextView*)view).font=font;}if([view isKindOfClass:UILabel.class]&&(lineHeight>0||letterSpacing!=0)){UILabel*l=(UILabel*)view;NSMutableParagraphStyle*paragraph=[NSMutableParagraphStyle new];if(lineHeight>0){paragraph.minimumLineHeight=lineHeight;paragraph.maximumLineHeight=lineHeight;}l.attributedText=[[NSAttributedString alloc]initWithString:l.text?:@"" attributes:@{NSKernAttributeName:@(letterSpacing),NSParagraphStyleAttributeName:paragraph}];}view.userInteractionEnabled=interaction[disabledOffset]==0;
}

static void GNConstrainSafeAreaChild(GNSafeAreaView *parent, UIView *view) {
    UILayoutGuide *guide=parent.safeAreaLayoutGuide;
    NSMutableArray<NSLayoutConstraint *> *constraints=[NSMutableArray arrayWithArray:@[[view.leadingAnchor constraintGreaterThanOrEqualToAnchor:guide.leadingAnchor],[view.trailingAnchor constraintLessThanOrEqualToAnchor:guide.trailingAnchor],[view.topAnchor constraintGreaterThanOrEqualToAnchor:guide.topAnchor],[view.bottomAnchor constraintLessThanOrEqualToAnchor:guide.bottomAnchor]]];
    if(parent.gnAlignment==1){[constraints addObject:[view.centerXAnchor constraintEqualToAnchor:guide.centerXAnchor]];[constraints addObject:[view.centerYAnchor constraintEqualToAnchor:guide.centerYAnchor]];}
    else if(parent.gnAlignment==2){[constraints addObject:[view.trailingAnchor constraintEqualToAnchor:guide.trailingAnchor]];[constraints addObject:[view.bottomAnchor constraintEqualToAnchor:guide.bottomAnchor]];}
    else{[constraints addObject:[view.leadingAnchor constraintEqualToAnchor:guide.leadingAnchor]];[constraints addObject:[view.topAnchor constraintEqualToAnchor:guide.topAnchor]];}
    [NSLayoutConstraint activateConstraints:constraints];
}

static void GNApplyComputedFrame(uint64_t nodeID,UIView *view) {
    if(!view||!view.superview||view.superview==GNRoot.view)return;
    NSValue *value=GNComputedFrames[@(nodeID)];if(!value)return;CGRect frame=value.CGRectValue;
    NSArray<NSLayoutConstraint *> *old=GNFrameConstraints[@(nodeID)];if(old.count)[NSLayoutConstraint deactivateConstraints:old];
    NSMutableArray<NSLayoutConstraint *> *legacySize=[NSMutableArray array];for(NSLayoutConstraint *constraint in view.constraints)if(constraint.firstItem==view&&(constraint.firstAttribute==NSLayoutAttributeWidth||constraint.firstAttribute==NSLayoutAttributeHeight))[legacySize addObject:constraint];if(legacySize.count)[NSLayoutConstraint deactivateConstraints:legacySize];
    NSMutableArray<NSLayoutConstraint *> *constraints=[NSMutableArray array];
    if(![view.superview isKindOfClass:UIStackView.class]){[constraints addObject:[view.leadingAnchor constraintEqualToAnchor:view.superview.leadingAnchor constant:frame.origin.x]];[constraints addObject:[view.topAnchor constraintEqualToAnchor:view.superview.topAnchor constant:frame.origin.y]];}
    if(frame.size.width>=0)[constraints addObject:[view.widthAnchor constraintEqualToConstant:frame.size.width]];
    if(frame.size.height>=0)[constraints addObject:[view.heightAnchor constraintEqualToConstant:frame.size.height]];
    [NSLayoutConstraint activateConstraints:constraints];GNFrameConstraints[@(nodeID)]=constraints;
}

typedef struct { uint8_t wrap,overflow,inputMode,inputKind,returnKey,capitalization,autocorrect,validation;uint32_t maxLines;BOOL selectable,secure,multiline,readOnly;NSData *rich;uint64_t linkHandler,submitHandler,selectionHandler;NSString *placeholder,*errorText;int32_t selectionStart,selectionEnd,maxLength; } GNTextContract;
static BOOL GNReadTextContract(GNReader *r,GNTextContract *v){
    if((NSUInteger)(r->end-r->p)<7)return NO;v->wrap=u8(r);v->overflow=u8(r);v->maxLines=u32(r);v->selectable=u8(r);NSData *rich=nil;if(!GNReadBoundedData(r,&rich)||(NSUInteger)(r->end-r->p)<8)return NO;v->rich=rich;v->linkHandler=u64(r);NSString *placeholder=nil;if(!GNReadBoundedString(r,&placeholder)||(NSUInteger)(r->end-r->p)<9)return NO;v->placeholder=placeholder;v->inputMode=u8(r);v->inputKind=u8(r);v->returnKey=u8(r);v->capitalization=u8(r);v->autocorrect=u8(r);v->secure=u8(r);v->multiline=u8(r);v->readOnly=u8(r);v->validation=u8(r);NSString *errorText=nil;if(!GNReadBoundedString(r,&errorText)||(NSUInteger)(r->end-r->p)<28)return NO;v->errorText=errorText;v->selectionStart=i32(r);v->selectionEnd=i32(r);v->maxLength=i32(r);v->submitHandler=u64(r);v->selectionHandler=u64(r);return YES;
}

static void GNApply(NSData *data){uint64_t started=GNNowNanos();if(data.length<14||data.length>16777216)return;GNReader r={(const uint8_t*)data.bytes,(const uint8_t*)data.bytes+data.length};if(u16(&r)!=10)return;uint32_t count=u32(&r);if(count>100000)return;uint64_t sequence=u64(&r);for(uint32_t op=0;op<count;op++){if((NSUInteger)(r.end-r.p)<67)return;GNMutation mutation=(GNMutation)u8(&r);GNNode kind=(GNNode)u8(&r);uint64_t nodeID=u64(&r),parentID=u64(&r);int32_t index=i32(&r),from=i32(&r);float width=f32(&r),height=f32(&r),padding=f32(&r),gap=f32(&r);uint8_t alignment=u8(&r);BOOL bold=u8(&r);float fontSize=f32(&r);uint64_t handler=u64(&r),changeHandler=u64(&r),toggleHandler=u64(&r);BOOL checked=u8(&r);float progress=f32(&r);NSString*text=nil,*accessibility=nil,*hint=nil,*imageSource=nil;if(!GNReadBoundedString(&r,&text)||!GNReadBoundedString(&r,&accessibility)||!GNReadBoundedString(&r,&hint)||(NSUInteger)(r.end-r.p)<3)return;uint8_t role=u8(&r);BOOL focused=u8(&r);BOOL scalesText=u8(&r);if(!GNReadBoundedString(&r,&imageSource)||(NSUInteger)(r.end-r.p)<6)return;uint8_t imageMode=u8(&r);BOOL horizontal=u8(&r);NSData *interactions=nil;if(!GNReadBoundedData(&r,&interactions))return;GNTextContract contract={0};if(!GNReadTextContract(&r,&contract))return;NSData *typedStyle=nil;if(!GNReadBoundedData(&r,&typedStyle)||(NSUInteger)(r.end-r.p)<17)return;BOOL hasFrame=u8(&r);float frameX=f32(&r),frameY=f32(&r),frameWidth=f32(&r),frameHeight=f32(&r);NSNumber*key=@(nodeID);if(hasFrame)GNComputedFrames[key]=[NSValue valueWithCGRect:CGRectMake(frameX,frameY,MAX(0,frameWidth),MAX(0,frameHeight))];else{[GNComputedFrames removeObjectForKey:key];NSArray<NSLayoutConstraint*>*old=GNFrameConstraints[key];if(old.count)[NSLayoutConstraint deactivateConstraints:old];[GNFrameConstraints removeObjectForKey:key];}UIView*view=GNViews[key];
    if(mutation==GNCreate){view=GNMake(kind,contract.multiline,contract.selectable||contract.rich.length>0);GNViews[key]=view;GNStyle(nodeID,view,kind,text,width,height,padding,gap,alignment,fontSize,bold,handler,changeHandler,toggleHandler,checked,progress,accessibility,hint,role,focused,scalesText,imageSource,imageMode,horizontal,contract.wrap,contract.overflow,contract.maxLines,contract.selectable,contract.rich,contract.linkHandler,contract.placeholder,contract.inputMode,contract.inputKind,contract.returnKey,contract.capitalization,contract.autocorrect,contract.secure,contract.multiline,contract.readOnly,contract.validation,contract.errorText,contract.selectionStart,contract.selectionEnd,contract.maxLength,contract.submitHandler,contract.selectionHandler,interactions,NO);GNApplyTypedStyle(view,typedStyle);if(!GNRoot.view.subviews.count){view.backgroundColor=UIColor.systemBackgroundColor;[GNRoot.view addSubview:view];[NSLayoutConstraint activateConstraints:@[[view.leadingAnchor constraintEqualToAnchor:GNRoot.view.safeAreaLayoutGuide.leadingAnchor],[view.trailingAnchor constraintEqualToAnchor:GNRoot.view.safeAreaLayoutGuide.trailingAnchor],[view.topAnchor constraintEqualToAnchor:GNRoot.view.safeAreaLayoutGuide.topAnchor],[view.bottomAnchor constraintEqualToAnchor:GNRoot.view.safeAreaLayoutGuide.bottomAnchor]]];}}
    else if(mutation==GNUpdate){GNStyle(nodeID,view,kind,text,width,height,padding,gap,alignment,fontSize,bold,handler,changeHandler,toggleHandler,checked,progress,accessibility,hint,role,focused,scalesText,imageSource,imageMode,horizontal,contract.wrap,contract.overflow,contract.maxLines,contract.selectable,contract.rich,contract.linkHandler,contract.placeholder,contract.inputMode,contract.inputKind,contract.returnKey,contract.capitalization,contract.autocorrect,contract.secure,contract.multiline,contract.readOnly,contract.validation,contract.errorText,contract.selectionStart,contract.selectionEnd,contract.maxLength,contract.submitHandler,contract.selectionHandler,interactions,YES);GNApplyTypedStyle(view,typedStyle);GNApplyComputedFrame(nodeID,view);}
    else if(mutation==GNInsert){UIView*parent=GNViews[@(parentID)];if([parent isKindOfClass:UIStackView.class]){UIStackView*s=(UIStackView*)parent;[s insertArrangedSubview:view atIndex:MIN((NSUInteger)MAX(index,0),s.arrangedSubviews.count)];}else{[parent insertSubview:view atIndex:MIN((NSUInteger)MAX(index,0),parent.subviews.count)];GNSafeAreaView*safe=[parent isKindOfClass:GNSafeAreaView.class]?(GNSafeAreaView*)parent:nil;if([parent isKindOfClass:UIScrollView.class]&&!hasFrame){UIScrollView*s=(UIScrollView*)parent;[NSLayoutConstraint activateConstraints:@[[view.leadingAnchor constraintEqualToAnchor:s.contentLayoutGuide.leadingAnchor],[view.trailingAnchor constraintEqualToAnchor:s.contentLayoutGuide.trailingAnchor],[view.topAnchor constraintEqualToAnchor:s.contentLayoutGuide.topAnchor],[view.bottomAnchor constraintEqualToAnchor:s.contentLayoutGuide.bottomAnchor],[view.widthAnchor constraintEqualToAnchor:s.frameLayoutGuide.widthAnchor]]];}else if(safe&&!hasFrame){GNConstrainSafeAreaChild(safe,view);}GNApplyComputedFrame(nodeID,view);}}
    else if(mutation==GNRemove){[view removeFromSuperview];}
    else if(mutation==GNMove){UIView*parent=GNViews[@(parentID)];if([parent isKindOfClass:UIStackView.class]){UIStackView*s=(UIStackView*)parent;[s removeArrangedSubview:view];[s insertArrangedSubview:view atIndex:MIN((NSUInteger)MAX(index,0),s.arrangedSubviews.count)];}}
    else if(mutation==GNDelete){NSArray<NSLayoutConstraint *>*frameConstraints=GNFrameConstraints[key];if(frameConstraints.count)[NSLayoutConstraint deactivateConstraints:frameConstraints];[view removeFromSuperview];[GNViews removeObjectForKey:key];[GNActions removeObjectForKey:key];[GNGestureActions removeObjectForKey:key];[GNComputedFrames removeObjectForKey:key];[GNFrameConstraints removeObjectForKey:key];}
    (void)from;
}GoNativeReportBatchApplied(sequence,GNNowNanos()-started);}

void GNApplyMutationBatch(const uint8_t *bytes,int32_t length){NSData*copy=[NSData dataWithBytes:bytes length:(NSUInteger)length];dispatch_async(dispatch_get_main_queue(),^{GNApply(copy);});}
