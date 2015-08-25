//
//  PayEngine.m
//  CloudCashier
//
//  Created by 司瑞华 on 15/4/23.
//  Copyright (c) 2015年 __VTPayment__. All rights reserved.
//

#import "PayEngine.h"
#import "SPostRequest.h"
#import "VTSHAAndMD5.h"


@implementation PayEngine


#pragma mark - 注册
+(void)registerAccountWithUserName:(NSString *)userName password:(NSString *)password succeedBlock:(void(^)(NSDictionary * receiveDict))complete
{    
    NSString *str = [NSString stringWithFormat:@"%@register",urlHost];
    ///设置标准时间格式
    NSDate * dateTime = [NSDate date];
    NSDateFormatter * dateFormatter = [[NSDateFormatter alloc]init];
    [dateFormatter setDateFormat:@"YYYYMMddHHmmss"];
    NSString * formatterTime = [dateFormatter stringFromDate:dateTime];
   ///// NSLog(@"formatterTime===\n %@",formatterTime);
    
    //对密码进行加密
    NSString * pwEncryptionStr = [VTSHAAndMD5 md5:password];
    
    //对参数进行SHA-1加密
    NSString * paraStr = [NSString stringWithFormat:@"password=%@&transtime=%@&username=%@%@",pwEncryptionStr,formatterTime,userName,ENKey];
    NSString * paraEncryptionStr = [VTSHAAndMD5 sha1:paraStr];
   //请求注册的参数
    NSMutableDictionary *dict =[[NSMutableDictionary alloc]initWithObjectsAndKeys:userName,@"username",pwEncryptionStr,@"password",formatterTime,@"transtime",paraEncryptionStr,@"sign", nil];
    
    [SPostRequest postRequestWithURL:str paramters:dict succeed:^(NSDictionary *receiveDict)
    {
        if (complete)
        {
            NSLog(@"输出 注册接收的数据 receiveDict------%@",receiveDict);
            complete(receiveDict);
        }
        
    } failed:^(NSString *dataString) {
        
        NSLog(@"【post】请求失败 的数据===\n %@",dataString);
        
    }];
}
#pragma mark - 登录
+(void)logPayViewWithUserName:(NSString *)userName password:(NSString *)password succeedBlock:(void(^)(NSDictionary * receiveDict))complete
{
    NSString *str = [NSString stringWithFormat:@"%@login",urlHost];
    ///设置标准时间格式
    NSDate * dateTime = [NSDate date];
    NSDateFormatter * dateFormatter = [[NSDateFormatter alloc]init];
    [dateFormatter setDateFormat:@"YYYYMMddHHmmss"];
    NSString * formatterTime = [dateFormatter stringFromDate:dateTime];
    ///// NSLog(@"formatterTime===\n %@",formatterTime);
    
    //对密码进行加密
    NSString * pwEncryptionStr = [VTSHAAndMD5 md5:password];
    
    //对参数进行SHA-1加密
    NSString * paraStr = [NSString stringWithFormat:@"password=%@&transtime=%@&username=%@%@",pwEncryptionStr,formatterTime,userName,ENKey];
    NSString * paraEncryptionStr = [VTSHAAndMD5 sha1:paraStr];
    //请求注册的参数
    NSMutableDictionary *dict =[[NSMutableDictionary alloc]initWithObjectsAndKeys:userName,@"username",pwEncryptionStr,@"password",formatterTime,@"transtime",paraEncryptionStr,@"sign", nil];
    
    [SPostRequest postRequestWithURL:str paramters:dict succeed:^(NSDictionary *receiveDict)
     {
         if (complete)
         {
             complete(receiveDict);
         }
         
     } failed:^(NSString *dataString) {
         
         NSLog(@"【post】请求失败 的数据===\n %@",dataString);
         
     }];

}
#pragma mark - 首次录入信息
+(void)firstLogPayViewWithUserName:(NSString *)userName password:(NSString *)password bankName:(NSString *)bankName accountName:(NSString *)accountName bankNum:(NSString *)bankNum phoneNum:(NSString *)phoneNum succeedBlock:(void(^)(NSDictionary * receiveDict))complete
{
    NSString *str = [NSString stringWithFormat:@"%@improveinfo",urlHost];
    ///设置标准时间格式
    NSDate * dateTime = [NSDate date];
    NSDateFormatter * dateFormatter = [[NSDateFormatter alloc]init];
    [dateFormatter setDateFormat:@"YYYYMMddHHmmss"];
    NSString * formatterTime = [dateFormatter stringFromDate:dateTime];
    
    // NSLog(@"formatterTime===\n %@",formatterTime);
    
    //对密码进行加密
    NSString * pwEncryptionStr = [VTSHAAndMD5 md5:password];
    
    //对参数进行SHA-1加密
    NSString * paraStr = [NSString stringWithFormat:@"bank_open=%@&password=%@&payee=%@&payee_card=%@&phone_num=%@&transtime=%@&username=%@%@",bankName,pwEncryptionStr,accountName,bankNum,phoneNum,formatterTime,userName,ENKey];
    NSString * paraEncryptionStr = [VTSHAAndMD5 sha1:paraStr];
    //请求注册的参数
    NSMutableDictionary *dict =[[NSMutableDictionary alloc]initWithObjectsAndKeys:userName,@"username",pwEncryptionStr,@"password",bankName,@"bank_open",accountName,@"payee",bankNum,@"payee_card",phoneNum,@"phone_num",formatterTime,@"transtime",paraEncryptionStr,@"sign", nil];
    
    [SPostRequest postRequestWithURL:str paramters:dict succeed:^(NSDictionary *receiveDict)
     {
         if (complete)
         {
             complete(receiveDict);
         }
         
     } failed:^(NSString *dataString) {
         
         NSLog(@"【post】请求失败 的数据===\n %@",dataString);
         
     }];

}
#pragma mark - 判断是否是邮箱
+(BOOL)validateEmail:(NSString *)email
{
    NSString *emailRegex = @"[A-Z0-9a-z._%+-]+@[A-Za-z0-9.-]+\\.[A-Za-z]{2,4}";
    NSPredicate *emailTest = [NSPredicate predicateWithFormat:@"SELF MATCHES %@", emailRegex];
    return [emailTest evaluateWithObject:email];
}

