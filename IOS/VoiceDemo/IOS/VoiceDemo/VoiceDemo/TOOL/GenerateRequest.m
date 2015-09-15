//
//  GenerateRequest.m
//  VoiceDemo
//
//  Created by 黄达能 on 15/9/7.
//  Copyright (c) 2015年 __VTPayment__. All rights reserved.
//

#import "GenerateRequest.h"
#import "CommonCrypto/CommonDigest.h"
#import "DetectRequest.h"
#import "RegisterTable.h"
#import "Request.h"

@interface GenerateRequest()


@end

@implementation GenerateRequest
- (void)connectionNet:(NSString *)UserKey
{
    NSString *appkey = @"13784381190000d5";
    NSString *userId = @"xunlian";
    NSMutableDictionary *dict= [RegisterTableDAO getObjectByName:[RegisterTableDAO getNameWhoIsUsing]];
    NSString *userid = [NSString stringWithFormat:@"%@_%@",[RegisterTableDAO getNameWhoIsUsing],[dict objectForKey:@"time"]];
    NSString *secretkey = @"e6ec2392200db4315bf5c9745546bd92";
    NSString *userkey=[NSString stringWithString:UserKey];
    
    //签名后的sig数据
    NSMutableString *sign=[[NSMutableString alloc]initWithString:appkey];
    
    NSDate *date=[NSDate date];
    NSString *timeSp = [NSString stringWithFormat:@"%ld", (long)[date timeIntervalSince1970]];
    [sign appendString:timeSp];
    [sign appendString:secretkey];
    
    [sign setString:[self sha1:sign]];
    
    NSString * string = [NSString stringWithFormat: @"{\"cmd\":\"start\",\"param\":{\"app\":{\"applicationId\":\"%@\",\"userId\":\"%@\",\"timestamp\":\"%@\",\"sig\":\"%@\"},\"audio\":{\"audioType\":\"wav\",\"channel\": 1,\"sampleBytes\":2,\"sampleRate\": 16000},\"request\": {\"coreType\":\"sv\",\"userid\":\"%@\",\"userkey\":\"%@\",\"svMode\":1}}}",appkey,userId,timeSp,sign,userid,userkey];
    NSData * data = [string dataUsingEncoding:NSUTF8StringEncoding];
    
    NSLog(@"str = \n %@  \n 结束 ",string);
    
    
    NSURL *url=[NSURL URLWithString:[NSString stringWithFormat:@"%@/api/v3.0/score",CONNECTION_URL]];
    //创建请求
    NSMutableURLRequest *request=[NSMutableURLRequest requestWithURL:url cachePolicy:NSURLRequestReloadIgnoringLocalCacheData timeoutInterval:60];
    [request setHTTPBody:data];
    [request setHTTPMethod:@"POST"];
    [request setValue:@"text/plain" forHTTPHeaderField:@"Content-Type"];
    [request setValue:@"Keep-Alive" forHTTPHeaderField:@"Connection"];
    VTConnectionRequest * connectiona = [[VTConnectionRequest alloc]initWithRequest:request delegate:self];
    connectiona.tag=1;
    [connectiona start];
    
}
- (void)connection:(VTConnectionRequest *)connection didReceiveResponse:(NSURLResponse *)response
{
    NSLog(@"%@",response);
}
- (void)connection:(VTConnectionRequest *)connection didReceiveData:(NSData *)data
{
    if (1)
    {
        NSDictionary * dict = [NSJSONSerialization JSONObjectWithData:data options:NSJSONReadingMutableLeaves error:nil];
        //        NSString * str = [[NSString alloc]initWithData:self.resultData encoding:NSUTF8StringEncoding];
        NSLog(@"-----%@",dict);
       // NSDictionary *dic=[dict objectForKey:@"result"];
        //第二次请求失败
//        if (![[dic objectForKey:@"svValue"] isEqualToNumber:[NSNumber numberWithInt:0]]) {
//            [[NSNotificationCenter defaultCenter]postNotificationName:@"RequestIsDefault" object:nil];
//        }
        NSLog(@"~~~~~~~~~~~~~~~~~~~~~~~~~~~~第二次结束~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~");
    }
}
- (void)connectionDidFinishLoading:(VTConnectionRequest *)connection
{
#if 0
    static dispatch_once_t onceToken;
    dispatch_once(&onceToken, ^{
        NSURL *url=[NSURL URLWithString:[NSString stringWithFormat:@"%@/api/v3.0/score",CONNECTION_URL]];
        //创建请求
        NSMutableURLRequest *request=[NSMutableURLRequest requestWithURL:url cachePolicy:NSURLRequestReloadIgnoringLocalCacheData timeoutInterval:60];
        NSString * string = [NSString stringWithFormat: @"{\"cmd\":\"stop\"}"];
        NSData * data = [string dataUsingEncoding:NSUTF8StringEncoding];
        [request setHTTPBody:data];
        [request setHTTPMethod:@"POST"];
        [request setValue:@"text/plain" forHTTPHeaderField:@"Content-Type"];
        [request setValue:@"Keep-Alive" forHTTPHeaderField:@"Connection"];
        NSURLConnection * connectionc = [[NSURLConnection alloc]initWithRequest:request delegate:self];
        [connectionc start];
    });
#endif
    if (connection.tag==1) {
        NSURL *url=[NSURL URLWithString:[NSString stringWithFormat:@"%@/api/v3.0/score",CONNECTION_URL]];
        //创建请求
        NSMutableURLRequest *request=[NSMutableURLRequest requestWithURL:url cachePolicy:NSURLRequestReloadIgnoringLocalCacheData timeoutInterval:60];
        NSString * string = [NSString stringWithFormat: @"{\"cmd\":\"stop\"}"];
        NSData * data = [string dataUsingEncoding:NSUTF8StringEncoding];
        [request setHTTPBody:data];
        [request setHTTPMethod:@"POST"];
        [request setValue:@"text/plain" forHTTPHeaderField:@"Content-Type"];
        [request setValue:@"Keep-Alive" forHTTPHeaderField:@"Connection"];
        VTConnectionRequest * connectionb = [[VTConnectionRequest alloc]initWithRequest:request delegate:self];
        connectionb.tag = 2;
        [connectionb start];
    }
}

#pragma mark sha1 加密
-(NSString*) sha1: (NSString *) inPutText
{
    //可以对中文进行加密
    const char * cstr = [inPutText UTF8String];
    //使用对应的CC_SHA1,CC_SHA256,CC_SHA384,CC_SHA512的长度分别是20,32,48,64
    unsigned char digest[CC_SHA1_DIGEST_LENGTH];
    //使用对应的CC_SHA256,CC_SHA384,CC_SHA512
    CC_SHA1(cstr,  (CC_LONG)strlen(cstr), digest);
    NSMutableString* output = [NSMutableString stringWithCapacity:CC_SHA1_DIGEST_LENGTH * 2];
    for(int i = 0; i < CC_SHA1_DIGEST_LENGTH; i++) {
        [output appendFormat:@"%02x", digest[i]];
    }
    return output;
}

@end
