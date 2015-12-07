package com.cardinfolink.yunshouyin.view;

import android.content.Context;
import android.content.Intent;
import android.view.LayoutInflater;
import android.view.View;
import android.widget.LinearLayout;

import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.activity.WapActivity;
import com.cardinfolink.yunshouyin.ui.SettingClikcItem;

/**
 * 第四个界面，就是设置界面
 * Created by mamh on 15-12-7.
 */
public class MySettingView extends LinearLayout {
    private static final String TAG = "MySettingView";

    private Context mContext;

    private SettingClikcItem mMyWap;

    public MySettingView(Context context) {
        super(context);
        mContext = context;
        View contentView = LayoutInflater.from(context).inflate(R.layout.my_setting_view, null);
        LinearLayout.LayoutParams layoutParams = new LinearLayout.LayoutParams(LinearLayout.LayoutParams.MATCH_PARENT, LinearLayout.LayoutParams.MATCH_PARENT);
        contentView.setLayoutParams(layoutParams);
        addView(contentView);

        mMyWap = (SettingClikcItem) contentView.findViewById(R.id.my_wap);
        mMyWap.setOnClickListener(new OnClickListener() {
            @Override
            public void onClick(View v) {
                Intent intent = new Intent(mContext, WapActivity.class);
                mContext.startActivity(intent);
            }
        });
    }
}
