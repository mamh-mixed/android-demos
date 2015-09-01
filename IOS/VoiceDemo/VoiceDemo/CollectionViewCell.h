//
//  CollectionViewCell.h
//  VTCheckstandDemo
//
//  Created by 司瑞华 on 15/3/12.
//  Copyright (c) 2015年 __VTPayment__. All rights reserved.
//

#import <UIKit/UIKit.h>
#import "ButtonPro.h"

@interface CollectionViewCell : UICollectionViewCell

@property (nonatomic,strong)UIImageView         * imgView;
@property (nonatomic,strong)UILabel             * titleLbl;
@property (nonatomic,strong)ButtonPro           * buyBtn;
@end
