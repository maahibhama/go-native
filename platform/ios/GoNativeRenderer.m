#import "GoNativeRenderer.h"
#import "counter.h"
#include <time.h>

typedef NS_ENUM(uint8_t, GNMutation) { GNCreate=1, GNDelete, GNUpdate, GNInsert, GNRemove, GNMove };
typedef NS_ENUM(uint8_t, GNNode) { GNView=1, GNText, GNButton, GNRow, GNColumn, GNSafeArea };

@interface GNSafeAreaView : UIView
@end
@implementation GNSafeAreaView
@end

@interface GNAction : NSObject
@property(nonatomic) uint64_t handler;
- (void)invoke;
@end
@implementation GNAction
- (void)invoke { GoNativeDispatchEvent(self.handler); }
@end

static NSMutableDictionary<NSNumber *,UIView *> *GNViews;
static NSMutableDictionary<NSNumber *,GNAction *> *GNActions;
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

static void GNStyle(uint64_t nodeID, UIView *view, GNNode kind, NSString *text, float width, float height, float padding, float gap, uint8_t alignment, float fontSize, BOOL bold, uint64_t handler, NSString *accessibility) {
    if ([view isKindOfClass:UILabel.class]) { UILabel*l=(UILabel*)view;l.text=text;l.font=bold?[UIFont boldSystemFontOfSize:fontSize>0?fontSize:17]:[UIFont systemFontOfSize:fontSize>0?fontSize:17]; }
    if ([view isKindOfClass:UIButton.class]) { [(UIButton*)view setTitle:text forState:UIControlStateNormal]; NSNumber*actionKey=@(nodeID);GNAction*a=GNActions[actionKey];if(!a&&handler){a=[GNAction new];GNActions[actionKey]=a;[(UIButton*)view addTarget:a action:@selector(invoke) forControlEvents:UIControlEventTouchUpInside];}a.handler=handler; }
    if ([view isKindOfClass:UIStackView.class]) { UIStackView*s=(UIStackView*)view;s.spacing=gap;s.layoutMarginsRelativeArrangement=YES;s.directionalLayoutMargins=NSDirectionalEdgeInsetsMake(padding,padding,padding,padding);s.alignment=alignment==1?UIStackViewAlignmentCenter:alignment==2?UIStackViewAlignmentTrailing:UIStackViewAlignmentLeading; }
    if(width>0)[view.widthAnchor constraintEqualToConstant:width].active=YES;if(height>0)[view.heightAnchor constraintEqualToConstant:height].active=YES;
    view.isAccessibilityElement=(kind==GNText||kind==GNButton);view.accessibilityLabel=accessibility.length?accessibility:text;
}

static UIView *GNMake(GNNode kind){UIView*v;if(kind==GNText){UILabel*l=[UILabel new];l.numberOfLines=0;v=l;}else if(kind==GNButton){UIButton*b=[UIButton buttonWithType:UIButtonTypeSystem];v=b;}else if(kind==GNRow||kind==GNColumn){UIStackView*s=[UIStackView new];s.axis=kind==GNRow?UILayoutConstraintAxisHorizontal:UILayoutConstraintAxisVertical;v=s;}else if(kind==GNSafeArea){v=[GNSafeAreaView new];}else{v=[UIView new];}v.translatesAutoresizingMaskIntoConstraints=NO;return v;}

static void GNApply(NSData *data){uint64_t started=GNNowNanos();GNReader r={(const uint8_t*)data.bytes,(const uint8_t*)data.bytes+data.length};if(u16(&r)!=2)return;uint32_t count=u32(&r);uint64_t sequence=u64(&r);for(uint32_t op=0;op<count;op++){GNMutation mutation=(GNMutation)u8(&r);GNNode kind=(GNNode)u8(&r);uint64_t nodeID=u64(&r),parentID=u64(&r);int32_t index=i32(&r),from=i32(&r);float width=f32(&r),height=f32(&r),padding=f32(&r),gap=f32(&r);uint8_t alignment=u8(&r);BOOL bold=u8(&r);float fontSize=f32(&r);uint64_t handler=u64(&r);NSString*text=str(&r);NSString*accessibility=str(&r);NSNumber*key=@(nodeID);UIView*view=GNViews[key];
    if(mutation==GNCreate){view=GNMake(kind);GNViews[key]=view;GNStyle(nodeID,view,kind,text,width,height,padding,gap,alignment,fontSize,bold,handler,accessibility);if(!GNRoot.view.subviews.count){[GNRoot.view addSubview:view];[NSLayoutConstraint activateConstraints:@[[view.leadingAnchor constraintEqualToAnchor:GNRoot.view.safeAreaLayoutGuide.leadingAnchor],[view.trailingAnchor constraintEqualToAnchor:GNRoot.view.safeAreaLayoutGuide.trailingAnchor],[view.topAnchor constraintEqualToAnchor:GNRoot.view.safeAreaLayoutGuide.topAnchor]]];}}
    else if(mutation==GNUpdate){GNStyle(nodeID,view,kind,text,width,height,padding,gap,alignment,fontSize,bold,handler,accessibility);}
    else if(mutation==GNInsert){UIView*parent=GNViews[@(parentID)];if([parent isKindOfClass:UIStackView.class]){UIStackView*s=(UIStackView*)parent;[s insertArrangedSubview:view atIndex:MIN((NSUInteger)MAX(index,0),s.arrangedSubviews.count)];}else{[parent insertSubview:view atIndex:MIN((NSUInteger)MAX(index,0),parent.subviews.count)];UILayoutGuide*guide=[parent isKindOfClass:GNSafeAreaView.class]?parent.safeAreaLayoutGuide:nil;if(guide){[NSLayoutConstraint activateConstraints:@[[view.leadingAnchor constraintEqualToAnchor:guide.leadingAnchor],[view.trailingAnchor constraintEqualToAnchor:guide.trailingAnchor],[view.topAnchor constraintEqualToAnchor:guide.topAnchor],[view.bottomAnchor constraintLessThanOrEqualToAnchor:guide.bottomAnchor]]];}}}
    else if(mutation==GNRemove){[view removeFromSuperview];}
    else if(mutation==GNMove){UIView*parent=GNViews[@(parentID)];if([parent isKindOfClass:UIStackView.class]){UIStackView*s=(UIStackView*)parent;[s removeArrangedSubview:view];[s insertArrangedSubview:view atIndex:MIN((NSUInteger)MAX(index,0),s.arrangedSubviews.count)];}}
    else if(mutation==GNDelete){[view removeFromSuperview];[GNViews removeObjectForKey:key];[GNActions removeObjectForKey:key];}
    (void)from;
}GoNativeReportBatchApplied(sequence,GNNowNanos()-started);}

void GNApplyMutationBatch(const uint8_t *bytes,int32_t length){NSData*copy=[NSData dataWithBytes:bytes length:(NSUInteger)length];dispatch_async(dispatch_get_main_queue(),^{GNApply(copy);});}

@implementation GNRootViewController
- (void)viewDidLoad {[super viewDidLoad];self.view.backgroundColor=UIColor.systemBackgroundColor;GNRoot=self;GNViews=[NSMutableDictionary dictionary];GNActions=[NSMutableDictionary dictionary];GoNativeStart();}
- (void)dealloc { GoNativeStop(); }
@end
