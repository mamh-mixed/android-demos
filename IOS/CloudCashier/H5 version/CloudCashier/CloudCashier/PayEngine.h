//
//  PayEngine.h
//  CloudCashier
//
//  Created by 司瑞华 on 15/4/23.
//  Copyright (c) 2015年 __VTPayment__. All rights reserved.
//

#import <Foundation/Foundation.h>

@interface PayEngine : NSObject

//
///**
// *  注册账户
// *
// *  @param userName       用户名           NSString类型
// *  @param password       密码             NSString类型
// *  @param succeedBlock   成功时的回调        返回NSDictionary类型
// */

+(void)registerAccountWithUserName:(NSString *)userName password:(NSString *)password succeedBlock:(void(^)(NSDictionary * receiveDict))complete;

//
///**
// *  登录账户
// *
// *  @param userName       用户名             NSString类型
// *  @param password       密码               NSString类型
// *  @param succeedBlock   成功时的回调        返回NSDictionary类型
// */

+(void)logPayViewWithUserName:(NSString *)userName password:(NSString *)password succeedBlock:(void(^)(NSDictionary * receiveDict))complete;

//
///**
// *  首次录入页面
// *
// *  @param accountName        开户行                NSString类型
// *  @param userName           用户名                NSString类型
// *  @param bankNum            银行卡号              NSString类型
// *  @param phoneNum           手机号                NSString类型
// *  @param succeedBlock       成功时的回调           返回NSDictionary类型
// */
+(void)firstLogPayViewWithUserName:(NSString *)userName password:(NSString *)password bankName:(NSString *)bankName accountName:(NSString *)accountName bankNum:(NSString *)bankNum phoneNum:(NSString *)phoneNum succeedBlock:(void(^)(NSDictionary * receiveDict))complete;

//
///**
// *  判断是否是邮箱
// * /
+(BOOL)validateEmail:(NSString *)email ;


//
/* *  根据邮箱返回验证码
 *    @param userName           用户名                NSString类型
 */
+(void)getAuthCodeWithUserName:(NSString *)userName succeedBlock:(void(^)(NSDictionary * receiveDict))complete;




//
/* *  激活账号
 *    @param userName           用户名                NSString类型
 */
+(void)activateAccountWithUserName:(NSString *)userName password:(NSString *)password succeedBlock:(void(^)(NSDictionary * receiveDict))complete;

//
///**
// *  忘记密码，重设密码
// *  @param userName           用户名                NSString类型
// *  @param code               验证码                NSString类型
// *  @param newPassword        新密码                NSString类型
// * /
+(void)resetPasswordWithUserName:(NSString *)username code:(NSString *)code newPassword:(NSString *)newpassword succeedBlock:(void(^)(NSDictionary * receiveDict))complete;


//
///**
// *  修改密码
// *  @param userName           用户名                NSString类型
// *  @param oldPassword        旧密码                NSString类型
// *  @param newPassword        新密码                NSString类型
// * /
+(void)updatePasswordWithUserName:(NSString *)username oldPassword:(NSString *)oldpassword newPassword:(NSString *)newpassword succeedBlock:(void(^)(NSDictionary * receiveDict))complete;


//
///**
// *  获取交易订单
// *  @param userName           用户名                 NSString类型
// *  @param password           密码                   NSString类型
// *  @param clientId           商户号                 NSString类型
// *  @param month              查询时间                NSString类型
// *  @param index              页码                   NSString类型
// * /

+(void)queryOrderWithUserName:(NSString *)username password:(NSString *)password clientId:(NSString *)clientid month:(NSString *)month index:(NSString *)index succeedBlock:(void(^)(NSDictionary * receiveDict))complete;



//
///**
// *  查询是否还有余额可以退款
// *  @param userName           用户名                NSString类型
// *  @param password           密码                  NSString类型
// *  @param clientId           商户号                NSString类型
// *  @param orderNum           订单号                NSString类型
// * /

+(void)checkBalanceWithUserName:(NSString *)username password:(NSString *)password clientId:(NSString *)clientid orderNum:(NSString *)orderNumStr succeedBlock:(void(^)(NSDictionary * receiveDict))complete;



//
///**
// *  修改账户
// *  @param username           用户名                NSString类型
// *  @param password           密码                  NSString类型
// *  @param bankopen           开户行                NSString类型
// *  @param payee              姓名                  NSString类型
// *  @param payeecard          银行卡号               NSString类型
// *  @param phonenum           电话号码               NSString类型
// * /

+(void)updateAccountWithUserName:(NSString *)username password:(NSString *)pwd bankOpen:(NSString *)bankopen payee:(NSString *)payee payeeCard:(NSString *)payeecard phoneNum:(NSString *)phonenum succeedBlock:(void(^)(NSDictionary * receiveDict))complete;


//
///**
// *  限额提升
// *  @param username           用户名                NSString类型
// *  @param password           密码                  NSString类型
// *  @param email              验证邮箱               NSString类型
// *  @param payee              姓名                  NSString类型
// *  @param phonenum           电话号码               NSString类型
// * /

+(void)limitinCreaseWithUserName:(NSString *)username password:(NSString *)pwd email:(NSString *)email payee:(NSString *)payee  phoneNum:(NSString *)phonenum succeedBlock:(void(^)(NSDictionary * receiveDict))complete;

//
///**
// *  获取一天内交易的金额量
// *  @param username           用户名                NSString类型
// *  @param password           密码                  NSString类型
// *  @param email              验证邮箱               NSString类型
// *  @param payee              姓名                  NSString类型
// *  @param phonenum           电话号码               NSString类型
// * /

+(void)getTotalWithUserName:(NSString *)username password:(NSString *)password clientId:(NSString *)clientid succeedBlock:(void(^)(NSDictionary * receiveDict))complete;
@end
