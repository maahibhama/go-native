#import "GNProtocolReader.h"

#include <string.h>

static const uint32_t GNMaximumFieldLength = 1024 * 1024;

uint8_t GNReadU8(GNReader *reader) {
    return reader->p < reader->end ? *reader->p++ : 0;
}

uint16_t GNReadU16(GNReader *reader) {
    uint16_t value = 0;
    if ((size_t)(reader->end - reader->p) >= sizeof(value)) {
        memcpy(&value, reader->p, sizeof(value));
        reader->p += sizeof(value);
    }
    return value;
}

uint32_t GNReadU32(GNReader *reader) {
    uint32_t value = 0;
    if ((size_t)(reader->end - reader->p) >= sizeof(value)) {
        memcpy(&value, reader->p, sizeof(value));
        reader->p += sizeof(value);
    }
    return value;
}

uint64_t GNReadU64(GNReader *reader) {
    uint64_t value = 0;
    if ((size_t)(reader->end - reader->p) >= sizeof(value)) {
        memcpy(&value, reader->p, sizeof(value));
        reader->p += sizeof(value);
    }
    return value;
}

int32_t GNReadI32(GNReader *reader) {
    return (int32_t)GNReadU32(reader);
}

float GNReadF32(GNReader *reader) {
    float value = 0;
    if ((size_t)(reader->end - reader->p) >= sizeof(value)) {
        memcpy(&value, reader->p, sizeof(value));
        reader->p += sizeof(value);
    }
    return value;
}

NSString *GNReadString(GNReader *reader) {
    uint32_t length = GNReadU32(reader);
    if ((size_t)(reader->end - reader->p) < length) {
        reader->p = reader->end;
        return @"";
    }
    NSString *value = [[NSString alloc] initWithBytes:reader->p
                                               length:length
                                             encoding:NSUTF8StringEncoding] ?: @"";
    reader->p += length;
    return value;
}

BOOL GNReadBoundedData(GNReader *reader, NSData **value) {
    if ((size_t)(reader->end - reader->p) < sizeof(uint32_t)) return NO;
    uint32_t length = GNReadU32(reader);
    if (length > GNMaximumFieldLength || (size_t)(reader->end - reader->p) < length) return NO;
    *value = [NSData dataWithBytes:reader->p length:length];
    reader->p += length;
    return YES;
}

BOOL GNReadBoundedString(GNReader *reader, NSString **value) {
    NSData *data = nil;
    if (!GNReadBoundedData(reader, &data)) return NO;
    *value = [[NSString alloc] initWithData:data encoding:NSUTF8StringEncoding] ?: @"";
    return YES;
}
