#import "GoNativeRenderer.h"
#import "counter.h"
#include <time.h>

typedef NS_ENUM(uint8_t, GNMutation) { GNCreate=1, GNDelete, GNUpdate, GNInsert, GNRemove, GNMove };
typedef NS_ENUM(uint8_t, GNNode) { GNView=1, GNText, GNButton, GNRow, GNColumn, GNSafeArea, GNTextInput, GNSwitch, GNProgressIndicator, GNImage, GNScrollView };

@interface GNSafeAreaView : UIView
@property(nonatomic) uint8_t gnAlignment;
@end
@implementation GNSafeAreaView
@end

@interface GNAction : NSObject
@property(nonatomic) uint64_t handler;
- (void)invoke;
- (void)change:(UITextField *)sender;
- (void)toggle:(UISwitch *)sender;
@end
@implementation GNAction
- (void)invoke { GoNativeDispatchEvent(self.handler); }
- (void)change:(UITextField *)sender { GoNativeDispatchValueEvent(self.handler, (char *)sender.text.UTF8String); }
- (void)toggle:(UISwitch *)sender { GoNativeDispatchBoolEvent(self.handler, sender.isOn ? 1 : 0); }
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
            cfg.title = text;
            cfg.baseBackgroundColor = [UIColor systemBlueColor];
            cfg.baseForegroundColor = [UIColor whiteColor];
            cfg.cornerStyle = UIButtonConfigurationCornerStyleMedium;
            b.configuration = cfg;
        } else {
            [b setTitle:text forState:UIControlStateNormal];
            b.titleLabel.font = btnFont;
            b.backgroundColor = [UIColor systemBlueColor];
            [b setTitleColor:UIColor.whiteColor forState:UIControlStateNormal];
            b.layer.cornerRadius = 8.0;
            b.clipsToBounds = YES;
        }
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
        NSNumber*actionKey=@(nodeID);GNAction*a=GNActions[actionKey];if(!a&&changeHandler){a=[GNAction new];GNActions[actionKey]=a;[f addTarget:a action:@selector(change:) forControlEvents:UIControlEventEditingChanged];}a.handler=changeHandler;
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
    view.isAccessibilityElement=(kind==GNText||kind==GNButton||kind==GNTextInput||kind==GNSwitch||kind==GNProgressIndicator||role!=0);view.accessibilityLabel=accessibility.length?accessibility:text;view.accessibilityHint=hint;UIAccessibilityTraits traits=UIAccessibilityTraitNone;if(role==2||kind==GNButton)traits|=UIAccessibilityTraitButton;if(role==3)traits|=UIAccessibilityTraitHeader;if(role==4)traits|=UIAccessibilityTraitImage;view.accessibilityTraits=traits;if(focused)UIAccessibilityPostNotification(UIAccessibilityScreenChangedNotification,view);
    GNConfigureInteractions(nodeID,view,interactions,animate);
}

static UIView *GNMake(GNNode kind){UIView*v;if(kind==GNText){UILabel*l=[UILabel new];l.numberOfLines=0;l.textColor=UIColor.labelColor;v=l;}else if(kind==GNButton){UIButton*b=[UIButton buttonWithType:UIButtonTypeSystem];v=b;}else if(kind==GNTextInput){UITextField*f=[UITextField new];f.borderStyle=UITextBorderStyleRoundedRect;v=f;}else if(kind==GNSwitch){v=[UISwitch new];}else if(kind==GNProgressIndicator){v=[[UIProgressView alloc]initWithProgressViewStyle:UIProgressViewStyleDefault];}else if(kind==GNImage){v=[UIImageView new];}else if(kind==GNScrollView){v=[UIScrollView new];}else if(kind==GNRow||kind==GNColumn){UIStackView*s=[UIStackView new];s.axis=kind==GNRow?UILayoutConstraintAxisHorizontal:UILayoutConstraintAxisVertical;v=s;}else if(kind==GNSafeArea){v=[GNSafeAreaView new];v.backgroundColor=UIColor.systemBackgroundColor;}else{v=[UIView new];v.backgroundColor=UIColor.systemBackgroundColor;}v.translatesAutoresizingMaskIntoConstraints=NO;return v;}

