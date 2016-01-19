package com.cardinfolink.yunshouyin.data;

import android.content.Context;
import android.database.sqlite.SQLiteDatabase;
import android.database.sqlite.SQLiteDatabase.CursorFactory;
import android.database.sqlite.SQLiteOpenHelper;

import com.cardinfolink.yunshouyin.util.Log;

public class MessageSQLiteOpenHelper extends SQLiteOpenHelper {

    public MessageSQLiteOpenHelper(Context context, String name, CursorFactory factory, int version) {
        super(context, name, factory, version);
    }

    public MessageSQLiteOpenHelper(Context context) {
        this(context, "message.db", null, 1);
    }

    @Override
    public void onCreate(SQLiteDatabase db) {
        db.execSQL("create table message (" +
                "msgId VARCHAR(50)," +
                "username VARCHAR(50)," +
                "title VARCHAR(50)," +
                "message VARCHAR(50)," +
                "pushtime VARCHAR(50)," +
                "updateTime VARCHAR(50)," +
                "status VARCHAR(1))");
    }

    @Override
    public void onUpgrade(SQLiteDatabase db, int oldVersion, int newVersion) {
        Log.d("MessageSQLiteOpenHelper", "数据库版本已经更新，新版本是：" + newVersion);
    }
}
