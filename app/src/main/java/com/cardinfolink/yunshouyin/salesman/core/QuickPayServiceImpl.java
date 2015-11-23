package com.cardinfolink.yunshouyin.salesman.core;

import android.os.AsyncTask;
import android.text.TextUtils;

import com.cardinfolink.yunshouyin.salesman.api.QuickPayApi;
import com.cardinfolink.yunshouyin.salesman.api.QuickPayApiImpl;
import com.cardinfolink.yunshouyin.salesman.api.QuickPayConfigStorage;
import com.cardinfolink.yunshouyin.salesman.api.QuickPayException;
import com.cardinfolink.yunshouyin.salesman.model.User;

import java.util.regex.Matcher;
import java.util.regex.Pattern;

/**
 * QuickPayService接口的实现子类
 * Created by mamh on 15-11-23.
 */
public class QuickPayServiceImpl implements QuickPayService {

    private QuickPayApi quickPayApi;
    private QuickPayConfigStorage quickPayConfigStorage;

    public QuickPayServiceImpl(QuickPayConfigStorage quickPayConfigStorage) {
        this.quickPayApi = new QuickPayApiImpl(quickPayConfigStorage);
        this.quickPayConfigStorage = quickPayConfigStorage;
    }

    private static boolean checkEmail(String email) {
        boolean flag = false;
        try {
            String check = "^([a-z0-9A-Z]+[-|_|\\.]?)+[a-z0-9A-Z]@([a-z0-9A-Z]+(-[a-z0-9A-Z]+)?\\.)+[a-zA-Z]{2,}$";
            Pattern regex = Pattern.compile(check);
            Matcher matcher = regex.matcher(email);
            flag = matcher.matches();
        } catch (Exception e) {
            flag = false;
        }
        return flag;
    }

    /**
     * @param username
     * @param password
     * @param listener
     */
    @Override
    public void loginAsync(final String username, final String password, final QuickPayCallbackListener<String> listener) {
        if (TextUtils.isEmpty(username)) {
            listener.onFailure(new QuickPayException("", "用户名不能为空!"));
            return;
        }

        if (TextUtils.isEmpty(password)) {
            listener.onFailure(new QuickPayException("", "密码不能为空!"));
            return;
        }

        new AsyncTask<Void, Integer, AsyncTaskResult<String>>() {
            @Override
            protected AsyncTaskResult<String> doInBackground(Void... params) {
                try {
                    String accessToken = quickPayApi.login(username, password);
                    return new AsyncTaskResult<String>(accessToken);
                } catch (QuickPayException ex) {
                    return new AsyncTaskResult<String>(ex);
                }
            }

            @Override
            protected void onPostExecute(AsyncTaskResult<String> result) {
                if (result.getException() != null) {
                    listener.onFailure(result.getException());
                } else {
                    listener.onSuccess(result.getResult());
                }
            }
        }.execute();
    }

    /**
     *
     * @param listener
     */
    @Override
    public void getUsersAsync(final QuickPayCallbackListener<User[]> listener) {
        new Thread(new Runnable() {
            @Override
            public void run() {
                try {
                    User[] users = quickPayApi.getUsers();
                    listener.onSuccess(users);
                } catch (QuickPayException ex) {
                    listener.onFailure(ex);
                }
            }
        }).start();
    }

    /**
     *
     * @param merchantId
     * @param imageType
     * @param listener
     */
    @Override
    public void getQrPostUrlAsync(final String merchantId, final String imageType, final QuickPayCallbackListener<String> listener) {
        new Thread(new Runnable() {
            @Override
            public void run() {
                try {
                    String url = quickPayApi.getQrPostUrl(merchantId, imageType);
                    listener.onSuccess(url);
                } catch (QuickPayException ex) {
                    listener.onFailure(ex);
                }
            }
        }).start();
    }


