package com.cardinfolink.yunshouyin.adapter;

import android.content.Context;
import android.content.Intent;
import android.os.Bundle;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.BaseAdapter;
import android.widget.LinearLayout;
import android.widget.TextView;

import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.activity.MessageDetailActivity;
import com.cardinfolink.yunshouyin.model.Message;

import java.util.List;

public class MessageAdapter extends BaseAdapter {
    public final static String SER_KEY = "com.cardinfolink.yunshouyin.adapter.ser";
    private List<Message> messageList;
    private Context mContext;

    public MessageAdapter(Context context, List<Message> messageList) {
        this.messageList = messageList;
        this.mContext = context;
    }

    @Override
    public int getCount() {
        return messageList.size();
    }

    @Override
    public Object getItem(int position) {
        return messageList.get(position);
    }

    @Override
    public long getItemId(int position) {
        return position;
    }

    public List<Message> getMessageList() {
        return messageList;
    }

    @Override
    public View getView(int position, View view, ViewGroup group) {
        ViewHolder holder = null;
        if (view == null) {
            holder = new ViewHolder();
            view = LayoutInflater.from(mContext).inflate(R.layout.message_list_item, null);

            holder.messageContent = (TextView) view.findViewById(R.id.message_content);
            holder.messageTime = (TextView) view.findViewById(R.id.message_time);
            holder.messageLinearLayout = (LinearLayout) view.findViewById(R.id.ll_message);
            view.setTag(holder);
        } else {
            holder = (ViewHolder) view.getTag();
        }
        final Message message = messageList.get(position);

        holder.messageContent.setText(message.getTitle());
        holder.messageTime.setText(message.getPushtime());

        if (!"0".equals(message.getStatus())) {
            holder.messageContent.setTextColor(mContext.getResources().getColor(R.color.message_read));
            holder.messageTime.setTextColor(mContext.getResources().getColor(R.color.message_read));
        } else {
            holder.messageContent.setTextColor(mContext.getResources().getColor(R.color.message_unread));
            holder.messageTime.setTextColor(mContext.getResources().getColor(R.color.message_unread));
        }
        final ViewHolder finalHolder = holder;
        holder.messageLinearLayout.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                finalHolder.messageContent.setTextColor(mContext.getResources().getColor(R.color.message_read));
                finalHolder.messageTime.setTextColor(mContext.getResources().getColor(R.color.message_read));

                Intent intent = new Intent(mContext, MessageDetailActivity.class);
                Bundle bundle = new Bundle();
                bundle.putSerializable(SER_KEY, message);
                intent.putExtras(bundle);

                mContext.startActivity(intent);
            }
        });

        return view;
    }

    private static final class ViewHolder {
        TextView messageContent;
        TextView messageTime;
        LinearLayout messageLinearLayout;
    }
}
