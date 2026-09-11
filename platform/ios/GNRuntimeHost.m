#import "GoNativeRenderer.h"
#import "GNViewRegistry.h"
#import "../abi/GoNativeApp.h"

@implementation GNRootViewController

- (void)viewDidLoad {
    [super viewDidLoad];
    self.view.backgroundColor = UIColor.systemBackgroundColor;
    GNResetViewRegistry(self);
    GoNativeStart();
    GoNativeSetLifecycle(1);

    NSNotificationCenter *notifications = NSNotificationCenter.defaultCenter;
    [notifications addObserver:self selector:@selector(gnDidBecomeActive) name:UIApplicationDidBecomeActiveNotification object:nil];
    [notifications addObserver:self selector:@selector(gnWillResignActive) name:UIApplicationWillResignActiveNotification object:nil];
    [notifications addObserver:self selector:@selector(gnDidEnterBackground) name:UIApplicationDidEnterBackgroundNotification object:nil];
    [notifications addObserver:self selector:@selector(gnWillEnterForeground) name:UIApplicationWillEnterForegroundNotification object:nil];
    [notifications addObserver:self selector:@selector(gnMemoryWarning) name:UIApplicationDidReceiveMemoryWarningNotification object:nil];
}

- (void)gnDidBecomeActive { GoNativeSetLifecycle(2); }
- (void)gnWillResignActive { GoNativeSetLifecycle(3); }
- (void)gnDidEnterBackground { GoNativeSetLifecycle(4); }
- (void)gnWillEnterForeground { GoNativeSetLifecycle(1); }
- (void)gnMemoryWarning { GoNativeSetLifecycle(5); }

- (void)viewDidLayoutSubviews {
    [super viewDidLayoutSubviews];
    CGRect viewport = self.view.safeAreaLayoutGuide.layoutFrame;
    GoNativeSetViewport((float)viewport.size.width,
                        (float)viewport.size.height,
                        (float)UIScreen.mainScreen.scale);
}

- (void)dealloc {
    [NSNotificationCenter.defaultCenter removeObserver:self];
    GoNativeSetLifecycle(6);
    GoNativeStop();
    GNClearViewRegistry();
}

@end
