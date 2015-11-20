package com.cardinfolink.yunshouyin.core;

import com.cardinfolink.yunshouyin.data.User;
import com.cardinfolink.yunshouyin.model.BankInfo;
import com.cardinfolink.yunshouyin.model.ServerPacket;

public interface QuickPayService {
    void registerAsync(String username, String password, String password_repeat, QuickPayCallbackListener<Void> quickPayCallbackListener);
    void loginAsync(String username, String password, QuickPayCallbackListener<User> quickPayCallbackListener);
    void activateAsync(String username, String password, QuickPayCallbackListener<Void> quickPayCallbackListener);

    //user logged in
    void updateInfoAsync(String province, String city, String bank_open, String branch_bank, String bankNo, String payee, String payee_card, String phone_num, QuickPayCallbackListener<Void> quickPayCallbackListener);
    void increaseLimitAsync(String payee, String phone_num, String email, QuickPayCallbackListener<Void> quickPayCallbackListener);
    void getBankInfoAsync(QuickPayCallbackListener<BankInfo> quickPayCallbackListener);
    void updatePasswordAsync(String oldPassword, String newPassword, String newPassword_repeat, QuickPayCallbackListener<Void> quickPayCallbackListener);

    void getHistoryBillsAsync(String month, long index, String status, QuickPayCallbackListener<ServerPacket> quickPayCallbackListener);
}
