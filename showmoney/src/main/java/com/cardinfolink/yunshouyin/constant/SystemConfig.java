package com.cardinfolink.yunshouyin.constant;

public class SystemConfig {
    public static final String APP_KEY = "eu1dr0c8znpa43blzy1wirzmk8jqdaon";

    /**
     * 生产环境的 host 地址
     */
    public static final String PRO_SERVER = "https://api.shou.money/app";

    /**
     * 测试环境的地址
     */
    public static final String TEST_SERVER = "http://test.quick.ipay.so/app";

    /**
     * 开发环境的地址
     */
    public static final String DEV_SERVER = "http://dev.quick.ipay.so/app";

    //公共数据平台秘钥和key
    public static final String BANKBASE_KEY = "20e786206dcf4aae8a63fe34553fd274";
    public static final String BANKBASE_URL = "http://211.144.213.120:443/bdp";

    /**
     * 用户系统服务器地址
     * 这个之后做个可以在界面中 更改的操作,所以这里不能是final类型
     */
    public static String SERVER;

    /**
     * SDK 环境
     * 这个之后做个可以在界面中 更改的操作,所以这里不能是final类型
     */
    public static boolean IS_PRODUCE;


    /**
     * 是否 是 debug 模式，这个可以作为log的开关，debug为true的时候打印log，为false的时候不打印，
     * 这个也可以在界面中去设置
     */
    public static boolean DEBUG = false;
}
