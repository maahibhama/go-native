#import <Foundation/Foundation.h>
#include <stdint.h>

NS_ASSUME_NONNULL_BEGIN

typedef struct {
    const uint8_t *p;
    const uint8_t *end;
} GNReader;

uint8_t GNReadU8(GNReader *reader);
uint16_t GNReadU16(GNReader *reader);
uint32_t GNReadU32(GNReader *reader);
uint64_t GNReadU64(GNReader *reader);
int32_t GNReadI32(GNReader *reader);
float GNReadF32(GNReader *reader);
NSString *GNReadString(GNReader *reader);

BOOL GNReadBoundedData(GNReader *reader, NSData * _Nullable * _Nonnull value);
BOOL GNReadBoundedString(GNReader *reader, NSString * _Nullable * _Nonnull value);

// Transitional private aliases keep the renderer diff mechanical while the
// remaining decoder code is split into focused files.
#define u8 GNReadU8
#define u16 GNReadU16
#define u32 GNReadU32
#define u64 GNReadU64
#define i32 GNReadI32
#define f32 GNReadF32
#define str GNReadString

NS_ASSUME_NONNULL_END
