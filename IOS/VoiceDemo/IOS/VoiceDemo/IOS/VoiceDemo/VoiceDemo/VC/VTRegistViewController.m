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
#import "sqlite3.h"
#import "RegisterTable.h"

#import "VTRecord_Uppay.h"
#import "VTRecord_Wxpay.h"

#define kDatabaseName @"database.sqlite3"

@interface VTRegistViewController ()
{
    MyTextField *name;
    MyTextField *password;
    MyTextField *passwordAgain;
}

@property (copy, nonatomic) NSString *databaseFilePath;

@end

@implementation VTRegistViewController

- (void)viewDidLoad {
    [super viewDidLoad];
    self.view.backgroundColor=[UIColor lightGrayColor];
    UIButton *back=[UIButton buttonWithType:UIButtonTypeCustom];
    back.frame=CGRectMake(10, 30, 12, 21);
    [back setImage:[UIImage imageNamed:@"back"] forState:UIControlStateNormal];
    [back addTarget:self action:@selector(back) forControlEvents:UIControlEventTouchUpInside];
    [self.view addSubview:back];
    
    NSArray *paths = NSSearchPathForDirectoriesInDomains(NSDocumentDirectory, NSUserDomainMask, YES);
    NSString *documentsDirectory = [paths objectAtIndex:0];
    self.databaseFilePath = [documentsDirectory stringByAppendingPathComponent:kDatabaseName];
    NSLog(@"%@         _____________________ ",self.databaseFilePath);
    
    sqlite3 *database;
    if (sqlite3_open([self.databaseFilePath UTF8String] , &database) != SQLITE_OK) {
        sqlite3_close(database);
        NSAssert(0, @"打开数据库失败！");
    }
    NSString *createSQL = @"CREATE TABLE IF NOT EXISTS UserList (username TEXT PRIMARY KEY, password TEXT ,isUsed TEXT , time TEXT);";
    char *errorMsg;
    if (sqlite3_exec(database, [createSQL UTF8String], NULL, NULL, &errorMsg) != SQLITE_OK) {
        sqlite3_close(database);
        NSAssert(0, @"创建数据库表错误: %s", errorMsg);
    }
    //关闭数据库
    sqlite3_close(database);
    
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
    CGFloat space=20;
    name=[[MyTextField alloc]initWithFrame:CGRectMake(x, 80,SCREENWIDTH-x*2, height) withImageName:@"name" withPlaceHolder:@"请输入用户名"];
    [self.view addSubview:name];
    password=[[MyTextField alloc]initWithFrame:CGRectMake(x,name.frame.origin.y+name.frame.size.height+space,SCREENWIDTH-x*2, height) withImageName:@"password" withPlaceHolder:@"请输入密码"];
    [self.view addSubview:password];
    passwordAgain=[[MyTextField alloc]initWithFrame:CGRectMake(x,password.frame.origin.y+password.frame.size.height+space,SCREENWIDTH-x*2, height) withImageName:@"password" withPlaceHolder:@"请再次输入密码"];
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
-(void)back
{
    [self dismissViewControllerAnimated:YES completion:nil];
}
#pragma mark- 注册
-(void)register
{
#if 1
    if([password.text isEqualToString:passwordAgain.text]&&password.text)
    {
        RegisterTable * table = [[RegisterTable alloc]init];
        NSDateFormatter *formatter=[[NSDateFormatter alloc]init];
        [formatter setDateFormat:@"yyyyMMddHHmmss"];
        NSString *dateTime=[formatter stringFromDate:[NSDate date]];
        if(![name.text isEqualToString:@""]){
            table.username = name.text;
            table.password = password.text;
            table.isUsed = @"1";
            table.time=dateTime;
        }
        else{
            [self alertWithMessage:@"用户名不能为空"];
            return;
        }
        [RegisterTableDAO insertObject:table complete:^(NSString *isExists) {
            if ([isExists isEqualToString:@"exists"])
            {
                [self alertWithMessage:@"用户已存在"];
                return ;
            }
            else if ([isExists isEqualToString:@"success"])
            {
                VTRecord_Alipay *alipay=[[VTRecord_Alipay alloc]init];
                [self presentViewController:alipay animated:YES completion:nil];
                [RegisterTableDAO changeisUsedByName:table.username];
                return;
            }
        }];
    }else
    {
        if(password.text){
            [self alertWithMessage:@"两次密码不一致"];
            return;
        }
        else{
            [self alertWithMessage:@"密码不能为空"];
            return;
        }
    }
#endif

}
-(void)alertWithMessage:(NSString *)message
{
    UIAlertView * alert = [[UIAlertView alloc]initWithTitle:@"温馨提示" message:message delegate:self cancelButtonTitle:@"我知道了" otherButtonTitles:nil, nil];
    [alert show];
}
-(void)touchesBegan:(NSSet *)touches withEvent:(UIEvent *)event
{
    [name resignFirstResponder];
    [password resignFirstResponder];
    [passwordAgain resignFirstResponder];
}
- (void)didReceiveMemoryWarning {
    [super didReceiveMemoryWarning];
}

@end
