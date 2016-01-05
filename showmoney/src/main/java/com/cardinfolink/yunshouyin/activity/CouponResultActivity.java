package com.cardinfolink.yunshouyin.activity;

import android.app.Activity;
import android.content.Context;
import android.content.Intent;
import android.os.Bundle;
import android.util.Log;
import android.view.View;
import android.widget.Button;
import android.widget.TextView;

import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.constant.Msg;
import com.cardinfolink.yunshouyin.data.Coupon;
import com.cardinfolink.yunshouyin.ui.SettingActionBarItem;

/**
 * Created by charles on 2015/12/29.
 */
public class CouponResultActivity extends Activity {

    private static final String TAG = "CouponResultActivity";
    private Context mContext;

    private TextView mCouponContent;
    private Button mPayByCash;//现金收款按钮
    private Button mPayByScanCode;//扫码付款按钮
    private SettingActionBarItem mActionBar;


    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_coupon_result);
        mContext = this;

        mCouponContent = (TextView) findViewById(R.id.tv_coupon_message);
        mPayByCash = (Button) findViewById(R.id.bt_money);
        mPayByScanCode = (Button) findViewById(R.id.bt_scancode);
        mActionBar = (SettingActionBarItem) findViewById(R.id.action_bar);

        Intent intent = getIntent();
        Bundle bundle = intent.getExtras();
        boolean isSuccess = bundle.getBoolean("check_coupon_result_flag", false);//核销成功失败的标记
        if (isSuccess) {
            //换物品
            if ("2".equals(Coupon.getInstance().getVoucherType())) {
                mCouponContent.setText(Coupon.getInstance().getCardId());
                mPayByScanCode.setVisibility(View.GONE);
                mPayByCash.setText(getString(R.string.coupon_confirm_ok));
            } else {//打折
                //指定扫码支付
                Boolean preferenceScancode = Coupon.getInstance().getVoucherType() != null && (Coupon.getInstance().getVoucherType().startsWith("4")
                        || Coupon.getInstance().getVoucherType().startsWith("5"));
                //指定扫码支付
                Log.e(TAG, "preferenceScancode:" + preferenceScancode);
                if (preferenceScancode) {
                    mPayByCash.setVisibility(View.INVISIBLE);
                }

                mCouponContent.setText(Coupon.getInstance().getCardId());
                mPayByScanCode.setOnClickListener(new View.OnClickListener() {
                    @Override
                    public void onClick(View v) {
                        MainActivity.getHandler().sendEmptyMessage(Msg.MSG_FROM_SERVER_COUPON_SUCCESS);
                        finish();
                    }
                });
            }
        } else {
            mCouponContent.setText(getString(R.string.coupon_verial_fail_info));
            mPayByScanCode.setVisibility(View.GONE);
            mPayByCash.setText(getString(R.string.coupon_goback));
        }

        //现金支付
        mPayByCash.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                Coupon.getInstance().clear();//清空卡券信息
                finish();

            }
        });
        //返回销券
        mActionBar.setLeftTextOnclickListner(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                Coupon.getInstance().clear();//清空卡券信息
                finish();
            }
        });

    }
}
