//
//  VTLoginViewController.m
//  VoiceDemo
//
//  Created by 司瑞华 on 15/8/26.
//  Copyright (c) 2015年 __VTPayment__. All rights reserved.
//

#import "VTLoginViewController.h"
#import "VTRegistViewController.h"
#import "MyTextField.h"

@interface VTLoginViewController ()
{
    MyTextField *login;
    MyTextField *password;
}
@end

@implementation VTLoginViewController

- (void)viewDidLoad {
    [super viewDidLoad];
    self.view.backgroundColor = [UIColor lightGrayColor];
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
    CGFloat x=30;
    
    login=[[MyTextField alloc]initWithFrame:CGRectMake(x, 100,SCREENWIDTH-x*2, 50) withImageName:@"name" withPlaceHolder:@"请输入用户名"];
    [self.view addSubview:login];
    
    password=[[MyTextField alloc]initWithFrame:CGRectMake(x,170,SCREENWIDTH-x*2, 50) withImageName:@"password" withPlaceHolder:@"请输入密码"];
    [self.view addSubview:password];
    
    UIButton *loginBtn=[UIButton buttonWithType:UIButtonTypeCustom];
    loginBtn.frame=CGRectMake(x,password.frame.origin.y+password.frame.size.height+20,90,50);
    loginBtn.layer.cornerRadius=loginBtn.frame.size.height/2;
    loginBtn.layer.masksToBounds=YES;
    loginBtn.backgroundColor=[UIColor orangeColor];
    [loginBtn setTitle:@"登入" forState:UIControlStateNormal];
    [loginBtn addTarget:self action:@selector(login) forControlEvents:UIControlEventTouchUpInside];
    [self.view addSubview:loginBtn];
    
    UIButton *registBtn=[UIButton buttonWithType:UIButtonTypeCustom];
    registBtn.layer.cornerRadius=loginBtn.frame.size.height/2;
    registBtn.layer.masksToBounds=YES;
    [registBtn setTitle:@"注册" forState:UIControlStateNormal];
    [registBtn addTarget:self action:@selector(regist) forControlEvents:UIControlEventTouchUpInside];
    registBtn.backgroundColor=[UIColor blueColor];
    registBtn.translatesAutoresizingMaskIntoConstraints=NO;
    [self.view addSubview:registBtn];
    
    [self.view addConstraint:[NSLayoutConstraint constraintWithItem:registBtn attribute:NSLayoutAttributeRight relatedBy:NSLayoutRelationEqual toItem:login attribute:NSLayoutAttributeRight multiplier:1 constant:0]];
    
    [self.view addConstraint:[NSLayoutConstraint constraintWithItem:registBtn attribute:NSLayoutAttributeWidth relatedBy:NSLayoutRelationEqual toItem:loginBtn attribute:NSLayoutAttributeWidth multiplier:1 constant:0]];
    
    [self.view addConstraint:[NSLayoutConstraint constraintWithItem:registBtn attribute:NSLayoutAttributeHeight relatedBy:NSLayoutRelationEqual toItem:loginBtn attribute:NSLayoutAttributeHeight multiplier:1 constant:0]];
    
    [self.view addConstraint:[NSLayoutConstraint constraintWithItem:registBtn attribute:NSLayoutAttributeBottom relatedBy:NSLayoutRelationEqual toItem:loginBtn attribute:NSLayoutAttributeBottom multiplier:1 constant:0]];
}
#pragma mark-登入 注册
-(void)login
{
}
-(void)regist
{
    VTRegistViewController *regist=[[VTRegistViewController alloc]init];
    [self presentViewController:regist animated:YES completion:nil];
}

-(void)touchesBegan:(NSSet *)touches withEvent:(UIEvent *)event
{
    [login resignFirstResponder];
    [password resignFirstResponder];
}
- (void)didReceiveMemoryWarning {
    [super didReceiveMemoryWarning];
    // Dispose of any resources that can be recreated.
}


@end
