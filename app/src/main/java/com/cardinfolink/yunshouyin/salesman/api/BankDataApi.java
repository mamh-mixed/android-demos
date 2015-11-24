package com.cardinfolink.yunshouyin.salesman.api;


import com.cardinfolink.yunshouyin.salesman.model.Bank;
import com.cardinfolink.yunshouyin.salesman.model.City;
import com.cardinfolink.yunshouyin.salesman.model.SubBank;

import java.util.List;
import java.util.Map;

public interface BankDataApi {
    List<String> getProvince();

    List<City> getCity(String province);

    Map<String, Bank> getBank();

    List<SubBank> search(String cityCode, String bankId);
}
