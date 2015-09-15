//
//  VTStartScannerViewController.m
//  CloudCashier
//
//  Created by 司瑞华 on 15/7/17.
//  Copyright (c) 2015年 __VTPayment__. All rights reserved.
//

#import "VTStartScannerViewController.h"
#import <AVFoundation/AVFoundation.h>
#import "VTScannerView.h"
#import "CloudCashierAPI.h"
#import "CustomAlertView.h"

#define SCREENWIDTH                              self.frame.size.width
#define SCREENHEIGHT                             self.frame.size.height

@interface VTStartScannerViewController ()<AVCaptureMetadataOutputObjectsDelegate,BackInfoDelegate,CustomAlertViewDelegate>
{
    VTScannerView                            * _vtRectView;
    NSString                                 * _errorDetail;
    id<BackInfoDelegate>                     delegate;
    NSString                                 * _scannerCode;
    BOOL                                     isOff;
    NSString                                 * _querybill;
    CustomAlertView                          * _alertView;
}

@property ( strong , nonatomic ) AVCaptureDevice              * device;
@property ( strong , nonatomic ) AVCaptureDeviceInput         * input;
@property ( strong , nonatomic ) AVCaptureMetadataOutput      * output;
@property ( strong , nonatomic ) AVCaptureSession             * session;
@property ( strong , nonatomic ) AVCaptureVideoPreviewLayer   * preview;

@end

@implementation VTStartScannerViewController

