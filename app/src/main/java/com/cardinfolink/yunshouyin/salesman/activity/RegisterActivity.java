package com.cardinfolink.yunshouyin.salesman.activity;

import android.content.Intent;
import android.content.SharedPreferences;
import android.os.Bundle;
import android.util.Log;
import android.view.View;
import android.widget.Button;
import android.widget.EditText;

import com.cardinfolink.yunshouyin.salesman.R;
import com.cardinfolink.yunshouyin.salesman.api.QuickPayException;
import com.cardinfolink.yunshouyin.salesman.core.QuickPayCallbackListener;
import com.cardinfolink.yunshouyin.salesman.model.SessonData;
import com.cardinfolink.yunshouyin.salesman.model.User;
import com.cardinfolink.yunshouyin.salesman.utils.ActivityCollector;
import com.cardinfolink.yunshouyin.salesman.utils.VerifyUtil;

public class RegisterActivity extends BaseActivity {
    private final String TAG = "RegisterActivity";

    private EditText mEmail;
    private EditText mPassword;
    private EditText mQrPassword;
    private Button mLogin;


    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.register_activity);
        initLayout();
        //每次进入三步创建环节,都新建一个静态用户变量,后面两部也会使用到
        SessonData.registerUser = new User();
    }

    private void initLayout() {
        mEmail = (EditText) findViewById(R.id.register_email);
        VerifyUtil.addEmailLimit(mEmail);
        mPassword = (EditText) findViewById(R.id.register_password);
        VerifyUtil.addEmailLimit(mPassword);
        mQrPassword = (EditText) findViewById(R.id.register_qr_password);
        VerifyUtil.addEmailLimit(mQrPassword);

        mLogin = (Button) findViewById(R.id.btnlogin);
        mLogin.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                startLoading();
                final String username = mEmail.getText().toString();
                final String password = mPassword.getText().toString();
                final String passwordRepeat = mQrPassword.getText().toString();

                quickPayService.registerUserAsync(username, password, passwordRepeat, new QuickPayCallbackListener<User>() {
                    @Override
                    public void onSuccess(User data) {
                        SessonData.registerUser.setUsername(username);
                        SessonData.registerUser.setPassword(password);

                        SharedPreferences.Editor editor = mRegisterSharedPreferences.edit();
                        editor.putInt("register_step_finish", 1);
                        editor.putString("register_username", username);
                        editor.putString("register_password", password);
                        editor.commit();

                        endLoading();
                        Intent intent = new Intent(RegisterActivity.this, RegisterNextActivity.class);
                        startActivity(intent);
                        finish();
                    }

                    @Override
                    public void onFailure(final QuickPayException ex) {
                        String error = ex.getErrorMsg();
                        endLoadingWithError(error);
                        if (ex.getErrorCode().equals(QuickPayException.ACCESSTOKEN_NOT_FOUND)) {
                            //关闭所有activity,除了登录框
                            ActivityCollector.goLoginAndFinishRest();
                        }
                    }
                });
            }
        });
    }
}
