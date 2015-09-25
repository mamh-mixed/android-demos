//
//  VTGenQRViewController.m
//  CloudCashier
//
//  Created by 司瑞华 on 15/7/17.
//  Copyright (c) 2015年 __VTPayment__. All rights reserved.
//

#import "VTGenQRViewController.h"
#import "CloudCashierAPI.h"
#import "CustomAlertView.h"

#define SCreenWidth                                  self.view.frame.size.width
#define SCreenHeight                                 self.view.frame.size.height

@interface VTGenQRViewController ()<CustomAlertViewDelegate,BackInfoDelegate>
{
    id<BackInfoDelegate>            delegate;//6个接口返回数据需要用到的代理
}
@property(strong,nonatomic) UIImageView                 * qrcodeImgView;//生成二维码页面的View
@end

@implementation VTGenQRViewController

- (void)viewDidLoad {
    [super viewDidLoad];
    
    delegate = self;
    [CloudCashierAPI transmitDelegate:delegate];

    UIView * genQRView = [[UIView alloc]initWithFrame:self.view.frame];
    genQRView.tag = 100083;
    genQRView.backgroundColor =  [UIColor whiteColor];
    [self.view addSubview:genQRView];
    UIView * headView = [[UIView alloc]initWithFrame:CGRectMake(0, 0, SCreenWidth, 64)];
    headView.backgroundColor = [UIColor colorWithRed:0/255.0 green:187/255.0 blue:211/255.0 alpha:1];
    [genQRView addSubview:headView];
    //返回按钮
    UIButton * backBtn = [UIButton buttonWithType:UIButtonTypeRoundedRect];
    backBtn.frame = CGRectMake(20, (64-37*15/21)/2+5, 15, 37*15/21);
    [backBtn setBackgroundImage:[UIImage imageNamed:@"return.png"] forState:UIControlStateNormal];
    backBtn.tag = 8;
    [backBtn addTarget:self action:@selector(btnClick:) forControlEvents:UIControlEventTouchUpInside];
    [headView addSubview:backBtn];
    
    //titleLBL
    UILabel * titleLbl = [[UILabel alloc]initWithFrame:CGRectMake(SCreenWidth/2-60, 17, 120, 40)];
    titleLbl.text = @"扫码支付";
    titleLbl.font = [UIFont boldSystemFontOfSize:24];
    titleLbl.textAlignment = NSTextAlignmentCenter;
    titleLbl.textColor = [UIColor whiteColor];
    [headView addSubview:titleLbl];
    
    UILabel * contentOneLbl = [[UILabel alloc]initWithFrame:CGRectMake(SCreenWidth/2-SCreenWidth/3, 80, SCreenWidth/3*2, 30)];
    NSString * lblText = [self.chcd isEqualToString:@"alipay.png"] ? @"请打开支付宝钱包":@"请打开微信";
    contentOneLbl.text = lblText;
    contentOneLbl.textAlignment = NSTextAlignmentCenter;
    contentOneLbl.font = [UIFont systemFontOfSize:24];
    [genQRView addSubview:contentOneLbl];
    UILabel * contentTwoLbl = [[UILabel alloc]initWithFrame:CGRectMake(SCreenWidth/2-40, 108, 80, 25)];
    contentTwoLbl.textAlignment = NSTextAlignmentCenter;
    contentTwoLbl.text = @"扫一扫";
    contentTwoLbl.font = [UIFont systemFontOfSize:24];
    [genQRView addSubview:contentTwoLbl];
    
    //二维码图片
    UIImage *qrcode = [self createNonInterpolatedUIImageFormCIImage:[self createQRForString:self.qrInfo] withSize:250.0f];
    UIImage *customQrcode = [self imageBlackToTransparent:qrcode withRed:60.0f andGreen:74.0f andBlue:89.0f];
    self.qrcodeImgView = [[UIImageView alloc]initWithFrame:CGRectMake((SCreenWidth-220)/2, 137, 220, 220)];
    self.qrcodeImgView.image = customQrcode;
    // set shadow
    self.qrcodeImgView.layer.shadowOffset = CGSizeMake(0, 2);
    self.qrcodeImgView.layer.shadowRadius = 2;
    self.qrcodeImgView.layer.shadowColor = [UIColor blackColor].CGColor;
    self.qrcodeImgView.layer.shadowOpacity = 0.5;
    [genQRView addSubview:self.qrcodeImgView];
    //二维码中间的小图标
    UIImageView * middleImgView = [[UIImageView alloc]init];
    middleImgView.center = self.qrcodeImgView.center;
    middleImgView.backgroundColor = [UIColor whiteColor];
    middleImgView.bounds = CGRectMake(0, 0, 43, 43);
    middleImgView.image = [UIImage imageNamed:_chcd];
    [genQRView addSubview:middleImgView];
    
    //金额lbl
    UILabel * amountLbl = [[UILabel alloc]initWithFrame:CGRectMake((SCreenWidth-SCreenWidth/4*3)/2, 390, SCreenWidth/4*3, 30)];
    amountLbl.text = [NSString stringWithFormat:@"￥%@",_amount];
    amountLbl.textAlignment = NSTextAlignmentCenter;
    amountLbl.font = [UIFont systemFontOfSize:24];
    [genQRView addSubview:amountLbl];
    
    UIButton * queryBtn = [UIButton buttonWithType:UIButtonTypeRoundedRect];
    queryBtn.layer.cornerRadius =20;
    queryBtn.layer.masksToBounds = YES;
    [queryBtn setTitle:@"交易查询" forState:UIControlStateNormal];
    [queryBtn setTitleColor:[UIColor whiteColor] forState:UIControlStateNormal];
    queryBtn.titleLabel.font = [UIFont boldSystemFontOfSize:20];
    queryBtn.backgroundColor = [UIColor colorWithRed:0/255.0 green:187/255.0 blue:211/255.0 alpha:1];
    queryBtn.frame = CGRectMake((SCreenWidth-SCreenWidth/5*4)/2, 440, SCreenWidth/5*4, 40);
    queryBtn.tag = 10;
    [queryBtn addTarget:self action:@selector(btnClick:) forControlEvents:UIControlEventTouchUpInside];
    [genQRView addSubview:queryBtn];
}
-(void)btnClick:(UIButton *)btn
{
    if (btn.tag == 8)
    {
        [self dismissViewControllerAnimated:NO completion:nil];
    }else if(btn.tag == 10)
    {
        Parameter * para = [[Parameter alloc]init];
        para.origOrderNum = _qureyOrderNum;
        [CloudCashierAPI queryWithpara:para];
        
    }
}
-(void)getResultDataWithBackParameter:(BackParameter *)backData errorCode:(NSInteger)errorNum
{
    if (errorNum == 1)
    {
        CustomAlertView * alertView = [[CustomAlertView alloc]initWithTitle:nil icon:[UIImage imageNamed:@"wrong.png"] message:@"交易失败" subtitleMsg:@"请重新购买！" type:2 delegate:self buttonTitles:@"历史交易",@"返回", nil];
        
        alertView.tag = 10;
        [alertView show];
    }else if (backData.tag == 3)
    {
        NSLog(@"----输出支付状态----%@",backData.respcd);
        if ([backData.respcd isEqualToString:@"09"])
        {
            CustomAlertView * alertView = [[CustomAlertView alloc]initWithTitle:nil icon:[UIImage imageNamed:@"wrong.png"] message:@"交易未支付，请支付后查询" subtitleMsg:nil type:2 delegate:self buttonTitles:@"查询交易",@"关闭", nil];
            alertView.tag = 9;
            [alertView show];
        }else if ([backData.respcd isEqualToString:@"58"] || [backData.respcd isEqualToString:@"12"])
        {
            CustomAlertView * alertView = [[CustomAlertView alloc]initWithTitle:nil icon:[UIImage imageNamed:@"wrong.png"] message:@"交易失败" subtitleMsg:@"请重新购买！" type:2 delegate:self buttonTitles:@"历史交易",@"返回", nil];

            alertView.tag = 10;
            [alertView show];
        }else if ([backData.respcd isEqualToString:@"00"])
        {
            CustomAlertView * alertView = [[CustomAlertView alloc]initWithTitle:nil icon:[UIImage imageNamed:@"right.png"] message:@"交易成功确认" subtitleMsg:@"感谢购买！" type:2 delegate:self buttonTitles:@"历史交易",@"返回", nil];
            alertView.tag = 20;
            [alertView show];
        }
    }
}
- (void)alertView:(CustomAlertView *)alertView clickedButtonAtIndex:(NSInteger)buttonIndex
{
        NSLog(@"%ld", (long)buttonIndex);
        if (alertView.tag == 9 )
        {
            if (buttonIndex == 0)
            {
                Parameter * para = [[Parameter alloc]init];
                para.origOrderNum = _qureyOrderNum;
                [CloudCashierAPI queryWithpara:para];
            }else
            {
                [self dismissViewControllerAnimated:NO completion:nil];
            }
        }else if (alertView.tag == 20 || alertView.tag == 10 )//历史交易页面
        {
            if (buttonIndex == 0)
            {
                VTScannerViewController * vc = [[VTScannerViewController alloc]init];
                vc.whichPage = 2;
                [self presentViewController:vc animated:NO completion:nil];
            }else
            {
                [self dismissViewControllerAnimated:NO completion:nil];
            }           
        }
}
#pragma mark - InterpolatedUIImage
- (UIImage *)createNonInterpolatedUIImageFormCIImage:(CIImage *)image withSize:(CGFloat) size {
    CGRect extent = CGRectIntegral(image.extent);
    CGFloat scale = MIN(size/CGRectGetWidth(extent), size/CGRectGetHeight(extent));
    // create a bitmap image that we'll draw into a bitmap context at the desired size;
    size_t width = CGRectGetWidth(extent) * scale;
    size_t height = CGRectGetHeight(extent) * scale;
    CGColorSpaceRef cs = CGColorSpaceCreateDeviceGray();
    CGContextRef bitmapRef = CGBitmapContextCreate(nil, width, height, 8, 0, cs, (CGBitmapInfo)kCGImageAlphaNone);
    CIContext *context = [CIContext contextWithOptions:nil];
    CGImageRef bitmapImage = [context createCGImage:image fromRect:extent];
    CGContextSetInterpolationQuality(bitmapRef, kCGInterpolationNone);
    CGContextScaleCTM(bitmapRef, scale, scale);
    CGContextDrawImage(bitmapRef, extent, bitmapImage);
    // Create an image with the contents of our bitmap
    CGImageRef scaledImage = CGBitmapContextCreateImage(bitmapRef);
    // Cleanup
    CGContextRelease(bitmapRef);
    CGImageRelease(bitmapImage);
    return [UIImage imageWithCGImage:scaledImage];
}

