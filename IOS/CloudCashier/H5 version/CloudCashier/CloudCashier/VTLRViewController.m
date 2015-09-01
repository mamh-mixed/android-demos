//
//  VTLRViewController.m
//  CloudCashier
//
//  Created by 司瑞华 on 15/7/7.
//  Copyright (c) 2015年 __VTPayment__. All rights reserved.
//

#import "VTLRViewController.h"
#import "CustomAlertView.h"
#import "PayEngine.h"
#import "CloudCashierAPI.h"
#import "Reachability.h"

#define SCreenWidth                                  self.view.frame.size.width
#define SCreenHeight                                 self.view.frame.size.height


@interface VTLRViewController ()<UIWebViewDelegate,CustomAlertViewDelegate>
{
    BOOL              isRegister;
    int               _type;//登录跳注册一：1  登录跳注册二：2  注册一跳注册二：3
    UIWebView         * _webView;
    NSString          * _userName;
    NSString          * _password ;
    NSString          * _basePath;
}

@end

@implementation VTLRViewController

- (void)viewDidLoad {
    [super viewDidLoad];
    self.view.backgroundColor = [UIColor colorWithRed:30/255.0 green:190/255.0  blue:214/255.0  alpha:1];
    isRegister = NO;
    
    _webView = [[UIWebView alloc]initWithFrame:CGRectMake(0, 20, SCreenWidth, SCreenHeight-20)];
    //加载本地登录页面的路径
    NSString * mainBundlePath = [[NSBundle mainBundle] bundlePath];
    _basePath = [NSString stringWithFormat:@"%@/cloudCashiercloud",mainBundlePath];
    NSURL * baseUrl = [NSURL fileURLWithPath:_basePath isDirectory:YES];
    NSString * htmlPath = [NSString stringWithFormat:@"%@/login.html",_basePath];
    NSString * htmlString = [NSString stringWithContentsOfFile:htmlPath encoding:NSUTF8StringEncoding error:nil];
    
    [_webView loadHTMLString:htmlString baseURL:baseUrl];
    _webView.scrollView.scrollEnabled = YES;
    _webView.tag = 10;//登录页面
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
#pragma mark - 状态栏字体变成白色
-(UIStatusBarStyle)preferredStatusBarStyle
{
    return UIStatusBarStyleLightContent;
}
#pragma mark - webView的协议方法

-(void)webViewDidFinishLoad:(UIWebView *)webView
{
    if (webView.tag == 10)
    {
        NSMutableDictionary * dict = [[NSMutableDictionary alloc]initWithCapacity:0];
        [dict setValue:@"ios" forKey:@"device"];
        NSUserDefaults * defaults = [NSUserDefaults standardUserDefaults];
        BOOL autoLoginStr = [defaults boolForKey:@"recordState"];
        NSString * autologin = autoLoginStr == 1?@"true":@"false";
        [dict setValue:autologin forKey:@"autologin"];
        _userName = [defaults valueForKey:@"username"];
        _password = [defaults valueForKey:@"password"];
        NSString * userna = !_userName ? @"":_userName;
        NSString * pasw = !_password ? @"":_password;
        [dict setValue:userna forKey:@"username"];
        [dict setValue:pasw forKey:@"password"];
        
        NSData * jsonData = [NSJSONSerialization dataWithJSONObject:dict options:NSJSONWritingPrettyPrinted error:nil];
        NSString * jsonString;
        if ([jsonData length] > 0 )
        {
            jsonString = [[NSString alloc] initWithData:jsonData encoding:NSUTF8StringEncoding];
            NSString * js = [NSString stringWithFormat:@"Login.setParameter(%@)",jsonString];
            [webView stringByEvaluatingJavaScriptFromString:js];
        }
    }
}
//这个方法是网页中的每一个请求都会被触发的
-(BOOL)webView:(UIWebView *)webView shouldStartLoadWithRequest:(NSURLRequest *)request navigationType:(UIWebViewNavigationType)navigationType
{
    NSString * urlString = [[request URL] absoluteString];
    NSString * decodeStr = [urlString stringByRemovingPercentEncoding];
    NSArray * components = [decodeStr componentsSeparatedByString:@"://"];
    if ([components count] && [[components objectAtIndex:0] isEqualToString:@"cloudcashier"])
    {
        NSArray * diffArray = [(NSString *)[components objectAtIndex:1] componentsSeparatedByString:@"/"];
        NSString * importantStr = [diffArray objectAtIndex:0];
        NSLog(@"------输出用来区别是什么事件的字符串---%@",importantStr);
        if ([importantStr isEqualToString:@"jump_register"])//跳到注册页面 jump_register
        {
            isRegister = YES;
            _type = 1;
            [self jumpToOneRegisterViewController];
        }else if ([importantStr isEqualToString:@"login"])//点击登录按钮
        {
            NSData * jsonData = [[diffArray objectAtIndex:1] dataUsingEncoding:NSUTF8StringEncoding];
            NSError * error;
            NSDictionary * dic = [NSJSONSerialization JSONObjectWithData:jsonData options:NSJSONReadingMutableContainers error:&error];
            NSLog(@"输出点击登录按钮时获取的数据---dic ---%@",dic);
            _userName = [dic objectForKey:@"username"];
            _password = [dic objectForKey:@"password"];
            if ([self isConnectNet] == YES)
            {
                [self aboutLoginWithUsername:[dic objectForKey:@"username"] password:[dic objectForKey:@"password"] autoLogin:[dic objectForKey:@"autologin"]];
            }else
            {
                [self alertViewWithMessage:@"请检查您的网络连接"];
            }
            
        }else if ([importantStr isEqualToString:@"back"])
        {
            if (_type == 1)
            {
                //加载本地登录页面的路径
                NSURL * baseUrl = [NSURL fileURLWithPath:_basePath isDirectory:YES];
                NSString * htmlPath = [NSString stringWithFormat:@"%@/login.html",_basePath];
                NSString * htmlString = [NSString stringWithContentsOfFile:htmlPath encoding:NSUTF8StringEncoding error:nil];
                
                [_webView loadHTMLString:htmlString baseURL:baseUrl];
                _webView.tag = 10;//登录页面
            }else if (_type == 2)
            {
                //加载本地登录页面的路径
                NSURL * baseUrl = [NSURL fileURLWithPath:_basePath isDirectory:YES];
                NSString * htmlPath = [NSString stringWithFormat:@"%@/login.html",_basePath];
                NSString * htmlString = [NSString stringWithContentsOfFile:htmlPath encoding:NSUTF8StringEncoding error:nil];
                [_webView loadHTMLString:htmlString baseURL:baseUrl];
                _webView.tag = 10;//登录页面
            }else if (_type == 3)
            {
                [self jumpToOneRegisterViewController];
            }
        }else if ([importantStr isEqualToString:@"register"])
        {
            NSData * jsonData = [[diffArray objectAtIndex:1] dataUsingEncoding:NSUTF8StringEncoding];
            NSError * error;
            NSDictionary * dic = [NSJSONSerialization JSONObjectWithData:jsonData options:NSJSONReadingMutableContainers error:&error];
            NSLog(@"输出点击注册页面按钮时获取的数据---dic ---%@",dic);
            _userName = [dic objectForKey:@"username"];
            _password = [dic objectForKey:@"password"];
            if (isRegister == YES)
            {
                if ([self isConnectNet] == YES)
                {
                    [PayEngine registerAccountWithUserName:_userName password:_password succeedBlock:^(NSDictionary *receiveDict) {
                        if ([[receiveDict objectForKey:@"error"] isEqualToString:@"username_format_error"])
                        {
                            [self alertViewWithMessage:@"邮箱格式错误"];
                        }
                        if ([[receiveDict objectForKey:@"state"] isEqualToString:@"success"])
                        {
                            isRegister = NO;
                            CustomAlertView * alertView = [[CustomAlertView alloc]initWithTitle:@"请激活账号" icon:nil message:@"激活链接将发送到该邮箱:" subtitleMsg:[NSString stringWithFormat:@"%@",[dic objectForKey:@"username"]] type:2 delegate:self buttonTitles:@"取消",@"确定", nil];
                            [alertView show];
                        }else if ([[receiveDict objectForKey:@"error"] isEqualToString:@"username_exist"] )
                        {
                            [self alertViewWithMessage:@"用户名已存在"];
                        }
                    }];

                }else
                {
                    [self alertViewWithMessage:@"请检查您的网络连接"];
                }
            }else if(isRegister == NO)
            {
                if ([self isConnectNet] == YES)
                {
                   [self aboutLoginWithUsername:[dic objectForKey:@"username"] password:[dic objectForKey:@"password"] autoLogin:nil];
                }else
                {
                    [self alertViewWithMessage:@"请检查您的网络连接"];
                }
            }
        }else if ([importantStr isEqualToString:@"improveinfo"])//improveinfo
        {
            if ([self isConnectNet] == YES)
            {
                NSData * jsonData = [[diffArray objectAtIndex:1] dataUsingEncoding:NSUTF8StringEncoding];
                NSError * error;
                NSDictionary * dic = [NSJSONSerialization JSONObjectWithData:jsonData options:NSJSONReadingMutableContainers error:&error];
                NSLog(@"输出点击完善信息页面按钮时获取的数据---dic ---%@",dic);
                [PayEngine firstLogPayViewWithUserName:_userName password:_password bankName:[dic objectForKey:@"bank_open"] accountName:[dic objectForKey:@"payee"] bankNum:[dic objectForKey:@"payee_card"] phoneNum:[dic objectForKey:@"phone_num"] succeedBlock:^(NSDictionary *receiveDict) {
                    NSLog(@"-------receiveDict----%@------",receiveDict);
                    if ([[receiveDict objectForKey:@"state"] isEqualToString:@"success"])
                    {
                        [CloudCashierAPI registerInscd:[[receiveDict objectForKey:@"user"] objectForKey:@"inscd"] mchntid:[[receiveDict objectForKey:@"user"] objectForKey:@"clientid"] signKey:[[receiveDict objectForKey:@"user"] objectForKey:@"signKey"] terminalid:@"dsfdsf" tradeFrom:@"app"];
                        VTScannerViewController * vc = [[VTScannerViewController alloc]init];
                        NSUserDefaults * defaults = [NSUserDefaults standardUserDefaults];
                        [defaults setObject:[receiveDict objectForKey:@"user"] forKey:@"dictionary"];
                        [defaults setValue:_password forKey:@"recordpw"];
                        [defaults synchronize];
                        vc.userName = _userName;
                        vc.password = _password;
                        [self presentViewController:vc animated:NO completion:nil];
                    }else if ([[receiveDict objectForKey:@"error"] isEqualToString:@"info_exist"])
                    {
                        [self alertViewWithMessage:@"用户信息已存在"];
                    }
                }];
            }else
            {
                [self alertViewWithMessage:@"请检查您的网络连接"];
            }
        }
    }
    return YES;
}
#pragma mark - 注册界面一
-(void)jumpToOneRegisterViewController
{
    NSURL * baseUrl = [NSURL fileURLWithPath:_basePath isDirectory:YES];
    NSString * htmlPath = [NSString stringWithFormat:@"%@/register.html",_basePath];
    NSString * htmlString = [NSString stringWithContentsOfFile:htmlPath encoding:NSUTF8StringEncoding error:nil];
    _webView.tag = 11;
    [_webView loadHTMLString:htmlString baseURL:baseUrl];
}
#pragma mark - 注册界面二
-(void)jumpToTwoRegisterViewController
{
    NSURL * baseUrl = [NSURL fileURLWithPath:_basePath isDirectory:YES];
    NSString * htmlPath = [NSString stringWithFormat:@"%@/register-next.html",_basePath];
    NSString * htmlString = [NSString stringWithContentsOfFile:htmlPath encoding:NSUTF8StringEncoding error:nil];
    [_webView loadHTMLString:htmlString baseURL:baseUrl];
}
#pragma mark - 登录判断
-(void)aboutLoginWithUsername:(NSString *)userName password:(NSString *)pword autoLogin:(NSString *)autologin
{
    _userName = userName;
    _password = pword;
    [PayEngine logPayViewWithUserName:userName password:pword succeedBlock:^(NSDictionary *receiveDict) {
        NSLog(@"输出点击完善信息页面按钮时获取的数据---dic ---%@",receiveDict);
        NSString * stateStr = [receiveDict objectForKey:@"state"];
        NSString * activate = [[receiveDict objectForKey:@"user"] objectForKey:@"activate"];
        NSString * clientid = [[receiveDict objectForKey:@"user"] objectForKey:@"clientid"];
        NSString * error = [receiveDict objectForKey:@"error"];
                
        if ([error isEqualToString:@"user_no_activate"])
        {
            CustomAlertView * alertView = [[CustomAlertView alloc]initWithTitle:@"请激活账号" icon:nil message:@"激活链接将发送到该邮箱:" subtitleMsg:userName type:2 delegate:self buttonTitles:@"取消",@"确定", nil];
            [alertView show];
        }else if ([error isEqualToString:@"username_password_error"])
        {
            [self alertViewWithMessage:@"用户名密码错误"];
        }else if ([error isEqualToString:@"username_no_exist"])
        {
            [self alertViewWithMessage:@"用户不存在"];
        }else if ([stateStr isEqualToString:@"success"] && [activate isEqualToString:@"true"])
        {
            if (!clientid)
            {
                _type = 2;
                NSUserDefaults * defaults = [NSUserDefaults standardUserDefaults];
                [defaults setObject:userName forKey:@"username"];
                [defaults setValue:pword forKey:@"password"];
                [defaults synchronize];
                [self jumpToTwoRegisterViewController];
            }else
            {
                [CloudCashierAPI registerInscd:[[receiveDict objectForKey:@"user"] objectForKey:@"inscd"] mchntid:[[receiveDict objectForKey:@"user"] objectForKey:@"clientid"] signKey:[[receiveDict objectForKey:@"user"] objectForKey:@"signKey"] terminalid:@"dsfdsf" tradeFrom:@"app"];
                if ([autologin isEqualToString:@"true"])//自动登录
                {
                    [self writePasswordToDefaultsWithUserName:userName password:pword];
                }else if ([autologin isEqualToString:@"false"])//取消自动登录
                {
                    [self deletePasswordInDefaultsWithUserName:userName];
                }else if (autologin == nil)
                {
                    
                }
                VTScannerViewController * vc = [[VTScannerViewController alloc]init];
                vc.userName = userName;
                vc.password = pword;
                NSUserDefaults * defaults = [NSUserDefaults standardUserDefaults];
                [defaults setObject:[receiveDict objectForKey:@"user"] forKey:@"dictionary"];
                [defaults setValue:pword forKey:@"recordpw"];
                [defaults synchronize];
                [self presentViewController:vc animated:NO completion:nil];
            }
        }
    }];
}
-(void)sendLinkToMail//向邮箱中发送激活邮件
{
    [PayEngine activateAccountWithUserName:_userName password:_password succeedBlock:^(NSDictionary *receiveDict) {
    }];
}

#pragma mark -  plist文件写入
- (void)writePasswordToDefaultsWithUserName:(NSString *)username password:(NSString *)pwd
{
    NSUserDefaults * defaults = [NSUserDefaults standardUserDefaults];
    [defaults setValue:username forKey:@"username"];
    [defaults setValue:pwd forKey:@"password"];
    [defaults setBool:YES forKey:@"recordState"];
    [defaults synchronize];
}
-(void)deletePasswordInDefaultsWithUserName:(NSString *)username
{
    NSUserDefaults * defaults = [NSUserDefaults standardUserDefaults];
    [defaults setValue:username forKey:@"username"];
    [defaults removeObjectForKey:@"password"];
    [defaults setBool:NO forKey:@"recordState"];
    [defaults synchronize];
}

#pragma mark - 自定义alertView代理方法
- (void)alertView:(CustomAlertView *)alertView clickedButtonAtIndex:(NSInteger)buttonIndex {
    NSLog(@"%ld", (long)buttonIndex);
    if (buttonIndex == 1)
    {
        [self sendLinkToMail];
    }
}
#pragma mark - alertView
-(void)alertViewWithMessage:(NSString *)message
{
    UIAlertView * alertView = [[UIAlertView alloc]initWithTitle:@"温馨提示" message:message delegate:self cancelButtonTitle:@"知道了" otherButtonTitles:nil, nil];
    [alertView show];
}

- (void)didReceiveMemoryWarning {
    [super didReceiveMemoryWarning];
    
}
@end
#pragma mark - 获取uiwebview
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