#pragma mark - 忘记密码输入邮箱获得验证码
+(void)getAuthCodeWithUserName:(NSString *)userName succeedBlock:(void(^)(NSDictionary * receiveDict))complete
{
    NSString *str = [NSString stringWithFormat:@"%@forgetpassword",urlHost];
    NSDate * dateTime = [NSDate date];
    NSDateFormatter * dateFormatter = [[NSDateFormatter alloc]init];
    [dateFormatter setDateFormat:@"YYYYMMddHHmmss"];
    NSString * formatterTime = [dateFormatter stringFromDate:dateTime];
    
    //对参数进行SHA-1加密
    NSString * paraStr = [NSString stringWithFormat:@"transtime=%@&username=%@%@",formatterTime,userName,ENKey];
    
    NSString * paraEncryptionStr = [VTSHAAndMD5 sha1:paraStr];
    //请求注册的参数
    NSMutableDictionary *dict =[[NSMutableDictionary alloc]initWithObjectsAndKeys:userName,@"username",formatterTime,@"transtime",paraEncryptionStr,@"sign", nil];
    
    [SPostRequest postRequestWithURL:str paramters:dict succeed:^(NSDictionary *receiveDict)
     {
         if (complete)
         {
             complete(receiveDict);
         }
     } failed:^(NSString *dataString) {
         
         NSLog(@"【post】请求失败 的数据===\n %@",dataString);
         
     }];

}
#pragma mark - 忘记密码重设密码
+(void)resetPasswordWithUserName:(NSString *)username code:(NSString *)code newPassword:(NSString *)newpassword succeedBlock:(void(^)(NSDictionary * receiveDict))complete
{
    NSString *str = [NSString stringWithFormat:@"%@resetpassword",urlHost];
    NSDate * dateTime = [NSDate date];
    NSDateFormatter * dateFormatter = [[NSDateFormatter alloc]init];
    [dateFormatter setDateFormat:@"YYYYMMddHHmmss"];
    NSString * formatterTime = [dateFormatter stringFromDate:dateTime];
    //对密码进行加密
    NSString * pwEncryptionStr = [VTSHAAndMD5 md5:newpassword];
    
    //对参数进行SHA-1加密
    NSString * paraStr = [NSString stringWithFormat:@"code=%@&newpassword=%@&transtime=%@&username=%@%@",code,pwEncryptionStr,formatterTime,username,ENKey];
    
    NSString * paraEncryptionStr = [VTSHAAndMD5 sha1:paraStr];
    //请求注册的参数
    NSMutableDictionary *dict =[[NSMutableDictionary alloc]initWithObjectsAndKeys:username,@"username",formatterTime,@"transtime",paraEncryptionStr,@"sign",code,@"code",pwEncryptionStr,@"newpassword", nil];
    
    [SPostRequest postRequestWithURL:str paramters:dict succeed:^(NSDictionary *receiveDict)
     {
         if (complete)
         {
             complete(receiveDict);
         }
     } failed:^(NSString *dataString) {
         
         NSLog(@"【post】请求失败 的数据===\n %@",dataString);
         
     }];

}


