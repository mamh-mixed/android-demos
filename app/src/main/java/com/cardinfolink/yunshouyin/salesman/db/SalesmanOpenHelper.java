package com.cardinfolink.yunshouyin.salesman.db;

import android.content.Context;
import android.database.sqlite.SQLiteDatabase;
import android.database.sqlite.SQLiteOpenHelper;

/**
 * Created by mamh on 15-11-26.
 */
public class SalesmanOpenHelper extends SQLiteOpenHelper {

    //创建省份的sql语句
    private static final String CREATE_PROVINCE =
            "create table province(" +
                    "_id integer primary key autoincrement, " +
                    "province_name text "+
                    ")";

    //创建市，县的sql语句,_id是表的id和城市县无关的
    private static final String CREATE_CITY =
            "create table city(" +
                    "_id integer primary key autoincrement, " +
                    "id text, " +
                    "city_code text, " +
                    "province_code text, " +
                    "city_name text, " +
                    "city_jb text, " +
                    "city text, " +
                    "province text " +
                    ")";

    private static final String CREATE_BANK=
            "create table bank(" +
                    "_id integer primary key autoincrement, " +
                    "id text, " +
                    "bank_name text " +
                    ")";

    //银行分行
    private static final String CREATE_BRANCH_BANK=
            "create table branchbank(" +
                    "_id integer primary key autoincrement, " +
                    "bank_name text, " +
                    "city_code text, " +
                    "one_bank_no text, " +
                    "two_bank_no text " +
                    ")";




    public SalesmanOpenHelper(Context context, String name, SQLiteDatabase.CursorFactory factory, int version) {
        super(context, name, factory, version);
    }


    @Override
    public void onCreate(SQLiteDatabase db) {
        db.execSQL(CREATE_PROVINCE);
        db.execSQL(CREATE_CITY);
        db.execSQL(CREATE_BANK);
        db.execSQL(CREATE_BRANCH_BANK);
    }

    @Override
    public void onUpgrade(SQLiteDatabase db, int oldVersion, int newVersion) {
        //当数据库更新时回调这个方法
    }
}
