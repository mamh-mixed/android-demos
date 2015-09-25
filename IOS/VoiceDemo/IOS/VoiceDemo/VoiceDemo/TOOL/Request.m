//
//  Request.m
//  VoiceDemo
//
//  Created by 黄达能 on 15/9/6.
//  Copyright (c) 2015年 __VTPayment__. All rights reserved.
//

#import "Request.h"
#import "CommonCrypto/CommonDigest.h"
#import "GenerateRequest.h"
#import "RegisterTable.h"

@implementation VTConnectionRequest


@end

@interface Request()

@property (strong, nonatomic) NSString           *UserKey;//保存connectionNet 传入的UserKey

@property (strong, nonatomic) NSMutableData      *resultData;

@property (strong, nonatomic) NSArray            *Voice_Array;//保存connectionNet 传入的数组

@end

@implementation Request

static Request *request=nil;

+(Request *)sharedRequest
{
    static dispatch_once_t onceTocken;
    dispatch_once(&onceTocken, ^{
        if (request==nil) {
            request=[[Request alloc]init];
            request.successTimes=0;
        }
    });
    return request;
}
-(void)connectionNet:(NSArray *)VoicePath andUserKey:(NSString *)UserKey
{
    _Voice_Array=[NSArray arrayWithArray:VoicePath];
    
    NSString *appkey = @"13784381190000d5";
    NSString *userId = @"xunlian";
    NSMutableDictionary *dict= [RegisterTableDAO getObjectByName:[RegisterTableDAO getNameWhoIsUsing]];
    NSString *userid = [NSString stringWithFormat:@"%@_%@",[RegisterTableDAO getNameWhoIsUsing],[dict objectForKey:@"time"]];
    NSString *secretkey = @"e6ec2392200db4315bf5c9745546bd92";
    NSString *userkey=[NSString stringWithString:UserKey];
    _UserKey=userkey;
    
    //签名后的sig数据
    NSMutableString *sign=[[NSMutableString alloc]initWithString:appkey];
    
    NSDate *date=[NSDate date];
    NSString *timeSp = [NSString stringWithFormat:@"%ld", (long)[date timeIntervalSince1970]];
    [sign appendString:timeSp];
    
    [sign appendString:secretkey];
    
    [sign setString:[self sha1:sign]];
    
    NSString * string = [NSString stringWithFormat: @"{\"cmd\":\"start\",\"param\":{\"app\":{\"applicationId\":\"%@\",\"userId\":\"%@\",\"timestamp\":\"%@\",\"sig\":\"%@\"},\"audio\":{\"audioType\":\"wav\",\"channel\": 1,\"sampleBytes\":2,\"sampleRate\": 16000},\"request\": {\"coreType\":\"sv\",\"userid\":\"%@\",\"userkey\":\"%@\",\"svMode\":0}}}",appkey,userId,timeSp,sign,userid,userkey];
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
    connectiona.tag = 1;
    [connectiona start];
    
    NSLog(@"~~~~~~~~~~~~~~~~~~~~~~~~~~~~第一次开始~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~");
}
- (void)connection:(VTConnectionRequest *)connection didReceiveResponse:(NSURLResponse *)response
{
    NSLog(@"%ld---------%@",connection.tag,response);
    NSHTTPURLResponse* httpResponse = (NSHTTPURLResponse*)response;
    NSInteger responseStatusCode = [httpResponse statusCode];
    if (connection.tag == 1)
    {
        if (responseStatusCode == 200)
        {
            NSURL *url=[NSURL URLWithString:[NSString stringWithFormat:@"%@/api/v3.0/score",CONNECTION_URL]];
            //创建请求
            NSString *path=[NSString stringWithString:_Voice_Array[_successTimes]];
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
    }else if (connection.tag == 2)
    {
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
        NSDictionary *dic=[dict objectForKey:@"result"];
        if ([[dic objectForKey:@"svValue"] isEqualToNumber:[NSNumber numberWithInt:0]]) {//语音发送成功
            _successTimes++;
            if (_successTimes==1) {//发送第二次语音
                [self connectionNet:_Voice_Array andUserKey:_UserKey];
            }
            else if(_successTimes==2){//发送第三次语音
                [self connectionNet:_Voice_Array andUserKey:_UserKey];
            }
            else if (_successTimes==3) {
                //所有的语音都发送成功
                dispatch_async(dispatch_get_main_queue(), ^{
#pragma mark- 发送一个请求成功的通知 （在VTRecord_Alipay等界面接收）
                    [[NSNotificationCenter defaultCenter] postNotificationName:@"RequestIsSuccess" object:nil];
                });
                GenerateRequest *gRequest=[[GenerateRequest alloc]init];
                [gRequest connectionNet:_UserKey];
            }
        }
        else{
            //语音发送失败
            dispatch_async(dispatch_get_main_queue(), ^{
#pragma mark- 发送一个请求失败的通知  带了一个成功的次数 利用成功的次数判断是哪次语音发送失败（在VTRecord_Alipay等界面接受收）
                [[NSNotificationCenter defaultCenter] postNotificationName:@"RequestIsDefault" object:[NSNumber numberWithInteger:_successTimes]];
            });
        }
    }
}
// 请求数据失败时触发
- (void)connection:(NSURLConnection *)connection didFailWithError:(NSError *)error
{
    NSLog(@"%s", __FUNCTION__);
}
#pragma mark -sha1 加密
- (NSString*) sha1: (NSString *) inPutText
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
