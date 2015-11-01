package com.cardinfolink.yunshouyin.util;

import android.app.Application;

public class ContextUtil extends Application {
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
    }
}
