//
//  MyTextField.h
//  云收银
//
//  Created by 黄达能 on 15/7/27.
//  Copyright (c) 2015年 黄达能. All rights reserved.
//

#import <UIKit/UIKit.h>


@interface MyTextField : UITextField

-(id)initWithFrame:(CGRect)frame;

-(id)initWithFrame:(CGRect)frame withImageName:(NSString *)imageName;

-(id)initWithFrame:(CGRect)frame withImageName:(NSString *)imageName withPlaceHolder:(NSString *)placeholder;

@end

