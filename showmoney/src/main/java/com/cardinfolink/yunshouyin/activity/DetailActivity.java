package com.cardinfolink.yunshouyin.activity;

import android.content.Intent;
import android.graphics.Color;
import android.os.Bundle;
import android.text.TextUtils;
import android.util.Log;
import android.view.View;
import android.view.animation.Animation;
import android.view.animation.AnimationUtils;
import android.widget.ImageView;
import android.widget.TextView;

import com.cardinfolink.cashiersdk.model.OrderData;
import com.cardinfolink.cashiersdk.util.TxamtUtil;
import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.api.QuickPayException;
import com.cardinfolink.yunshouyin.core.QuickPayCallbackListener;
import com.cardinfolink.yunshouyin.data.SessonData;
import com.cardinfolink.yunshouyin.data.TradeBill;
import com.cardinfolink.yunshouyin.model.QRequest;
import com.cardinfolink.yunshouyin.model.ServerPacket;
import com.cardinfolink.yunshouyin.model.Txn;
import com.cardinfolink.yunshouyin.ui.ResultInfoItem;
import com.cardinfolink.yunshouyin.ui.SettingActionBarItem;

import java.math.BigDecimal;
import java.text.ParseException;
import java.text.SimpleDateFormat;
import java.util.Date;

/**
 * 这个是显示账单 详情界面，只会卡券还有个详情界面
 */
