package com.cardinfolink.yunshouyin.activity;

import android.content.Context;
import android.content.Intent;
import android.os.Bundle;
import android.app.Activity;
import android.os.Handler;
import android.os.Message;
import android.text.TextUtils;
import android.view.View;
import android.widget.Button;
import android.widget.EditText;
import android.widget.Toast;

import com.cardinfolink.cashiersdk.listener.CashierListener;
import com.cardinfolink.cashiersdk.model.OrderData;
import com.cardinfolink.cashiersdk.model.ResultData;
import com.cardinfolink.cashiersdk.sdk.CashierSdk;
import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.constant.Msg;
import com.cardinfolink.yunshouyin.data.SessonData;
import com.cardinfolink.yunshouyin.view.TradingCustomDialog;

import java.text.SimpleDateFormat;
import java.util.Date;
import java.util.Random;

public class ScanCodeGunActivity extends Activity {

    private EditText mInputScan;
    private Intent intent;
    private String chcd;
    private String total;
    private String mOrderNum;
    private ResultData mResultData;
    private Handler mHandler;
    private TradingCustomDialog mCustomDialog;
    private Button mPay;
    private String scancode;
    private Context mContext;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_scan_code_gun);
        mContext = this;
        mInputScan = (EditText) findViewById(R.id.et_input_box);
        intent = getIntent();
        chcd = intent.getStringExtra("chcd");
        total = intent.getStringExtra("total");
        mPay = (Button) findViewById(R.id.bt_pay);
        Date now = new Date();
        SimpleDateFormat spf = new SimpleDateFormat("yyMMddHHmmss");
        mOrderNum = spf.format(now);
        Random random = new Random();
        for (int i = 0; i < 5; i++) {
            mOrderNum = mOrderNum + random.nextInt(10);
        }
        initHandler();
        mCustomDialog = new TradingCustomDialog(mContext, mHandler,
                findViewById(R.id.trading_custom_dialog), mOrderNum);
        mPay.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                scancode = mInputScan.getText().toString();
                scancode = scancode.replace("\n", "");
                if (TextUtils.isEmpty(scancode)) {
                    Toast.makeText(ScanCodeGunActivity.this, "QR不能为空", Toast.LENGTH_SHORT).show();
                } else {
                    Message msg = Message.obtain();
                    msg.what = Msg.MSG_FROM_SCANCODE_SUCCESS;
                    msg.obj = scancode;
                    mHandler.sendMessage(msg);
                }

            }
        });
    }

    public void initHandler() {
        mHandler = new Handler() {
            @Override
            public void handleMessage(Message msg) {
                switch (msg.what) {
                    case Msg.MSG_FROM_SCANCODE_SUCCESS: {
                        mCustomDialog.loading();
                        final OrderData orderData = new OrderData();
                        orderData.orderNum = mOrderNum;
                        orderData.txamt = total;
                        orderData.currency = CashierSdk.SDK_CURRENCY;
                        orderData.chcd = chcd;
                        orderData.scanCodeId = (String) msg.obj;
                        // /orderData.scanCodeId="13241252555";
                        CashierSdk.startPay(orderData, new CashierListener() {

                            @Override
                            public void onResult(ResultData resultData) {


                                mResultData = resultData;
                                if (mResultData.respcd.equals("00")) {
                                    mHandler.sendEmptyMessage(Msg.MSG_FROM_SERVER_TRADE_SUCCESS);
                                } else if (mResultData.respcd.equals("09")) {
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

                        break;
                    }

                    case Msg.MSG_FROM_DIGLOG_CLOSE: {
                        ScanCodeGunActivity.this.finish();
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
                        mCustomDialog.nopay();
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


}
