package com.cardinfolink.yunshouyin.activity;

import android.app.Activity;
import android.content.Intent;
import android.graphics.Color;
import android.os.Bundle;
import android.util.Log;
import android.view.View;
import android.widget.Button;
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
import com.cardinfolink.yunshouyin.ui.ResultInfoItem;
import com.cardinfolink.yunshouyin.ui.SettingActionBarItem;
import com.cardinfolink.yunshouyin.view.HintDialog;

import java.math.BigDecimal;

/**
 * Created by charles on 2015/12/20.
 */
public class PayResultActivity extends Activity {
    private SettingActionBarItem mActionBar;

    private TextView mPayResult;
    private ImageView mPayResultPhoto;
    private TextView mResultExplain;
    private ResultInfoItem mPersonAccount;
    private ResultInfoItem mMakeDealTime;
    private ResultInfoItem mBillOrderNum;
    private TextView mReceiveMoney;
    private Button mConfirm;
    private ImageView mPayAccess;
    private TextView mReceiveMoneyStatus;
    private ResultInfoItem mCouponContent;
    private ResultInfoItem mActualTotalMoney;
    private ResultInfoItem mActualDiscount;
    private HintDialog mHintDialog;


    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_pay_result);

        mActionBar = (SettingActionBarItem) findViewById(R.id.action_bar);
        mHintDialog = new HintDialog(PayResultActivity.this, findViewById(R.id.hint_dialog));

        //获取交易参数
        Intent intent = getIntent();
        Bundle bundle = intent.getExtras();
        String txamt = "";
        String orderNum = "";
        String chcd = "";
        Boolean result = false;
        String currentTime = "";
        String resultExplain = "";
        String originaltotal = "";
        String total = "";
        if (bundle != null) {
            txamt = bundle.getString("txamt");
            orderNum = bundle.getString("orderNum");
            chcd = bundle.getString("chcd");
            currentTime = bundle.getString("mCurrentTime");
            result = bundle.getBoolean("result");
            resultExplain = bundle.getString("errorDetail");
            originaltotal = bundle.getString("originaltotal");
            total = bundle.getString("total");
        }

        mPayResult = (TextView) findViewById(R.id.tv_pay_result);
        mPayResultPhoto = (ImageView) findViewById(R.id.iv_pay_result);
        mResultExplain = (TextView) findViewById(R.id.tv_result_explain);

        mPersonAccount = (ResultInfoItem) findViewById(R.id.pay_account);
        mMakeDealTime = (ResultInfoItem) findViewById(R.id.pay_datetime);
        mBillOrderNum = (ResultInfoItem) findViewById(R.id.order_number);
        mCouponContent = (ResultInfoItem) findViewById(R.id.rif_coupon_content);
        mActualTotalMoney = (ResultInfoItem) findViewById(R.id.rif_coupon_total);
        mActualDiscount = (ResultInfoItem) findViewById(R.id.rif_coupon_discount);

        mReceiveMoney = (TextView) findViewById(R.id.total_money);
        mReceiveMoneyStatus = (TextView) findViewById(R.id.total_state);
        mConfirm = (Button) findViewById(R.id.btnconfirm);
        mPayAccess = (ImageView) findViewById(R.id.pay_chcd);
        boolean hasCouponDiscount = SessonData.loginUser.getResultData().saleDiscount != null &&
                !"0".equals(SessonData.loginUser.getResultData().saleDiscount);
        //判断是否有卡券优惠
        if (hasCouponDiscount) {
            mCouponContent.setVisibility(View.VISIBLE);
            mActualTotalMoney.setVisibility(View.VISIBLE);
            mActualDiscount.setVisibility(View.VISIBLE);

            mCouponContent.setRightText(SessonData.loginUser.getResultData().cardId);//卡券内容
            mActualTotalMoney.setRightText(originaltotal);//消费金额
            mActualDiscount.setRightText(String.valueOf(new BigDecimal(originaltotal).subtract(new BigDecimal(total)).doubleValue()));//优惠金额
        } else {
            mCouponContent.setVisibility(View.GONE);
            mActualTotalMoney.setVisibility(View.GONE);
            mActualDiscount.setVisibility(View.GONE);
        }

        //判断支付是否成功
        if (result) {
            mPayResult.setText(R.string.pay_result_success);
            mPayResultPhoto.setImageResource(R.drawable.pay_result_succeed);
            mResultExplain.setText("");
            //标题收款键
            mActionBar.setLeftTextOnclickListner(new View.OnClickListener() {
                @Override
                public void onClick(View v) {
                    SessonData.loginUser.setResultData(null);
                    finish();
                }
            });
        } else {
            mPayResult.setText(R.string.pay_result_fail);
            mPayResult.setTextColor(Color.RED);
            mPayResultPhoto.setImageResource(R.drawable.pay_result_fail);
            mResultExplain.setText(resultExplain);
            mReceiveMoneyStatus.setTextColor(Color.RED);
            mReceiveMoney.setTextColor(Color.RED);
            mReceiveMoneyStatus.setText(R.string.pay_total_state_fail);
            if (hasCouponDiscount) {
                mActionBar.setLeftText(getString(R.string.coupon_getmoney));
                mActionBar.setLeftTextOnclickListner(new View.OnClickListener() {
                    @Override
                    public void onClick(View v) {
                        showPayFail();
                    }
                });
            } else {
                mActionBar.setLeftTextOnclickListner(new View.OnClickListener() {
                    @Override
                    public void onClick(View v) {
                        finish();
                    }
                });
            }
        }
        mPersonAccount.setRightText(SessonData.loginUser.getUsername());
        mMakeDealTime.setRightText(currentTime);
        mBillOrderNum.setRightText(orderNum);
        mReceiveMoney.setText(txamt);

        if ("ALP".equals(chcd)) {
            mPayAccess.setImageResource(R.drawable.scan_alipay);
        } else if ("WXP".equals(chcd)) {
            mPayAccess.setImageResource(R.drawable.scan_wechat);
        } else {
            mPayAccess.setImageResource(R.drawable.pay_result_fail);
        }

        mConfirm.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                SessonData.loginUser.setResultData(null);
                finish();
            }
        });
    }

    //弹出是否放弃本子核销对话框
    public void showPayFail() {
        mHintDialog.setText(getString(R.string.coupon_abandom_verial_or_not), getString(R.string.coupon_abandom), getString(R.string.coupon_pay_again));
        mHintDialog.setCancelOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                //卡券冲正
                OrderData orderData = new OrderData();
                orderData.orderNum = SessonData.loginUser.getResultData().orderNum;//订单号
                orderData.origOrderNum = "";//原始订单号
                CashierSdk.startReversal(orderData, new CashierListener() {
                    @Override
                    public void onResult(ResultData resultData) {
                        if ("00".equals(resultData.respcd)) {
                            //冲正成功
                            SessonData.loginUser.setResultData(null);
                            runOnUiThread(new Runnable() {
                                @Override
                                public void run() {
                                    Toast.makeText(PayResultActivity.this, getString(R.string.coupon_verial_success), Toast.LENGTH_SHORT).show();
                                    finish();
                                }
                            });
                        } else {
                            //冲正失败
                            runOnUiThread(new Runnable() {
                                @Override
                                public void run() {
                                    Toast.makeText(PayResultActivity.this, getString(R.string.coupon_verial_fail), Toast.LENGTH_SHORT).show();
                                }
                            });
                        }
                    }

                    @Override
                    public void onError(int errorCode) {

                    }
                });
                mHintDialog.hide();
            }
        });
        mHintDialog.setOkOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                //重新支付
                MainActivity.getHandler().sendEmptyMessage(Msg.MSG_FROM_SERVER_COUPON_SUCCESS);
                mHintDialog.hide();
                finish();
            }

        });
        mHintDialog.show();
    }
}
