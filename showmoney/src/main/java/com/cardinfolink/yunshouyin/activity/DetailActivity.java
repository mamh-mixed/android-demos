package com.cardinfolink.yunshouyin.activity;

import android.annotation.SuppressLint;
import android.content.Intent;
import android.graphics.BitmapFactory;
import android.graphics.Color;
import android.os.Bundle;
import android.os.Handler;
import android.os.Message;
import android.view.View;
import android.widget.Button;
import android.widget.ImageView;
import android.widget.TextView;

import com.cardinfolink.cashiersdk.listener.CashierListener;
import com.cardinfolink.cashiersdk.model.OrderData;
import com.cardinfolink.cashiersdk.model.ResultData;
import com.cardinfolink.cashiersdk.sdk.CashierSdk;
import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.constant.Msg;
import com.cardinfolink.yunshouyin.data.SessonData;
import com.cardinfolink.yunshouyin.data.TradeBill;
import com.cardinfolink.yunshouyin.util.CommunicationListener;
import com.cardinfolink.yunshouyin.util.ErrorUtil;
import com.cardinfolink.yunshouyin.util.HttpCommunicationUtil;
import com.cardinfolink.yunshouyin.util.JsonUtil;
import com.cardinfolink.yunshouyin.util.ParamsUtil;
import com.cardinfolink.yunshouyin.view.RefdDialog;

import java.text.ParseException;
import java.text.SimpleDateFormat;
import java.util.Date;

