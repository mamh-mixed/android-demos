//
//  VTScannerViewController.m
//  CloudCashier
//
//  Created by 司瑞华 on 15/7/7.
//  Copyright (c) 2015年 __VTPayment__. All rights reserved.
//

#import "VTScannerViewController.h"
#import "CustomAlertView.h"
#import "CloudCashierAPI.h"
#import "VTLRViewController.h"
#import "PayEngine.h"
#import "VTScannerView.h"
#import "VTStartScannerViewController.h"
#import "VTGenQRViewController.h"
#import <AVFoundation/AVFoundation.h>
#import "Reachability.h"



#define SCreenWidth                                  self.view.frame.size.width
#define SCreenHeight                                 self.view.frame.size.height


@interface VTScannerViewController ()<AVCaptureMetadataOutputObjectsDelegate,UIWebViewDelegate,BackInfoDelegate,CustomAlertViewDelegate>
{
    VTScannerView                   * _vtRectView;//扫码页面中嵌套的View
    UIWebView                       * _webView;
    NSString                        * _basePath;//加载本地H5的基本路径
    id<BackInfoDelegate>            delegate;//6个接口返回数据需要用到的代理
    double                          _refdAmount;//退款用的金额
    NSString                        * _origeOrderNum ;//用于退款的原订单号
    BOOL                            _isLimit;//是否受限额控制

}

@end

@implementation VTScannerViewController

- (void)viewDidLoad {
    [super viewDidLoad];
    self.view.backgroundColor = [UIColor whiteColor];
    
    UIView * bgView = [[UIView alloc]init];
    bgView.backgroundColor = [UIColor colorWithRed:30/255.0 green:190/255.0  blue:214/255.0  alpha:1];
    bgView.frame = CGRectMake(0, 0, SCreenWidth, 200);
    [self.view addSubview:bgView];
    
    _webView = [[UIWebView alloc]initWithFrame:CGRectMake(0, 20, SCreenWidth, SCreenHeight-20)];
    
    
    NSUserDefaults * defaults = [NSUserDefaults standardUserDefaults];
    self.dictionary = [NSDictionary dictionaryWithDictionary:[defaults objectForKey:@"dictionary"]];
    
    //加载本地登录页面的路径
    NSString * mainBundlePath = [[NSBundle mainBundle] bundlePath];
    _basePath = [NSString stringWithFormat:@"%@/cloudCashiercloud",mainBundlePath];
    NSURL * baseUrl = [NSURL fileURLWithPath:_basePath isDirectory:YES];
    NSString * htmlPath = [NSString stringWithFormat:@"%@/index.html",_basePath];
    NSString * htmlString = [NSString stringWithContentsOfFile:htmlPath encoding:NSUTF8StringEncoding error:nil];
    
    
    [_webView loadHTMLString:htmlString baseURL:baseUrl];
    _webView.scrollView.scrollEnabled = YES;
    if (self.whichPage == 2)
    {
        _webView.tag = 2;
    }else
    {
         _webView.tag = 10;//登录页面
    }
    _webView.backgroundColor = [UIColor clearColor];
    _webView.delegate = self;
    [self.view addSubview:_webView];
}
-(BOOL)isConnectNet
{
    BOOL isExistenceNetwork = '\1' ;
    Reachability *reach = [Reachability reachabilityWithHostName:@"www.baidu.com"];
    switch ([reach currentReachabilityStatus])
    {
        case NotReachable:
            isExistenceNetwork = NO;
            NSLog(@"notReachable");
            break;
        case ReachableViaWiFi:
            isExistenceNetwork = YES;
            NSLog(@"WIFI");
            break;
        case ReachableViaWWAN:
            isExistenceNetwork = YES;
            NSLog(@"3G");
            break;
    }
    return isExistenceNetwork;
}


