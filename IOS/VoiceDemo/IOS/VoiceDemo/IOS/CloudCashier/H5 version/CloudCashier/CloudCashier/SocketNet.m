//
//  SocketNet.m
//  CloudCashierAPI
//
//  Created by 司瑞华 on 15/7/9.
//  Copyright (c) 2015年 __VTPayment__. All rights reserved.
//

#import "SocketNet.h"
#import "VTSHAAndMD5.h"
#import "VTSHAAndMD5.h"


static SocketNet * socketNet ;
static id<BackInfoDelegate> _delegate;
static NSString * signKey;
NSInputStream                   * inputStream;
NSInteger typeNum;

@implementation SocketNet
//字典内排序
+(NSString *)getSignStrWithDiction:(NSDictionary *)dict signKey:(NSString *)signKey
{
    NSMutableArray * array = [NSMutableArray arrayWithArray:dict.allKeys];
    NSMutableString * encryptionStr = [[NSMutableString alloc]initWithFormat:@""] ;
    for (int i =0; i<array.count-1; i++)
    {
        for (int j = 0; j <array.count-1-i; j++)
        {     //每趟比较的次数
            if ([array[j] compare: array[j + 1]] == NSOrderedDescending)
            {
                
                NSString * temp = array[j];
                array[j] = array[j+1];
                array[j+1] = temp;
            }
        }
    }
    for (int i = 0; i < array.count; i++)
    {
        NSString * valueStr;
        NSString * key = [array objectAtIndex:i];
        if ([key isEqualToString:@"goodsInfo"])
        {
            NSString * info = [dict objectForKey:key];
            NSString * goods = [info stringByReplacingPercentEscapesUsingEncoding:NSUTF8StringEncoding];
            valueStr = goods;
        }else
        {
            valueStr = [dict objectForKey:key];
        }
        NSString * KVStr = [NSString stringWithFormat:@"%@=%@",key,valueStr];
        if (i == array.count-1)
        {
            [encryptionStr insertString:[NSString stringWithFormat:@"%@%@",KVStr,signKey] atIndex:encryptionStr.length];
        }else
        {
            [encryptionStr insertString:[NSString stringWithFormat:@"%@&",KVStr] atIndex:encryptionStr.length];
        }
    }
    
    NSString * sign = [VTSHAAndMD5 sha1:encryptionStr];
    return sign;

}

//socket请求
+(void)socketWithTransitionStr:(NSString *)transitionStr withType:(NSInteger)type
{
    socketNet = [[SocketNet alloc]init];
    
    typeNum = type;
    // 创建CF下的读入流
    CFReadStreamRef readStream;
    // 创建CF下的写出流
    CFWriteStreamRef writeStream;
    UInt32  port = SOCKETPORT;
    NSString * host = @"211.147.72.70";
    // 创建流
    CFStreamCreatePairWithSocketToHost(NULL, (__bridge CFStringRef)(host), port, &readStream, &writeStream);
    
    
    NSOutputStream                  * outputStream;
    
    // 将CFXXX流和NSXXX流建立对应关系
    inputStream = (__bridge NSInputStream *)(readStream);
    outputStream = (__bridge NSOutputStream *)(writeStream);
    
    // 设置通信过程中的代理
    inputStream.delegate = socketNet;
    outputStream.delegate = socketNet;
    
    // 将流对象添加到主运行循环(如果不加到主循环,Socket流是不会工作的)
    [inputStream scheduleInRunLoop:[NSRunLoop mainRunLoop] forMode:NSDefaultRunLoopMode];
    [outputStream scheduleInRunLoop:[NSRunLoop mainRunLoop] forMode:NSDefaultRunLoopMode];
    
    // 打开流
    [inputStream open];
    [outputStream open];
    
    NSStringEncoding encode = CFStringConvertEncodingToNSStringEncoding(kCFStringEncodingGB_18030_2000);
    NSData *data = [transitionStr dataUsingEncoding:encode];
    [outputStream write:data.bytes maxLength:data.length];
}

