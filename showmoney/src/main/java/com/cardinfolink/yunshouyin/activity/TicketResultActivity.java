package com.cardinfolink.yunshouyin.activity;

import android.app.Activity;
import android.content.Intent;
import android.graphics.Color;
import android.os.Bundle;
import android.view.View;
import android.widget.Button;
import android.widget.LinearLayout;
import android.widget.TextView;

import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.ui.ResultInfoItem;
import com.cardinfolink.yunshouyin.ui.SettingActionBarItem;

/**
 * Created by charles on 2015/12/25.
 */
public class TicketResultActivity extends Activity {


    private SettingActionBarItem mActionBar;

    private Button mCouponScanAgain;

    private TextView mCouponRes;
    private TextView mCoupoxnFaiRes;
    private TextView mCouponInf;
    private LinearLayout mLCouponMoSi;

    private ResultInfoItem mCouponOrSu;
    private ResultInfoItem mCouponAccReMo;
    private ResultInfoItem mCouponOrNu;
    private ResultInfoItem mCouponPaNum;
    private ResultInfoItem mCouponAcc;
    private ResultInfoItem mCouponUsTim;
    private ResultInfoItem mCouponTerm;
    private ResultInfoItem mCouponDis;

    private String mOrigMon;
    private String mDiscountMon;
    private String mActualMon;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_ticket_result);
        mActionBar = (SettingActionBarItem) findViewById(R.id.sabi_coupon_action_bar);//标题栏
        mCouponScanAgain = (Button) findViewById(R.id.btn_coupon_scan_again);//重新扫码支付
        mCouponRes = (TextView) findViewById(R.id.tv_coupon_result);//支付结果
        mCoupoxnFaiRes = (TextView) findViewById(R.id.tv_coupon_fail_reason);//支付失败原因
        mCouponInf = (TextView) findViewById(R.id.tv_coupon_info);//风云再起满50减40
        mLCouponMoSi = (LinearLayout) findViewById(R.id.ll_coupon_money_situation);//控制原，卡券，应到账金额控制
        mCouponOrSu = (ResultInfoItem) findViewById(R.id.sii_coupon_original_sum);//原金额
        mCouponDis = (ResultInfoItem) findViewById(R.id.sii_coupon_coupon_discount);//卡券抵消金额
        mCouponAccReMo = (ResultInfoItem) findViewById(R.id.sii_coupon_account_receive_money);//应到账金额
        mCouponOrNu = (ResultInfoItem) findViewById(R.id.sii_coupon_order_number);//核销订单号
        mCouponPaNum = (ResultInfoItem) findViewById(R.id.sii_coupon_pay_number);//支付订单号
        mCouponAcc = (ResultInfoItem) findViewById(R.id.sii_coupon_access);//卡卷渠道
        mCouponUsTim = (ResultInfoItem) findViewById(R.id.sii_tv_coupon_use_time);//核销时间
        mCouponTerm = (ResultInfoItem) findViewById(R.id.sii_terminal);//操作终端


        Intent intent = getIntent();
        Bundle bundle = intent.getExtras();
        Boolean flag = bundle.getBoolean("flag");
        mOrigMon = bundle.getString("originalmoney");
        mDiscountMon = bundle.getString("dicountmoney");
        mActualMon = bundle.getString("actualpaymoney");
        if (flag) {
            //核销成功
            verifiForDiscountSucc();
        } else {
            //核销失败
            verifiForDiscoFailSpecifyAcc();
        }

    }

    //兑换劵成功页面
    public void verifiForItemSucc() {
        mActionBar.setLeftText(getResources().getString(R.string.coupon_cancel_return));//返回核销
        mActionBar.setLeftTextOnclickListner(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                //title 返回核销
                finish();
            }
        });
        mCouponRes.setText(getResources().getString(R.string.pay_result_success));
        mCouponRes.setTextColor(Color.BLUE);
        mCoupoxnFaiRes.setVisibility(View.INVISIBLE);

        mCouponOrSu.setRightText(mOrigMon);//设置原金额
        mCouponDis.setRightText(mDiscountMon);//设置抵消金额
        mCouponAccReMo.setRightText(mActualMon);//设置应到账金额
        mCouponOrNu.setRightText("");//设置核销订单号
        mCouponPaNum.setRightText("");//支付订单号
        mCouponAcc.setRightText("");//卡卷渠道
        mCouponUsTim.setRightText("");//核销时间
        mCouponTerm.setRightText("");//操作终端
    }

    //折扣券支付成功页面
    public void verifiForDiscountSucc() {


    }

    //折扣券指定支付方式失败页面
    public void verifiForDiscoFailSpecifyAcc() {
        mActionBar.setLeftText(getResources().getString(R.string.coupon_cancel_return));//返回核销
        mActionBar.setLeftTextOnclickListner(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                //title 返回核销
                finish();
            }
        });
        mCouponRes.setText(getResources().getString(R.string.pay_result_fail));
        mCouponRes.setTextColor(Color.RED);

        mCouponOrSu.setRightText(mOrigMon);//设置原金额
        mCouponDis.setRightText(mDiscountMon);//设置抵消金额
        mCouponAccReMo.setRightText(mActualMon);//设置应到账金额
        mCouponOrNu.setRightText("");//设置核销订单号
        mCouponPaNum.setRightText("");//支付订单号
        mCouponAcc.setRightText("");//卡卷渠道
        mCouponUsTim.setRightText("");//核销时间
        mCouponTerm.setRightText("");//操作终端


    }


    //折扣券未指定支付方式失败页面
    public void verifiForDiscoFailNoAcc() {
        mActionBar.setLeftText(getResources().getString(R.string.coupon_cancel_return));////现金收款

        mCouponRes.setText(getResources().getString(R.string.pay_result_fail));
        mCouponRes.setTextColor(Color.RED);

        mCouponOrSu.setRightText(mOrigMon);//设置原金额
        mCouponDis.setRightText(mDiscountMon);//设置抵消金额
        mCouponAccReMo.setRightText(mActualMon);//设置应到账金额
        mCouponOrNu.setRightText("");//设置核销订单号
        mCouponPaNum.setRightText("");//支付订单号
        mCouponAcc.setRightText("");//卡卷渠道
        mCouponUsTim.setRightText("");//核销时间
        mCouponTerm.setRightText("");//操作终端


        mCouponScanAgain.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                //重新扫码支付
            }
        });

        mActionBar.setLeftTextOnclickListner(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                //现金收款
            }
        });

    }
}
