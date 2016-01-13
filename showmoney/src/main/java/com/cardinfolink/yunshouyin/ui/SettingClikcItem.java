package com.cardinfolink.yunshouyin.ui;

import android.content.Context;
import android.content.res.TypedArray;
import android.graphics.drawable.Drawable;
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
public class SettingClikcItem extends RelativeLayout {
    protected ImageView mImageView;
    protected TextView mTitle;
    protected TextView mRightText;


    public SettingClikcItem(Context context) {
        super(context);
        initView(context);
    }

    public SettingClikcItem(Context context, AttributeSet attrs) {
        super(context, attrs);
        initView(context);

        TypedArray typeArray = context.obtainStyledAttributes(attrs, R.styleable.SettingItemView);
        String title = typeArray.getString(R.styleable.SettingItemView_title);
        String right = typeArray.getString(R.styleable.SettingItemView_right_text);
        Drawable drawable = typeArray.getDrawable(R.styleable.SettingItemView_android_src);
        typeArray.recycle();

        setRightText(right);
        setTitle(title);
        setmImageView(drawable);
    }

    public SettingClikcItem(Context context, AttributeSet attrs, int defStyle) {
        super(context, attrs, defStyle);
        initView(context);
    }

    protected void initView(Context context) {
        View.inflate(context, R.layout.setting_click_item, this);
        mImageView = (ImageView) this.findViewById(R.id.iv_setting);
        mTitle = (TextView) this.findViewById(R.id.tv_title);
        mRightText = (TextView) this.findViewById(R.id.tv_right);
    }


    public void setTitle(String title) {
        if (mTitle != null) {
            mTitle.setText(title);
        }
    }

    public String getTitle() {
        return mTitle.getText().toString();
    }

    public void setRightText(String text) {
        mRightText.setText(text);
    }

    public String getRightText() {
        return mRightText.getText().toString();
    }

    public void setImageResource(int id) {
        mImageView.setImageResource(id);
    }

    public void setmImageView(Drawable drawable) {
        mImageView.setImageDrawable(drawable);
    }

}
