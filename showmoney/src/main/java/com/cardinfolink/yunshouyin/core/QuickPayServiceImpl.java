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
import com.cardinfolink.yunshouyin.model.ServerPacketOrder;

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
    public void registerAsync(final String username, final String password, String password_repeat, final QuickPayCallbackListener<Void> listener) {

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
            protected void onPostExecute(AsyncTaskResult<Void> result) {
                if (result == null) {
                    listener.onSuccess(null);
                } else if (result.getException() != null) {
                    listener.onFailure(result.getException());
                }
            }
        }.execute();
    }

    @Override
    public void loginAsync(final String username, final String password, final QuickPayCallbackListener<User> listener) {

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
            protected void onPostExecute(AsyncTaskResult<User> result) {
                if (result.getException() != null) {
                    listener.onFailure(result.getException());
                } else {
                    listener.onSuccess(result.getResult());
                }
            }
        }.execute();
    }

    @Override
    public void activateAsync(final String username, final String password, final QuickPayCallbackListener<Void> listener) {
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
            protected void onPostExecute(AsyncTaskResult<Void> result) {
                if (result == null) {
                    listener.onSuccess(null);
                } else if (result.getException() != null) {
                    listener.onFailure(result.getException());
                }
            }
        }.execute();
    }

    @Override
    public void improveInfoAsync(final User user, final QuickPayCallbackListener<User> listener) {

        new AsyncTask<Void, Integer, AsyncTaskResult<User>>() {
            @Override
            protected AsyncTaskResult<User> doInBackground(Void... params) {
                try {
                    User newUser = quickPayApi.improveInfo(user);
                    return new AsyncTaskResult<User>(user, null);
                } catch (QuickPayException ex) {
                    return new AsyncTaskResult<User>(null, ex);
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

    @Override
    public void updateInfoAsync(final User user, final QuickPayCallbackListener<User> listener) {
        new AsyncTask<Void, Integer, AsyncTaskResult<User>>() {
            @Override
            protected AsyncTaskResult<User> doInBackground(Void... params) {
                try {
                    //注意这里和上面那个improveInfoAsync（）里面的不同。
                    User newUser = quickPayApi.updateInfo(user);
                    return new AsyncTaskResult<User>(user, null);
                } catch (QuickPayException ex) {
                    return new AsyncTaskResult<User>(null, ex);
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

    @Override
    public void increaseLimitAsync(final String payee, final String phone_num, final String email, final QuickPayCallbackListener<Void> listener) {
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
            protected void onPostExecute(AsyncTaskResult<Void> result) {
                if (result.getException() != null) {
                    listener.onFailure(result.getException());
                } else {
                    listener.onSuccess(null);
                }
            }
        }.execute();
    }

    @Override
    public void getBankInfoAsync(final User user, final QuickPayCallbackListener<BankInfo> listener) {
        new AsyncTask<Void, Integer, AsyncTaskResult<BankInfo>>() {
            @Override
            protected AsyncTaskResult<BankInfo> doInBackground(Void... params) {
                try {
                    BankInfo bankInfo = quickPayApi.getBankInfo(user);
                    return new AsyncTaskResult<BankInfo>(bankInfo);
                } catch (QuickPayException ex) {

                    return new AsyncTaskResult<BankInfo>(null, ex);
                }
            }

            @Override
            protected void onPostExecute(AsyncTaskResult<BankInfo> result) {
                if (result.getException() != null) {
                    listener.onFailure(result.getException());
                } else {
                    listener.onSuccess(result.getResult());
                }
            }
        }.execute();
    }

    @Override
    public void updatePasswordAsync(final String username, final String oldPassword, final String newPassword, final QuickPayCallbackListener<Void> listener) {
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
            protected void onPostExecute(AsyncTaskResult<Void> result) {
                if (result == null) {//等于null表示成功
                    listener.onSuccess(null);
                } else if (result.getException() != null) {
                    listener.onFailure(result.getException());
                }
            }
        }.execute();
    }

    @Override
    public void getHistoryBillsAsync(final String month, final long index, final String status, final QuickPayCallbackListener<ServerPacket> listener) {
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
            protected void onPostExecute(AsyncTaskResult<ServerPacket> result) {
                if (result.getException() != null) {
                    listener.onFailure(result.getException());
                } else {
                    listener.onSuccess(result.getResult());
                }
            }
        }.execute();
    }


    @Override
    public void getTotalAsync(final User user, final String date, final QuickPayCallbackListener<String> listener) {
        new AsyncTask<Void, Integer, AsyncTaskResult<String>>() {
            @Override
            protected AsyncTaskResult<String> doInBackground(Void... params) {
                try {
                    String total = quickPayApi.getTotal(user, date);
                    return new AsyncTaskResult<String>(total);
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

    @Override
    public void getOrderAsync(final User user, final String orderNum, final QuickPayCallbackListener<ServerPacketOrder> listener) {
        new AsyncTask<Void, Integer, AsyncTaskResult<ServerPacketOrder>>() {

            @Override
            protected AsyncTaskResult<ServerPacketOrder> doInBackground(Void... params) {
                ServerPacketOrder serverPacket = quickPayApi.getOrder(user, orderNum);
                return new AsyncTaskResult<ServerPacketOrder>(serverPacket, null);
            }

            @Override
            protected void onPostExecute(AsyncTaskResult<ServerPacketOrder> result) {
                if (result.getException() != null) {
                    listener.onFailure(result.getException());
                } else {
                    listener.onSuccess(result.getResult());
                }
            }
        }.execute();
    }

    @Override
    public void getRefdAsync(final User user, final String orderNum, final QuickPayCallbackListener<ServerPacket> listener) {
        new AsyncTask<Void, Integer, AsyncTaskResult<ServerPacket>>() {

            @Override
            protected AsyncTaskResult<ServerPacket> doInBackground(Void... params) {
                ServerPacket serverPacket = quickPayApi.getRefd(user, orderNum);//退款
                return new AsyncTaskResult<ServerPacket>(serverPacket, null);
            }

            @Override
            protected void onPostExecute(AsyncTaskResult<ServerPacket> result) {
                if (result.getException() != null) {
                    listener.onFailure(result.getException());
                } else {
                    listener.onSuccess(result.getResult());
                }
            }
        }.execute();
    }
}
