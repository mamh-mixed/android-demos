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

    //公共数据平台秘钥和key

    public static final String bankbase_key = "20e786206dcf4aae8a63fe34553fd274";
    public static final String bankbase_url = "http://211.144.213.120:443/bdp";

}