#pragma mark - 调用js 方法
-(void)webViewDidFinishLoad:(UIWebView *)webView
{
    NSUserDefaults * defaults = [NSUserDefaults standardUserDefaults];
    NSMutableDictionary * dict = [[NSMutableDictionary alloc]initWithCapacity:0];
    [dict setObject:[self.dictionary objectForKey:@"username"] forKey:@"username"];
    [dict setObject:[defaults objectForKey:@"recordpw"] forKey:@"password"];
    [dict setObject:@"eu1dr0c8znpa43blzy1wirzmk8jqdaon" forKey:@"key"];
    [dict setObject:@"ios" forKey:@"device"];
    [dict setObject:[self.dictionary objectForKey:@"clientid"] forKey:@"clientid"];
    if (webView.tag == 2)
    {
        [dict setObject:@"transManage" forKey:@"target"];
    }if (webView.tag == 10)
    {
        [dict setObject:@"scanPage" forKey:@"target"];
    }
    NSData * jsonData = [NSJSONSerialization dataWithJSONObject:dict options:NSJSONWritingPrettyPrinted error:nil];
    NSString * jsonString ;
    if ([jsonData length] > 0)
    {
        jsonString = [[NSString alloc]initWithData:jsonData encoding:NSUTF8StringEncoding];
        NSString * js = [NSString stringWithFormat:@"CloudCashierBridge.saveUserData(%@)", jsonString];
        [webView stringByEvaluatingJavaScriptFromString:js];
    }
}
#pragma mark - 捕捉H5事件并做相应地处理
-(BOOL)webView:(UIWebView *)webView shouldStartLoadWithRequest:(NSURLRequest *)request navigationType:(UIWebViewNavigationType)navigationType
{
    NSString * urlString = [[request URL] absoluteString];
    NSString * decodeStr = [urlString stringByRemovingPercentEncoding];
    NSArray * components = [decodeStr componentsSeparatedByString:@"://"];
    if ([self isConnectNet] == YES)
    {
        if ([components count] && [[components objectAtIndex:0] isEqualToString:@"cloudcashier"])
        {
            NSArray * diffArray = [(NSString *)[components objectAtIndex:1] componentsSeparatedByString:@"/"];
            NSString * importantStr = [diffArray objectAtIndex:0];
            NSLog(@"------输出用来区别是什么事件的字符串---%@",importantStr);
            /*  扫码支付主页面包含
             *  下单支付
             *  预下单支付
             */
            if ([importantStr isEqualToString:@"scancode"])
            {
                NSData * jsonData = [[diffArray objectAtIndex:1] dataUsingEncoding:NSUTF8StringEncoding];
                NSError * error;
                NSDictionary * dic = [NSJSONSerialization JSONObjectWithData:jsonData options:NSJSONReadingMutableContainers error:&error];
                NSLog(@"-下单支付||预下单支付--输出dic 查看所需要的值---%@",dic);
                if ([[dic objectForKey:@"busicd"] isEqualToString:@"PURC"])
                {
                    /*
                     *  下单支付
                     *
                     */
                    [self getTotal];
                    if (_isLimit == NO)
                    {
                        if ([self validateCamera])
                        {
                            NSString *mediaType = AVMediaTypeVideo;
                            
                            AVAuthorizationStatus authStatus = [AVCaptureDevice authorizationStatusForMediaType:mediaType];
                            
                            if(authStatus == AVAuthorizationStatusRestricted || authStatus == AVAuthorizationStatusDenied)
                            {
                                [self alertViewWithMessage:@"相机权限受限,请在“设置-隐私-相机”选项中,允许云收银访问您的相机"];
                            }else
                            {
                                VTStartScannerViewController * vc = [[VTStartScannerViewController alloc]init];
                                vc.countMoney = [dic objectForKey:@"sum"];
                                [self presentViewController:vc animated:NO completion:nil];
                            }
                        }else
                        {
                            [self alertViewWithMessage:@"没有摄像头或摄像头不可用"];
                        }
                        
                    }else
                    {
                        [self alertViewWithMessage:@"您今天的交易已经达到限额，若想继续请提升限额"];
                    }
                }else if ([[dic objectForKey:@"busicd"] isEqualToString:@"PAUT"])
                {
                    /*
                     *
                     *  预下单支付
                     */
                    [self getTotal];
                    if (_isLimit == NO)
                    {
                        NSDate * datenow = [NSDate date];
                        NSDateFormatter * formatter = [[NSDateFormatter alloc]init];
                        [formatter setDateFormat:@"YYMMddHHmmss"];
                        NSString * timeStr = [formatter stringFromDate:datenow];
                        NSString * orderNumStr = [NSString stringWithFormat:@"%@%u%u%u%u%u%u%u%u%u%u",timeStr,arc4random()%10,arc4random()%10,arc4random()%10,arc4random()%10,arc4random()%10,arc4random()%10,arc4random()%10,arc4random()%10,arc4random()%10,arc4random()%10];
                        delegate = self;
                        [CloudCashierAPI transmitDelegate:delegate];
                        Parameter * para = [[Parameter alloc]init];
                        para.orderNum = orderNumStr;
                        para.txamt = [dic objectForKey:@"sum"];
                        para.chcd = [dic objectForKey:@"chcd"];
                        para.currency = @"CNY";
                        para.goodsInfo = @"讯联云收银在线支付";
                        [CloudCashierAPI preOrderPayWithpara:para];
                        
                    }else
                    {
                        [self alertViewWithMessage:@"您今天的交易已经达到限额，若想继续请提升限额"];
                    }
                }
            }else if ([importantStr isEqualToString:@"openwapbill"])
            {
                /*
                 *
                 *  打开网页版的账单
                 */
                NSLog(@"打开网页版的账单");
                NSString * str = [NSString stringWithFormat:NETLIST,[self.dictionary objectForKey:@"objectId"]];
                [[UIApplication sharedApplication] openURL:[NSURL URLWithString:str]];
            }else if ([importantStr isEqualToString:@"safeexit"])
            {
                /*
                 *
                 *  安全退出
                 */
                [self deletePasswordInDefaults];
                VTLRViewController * vc = [[VTLRViewController alloc]init];
                [self presentViewController:vc animated:NO completion:nil];
            }else if ([importantStr isEqualToString:@"limitincrease"])//限额提升页面
            {
                /*
                 *
                 *  限额提升页面
                 */
                NSLog(@"-------输出从其他页面传送过来的dictionary---%@",self.dictionary);
                if ([[self.dictionary objectForKey:@"limit"] isEqualToString:@"true"])
                {
                    NSData * jsonData = [[diffArray objectAtIndex:1] dataUsingEncoding:NSUTF8StringEncoding];
                    NSError * error;
                    NSDictionary * dic = [NSJSONSerialization JSONObjectWithData:jsonData options:NSJSONReadingMutableContainers error:&error];
                    NSLog(@"---输出dic 查看所需要的值---%@",dic);
                    [PayEngine limitinCreaseWithUserName:self.userName password:self.password email:[dic objectForKey:@"email"] payee:[dic objectForKey:@"payee"] phoneNum:[dic objectForKey:@"phone_num"] succeedBlock:^(NSDictionary *receiveDict) {
                        
                        if ([[receiveDict objectForKey:@"state"] isEqualToString:@"success"])
                        {
                            [self alertViewWithMessage:@"申请成功"];
                        }else
                        {
                            [self alertViewWithMessage:@"申请失败"];
                        }
                    }];
                }else
                {
                    [self alertViewWithMessage:@"您的限额已经提升!"];
                }
            }else if ([importantStr isEqualToString:@"updateaccount"])
            {
                /*
                 *
                 *  修改账户页面
                 */
                NSData * jsonData = [[diffArray objectAtIndex:1] dataUsingEncoding:NSUTF8StringEncoding];
                NSError * error;
                NSDictionary * dic = [NSJSONSerialization JSONObjectWithData:jsonData options:NSJSONReadingMutableContainers error:&error];
                NSLog(@"---输出dic 查看所需要的值---%@",dic);
                NSUserDefaults * userDefaults = [NSUserDefaults standardUserDefaults];
                [PayEngine updateAccountWithUserName:self.userName password:[userDefaults objectForKey:@"recordpw"] bankOpen:[dic objectForKey:@"bank_open"] payee:[dic objectForKey:@"payee"] payeeCard:[dic objectForKey:@"payee_card"] phoneNum:[dic objectForKey:@"phone_num"] succeedBlock:^(NSDictionary *receiveDict) {
                    if ([[receiveDict objectForKey:@"state"] isEqualToString:@"success"])
                    {
                        [self alertViewWithMessage:@"修改成功"];
                    }else
                    {
                        [self alertViewWithMessage:@"修改失败"];
                    }
                }];
            }else if ([importantStr isEqualToString:@"updatepassword"])
            {
                /*
                 *
                 *  修改密码页面
                 */
                NSData * jsonData = [[diffArray objectAtIndex:1] dataUsingEncoding:NSUTF8StringEncoding];
                NSError * error;
                NSDictionary * dic = [NSJSONSerialization JSONObjectWithData:jsonData options:NSJSONReadingMutableContainers error:&error];
                NSLog(@"---输出dic 判断调什么界面---%@",dic);
                [PayEngine updatePasswordWithUserName:self.userName oldPassword:[dic objectForKey:@"oldpwd"] newPassword:[dic objectForKey:@"newpwd"] succeedBlock:^(NSDictionary *receiveDict) {
                    NSLog(@"---------修改密码返回的数据----%@",receiveDict);
                    if ([[receiveDict objectForKey:@"state"] isEqualToString:@"success"])
                    {
                        NSUserDefaults * userDefaults = [NSUserDefaults standardUserDefaults];
                        [userDefaults setObject:[dic objectForKey:@"newpwd"] forKey:@"recordpw"];
                        if (![[userDefaults objectForKey:@"password"] isEqualToString:@""] && [userDefaults objectForKey:@"password"]  != nil)
                        {
                            [userDefaults setObject:[dic objectForKey:@"newpwd"] forKey:@"password"];
                        }
                        [userDefaults synchronize];
                        [self alertViewWithMessage:@"修改成功"];
                        
                        NSString * jsonString = [dic objectForKey:@"newpwd"];
                        if ([jsonString length] > 0)
                        {
                            NSString * js = [NSString stringWithFormat:@"Util.setPassword(%@)", jsonString];
                            [webView stringByEvaluatingJavaScriptFromString:js];
                        }
                    }else if ([[receiveDict objectForKey:@"error"] isEqualToString:@"username_password_error"])
                    {
                        [self alertViewWithMessage:@"原密码错误"];
                    }
                }];
            }else if ([importantStr isEqualToString:@"refd"])
            {
                /*
                 *
                 *  //退款按钮
                 */
                NSData * jsonData = [[diffArray objectAtIndex:1] dataUsingEncoding:NSUTF8StringEncoding];
                NSError * error;
                NSDictionary * dic = [NSJSONSerialization JSONObjectWithData:jsonData options:NSJSONReadingMutableContainers error:&error];
                NSLog(@"--退款界面-输出dic 查看所需要的值---%@",dic);
                _origeOrderNum = [dic objectForKey:@"orderNum"];
                NSString * amountTotalStr = [dic objectForKey:@"total"];
                [PayEngine checkBalanceWithUserName:[self.dictionary objectForKey:@"username"] password:[self.dictionary objectForKey:@"password"] clientId:[self.dictionary objectForKey:@"clientid"] orderNum:[dic objectForKey:@"orderNum"] succeedBlock:^(NSDictionary *receiveDict) {
                    NSLog(@"---------余额数据----%@",receiveDict);
                    
                    NSString * hasRefdStr = [receiveDict objectForKey:@"refdtotal"];
                    double amounTotal = [amountTotalStr doubleValue];
                    double hasRefd = [hasRefdStr doubleValue];
                    _refdAmount = amounTotal - hasRefd;
                    
                    if (_refdAmount > 0)
                    {
                        NSString * str = [NSString stringWithFormat:@"本次可退款额度:￥%.2f",_refdAmount];
                        CustomAlertView * alertView = [[CustomAlertView alloc]initWithTitle:str message:@"请输入退款金额:" subtitleMsg:@"请输入登录密码:" delegate:self buttonTitles:@"取消",@"确认", nil];
                        alertView.tag = 100;
                        [alertView show];
                    }else
                    {
                        [self alertViewWithMessage:@"交易已经被退款，谢谢!"];
                    }
                }];
            }
        }
    }else
    {
        [self alertViewWithMessage:@"请检查您的网络连接"];
    }
    return YES;
}

