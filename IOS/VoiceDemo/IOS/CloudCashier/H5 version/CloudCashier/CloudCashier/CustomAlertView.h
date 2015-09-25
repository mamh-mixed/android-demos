//
//  CustomAlertView.h
//  TESttt
//
//  Created by 司瑞华 on 15/6/23.
//  Copyright (c) 2015年 __VTPayment__. All rights reserved.
//

#import <UIKit/UIKit.h>

@protocol CustomAlertViewDelegate ;

@interface CustomAlertView : UIView
@property (strong, nonatomic) UIView           * contentView;
@property (strong, nonatomic) UIScrollView     * bgScrolView;
@property (strong, nonatomic) UIImage          * icon;
@property (strong, nonatomic) NSString         * title;
@property (strong, nonatomic) NSString         * message;
@property (strong, nonatomic) NSString         * subtitleMsg;
@property (strong, nonatomic) UITextField      * amountTF;
@property (strong, nonatomic) UITextField      * passwTF;
@property (weak, nonatomic) id<CustomAlertViewDelegate> delegate;

- (instancetype)initWithTitle:(NSString *)title icon:(UIImage *)icon message:(NSString *)message subtitleMsg:(NSString *)submsg type:(NSInteger)type delegate:(id<CustomAlertViewDelegate>)delegate buttonTitles:(NSString *)buttonTitles, ... NS_REQUIRES_NIL_TERMINATION;

- (instancetype)initWithTitle:(NSString *)title message:(NSString *)message subtitleMsg:(NSString *)submsg delegate:(id<CustomAlertViewDelegate>)delegate buttonTitles:(NSString *)buttonTitles, ... NS_REQUIRES_NIL_TERMINATION;

// Show the alert view in current window
- (void)show;

// Hide the alert view
- (void)hide;

// Set the color and font size of title, if color is nil, default is black. if fontsize is 0, default is 14
- (void)setTitleColor:(UIColor *)color fontSize:(CGFloat)size;

// Set the color and font size of message, if color is nil, default is black. if fontsize is 0, default is 12
- (void)setMessageColor:(UIColor *)color fontSize:(CGFloat)size;

// Set the color and font size of button at the index, if color is nil, default is black. if fontsize is 0, default is 16
- (void)setButtonTitleColor:(UIColor *)color fontSize:(CGFloat)size atIndex:(NSInteger)index;


@end

@protocol CustomAlertViewDelegate <NSObject>

- (void)alertView:(CustomAlertView *)alertView clickedButtonAtIndex:(NSInteger)buttonIndex;

@end








