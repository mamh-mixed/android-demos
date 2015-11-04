package com.cardinfolink.yunshouyin.salesman.models;

import android.util.Log;

public class SystemConfig {

    private static final String TAG = "SystemConfig";

    //dev, test, pro 是一样的
    public static final String APP_KEY = "eu1dr0c8znpa43blzy1wirzmk8jqdaon";

    //用户系统服务器地址,默认是test
    private static String APP_Server = "http://test.quick.ipay.so";
    private static String Server_Tool = APP_Server + "/app/tools";

    public static String getServer_Tool() {
        return Server_Tool;
    }

    // 行号信息
    public static final String bankbase_key = "20e786206dcf4aae8a63fe34553fd274";
    public static final String bankbase_url = "http://211.144.213.120:443/bdp";

    public static void initEnvironment(String env) {
        switch (env) {
            case "dev":
                APP_Server = "http://dev.quick.ipay.so";

                break;
            case "test":
                APP_Server = "http://test.quick.ipay.so";
                break;
            case "pro":
                APP_Server = "https://api.shou.money";
                break;
            default:
                break;
        }
        Server_Tool = APP_Server + "/app/tools";
        Log.d(TAG, "APP_Server is " + APP_Server);
        Log.d(TAG, "Server_Tool is " + Server_Tool);
    }
}
