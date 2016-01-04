package com.cardinfolink.yunshouyin.activity;

import android.content.Intent;
import android.os.AsyncTask;
import android.os.Bundle;
import android.view.View;
import android.widget.Button;
import android.widget.ListView;

import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.adapter.MessageAdapter;
import com.cardinfolink.yunshouyin.api.QuickPayException;
import com.cardinfolink.yunshouyin.core.QuickPayCallbackListener;
import com.cardinfolink.yunshouyin.data.MessageDB;
import com.cardinfolink.yunshouyin.data.SessonData;
import com.cardinfolink.yunshouyin.model.Message;
import com.cardinfolink.yunshouyin.model.ServerPacket;
import com.cardinfolink.yunshouyin.ui.SettingActionBarItem;
import com.handmark.pulltorefresh.library.PullToRefreshBase;
import com.handmark.pulltorefresh.library.PullToRefreshListView;

import java.text.SimpleDateFormat;
import java.util.ArrayList;
import java.util.Date;
import java.util.List;

public class MessageActivity extends BaseActivity {

    private SettingActionBarItem mActionBar;

    private PullToRefreshListView mUnreadMessageListView;

    private MessageAdapter mAdapter;

    private String lastTime;

    private MessageDB mMessageDB;

    private static final String PULL_UP = "pull_up";
    private static final String PULL_DOWN = "pull_down";

    private static final String PAGE_SIZE = "10";

    private static final String UNREAD_OR_UNDELETED = "0";
    private static final String READ_OR_UNDELETED = "1";
    private static final String UNREAD_OR_DELETED = "2";
    private static final String READ_OR_DELETE = "3";


