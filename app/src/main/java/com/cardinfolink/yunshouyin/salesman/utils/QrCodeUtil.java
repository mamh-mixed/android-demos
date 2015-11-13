package com.cardinfolink.yunshouyin.salesman.utils;

import android.graphics.Bitmap;
import android.graphics.BitmapFactory;
import android.graphics.Canvas;
import android.graphics.Color;
import android.graphics.Paint;
import android.graphics.Rect;
import android.graphics.Typeface;

import com.cardinfolink.yunshouyin.salesman.R;
import com.google.zxing.BarcodeFormat;
import com.google.zxing.EncodeHintType;
import com.google.zxing.WriterException;
import com.google.zxing.common.BitMatrix;
import com.google.zxing.qrcode.QRCodeWriter;
import com.google.zxing.qrcode.decoder.ErrorCorrectionLevel;

import java.util.HashMap;
import java.util.Map;

public class QrCodeUtil {
    public static Bitmap getQRPostA(String qrContent, String title) {
        return getQRPostCommon(R.drawable.template_a, qrContent, title);
    }

    public static Bitmap getQRPostB(String qrContent, String title) {
        return getQRPostCommon(R.drawable.template_b, qrContent, title);
    }

    /**
     * 这里是根据UI给的标注图设定的,TODO:缺少title的位置
     *
     * @param templateResId
     * @param qrContent
     * @param title
     * @return
     */
    private static Bitmap getQRPostCommon(int templateResId, String qrContent, String title) {
        if (title == null || title.equals("")) {
            title = "没有名字";
        }
        BitmapFactory.Options options = new BitmapFactory.Options();
        options.inScaled = false;
        // mutable template
        Bitmap template = BitmapFactory.decodeResource(SalesmanApplication.getInstance().getContext().getResources(), templateResId, options).copy(Bitmap.Config.ARGB_8888, true);
        int screenWidth = template.getWidth();
        int screenHeight = template.getHeight();
        int qrWidth = 90 * screenWidth / 148;
        int logoWidth = 20 * screenWidth / 148;
        int qrPosX = 28 * screenWidth / 148;
        int qrPosY = 28 * screenWidth / 148;

        template = drawQRCode(template, qrContent, qrPosX, qrPosY, qrWidth, logoWidth);
        //font size

        int txtCenterX = screenWidth / 2;
        int txtCenterY = screenHeight - 80;
        template = drawText(template, title, txtCenterX, txtCenterY);
        return template;
    }


    /**
     * draw text
     *
     * @param template
     * @param title
     * @param txtCenterX
     * @param txtCenterY
     * @return
     */
    private static Bitmap drawText(Bitmap template, String title, int txtCenterX, int txtCenterY) {
        Canvas canvas = new Canvas(template);

        Paint paint = new Paint();
        paint.setColor(Color.WHITE);
        paint.setTextSize(50);

        Typeface tf = Typeface.createFromAsset(SalesmanApplication.getInstance().getContext().getAssets(), "fonts/simhei.ttf");
        paint.setTypeface(tf);

        Rect bounds = new Rect();
        paint.getTextBounds(title, 0, title.length(), bounds);
        int x = txtCenterX - bounds.width() / 2;
        int y = txtCenterY - bounds.height() / 2;
        canvas.drawText(title, x, y, paint);

        return template;
    }

    /**
     * logo预先已经生成好了,所以QRCode不要覆盖掉logo的位置
     * logo居中于QRCode
     *
     * @param template
     * @param qrContent
     * @param qrWidth
     * @param logoWidth
     * @return
     */
    private static Bitmap drawQRCode(Bitmap template, String qrContent, int qrPosX, int qrPosY, int qrWidth, int logoWidth) {
        Canvas canvas = new Canvas(template);

        try {
            Bitmap qrCode = getQRBitmap(qrContent, qrWidth, logoWidth);
            //NOTE: 必须保持两图density一致,不然android会自动缩放
            qrCode.setDensity(template.getDensity());
            canvas.drawBitmap(qrCode, qrPosX, qrPosY, null);

        } catch (WriterException e) {
            //TODO: handle exception
            e.printStackTrace();
        }

        canvas.save(Canvas.ALL_SAVE_FLAG);
        canvas.restore();
        return template;
    }


    private static Bitmap getQRBitmap(String qrContent, int qrWidth, int logoWidth) throws WriterException {
        QRCodeWriter writer = new QRCodeWriter();
        String charset = "UTF-8";
        Map hintMap = new HashMap<>();
        hintMap.put(EncodeHintType.ERROR_CORRECTION, ErrorCorrectionLevel.H);
        hintMap.put(EncodeHintType.CHARACTER_SET, "utf-8");
        hintMap.put(EncodeHintType.MARGIN, 1);
        BitMatrix bitMatrix = writer.encode(qrContent, BarcodeFormat.QR_CODE, qrWidth, qrWidth, hintMap);
        Bitmap result = zxing2Bitmap(bitMatrix, logoWidth);
        Canvas canvas = new Canvas(result);
        canvas.save(Canvas.ALL_SAVE_FLAG);
        canvas.restore();
        return result;
    }

    // Ref: http://codeisland.org/2013/generating-qr-codes-with-zxing/
    private static Bitmap zxing2Bitmap(BitMatrix matrix, int logoWidth) {
        int height = matrix.getHeight();
        int width = matrix.getWidth();
        BitmapFactory.Options options = new BitmapFactory.Options();
        options.inScaled = false;
        //320
        Bitmap bmp = Bitmap.createBitmap(width, height, Bitmap.Config.ARGB_8888);
        for (int x = 0; x < width; x++) {
            for (int y = 0; y < height; y++) {
                /**
                 * 中间logo位置不画
                 */
                if (x >= (width - logoWidth) / 2 && x <= (width + logoWidth) / 2 &&
                        y >= (height - logoWidth) / 2 && y <= (height + logoWidth) / 2) {
                    continue;
                }
                bmp.setPixel(x, y, matrix.get(x, y) ? Color.BLACK : Color.WHITE);
            }
        }
        return bmp;
    }
}
