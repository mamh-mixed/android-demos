package com.cardinfolink.yunshouyin.view;

import android.content.Context;
import android.graphics.Bitmap;
import android.graphics.BitmapFactory;
import android.os.Handler;
import android.text.TextUtils;
import android.view.MotionEvent;
import android.view.View;
import android.view.View.OnClickListener;
import android.view.View.OnTouchListener;
import android.widget.EditText;
import android.widget.TextView;

import com.cardinfolink.cashiersdk.listener.CashierListener;
import com.cardinfolink.cashiersdk.model.OrderData;
import com.cardinfolink.cashiersdk.model.ResultData;
import com.cardinfolink.cashiersdk.sdk.CashierSdk;
import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.activity.BaseActivity;
import com.cardinfolink.yunshouyin.data.SessonData;
import com.cardinfolink.yunshouyin.util.ShowMoneyApp;

import java.math.BigDecimal;
import java.text.SimpleDateFormat;
import java.util.Date;
import java.util.Random;


public class RefdDialog {
    private Bitmap wronBitmap;
    private Bitmap rightBitmap;

    private EditText refdValue;
    private EditText refdPassword;
    private TextView refdTitle;

    private Context mContext;

    private Handler mHandler;

    private View dialogView;

    private double maxRefd = 0;

    private String mOrderNum;

    private BaseActivity mBaseActivity;

    public RefdDialog(Context context, Handler handler, View view, String orderNum, String refdTotal, String total) {
        mContext = context;
        mBaseActivity = (BaseActivity) mContext;
        mHandler = handler;
        dialogView = view;
        mOrderNum = orderNum;
        maxRefd = Double.parseDouble(total) - Double.parseDouble(refdTotal);
        BigDecimal b = new BigDecimal(maxRefd);
        maxRefd = b.setScale(2, BigDecimal.ROUND_HALF_UP).doubleValue();

        rightBitmap = BitmapFactory.decodeResource(mContext.getResources(), R.drawable.right);
        wronBitmap = BitmapFactory.decodeResource(mContext.getResources(), R.drawable.wrong);

    }

    public void show() {
        refdTitle = (TextView) dialogView.findViewById(R.id.refd_title);
        refdTitle.setText(ShowMoneyApp.getResString(R.string.refd_dialog_refd_max) + maxRefd);

        refdValue = (EditText) dialogView.findViewById(R.id.refd_value_edit);
        refdPassword = (EditText) dialogView.findViewById(R.id.refd_password_edit);

        dialogView.setVisibility(View.VISIBLE);
        dialogView.setOnTouchListener(new OnTouchListener() {

            @Override
            public boolean onTouch(View v, MotionEvent event) {
                return true;
            }
        });

        dialogView.findViewById(R.id.refd_dialog_cancel).setOnClickListener(new OnClickListener() {

            @Override
            public void onClick(View v) {
                dialogView.setVisibility(View.GONE);
                refdPassword.setText("");
                refdValue.setText("");
            }
        });


        //退款
        dialogView.findViewById(R.id.refd_dialog_ok).setOnClickListener(new OnClickListener() {

            @Override
            public void onClick(View v) {
                refdOnClick(v);
            }
        });
    }

