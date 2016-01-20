package com.cardinfolink.yunshouyin.activity;

import android.os.Bundle;
import android.view.View;
import android.widget.TextView;

import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.adapter.MessageAdapter;
import com.cardinfolink.yunshouyin.api.QuickPayException;
import com.cardinfolink.yunshouyin.core.QuickPayCallbackListener;
import com.cardinfolink.yunshouyin.data.MessageDB;
import com.cardinfolink.yunshouyin.data.SessionData;
import com.cardinfolink.yunshouyin.model.Message;
import com.cardinfolink.yunshouyin.model.ServerPacket;
import com.cardinfolink.yunshouyin.ui.SettingActionBarItem;

public class MessageDetailActivity extends BaseActivity {
    private SettingActionBarItem mActionBar;

    private Message message;

    private static final String PULL_UP = "pull_up";
    private static final String PULL_DOWN = "pull_down";

    private static final String PAGE_SIZE = "10";

    private static final String UNREAD_OR_UNDELETED = "0";
    private static final String READ_OR_UNDELETED = "1";
    private static final String UNREAD_OR_DELETED = "2";
    private static final String READ_OR_DELETE = "3";


    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_message_detail);

        //初始化ActionBar
        mActionBar = (SettingActionBarItem) findViewById(R.id.action_bar);
        mActionBar.setLeftTextOnclickListner(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                finish();
            }
        });

        //得到当前消息对象
        message = (Message) this.getIntent().getSerializableExtra(MessageAdapter.SER_KEY);

        //设置消息界面
        TextView titleView = (TextView) findViewById(R.id.title);
        TextView pushTimeView = (TextView) findViewById(R.id.pushtime);
        TextView messageView = (TextView) findViewById(R.id.message);
        titleView.setText(message.getTitle());
        pushTimeView.setText(message.getPushtime());
        messageView.setText(message.getMessage());

        //更新消息状态：包括服务器端和本地客户端
        updateMessageStatus();
    }

    protected void updateMessageStatus() {
        String username = SessionData.loginUser.getUsername();
        String password = SessionData.loginUser.getPassword();
        message.setStatus("1");
        quickPayService.updateMessageAsync(username, password, new Message[]{message}, READ_OR_UNDELETED, new QuickPayCallbackListener<ServerPacket>() {
            @Override
            public void onSuccess(ServerPacket data) {
                String state = data.getState();
                if ("".equals(state) || "success".equals(state)) {
                    //更新本地消息状态
                    MessageDB messageService = new MessageDB(MessageDetailActivity.this);
                    message.setStatus(READ_OR_UNDELETED);
                    messageService.update(message);
                } else {
                }
            }

            @Override
            public void onFailure(QuickPayException ex) {
            }
        });
    }

}
