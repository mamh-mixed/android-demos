package com.cardinfolink.yunshouyin.activity;

import android.content.Intent;
import android.graphics.Bitmap;
import android.graphics.BitmapFactory;
import android.os.Bundle;
import android.text.Editable;
import android.text.TextUtils;
import android.text.TextWatcher;
import android.util.Log;
import android.view.View;
import android.view.animation.Animation;
import android.view.animation.AnimationUtils;
import android.widget.Button;
import android.widget.CheckBox;
import android.widget.ImageView;
import android.widget.TextView;

import com.cardinfolink.cashiersdk.model.InitData;
import com.cardinfolink.cashiersdk.sdk.CashierSdk;
import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.api.QuickPayException;
import com.cardinfolink.yunshouyin.constant.SystemConfig;
import com.cardinfolink.yunshouyin.core.QuickPayCallbackListener;
import com.cardinfolink.yunshouyin.data.SaveData;
import com.cardinfolink.yunshouyin.data.SessonData;
import com.cardinfolink.yunshouyin.data.User;
import com.cardinfolink.yunshouyin.ui.EditTextClear;
import com.cardinfolink.yunshouyin.util.TelephonyManagerUtil;
import com.cardinfolink.yunshouyin.util.VerifyUtil;
import com.cardinfolink.yunshouyin.view.ActivateDialog;
import com.cardinfolink.yunshouyin.view.HintDialog;

public class LoginActivity extends BaseActivity {
    private static final String TAG = "LoginActivity";

    private EditTextClear mUsernameEdit;
    private EditTextClear mPasswordEdit;
    private CheckBox mAutoLogin;
    private TextView mRegister;//用来注册
    private ImageView mIncrease;//用来添加用户
    private ImageView mHelp;//显示帮助对话框
    private ImageView mLoad;//loading的图片

