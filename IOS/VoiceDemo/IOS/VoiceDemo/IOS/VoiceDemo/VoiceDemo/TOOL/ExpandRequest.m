//
//  ExpandRequest.m
//  VoiceDemo
//
//  Created by 黄达能 on 15/9/7.
//  Copyright (c) 2015年 __VTPayment__. All rights reserved.
//

#import "ExpandRequest.h"
#import "CommonCrypto/CommonDigest.h"
#import "Request.h"
#import "RegisterTable.h"

@interface ExpandRequest()

@property (strong, nonatomic) NSMutableData      *resultData;

@property (strong, nonatomic) NSMutableString           *Voice_Path;

@end

@implementation ExpandRequest
-(void)connectionNet:(NSString *)VoicePath andUserKey:(NSString *)userkey
{
    _Voice_Path=[NSMutableString stringWithString:VoicePath];
    
    NSString *appkey = @"13784381190000d5";
    NSString *userId = @"xunlian";
    NSString *secretkey = @"e6ec2392200db4315bf5c9745546bd92";
    NSMutableDictionary *dict= [RegisterTableDAO getObjectByName:[RegisterTableDAO getNameWhoIsUsing]];
    NSString *userid = [NSString stringWithFormat:@"%@_%@",[RegisterTableDAO getNameWhoIsUsing],[dict objectForKey:@"time"]];
    
    
    //签名后的sig数据
    NSMutableString *sign=[[NSMutableString alloc]initWithString:appkey];
    
    NSDate *date=[NSDate date];
    NSString *timeSp = [NSString stringWithFormat:@"%ld", (long)[date timeIntervalSince1970]];
    [sign appendString:timeSp];
    [sign appendString:secretkey];
    
    [sign setString:[self sha1:sign]];
    
    NSString * string = [NSString stringWithFormat: @"{\"cmd\":\"start\",\"param\":{\"app\":{\"applicationId\":\"%@\",\"userId\":\"%@\",\"timestamp\":\"%@\",\"sig\":\"%@\"},\"audio\":{\"audioType\":\"wav\",\"channel\": 1,\"sampleBytes\":2,\"sampleRate\": 16000},\"request\": {\"coreType\":\"sv\",\"userid\":\"%@\",\"userkey\":\"%@\",\"svMode\":3,\"threshold\":0.5}}}",appkey,userId,timeSp,sign,userid,userkey];
    NSData * data = [string dataUsingEncoding:NSUTF8StringEncoding];
    
    //NSLog(@"str = \n %@  \n 结束 ",string);
    
    NSURL *url=[NSURL URLWithString:[NSString stringWithFormat:@"%@/api/v3.0/score",CONNECTION_URL]];
    //创建请求
    NSMutableURLRequest *request=[NSMutableURLRequest requestWithURL:url cachePolicy:NSURLRequestReloadIgnoringLocalCacheData timeoutInterval:60];
    [request setHTTPBody:data];
    [request setHTTPMethod:@"POST"];
    [request setValue:@"text/plain" forHTTPHeaderField:@"Content-Type"];
    [request setValue:@"Keep-Alive" forHTTPHeaderField:@"Connection"];
    VTConnectionRequest * connectiona = [[VTConnectionRequest alloc]initWithRequest:request delegate:self];
    connectiona.tag = 1;
    [connectiona start];
}
- (void)connection:(VTConnectionRequest *)connection didReceiveResponse:(NSURLResponse *)response
{
    NSLog(@"%ld---------%@",connection.tag,response);
    NSHTTPURLResponse* httpResponse = (NSHTTPURLResponse*)response;
    NSInteger responseStatusCode = [httpResponse statusCode];
    if (connection.tag == 1)
    {
        NSLog(@"connection.tag == 1 状态码 =  %ld   \n ",responseStatusCode);
        if (responseStatusCode == 200)
        {
            NSURL *url=[NSURL URLWithString:[NSString stringWithFormat:@"%@/api/v3.0/score",CONNECTION_URL]];
            //创建请求
            NSString *path=[NSString stringWithString:_Voice_Path];
            NSData * data = [NSData dataWithContentsOfFile:path];
            NSMutableURLRequest *request=[NSMutableURLRequest requestWithURL:url cachePolicy:NSURLRequestReloadIgnoringLocalCacheData timeoutInterval:60];
            
            [request setHTTPBody:data];
            [request setHTTPMethod:@"POST"];
            [request setValue:@"application/octet-stream" forHTTPHeaderField:@"Content-Type"];
            [request setValue:@"Keep-Alive" forHTTPHeaderField:@"Connection"];
            VTConnectionRequest * connectionb = [[VTConnectionRequest alloc]initWithRequest:request delegate:self];
            connectionb.tag = 2;
            [connectionb start];
        }
    }
    else if (connection.tag == 2)
    {
        NSLog(@"connection.tag == 2 状态码 =  %ld   \n ",responseStatusCode);
        if (responseStatusCode == 200)
        {
            NSURL *url=[NSURL URLWithString:[NSString stringWithFormat:@"%@/api/v3.0/score",CONNECTION_URL]];
            //创建请求
            NSMutableURLRequest *request=[NSMutableURLRequest requestWithURL:url cachePolicy:NSURLRequestReloadIgnoringLocalCacheData timeoutInterval:60];
            NSString * string = [NSString stringWithFormat: @"{\"cmd\":\"stop\"}"];
            NSData * data = [string dataUsingEncoding:NSUTF8StringEncoding];
            [request setHTTPBody:data];
            [request setHTTPMethod:@"POST"];
            [request setValue:@"text/plain" forHTTPHeaderField:@"Content-Type"];
            [request setValue:@"Keep-Alive" forHTTPHeaderField:@"Connection"];
            VTConnectionRequest * connectionc = [[VTConnectionRequest alloc]initWithRequest:request delegate:self];
            connectionc.tag = 3;
            [connectionc start];
        }
    }else if (connection.tag == 3)
    {
        NSLog(@"connection.tag == 3 状态码 =  %ld   \n ",responseStatusCode);
        
        if (!self.resultData)
        {
            self.resultData = [[NSMutableData alloc]init];
        }else
        {
            [self.resultData setLength:0];
        }
    }
}
- (void)connection:(VTConnectionRequest *)connection didReceiveData:(NSData *)data
{
    if (connection.tag == 3)
    {
        [self.resultData appendData:data];
    }
}
- (void)connectionDidFinishLoading:(VTConnectionRequest *)connection
{
    if (connection.tag == 3)
    {
        NSDictionary * dict = [NSJSONSerialization JSONObjectWithData:self.resultData options:NSJSONReadingMutableLeaves error:nil];
        //        NSString * str = [[NSString alloc]initWithData:self.resultData encoding:NSUTF8StringEncoding];
        NSLog(@"-----%@",dict);
        NSLog(@"~~~~~~~~~~~~~~~~~~~~~~~~~~~~第四次结束~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~");
        //        NSLog(@"self.resultData   --- \n  dict --=  \n %@",[dict objectForKey:@"result"]);
        //        NSDictionary * diction = [dict objectForKey:@"result"];
        //        NSString * str = [diction objectForKey:@"svValue"];
        //        NSLog(@" \n  str --=  \n %@",str);
    }
}
#pragma mark sha1 加密
- (NSString *) sha1: (NSString *) inPutText
{
    //可以对中文进行加密
    const char * cstr = [inPutText UTF8String];
    //使用对应的CC_SHA1,CC_SHA256,CC_SHA384,CC_SHA512的长度分别是20,32,48,64
    unsigned char digest[CC_SHA1_DIGEST_LENGTH];
    //使用对应的CC_SHA256,CC_SHA384,CC_SHA512
    CC_SHA1(cstr,  (CC_LONG)strlen(cstr), digest);
    NSMutableString *output = [NSMutableString stringWithCapacity:CC_SHA1_DIGEST_LENGTH * 2];
    for(int i = 0; i < CC_SHA1_DIGEST_LENGTH; i++) {
        [output appendFormat:@"%02x", digest[i]];
    }
    return output;
}
@end
