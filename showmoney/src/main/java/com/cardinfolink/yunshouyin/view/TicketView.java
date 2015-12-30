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
import com.cardinfolink.yunshouyin.activity.CouponResultActivity;
import com.cardinfolink.yunshouyin.activity.TicketResultActivity;
import com.cardinfolink.yunshouyin.constant.Msg;
import com.cardinfolink.yunshouyin.data.SessonData;

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

    private Handler mMainActivityHandler;//这个是mainActivity类里面的handler

    private ResultData mResultData;

    private Handler mHandler;//这个是本类里面自有的handler


    public TicketView(Context context) {
        this(context, null);
    }

    //添加了一个新的构造方法
    public TicketView(Context context, Handler handler) {
        super(context);
        mContext = context;
        mMainActivityHandler = handler;

        contentView = LayoutInflater.from(context).inflate(R.layout.ticket_view, null);
        LinearLayout.LayoutParams layoutParams = new LayoutParams(LayoutParams.MATCH_PARENT, LayoutParams.MATCH_PARENT);
        contentView.setLayoutParams(layoutParams);
        addView(contentView);

        initLayout();
        initHandler();//初始化handler
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
                        intent = new Intent(mContext, CouponResultActivity.class);
                        bundle = new Bundle();
                        bundle.putBoolean("check_coupon_result_flag", true);
                        intent.putExtras(bundle);
                        mContext.startActivity(intent);
                        break;
                    case Msg.MSG_FROM_SERVER_COUPON_FAIL:
                        //核销失败
                        intent = new Intent(mContext, CouponResultActivity.class);
                        bundle = new Bundle();
                        bundle.putBoolean("check_coupon_result_flag", false);
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
                final OrderData orderData = new OrderData();
                orderData.orderNum = geneOrderNumber();
                orderData.scanCodeId = scancode;

                CashierSdk.startVeri(orderData, new CashierListener() {

                    @Override
                    public void onResult(ResultData resultData) {
                        mResultData = resultData;
                        SessonData.loginUser.setResultData(resultData);
                        if ("00".equals(mResultData.respcd)) {
                            mHandler.sendEmptyMessage(Msg.MSG_FROM_SERVER_COUPON_SUCCESS);
                        } else {
                            //核销失败
                            mHandler.sendEmptyMessage(Msg.MSG_FROM_SERVER_COUPON_FAIL);
                        }
                    }

                    @Override
                    public void onError(int errorCode) {
                        mHandler.sendEmptyMessage(Msg.MSG_FROM_SERVER_COUPON_FAIL);
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
