//
//  ViewController.m
//  CloudCashier
//
//  Created by 司瑞华 on 15/7/7.
//  Copyright (c) 2015年 __VTPayment__. All rights reserved.
//

#import "VTViewController.h"
#import "PayEngine.h"
#import "Reachability.h"
#import "VTLRViewController.h"
#import "CloudCashierAPI.h"


//#import "VTRefdAlertView.h"


@interface VTViewController ()<UIWebViewDelegate,NSXMLParserDelegate,UIAlertViewDelegate>
{
    NSMutableArray                * parseArray;//xml解析出来的字符串
}

@end

@implementation VTViewController

-(void)viewDidLoad {
    [super viewDidLoad];
    self.view.backgroundColor = [UIColor colorWithRed:30/255.0 green:190/255.0  blue:214/255.0  alpha:1];
    UIWebView * webView = [[UIWebView alloc]initWithFrame:self.view.frame];
    NSString * mainBundlePath = [[NSBundle mainBundle] bundlePath];
    NSString * basePath = [NSString stringWithFormat:@"%@/cloudCashiercloud",mainBundlePath];
    NSURL * baseUrl = [NSURL fileURLWithPath:basePath isDirectory:YES];
    NSString * htmlPath = [NSString stringWithFormat:@"%@/home-page.html",basePath];
    NSString * htmlString = [NSString stringWithContentsOfFile:htmlPath encoding:NSUTF8StringEncoding error:nil];
    
    [webView loadHTMLString:htmlString baseURL:baseUrl];    
    webView.backgroundColor = [UIColor clearColor];
    webView.delegate = self;
    [self.view addSubview:webView];
}

-(BOOL)prefersStatusBarHidden//隐藏状态栏
{
    return YES;
}
-(void)webViewDidFinishLoad:(UIWebView *)webView
{
    
    [self performSelector:@selector(goToLoginView) withObject:self afterDelay:0.2f];//延长2秒调用去登录页面
}
-(void)goToLoginView
{
    NSURL * url = [NSURL URLWithString:@"http://qrcode.cardinfolink.net/app/version/ios_testversion.xml"];
    NSURLRequest * request = [NSURLRequest requestWithURL:url cachePolicy:NSURLRequestUseProtocolCachePolicy timeoutInterval:10.0f];
    self.currentText = [[NSMutableString alloc]init];
    parseArray = [[NSMutableArray alloc]initWithCapacity:0];
    NSURLResponse * response = nil;
    NSError * error = nil;
    NSData * data = [NSURLConnection sendSynchronousRequest:request returningResponse:&response error:&error];
    
    NSXMLParser * parser = [[NSXMLParser alloc]initWithData:data];
    [parser setDelegate:self];
    [parser parse];
    
    NSString *version = [[[NSBundle mainBundle] infoDictionary] objectForKey:(NSString *)kCFBundleVersionKey];
    NSLog(@"-------- [parseArray objectAtIndex:1] _%@----%@-----[parseArray objectAtIndex:3]_%@---",[parseArray objectAtIndex:0],version,[parseArray objectAtIndex:1]);
    if (![version isEqualToString:[parseArray objectAtIndex:1]])
    {
        NSLog(@"提醒更新");
        UIAlertView * alertView = [[UIAlertView alloc]initWithTitle:@"发现新的版本" message:[parseArray objectAtIndex:2] delegate:self cancelButtonTitle:@"取消" otherButtonTitles:@"去下载页", nil];
        alertView.delegate = self;
        [alertView show];
    }else
    {
        [self startGoToNextView];
    }
 }
-(BOOL)isConnection//判断是否有网络连接
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

#pragma mark - 调用alertView自带的方法
-(void)alertView:(UIAlertView *)alertView clickedButtonAtIndex:(NSInteger)buttonIndex
{
    if (buttonIndex == 1)
    {
        [[UIApplication sharedApplication] openURL:[NSURL URLWithString:[parseArray objectAtIndex:0]]];
    }else
    {
        [self startGoToNextView];
    }
}

-(void)startGoToNextView
{
    BOOL isExistenceNetwork = [self isConnection];
    VTLRViewController * loginVc = [[VTLRViewController alloc]init];
    if (isExistenceNetwork == YES)
    {
        NSUserDefaults * defaults = [NSUserDefaults standardUserDefaults];
        BOOL status = [defaults boolForKey:@"recordState"];
        if (status == YES)
        {
            //向服务器发送请求，登录
            [PayEngine logPayViewWithUserName:[defaults objectForKey:@"username"] password:[defaults objectForKey:@"password"] succeedBlock:^(NSDictionary *receiveDict) {
                NSString * stateStr = [receiveDict objectForKey:@"state"];
                NSString * error = [receiveDict objectForKey:@"error"];
                if ([stateStr isEqualToString:@"success"])//如果状态是成功，则跳转到扫码页面
                {
                    [CloudCashierAPI registerInscd:[[receiveDict objectForKey:@"user"] objectForKey:@"inscd"] mchntid:[[receiveDict objectForKey:@"user"] objectForKey:@"clientid"] signKey:[[receiveDict objectForKey:@"user"] objectForKey:@"signKey"] terminalid:@"dsfdsf" tradeFrom:@"app"];
                    VTScannerViewController * vc = [[VTScannerViewController alloc]init];
                    vc.userName = [defaults objectForKey:@"username"];
                    vc.password = [defaults objectForKey:@"password"];
                    [defaults setObject:[receiveDict objectForKey:@"user"] forKey:@"dictionary"];
                    [defaults setValue:[defaults objectForKey:@"password"] forKey:@"recordpw"];
                    [defaults synchronize];
                    [self presentViewController:vc animated:NO completion:nil];
                }else
                {
                    if ([error isEqualToString:@"username_password_error"])//用户名错误
                    {
                        [self alertViewWithMessage:@"用户名密码错误"];
                    }else if ([error isEqualToString:@"username_no_exist"])
                    {
                        [self alertViewWithMessage:@"用户不存在"];
                    }
                    [self presentViewController:loginVc animated:NO completion:nil];//不论是什么错误都跳到登录页面
                }
            }];
        }else
        {
            [self presentViewController:loginVc animated:NO completion:nil];
        }
    }else//没有网的情况下，也是跳到登录页面
    {
        [self alertViewWithMessage:@"请检查您的网络连接"];
        [self presentViewController:loginVc animated:NO completion:nil];
    }
}


#pragma mark - XML解析器代理方法
//2开始解析一个元素,新的节点开始了。
-(void)parser:(NSXMLParser *)parser didStartElement:(NSString *)elementName namespaceURI:(NSString *)namespaceURI qualifiedName:(NSString *)qName attributes:(NSDictionary *)attributeDict
{
    
}

//3接收元素的数据(发现字符，这个方法会因为元素内容过大，此方法会被重复调用，需要拼接数据)
- (void)parser:(NSXMLParser *)parser foundCharacters:(NSString *)string
{
    if (![string isEqualToString:@"\n"])
    {
        [parseArray addObject:string];
    }
}



#pragma mark - alertView
-(void)alertViewWithMessage:(NSString *)message
{
    UIAlertView * alertView = [[UIAlertView alloc]initWithTitle:@"温馨提示" message:message delegate:self cancelButtonTitle:@"知道了" otherButtonTitles:nil, nil];
    [alertView show];
}


@end
