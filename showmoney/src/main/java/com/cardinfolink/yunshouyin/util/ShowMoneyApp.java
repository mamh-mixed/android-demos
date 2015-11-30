package com.cardinfolink.yunshouyin.util;

import android.app.Application;

import com.cardinfolink.yunshouyin.BuildConfig;
import com.cardinfolink.yunshouyin.api.QuickPayConfigStorage;
import com.cardinfolink.yunshouyin.constant.SystemConfig;
import com.cardinfolink.yunshouyin.core.BankDataService;
import com.cardinfolink.yunshouyin.core.BankDataServiceImpl;
import com.cardinfolink.yunshouyin.core.QuickPayService;
import com.cardinfolink.yunshouyin.core.QuickPayServiceImpl;

public class ShowMoneyApp extends Application {
    private static final String ENVIRONMENT = BuildConfig.ENVIRONMENT;

    private static ShowMoneyApp instance;

    private QuickPayConfigStorage quickPayConfigStorage;
    private QuickPayService quickPayService;
    private BankDataService bankDataService;

    public static ShowMoneyApp getInstance() {
        return instance;
    }

    public static String getResString(int id) {
        return instance.getResources().getString(id);
    }


    public QuickPayService getQuickPayService() {
        return quickPayService;
    }

    public BankDataService getBankDataService() {
        return bankDataService;
    }

    @Override
    public void onCreate() {
        // TODO Auto-generated method stub
        super.onCreate();
        instance = this;
        initEnvironment();
    }


    private void initEnvironment() {
        quickPayConfigStorage = new QuickPayConfigStorage();
        //dev, test, pro 是一样的
        quickPayConfigStorage.setAppKey("eu1dr0c8znpa43blzy1wirzmk8jqdaon");

        //default is pro
        SystemConfig.IS_PRODUCE = true;
        SystemConfig.Server = "https://api.shou.money/app";
        SystemConfig.WEB_BILL_URL = "http://qrcode.cardinfolink.net/payment/trade.html";
        switch (ENVIRONMENT) {
            case "pro":
                SystemConfig.IS_PRODUCE = true;
                SystemConfig.Server = "https://api.shou.money/app";
                SystemConfig.WEB_BILL_URL = "http://qrcode.cardinfolink.net/payment/trade.html";
                break;
            case "test":
                SystemConfig.IS_PRODUCE = false;
                SystemConfig.Server = "http://test.quick.ipay.so/app";
                SystemConfig.WEB_BILL_URL = "http://qrcode.cardinfolink.net/agent/trade.html";
                break;
            case "dev":
                SystemConfig.IS_PRODUCE = false;
                SystemConfig.Server = "http://dev.quick.ipay.so/app";
                SystemConfig.WEB_BILL_URL = "http://qrcode.cardinfolink.net/agent/trade.html";
                break;
            default:
                break;
        }

        quickPayConfigStorage.setUrl(SystemConfig.Server);

        quickPayConfigStorage.setBankbaseKey(SystemConfig.bankbase_key);
        quickPayConfigStorage.setBankbaseUrl(SystemConfig.bankbase_url);

        quickPayService = new QuickPayServiceImpl(quickPayConfigStorage);
        bankDataService = new BankDataServiceImpl(quickPayConfigStorage);



    }
}