#pragma mark - 修改密码
+(void)updatePasswordWithUserName:(NSString *)username oldPassword:(NSString *)oldpassword newPassword:(NSString *)newpassword succeedBlock:(void(^)(NSDictionary * receiveDict))complete
{
    NSString *str = [NSString stringWithFormat:@"%@updatepassword",urlHost];
    NSDate * dateTime = [NSDate date];
    NSDateFormatter * dateFormatter = [[NSDateFormatter alloc]init];
    [dateFormatter setDateFormat:@"YYYYMMddHHmmss"];
    NSString * formatterTime = [dateFormatter stringFromDate:dateTime];
    //对密码进行加密
    NSString * npwEncryptionStr = [VTSHAAndMD5 md5:newpassword];
    NSString * opwEncStr = [VTSHAAndMD5 md5:oldpassword];
    
    //对参数进行SHA-1加密
    NSString * paraStr = [NSString stringWithFormat:@"newpassword=%@&oldpassword=%@&transtime=%@&username=%@%@",npwEncryptionStr,opwEncStr,formatterTime,username,ENKey];
    
    NSString * paraEncryptionStr = [VTSHAAndMD5 sha1:paraStr];
    //请求注册的参数
    NSMutableDictionary *dict =[[NSMutableDictionary alloc]initWithObjectsAndKeys:username,@"username",formatterTime,@"transtime",paraEncryptionStr,@"sign",opwEncStr,@"oldpassword",npwEncryptionStr,@"newpassword", nil];
    
    [SPostRequest postRequestWithURL:str paramters:dict succeed:^(NSDictionary *receiveDict)
     {
         if (complete)
         {
             complete(receiveDict);
         }
     } failed:^(NSString *dataString) {
         
         NSLog(@"【post】请求失败 的数据===\n %@",dataString);
         
     }];

}
#pragma mark - 查询账单
+(void)queryOrderWithUserName:(NSString *)username password:(NSString *)password clientId:(NSString *)clientid month:(NSString *)month index:(NSString *)index succeedBlock:(void(^)(NSDictionary * receiveDict))complete
{
    NSString *str = [NSString stringWithFormat:@"%@bill",urlHost];
    NSDate * dateTime = [NSDate date];
    NSDateFormatter * dateFormatter = [[NSDateFormatter alloc]init];
    [dateFormatter setDateFormat:@"YYYYMMddHHmmss"];
    NSString * formatterTime = [dateFormatter stringFromDate:dateTime];
    //对密码进行加密
    NSString * pwEncStr = [VTSHAAndMD5 md5:password];
    
    
    //对参数进行SHA-1加密
    NSString * paraStr = [NSString stringWithFormat:@"clientid=%@&index=%@&month=%@&password=%@&status=all&transtime=%@&username=%@%@",clientid,index,month,pwEncStr,formatterTime,username,ENKey];
    
    NSString * paraEncryptionStr = [VTSHAAndMD5 sha1:paraStr];
    //请求注册的参数
    NSMutableDictionary *dict =[[NSMutableDictionary alloc]initWithObjectsAndKeys:username,@"username",formatterTime,@"transtime",paraEncryptionStr,@"sign",pwEncStr,@"password",clientid,@"clientid",month,@"month",index,@"index",@"all",@"status", nil];
    
    [SPostRequest postRequestWithURL:str paramters:dict succeed:^(NSDictionary *receiveDict)
     {
         if (complete)
         {
             complete(receiveDict);
         }
     } failed:^(NSString *dataString) {
         
         NSLog(@"【post】请求失败 的数据===\n %@",dataString);
         
     }];

}
#pragma mark - 激活账号
+(void)activateAccountWithUserName:(NSString *)userName password:(NSString *)password succeedBlock:(void(^)(NSDictionary * receiveDict))complete
{
    NSString *str = [NSString stringWithFormat:@"%@request_activate",urlHost];
    NSDate * dateTime = [NSDate date];
    NSDateFormatter * dateFormatter = [[NSDateFormatter alloc]init];
    [dateFormatter setDateFormat:@"YYYYMMddHHmmss"];
    NSString * formatterTime = [dateFormatter stringFromDate:dateTime];
     NSString * pwEncryptionStr = [VTSHAAndMD5 md5:password];
    //对参数进行SHA-1加密
    NSString * paraStr = [NSString stringWithFormat:@"password=%@&transtime=%@&username=%@%@",pwEncryptionStr,formatterTime,userName,ENKey];
    NSString * paraEncryptionStr = [VTSHAAndMD5 sha1:paraStr];
    //请求注册的参数
    NSMutableDictionary *dict =[[NSMutableDictionary alloc]initWithObjectsAndKeys:userName,@"username",pwEncryptionStr,@"password",formatterTime,@"transtime",paraEncryptionStr,@"sign", nil];
    
    [SPostRequest postRequestWithURL:str paramters:dict succeed:^(NSDictionary *receiveDict)
     {
         if (complete)
         {
             complete(receiveDict);
         }
         
     } failed:^(NSString *dataString) {
         
         NSLog(@"【post】请求失败 的数据===\n %@",dataString);
         
     }];

}
#pragma mark - 查询是否还有余额可以退款
+(void)checkBalanceWithUserName:(NSString *)username password:(NSString *)password clientId:(NSString *)clientid orderNum:(NSString *)orderNumStr succeedBlock:(void(^)(NSDictionary * receiveDict))complete
{
    NSString *str = [NSString stringWithFormat:@"%@getrefd",urlHost];
    NSDate * dateTime = [NSDate date];
    NSDateFormatter * dateFormatter = [[NSDateFormatter alloc]init];
    [dateFormatter setDateFormat:@"YYYYMMddHHmmss"];
    NSString * formatterTime = [dateFormatter stringFromDate:dateTime];
    //对密码进行加密
    NSString * pwEncStr = [VTSHAAndMD5 md5:password];
    
    
    //对参数进行SHA-1加密
    NSString * paraStr = [NSString stringWithFormat:@"clientid=%@&orderNum=%@&password=%@&transtime=%@&username=%@%@",clientid,orderNumStr,pwEncStr,formatterTime,username,ENKey];
    
    NSString * paraEncryptionStr = [VTSHAAndMD5 sha1:paraStr];
    //请求注册的参数
    NSMutableDictionary *dict =[[NSMutableDictionary alloc]initWithObjectsAndKeys:username,@"username",formatterTime,@"transtime",paraEncryptionStr,@"sign",pwEncStr,@"password",clientid,@"clientid",orderNumStr,@"orderNum", nil];
    
    [SPostRequest postRequestWithURL:str paramters:dict succeed:^(NSDictionary *receiveDict)
     {
         if (complete)
         {
             complete(receiveDict);
         }
     } failed:^(NSString *dataString) {
         
         NSLog(@"【post】请求失败 的数据===\n %@",dataString);
         
     }];
}

