#import <UIKit/UIKit.h>
#import "GNProtocolReader.h"
#include <stdint.h>

NS_ASSUME_NONNULL_BEGIN

typedef NS_ENUM(uint8_t, GNNode) {
    GNView = 1,
    GNText,
    GNButton,
    GNRow,
    GNColumn,
    GNSafeArea,
    GNTextInput,
    GNSwitch,
    GNProgressIndicator,
    GNImage,
    GNScrollView
};

@interface GNSafeAreaView : UIView
@property(nonatomic) uint8_t gnAlignment;
@end

@interface GNTextView : UITextView
@property(nonatomic, strong) UILabel *gnPlaceholder;
@end

@interface GNAction : NSObject <UITextFieldDelegate, UITextViewDelegate>
@property(nonatomic) uint64_t handler;
@property(nonatomic) uint64_t nodeID;
@property(nonatomic) uint64_t submitHandler;
@property(nonatomic) uint64_t selectionHandler;
@property(nonatomic) uint64_t linkHandler;
@property(nonatomic) int32_t maxLength;
@property(nonatomic) BOOL submitOnReturn;
- (void)invoke;
- (void)change:(UITextField *)sender;
- (void)toggle:(UISwitch *)sender;
- (void)focus:(UITextField *)sender;
- (void)blur:(UITextField *)sender;
@end

@interface GNGestureAction : NSObject
@property(nonatomic) uint64_t handler;
@property(nonatomic) uint8_t kind;
- (void)recognize:(UIGestureRecognizer *)recognizer;
@end

typedef struct {
    uint8_t wrap, overflow, inputMode, inputKind, returnKey, capitalization, autocorrect, validation;
    uint32_t maxLines;
    BOOL selectable, secure, multiline, readOnly;
    NSData * _Nullable rich;
    uint64_t linkHandler, submitHandler, selectionHandler;
    NSString * _Nullable placeholder;
    NSString * _Nullable errorText;
    int32_t selectionStart, selectionEnd, maxLength;
} GNTextContract;

BOOL GNReadTextContract(GNReader *r, GNTextContract *v);

UIView *GNMake(GNNode kind, BOOL multiline, BOOL richOrSelectable);

void GNStyle(uint64_t nodeID, UIView *view, GNNode kind, NSString * _Nullable text,
             float width, float height, float padding, float gap, uint8_t alignment,
             float fontSize, BOOL bold, uint64_t handler, uint64_t changeHandler,
             uint64_t toggleHandler, BOOL checked, float progress,
             NSString * _Nullable accessibility, NSString * _Nullable hint, uint8_t role,
             BOOL focused, BOOL scalesText, NSString * _Nullable imageSource, uint8_t imageMode,
             BOOL horizontal, uint8_t textWrap, uint8_t textOverflow, uint32_t maxLines,
             BOOL selectable, NSData * _Nullable richText, uint64_t linkHandler,
             NSString * _Nullable placeholder, uint8_t inputMode, uint8_t inputKind,
             uint8_t returnKey, uint8_t capitalization, uint8_t autocorrect, BOOL secure,
             BOOL multiline, BOOL readOnly, uint8_t validation, NSString * _Nullable errorText,
             int32_t selectionStart, int32_t selectionEnd, int32_t maxLength,
             uint64_t submitHandler, uint64_t selectionHandler,
             NSData * _Nullable interactions, BOOL animate);

void GNApplyTypedStyle(UIView *view, NSData * _Nullable payload);

void GNApplyComputedFrame(uint64_t nodeID, UIView *view);

void GNConstrainSafeAreaChild(GNSafeAreaView *parent, UIView *view);

NS_ASSUME_NONNULL_END
