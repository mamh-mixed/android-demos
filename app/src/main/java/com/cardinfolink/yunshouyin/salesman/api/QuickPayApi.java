package com.cardinfolink.yunshouyin.salesman.api;

import com.cardinfolink.yunshouyin.salesman.model.User;

public interface QuickPayApi {
    /**
     * get accessToken
     * accessToken will be used in following request
     *
     * @param username
     * @param password
     * @return
     */
    String login(String username, String password);

    /**
     * get qiniu uploadToken
     *
     * @return
     */
    String getUploadToken();

    User[] getUsers();

    User registerUser(String username, String password);

    User updateUser(User user);

    User activateUser(String username);

    /**
     * imageType: Bill, Pay
     *
     * @param merchantId
     * @param imageType
     * @return
     */
    String getQrPostUrl(String merchantId, String imageType);
}
