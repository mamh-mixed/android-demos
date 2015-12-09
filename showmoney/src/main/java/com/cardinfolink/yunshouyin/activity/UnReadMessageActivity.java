package com.cardinfolink.yunshouyin.activity;

import android.content.Intent;
import android.os.Bundle;
import android.app.Activity;
import android.view.View;
import android.widget.Button;
import android.widget.ListView;

import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.ui.SettingActionBarItem;

public class UnReadMessageActivity extends Activity {

    private SettingActionBarItem mUnreadMessage;

    private Button mSetAllMessageRead;

    private ListView mUnreadMessageListView;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_un_read_message);

        mUnreadMessage = (SettingActionBarItem) findViewById(R.id.sabi_unread_message);
        mUnreadMessage.setLeftTextOnclickListner(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                Intent intent = new Intent(UnReadMessageActivity.this, MessageActivity.class);
                startActivity(intent);
                finish();
            }
        });

        mSetAllMessageRead = (Button) findViewById(R.id.set_all_message_read);
    }

}
