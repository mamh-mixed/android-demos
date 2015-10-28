package com.cardinfolink.yunshouyin.salesman.models;

public class SystemConfig {

    //dev, test, pro 是一样的
    public static final String APP_KEY = "eu1dr0c8znpa43blzy1wirzmk8jqdaon";

    //用户系统服务器地址
    // dev
//    public static final String APP_Server ="http://dev.quick.ipay.so";
    // test
//    public static final String APP_Server = "http://test.quick.ipay.so";
    // pro
    public static final String APP_Server = "https://api.shou.money";
    public static final String Server_Tool = APP_Server + "/app/tools";

    // 行号信息
    public static final String bankbase_key = "20e786206dcf4aae8a63fe34553fd274";
    public static final String bankbase_url = "http://211.144.213.120:443/bdp";

    //扫固定码支付网页订单支付
//    public static final String WEB_PAY_URL ="http://qrcode.cardinfolink.net/agent/pay1.html?merchantCode=";
//    public static final String WEB_BILL_URL="http://qrcode.cardinfolink.net/agent/trade1.html?merchantCode=";
}
