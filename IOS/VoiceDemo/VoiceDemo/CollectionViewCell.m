//
//  CollectionViewCell.m
//  VTCheckstandDemo
//
//  Created by 司瑞华 on 15/3/12.
//  Copyright (c) 2015年 __VTPayment__. All rights reserved.
//

#import "CollectionViewCell.h"


@implementation CollectionViewCell

-(id)initWithFrame:(CGRect)frame
{
    self = [super initWithFrame:frame];
    if (self)
    {
        self.imgView = [[UIImageView alloc]init];
        [self.contentView addSubview:_imgView];
        _titleLbl = [[UILabel alloc]init];
        _titleLbl.textColor = [UIColor blackColor];
        _titleLbl.numberOfLines = 0;
        _titleLbl.font = [UIFont systemFontOfSize:20.0f];
        _titleLbl.textAlignment = NSTextAlignmentLeft;
        [self.contentView addSubview:_titleLbl];
        _buyBtn = [ButtonPro buttonWithType:UIButtonTypeRoundedRect];
        [_buyBtn setImage:[UIImage imageNamed:@"qiugou@2x.png"] forState:UIControlStateNormal];
        [self.contentView addSubview:_buyBtn];
    }
    return self;
}
-(void)layoutSubviews
{
    [super layoutSubviews];
    float y = self.contentView.bounds.size.height-45.0f;
    _imgView.frame = CGRectMake(10.0f, 0.0f, self.contentView.bounds.size.width-20, self.contentView.bounds.size.width-20);
    _titleLbl.frame = CGRectMake(10.0f, y+15, 60.0f, 30.0f);
    _buyBtn.frame = CGRectMake(self.contentView.bounds.size.width-60, y+15, 48.0f, 28.0f);
}

@end
