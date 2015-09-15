//
//  VTListViewController.m
//  VoiceDemo
//
//  Created by 司瑞华 on 15/8/26.
//  Copyright (c) 2015年 __VTPayment__. All rights reserved.
//

#import "VTListViewController.h"
#import "TableViewCell.h"
#import "Model.h"
#import "VTPayViewController.h"

@interface VTListViewController ()<UITableViewDataSource,UITableViewDelegate>
{
    UIScrollView *scrollView;
    int GoodNum;//购买的产品数量
}
@end

@implementation VTListViewController

- (void)viewDidLoad {
    [super viewDidLoad];
    UIImageView *image=[[UIImageView alloc]initWithFrame:CGRectMake(0, 0, SCREENWIDTH, SCREENHEIGHT)];
    image.image=[UIImage imageNamed:@"paybg"];
    [self.view addSubview:image];
    GoodNum=0;
    for (int i=0; i<_cellContentArray.count; i++) {
        if(![_cellContentArray[i] isEqualToString:@"0"])
        {
            GoodNum++;
        }
    }
    [self createUI];
    [self createTableView];
}
-(void)createUI
{
    UIImageView *header=[[UIImageView alloc]initWithFrame:CGRectMake(0, 20, SCREENWIDTH, 44)];
    header.userInteractionEnabled=YES;
    UILabel *label=[[UILabel alloc]initWithFrame:CGRectMake(0, 0, SCREENWIDTH, 44)];
    label.backgroundColor=[UIColor whiteColor];
    label.text=@"购物车";
    label.textAlignment=NSTextAlignmentCenter;
    [header addSubview:label];
    UIButton *back=[UIButton buttonWithType:UIButtonTypeCustom];
    [back setImage:[UIImage imageNamed:@"back"] forState:UIControlStateNormal];
    [back addTarget:self action:@selector(back) forControlEvents:UIControlEventTouchUpInside];
    back.frame=CGRectMake(10, 10, 12, 21);
    [header addSubview:back];
    [self.view addSubview:header];
    
    scrollView=[[UIScrollView alloc]initWithFrame:CGRectMake(0, 64, SCREENWIDTH, SCREENHEIGHT-64)];
    UIImageView *image=[[UIImageView alloc]initWithFrame:CGRectMake(0, 0, scrollView.frame.size.width, scrollView.frame.size.height)];
    image.image=[UIImage imageNamed:@"paybg"];
    if (SCREENHEIGHT<600) {
        scrollView.contentSize=CGSizeMake(SCREENWIDTH, 50*GoodNum+120);
    }
    else{
        scrollView.contentSize=CGSizeMake(SCREENWIDTH, 60*GoodNum+120);
    }
    scrollView.showsVerticalScrollIndicator=NO;
    [scrollView addSubview:image];
    [self.view addSubview:scrollView];
}
-(void)createTableView
{
    UITableView *table;
    CGFloat height;
    if (SCREENHEIGHT<600) {
        height=50;
    }
    else{
        height=60;
    }
    table=[[UITableView alloc]initWithFrame:CGRectMake(0, 0, SCREENWIDTH, height*GoodNum) style:UITableViewStylePlain];
    
    UIImageView *image=[[UIImageView alloc]initWithFrame:CGRectMake(0, 0,scrollView.frame.size.width, scrollView.frame.size.height)];
    image.image=[UIImage imageNamed:@"paybg"];
    table.backgroundView=image;
    table.bounces=NO;
    table.dataSource=self;
    table.delegate=self;
    [scrollView addSubview:table];
    
    UILabel *countPrice=[[UILabel alloc]initWithFrame:CGRectMake(5, table.frame.size.height+2, SCREENWIDTH-5, 20)];
    countPrice.textColor=[UIColor redColor];
    countPrice.text=[NSString stringWithFormat:@"CNY%.2f",_countMoney];
    [scrollView addSubview:countPrice];
    
    UIButton *btn=[UIButton buttonWithType:UIButtonTypeCustom];
    btn.frame=CGRectMake(30, countPrice.frame.origin.y+60, SCREENWIDTH-30*2, height-10);
    btn.layer.cornerRadius=5;
    btn.layer.masksToBounds=YES;
    btn.backgroundColor=[UIColor whiteColor];
    [btn addTarget:self action:@selector(pay) forControlEvents:UIControlEventTouchUpInside];
    [btn setTitle:@"一键支付" forState:UIControlStateNormal];
    [btn setTitleColor:[UIColor blackColor] forState:UIControlStateNormal];
    [scrollView addSubview:btn];
}
#pragma mark -UITableViewDelegate
-(NSInteger)tableView:(UITableView *)tableView numberOfRowsInSection:(NSInteger)section
{
    return GoodNum;
}
-(CGFloat)tableView:(UITableView *)tableView heightForRowAtIndexPath:(NSIndexPath *)indexPath
{
    if (SCREENHEIGHT<600) {
        return 50;
    }
    return 60;
}
-(UITableViewCell *)tableView:(UITableView *)tableView cellForRowAtIndexPath:(NSIndexPath *)indexPath
{
    static NSString *cellID=@"cellID";
    TableViewCell *cell=[tableView dequeueReusableCellWithIdentifier:cellID];
    if (cell==nil) {
        cell=[[[NSBundle mainBundle]loadNibNamed:@"TableViewCell" owner:self options:nil]lastObject];
    }
    cell.selectionStyle=UITableViewCellSelectionStyleNone;
    Model *model=_dataArray[indexPath.row];
    [cell configUI:model];
    return cell;
}

#pragma mark- 一键支付
-(void)pay
{
    VTPayViewController *pay=[[VTPayViewController alloc]init];
    [self presentViewController:pay animated:YES completion:nil];
}
- (void)didReceiveMemoryWarning {
    [super didReceiveMemoryWarning];
}
-(void)back
{
    [self dismissViewControllerAnimated:YES completion:nil];
}
@end
