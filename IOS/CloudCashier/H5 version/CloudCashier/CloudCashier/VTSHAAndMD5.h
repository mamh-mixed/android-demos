//
//  VTSHAAndMD5.h
//  CloudCashierAPI
//
//  Created by 司瑞华 on 15/7/9.
//  Copyright (c) 2015年 __VTPayment__. All rights reserved.
//

#import <Foundation/Foundation.h>

@interface VTSHAAndMD5 : NSObject

+ (NSString*) sha1: (NSString *) inPutText;
+(NSString *) md5: (NSString *) inPutText;

@end
