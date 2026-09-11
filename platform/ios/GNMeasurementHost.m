#import "GNMeasurementHost.h"
#import "GNProtocolReader.h"

#import <UIKit/UIKit.h>
#include <math.h>

typedef NS_ENUM(uint8_t, GNMeasurementNode) {
    GNMeasurementView=1, GNMeasurementText, GNMeasurementButton, GNMeasurementRow,
    GNMeasurementColumn, GNMeasurementSafeArea, GNMeasurementTextInput,
    GNMeasurementSwitch, GNMeasurementProgressIndicator, GNMeasurementImage,
    GNMeasurementScrollView
};

static NSLineBreakMode GNMeasurementLineBreak(uint8_t wrap,uint8_t overflow){
    if(overflow==1)return NSLineBreakByTruncatingHead;if(overflow==2)return NSLineBreakByTruncatingMiddle;if(overflow==3)return NSLineBreakByTruncatingTail;return wrap==2?NSLineBreakByClipping:NSLineBreakByWordWrapping;
}
static UIKeyboardType GNMeasurementKeyboardType(uint8_t kind){switch(kind){case 1:return UIKeyboardTypeEmailAddress;case 2:return UIKeyboardTypePhonePad;case 3:return UIKeyboardTypeURL;case 4:return UIKeyboardTypeNumberPad;case 5:return UIKeyboardTypeDecimalPad;case 6:return UIKeyboardTypeWebSearch;default:return UIKeyboardTypeDefault;}}
static UIColor *GNMeasurementColor(const uint8_t *p){return [UIColor colorWithRed:p[0]/255.0 green:p[1]/255.0 blue:p[2]/255.0 alpha:p[3]/255.0];}
static NSUInteger GNMeasurementStyleSize(const uint8_t *style,const uint8_t *end){if(style+185>end)return 0;uint32_t fontLength=0;memcpy(&fontLength,style+181,4);NSUInteger size=220+(NSUInteger)fontLength;return style+size<=end?size:0;}
static UIView *GNMeasurementMakeView(GNMeasurementNode kind,BOOL multiline,BOOL richOrSelectable){
    if(kind==GNMeasurementText){if(richOrSelectable){UITextView *view=[UITextView new];view.editable=NO;view.scrollEnabled=NO;view.backgroundColor=UIColor.clearColor;return view;}UILabel *label=[UILabel new];label.numberOfLines=0;label.textColor=UIColor.labelColor;return label;}
    if(kind==GNMeasurementButton)return [UIButton buttonWithType:UIButtonTypeSystem];
    if(kind==GNMeasurementTextInput)return multiline?[UITextView new]:[UITextField new];
    if(kind==GNMeasurementSwitch)return [UISwitch new];
    if(kind==GNMeasurementProgressIndicator)return [[UIProgressView alloc]initWithProgressViewStyle:UIProgressViewStyleDefault];
    if(kind==GNMeasurementImage)return [UIImageView new];
    if(kind==GNMeasurementScrollView)return [UIScrollView new];
    if(kind==GNMeasurementRow||kind==GNMeasurementColumn){UIStackView *stack=[UIStackView new];stack.axis=kind==GNMeasurementRow?UILayoutConstraintAxisHorizontal:UILayoutConstraintAxisVertical;return stack;}
    return [UIView new];
}
static NSAttributedString *GNMeasurementRichText(NSData *payload,UIFont *fallback,NSLineBreakMode lineBreak){
    if(payload.length<6)return nil;GNReader r={(const uint8_t *)payload.bytes,(const uint8_t *)payload.bytes+payload.length};if(u16(&r)!=1)return nil;uint32_t count=u32(&r);if(count>100000)return nil;NSMutableAttributedString *result=[NSMutableAttributedString new];NSMutableParagraphStyle *paragraph=[NSMutableParagraphStyle new];paragraph.lineBreakMode=lineBreak;
    for(uint32_t i=0;i<count;i++){if(r.p+4>r.end)return nil;NSString *text=str(&r);if(r.p+4>r.end)return nil;NSString *link=str(&r);if(r.p+12>r.end)return nil;float size=f32(&r);uint16_t weight=u16(&r);uint8_t flags=u8(&r),hasColor=u8(&r);const uint8_t *rgba=r.p;r.p+=4;if(r.p>r.end)return nil;UIFont *font=[UIFont systemFontOfSize:size>0?size:fallback.pointSize weight:weight>=600?UIFontWeightBold:UIFontWeightRegular];if(flags&2)font=[UIFont italicSystemFontOfSize:font.pointSize];NSMutableDictionary *attrs=[@{NSFontAttributeName:font,NSParagraphStyleAttributeName:paragraph} mutableCopy];if(flags&1)attrs[NSUnderlineStyleAttributeName]=@(NSUnderlineStyleSingle);if(hasColor)attrs[NSForegroundColorAttributeName]=GNMeasurementColor(rgba);if(link.length)attrs[NSLinkAttributeName]=link;[result appendAttributedString:[[NSAttributedString alloc]initWithString:text attributes:attrs]];}
    return r.p==r.end?result:nil;
}
static UIFont *GNMeasurementFont(NSData *typedStyle,CGFloat fallback){
    const uint8_t *record=typedStyle.bytes,*end=record+typedStyle.length;if(typedStyle.length<2)return [UIFont systemFontOfSize:fallback];uint16_t version=0;memcpy(&version,record,2);if(version!=1)return [UIFont systemFontOfSize:fallback];
    const uint8_t *style=record+2;NSUInteger styleSize=GNMeasurementStyleSize(style,end);if(!styleSize)return [UIFont systemFontOfSize:fallback];uint32_t familyLength=0;memcpy(&familyLength,style+181,4);NSUInteger fontOffset=185+(NSUInteger)familyLength;if(style+fontOffset+6>end)return [UIFont systemFontOfSize:fallback];
    float fontSize=0;uint16_t weight=0;memcpy(&fontSize,style+fontOffset,4);memcpy(&weight,style+fontOffset+4,2);NSString *family=[[NSString alloc]initWithBytes:style+185 length:familyLength encoding:NSUTF8StringEncoding]?:@"";CGFloat size=fontSize>0?fontSize:fallback;UIFont *font=family.length?[UIFont fontWithName:family size:size]:nil;return font?:[UIFont systemFontOfSize:size weight:weight>=600?UIFontWeightBold:UIFontWeightRegular];
}
static CGSize GNMeasureControl(GNMeasurementNode kind,NSString *text,NSString *imageSource,NSData *typedStyle,CGSize constraint,uint8_t wrap,uint8_t overflow,uint32_t maxLines,BOOL selectable,NSData *richText,NSString *placeholder,uint8_t inputKind,BOOL secure,BOOL multiline,int32_t maxLength){
    UIView *view=GNMeasurementMakeView(kind,multiline,selectable||richText.length>0);UIFont *font=GNMeasurementFont(typedStyle,kind==GNMeasurementButton||kind==GNMeasurementTextInput?16:17);NSLineBreakMode lineBreak=GNMeasurementLineBreak(wrap,overflow);
    if([view isKindOfClass:UILabel.class]){UILabel *label=(UILabel *)view;label.text=text;label.font=font;label.numberOfLines=maxLines>0?(NSInteger)maxLines:(wrap==2?1:0);label.lineBreakMode=lineBreak;NSAttributedString *rich=GNMeasurementRichText(richText,font,lineBreak);if(rich)label.attributedText=rich;}
    else if([view isKindOfClass:UIButton.class]){CGRect title=[(text?:@"") boundingRectWithSize:CGSizeMake(CGFLOAT_MAX,CGFLOAT_MAX) options:NSStringDrawingUsesLineFragmentOrigin|NSStringDrawingUsesFontLeading attributes:@{NSFontAttributeName:font} context:nil];return CGSizeMake(MAX(44,ceil(title.size.width)+40),MAX(44,ceil(title.size.height)+24));}
    else if([view isKindOfClass:UITextField.class]){UITextField *field=(UITextField *)view;field.font=font;field.placeholder=placeholder;field.keyboardType=GNMeasurementKeyboardType(inputKind);field.secureTextEntry=secure;NSString *content=text.length?text:(placeholder.length?placeholder:@"M");CGRect bounds=[content boundingRectWithSize:CGSizeMake(CGFLOAT_MAX,CGFLOAT_MAX) options:NSStringDrawingUsesLineFragmentOrigin|NSStringDrawingUsesFontLeading attributes:@{NSFontAttributeName:font} context:nil];return CGSizeMake(MAX(240,ceil(bounds.size.width)+32),MAX(44,ceil(bounds.size.height)+20));}
    else if([view isKindOfClass:UITextView.class]){UITextView *input=(UITextView *)view;input.font=font;input.text=text.length?text:placeholder;input.textContainer.lineBreakMode=lineBreak;input.textContainer.maximumNumberOfLines=maxLines;NSAttributedString *rich=GNMeasurementRichText(richText,font,lineBreak);if(rich)input.attributedText=rich;CGSize size=[input sizeThatFits:constraint];return CGSizeMake(MAX(kind==GNMeasurementTextInput?240:0,ceil(size.width)),MAX(kind==GNMeasurementTextInput?88:0,ceil(size.height)));}
    else if([view isKindOfClass:UIImageView.class]){UIImage *image=[UIImage imageNamed:imageSource]?:[UIImage systemImageNamed:imageSource];((UIImageView *)view).image=image;}
    (void)maxLength;CGSize size=[view sizeThatFits:constraint];if(size.width<=0||size.height<=0)size=view.intrinsicContentSize;if(size.width==UIViewNoIntrinsicMetric)size.width=0;if(size.height==UIViewNoIntrinsicMetric)size.height=0;return size;
}
static void GNAppendU16(NSMutableData *data,uint16_t value){[data appendBytes:&value length:sizeof(value)];}
static void GNAppendU32(NSMutableData *data,uint32_t value){[data appendBytes:&value length:sizeof(value)];}
static void GNAppendU64(NSMutableData *data,uint64_t value){[data appendBytes:&value length:sizeof(value)];}
static void GNAppendF32(NSMutableData *data,float value){[data appendBytes:&value length:sizeof(value)];}

