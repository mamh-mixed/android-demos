//
//  DataBaseUtility.h
//  OrderDishes
//
//  Created by ZhenFan on 14-8-21.
//  Copyright (c) 2014年 ZhenFan. All rights reserved.
//

#import <Foundation/Foundation.h>
#import "FMDatabase.h"
#import "FMDatabaseAdditions.h"
#import "FMResultSet.h"
@interface DataBaseUtility : NSObject
+(FMDatabase*)getDataBase ;
@end