- (void)viewDidLoad {
    [super viewDidLoad];
    
    _device = [AVCaptureDevice defaultDeviceWithMediaType:AVMediaTypeVideo];
    
    // Input
    _input = [AVCaptureDeviceInput deviceInputWithDevice:self.device error:nil];
    
    // Output
    _output = [[AVCaptureMetadataOutput alloc]init];
    [_output setMetadataObjectsDelegate:self queue:dispatch_get_main_queue()];
    
    // Session
    _session = [[AVCaptureSession alloc]init];
    [_session setSessionPreset:AVCaptureSessionPresetHigh];
    if([_session canAddInput:self.input])
    {
        [_session addInput:self.input];
    }
    if([_session canAddOutput:self.output])
    {
        [_session addOutput:self.output];
    }
    // 条码类型 AVMetadataObjectTypeQRCode
    _output.metadataObjectTypes = @[ AVMetadataObjectTypeQRCode,AVMetadataObjectTypeEAN13Code,AVMetadataObjectTypeEAN8Code,AVMetadataObjectTypeCode128Code] ;
    
    // Preview
    _preview = [AVCaptureVideoPreviewLayer layerWithSession:_session];
    //    _preview.videoGravity = AVLayerVideoGravityResizeAspectFill;
    _preview.videoGravity =AVLayerVideoGravityResize;
    _preview.frame = self.view.layer.bounds ;
    [ self.view.layer insertSublayer:_preview atIndex:0];
    // Start
    [ _session startRunning];
    
    CGRect screenRect = [UIScreen mainScreen].bounds;
    _vtRectView = [[VTScannerView alloc] initWithFrame:screenRect];
    _vtRectView.transparentArea = CGSizeMake(250, 250);
    _vtRectView.backgroundColor = [UIColor clearColor];
    _vtRectView.center = CGPointMake(self.view.frame.size.width / 2, self.view.frame.size.height / 2);
    [self.view addSubview:_vtRectView];
    
    
    CGFloat width = self.view.frame.size.width;
    UIButton * backBtn = [UIButton buttonWithType:UIButtonTypeRoundedRect];
    backBtn.frame = CGRectMake(width/10, 55, width/10, width/10);
    [backBtn setBackgroundImage:[UIImage imageNamed:@"1.png"] forState:UIControlStateNormal];
    backBtn.tag = 11;///扫码页面的返回按钮
    [backBtn addTarget:self action:@selector(btnClick:) forControlEvents:UIControlEventTouchUpInside];
    [self.view addSubview:backBtn];
    
    isOff = NO;
    UIButton * lightBtn = [UIButton buttonWithType:UIButtonTypeRoundedRect];
    [lightBtn setBackgroundImage:[UIImage imageNamed:@"2.png"] forState:UIControlStateNormal];
    lightBtn.frame = CGRectMake(250, 55, width/10, width/10);
    lightBtn.tag = 12;//扫码页面的闪光灯按钮
    [lightBtn addTarget:self action:@selector(btnClick:) forControlEvents:UIControlEventTouchUpInside];
    [self.view addSubview:lightBtn];

    NSString *mediaType = AVMediaTypeVideo;
    
    AVAuthorizationStatus authStatus = [AVCaptureDevice authorizationStatusForMediaType:mediaType];
    
    if(authStatus == AVAuthorizationStatusRestricted || authStatus == AVAuthorizationStatusDenied)
    {
        [self alertViewWithMessage:@"相机权限受限,请在“设置-隐私-相机”选项中,允许云收银访问您的相机"];
    }

    //修正扫描区域
    CGFloat screenHeight = self.view.frame.size.height;
    CGFloat screenWidth = self.view.frame.size.width;
    CGRect cropRect = CGRectMake((screenWidth - _vtRectView.transparentArea.width) / 2,
                                 (screenHeight - _vtRectView.transparentArea.height) / 2,
                                 _vtRectView.transparentArea.width,
                                 _vtRectView.transparentArea.height);
    
    [_output setRectOfInterest:CGRectMake(cropRect.origin.y / screenHeight,
                                          cropRect.origin.x / screenWidth,
                                          cropRect.size.height / screenHeight,
                                          cropRect.size.width / screenWidth)];

}
-(void)btnClick:(UIButton *)btn
{
    if (btn.tag == 11)///点击扫码页面btn 扫码页面移除
    {
        [self dismissViewControllerAnimated:NO completion:nil];
    }else if (btn.tag == 12)
    {
        AVCaptureDevice * device = [AVCaptureDevice defaultDeviceWithMediaType:AVMediaTypeVideo];
        if (![device hasTorch])
        {
            UIAlertView *alert = [[UIAlertView alloc] initWithTitle:@"闪光灯" message:@"抱歉，该设备没有闪光灯而无法使用手电筒功能！" delegate:nil
                                                  cancelButtonTitle:@"确定" otherButtonTitles:nil];
            
            [alert show];
        }else
        {
            if (isOff == NO)
            {
                [device lockForConfiguration:nil];
                [device setTorchMode: AVCaptureTorchModeOn];
                [device unlockForConfiguration];
                isOff = YES;
            }else if (isOff == YES)
            {
                [device lockForConfiguration:nil];
                [device setTorchMode: AVCaptureTorchModeOff];
                [device unlockForConfiguration];
                isOff = NO;
            }
        }
    }
}
#pragma mark AVCaptureMetadataOutputObjectsDelegate
- (void)captureOutput:(AVCaptureOutput *)captureOutput didOutputMetadataObjects:(NSArray *)metadataObjects fromConnection:(AVCaptureConnection *)connection
{
    NSString *stringValue;
    if ([metadataObjects count] >0)
    {
        //停止扫描
        [_session stopRunning];
        AVMetadataMachineReadableCodeObject * metadataObject = [metadataObjects objectAtIndex:0];
        stringValue = metadataObject.stringValue;
        [_vtRectView.timer invalidate];
        _vtRectView.timer = nil;
    }
    NSDate * datenow = [NSDate date];
    NSDateFormatter * formatter = [[NSDateFormatter alloc]init];
    [formatter setDateFormat:@"YYMMddHHmmss"];
    NSString * str = [formatter stringFromDate:datenow];
    NSString * orderNumStr = [NSString stringWithFormat:@"%@%d%d%d%d%d",str,arc4random()%10,arc4random()%10,arc4random()%10,arc4random()%10,arc4random()%10];
    
    delegate = self;
    [CloudCashierAPI transmitDelegate:delegate];
    Parameter * para = [[Parameter alloc] init];
    para.txamt = self.countMoney;
    //para.goodsInfo = @"花生,1,30";
    para.orderNum = orderNumStr;
    para.scanCodeId = stringValue;
    para.currency = @"CNY";
    [CloudCashierAPI scannerPayWithpara:para];
    
    _alertView = [[CustomAlertView alloc]initWithTitle:nil icon:[UIImage imageNamed:@"londing.png"] message:@"交易读取中，请稍后。。。" subtitleMsg:@"0s" type:1 delegate:self buttonTitles:@"查询结果",@"关闭", nil];
    _alertView.tag = 1005;
    [_alertView show];
    
    NSLog(@"--------- %@",stringValue);
}
-(void)getResultDataWithBackParameter:(BackParameter *)backData errorCode:(NSInteger)errorNum
{
    [_alertView hide];
    if (backData.tag == 1 || backData.tag == 3)
    {
        if ([backData.respcd isEqualToString:@"09"] )
        {
            _querybill = backData.orderNum;
            CustomAlertView * alertView = [[CustomAlertView alloc]initWithTitle:nil icon:[UIImage imageNamed:@"wrong.png"] message:@"交易未支付，请支付后查询" subtitleMsg:nil type:2 delegate:self buttonTitles:@"查询交易",@"关闭", nil];
            alertView.tag = 19;
            [alertView show];
        }else if ([backData.respcd isEqualToString:@"00"])
        {
            CustomAlertView * alertView = [[CustomAlertView alloc]initWithTitle:nil icon:[UIImage imageNamed:@"right.png"] message:@"交易成功确认" subtitleMsg:@"感谢购买！" type:2 delegate:self buttonTitles:@"历史交易",@"返回", nil];
            alertView.tag = 20;
            [alertView show];
        }else if ([backData.respcd isEqualToString:@"58"] || [backData.respcd isEqualToString:@"12"]|| [backData.respcd isEqualToString:@"96"])
        {
            CustomAlertView * alertView = [[CustomAlertView alloc]initWithTitle:nil icon:[UIImage imageNamed:@"wrong.png"] message:@"交易失败" subtitleMsg:@"请重新购买！" type:2 delegate:self buttonTitles:@"历史交易",@"返回", nil];
            alertView.tag = 21;
            [alertView show];
        }
    }
}

#pragma mark - 自定义alertView代理方法
- (void)alertView:(CustomAlertView *)alertView clickedButtonAtIndex:(NSInteger)buttonIndex {
    NSLog(@"%ld", (long)buttonIndex);
    
    if(alertView.tag == 19 || alertView.tag == 1005)
    {
        if (buttonIndex == 0)//查询交易
        {
            Parameter * para = [[Parameter alloc]init];
            para.origOrderNum = _querybill;
            [CloudCashierAPI queryWithpara:para];
        }else
        {
            [self dismissViewControllerAnimated:NO completion:nil];
        }
    }else if (alertView.tag == 20 || alertView.tag == 21)//历史交易页面
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

#pragma mark - alertView
-(void)alertViewWithMessage:(NSString *)message
{
    UIAlertView * alertView = [[UIAlertView alloc]initWithTitle:@"温馨提示" message:message delegate:self cancelButtonTitle:@"我知道了" otherButtonTitles:nil, nil];
    [alertView show];
}



@end
