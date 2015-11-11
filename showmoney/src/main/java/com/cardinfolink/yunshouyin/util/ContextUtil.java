package com.cardinfolink.yunshouyin.util;

import android.app.Application;

import com.cardinfolink.yunshouyin.BuildConfig;
import com.cardinfolink.yunshouyin.constant.SystemConfig;

public class ContextUtil extends Application {
    private static final String ENVIRONMENT = BuildConfig.ENVIRONMENT;

    private static ContextUtil instance;

    public static ContextUtil getInstance() {
        return instance;
    }

    public static String getResString(int id) {
        return instance.getResources().getString(id);
    }

    @Override
    public void onCreate() {
        // TODO Auto-generated method stub
        super.onCreate();
        instance = this;
        CrashHandler handler = CrashHandler.getInstance();
        handler.init(getApplicationContext());

        initEnvironment();
    }


    private void initEnvironment() {
        //default is pro
        SystemConfig.IS_PRODUCE =true;
        SystemConfig.Server="https://api.shou.money/app";
        SystemConfig.WEB_BILL_URL="http://qrcode.cardinfolink.net/payment/trade.html";
        switch (ENVIRONMENT){
            case "pro":
                SystemConfig.IS_PRODUCE =true;
                SystemConfig.Server="https://api.shou.money/app";
                SystemConfig.WEB_BILL_URL="http://qrcode.cardinfolink.net/payment/trade.html";
                break;
            case "test":
                SystemConfig.IS_PRODUCE =false;
                SystemConfig.Server="http://test.quick.ipay.so/app";
                SystemConfig.WEB_BILL_URL="http://qrcode.cardinfolink.net/agent/trade.html";
                break;
            case "dev":
                SystemConfig.IS_PRODUCE =false;
                SystemConfig.Server="http://dev.quick.ipay.so/app";
                SystemConfig.WEB_BILL_URL="http://qrcode.cardinfolink.net/agent/trade.html";
                break;
            default:
                break;
        }
    }
}
