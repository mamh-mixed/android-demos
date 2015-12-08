package com.cardinfolink.yunshouyin.activity;

import android.content.Intent;
import android.os.Bundle;
import android.app.Activity;
import android.view.View;

import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.ui.SettingActionBarItem;
import com.cardinfolink.yunshouyin.ui.SettingClikcItem;

public class AccountSecurityActivity extends Activity implements View.OnClickListener {

    private SettingClikcItem mAccountInfo;
    private SettingClikcItem mUpdatePassword;
    private SettingActionBarItem mAccountSecurity;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_account_security);

        mAccountInfo = (SettingClikcItem) findViewById(R.id.account_info);
        mUpdatePassword = (SettingClikcItem) findViewById(R.id.update_password);
        mAccountInfo.setOnClickListener(this);
        mUpdatePassword.setOnClickListener(this);

        mAccountSecurity = (SettingActionBarItem) findViewById(R.id.sabi_account_security);
        mAccountSecurity.setLeftTextOnclickListner(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                finish();
            }
        });
    }

    @Override
    public void onClick(View v) {
        Intent intent = null;
        switch (v.getId()) {
            case R.id.account_info:
                intent = new Intent(AccountSecurityActivity.this, AccountInfoActivity.class);
                startActivity(intent);
                break;
            case R.id.update_password:
                intent = new Intent(AccountSecurityActivity.this, UpdatePasswordActivity.class);
                startActivity(intent);
                break;
        }
        finish();//结束当前的activity
    }
}
