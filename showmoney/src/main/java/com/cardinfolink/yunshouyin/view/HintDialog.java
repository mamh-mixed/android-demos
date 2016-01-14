package com.cardinfolink.yunshouyin.view;

import android.content.Context;
import android.view.MotionEvent;
import android.view.View;
import android.view.View.OnClickListener;
import android.view.View.OnTouchListener;
import android.widget.TextView;

import com.cardinfolink.yunshouyin.R;

/**
 * 提示对话框。hint表示提示的意思，下面两个按钮，上面一段文件的对话框。
 */
public class HintDialog {
    private Context mContext;
    private View dialogView;

    private TextView mTitle;
    private TextView mOk;
    private TextView mCancel;

    private OnClickListener mOkOnClickListener;
    private OnClickListener mCancelOnClickListener;

    public HintDialog(Context context, View view) {
        mContext = context;
        dialogView = view;

        mTitle = (TextView) dialogView.findViewById(R.id.hint_messsage);
        mOk = (TextView) dialogView.findViewById(R.id.hint_dialog_ok);
        mCancel = (TextView) dialogView.findViewById(R.id.hint_dialog_cancel);

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
                dialogView.setVisibility(View.GONE);
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
     * @param title
     * @param ok
     * @param cancel
     */
    public void show(String title, String ok, String cancel) {
        setText(title, ok, cancel);
        show();
    }

    public void show(String title, String ok, String cancel, View.OnClickListener okOnClickListener, View.OnClickListener cancelOnClickListener) {
        setOkOnClickListener(okOnClickListener);
        setCancelOnClickListener(cancelOnClickListener);
        show(title, ok, cancel);
    }

    public void show(String ok, String cancel, View.OnClickListener okOnClickListener, View.OnClickListener cancelOnClickListener) {
        String title = mTitle.getText().toString();
        show(title, ok, cancel, okOnClickListener, cancelOnClickListener);
    }

    public void hide() {
        dialogView.setVisibility(View.GONE);
    }

    /**
     * 设置对话框的 显示的文本
     *
     * @param title
     * @param ok
     * @param cancel
     */
    public void setText(String title, String ok, String cancel) {
        setTitle(title);
        setOkText(ok);
        setCancelText(cancel);
    }

    /**
     * 右边 一般显示 “确认”按钮
     * 设置对话框中间显示的文本
     *
     * @param title
     */
    public void setTitle(String title) {
        mTitle.setText(title);
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

    /**
     * ok是右边按钮
     *
     * @param v
     */
    public void setOkVisibility(int v) {
        mCancel.setVisibility(v);
    }

    /**
     * cancel是左边按钮
     *
     * @param v
     */
    public void setCancelVisibility(int v) {
        mCancel.setVisibility(v);
    }
}
