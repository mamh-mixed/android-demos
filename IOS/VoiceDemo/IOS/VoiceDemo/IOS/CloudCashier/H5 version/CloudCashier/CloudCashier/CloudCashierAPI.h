//
//  CloudCashierAPI.h
//  CloudCashierAPI
//
//  Created by 司瑞华 on 15/7/9.
//  Copyright (c) 2015年 __VTPayment__. All rights reserved.
//

#import <Foundation/Foundation.h>
#import "BackInfoDelegate.h"

@interface Parameter : NSObject

@property (nonatomic,strong) NSString                *txamt;//订单金额
@property (nonatomic,strong) NSString                *orderNum;//订单号
@property (nonatomic,strong) NSString                *scanCodeId;//扫码号
@property (nonatomic,strong) NSString                *goodsInfo;//商品名称
@property (nonatomic,strong) NSString                *currency;//币种
@property (nonatomic,strong) NSString                *chcd;//渠道
@property (nonatomic,strong) NSString                *origOrderNum;//原订单号

@end



@interface CloudCashierAPI : NSObject

/*/! @brief 第一次调用云支付SDK时需要先向云支付注册以下几个参数
 *
 * Inscd     : 机构号，商户所属机构标识
 * mchntid   : 商户号 ，由讯联数据分配
 * signKey   : 双方约定的签名密钥
 * terminalid: 终端号
 */
+(void)registerInscd:(NSString *)inscdStr mchntid:(NSString *)mchntidStr signKey:(NSString *)signKeyStr terminalid:(NSString *)terminalidStr tradeFrom:(NSString *)tradeFrom;


/**
 *  //下单支付
 * txamt : 订单金额
 * orderNum : 订单号
 * scanCodeId : 扫码号
 * goodsInfo : 商品信息
 * currency: 币种
 *
 * **
 */
+(void)scannerPayWithpara:(Parameter *)para ;

/**
 *  //预下单支付
 * chcd : 渠道
 * txamt : 订单金额
 * goodsInfo :商品信息
 * orderNum : 订单号
 * currency : 币种
 * **
 */
+(void)preOrderPayWithpara:(Parameter *)para ;

/**
 *  //查询订单
 * origOrderNum : 原订单号
 * **
 */
+(void)queryWithpara:(Parameter *)para ;

/**
 *  //撤销
 * orderNum : 订单号
 * origOrderNum : 原订单号
 *
 * **
 */
+(void)undoWithpara:(Parameter *)para ;

/**
 *  //退款
 * txamt : 订单金额
 * orderNum : 订单号
 * origOrderNum : 原订单号
 * currency : 币种
 * **
 */
+(void)refundWithpara:(Parameter *)para ;

/**
 *  //卡券核销
 * orderNum : 订单号
 * scanCodeId : 扫码号
 * **
 */
+(void)cardWithpara:(Parameter *)para ;



+(void)transmitDelegate:(id<BackInfoDelegate>)delegate;

@end
