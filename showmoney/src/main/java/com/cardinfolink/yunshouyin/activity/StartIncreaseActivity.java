package com.cardinfolink.yunshouyin.activity;

import android.app.Activity;
import android.os.Bundle;
import android.view.View;

import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.ui.SettingActionBarItem;

/**
 * 提升限额 页面，也是 免费升级页面
 */
public class StartIncreaseActivity extends Activity {
    private SettingActionBarItem mActionBar;// action bar

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_start_increase);

        mActionBar = (SettingActionBarItem) findViewById(R.id.action_bar);
        mActionBar.setLeftTextOnclickListner(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                finish();
            }
        });


    }
}
