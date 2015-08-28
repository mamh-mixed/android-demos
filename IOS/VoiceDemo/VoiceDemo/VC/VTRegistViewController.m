//
//  VTRegistViewController.m
//  VoiceDemo
//
//  Created by 黄达能 on 15/8/28.
//  Copyright (c) 2015年 __VTPayment__. All rights reserved.
//

#import "VTRegistViewController.h"
#import "VTRecord_Alipay.h"
#import "MyTextField.h"


@interface VTRegistViewController ()
{
    MyTextField *name;
    MyTextField *password;
    MyTextField *passwordAgain;
}
@end

@implementation VTRegistViewController

- (void)viewDidLoad {
    [super viewDidLoad];
    self.view.backgroundColor=[UIColor lightGrayColor];
    
    CGFloat x=30;
    CGFloat space=20;
    name=[[MyTextField alloc]initWithFrame:CGRectMake(x, 80,SCREENWIDTH-x*2, 50) withImageName:@"name" withPlaceHolder:@"请输入用户名"];
    [self.view addSubview:name];
    password=[[MyTextField alloc]initWithFrame:CGRectMake(x,name.frame.origin.y+name.frame.size.height+space,SCREENWIDTH-x*2, 50) withImageName:@"password" withPlaceHolder:@"请输入密码"];
    [self.view addSubview:password];
    passwordAgain=[[MyTextField alloc]initWithFrame:CGRectMake(x,password.frame.origin.y+password.frame.size.height+space,SCREENWIDTH-x*2, 50) withImageName:@"password" withPlaceHolder:@"请再次输入密码"];
    [self.view addSubview:passwordAgain];
    
    UIButton *resgist=[UIButton buttonWithType:UIButtonTypeCustom];
    resgist.frame=CGRectMake((SCREENWIDTH-100)/2,passwordAgain.frame.origin.y+passwordAgain.frame.size.height+space,100,50);
    resgist.layer.cornerRadius=resgist.frame.size.height/2;
    resgist.layer.masksToBounds=YES;
    resgist.backgroundColor=[UIColor orangeColor];
    [resgist setTitle:@"注册" forState:UIControlStateNormal];
    [resgist addTarget:self action:@selector(register) forControlEvents:UIControlEventTouchUpInside];
    [self.view addSubview:resgist];
}
#pragma mark- 注册
-(void)register
{
    VTRecord_Alipay *alipay=[[VTRecord_Alipay alloc]init];
    [self presentViewController:alipay animated:YES completion:nil];
}
-(void)touchesBegan:(NSSet *)touches withEvent:(UIEvent *)event
{
    [name resignFirstResponder];
    [password resignFirstResponder];
    [passwordAgain resignFirstResponder];
}
- (void)didReceiveMemoryWarning {
    [super didReceiveMemoryWarning];
    // Dispose of any resources that can be recreated.
}



@end
