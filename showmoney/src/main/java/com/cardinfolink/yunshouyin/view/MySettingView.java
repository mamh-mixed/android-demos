package com.cardinfolink.yunshouyin.view;

import android.app.Activity;
import android.content.Context;
import android.content.Intent;
import android.view.LayoutInflater;
import android.view.View;
import android.widget.Button;
import android.widget.ImageView;
import android.widget.LinearLayout;
import android.widget.TextView;

import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.activity.AboutActivity;
import com.cardinfolink.yunshouyin.activity.AccountSecurityActivity;
import com.cardinfolink.yunshouyin.activity.LimitIncreaseActivity;
import com.cardinfolink.yunshouyin.activity.LoginActivity;
import com.cardinfolink.yunshouyin.activity.MyChannelActivity;
import com.cardinfolink.yunshouyin.activity.UnReadMessageActivity;
import com.cardinfolink.yunshouyin.activity.WapActivity;
import com.cardinfolink.yunshouyin.api.QuickPayException;
import com.cardinfolink.yunshouyin.core.QuickPayCallbackListener;
import com.cardinfolink.yunshouyin.core.QuickPayService;
import com.cardinfolink.yunshouyin.data.SessonData;
import com.cardinfolink.yunshouyin.data.User;
import com.cardinfolink.yunshouyin.ui.SettingClikcItem;
import com.cardinfolink.yunshouyin.util.ShowMoneyApp;

import java.text.SimpleDateFormat;
import java.util.Date;

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

    private TextView mEmail;//账户名
    private TextView mLimit;//显示限额的一些信息的

    private ImageView mMessage;

    public MySettingView(Context context) {
        super(context);
        mContext = context;
        View contentView = LayoutInflater.from(context).inflate(R.layout.my_setting_view, null);
        LinearLayout.LayoutParams layoutParams = new LinearLayout.LayoutParams(LinearLayout.LayoutParams.MATCH_PARENT, LinearLayout.LayoutParams.MATCH_PARENT);
        contentView.setLayoutParams(layoutParams);
        addView(contentView);

        mExit = (Button) contentView.findViewById(R.id.btn_exit);
        mIncreaseLimit = (Button) contentView.findViewById(R.id.btn_limit);//只有限额的用户会显示这个按钮

        mAccountAndSecurity = (SettingClikcItem) contentView.findViewById(R.id.account_security);
        mSupportChannel = (SettingClikcItem) contentView.findViewById(R.id.support_channel);
        mMyWap = (SettingClikcItem) contentView.findViewById(R.id.my_wap);
        mAbout = (SettingClikcItem) contentView.findViewById(R.id.about);

        mEmail = (TextView) contentView.findViewById(R.id.tv_email);//账户名
        mEmail.setText(SessonData.loginUser.getUsername());//通过sessonData设置一下用户名

        mLimit = (TextView) contentView.findViewById(R.id.tv_limit_info);//显示限额的一些信息的
        mMessage = (ImageView) contentView.findViewById(R.id.iv_message);//右上角显示是否有未读消息的图片

        mExit.setOnClickListener(this);
        mIncreaseLimit.setOnClickListener(this);
        mAccountAndSecurity.setOnClickListener(this);
        mSupportChannel.setOnClickListener(this);
        mMyWap.setOnClickListener(this);
        mAbout.setOnClickListener(this);
        mMessage.setOnClickListener(this);

        checkLimit();//发送http请求来检查当日的限额数。
    }

    private void checkLimit() {
        QuickPayService quickPayService = ShowMoneyApp.getInstance().getQuickPayService();
        String date = (new SimpleDateFormat("yyyyMMdd")).format(new Date());
        User user = SessonData.loginUser;
        if (user.getLimit().equals("true")) {//这里等于true表示这个用户有限额。
            quickPayService.getTotalAsync(user, date, new QuickPayCallbackListener<String>() {
                @Override
                public void onSuccess(String data) {
                    //如果有限额的话
                    double limit = Double.parseDouble(data);
                    if (limit > 0) {//大于零表示有限额了
                        String limitMsg = getResources().getString(R.string.setting_limit_message);
                        limitMsg = String.format(limitMsg, data);
                        mLimit.setText(limitMsg);//这里设置限额多少的提示文本
                        mIncreaseLimit.setVisibility(VISIBLE);//把提升限额的按钮显示出来
                    }else{
                        //这里表示没有限额
                    }

                }

                @Override
                public void onFailure(QuickPayException ex) {
                    mLimit.setText(ex.getErrorMsg());
                }
            });
        }else{
            //else这里表示用户没有限额
        }
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
            case R.id.iv_message:
                //这里要加个判断 是否有 未读消息，有的话跳转未读消息界面。
                intent = new Intent(mContext, UnReadMessageActivity.class);
                mContext.startActivity(intent);
                break;
        }
    }

}
