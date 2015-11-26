package com.cardinfolink.yunshouyin.salesman.db;

import android.content.ContentValues;
import android.content.Context;
import android.database.Cursor;
import android.database.sqlite.SQLiteDatabase;
import android.os.Build;

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
        db.insert(PROVINCE_TABLE, null, values);
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
        db.insert(CITY_TABLE, null, values);
    }

    public List<City> loadCity() {
        List<City> list = new ArrayList<City>();

        Cursor cursor = db.query(CITY_TABLE, null, null, null, null, null, null);
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
        db.insert(BANK_TABLE, null, values);
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
     * "two_bank_no text " +
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
        db.insert(BRANCH_BANK_TABLE, null, values);
    }

    public List<SubBank> loadBranchBank() {
        List<SubBank> list = new ArrayList<SubBank>();

        Cursor cursor = db.query(BRANCH_BANK_TABLE, null, null, null, null, null, null);
        if (cursor.moveToFirst()) {
            do {
                String bankName = cursor.getString(cursor.getColumnIndex("bank_name"));
                String cityCode = cursor.getString(cursor.getColumnIndex("city_code"));
                String oneBNo = cursor.getString(cursor.getColumnIndex("one_bank_no"));
                String twoBNo = cursor.getString(cursor.getColumnIndex("two_bnak_no"));
                SubBank subBank = new SubBank(bankName, cityCode, oneBNo, twoBNo);
                list.add(subBank);
            } while (cursor.moveToNext());
        }
        if (cursor != null) {
            cursor.close();
        }

        return list;
    }

}
