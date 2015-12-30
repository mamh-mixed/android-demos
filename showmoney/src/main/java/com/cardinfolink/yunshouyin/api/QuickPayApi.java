package com.cardinfolink.yunshouyin.api;

import com.cardinfolink.yunshouyin.data.User;
import com.cardinfolink.yunshouyin.model.BankInfo;
import com.cardinfolink.yunshouyin.model.ServerPacket;
import com.cardinfolink.yunshouyin.model.ServerPacketOrder;
import com.cardinfolink.yunshouyin.model.Txn;

import java.util.Map;

public interface QuickPayApi {
    // user related
    void register(String username, String password);

    void register(String username, String password, String invite);

    User login(String username, String password);

    void updatePassword(String username, String oldPassword, String newPassword);

    ServerPacket forgetPassword(String username);

    void resetPassword(String username, String code, String newPassword);

    void activate(String username, String password);

    // bank related
    User improveInfo(User user);

    User updateInfo(User user); //这个和上面那个improveinfo很容易弄混。


    void increaseLimit(User user);

    BankInfo getBankInfo(User user);


    // txn related
    ServerPacket getHistoryBills(User user, String month, String index, String status);

    ServerPacket getHistoryBills(User user, String month, String index, String size, String status);

    ServerPacket findOrder(User user, String orderNum);

    ServerPacket findOrder(User user, String index, String size, String recType, String payType, String txnStatus);

    ServerPacket findOrder(User user, String index, String size, String orderNum, String recType, String payType, String txnStatus);

    String getTotal(User user, String date);

    ServerPacketOrder getOrder(User user, String orderNum);

    ServerPacket getRefd(User user, String orderNum);

    //获取七牛的token的方法
    String getUploadToken(User user);

    void improveCertInfo(User user, Map<String, String> imageMap);
}
