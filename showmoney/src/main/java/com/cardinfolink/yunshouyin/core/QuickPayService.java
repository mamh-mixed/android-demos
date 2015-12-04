package com.cardinfolink.yunshouyin.core;

import com.cardinfolink.yunshouyin.data.User;
import com.cardinfolink.yunshouyin.model.BankInfo;
import com.cardinfolink.yunshouyin.model.ServerPacket;
import com.cardinfolink.yunshouyin.model.ServerPacketOrder;

public interface QuickPayService {
    void registerAsync(String username, String password, String password_repeat, QuickPayCallbackListener<Void> quickPayCallbackListener);

    void loginAsync(String username, String password, QuickPayCallbackListener<User> quickPayCallbackListener);

    void activateAsync(String username, String password, QuickPayCallbackListener<Void> quickPayCallbackListener);

    //user logged in
    void updateInfoAsync(User user, QuickPayCallbackListener<User> quickPayCallbackListener);

    void increaseLimitAsync(String payee, String phone_num, String email, QuickPayCallbackListener<Void> quickPayCallbackListener);

    void getBankInfoAsync(QuickPayCallbackListener<BankInfo> quickPayCallbackListener);

    void updatePasswordAsync(String oldPassword, String newPassword, String newPassword_repeat, QuickPayCallbackListener<Void> quickPayCallbackListener);

    void getHistoryBillsAsync(String month, long index, String status, QuickPayCallbackListener<ServerPacket> quickPayCallbackListener);

    void getTotalAsync(User user, String date, QuickPayCallbackListener<String> listener);

    void getOrderAsync(User user, String orderNum, QuickPayCallbackListener<ServerPacketOrder> listener);
}
