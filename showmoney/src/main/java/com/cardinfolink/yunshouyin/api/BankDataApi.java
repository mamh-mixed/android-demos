package com.cardinfolink.yunshouyin.api;

import com.cardinfolink.yunshouyin.model.Bank;
import com.cardinfolink.yunshouyin.model.City;
import com.cardinfolink.yunshouyin.model.Province;
import com.cardinfolink.yunshouyin.model.SubBank;

import java.util.List;
import java.util.Map;

public interface BankDataApi {
    List<Province> getProvince();

    List<City> getCity(String province);

    Map<String, Bank> getBank();

    List<SubBank> getBranchBank(String city_code, String bank_id);
}
