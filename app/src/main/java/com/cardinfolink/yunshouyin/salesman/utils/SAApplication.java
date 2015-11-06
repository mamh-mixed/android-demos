package com.cardinfolink.yunshouyin.salesman.utils;

import android.app.Application;
import android.content.Context;
import android.content.pm.ApplicationInfo;
import android.content.pm.PackageManager;
import android.os.Bundle;
import android.util.Log;

import com.cardinfolink.yunshouyin.salesman.api.QuickPayConfigStorage;
import com.cardinfolink.yunshouyin.salesman.core.QiniuMultiUploadService;
import com.cardinfolink.yunshouyin.salesman.core.QuickPayService;
import com.cardinfolink.yunshouyin.salesman.model.SystemConfig;

public class SAApplication extends Application {
    private static final String TAG = "SAApplication";
    private static SAApplication singleton;

    public static SAApplication getInstance(){
        return singleton;
    }

    private Context context;
    private QuickPayConfigStorage quickPayConfigStorage;
    private QuickPayService quickPayService;
    private QiniuMultiUploadService qiniuMultiUploadService;

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
            ApplicationInfo ai = getPackageManager().getApplicationInfo(
                    getPackageName(), PackageManager.GET_META_DATA);
            Bundle bundle = ai.metaData;
            String environment = bundle.getString("ENVIRONMENT");
            Log.d(TAG, "ENVIRONMENT is " + environment);
            switch (environment) {
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
        }
        catch (Exception e) {
            Log.e(TAG, "Failed to load meta-data: " + e.getMessage());
        }

        quickPayService = new QuickPayService(quickPayConfigStorage);
        qiniuMultiUploadService = new QiniuMultiUploadService(quickPayService);
    }
}
