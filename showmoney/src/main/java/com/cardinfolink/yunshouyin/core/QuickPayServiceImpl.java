package com.cardinfolink.yunshouyin.core;

import android.os.AsyncTask;

import com.cardinfolink.yunshouyin.api.QuickPayApi;
import com.cardinfolink.yunshouyin.api.QuickPayApiImpl;
import com.cardinfolink.yunshouyin.api.QuickPayConfigStorage;
import com.cardinfolink.yunshouyin.api.QuickPayException;
import com.cardinfolink.yunshouyin.data.SessonData;
import com.cardinfolink.yunshouyin.data.User;
import com.cardinfolink.yunshouyin.model.BankInfo;
import com.cardinfolink.yunshouyin.model.ServerPacket;

import java.util.regex.Matcher;
import java.util.regex.Pattern;

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

    @Override
    public void registerAsync(final String username, final String password, String password_repeat, final QuickPayCallbackListener<Void> quickPayCallbackListener) {
        //TODO: move validation here

        new AsyncTask<Void, Integer, AsyncTaskResult<Void>>() {
            @Override
            protected AsyncTaskResult<Void> doInBackground(Void... params) {
                try {
                    quickPayApi.register(username, password);
                    return null;
                } catch (QuickPayException ex) {

                    return new AsyncTaskResult<Void>(null, ex);
                }
            }

            @Override
            protected void onPostExecute(AsyncTaskResult<Void> stringAsyncTaskResult) {
                if (stringAsyncTaskResult.getException() != null) {
                    quickPayCallbackListener.onFailure(stringAsyncTaskResult.getException());
                } else {
                    quickPayCallbackListener.onSuccess(null);
                }
            }
        }.execute();
    }

    @Override
    public void loginAsync(final String username, final String password, final QuickPayCallbackListener<User> quickPayCallbackListener) {
        //TODO: move validation here

        new AsyncTask<Void, Integer, AsyncTaskResult<User>>() {
            @Override
            protected AsyncTaskResult<User> doInBackground(Void... params) {
                try {
                    User user = quickPayApi.login(username, password);
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
    public void activateAsync(final String username, final String password, final QuickPayCallbackListener<Void> quickPayCallbackListener) {
        //TODO: move validation here
        new AsyncTask<Void, Integer, AsyncTaskResult<Void>>() {
            @Override
            protected AsyncTaskResult<Void> doInBackground(Void... params) {
                try {
                    quickPayApi.activate(username, password);
                    return null;
                } catch (QuickPayException ex) {

                    return new AsyncTaskResult<Void>(null, ex);
                }
            }

            @Override
            protected void onPostExecute(AsyncTaskResult<Void> stringAsyncTaskResult) {
                if (stringAsyncTaskResult.getException() != null) {
                    quickPayCallbackListener.onFailure(stringAsyncTaskResult.getException());
                } else {
                    quickPayCallbackListener.onSuccess(null);
                }
            }
        }.execute();
    }

    @Override
    public void updateInfoAsync(final String province, final String city, final String bank_open, final String branch_bank, final String bankNo, final String payee, final String payee_card, final String phone_num, final QuickPayCallbackListener<Void> quickPayCallbackListener) {
        //TODO: move validation here
        //TODO: get username, password from login user
        final String username = "";
        final String password = "";

        new AsyncTask<Void, Integer, AsyncTaskResult<Void>>() {
            @Override
            protected AsyncTaskResult<Void> doInBackground(Void... params) {
                try {
                    quickPayApi.updateInfo(username, password, province, city, bank_open, branch_bank, bankNo, payee, payee_card, phone_num);
                    return null;
                } catch (QuickPayException ex) {

                    return new AsyncTaskResult<Void>(null, ex);
                }
            }

            @Override
            protected void onPostExecute(AsyncTaskResult<Void> stringAsyncTaskResult) {
                if (stringAsyncTaskResult.getException() != null) {
                    quickPayCallbackListener.onFailure(stringAsyncTaskResult.getException());
                } else {
                    quickPayCallbackListener.onSuccess(null);
                }
            }
        }.execute();
    }

    @Override
    public void increaseLimitAsync(final String payee, final String phone_num, final String email, final QuickPayCallbackListener<Void> quickPayCallbackListener) {
        //TODO: move validation here
        //TODO: get username, password from login user
        final String username = "";
        final String password = "";

        new AsyncTask<Void, Integer, AsyncTaskResult<Void>>() {
            @Override
            protected AsyncTaskResult<Void> doInBackground(Void... params) {
                try {
                    quickPayApi.increaseLimit(username, password, payee, phone_num, email);
                    return null;
                } catch (QuickPayException ex) {

                    return new AsyncTaskResult<Void>(null, ex);
                }
            }

            @Override
            protected void onPostExecute(AsyncTaskResult<Void> stringAsyncTaskResult) {
                if (stringAsyncTaskResult.getException() != null) {
                    quickPayCallbackListener.onFailure(stringAsyncTaskResult.getException());
                } else {
                    quickPayCallbackListener.onSuccess(null);
                }
            }
        }.execute();
    }

    @Override
    public void getBankInfoAsync(final QuickPayCallbackListener<BankInfo> quickPayCallbackListener) {
        //TODO: move validation here
        //TODO: get username, password from login user
        final String username = "";
        final String password = "";

        new AsyncTask<Void, Integer, AsyncTaskResult<Void>>() {
            @Override
            protected AsyncTaskResult<Void> doInBackground(Void... params) {
                try {
                    quickPayApi.getBankInfo(username, password);
                    return null;
                } catch (QuickPayException ex) {

                    return new AsyncTaskResult<Void>(null, ex);
                }
            }

            @Override
            protected void onPostExecute(AsyncTaskResult<Void> stringAsyncTaskResult) {
                if (stringAsyncTaskResult.getException() != null) {
                    quickPayCallbackListener.onFailure(stringAsyncTaskResult.getException());
                } else {
                    quickPayCallbackListener.onSuccess(null);
                }
            }
        }.execute();
    }

    @Override
    public void updatePasswordAsync(final String oldPassword, final String newPassword, String newPassword_repeat, final QuickPayCallbackListener<Void> quickPayCallbackListener) {
        //TODO: move validation here
        //TODO: get username, password from login user
        final String username = "";
        final String password = "";
        // TODO: compare with oldPassword

        new AsyncTask<Void, Integer, AsyncTaskResult<Void>>() {
            @Override
            protected AsyncTaskResult<Void> doInBackground(Void... params) {
                try {
                    quickPayApi.updatePassword(username, oldPassword, newPassword);
                    return null;
                } catch (QuickPayException ex) {

                    return new AsyncTaskResult<Void>(null, ex);
                }
            }

            @Override
            protected void onPostExecute(AsyncTaskResult<Void> stringAsyncTaskResult) {
                if (stringAsyncTaskResult.getException() != null) {
                    quickPayCallbackListener.onFailure(stringAsyncTaskResult.getException());
                } else {
                    quickPayCallbackListener.onSuccess(null);
                }
            }
        }.execute();
    }

    @Override
    public void getHistoryBillsAsync(final String month, final long index, final String status, final QuickPayCallbackListener<ServerPacket> quickPayCallbackListener) {
        final User loginUser = SessonData.loginUser;

        new AsyncTask<Void, Integer, AsyncTaskResult<ServerPacket>>() {
            @Override
            protected AsyncTaskResult<ServerPacket> doInBackground(Void... params) {
                try {
                    ServerPacket serverPacket = quickPayApi.getHistoryBills(loginUser.getUsername(), loginUser.getPassword(), loginUser.getClientid(), month, index, status);
                    return new AsyncTaskResult<ServerPacket>(serverPacket, null);
                } catch (QuickPayException ex) {
                    return new AsyncTaskResult<ServerPacket>(null, ex);
                }
            }

            @Override
            protected void onPostExecute(AsyncTaskResult<ServerPacket> serverPacketAsyncTaskResult) {
                if (serverPacketAsyncTaskResult.getException() != null) {
                    quickPayCallbackListener.onFailure(serverPacketAsyncTaskResult.getException());
                } else {
                    quickPayCallbackListener.onSuccess(serverPacketAsyncTaskResult.getResult());
                }
            }
        }.execute();
    }

}
