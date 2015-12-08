package com.cardinfolink.yunshouyin.activity;

import android.os.Bundle;
import android.app.Activity;
import android.view.View;

import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.ui.SettingActionBarItem;

public class UpdatePasswordActivity extends Activity {

    private SettingActionBarItem mUpdatePassword;//修改密码界面的 标题栏

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_update_password);

        mUpdatePassword = (SettingActionBarItem) findViewById(R.id.sabi_update_password);
        mUpdatePassword.setLeftTextOnclickListner(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                finish();
            }
        });
    }

}
