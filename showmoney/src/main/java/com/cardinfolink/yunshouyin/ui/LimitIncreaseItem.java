package com.cardinfolink.yunshouyin.ui;

import android.content.Context;
import android.content.res.TypedArray;
import android.graphics.drawable.Drawable;
import android.util.AttributeSet;
import android.view.View;
import android.widget.ImageView;
import android.widget.TextView;

import com.cardinfolink.yunshouyin.R;

/**
 * Created by mamh on 16-1-13.
 */
public class LimitIncreaseItem extends SettingClikcItem {

    public LimitIncreaseItem(Context context) {
        super(context);
    }

    public LimitIncreaseItem(Context context, AttributeSet attrs) {
        super(context, attrs);
    }

    public LimitIncreaseItem(Context context, AttributeSet attrs, int defStyle) {
        super(context, attrs, defStyle);

        TypedArray typeArray = context.obtainStyledAttributes(attrs, R.styleable.SettingItemView);
        String title = typeArray.getString(R.styleable.SettingItemView_title);
        String right = typeArray.getString(R.styleable.SettingItemView_right_text);
        Drawable drawable = typeArray.getDrawable(R.styleable.SettingItemView_android_src);
        typeArray.recycle();

        setRightText(right);
        setTitle(title);
        setmImageView(drawable);
    }

    @Override
    protected void initView(Context context) {
        View.inflate(context, R.layout.limit_increase_item, this);
        mImageView = (ImageView) this.findViewById(R.id.iv_setting);
        mTitle = (TextView) this.findViewById(R.id.tv_title);
        mRightText = (TextView) this.findViewById(R.id.tv_right);
    }

}
