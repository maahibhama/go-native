#import "GNViewRegistry.h"

static NSMutableDictionary<NSNumber *, UIView *> *views;
static NSMutableDictionary<NSNumber *, id> *actions;
static NSMutableDictionary<NSNumber *, id> *gestureActions;
static NSMutableDictionary<NSNumber *, NSValue *> *computedFrames;
static NSMutableDictionary<NSNumber *, NSArray<NSLayoutConstraint *> *> *frameConstraints;
static __weak UIViewController *rootViewController;

NSMutableDictionary<NSNumber *, UIView *> *GNViewRegistryViews(void) { return views; }
NSMutableDictionary<NSNumber *, id> *GNViewRegistryActions(void) { return actions; }
NSMutableDictionary<NSNumber *, id> *GNViewRegistryGestureActions(void) { return gestureActions; }
NSMutableDictionary<NSNumber *, NSValue *> *GNViewRegistryComputedFrames(void) { return computedFrames; }
NSMutableDictionary<NSNumber *, NSArray<NSLayoutConstraint *> *> *GNViewRegistryFrameConstraints(void) { return frameConstraints; }
UIViewController *GNViewRegistryRoot(void) { return rootViewController; }

void GNResetViewRegistry(UIViewController *root) {
    rootViewController = root;
    views = [NSMutableDictionary dictionary];
    actions = [NSMutableDictionary dictionary];
    gestureActions = [NSMutableDictionary dictionary];
    computedFrames = [NSMutableDictionary dictionary];
    frameConstraints = [NSMutableDictionary dictionary];
}

void GNClearViewRegistry(void) {
    for (NSArray<NSLayoutConstraint *> *constraints in frameConstraints.allValues) {
        [NSLayoutConstraint deactivateConstraints:constraints];
    }
    [views removeAllObjects];
    [actions removeAllObjects];
    [gestureActions removeAllObjects];
    [computedFrames removeAllObjects];
    [frameConstraints removeAllObjects];
    rootViewController = nil;
}