@SuppressLint("SimpleDateFormat")
public class DetailActivity extends BaseActivity {
    private TradeBill mTradeBill;
    private ImageView mPaylogoImage;
    private TextView mTradeFromText;
    private TextView mTradeDateText;
    private TextView mTradeStatusText;
    private TextView mConsumerAccount;
    private TextView mTradeAmountText;
    private TextView mOrderNumText;
    private TextView mGoodInfoText;
    private Button mRefdButton;
    private Handler mHandler;


    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.detail_activity);
        Intent intent = getIntent();
        Bundle billBundle = intent.getBundleExtra("BillBundle");
        mTradeBill = (TradeBill) billBundle.get("TradeBill");
        initLayout();
        initData();
        initHandler();
    }

    private void initLayout() {
        mPaylogoImage = (ImageView) findViewById(R.id.paylogo);
        mTradeFromText = (TextView) findViewById(R.id.tradefrom);
        mTradeDateText = (TextView) findViewById(R.id.tradedate);
        mTradeStatusText = (TextView) findViewById(R.id.tradestatus);
        mConsumerAccount = (TextView) findViewById(R.id.consumer_account);
        mTradeAmountText = (TextView) findViewById(R.id.tradeamount);
        mOrderNumText = (TextView) findViewById(R.id.ordernum);
        mGoodInfoText = (TextView) findViewById(R.id.goodinfo);
        mRefdButton = (Button) findViewById(R.id.detail_btn_refd);
    }

    @SuppressLint("NewApi")
    private void initData() {
        if (mTradeBill.chcd.equals("WXP")) {
            mPaylogoImage.setImageResource(R.drawable.wpay);
        } else {
            mPaylogoImage.setImageResource(R.drawable.apay);
        }
        SimpleDateFormat spf1 = new SimpleDateFormat("yyyyMMddHHmmss");
        SimpleDateFormat spf2 = new SimpleDateFormat("yyyy-MM-dd HH:mm:ss");
        try {
            Date tandeDate = spf1.parse(mTradeBill.tandeDate);
            mTradeDateText.setText(spf2.format(tandeDate));
        } catch (ParseException e) {
            // TODO Auto-generated catch block
            e.printStackTrace();
        }
        String tradeFrom = "PC";
        if (!mTradeBill.tradeFrom.isEmpty()) {
            tradeFrom = mTradeBill.tradeFrom;
        }
        String busicd = getResources().getString(R.string.detail_activity_busicd_pay);
        if (mTradeBill.busicd.equals("REFD")) {
            busicd = getResources().getString(R.string.detail_activity_busicd_refd);
        }

        mTradeFromText.setText(tradeFrom + busicd);
        String tradeStatus = getResources().getString(R.string.detail_activity_trade_status_success);
        ;
        if (mTradeBill.response.equals("00")) {
            tradeStatus = getResources().getString(R.string.detail_activity_trade_status_success);
            ;
            ;
            mTradeStatusText.setTextColor(Color.parseColor("#888888"));
        } else if (mTradeBill.response.equals("09")) {
            tradeStatus = getResources().getString(R.string.detail_activity_trade_status_nopay);
            ;
            ;
            mTradeStatusText.setTextColor(Color.RED);
        } else {
            tradeStatus = getResources().getString(R.string.detail_activity_trade_status_fail);
            ;
            mTradeStatusText.setTextColor(Color.RED);
        }
        mTradeStatusText.setText(tradeStatus);
        mConsumerAccount.setText(mTradeBill.consumerAccount);
        mTradeAmountText.setText("￥" + mTradeBill.amount);
        mOrderNumText.setText(mTradeBill.orderNum);
        mGoodInfoText.setText(mTradeBill.goodsInfo);
        if (mTradeBill.busicd.equals("REFD")
                || !mTradeBill.response.equals("00")) {
            mRefdButton.setVisibility(View.INVISIBLE);
        } else {
            mRefdButton.setVisibility(View.VISIBLE);
        }
    }

    public void BtnBackOnClick(View view) {

        DetailActivity.this.finish();

    }

    public void BtnRefreshOnClick(View view) {

        OrderData orderData = new OrderData();
        orderData.origOrderNum = mTradeBill.orderNum;
        startLoading();
        CashierSdk.startQy(orderData, new CashierListener() {

            @Override
            public void onResult(ResultData resultData) {

                HttpCommunicationUtil.sendDataToServer(ParamsUtil.getOrder(SessonData.loginUser, mTradeBill.orderNum), new CommunicationListener() {

                    @Override
                    public void onResult(String result) {
                        String state = JsonUtil.getParam(result, "state");
                        if (state.equals("success")) {
                            mTradeBill.response = JsonUtil.getParam(JsonUtil.getParam(result, "txn"), "response");

                        }
                        runOnUiThread(new Runnable() {

                            @Override
                            public void run() {
                                // 更新UI
                                endLoading();
                                initData();
                            }

                        });

                    }

                    @Override
                    public void onError(String error) {
                        runOnUiThread(new Runnable() {

                            @Override
                            public void run() {
                                // 更新UI
                                endLoading();
                                initData();
                            }

                        });

                    }
                });


            }

            @Override
            public void onError(int errorCode) {

            }

        });


    }


    public void BtnRefdOnClick(View view) {
        startLoading();
        HttpCommunicationUtil.sendDataToServer(
                ParamsUtil.getRefd(SessonData.loginUser, mTradeBill.orderNum),
                new CommunicationListener() {

                    @Override
                    public void onResult(final String result) {
                        String state = JsonUtil.getParam(result, "state");
                        if (state.equals("success")) {
                            final String refdtotal = JsonUtil.getParam(result,
                                    "refdtotal");
                            runOnUiThread(new Runnable() {

                                @Override
                                public void run() {
                                    // 更新UI
                                    endLoading();
                                    RefdDialog refd_Dialog = new RefdDialog(
                                            DetailActivity.this, mHandler,
                                            findViewById(R.id.refd_dialog),
                                            mTradeBill.orderNum, refdtotal,
                                            mTradeBill.amount);
                                    refd_Dialog.show();
                                }

                            });

                        } else {
                            runOnUiThread(new Runnable() {

                                @Override
                                public void run() {
                                    // 更新UI
                                    endLoading();
                                    mAlertDialog.show(ErrorUtil
                                                    .getErrorString(JsonUtil.getParam(
                                                            result, "error")),
                                            BitmapFactory.decodeResource(
                                                    mContext.getResources(),
                                                    R.drawable.wrong));
                                }

                            });
                        }

                    }

                    @Override
                    public void onError(final String error) {
                        runOnUiThread(new Runnable() {

                            @Override
                            public void run() {
                                // 更新UI
                                endLoading();
                                mAlertDialog.show(error, BitmapFactory
                                        .decodeResource(
                                                mContext.getResources(),
                                                R.drawable.wrong));
                            }

                        });

                    }
                });

    }

    private void initHandler() {
        mHandler = new Handler() {
            @Override
            public void handleMessage(Message msg) {
                switch (msg.what) {
                    case Msg.MSG_FROM_CLIENT_ALERT_OK: {
                        SessonData.positionView = 1;
                        setResult(101);
                        finish();
                    }
                }


                super.handleMessage(msg);
            }
        };
    }

}
