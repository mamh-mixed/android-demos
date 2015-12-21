package com.cardinfolink.yunshouyin.activity;

import android.content.Intent;
import android.os.Bundle;
import android.app.Activity;
import android.view.View;
import android.widget.Button;
import android.widget.ListView;

import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.ui.SettingActionBarItem;

public class MessageActivity extends Activity {

    private SettingActionBarItem mActionBar;

    private Button mSetAllMessageRead;

    private ListView mAllMessageListView;


    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_message);

        mActionBar = (SettingActionBarItem) findViewById(R.id.action_bar);
        mActionBar.setLeftTextOnclickListner(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                //返回
                finish();
            }
        });
        mActionBar.setRightTextOnclickListner(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                //进入未读消息界面
                Intent intent = new Intent(MessageActivity.this, UnReadMessageActivity.class);
                startActivity(intent);
                finish();
            }
        });

        mSetAllMessageRead = (Button) findViewById(R.id.set_all_message_read);
        mAllMessageListView = (ListView) findViewById(R.id.all_message);
    }

}
