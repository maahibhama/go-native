#import "GoNativeRenderer.h"
#import "counter.h"
#include <time.h>
#include <math.h>

typedef NS_ENUM(uint8_t, GNMutation) { GNCreate=1, GNDelete, GNUpdate, GNInsert, GNRemove, GNMove };
typedef NS_ENUM(uint8_t, GNNode) { GNView=1, GNText, GNButton, GNRow, GNColumn, GNSafeArea, GNTextInput, GNSwitch, GNProgressIndicator, GNImage, GNScrollView };

@interface GNSafeAreaView : UIView
@property(nonatomic) uint8_t gnAlignment;
@end
@implementation GNSafeAreaView
@end

@interface GNAction : NSObject
@property(nonatomic) uint64_t handler;
@property(nonatomic) uint64_t nodeID;
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

static NSMutableDictionary<NSNumber *,UIView *> *GNViews;
static NSMutableDictionary<NSNumber *,GNAction *> *GNActions;
static NSMutableDictionary<NSNumber *,NSArray<GNGestureAction *> *> *GNGestureActions;
static NSMutableDictionary<NSNumber *,NSValue *> *GNComputedFrames;
static NSMutableDictionary<NSNumber *,NSArray<NSLayoutConstraint *> *> *GNFrameConstraints;
static __weak GNRootViewController *GNRoot;

typedef struct { const uint8_t *p; const uint8_t *end; } GNReader;
static uint8_t u8(GNReader *r){return r->p<r->end?*r->p++:0;}
static uint16_t u16(GNReader*r){uint16_t v=0;if(r->p+2<=r->end){memcpy(&v,r->p,2);r->p+=2;}return v;}
static uint32_t u32(GNReader*r){uint32_t v=0;if(r->p+4<=r->end){memcpy(&v,r->p,4);r->p+=4;}return v;}
static uint64_t u64(GNReader*r){uint64_t v=0;if(r->p+8<=r->end){memcpy(&v,r->p,8);r->p+=8;}return v;}
static int32_t i32(GNReader*r){return (int32_t)u32(r);}
static float f32(GNReader*r){float v=0;if(r->p+4<=r->end){memcpy(&v,r->p,4);r->p+=4;}return v;}
static NSString *str(GNReader*r){uint32_t n=u32(r);if(r->p+n>r->end){r->p=r->end;return @"";}NSString*s=[[NSString alloc]initWithBytes:r->p length:n encoding:NSUTF8StringEncoding]?:@"";r->p+=n;return s;}
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

static void GNStyle(uint64_t nodeID, UIView *view, GNNode kind, NSString *text, float width, float height, float padding, float gap, uint8_t alignment, float fontSize, BOOL bold, uint64_t handler, uint64_t changeHandler, uint64_t toggleHandler, BOOL checked, float progress, NSString *accessibility, NSString *hint, uint8_t role, BOOL focused, BOOL scalesText, NSString *imageSource, uint8_t imageMode, BOOL horizontal, NSData *interactions, BOOL animate) {
    if ([view isKindOfClass:UILabel.class]) { UILabel*l=(UILabel*)view;l.text=text;UIFont*base=bold?[UIFont boldSystemFontOfSize:fontSize>0?fontSize:17]:[UIFont systemFontOfSize:fontSize>0?fontSize:17];l.font=scalesText?[[UIFontMetrics defaultMetrics] scaledFontForFont:base]:base;l.adjustsFontForContentSizeCategory=scalesText; }
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
        if(![f.text isEqualToString:text]) f.text=text;
        if(hint.length) f.placeholder = hint;
        f.font = bold?[UIFont boldSystemFontOfSize:fontSize>0?fontSize:16]:[UIFont systemFontOfSize:fontSize>0?fontSize:16];
        f.borderStyle = UITextBorderStyleRoundedRect;
        f.autocapitalizationType = UITextAutocapitalizationTypeNone;
        f.autocorrectionType = UITextAutocorrectionTypeNo;
        NSNumber*actionKey=@(nodeID);GNAction*a=GNActions[actionKey];if(!a){a=[GNAction new];a.nodeID=nodeID;GNActions[actionKey]=a;[f addTarget:a action:@selector(change:) forControlEvents:UIControlEventEditingChanged];[f addTarget:a action:@selector(focus:) forControlEvents:UIControlEventEditingDidBegin];[f addTarget:a action:@selector(blur:) forControlEvents:UIControlEventEditingDidEnd];}a.handler=changeHandler;
    }
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
    view.isAccessibilityElement=(kind==GNText||kind==GNButton||kind==GNTextInput||kind==GNSwitch||kind==GNProgressIndicator||role!=0);view.accessibilityLabel=accessibility.length?accessibility:text;view.accessibilityHint=hint;UIAccessibilityTraits traits=UIAccessibilityTraitNone;if(role==2||kind==GNButton)traits|=UIAccessibilityTraitButton;if(role==3)traits|=UIAccessibilityTraitHeader;if(role==4)traits|=UIAccessibilityTraitImage;view.accessibilityTraits=traits;if(focused){if(view.canBecomeFirstResponder&&![view isFirstResponder])[view becomeFirstResponder];UIAccessibilityPostNotification(UIAccessibilityScreenChangedNotification,view);}else if([view isFirstResponder]){[view resignFirstResponder];}
    GNConfigureInteractions(nodeID,view,interactions,animate);
}

