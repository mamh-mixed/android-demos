package com.cardinfolink.yunshouyin.activity;

import android.app.Activity;
import android.os.Bundle;
import android.view.View;

import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.ui.SettingActionBarItem;
import com.cardinfolink.yunshouyin.ui.SettingInputItem;
import com.cardinfolink.yunshouyin.view.ResetDialog;

/**
 * 这个是忘记密码的界面
 */
public class ForgetPasswordActivity extends Activity {

    private SettingActionBarItem mActionBar;//标题栏，自定义的标题栏

    private SettingInputItem mEmail;

    private ResetDialog mResetDialog;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_forget_password);

        mActionBar = (SettingActionBarItem) findViewById(R.id.action_bar);
        mActionBar.setLeftTextOnclickListner(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                finish();
            }
        });
        mEmail = (SettingInputItem) findViewById(R.id.email);
        mEmail.setImageViewOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                String title = getResources().getString(R.string.forget_account_message);
                String ok = getResources().getString(R.string.forget_i_known);
                String cancel = getResources().getString(R.string.forget_cancel);
                mResetDialog.show(title, ok, cancel);
            }
        });

        mResetDialog = new ResetDialog(this, findViewById(R.id.reset_dialog));

    }

}
