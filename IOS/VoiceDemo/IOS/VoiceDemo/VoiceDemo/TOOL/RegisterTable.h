//
//  RegisterTableDAO.h
//  888888
//
//  Created by 司瑞华 on 15/8/26.
//  Copyright (c) 2015年 __VTPayment__. All rights reserved.
//

#import <Foundation/Foundation.h>

@interface RegisterTable : NSObject

@property(nonatomic,strong) NSString    * username , * password , * isUsed , * time;

@end

@interface RegisterTableDAO : NSObject

+(void)insertObject:(RegisterTable * )object complete:(void(^)(NSString * isExists))complete ;

+(NSMutableDictionary *)getObjectByName:(NSString *)name;

+(void)changeisUsedByName:(NSString *)name;

//判断正在使用的用户名
+(NSString *)getNameWhoIsUsing;

//用户名是否存在
+(BOOL)isExistTheName:(NSString *)name;

@end