static UIView *GNMake(GNNode kind){UIView*v;if(kind==GNText){UILabel*l=[UILabel new];l.numberOfLines=0;l.textColor=UIColor.labelColor;v=l;}else if(kind==GNButton){UIButton*b=[UIButton buttonWithType:UIButtonTypeSystem];v=b;}else if(kind==GNTextInput){UITextField*f=[UITextField new];f.borderStyle=UITextBorderStyleRoundedRect;v=f;}else if(kind==GNSwitch){v=[UISwitch new];}else if(kind==GNProgressIndicator){v=[[UIProgressView alloc]initWithProgressViewStyle:UIProgressViewStyleDefault];}else if(kind==GNImage){v=[UIImageView new];}else if(kind==GNScrollView){v=[UIScrollView new];}else if(kind==GNRow||kind==GNColumn){UIStackView*s=[UIStackView new];s.axis=kind==GNRow?UILayoutConstraintAxisHorizontal:UILayoutConstraintAxisVertical;v=s;}else if(kind==GNSafeArea){v=[GNSafeAreaView new];v.backgroundColor=UIColor.systemBackgroundColor;}else{v=[UIView new];v.backgroundColor=UIColor.systemBackgroundColor;}v.translatesAutoresizingMaskIntoConstraints=NO;return v;}

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
    uint32_t fontLength=0;memcpy(&fontLength,text+181,4);NSUInteger disabledOffset=203;uint32_t interactionFontLength=0;memcpy(&interactionFontLength,interaction+181,4);disabledOffset+=(NSUInteger)interactionFontLength;if(interaction+disabledOffset>=end)return;NSUInteger fontOffset=185+(NSUInteger)fontLength;float fontSize=0,lineHeight=0,letterSpacing=0;uint16_t fontWeight=0;memcpy(&fontSize,text+fontOffset,4);memcpy(&fontWeight,text+fontOffset+4,2);memcpy(&lineHeight,text+fontOffset+6,4);memcpy(&letterSpacing,text+fontOffset+10,4);NSString*family=[[NSString alloc]initWithBytes:text+185 length:fontLength encoding:NSUTF8StringEncoding]?:@"";UIFont*font=family.length?[UIFont fontWithName:family size:fontSize]:nil;if(!font&&fontSize>0)font=[UIFont systemFontOfSize:fontSize weight:fontWeight>=600?UIFontWeightBold:UIFontWeightRegular];if(font){if([view isKindOfClass:UILabel.class])((UILabel*)view).font=font;else if([view isKindOfClass:UIButton.class])((UIButton*)view).titleLabel.font=font;else if([view isKindOfClass:UITextField.class])((UITextField*)view).font=font;}if([view isKindOfClass:UILabel.class]&&(lineHeight>0||letterSpacing!=0)){UILabel*l=(UILabel*)view;NSMutableParagraphStyle*paragraph=[NSMutableParagraphStyle new];if(lineHeight>0){paragraph.minimumLineHeight=lineHeight;paragraph.maximumLineHeight=lineHeight;}l.attributedText=[[NSAttributedString alloc]initWithString:l.text?:@"" attributes:@{NSKernAttributeName:@(letterSpacing),NSParagraphStyleAttributeName:paragraph}];}view.userInteractionEnabled=interaction[disabledOffset]==0;
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

