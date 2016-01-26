package com.cardinfolink.yunshouyin.activity;

import android.content.Context;
import android.content.Intent;
import android.os.Bundle;
import android.view.KeyEvent;
import android.view.View;
import android.widget.Button;
import android.widget.TextView;
import android.widget.Toast;

import com.cardinfolink.cashiersdk.listener.CashierListener;
import com.cardinfolink.cashiersdk.model.OrderData;
import com.cardinfolink.cashiersdk.model.ResultData;
import com.cardinfolink.cashiersdk.sdk.CashierSdk;
import com.cardinfolink.cashiersdk.util.TxamtUtil;
import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.data.Coupon;
import com.cardinfolink.yunshouyin.ui.SettingActionBarItem;
import com.cardinfolink.yunshouyin.util.ErrorUtil;
import com.cardinfolink.yunshouyin.util.Utility;
import com.cardinfolink.yunshouyin.view.HintDialog;

import java.math.BigDecimal;

/**
 * Created by charles on 2015/12/29.
 */
public class CouponResultActivity extends BaseActivity {

    private static final String TAG = "CouponResultActivity";
    private Context mContext;

    private TextView mCouponContent;
    private TextView mCouponMessage;

    private Button mPayByCash;//现金收款按钮
    private Button mPayByScanCode;//扫码付款按钮
    private SettingActionBarItem mActionBar;
    private HintDialog mHintDialog;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_coupon_result);
        mContext = this;
        mHintDialog = new HintDialog(CouponResultActivity.this, findViewById(R.id.hint_dialog));
        mCouponContent = (TextView) findViewById(R.id.tv_coupon_content);
        mCouponMessage = (TextView) findViewById(R.id.tv_coupon_message);

        mPayByCash = (Button) findViewById(R.id.bt_money);
        mPayByScanCode = (Button) findViewById(R.id.bt_scancode);
        mActionBar = (SettingActionBarItem) findViewById(R.id.action_bar);
        mActionBar.setLeftText(getString(R.string.coupon_vertical));

        Intent intent = getIntent();
        Bundle bundle = intent.getExtras();
        boolean isSuccess = bundle.getBoolean("check_coupon_result_flag", false);
        String failRespcd=bundle.getString("respcd");
        //核销成功
        if (isSuccess) {
            mCouponMessage.setText(R.string.coupon_verial_success_info);

            //换物品
            if ("2".equals(Coupon.getInstance().getVoucherType())) {
                mCouponContent.setText(Coupon.getInstance().getCardId());
                mPayByScanCode.setVisibility(View.GONE);
                mPayByCash.setText(getString(R.string.coupon_confirm_ok));
            } else {
                //打折
                if (Coupon.getInstance().getVoucherType().startsWith("4")) {
                    // 4开头的指定微信支付。
                    mCouponMessage.setText(R.string.coupon_verial_success_info_weixin);
                } else if (Coupon.getInstance().getVoucherType().startsWith("5")) {
                    // 5开头的表示 指定支付宝的
                    mCouponMessage.setText(R.string.coupon_verial_success_info_ali);
                }

                //指定扫码支付,5开头的表示 指定支付宝的，4开头的指定微信支付。
                Boolean isHidePayByCash = Coupon.getInstance().getVoucherType() != null && (Coupon.getInstance().getVoucherType().startsWith("4") || Coupon.getInstance().getVoucherType().startsWith("5"));
                //指定扫码支付,这里判断是否要隐藏 现金支付 这个按钮
                if (isHidePayByCash) {
                    mPayByCash.setVisibility(View.INVISIBLE);
                }

                mPayByScanCode.setOnClickListener(new View.OnClickListener() {
                    @Override
                    public void onClick(View v) {
                        Intent intent = new Intent(CouponResultActivity.this, ScanCodeActivity.class);
                        startActivity(intent);
                        finish();
                    }
                });
                String saleMinAccount = TxamtUtil.getNormal(Coupon.getInstance().getSaleMinAmount());
                String discount = TxamtUtil.getNormal(Coupon.getInstance().getSaleDiscount());
                if (Coupon.getInstance().getVoucherType().endsWith("3")) {
                    //满折券
                    discount = new BigDecimal(discount).multiply(new BigDecimal("10")).toString();
                    mCouponContent.setText(Coupon.getInstance().getCardId() + getString(R.string.coupon_man) + saleMinAccount + getString(R.string.coupon_yuan) + getString(R.string.coupon_da) + discount + getString(R.string.coupon_zhe));
                } else if (Coupon.getInstance().getVoucherType().endsWith("1")) {
                    //满减券
                    mCouponContent.setText(Coupon.getInstance().getCardId() + getString(R.string.coupon_man) + saleMinAccount + getString(R.string.coupon_yuan) + getString(R.string.coupon_jian) + discount + getString(R.string.coupon_yuan));
                } else {
                    mCouponContent.setText(Coupon.getInstance().getCardId() + getString(R.string.coupon_jian) + discount + getString(R.string.coupon_yuan));
                    //返回销券
                }
            }
        } else {//核销失败

            mCouponMessage.setText(ErrorUtil.getErrorString(failRespcd));
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

        if (Coupon.getInstance().getVoucherType() != null) {
            if (Coupon.getInstance().getVoucherType().startsWith("4") || Coupon.getInstance().getVoucherType().startsWith("5")) {
                mActionBar.setLeftTextOnclickListner(new View.OnClickListener() {
                    @Override
                    public void onClick(View v) {
                        showPayFailPref();
                    }
                });


            } else {
                //返回销券
                mActionBar.setLeftTextOnclickListner(new View.OnClickListener() {
                    @Override
                    public void onClick(View v) {
                        Coupon.getInstance().clear();//清空卡券信息
                        finish();
                    }
                });
            }
        } else {
            mActionBar.setLeftTextOnclickListner(new View.OnClickListener() {
                @Override
                public void onClick(View v) {
                    Coupon.getInstance().clear();//清空卡券信息
                    finish();
                }
            });
        }

    }

    //弹出是否放弃本子核销对话框
    public void showPayFailPref() {
        mHintDialog.setText(getString(R.string.coupon_abandom_verial_or_not), getString(R.string.coupon_pay_again), getString(R.string.coupon_abandom));
        mHintDialog.setOkVisibility(View.GONE);
        mHintDialog.setCancelOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                //卡券冲正
                OrderData orderData = new OrderData();
                orderData.orderNum = Utility.geneOrderNumber();//订单号
                orderData.origOrderNum = Coupon.getInstance().getOrderNum();//设置原始订单号
                CashierSdk.startReversal(orderData, new CashierListener() {
                    @Override
                    public void onResult(ResultData resultData) {
                        if ("00".equals(resultData.respcd)) {
                            //冲正成功
                            Coupon.getInstance().clear();
                            runOnUiThread(new Runnable() {
                                @Override
                                public void run() {
                                    Toast.makeText(CouponResultActivity.this, getString(R.string.coupon_verial_success), Toast.LENGTH_SHORT).show();
                                    finish();
                                }
                            });
                        } else {
                            //冲正失败
                            runOnUiThread(new Runnable() {
                                @Override
                                public void run() {
                                    Toast.makeText(CouponResultActivity.this, getString(R.string.coupon_verial_fail), Toast.LENGTH_SHORT).show();
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
                mHintDialog.hide();
                finish();
            }

        });
        mHintDialog.show();
    }


    @Override
    public boolean onKeyDown(int keyCode, KeyEvent event) {
        if (keyCode == KeyEvent.KEYCODE_BACK) {
            if (Coupon.getInstance().getVoucherType() != null) {
                if (Coupon.getInstance().getVoucherType().startsWith("4") || Coupon.getInstance().getVoucherType().startsWith("5")) {
                    showPayFailPref();
                } else {
                    Coupon.getInstance().clear();//清空卡券信息
                    finish();
                }
            } else {
                Coupon.getInstance().clear();//清空卡券信息
                finish();
            }
        }
        return false;
    }
}
