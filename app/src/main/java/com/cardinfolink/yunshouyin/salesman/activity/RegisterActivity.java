package com.cardinfolink.yunshouyin.salesman.activity;

import android.annotation.SuppressLint;
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
import com.cardinfolink.yunshouyin.salesman.utils.ErrorUtil;
import com.cardinfolink.yunshouyin.salesman.utils.VerifyUtil;

import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class RegisterActivity extends BaseActivity {
    private final String TAG = "RegisterActivity";
    private EditText mEmailEdit;
    private EditText mPasswordEdit;
    private EditText mQrPasswordEdit;

    private Button btnLogin;

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

        initLayout();

        //每次进入三步创建环节,都新建一个静态用户变量,后面两部也会使用到
        SessonData.registerUser = new User();
    }

    private void initLayout() {
        mEmailEdit = (EditText) findViewById(R.id.register_email);
        VerifyUtil.addEmailLimit(mEmailEdit);
        mPasswordEdit = (EditText) findViewById(R.id.register_password);
        VerifyUtil.addEmailLimit(mPasswordEdit);
        mQrPasswordEdit = (EditText) findViewById(R.id.register_qr_password);
        VerifyUtil.addEmailLimit(mQrPasswordEdit);


        btnLogin = (Button) findViewById(R.id.btnlogin);
        btnLogin.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                if (!validate()) {
                    return;
                }
                startLoading();
                final String username = mEmailEdit.getText().toString();
                final String password = mPasswordEdit.getText().toString();
                application.getQuickPayService().registerUserAsync(username, password, new QuickPayCallbackListener<User>() {
                    @Override
                    public void onSuccess(User data) {
                        SessonData.registerUser.setUsername(username);
                        SessonData.registerUser.setPassword(password);
                        Log.d("register user", SessonData.registerUser.getJsonString());
                        runOnUiThread(new Runnable() {
                            @Override
                            public void run() {
                                endLoading();
                                intentToActivity(RegisterNextActivity.class);
                            }
                        });
                    }

                    @Override
                    public void onFailure(final QuickPayException ex) {
                        runOnUiThread(new Runnable() {
                            @Override
                            public void run() {
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
        });
    }

    @SuppressLint("NewApi")
    private boolean validate() {
        String email, password, qr_password;
        email = mEmailEdit.getText().toString();
        password = mPasswordEdit.getText().toString();
        qr_password = mQrPasswordEdit.getText().toString();

        if (email.isEmpty()) {
            alertError("邮箱不能为空!");
            return false;
        }
        if (!checkEmail(email)) {
            alertError("邮箱格式不正确!");
            return false;
        }
        if (password.isEmpty()) {
            alertError("密码不能为空!");
            return false;
        }
        if (password.length() < 6) {
            alertError("密码不能小于六位!");
            return false;
        }
        if (!password.equals(qr_password)) {
            alertError("确认密码不一致!");
            return false;
        }

        return true;
    }

}