static void GNApply(NSData *data){uint64_t started=GNNowNanos();GNReader r={(const uint8_t*)data.bytes,(const uint8_t*)data.bytes+data.length};if(u16(&r)!=9)return;uint32_t count=u32(&r);uint64_t sequence=u64(&r);for(uint32_t op=0;op<count;op++){GNMutation mutation=(GNMutation)u8(&r);GNNode kind=(GNNode)u8(&r);uint64_t nodeID=u64(&r),parentID=u64(&r);int32_t index=i32(&r),from=i32(&r);float width=f32(&r),height=f32(&r),padding=f32(&r),gap=f32(&r);uint8_t alignment=u8(&r);BOOL bold=u8(&r);float fontSize=f32(&r);uint64_t handler=u64(&r),changeHandler=u64(&r),toggleHandler=u64(&r);BOOL checked=u8(&r);float progress=f32(&r);NSString*text=str(&r);NSString*accessibility=str(&r);NSString*hint=str(&r);uint8_t role=u8(&r);BOOL focused=u8(&r);BOOL scalesText=u8(&r);NSString*imageSource=str(&r);uint8_t imageMode=u8(&r);BOOL horizontal=u8(&r);uint32_t interactionLength=u32(&r);NSData *interactions;if(r.p+interactionLength<=r.end){interactions=[NSData dataWithBytes:r.p length:interactionLength];r.p+=interactionLength;}else{r.p=r.end;interactions=[NSData data];}uint32_t styleLength=u32(&r);if(styleLength>1048576||r.p+styleLength>r.end)return;NSData *typedStyle=[NSData dataWithBytes:r.p length:styleLength];r.p+=styleLength;if(r.p+17>r.end)return;BOOL hasFrame=u8(&r);float frameX=f32(&r),frameY=f32(&r),frameWidth=f32(&r),frameHeight=f32(&r);NSNumber*key=@(nodeID);if(hasFrame)GNComputedFrames[key]=[NSValue valueWithCGRect:CGRectMake(frameX,frameY,MAX(0,frameWidth),MAX(0,frameHeight))];else{[GNComputedFrames removeObjectForKey:key];NSArray<NSLayoutConstraint*>*old=GNFrameConstraints[key];if(old.count)[NSLayoutConstraint deactivateConstraints:old];[GNFrameConstraints removeObjectForKey:key];}UIView*view=GNViews[key];
    if(mutation==GNCreate){view=GNMake(kind);GNViews[key]=view;GNStyle(nodeID,view,kind,text,width,height,padding,gap,alignment,fontSize,bold,handler,changeHandler,toggleHandler,checked,progress,accessibility,hint,role,focused,scalesText,imageSource,imageMode,horizontal,interactions,NO);GNApplyTypedStyle(view,typedStyle);if(!GNRoot.view.subviews.count){view.backgroundColor=UIColor.systemBackgroundColor;[GNRoot.view addSubview:view];[NSLayoutConstraint activateConstraints:@[[view.leadingAnchor constraintEqualToAnchor:GNRoot.view.safeAreaLayoutGuide.leadingAnchor],[view.trailingAnchor constraintEqualToAnchor:GNRoot.view.safeAreaLayoutGuide.trailingAnchor],[view.topAnchor constraintEqualToAnchor:GNRoot.view.safeAreaLayoutGuide.topAnchor],[view.bottomAnchor constraintEqualToAnchor:GNRoot.view.safeAreaLayoutGuide.bottomAnchor]]];}}
    else if(mutation==GNUpdate){GNStyle(nodeID,view,kind,text,width,height,padding,gap,alignment,fontSize,bold,handler,changeHandler,toggleHandler,checked,progress,accessibility,hint,role,focused,scalesText,imageSource,imageMode,horizontal,interactions,YES);GNApplyTypedStyle(view,typedStyle);GNApplyComputedFrame(nodeID,view);}
    else if(mutation==GNInsert){UIView*parent=GNViews[@(parentID)];if([parent isKindOfClass:UIStackView.class]){UIStackView*s=(UIStackView*)parent;[s insertArrangedSubview:view atIndex:MIN((NSUInteger)MAX(index,0),s.arrangedSubviews.count)];}else{[parent insertSubview:view atIndex:MIN((NSUInteger)MAX(index,0),parent.subviews.count)];GNSafeAreaView*safe=[parent isKindOfClass:GNSafeAreaView.class]?(GNSafeAreaView*)parent:nil;if([parent isKindOfClass:UIScrollView.class]&&!hasFrame){UIScrollView*s=(UIScrollView*)parent;[NSLayoutConstraint activateConstraints:@[[view.leadingAnchor constraintEqualToAnchor:s.contentLayoutGuide.leadingAnchor],[view.trailingAnchor constraintEqualToAnchor:s.contentLayoutGuide.trailingAnchor],[view.topAnchor constraintEqualToAnchor:s.contentLayoutGuide.topAnchor],[view.bottomAnchor constraintEqualToAnchor:s.contentLayoutGuide.bottomAnchor],[view.widthAnchor constraintEqualToAnchor:s.frameLayoutGuide.widthAnchor]]];}else if(safe&&!hasFrame){GNConstrainSafeAreaChild(safe,view);}GNApplyComputedFrame(nodeID,view);}}
    else if(mutation==GNRemove){[view removeFromSuperview];}
    else if(mutation==GNMove){UIView*parent=GNViews[@(parentID)];if([parent isKindOfClass:UIStackView.class]){UIStackView*s=(UIStackView*)parent;[s removeArrangedSubview:view];[s insertArrangedSubview:view atIndex:MIN((NSUInteger)MAX(index,0),s.arrangedSubviews.count)];}}
    else if(mutation==GNDelete){NSArray<NSLayoutConstraint *>*frameConstraints=GNFrameConstraints[key];if(frameConstraints.count)[NSLayoutConstraint deactivateConstraints:frameConstraints];[view removeFromSuperview];[GNViews removeObjectForKey:key];[GNActions removeObjectForKey:key];[GNGestureActions removeObjectForKey:key];[GNComputedFrames removeObjectForKey:key];[GNFrameConstraints removeObjectForKey:key];}
    (void)from;
}GoNativeReportBatchApplied(sequence,GNNowNanos()-started);}

