package com.cardinfolink.yunshouyin.view;

import android.content.Context;
import android.view.MotionEvent;
import android.view.View;
import android.view.View.OnClickListener;
import android.view.View.OnTouchListener;
import android.widget.TextView;

import com.cardinfolink.yunshouyin.R;

/**
 * 提示对话框。SelectDialog，为选择省份城市 特定的对话框
 * 上面两个按钮，中间一个搜索提示框，下面两个齿轮组件。
 */
public class SelectDialog {
    private Context mContext;
    private View dialogView;

    private TextView mOk;
    private TextView mCancel;

    private OnClickListener mOkOnClickListener;
    private OnClickListener mCancelOnClickListener;

    public SelectDialog(Context context, View view) {
        mContext = context;
        dialogView = view;

        mOk = (TextView) dialogView.findViewById(R.id.select_dialog_ok);
        mCancel = (TextView) dialogView.findViewById(R.id.select_dialog_cancel);

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

    }

    /**
     * 默认显示的对话框，按钮也是默认的行为。
     */
    public void show() {
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
}