    private Button mSetMessageRead;


    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_message);

        mMessageDB = new MessageDB(this);

        //初始化ActionBar
        mActionBar = (SettingActionBarItem) findViewById(R.id.action_bar);
        mActionBar.setLeftTextOnclickListner(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                finish();
            }
        });
        mActionBar.setRightTextOnclickListner(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                Intent intent = new Intent(MessageActivity.this, UnReadMessageActivity.class);
                startActivity(intent);
                finish();
            }
        });

        //获取最后一次推送的时间
        lastTime = mMessageDB.getLastTime(SessonData.loginUser.getUsername());

        //添加消息重置为已读状态事件
        mSetMessageRead = (Button) findViewById(R.id.set_all_message_read);
        mSetMessageRead.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                //查询所有未读消息
                MessageDB messageService = new MessageDB(MessageActivity.this);
                List<Message> unreadMessageList = messageService.getUnreadedMessages(SessonData.loginUser.getUsername());
                //与服务端同步消息状态（所有消息设置成已读状态）
                Message[] messages = new Message[unreadMessageList.size()];
                updateMessageStatus(unreadMessageList.toArray(messages));
                mActionBar.setRightText("未读消息(" + 0 + ")");
            }
        });

        mUnreadMessageListView = (PullToRefreshListView) findViewById(R.id.all_message_list_view);
        mUnreadMessageListView.setMode(PullToRefreshBase.Mode.BOTH);
        mUnreadMessageListView.setOnRefreshListener(new PullToRefreshBase.OnRefreshListener2<ListView>() {
            @Override
            public void onPullDownToRefresh(PullToRefreshBase<ListView> refreshView) {
                initMessageList(PULL_DOWN);
            }

            @Override
            public void onPullUpToRefresh(PullToRefreshBase<ListView> refreshView) {
                initMessageList(PULL_UP);
            }
        });

        initMessageList(null);
    }

    @Override
    protected void onResume() {
        super.onResume();
        //查询未读消息数量
        int count = mMessageDB.countUnreadedMessages(SessonData.loginUser.getUsername());
        mActionBar.setRightText("未读消息(" + count + ")");
    }


    private void updateMessageStatus(Message[] messages) {
        String username = SessonData.loginUser.getUsername();
        String password = SessonData.loginUser.getPassword();
        quickPayService.updateMessageAsync(username, password, messages, READ_OR_UNDELETED, new QuickPayCallbackListener<ServerPacket>() {
            @Override
            public void onSuccess(ServerPacket data) {
                String state = data.getState();
                if ("".equals(state) || "success".equals(state)) {
                    //批量更新本地消息为已读状态
                    Message message = new Message();
                    message.setStatus(READ_OR_UNDELETED);
                    message.setUsername(SessonData.loginUser.getUsername());
                    mMessageDB.setAllMessageReaded(message);
                    notifyChange(READ_OR_UNDELETED);
                } else {
                }
            }

            @Override
            public void onFailure(QuickPayException ex) {
            }
        });
    }

    /**
     * 填充消息数据
     *
     * @param type 操作类型：down-下拉 up-上拉 null-首次加载
     */
    private void setMessageToView(List<Message> messageListTemp, String type) {
        List<Message> messageList;
        if (type == null) {
            messageList = new ArrayList<>();
            messageList.addAll(0, messageListTemp);
            mAdapter = new MessageAdapter(this, messageList);
            mUnreadMessageListView.setAdapter(mAdapter);
        } else if (PULL_DOWN.equals(type)) {
            messageList = mAdapter.getMessageList();
            messageList.addAll(0, messageListTemp);
            mAdapter.notifyDataSetChanged();
        } else {
            messageList = mAdapter.getMessageList();
            messageList.addAll(messageList.size(), messageListTemp);
            mAdapter.notifyDataSetChanged();
        }
        if (messageList.size() > 0) {
            lastTime = messageList.get(messageList.size() - 1).getPushtime();
        }
        mUnreadMessageListView.onRefreshComplete();
    }


    private void notifyChange(String status) {
        List<Message> messageList = mAdapter.getMessageList();
        for (Message msg : messageList) {
            msg.setStatus(status);
        }
        mAdapter.notifyDataSetChanged();
        mUnreadMessageListView.onRefreshComplete();
    }

    /**
     * 拉取消息数据：首先从本地数据库加载，本地没有，从服务器获取
     */
    private void initMessageList(final String type) {
        String username = SessonData.loginUser.getUsername();
        String password = SessonData.loginUser.getPassword();
        if (type == null || PULL_DOWN.equals(type)) { //查询最新消息
            //获取系统当前时间
            SimpleDateFormat format = new SimpleDateFormat("yyyy-MM-dd HH:mm:ss");
            String currentTime = format.format(new Date());
            final String lastTime = mMessageDB.getLastTime(username);

            //从服务器端查询最新的消息
            quickPayService.pullinfoAsync(username, password, "0", lastTime, currentTime, new QuickPayCallbackListener<ServerPacket>() {
                @Override
                public void onSuccess(ServerPacket data) {
                    List<Message> messageList = new ArrayList<>();
                    if (data.getCount() > 0) {
                        Message[] messages = data.getMessage();
                        if (messages != null && messages.length > 0) {
                            for (Message message : messages) {
                                messageList.add(message);
                            }
                        }
                        //数据写入到本地
                        mMessageDB.add(messageList);
                        //更新未读消息数量
                        String unReadMessageText = mActionBar.getRightText();
                        int count = Integer.valueOf(unReadMessageText.substring(unReadMessageText.indexOf("(") + 1, unReadMessageText.length() - 1));
                        mActionBar.setRightText("未读消息(" + (count + messageList.size()) + ")");
                    }
                    if (data.getCount() == 0 && type == null) { //服务器没有获取数据，重本地读取
                        messageList = getLocalMessages(lastTime, null);
                    }
                    //更新界面内容
                    setMessageToView(messageList, type);
                }

                @Override
                public void onFailure(QuickPayException ex) {
                }
            });
        } else { //查询本地数据库消息
            new AsyncTask<Void, Integer, List<Message>>() {
                @Override
                protected List<Message> doInBackground(Void... params) {
                    try {
                        return getLocalMessages(lastTime, null);
                    } catch (QuickPayException ex) {
                        ex.printStackTrace();
                        return null;
                    }
                }

                @Override
                protected void onPostExecute(List<Message> messageList) {
                    setMessageToView(messageList, type);
                }
            }.execute();
        }
    }


    private List<Message> getLocalMessages(String pushTime, String status) {
        Message message = new Message();
        message.setUsername(SessonData.loginUser.getUsername());
        message.setPushtime(pushTime);
        message.setStatus(status);
        return mMessageDB.getLocalMessages(message, PAGE_SIZE);
    }
}
