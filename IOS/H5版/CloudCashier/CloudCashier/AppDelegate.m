//
//  AppDelegate.m
//  CloudCashier
//
//  Created by 司瑞华 on 15/7/7.
//  Copyright (c) 2015年 __VTPayment__. All rights reserved.
//

#import "AppDelegate.h"
#import "VTViewController.h"

@interface AppDelegate ()<NSXMLParserDelegate,UIAlertViewDelegate>
{
    
}

@end

@implementation AppDelegate


- (BOOL)application:(UIApplication *)application didFinishLaunchingWithOptions:(NSDictionary *)launchOptions {
    // Override point for customization after application launch.
    self.window = [[UIWindow alloc]initWithFrame:[[UIScreen mainScreen] bounds]];
   // self.window.backgroundColor = [UIColor colorWithRed:30/255.0 green:190/255.0  blue:214/255.0  alpha:1];
    VTViewController * vc = [[VTViewController alloc]init];
    self.window.rootViewController = vc;
    
       
    [self.window makeKeyAndVisible];
    return YES;
}



@end










