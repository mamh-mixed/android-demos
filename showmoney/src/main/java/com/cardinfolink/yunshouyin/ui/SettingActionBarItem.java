package com.cardinfolink.yunshouyin.ui;

import android.content.Context;
import android.content.res.TypedArray;
import android.graphics.Color;
import android.util.AttributeSet;
import android.view.View;
import android.widget.RelativeLayout;
import android.widget.TextView;

import com.cardinfolink.yunshouyin.R;


/**
 * Created by mamh on 15-11-1.
 * 自定义的组合控件
 */
public class SettingActionBarItem extends RelativeLayout {
    private TextView mLeftText;
    private TextView mTitle;
    private TextView mRightText;
    private RelativeLayout relativeLayout;


    public SettingActionBarItem(Context context) {
        super(context);
        initView(context);
    }

    public SettingActionBarItem(Context context, AttributeSet attrs) {
        super(context, attrs);
        initView(context);
        TypedArray typeArray = context.obtainStyledAttributes(attrs, R.styleable.SettingItemView);
        String title = typeArray.getString(R.styleable.SettingItemView_title);
        String left = typeArray.getString(R.styleable.SettingItemView_left_text);
        String right = typeArray.getString(R.styleable.SettingItemView_right_text);
        String leftTextColor = typeArray.getString(R.styleable.SettingItemView_left_text_color);
        typeArray.recycle();
        mTitle.setText(title);
        mLeftText.setText(left);
        if (leftTextColor != null) {
            mLeftText.setTextColor(Color.parseColor(leftTextColor));
        }
        mRightText.setText(right);
    }

    public SettingActionBarItem(Context context, AttributeSet attrs, int defStyle) {
        super(context, attrs, defStyle);
        initView(context);
    }

    private void initView(Context context) {
        View.inflate(context, R.layout.setting_action_bar_item, this);
        mTitle = (TextView) this.findViewById(R.id.tv_title);
        mLeftText = (TextView) this.findViewById(R.id.tv_left);
        mRightText = (TextView) this.findViewById(R.id.tv_right);
        relativeLayout = (RelativeLayout) this.findViewById(R.id.rl_action_bar);
    }


    public void setTitle(String title) {
        if (mTitle != null) {
            mTitle.setText(title);
        }
    }

    public void setLeftText(String str) {
        mLeftText.setText(str);
    }

    public void setRightText(String str) {
        mRightText.setText(str);
    }

    public String getLeftText() {
        return mLeftText.getText().toString();
    }

    public String getRightText() {
        return mRightText.getText().toString();
    }

    public void setBackgroundColor(int color) {
        relativeLayout.setBackgroundColor(color);
        mLeftText.setBackgroundColor(color);
        mRightText.setBackgroundColor(color);
        mTitle.setBackgroundColor(color);
    }

    public void setTitleColor(int color) {
        mTitle.setTextColor(color);
    }

    public void setLeftTextColor(int color) {
        mLeftText.setTextColor(color);
    }


    public String getTitle() {
        return mTitle.getText().toString();
    }

    public void setLeftTextOnclickListner(OnClickListener l) {
        mLeftText.setOnClickListener(l);
    }

    public void setRightTextOnclickListner(OnClickListener l) {
        mRightText.setOnClickListener(l);
    }

    public void setLeftTextVisibility(int visibility) {
        mLeftText.setVisibility(visibility);
    }

    public void setRightTextVisibility(int visibility) {
        mRightText.setVisibility(visibility);
    }

    public void setTitleVisibility(int visibility) {
        mTitle.setVisibility(visibility);
    }
}
