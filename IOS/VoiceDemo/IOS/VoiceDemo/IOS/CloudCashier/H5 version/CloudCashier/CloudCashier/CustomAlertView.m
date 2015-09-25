//
//  CustomAlertView.m
//  TESttt
//
//  Created by 司瑞华 on 15/6/23.
//  Copyright (c) 2015年 __VTPayment__. All rights reserved.
//

#import "CustomAlertView.h"

#define TITLE_FONT_SIZE 18
#define MESSAGE_FONT_SIZE 16
#define BUTTON_FONT_SIZE 16
#define MARGIN_TOP 20
#define MARGIN_LEFT_LARGE 30
#define MARGIN_LEFT_SMALL 15
#define MARGIN_RIGHT_LARGE 30
#define MARGIN_RIGHT_SMALL 15
#define SPACE_LARGE 20
#define SPACE_SMALL 5
#define MESSAGE_LINE_SPACE 5

#define RGBA(R, G, B, A) [UIColor colorWithRed:R / 255.0 green:G / 255.0 blue:B / 255.0 alpha:A]

@interface CustomAlertView()<UITextFieldDelegate,UIScrollViewDelegate>

@property (strong, nonatomic) UIView             * backgroundView;
@property (strong, nonatomic) UIView             * titleView;
@property (strong, nonatomic) UIImageView        * iconImageView;
@property (strong, nonatomic) UILabel            * titleLabel;
@property (strong, nonatomic) UILabel            * messageLabel;
@property (strong, nonatomic) UILabel            * subtitleLbl;

@property (strong, nonatomic) NSMutableArray     * buttonArray;
@property (strong, nonatomic) NSMutableArray     * buttonTitleArray;
@property (strong, nonatomic) NSTimer            * timeCumulative ;
@property (assign, nonatomic) NSInteger          second ;


@end

CGFloat contentViewWidth;
CGFloat contentViewHeight;

NSInteger transType ;
@implementation CustomAlertView

- (instancetype)init {
    if (self = [super initWithFrame:[UIScreen mainScreen].bounds]) {
        self.backgroundColor = [UIColor clearColor];
        
        _backgroundView = [[UIView alloc] initWithFrame:self.frame];
        _backgroundView.backgroundColor = [UIColor blackColor];
        [self addSubview:_backgroundView];
    }
    return self;
}

-(instancetype)initWithTitle:(NSString *)title icon:(UIImage *)icon message:(NSString *)message subtitleMsg:(NSString *)submsg type:(NSInteger)type delegate:(id<CustomAlertViewDelegate>)delegate buttonTitles:(NSString *)buttonTitles, ... {
    if (self = [super initWithFrame:[UIScreen mainScreen].bounds]) {
        
        _icon = icon;
        _title = title;
        _message = message;
        _delegate = delegate;
        _subtitleMsg = submsg;
        _buttonArray = [NSMutableArray array];
        _buttonTitleArray = [NSMutableArray array];
        transType = type;
        
        va_list args;
        va_start(args, buttonTitles);
        if (buttonTitles)
        {
            [_buttonTitleArray addObject:buttonTitles];
            while (1)
            {
                NSString *  otherButtonTitle = va_arg(args, NSString *);
                if(otherButtonTitle == nil) {
                    break;
                } else {
                    [_buttonTitleArray addObject:otherButtonTitle];
                }
            }
        }
        va_end(args);
        
        self.backgroundColor = [UIColor clearColor];
        _backgroundView = [[UIView alloc] initWithFrame:self.frame];
        _backgroundView.backgroundColor =  [UIColor blackColor];
        [self addSubview:_backgroundView];
        [self initContentView];
    }
    return self;
}
// Init the content of content view
- (void)initContentView {
    contentViewWidth = 280 * self.frame.size.width / 320-20;
    contentViewHeight = MARGIN_TOP;
    
    _contentView = [[UIView alloc] init];
    _contentView.backgroundColor = [UIColor whiteColor];
    _contentView.layer.cornerRadius = 15.0;
    _contentView.layer.masksToBounds = YES;
    
    [self initTitleAndIcon];
    [self initMessage];
    [self initAllButtons];
    
    _contentView.frame = CGRectMake(0, 0, contentViewWidth, contentViewHeight);
    _contentView.center = self.center;
    [self addSubview:_contentView];
}

