//
//  CloudCashierAPI.m
//  CloudCashierAPI
//
//  Created by 司瑞华 on 15/7/9.
//  Copyright (c) 2015年 __VTPayment__. All rights reserved.
//

#import "CloudCashierAPI.h"
#import "VTSHAAndMD5.h"
#import "SocketNet.h"
#import <UIKit/UIKit.h>


static NSMutableArray * dataArray ;
static CloudCashierAPI * api;

@implementation Parameter


@end


@implementation CloudCashierAPI

//几个固定信息
+(void)registerInscd:(NSString *)inscdStr mchntid:(NSString *)mchntidStr signKey:(NSString *)signKeyStr terminalid:(NSString *)terminalidStr tradeFrom:(NSString *)tradeFrom
{
    [SocketNet byValueSignKey:signKeyStr];
    api = [[CloudCashierAPI alloc]init];
    dataArray = [[NSMutableArray alloc]initWithObjects:inscdStr,mchntidStr,signKeyStr,terminalidStr, tradeFrom,nil];
}


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
+(void)scannerPayWithpara:(Parameter *)para
{
    NSMutableDictionary * dic = [[NSMutableDictionary alloc] init];
    
    [dic setValue:para.scanCodeId forKey:@"scanCodeId"];
    [dic setValue:[dataArray objectAtIndex:1] forKey:@"mchntid"];
    [dic setValue:[dataArray objectAtIndex:4] forKey:@"tradeFrom"];
    [dic setValue:para.orderNum forKey:@"orderNum"];
    [dic setValue:@"PURC" forKey:@"busicd"];
    [dic setValue:[dataArray objectAtIndex:0] forKey:@"inscd"];
    
    [dic setValue:@"Q" forKey:@"txndir"];
    [dic setValue:[dataArray objectAtIndex:3] forKey:@"terminalid"];
    [dic setValue:para.currency forKey:@"currency"];
    if ([para.currency isEqualToString:@""] || ![para.currency isEqualToString:@"CNY"])
    {
        UIAlertView * alertView = [[UIAlertView alloc]initWithTitle:@"币种错误" message:nil delegate:self cancelButtonTitle:@"我知道了" otherButtonTitles:nil, nil];
        [alertView show];
    }else if(para.currency)
    {
        double txamt = [para.txamt doubleValue];
        NSMutableString * txamtStr = [NSMutableString stringWithFormat:@"%.2f",txamt];
        NSRange substr;
        substr = [txamtStr rangeOfString:@"."];
        if (substr.location != NSNotFound)
        {
            [txamtStr deleteCharactersInRange:substr];
            for (int i = 1; i < 13-txamtStr.length; )
            {
                [txamtStr insertString:@"0" atIndex:0];
                if (txamtStr.length > 12)
                {
                    break;
                }
            }
        }
        
        [dic setValue:txamtStr forKey:@"txamt"];
        [SocketNet socketWithTransitionStr:[SocketNet pinyinSort:dic signKey:[dataArray objectAtIndex:2]] withType:1];
    }
}

/**
 *  //预下单支付
 * chcd : 渠道
 * txamt : 订单金额
 * goodsInfo :商品信息
 * orderNum : 订单号
 * currency : 币种
 * **
 */
+(void)preOrderPayWithpara:(Parameter *)para
{
    NSMutableDictionary * dic = [[NSMutableDictionary alloc] init];
    [dic setValue:@"Q" forKey:@"txndir"];
    [dic setValue:@"PAUT" forKey:@"busicd"];
    [dic setValue:[dataArray objectAtIndex:0] forKey:@"inscd"];
    [dic setValue:para.chcd forKey:@"chcd"];
    [dic setValue:[dataArray objectAtIndex:4] forKey:@"tradeFrom"];
    [dic setValue:[dataArray objectAtIndex:1] forKey:@"mchntid"];
    [dic setValue:para.orderNum forKey:@"orderNum"];
    [dic setValue:para.currency forKey:@"currency"];
    [dic setValue:[dataArray objectAtIndex:3] forKey:@"terminalid"];
    
    if ([para.currency isEqualToString:@""] || ![para.currency isEqualToString:@"CNY"])
    {
        UIAlertView * alertView = [[UIAlertView alloc]initWithTitle:@"币种错误" message:nil delegate:self cancelButtonTitle:@"我知道了" otherButtonTitles:nil, nil];
        [alertView show];
    }else if(para.currency)
    {
        double txamt = [para.txamt doubleValue];
        NSMutableString * txamtStr = [NSMutableString stringWithFormat:@"%.2f",txamt];
        NSRange substr;
        substr = [txamtStr rangeOfString:@"."];
        if (substr.location != NSNotFound)
        {
            [txamtStr deleteCharactersInRange:substr];
            for (int i = 1; i < 13-txamtStr.length; )
            {
                [txamtStr insertString:@"0" atIndex:0];
                if (txamtStr.length > 12)
                {
                    break;
                }
            }
        }
        
        [dic setValue:txamtStr forKey:@"txamt"];
        [SocketNet socketWithTransitionStr:[SocketNet pinyinSort:dic signKey:[dataArray objectAtIndex:2]] withType:2];
    }

}

/**
 *  //查询订单
 * origOrderNum : 原订单号
 * **
 */
