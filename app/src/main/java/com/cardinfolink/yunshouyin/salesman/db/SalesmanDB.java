package com.cardinfolink.yunshouyin.salesman.db;

import android.content.ContentValues;
import android.content.Context;
import android.database.Cursor;
import android.database.SQLException;
import android.database.sqlite.SQLiteDatabase;

import com.cardinfolink.yunshouyin.salesman.model.Bank;
import com.cardinfolink.yunshouyin.salesman.model.City;
import com.cardinfolink.yunshouyin.salesman.model.Province;
import com.cardinfolink.yunshouyin.salesman.model.SubBank;

import java.util.ArrayList;
import java.util.List;

/**
 * Created by mamh on 15-11-26.
 */
public class SalesmanDB {
    public static final String DB_NAME = "sales_man";

    public static final int VERSION = 1;
    private static final String TAG = "SalesmanDB";

    private static SalesmanDB salesmanDB;

    private SQLiteDatabase db;

    private static final String PROVINCE_TABLE = "province";
    private static final String CITY_TABLE = "city";
    private static final String BANK_TABLE = "bank";
    private static final String BRANCH_BANK_TABLE = "branchbank";

    private SalesmanDB(Context context) {
        SalesmanOpenHelper dbHelper = new SalesmanOpenHelper(context, DB_NAME, null, VERSION);
        db = dbHelper.getWritableDatabase();
    }

    /**
     * 多线程读写
     * <p/>
     * SQLite实质上是将数据写入一个文件，通常情况下，在应用的包名下面都能找到xxx.db的文件，
     * 拥有root权限的手机，可以通过adb shell，看到data/data/packagename/databases/xxx.db这样的文件。
     * <p/>
     * 我们可以得知SQLite是文件级别的锁：多个线程可以同时读，但是同时只能有一个线程写。
     * Android提供了SqliteOpenHelper类，加入Java的锁机制以便调用。
     * <p/>
     * 如果多线程同时读写（这里的指不同的线程用使用的是不同的Helper实例），
     * 后面的就会遇到android.database.sqlite.SQLiteException: database is locked这样的异常。
     * 对于这样的问题，解决的办法就是keep single sqlite connection，
     * 保持单个SqliteOpenHelper实例，同时对所有数据库操作的方法添加synchronized关键字。
     */
    public synchronized static SalesmanDB getInstance(Context context) {
        if (salesmanDB == null) {
            salesmanDB = new SalesmanDB(context);
        }
        return salesmanDB;
    }


    /**
     * "create table province(" +
     * "_id integer primary key autoincrement, " +
     * "province_name text
     * )";
     */

    public void saveProvince(Province province) {
        if (province == null) {
            return;
        }

        ContentValues values = new ContentValues();
        values.put("province_name", province.getProvinceName());
        try {
            db.insertOrThrow(PROVINCE_TABLE, null, values);
        } catch (SQLException e) {

        }
    }

    public List<Province> loadProvince() {
        List<Province> list = new ArrayList<Province>();

        String[] columns = new String[]{"province_name"};//要从数据库表格中查询的列
        Cursor cursor = db.query(PROVINCE_TABLE, columns, null, null, null, null, null);
        if (cursor.moveToFirst()) {
            do {
                String provinceName = cursor.getString(cursor.getColumnIndex("province_name"));
                Province p = new Province(provinceName);
                list.add(p);
            } while (cursor.moveToNext());
        }
        if (cursor != null) {
            cursor.close();
        }

        return list;
    }

    /**
     * "create table city(" +
     * "_id integer primary key autoincrement, " +
     * "id text, " +
     * "city_code text, " +
     * "province_code text, " +
     * "city_name text, " +
     * "city_jb text, " +
     * "city text, " +
     * "province text " +
     * ")";
     */
    public void saveCity(City city) {
        if (city == null) {
            return;
        }

        ContentValues values = new ContentValues();
        values.put("id", city.getId());
        values.put("city_code", city.getCityCode());
        values.put("province_code", city.getProvinceCode());
        values.put("city_name", city.getCityName());
        values.put("city_jb", city.getCityJb());
        values.put("city", city.getCity());
        values.put("province", city.getProvince());
        try {
            db.insertOrThrow(CITY_TABLE, null, values);
        } catch (Exception e) {

        }
    }