-(void)stream:(NSStream *)aStream handleEvent:(NSStreamEvent)eventCode
{
    switch (eventCode) {
        case NSStreamEventOpenCompleted:
            
            break;
        case NSStreamEventHasBytesAvailable:
        {
            uint8_t buffer[1024];
            NSMutableString *mstr = [[NSMutableString alloc]init];
            NSInteger len;// = [inputStream read:buffer maxLength:sizeof(buffer)];
            do{
                len =  [inputStream read:buffer maxLength:sizeof(buffer)];
                
                NSStringEncoding encode = CFStringConvertEncodingToNSStringEncoding(kCFStringEncodingGB_18030_2000);
                NSString *s = [[NSString alloc] initWithBytes:buffer length:len encoding:encode];
                if (s != nil)
                {
                    [mstr appendString:s];
                }
            }while (len == sizeof(buffer));
            NSInteger errorNum = 0 ;
            if (mstr == nil)
            {
                errorNum = 5;//返回数据为空
            }else
            {
                if (mstr.length >4)
                {
                    [mstr deleteCharactersInRange:NSMakeRange(0, 4)];
                    NSLog(@"-----mstr --- %@",mstr);
                    //NSMutableString * resultMuStr = [NSMutableString stringWithString:mstr];
                    if (![mstr isEqualToString:@""])
                    {
                        NSData *jsonData = [mstr dataUsingEncoding:NSUTF8StringEncoding];
                        NSError *err;
                        NSMutableDictionary *diction = [NSJSONSerialization JSONObjectWithData:jsonData options:NSJSONReadingMutableContainers error:&err];
                        NSString * resultSign = [diction objectForKey:@"sign"];
                        [diction removeObjectForKey:@"sign"];
                        NSString * sign = [SocketNet getSignStrWithDiction:diction signKey:signKey];
                        if (![sign isEqualToString:resultSign])
                        {
                            errorNum = 4;//签名错误
                        }
                    }
                    [_delegate getResultDataWithBackParameter:[SocketNet dataWithResultStr:mstr type:typeNum] errorCode:errorNum];
                }
            }
        }
            break;
        case NSStreamEventHasSpaceAvailable:
            break;
        case NSStreamEventErrorOccurred:
            
            [_delegate getResultDataWithBackParameter:nil errorCode:1];
            break;
        case NSStreamEventEndEncountered:
            // 做善后工作
            // 关闭流的同时，将流从主运行循环中删除
            [aStream close];
            [aStream removeFromRunLoop:[NSRunLoop mainRunLoop] forMode:NSDefaultRunLoopMode];
        default:
            break;
    }
    
}

//向服务器传的参数
+(NSString *)pinyinSort:(NSDictionary *)dict signKey:(NSString *)signkey
{
    NSString * sign = [SocketNet getSignStrWithDiction:dict signKey:signkey];
    [dict setValue:sign forKey:@"sign"];
    
    NSString * appendStr;
    if ([NSJSONSerialization isValidJSONObject:dict])
    {
        NSError *error;
        NSData *registerData = [NSJSONSerialization dataWithJSONObject:dict options:NSJSONWritingPrettyPrinted error:&error];
        NSString * str = [[NSString alloc] initWithData:registerData encoding:NSUTF8StringEncoding];
        appendStr = [NSString stringWithFormat:@"0%lu%@",(unsigned long)str.length,str];
        NSLog(@"Register JSON:%@",appendStr);
    }
    
    // NSLog(@"====%@====",encryptionStr);
    return  appendStr;
}

//传代理
+(void)setDelegate:(id<BackInfoDelegate>)delegate
{
    _delegate = delegate;
}

//获取密钥
+(void)byValueSignKey:(NSString *)signKeyStr
{
    signKey = signKeyStr;
}

