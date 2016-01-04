package com.cardinfolink.yunshouyin.data;

public class SessonData {
    public static User loginUser = new User();

    //建一个静态对象存放卡券优惠相关信息,使用了单例模式
    public static Coupon coupon = Coupon.getInstance();
}
