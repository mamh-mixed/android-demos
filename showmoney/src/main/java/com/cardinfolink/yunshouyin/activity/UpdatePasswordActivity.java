package com.cardinfolink.yunshouyin.activity;

import android.os.Bundle;
import android.text.TextUtils;
import android.view.View;
import android.widget.Button;

import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.api.QuickPayException;
import com.cardinfolink.yunshouyin.core.QuickPayCallbackListener;
import com.cardinfolink.yunshouyin.data.SaveData;
import com.cardinfolink.yunshouyin.data.SessionData;
import com.cardinfolink.yunshouyin.data.User;
import com.cardinfolink.yunshouyin.ui.SettingActionBarItem;
import com.cardinfolink.yunshouyin.ui.SettingPasswordItem;
import com.cardinfolink.yunshouyin.util.Log;
import com.cardinfolink.yunshouyin.util.ShowMoneyApp;
import com.cardinfolink.yunshouyin.util.VerifyUtil;
import com.cardinfolink.yunshouyin.view.YellowTips;

public class UpdatePasswordActivity extends BaseActivity {

    private SettingActionBarItem mActionBar;//修改密码界面的 标题栏

    private SettingPasswordItem mOriginPassword;//原始密码
    private SettingPasswordItem mNewPassword;//新密码

    private Button mUpdate;

    private YellowTips mYellowTips;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_update_password);
        //修改密码界面的 标题栏
        mActionBar = (SettingActionBarItem) findViewById(R.id.action_bar);
        mActionBar.setLeftTextOnclickListner(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                finish();
            }
        });

        mOriginPassword = (SettingPasswordItem) findViewById(R.id.orgin_password);
        mNewPassword = (SettingPasswordItem) findViewById(R.id.new_password);

        mUpdate = (Button) findViewById(R.id.btn_update_password);
        mUpdate.setOnClickListener(new View.OnClickListener() {

            @Override
            public void onClick(View v) {
                updatePasswordOnClick(v);
            }
        });


        mYellowTips = new YellowTips(this, findViewById(R.id.yellow_tips));
    }

    private void updatePasswordOnClick(View v) {
        final String originPwd = mOriginPassword.getPassword().replace(" ", "");//注意这里把所有的空格都删除了
        final String newPwd = mNewPassword.getPassword().replace(" ", "");

        if (!validate(originPwd, newPwd)) {
            return;
        }
        startLoading();

        quickPayService.updatePasswordAsync(SessionData.loginUser.getUsername(), originPwd, newPwd, new QuickPayCallbackListener<Void>() {
            @Override
            public void onSuccess(Void data) {
                //更新一下UI
                endLoading();
                String alertMsg = getResources().getString(R.string.alert_update_success);
                mYellowTips.show(alertMsg);
                mOriginPassword.setPassword("");
                mNewPassword.setPassword("");

                //保存密码成功之后 就要设置一下这个新的密码
                SessionData.loginUser.setPassword(newPwd);
                User user = new User();
                user.setUsername(SessionData.loginUser.getUsername());
                user.setPassword(SessionData.loginUser.getPassword());
                user.setAutoLogin(SessionData.loginUser.isAutoLogin());
                SaveData.setUser(mContext, user);
            }

            @Override
            public void onFailure(QuickPayException ex) {
                endLoading();
                String error = ex.getErrorMsg();
                mYellowTips.show(error);
            }
        });
    }


    private boolean validate(String originPwd, String newPwd) {
        String alertMsg = "";
        if (TextUtils.isEmpty(originPwd)) {
            alertMsg = ShowMoneyApp.getResString(R.string.alert_error_old_password_cannot_empty);
            mYellowTips.show(alertMsg);
            return false;
        }

        if (TextUtils.isEmpty(newPwd)) {
            alertMsg = ShowMoneyApp.getResString(R.string.alert_error_new_password_cannot_empty);
            mYellowTips.show(alertMsg);
            return false;
        }
        //新增要求密码长度不能少于八位
        if (newPwd.length() < 8) {
            alertMsg = ShowMoneyApp.getResString(R.string.alert_error_new_password_too_short);
            mYellowTips.show(alertMsg);
            return false;
        }
        if (newPwd.length() > 30) {
            alertMsg = ShowMoneyApp.getResString(R.string.alert_error_new_password_too_long);
            mYellowTips.show(alertMsg);
            return false;
        }
        //检查密码等级返回一个整数
        int level = VerifyUtil.checkPasswordLevel(newPwd);
        if (level < 2) {
            alertMsg = ShowMoneyApp.getResString(R.string.alert_error_password_should_contain);
            mYellowTips.show(alertMsg);
            return false;
        }

        return true;
    }
}
