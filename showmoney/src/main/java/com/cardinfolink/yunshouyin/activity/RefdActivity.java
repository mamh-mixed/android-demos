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
import com.cardinfolink.yunshouyin.data.SessonData;
import com.cardinfolink.yunshouyin.data.TradeBill;
import com.cardinfolink.yunshouyin.ui.EditTextClear;
import com.cardinfolink.yunshouyin.ui.SettingActionBarItem;
import com.cardinfolink.yunshouyin.ui.SettingInputItem;
import com.cardinfolink.yunshouyin.ui.SettingPasswordItem;
import com.cardinfolink.yunshouyin.util.ShowMoneyApp;
import com.cardinfolink.yunshouyin.util.Utility;

import java.math.BigDecimal;

public class RefdActivity extends BaseActivity {
    private static final String TAG = "RefdActivity";


    private TradeBill mTradeBill;

    private SettingActionBarItem mActionBar;
    private TextView mRefdAmount;
    private SettingInputItem mRefdMoney;
    private EditTextClear mPassword;
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
        mRefdMoney.setImageViewDrawable(null);
        mRefdMoney.setInputType(InputType.TYPE_CLASS_NUMBER | InputType.TYPE_NUMBER_FLAG_DECIMAL);//限制输入法只能是数字
        mPassword = (EditTextClear) findViewById(R.id.refd_password);

        mRefdButton = (Button) findViewById(R.id.btnrefd);
        mRefdButton.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                Utility.hideInput(mContext, v);
                refdOnClick(v);
            }
        });

        //获取可退款的金额
        getRefdTotal();

    }

    private void getRefdTotal() {
        try {

            BigDecimal a = new BigDecimal(mTradeBill.amount);
            BigDecimal b = new BigDecimal(mTradeBill.refundAmt);

            BigDecimal c = a.subtract(b);

            //本次可退款的金额
            maxRefd = c.setScale(2, BigDecimal.ROUND_HALF_UP).doubleValue();
        } catch (Exception e) {
            maxRefd = 0;
        }
        mRefdAmount.setText(String.valueOf(maxRefd));
    }

    //退款按钮的响应事件处理方法
    public void refdOnClick(View view) {
        if (!validate()) {
            return;
        }

        startLoading();
        OrderData orderData = new OrderData();
        orderData.origOrderNum = mTradeBill.orderNum;
        orderData.orderNum = Utility.geneOrderNumber();//生成一个新的订单号
        orderData.currency = CashierSdk.SDK_CURRENCY;
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
        if (!mPassword.getText().toString().equals(SessonData.loginUser.getPassword())) {
            alertMsg = ShowMoneyApp.getResString(R.string.refd_dialog__password_error);
            mAlertDialog.show(alertMsg, bitmap);
            return false;
        }

        return true;
    }
}
