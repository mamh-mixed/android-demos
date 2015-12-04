package com.cardinfolink.yunshouyin.api;

import com.cardinfolink.yunshouyin.data.User;
import com.cardinfolink.yunshouyin.model.BankInfo;
import com.cardinfolink.yunshouyin.model.ServerPacket;
import com.cardinfolink.yunshouyin.model.ServerPacketOrder;
import com.cardinfolink.yunshouyin.model.Txn;

public interface QuickPayApi {
    // user related
    void register(String username, String password);

    User login(String username, String password);

    void updatePassword(String username, String oldPassword, String newPassword);

    void forgetPassword(String username);

    void resetPassword(String username, String code, String newPassword);

    void activate(String username, String password);

    // bank related
    User improveInfo(User user);

    User updateInfo(User user); //这个和上面那个improveinfo很容易弄混。


    void increaseLimit(String username, String password, String payee, String phone_num, String email);

    BankInfo getBankInfo(User user);


    // txn related
    ServerPacket getHistoryBills(String username, String password, String clientid, String month, long index, String status);

    String getTotal(User user, String date);

    ServerPacketOrder getOrder(User user, String orderNum);

    ServerPacket getRefd(User user, String orderNum);
}
