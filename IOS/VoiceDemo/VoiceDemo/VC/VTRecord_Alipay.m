//
//  VTRecord_Alipay.m
//  VoiceDemo
//
//  Created by 黄达能 on 15/8/28.
//  Copyright (c) 2015年 __VTPayment__. All rights reserved.
//

#import "VTRecord_Alipay.h"

#define SCREENWIDTH [UIScreen mainScreen].bounds.size.width
#define SCREENHEIGHT [UIScreen mainScreen].bounds.size.height

@interface VTRecord_Alipay ()

@end

@implementation VTRecord_Alipay

@synthesize recorder;
@synthesize player;

- (void)viewDidLoad {
    [super viewDidLoad];
    UIImageView *imageView=[[UIImageView alloc]initWithFrame:[UIScreen mainScreen].bounds];
    imageView.image=[UIImage imageNamed:@"paybg"];
    [self.view addSubview:imageView];
    UILabel *label=[[UILabel alloc]initWithFrame:CGRectMake(0, 40, SCREENWIDTH, 30)];
    label.textAlignment=NSTextAlignmentCenter;
    label.text=@"按住话筒录音，读出下述文字:";
    [self.view addSubview:label];
    
    UILabel *lbl=[[UILabel alloc]initWithFrame:CGRectMake(0, 80, SCREENWIDTH, 30)];
    lbl.textAlignment=NSTextAlignmentCenter;
    lbl.text=@"支付宝";
    lbl.textColor=[UIColor redColor];
    [self.view addSubview:lbl];
    
    CGFloat height=50.0f;
    for (int i=0; i<3; i++) {//录音按钮的tag 值为 0 1 2
        UIButton *record=[UIButton buttonWithType:UIButtonTypeCustom];
        record.tag=i;
        record.frame=CGRectMake(0, SCREENHEIGHT-4*height+i*(height + 10 )-70, SCREENWIDTH, height-1);
        record.backgroundColor=[UIColor whiteColor];
        UIImageView *image=[[UIImageView alloc]initWithFrame:CGRectMake((SCREENWIDTH-20)/2,(height-32)/2, 20, 32)];
        image.image=[UIImage imageNamed:@"microphone"];
        [record addSubview:image];
        [record addTarget:self action:@selector(record:) forControlEvents:UIControlEventTouchDown];
        [record addTarget:self action:@selector(touchUpInside:) forControlEvents:UIControlEventTouchUpInside];
        [self.view addSubview:record];
    }
    
    UIButton *btn=[UIButton buttonWithType:UIButtonTypeCustom];
    btn.frame=CGRectMake(0, SCREENHEIGHT-height, SCREENWIDTH, height-1);
    [btn setTitle:@"提交" forState:UIControlStateNormal];
    [btn addTarget:self action:@selector(sumbit) forControlEvents:UIControlEventTouchUpInside];
}
#pragma mark -录音
-(void)record:(UIButton *)sender
{
    sender.backgroundColor=[UIColor orangeColor];
    AVAudioSession *session=[AVAudioSession sharedInstance];
    [session setCategory:AVAudioSessionCategoryRecord error:nil];
    
    NSMutableDictionary *recordSettings=[[NSMutableDictionary alloc]initWithCapacity:10];
    [recordSettings setObject:[NSNumber numberWithInt: kAudioFormatLinearPCM] forKey: AVFormatIDKey];
    [recordSettings setObject:[NSNumber numberWithFloat:8000.0] forKey: AVSampleRateKey];
    [recordSettings setObject:[NSNumber numberWithInt:1] forKey:AVNumberOfChannelsKey];
    [recordSettings setObject:[NSNumber numberWithInt:16] forKey:AVLinearPCMBitDepthKey];
    [recordSettings setObject:[NSNumber numberWithInt:AVAudioQualityHigh] forKey:AVEncoderAudioQualityKey];
    //    [recordSettings setObject:[NSNumber numberWithBool:NO] forKey:AVLinearPCMIsBigEndianKey];
    //    [recordSettings setObject:[NSNumber numberWithBool:NO] forKey:AVLinearPCMIsFloatKey];
    
    NSURL *url = [NSURL fileURLWithPath:[NSTemporaryDirectory() stringByAppendingPathComponent: [NSString stringWithFormat: @"%ld.%@",(long)sender.tag,@"wav"]]];//默认acf格式 转成wav格式 方便后面的api解析
    recorder=[[AVAudioRecorder alloc]initWithURL:url settings:recordSettings error:nil];
    [recorder record];//开始录音
}
-(void)touchUpInside:(UIButton *)sender
{
    sender.backgroundColor=[UIColor blueColor];
    sender.alpha=0;
    [recorder stop];
    
    UIButton *btn=[UIButton buttonWithType:UIButtonTypeCustom];
    btn.frame=CGRectMake(sender.frame.origin.x, sender.frame.origin.y, sender.frame.size.width-50, sender.frame.size.height);
    btn.backgroundColor=[UIColor blueColor];
    btn.tag=sender.tag;
    [btn addTarget:self action:@selector(play:) forControlEvents:UIControlEventTouchUpInside];
    [self.view addSubview:btn];
    UIImageView *image=[[UIImageView alloc]initWithImage:[UIImage imageNamed:@"play"]];
    image.frame=CGRectMake((btn.frame.size.width-20)/2,(btn.frame.size.height-20)/2, 20, 20);
    [btn addSubview:image];
}
-(void)play:(UIButton *)sender
{
    NSURL *url=[NSURL URLWithString:[NSTemporaryDirectory() stringByAppendingString:[NSString stringWithFormat:@"%ld.wav",(long)sender.tag]]];
    NSLog(@"%@",url);
    player=[[AVAudioPlayer alloc]initWithContentsOfURL:url error:nil];
    player.volume=1.5f;
    [player prepareToPlay];
    [player play];
}
#pragma mark -提交
-(void)sumbit
{
    
}
- (void)didReceiveMemoryWarning {
    [super didReceiveMemoryWarning];
    // Dispose of any resources that can be recreated.
}



@end
