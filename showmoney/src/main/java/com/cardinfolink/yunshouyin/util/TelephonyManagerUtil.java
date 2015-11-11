package com.cardinfolink.yunshouyin.util;

import android.content.Context;
import android.telephony.TelephonyManager;

public class TelephonyManagerUtil {
    public static String getDeviceId(Context context) {
        TelephonyManager tm = (TelephonyManager) context.getSystemService(Context.TELEPHONY_SERVICE);
        return tm.getDeviceId();
    }
}
