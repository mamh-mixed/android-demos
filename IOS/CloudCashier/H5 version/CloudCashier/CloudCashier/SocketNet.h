//
//  SocketNet.h
//  CloudCashierAPI
//
//  Created by 司瑞华 on 15/7/9.
//  Copyright (c) 2015年 __VTPayment__. All rights reserved.
//

#import <Foundation/Foundation.h>
#import "BackInfoDelegate.h"

@interface SocketNet : NSObject<NSStreamDelegate,BackInfoDelegate>

//字典内排序
+(NSString *)getSignStrWithDiction:(NSDictionary *)dict signKey:(NSString *)signKey;

//socket请求
+(void)socketWithTransitionStr:(NSString *)transitionStr withType:(NSInteger)type;

//向服务器传的参数
+(NSString *)pinyinSort:(NSDictionary *)dict signKey:(NSString *)signkey;

//传代理
+(void)setDelegate:(id<BackInfoDelegate>)delegate;

//获取密钥
+(void)byValueSignKey:(NSString *)signKeyStr;

//返回给用户的数据对象
+(BackParameter *)dataWithResultStr:(NSString *)str type:(NSInteger)typeNum;
@end