#pragma mark - 删除本地存储的数据
-(void)deletePasswordInDefaults
{
    NSUserDefaults * defaults = [NSUserDefaults standardUserDefaults];
    [defaults removeObjectForKey:@"password"];
    [defaults setBool:NO forKey:@"recordState"];
    [defaults synchronize];
}
#pragma mark - 获取一天内交易的总金额
-(void)getTotal
{
    NSUserDefaults * defaults = [NSUserDefaults standardUserDefaults];
    self.dictionary = [NSDictionary dictionaryWithDictionary:[defaults objectForKey:@"dictionary"]];
    if ([[self.dictionary objectForKey:@"limit"] isEqualToString:@"true"])
    {
        NSUserDefaults * defaults = [NSUserDefaults standardUserDefaults];
        [PayEngine getTotalWithUserName:[defaults objectForKey:@"username"] password:[defaults objectForKey:@"recordpw"] clientId:[self.dictionary objectForKey:@"clientid"] succeedBlock:^(NSDictionary *receiveDict) {
            //            NSLog(@"---一天内的交易总额返回数据---receiveDict---%@--",receiveDict);
            NSString * totalStr = [receiveDict objectForKey:@"total"];
            double total = [totalStr doubleValue];
            if (total > 500)
            {
                _isLimit = YES;
            }else
            {
                _isLimit = NO;
            }
        }];
    }else
    {
        _isLimit = NO;
    }
}

