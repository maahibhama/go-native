#import "GoNativeRenderer.h"
#import "GNControls.h"
#import "GNProtocolReader.h"
#import "GNViewRegistry.h"
#import "../abi/GoNativeApp.h"
#include <time.h>

typedef NS_ENUM(uint8_t, GNMutation) { GNCreate=1, GNDelete, GNUpdate, GNInsert, GNRemove, GNMove };

static uint64_t GNNowNanos(void){
    struct timespec t;
    clock_gettime(CLOCK_MONOTONIC_RAW,&t);
    return (uint64_t)t.tv_sec*1000000000ull+(uint64_t)t.tv_nsec;
}

static void GNApply(NSData *data){
    uint64_t started=GNNowNanos();
    if(data.length<14||data.length>16777216)return;
    GNReader r={(const uint8_t*)data.bytes,(const uint8_t*)data.bytes+data.length};
    if(u16(&r)!=10)return;
    uint32_t count=u32(&r);
    if(count>100000)return;
    uint64_t sequence=u64(&r);
    for(uint32_t op=0;op<count;op++){
        if((NSUInteger)(r.end-r.p)<67)return;
        GNMutation mutation=(GNMutation)u8(&r);
        GNNode kind=(GNNode)u8(&r);
        uint64_t nodeID=u64(&r),parentID=u64(&r);
        int32_t index=i32(&r),from=i32(&r);
        float width=f32(&r),height=f32(&r),padding=f32(&r),gap=f32(&r);
        uint8_t alignment=u8(&r);
        BOOL bold=u8(&r);
        float fontSize=f32(&r);
        uint64_t handler=u64(&r),changeHandler=u64(&r),toggleHandler=u64(&r);
        BOOL checked=u8(&r);
        float progress=f32(&r);
        NSString*text=nil,*accessibility=nil,*hint=nil,*imageSource=nil;
        if(!GNReadBoundedString(&r,&text)||!GNReadBoundedString(&r,&accessibility)||!GNReadBoundedString(&r,&hint)||(NSUInteger)(r.end-r.p)<3)return;
        uint8_t role=u8(&r);
        BOOL focused=u8(&r);
        BOOL scalesText=u8(&r);
        if(!GNReadBoundedString(&r,&imageSource)||(NSUInteger)(r.end-r.p)<6)return;
        uint8_t imageMode=u8(&r);
        BOOL horizontal=u8(&r);
        NSData *interactions=nil;
        if(!GNReadBoundedData(&r,&interactions))return;
        GNTextContract contract={0};
        if(!GNReadTextContract(&r,&contract))return;
        NSData *typedStyle=nil;
        if(!GNReadBoundedData(&r,&typedStyle)||(NSUInteger)(r.end-r.p)<17)return;
        BOOL hasFrame=u8(&r);
        float frameX=f32(&r),frameY=f32(&r),frameWidth=f32(&r),frameHeight=f32(&r);
        NSNumber*key=@(nodeID);
        if(hasFrame)GNComputedFrames[key]=[NSValue valueWithCGRect:CGRectMake(frameX,frameY,MAX(0,frameWidth),MAX(0,frameHeight))];
        else{
            [GNComputedFrames removeObjectForKey:key];
            NSArray<NSLayoutConstraint*>*old=GNFrameConstraints[key];
            if(old.count)[NSLayoutConstraint deactivateConstraints:old];
            [GNFrameConstraints removeObjectForKey:key];
        }
        UIView*view=GNViews[key];
        if(mutation==GNCreate){
            view=GNMake(kind,contract.multiline,contract.selectable||contract.rich.length>0);
            GNViews[key]=view;
            GNStyle(nodeID,view,kind,text,width,height,padding,gap,alignment,fontSize,bold,handler,changeHandler,toggleHandler,checked,progress,accessibility,hint,role,focused,scalesText,imageSource,imageMode,horizontal,contract.wrap,contract.overflow,contract.maxLines,contract.selectable,contract.rich,contract.linkHandler,contract.placeholder,contract.inputMode,contract.inputKind,contract.returnKey,contract.capitalization,contract.autocorrect,contract.secure,contract.multiline,contract.readOnly,contract.validation,contract.errorText,contract.selectionStart,contract.selectionEnd,contract.maxLength,contract.submitHandler,contract.selectionHandler,interactions,NO);
            GNApplyTypedStyle(view,typedStyle);
            if(!GNRoot.view.subviews.count){
                view.backgroundColor=UIColor.systemBackgroundColor;
                [GNRoot.view addSubview:view];
                [NSLayoutConstraint activateConstraints:@[
                    [view.leadingAnchor constraintEqualToAnchor:GNRoot.view.safeAreaLayoutGuide.leadingAnchor],
                    [view.trailingAnchor constraintEqualToAnchor:GNRoot.view.safeAreaLayoutGuide.trailingAnchor],
                    [view.topAnchor constraintEqualToAnchor:GNRoot.view.safeAreaLayoutGuide.topAnchor],
                    [view.bottomAnchor constraintEqualToAnchor:GNRoot.view.safeAreaLayoutGuide.bottomAnchor]
                ]];
            }
        }
        else if(mutation==GNUpdate){
            GNStyle(nodeID,view,kind,text,width,height,padding,gap,alignment,fontSize,bold,handler,changeHandler,toggleHandler,checked,progress,accessibility,hint,role,focused,scalesText,imageSource,imageMode,horizontal,contract.wrap,contract.overflow,contract.maxLines,contract.selectable,contract.rich,contract.linkHandler,contract.placeholder,contract.inputMode,contract.inputKind,contract.returnKey,contract.capitalization,contract.autocorrect,contract.secure,contract.multiline,contract.readOnly,contract.validation,contract.errorText,contract.selectionStart,contract.selectionEnd,contract.maxLength,contract.submitHandler,contract.selectionHandler,interactions,YES);
            GNApplyTypedStyle(view,typedStyle);
            GNApplyComputedFrame(nodeID,view);
        }
        else if(mutation==GNInsert){
            UIView*parent=GNViews[@(parentID)];
            if([parent isKindOfClass:UIStackView.class]){
                UIStackView*s=(UIStackView*)parent;
                [s insertArrangedSubview:view atIndex:MIN((NSUInteger)MAX(index,0),s.arrangedSubviews.count)];
            }else{
                [parent insertSubview:view atIndex:MIN((NSUInteger)MAX(index,0),parent.subviews.count)];
                GNSafeAreaView*safe=[parent isKindOfClass:GNSafeAreaView.class]?(GNSafeAreaView*)parent:nil;
                if([parent isKindOfClass:UIScrollView.class]&&!hasFrame){
                    UIScrollView*s=(UIScrollView*)parent;
                    [NSLayoutConstraint activateConstraints:@[
                        [view.leadingAnchor constraintEqualToAnchor:s.contentLayoutGuide.leadingAnchor],
                        [view.trailingAnchor constraintEqualToAnchor:s.contentLayoutGuide.trailingAnchor],
                        [view.topAnchor constraintEqualToAnchor:s.contentLayoutGuide.topAnchor],
                        [view.bottomAnchor constraintEqualToAnchor:s.contentLayoutGuide.bottomAnchor],
                        [view.widthAnchor constraintEqualToAnchor:s.frameLayoutGuide.widthAnchor]
                    ]];
                }else if(safe&&!hasFrame){
                    GNConstrainSafeAreaChild(safe,view);
                }
                GNApplyComputedFrame(nodeID,view);
            }
        }
        else if(mutation==GNRemove){
            [view removeFromSuperview];
        }
        else if(mutation==GNMove){
            UIView*parent=GNViews[@(parentID)];
            if([parent isKindOfClass:UIStackView.class]){
                UIStackView*s=(UIStackView*)parent;
                [s removeArrangedSubview:view];
                [s insertArrangedSubview:view atIndex:MIN((NSUInteger)MAX(index,0),s.arrangedSubviews.count)];
            }
        }
        else if(mutation==GNDelete){
            NSArray<NSLayoutConstraint *>*frameConstraints=GNFrameConstraints[key];
            if(frameConstraints.count)[NSLayoutConstraint deactivateConstraints:frameConstraints];
            [view removeFromSuperview];
            [GNViews removeObjectForKey:key];
            [GNActions removeObjectForKey:key];
            [GNGestureActions removeObjectForKey:key];
            [GNComputedFrames removeObjectForKey:key];
            [GNFrameConstraints removeObjectForKey:key];
        }
        (void)from;
    }
    GoNativeReportBatchApplied(sequence,GNNowNanos()-started);
}

void GNApplyMutationBatch(const uint8_t *bytes,int32_t length){
    NSData*copy=[NSData dataWithBytes:bytes length:(NSUInteger)length];
    dispatch_async(dispatch_get_main_queue(),^{GNApply(copy);});
}
