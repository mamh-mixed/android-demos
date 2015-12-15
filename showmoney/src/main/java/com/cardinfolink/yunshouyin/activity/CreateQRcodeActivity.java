package com.cardinfolink.yunshouyin.activity;

import android.content.Intent;
import android.graphics.Bitmap;
import android.graphics.BitmapFactory;
import android.os.Bundle;
import android.os.Handler;
import android.os.Message;
import android.text.TextUtils;
import android.view.View;
import android.view.View.OnClickListener;
import android.widget.ImageView;
import android.widget.TextView;
import android.widget.Toast;

import com.cardinfolink.cashiersdk.listener.CashierListener;
import com.cardinfolink.cashiersdk.model.OrderData;
import com.cardinfolink.cashiersdk.model.ResultData;
import com.cardinfolink.cashiersdk.sdk.CashierSdk;
import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.constant.Msg;
import com.cardinfolink.yunshouyin.data.SessonData;
import com.cardinfolink.yunshouyin.util.Untilly;
import com.cardinfolink.yunshouyin.view.TradingCustomDialog;
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

public class CreateQRcodeActivity extends BaseActivity {
    private static final String TAG = "CreateQRcodeActivity";
    private static final int IMAGE_HALFWIDTH = 40;
    int FOREGROUND_COLOR = 0xff000000;
    int BACKGROUND_COLOR = 0xffffffff;
    private ImageView mQrcodeImage;
    private ResultData mResultData;
    private Handler mHandler;
    private String total;
    private String chcd;
    private TextView mPayMoneyText;
    private String mOrderNum;
    private TextView mScanText;
    private TradingCustomDialog mCustomDialog;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.create_qrcode);
        Intent intent = getIntent();
        total = intent.getStringExtra("total");
        chcd = intent.getStringExtra("chcd");

        Date now = new Date();
        SimpleDateFormat spf = new SimpleDateFormat("yyMMddHHmmss");
        mOrderNum = spf.format(now);
        Random random = new Random();
        for (int i = 0; i < 5; i++) {
            mOrderNum = mOrderNum + random.nextInt(10);
        }

        final OrderData orderData = new OrderData();
        orderData.orderNum = mOrderNum;
        orderData.txamt = total;
        orderData.currency = "156";
        orderData.chcd = chcd;

        initHandler();
        initLayout();

        initListener();

        CashierSdk.startPrePay(orderData, new CashierListener() {

            @Override
            public void onResult(ResultData resultData) {
                mResultData = resultData;
                Message msg = new Message();
                msg.what = 1;
                mHandler.sendMessageDelayed(msg, 0);
            }

            @Override
            public void onError(int errorCode) {

            }
        });
    }

    private void initLayout() {
        mQrcodeImage = (ImageView) findViewById(R.id.qrcode_img);
        mPayMoneyText = (TextView) findViewById(R.id.pay_money);
        mPayMoneyText.setText(getResources().getString(R.string.create_qrcode_activity_amount) + total);
        mScanText = (TextView) findViewById(R.id.scan_text);
        mCustomDialog = new TradingCustomDialog(mContext, mHandler,
                findViewById(R.id.trading_custom_dialog), mOrderNum);

    }

    private void initListener() {
        findViewById(R.id.qy).setOnClickListener(new OnClickListener() {

            @Override
            public void onClick(View v) {
                OrderData orderData = new OrderData();
                orderData.origOrderNum = mOrderNum;
                CashierSdk.startQy(orderData, new CashierListener() {

                    @Override
                    public void onResult(ResultData resultData) {
                        if (resultData.respcd.equals("00")) {
                            mHandler.sendEmptyMessage(Msg.MSG_FROM_SERVER_TRADE_SUCCESS);
                        } else if (resultData.respcd.equals("09")) {
                            mHandler.sendEmptyMessage(Msg.MSG_FROM_SERVER_TRADE_NOPAY);
                        } else {
                            mHandler.sendEmptyMessage(Msg.MSG_FROM_SERVER_TRADE_FAIL);
                        }
                    }

                    @Override
                    public void onError(int errorCode) {
                        mHandler.sendEmptyMessage(Msg.MSG_FROM_SERVER_TIMEOUT);
                    }

                });
            }

        });

        findViewById(R.id.back).setOnClickListener(new OnClickListener() {

            @Override
            public void onClick(View v) {
                CreateQRcodeActivity.this.finish();

            }
        });

    }

    private void updateLayout() {
        Bitmap icon = null;
        if (mResultData.chcd.equals("WXP")) {
            icon = BitmapFactory.decodeResource(getResources(), R.drawable.wpay);
            mScanText.setText(getResources().getString(R.string.create_qrcode_activity_open_wx));
        } else {
            icon = BitmapFactory.decodeResource(getResources(), R.drawable.apay);
            mScanText.setText(getResources().getString(R.string.create_qrcode_activity_open_ali));
        }
        Bitmap bitmap;
        try {
            if (!TextUtils.isEmpty(mResultData.qrcode)) {
                bitmap = cretaeBitmap(mResultData.qrcode, icon);
                mQrcodeImage.setImageBitmap(bitmap);
            } else {
                mScanText.setText(mResultData.errorDetail);
            }
        } catch (WriterException e) {
            e.printStackTrace();
        }

    }

    private void initHandler() {
        mHandler = new Handler() {
            @Override
            public void handleMessage(Message msg) {
                switch (msg.what) {
                    case 1: {

                        if (mResultData != null) {
                            if (mResultData.respcd.equals("00") || mResultData.respcd.equals("09")) {
                                updateLayout();
                            }
                        }
                        break;
                    }

                    case 2: {
                        if (mResultData != null) {
                            if (mResultData.respcd.equals("00")) {

                                CreateQRcodeActivity.this.finish();
                            } else {

                            }
                        }
                        break;
                    }
                    case 3: {
                        Toast toast = Toast.makeText(getApplicationContext(), getResources().getString(R.string.server_timeout), Toast.LENGTH_SHORT);
                        toast.show();
                    }
                    case Msg.MSG_FROM_DIGLOG_CLOSE: {
                        CreateQRcodeActivity.this.finish();
                        break;
                    }
                    case Msg.MSG_FROM_SERVER_TRADE_SUCCESS: {

                        mCustomDialog.success();
                        break;
                    }
                    case Msg.MSG_FROM_SERVER_TRADE_FAIL: {
                        mCustomDialog.fail();
                        break;
                    }
                    case Msg.MSG_FROM_SERVER_TRADE_NOPAY: {
                        mAlertDialog.show("" + getResources().getString(R.string.dialog_trade_nopay), BitmapFactory
                                .decodeResource(
                                        mContext.getResources(),
                                        R.drawable.wrong));
                        //mCustomDialog.nopay();
                        break;
                    }
                    case Msg.MSG_FROM_SUCCESS_DIGLOG_HISTORY: {
                        SessonData.positionView = 1;
                        setResult(101);
                        finish();
                        break;
                    }

                }
                super.handleMessage(msg);
            }
        };
    }


    public Bitmap cretaeBitmap(String str, Bitmap icon) throws WriterException {
        icon = Untilly.zoomBitmap(icon, IMAGE_HALFWIDTH);
        Hashtable<EncodeHintType, Object> hints = new Hashtable<EncodeHintType, Object>();
        hints.put(EncodeHintType.ERROR_CORRECTION, ErrorCorrectionLevel.H);
        hints.put(EncodeHintType.CHARACTER_SET, "utf-8");
        hints.put(EncodeHintType.MARGIN, 1);
        BitMatrix matrix = new MultiFormatWriter().encode(str,
                BarcodeFormat.QR_CODE, 300, 300, hints);
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
                    pixels[y * width + x] = icon.getPixel(x - halfW
                            + IMAGE_HALFWIDTH, y - halfH + IMAGE_HALFWIDTH);
                } else {
                    if (matrix.get(x, y)) {
                        pixels[y * width + x] = FOREGROUND_COLOR;
                    } else {
                        pixels[y * width + x] = BACKGROUND_COLOR;
                    }
                }

            }
        }
        Bitmap bitmap = Bitmap.createBitmap(width, height,
                Bitmap.Config.ARGB_8888);
        bitmap.setPixels(pixels, 0, width, 0, 0, width, height);

        return bitmap;
    }


}