-(void)alertView:(CustomAlertView *)alertView clickedButtonAtIndex:(NSInteger)buttonIndex
{
    if (alertView.tag == 100)
    {
        if (buttonIndex == 0)
        {
            NSLog(@"0");
        }else
        {
            if (_refdAmount > 0)
            {
                NSUserDefaults * userDefaults = [NSUserDefaults standardUserDefaults];
                if ([[userDefaults objectForKey:@"recordpw"] isEqualToString:alertView.passwTF.text])
                {
                    NSDate * datenow = [NSDate date];
                    NSDateFormatter * formatter = [[NSDateFormatter alloc]init];
                    [formatter setDateFormat:@"YYMMddHHmmss"];
                    NSString * timeStr = [formatter stringFromDate:datenow];
                    NSString * orderNumStr = [NSString stringWithFormat:@"%@%u%u%u%u%u%u%u%u%u%u",timeStr,arc4random()%10,arc4random()%10,arc4random()%10,arc4random()%10,arc4random()%10,arc4random()%10,arc4random()%10,arc4random()%10,arc4random()%10,arc4random()%10];
                    
                    delegate = self;
                    [CloudCashierAPI transmitDelegate:delegate];

                    NSLog(@"1-------%@",alertView.amountTF.text);
                    Parameter * para = [[Parameter alloc]init];
                    para.orderNum = orderNumStr;
                    para.txamt = alertView.amountTF.text;
                    para.origOrderNum = _origeOrderNum;
                    para.currency = @"CNY";
                    para.goodsInfo = @"讯联云收银在线支付";
                    [CloudCashierAPI refundWithpara:para];
                }else
                {
                    [self alertViewWithMessage:@"密码错误"];
                }

            }else
            {
                [self alertViewWithMessage:@"您已超过最大退款额度"];
            }
        }
    }
}