+(void)queryWithpara:(Parameter *)para
{
    NSMutableDictionary * dic = [[NSMutableDictionary alloc] init];
    [dic setValue:@"Q" forKey:@"txndir"];
    [dic setValue:@"INQY" forKey:@"busicd"];
    [dic setValue:[dataArray objectAtIndex:0] forKey:@"inscd"];
    [dic setValue:[dataArray objectAtIndex:4] forKey:@"tradeFrom"];
    [dic setValue:[dataArray objectAtIndex:1] forKey:@"mchntid"];
    [dic setValue:para.origOrderNum forKey:@"origOrderNum"];
    [dic setValue:[dataArray objectAtIndex:3] forKey:@"terminalid"];
    [SocketNet socketWithTransitionStr:[SocketNet pinyinSort:dic signKey:[dataArray objectAtIndex:2]] withType:3];
}

/**
 *  //撤销
 * orderNum : 订单号
 * origOrderNum : 原订单号
 *
 * **
 */
+(void)undoWithpara:(Parameter *)para
{
    NSMutableDictionary * dic = [[NSMutableDictionary alloc] init];
    [dic setValue:@"Q" forKey:@"txndir"];
    [dic setValue:@"VOID" forKey:@"busicd"];
    [dic setValue:[dataArray objectAtIndex:0] forKey:@"inscd"];
    [dic setValue:[dataArray objectAtIndex:4] forKey:@"tradeFrom"];
    [dic setValue:[dataArray objectAtIndex:1] forKey:@"mchntid"];
    [dic setValue:para.orderNum forKey:@"orderNum"];
    [dic setValue:para.origOrderNum forKey:@"origOrderNum"];
    [dic setValue:[dataArray objectAtIndex:3] forKey:@"terminalid"];
    [SocketNet socketWithTransitionStr:[SocketNet pinyinSort:dic signKey:[dataArray objectAtIndex:2]] withType:4];
}

/**
 *  //退款
 * txamt : 订单金额
 * orderNum : 订单号
 * origOrderNum : 原订单号
 * currency : 币种
 * **
 */
+(void)refundWithpara:(Parameter *)para
{
    NSMutableDictionary * dic = [[NSMutableDictionary alloc] init];
    [dic setValue:@"Q" forKey:@"txndir"];
    [dic setValue:@"REFD" forKey:@"busicd"];
    [dic setValue:[dataArray objectAtIndex:0] forKey:@"inscd"];
    [dic setValue:[dataArray objectAtIndex:1] forKey:@"mchntid"];
    [dic setValue:[dataArray objectAtIndex:4] forKey:@"tradeFrom"];
    [dic setValue:para.orderNum forKey:@"orderNum"];
    [dic setValue:para.origOrderNum forKey:@"origOrderNum"];
    [dic setValue:para.currency forKey:@"currency"];
    [dic setValue:[dataArray objectAtIndex:3] forKey:@"terminalid"];
    if ([para.currency isEqualToString:@""] || ![para.currency isEqualToString:@"CNY"])
    {
        UIAlertView * alertView = [[UIAlertView alloc]initWithTitle:@"币种错误" message:nil delegate:self cancelButtonTitle:@"我知道了" otherButtonTitles:nil, nil];
        [alertView show];
        
    }else if(para.currency)
    {
        double txamt = [para.txamt doubleValue];
        NSMutableString * txamtStr = [NSMutableString stringWithFormat:@"%.2f",txamt];
        NSRange substr;
        substr = [txamtStr rangeOfString:@"."];
        if (substr.location != NSNotFound)
        {
            [txamtStr deleteCharactersInRange:substr];
            for (int i = 1; i < 13-txamtStr.length; )
            {
                [txamtStr insertString:@"0" atIndex:0];
                if (txamtStr.length > 12)
                {
                    break;
                }
            }
        }
        
        [dic setValue:txamtStr forKey:@"txamt"];
        [SocketNet socketWithTransitionStr:[SocketNet pinyinSort:dic signKey:[dataArray objectAtIndex:2]] withType:5];
    }
}

/**
 *  //卡券核销
 * orderNum : 订单号
 * scanCodeId : 扫码号
 * **
 */
+(void)cardWithpara:(Parameter *)para
{
    NSMutableDictionary * dic = [[NSMutableDictionary alloc] init];
    [dic setValue:@"Q" forKey:@"txndir"];
    [dic setValue:@"CERI" forKey:@"busicd"];
    [dic setValue:[dataArray objectAtIndex:0] forKey:@"inscd"];
    [dic setValue:[dataArray objectAtIndex:1] forKey:@"mchntid"];
    [dic setValue:[dataArray objectAtIndex:4] forKey:@"tradeFrom"];
    [dic setValue:para.scanCodeId forKey:@"scanCodeId"];
    [dic setValue:para.orderNum forKey:@"orderNum"];
    [dic setValue:[dataArray objectAtIndex:3] forKey:@"terminalid"];
    [SocketNet socketWithTransitionStr:[SocketNet pinyinSort:dic signKey:[dataArray objectAtIndex:2]] withType:6];
}

//传送代理
+(void)transmitDelegate:(id<BackInfoDelegate>)delegate
{
    [SocketNet setDelegate:delegate];
}

@end
