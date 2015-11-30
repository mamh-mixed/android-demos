package com.cardinfolink.yunshouyin.view;

import android.annotation.SuppressLint;
import android.app.Activity;
import android.content.Context;
import android.graphics.BitmapFactory;
import android.view.LayoutInflater;
import android.view.View;
import android.widget.EditText;
import android.widget.LinearLayout;

import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.activity.BaseActivity;
import com.cardinfolink.yunshouyin.data.SaveData;
import com.cardinfolink.yunshouyin.data.SessonData;
import com.cardinfolink.yunshouyin.data.User;
import com.cardinfolink.yunshouyin.util.CommunicationListener;
import com.cardinfolink.yunshouyin.util.ShowMoneyApp;
import com.cardinfolink.yunshouyin.util.ErrorUtil;
import com.cardinfolink.yunshouyin.util.HttpCommunicationUtil;
import com.cardinfolink.yunshouyin.util.JsonUtil;
import com.cardinfolink.yunshouyin.util.ParamsUtil;
import com.cardinfolink.yunshouyin.util.VerifyUtil;

public class PasswordUpdateView extends LinearLayout {
    private Context mContext;
    private EditText mOldPwdEdit;
    private EditText mNewPwdEdit;
    private EditText mQrPwdEdit;
    private BaseActivity mBaseActivity;

    public PasswordUpdateView(Context context) {
        super(context);
        mContext = context;
        mBaseActivity = (BaseActivity) mContext;
        View contentView = LayoutInflater.from(context).inflate(
                R.layout.password_update_view, null);
        LinearLayout.LayoutParams layoutParams = new LayoutParams(
                LayoutParams.MATCH_PARENT, LayoutParams.MATCH_PARENT);
        contentView.setLayoutParams(layoutParams);
        addView(contentView);
        mOldPwdEdit = (EditText) contentView.findViewById(R.id.update_password_oldpwd);
        VerifyUtil.addEmailLimit(mOldPwdEdit);
        mNewPwdEdit = (EditText) contentView.findViewById(R.id.update_password_newpwd);
        VerifyUtil.addEmailLimit(mNewPwdEdit);
        mQrPwdEdit = (EditText) contentView.findViewById(R.id.update_password_qr_newpwd);
        VerifyUtil.addEmailLimit(mQrPwdEdit);


        contentView.findViewById(R.id.btn_update_password).setOnClickListener(new OnClickListener() {

            @Override
            public void onClick(View v) {

                if (validate()) {
                    mBaseActivity.startLoading();
                    final String oldpwd = mOldPwdEdit.getText().toString().replace(" ", "");
                    final String newpwd = mNewPwdEdit.getText().toString().replace(" ", "");
                    ;
                    String qrpwd = mQrPwdEdit.getText().toString().replace(" ", "");
                    ;
                    HttpCommunicationUtil.sendDataToServer(ParamsUtil.getUpdate(SessonData.loginUser.getUsername(), oldpwd, newpwd), new CommunicationListener() {

                        @Override
                        public void onResult(String result) {
                            String state = JsonUtil.getParam(result, "state");
                            if (state.equals("success")) {
                                User user = SaveData.getUser(mContext);
                                if (user.isAutoLogin()) {
                                    user.setPassword(newpwd);
                                }
                                SaveData.setUser(mContext, user);
                                SessonData.loginUser.setPassword(newpwd);
                                ((Activity) mContext).runOnUiThread(new Runnable() {

                                    @Override
                                    public void run() {
                                        //更新UI
                                        mBaseActivity.endLoading();
                                        AlertDialog alert_Dialog = new AlertDialog(mContext, null, ((Activity) mContext).findViewById(R.id.alert_dialog), getResources().getString(R.string.alert_update_success), BitmapFactory.decodeResource(mContext.getResources(), R.drawable.right));
                                        alert_Dialog.show();
                                        mOldPwdEdit.setText("");
                                        mNewPwdEdit.setText("");
                                        mQrPwdEdit.setText("");
                                    }

                                });

                            } else {
                                final String error = JsonUtil.getParam(result, "error");
                                ((Activity) mContext).runOnUiThread(new Runnable() {

                                    @Override
                                    public void run() {
                                        mBaseActivity.endLoading();
                                        AlertDialog alert_Dialog = new AlertDialog(mContext, null, ((Activity) mContext).findViewById(R.id.alert_dialog), ErrorUtil.getErrorString(error), BitmapFactory.decodeResource(mContext.getResources(), R.drawable.wrong));
                                        alert_Dialog.show();

                                    }

                                });

                            }

                        }

                        @Override
                        public void onError(final String error) {
                            ((Activity) mContext).runOnUiThread(new Runnable() {

                                @Override
                                public void run() {
                                    mBaseActivity.endLoading();
                                    AlertDialog alert_Dialog = new AlertDialog(mContext, null, ((Activity) mContext).findViewById(R.id.alert_dialog), ErrorUtil.getErrorString(error), BitmapFactory.decodeResource(mContext.getResources(), R.drawable.wrong));
                                    alert_Dialog.show();

                                }

                            });

                        }
                    });
                }

            }
        });

    }


    @SuppressLint("NewApi")
    private boolean validate() {
        String oldpwd = mOldPwdEdit.getText().toString().replace(" ", "");
        String newpwd = mNewPwdEdit.getText().toString().replace(" ", "");
        ;
        String qrpwd = mQrPwdEdit.getText().toString().replace(" ", "");
        ;


        if (oldpwd.isEmpty()) {
            mBaseActivity.alertShow(ShowMoneyApp.getResString(R.string.alert_error_old_password_cannot_empty), BitmapFactory.decodeResource(this.getResources(), R.drawable.wrong));
            return false;
        }

        if (oldpwd.length() < 6) {
            mBaseActivity.alertShow(ShowMoneyApp.getResString(R.string.alert_error_old_password_short_six), BitmapFactory.decodeResource(this.getResources(), R.drawable.wrong));
            return false;
        }


        if (newpwd.isEmpty()) {
            mBaseActivity.alertShow(ShowMoneyApp.getResString(R.string.alert_error_new_password_cannot_empty), BitmapFactory.decodeResource(this.getResources(), R.drawable.wrong));
            return false;
        }

        if (newpwd.length() < 6) {
            mBaseActivity.alertShow(ShowMoneyApp.getResString(R.string.alert_error_new_password_short_six), BitmapFactory.decodeResource(this.getResources(), R.drawable.wrong));
            return false;
        }

        if (!qrpwd.equals(newpwd)) {
            mBaseActivity.alertShow(ShowMoneyApp.getResString(R.string.alert_error_qrpassword_error), BitmapFactory.decodeResource(this.getResources(), R.drawable.wrong));
            return false;
        }

        return true;
    }

}
