package com.cardinfolink.yunshouyin.activity;

import android.app.Activity;
import android.os.Bundle;

import com.cardinfolink.yunshouyin.core.QuickPayService;
import com.cardinfolink.yunshouyin.data.MessageDB;
import com.cardinfolink.yunshouyin.data.SessonData;
import com.cardinfolink.yunshouyin.data.User;
import com.cardinfolink.yunshouyin.model.Message;
import com.cardinfolink.yunshouyin.util.ShowMoneyApp;

import java.text.SimpleDateFormat;
import java.util.List;


public class BasicMessageActivity extends Activity {

    protected QuickPayService quickPayService;
    protected MessageDB messageDB;

    protected static final String PULL_UP = "pull_up";
    protected static final String PULL_DOWN = "pull_down";

    protected static final String PAGE_SIZE = "10";

    protected static final String UNREAD_OR_UNDELETED = "0";
    protected static final String READ_OR_UNDELETED = "1";
    protected static final String UNREAD_OR_DELETED = "2";
    protected static final String READ_OR_DELETE = "3";

    protected User user = SessonData.loginUser;

    protected SimpleDateFormat format = new SimpleDateFormat("yyyy-MM-dd HH:mm:ss");

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);

        ShowMoneyApp yunApplication = (ShowMoneyApp) getApplication();
        quickPayService = yunApplication.getQuickPayService();
        messageDB = new MessageDB(this);
    }

    protected List<Message> getLocalMessages(String pushTime, String status) {
        Message message = new Message();
        message.setUsername(user.getUsername());
        message.setPushtime(pushTime);
        message.setStatus(status);
        return messageDB.getLocalMessages(message, PAGE_SIZE);
    }

}