    private Button mLoginButton;//登录按钮
    private HintDialog mHintDialog;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_login);
        initView();
        initLisenter();

        User user = SaveData.getUser(mContext);
        mAutoLogin.setChecked(user.isAutoLogin());
        mUsernameEdit.setText(user.getUsername());
        mPasswordEdit.setText(user.getPassword());
        if (user.isAutoLogin()) {
            startLoading();
            login();
        }
    }


    public void startLoading() {
        Animation loadingAnimation = AnimationUtils.loadAnimation(mContext, R.anim.loading_animation);
        mLoad.startAnimation(loadingAnimation);
        mLoad.setVisibility(View.VISIBLE);
    }

    public void endLoading() {
        mLoad.clearAnimation();
        mLoad.setVisibility(View.GONE);
    }

    /**
     * 初始化各个 组件
     */
    private void initView() {
        mIncrease = (ImageView) findViewById(R.id.iv_increase);//用来添加用户
        mHelp = (ImageView) findViewById(R.id.iv_help);//显示帮助对话框

        mRegister = (TextView) findViewById(R.id.tv_register);//用来注册

        mUsernameEdit = (EditTextClear) findViewById(R.id.login_username);
        VerifyUtil.addEmailLimit(mUsernameEdit);

        mPasswordEdit = (EditTextClear) findViewById(R.id.login_password);
        VerifyUtil.addEmailLimit(mPasswordEdit);


        mAutoLogin = (CheckBox) findViewById(R.id.login_auto);

        mLoginButton = (Button) findViewById(R.id.btnlogin);

        mHintDialog = new HintDialog(this, findViewById(R.id.hint_dialog));

        mLoad = (ImageView) findViewById(R.id.iv_loading);

    }

    /**
     * 初始化相应组件的 事件监听
     */
    private void initLisenter() {
        mUsernameEdit.addTextChangedListener(new LoginTextWatcher(mUsernameEdit));

        mPasswordEdit.addTextChangedListener(new LoginTextWatcher(mPasswordEdit));

        mLoginButton.setOnClickListener(new LoginOnClickListener());

        mHelp.setOnClickListener(new LoginOnClickListener());

        mRegister.setOnClickListener(new LoginOnClickListener());
    }


    private boolean validate() {
        String username, password;
        username = mUsernameEdit.getText().toString();
        password = mPasswordEdit.getText().toString();

        Bitmap wrongBitmap = BitmapFactory.decodeResource(this.getResources(), R.drawable.wrong);

        if (TextUtils.isEmpty(username)) {
            String alertMsg = getResources().getString(R.string.alert_error_username_cannot_empty);
            mAlertDialog.show(alertMsg, wrongBitmap);
            return false;
        }

        if (TextUtils.isEmpty(password)) {
            String alertMsg = getResources().getString(R.string.alert_error_password_cannot_empty);
            mAlertDialog.show(alertMsg, wrongBitmap);
            return false;
        }
        return true;
    }


    private void login() {
        if (!validate()) {
            endLoading();
            return;
        }

        final String username = mUsernameEdit.getText().toString();
        final String password = mPasswordEdit.getText().toString();

        quickPayService.loginAsync(username, password, new QuickPayCallbackListener<User>() {
            @Override
            public void onSuccess(User data) {
                endLoading();

                User user = new User();
                user.setUsername(username);
                if (mAutoLogin.isChecked()) {
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
                    SessonData.positionView = 0;
                    Intent intent = new Intent(mContext, MainActivity.class);
                    intent.setFlags(Intent.FLAG_ACTIVITY_CLEAR_TOP);
                    mContext.startActivity(intent);
                }
            }

            @Override
            public void onFailure(QuickPayException ex) {
                endLoading();
                String errorCode = ex.getErrorCode();
                String errorMsg = ex.getErrorMsg();
                Log.e(TAG, "login onFailure: " + errorCode + " = " + errorMsg);
                User user = new User();
                user.setUsername(username);
                user.setPassword(password);
                SaveData.setUser(mContext, user);
                SessonData.loginUser.setUsername(username);
                SessonData.loginUser.setPassword(password);
                if (errorCode.equals("user_no_activate")) {
                    //更新UI
                    View view = findViewById(R.id.activate_dialog);
                    String eMail = SessonData.loginUser.getUsername();
                    ActivateDialog activateDialog = new ActivateDialog(mContext, view, eMail);
                    activateDialog.show();
                } else {
                    mAlertDialog.show(errorMsg, BitmapFactory.decodeResource(mContext.getResources(), R.drawable.wrong));
                    if (errorCode.equals("username_password_error")) {
                        mPasswordEdit.setText("");
                    }
                }
            }
        });
    }


    @Override
    protected void onResume() {
        super.onResume();
        User user = SaveData.getUser(mContext);
        mAutoLogin.setChecked(user.isAutoLogin());
        mUsernameEdit.setText(user.getUsername());
        mPasswordEdit.setText(user.getPassword());
    }

    private class LoginOnClickListener implements View.OnClickListener {

        @Override
        public void onClick(View v) {
            switch (v.getId()) {
                case R.id.btnlogin:
                    Log.e(TAG, "onClick to longin");
                    startLoading();
                    login();
                    break;
                case R.id.iv_help:
                    mHintDialog.setOkOnClickListener(new View.OnClickListener() {
                        @Override
                        public void onClick(View v) {
                            Intent intent = new Intent(LoginActivity.this, ForgetPasswordActivity.class);
                            startActivity(intent);
                            mHintDialog.hide();
                        }
                    });
                    mHintDialog.show();
                    break;
                case R.id.tv_register:
                    Intent intent = new Intent(LoginActivity.this, RegisterActivity.class);
                    startActivity(intent);
                    break;
            }
        }
    }

    private class LoginTextWatcher implements TextWatcher {
        private View view;

        public LoginTextWatcher(View view) {
            this.view = view;
        }

        @Override
        public void beforeTextChanged(CharSequence s, int start, int count, int after) {

        }

        @Override
        public void onTextChanged(CharSequence s, int start, int before, int count) {

        }

        @Override
        public void afterTextChanged(Editable s) {
            switch (view.getId()) {
                case R.id.login_username:
                    if (TextUtils.isEmpty(mUsernameEdit.getText())) {
                        mIncrease.setVisibility(View.VISIBLE);
                    } else {
                        mIncrease.setVisibility(View.INVISIBLE);
                    }
                    break;
                case R.id.login_password:
                    if (TextUtils.isEmpty(mPasswordEdit.getText())) {
                        mHelp.setVisibility(View.VISIBLE);
                    } else {
                        mHelp.setVisibility(View.INVISIBLE);
                    }
                    break;
            }
        }
    }

}
