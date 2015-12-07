package com.cardinfolink.yunshouyin.ui;

import android.content.Context;
import android.content.res.TypedArray;
import android.util.AttributeSet;
import android.view.View;
import android.widget.EditText;
import android.widget.ImageView;
import android.widget.RelativeLayout;
import android.widget.TextView;

import com.cardinfolink.yunshouyin.R;


/**
 * Created by mamh on 15-11-1.
 * 自定义的组合控件
 */
public class SettingPasswordItem extends RelativeLayout {
    private ImageView mImageView;
    private TextView mTitle;
    private EditText mPassword;


    public SettingPasswordItem(Context context) {
        super(context);
        initView(context);
    }

    public SettingPasswordItem(Context context, AttributeSet attrs) {
        super(context, attrs);
        initView(context);
        TypedArray typeArray = context.obtainStyledAttributes(attrs, R.styleable.SettingItemView);
        String title = typeArray.getString(R.styleable.SettingItemView_title);
        String hint = typeArray.getString(R.styleable.SettingItemView_password_hint);
        typeArray.recycle();
        mPassword.setHint(hint);
        mTitle.setText(title);
    }

    public SettingPasswordItem(Context context, AttributeSet attrs, int defStyle) {
        super(context, attrs, defStyle);
        initView(context);
    }

    private void initView(Context context) {
        View.inflate(context, R.layout.setting_password_item, this);
        mImageView = (ImageView) this.findViewById(R.id.iv_show);
        mTitle = (TextView) this.findViewById(R.id.tv_title);
        mPassword = (EditText) this.findViewById(R.id.et_password);
    }


    public void setTitle(String title) {
        if (mTitle != null) {
            mTitle.setText(title);
        }
    }

    public String getTitle() {
        return mTitle.getText().toString();
    }

    public void setPassword(String password) {
        if (mPassword != null) {
            mPassword.setText(password);
        }
    }

    public String getPassword() {
        return mPassword.getText().toString();
    }
}
