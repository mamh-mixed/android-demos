package com.cardinfolink.yunshouyin.salesman.activity;

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
        VerifyUtil.addEmailLimit(mEmailEdit);
        mPasswordEdit = (EditText) findViewById(R.id.register_password);
        VerifyUtil.addEmailLimit(mPasswordEdit);
        mQrPasswordEdit = (EditText) findViewById(R.id.register_qr_password);
        VerifyUtil.addEmailLimit(mQrPasswordEdit);


        btnLogin = (Button) findViewById(R.id.btnlogin);
        btnLogin.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                startLoading();
                final String username = mEmailEdit.getText().toString();
                final String password = mPasswordEdit.getText().toString();
                final String password_repeat = mQrPasswordEdit.getText().toString();

                application.getQuickPayService().registerUserAsync(username, password, password_repeat, new QuickPayCallbackListener<User>() {
                    @Override
                    public void onSuccess(User data) {
                        SessonData.registerUser.setUsername(username);
                        SessonData.registerUser.setPassword(password);
                        Log.d(TAG, SessonData.registerUser.getJsonString());

                        SharedPreferences.Editor editor = mSharedPreferences.edit();
                        editor.putInt("register_step_finish", 1);
                        editor.putString("register_username", username);
                        editor.putString("register_password", password);
                        editor.commit();

                        endLoading();
                        intentToActivity(RegisterNextActivity.class);
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