void GNApplyMutationBatch(const uint8_t *bytes,int32_t length){NSData*copy=[NSData dataWithBytes:bytes length:(NSUInteger)length];dispatch_async(dispatch_get_main_queue(),^{GNApply(copy);});}

static void GNAppendU16(NSMutableData *data,uint16_t value){[data appendBytes:&value length:sizeof(value)];}
static void GNAppendU32(NSMutableData *data,uint32_t value){[data appendBytes:&value length:sizeof(value)];}
static void GNAppendU64(NSMutableData *data,uint64_t value){[data appendBytes:&value length:sizeof(value)];}
static void GNAppendF32(NSMutableData *data,float value){[data appendBytes:&value length:sizeof(value)];}
static BOOL GNReadStringChecked(GNReader *reader,NSString **value){if(reader->p+4>reader->end)return NO;uint32_t length=u32(reader);if(length>1048576||reader->p+length>reader->end)return NO;*value=[[NSString alloc]initWithBytes:reader->p length:length encoding:NSUTF8StringEncoding]?:@"";reader->p+=length;return YES;}
static UIFont *GNMeasurementFont(NSData *typedStyle,CGFloat fallback){
    const uint8_t *record=typedStyle.bytes,*end=record+typedStyle.length;if(typedStyle.length<2)return [UIFont systemFontOfSize:fallback];uint16_t version=0;memcpy(&version,record,2);if(version!=1)return [UIFont systemFontOfSize:fallback];
    const uint8_t *style=record+2;NSUInteger styleSize=GNStyleSize(style,end);if(!styleSize)return [UIFont systemFontOfSize:fallback];uint32_t familyLength=0;memcpy(&familyLength,style+181,4);NSUInteger fontOffset=185+(NSUInteger)familyLength;if(style+fontOffset+6>end)return [UIFont systemFontOfSize:fallback];
    float fontSize=0;uint16_t weight=0;memcpy(&fontSize,style+fontOffset,4);memcpy(&weight,style+fontOffset+4,2);NSString *family=[[NSString alloc]initWithBytes:style+185 length:familyLength encoding:NSUTF8StringEncoding]?:@"";CGFloat size=fontSize>0?fontSize:fallback;UIFont *font=family.length?[UIFont fontWithName:family size:size]:nil;return font?:[UIFont systemFontOfSize:size weight:weight>=600?UIFontWeightBold:UIFontWeightRegular];
}
static CGSize GNMeasureControl(GNNode kind,NSString *text,NSString *imageSource,NSData *typedStyle,CGSize constraint){
    UIView *view=GNMake(kind);UIFont *font=GNMeasurementFont(typedStyle,kind==GNButton||kind==GNTextInput?16:17);
    if([view isKindOfClass:UILabel.class]){UILabel *label=(UILabel *)view;label.text=text;label.font=font;label.numberOfLines=0;}
    else if([view isKindOfClass:UIButton.class]){CGRect title=[(text?:@"") boundingRectWithSize:CGSizeMake(CGFLOAT_MAX,CGFLOAT_MAX) options:NSStringDrawingUsesLineFragmentOrigin|NSStringDrawingUsesFontLeading attributes:@{NSFontAttributeName:font} context:nil];return CGSizeMake(MAX(44,ceil(title.size.width)+40),MAX(44,ceil(title.size.height)+24));}
    else if([view isKindOfClass:UITextField.class]){CGRect content=[(text.length?text:@"M") boundingRectWithSize:CGSizeMake(CGFLOAT_MAX,CGFLOAT_MAX) options:NSStringDrawingUsesLineFragmentOrigin|NSStringDrawingUsesFontLeading attributes:@{NSFontAttributeName:font} context:nil];return CGSizeMake(MAX(240,ceil(content.size.width)+32),MAX(44,ceil(content.size.height)+20));}
    else if([view isKindOfClass:UIImageView.class]){UIImage *image=[UIImage imageNamed:imageSource]?:[UIImage systemImageNamed:imageSource];((UIImageView *)view).image=image;}
    CGSize size=[view sizeThatFits:constraint];if(size.width<=0||size.height<=0)size=view.intrinsicContentSize;if(size.width==UIViewNoIntrinsicMetric)size.width=0;if(size.height==UIViewNoIntrinsicMetric)size.height=0;return size;
}
int32_t GNMeasureNativeBatch(const uint8_t *bytes,int32_t length,uint8_t **results,int32_t *resultLength){
    if(!results||!resultLength){return 1;}*results=NULL;*resultLength=0;if(!bytes||length<6||length>16777216)return 2;NSData *input=[NSData dataWithBytes:bytes length:(NSUInteger)length];__block NSMutableData *output=nil;__block int32_t status=0;
    void (^measure)(void)=^{GNReader reader={(const uint8_t *)input.bytes,(const uint8_t *)input.bytes+input.length};uint16_t version=u16(&reader);uint32_t count=u32(&reader);if(version!=1||count>100000){status=3;return;}output=[NSMutableData data];GNAppendU16(output,1);GNAppendU32(output,count);
        for(uint32_t index=0;index<count;index++){if(reader.p+25>reader.end){status=4;return;}uint64_t requestID=u64(&reader);GNNode kind=(GNNode)u8(&reader);float minWidth=f32(&reader),maxWidth=f32(&reader),minHeight=f32(&reader),maxHeight=f32(&reader);NSString *text=@"",*image=@"";if(!GNReadStringChecked(&reader,&text)||!GNReadStringChecked(&reader,&image)||reader.p+4>reader.end){status=4;return;}uint32_t styleLength=u32(&reader);if(styleLength>1048576||reader.p+styleLength>reader.end){status=4;return;}NSData *style=[NSData dataWithBytes:reader.p length:styleLength];reader.p+=styleLength;CGFloat width=isfinite(maxWidth)&&maxWidth>0?maxWidth:CGFLOAT_MAX,height=isfinite(maxHeight)&&maxHeight>0?maxHeight:CGFLOAT_MAX;CGSize size=GNMeasureControl(kind,text,image,style,CGSizeMake(width,height));size.width=MAX(minWidth,MIN(size.width,width));size.height=MAX(minHeight,MIN(size.height,height));GNAppendU64(output,requestID);GNAppendF32(output,(float)size.width);GNAppendF32(output,(float)size.height);GNAppendU32(output,0);}
        if(reader.p!=reader.end||output.length>16777216){status=5;output=nil;}
    };if(NSThread.isMainThread)measure();else dispatch_sync(dispatch_get_main_queue(),measure);if(status!=0||!output)return status?:6;void *buffer=malloc(output.length);if(!buffer)return 7;memcpy(buffer,output.bytes,output.length);*results=buffer;*resultLength=(int32_t)output.length;return 0;
}
void GNFreeNativeBuffer(void *buffer){free(buffer);}

