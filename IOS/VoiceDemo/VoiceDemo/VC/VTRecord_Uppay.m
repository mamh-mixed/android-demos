//
//  VTRecord_Uppay.m
//  VoiceDemo
//
//  Created by 黄达能 on 15/9/9.
//  Copyright (c) 2015年 __VTPayment__. All rights reserved.
//

#import "VTRecord_Uppay.h"
#import "Request.h"
#import "VTRecord_Wxpay.h"
#import "RegisterTable.h"

@interface VTRecord_Uppay ()
{
    NSString *path1;//音频存放的路径
    NSString *path2;
    NSString *path3;
    
    UIImageView *progressImage;
}
@end

@implementation VTRecord_Uppay

@synthesize recorder;
@synthesize player;



- (void)viewDidLoad {
    [super viewDidLoad];
    
    
    UIImageView *imageView=[[UIImageView alloc]initWithFrame:[UIScreen mainScreen].bounds];
    imageView.image=[UIImage imageNamed:@"paybg"];
    [self.view addSubview:imageView];
    UIButton *back=[UIButton buttonWithType:UIButtonTypeCustom];
    back.frame=CGRectMake(10, 30, 12, 21);
    [back setImage:[UIImage imageNamed:@"back"] forState:UIControlStateNormal];
    [back addTarget:self action:@selector(back) forControlEvents:UIControlEventTouchUpInside];
    [self.view addSubview:back];
    
    UILabel *label=[[UILabel alloc]initWithFrame:CGRectMake(0, 40, SCREENWIDTH, 30)];
    label.textAlignment=NSTextAlignmentCenter;
    label.text=@"按住话筒录音，读出下述文字:";
    [self.view addSubview:label];
    
    UILabel *lbl=[[UILabel alloc]initWithFrame:CGRectMake(0, 80, SCREENWIDTH, 30)];
    lbl.textAlignment=NSTextAlignmentCenter;
    lbl.text=@"银联";
    lbl.textColor=[UIColor redColor];
    [self.view addSubview:lbl];
    
    CGFloat height=50.0f;
    CGFloat space=10.0f;
    for (int i=1; i<4; i++) {//录音按钮的tag 值为 1 2 3
        UIButton *record=[UIButton buttonWithType:UIButtonTypeCustom];
        record.tag = i;
        record.frame=CGRectMake(0, SCREENHEIGHT-4*height-3*space+(i-1)*(space+height), SCREENWIDTH, height);
        record.backgroundColor=[UIColor whiteColor];
        
        UIImageView *image=[[UIImageView alloc]initWithFrame:CGRectMake((SCREENWIDTH-20)/2,(height-32)/2, 20, 32)];
        image.image=[UIImage imageNamed:@"microphone"];
        [record addSubview:image];
        [record addTarget:self action:@selector(record:) forControlEvents:UIControlEventTouchDown];
        [record addTarget:self action:@selector(touchUpInside:) forControlEvents:UIControlEventTouchUpInside];
        [self.view addSubview:record];
    }
    
    UIButton *btn=[UIButton buttonWithType:UIButtonTypeCustom];
    btn.frame=CGRectMake(0, SCREENHEIGHT-height, SCREENWIDTH, height);
    [btn setTitle:@"提交" forState:UIControlStateNormal];
    [btn addTarget:self action:@selector(sumbit) forControlEvents:UIControlEventTouchUpInside];
    btn.backgroundColor=[UIColor blueColor];
    [self.view addSubview:btn];
}
-(void)back
{
    [self dismissViewControllerAnimated:YES completion:nil];
}
#pragma mark -录音
-(void)record:(UIButton *)sender
{
    sender.backgroundColor=[UIColor orangeColor];
    AVAudioSession *session=[AVAudioSession sharedInstance];
    [session setCategory:AVAudioSessionCategoryPlayAndRecord error:nil];//设置类别,表示该应用同时支持播放和录音
    
    NSMutableDictionary *recordSettings=[[NSMutableDictionary alloc]initWithCapacity:10];
    [recordSettings setObject:[NSNumber numberWithInt: kAudioFormatLinearPCM] forKey: AVFormatIDKey];
    [recordSettings setObject:[NSNumber numberWithFloat:8000.0] forKey: AVSampleRateKey];
    [recordSettings setObject:[NSNumber numberWithInt:1] forKey:AVNumberOfChannelsKey];
    [recordSettings setObject:[NSNumber numberWithInt:16] forKey:AVLinearPCMBitDepthKey];
    [recordSettings setObject:[NSNumber numberWithInt:AVAudioQualityHigh] forKey:AVEncoderAudioQualityKey];
    
    NSURL *url = [NSURL fileURLWithPath:[NSTemporaryDirectory() stringByAppendingPathComponent: [NSString stringWithFormat: @"%ld.%@",(long)sender.tag,@"wav"]]];//默认acf格式 转成wav格式 方便后面的api解析
    //NSLog(@" url  -----%@",url);
    if (sender.tag==1) {
        path1=[NSTemporaryDirectory() stringByAppendingPathComponent: [NSString stringWithFormat: @"%ld.%@",(long)sender.tag,@"wav"]];
        //NSLog(@"path1 --------%@",path1);
    }
    else if (sender.tag==2){
        path2=[NSTemporaryDirectory() stringByAppendingPathComponent: [NSString stringWithFormat: @"%ld.%@",(long)sender.tag,@"wav"]];
    }
    else if (sender.tag==3){
        path3=[NSTemporaryDirectory() stringByAppendingPathComponent: [NSString stringWithFormat: @"%ld.%@",(long)sender.tag,@"wav"]];
    }
    recorder=[[AVAudioRecorder alloc]initWithURL:url settings:recordSettings error:nil];
    [recorder record];//开始录音
}
-(void)touchUpInside:(UIButton *)sender
{
    sender.backgroundColor=[UIColor whiteColor];
    sender.alpha=0;
    [recorder stop];
    //判断录音的时长
    dispatch_async(dispatch_get_main_queue(), ^{
        if ([self testAudioDuration:sender.tag]) {
            //4 5 6
            UIButton *btn=[UIButton buttonWithType:UIButtonTypeCustom];
            btn.frame=CGRectMake(sender.frame.origin.x, sender.frame.origin.y, sender.frame.size.width-50, sender.frame.size.height);
            btn.backgroundColor=[UIColor blueColor];
            btn.tag=sender.tag+3;
            [btn addTarget:self action:@selector(play:) forControlEvents:UIControlEventTouchUpInside];
            [self.view addSubview:btn];
            UIImageView *image=[[UIImageView alloc]initWithImage:[UIImage imageNamed:@"play"]];
            image.frame=CGRectMake((btn.frame.size.width-30)/2,(btn.frame.size.height-30)/2, 30, 30);
            [btn addSubview:image];
            //7 8 9
            UIButton *cancel=[UIButton buttonWithType:UIButtonTypeCustom];
            cancel.frame=CGRectMake(btn.frame.size.width+1, btn.frame.origin.y, SCREENWIDTH-btn.frame.size.width-1, btn.frame.size.height);
            [cancel addTarget:self action:@selector(cancelClick:) forControlEvents:UIControlEventTouchUpInside];
            cancel.tag=btn.tag+3;
            UIImageView *ige=[[UIImageView alloc]initWithFrame:CGRectMake((cancel.frame.size.width-20)/2, (cancel.frame.size.height-20)/2, 20, 20)];
            ige.image=[UIImage imageNamed:@"block"];
            cancel.backgroundColor=[UIColor blueColor];
            [cancel addSubview:ige];
            [self.view addSubview:cancel];
        }
        else{
            sender.alpha=1;
        }
    });
}
//测试音频的时长
-(BOOL)testAudioDuration:(NSUInteger )tag
{
    AVURLAsset *audioAsset;
    switch (tag) {
        case 1:audioAsset=[AVURLAsset URLAssetWithURL:[NSURL URLWithString:[NSString stringWithFormat:@"file://%@",path1]] options:nil];
            break;
        case 2:audioAsset=[AVURLAsset URLAssetWithURL:[NSURL URLWithString:[NSString stringWithFormat:@"file://%@",path2]] options:nil];
            break;
        case 3:audioAsset=[AVURLAsset URLAssetWithURL:[NSURL URLWithString:[NSString stringWithFormat:@"file://%@",path3]] options:nil];
            break;
        default:break;
    }
    CMTime audioDuration =audioAsset.duration;
    float audioDurationSeconds=CMTimeGetSeconds(audioDuration);
    if (audioDurationSeconds<0.6) {
        dispatch_async(dispatch_get_main_queue(), ^{
            UIAlertView *alert=[[UIAlertView alloc]initWithTitle:nil message:@"录音时间过短" delegate:self cancelButtonTitle:@"我知道了" otherButtonTitles: nil];
            [alert show];
        });
        return NO;
    }
    return YES;
}
-(void)cancelClick:(UIButton *)sender
{
    [sender removeFromSuperview];
    UIButton *btn=(UIButton *)[self.view viewWithTag:sender.tag-3];
    [btn removeFromSuperview];
    UIButton *button=(UIButton *)[self.view viewWithTag:sender.tag-6];
    button.alpha=1;
}
-(void)play:(UIButton *)sender
{
    NSString *path=[NSTemporaryDirectory() stringByAppendingPathComponent:[NSString stringWithFormat:@"%ld.wav",sender.tag-3]];
    player=[[AVAudioPlayer alloc]initWithContentsOfURL:[NSURL URLWithString:path] error:nil];
    player.volume=20.0f;
    [player prepareToPlay];
    [player play];
}
#pragma mark -提交
-(void)sumbit
{
#if 1
    //判断三次的录音是否已经完成
    UIButton *btn1=(UIButton *)[self.view viewWithTag:4];
    UIButton *btn2=(UIButton *)[self.view viewWithTag:5];
    UIButton *btn3=(UIButton *)[self.view viewWithTag:6];
    if (!(btn1&&btn2&&btn3)) {
        UIAlertView *alert=[[UIAlertView alloc]initWithTitle:nil message:@"请完成三次录音后在提交" delegate:self cancelButtonTitle:@"取消" otherButtonTitles:nil];
        [alert show];
        return;
    }
#endif
#if 1
    //训练模式
    //训练模式分为两种模式，模式一，清空服务器之前训练过的缓存，调用start之后直接调用end，不发送音频。模式二，正常训练
    NSDictionary *dict =[RegisterTableDAO getObjectByName:[RegisterTableDAO getNameWhoIsUsing]];
    NSMutableString *userkey=[[NSMutableString alloc]initWithString:[RegisterTableDAO getNameWhoIsUsing]];
    [userkey appendString:@"_uppay_"];
    [userkey appendString:[dict objectForKey:@"time"]];
    NSArray *path_array=[NSArray arrayWithObjects:path1,path2,path3,nil];
    dispatch_async(dispatch_get_main_queue(), ^{
        [Request sharedRequest].successTimes=0;
        [[Request sharedRequest] connectionNet:path_array andUserKey:userkey];
    });
#endif
    [[NSNotificationCenter defaultCenter] addObserver:self selector:@selector(requestIsSuccess) name:@"RequestIsSuccess" object:nil];//监听是否成功
    [[NSNotificationCenter defaultCenter] addObserver:self selector:@selector(requestIsDefault) name:@"RequestIsDefault" object:nil];//监听是否失败
    
    progressImage=[[UIImageView alloc]initWithFrame:CGRectMake(0, 0, SCREENWIDTH, SCREENHEIGHT)];
    [self.view addSubview:progressImage];
    progressImage.backgroundColor=[UIColor blackColor];
    progressImage.alpha=0.8;
    progressImage.userInteractionEnabled=YES;
    UIActivityIndicatorView *activity=[[UIActivityIndicatorView alloc]initWithActivityIndicatorStyle:UIActivityIndicatorViewStyleWhite];
    activity.center=progressImage.center;
    [progressImage addSubview:activity];
    [activity startAnimating];
}
-(void)requestIsDefault
{
    [progressImage removeFromSuperview];
    UIAlertView *alert=[[UIAlertView alloc]initWithTitle:nil message:@"语音发送失败 请将所有语音删除后 重试" delegate:self cancelButtonTitle:@"我知道了" otherButtonTitles: nil];
    [alert show];
    [[NSNotificationCenter defaultCenter]removeObserver:self name:@"RequestIsDefault" object:nil];
}
-(void)requestIsSuccess
{
    [progressImage removeFromSuperview];
    VTRecord_Wxpay *wxpay=[[VTRecord_Wxpay alloc]init];
    [self presentViewController:wxpay animated:YES completion:nil];
}
- (void)didReceiveMemoryWarning {
    [super didReceiveMemoryWarning];
}

@end
