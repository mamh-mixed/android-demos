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
#import "RegisterTable.h"

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
    CGFloat x;
    CGFloat height;
    if (SCREENHEIGHT<600) {
        x=30;
        height=50;
    }
    else{
        x=40;
        height=60;
    }
    
    login=[[MyTextField alloc]initWithFrame:CGRectMake(x, 100,SCREENWIDTH-x*2, height) withImageName:@"name" withPlaceHolder:@"请输入用户名"];
    [self.view addSubview:login];
    
    password=[[MyTextField alloc]initWithFrame:CGRectMake(x,login.frame.origin.y+login.frame.size.height+30,SCREENWIDTH-x*2, height) withImageName:@"password" withPlaceHolder:@"请输入密码"];
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
    NSString *message;
    if ([login.text isEqualToString:@""]) {
        message=@"用户名不能为空";
        [self alertWithMessage:message];
        return;
    }
    else if([password.text isEqualToString:@""]){
        message=@"密码不能为空";
        [self alertWithMessage:message];
        return;
    }
    if ([RegisterTableDAO isExistTheName:login.text]) {
        NSMutableDictionary *dict=[RegisterTableDAO getObjectByName:login.text];
        if ([[dict objectForKey:@"password"]isEqualToString:password.text]) {
            //登入成功
            dispatch_async(dispatch_get_global_queue(0, 0), ^{
                [RegisterTableDAO changeisUsedByName:login.text];
                dispatch_async(dispatch_get_main_queue(), ^{
                    [self dismissViewControllerAnimated:YES completion:nil];
                });
            });
            return;
        }
        else{
            message=@"密码不正确";
            [self alertWithMessage:message];
            return;
        }
    }
    else{
        message=@"用户名不存在";
        [self alertWithMessage:message];
        return;
    }
}
-(void)alertWithMessage:(NSString *)message
{
    UIAlertView * alert = [[UIAlertView alloc]initWithTitle:@"温馨提示" message:message delegate:self cancelButtonTitle:@"我知道了" otherButtonTitles:nil, nil];
    [alert show];
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
