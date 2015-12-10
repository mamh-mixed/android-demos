package com.cardinfolink.yunshouyin.view;

import android.app.Activity;
import android.content.Context;
import android.content.Intent;
import android.text.TextUtils;
import android.view.MotionEvent;
import android.view.View;
import android.view.View.OnClickListener;
import android.view.View.OnTouchListener;
import android.widget.TextView;

import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.activity.LoginActivity;
import com.cardinfolink.yunshouyin.api.QuickPayException;
import com.cardinfolink.yunshouyin.core.QuickPayCallbackListener;
import com.cardinfolink.yunshouyin.core.QuickPayService;
import com.cardinfolink.yunshouyin.data.SessonData;
import com.cardinfolink.yunshouyin.util.ShowMoneyApp;

public class ActivateDialog {
    private Context mContext;
    private View dialogView;
    private String mEmail;

    private TextView mTitle;//对话框标题，在上边。
    private TextView mBodyText;//对话框中间显示的文本。
    private TextView mOk;//默认右边是确认按钮当然 按钮的行为可以自定义
    private TextView mCancel;//默认左边是取消按钮.当然 按钮的行为可以自定义

    private OnClickListener mOkOnClickListener;
    private OnClickListener mCancelOnClickListener;


    public ActivateDialog(Context context, View view, String email) {
        mContext = context;
        dialogView = view;
        mEmail = email;

        //对话框中间显示的文本。
        String body = ShowMoneyApp.getResString(R.string.activate_sentto_email) + mEmail + "";
        mBodyText = (TextView) dialogView.findViewById(R.id.email);
        setBodyText(body);//对话框中间显示的文本。

        mOk = (TextView) dialogView.findViewById(R.id.activate_dialog_ok);
        mCancel = (TextView) dialogView.findViewById(R.id.activate_dialog_cancel);

        dialogView.setOnTouchListener(new OnTouchListener() {
            @Override
            public boolean onTouch(View v, MotionEvent event) {
                return true;
            }
        });

        //设置取消按钮的行为，取消的话就跳转到 登录界面 了
        setCancelOnClickListener(new OnClickListener() {
            @Override
            public void onClick(View v) {
                enterLoginActivity();
            }
        });

        //设置确认按钮的行为
        setOkOnClickListener(new OnClickListener() {
            @Override
            public void onClick(View v) {
                activate();// 激活用户
            }
        });
    }

    public void show() {
        dialogView.setVisibility(View.VISIBLE);
    }

    /**
     * 激活用户账户的
     */
    private void activate() {
        String username = SessonData.loginUser.getUsername();
        String password = SessonData.loginUser.getPassword();
        if (TextUtils.isEmpty(username) || TextUtils.isEmpty(password)) {
            //sessondata如果没有保存的情况下会出现。在loginActivity里没设置的话会出现此问题。
            //强制停止应用，再次登录会出现此问题。因为此时
            //sessondata是空的没有设置username和password。
            enterLoginActivity();//进入新的界面了
            return;
        }
        QuickPayService quick = ShowMoneyApp.getInstance().getQuickPayService();
        quick.activateAsync(username, password, new QuickPayCallbackListener<Void>() {
            @Override
            public void onSuccess(Void data) {
                enterLoginActivity();//进入新的界面了
            }

            @Override
            public void onFailure(QuickPayException ex) {
                hide();//隐藏对话框了
            }
        });

    }

    private void enterLoginActivity() {
        hide();

        Intent intent = new Intent(mContext, LoginActivity.class);
        mContext.startActivity(intent);
        if (!(mContext instanceof LoginActivity)) {
            ((Activity) mContext).finish();
        }
    }


    public void hide() {
        dialogView.setVisibility(View.GONE);
    }

    /**
     * 设置对话框的 显示的文本
     *
     * @param body
     * @param ok
     * @param cancel
     */
    public void setText(String body, String ok, String cancel) {
        setBodyText(body);
        setOkText(ok);
        setCancelText(cancel);
    }

    public void setText(String title, String body, String ok, String cancel) {
        setTitle(title);
        setText(body, ok, cancel);
    }

    /**
     * 右边 一般显示 “确认”按钮
     * 设置对话框中间上边显示的文本
     *
     * @param title
     */
    public void setTitle(String title) {
        mTitle.setText(title);
    }

    /**
     * 右边 一般显示 “确认”按钮
     * 设置对话框中间显示的文本
     *
     * @param title
     */
    public void setBodyText(String title) {
        mBodyText.setText(title);
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