    public List<City> loadCity(String whichProvince) {
        List<City> list = new ArrayList<City>();
        String[] selectArgs = new String[]{
                whichProvince
        };
        String selection = "province = ? ";
        Cursor cursor = db.query(CITY_TABLE, null, selection, selectArgs, null, null, null);
        if (cursor.moveToFirst()) {
            do {
                String id = cursor.getString(cursor.getColumnIndex("id"));
                String cityCode = cursor.getString(cursor.getColumnIndex("city_code"));
                String provinceCode = cursor.getString(cursor.getColumnIndex("province_code"));
                String cityName = cursor.getString(cursor.getColumnIndex("city_name"));
                String cityJb = cursor.getString(cursor.getColumnIndex("city_jb"));
                String city = cursor.getString(cursor.getColumnIndex("city"));
                String province = cursor.getString(cursor.getColumnIndex("province"));
                City c = new City(id, cityCode, provinceCode, cityName, cityJb, city, province);
                list.add(c);
            } while (cursor.moveToNext());
        }
        if (cursor != null) {
            cursor.close();
        }

        return list;
    }


    /**
     * "create table bank(" +
     * "_id integer primary key autoincrement, " +
     * "id text, " +
     * "bank_name text " +
     * ")";
     */
    public void saveBank(Bank bank) {
        if (bank == null) {
            return;
        }

        ContentValues values = new ContentValues();
        values.put("id", bank.getId());
        values.put("bank_name", bank.getBankName());
        try {
            db.insertOrThrow(BANK_TABLE, null, values);
        } catch (Exception e) {

        }
    }

    public List<Bank> loadBank() {
        List<Bank> list = new ArrayList<Bank>();

        Cursor cursor = db.query(BANK_TABLE, null, null, null, null, null, null);
        if (cursor.moveToFirst()) {
            do {
                String id = cursor.getString(cursor.getColumnIndex("id"));
                String bankName = cursor.getString(cursor.getColumnIndex("bank_name"));
                Bank bank = new Bank(id, bankName);
                list.add(bank);
            } while (cursor.moveToNext());
        }
        if (cursor != null) {
            cursor.close();
        }

        return list;
    }


    /**
     * "create table branchbank(" +
     * "_id integer primary key autoincrement, " +
     * "bank_name text, " +
     * "city_code text, " +
     * "one_bank_no text, " +
     * "two_bank_no text, " +
     * "bank_id " +
     * ")";
     */
    public void saveBranchBank(SubBank sBank) {
        if (sBank == null) {
            return;
        }

        ContentValues values = new ContentValues();
        values.put("bank_name", sBank.getBankName());
        values.put("city_code", sBank.getCityCode());
        values.put("one_bank_no", sBank.getOneBankNo());
        values.put("two_bank_no", sBank.getTwoBankNo());
        values.put("bank_id", sBank.getBankId());
        try {
            db.insertOrThrow(BRANCH_BANK_TABLE, null, values);
        } catch (Exception e) {

        }
    }

    public List<SubBank> loadBranchBank(String whichCityCode, String whichBankId) {
        List<SubBank> list = new ArrayList<SubBank>();
        String[] selectArgs = new String[]{
                whichCityCode, whichBankId
        };
        //通过cityCode 和 大银行号 bankId来查询分行的信息
        String selection = "city_code = ? and bank_id = ? ";
        Cursor cursor = db.query(BRANCH_BANK_TABLE, null, selection, selectArgs, null, null, null);
        if (cursor.moveToFirst()) {
            do {
                String bankName = cursor.getString(cursor.getColumnIndex("bank_name"));
                String cityCode = cursor.getString(cursor.getColumnIndex("city_code"));
                String oneBNo = cursor.getString(cursor.getColumnIndex("one_bank_no"));
                String twoBNo = cursor.getString(cursor.getColumnIndex("two_bank_no"));
                String bankId = cursor.getString(cursor.getColumnIndex("bank_id"));
                SubBank subBank = new SubBank(bankName, cityCode, oneBNo, twoBNo, bankId);
                list.add(subBank);
            } while (cursor.moveToNext());
        }
        if (cursor != null) {
            cursor.close();
        }

        return list;
    }

}
