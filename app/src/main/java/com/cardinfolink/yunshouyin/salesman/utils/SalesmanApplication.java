package com.cardinfolink.yunshouyin.salesman.utils;

import android.app.Activity;
import android.app.Application;
import android.content.Context;
import android.content.SharedPreferences;
import android.util.Log;

import com.cardinfolink.yunshouyin.salesman.BuildConfig;
import com.cardinfolink.yunshouyin.salesman.api.QuickPayConfigStorage;
import com.cardinfolink.yunshouyin.salesman.core.BankDataService;
import com.cardinfolink.yunshouyin.salesman.core.BankDataServiceImpl;
import com.cardinfolink.yunshouyin.salesman.core.QiniuMultiUploadService;
import com.cardinfolink.yunshouyin.salesman.core.QuickPayService;
import com.cardinfolink.yunshouyin.salesman.core.QuickPayServiceImpl;
import com.cardinfolink.yunshouyin.salesman.model.SystemConfig;
import com.cardinfolink.yunshouyin.salesman.model.User;

public class SalesmanApplication extends Application {
    private static final String TAG = "SalesmanApplication";
    private static final String ENV = BuildConfig.ENVIRONMENT;

    private static SalesmanApplication singleton;
    private Context context;

    private QuickPayConfigStorage quickPayConfigStorage;

    private QuickPayService quickPayService;
    private BankDataServiceImpl bankDataService;

    private QiniuMultiUploadService qiniuMultiUploadService;


    public static SalesmanApplication getInstance() {
        return singleton;
    }

    public Context getContext() {
        return context;
    }

    public QuickPayService getQuickPayService() {
        return quickPayService;
    }

    public BankDataService getBankDataService() {
        return bankDataService;
    }

    public QiniuMultiUploadService getQiniuMultiUploadService() {
        return qiniuMultiUploadService;
    }

    public QuickPayConfigStorage getQuickPayConfigStorage() {
        return quickPayConfigStorage;
    }

    @Override
    public void onCreate() {
        super.onCreate();
        context = getApplicationContext();
        singleton = this;
        initEnvironment();
    }

    private void initEnvironment() {
        quickPayConfigStorage = new QuickPayConfigStorage();
        //dev, test, pro 是一样的
        quickPayConfigStorage.setAppKey(SystemConfig.APP_KEY);//app用户系统交互key
        quickPayConfigStorage.setUrl(SystemConfig.URL);//默认设置未测试的url地址。
        //设置bank 需要的key和url
        quickPayConfigStorage.setBankbaseKey(SystemConfig.BANKBASE_KEY);
        quickPayConfigStorage.setBankbaseUrl(SystemConfig.BANKBASE_URL);
        try {
            //这里根据environment值的不同选择不同的url地址
            Log.d(TAG, "ENVIRONMENT is " + ENV);
            switch (ENV) {
                case "dev":
                    quickPayConfigStorage.setUrl("http://dev.quick.ipay.so/app/tools");
                    break;
                case "test":
                    quickPayConfigStorage.setUrl("http://test.quick.ipay.so/app/tools");
                    break;
                case "pro":
                    quickPayConfigStorage.setUrl("https://api.shou.money/app/tools");
                    break;
                default:
                    break;
            }
        } catch (Exception e) {
            Log.e(TAG, "Failed to load meta-data: " + e.getMessage());
        }

        quickPayService = new QuickPayServiceImpl(quickPayConfigStorage);
        bankDataService = new BankDataServiceImpl(quickPayConfigStorage);
        qiniuMultiUploadService = new QiniuMultiUploadService(quickPayService);
    }


}
