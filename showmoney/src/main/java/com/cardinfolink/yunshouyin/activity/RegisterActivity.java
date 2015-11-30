package com.cardinfolink.yunshouyin.activity;

import android.annotation.SuppressLint;
import android.graphics.BitmapFactory;
import android.os.Bundle;
import android.view.View;
import android.widget.EditText;

import com.cardinfolink.yunshouyin.R;
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
import com.cardinfolink.yunshouyin.view.ActivateDialog;

import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class RegisterActivity extends BaseActivity {
    private EditText mEmailEdit;
    private EditText mPasswordEdit;
    private EditText mQrPasswordEdit;

    /**
     * 验证邮箱
     *
     * @param email
     * @return
     */
    public static boolean checkEmail(String email) {
        boolean flag = false;
        try {
            String check = "^([a-z0-9A-Z]+[-|_|\\.]?)+[a-z0-9A-Z]@([a-z0-9A-Z]+(-[a-z0-9A-Z]+)?\\.)+[a-zA-Z]{2,}$";
            Pattern regex = Pattern.compile(check);
            Matcher matcher = regex.matcher(email);
            flag = matcher.matches();
        } catch (Exception e) {
            flag = false;
        }
        return flag;
    }

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.register_activity);
        mContext = this;
        initLayout();
    }

    private void initLayout() {
        mEmailEdit = (EditText) findViewById(R.id.register_email);
        VerifyUtil.addEmailLimit(mEmailEdit);
        mPasswordEdit = (EditText) findViewById(R.id.register_password);
        VerifyUtil.addEmailLimit(mPasswordEdit);
        mQrPasswordEdit = (EditText) findViewById(R.id.register_qr_password);
        VerifyUtil.addEmailLimit(mQrPasswordEdit);


    }

    public void BtnRegisterNextOnClick(View view) {
        if (validate()) {
            mLoadingDialog.startLoading();
            final String username = mEmailEdit.getText().toString();
            final String password = mPasswordEdit.getText().toString();
            HttpCommunicationUtil.sendDataToServer(ParamsUtil.getRegister(username, password), new CommunicationListener() {

                @Override
                public void onResult(final String result) {

                    String state = JsonUtil.getParam(result, "state");
                    if (state.equals("success")) {
                        User user = new User();
                        user.setUsername(username);
                        user.setPassword(password);
                        SaveData.setUser(mContext, user);

                        SessonData.loginUser.setUsername(username);
                        SessonData.loginUser.setPassword(password);
                        runOnUiThread(new Runnable() {

                            @Override
                            public void run() {
                                //更新UI
                                mLoadingDialog.endLoading();
                                ActivateDialog activate_dialog = new ActivateDialog(mContext, RegisterActivity.this.findViewById(R.id.activate_dialog), SessonData.loginUser.getUsername());
                                activate_dialog.show();
                            }

                        });

                    } else {
                        runOnUiThread(new Runnable() {

                            @Override
                            public void run() {
                                //更新UI
                                mLoadingDialog.endLoading();
                                mAlertDialog.show(ErrorUtil.getErrorString(JsonUtil.getParam(result, "error")), BitmapFactory.decodeResource(mContext.getResources(), R.drawable.wrong));
                            }

                        });

                    }
                }

                @Override
                public void onError(final String error) {
                    runOnUiThread(new Runnable() {

                        @Override
                        public void run() {
                            //更新UI
                            mLoadingDialog.endLoading();
                            mAlertDialog.show(error, BitmapFactory.decodeResource(mContext.getResources(), R.drawable.wrong));
                        }

                    });
                }
            });
        }

    }

    @SuppressLint("NewApi")
    private boolean validate() {
        String email, password, qr_password;
        email = mEmailEdit.getText().toString();
        password = mPasswordEdit.getText().toString();
        qr_password = mQrPasswordEdit.getText().toString();

        if (email.isEmpty()) {
            mAlertDialog.show(ShowMoneyApp.getResString(R.string.alert_error_email_cannot_empty), BitmapFactory.decodeResource(this.getResources(), R.drawable.wrong));
            return false;
        }
        if (!checkEmail(email)) {
            mAlertDialog.show(ShowMoneyApp.getResString(R.string.alert_error_email_format_error), BitmapFactory.decodeResource(this.getResources(), R.drawable.wrong));
            return false;
        }

        if (password.isEmpty()) {
            mAlertDialog.show(ShowMoneyApp.getResString(R.string.alert_error_password_cannot_empty), BitmapFactory.decodeResource(this.getResources(), R.drawable.wrong));
            return false;
        }
        if (password.length() < 6) {
            mAlertDialog.show(ShowMoneyApp.getResString(R.string.alert_error_password_short_six), BitmapFactory.decodeResource(this.getResources(), R.drawable.wrong));
            return false;
        }
        if (!password.equals(qr_password)) {
            mAlertDialog.show(ShowMoneyApp.getResString(R.string.alert_error_qrpassword_error), BitmapFactory.decodeResource(this.getResources(), R.drawable.wrong));
            return false;
        }


        return true;
    }

}
