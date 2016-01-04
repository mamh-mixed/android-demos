package com.cardinfolink.yunshouyin.data;

import android.content.Context;
import android.database.Cursor;
import android.database.sqlite.SQLiteDatabase;
import android.database.sqlite.SQLiteStatement;
import android.util.Log;

import com.cardinfolink.yunshouyin.model.Message;

import java.util.ArrayList;
import java.util.List;

public class MessageDB {

    private MessageSQLiteOpenHelper helper = null;

    public MessageDB(Context context) {
        this.helper = new MessageSQLiteOpenHelper(context);
    }

    /**
     * 批量添加消息
     */
    public void add(List<Message> messages) {
        SQLiteDatabase database = helper.getWritableDatabase();
        String sql = "insert into message(msgId,username,title,message,pushtime,updateTime,status) values(?,?,?,?,?,?,?)";
        database.beginTransaction();
        SQLiteStatement statement = database.compileStatement(sql);
        for (int i = 0; i < messages.size(); i++) {
            Message message = messages.get(i);
            statement.bindString(1, message.getMsgId());
            statement.bindString(2, message.getUsername());
            statement.bindString(3, message.getTitle());
            statement.bindString(4, message.getMessage());
            statement.bindString(5, message.getPushtime());
            statement.bindString(6, message.getUpdateTime());
            statement.bindString(7, message.getStatus());
            statement.execute();
            statement.clearBindings();
        }
        database.setTransactionSuccessful();
        database.endTransaction();
        database.close();
    }

    /**
     * 修改消息状态为已读或者删除（针对单条消息）
     */
    public void update(Message message) {
        SQLiteDatabase database = helper.getWritableDatabase();
        database.execSQL("update message set status=? where msgId=? and username=?",
                new Object[]{message.getStatus(), message.getMsgId(), message.getUsername()});
        database.close();
    }

    /**
     * 查询最近一次的推送时间
     */
    public String getLastTime(String username) {
        String lastTime = null;
        SQLiteDatabase database = helper.getReadableDatabase();
        Cursor cursor = database.rawQuery("select max(pushtime) from message where username=?", new String[]{username});
        if (cursor.moveToNext()) {
            lastTime = cursor.getString(0);
        }
        cursor.close();
        database.close();
        return lastTime;
    }

    /**
     * 查询本地历史数据
     */
    public List<Message> getLocalMessages(Message message, String size) {
        List<Message> messageList = new ArrayList<>();
        SQLiteDatabase database = helper.getReadableDatabase();
        Cursor cursor;
        if (message.getStatus() == null) { //查询全部消息（包括删除的）
            cursor = database.rawQuery("select * from message where username=? and pushtime<=? order by pushtime desc limit ?",
                    new String[]{message.getUsername(), message.getPushtime(), size});
        } else {
            cursor = database.rawQuery("select * from message where username=? and pushtime<? and status=? order by pushtime desc limit ?",
                    new String[]{message.getUsername(), message.getPushtime(), message.getStatus(), size});
        }
        Log.e(this.getClass().getName(), String.valueOf("数据条数：" + cursor.getCount()));
        while (cursor.moveToNext()) {
            String msgId = cursor.getString(cursor.getColumnIndex("msgId"));
            String username = cursor.getString(cursor.getColumnIndex("username"));
            String title = cursor.getString(cursor.getColumnIndex("title"));
            String messageContent = cursor.getString(cursor.getColumnIndex("message"));
            String pushTime = cursor.getString(cursor.getColumnIndex("pushtime"));
            String updateTime = cursor.getString(cursor.getColumnIndex("updateTime"));
            String status = cursor.getString(cursor.getColumnIndex("status"));
            message = new Message(msgId, username, title, messageContent, pushTime, updateTime, status);
            messageList.add(message);
        }
        cursor.close();
        database.close();
        return messageList;
    }

    /**
     * 查询所有未读消息
     */
    public List<Message> getUnreadedMessages(String username) {
        List<Message> messageList = new ArrayList<>();
        Message message;
        SQLiteDatabase database = helper.getReadableDatabase();
        Cursor cursor = database.rawQuery("select * from message where username=? and status='0'", new String[]{username});
        Log.e(this.getClass().getName(), String.valueOf("数据条数：" + cursor.getCount()));
        while (cursor.moveToNext()) {
            String msgId = cursor.getString(cursor.getColumnIndex("msgId"));
            username = cursor.getString(cursor.getColumnIndex("username"));
            String title = cursor.getString(cursor.getColumnIndex("title"));
            String messageContent = cursor.getString(cursor.getColumnIndex("message"));
            String pushTime = cursor.getString(cursor.getColumnIndex("pushtime"));
            String updateTime = cursor.getString(cursor.getColumnIndex("updateTime"));
            String status = cursor.getString(cursor.getColumnIndex("status"));
            message = new Message(msgId, username, title, messageContent, pushTime, updateTime, status);
            messageList.add(message);
        }
        cursor.close();
        database.close();
        return messageList;
    }

    /**
     * 查询未读消息数量
     */
    public int countUnreadedMessages(String username) {
        SQLiteDatabase database = helper.getReadableDatabase();
        Cursor cursor = database.rawQuery("select count(*) from message where username=? and status='0'", new String[]{username});
        int count = 0;
        if (cursor.moveToFirst()) {
            count = cursor.getInt(0);
        }
        cursor.close();
        database.close();
        return count;
    }

    /**
     * 删除消息（物理删除）
     */
    public void delete(Message message) {
        SQLiteDatabase database = helper.getWritableDatabase();
        database.execSQL("delete from message where msgId=? ", new Object[]{message.getMsgId()});
        database.close();
    }

    /**
     * 所有消息设置成已读状态
     */
    public void setAllMessageReaded(Message message) {
        SQLiteDatabase database = helper.getWritableDatabase();
        database.execSQL("update message set status=?,username=? ", new Object[]{message.getStatus(), message.getUsername()});
        database.close();
    }

}
