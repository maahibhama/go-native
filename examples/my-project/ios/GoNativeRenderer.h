#import <UIKit/UIKit.h>

@interface GNRootViewController : UIViewController
@end

void GNApplyMutationBatch(const uint8_t *bytes, int32_t length);
int32_t GNMeasureNativeBatch(const uint8_t *bytes, int32_t length, uint8_t **results, int32_t *resultLength);
void GNFreeNativeBuffer(void *buffer);
