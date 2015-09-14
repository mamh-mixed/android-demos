//
//  VTPayViewController.m
//  VoiceDemo
//
//  Created by 黄达能 on 15/9/9.
//  Copyright (c) 2015年 __VTPayment__. All rights reserved.
//

#import "VTPayViewController.h"
#import "DetectRequest.h"
#import "RegisterTable.h"

@interface VTPayViewController ()

@property (nonatomic,strong) NSString *path;//音频的路径

@end

@implementation VTPayViewController

@synthesize recorder;

- (void)viewDidLoad {
    [super viewDidLoad];
    self.view.backgroundColor=[UIColor whiteColor];
    UIImageView *imageView=[[UIImageView alloc]initWithFrame:[UIScreen mainScreen].bounds];
    imageView.image=[UIImage imageNamed:@"paybg"];
    [self.view addSubview:imageView];
    
    _label=[[UILabel alloc]initWithFrame:CGRectMake(0, 140, SCREENWIDTH, 50)];
    _label.backgroundColor=[UIColor blackColor];
    _label.textAlignment=NSTextAlignmentCenter;
    _label.textColor=[UIColor whiteColor];
    [self.view addSubview:_label];
    
    UIButton *back=[UIButton buttonWithType:UIButtonTypeCustom];
    back.frame=CGRectMake(10, 30, 12, 21);
    [back setImage:[UIImage imageNamed:@"back"] forState:UIControlStateNormal];
    [back addTarget:self action:@selector(back) forControlEvents:UIControlEventTouchUpInside];
    [self.view addSubview:back];
    [self createUI];
    
}
-(void)back
{
    [self dismissViewControllerAnimated:YES completion:nil];
}
-(void)createUI
{
    UILabel *label=[[UILabel alloc]initWithFrame:CGRectMake(0, 40, SCREENWIDTH, 30)];
    label.textAlignment=NSTextAlignmentCenter;
    label.text=@"按住话筒录音，选择支付方式 :";
    [self.view addSubview:label];
    
    UILabel *lbl=[[UILabel alloc]initWithFrame:CGRectMake(0, 80, SCREENWIDTH, 30)];
    lbl.textAlignment=NSTextAlignmentCenter;
    lbl.text=@"支付宝，银联或者微信";
    lbl.textColor=[UIColor redColor];
    [self.view addSubview:lbl];
    
    CGFloat height=50;
    UIButton *record=[UIButton buttonWithType:UIButtonTypeCustom];
    record.frame=CGRectMake(0,SCREENHEIGHT-200 , SCREENWIDTH, height);
    record.backgroundColor=[UIColor whiteColor];
    
    UIImageView *image=[[UIImageView alloc]initWithFrame:CGRectMake((SCREENWIDTH-20)/2,(height-32)/2, 20, 32)];
    image.image=[UIImage imageNamed:@"microphone"];
    [record addSubview:image];
    [record addTarget:self action:@selector(record:) forControlEvents:UIControlEventTouchDown];
    [record addTarget:self action:@selector(touchUpInside:) forControlEvents:UIControlEventTouchUpInside];
    [self.view addSubview:record];
}
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
    _path=[NSTemporaryDirectory() stringByAppendingPathComponent:[NSString stringWithFormat: @"识别录音.%@",@"wav"]];
    NSURL *url = [NSURL fileURLWithPath:_path];
    recorder=[[AVAudioRecorder alloc]initWithURL:url settings:recordSettings error:nil];
    [recorder record];
}
-(void)touchUpInside:(UIButton *)sender
{
    sender.backgroundColor=[UIColor whiteColor];
    [recorder stop];
    if ([self testAudioDuration]) {//音频时长超过0.6s
        _label.text=nil;
        //发送检测的请求
        DetectRequest *dRequest=[[DetectRequest alloc]init];
        dRequest.viewController=self;
        [dRequest connectionNet:_path];
    }
}
-(BOOL)testAudioDuration
{
    AVURLAsset *audioAsset;
    audioAsset=[AVURLAsset URLAssetWithURL:[NSURL URLWithString:[NSString stringWithFormat:@"file://%@",_path]] options:nil];
    CMTime audioDuration =audioAsset.duration;
    float audioDurationSeconds=CMTimeGetSeconds(audioDuration);
    if (audioDurationSeconds<0.6) {
        UIAlertView *alert=[[UIAlertView alloc]initWithTitle:nil message:@"录音时间过短请重录" delegate:self cancelButtonTitle:@"取消" otherButtonTitles: nil];
        [alert show];
        return NO;
    }
    return YES;
}
- (void)didReceiveMemoryWarning {
    [super didReceiveMemoryWarning];
}

@end
