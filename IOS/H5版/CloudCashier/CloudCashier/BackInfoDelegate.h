//
//  BackInfoDelegate.h
//  CloudCashierAPI
//
//  Created by 司瑞华 on 15/7/9.
//  Copyright (c) 2015年 __VTPayment__. All rights reserved.
//

#import <Foundation/Foundation.h>

@interface BackParameter : NSObject

@property (nonatomic,strong) NSString                * txndir;//交易方向
@property (nonatomic,strong) NSString                * busicd;//交易类型
@property (nonatomic,strong) NSString                * respcd;//交易结果
@property (nonatomic,strong) NSString                * inscd;//机构号
@property (nonatomic,strong) NSString                * chcd;//渠道
@property (nonatomic,strong) NSString                * mchntid;//商户号
@property (nonatomic,strong) NSString                * txamt;//订单金额
@property (nonatomic,strong) NSString                * channelOrderNum;//渠道交易号
@property (nonatomic,strong) NSString                * consumerAccount;//渠道账号
@property (nonatomic,strong) NSString                * consumerId;//渠道账号ID
@property (nonatomic,strong) NSString                * errorDetail;//错误信息
@property (nonatomic,strong) NSString                * origOrderNum;//原订单号
@property (nonatomic,strong) NSString                * orderNum;//订单号
@property (nonatomic,strong) NSString                * qrcode;//二维码信息
@property (nonatomic,strong) NSString                * sign;//签名
@property (nonatomic,strong) NSString                * chcdDiscount;//渠道优惠
@property (nonatomic,strong) NSString                * merDiscount;//商户优惠
@property (nonatomic,strong) NSString                * scanCodeId;//扫码号
@property (nonatomic,strong) NSString                * cardId;//卡券类型
@property (nonatomic,strong) NSString                * cardInfo;//卡券详情
@property (nonatomic,assign) NSInteger               tag;//判断是哪种类型发出的请求



@end


@protocol BackInfoDelegate <NSObject>

@optional

-(void)getResultDataWithBackParameter:(BackParameter *)backData errorCode:(NSInteger) errorNum;


@end