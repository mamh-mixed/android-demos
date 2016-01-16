package com.cardinfolink.yunshouyin.util;

/**
 * Created by mamh on 16-1-16.
 * 日志的工具类
 */
public class Log {
    public static boolean LOG = false;

    private Log() {
    }

    public static void i(String tag, String msg) {
        if (LOG) {
            android.util.Log.i(tag, msg);
        }
    }


    public static void d(String tag, String msg) {
        if (LOG) {
            android.util.Log.d(tag, msg);
        }
    }


    public static void v(String tag, String msg) {
        if (LOG) {
            android.util.Log.v(tag, msg);
        }
    }


    public static void w(String tag, String msg) {
        if (LOG) {
            android.util.Log.w(tag, msg);
        }
    }

    public static void e(String tag, String msg) {
        if (LOG) {
            android.util.Log.e(tag, msg);
        }
    }

}
