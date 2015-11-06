package com.cardinfolink.yunshouyin.salesman.activity;

import android.annotation.SuppressLint;
import android.content.Intent;
import android.content.pm.ApplicationInfo;
import android.content.pm.PackageManager;
import android.os.Bundle;
import android.util.Log;
import android.view.View;
import android.widget.CheckBox;
import android.widget.EditText;

import com.cardinfolink.yunshouyin.salesman.R;
import com.cardinfolink.yunshouyin.salesman.api.QuickPayException;
import com.cardinfolink.yunshouyin.salesman.core.QuickPayCallbackListener;
import com.cardinfolink.yunshouyin.salesman.model.SaveData;
import com.cardinfolink.yunshouyin.salesman.model.SessonData;
import com.cardinfolink.yunshouyin.salesman.model.SystemConfig;
import com.cardinfolink.yunshouyin.salesman.model.User;
import com.cardinfolink.yunshouyin.salesman.utils.VerifyUtil;
import com.umeng.update.UmengUpdateAgent;


public class LoginActivity extends BaseActivity {
    private final String TAG = "LoginActivity";
    private EditText mUsernameEdit;
    private EditText mPasswordEdit;
    private CheckBox mAutoLoginCheckBox;

    @SuppressLint("NewApi")
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        // check update
        UmengUpdateAgent.update(this);
        setContentView(R.layout.login_activity);

        mUsernameEdit = (EditText) findViewById(R.id.login_username);
        VerifyUtil.addEmailLimit(mUsernameEdit);
        mPasswordEdit = (EditText) findViewById(R.id.login_password);
        VerifyUtil.addEmailLimit(mPasswordEdit);
        mAutoLoginCheckBox = (CheckBox) findViewById(R.id.checkbox_auto_login);
        User user = SaveData.getUser(mContext);
        mAutoLoginCheckBox.setChecked(user.isAutoLogin());
        mUsernameEdit.setText(user.getUsername());
        mPasswordEdit.setText(user.getPassword());

        if (user.isAutoLogin()) {
            login();
        }
    }

    public void BtnLoginOnClick(View view) {
        login();
    }

    /**
     * 验证用户名密码是否有值
     *
     * @return
     */
    @SuppressLint("NewApi")
    private boolean validate() {
        String username, password;
        username = mUsernameEdit.getText().toString();
        password = mPasswordEdit.getText().toString();
        if (username.isEmpty()) {
            alertError("用户名不能为空!");
            return false;
        }

        if (password.isEmpty()) {
            alertError("密码不能为空!");
            return false;
        }
        return true;
    }


    private void login() {
        Log.d(TAG, "======================login========================");

        if (validate()) {
            startLoading();

            final String username = mUsernameEdit.getText().toString();
            final String password = mPasswordEdit.getText().toString();

            /**
             * save to share preference
             */
            User user = new User();
            user.setUsername(username);
            user.setPassword(password);

            // 自动登录checkbox 保存密码
            if (mAutoLoginCheckBox.isChecked()) {
                user.setAutoLogin(true);
                user.setPassword(password);
            }

            SaveData.setUser(mContext, user);


            /**
             * save to session
             */
            SessonData.loginUser.setUsername(username);
            SessonData.loginUser.setPassword(password);


            /**
             * async network call and callbacks
             */
            application.getQuickPayService().loginAsync(username, password, new QuickPayCallbackListener<String>() {
                @Override
                public void onSuccess(String data) {
                    SessonData.loginUser.setAccessToken(data);
                    runOnUiThread(new Runnable() {
                        @Override
                        public void run() {
                            endLoading();
                            Intent intent = new Intent(mContext, SAMerchantListActivity.class);
                            mContext.startActivity(intent);
                        }
                    });
                }

                @Override
                public void onFailure(final QuickPayException ex) {
                    runOnUiThread(new Runnable() {
                        @Override
                        public void run() {
                            if (ex.getErrorCode().equals("username_password_error")) {
                                mPasswordEdit.setText("");
                            }
                            String errorStr = ex.getErrorMsg();
                            endLoadingWithError(errorStr);
                        }
                    });
                }
            });
        }
    }


    @Override
    protected void onResume() {
        super.onResume();
        User user = SaveData.getUser(mContext);
        mAutoLoginCheckBox.setChecked(user.isAutoLogin());
        mUsernameEdit.setText(user.getUsername());
        mPasswordEdit.setText(user.getPassword());
    }
}
