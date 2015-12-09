package com.cardinfolink.yunshouyin.core;

import com.cardinfolink.yunshouyin.data.User;
import com.cardinfolink.yunshouyin.model.BankInfo;
import com.cardinfolink.yunshouyin.model.ServerPacket;
import com.cardinfolink.yunshouyin.model.ServerPacketOrder;

public interface QuickPayService {
    void registerAsync(String username, String password, String password_repeat, QuickPayCallbackListener<Void> listener);

    void loginAsync(String username, String password, QuickPayCallbackListener<User> listener);

    void activateAsync(String username, String password, QuickPayCallbackListener<Void> listener);

    //user logged in
    void improveInfoAsync(User user, QuickPayCallbackListener<User> listener);

    void updateInfoAsync(User user, QuickPayCallbackListener<User> listener);

    void increaseLimitAsync(User user, QuickPayCallbackListener<Void> listener);

    void getBankInfoAsync(User user, QuickPayCallbackListener<BankInfo> quickPayCallbackListener);

    void updatePasswordAsync(String oldPassword, String newPassword, String newPassword_repeat, QuickPayCallbackListener<Void> listener);

    void getHistoryBillsAsync(User user, String month, String index, String status, QuickPayCallbackListener<ServerPacket> listener);

    void getTotalAsync(User user, String date, QuickPayCallbackListener<String> listener);

    void getOrderAsync(User user, String orderNum, QuickPayCallbackListener<ServerPacketOrder> listener);

    void getRefdAsync(User user, String orderNum, QuickPayCallbackListener<ServerPacket> listener);

    void forgetPasswordAsync(String user, QuickPayCallbackListener<ServerPacket> listener);
}
