package com.cardinfolink.yunshouyin.salesman.core;


import com.cardinfolink.yunshouyin.salesman.model.User;

public interface QuickPayService {
    void loginAsync(String username, String password, QuickPayCallbackListener<String> quickPayCallbackListener);

    void getUsersAsync(QuickPayCallbackListener<User[]> quickPayCallbackListener);

    void getQrPostUrlAsync(String merchantId, String imageType, QuickPayCallbackListener<String> quickPayCallbackListener);


    void getUploadToken(QuickPayCallbackListener<String> quickPayCallbackListener);

    void getUploadTokenAsync(QuickPayCallbackListener<String> quickPayCallbackListener);

    void registerUserAsync(String email, String password, String password_repeat, QuickPayCallbackListener<User> quickPayCallbackListener);

    void updateUserAsync(User user, QuickPayCallbackListener<User> quickPayCallbackListener);

    void updateUser(User user, QuickPayCallbackListener<User> quickPayCallbackListener);

    void activateUserAsync(String username, QuickPayCallbackListener<User> quickPayCallbackListener);

    void activateUser(String username, QuickPayCallbackListener<User> quickPayCallbackListener);
}