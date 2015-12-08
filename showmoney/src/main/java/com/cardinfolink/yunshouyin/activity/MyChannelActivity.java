package com.cardinfolink.yunshouyin.activity;

import android.os.Bundle;
import android.app.Activity;
import android.view.View;

import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.ui.SettingActionBarItem;

public class MyChannelActivity extends Activity {

    private SettingActionBarItem mMyChannel;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_my_channel);

        mMyChannel = (SettingActionBarItem) findViewById(R.id.sabi_my_channel);
        mMyChannel.setLeftTextOnclickListner(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                finish();
            }
        });
    }

}
