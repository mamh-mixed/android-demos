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
    private String mEmali;

    public ActivateDialog(Context context, View view, String email) {
        mContext = context;
        dialogView = view;
        mEmali = email;
    }

    public void show() {
        TextView textView = (TextView) dialogView.findViewById(R.id.email);
        textView.setText(ShowMoneyApp.getResString(R.string.activate_sentto_email) + mEmali + "");
        dialogView.setVisibility(View.VISIBLE);
        dialogView.setOnTouchListener(new OnTouchListener() {

            @Override
            public boolean onTouch(View v, MotionEvent event) {
                return true;
            }
        });
        dialogView.findViewById(R.id.activate_dialog_cancel).setOnClickListener(new OnClickListener() {

            @Override
            public void onClick(View v) {
                enterLoginActivity();
            }
        });

        dialogView.findViewById(R.id.activate_dialog_ok).setOnClickListener(new OnClickListener() {

            @Override
            public void onClick(View v) {
                String username = SessonData.loginUser.getUsername();
                String password = SessonData.loginUser.getPassword();
                if (TextUtils.isEmpty(username) || TextUtils.isEmpty(password)) {
                    //sessondata如果没有保存的情况下会出现。在loginActivity里没设置的话会出现此问题。
                    //强制停止应用，再次登录会出现此问题。因为此时
                    //sessondata是空的没有设置username和password。
                    enterLoginActivity();
                    return;
                }
                QuickPayService quick = ShowMoneyApp.getInstance().getQuickPayService();
                quick.activateAsync(username, password, new QuickPayCallbackListener<Void>() {
                    @Override
                    public void onSuccess(Void data) {
                        enterLoginActivity();
                    }

                    @Override
                    public void onFailure(QuickPayException ex) {

                    }
                });

            }
        });
    }

    private void enterLoginActivity() {
        dialogView.setVisibility(View.GONE);
        Intent intent = new Intent(mContext, LoginActivity.class);
        mContext.startActivity(intent);
        if (!(mContext instanceof LoginActivity)) {
            ((Activity) mContext).finish();
        }
    }
}
