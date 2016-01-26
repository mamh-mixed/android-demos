package com.cardinfolink.yunshouyin.util;

import android.app.Activity;

import com.cardinfolink.yunshouyin.activity.LoginActivity;

import java.util.ArrayList;
import java.util.List;

/**
 * Created by mamh on 16-1-26.
 */
public class ActivityCollector {


    public static List<Activity> activityList = new ArrayList<>();

    public static void addActivity(Activity activity) {
        activityList.add(activity);
    }

    public static void removeActivity(Activity activity) {
        activityList.remove(activity);
    }

    public static void finishAll() {
        for (Activity activity : activityList) {
            if (!activity.isFinishing()) {
                activity.finish();
            }
        }
    }

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
}