// Init the title and icon
- (void)initTitleAndIcon {
    _titleView = [[UIView alloc] init];
    if (_icon != nil) {
        _iconImageView = [[UIImageView alloc] init];
        _iconImageView.image = _icon;
        _iconImageView.frame = CGRectMake(0, 0, 30, 35*165/219);
        [_titleView addSubview:_iconImageView];
        if (transType == 1)
        {
            _iconImageView.frame = CGRectMake(0, 0, 35, 35);
            CABasicAnimation *animation = [ CABasicAnimation
                                           animationWithKeyPath: @"transform" ];
            animation.fromValue = [NSValue valueWithCATransform3D:CATransform3DIdentity];
            
            //围绕Z轴旋转，垂直与屏幕
            animation.toValue = [ NSValue valueWithCATransform3D:CATransform3DMakeRotation(M_PI, 0.0, 0.0, 1.0)];
            animation.duration = 0.5;
            //旋转效果累计，先转180度，接着再旋转180度，从而实现360旋转
            animation.cumulative = YES;
            animation.repeatCount = 1000;
            
            //在图片边缘添加一个像素的透明区域，去图片锯齿
            CGRect imageRrect = CGRectMake(0, 0,_iconImageView.frame.size.width, _iconImageView.frame.size.height);
            UIGraphicsBeginImageContext(imageRrect.size);
            [_iconImageView.image drawInRect:CGRectMake(1,1,_iconImageView.frame.size.width-2,_iconImageView.frame.size.height-2)];
            _iconImageView.image = UIGraphicsGetImageFromCurrentImageContext();
            UIGraphicsEndImageContext();
            
            [_iconImageView.layer addAnimation:animation forKey:nil];
        }
    }
    
    CGSize titleSize = [self getTitleSize];
    if (_title != nil && ![_title isEqualToString:@""]) {
        _titleLabel = [[UILabel alloc] init];
        _titleLabel.text = _title;
        _titleLabel.textColor = RGBA(28, 28, 28, 1.0);
        _titleLabel.textAlignment = NSTextAlignmentCenter;
        _titleLabel.font = [UIFont systemFontOfSize:TITLE_FONT_SIZE];
        _titleLabel.numberOfLines = 0;
        _titleLabel.lineBreakMode = NSLineBreakByWordWrapping;
        _titleLabel.frame = CGRectMake(_iconImageView.frame.origin.x + _iconImageView.frame.size.width + SPACE_SMALL, 1, titleSize.width, titleSize.height);
        [_titleView addSubview:_titleLabel];
    }
    
    _titleView.frame = CGRectMake(0, MARGIN_TOP, _iconImageView.frame.size.width + SPACE_SMALL + titleSize.width, MAX(_iconImageView.frame.size.height, titleSize.height));
    _titleView.center = CGPointMake(contentViewWidth / 2, MARGIN_TOP + _titleView.frame.size.height / 2);
    [_contentView addSubview:_titleView];
    contentViewHeight += _titleView.frame.size.height;
}