public class DetailActivity extends BaseActivity {
    private static final String TAG = "DetailActivity";

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
                mPayChcd.setRightText(getString(R.string.detail_activity_chcd_type1));
            } else if ("ALP".equals(mTradeBill.chcd)) {
                //支付宝
                mPayChcd.setRightText(getString(R.string.detail_activity_chcd_type2));
            } else {
                //其他支付渠道
                mPayChcd.setRightText(getString(R.string.detail_activity_chcd_type3));
            }
        } else {
            //这里表明没有渠道，默认先设置为 其他支付渠道
            mPayChcd.setRightText(getString(R.string.detail_activity_chcd_type3));
        }

        SimpleDateFormat spf1 = new SimpleDateFormat("yyyyMMddHHmmss");
        SimpleDateFormat spf2 = new SimpleDateFormat("yyyy-MM-dd HH:mm:ss");
        try {
            Date tandeDate = spf1.parse(mTradeBill.tandeDate);
            mPayDatetime.setRightText(spf2.format(tandeDate));
        } catch (ParseException e) {
            mPayDatetime.setRightText("");
        }

        //支付终端
        mPayTerminator.setRightText(SessonData.loginUser.getUsername());

        //支付方式
        if (!TextUtils.isEmpty(mTradeBill.tradeFrom)) {
            if ("android".equals(mTradeBill.tradeFrom) || "ios".equals(mTradeBill.tradeFrom)) {
                //app收款
                mPayType.setRightText(getString(R.string.detail_activity_pay_type1));
            } else if ("wap".equals(mTradeBill.tradeFrom)) {
                //web 收款
                mPayType.setRightText(getString(R.string.detail_activity_pay_type2));
            } else {
                //其他收款
                mPayType.setRightText(getString(R.string.detail_activity_pay_type3));
            }
        } else {
            //tradeFrom 是 null 清空
            mPayType.setRightText(getString(R.string.detail_activity_pay_type4));
        }

        //设置订单号
        mPayOrder.setRightText(mTradeBill.orderNum);

        //设置交易的结果状态
        if ("10".equals(mTradeBill.transStatus)) {
            //处理中
            mPayResult.setText(getString(R.string.detail_activity_trade_status_nopay));
            mPayResultImage.setImageResource(R.drawable.bill_fresh);
            mPayResultImage.setOnClickListener(new View.OnClickListener() {
                @Override
                public void onClick(View v) {
                    refreshOnclick(v);
                }
            });
        } else if ("30".equals(mTradeBill.transStatus)) {
            double amt = Double.parseDouble(mTradeBill.refundAmt);
            if (amt == 0) {
                //成功的
                mPayResult.setTextColor(Color.parseColor("#00bbd3"));
                mPayResult.setText(getString(R.string.detail_activity_trade_status_success));
                mPayResultImage.setImageResource(R.drawable.pay_result_succeed);
            } else {
                //部分退款的
                mPayResult.setText(getString(R.string.detail_activity_trade_status_partrefd));
            }
        } else if ("40".equals(mTradeBill.transStatus)) {
            if ("09".equals(mTradeBill.response)) {
                //已关闭
                mPayResult.setText(getString(R.string.detail_activity_trade_status_closed));
            } else {
                //全额退款
                mPayResult.setText(getString(R.string.detail_activity_trade_status_partrefd));
            }
        } else {
            //失败的
            mPayResult.setText(getString(R.string.detail_activity_trade_status_fail));
            mPayResult.setTextColor(Color.RED);
            mPayResultImage.setImageResource(R.drawable.pay_result_fail);
        }


        String amount = "";
        String refd = "";
        String arriavl = "";
        String discount = "";
        try {
            BigDecimal bg0 = new BigDecimal("0");//这个是数字0

            BigDecimal txamtBD = new BigDecimal(mTradeBill.amount);//这个是txamt传来的，就是交易时给的金额

            BigDecimal refdBD = new BigDecimal(mTradeBill.refundAmt);//退款金额

            //txamt - refd == arriavl 交易的 - 退款的 == 到账的
            arriavl = txamtBD.subtract(refdBD).setScale(2, BigDecimal.ROUND_HALF_UP).toString();

            //退款金额
            refd = refdBD.setScale(2, BigDecimal.ROUND_HALF_UP).toString();

            //卡券优惠金额
            BigDecimal discountBD = new BigDecimal(mTradeBill.couponDiscountAmt);
            discount = discountBD.setScale(2, BigDecimal.ROUND_HALF_UP).toString();

            amount = txamtBD.add(discountBD).setScale(2, BigDecimal.ROUND_HALF_UP).toString();
            if (discountBD.compareTo(bg0) > 0) {
                //大于零说明有优惠金额

            }
        } catch (Exception e) {
            e.printStackTrace();
        }

        mPayMoney.setText(amount);//金额，最上面显示的那个数字
        mRefdMoney.setRightText(refd);//退款金额
        mArriavlMoney.setRightText(arriavl);//到款金额
        mCardDiscount.setRightText(discount);//卡券金额

    }


    //刷新按钮点击事件处理方法
    public void refreshOnclick(View view) {
        OrderData orderData = new OrderData();
        orderData.origOrderNum = mTradeBill.orderNum;
        startLoading(mPayResultImage);

        quickPayService.getOrderAsync(SessonData.loginUser, mTradeBill.orderNum, new QuickPayCallbackListener<ServerPacket>() {
            @Override
            public void onSuccess(ServerPacket data) {
                Log.e(TAG, "[getOrderAsync][onSuccess]  data = " + data);
                //注意这里使用了findOrder的新的的接口，这里txn返回的数组，不再是一个字段了，
                // 这时候就没必要使用新的com.cardinfolink.yunshouyin.model.ServerPacketOrder的这个类了
                Txn[] txn = data.getTxn();
                //获取txn数组，判断是否为null，且长度是否为1.这里是精确查找账单的txn返回的数组必须是1个长度的
                if (txn != null && txn.length == 1) {
                    mTradeBill.response = txn[0].getResponse();
                    mTradeBill.tandeDate = txn[0].getSystemDate();
                    mTradeBill.consumerAccount = txn[0].getConsumerAccount();
                    mTradeBill.transStatus = txn[0].getTransStatus();
                    mTradeBill.refundAmt = TxamtUtil.getNormal(txn[0].getRefundAmt());

                    QRequest req = txn[0].getmRequest();
                    if (req != null) {
                        mTradeBill.orderNum = req.getOrderNum();
                        mTradeBill.amount = TxamtUtil.getNormal(req.getTxamt());
                        mTradeBill.busicd = req.getBusicd();

                        //使用/v3/bill接口 退款的好像也没有拉取到
                        if ("REFD".equals(mTradeBill.busicd)) {
                            mTradeBill.amount = "-" + mTradeBill.amount;
                        }
                        mTradeBill.chcd = req.getChcd();
                        mTradeBill.tradeFrom = req.getTradeFrom();
                        mTradeBill.goodsInfo = req.getGoodsInfo();
                    }
                }//end if()
                endLoading(mPayResultImage);
                initData();
            }

            @Override
            public void onFailure(QuickPayException ex) {
                Log.e(TAG, "[getOrderAsync][onFailure]  ex = " + ex);
                endLoading(mPayResultImage);
                initData();
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
