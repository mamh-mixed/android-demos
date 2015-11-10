package com.cardinfolink.yunshouyin.salesman.utils;

import android.app.Activity;
import android.app.Application;
import android.content.Context;
import android.content.SharedPreferences;
import android.content.pm.ApplicationInfo;
import android.content.pm.PackageManager;
import android.os.Bundle;
import android.util.Log;

import com.cardinfolink.yunshouyin.salesman.BuildConfig;
import com.cardinfolink.yunshouyin.salesman.api.QuickPayConfigStorage;
import com.cardinfolink.yunshouyin.salesman.core.QiniuMultiUploadService;
import com.cardinfolink.yunshouyin.salesman.core.QuickPayService;
import com.cardinfolink.yunshouyin.salesman.model.User;

public class SAApplication extends Application {
    private static final String TAG = "SAApplication";
    private static final String ENV= BuildConfig.ENVIRONMENT;

    private static SAApplication singleton;
    private Context context;
    private QuickPayConfigStorage quickPayConfigStorage;
    private QuickPayService quickPayService;
    private QiniuMultiUploadService qiniuMultiUploadService;
    private User loginUser = new User();

    public static SAApplication getInstance() {
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

    public User getLoginUser() {
        SharedPreferences mySharedPreferences = context.getSharedPreferences("savedata", Activity.MODE_PRIVATE);
        loginUser.setUsername(mySharedPreferences.getString("username", ""));
        loginUser.setPassword(mySharedPreferences.getString("password", ""));
        loginUser.setAutoLogin(mySharedPreferences.getBoolean("autologin", false));
        System.out.println(loginUser.getUsername());
        return loginUser;
    }

    public void setLoginUser(User user) {
        this.loginUser = user;
        SharedPreferences mySharedPreferences = context.getSharedPreferences("savedata", Activity.MODE_PRIVATE);
        SharedPreferences.Editor editor = mySharedPreferences.edit();
        editor.putString("username", user.getUsername());
        editor.putString("password", user.getPassword());
        editor.putBoolean("autologin", user.isAutoLogin());
        editor.commit();
    }
}
