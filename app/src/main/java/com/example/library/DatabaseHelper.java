package com.example.library;

import android.content.Context;
import android.database.sqlite.SQLiteDatabase;
import android.database.sqlite.SQLiteDatabase.CursorFactory;
import android.database.sqlite.SQLiteOpenHelper;

public class DatabaseHelper extends SQLiteOpenHelper {
    private static final String DB_NAME = "eating.db";

    private static final String CREATEDAILYEAT = "CREATE TABLE 'ew_dailyeat' " +
            "('dailyeat_id' INTEGER NOT NULL," +
            "'dailyeat_meals' TEXT," +
            "'dailyeat_datetime' TEXT," +
            "'dailyeat_restaurant_id' INTEGER," +
            "'dailyeat_menu_id' INTEGER," +
            "'dailyeat_description' TEXT);";

    private static final String CREATERESTAURANT = "CREATE TABLE 'ew_restaurant' " +
            "('restaurant_id' INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL," +
            "'restaurant_name' TEXT," +
            "'restaurant_address' TEXT," +
            "'restaurant_phone' TEXT," +
            "'restaurant_description' TEXT);";

    private static final String CREATEMENU = "CREATE TABLE 'ew_menu'" +
            " ('menu_id'  INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL," +
            "'menu_name' TEXT," +
            "'restaurant_id' INTEGER," +
            " 'menu_description' TEXT);";

    private static final int VERSION = 1;

    public DatabaseHelper(Context context) {
        super(context, DB_NAME, null, VERSION);
    }

    public DatabaseHelper(Context context, String name, CursorFactory factory, int version) {
        super(context, name, factory, version);
    }

    @Override
    public void onCreate(SQLiteDatabase db) {
        db.beginTransaction();
        try {
            createTables(db);
            insertIntoTables(db);
            db.setTransactionSuccessful();
        } catch (Exception e) {
            e.printStackTrace();
        } finally {
            db.endTransaction();
        }
    }

    @Override
    public void onUpgrade(SQLiteDatabase arg0, int arg1, int arg2) {

    }

    /*
     * create database tables
     */
    private void createTables(SQLiteDatabase db) {
        db.execSQL(CREATEDAILYEAT);

        db.execSQL(CREATERESTAURANT);

        db.execSQL(CREATEMENU);
    }

    /*
     * insert some datas to tables
     */
    private void insertIntoTables(SQLiteDatabase db) {
        db.execSQL("insert into ew_restaurant(restaurant_name,restaurant_address,restaurant_phone,restaurant_description) " +
                "values(\"���㷻����\",\"�ֶ������潭·135Ū88��(��ʢ��·)\",\"(021)33922848\",\"��ǩ����ʳ ���� �Ž�\")");
    }
}
