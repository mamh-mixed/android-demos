package com.cardinfolink.yunshouyin.salesman.activity;

import android.content.Intent;
import android.os.Bundle;
import android.view.View;
import android.widget.Button;
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

    private EditText mUsername;
    private EditText mPassword;
    private CheckBox mAutoLogin;
    private Button mLogin;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        // check update
        UmengUpdateAgent.update(this);
        setContentView(R.layout.login_activity);

        mUsername = (EditText) findViewById(R.id.login_username);
        VerifyUtil.addEmailLimit(mUsername);

        mPassword = (EditText) findViewById(R.id.login_password);
        VerifyUtil.addEmailLimit(mPassword);

        mAutoLogin = (CheckBox) findViewById(R.id.checkbox_auto_login);

        User user = getLoginUser();
        mAutoLogin.setChecked(user.isAutoLogin());
        mUsername.setText(user.getUsername());
        mPassword.setText(user.getPassword());

        if (user.isAutoLogin()) {
            login();
        }

        mLogin = (Button) findViewById(R.id.btnlogin);
        mLogin.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                login();
            }
        });
    }


    private void login() {
        startLoading();

        final String username = mUsername.getText().toString();
        final String password = mPassword.getText().toString();

        quickPayService.loginAsync(username, password, new QuickPayCallbackListener<String>() {
            @Override
            public void onSuccess(String data) {
                //save to share preference
                User user = new User();
                user.setUsername(username);
                user.setPassword(password);
                // 自动登录checkbox 保存密码
                if (mAutoLogin.isChecked()) {
                    user.setAutoLogin(true);
                    user.setPassword(password);
                }
                setLoginUser(user);
                endLoading();
                Intent intent = new Intent(LoginActivity.this, MerchantListActivity.class);
                startActivity(intent);
            }

            @Override
            public void onFailure(final QuickPayException ex) {
                if (ex.getErrorCode().equals("username_password_error")) {
                    mPassword.setText("");
                }
                String errorStr = ex.getErrorMsg();
                endLoadingWithError(errorStr);
            }
        });
    }


    @Override
    protected void onResume() {
        super.onResume();
        User user = getLoginUser();
        mAutoLogin.setChecked(user.isAutoLogin());
        mUsername.setText(user.getUsername());
        mPassword.setText(user.getPassword());
    }
}
