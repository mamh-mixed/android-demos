//
//  RegisterTableDAO.m
//  888888
//
//  Created by 司瑞华 on 15/8/26.
//  Copyright (c) 2015年 __VTPayment__. All rights reserved.
//

#import "RegisterTable.h"

#import "DataBaseUtility.h"

@implementation RegisterTable

@end


@implementation RegisterTableDAO

+(void)insertObject:(RegisterTable * )object complete:(void(^)(NSString * isExists))complete
{
    FMDatabase * db = [DataBaseUtility getDataBase];
    if ([db open])
    {
        //遍历表单，看是否由menuName同名的记录，如果有，执行更新操作，改变menuNum;如果没有，插入一条记录
        //获得表中所有记录
        FMResultSet * set = [db executeQuery:@"select * from UserList"];
        //遍历这些记录
        while ([set next])
        {
            //获得记录中的menuName
            NSString * name = [set stringForColumn:@"username"];
            //比较menuName与object.menuName是否相同
            if ([name isEqualToString:object.username] == YES)
            {
                //返回
                if (complete)
                {
                    complete(@"exists");
                }
                return;
            }
        }
        [set close];
        //没有找到重名的  添加操作
          [db executeUpdate:@"insert into UserList (username,password,isUsed,time) values(?,?,?,?)",object.username,object.password,object.isUsed,object.time];
        if (complete)
        {
            complete(@"success");
        }
    }
    [db close];
}
+(NSMutableDictionary *)getObjectByName:(NSString *)name
{
    NSMutableDictionary * dict = [[NSMutableDictionary alloc]initWithCapacity:0];
    FMDatabase * db = [DataBaseUtility getDataBase];
    if ([db open])
    {
        FMResultSet * set = [db executeQuery:@"select * from UserList"];
        while ([set next])
        {
            NSString * nameStr = [set stringForColumn:@"username"];
            if ([name isEqualToString:nameStr] == YES)
            {
                [dict setValue:[set stringForColumn:@"isUsed"] forKey:@"isUsed"];
                [dict setValue:nameStr forKey:@"name"];
                [dict setValue:[set stringForColumn:@"password"] forKey:@"password"];
                [dict setValue:[set stringForColumn:@"time"] forKey:@"time"];
            }
        }
        [set close];
    }
    [db close];
    return dict;
}
+(void)changeisUsedByName:(NSString *)name
{
    FMDatabase * db = [DataBaseUtility getDataBase];
    if ([db open])
    {
        FMResultSet * set = [db executeQuery:@"select * from UserList"];
        //遍历这些记录
        while ([set next])
        {
            //获得记录中的menuName
            NSString * namestr = [set stringForColumn:@"username"];
            //比较menuName与object.menuName是否相同
            if ([namestr isEqualToString:name] != YES)
            {
                [db executeUpdate:@"update UserList set isUsed = ? WHERE username = ? ",@"0",namestr];
            }else
            {
                [db executeUpdate:@"update UserList set isUsed = ? WHERE username = ? ",@"1",namestr];
            }
        }
        [set close];
    }
    [db close];
}
+(NSString *)getNameWhoIsUsing
{
    FMDatabase *db=[DataBaseUtility getDataBase];
    if([db open]){
        FMResultSet *set=[db executeQuery:@"select * from UserList"];
        while ([set next]) {
            NSString *isUsed=[set stringForColumn:@"isUsed"];
            if([isUsed isEqualToString:@"1"])
            {
                return [set stringForColumn:@"username"];
            }
        }
        [set close];
    }
    [db close];
    return nil;
}
+(BOOL)isExistTheName:(NSString *)name
{
    FMDatabase *db=[DataBaseUtility getDataBase];
    if([db open]){
        FMResultSet *set=[db executeQuery:@"select * from UserList"];
        while ([set next]) {
            if ([[set stringForColumn:@"username"] isEqualToString:name]) {
                return YES;
            }
        }
    }
    return NO;
}
@end