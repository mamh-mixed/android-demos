package com.cardinfolink.yunshouyin.salesman.utils;

import android.app.Activity;
import android.app.Application;
import android.content.Context;
import android.content.SharedPreferences;
import android.util.Log;

import com.cardinfolink.yunshouyin.salesman.BuildConfig;
import com.cardinfolink.yunshouyin.salesman.api.QuickPayConfigStorage;
import com.cardinfolink.yunshouyin.salesman.core.QiniuMultiUploadService;
import com.cardinfolink.yunshouyin.salesman.core.QuickPayService;
import com.cardinfolink.yunshouyin.salesman.model.User;

public class SalesmanApplication extends Application {
    private static final String TAG = "SalesmanApplication";
    private static final String ENV= BuildConfig.ENVIRONMENT;

    private static SalesmanApplication singleton;
    private Context context;
    private QuickPayConfigStorage quickPayConfigStorage;
    private QuickPayService quickPayService;
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

    public QiniuMultiUploadService getQiniuMultiUploadService() {
        return qiniuMultiUploadService;
    }

    public QuickPayConfigStorage getQuickPayConfigStorage() {
        return quickPayConfigStorage;
    }

    @Override
    public void onCreate() {
        super.onCreate();

        /**
         * initial setup at the beginning of startup
         */
        context = getApplicationContext();
        singleton = this;

        initEnvironment();
    }

    private void initEnvironment() {
        quickPayConfigStorage = new QuickPayConfigStorage();
        //dev, test, pro 是一样的
        quickPayConfigStorage.setAppKey("eu1dr0c8znpa43blzy1wirzmk8jqdaon");
        quickPayConfigStorage.setUrl("http://test.quick.ipay.so/app/tools");

        try {
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

        quickPayService = new QuickPayService(quickPayConfigStorage);
        qiniuMultiUploadService = new QiniuMultiUploadService(quickPayService);
    }


}
