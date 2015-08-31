//
//  ViewController.m
//  VoiceDemo
//
//  Created by 司瑞华 on 15/8/26.
//  Copyright (c) 2015年 __VTPayment__. All rights reserved.
//

#import "ViewController.h"
#import "UICollectionViewWaterfallLayout.h"
#import "CollectionViewCell.h"
#import "ButtonPro.h"
#import "VTListViewController.h"
#import "VTLoginViewController.h"

#define KCellIdentifier         @"identifier"
#define CELL_WIDTH              self.view.frame.size.width/2
#define CELL_COLCount           2
#define VIEWSIZE                self.view.frame.size

@interface ViewController ()<UICollectionViewDelegateWaterfallLayout,UICollectionViewDataSource,UICollectionViewDelegate>
{
    UICollectionView        * _collectionView;
    float                   countMoney;
    NSMutableArray          * btnRecordArray;
    CollectionViewCell      * cell ;
}

@end

@implementation ViewController

- (void)viewDidLoad {
    [super viewDidLoad];
    clickNum = 0;
    countMoney = 0;
    
    [self createRecordeClickString];
    
    self.view.backgroundColor = [UIColor colorWithRed:247/255.0 green:246/255.0 blue:241/255.0 alpha:1];
    [self creatView];
    
    UICollectionViewWaterfallLayout * layout = [[UICollectionViewWaterfallLayout alloc]init];
    layout.delegate = self;
    layout.itemWidth = CELL_WIDTH;
    layout.columnCount = CELL_COLCount;
    _collectionView = [[UICollectionView alloc]initWithFrame:CGRectMake(0.0f, 64.0f, self.view.frame.size.width, self.view.frame.size.height-64) collectionViewLayout:layout];
    _collectionView.delegate = self;
    _collectionView.dataSource = self;
    
    _collectionView.backgroundColor = [UIColor clearColor];
    [_collectionView registerClass:[CollectionViewCell class] forCellWithReuseIdentifier:KCellIdentifier];
    [self.view addSubview:_collectionView];
}
-(void)goPayView
{
    VTListViewController * vc = [[VTListViewController alloc]init];
    vc.countMoney = countMoney;
    vc.countNum = clickNum;
    vc.cellContentArray = [NSMutableArray arrayWithArray:btnRecordArray];
    [self presentViewController:vc animated:YES completion:nil];
}
-(void)goLoginView
{
    VTLoginViewController * vc = [[VTLoginViewController alloc]init];
    [self presentViewController:vc animated:NO completion:nil];
}
//导航头上的控件
-(void)creatView
{
     
    UILabel * lbl = [[UILabel alloc]initWithFrame:CGRectMake(0, 20, self.view.frame.size.width, 44)];
    lbl.backgroundColor = [UIColor whiteColor];
    lbl.text = @"精品特卖";
    lbl.textAlignment = NSTextAlignmentCenter;
    lbl.textColor = [UIColor colorWithRed:194/255.0 green:176/255.0 blue:146/255.0 alpha:1];
    [self.view addSubview:lbl];
    
    UIButton * logBtn = [UIButton buttonWithType:UIButtonTypeRoundedRect];
    logBtn.frame = CGRectMake(10, 30, 45, 30) ;
    [logBtn setTitle:@"登录" forState:UIControlStateNormal];
    [logBtn setTitleColor:[UIColor redColor] forState:UIControlStateNormal];
    [logBtn addTarget:self action:@selector(goLoginView) forControlEvents:UIControlEventTouchUpInside];
    [self.view addSubview:logBtn];
    
    UIButton * btn = [UIButton buttonWithType:UIButtonTypeRoundedRect];
    [btn setImage:[UIImage imageNamed:@"buy@2x.png"] forState:UIControlStateNormal];
    btn.frame = CGRectMake(VIEWSIZE.width-60, 23, 44, 44);
    [btn addTarget:self action:@selector(goPayView) forControlEvents:UIControlEventTouchUpInside];
    [self.view addSubview:btn];
    UILabel * numLbl = [[UILabel alloc]initWithFrame:CGRectMake(VIEWSIZE.width-40, 25, 20, 20)];
    numLbl.backgroundColor = [UIColor colorWithRed:231/255.0 green:76/255.0 blue:60/255.0 alpha:1];
    numLbl.layer.cornerRadius = 10;
    numLbl.layer.masksToBounds = YES;
    numLbl.textColor = [UIColor whiteColor];
    numLbl.font = [UIFont systemFontOfSize:14];
    numLbl.text = [NSString stringWithFormat:@"%d",clickNum];
    //记录购物件数小lbl
    if ([numLbl.text isEqualToString:@"0"])
    {
        numLbl.hidden = YES;
    }else
    {
        numLbl.hidden = NO;
    }
    numLbl.textAlignment = NSTextAlignmentCenter;
    [self.view addSubview:numLbl];
}
#pragma mark - UICollectionViewDelegate and dataSource
-(NSInteger)collectionView:(UICollectionView *)collectionView numberOfItemsInSection:(NSInteger)section
{
    return 10;
}
//每个item内容赋值
-(UICollectionViewCell *)collectionView:(UICollectionView *)collectionView cellForItemAtIndexPath:(NSIndexPath *)indexPath
{
    static NSString * cellIdentifier = KCellIdentifier;
    
    cell = (CollectionViewCell *) [collectionView dequeueReusableCellWithReuseIdentifier:cellIdentifier forIndexPath:indexPath];
    cell.titleLbl.text = [NSString stringWithFormat:@"%@",[[self creatArray]objectAtIndex:indexPath.row]];
    cell.imgView.image = [UIImage imageNamed:[NSString stringWithFormat:@"n%ld.jpg",indexPath.row+1]];
    [cell.buyBtn addTarget:self action:@selector(btnClick:) forControlEvents:UIControlEventTouchUpInside];
    cell.buyBtn.tag = indexPath.row;
    NSMutableString * str = [btnRecordArray objectAtIndex:indexPath.row];
    cell.buyBtn.isAdd = NO;
    cell.buyBtn.btnClickNum = 0;
    str =[NSMutableString stringWithFormat:@"%ld",(long)cell.buyBtn.btnClickNum];
    return cell;
}
- (CGFloat)collectionView:(UICollectionView *)collectionView
                   layout:(UICollectionViewWaterfallLayout *)collectionViewLayout
 heightForItemAtIndexPath:(NSIndexPath *)indexPath
{
    return CELL_WIDTH + 15.0;
}

