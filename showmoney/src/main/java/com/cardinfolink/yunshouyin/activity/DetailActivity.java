package com.cardinfolink.yunshouyin.activity;

import android.content.Intent;
import android.graphics.Color;
import android.os.Bundle;
import android.text.TextUtils;
import android.view.View;
import android.view.animation.Animation;
import android.view.animation.AnimationUtils;
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
import com.cardinfolink.yunshouyin.ui.ResultInfoItem;
import com.cardinfolink.yunshouyin.ui.SettingActionBarItem;

import java.math.BigDecimal;
import java.text.ParseException;
import java.text.SimpleDateFormat;
import java.util.Date;

public class DetailActivity extends BaseActivity {

    private TradeBill mTradeBill;
    private SettingActionBarItem mActionBar;
    private TextView mPayResult;
    private ImageView mPayResultImage;

    private TextView mPayMoney;

    private ResultInfoItem mCardDiscount;//卡券折扣
    private ResultInfoItem mRefdMoney;//退款金额
    private ResultInfoItem mArriavlMoney;//到账金额


    private ResultInfoItem mPayChcd;//支付渠道
    private ResultInfoItem mPayTerminator;//操作终端
    private ResultInfoItem mPayDatetime;//支付时间
    private ResultInfoItem mPayOrder;//支付订单号
    private ResultInfoItem mPayType;//支付方式

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
                    Intent intent = new Intent(DetailActivity.this, RefdActivity.class);
                    Bundle bundle = new Bundle();
                    bundle.putSerializable("TradeBill", mTradeBill);
                    intent.putExtra("BillBundle", bundle);
                    startActivity(intent);
                }
            }
        });

        mPayResult = (TextView) findViewById(R.id.tv_pay_result);
        mPayResultImage = (ImageView) findViewById(R.id.iv_pay_result);

        mPayMoney = (TextView) findViewById(R.id.pay_money);//到账金额

        mCardDiscount = (ResultInfoItem) findViewById(R.id.card_discount);//卡券折扣
        mRefdMoney = (ResultInfoItem) findViewById(R.id.refd_money);
        mArriavlMoney = (ResultInfoItem) findViewById(R.id.pay_arrival_money);

        mPayChcd = (ResultInfoItem) findViewById(R.id.pay_chcd);
        mPayTerminator = (ResultInfoItem) findViewById(R.id.pay_terminator);
        mPayDatetime = (ResultInfoItem) findViewById(R.id.pay_datetime);
        mPayOrder = (ResultInfoItem) findViewById(R.id.pay_order);
        mPayType = (ResultInfoItem) findViewById(R.id.pay_type);

    }

    private void initData() {
        //这里是设置支付渠道
        if (!TextUtils.isEmpty(mTradeBill.chcd)) {
            if ("WXP".equals(mTradeBill.chcd)) {
                mPayChcd.setRightText(getString(R.string.detail_activity_pay_type1));
            } else {
                //支付宝
                mPayChcd.setRightText(getString(R.string.detail_activity_pay_type2));
            }
        } else {
            //其他支付渠道
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

        mPayTerminator.setRightText(SessonData.loginUser.getUsername());

        String tradeFrom = "";
        if (!TextUtils.isEmpty(mTradeBill.tradeFrom)) {
            tradeFrom = mTradeBill.tradeFrom;
        }

        mPayType.setRightText(tradeFrom + " 收款");

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
            mPayResultImage.setImageResource(R.drawable.bill_fresh);
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

        mPayMoney.setText(mTradeBill.amount);
        mRefdMoney.setRightText(mTradeBill.refundAmt);

    }


    //刷新按钮点击事件处理方法
    public void refreshOnclick(View view) {
        OrderData orderData = new OrderData();
        orderData.origOrderNum = mTradeBill.orderNum;
        startLoading(mPayResultImage);
        CashierSdk.startQy(orderData, new CashierListener() {

            @Override
            public void onResult(ResultData resultData) {
                quickPayService.getOrderAsync(SessonData.loginUser, mTradeBill.orderNum, new QuickPayCallbackListener<ServerPacket>() {
                    @Override
                    public void onSuccess(ServerPacket data) {
                        //注意这里使用了findOrder的新的的接口，这里txn返回的数组，不再是一个字段了，
                        // 这时候就没必要使用新的com.cardinfolink.yunshouyin.model.ServerPacketOrder的这个类了
                        mTradeBill.response = data.getTxn()[0].getResponse();
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


}