    /**
     * 退款
     *
     * @param view
     */
    private void refdOnClick(View view) {
        dialogView.setVisibility(View.GONE);

        if (!validate()) {
            refdPassword.setText("");
            refdValue.setText("");
            return;
        }

        mBaseActivity.startLoading();
        OrderData orderData = new OrderData();
        orderData.origOrderNum = mOrderNum;
        Date now = new Date();
        SimpleDateFormat spf = new SimpleDateFormat("yyMMddHHmmss");
        String orderNmuber = spf.format(now);
        Random random = new Random();
        for (int i = 0; i < 5; i++) {
            orderNmuber = orderNmuber + random.nextInt(10);
        }
        orderData.orderNum = orderNmuber;
        orderData.currency = "156";
        orderData.txamt = refdValue.getText().toString();

        CashierSdk.startRefd(orderData, new CashierListener() {

            @Override
            public void onResult(final ResultData resultData) {
                mBaseActivity.runOnUiThread(new Runnable() {

                    @Override
                    public void run() {
                        mBaseActivity.endLoading();
                        View alertView = mBaseActivity.findViewById(R.id.alert_dialog);
                        String alertMsg = "";
                        if (resultData.respcd.equals("00")) {
                            alertMsg = ShowMoneyApp.getResString(R.string.refd_dialog_refd_success);
                            AlertDialog alertDialog = new AlertDialog(mContext, mHandler, alertView, alertMsg, rightBitmap);
                            alertDialog.show();
                        } else if (resultData.respcd.equals("R6")) {
                            alertMsg = ShowMoneyApp.getResString(R.string.refd_dialog_nextday_not_refd);
                            AlertDialog alertDialog = new AlertDialog(mContext, null, alertView, alertMsg, wronBitmap);
                            alertDialog.show();
                        } else {
                            alertMsg = ShowMoneyApp.getResString(R.string.refd_dialog_refd_fail);
                            AlertDialog alertDialog = new AlertDialog(mContext, null, alertView, alertMsg, wronBitmap);
                            alertDialog.show();
                        }
                    }

                });
            }

            @Override
            public void onError(int errorCode) {
                mBaseActivity.runOnUiThread(new Runnable() {

                    @Override
                    public void run() {
                        mBaseActivity.endLoading();
                        View alertView = mBaseActivity.findViewById(R.id.alert_dialog);
                        String alertMsg = ShowMoneyApp.getResString(R.string.refd_dialog_refd_fail);
                        AlertDialog alertDialog = new AlertDialog(mContext, null, alertView, alertMsg, wronBitmap);
                        alertDialog.show();
                    }
                });
            }
        });

        refdPassword.setText("");
        refdValue.setText("");
    }

    private boolean validate() {
        View alertView = mBaseActivity.findViewById(R.id.alert_dialog);
        String alertMsg = "";

        String refdStr = refdValue.getText().toString();
        if (TextUtils.isEmpty(refdStr)) {
            alertMsg = ShowMoneyApp.getResString(R.string.refd_dialog_amount_cannot_empty);
            AlertDialog alertDialog = new AlertDialog(mContext, null, alertView, alertMsg, wronBitmap);
            alertDialog.show();
            return false;
        }

        double refd = 0;
        try {
            refd = Double.parseDouble(refdStr);
        } catch (Exception e) {
            alertMsg = ShowMoneyApp.getResString(R.string.refd_dialog_amount_foramt_error);
            AlertDialog alertDialog = new AlertDialog(mContext, null, alertView, alertMsg, wronBitmap);
            alertDialog.show();
            return false;
        }

        if (refd < 0.01) {
            alertMsg = ShowMoneyApp.getResString(R.string.refd_dialog_amount_not_enough);
            AlertDialog alertDialog = new AlertDialog(mContext, null, alertView, alertMsg, wronBitmap);
            alertDialog.show();
            return false;
        }

        if (refd > maxRefd) {
            alertMsg = String.format(ShowMoneyApp.getResString(R.string.refd_dialog_amount_not_exceeds_max), maxRefd);
            AlertDialog alertDialog = new AlertDialog(mContext, null, alertView, alertMsg, wronBitmap);
            alertDialog.show();
            return false;
        }
        if (!refdPassword.getText().toString().equals(SessonData.loginUser.getPassword())) {
            alertMsg = ShowMoneyApp.getResString(R.string.refd_dialog__password_error);
            AlertDialog alertDialog = new AlertDialog(mContext, null, alertView, alertMsg, wronBitmap);
            alertDialog.show();
            return false;
        }

        return true;
    }


}
