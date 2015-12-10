package com.cardinfolink.yunshouyin.activity;

import android.app.Activity;
import android.os.Bundle;
import android.util.Log;
import android.view.View;
import android.widget.ImageView;

import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.ui.SettingActionBarItem;
import com.cardinfolink.yunshouyin.ui.SettingInputItem;
import com.cardinfolink.yunshouyin.view.ResetDialog;

public class ResetPasswordActivity extends Activity {

    private SettingActionBarItem mResetPassword;

    private SettingInputItem mEmail;

    private ResetDialog mResetDialog;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_reset_password);

        mResetPassword = (SettingActionBarItem) findViewById(R.id.sabi_reset_password);
        mResetPassword.setLeftTextOnclickListner(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                finish();
            }
        });
        mEmail = (SettingInputItem) findViewById(R.id.reset_email);
        mEmail.setImageViewOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                String title = getResources().getString(R.string.reset_forget_account_message);
                String ok = getResources().getString(R.string.reset_i_known);
                String cancel = getResources().getString(R.string.reset_cancel);
                mResetDialog.show(title, ok, cancel);
            }
        });

        mResetDialog = new ResetDialog(this, findViewById(R.id.reset_dialog));

    }

}