#pragma mark - 修改账户页面
+(void)updateAccountWithUserName:(NSString *)username password:(NSString *)pwd bankOpen:(NSString *)bankopen payee:(NSString *)payee payeeCard:(NSString *)payeecard phoneNum:(NSString *)phonenum succeedBlock:(void(^)(NSDictionary * receiveDict))complete
{
    NSString *str = [NSString stringWithFormat:@"%@updateinfo",urlHost];
    NSDate * dateTime = [NSDate date];
    NSDateFormatter * dateFormatter = [[NSDateFormatter alloc]init];
    [dateFormatter setDateFormat:@"YYYYMMddHHmmss"];
    NSString * formatterTime = [dateFormatter stringFromDate:dateTime];
    NSString * pwEncryptionStr = [VTSHAAndMD5 md5:pwd];
    //对参数进行SHA-1加密
    NSString * paraStr = [NSString stringWithFormat:@"bank_open=%@&password=%@&payee=%@&payee_card=%@&phone_num=%@&transtime=%@&username=%@%@",bankopen,pwEncryptionStr,payee,payeecard,phonenum,formatterTime,username,ENKey];
    NSString * paraEncryptionStr = [VTSHAAndMD5 sha1:paraStr];
    //请求注册的参数
    NSMutableDictionary *dict =[[NSMutableDictionary alloc]initWithObjectsAndKeys:username,@"username",pwEncryptionStr,@"password",bankopen,@"bank_open",payee,@"payee",payeecard,@"payee_card",phonenum,@"phone_num",formatterTime,@"transtime",paraEncryptionStr,@"sign", nil];
    
    [SPostRequest postRequestWithURL:str paramters:dict succeed:^(NSDictionary *receiveDict)
     {
         if (complete)
         {
             complete(receiveDict);
         }
         
     } failed:^(NSString *dataString) {
         
         NSLog(@"【post】请求失败 的数据===\n %@",dataString);
         
     }];

}
#pragma mark - 限额提升
+(void)limitinCreaseWithUserName:(NSString *)username password:(NSString *)pwd email:(NSString *)email payee:(NSString *)payee  phoneNum:(NSString *)phonenum succeedBlock:(void(^)(NSDictionary * receiveDict))complete
{
    NSString *str = [NSString stringWithFormat:@"%@limitincrease",urlHost];
    NSDate * dateTime = [NSDate date];
    NSDateFormatter * dateFormatter = [[NSDateFormatter alloc]init];
    [dateFormatter setDateFormat:@"YYYYMMddHHmmss"];
    NSString * formatterTime = [dateFormatter stringFromDate:dateTime];
    NSString * pwEncryptionStr = [VTSHAAndMD5 md5:pwd];
    //对参数进行SHA-1加密
    NSString * paraStr = [NSString stringWithFormat:@"email=%@&password=%@&payee=%@&phone_num=%@&transtime=%@&username=%@%@",email,pwEncryptionStr,payee,phonenum,formatterTime,username,ENKey];
    NSString * paraEncryptionStr = [VTSHAAndMD5 sha1:paraStr];
    //请求注册的参数
    NSMutableDictionary *dict =[[NSMutableDictionary alloc]initWithObjectsAndKeys:username,@"username",pwEncryptionStr,@"password",email,@"email",payee,@"payee",phonenum,@"phone_num",formatterTime,@"transtime",paraEncryptionStr,@"sign", nil];
    
    [SPostRequest postRequestWithURL:str paramters:dict succeed:^(NSDictionary *receiveDict)
     {
         if (complete)
         {
             complete(receiveDict);
         }
         
     } failed:^(NSString *dataString) {
         
         NSLog(@"【post】请求失败 的数据===\n %@",dataString);
         
     }];
}

