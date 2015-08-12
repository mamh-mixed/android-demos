//
//  VTGenQRViewController.h
//  CloudCashier
//
//  Created by 司瑞华 on 15/7/17.
//  Copyright (c) 2015年 __VTPayment__. All rights reserved.
//

#import <UIKit/UIKit.h>

@interface VTGenQRViewController : UIViewController


@property (strong ,nonatomic) NSString                        *chcd;//渠道号
@property (strong ,nonatomic) NSString                        *amount;
@property (strong ,nonatomic) NSString                        *qrInfo;
@property (strong ,nonatomic) NSString                        * qureyOrderNum;//用来查询订单的订单号

@end
