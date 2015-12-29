package com.cardinfolink.yunshouyin.activity;

import android.content.Intent;
import android.graphics.Bitmap;
import android.graphics.BitmapFactory;
import android.graphics.Color;
import android.os.Bundle;
import android.text.TextUtils;
import android.view.View;
import android.view.animation.Animation;
import android.view.animation.AnimationUtils;
import android.widget.Button;
import android.widget.ImageView;
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
import com.cardinfolink.yunshouyin.model.ServerPacketOrder;
import com.cardinfolink.yunshouyin.ui.ResultInfoItem;
import com.cardinfolink.yunshouyin.ui.SettingActionBarItem;
import com.cardinfolink.yunshouyin.view.RefdDialog;

import java.math.BigDecimal;
import java.text.ParseException;
import java.text.SimpleDateFormat;
import java.util.Date;

public class DetailActivity extends BaseActivity {

    private TradeBill mTradeBill;
    private SettingActionBarItem mActionBar;
    private TextView mPayResult;
    private ImageView mPayResultImage;
    private ResultInfoItem mRealPayMoney;
    private ResultInfoItem mPayMoney;
    private ResultInfoItem mRefdMoney;
    private ResultInfoItem mPayChcd;
    private ResultInfoItem mPayAccount;
    private ResultInfoItem mPayDatetime;
    private ResultInfoItem mPayOrder;
    private ResultInfoItem mTradeFrom;
    private Button mClose;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_detail);
        Intent intent = getIntent();
        Bundle billBundle = intent.getBundleExtra("BillBundle");
        mTradeBill = (TradeBill) billBundle.get("TradeBill");
        initLayout();
        initData();
    }

    private void initLayout() {
        mActionBar = (SettingActionBarItem) findViewById(R.id.action_bar);
        mActionBar.setLeftTextOnclickListner(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                finish();
            }
        });
        mActionBar.setRightTextOnclickListner(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                //退款按钮
                if ("REFD".equals(mTradeBill.busicd) || "CANC".equals(mTradeBill.busicd) || !"00".equals(mTradeBill.response)) {
                } else {
                    refdOnClick(v);//退款按钮
                }
            }
        });

        mPayResult = (TextView) findViewById(R.id.tv_pay_result);
        mPayResultImage = (ImageView) findViewById(R.id.iv_pay_result);
        mRealPayMoney = (ResultInfoItem) findViewById(R.id.real_pay_money);
        mPayMoney = (ResultInfoItem) findViewById(R.id.pay_money);
        mRefdMoney = (ResultInfoItem) findViewById(R.id.refd_money);

        mPayChcd = (ResultInfoItem) findViewById(R.id.pay_chcd);
        mPayAccount = (ResultInfoItem) findViewById(R.id.pay_account);
        mPayDatetime = (ResultInfoItem) findViewById(R.id.pay_datetime);
        mPayOrder = (ResultInfoItem) findViewById(R.id.pay_order);
        mTradeFrom = (ResultInfoItem) findViewById(R.id.pay_tradefrom);

        mClose = (Button) findViewById(R.id.btnclose);
        mClose.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                finish();
            }
        });
    }

    private void initData() {
        //这里是设置支付渠道
        if (!TextUtils.isEmpty(mTradeBill.chcd)) {
            if ("WXP".equals(mTradeBill.chcd)) {
                mPayChcd.setRightText("微信");
            } else {
                mPayChcd.setRightText("支付宝");
            }
        } else {
            mPayChcd.setRightText("");
        }

        SimpleDateFormat spf1 = new SimpleDateFormat("yyyyMMddHHmmss");
        SimpleDateFormat spf2 = new SimpleDateFormat("yyyy-MM-dd HH:mm:ss");
        try {
            Date tandeDate = spf1.parse(mTradeBill.tandeDate);
            mPayDatetime.setRightText(spf2.format(tandeDate));
        } catch (ParseException e) {
            mPayDatetime.setRightText("");
        }

        mPayAccount.setRightText(SessonData.loginUser.getUsername());

        String tradeFrom = "";
        if (!TextUtils.isEmpty(mTradeBill.tradeFrom)) {
            tradeFrom = mTradeBill.tradeFrom;
        }

        mTradeFrom.setRightText(tradeFrom + " 收款");

        mPayOrder.setRightText(mTradeBill.orderNum);

        //设置交易的结果状态
        String tradeStatus = "";
        if ("00".equals(mTradeBill.response)) {
            //成功的
            tradeStatus = getResources().getString(R.string.detail_activity_trade_status_success);
            mPayResult.setTextColor(Color.parseColor("#00bbd3"));
            mPayResultImage.setImageResource(R.drawable.pay_result_succeed);
        } else if ("09".equals(mTradeBill.response)) {
            //正在处理中的交易
            tradeStatus = getResources().getString(R.string.detail_activity_trade_status_nopay);
            mPayResult.setTextColor(Color.BLACK);
            mPayResultImage.setImageResource(R.drawable.refresh);
            mPayResultImage.setOnClickListener(new View.OnClickListener() {
                @Override
                public void onClick(View v) {
                    refreshOnclick(v);
                }
            });
        } else {
            //其他
            tradeStatus = getResources().getString(R.string.detail_activity_trade_status_fail);
            mPayResult.setTextColor(Color.RED);
            mPayResultImage.setImageResource(R.drawable.pay_result_fail);
        }
        mPayResult.setText(tradeStatus);

        mPayMoney.setRightText(mTradeBill.amount);
        mRefdMoney.setRightText(mTradeBill.refundAmt);
        BigDecimal bgPay = new BigDecimal(mTradeBill.amount);
        BigDecimal bgRefd = new BigDecimal(mTradeBill.refundAmt);
        String realPay = bgPay.subtract(bgRefd).toString();
        mRealPayMoney.setRightText(realPay);


    }


    //刷新按钮点击事件处理方法
    public void refreshOnclick(View view) {
        OrderData orderData = new OrderData();
        orderData.origOrderNum = mTradeBill.orderNum;
        startLoading(mPayResultImage);
        CashierSdk.startQy(orderData, new CashierListener() {

            @Override
            public void onResult(ResultData resultData) {
                quickPayService.getOrderAsync(SessonData.loginUser, mTradeBill.orderNum, new QuickPayCallbackListener<ServerPacketOrder>() {
                    @Override
                    public void onSuccess(ServerPacketOrder data) {
                        mTradeBill.response = data.getTxn().getResponse();
                        endLoading(mPayResultImage);
                        initData();
                    }

                    @Override
                    public void onFailure(QuickPayException ex) {
                        endLoading(mPayResultImage);
                        initData();
                    }
                });
            }

            @Override
            public void onError(int errorCode) {
                endLoading(mPayResultImage);
            }

        });


    }


    public void startLoading(View view) {
        Animation loadingAnimation = AnimationUtils.loadAnimation(mContext, R.anim.loading_animation);
        view.startAnimation(loadingAnimation);
    }

    public void endLoading(View view) {
        view.clearAnimation();
    }

    //退款按钮的响应事件处理方法
    public void refdOnClick(View view) {
        startLoading();
        final String orderNum = mTradeBill.orderNum;
        final String amount = mTradeBill.amount;

        quickPayService.getRefdAsync(SessonData.loginUser, orderNum, new QuickPayCallbackListener<ServerPacket>() {
            @Override
            public void onSuccess(ServerPacket data) {
                String refdtotal = data.getRefdtotal();
                // 更新UI
                endLoading();
                View refdView = findViewById(R.id.refd_dialog);
                RefdDialog refdDialog = new RefdDialog(mContext, null, refdView, orderNum, refdtotal, amount);

                refdDialog.show();
            }

            @Override
            public void onFailure(QuickPayException ex) {
                // 更新UI,显示提示对话框
                endLoading();
                String error = ex.getErrorMsg();
                Bitmap bitmap = BitmapFactory.decodeResource(mContext.getResources(), R.drawable.wrong);
                mAlertDialog.show(error, bitmap);
            }
        });
    }


}
