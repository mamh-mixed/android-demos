package com.cardinfolink.yunshouyin.api;

import com.cardinfolink.yunshouyin.data.User;
import com.cardinfolink.yunshouyin.model.BankInfo;
import com.cardinfolink.yunshouyin.model.ServerPacket;
import com.cardinfolink.yunshouyin.model.Txn;

public interface QuickPayApi {
    // user related
    void register(String username, String password);
    User login(String username, String password);
    void updatePassword(String username, String oldPassword, String newPassword);
    void forgetPassword(String username);
    void resetPassword(String username,String code, String newPassword);
    void activate(String username, String password);

    // bank related
    void updateInfo(String username, String password, String province, String city, String bank_open, String branch_bank, String bankNo, String payee, String payee_card, String phone_num);
    void increaseLimit(String username, String password, String payee, String phone_num, String email);
    BankInfo getBankInfo(String username, String password);


    // txn related
    // Txn getOrder(String username, String password, String orderNum, String clientId);
    ServerPacket getHistoryBills(String username, String password, String clientid, String month, long index, String status);
}
