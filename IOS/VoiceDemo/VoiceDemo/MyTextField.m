//
//  MyTextField.m
//  云收银
//
//  Created by 黄达能 on 15/7/27.
//  Copyright (c) 2015年 黄达能. All rights reserved.
//

#import "MyTextField.h"


@implementation MyTextField
-(id)initWithFrame:(CGRect)frame
{
    if(self=[super initWithFrame:frame])
    {
        self.backgroundColor=[UIColor lightGrayColor];
        self.layer.cornerRadius=frame.size.height/2;
        self.layer.masksToBounds=YES;
        self.layer.borderColor=[UIColor whiteColor].CGColor;
        self.layer.borderWidth=3.0;
        if(frame.size.height<40)
        {
            self.font=[UIFont fontWithName:@"Arial" size:11.0f];
        }
        else
        {
            self.font=[UIFont fontWithName:@"Arial" size:13.0f];
        }
        self.textColor=[UIColor whiteColor];
    }
    return self;
}
//placehold的起始位置
-(CGRect)textRectForBounds:(CGRect)bounds
{
    CGRect inset=CGRectMake(bounds.origin.x+55, bounds.origin.y+1, bounds.size.width-55, bounds.size.height);
    return inset;
}
//编辑时的起始位置
-(CGRect)editingRectForBounds:(CGRect)bounds
{
    CGRect inset=CGRectMake(bounds.origin.x+55, bounds.origin.y+1, bounds.size.width-55, bounds.size.height);
    return inset;
}
//设置textfield的leftView属性
-(CGRect)leftViewRectForBounds:(CGRect)bounds
{
    CGRect inset=CGRectMake(bounds.origin.x+25,bounds.origin.y+(bounds.size.height-13)/2, 15, 15);
    return inset;
}


-(id)initWithFrame:(CGRect)frame withImageName:(NSString *)imageName
{
    if(imageName)
    {
        UIImageView *image=[[UIImageView alloc]init];
        image.image=[UIImage imageNamed:imageName];
        self.leftView=image;
        self.leftViewMode=UITextFieldViewModeAlways;
    }
    return [self initWithFrame:frame];
}

-(id)initWithFrame:(CGRect)frame withImageName:(NSString *)imageName withPlaceHolder:(NSString *)placeholder
{
//    UIColor *color=[UIColor whiteColor];
//    self.attributedPlaceholder=[[NSAttributedString alloc]initWithString:placeholder attributes:@{NSForegroundColorAttributeName:color}];
    self.placeholder=placeholder;
    [self setValue:[UIColor whiteColor] forKeyPath:@"_placeholderLabel.textColor"];
    return [self initWithFrame:frame withImageName:imageName];
}


@end
