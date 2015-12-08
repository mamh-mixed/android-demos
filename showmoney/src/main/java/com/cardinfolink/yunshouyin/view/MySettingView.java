package com.cardinfolink.yunshouyin.view;

import android.app.Activity;
import android.content.Context;
import android.content.Intent;
import android.view.LayoutInflater;
import android.view.View;
import android.widget.Button;
import android.widget.LinearLayout;

import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.activity.AboutActivity;
import com.cardinfolink.yunshouyin.activity.AccountSecurityActivity;
import com.cardinfolink.yunshouyin.activity.LimitIncreaseActivity;
import com.cardinfolink.yunshouyin.activity.LoginActivity;
import com.cardinfolink.yunshouyin.activity.MyChannelActivity;
import com.cardinfolink.yunshouyin.activity.WapActivity;
import com.cardinfolink.yunshouyin.ui.SettingClikcItem;

/**
 * 第四个界面，就是设置界面
 * Created by mamh on 15-12-7.
 */
public class MySettingView extends LinearLayout implements View.OnClickListener {
    private static final String TAG = "MySettingView";

    private Context mContext;

    private SettingClikcItem mAccountAndSecurity;//账户与安全
    private SettingClikcItem mSupportChannel;//支持的渠道
    private SettingClikcItem mMyWap;//我的网页版
    private SettingClikcItem mAbout;//关于云收银

    private Button mExit;
    private Button mIncreaseLimit;//提升限额

    public MySettingView(Context context) {
        super(context);
        mContext = context;
        View contentView = LayoutInflater.from(context).inflate(R.layout.my_setting_view, null);
        LinearLayout.LayoutParams layoutParams = new LinearLayout.LayoutParams(LinearLayout.LayoutParams.MATCH_PARENT, LinearLayout.LayoutParams.MATCH_PARENT);
        contentView.setLayoutParams(layoutParams);
        addView(contentView);

        mExit = (Button) contentView.findViewById(R.id.btn_exit);
        mIncreaseLimit = (Button) contentView.findViewById(R.id.btn_limit);

        mAccountAndSecurity = (SettingClikcItem) contentView.findViewById(R.id.account_security);
        mSupportChannel = (SettingClikcItem) contentView.findViewById(R.id.support_channel);
        mMyWap = (SettingClikcItem) contentView.findViewById(R.id.my_wap);
        mAbout = (SettingClikcItem) contentView.findViewById(R.id.about);

        mExit.setOnClickListener(this);
        mIncreaseLimit.setOnClickListener(this);
        mAccountAndSecurity.setOnClickListener(this);
        mSupportChannel.setOnClickListener(this);
        mMyWap.setOnClickListener(this);
        mAbout.setOnClickListener(this);

    }


    @Override
    public void onClick(View v) {
        Intent intent = null;
        switch (v.getId()) {
            case R.id.btn_exit:
                intent = new Intent(mContext, LoginActivity.class);
                mContext.startActivity(intent);
                ((Activity) mContext).finish();
                break;
            case R.id.btn_limit:
                intent = new Intent(mContext, LimitIncreaseActivity.class);
                mContext.startActivity(intent);
                break;
            case R.id.account_security:
                //账户与安全
                intent = new Intent(mContext, AccountSecurityActivity.class);
                mContext.startActivity(intent);
                break;
            case R.id.support_channel:
                //支持的渠道
                intent = new Intent(mContext, MyChannelActivity.class);
                mContext.startActivity(intent);
                break;
            case R.id.my_wap:
                //我的网页版
                intent = new Intent(mContext, WapActivity.class);
                mContext.startActivity(intent);
                break;
            case R.id.about:
                //关于云收银
                intent = new Intent(mContext, AboutActivity.class);
                mContext.startActivity(intent);
                break;
        }
    }

}
