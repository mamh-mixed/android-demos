package com.cardinfolink.yunshouyin.activity;

import android.annotation.SuppressLint;
import android.content.Intent;
import android.graphics.BitmapFactory;
import android.os.Bundle;
import android.text.TextUtils;
import android.util.Log;
import android.view.View;
import android.widget.CheckBox;
import android.widget.EditText;

import com.cardinfolink.cashiersdk.model.InitData;
import com.cardinfolink.cashiersdk.sdk.CashierSdk;
import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.api.QuickPayException;
import com.cardinfolink.yunshouyin.constant.SystemConfig;
import com.cardinfolink.yunshouyin.core.QuickPayCallbackListener;
import com.cardinfolink.yunshouyin.data.SaveData;
import com.cardinfolink.yunshouyin.data.SessonData;
import com.cardinfolink.yunshouyin.data.User;
import com.cardinfolink.yunshouyin.util.CommunicationListener;
import com.cardinfolink.yunshouyin.util.ErrorUtil;
import com.cardinfolink.yunshouyin.util.HttpCommunicationUtil;
import com.cardinfolink.yunshouyin.util.JsonUtil;
import com.cardinfolink.yunshouyin.util.ParamsUtil;
import com.cardinfolink.yunshouyin.util.TelephonyManagerUtil;
import com.cardinfolink.yunshouyin.util.VerifyUtil;
import com.cardinfolink.yunshouyin.view.ActivateDialog;

public class LoginActivity extends BaseActivity {
    private static final String TAG = "LoginActivity";

    private EditText mUsernameEdit;
    private EditText mPasswordEdit;
    private CheckBox mAutoLoginCheckBox;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
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

    @SuppressLint("NewApi")
    private boolean validate() {
        String username, password;
        username = mUsernameEdit.getText().toString();
        password = mPasswordEdit.getText().toString();
        if (username.isEmpty()) {
            mAlertDialog.show(getResources().getString(R.string.alert_error_username_cannot_empty), BitmapFactory.decodeResource(this.getResources(), R.drawable.wrong));
            return false;
        }

        if (password.isEmpty()) {
            mAlertDialog.show(getResources().getString(R.string.alert_error_password_cannot_empty), BitmapFactory.decodeResource(this.getResources(), R.drawable.wrong));
            return false;
        }
        return true;
    }


    private void login() {
        if (!validate()) {
            return;
        }


        final String username = mUsernameEdit.getText().toString();
        final String password = mPasswordEdit.getText().toString();

        quickPayService.loginAsync(username, password, new QuickPayCallbackListener<User>() {
            @Override
            public void onSuccess(User data) {
                User user = new User();
                user.setUsername(username);
                if (mAutoLoginCheckBox.isChecked()) {
                    user.setPassword(password);
                    user.setAutoLogin(true);
                }
                SaveData.setUser(mContext, user);
                SessonData.loginUser.setUsername(username);
                SessonData.loginUser.setPassword(password);
                SessonData.loginUser.setClientid(data.getClientid());
                SessonData.loginUser.setObjectId(data.getObjectId());
                SessonData.loginUser.setLimit(data.getLimit());

                if (TextUtils.isEmpty(SessonData.loginUser.getClientid())) {
                    // clientid为空,跳转到完善信息页面
                    mLoadingDialog.endLoading();
                    Intent intent = new Intent(mContext, RegisterNextActivity.class);
                    mContext.startActivity(intent);
                } else {
                    InitData initData = new InitData();
                    initData.setMchntid(data.getClientid());
                    initData.setInscd(data.getInscd());
                    initData.setSignKey(data.getSignKey());
                    initData.setTerminalid(TelephonyManagerUtil.getDeviceId(mContext));
                    initData.setIsProduce(SystemConfig.IS_PRODUCE);
                    CashierSdk.init(initData);//初始化sdk
                    //更新UI
                    mLoadingDialog.endLoading();
                    SessonData.positionView = 0;
                    Intent intent = new Intent(mContext, MainActivity.class);
                    intent.setFlags(Intent.FLAG_ACTIVITY_CLEAR_TOP);
                    mContext.startActivity(intent);
                }
            }

            @Override
            public void onFailure(QuickPayException ex) {
                String errorCode = ex.getErrorCode();
                String errorMsg = ex.getErrorMsg();
                Log.e(TAG, "login onFailure: " + errorCode + " = " + errorMsg);
                User user = new User();
                user.setUsername(username);
                user.setPassword(password);
                SaveData.setUser(mContext, user);
                if (errorCode.equals("user_no_activate")) {
                    //更新UI
                    mLoadingDialog.endLoading();
                    View view = findViewById(R.id.activate_dialog);
                    String eMail = SessonData.loginUser.getUsername();
                    ActivateDialog activateDialog = new ActivateDialog(mContext, view, eMail);
                    activateDialog.show();
                } else {
                    mLoadingDialog.endLoading();
                    mAlertDialog.show(errorMsg, BitmapFactory.decodeResource(mContext.getResources(), R.drawable.wrong));
                    if (errorCode.equals("username_password_error")) {
                        mPasswordEdit.setText("");
                    }
                }
            }
        });
    }

    public void BtnRegisterOnClick(View view) {
        Intent intent = new Intent(LoginActivity.this, RegisterActivity.class);
        LoginActivity.this.startActivity(intent);
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
