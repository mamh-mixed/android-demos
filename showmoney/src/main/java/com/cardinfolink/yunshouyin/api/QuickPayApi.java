package com.cardinfolink.yunshouyin.api;

import com.cardinfolink.yunshouyin.data.User;
import com.cardinfolink.yunshouyin.model.BankInfo;
import com.cardinfolink.yunshouyin.model.Message;
import com.cardinfolink.yunshouyin.model.ServerPacket;
import com.cardinfolink.yunshouyin.model.ServerPacketOrder;

import java.util.Map;

public interface QuickPayApi {
    // user related
    void register(String username, String password, String invite);

    User login(String username, String password, String deviceToken);

    void updatePassword(String username, String oldPassword, String newPassword);

    ServerPacket forgetPassword(String username);


    void activate(String username, String password);

    // bank related
    User improveInfo(User user);


    BankInfo getBankInfo(User user);


    // txn related
    ServerPacket getHistoryBills(User user, String month, String index, String size, String status);

    ServerPacket findOrder(User user, String orderNum);

    ServerPacket findOrder(User user, String index, String size, String recType, String payType, String txnStatus);

    ServerPacket findOrder(User user, String index, String size, String orderNum, String recType, String payType, String txnStatus);

    String getTotal(User user, String date);


    ServerPacket getSummaryDay(User user, String date, String reportType);

    ServerPacket getRefd(User user, String orderNum);

    //获取七牛的token的方法
    String getUploadToken(User user);

    void improveCertInfo(User user, String certName, String certAddr, Map<String, String> imageMap);

    ServerPacket pullinfo(String username, String password, String size, String lasttime, String maxtime);

    ServerPacket updateMessage(String username, String password, String status, Message[] messages);
}
