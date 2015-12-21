package com.cardinfolink.yunshouyin.activity;

import android.graphics.Bitmap;
import android.graphics.BitmapFactory;
import android.os.Bundle;
import android.text.TextUtils;
import android.view.View;
import android.widget.Button;

import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.api.QuickPayException;
import com.cardinfolink.yunshouyin.core.QuickPayCallbackListener;
import com.cardinfolink.yunshouyin.data.SessonData;
import com.cardinfolink.yunshouyin.ui.SettingActionBarItem;
import com.cardinfolink.yunshouyin.ui.SettingPasswordItem;
import com.cardinfolink.yunshouyin.util.ShowMoneyApp;

public class UpdatePasswordActivity extends BaseActivity {

    private SettingActionBarItem mActionBar;//修改密码界面的 标题栏

    private SettingPasswordItem mOriginPassword;//原始密码
    private SettingPasswordItem mNewPassword;//新密码
    private SettingPasswordItem mConfirmPassword;//确认密码

    private Button mUpdate;

    private Bitmap mWrongBitmap;//图片，错误的叉的图片，alert对话框上用的
    private Bitmap mRightBitmap;

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
        mConfirmPassword = (SettingPasswordItem) findViewById(R.id.confirm_password);

        mUpdate = (Button) findViewById(R.id.btn_update_password);
        mUpdate.setOnClickListener(new View.OnClickListener() {

            @Override
            public void onClick(View v) {
                updatePasswordOnClick(v);
            }
        });

        mWrongBitmap = BitmapFactory.decodeResource(this.getResources(), R.drawable.wrong);
        mRightBitmap = BitmapFactory.decodeResource(this.getResources(), R.drawable.right);

    }

    private void updatePasswordOnClick(View v) {
        final String originPwd = mOriginPassword.getPassword().replace(" ", "");//注意这里把所有的空格都删除了
        final String newPwd = mNewPassword.getPassword().replace(" ", "");
        final String confirmPwd = mConfirmPassword.getPassword().replace(" ", "");

        if (!validate(originPwd, newPwd, confirmPwd)) {
            return;
        }
        startLoading();

        quickPayService.updatePasswordAsync(SessonData.loginUser.getUsername(), originPwd, newPwd, new QuickPayCallbackListener<Void>() {
            @Override
            public void onSuccess(Void data) {
                //更新一下UI
                endLoading();
                String alertMsg = getResources().getString(R.string.alert_update_success);
                alertShow(alertMsg, mRightBitmap);//调用父类的方法了
                mOriginPassword.setPassword("");
                mNewPassword.setPassword("");
                mConfirmPassword.setPassword("");
            }

            @Override
            public void onFailure(QuickPayException ex) {
                endLoading();
                String error = ex.getErrorMsg();
                alertShow(error, mWrongBitmap);//调用父类的方法了
            }
        });
    }


    private boolean validate(String originPwd, String newPwd, String confirmPwd) {
        String alertMsg = "";
        if (TextUtils.isEmpty(originPwd)) {
            alertMsg = ShowMoneyApp.getResString(R.string.alert_error_old_password_cannot_empty);
            alertShow(alertMsg, mWrongBitmap);
            return false;
        }

        if (originPwd.length() < 6) {
            alertMsg = ShowMoneyApp.getResString(R.string.alert_error_old_password_short_six);
            alertShow(alertMsg, mWrongBitmap);
            return false;
        }

        if (TextUtils.isEmpty(newPwd)) {
            alertMsg = ShowMoneyApp.getResString(R.string.alert_error_new_password_cannot_empty);
            alertShow(alertMsg, mWrongBitmap);
            return false;
        }

        if (newPwd.length() < 6) {
            alertMsg = ShowMoneyApp.getResString(R.string.alert_error_new_password_short_six);
            alertShow(alertMsg, mWrongBitmap);
            return false;
        }

        if (!newPwd.equals(confirmPwd)) {
            alertMsg = ShowMoneyApp.getResString(R.string.alert_error_qrpassword_error);
            alertShow(alertMsg, mWrongBitmap);
            return false;
        }

        return true;
    }
}
