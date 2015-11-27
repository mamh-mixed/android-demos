package com.cardinfolink.yunshouyin.salesman.core;



import com.cardinfolink.yunshouyin.salesman.model.Bank;
import com.cardinfolink.yunshouyin.salesman.model.City;
import com.cardinfolink.yunshouyin.salesman.model.Province;
import com.cardinfolink.yunshouyin.salesman.model.SubBank;

import java.util.List;
import java.util.Map;

public interface BankDataService {

    void getProvince(QuickPayCallbackListener<List<Province>> quickPayCallbackListener);

    void getCity(String province, QuickPayCallbackListener<List<City>> quickPayCallbackListener);

    void getBank(QuickPayCallbackListener<Map<String, Bank>> quickPayCallbackListener);

    void getBranchBank(String city_code, String bank_id, QuickPayCallbackListener<List<SubBank>> quickPayCallbackListener);
}
