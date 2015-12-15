package com.cardinfolink.yunshouyin.activity;

import android.annotation.SuppressLint;
import android.content.Intent;
import android.graphics.Bitmap;
import android.graphics.BitmapFactory;
import android.graphics.Color;
import android.os.Bundle;
import android.os.Handler;
import android.os.Message;
import android.text.TextUtils;
import android.view.View;
import android.widget.Button;
import android.widget.ImageView;
import android.widget.TextView;

import com.cardinfolink.cashiersdk.listener.CashierListener;
import com.cardinfolink.cashiersdk.model.OrderData;
import com.cardinfolink.cashiersdk.model.ResultData;
import com.cardinfolink.cashiersdk.sdk.CashierSdk;
import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.api.QuickPayException;
import com.cardinfolink.yunshouyin.constant.Msg;
import com.cardinfolink.yunshouyin.core.QuickPayCallbackListener;
import com.cardinfolink.yunshouyin.data.SessonData;
import com.cardinfolink.yunshouyin.data.TradeBill;
import com.cardinfolink.yunshouyin.model.ServerPacket;
import com.cardinfolink.yunshouyin.model.ServerPacketOrder;
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

    private Button mRefdButton; //退款按钮

    private Handler mHandler;

    private ImageView mBack;    //返回，这个是iamgeview
    private Button mRefreshButton;  //刷新按钮

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

        mRefdButton = (Button) findViewById(R.id.btn_refd);
        mRefdButton.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                btnRefdOnClick(v);//退款按钮的响应事件处理方法
            }
        });

        mBack = (ImageView) findViewById(R.id.iv_back);
        mBack.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                finish();
            }
        });
        mRefreshButton = (Button) findViewById(R.id.refresh);
        mRefreshButton.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                btnRefreshOnClick(v);//刷新按钮点击事件处理方法
            }
        });
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
            e.printStackTrace();
        }
        String tradeFrom = "PC";
        if (!TextUtils.isEmpty(mTradeBill.tradeFrom)) {
            tradeFrom = mTradeBill.tradeFrom;
        }
        String busicd = getResources().getString(R.string.detail_activity_busicd_pay);
        if (mTradeBill.busicd.equals("REFD")) {
            busicd = getResources().getString(R.string.detail_activity_busicd_refd);
        }

        mTradeFromText.setText(tradeFrom + busicd);
        String tradeStatus = getResources().getString(R.string.detail_activity_trade_status_success);

        if (mTradeBill.response.equals("00")) {
            tradeStatus = getResources().getString(R.string.detail_activity_trade_status_success);
            mTradeStatusText.setTextColor(Color.parseColor("#888888"));
        } else if (mTradeBill.response.equals("09")) {
            tradeStatus = getResources().getString(R.string.detail_activity_trade_status_nopay);
            mTradeStatusText.setTextColor(Color.RED);
        } else {
            tradeStatus = getResources().getString(R.string.detail_activity_trade_status_fail);
            mTradeStatusText.setTextColor(Color.RED);
        }
        mTradeStatusText.setText(tradeStatus);
        mConsumerAccount.setText(mTradeBill.consumerAccount);
        mTradeAmountText.setText("￥" + mTradeBill.amount);
        mOrderNumText.setText(mTradeBill.orderNum);
        mGoodInfoText.setText(mTradeBill.goodsInfo);
        if (mTradeBill.busicd.equals("REFD") || !mTradeBill.response.equals("00")) {
            mRefdButton.setVisibility(View.INVISIBLE);
        } else {
            mRefdButton.setVisibility(View.VISIBLE);
        }
    }

    //刷新按钮点击事件处理方法
    public void btnRefreshOnClick(View view) {
        OrderData orderData = new OrderData();
        orderData.origOrderNum = mTradeBill.orderNum;
        startLoading();
        CashierSdk.startQy(orderData, new CashierListener() {

            @Override
            public void onResult(ResultData resultData) {
                quickPayService.getOrderAsync(SessonData.loginUser, mTradeBill.orderNum, new QuickPayCallbackListener<ServerPacketOrder>() {
                    @Override
                    public void onSuccess(ServerPacketOrder data) {
                        mTradeBill.response = data.getTxn().getResponse();
                        endLoading();
                        initData();
                    }

                    @Override
                    public void onFailure(QuickPayException ex) {
                        endLoading();
                        initData();
                    }
                });
            }

            @Override
            public void onError(int errorCode) {

            }

        });


    }


    //退款按钮的响应事件处理方法
    public void btnRefdOnClick(View view) {
        startLoading();
        final String orderNum = mTradeBill.orderNum;
        final String amount = mTradeBill.amount;

        quickPayService.getRefdAsync(SessonData.loginUser, orderNum, new QuickPayCallbackListener<ServerPacket>() {
            @Override
            public void onSuccess(ServerPacket data) {
                String refdtotal = data.getRefdtotal();
                // 更新UI
                endLoading();

                //退款对话框
                View refdView = findViewById(R.id.refd_dialog);

                RefdDialog refdDialog = new RefdDialog(mContext, mHandler, refdView, orderNum, refdtotal, amount);
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