//返回给用户的数据对象
+(BackParameter *)dataWithResultStr:(NSString *)str type:(NSInteger)typeNum
{
    NSData *jsonData = [str dataUsingEncoding:NSUTF8StringEncoding];
    NSError *err;
    NSMutableDictionary *diction = [NSJSONSerialization JSONObjectWithData:jsonData options:NSJSONReadingMutableContainers error:&err];
    BackParameter * backData = [[BackParameter alloc]init];
    backData.busicd = [diction objectForKey:@"busicd"];
    backData.respcd = [diction objectForKey:@"respcd"];
    backData.chcd = [diction objectForKey:@"chcd"];
    backData.errorDetail = [diction objectForKey:@"errorDetail"];
    backData.sign = [diction objectForKey:@"sign"];
    if (typeNum == 1)
    {
        NSMutableString * string =[NSMutableString stringWithString:[diction objectForKey:@"txamt"]];
        [string insertString:@"." atIndex:string.length-2];
        double txamt = [string doubleValue];
        backData.txamt = [NSString stringWithFormat:@"%.2f",txamt];
        backData.channelOrderNum = [diction objectForKey:@"channelOrderNum"];
        backData.consumerAccount = [diction objectForKey:@"consumerAccount"];
        backData.consumerId = [diction objectForKey:@"consumerId"];
        backData.orderNum = [diction objectForKey:@"orderNum"];
        backData.chcdDiscount = [diction objectForKey:@"chcdDiscount"];
        backData.merDiscount = [diction objectForKey:@"merDiscount"];
        backData.tag = 1;
        
    }else if (typeNum == 2)
    {
        NSMutableString * string = [NSMutableString stringWithString:[diction objectForKey:@"txamt"]];
        [string insertString:@"." atIndex:string.length-2];
        double txamt = [string doubleValue];
        backData.txamt = [NSString stringWithFormat:@"%.2f",txamt];
        backData.channelOrderNum = [diction objectForKey:@"channelOrderNum"];
        backData.orderNum = [diction objectForKey:@"orderNum"];
        backData.qrcode = [diction objectForKey:@"qrcode"];
        backData.tag = 2;
    }else if (typeNum == 3)
    {
        backData.channelOrderNum = [diction objectForKey:@"channelOrderNum"];
        backData.consumerAccount = [diction objectForKey:@"consumerAccount"];
        backData.consumerId = [diction objectForKey:@"consumerId"];
        backData.origOrderNum = [diction objectForKey:@"origOrderNum"];
        backData.chcdDiscount = [diction objectForKey:@"chcdDiscount"];
        backData.merDiscount = [diction objectForKey:@"merDiscount"];
        backData.tag = 3;
    }else if(typeNum == 4)
    {
        backData.consumerAccount = [diction objectForKey:@"consumerAccount"];
        backData.consumerId = [diction objectForKey:@"consumerId"];
        backData.orderNum = [diction objectForKey:@"orderNum"];
        backData.origOrderNum = [diction objectForKey:@"origOrderNum"];
        backData.chcdDiscount = [diction objectForKey:@"chcdDiscount"];
        backData.merDiscount = [diction objectForKey:@"merDiscount"];
        backData.tag = 4;
    }else if (typeNum == 5)
    {
        backData.channelOrderNum = [diction objectForKey:@"channelOrderNum"];
        backData.consumerAccount = [diction objectForKey:@"consumerAccount"];
        backData.consumerId = [diction objectForKey:@"consumerId"];
        backData.orderNum = [diction objectForKey:@"orderNum"];
        backData.origOrderNum = [diction objectForKey:@"origOrderNum"];
        backData.chcdDiscount = [diction objectForKey:@"chcdDiscount"];
        backData.merDiscount = [diction objectForKey:@"merDiscount"];
        backData.tag = 5;
    }else if (typeNum == 6)
    {
        backData.orderNum = [diction objectForKey:@"orderNum"];
        backData.scanCodeId = [diction objectForKey:@"scanCodeId"];
        backData.cardId = [diction objectForKey:@"cardId"];
        backData.cardInfo = [diction objectForKey:@"cardInfo"];
        backData.tag = 6;
    }
    return backData;

}

@end
