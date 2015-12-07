package com.cardinfolink.yunshouyin.ui;

import android.content.Context;
import android.content.res.TypedArray;
import android.util.AttributeSet;
import android.view.View;
import android.widget.ImageView;
import android.widget.RelativeLayout;
import android.widget.TextView;

import com.cardinfolink.yunshouyin.R;


/**
 * Created by mamh on 15-11-1.
 * 自定义的组合控件
 */
public class SettingClikcView extends RelativeLayout {
    private ImageView mImageView;
    private TextView mTitle;


    public SettingClikcView(Context context) {
        super(context);
        initView(context);
    }

    public SettingClikcView(Context context, AttributeSet attrs) {
        super(context, attrs);
        initView(context);
        TypedArray typeArray = context.obtainStyledAttributes(attrs, R.styleable.SettingItemView);
        String title = typeArray.getString(R.styleable.SettingItemView_title);
        typeArray.recycle();
        mTitle.setText(title);
    }

    public SettingClikcView(Context context, AttributeSet attrs, int defStyle) {
        super(context, attrs, defStyle);
        initView(context);
    }

    private void initView(Context context) {
        View.inflate(context, R.layout.setting_click_item, this);
        mImageView = (ImageView) this.findViewById(R.id.iv_setting);
        mTitle = (TextView) this.findViewById(R.id.tv_title);
    }


    public void setTitle(String title) {
        if (mTitle != null) {
            mTitle.setText(title);
        }
    }

    public String getTitle() {
        return mTitle.getText().toString();
    }

}
