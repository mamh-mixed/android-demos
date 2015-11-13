package com.cardinfolink.yunshouyin.salesman.activity;

import android.annotation.SuppressLint;
import android.content.Intent;
import android.os.Bundle;
import android.util.Log;
import android.view.View;
import android.widget.CheckBox;
import android.widget.EditText;

import com.cardinfolink.yunshouyin.salesman.R;
import com.cardinfolink.yunshouyin.salesman.api.QuickPayException;
import com.cardinfolink.yunshouyin.salesman.core.QuickPayCallbackListener;
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

        User user = application.getLoginUser();
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


    private void login() {
        startLoading();

        final String username = mUsernameEdit.getText().toString();
        final String password = mPasswordEdit.getText().toString();

        application.getQuickPayService().loginAsync(username, password, new QuickPayCallbackListener<String>() {
            @Override
            public void onSuccess(String data) {
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
                application.setLoginUser(user);
                endLoading();
                Intent intent = new Intent(mContext, MerchantListActivity.class);
                mContext.startActivity(intent);
            }

            @Override
            public void onFailure(final QuickPayException ex) {
                if (ex.getErrorCode().equals("username_password_error")) {
                    mPasswordEdit.setText("");
                }
                String errorStr = ex.getErrorMsg();
                endLoadingWithError(errorStr);
            }
        });
    }


    @Override
    protected void onResume() {
        super.onResume();
        User user = application.getLoginUser();
        mAutoLoginCheckBox.setChecked(user.isAutoLogin());
        mUsernameEdit.setText(user.getUsername());
        mPasswordEdit.setText(user.getPassword());
    }
}