#pragma mark - QRCodeGenerator根据信息生成二维码
- (CIImage *)createQRForString:(NSString *)qrString
{
    // Need to convert the string to a UTF-8 encoded NSData object
    NSData *stringData = [qrString dataUsingEncoding:NSUTF8StringEncoding];
    // Create the filter
    CIFilter *qrFilter = [CIFilter filterWithName:@"CIQRCodeGenerator"];
    // Set the message content and error-correction level
    [qrFilter setValue:stringData forKey:@"inputMessage"];
    [qrFilter setValue:@"M" forKey:@"inputCorrectionLevel"];
    // Send the image back
    return qrFilter.outputImage;
}
#pragma mark - imageToTransparent
void ProviderReleaseData (void *info, const void *data, size_t size){
    free((void*)data);
}
- (UIImage*)imageBlackToTransparent:(UIImage*)image withRed:(CGFloat)red andGreen:(CGFloat)green andBlue:(CGFloat)blue{
    const int imageWidth = image.size.width;
    const int imageHeight = image.size.height;
    size_t      bytesPerRow = imageWidth * 4;
    uint32_t* rgbImageBuf = (uint32_t*)malloc(bytesPerRow * imageHeight);
    // create context
    CGColorSpaceRef colorSpace = CGColorSpaceCreateDeviceRGB();
    CGContextRef context = CGBitmapContextCreate(rgbImageBuf, imageWidth, imageHeight, 8, bytesPerRow, colorSpace,
                                                 kCGBitmapByteOrder32Little | kCGImageAlphaNoneSkipLast);
    CGContextDrawImage(context, CGRectMake(0, 0, imageWidth, imageHeight), image.CGImage);
    // traverse pixe
    int pixelNum = imageWidth * imageHeight;
    uint32_t* pCurPtr = rgbImageBuf;
    for (int i = 0; i < pixelNum; i++, pCurPtr++){
        if ((*pCurPtr & 0xFFFFFF00) < 0x99999900){
            // change color
            uint8_t* ptr = (uint8_t*)pCurPtr;
            ptr[3] = red; //0~255
            ptr[2] = green;
            ptr[1] = blue;
        }else{
            uint8_t* ptr = (uint8_t*)pCurPtr;
            ptr[0] = 0;
        }
    }
    // context to image
    CGDataProviderRef dataProvider = CGDataProviderCreateWithData(NULL, rgbImageBuf, bytesPerRow * imageHeight, ProviderReleaseData);
    CGImageRef imageRef = CGImageCreate(imageWidth, imageHeight, 8, 32, bytesPerRow, colorSpace,
                                        kCGImageAlphaLast | kCGBitmapByteOrder32Little, dataProvider,
                                        NULL, true, kCGRenderingIntentDefault);
    CGDataProviderRelease(dataProvider);
    UIImage* resultUIImage = [UIImage imageWithCGImage:imageRef];
    // release
    CGImageRelease(imageRef);
    CGContextRelease(context);
    CGColorSpaceRelease(colorSpace);
    return resultUIImage;
}

#pragma mark - alertView
-(void)alertViewWithMessage:(NSString *)message
{
    UIAlertView * alertView = [[UIAlertView alloc]initWithTitle:@"温馨提示" message:message delegate:self cancelButtonTitle:@"知道了" otherButtonTitles:nil, nil];
    [alertView show];
}


@end