@implementation GNRootViewController
- (void)viewDidLoad {[super viewDidLoad];self.view.backgroundColor=UIColor.systemBackgroundColor;GNRoot=self;GNViews=[NSMutableDictionary dictionary];GNActions=[NSMutableDictionary dictionary];GNGestureActions=[NSMutableDictionary dictionary];GNComputedFrames=[NSMutableDictionary dictionary];GNFrameConstraints=[NSMutableDictionary dictionary];GoNativeStart();GoNativeSetLifecycle(1);[[NSNotificationCenter defaultCenter]addObserver:self selector:@selector(gnDidBecomeActive) name:UIApplicationDidBecomeActiveNotification object:nil];[[NSNotificationCenter defaultCenter]addObserver:self selector:@selector(gnWillResignActive) name:UIApplicationWillResignActiveNotification object:nil];[[NSNotificationCenter defaultCenter]addObserver:self selector:@selector(gnDidEnterBackground) name:UIApplicationDidEnterBackgroundNotification object:nil];[[NSNotificationCenter defaultCenter]addObserver:self selector:@selector(gnWillEnterForeground) name:UIApplicationWillEnterForegroundNotification object:nil];[[NSNotificationCenter defaultCenter]addObserver:self selector:@selector(gnMemoryWarning) name:UIApplicationDidReceiveMemoryWarningNotification object:nil];}
- (void)gnDidBecomeActive { GoNativeSetLifecycle(2); }
- (void)gnWillResignActive { GoNativeSetLifecycle(3); }
- (void)gnDidEnterBackground { GoNativeSetLifecycle(4); }
- (void)gnWillEnterForeground { GoNativeSetLifecycle(1); }
- (void)gnMemoryWarning { GoNativeSetLifecycle(5); }
- (void)viewDidLayoutSubviews {[super viewDidLayoutSubviews];CGRect viewport=self.view.safeAreaLayoutGuide.layoutFrame;GoNativeSetViewport((float)viewport.size.width,(float)viewport.size.height,(float)UIScreen.mainScreen.scale);}
- (void)dealloc { [[NSNotificationCenter defaultCenter]removeObserver:self];GoNativeSetLifecycle(6);GoNativeStop(); }
@end
