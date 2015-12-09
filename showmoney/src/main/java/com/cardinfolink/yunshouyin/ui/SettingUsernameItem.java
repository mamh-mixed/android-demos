package com.cardinfolink.yunshouyin.ui;

import android.content.Context;
import android.content.res.TypedArray;
import android.util.AttributeSet;
import android.view.View;
import android.widget.RelativeLayout;
import android.widget.TextView;

import com.cardinfolink.yunshouyin.R;


/**
 * Created by mamh on 15-11-1.
 * 自定义的组合控件,用来输入文本的和输入密码的多少不太一样
 */
public class SettingUsernameItem extends RelativeLayout {
    private TextView mTitle;
    private EditTextClear mUsername;

    public SettingUsernameItem(Context context) {
        super(context);
        initView(context);
    }

    public SettingUsernameItem(Context context, AttributeSet attrs) {
        super(context, attrs);
        initView(context);
        TypedArray typeArray = context.obtainStyledAttributes(attrs, R.styleable.SettingItemView);
        String title = typeArray.getString(R.styleable.SettingItemView_title);
        String hint = typeArray.getString(R.styleable.SettingItemView_username_hint);
        typeArray.recycle();
        mUsername.setHint(hint);
        mTitle.setText(title);
    }

    public SettingUsernameItem(Context context, AttributeSet attrs, int defStyle) {
        super(context, attrs, defStyle);
        initView(context);
    }

    private void initView(Context context) {
        View.inflate(context, R.layout.setting_username_item, this);
        mTitle = (TextView) this.findViewById(R.id.tv_title);
        mUsername = (EditTextClear) this.findViewById(R.id.et_username);
    }


    public void setTitle(String title) {
        if (mTitle != null) {
            mTitle.setText(title);
        }
    }

    public String getTitle() {
        return mTitle.getText().toString();
    }

    public void setUsername(String username) {
        if (mUsername != null) {
            mUsername.setText(username);
        }
    }

    public String getUsername() {
        return mUsername.getText().toString();
    }

    public void setUsernameIconVisible(boolean visible) {
        mUsername.setClearIconVisible(visible);
    }

    public void setShakeAnimation() {
        mUsername.setShakeAnimation();
    }


}
