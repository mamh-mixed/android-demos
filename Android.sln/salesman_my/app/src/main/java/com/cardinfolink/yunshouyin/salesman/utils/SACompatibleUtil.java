package com.cardinfolink.yunshouyin.salesman.utils;

import android.os.Environment;

import java.io.File;

public class SACompatibleUtil {

    // bug workaround: https://code.google.com/p/android/issues/detail?id=75447
    public static void fixMediaDir() {
        File sdcard = Environment.getExternalStorageDirectory();
        if (sdcard == null) {
            return;
        }
        File dcim = new File(sdcard, "DCIM");
        if (dcim == null) {
            return;
        }
        File camera = new File(dcim, "Camera");
        if (camera.exists()) {
            return;
        }
        camera.mkdir();
    }

}
