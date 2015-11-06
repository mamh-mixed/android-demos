package com.cardinfolink.yunshouyin.salesman.core;

import com.cardinfolink.yunshouyin.salesman.api.QuickPayApi;
import com.cardinfolink.yunshouyin.salesman.api.QuickPayApiImpl;
import com.cardinfolink.yunshouyin.salesman.api.QuickPayConfigStorage;
import com.cardinfolink.yunshouyin.salesman.api.QuickPayException;
import com.cardinfolink.yunshouyin.salesman.model.User;

public class QuickPayService {
    private QuickPayApi quickPayApi;
    private QuickPayConfigStorage quickPayConfigStorage;

    public QuickPayService(QuickPayConfigStorage quickPayConfigStorage) {
        quickPayApi = new QuickPayApiImpl(quickPayConfigStorage);
        this.quickPayConfigStorage = quickPayConfigStorage;
    }

    /**
     * TODO: add validation
     *
     * @param username
     * @param password
     * @param quickPayCallbackListener
     */
    public void loginAsync(final String username, final String password, final QuickPayCallbackListener<String> quickPayCallbackListener) {
        new Thread(new Runnable() {
            @Override
            public void run() {
                try {
                    String accessToken = quickPayApi.login(username, password);
                    quickPayCallbackListener.onSuccess(accessToken);
                } catch (QuickPayException ex) {
                    quickPayCallbackListener.onFailure(ex);
                }
            }
        }).start();
    }

    public void getUsersAsync(final QuickPayCallbackListener<User[]> quickPayCallbackListener) {
        new Thread(new Runnable() {
            @Override
            public void run() {
                try {
                    User[] users = quickPayApi.getUsers();
                    quickPayCallbackListener.onSuccess(users);
                } catch (QuickPayException ex) {
                    quickPayCallbackListener.onFailure(ex);
                }
            }
        }).start();
    }

    public void getQrPostUrlAsync(final String merchantId, final String imageType, final QuickPayCallbackListener<String> quickPayCallbackListener) {
        new Thread(new Runnable() {
            @Override
            public void run() {
                try {
                    String url = quickPayApi.getQrPostUrl(merchantId, imageType);
                    quickPayCallbackListener.onSuccess(url);
                } catch (QuickPayException ex) {
                    quickPayCallbackListener.onFailure(ex);
                }
            }
        }).start();
    }

    public void getUploadToken(QuickPayCallbackListener<String> quickPayCallbackListener) {
        if (quickPayConfigStorage.getUploadToken() != null && "".equals(quickPayConfigStorage.getUploadToken())) {
            quickPayCallbackListener.onSuccess(quickPayConfigStorage.getUploadToken());
            return;
        }

        try {
            String uploadToken = quickPayApi.getUploadToken();
            // cache
            // TODO: what if uploadToken expires
            quickPayConfigStorage.setUploadToken(uploadToken);
            quickPayCallbackListener.onSuccess(uploadToken);

        } catch (QuickPayException ex) {
            quickPayCallbackListener.onFailure(ex);
        }
    }

    public void getUploadTokenAsync(final QuickPayCallbackListener<String> quickPayCallbackListener) {
        new Thread(new Runnable() {
            @Override
            public void run() {
                try {
                    String url = quickPayApi.getUploadToken();
                    quickPayCallbackListener.onSuccess(url);
                } catch (QuickPayException ex) {
                    quickPayCallbackListener.onFailure(ex);
                }
            }
        }).start();
    }

    public void registerUserAsync(final String username, final String password, final QuickPayCallbackListener<User> quickPayCallbackListener) {
        new Thread(new Runnable() {
            @Override
            public void run() {
                try {
                    User user = quickPayApi.registerUser(username, password);
                    quickPayCallbackListener.onSuccess(user);
                } catch (QuickPayException ex) {
                    quickPayCallbackListener.onFailure(ex);
                }
            }
        }).start();
    }

    public void updateUserAsync(final User user, final QuickPayCallbackListener<User> quickPayCallbackListener) {
        new Thread(new Runnable() {
            @Override
            public void run() {
                try {
                    User res = quickPayApi.updateUser(user);
                    quickPayCallbackListener.onSuccess(res);
                } catch (QuickPayException ex) {
                    quickPayCallbackListener.onFailure(ex);
                }
            }
        }).start();
    }

    public void updateUser(final User user, QuickPayCallbackListener<User> quickPayCallbackListener) {
        try {
            User res = quickPayApi.updateUser(user);
            quickPayCallbackListener.onSuccess(res);
        } catch (QuickPayException ex) {
            quickPayCallbackListener.onFailure(ex);
        }
    }

    public void activateUserAsync(final String username, final QuickPayCallbackListener<User> quickPayCallbackListener) {
        new Thread(new Runnable() {
            @Override
            public void run() {
                try {
                    User user = quickPayApi.activateUser(username);
                    quickPayCallbackListener.onSuccess(user);
                } catch (QuickPayException ex) {
                    quickPayCallbackListener.onFailure(ex);
                }
            }
        }).start();
    }

    public void activateUser(final String username, final QuickPayCallbackListener<User> quickPayCallbackListener) {
        try {
            User user = quickPayApi.activateUser(username);
            quickPayCallbackListener.onSuccess(user);
        } catch (QuickPayException ex) {
            quickPayCallbackListener.onFailure(ex);
        }
    }
}
