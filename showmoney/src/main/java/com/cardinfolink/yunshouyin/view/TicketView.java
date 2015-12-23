package com.cardinfolink.yunshouyin.view;

import android.content.Context;
import android.content.Intent;
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

    public TicketView(Context context) {
        super(context);
        mContext = context;
        contentView = LayoutInflater.from(context).inflate(R.layout.ticket_view, null);
        LinearLayout.LayoutParams layoutParams = new LayoutParams(LayoutParams.MATCH_PARENT, LayoutParams.MATCH_PARENT);
        contentView.setLayoutParams(layoutParams);
        addView(contentView);
        initLayout();

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


    @Override
    public void onClick(View v) {
        switch (v.getId()) {
            case R.id.iv_scan_code:
                //扫码
                Intent intent = new Intent(mContext, CaptureActivity.class);
                mContext.startActivity(intent);
                break;
            case R.id.bt_confirm:
                //核销
                OrderData orderData = new OrderData();
                CashierSdk.startVeri(orderData, new CashierListener() {
                    @Override
                    public void onResult(ResultData resultData) {

                    }

                    @Override
                    public void onError(int errorCode) {

                    }
                });
                break;
        }

    }
}