//是否有可利用的摄像头
- (BOOL)validateCamera {
    
    return [UIImagePickerController isSourceTypeAvailable:UIImagePickerControllerSourceTypeCamera] &&
    [UIImagePickerController isCameraDeviceAvailable:UIImagePickerControllerCameraDeviceRear];
}


#pragma mark - 支付、查询的代理方法
-(void)getResultDataWithBackParameter:(BackParameter *)backData errorCode:(NSInteger)errorNum
{
    if (backData.tag == 2)
    {
        NSString * chcd = [backData.chcd isEqualToString:@"ALP"]?@"alipay.png":@"wechat.png";
        VTGenQRViewController * vc = [[VTGenQRViewController alloc] init];
        vc.chcd = chcd;
        vc.amount = backData.txamt;
        vc.qrInfo = backData.qrcode;
        vc.qureyOrderNum = backData.orderNum;
        [self presentViewController:vc animated:NO completion:nil];
    }else if (backData.tag == 5)
    {
        NSLog(@"----输出支付状态----%@---%@",backData.respcd,backData.errorDetail);
        if ([backData.respcd isEqualToString:@"00"])
        {
            CustomAlertView * alertView = [[CustomAlertView alloc]initWithTitle:nil icon:[UIImage imageNamed:@"right.png"] message:@"退款成功!" subtitleMsg:nil  type:2 delegate:self buttonTitles:@"确定", nil];
            alertView.tag = 20;
            [alertView show];
        }else if ([backData.respcd isEqualToString:@"58"] || [backData.respcd isEqualToString:@"12"]|| [backData.respcd isEqualToString:@"96"])
        {
            CustomAlertView * alertView = [[CustomAlertView alloc]initWithTitle:nil icon:[UIImage imageNamed:@"wrong.png"] message:@"退款失败!" subtitleMsg:nil type:2 delegate:self buttonTitles:@"确定", nil];
            alertView.tag = 20;
            [alertView show];
        }
    }
}
#pragma mark - alertView
-(void)alertViewWithMessage:(NSString *)message
{
    UIAlertView * alertView = [[UIAlertView alloc]initWithTitle:@"温馨提示" message:message delegate:self cancelButtonTitle:@"我知道了" otherButtonTitles:nil, nil];
    [alertView show];
}


-(UIStatusBarStyle)preferredStatusBarStyle
{
    return UIStatusBarStyleLightContent;
}

@end

#pragma mark - 获取uiwebview中的alertView
@interface UIWebView (JavaScriptAlert)

- (void)webView:(UIWebView *)sender runJavaScriptAlertPanelWithMessage:(NSString *)message initiatedByFrame:(CGRect *)frame;

@end

@implementation UIWebView (JavaScriptAlert)

- (void)webView:(UIWebView *)sender runJavaScriptAlertPanelWithMessage:(NSString *)message initiatedByFrame:(CGRect *)frame {
    
    
    UIAlertView* customAlert = [[UIAlertView alloc] initWithTitle:@"温馨提示"
                                                          message:message
                                                         delegate:nil
                                                cancelButtonTitle:@"我知道了"
                                                otherButtonTitles:nil];
    [customAlert show];
}
@end

