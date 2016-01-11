package com.cardinfolink.yunshouyin.activity;

import android.content.Intent;
import android.graphics.Bitmap;
import android.graphics.BitmapFactory;
import android.net.Uri;
import android.os.Bundle;
import android.provider.MediaStore;
import android.text.TextUtils;
import android.view.View;
import android.widget.Button;
import android.widget.ImageView;
import android.widget.Toast;

import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.data.SessonData;
import com.cardinfolink.yunshouyin.ui.SettingActionBarItem;
import com.cardinfolink.yunshouyin.util.EncoderUtil;
import com.cardinfolink.yunshouyin.util.Utility;
import com.google.zxing.WriterException;

import java.io.File;
import java.io.FileNotFoundException;
import java.io.FileOutputStream;


/**
 * 我的网页版，这里显示一个二维码图片。通过payUrl来显示二维码图片。
 */
public class MyWebActivity extends BaseActivity {
    private static final String TAG = "MyWebActivity";

    //二维码图片的宽高
    private static final int QR_WIDTH = 300;
    private static final int QR_HEIGHT = 300;

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
        mQRBitmap = null;
        try {
            String qrcode = SessonData.loginUser.getPayUrl();
            if (!TextUtils.isEmpty(qrcode)) {
                Bitmap icon = BitmapFactory.decodeResource(getResources(), R.drawable.scan_wechat);
                mQRBitmap = Utility.cretaeBitmap(qrcode, icon, QR_WIDTH, QR_HEIGHT);
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
     * 存在外部存储里,图库里是不显示的
     */
    public void saveImageToExternalStorage(Bitmap bitmap) {
        //get app's album folder
        File cacheDir = mContext.getExternalCacheDir();
        //如果外部的不能用，就调用内部的
        if (cacheDir == null) {
            cacheDir = mContext.getCacheDir();
        }
        if (!cacheDir.exists()) {
            cacheDir.mkdirs();
        }

        //generate file name
        String payUrl = SessonData.loginUser.getPayUrl();
        String filename = "myweb.jpg";
        if (!TextUtils.isEmpty(payUrl)) {
            filename = EncoderUtil.Encrypt(payUrl, "MD5");
        }
        File file = new File(cacheDir, filename + "_myweb.jpg");

        //bitmap to png
        try {
            FileOutputStream outputStream = new FileOutputStream(file);
            bitmap.compress(Bitmap.CompressFormat.JPEG, 100, outputStream);
            outputStream.flush();
            outputStream.close();

            //show in gallery
            saveImageToSystemGallery(file);
        } catch (Exception ex) {
            ex.printStackTrace();
        }
    }


    private void saveImageToSystemGallery(File file) {
        try {
            MediaStore.Images.Media.insertImage(mContext.getContentResolver(), file.getAbsolutePath(), file.getName(), null);
            mContext.sendBroadcast(new Intent(Intent.ACTION_MEDIA_SCANNER_SCAN_FILE, Uri.parse("file://" + file.getPath())));
        } catch (FileNotFoundException e) {
            e.printStackTrace();
        }

        String toastMsg = getString(R.string.my_web_activity_save_success);
        Toast.makeText(mContext, toastMsg, Toast.LENGTH_SHORT).show();
    }
}
