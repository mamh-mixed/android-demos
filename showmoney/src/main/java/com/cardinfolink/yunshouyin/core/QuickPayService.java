package com.cardinfolink.yunshouyin.core;

import com.cardinfolink.yunshouyin.data.User;
import com.cardinfolink.yunshouyin.model.BankInfo;
import com.cardinfolink.yunshouyin.model.Message;
import com.cardinfolink.yunshouyin.model.ServerPacket;

import java.util.Map;

public interface QuickPayService {
    void registerAsync(String username, String password, String invite, QuickPayCallbackListener<Void> listener);

    void loginAsync(String username, String password, QuickPayCallbackListener<User> listener);

    void loginAsync(String username, String password, String deviceToken, QuickPayCallbackListener<User> listener);

    void activateAsync(String username, String password, QuickPayCallbackListener<Void> listener);

    //user logged in
    void improveInfoAsync(User user, QuickPayCallbackListener<User> listener);


    void getBankInfoAsync(User user, QuickPayCallbackListener<BankInfo> quickPayCallbackListener);

    void updatePasswordAsync(String oldPassword, String newPassword, String newPassword_repeat, QuickPayCallbackListener<Void> listener);

    void getHistoryBillsAsync(User user, String month, String index, String status, QuickPayCallbackListener<ServerPacket> listener);

    void getHistoryBillsAsync(User user, String month, String index, String size, String status, QuickPayCallbackListener<ServerPacket> listener);


    void findOrderAsync(User user, String index, String size, String recType, String payType, String txnStatus, QuickPayCallbackListener<ServerPacket> listener);

    void getTotalAsync(User user, String date, QuickPayCallbackListener<String> listener);

    void getOrderAsync(User user, String orderNum, QuickPayCallbackListener<ServerPacket> listener);

    void getRefdAsync(User user, String orderNum, QuickPayCallbackListener<ServerPacket> listener);

    void forgetPasswordAsync(String user, QuickPayCallbackListener<ServerPacket> listener);

    //获取七牛上传图片时用的token
    void getUploadTokenAsync(User user, QuickPayCallbackListener<String> listener);

    void improveCertInfoAsync(User user, String certName, String certAddr, Map<String, String> imageMap, QuickPayCallbackListener<Void> listener);

    void pullinfoAsync(String username, String password, String size, String lasttime, String maxtime, QuickPayCallbackListener<ServerPacket> listener);

    void updateMessageAsync(String username, String password, Message[] messages, String status, QuickPayCallbackListener<ServerPacket> listener);
}
