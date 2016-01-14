package com.cardinfolink.yunshouyin.util;

import android.content.Context;
import android.graphics.Bitmap;
import android.graphics.Matrix;
import android.view.inputmethod.InputMethodManager;

import com.google.zxing.BarcodeFormat;
import com.google.zxing.EncodeHintType;
import com.google.zxing.MultiFormatWriter;
import com.google.zxing.WriterException;
import com.google.zxing.common.BitMatrix;
import com.google.zxing.qrcode.decoder.ErrorCorrectionLevel;

import java.text.SimpleDateFormat;
import java.util.Date;
import java.util.Hashtable;
import java.util.Random;


/**
 * 工具类 方法，这里面存放一些通用的工具类方法
 */
public class Utility {
    private static final int IMAGE_HALFWIDTH = 30;
    private static final int FOREGROUND_COLOR = 0xff000000;
    private static final int BACKGROUND_COLOR = 0xffffffff;


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

    public static Bitmap cretaeBitmap(String str, int widthx, int heighty) throws WriterException {
        return cretaeBitmap(str, null, widthx, heighty);
    }

    /**
     * //生成bitmap 二维码图片,生成一个固定长宽都是width和height的二维码图片
     *
     * @param str
     * @param icon
     * @param widthx
     * @param heighty
     * @return
     * @throws WriterException
     */
    public static Bitmap cretaeBitmap(String str, Bitmap icon, int widthx, int heighty) throws WriterException {
        if (icon != null) {
            icon = zoomBitmap(icon, IMAGE_HALFWIDTH);
        }

        Hashtable<EncodeHintType, Object> hints = new Hashtable<EncodeHintType, Object>();
        hints.put(EncodeHintType.ERROR_CORRECTION, ErrorCorrectionLevel.H);
        hints.put(EncodeHintType.CHARACTER_SET, "utf-8");
        hints.put(EncodeHintType.MARGIN, 1);
        //调用com.google.zxing里面的生成二维码的方法
        BitMatrix matrix = new MultiFormatWriter().encode(str, BarcodeFormat.QR_CODE, widthx, heighty, hints);

        int width = matrix.getWidth();
        int height = matrix.getHeight();

        int halfW = width / 2;
        int halfH = height / 2;
        int[] pixels = new int[width * height];
        for (int y = 0; y < height; y++) {
            for (int x = 0; x < width; x++) {
                if (icon != null && x > halfW - IMAGE_HALFWIDTH && x < halfW + IMAGE_HALFWIDTH
                        && y > halfH - IMAGE_HALFWIDTH
                        && y < halfH + IMAGE_HALFWIDTH) {
                    pixels[y * width + x] = icon.getPixel(x - halfW + IMAGE_HALFWIDTH, y - halfH + IMAGE_HALFWIDTH);
                } else {
                    if (matrix.get(x, y)) {
                        pixels[y * width + x] = FOREGROUND_COLOR;
                    } else {
                        pixels[y * width + x] = BACKGROUND_COLOR;
                    }
                }

            }
        }
        Bitmap bitmap = Bitmap.createBitmap(width, height, Bitmap.Config.ARGB_8888);
        bitmap.setPixels(pixels, 0, width, 0, 0, width, height);

        return bitmap;
    }


    public static void hideInput(Context context) {
        try {
            InputMethodManager inputMethodManager = (InputMethodManager) context.getSystemService(Context.INPUT_METHOD_SERVICE);
            inputMethodManager.toggleSoftInput(0, InputMethodManager.HIDE_NOT_ALWAYS);
        } catch (Exception e) {
            e.printStackTrace();
        }

    }

}