static UIColor *GNColor(const uint8_t *p){return [UIColor colorWithRed:p[0]/255.0 green:p[1]/255.0 blue:p[2]/255.0 alpha:p[3]/255.0];}
static void GNApplyTypedStyle(UIView *view, NSData *payload){
    if(!view||payload.length<187)return;const uint8_t*p=payload.bytes;uint16_t version=0;memcpy(&version,p,2);if(version!=1)return;
    float borderWidth=0,cornerRadius=0,opacity=0;memcpy(&borderWidth,p+122,4);memcpy(&cornerRadius,p+130,4);memcpy(&opacity,p+158,4);
    if(p[117]>0)view.backgroundColor=GNColor(p+114);
    if([view isKindOfClass:UILabel.class]&&p[121]>0)((UILabel*)view).textColor=GNColor(p+118);
    if([view isKindOfClass:UIButton.class]&&p[121]>0)[((UIButton*)view) setTitleColor:GNColor(p+118) forState:UIControlStateNormal];
    if(borderWidth>0){view.layer.borderWidth=borderWidth;view.layer.borderColor=GNColor(p+126).CGColor;}
    if(cornerRadius>0){view.layer.cornerRadius=cornerRadius;view.clipsToBounds=YES;}
    if(opacity>0)view.alpha=MIN(1,opacity);view.hidden=p[182]!=0;
    uint32_t fontLength=0;memcpy(&fontLength,p+183,4);NSUInteger disabledOffset=205+(NSUInteger)fontLength;if(disabledOffset<payload.length)view.userInteractionEnabled=p[disabledOffset]==0;
}

static void GNConstrainSafeAreaChild(GNSafeAreaView *parent, UIView *view) {
    UILayoutGuide *guide=parent.safeAreaLayoutGuide;
    NSMutableArray<NSLayoutConstraint *> *constraints=[NSMutableArray arrayWithArray:@[[view.leadingAnchor constraintGreaterThanOrEqualToAnchor:guide.leadingAnchor],[view.trailingAnchor constraintLessThanOrEqualToAnchor:guide.trailingAnchor],[view.topAnchor constraintGreaterThanOrEqualToAnchor:guide.topAnchor],[view.bottomAnchor constraintLessThanOrEqualToAnchor:guide.bottomAnchor]]];
    if(parent.gnAlignment==1){[constraints addObject:[view.centerXAnchor constraintEqualToAnchor:guide.centerXAnchor]];[constraints addObject:[view.centerYAnchor constraintEqualToAnchor:guide.centerYAnchor]];}
    else if(parent.gnAlignment==2){[constraints addObject:[view.trailingAnchor constraintEqualToAnchor:guide.trailingAnchor]];[constraints addObject:[view.bottomAnchor constraintEqualToAnchor:guide.bottomAnchor]];}
    else{[constraints addObject:[view.leadingAnchor constraintEqualToAnchor:guide.leadingAnchor]];[constraints addObject:[view.topAnchor constraintEqualToAnchor:guide.topAnchor]];}
    [NSLayoutConstraint activateConstraints:constraints];
}

