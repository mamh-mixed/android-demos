package com.cardinfolink.yunshouyin.ui;

import android.content.Context;
import android.content.res.TypedArray;
import android.text.Editable;
import android.text.InputType;
import android.text.TextUtils;
import android.text.TextWatcher;
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
public class SettingPasswordItem extends RelativeLayout implements View.OnClickListener, TextWatcher {
    private ImageView mImageView;
    private TextView mTitle;
    private EditTextClear mPassword;

    private boolean mPassowrdIsVisible = false;//初始时 密码是不显示的


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
        mPassword = (EditTextClear) this.findViewById(R.id.et_password);
        mPassword.addTextChangedListener(this);//这里设置一些密码输入框变化的监听事件
        mImageView.setOnClickListener(this);//这里设置图片的点击事件
        if (TextUtils.isEmpty(getPassword())) {
            mImageView.setVisibility(View.INVISIBLE);
        }
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

    public void setPasswordIconVisible(boolean visible) {
        mPassword.setClearIconVisible(visible);
    }

    public void setShakeAnimation() {
        mPassword.setShakeAnimation();
    }

    public void setImageViewOnClickListener(OnClickListener l) {
        mImageView.setOnClickListener(l);
    }


    @Override
    public void onClick(View v) {
        switch (v.getId()) {
            case R.id.iv_show:
                //判断来设置密码输入框是否显示密码还是隐藏密码
                if (mPassowrdIsVisible) {
                    // 隐藏密码
                    mPassword.setInputType(InputType.TYPE_CLASS_TEXT | InputType.TYPE_TEXT_VARIATION_PASSWORD);
                } else {
                    //显示
                    mPassword.setInputType(InputType.TYPE_TEXT_VARIATION_VISIBLE_PASSWORD);
                }
                mPassowrdIsVisible = !mPassowrdIsVisible;
                break;
            case R.id.tv_title:
                break;
        }//end switch()
    }


    @Override
    public void beforeTextChanged(CharSequence s, int start, int count, int after) {

    }

    @Override
    public void onTextChanged(CharSequence s, int start, int before, int count) {

    }

    @Override
    public void afterTextChanged(Editable s) {
        //密码框输入变化后调用，如果密码是空的就不显示此图片了。
        if (TextUtils.isEmpty(getPassword())) {
            mImageView.setVisibility(INVISIBLE);
        } else {
            mImageView.setVisibility(VISIBLE);
        }
    }
}
