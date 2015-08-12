//
//  VTScannerView.m
//  CloudCashier
//
//  Created by 司瑞华 on 15/7/14.
//  Copyright (c) 2015年 __VTPayment__. All rights reserved.
//

#import "VTScannerView.h"
#import "VTMenuView.h"

static NSTimeInterval kRhLineanimateDuration = 0.02;
@implementation VTScannerView
{
    UIImageView               * _rhLineImgView;
    CGFloat                   _rhLineY;
    VTMenuView                * _vtMenuView;
}

-(instancetype)initWithFrame:(CGRect)frame
{
    self = [super initWithFrame:frame];
    if (self)
    {
        
    }
    return self;
}
-(void)layoutSubviews
{
    [super layoutSubviews];
    if (!_rhLineImgView)
    {
        [self initRhLineImgView];
        self.timer = [NSTimer scheduledTimerWithTimeInterval:kRhLineanimateDuration target:self selector:@selector(showTime) userInfo:nil repeats:YES];
        [_timer fire];
    }
    if (!_vtMenuView)
    {
    
    }
}

-(void)initRhLineImgView
{
    _rhLineImgView  = [[UIImageView alloc] initWithFrame:CGRectMake(self.bounds.size.width / 2 - self.transparentArea.width / 2, self.bounds.size.height / 2 - self.transparentArea.height / 2, self.transparentArea.width, 2)];
    _rhLineImgView.image = [UIImage imageNamed:@"qr_scan_line"];
    _rhLineImgView.contentMode = UIViewContentModeScaleAspectFill;
    [self addSubview:_rhLineImgView];
    _rhLineY = _rhLineImgView.frame.origin.y;
    
}

//扫描页面周围的蒙版 及视图
-(void)initRhMenuView
{
}
-(void)showTime
{
    //一直循环运动
    [UIView animateWithDuration:kRhLineanimateDuration animations:^{
        CGRect rect = _rhLineImgView.frame;
        rect.origin.y = _rhLineY;
        _rhLineImgView.frame = rect;
    } completion:^(BOOL finished) {
        CGFloat maxBorder = self.frame.size.height/2 + self.transparentArea.height/2 - 4;
        if (_rhLineY > maxBorder)//如果中间的扫描线已经到低端，就重置位置
        {
            _rhLineY = self.frame.size.height / 2 - self.transparentArea.height /2;
        }
        _rhLineY++;
    }];
    
}

- (void)drawRect:(CGRect)rect {
    
    //整个二维码扫描界面的颜色
    CGSize screenSize =[UIScreen mainScreen].bounds.size;
    CGRect screenDrawRect =CGRectMake(0, 0, screenSize.width,screenSize.height);
    
    //中间清空的矩形框
    CGRect clearDrawRect = CGRectMake(screenDrawRect.size.width / 2 - self.transparentArea.width / 2,
                                      screenDrawRect.size.height / 2 - self.transparentArea.height / 2,
                                      self.transparentArea.width,self.transparentArea.height);
    
    CGContextRef ctx = UIGraphicsGetCurrentContext();
    [self addScreenFillRect:ctx rect:screenDrawRect];
    
    [self addCenterClearRect:ctx rect:clearDrawRect];
    
    [self addWhiteRect:ctx rect:clearDrawRect];
    
    [self addCornerLineWithContext:ctx rect:clearDrawRect];
}
- (void)addScreenFillRect:(CGContextRef)ctx rect:(CGRect)rect {
    
    CGContextSetRGBFillColor(ctx, 40 / 255.0,40 / 255.0,40 / 255.0,0.5);
    CGContextFillRect(ctx, rect);   //draw the transparent layer
}
- (void)addCenterClearRect :(CGContextRef)ctx rect:(CGRect)rect {
    
    CGContextClearRect(ctx, rect);  //clear the center rect  of the layer
}
- (void)addWhiteRect:(CGContextRef)ctx rect:(CGRect)rect {
    
    CGContextStrokeRect(ctx, rect);
    CGContextSetRGBStrokeColor(ctx, 1, 1, 1, 1);
    CGContextSetLineWidth(ctx, 0.8);
    CGContextAddRect(ctx, rect);
    CGContextStrokePath(ctx);
}

- (void)addCornerLineWithContext:(CGContextRef)ctx rect:(CGRect)rect{
    
    //画四个边角
    CGContextSetLineWidth(ctx, 2);
    CGContextSetRGBStrokeColor(ctx, 83 /255.0, 239/255.0, 111/255.0, 1);//绿色
    
    //左上角
    CGPoint poinsTopLeftA[] = {
        CGPointMake(rect.origin.x+0.7, rect.origin.y),
        CGPointMake(rect.origin.x+0.7 , rect.origin.y + 15)
    };
    
    CGPoint poinsTopLeftB[] = {CGPointMake(rect.origin.x, rect.origin.y +0.7),CGPointMake(rect.origin.x + 15, rect.origin.y+0.7)};
    [self addLine:poinsTopLeftA pointB:poinsTopLeftB ctx:ctx];
    
    //左下角
    CGPoint poinsBottomLeftA[] = {CGPointMake(rect.origin.x+ 0.7, rect.origin.y + rect.size.height - 15),CGPointMake(rect.origin.x +0.7,rect.origin.y + rect.size.height)};
    CGPoint poinsBottomLeftB[] = {CGPointMake(rect.origin.x , rect.origin.y + rect.size.height - 0.7) ,CGPointMake(rect.origin.x+0.7 +15, rect.origin.y + rect.size.height - 0.7)};
    [self addLine:poinsBottomLeftA pointB:poinsBottomLeftB ctx:ctx];
    
    //右上角
    CGPoint poinsTopRightA[] = {CGPointMake(rect.origin.x+ rect.size.width - 15, rect.origin.y+0.7),CGPointMake(rect.origin.x + rect.size.width,rect.origin.y +0.7 )};
    CGPoint poinsTopRightB[] = {CGPointMake(rect.origin.x+ rect.size.width-0.7, rect.origin.y),CGPointMake(rect.origin.x + rect.size.width-0.7,rect.origin.y + 15 +0.7 )};
    [self addLine:poinsTopRightA pointB:poinsTopRightB ctx:ctx];
    
    CGPoint poinsBottomRightA[] = {CGPointMake(rect.origin.x+ rect.size.width -0.7 , rect.origin.y+rect.size.height+ -15),CGPointMake(rect.origin.x-0.7 + rect.size.width,rect.origin.y +rect.size.height )};
    CGPoint poinsBottomRightB[] = {CGPointMake(rect.origin.x+ rect.size.width - 15 , rect.origin.y + rect.size.height-0.7),CGPointMake(rect.origin.x + rect.size.width,rect.origin.y + rect.size.height - 0.7 )};
    [self addLine:poinsBottomRightA pointB:poinsBottomRightB ctx:ctx];
    CGContextStrokePath(ctx);
}

- (void)addLine:(CGPoint[])pointA pointB:(CGPoint[])pointB ctx:(CGContextRef)ctx {
    CGContextAddLines(ctx, pointA, 2);
    CGContextAddLines(ctx, pointB, 2);
}
@end
