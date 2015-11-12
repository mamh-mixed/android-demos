package com.cardinfolink.yunshouyin.core;


import com.cardinfolink.yunshouyin.model.Bank;
import com.cardinfolink.yunshouyin.model.City;
import com.cardinfolink.yunshouyin.model.SubBank;

import java.util.List;
import java.util.Map;

//TODO: 参照QuickPayServiceImpl,在这里加入缓存效果
public class BankDataServiceImpl implements BankDataService {
    @Override
    public void getProvince(QuickPayCallbackListener<List<String>> quickPayCallbackListener) {

    }

    @Override
    public void getCity(String province, QuickPayCallbackListener<List<City>> quickPayCallbackListener) {

    }

    @Override
    public void getBank(QuickPayCallbackListener<Map<String, Bank>> quickPayCallbackListener) {

    }

    @Override
    public void search(String city_code, String bank_id, QuickPayCallbackListener<List<SubBank>> quickPayCallbackListener) {

    }
}
