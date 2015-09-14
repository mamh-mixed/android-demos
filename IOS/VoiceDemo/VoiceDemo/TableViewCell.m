//
//  TableViewCell.m
//  VoiceDemo
//
//  Created by 黄达能 on 15/8/31.
//  Copyright (c) 2015年 __VTPayment__. All rights reserved.
//

#import "TableViewCell.h"

@interface TableViewCell()

@property (weak, nonatomic) IBOutlet UIImageView *image;

@property (weak, nonatomic) IBOutlet UILabel *priceAndCount;

@property (weak, nonatomic) IBOutlet UILabel *totalMoney;

@end
@implementation TableViewCell

- (void)awakeFromNib {
    UIImageView *bgImage=[[UIImageView alloc]initWithFrame:CGRectMake(0, 0, SCREENWIDTH, 50)];
    bgImage.image=[UIImage imageNamed:@"paybg"];
    self.backgroundView=bgImage;
}

- (void)setSelected:(BOOL)selected animated:(BOOL)animated {
    [super setSelected:selected animated:animated];
}
-(void)configUI:(Model *)model
{
    self.image.image=[UIImage imageNamed:model.image];
    self.priceAndCount.text=[NSString stringWithFormat:@"%@*%@",model.price,model.clickNum];
    self.totalMoney.text=[NSString stringWithFormat:@"%.2f元",[model.price floatValue]*[model.clickNum integerValue]];
}
@end
