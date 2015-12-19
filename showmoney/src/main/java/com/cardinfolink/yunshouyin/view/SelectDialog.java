package com.cardinfolink.yunshouyin.view;

import android.content.Context;
import android.view.MotionEvent;
import android.view.View;
import android.view.View.OnClickListener;
import android.view.View.OnTouchListener;
import android.view.animation.Animation;
import android.view.animation.TranslateAnimation;
import android.widget.TextView;

import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.ui.EditTextClear;

import kankan.wheel.widget.OnWheelChangedListener;
import kankan.wheel.widget.OnWheelScrollListener;
import kankan.wheel.widget.WheelView;
import kankan.wheel.widget.adapters.AbstractWheelTextAdapter;

/**
 * 提示对话框。SelectDialog，为选择省份城市 特定的对话框
 * 上面两个按钮，中间一个搜索提示框，下面两个齿轮组件。
 */
public class SelectDialog {
    private WheelView mWheelLeft;//左边的滚轮组件，显示省份，显示银行
    private WheelView mWheelRight;//右边的滚轮组件显示城市，显示分行

    private OnWheelScrollListener mOnWheelScrollLeftListener;

    private TranslateAnimation mShowAnimation;//显示的动画，
    private TranslateAnimation mHideAnimation;//隐藏的动画

    private Context mContext;
    private View dialogView;

    private TextView mOk;
    private TextView mCancel;
    private EditTextClear mSearch;

    private OnClickListener mOkOnClickListener;
    private OnClickListener mCancelOnClickListener;

    public SelectDialog(Context context, View view) {
        mContext = context;
        dialogView = view;

        mOk = (TextView) dialogView.findViewById(R.id.select_dialog_ok);
        mCancel = (TextView) dialogView.findViewById(R.id.select_dialog_cancel);

        mSearch = (EditTextClear) dialogView.findViewById(R.id.select_search);

        mWheelLeft = (WheelView) dialogView.findViewById(R.id.wheel_left);
        mWheelRight = (WheelView) dialogView.findViewById(R.id.wheel_right);
        mWheelLeft.setCyclic(true);
        mWheelRight.setCyclic(true);
        dialogView.setOnTouchListener(new OnTouchListener() {

            @Override
            public boolean onTouch(View v, MotionEvent event) {
                return true;
            }
        });

        //cancel默认是关闭对话框的行为
        mCancelOnClickListener = new OnClickListener() {
            @Override
            public void onClick(View v) {
                hide();
            }
        };
        mCancel.setOnClickListener(mCancelOnClickListener);

        //把ok按钮默认也设置为关闭的行为
        mOkOnClickListener = mCancelOnClickListener;
        mOk.setOnClickListener(mOkOnClickListener);

        mShowAnimation = new TranslateAnimation(
                Animation.RELATIVE_TO_SELF,
                0.0f,
                Animation.RELATIVE_TO_SELF,
                0.0f,
                Animation.RELATIVE_TO_SELF,
                1.0f,
                Animation.RELATIVE_TO_SELF,
                0.0f);
        mHideAnimation = new TranslateAnimation(
                Animation.RELATIVE_TO_SELF,
                0.0f,
                Animation.RELATIVE_TO_SELF,
                0.0f,
                Animation.RELATIVE_TO_SELF,
                0.0f,
                Animation.RELATIVE_TO_SELF,
                1.0f);
        mShowAnimation.setDuration(800);
        mHideAnimation.setDuration(800);
    }

    /**
     * 默认显示的对话框，按钮也是默认的行为。
     */
    public void show() {
        dialogView.startAnimation(mShowAnimation);
        dialogView.setVisibility(View.VISIBLE);
    }


    /**
     * 两个按钮一个显示消息的对话框，两个按钮都是cancel的功能。
     * 可以传人不同的 文本。
     * 可以设置不同文本用来显示不同内容，按钮的行为还是默认的。
     *
     * @param ok
     * @param cancel
     */
    public void show(String ok, String cancel) {
        setText(ok, cancel);
        show();
    }


    public void show(OnClickListener okListener, OnClickListener cancelListner) {
        setOkOnClickListener(okListener);
        setCancelOnClickListener(cancelListner);
        show();
    }

    public void show(String ok, String cancel, OnClickListener okListener, OnClickListener cancelListner) {
        setText(ok, cancel);
        setOkOnClickListener(okListener);
        setCancelOnClickListener(cancelListner);
        show();
    }

    public void hide() {
        dialogView.startAnimation(mHideAnimation);
        dialogView.setVisibility(View.GONE);
    }


    /**
     * 设置对话框的 显示的文本
     *
     * @param ok
     * @param cancel
     */
    public void setText(String ok, String cancel) {
        setOkText(ok);
        setCancelText(cancel);
    }


    /**
     * 设置对话框右边按钮显示的文本
     *
     * @param ok
     */
    public void setOkText(String ok) {
        mOk.setText(ok);
    }

    /**
     * zuo边一般显示 “取消” 按钮
     * 设置对话框左边按钮显示的文本
     *
     * @param cancelText
     */
    public void setCancelText(String cancelText) {
        mCancel.setText(cancelText);
    }

    //设置ok按钮的点击事件
    public void setOkOnClickListener(OnClickListener l) {
        mOkOnClickListener = l;
        mOk.setOnClickListener(l);
    }

    public void setCancelOnClickListener(OnClickListener l) {
        mCancelOnClickListener = l;
        mCancel.setOnClickListener(l);
    }

    public void addLeftScrollingListener(OnWheelScrollListener listener) {
        mWheelLeft.removeScrollingListener(mOnWheelScrollLeftListener);
        mWheelLeft.addScrollingListener(listener);
        mOnWheelScrollLeftListener = listener;
    }

    public void setWheelLeftAdapter(AbstractWheelTextAdapter adapter) {
        mWheelLeft.setViewAdapter(adapter);
    }

    public void setWheelRigthAdapter(AbstractWheelTextAdapter adapter) {
        mWheelRight.setViewAdapter(adapter);
    }

    public int getWheelLeftCurrentItem() {
        return mWheelLeft.getCurrentItem();
    }

    public int getWheelRightCurrentItem() {
        return mWheelRight.getCurrentItem();
    }

    public void setWheelLeftCurrentItem(int index) {
        mWheelLeft.setCurrentItem(index);
    }

    public void setWheelRightCurrentItem(int index) {
        mWheelRight.setCurrentItem(index);
    }


    public void setSearchText(String text) {
        mSearch.setText(text);
    }


}
