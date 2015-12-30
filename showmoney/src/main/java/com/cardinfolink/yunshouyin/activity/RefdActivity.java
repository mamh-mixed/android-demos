package com.cardinfolink.yunshouyin.activity;

import android.content.Intent;
import android.graphics.Bitmap;
import android.graphics.BitmapFactory;
import android.os.Bundle;
import android.text.InputType;
import android.text.TextUtils;
import android.view.View;
import android.widget.Button;
import android.widget.TextView;

import com.cardinfolink.cashiersdk.listener.CashierListener;
import com.cardinfolink.cashiersdk.model.OrderData;
import com.cardinfolink.cashiersdk.model.ResultData;
import com.cardinfolink.cashiersdk.sdk.CashierSdk;
import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.api.QuickPayException;
import com.cardinfolink.yunshouyin.core.QuickPayCallbackListener;
import com.cardinfolink.yunshouyin.data.SessonData;
import com.cardinfolink.yunshouyin.data.TradeBill;
import com.cardinfolink.yunshouyin.model.ServerPacket;
import com.cardinfolink.yunshouyin.ui.SettingActionBarItem;
import com.cardinfolink.yunshouyin.ui.SettingInputItem;
import com.cardinfolink.yunshouyin.ui.SettingPasswordItem;
import com.cardinfolink.yunshouyin.util.ShowMoneyApp;

import java.math.BigDecimal;
import java.text.SimpleDateFormat;
import java.util.Date;
import java.util.Random;

