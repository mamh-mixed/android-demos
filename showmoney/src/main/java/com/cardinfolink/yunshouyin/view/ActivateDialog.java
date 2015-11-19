package com.cardinfolink.yunshouyin.view;

import android.app.Activity;
import android.content.Context;
import android.content.Intent;
import android.view.MotionEvent;
import android.view.View;
import android.view.View.OnClickListener;
import android.view.View.OnTouchListener;
import android.widget.TextView;

import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.activity.LoginActivity;
import com.cardinfolink.yunshouyin.data.SessonData;
import com.cardinfolink.yunshouyin.util.CommunicationListener;
import com.cardinfolink.yunshouyin.util.ShowMoneyApp;
import com.cardinfolink.yunshouyin.util.HttpCommunicationUtil;
import com.cardinfolink.yunshouyin.util.JsonUtil;
import com.cardinfolink.yunshouyin.util.ParamsUtil;

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
                // TODO Auto-generated method stub
                return true;
            }
        });
        dialogView.findViewById(R.id.activate_dialog_cancel)
                .setOnClickListener(new OnClickListener() {

                    @Override
                    public void onClick(View v) {
                        dialogView.setVisibility(View.GONE);
                        Intent intent = new Intent(mContext,
                                LoginActivity.class);
                        mContext.startActivity(intent);
                        if (!(mContext instanceof LoginActivity)) {
                            ((Activity) mContext).finish();
                        }

                    }
                });

        dialogView.findViewById(R.id.activate_dialog_ok).setOnClickListener(
                new OnClickListener() {

                    @Override
                    public void onClick(View v) {

                        HttpCommunicationUtil.sendDataToServer(ParamsUtil
                                        .getRequestActivate(
                                                SessonData.loginUser.getUsername(),
                                                SessonData.loginUser.getPassword()),
                                new CommunicationListener() {

                                    @Override
                                    public void onResult(String result) {
                                        String state = JsonUtil.getParam(
                                                result, "state");
                                        if (state.equals("success")) {

                                            ((Activity) mContext).runOnUiThread(new Runnable() {

                                                @Override
                                                public void run() {
                                                    //更新UI
                                                    dialogView.setVisibility(View.GONE);
                                                    Intent intent = new Intent(mContext,
                                                            LoginActivity.class);
                                                    mContext.startActivity(intent);
                                                    if (!(mContext instanceof LoginActivity)) {
                                                        ((Activity) mContext).finish();
                                                    }
                                                }

                                            });

                                        }
                                    }

                                    @Override
                                    public void onError(String error) {

                                    }
                                });


                    }
                });
    }
}
