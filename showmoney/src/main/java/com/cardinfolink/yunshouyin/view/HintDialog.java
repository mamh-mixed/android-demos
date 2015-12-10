package com.cardinfolink.yunshouyin.view;

import android.content.Context;
import android.content.Intent;
import android.view.MotionEvent;
import android.view.View;
import android.view.View.OnClickListener;
import android.view.View.OnTouchListener;
import android.widget.TextView;

import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.activity.ForgetPasswordActivity;

/**
 * 提示对话框。hint表示提示的意思，下面两个按钮，上面一段文件的对话框。
 */
public class HintDialog {
    private Context mContext;
    private View dialogView;

    private TextView mTitle;
    private TextView mOk;
    private TextView mCancel;

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
        mCancel.setOnClickListener(new OnClickListener() {
            @Override
            public void onClick(View v) {
                dialogView.setVisibility(View.GONE);
            }
        });
    }

    /**
     * 重置密码对话框的默认行为，显示的是默认的文本。
     */
    public void show() {
        dialogView.setVisibility(View.VISIBLE);

        mOk.setOnClickListener(new OnClickListener() {

            @Override
            public void onClick(View v) {
                Intent intent = new Intent(mContext, ForgetPasswordActivity.class);
                mContext.startActivity(intent);
                dialogView.setVisibility(View.GONE);
            }
        });
    }

    /**
     * 两个按钮一个显示消息的对话框，两个按钮都是cancel的功能。
     * 可以传人不同的 文本。
     * @param title
     * @param ok
     * @param cancel
     */
    public void show(String title, String ok, String cancel) {
        dialogView.setVisibility(View.VISIBLE);
        mTitle.setText(title);
        mOk.setText(ok);
        mCancel.setText(cancel);
        mOk.setOnClickListener(new OnClickListener() {

            @Override
            public void onClick(View v) {
                dialogView.setVisibility(View.GONE);
            }
        });


    }
}
