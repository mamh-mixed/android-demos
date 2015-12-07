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
public class SettingDetailView extends RelativeLayout {
    private TextView mTitle;
    private TextView mDetail;


    public SettingDetailView(Context context) {
        super(context);
        initView(context);
    }

    public SettingDetailView(Context context, AttributeSet attrs) {
        super(context, attrs);
        initView(context);
        TypedArray typeArray = context.obtainStyledAttributes(attrs, R.styleable.SettingItemView);
        String title = typeArray.getString(R.styleable.SettingItemView_title);
        String detail = typeArray.getString(R.styleable.SettingItemView_detail);
        typeArray.recycle();
        mTitle.setText(title);
        mDetail.setText(detail);
    }

    public SettingDetailView(Context context, AttributeSet attrs, int defStyle) {
        super(context, attrs, defStyle);
        initView(context);
    }

    private void initView(Context context) {
        View.inflate(context, R.layout.setting_detail_item, this);
        mTitle = (TextView) this.findViewById(R.id.tv_title);
        mDetail = (TextView) this.findViewById(R.id.tv_detail);
    }


    public void setTitle(String title) {
        if (mTitle != null) {
            mTitle.setText(title);
        }
    }

    public String getTitle() {
        return mTitle.getText().toString();
    }

    public void setDetail(String detail) {
        if (mDetail != null) {
            mDetail.setText(detail);
        }
    }

    public String getDetail() {
        return mDetail.getText().toString();
    }
}
