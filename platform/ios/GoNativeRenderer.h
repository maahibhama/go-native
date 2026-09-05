#import <UIKit/UIKit.h>
#include <stdint.h>

NS_ASSUME_NONNULL_BEGIN
void GNApplyMutationBatch(const uint8_t *bytes, int32_t length);
// Synchronously measures intrinsic UIKit sizes. The caller owns the returned
// buffer and releases it with GNFreeNativeBuffer.
int32_t GNMeasureNativeBatch(const uint8_t *bytes, int32_t length, uint8_t * _Nullable * _Nonnull results, int32_t *resultLength);
void GNFreeNativeBuffer(void *buffer);
@interface GNRootViewController : UIViewController
@end
NS_ASSUME_NONNULL_END
