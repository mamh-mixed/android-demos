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
import android.widget.LinearLayout;
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
    private ImageView mDoFold;

    private View mPayInfo;

    private TextView mPayMoneyText;
    private TextView mPayMoney;

    private ResultInfoItem mCardDiscount;//卡券折扣
    private ResultInfoItem mRefdMoney;//退款金额
    private ResultInfoItem mArriavlMoney;//到账金额


    private ResultInfoItem mPayChcd;//支付渠道
    private ResultInfoItem mPayTerminator;//操作终端
    private ResultInfoItem mPayDatetime;//支付时间
    private ResultInfoItem mPayOrder;//支付订单号
    private ResultInfoItem mPayType;//支付方式
    private ResultInfoItem mPayComments;
    private LinearLayout mPayCommentsLayout;

    //显示卡券名字
    private TextView mCouponName;
    //核销订单号，这个只在卡券详情的时候显示
    private ResultInfoItem mVeriOrder;


    private ResultInfoItem mPayCheckCode;
    private ResultInfoItem mPaySmallTicket;


    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_detail);
        Intent intent = getIntent();
        Bundle billBundle = intent.getBundleExtra("BillBundle");
        mTradeBill = (TradeBill) billBundle.get("TradeBill");

        //初始化 通用的几个view，就是卡券和支付 详情都会用的公共的view在这里初始化
        initLayout();

        if (TradeBill.COUPON_TYPE.equals(mTradeBill.billType)) {
            //这里说明是卡券账单
            initCouponLayout();
            initCouponData();
        } else {
            initBillLayout();
            initBillData();
        }
    }


    private void initCouponData() {

        SimpleDateFormat spf1 = new SimpleDateFormat("yyyyMMddHHmmss");
        SimpleDateFormat spf2 = new SimpleDateFormat("yyyy-MM-dd HH:mm:ss");
        try {
            Date tandeDate = spf1.parse(mTradeBill.tandeDate);
            mPayDatetime.setRightText(spf2.format(tandeDate));
        } catch (ParseException e) {
            mPayDatetime.setRightText("");
        }


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
        mPayOrder.setRightText(mTradeBill.orderNum);//支付的订单号
        mVeriOrder.setRightText(mTradeBill.couponOrderNum);//核销订单
        mPayChcd.setRightText(mTradeBill.couponChannel);//卡券渠道
        mPayTerminator.setRightText(mTradeBill.terminalid);

        String amount = "";//原金额
        String arriavl = "";//到账金额
        String discount = "";//卡券金额
        try {
            BigDecimal bg0 = new BigDecimal("0");//这个是数字0
            BigDecimal txamtBD;
            if (TextUtils.isEmpty(mTradeBill.amount)) {
                txamtBD = new BigDecimal("0");//这个是txamt传来的，就是交易时给的金额
            } else {
                txamtBD = new BigDecimal(mTradeBill.amount);//这个是txamt传来的，就是交易时给的金额
            }
            arriavl = txamtBD.setScale(2, BigDecimal.ROUND_HALF_UP).toString();

            //卡券优惠金额
            BigDecimal discountBD;
            if (TextUtils.isEmpty(mTradeBill.couponDiscountAmt)) {
                //卡券优惠是空的时候
                discountBD = new BigDecimal("0.00");
            } else {
                discountBD = new BigDecimal(mTradeBill.couponDiscountAmt);
            }
            discount = discountBD.setScale(2, BigDecimal.ROUND_HALF_UP).toString();

            amount = txamtBD.add(discountBD).setScale(2, BigDecimal.ROUND_HALF_UP).toString();
            if (discountBD.compareTo(bg0) <= 0) {
                //大于零说明有优惠金额
                mCardDiscount.setVisibility(View.GONE);
            }
            if (txamtBD.compareTo(bg0) <= 0) {
                mPayMoneyText.setVisibility(View.GONE);
                mPayMoney.setVisibility(View.GONE);
                mDoFold.setVisibility(View.GONE);
                mPayInfo.setVisibility(View.GONE);
            }
        } catch (Exception e) {
        }

        mPayMoney.setText(amount);
        mCardDiscount.setRightText("-" + discount);
        mArriavlMoney.setRightText(arriavl);
    }

    /**
     * 公共的 布局view放在这里
     */
    private void initLayout() {
        mActionBar = (SettingActionBarItem) findViewById(R.id.action_bar);
        mActionBar.setLeftTextOnclickListner(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                finish();
            }
        });

        mCouponName = (TextView) findViewById(R.id.tv_coupon_name);
        mPayResult = (TextView) findViewById(R.id.tv_pay_result);
        mPayResultImage = (ImageView) findViewById(R.id.iv_pay_result);
        mDoFold = (ImageView) findViewById(R.id.iv_dofold);
        mDoFold.setOnClickListener(new View.OnClickListener() {
            private boolean isFold = false;

            @Override
            public void onClick(View v) {
                if (isFold) {//如果是 折叠，这里让其显示
                    mPayInfo.setVisibility(View.VISIBLE);
                    mDoFold.setImageResource(R.drawable.bill_up);
                } else {
                    mPayInfo.setVisibility(View.GONE);
                    mDoFold.setImageResource(R.drawable.bill_down);
                }
                isFold = !isFold;
            }
        });

        mPayMoneyText = (TextView) findViewById(R.id.pay_money_text);
        mPayMoney = (TextView) findViewById(R.id.pay_money);//到账金额

        mCardDiscount = (ResultInfoItem) findViewById(R.id.card_discount);//卡券折扣
        mRefdMoney = (ResultInfoItem) findViewById(R.id.refd_money);
        mArriavlMoney = (ResultInfoItem) findViewById(R.id.pay_arrival_money);

        mCardDiscount.setTextColor(R.color.gray3);
        mRefdMoney.setTextColor(R.color.gray3);
        mArriavlMoney.setTextColor(R.color.gray3);

        mPayInfo = findViewById(R.id.ll_pay_info);

        //核销订单号，这个只在卡券详情的时候显示
        mVeriOrder = (ResultInfoItem) findViewById(R.id.veri_order);
        //六个
        mPayOrder = (ResultInfoItem) findViewById(R.id.pay_order);
        mPayChcd = (ResultInfoItem) findViewById(R.id.pay_chcd);
        mPayType = (ResultInfoItem) findViewById(R.id.pay_type);
        mPayDatetime = (ResultInfoItem) findViewById(R.id.pay_datetime);
        mPayTerminator = (ResultInfoItem) findViewById(R.id.pay_terminator);
        mPayComments = (ResultInfoItem) findViewById(R.id.pay_comments);

        mPayCheckCode = (ResultInfoItem) findViewById(R.id.pay_check_code);
        mPaySmallTicket = (ResultInfoItem) findViewById(R.id.pay_small_ticket_number);
        mPayCheckCode.setTextColor(R.color.gray3);
        mPaySmallTicket.setTextColor(R.color.gray3);


        mPayCommentsLayout = (LinearLayout) findViewById(R.id.ll_pay_comments);
    }

    private void initCouponLayout() {
        mActionBar.setRightText("");

        mPayCommentsLayout.setVisibility(View.GONE);

        mPayResult.setText(getString(R.string.detail_activity_veri_success));

        mCouponName.setVisibility(View.VISIBLE);
        mCouponName.setText(mTradeBill.couponName);

        mPayMoneyText.setText(getString(R.string.detail_activity_pay_money1));
        mCardDiscount.setLeftText(getString(R.string.detail_activity_coupone_dikou));
        mRefdMoney.setVisibility(View.GONE);

        //核销订单号，这个只在卡券详情的时候显示
        mVeriOrder.setVisibility(View.VISIBLE);
        //六个
        mPayOrder.setLeftText(getString(R.string.detail_activity_order_number1));
        mPayChcd.setLeftText(getString(R.string.detail_activity_coupon_chcd));
        mPayType.setLeftText(getString(R.string.detail_activity_veri_type));
        mPayDatetime.setLeftText(getString(R.string.detail_activity_veri_datetime));
        mPayTerminator.setLeftText(getString(R.string.detail_activity_terminator));
        mPayComments.setVisibility(View.GONE);
    }

    private void initBillLayout() {
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
                    finish();
                }
            }
        });


    }

    private void initBillData() {
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
        if (TextUtils.isEmpty(mTradeBill.terminalid)) {
            mPayTerminator.setVisibility(View.GONE);
        } else {
            mPayTerminator.setVisibility(View.VISIBLE);
            mPayTerminator.setRightText(mTradeBill.terminalid);
        }

        //检验码
        if (TextUtils.isEmpty(mTradeBill.checkCode)) {
            mPayCheckCode.setVisibility(View.GONE);
        } else {
            mPayCheckCode.setVisibility(View.VISIBLE);
            mPayCheckCode.setRightText(mTradeBill.checkCode);
        }
        //小票号
        if (TextUtils.isEmpty(mTradeBill.smallTicketNumber)) {
            mPaySmallTicket.setVisibility(View.GONE);
        } else {
            mPaySmallTicket.setVisibility(View.VISIBLE);
            mPaySmallTicket.setRightText(mTradeBill.smallTicketNumber);
        }
        //如果都为空，就都隐藏了
        if (TextUtils.isEmpty(mTradeBill.checkCode) && TextUtils.isEmpty(mTradeBill.smallTicketNumber)) {
            mPayCommentsLayout.setVisibility(View.GONE);
            mPayComments.setVisibility(View.GONE);
        }

        //支付方式
        if (!TextUtils.isEmpty(mTradeBill.tradeFrom)) {
            if ("android".equals(mTradeBill.tradeFrom)) {
                //app收款
                mPayType.setRightText(getString(R.string.detail_activity_pay_type2));

            } else if ("ios".equals(mTradeBill.tradeFrom)) {
                //app收款
                mPayType.setRightText(getString(R.string.detail_activity_pay_type3));
            } else if ("wap".equals(mTradeBill.tradeFrom)) {
                //web 收款
                mPayType.setRightText(getString(R.string.detail_activity_pay_type4));
            } else {
                //其他收款
                mPayType.setRightText(getString(R.string.detail_activity_pay_type1));
            }
        } else {
            //tradeFrom 是 null 清空
            mPayType.setRightText(getString(R.string.detail_activity_pay_type5));
        }

        //设置订单号
        mPayOrder.setRightText(mTradeBill.orderNum);

        //设置交易的结果状态
        if ("10".equals(mTradeBill.transStatus)) {
            //处理中
            mPayResult.setText(getString(R.string.detail_activity_trade_status_nopay));
            mPayResultImage.setVisibility(View.VISIBLE);
            mPayResultImage.setImageResource(R.drawable.bill_fresh);
            mPayResultImage.setOnClickListener(new View.OnClickListener() {
                @Override
                public void onClick(View v) {
                    refreshOnclick();
                }
            });
        } else if ("30".equals(mTradeBill.transStatus)) {
            try {
                BigDecimal refdBD = new BigDecimal(mTradeBill.refundAmt);//退款金额
                if (refdBD.compareTo(new BigDecimal("0")) == 0) {
                    //成功的
                    mPayResult.setTextColor(getResources().getColor(R.color.textview_textcolor_pay_success));
                    mPayResult.setText(getString(R.string.detail_activity_trade_status_success));
                    mPayResultImage.setVisibility(View.GONE);
                    mPayResultImage.setImageResource(R.drawable.pay_result_succeed);
                } else {
                    //部分退款的
                    mPayResult.setText(getString(R.string.detail_activity_trade_status_partrefd));
                }
            } catch (Exception e) {
                mPayResult.setText(getString(R.string.detail_activity_trade_status_partrefd));
            }
        } else if ("40".equals(mTradeBill.transStatus)) {
            if ("09".equals(mTradeBill.response)) {
                //已关闭
                mPayResult.setText(getString(R.string.detail_activity_trade_status_closed));
            } else {
                //全额退款
                mPayResult.setText(getString(R.string.detail_activity_trade_status_allrefd));
            }
        } else {
            //失败的
            mPayResult.setText(getString(R.string.detail_activity_trade_status_fail));
            mPayResult.setTextColor(Color.RED);
            mPayResultImage.setVisibility(View.GONE);
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
            BigDecimal discountBD;
            if (TextUtils.isEmpty(mTradeBill.couponDiscountAmt)) {
                //卡券优惠是空的时候
                discountBD = new BigDecimal("0.00");
            } else {
                discountBD = new BigDecimal(mTradeBill.couponDiscountAmt);

            }
            discount = discountBD.setScale(2, BigDecimal.ROUND_HALF_UP).toString();

            amount = txamtBD.add(discountBD).setScale(2, BigDecimal.ROUND_HALF_UP).toString();
            if (discountBD.compareTo(bg0) <= 0) {
                //没有优惠的时候不显示，这个
                mCardDiscount.setVisibility(View.GONE);
            }
            if (refdBD.compareTo(bg0) <= 0) {
                mRefdMoney.setVisibility(View.GONE);
            }

            //"当没有退款和折扣时，收款金额需要显示为原金额；当有退款或者折扣时，收款金额需要显示为金额"
            if (discountBD.compareTo(bg0) <= 0 && refdBD.compareTo(bg0) <= 0) {
                mPayMoneyText.setText(R.string.detail_activity_pay_money);//金额
                mPayInfo.setVisibility(View.GONE);
                mDoFold.setVisibility(View.GONE);
            } else {
                mPayMoneyText.setText(R.string.detail_activity_pay_money1);//原金额
                mPayInfo.setVisibility(View.VISIBLE);
                mDoFold.setVisibility(View.VISIBLE);
            }
        } catch (Exception e) {

        }

        mPayMoney.setText(amount);//金额，最上面显示的那个数字
        mRefdMoney.setRightText("-" + refd);//退款金额
        mArriavlMoney.setRightText(arriavl);//到款金额
        mCardDiscount.setRightText("-" + discount);//卡券金额

    }


    //刷新按钮点击事件处理方法
    public void refreshOnclick() {
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
                initBillData();
            }

            @Override
            public void onFailure(QuickPayException ex) {
                Log.e(TAG, "[getOrderAsync][onFailure]  ex = " + ex);
                endLoading(mPayResultImage);
                initBillData();
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