// Init the message
- (void)initMessage {
    if (_message != nil)
    {
        _messageLabel = [[UILabel alloc] init];
        _messageLabel.text = _message;
        _messageLabel.textColor = [UIColor blackColor];
        _messageLabel.numberOfLines = 0;
        _messageLabel.font = [UIFont systemFontOfSize:MESSAGE_FONT_SIZE];
        NSMutableParagraphStyle *paragraphStyle = [[NSMutableParagraphStyle alloc]init];
        paragraphStyle.lineSpacing = MESSAGE_LINE_SPACE;
        NSDictionary *attributes = @{NSParagraphStyleAttributeName:paragraphStyle};
        _messageLabel.attributedText = [[NSAttributedString alloc]initWithString:_message attributes:attributes];
        _messageLabel.textAlignment = NSTextAlignmentCenter;
        
        CGSize messageSize = [self getMessageSizeWithText:_messageLabel.text];
        _messageLabel.frame = CGRectMake(MARGIN_LEFT_LARGE, _titleView.frame.origin.y + _titleView.frame.size.height + SPACE_LARGE, MAX(contentViewWidth - MARGIN_LEFT_LARGE - MARGIN_RIGHT_LARGE, messageSize.width), messageSize.height);
        [_contentView addSubview:_messageLabel];
        contentViewHeight += SPACE_LARGE + _messageLabel.frame.size.height;
    }
    if (_subtitleMsg != nil)
    {
        _subtitleLbl = [[UILabel alloc] init];
        _subtitleLbl.text = _subtitleMsg;
        _subtitleLbl.textColor = [UIColor blackColor];
        _subtitleLbl.numberOfLines = 0;
        _subtitleLbl.font = [UIFont systemFontOfSize:MESSAGE_FONT_SIZE];
        NSMutableParagraphStyle *paragraphStyle = [[NSMutableParagraphStyle alloc]init];
        paragraphStyle.lineSpacing = MESSAGE_LINE_SPACE;
        NSDictionary *attributes = @{NSParagraphStyleAttributeName:paragraphStyle};
        _subtitleLbl.attributedText = [[NSAttributedString alloc]initWithString:_subtitleMsg attributes:attributes];
        _subtitleLbl.textAlignment = NSTextAlignmentCenter;
        
        _timeCumulative = [NSTimer scheduledTimerWithTimeInterval:1.0 target:self selector:@selector(timeCumulative) userInfo:nil repeats:NO];
        CGSize messageSize = [self getMessageSizeWithText:_subtitleLbl.text];
        _subtitleLbl.frame = CGRectMake(MARGIN_LEFT_LARGE, _messageLabel.frame.origin.y + _messageLabel.frame.size.height-20 + SPACE_LARGE, MAX(contentViewWidth - MARGIN_LEFT_LARGE - MARGIN_RIGHT_LARGE, messageSize.width), messageSize.height);
        [_contentView addSubview:_subtitleLbl];
        contentViewHeight += SPACE_LARGE + _subtitleLbl.frame.size.height;
    }
}
//定时器方法
-(NSTimer *)timeCumulative
{
    if (transType == 1)
    {
        _second++;
        [_subtitleLbl removeFromSuperview];
        _subtitleMsg = [NSString stringWithFormat:@"%lds",(long)_second];
        [self initMessage];
        return _timeCumulative;
    }
    return nil;
}

// Init all the buttons according to button titles
- (void)initAllButtons {
    if (_buttonTitleArray.count > 0) {
        contentViewHeight += SPACE_LARGE + 60;
        CGFloat buttonWidth = contentViewWidth / _buttonTitleArray.count;
        for (NSString *buttonTitle in _buttonTitleArray) {
            NSInteger index = [_buttonTitleArray indexOfObject:buttonTitle];
            UIButton *button = [[UIButton alloc]init];//[UIButton buttonWithType:UIButtonTypeRoundedRect];
            if (_subtitleLbl.text != nil)
            {
                button.frame = CGRectMake(15+index * (buttonWidth), _subtitleLbl.frame.origin.y + _subtitleLbl.frame.size.height + SPACE_LARGE, buttonWidth-30, 44);
            }else if(_subtitleLbl == nil)
            {
                button.frame = CGRectMake(15+index * (buttonWidth), _messageLabel.frame.origin.y + _messageLabel.frame.size.height + SPACE_LARGE, buttonWidth-30, 44);
            }
            
            button.layer.cornerRadius = 15;
            button.layer.masksToBounds = YES;
            button.titleLabel.font = [UIFont systemFontOfSize:BUTTON_FONT_SIZE];
            [button setTitle:buttonTitle forState:UIControlStateNormal];
            [button setTitleColor:[UIColor blackColor] forState:UIControlStateNormal];
            [button setBackgroundImage:[self imageWithColor:[UIColor colorWithRed:140/255.0 green:221/255.0 blue:233/255.0 alpha:1]] forState:UIControlStateNormal];
            [button setBackgroundImage:[self imageWithColor:[UIColor colorWithRed:0/255.0 green:187/255.0 blue:211/255.0 alpha:1]] forState:UIControlStateHighlighted];
            [button addTarget:self action:@selector(buttonWithPressed:) forControlEvents:UIControlEventTouchUpInside];
            [_buttonArray addObject:button];
            [_contentView addSubview:button];
        }
    }
}
//根据颜色生成不同的背景色图片
-(UIImage *)imageWithColor:(UIColor *)color {
    CGRect rect = CGRectMake(0.0f, 0.0f, 1.0f, 1.0f);
    UIGraphicsBeginImageContext(rect.size);
    CGContextRef context = UIGraphicsGetCurrentContext();
    
    CGContextSetFillColorWithColor(context, [color CGColor]);
    CGContextFillRect(context, rect);
    
    UIImage *image = UIGraphicsGetImageFromCurrentImageContext();
    UIGraphicsEndImageContext();
    
    return image;
}

