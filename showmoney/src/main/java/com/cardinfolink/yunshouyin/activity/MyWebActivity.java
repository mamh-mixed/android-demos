package com.cardinfolink.yunshouyin.activity;

import android.os.Bundle;
import android.view.View;

import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.ui.SettingActionBarItem;


/**
 * 我的网页版，这里显示一个二维码图片。通过payUrl来显示二维码图片。
 */
public class MyWebActivity extends BaseActivity {
    private static final String TAG = "MyWebActivity";
    private SettingActionBarItem mActionBar;

    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_my_web);

        mActionBar = (SettingActionBarItem) findViewById(R.id.action_bar);
        mActionBar.setLeftTextOnclickListner(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                finish();
            }
        });



    }

}
