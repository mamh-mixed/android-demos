package com.cardinfolink.yunshouyin.ui;

import android.content.Context;
import android.content.res.TypedArray;
import android.util.AttributeSet;
import android.view.View;
import android.widget.RelativeLayout;
import android.widget.TextView;

import com.cardinfolink.yunshouyin.R;

/**
 * Created by charles on 2015/12/25.
 */
public class ResultInfoItem extends RelativeLayout {
    private TextView mLeftText;
    private TextView mRightText;

    public ResultInfoItem(Context context) {
        super(context);
        initView(context);
    }

    public ResultInfoItem(Context context, AttributeSet attrs) {
        super(context, attrs);
        initView(context);
        TypedArray typeArray = context.obtainStyledAttributes(attrs, R.styleable.SettingItemView);
        String left_text = typeArray.getString(R.styleable.SettingItemView_left_text);
        String right_text = typeArray.getString(R.styleable.SettingItemView_right_text);
        typeArray.recycle();
        mLeftText.setText(left_text);
        mRightText.setText(right_text);
    }

    public void initView(Context context) {
        View.inflate(context, R.layout.result_info_item, this);
        mLeftText = (TextView) this.findViewById(R.id.tv_left);
        mRightText = (TextView) this.findViewById(R.id.tv_right);

    }

    public void setLeftText(String str1) {
        mLeftText.setText(str1);
    }

    public void setRightText(String str2) {
        mRightText.setText(str2);
    }

    public void setTextColor(int id) {
        mLeftText.setTextColor(id);
        mRightText.setTextColor(id);
    }
}
