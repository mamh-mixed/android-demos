package com.cardinfolink.yunshouyin.salesman.utils;

import android.app.Activity;

import com.cardinfolink.yunshouyin.salesman.activity.LoginActivity;
import com.cardinfolink.yunshouyin.salesman.activity.MerchantListActivity;

import java.util.ArrayList;
import java.util.List;

public class ActivityCollector {
    public static List<Activity> activityList = new ArrayList<>();

    public static void addActivity(Activity activity) {
        activityList.add(activity);
    }

    public static void removeActivity(Activity activity) {
        activityList.remove(activity);
    }

    /**
     * Exit the App
     */
    public static void finishAll() {
        for (Activity activity : activityList) {
            if (!activity.isFinishing()) {
                activity.finish();
            }
        }
    }

    /**
     * Used if accessToken expired
     */
    public static void goLoginAndFinishRest() {
        for (Activity activity : activityList) {
            if (activity instanceof LoginActivity) {
                continue;
            }
            if (!activity.isFinishing()) {
                activity.finish();
            }
        }
    }

    /**
     * Used if registration finished or cancelled
     */
    public static void goHomeAndFinishRest() {
        for (Activity activity : activityList) {
            if (activity instanceof LoginActivity) {
                continue;
            }
            if (activity instanceof MerchantListActivity) {
                continue;
            }
            if (!activity.isFinishing()) {
                activity.finish();
            }
        }
    }
}
