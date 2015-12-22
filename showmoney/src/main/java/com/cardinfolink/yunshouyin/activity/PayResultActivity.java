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

import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.data.SessonData;
import com.cardinfolink.yunshouyin.ui.SettingActionBarItem;

/**
 * Created by charles on 2015/12/20.
 */
public class PayResultActivity extends Activity {
    private SettingActionBarItem mActionBar;

    private TextView mPayResult;
    private ImageView mPayResultPhoto;
    private TextView mResultExplain;
    private TextView mPersonAccount;
    private TextView mMakeDealTime;
    private TextView mBillOrderNum;
    private TextView mReceiveMoney;
    private Button mConfirm;
    private ImageView mPayAccess;
    private TextView mReceiveMoneyStatus;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_pay_result);

        mActionBar = (SettingActionBarItem) findViewById(R.id.action_bar);
        mActionBar.setLeftTextOnclickListner(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                finish();
            }
        });

        //获取交易参数
        Intent intent = getIntent();
        Bundle bundle = intent.getExtras();
        String txamt = "";
        String orderNum = "";
        String chcd = "";
        Boolean result = false;
        String currentTime = "";
        String resultExplain = "";
        if (bundle != null) {
            txamt = bundle.getString("txamt");
            orderNum = bundle.getString("orderNum");
            chcd = bundle.getString("chcd");
            currentTime = bundle.getString("mCurrentTime");
            result = bundle.getBoolean("result");
            resultExplain = bundle.getString("errorDetail");
        }

        mPayResult = (TextView) findViewById(R.id.tv_pay_result);
        mPayResultPhoto = (ImageView) findViewById(R.id.iv_pay_result);
        mResultExplain = (TextView) findViewById(R.id.tv_result_explain);
        mPersonAccount = (TextView) findViewById(R.id.pay_account);
        mMakeDealTime = (TextView) findViewById(R.id.pay_datetime);
        mBillOrderNum = (TextView) findViewById(R.id.order_number);
        mReceiveMoney = (TextView) findViewById(R.id.total_money);
        mReceiveMoneyStatus = (TextView) findViewById(R.id.total_state);
        mConfirm = (Button) findViewById(R.id.btnconfirm);
        mPayAccess = (ImageView) findViewById(R.id.pay_chcd);

        if (result) {
            mPayResult.setText(R.string.pay_result_success);
            mPayResult.setTextColor(Color.BLUE);
            mPayResultPhoto.setImageResource(R.drawable.right);
            mResultExplain.setText("");
            mReceiveMoneyStatus.setTextColor(Color.BLUE);
            mReceiveMoney.setTextColor(Color.BLUE);
            mReceiveMoneyStatus.setText(R.string.pay_total_state_success);
        } else {
            mPayResult.setText(R.string.pay_result_fail);
            mPayResult.setTextColor(Color.RED);
            mPayResultPhoto.setImageResource(R.drawable.wrong);
            mResultExplain.setText(resultExplain);
            mReceiveMoneyStatus.setTextColor(Color.RED);
            mReceiveMoney.setTextColor(Color.RED);
            mReceiveMoneyStatus.setText(R.string.pay_total_state_fail);
        }
        mPersonAccount.setText(SessonData.loginUser.getUsername());
        mMakeDealTime.setText(currentTime);
        mBillOrderNum.setText(orderNum);
        mReceiveMoney.setText(txamt);

        if ("ALP".equals(chcd)) {
            mPayAccess.setImageResource(R.drawable.scan_alipay);
        } else if ("WXP".equals(chcd)) {
            mPayAccess.setImageResource(R.drawable.scan_wechat);
        } else {
            mPayAccess.setImageResource(R.drawable.wrong);
        }

        mConfirm.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                finish();
            }
        });
    }
}