// Get the size fo title
- (CGSize)getTitleSize {
    UIFont *font = [UIFont systemFontOfSize:TITLE_FONT_SIZE];
    
    NSMutableParagraphStyle *paragraphStyle = [[NSMutableParagraphStyle alloc] init];
    paragraphStyle.lineBreakMode = NSLineBreakByWordWrapping;
    NSDictionary *attributes = @{NSFontAttributeName:font, NSParagraphStyleAttributeName:paragraphStyle.copy};
    
    CGSize size = [_title boundingRectWithSize:CGSizeMake(contentViewWidth - (MARGIN_LEFT_SMALL + MARGIN_RIGHT_SMALL + _iconImageView.frame.size.width + SPACE_SMALL), 2000)
                                       options:NSStringDrawingUsesLineFragmentOrigin
                                    attributes:attributes context:nil].size;
    
    size.width = ceil(size.width);
    size.height = ceil(size.height);
    
    return size;
}


// Get the size of message
- (CGSize)getMessageSizeWithText:(NSString *)msg {
    UIFont *font = [UIFont systemFontOfSize:MESSAGE_FONT_SIZE];
    
    NSMutableParagraphStyle *paragraphStyle = [[NSMutableParagraphStyle alloc] init];
    paragraphStyle.lineSpacing = MESSAGE_LINE_SPACE;
    NSDictionary *attributes = @{NSFontAttributeName:font, NSParagraphStyleAttributeName:paragraphStyle.copy};
    
    CGSize size = [msg boundingRectWithSize:CGSizeMake(contentViewWidth - (MARGIN_LEFT_LARGE + MARGIN_RIGHT_LARGE), 2000)
                                         options:NSStringDrawingUsesLineFragmentOrigin
                                      attributes:attributes context:nil].size;
    
    size.width = ceil(size.width);
    size.height = ceil(size.height);
    
    return size;
}


- (void)show {
    UIWindow *window = [[UIApplication sharedApplication] keyWindow];
    NSArray *windowViews = [window subviews];
    if(windowViews && [windowViews count] > 0){
        UIView *subView = [windowViews objectAtIndex:[windowViews count]-1];
        for(UIView *aSubView in subView.subviews)
        {
            [aSubView.layer removeAllAnimations];
        }
        [subView addSubview:self];
        [self showBackground];
        [self showAlertAnimation];
    }
}

- (void)hide {
    [_contentView removeFromSuperview];
    _contentView.hidden = YES;
    if (transType == 1)
    {
        [_timeCumulative fire];
    }

    [self hideAlertAnimation];
    [self removeFromSuperview];
}

- (void)setTitleColor:(UIColor *)color fontSize:(CGFloat)size {
    if (color != nil) {
        _titleLabel.textColor = color;
    }
    
    if (size > 0) {
        _titleLabel.font = [UIFont systemFontOfSize:size];
    }
}

- (void)setMessageColor:(UIColor *)color fontSize:(CGFloat)size {
    if (color != nil) {
        _messageLabel.textColor = color;
    }
    
    if (size > 0) {
        _messageLabel.font = [UIFont systemFontOfSize:size];
    }
}

- (void)setButtonTitleColor:(UIColor *)color fontSize:(CGFloat)size atIndex:(NSInteger)index {
    UIButton *button = _buttonArray[index];
    if (color != nil) {
        [button setTitleColor:color forState:UIControlStateNormal];
    }
    
    if (size > 0) {
        button.titleLabel.font = [UIFont systemFontOfSize:size];
    }
}

- (void)showBackground
{
    _backgroundView.alpha = 0;
    [UIView beginAnimations:@"fadeIn" context:nil];
    [UIView setAnimationDuration:0.35];
    _backgroundView.alpha = 0.6;
    [UIView commitAnimations];
}

-(void)showAlertAnimation
{
    CAKeyframeAnimation * animation;
    animation = [CAKeyframeAnimation animationWithKeyPath:@"transform"];
    animation.duration = 0.30;
    animation.removedOnCompletion = YES;
    animation.fillMode = kCAFillModeForwards;
    NSMutableArray *values = [NSMutableArray array];
    [values addObject:[NSValue valueWithCATransform3D:CATransform3DMakeScale(0.9, 0.9, 1.0)]];
    [values addObject:[NSValue valueWithCATransform3D:CATransform3DMakeScale(1.1, 1.1, 1.0)]];
    [values addObject:[NSValue valueWithCATransform3D:CATransform3DMakeScale(1.0, 1.0, 1.0)]];
    animation.values = values;
    [_contentView.layer addAnimation:animation forKey:nil];
}

