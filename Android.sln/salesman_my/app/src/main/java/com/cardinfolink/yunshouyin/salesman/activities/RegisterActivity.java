package com.cardinfolink.yunshouyin.salesman.activities;

import android.annotation.SuppressLint;
import android.os.Bundle;
import android.util.Log;
import android.view.View;
import android.widget.Button;
import android.widget.EditText;

import com.cardinfolink.yunshouyin.salesman.R;
import com.cardinfolink.yunshouyin.salesman.models.SAServerPacket;
import com.cardinfolink.yunshouyin.salesman.models.SessonData;
import com.cardinfolink.yunshouyin.salesman.models.User;
import com.cardinfolink.yunshouyin.salesman.utils.ActivityCollector;
import com.cardinfolink.yunshouyin.salesman.utils.CommunicationListener;
import com.cardinfolink.yunshouyin.salesman.utils.HttpCommunicationUtil;
import com.cardinfolink.yunshouyin.salesman.utils.ParamsUtil;
import com.cardinfolink.yunshouyin.salesman.utils.VerifyUtil;

import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class RegisterActivity extends BaseActivity {
    private EditText mEmailEdit;
    private EditText mPasswordEdit;
    private EditText mQrPasswordEdit;

    private Button btnLogin;

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
        VerifyUtil.addEmialLimit(mEmailEdit);
        mPasswordEdit = (EditText) findViewById(R.id.register_password);
        VerifyUtil.addEmialLimit(mPasswordEdit);
        mQrPasswordEdit = (EditText) findViewById(R.id.register_qr_password);
        VerifyUtil.addEmialLimit(mQrPasswordEdit);
        btnLogin = (Button)findViewById(R.id.btnlogin);
        btnLogin.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                if (validate()) {
                    startLoading();
                    final String username = mEmailEdit.getText().toString();
                    final String password = mPasswordEdit.getText().toString();
                    HttpCommunicationUtil.sendDataToServer(ParamsUtil.getRegister_SA(SessonData.getAccessToken(), username, password), new CommunicationListener() {

                        @Override
                        public void onResult(final String result) {
                            final SAServerPacket serverPacket = SAServerPacket.getServerPacketFrom(result);

                            if (serverPacket.getState().equals("success")) {
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

                            } else {
                                runOnUiThread(new Runnable() {
                                    @Override
                                    public void run() {
                                        String error = serverPacket.getError();
                                        endLoadingWithError(error);

                                        if (error.equals("accessToken_error")) {
                                            //关闭所有activity,除了登录框
                                            ActivityCollector.goLoginAndFinishRest();
                                        }
                                    }
                                });
                            }
                        }

                        @Override
                        public void onError(final String error) {
                            runOnUiThread(new Runnable() {

                                @Override
                                public void run() {
                                    endLoadingWithError(error);
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

}
