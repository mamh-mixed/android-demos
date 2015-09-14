//
//  DataBaseUtility.m
//  OrderDishes
//
//  Created by ZhenFan on 14-8-21.
//  Copyright (c) 2014年 ZhenFan. All rights reserved.
//

#import "DataBaseUtility.h"
static FMDatabase * _db = nil ;
@implementation DataBaseUtility
+(FMDatabase*)getDataBase
{
    if (_db == nil)
    {
         NSString * toPath = [NSHomeDirectory() stringByAppendingPathComponent:@"Documents/database.sqlite3"];
        _db = [[FMDatabase alloc]initWithPath:toPath];
    }
    return _db ;
}
@end








