#import "AppDelegate.h"
#import "GoNativeRenderer.h"

@implementation AppDelegate

- (BOOL)application:(UIApplication *)application didFinishLaunchingWithOptions:(NSDictionary *)launchOptions {
    self.window = [[UIWindow alloc] initWithFrame:[[UIScreen mainScreen] bounds]];
    self.window.rootViewController = [GNRootViewController new];
    [self.window makeKeyAndVisible];
    return YES;
}

@end