int32_t GNMeasureNativeBatch(const uint8_t *bytes,int32_t length,uint8_t **results,int32_t *resultLength){
    if(!results||!resultLength){return 1;}*results=NULL;*resultLength=0;if(!bytes||length<6||length>16777216)return 2;NSData *input=[NSData dataWithBytes:bytes length:(NSUInteger)length];__block NSMutableData *output=nil;__block int32_t status=0;
    void (^measure)(void)=^{GNReader reader={(const uint8_t *)input.bytes,(const uint8_t *)input.bytes+input.length};uint16_t version=u16(&reader);uint32_t count=u32(&reader);if(version!=2||count>100000){status=3;return;}output=[NSMutableData data];GNAppendU16(output,2);GNAppendU32(output,count);
        for(uint32_t index=0;index<count;index++){if(reader.p+25>reader.end){status=4;return;}uint64_t requestID=u64(&reader);GNMeasurementNode kind=(GNMeasurementNode)u8(&reader);float minWidth=f32(&reader),maxWidth=f32(&reader),minHeight=f32(&reader),maxHeight=f32(&reader);NSString *text=@"",*image=@"",*placeholder=@"";NSData *rich=nil;if(!GNReadBoundedString(&reader,&text)||!GNReadBoundedString(&reader,&image)||(NSUInteger)(reader.end-reader.p)<7){status=4;return;}uint8_t wrap=u8(&reader),overflow=u8(&reader);uint32_t maxLines=u32(&reader);BOOL selectable=u8(&reader);if(!GNReadBoundedData(&reader,&rich)||!GNReadBoundedString(&reader,&placeholder)||(NSUInteger)(reader.end-reader.p)<11){status=4;return;}uint8_t inputKind=u8(&reader);BOOL secure=u8(&reader),multiline=u8(&reader);int32_t maxLength=i32(&reader);NSData *style=nil;if(!GNReadBoundedData(&reader,&style)){status=4;return;}CGFloat width=isfinite(maxWidth)&&maxWidth>0?maxWidth:CGFLOAT_MAX,height=isfinite(maxHeight)&&maxHeight>0?maxHeight:CGFLOAT_MAX;CGSize size=GNMeasureControl(kind,text,image,style,CGSizeMake(width,height),wrap,overflow,maxLines,selectable,rich,placeholder,inputKind,secure,multiline,maxLength);size.width=MAX(minWidth,MIN(size.width,width));size.height=MAX(minHeight,MIN(size.height,height));GNAppendU64(output,requestID);GNAppendF32(output,(float)size.width);GNAppendF32(output,(float)size.height);GNAppendU32(output,0);}
        if(reader.p!=reader.end||output.length>16777216){status=5;output=nil;}
    };if(NSThread.isMainThread)measure();else dispatch_sync(dispatch_get_main_queue(),measure);if(status!=0||!output)return status?:6;void *buffer=malloc(output.length);if(!buffer)return 7;memcpy(buffer,output.bytes,output.length);*results=buffer;*resultLength=(int32_t)output.length;return 0;
}
void GNFreeNativeBuffer(void *buffer){free(buffer);}
