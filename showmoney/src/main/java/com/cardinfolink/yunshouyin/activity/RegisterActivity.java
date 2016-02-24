package com.cardinfolink.yunshouyin.activity;

import android.content.Intent;
import android.os.Bundle;
import android.text.TextUtils;
import android.view.View;
import android.widget.Button;
import android.widget.TextView;

import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.api.QuickPayException;
import com.cardinfolink.yunshouyin.core.QuickPayCallbackListener;
import com.cardinfolink.yunshouyin.data.SaveData;
import com.cardinfolink.yunshouyin.data.SessionData;
import com.cardinfolink.yunshouyin.data.User;
import com.cardinfolink.yunshouyin.ui.SettingActionBarItem;
import com.cardinfolink.yunshouyin.ui.SettingInputItem;
import com.cardinfolink.yunshouyin.ui.SettingPasswordItem;
import com.cardinfolink.yunshouyin.util.Log;
import com.cardinfolink.yunshouyin.util.ShowMoneyApp;
import com.cardinfolink.yunshouyin.util.Utility;
import com.cardinfolink.yunshouyin.util.VerifyUtil;
import com.cardinfolink.yunshouyin.view.ActivateDialog;
import com.cardinfolink.yunshouyin.view.YellowTips;

public class RegisterActivity extends BaseActivity {
    private static final String TAG = "RegisterActivity";

    private SettingActionBarItem mActionBar;//注册页面的标题栏
    private SettingInputItem mEmailEdit;
    private SettingPasswordItem mPasswordEdit;
    private SettingInputItem mInviteCode;//邀请码
    private Button mRegisterNext;
    private TextView mAgreement;

    private YellowTips mYellowTips;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_register);
        mContext = this;
        initLayout();
    }

    private void initLayout() {

        mYellowTips = new YellowTips(this, findViewById(R.id.yellow_tips));
        mActionBar = (SettingActionBarItem) findViewById(R.id.action_bar);//注册页面标题栏
        mEmailEdit = (SettingInputItem) findViewById(R.id.register_email);
        mEmailEdit.setImageViewDrawable(null);

        mPasswordEdit = (SettingPasswordItem) findViewById(R.id.register_password);


        mInviteCode = (SettingInputItem) findViewById(R.id.register_invite_code);//邀请码
        mInviteCode.setImageViewDrawable(null);

        mAgreement = (TextView) findViewById(R.id.tv_agreement);

        //注册页面标题栏添加返回事件监听
        mActionBar.setLeftTextOnclickListner(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                Utility.hideInput(mContext, v);
                finish();
            }
        });

        mRegisterNext = (Button) findViewById(R.id.btnregister);
        mRegisterNext.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                Utility.hideInput(mContext, v);
                btnRegisterNextOnClick(v);
            }
        });

        mAgreement.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                Utility.hideInput(mContext, v);
                Intent intent = new Intent(RegisterActivity.this, AgreementActivity.class);
                startActivity(intent);
            }
        });
    }

    public void btnRegisterNextOnClick(View view) {
        final String username = mEmailEdit.getText(); //用户名
        final String password = mPasswordEdit.getPassword(); //密码，第一次输入的
        final String invite = mInviteCode.getText();//邀请码

        if (!validate(username, password)) {
            return;
        }

        mLoadingDialog.startLoading();
        quickPayService.registerAsync(username, password, invite, new QuickPayCallbackListener<Void>() {
            @Override
            public void onSuccess(Void data) {
                //没有返回值的,那边返回的是null.
                User user = new User();
                user.setUsername(username);
                user.setPassword(password);
                SaveData.setUser(mContext, user);//保存到本地文件

                //这里保存到sessondata中，下面一个activity中会用到
                SessionData.loginUser.setUsername(username);//保存到sessondata里面
                SessionData.loginUser.setPassword(password);//保存到sessondata里面

                mLoadingDialog.endLoading();

                activate();
            }

            @Override
            public void onFailure(QuickPayException ex) {
                String errorMsg = ex.getErrorMsg();
                mLoadingDialog.endLoading();
                mYellowTips.show(errorMsg);
            }
        });
    }

    private void activate() {
        String username = SessionData.loginUser.getUsername();
        String password = SessionData.loginUser.getPassword();
        quickPayService.activateAsync(username, password, new QuickPayCallbackListener<Void>() {
            @Override
            public void onSuccess(Void data) {
                //走到这里是注册成功了 ，就是注册的第一步成功了，之后会进入完善信息的页面。
                //如果注册成功 点击 了 已激活 按钮 就跳转到 注册的第二个页面。
                View activateView = findViewById(R.id.activate_dialog);
                String username = SessionData.loginUser.getUsername();
                ActivateDialog activateDialog = new ActivateDialog(mContext, activateView, username);
                activateDialog.setBodyText(getResources().getString(R.string.activate_message) + username);
                activateDialog.setCancelText(getResources().getString(R.string.activate_after));
                activateDialog.setOkText(getResources().getString(R.string.activate_had_activated));
                //这里重新实现右边确定按钮的点击事件。默认是进入login界面，这里要进入RegisterActivateActivity界面
                activateDialog.setOkOnClickListener(new View.OnClickListener() {
                    @Override
                    public void onClick(View v) {
                        Intent intent = new Intent(RegisterActivity.this, RegisterActivateActivity.class);
                        startActivity(intent);//跳转到RegisterActivateActivity页面
                        //在RegisterActivateActivity里面判断一下是否真的激活了。是的话就进入RegisterNext页面了。
                        finish();//结束当前的页面
                    }
                });

                activateDialog.show();//显示激活对话框，这里文本都是自定义的。
            }

            @Override
            public void onFailure(QuickPayException ex) {
                String errorMsg = ex.getErrorMsg();
                mYellowTips.show(errorMsg);
            }
        });
    }

    private boolean validate(String email, String password) {
        String alertMsg = "";
        if (TextUtils.isEmpty(email)) {
            alertMsg = ShowMoneyApp.getResString(R.string.alert_error_email_cannot_empty);
            mYellowTips.show(alertMsg);
            return false;
        }
        if (!VerifyUtil.checkEmail(email)) {
            alertMsg = ShowMoneyApp.getResString(R.string.alert_error_email_format_error);
            mYellowTips.show(alertMsg);
            return false;
        }
        if (TextUtils.isEmpty(password)) {
            alertMsg = ShowMoneyApp.getResString(R.string.alert_error_password_cannot_empty);
            mYellowTips.show(alertMsg);
            return false;
        }
        if (password.length() < 8) {
            alertMsg = ShowMoneyApp.getResString(R.string.alert_error_password_too_short);
            mYellowTips.show(alertMsg);
            return false;
        }
        if (password.length() > 30) {
            alertMsg = ShowMoneyApp.getResString(R.string.alert_error_password_too_long);
            mYellowTips.show(alertMsg);
            return false;
        }
        //检查密码等级返回一个整数
        int level = VerifyUtil.checkPasswordLevel(password);
        if (level < 2) {
            alertMsg = ShowMoneyApp.getResString(R.string.alert_error_password_should_contain);
            mYellowTips.show(alertMsg);
            return false;
        }

        return true;
    }

}
