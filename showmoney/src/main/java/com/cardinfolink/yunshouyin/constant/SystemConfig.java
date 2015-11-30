package com.cardinfolink.yunshouyin.constant;

import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.util.ShowMoneyApp;

public class SystemConfig {


    public static String APP_KEY = ShowMoneyApp.getInstance().getResources().getString(R.string.app_key);//app用户系统交互key

    //用户系统服务器地址
    public static String Server;

    //扫固定码支付网页订单支付
    public static String WEB_BILL_URL;

    // SDK 环境
    public static boolean IS_PRODUCE;

}
