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

#define kDatabaseName @"database.sqlite3"

#define SCREENWIDTH [UIScreen mainScreen].bounds.size.width
#define SCREENHEIGHT [UIScreen mainScreen].bounds.size.height

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
    
    NSArray *paths = NSSearchPathForDirectoriesInDomains(NSDocumentDirectory, NSUserDomainMask, YES);
    NSString *documentsDirectory = [paths objectAtIndex:0];
    self.databaseFilePath = [documentsDirectory stringByAppendingPathComponent:kDatabaseName];
    NSLog(@"%@         _____________________ ",self.databaseFilePath);
    
    sqlite3 *database;
    if (sqlite3_open([self.databaseFilePath UTF8String] , &database) != SQLITE_OK) {
        sqlite3_close(database);
        NSAssert(0, @"打开数据库失败！");
    }
    NSString *createSQL = @"CREATE TABLE IF NOT EXISTS UserList (用户名 TEXT PRIMARY KEY, 密码 TEXT);";
    char *errorMsg;
    if (sqlite3_exec(database, [createSQL UTF8String], NULL, NULL, &errorMsg) != SQLITE_OK) {
        sqlite3_close(database);
        NSAssert(0, @"创建数据库表错误: %s", errorMsg);
    }
    //关闭数据库
    sqlite3_close(database);
    
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
//    if([password.text isEqualToString:passwordAgain.text])
//    {
//        RegisterTable * table = [[RegisterTable alloc]init];
//        table.username = name.text;
//        table.password = password.text;
//        [RegisterTableDAO insertObject:table complete:^(NSString *isExists) {
//            if ([isExists isEqualToString:@"exists"])
//            {
//                [self alertWithMessage:@"用户已存在"];
//            }else if ([isExists isEqualToString:@"success"])
//            {
//                
//            }
//        }];
// 
//    }else
//    {
//        [self alertWithMessage:@"密码不一致"];
//    }
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
    // Dispose of any resources that can be recreated.
}



@end
