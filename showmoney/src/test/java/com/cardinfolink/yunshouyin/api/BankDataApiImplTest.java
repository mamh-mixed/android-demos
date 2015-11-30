package com.cardinfolink.yunshouyin.api;

import com.cardinfolink.yunshouyin.model.Bank;
import com.cardinfolink.yunshouyin.model.City;
import com.cardinfolink.yunshouyin.model.SubBank;

import org.junit.Before;
import org.junit.Test;

import java.util.List;
import java.util.Map;

import static org.junit.Assert.*;

public class BankDataApiImplTest {
    private BankDataApi bankDataApi;
    private QuickPayConfigStorage quickPayConfigStorage;
    @Before
    public void setUp() throws Exception {
        QuickPayConfigStorage quickPayConfigStorage = new QuickPayConfigStorage();
        quickPayConfigStorage.setBankbaseKey("20e786206dcf4aae8a63fe34553fd274");
        quickPayConfigStorage.setBankbaseUrl("http://211.144.213.120:443/bdp");
        quickPayConfigStorage.setProxyUrl("127.0.0.1");
        quickPayConfigStorage.setProxyPort(8888);

        bankDataApi = new BankDataApiImpl(quickPayConfigStorage);
    }

    @Test
    public void testGetProvince() throws Exception {
        List<String> provinceList = bankDataApi.getProvince();
    }

    @Test
    public void testGetCity() throws Exception {
        List<City> cities = bankDataApi.getCity("浙江省");
    }

    @Test
    public void testSearch() throws Exception {
        List<SubBank> subBanks = bankDataApi.search("3310", "102");
    }

    @Test
    public void testGetBank() throws Exception {
        Map<String, Bank>  bankMap = bankDataApi.getBank();
    }
}