-(void)btnClick:(ButtonPro *)btn
{
    clickNum++;
    btn.isAdd = YES;
    btn.btnClickNum ++;
    
    btn.isAdd = NO;
    NSString * str = [NSString stringWithFormat:@"%ld",(long)btn.btnClickNum];
    [btnRecordArray replaceObjectAtIndex:btn.tag withObject:str];
    
    float currentMoney = [[[self creatArray]objectAtIndex:btn.tag] floatValue];
    countMoney = countMoney + currentMoney;
    
    [self creatView];
    //NSLog(@"djfksfjksl点我啦 %d  ----%ld ---------array  %@",clickNum,(long)btn.tag,btnRecordArray);
}
//两个array 分别放了每个item上的图片名称以及价格
-(NSArray *)creatArray
{
    NSArray * mArray = [NSArray arrayWithObjects:@"0.01",@"10",@"5",@"3",@"99.9",@"1",@"8.88",@"0.02",@"16.8",@"20", nil];
    return mArray;
}
//创建记录各个cell点击次数的字符串，并加到数组中
-(void)createRecordeClickString
{
    btnRecordArray = [[NSMutableArray alloc]initWithObjects:@"0",@"0",@"0",@"0",@"0",@"0",@"0",@"0",@"0",@"0", nil];
}


- (void)didReceiveMemoryWarning {
    [super didReceiveMemoryWarning];
    
}

@end
