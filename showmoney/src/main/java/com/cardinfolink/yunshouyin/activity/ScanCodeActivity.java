package com.cardinfolink.yunshouyin.activity;

import android.app.Activity;
import android.os.Bundle;
import android.os.Handler;
import android.os.Message;
import android.view.KeyEvent;
import android.view.View;
import android.widget.LinearLayout;
import android.widget.Toast;

import com.cardinfolink.cashiersdk.listener.CashierListener;
import com.cardinfolink.cashiersdk.model.OrderData;
import com.cardinfolink.cashiersdk.model.ResultData;
import com.cardinfolink.cashiersdk.sdk.CashierSdk;
import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.constant.Msg;
import com.cardinfolink.yunshouyin.data.Coupon;
import com.cardinfolink.yunshouyin.util.Utility;
import com.cardinfolink.yunshouyin.view.HintDialog;
import com.cardinfolink.yunshouyin.view.ScanCodeView;

public class ScanCodeActivity extends Activity {

    private static Handler scanCodehandler;
    private HintDialog mHintDialog;
    private ScanCodeView mScanCodeView;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_scan_code);

        LinearLayout.LayoutParams layoutParams = new LinearLayout.LayoutParams(LinearLayout.LayoutParams.MATCH_PARENT, LinearLayout.LayoutParams.MATCH_PARENT);
        mScanCodeView = new ScanCodeView(this, MainActivity.getHandler());
        mScanCodeView.setLayoutParams(layoutParams);

        addContentView(mScanCodeView, layoutParams);

        mHintDialog = mScanCodeView.getmHintDialog();

        scanCodehandler = new Handler() {
            @Override
            public void handleMessage(Message msg) {
                super.handleMessage(msg);
                switch (msg.what) {
                    case Msg.MSG_FINISH_BIG_SCANCODEVIEW:
                        finish();
                        break;
                }
            }
        };

    }

    public static Handler getScanCodehandler() {
        return scanCodehandler;
    }

    public static void setScanCodehandler(Handler scanCodehandler) {
        ScanCodeActivity.scanCodehandler = scanCodehandler;
    }


    public void showPayFailPref() {
        mHintDialog.setText(getString(R.string.coupon_abandom_verial_or_not), getString(R.string.coupon_pay_again), getString(R.string.coupon_abandom));
        mHintDialog.setCancelOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                //卡券冲正
                OrderData orderData = new OrderData();
                orderData.orderNum = Utility.geneOrderNumber();//订单号
                orderData.origOrderNum = Coupon.getInstance().getOrderNum();//设置原始订单号
                CashierSdk.startReversal(orderData, new CashierListener() {
                    @Override
                    public void onResult(ResultData resultData) {
                        if ("00".equals(resultData.respcd)) {
                            //冲正成功
                            Coupon.getInstance().clear();
                            runOnUiThread(new Runnable() {
                                @Override
                                public void run() {
                                    Toast.makeText(ScanCodeActivity.this, getString(R.string.coupon_verial_success), Toast.LENGTH_SHORT).show();
                                    finish();
                                }
                            });
                        } else {
                            //冲正失败
                            runOnUiThread(new Runnable() {
                                @Override
                                public void run() {
                                    Toast.makeText(ScanCodeActivity.this, getString(R.string.coupon_verial_fail), Toast.LENGTH_SHORT).show();
                                }
                            });
                        }
                    }

                    @Override
                    public void onError(int errorCode) {

                    }
                });
                if (ScanCodeActivity.getScanCodehandler() != null) {
                    ScanCodeActivity.getScanCodehandler().sendEmptyMessage(Msg.MSG_FINISH_BIG_SCANCODEVIEW);
                }


                mHintDialog.hide();
            }
        });
        mHintDialog.setOkOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                //重新支付
                MainActivity.getHandler().sendEmptyMessage(Msg.MSG_FROM_SERVER_COUPON_SUCCESS);
                mHintDialog.hide();
                mScanCodeView.clearValue();

            }

        });
        mHintDialog.show();
    }

    @Override
    public boolean onKeyDown(int keyCode, KeyEvent event) {
        if (keyCode == KeyEvent.KEYCODE_BACK) {
            if (Coupon.getInstance().getVoucherType() != null) {
                if (Coupon.getInstance().getVoucherType().startsWith("4") || Coupon.getInstance().getVoucherType().startsWith("5")) {
                    showPayFailPref();
                } else {
                    finish();
                }
            } else {
                finish();
            }
        }
        return false;
    }
}
