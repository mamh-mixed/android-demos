//
//  Request.h
//  VoiceDemo
//
//  Created by 黄达能 on 15/9/6.
//  Copyright (c) 2015年 __VTPayment__. All rights reserved.
//

#import <Foundation/Foundation.h>

@interface VTConnectionRequest : NSURLConnection

@property (nonatomic , assign) NSInteger                            tag;

@end

@interface Request : NSObject  //训练模式下

@property (nonatomic , assign) NSInteger                          successTimes;//已经成功发送多少次Request请求（不同的界面时 需要先设置为0）
+(Request *)sharedRequest;//单例

-(void)connectionNet:(NSArray *)VoicePath andUserKey:(NSString *)UserKey;

@end
