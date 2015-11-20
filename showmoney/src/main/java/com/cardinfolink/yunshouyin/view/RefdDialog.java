package com.cardinfolink.yunshouyin.view;

import android.annotation.SuppressLint;
import android.app.Activity;
import android.content.Context;
import android.graphics.BitmapFactory;
import android.os.Handler;
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
    EditText refdValue;
    EditText refdPassword;
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
    }

    public void show() {

        TextView textView = (TextView) dialogView.findViewById(R.id.refd_title);
        textView.setText(ShowMoneyApp.getResString(R.string.refd_dialog_refd_max) + maxRefd);
        refdValue = (EditText) dialogView.findViewById(R.id.refd_value_edit);
        refdPassword = (EditText) dialogView.findViewById(R.id.refd_password_edit);


        dialogView.setVisibility(View.VISIBLE);
        dialogView.setOnTouchListener(new OnTouchListener() {

            @Override
            public boolean onTouch(View v, MotionEvent event) {
                // TODO Auto-generated method stub
                return true;
            }
        });

        dialogView.findViewById(R.id.refd_dialog_cancel).setOnClickListener(new OnClickListener() {

            @Override
            public void onClick(View v) {
                //    DeviceManageUtil.hideInput(mContext);
                dialogView.setVisibility(View.GONE);
                refdPassword.setText("");
                refdValue.setText("");
            }
        });


        dialogView.findViewById(R.id.refd_dialog_ok).setOnClickListener(new OnClickListener() {

            @Override
            public void onClick(View v) {
                //	 DeviceManageUtil.hideInput(mContext);
                dialogView.setVisibility(View.GONE);

                if (check()) {
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
                    ;
                    orderData.orderNum = orderNmuber;
                    orderData.currency = "156";
                    orderData.txamt = refdValue.getText().toString();

                    CashierSdk.startRefd(orderData, new CashierListener() {

                        @Override
                        public void onResult(final ResultData resultData) {
                            ((Activity) mContext).runOnUiThread(new Runnable() {

                                @Override
                                public void run() {
                                    mBaseActivity.endLoading();
                                    if (resultData.respcd.equals("00")) {
                                        AlertDialog alert_Dialog = new AlertDialog(mContext, mHandler, ((Activity) mContext).findViewById(R.id.alert_dialog), ShowMoneyApp.getResString(R.string.refd_dialog_refd_success), BitmapFactory.decodeResource(mContext.getResources(), R.drawable.right));
                                        alert_Dialog.show();

                                    } else if (resultData.respcd.equals("R6")) {
                                        AlertDialog alert_Dialog = new AlertDialog(mContext, null, ((Activity) mContext).findViewById(R.id.alert_dialog), ShowMoneyApp.getResString(R.string.refd_dialog_nextday_not_refd), BitmapFactory.decodeResource(mContext.getResources(), R.drawable.right));
                                        alert_Dialog.show();
                                    } else {
                                        AlertDialog alert_Dialog = new AlertDialog(mContext, null, ((Activity) mContext).findViewById(R.id.alert_dialog), ShowMoneyApp.getResString(R.string.refd_dialog_refd_fail), BitmapFactory.decodeResource(mContext.getResources(), R.drawable.right));
                                        alert_Dialog.show();
                                    }

                                }

                            });


                        }

                        @Override
                        public void onError(int errorCode) {
                            ((Activity) mContext).runOnUiThread(new Runnable() {

                                @Override
                                public void run() {
                                    mBaseActivity.endLoading();
                                    AlertDialog alert_Dialog = new AlertDialog(mContext, null, ((Activity) mContext).findViewById(R.id.alert_dialog), ShowMoneyApp.getResString(R.string.refd_dialog_refd_fail), BitmapFactory.decodeResource(mContext.getResources(), R.drawable.right));
                                    alert_Dialog.show();

                                }

                            });


                        }
                    });


                }

                refdPassword.setText("");
                refdValue.setText("");
            }
        });
    }

    @SuppressLint("NewApi")
    private boolean check() {

        if (refdValue.getText().toString().isEmpty()) {
            AlertDialog alert_Dialog = new AlertDialog(mContext, null, ((Activity) mContext).findViewById(R.id.alert_dialog), ShowMoneyApp.getResString(R.string.refd_dialog_amount_cannot_empty), BitmapFactory.decodeResource(mContext.getResources(), R.drawable.wrong));
            alert_Dialog.show();
            return false;
        }

        double refd = 0;
        try {

            refd = Double.parseDouble(refdValue.getText().toString());
        } catch (Exception e) {
            AlertDialog alert_Dialog = new AlertDialog(mContext, null, ((Activity) mContext).findViewById(R.id.alert_dialog), ShowMoneyApp.getResString(R.string.refd_dialog_amount_foramt_error), BitmapFactory.decodeResource(mContext.getResources(), R.drawable.wrong));
            alert_Dialog.show();
            return false;
        }

        if (refd < 0.01) {
            AlertDialog alert_Dialog = new AlertDialog(mContext, null, ((Activity) mContext).findViewById(R.id.alert_dialog), ShowMoneyApp.getResString(R.string.refd_dialog_amount_not_enough), BitmapFactory.decodeResource(mContext.getResources(), R.drawable.wrong));
            alert_Dialog.show();
            return false;
        }

        if (refd > maxRefd) {
            AlertDialog alert_Dialog = new AlertDialog(mContext, null, ((Activity) mContext).findViewById(R.id.alert_dialog), String.format(ShowMoneyApp.getResString(R.string.refd_dialog_amount_not_exceeds_max), maxRefd), BitmapFactory.decodeResource(mContext.getResources(), R.drawable.wrong));
            alert_Dialog.show();
            return false;
        }
        if (!refdPassword.getText().toString().equals(SessonData.loginUser.getPassword())) {
            AlertDialog alert_Dialog = new AlertDialog(mContext, null, ((Activity) mContext).findViewById(R.id.alert_dialog), ShowMoneyApp.getResString(R.string.refd_dialog__password_error), BitmapFactory.decodeResource(mContext.getResources(), R.drawable.wrong));
            alert_Dialog.show();
            return false;
        }


        return true;
    }


}
