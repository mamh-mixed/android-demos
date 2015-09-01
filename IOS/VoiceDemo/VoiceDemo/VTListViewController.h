//
//  VTListViewController.h
//  VoiceDemo
//
//  Created by 司瑞华 on 15/8/26.
//  Copyright (c) 2015年 __VTPayment__. All rights reserved.
//

#import <UIKit/UIKit.h>

@interface VTListViewController : UIViewController

@property(nonatomic,assign)float                       countMoney;
@property(nonatomic,assign)int                         countNum;
@property(nonatomic,strong)NSMutableArray              * cellContentArray;

@property(nonatomic,strong)NSMutableArray              * dataArray;//存放cell的数据

@end
