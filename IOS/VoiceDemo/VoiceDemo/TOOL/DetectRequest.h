//
//  DetectRequest.h
//  VoiceDemo
//
//  Created by 黄达能 on 15/9/7.
//  Copyright (c) 2015年 __VTPayment__. All rights reserved.
//

#import <Foundation/Foundation.h>
#import "VTPayViewController.h"
@interface DetectRequest : NSObject

-(void)connectionNet:(NSString *)VoicePath;

@property (nonatomic,strong) VTPayViewController *viewController;

@end
