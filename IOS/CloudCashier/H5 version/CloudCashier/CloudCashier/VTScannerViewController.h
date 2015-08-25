//
//  VTScannerViewController.h
//  CloudCashier
//
//  Created by 司瑞华 on 15/7/7.
//  Copyright (c) 2015年 __VTPayment__. All rights reserved.
//

#import <UIKit/UIKit.h>

@interface VTScannerViewController : UIViewController

@property (nonatomic,strong)NSString                       * userName;
@property (nonatomic,strong)NSString                       * password;
@property (nonatomic,strong)NSDictionary                   * dictionary;
@property (nonatomic,assign)NSInteger                      whichPage;

@end
