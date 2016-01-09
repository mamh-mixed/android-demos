package com.cardinfolink.yunshouyin.activity;

import android.content.ContentValues;
import android.content.Context;
import android.graphics.Bitmap;
import android.graphics.BitmapFactory;
import android.os.Bundle;
import android.os.Environment;
import android.provider.MediaStore;
import android.text.TextUtils;
import android.view.View;
import android.widget.Button;
import android.widget.ImageView;

import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.data.SessonData;
import com.cardinfolink.yunshouyin.ui.SettingActionBarItem;
import com.cardinfolink.yunshouyin.util.ShowMoneyApp;
import com.cardinfolink.yunshouyin.util.Utility;
import com.google.zxing.BarcodeFormat;
import com.google.zxing.EncodeHintType;
import com.google.zxing.MultiFormatWriter;
import com.google.zxing.WriterException;
import com.google.zxing.common.BitMatrix;
import com.google.zxing.qrcode.decoder.ErrorCorrectionLevel;

import java.io.File;
import java.io.FileOutputStream;
import java.util.Hashtable;
import java.util.UUID;


/**
 * 我的网页版，这里显示一个二维码图片。通过payUrl来显示二维码图片。
 */
public class MyWebActivity extends BaseActivity {
    private static final String TAG = "MyWebActivity";
    private static final int IMAGE_HALFWIDTH = 40;
    private static final int FOREGROUND_COLOR = 0xff000000;
    private static final int BACKGROUND_COLOR = 0xffffffff;


    private SettingActionBarItem mActionBar;
    private ImageView mQRCodeImage;
    private Button mSaveQR;
    private Bitmap mQRBitmap;

    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_my_web);

        mActionBar = (SettingActionBarItem) findViewById(R.id.action_bar);
        mActionBar.setLeftTextOnclickListner(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                finish();
            }
        });


        mQRCodeImage = (ImageView) findViewById(R.id.iv_qrcode);
        Bitmap icon = icon = BitmapFactory.decodeResource(getResources(), R.drawable.scan_wechat);
        mQRBitmap = null;
        try {
            String qrcode = SessonData.loginUser.getPayUrl();
            if (!TextUtils.isEmpty(qrcode)) {
                mQRBitmap = cretaeBitmap(qrcode, icon, 300, 300);
            }
        } catch (WriterException e) {
            e.printStackTrace();
        }
        mQRCodeImage.setImageBitmap(mQRBitmap);
        mSaveQR = (Button) findViewById(R.id.btnsave);
        mSaveQR.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                if (mQRBitmap != null) {
                    saveImageToExternalStorage(mQRBitmap);
                }
            }
        });

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
    private Bitmap cretaeBitmap(String str, Bitmap icon, int widthx, int heighty) throws WriterException {
        icon = Utility.zoomBitmap(icon, IMAGE_HALFWIDTH);
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
                if (x > halfW - IMAGE_HALFWIDTH && x < halfW + IMAGE_HALFWIDTH
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


    /**
     * 存在外部存储里,图库里是不显示的
     */
    public void saveImageToExternalStorage(Bitmap bitmap) {
        //get app's album folder
        File albumDir = getAlbumStorageDir(this, getString(R.string.app_name));

        //generate file name
        String filename = UUID.randomUUID().toString().replace("-", "") + ".jpg";
        File file = new File(albumDir, filename);

        //bitmap to png
        try {
            FileOutputStream outputStream = new FileOutputStream(file);
            bitmap.compress(Bitmap.CompressFormat.JPEG, 100, outputStream);
            outputStream.flush();
            outputStream.close();

            //show in gallery
            saveImageToSystemGallery(file.getAbsolutePath());
        } catch (Exception ex) {
            ex.printStackTrace();
        }
    }

    public File getAlbumStorageDir(Context context, String albumName) {
        // Get the directory for the app's private pictures directory.
        File file = new File(context.getExternalFilesDir(Environment.DIRECTORY_PICTURES), albumName);
        if (!file.mkdirs()) {

        }
        return file;
    }

    /**
     * 所有应用都是直接取gallery里的照片的
     * meta data setup for gallery
     * ref: http://stackoverflow.com/questions/20859584/how-save-image-in-android-gallery
     */
    private void saveImageToSystemGallery(final String filePath) {
        ContentValues values = new ContentValues();

        values.put(MediaStore.Images.Media.DATE_TAKEN, System.currentTimeMillis());
        values.put(MediaStore.Images.Media.MIME_TYPE, "image/jpeg");
        values.put(MediaStore.MediaColumns.DATA, filePath);

        ShowMoneyApp.getInstance().getContentResolver().insert(MediaStore.Images.Media.EXTERNAL_CONTENT_URI, values);
    }
}