public class RefdActivity extends BaseActivity {
    private TradeBill mTradeBill;
    private SettingActionBarItem mActionBar;
    private TextView mRefdAmount;
    private SettingInputItem mRefdMoney;
    private SettingPasswordItem mPassword;
    private Button mRefdButton;
    private double maxRefd = 0;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_refd);
        Intent intent = getIntent();
        Bundle billBundle = intent.getBundleExtra("BillBundle");
        mTradeBill = (TradeBill) billBundle.get("TradeBill");
        initLayout();
    }

    private void initLayout() {
        mActionBar = (SettingActionBarItem) findViewById(R.id.action_bar);
        mActionBar.setLeftTextOnclickListner(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                finish();
            }
        });

        mRefdAmount = (TextView) findViewById(R.id.refd_amount);
        mRefdMoney = (SettingInputItem) findViewById(R.id.refd_money);
        mRefdMoney.setInputType(InputType.TYPE_CLASS_NUMBER);//限制输入法只能是数字
        mPassword = (SettingPasswordItem) findViewById(R.id.refd_password);

        mRefdButton = (Button) findViewById(R.id.btnrefd);
        mRefdButton.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                refdOnClick(v);
            }
        });

        //获取可退款的金额
        getRefdTotal();

    }

    private void getRefdTotal() {
        startLoading();
        final String orderNum = mTradeBill.orderNum;
        final String amount = mTradeBill.amount;

        quickPayService.getRefdAsync(SessonData.loginUser, orderNum, new QuickPayCallbackListener<ServerPacket>() {
            @Override
            public void onSuccess(ServerPacket data) {
                String refdtotal = data.getRefdtotal();

                BigDecimal a = new BigDecimal(amount);
                BigDecimal b = new BigDecimal(refdtotal);
                BigDecimal c = a.subtract(b);
                //本次可退款的金额
                maxRefd = c.setScale(2, BigDecimal.ROUND_HALF_UP).doubleValue();

                mRefdAmount.setText(c.toString());
                endLoading();
            }

            @Override
            public void onFailure(QuickPayException ex) {
                endLoading();
                // 更新UI,显示提示对话框
                endLoading();
                String error = ex.getErrorMsg();
                Bitmap bitmap = BitmapFactory.decodeResource(mContext.getResources(), R.drawable.wrong);
                mAlertDialog.show(error, bitmap);
            }
        });
    }

    //退款按钮的响应事件处理方法
    public void refdOnClick(View view) {
        if (!validate()) {
            return;
        }

        startLoading();
        OrderData orderData = new OrderData();
        orderData.origOrderNum = mTradeBill.orderNum;
        Date now = new Date();
        SimpleDateFormat spf = new SimpleDateFormat("yyMMddHHmmss");
        String orderNmuber = spf.format(now);
        Random random = new Random();
        for (int i = 0; i < 5; i++) {
            orderNmuber = orderNmuber + random.nextInt(10);
        }
        orderData.orderNum = orderNmuber;
        orderData.currency = "156";
        orderData.txamt = mRefdMoney.getText();

        CashierSdk.startRefd(orderData, new CashierListener() {

            @Override
            public void onResult(final ResultData resultData) {
                runOnUiThread(new Runnable() {

                    @Override
                    public void run() {
                        endLoading();
                        String alertMsg = "";
                        if (resultData.respcd.equals("00")) {
                            alertMsg = ShowMoneyApp.getResString(R.string.refd_dialog_refd_success);
                            mAlertDialog.setTitle(alertMsg);
                            mAlertDialog.setImageViewResource(R.drawable.right);
                            mAlertDialog.setOkOnClickListener(new View.OnClickListener() {
                                @Override
                                public void onClick(View v) {
                                    mAlertDialog.hide();
                                    finish();
                                }
                            });
                            mAlertDialog.show();
                        } else if (resultData.respcd.equals("R6")) {
                            alertMsg = ShowMoneyApp.getResString(R.string.refd_dialog_nextday_not_refd);
                            mAlertDialog.setTitle(alertMsg);
                            mAlertDialog.setImageViewResource(R.drawable.wrong);
                            mAlertDialog.setOkOnClickListener(new View.OnClickListener() {
                                @Override
                                public void onClick(View v) {
                                    mAlertDialog.hide();
                                    finish();
                                }
                            });
                            mAlertDialog.show();
                        } else {
                            alertMsg = ShowMoneyApp.getResString(R.string.refd_dialog_refd_fail);
                            mAlertDialog.setTitle(alertMsg);
                            mAlertDialog.setImageViewResource(R.drawable.wrong);
                            mAlertDialog.setOkOnClickListener(new View.OnClickListener() {
                                @Override
                                public void onClick(View v) {
                                    mAlertDialog.hide();
                                    finish();
                                }
                            });
                            mAlertDialog.show();
                        }
                    }

                });
            }

            @Override
            public void onError(int errorCode) {
                runOnUiThread(new Runnable() {

                    @Override
                    public void run() {
                        endLoading();
                        String alertMsg = ShowMoneyApp.getResString(R.string.refd_dialog_refd_fail);
                        mAlertDialog.setTitle(alertMsg);
                        mAlertDialog.setImageViewResource(R.drawable.wrong);
                        mAlertDialog.setOkOnClickListener(new View.OnClickListener() {
                            @Override
                            public void onClick(View v) {
                                mAlertDialog.hide();
                                finish();
                            }
                        });
                        mAlertDialog.show();
                    }
                });
            }
        });


    }

    private boolean validate() {
        Bitmap bitmap = BitmapFactory.decodeResource(mContext.getResources(), R.drawable.wrong);
        String alertMsg = "";

        String refdStr = mRefdMoney.getText();
        if (TextUtils.isEmpty(refdStr)) {
            alertMsg = ShowMoneyApp.getResString(R.string.refd_dialog_amount_cannot_empty);
            mAlertDialog.show(alertMsg, bitmap);
            return false;
        }

        double refd = 0;
        try {
            refd = Double.parseDouble(refdStr);
        } catch (Exception e) {
            alertMsg = ShowMoneyApp.getResString(R.string.refd_dialog_amount_foramt_error);
            mAlertDialog.show(alertMsg, bitmap);
            return false;
        }

        if (refd < 0.01) {
            alertMsg = ShowMoneyApp.getResString(R.string.refd_dialog_amount_not_enough);
            mAlertDialog.show(alertMsg, bitmap);
            return false;
        }

        if (refd > maxRefd) {
            alertMsg = String.format(ShowMoneyApp.getResString(R.string.refd_dialog_amount_not_exceeds_max), maxRefd);
            mAlertDialog.show(alertMsg, bitmap);
            return false;
        }
        if (!mPassword.getPassword().equals(SessonData.loginUser.getPassword())) {
            alertMsg = ShowMoneyApp.getResString(R.string.refd_dialog__password_error);
            mAlertDialog.show(alertMsg, bitmap);
            return false;
        }

        return true;
    }
}
