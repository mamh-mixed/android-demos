package com.cardinfolink.yunshouyin.activity;

import android.app.Activity;
import android.os.Bundle;
import android.widget.Button;
import android.widget.TextView;

import com.cardinfolink.yunshouyin.R;

/**
 * Created by charles on 2015/12/29.
 */
public class CouponResultActivity extends Activity {

    private TextView mCouponContent;
    private Button mMoney;
    private Button mScanCode;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_coupon_result);

        mCouponContent = (TextView) findViewById(R.id.tv_coupon_message);
        mMoney = (Button) findViewById(R.id.bt_money);
        mScanCode = (Button) findViewById(R.id.bt_scancode);

    }
}