#pragma mark - 获取一日内交易的总金额
+(void)getTotalWithUserName:(NSString *)username password:(NSString *)password clientId:(NSString *)clientid succeedBlock:(void(^)(NSDictionary * receiveDict))complete
{
    NSString *str = [NSString stringWithFormat:@"%@getTotal",urlHost];
    NSDate * dateTime = [NSDate date];
    NSDateFormatter * dateFormatter = [[NSDateFormatter alloc]init];
    [dateFormatter setDateFormat:@"YYYYMMddHHmmss"];
    NSString * formatterTime = [dateFormatter stringFromDate:dateTime];
    
    [dateFormatter setDateFormat:@"YYYYMMdd"];
    NSString * nowDate = [dateFormatter stringFromDate:dateTime];
    NSLog(@"当天时间---%@------当天日期----%@",formatterTime,nowDate);
    //对密码进行加密
    NSString * pwEncStr = [VTSHAAndMD5 md5:password];
    
    
    //对参数进行SHA-1加密
    NSString * paraStr = [NSString stringWithFormat:@"clientid=%@&date=%@&password=%@&transtime=%@&username=%@%@",clientid,nowDate,pwEncStr,formatterTime,username,ENKey];
    
    NSString * paraEncryptionStr = [VTSHAAndMD5 sha1:paraStr];
    //请求注册的参数
    NSMutableDictionary *dict =[[NSMutableDictionary alloc]initWithObjectsAndKeys:username,@"username",formatterTime,@"transtime",paraEncryptionStr,@"sign",pwEncStr,@"password",clientid,@"clientid",nowDate,@"date", nil];
    
    [SPostRequest postRequestWithURL:str paramters:dict succeed:^(NSDictionary *receiveDict)
     {
         if (complete)
         {
             complete(receiveDict);
         }
     } failed:^(NSString *dataString) {
         
         NSLog(@"【post】请求失败 的数据===\n %@",dataString);
         
     }];

}


@end

