    @Override
    public void getUploadToken(QuickPayCallbackListener<String> listener) {
        if (quickPayConfigStorage.getUploadToken() != null && "".equals(quickPayConfigStorage.getUploadToken())) {
            listener.onSuccess(quickPayConfigStorage.getUploadToken());
            return;
        }

        try {
            String uploadToken = quickPayApi.getUploadToken();
            quickPayConfigStorage.setUploadToken(uploadToken);
            listener.onSuccess(uploadToken);

        } catch (QuickPayException ex) {
            listener.onFailure(ex);
        }
    }

    @Override
    public void getUploadTokenAsync(final QuickPayCallbackListener<String> listener) {
        new Thread(new Runnable() {
            @Override
            public void run() {
                try {
                    String url = quickPayApi.getUploadToken();
                    listener.onSuccess(url);
                } catch (QuickPayException ex) {
                    listener.onFailure(ex);
                }
            }
        }).start();
    }

    /**
     *
     * @param email
     * @param password
     * @param passwordRepeat
     * @param listener
     */
    @Override
    public void registerUserAsync(final String email, final String password, final String passwordRepeat, final QuickPayCallbackListener<User> listener) {
        if (email.equals("")) {
            listener.onFailure(new QuickPayException("", "邮箱不能为空!"));
            return;
        }
        if (!checkEmail(email)) {
            listener.onFailure(new QuickPayException("", "邮箱格式不正确!"));
            return;
        }
        if (password.equals("")) {
            listener.onFailure(new QuickPayException("", "密码不能为空!"));
            return;
        }
        if (password.length() < 6) {
            listener.onFailure(new QuickPayException("", "密码不能小于六位!"));
            return;
        }
        if (!password.equals(passwordRepeat)) {
            listener.onFailure(new QuickPayException("", "确认密码不一致!"));
            return;
        }

        new AsyncTask<Void, Integer, AsyncTaskResult<User>>() {
            @Override
            protected AsyncTaskResult<User> doInBackground(Void... params) {
                try {
                    User user = quickPayApi.registerUser(email, password);
                    return new AsyncTaskResult<User>(user);
                } catch (QuickPayException ex) {

                    return new AsyncTaskResult<User>(ex);
                }
            }

            @Override
            protected void onPostExecute(AsyncTaskResult<User> result) {
                if (result.getException() != null) {
                    listener.onFailure(result.getException());
                } else {
                    listener.onSuccess(result.getResult());
                }
            }
        }.execute();
    }

    /**
     *
     * @param user
     * @param listener
     */
    @Override
    public void updateUserAsync(final User user, final QuickPayCallbackListener<User> listener) {
        new Thread(new Runnable() {
            @Override
            public void run() {
                try {
                    User res = quickPayApi.updateUser(user);
                    listener.onSuccess(res);
                } catch (QuickPayException ex) {
                    listener.onFailure(ex);
                }
            }
        }).start();
    }

    /**
     *
     * @param user
     * @param listener
     */
    @Override
    public void updateUser(final User user, QuickPayCallbackListener<User> listener) {
        try {
            User res = quickPayApi.updateUser(user);
            listener.onSuccess(res);
        } catch (QuickPayException ex) {
            listener.onFailure(ex);
        }
    }

    /**
     *
     * @param username
     * @param listener
     */
    @Override
    public void activateUserAsync(final String username, final QuickPayCallbackListener<User> listener) {
        new Thread(new Runnable() {
            @Override
            public void run() {
                try {
                    User user = quickPayApi.activateUser(username);
                    listener.onSuccess(user);
                } catch (QuickPayException ex) {
                    listener.onFailure(ex);
                }
            }
        }).start();
    }

    /**
     *
     * @param username
     * @param listener
     */
    @Override
    public void activateUser(final String username, final QuickPayCallbackListener<User> listener) {
        try {
            User user = quickPayApi.activateUser(username);
            listener.onSuccess(user);
        } catch (QuickPayException ex) {
            listener.onFailure(ex);
        }
    }
}
