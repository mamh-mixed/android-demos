package com.cardinfolink.yunshouyin.view;

import android.content.Context;
import android.content.Intent;
import android.os.Bundle;
import android.os.Handler;
import android.os.Message;
import android.util.Log;
import android.view.LayoutInflater;
import android.view.View;
import android.widget.Button;
import android.widget.EditText;
import android.widget.ImageView;
import android.widget.LinearLayout;
import android.widget.TextView;

import com.cardinfolink.cashiersdk.listener.CashierListener;
import com.cardinfolink.cashiersdk.model.OrderData;
import com.cardinfolink.cashiersdk.model.ResultData;
import com.cardinfolink.cashiersdk.sdk.CashierSdk;
import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.activity.CaptureActivity;
import com.cardinfolink.yunshouyin.activity.TicketResultActivity;
import com.cardinfolink.yunshouyin.constant.Msg;

import java.text.SimpleDateFormat;
import java.util.Date;
import java.util.Random;

/**
 * 销券的界面
 * Created by mamh on 15-12-7.
 */
public class TicketView extends LinearLayout implements View.OnClickListener {
    private static final String TAG = "TicketView";

    private Context mContext;
    private View contentView;
    private Button mConfirm;
    private ImageView mCamera;
    private EditText mCouponCode;
    private ImageView mInfo;
    private TextView mAccount;
    private Handler mHandler;
    private ResultData mResultData;

    public TicketView(Context context) {
        super(context);
        mContext = context;
        contentView = LayoutInflater.from(context).inflate(R.layout.ticket_view, null);
        LinearLayout.LayoutParams layoutParams = new LayoutParams(LayoutParams.MATCH_PARENT, LayoutParams.MATCH_PARENT);
        contentView.setLayoutParams(layoutParams);
        addView(contentView);

        initLayout();
        initHandler();
    }

    public void initLayout() {
        mConfirm = (Button) contentView.findViewById(R.id.bt_confirm);
        mCouponCode = (EditText) contentView.findViewById(R.id.et_input_coupon_code);
        mCamera = (ImageView) contentView.findViewById(R.id.iv_scan_code);
        mInfo = (ImageView) contentView.findViewById(R.id.iv_coupon_info);
        mAccount = (TextView) contentView.findViewById(R.id.tv_coupon_account);

        mConfirm.setOnClickListener(this);
        mCamera.setOnClickListener(this);
    }

    public void initHandler() {
        mHandler = new Handler() {

            @Override
            public void handleMessage(Message msg) {
                Intent intent;
                Bundle bundle;
                super.handleMessage(msg);
                switch (msg.what) {
                    case Msg.MSG_FROM_SERVER_COUPON_SUCCESS:
                        //核销成功
                        Log.e("xxxxxxxxx", "核销成功");
                        intent = new Intent(mContext, TicketResultActivity.class);
                        bundle = new Bundle();
                        bundle.putBoolean("flag", true);
                        intent.putExtras(bundle);
                        mContext.startActivity(intent);
                        break;
                    case Msg.MSG_FROM_SERVER_COUPON_FAIL:
                        //核销失败
                        // ResultData{
                        // respcd='09', busicd='VERI', chcd='ULIVE', txamt='null',
                        // channelOrderNum='null', consumerAccount='null', consumerId='null',
                        // errorDetail='null', orderNum='15122617231557530', chcdDiscount='null',
                        // merDiscount='null', qrcode='null', origOrderNum='null', scanCodeId='1810800032000019',
                        // cardId='微信支付固定金额券', cardInfo='微信支付固定金额券',
                        // voucherType='null', saleMinAmount='null', saleDiscount='null',
                        // actualPayAmount='null', maxDiscountAmt='null'}
                        Log.e("xxxxxxxxx", "核销失败");
                        intent = new Intent(mContext, TicketResultActivity.class);
                        bundle = new Bundle();
                        bundle.putBoolean("flag", false);
                        bundle.putString("cardInfo", mResultData.cardInfo);
                        bundle.putString("originalmoney", 15 + "");
                        bundle.putString("dicountmoney", 12 + "");

                        intent.putExtras(bundle);
                        mContext.startActivity(intent);
                        break;
                }
            }
        };
    }


    @Override
    public void onClick(View v) {
        switch (v.getId()) {
            case R.id.iv_scan_code:
                //扫码
                Intent intent = new Intent(mContext, CaptureActivity.class);
                Bundle bundle = new Bundle();
                bundle.putString("original", "ticketview");

                intent.putExtras(bundle);
                mContext.startActivity(intent);

                break;
            case R.id.bt_confirm:
                //核销
                String scancode = mCouponCode.getText().toString();
                Date now = new Date();
                SimpleDateFormat spf = new SimpleDateFormat("yyMMddHHmmss");
                String mOrderNum = spf.format(now);
                Random random = new Random();
                for (int i = 0; i < 5; i++) {
                    mOrderNum = mOrderNum + random.nextInt(10);
                }
                final OrderData orderData = new OrderData();
                orderData.orderNum = mOrderNum;
                orderData.scanCodeId = scancode;

                CashierSdk.startVeri(orderData, new CashierListener() {

                    @Override
                    public void onResult(ResultData resultData) {
                        Log.e("scanCode", resultData.toString());
                        mResultData = resultData;
                        if (mResultData.respcd.equals("00")) {
                            mHandler.sendEmptyMessage(Msg.MSG_FROM_SERVER_COUPON_SUCCESS);
                        } else {
                            //核销失败
                            mHandler.sendEmptyMessage(Msg.MSG_FROM_SERVER_COUPON_FAIL);
                        }
                    }

                    @Override
                    public void onError(int errorCode) {
                        Log.e("scanCode", errorCode + "");
                    }
                });
                break;
        }

    }

    /**
     * 生成账单号  时间加上一个随机数
     *
     * @return
     */
    private String geneOrderNumber() {
        String mOrderNum;

        Date now = new Date();
        SimpleDateFormat spf = new SimpleDateFormat("yyMMddHHmmss");
        mOrderNum = spf.format(now);
        Random random = new Random();//订单号末尾随机的生成一个数
        for (int i = 0; i < 5; i++) {
            mOrderNum = mOrderNum + random.nextInt(10);
        }
        return mOrderNum;
    }
}
