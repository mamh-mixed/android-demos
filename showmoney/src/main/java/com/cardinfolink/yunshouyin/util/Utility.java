package com.cardinfolink.yunshouyin.util;

import android.graphics.Bitmap;
import android.graphics.Matrix;

import java.text.SimpleDateFormat;
import java.util.Date;
import java.util.Random;


/**
 * 工具类 方法，这里面存放一些通用的工具类方法
 */
public class Utility {
    public static Bitmap zoomBitmap(Bitmap icon, int h) {
        Matrix m = new Matrix();
        float sx = (float) 2 * h / icon.getWidth();
        float sy = (float) 2 * h / icon.getHeight();
        m.setScale(sx, sy);
        return Bitmap.createBitmap(icon, 0, 0, icon.getWidth(), icon.getHeight(), m, false);
    }

    /**
     * 生成一个新的订单号
     *
     * @return
     */
    public static String geneOrderNumber() {
        String mOrderNum;

        Date now = new Date();
        SimpleDateFormat spf = new SimpleDateFormat("yyMMddHHmmss");
        mOrderNum = spf.format(now);
        Random random = new Random();//订单号末尾随机的生成一个数
        for (int i = 0; i < 5; i++) {
            mOrderNum = mOrderNum + random.nextInt(10);
        }
        return mOrderNum;
    }

}
