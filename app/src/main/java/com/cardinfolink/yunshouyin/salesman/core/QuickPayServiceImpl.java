package com.cardinfolink.yunshouyin.salesman.core;

import android.os.AsyncTask;

import com.cardinfolink.yunshouyin.salesman.api.QuickPayApi;
import com.cardinfolink.yunshouyin.salesman.api.QuickPayApiImpl;
import com.cardinfolink.yunshouyin.salesman.api.QuickPayConfigStorage;
import com.cardinfolink.yunshouyin.salesman.api.QuickPayException;
import com.cardinfolink.yunshouyin.salesman.model.User;

import java.util.regex.Matcher;
import java.util.regex.Pattern;

/**
 * Created by mamh on 15-11-23.
 */
public class QuickPayServiceImpl implements QuickPayService {

    private QuickPayApi quickPayApi;
    private QuickPayConfigStorage quickPayConfigStorage;

    public QuickPayServiceImpl(QuickPayConfigStorage quickPayConfigStorage) {
        quickPayApi = new QuickPayApiImpl(quickPayConfigStorage);
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
     * @param quickPayCallbackListener
     */
    @Override
    public void loginAsync(final String username, final String password, final QuickPayCallbackListener<String> quickPayCallbackListener) {
        if (username.equals("")) {
            quickPayCallbackListener.onFailure(new QuickPayException("", "用户名不能为空!"));
            return;
        }

        if (password.equals("")) {
            quickPayCallbackListener.onFailure(new QuickPayException("", "密码不能为空!"));
            return;
        }

        new AsyncTask<Void, Integer, AsyncTaskResult<String>>() {
            @Override
            protected AsyncTaskResult<String> doInBackground(Void... params) {
                try {
                    String accessToken = quickPayApi.login(username, password);

                    return new AsyncTaskResult<String>(accessToken, null);
                } catch (QuickPayException ex) {

                    return new AsyncTaskResult<String>(null, ex);
                }
            }

            @Override
            protected void onPostExecute(AsyncTaskResult<String> stringAsyncTaskResult) {
                if (stringAsyncTaskResult.getException() != null) {
                    quickPayCallbackListener.onFailure(stringAsyncTaskResult.getException());
                } else {
                    quickPayCallbackListener.onSuccess(stringAsyncTaskResult.getResult());
                }
            }
        }.execute();
    }

    @Override
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

    @Override
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

    @Override
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

    @Override
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

    @Override
    public void registerUserAsync(final String email, final String password, final String password_repeat, final QuickPayCallbackListener<User> quickPayCallbackListener) {
        if (email.equals("")) {
            quickPayCallbackListener.onFailure(new QuickPayException("", "邮箱不能为空!"));
            return;
        }
        if (!checkEmail(email)) {
            quickPayCallbackListener.onFailure(new QuickPayException("", "邮箱格式不正确!"));
            return;
        }
        if (password.equals("")) {
            quickPayCallbackListener.onFailure(new QuickPayException("", "密码不能为空!"));
            return;
        }
        if (password.length() < 6) {
            quickPayCallbackListener.onFailure(new QuickPayException("", "密码不能小于六位!"));
            return;
        }
        if (!password.equals(password_repeat)) {
            quickPayCallbackListener.onFailure(new QuickPayException("", "确认密码不一致!"));
            return;
        }

        new AsyncTask<Void, Integer, AsyncTaskResult<User>>() {
            @Override
            protected AsyncTaskResult<User> doInBackground(Void... params) {
                try {
                    User user = quickPayApi.registerUser(email, password);
                    return new AsyncTaskResult<User>(user, null);
                } catch (QuickPayException ex) {

                    return new AsyncTaskResult<User>(null, ex);
                }
            }

            @Override
            protected void onPostExecute(AsyncTaskResult<User> stringAsyncTaskResult) {
                if (stringAsyncTaskResult.getException() != null) {
                    quickPayCallbackListener.onFailure(stringAsyncTaskResult.getException());
                } else {
                    quickPayCallbackListener.onSuccess(stringAsyncTaskResult.getResult());
                }
            }
        }.execute();
    }

    @Override
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

    @Override
    public void updateUser(final User user, QuickPayCallbackListener<User> quickPayCallbackListener) {
        try {
            User res = quickPayApi.updateUser(user);
            quickPayCallbackListener.onSuccess(res);
        } catch (QuickPayException ex) {
            quickPayCallbackListener.onFailure(ex);
        }
    }

    @Override
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

    @Override
    public void activateUser(final String username, final QuickPayCallbackListener<User> quickPayCallbackListener) {
        try {
            User user = quickPayApi.activateUser(username);
            quickPayCallbackListener.onSuccess(user);
        } catch (QuickPayException ex) {
            quickPayCallbackListener.onFailure(ex);
        }
    }
}