- (void)hideAlertAnimation {
    if (transType == 1)
    {
        [_timeCumulative fire];
    }
    [UIView beginAnimations:@"fadeIn" context:nil];
    [UIView setAnimationDuration:0.35];
    _backgroundView.alpha = 0.0;
    [UIView commitAnimations];
}

#pragma mark - 带输入框的警告框
- (instancetype)initWithTitle:(NSString *)title message:(NSString *)message subtitleMsg:(NSString *)submsg delegate:(id<CustomAlertViewDelegate>)delegate buttonTitles:(NSString *)buttonTitles, ... NS_REQUIRES_NIL_TERMINATION
{
    if (self = [super initWithFrame:[UIScreen mainScreen].bounds]) {
        _title = title;
        _message = message;
        _delegate = delegate;
        _subtitleMsg = submsg;
        _buttonArray = [NSMutableArray array];
        _buttonTitleArray = [NSMutableArray array];
        
        va_list args;
        va_start(args, buttonTitles);
        if (buttonTitles)
        {
            [_buttonTitleArray addObject:buttonTitles];
            while (1)
            {
                NSString *  otherButtonTitle = va_arg(args, NSString *);
                if(otherButtonTitle == nil) {
                    break;
                } else {
                    [_buttonTitleArray addObject:otherButtonTitle];
                }
            }
        }
        va_end(args);
        
        self.backgroundColor = [UIColor clearColor];
        _backgroundView = [[UIView alloc] initWithFrame:self.frame];
        _backgroundView.backgroundColor =  [UIColor blackColor];
        [self addSubview:_backgroundView];
        [self initTextFiledContentView];
    }
    return self;

}
-(void)initTextFiledContentView
{
    contentViewWidth = 280 * self.frame.size.width / 320;
    contentViewHeight = MARGIN_TOP;
    
    _bgScrolView = [[UIScrollView alloc]init];
    
    
    
    _contentView = [[UIView alloc] init];
    _contentView.backgroundColor = [UIColor whiteColor];
    _contentView.layer.cornerRadius = 15.0;
    _contentView.layer.masksToBounds = YES;
    
    [self initTextFiledTitle];
    [self initMessageWithTextFiled];
    [self initAllButtons];
    
    
    _bgScrolView.frame = self.frame;
    _bgScrolView.backgroundColor = [UIColor clearColor];
    [self addSubview:_bgScrolView];
    
    _contentView.frame = CGRectMake(0, 0, contentViewWidth, contentViewHeight);
    _contentView.center = self.center;
    [_bgScrolView addSubview:_contentView];
    
    //添加手势
    UITapGestureRecognizer *tap=[[UITapGestureRecognizer alloc]initWithTarget:self action:@selector(closeKeyboard:)];
    
    [self.bgScrolView addGestureRecognizer:tap];
    
}
- (BOOL)textFieldShouldBeginEditing:(UITextField *)textField
{
    [UIView animateWithDuration:0.33 animations:^{
        self.bgScrolView.contentOffset=CGPointMake(0,70);
    }];
    return YES;
}



#pragma mark 触摸背景 关闭键盘
-(void)closeKeyboard:(UIGestureRecognizer *)recognizer
{
    [UIView animateWithDuration:0.33 animations:^{
        self.bgScrolView.contentOffset=CGPointMake(0,-20);
    }];
    [self endEditing:YES];
    
}

