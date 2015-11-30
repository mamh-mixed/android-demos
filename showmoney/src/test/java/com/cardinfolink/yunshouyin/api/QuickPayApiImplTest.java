package com.cardinfolink.yunshouyin.api;

import com.cardinfolink.yunshouyin.data.User;

import org.junit.Before;
import org.junit.Test;

public class QuickPayApiImplTest {

    private QuickPayApi quickPayApi;

    @Before
    public void setUp() throws Exception {
        QuickPayConfigStorage quickPayConfigStorage = new QuickPayConfigStorage();
        quickPayConfigStorage.setAppKey("eu1dr0c8znpa43blzy1wirzmk8jqdaon");
        quickPayConfigStorage.setUrl("http://test.quick.ipay.so/app");
        quickPayConfigStorage.setProxyUrl("127.0.0.1");
        quickPayConfigStorage.setProxyPort(8888);

        quickPayApi = new QuickPayApiImpl(quickPayConfigStorage);
    }

    @Test
    public void testLogin() throws Exception {
        User user =  quickPayApi.login("453481716@qq.com", "123456");
    }

    @Test
    public void testRegister() throws Exception {
        quickPayApi.register("453481716@qqq.com", "123456");
    }

    @Test
    public void testUpdateInfo() throws Exception {
        quickPayApi.updateInfo("453481716@qq.com", "123456","上海","上海","上海浦发银行","张江支行","123456","哈哈过着","6228482371938777011","18516566509");
    }

    @Test
    public void testGetBankInfo() throws Exception {
        quickPayApi.getBankInfo("453481716@qq.com", "123456");
    }

    @Test
    public void testUpdatePassword() throws Exception {
        quickPayApi.updatePassword("453481716@qq.com", "1234567", "123456");
    }

    @Test
    public void testActivate() throws Exception {
        quickPayApi.activate("453481716@qq.com", "123456");
        quickPayApi.activate("453481716@qq.com", "1234567");
    }

    @Test
    public void testIncreaseLimit() throws Exception {
        quickPayApi.increaseLimit("453481716@qq.com", "123456", "Tom","18516566509","john.xu@cardinfolink.com");
    }

    @Test
    public void testForgetPassword() throws Exception {
        quickPayApi.forgetPassword("453481716@qq.com");
    }

    @Test
    public void testResetPassword() throws Exception {

    }
}
