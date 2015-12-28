package com.cardinfolink.yunshouyin.ui;

import android.annotation.TargetApi;
import android.content.Context;
import android.os.Build;
import android.util.AttributeSet;
import android.widget.ListView;

/**
 * Created by mamh on 15-12-28.
 */
public class SubListView extends ListView {

    public SubListView(Context context) {
        super(context);
    }

    public SubListView(Context context, AttributeSet attrs) {
        super(context, attrs);
    }

    public SubListView(Context context, AttributeSet attrs, int defStyleAttr) {
        super(context, attrs, defStyleAttr);
    }

    @TargetApi(Build.VERSION_CODES.LOLLIPOP)
    public SubListView(Context context, AttributeSet attrs, int defStyleAttr, int defStyleRes) {
        super(context, attrs, defStyleAttr, defStyleRes);
    }


    // 主要是重写onMeasure方法，这里将heightMeasureSpec参数设大，否则，嵌套的ListView会显示不全。
    public void onMeasure(int widthMeasureSpec, int heightMeasureSpec) {
        int expandSpec = MeasureSpec.makeMeasureSpec(Integer.MAX_VALUE >> 2, MeasureSpec.AT_MOST);
        super.onMeasure(widthMeasureSpec, expandSpec);

    }
}