package com.cardinfolink.yunshouyin.core;

import com.cardinfolink.yunshouyin.model.Bank;
import com.cardinfolink.yunshouyin.model.City;
import com.cardinfolink.yunshouyin.model.Province;
import com.cardinfolink.yunshouyin.model.SubBank;

import java.util.List;
import java.util.Map;

public interface BankDataService {
    void getProvince(QuickPayCallbackListener<List<Province>> quickPayCallbackListener);

    void getCity(String province, QuickPayCallbackListener<List<City>> quickPayCallbackListener);

    void getBank(QuickPayCallbackListener<List<Bank>> quickPayCallbackListener);

    void getBranchBank(String city_code, String bank_id, QuickPayCallbackListener<List<SubBank>> quickPayCallbackListener);
}