static void GNApply(NSData *data){uint64_t started=GNNowNanos();GNReader r={(const uint8_t*)data.bytes,(const uint8_t*)data.bytes+data.length};if(u16(&r)!=8)return;uint32_t count=u32(&r);uint64_t sequence=u64(&r);for(uint32_t op=0;op<count;op++){GNMutation mutation=(GNMutation)u8(&r);GNNode kind=(GNNode)u8(&r);uint64_t nodeID=u64(&r),parentID=u64(&r);int32_t index=i32(&r),from=i32(&r);float width=f32(&r),height=f32(&r),padding=f32(&r),gap=f32(&r);uint8_t alignment=u8(&r);BOOL bold=u8(&r);float fontSize=f32(&r);uint64_t handler=u64(&r),changeHandler=u64(&r),toggleHandler=u64(&r);BOOL checked=u8(&r);float progress=f32(&r);NSString*text=str(&r);NSString*accessibility=str(&r);NSString*hint=str(&r);uint8_t role=u8(&r);BOOL focused=u8(&r);BOOL scalesText=u8(&r);NSString*imageSource=str(&r);uint8_t imageMode=u8(&r);BOOL horizontal=u8(&r);uint32_t interactionLength=u32(&r);NSData *interactions;if(r.p+interactionLength<=r.end){interactions=[NSData dataWithBytes:r.p length:interactionLength];r.p+=interactionLength;}else{r.p=r.end;interactions=[NSData data];}uint32_t styleLength=u32(&r);if(styleLength>1048576||r.p+styleLength>r.end)return;NSData *typedStyle=[NSData dataWithBytes:r.p length:styleLength];r.p+=styleLength;NSNumber*key=@(nodeID);UIView*view=GNViews[key];
    if(mutation==GNCreate){view=GNMake(kind);GNViews[key]=view;GNStyle(nodeID,view,kind,text,width,height,padding,gap,alignment,fontSize,bold,handler,changeHandler,toggleHandler,checked,progress,accessibility,hint,role,focused,scalesText,imageSource,imageMode,horizontal,interactions,NO);if(!GNRoot.view.subviews.count){view.backgroundColor=UIColor.systemBackgroundColor;[GNRoot.view addSubview:view];[NSLayoutConstraint activateConstraints:@[[view.leadingAnchor constraintEqualToAnchor:GNRoot.view.safeAreaLayoutGuide.leadingAnchor],[view.trailingAnchor constraintEqualToAnchor:GNRoot.view.safeAreaLayoutGuide.trailingAnchor],[view.topAnchor constraintEqualToAnchor:GNRoot.view.safeAreaLayoutGuide.topAnchor],[view.bottomAnchor constraintEqualToAnchor:GNRoot.view.safeAreaLayoutGuide.bottomAnchor]]];}}
    else if(mutation==GNUpdate){GNStyle(nodeID,view,kind,text,width,height,padding,gap,alignment,fontSize,bold,handler,changeHandler,toggleHandler,checked,progress,accessibility,hint,role,focused,scalesText,imageSource,imageMode,horizontal,interactions,YES);}
    else if(mutation==GNInsert){UIView*parent=GNViews[@(parentID)];if([parent isKindOfClass:UIStackView.class]){UIStackView*s=(UIStackView*)parent;[s insertArrangedSubview:view atIndex:MIN((NSUInteger)MAX(index,0),s.arrangedSubviews.count)];}else{[parent insertSubview:view atIndex:MIN((NSUInteger)MAX(index,0),parent.subviews.count)];GNSafeAreaView*safe=[parent isKindOfClass:GNSafeAreaView.class]?(GNSafeAreaView*)parent:nil;if([parent isKindOfClass:UIScrollView.class]){UIScrollView*s=(UIScrollView*)parent;[NSLayoutConstraint activateConstraints:@[[view.leadingAnchor constraintEqualToAnchor:s.contentLayoutGuide.leadingAnchor],[view.trailingAnchor constraintEqualToAnchor:s.contentLayoutGuide.trailingAnchor],[view.topAnchor constraintEqualToAnchor:s.contentLayoutGuide.topAnchor],[view.bottomAnchor constraintEqualToAnchor:s.contentLayoutGuide.bottomAnchor],[view.widthAnchor constraintEqualToAnchor:s.frameLayoutGuide.widthAnchor]]];}else if(safe){GNConstrainSafeAreaChild(safe,view);}}}
    else if(mutation==GNRemove){[view removeFromSuperview];}
    else if(mutation==GNMove){UIView*parent=GNViews[@(parentID)];if([parent isKindOfClass:UIStackView.class]){UIStackView*s=(UIStackView*)parent;[s removeArrangedSubview:view];[s insertArrangedSubview:view atIndex:MIN((NSUInteger)MAX(index,0),s.arrangedSubviews.count)];}}
    else if(mutation==GNDelete){[view removeFromSuperview];[GNViews removeObjectForKey:key];[GNActions removeObjectForKey:key];[GNGestureActions removeObjectForKey:key];}
    (void)from;
}GoNativeReportBatchApplied(sequence,GNNowNanos()-started);}

void GNApplyMutationBatch(const uint8_t *bytes,int32_t length){NSData*copy=[NSData dataWithBytes:bytes length:(NSUInteger)length];dispatch_async(dispatch_get_main_queue(),^{GNApply(copy);});}

@implementation GNRootViewController
- (void)viewDidLoad {[super viewDidLoad];self.view.backgroundColor=UIColor.systemBackgroundColor;GNRoot=self;GNViews=[NSMutableDictionary dictionary];GNActions=[NSMutableDictionary dictionary];GNGestureActions=[NSMutableDictionary dictionary];GoNativeStart();}
- (void)dealloc { GoNativeStop(); }
@end
