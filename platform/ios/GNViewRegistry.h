#import <UIKit/UIKit.h>

NS_ASSUME_NONNULL_BEGIN

// The registry owns native objects by integer node identity. Values are kept as
// id here so renderer-private action classes do not leak into the framework API.
NSMutableDictionary<NSNumber *, UIView *> *GNViewRegistryViews(void);
NSMutableDictionary<NSNumber *, id> *GNViewRegistryActions(void);
NSMutableDictionary<NSNumber *, id> *GNViewRegistryGestureActions(void);
NSMutableDictionary<NSNumber *, NSValue *> *GNViewRegistryComputedFrames(void);
NSMutableDictionary<NSNumber *, NSArray<NSLayoutConstraint *> *> *GNViewRegistryFrameConstraints(void);
UIViewController * _Nullable GNViewRegistryRoot(void);

void GNResetViewRegistry(UIViewController *root);
void GNClearViewRegistry(void);

// Private migration aliases let the renderer move to a framework-owned
// registry without mixing that structural change with protocol behavior.
#define GNViews GNViewRegistryViews()
#define GNActions GNViewRegistryActions()
#define GNGestureActions GNViewRegistryGestureActions()
#define GNComputedFrames GNViewRegistryComputedFrames()
#define GNFrameConstraints GNViewRegistryFrameConstraints()
#define GNRoot GNViewRegistryRoot()

NS_ASSUME_NONNULL_END