-(void)initTextFiledTitle
{
    _titleView = [[UIView alloc] init];
    CGSize titleSize = [self getTitleSize];
    if (_title != nil && ![_title isEqualToString:@""]) {
        _titleLabel = [[UILabel alloc] init];
        _titleLabel.text = _title;
        _titleLabel.textColor = [UIColor blackColor];
        _titleLabel.textAlignment = NSTextAlignmentCenter;
        _titleLabel.font = [UIFont boldSystemFontOfSize:TITLE_FONT_SIZE];
        _titleLabel.numberOfLines = 0;
        _titleLabel.lineBreakMode = NSLineBreakByWordWrapping;
        _titleLabel.frame = CGRectMake(SPACE_SMALL, 1, titleSize.width, titleSize.height);
        [_titleView addSubview:_titleLabel];
    }
    
    _titleView.frame = CGRectMake(0, MARGIN_TOP,  SPACE_SMALL + titleSize.width,  titleSize.height);
    _titleView.center = CGPointMake(contentViewWidth / 2, MARGIN_TOP + _titleView.frame.size.height / 2);
    [_contentView addSubview:_titleView];
    contentViewHeight += _titleView.frame.size.height;
}
-(void)initMessageWithTextFiled
{
    if (_message != nil)
    {
        _messageLabel = [[UILabel alloc] init];
        _messageLabel.text = _message;
        _messageLabel.textColor = [UIColor blackColor];
        _messageLabel.numberOfLines = 0;
        _messageLabel.font = [UIFont systemFontOfSize:20];
        NSMutableParagraphStyle *paragraphStyle = [[NSMutableParagraphStyle alloc]init];
        paragraphStyle.lineSpacing = MESSAGE_LINE_SPACE;
        NSDictionary *attributes = @{NSParagraphStyleAttributeName:paragraphStyle};
        _messageLabel.attributedText = [[NSAttributedString alloc]initWithString:_message attributes:attributes];
        _messageLabel.textAlignment = NSTextAlignmentLeft;
        
        CGSize messageSize = [self getMessageSizeWithText:_messageLabel.text];
        _messageLabel.frame = CGRectMake(18, _titleView.frame.origin.y + _titleView.frame.size.height + SPACE_LARGE, MAX(contentViewWidth - MARGIN_LEFT_LARGE - MARGIN_RIGHT_LARGE-60, messageSize.width), messageSize.height+20);
        [_contentView addSubview:_messageLabel];
        
        _amountTF = [[UITextField alloc]init];
        _amountTF.frame = CGRectMake(_messageLabel.frame.origin.x+_messageLabel.frame.size.width-15, _messageLabel.frame.origin.y+3, 100, _messageLabel.frame.size.height-6);
        _amountTF.backgroundColor = RGBA(213, 213, 213, 1);
        _amountTF.delegate = self;
        [_contentView addSubview:_amountTF];
        contentViewHeight += SPACE_LARGE + _messageLabel.frame.size.height;
    }
    if (_subtitleMsg != nil)
    {
        _subtitleLbl = [[UILabel alloc] init];
        _subtitleLbl.text = _subtitleMsg;
        _subtitleLbl.textColor = [UIColor blackColor];
        _subtitleLbl.numberOfLines = 0;
        _subtitleLbl.font = [UIFont systemFontOfSize:20];
        NSMutableParagraphStyle *paragraphStyle = [[NSMutableParagraphStyle alloc]init];
        paragraphStyle.lineSpacing = MESSAGE_LINE_SPACE;
        NSDictionary *attributes = @{NSParagraphStyleAttributeName:paragraphStyle};
        _subtitleLbl.attributedText = [[NSAttributedString alloc]initWithString:_subtitleMsg attributes:attributes];
        _subtitleLbl.textAlignment = NSTextAlignmentLeft;
        
        _timeCumulative = [NSTimer scheduledTimerWithTimeInterval:1.0 target:self selector:@selector(timeCumulative) userInfo:nil repeats:NO];
        CGSize messageSize = [self getMessageSizeWithText:_subtitleLbl.text];
        _subtitleLbl.frame = CGRectMake(18, _messageLabel.frame.origin.y + _messageLabel.frame.size.height-20 + SPACE_LARGE, MAX(contentViewWidth - MARGIN_LEFT_LARGE - MARGIN_RIGHT_LARGE-60, messageSize.width), messageSize.height+20);
        [_contentView addSubview:_subtitleLbl];
        
        _passwTF = [[UITextField alloc]init];
        _passwTF.frame = CGRectMake(_subtitleLbl.frame.origin.x+_subtitleLbl.frame.size.width-15, _subtitleLbl.frame.origin.y+3, 100, _subtitleLbl.frame.size.height-6);
        _passwTF.delegate = self;
        _passwTF.secureTextEntry = YES;
        _passwTF.backgroundColor = RGBA(213, 213, 213, 1);
        [_contentView addSubview:_passwTF];
        contentViewHeight += SPACE_LARGE + _subtitleLbl.frame.size.height;
    }

}

- (void)buttonWithPressed:(UIButton *)button {
    if (_delegate && [_delegate respondsToSelector:@selector(alertView:clickedButtonAtIndex:)]) {
        NSInteger index = [_buttonTitleArray indexOfObject:button.titleLabel.text];
        [_delegate alertView:self clickedButtonAtIndex:index];
    }
    [self performSelector:@selector(hide) withObject:self afterDelay:0.2];
    
}

@end
