package com.cardinfolink.yunshouyin.salesman.api;

import org.junit.Before;
import org.junit.Test;

public class QuickPayApiImplTest {
    private QuickPayApi quickPayApi;

    @Test
    public void testLogin() throws Exception {
    }

    @Before
    public void setUp() throws Exception {
        QuickPayConfigStorage quickPayConfigStorage = new QuickPayConfigStorage();
        quickPayConfigStorage.setAppKey("eu1dr0c8znpa43blzy1wirzmk8jqdaon");
        quickPayConfigStorage.setUrl("http://test.quick.ipay.so/app/tools");
        quickPayConfigStorage.setProxy_url("127.0.0.1");
        quickPayConfigStorage.setProxy_port(8888);

        quickPayApi = new QuickPayApiImpl(quickPayConfigStorage);

        String accessToken = quickPayApi.login("toolstest", "Yun#1016");
        System.out.println(accessToken);
    }

    @Test
    public void testGetUsers() throws Exception {
        quickPayApi.getUsers();
    }
}