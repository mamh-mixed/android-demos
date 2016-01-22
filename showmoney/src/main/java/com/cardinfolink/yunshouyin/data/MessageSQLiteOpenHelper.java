package com.cardinfolink.yunshouyin.data;

import android.content.Context;
import android.database.sqlite.SQLiteDatabase;
import android.database.sqlite.SQLiteDatabase.CursorFactory;
import android.database.sqlite.SQLiteOpenHelper;

public class MessageSQLiteOpenHelper extends SQLiteOpenHelper {
    private static final String DB_NAME = "message.db";
    private static final int DB_VERSION = 1;

    private static final String CREATE_MESSAGE_TABLE = "" +
            "create table message (" +
                "msgId VARCHAR(50)," +
                "username VARCHAR(50)," +
                "title VARCHAR(50)," +
                "message VARCHAR(50)," +
                "pushtime VARCHAR(50)," +
                "updateTime VARCHAR(50)," +
                "status VARCHAR(1)" +
            ")";

    public MessageSQLiteOpenHelper(Context context, String name, CursorFactory factory, int version) {
        super(context, name, factory, version);
    }

    public MessageSQLiteOpenHelper(Context context) {
        this(context, DB_NAME, null, DB_VERSION);
    }

    @Override
    public void onCreate(SQLiteDatabase db) {
        db.execSQL(CREATE_MESSAGE_TABLE);
    }

    @Override
    public void onUpgrade(SQLiteDatabase db, int oldVersion, int newVersion) {
    }
}
