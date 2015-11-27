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
                    "province_name text " +
                    ")";
    /**
     * 创建唯一约束
     * CREATE UNIQUE INDEX index_name
     * ON table_name ( column1, column2,...columnN);
     */
    private static final String CREATE_UNIQUE_INDEX_PROVINCE =
            "CREATE UNIQUE INDEX UNIQUE_INDEX_PROVINCE" +
                    " on province (province_name)";

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

    private static final String CREATE_UNIQUE_INDEX_CITY =
            "create unique index unique_index_city" +
                    " on city (id, city_code, province_code," +
                    "city_name, city_jb, city, province)";

    private static final String CREATE_BANK =
            "create table bank(" +
                    "_id integer primary key autoincrement, " +
                    "id text, " +
                    "bank_name text " +
                    ")";

    private static final String CREATE_UNIQUE_INDEX_BANK =
            "create unique index unique_index_bank" +
                    " on bank (id, bank_name)";
    //银行分行
    private static final String CREATE_BRANCH_BANK =
            "create table branchbank(" +
                    "_id integer primary key autoincrement, " +
                    "bank_name text, " + //银行名称
                    "city_code text, " +
                    "one_bank_no text, " + //一级行号
                    "two_bank_no text " + //1级行号
                    ")";
    private static final String CREATE_UNIQUE_INDEX_BRANCH_BANK =
            "create unique index unique_index_branch_bank" +
                    " on branchbank (bank_name, city_code, one_bank_no, two_bank_no, bank_id)";


    public SalesmanOpenHelper(Context context, String name, SQLiteDatabase.CursorFactory factory, int version) {
        super(context, name, factory, version);
    }

    @Override
    public void onCreate(SQLiteDatabase db) {
        db.execSQL(CREATE_PROVINCE);
        db.execSQL(CREATE_UNIQUE_INDEX_PROVINCE);

        db.execSQL(CREATE_CITY);
        db.execSQL(CREATE_UNIQUE_INDEX_CITY);

        db.execSQL(CREATE_BANK);
        db.execSQL(CREATE_UNIQUE_INDEX_BANK);

        db.execSQL(CREATE_BRANCH_BANK);
        db.execSQL(CREATE_UNIQUE_INDEX_BRANCH_BANK);
    }

    @Override
    public void onUpgrade(SQLiteDatabase db, int oldVersion, int newVersion) {
        //当数据库更新时回调这个方法
    }
}
