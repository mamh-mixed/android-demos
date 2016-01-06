package com.cardinfolink.yunshouyin.core;

import android.os.AsyncTask;

import com.cardinfolink.yunshouyin.api.QuickPayApi;
import com.cardinfolink.yunshouyin.api.QuickPayApiImpl;
import com.cardinfolink.yunshouyin.api.QuickPayConfigStorage;
import com.cardinfolink.yunshouyin.api.QuickPayException;
import com.cardinfolink.yunshouyin.data.User;
import com.cardinfolink.yunshouyin.model.BankInfo;
import com.cardinfolink.yunshouyin.model.Message;
import com.cardinfolink.yunshouyin.model.ServerPacket;

import java.util.Map;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

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

    @Override
    public void registerAsync(final String username, final String password, final String invite, final QuickPayCallbackListener<Void> listener) {

        new AsyncTask<Void, Integer, AsyncTaskResult<Void>>() {
            @Override
            protected AsyncTaskResult<Void> doInBackground(Void... params) {
                try {
                    quickPayApi.register(username, password, invite);
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
        loginAsync(username, password, null, listener);
    }


    @Override
    public void loginAsync(final String username, final String password, final String deviceToken, final QuickPayCallbackListener<User> listener) {

        new AsyncTask<Void, Integer, AsyncTaskResult<User>>() {
            @Override
            protected AsyncTaskResult<User> doInBackground(Void... params) {
                try {
                    User user = quickPayApi.login(username, password, deviceToken);
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
    public void getHistoryBillsAsync(final User user, final String month, final String index, final String status, final QuickPayCallbackListener<ServerPacket> listener) {
        getHistoryBillsAsync(user, month, index, "50", status, listener);
    }

    @Override
    public void getHistoryBillsAsync(final User user, final String month, final String index, final String size, final String status, final QuickPayCallbackListener<ServerPacket> listener) {
        new AsyncTask<Void, Integer, AsyncTaskResult<ServerPacket>>() {
            @Override
            protected AsyncTaskResult<ServerPacket> doInBackground(Void... params) {
                try {
                    ServerPacket serverPacket = quickPayApi.getHistoryBills(user, month, index, size, status);
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
    public void getHistoryCouponsAsync(final User user, final String month, final String index, final String size, final QuickPayCallbackListener<ServerPacket> listener) {
        new AsyncTask<Void, Integer, AsyncTaskResult<ServerPacket>>() {
            @Override
            protected AsyncTaskResult<ServerPacket> doInBackground(Void... params) {
                try {
                    ServerPacket serverPacket = quickPayApi.getHistoryCoupons(user, month, index, size);
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


    //这个是通过条件来查找订单的，里面使用的是新的接口findOrder
    @Override
    public void findOrderAsync(final User user, final String index, final String size, final String recType, final String payType, final String txnStatus, final QuickPayCallbackListener<ServerPacket> listener) {
        new AsyncTask<Void, Integer, AsyncTaskResult<ServerPacket>>() {
            @Override
            protected AsyncTaskResult<ServerPacket> doInBackground(Void... params) {
                try {
                    ServerPacket serverPacket = quickPayApi.findOrder(user, index, size, recType, payType, txnStatus);
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

    //这个是通过订单号来查找订单的，里面使用的新的的接口findOrder
    @Override
    public void getOrderAsync(final User user, final String orderNum, final QuickPayCallbackListener<ServerPacket> listener) {
        new AsyncTask<Void, Integer, AsyncTaskResult<ServerPacket>>() {

            @Override
            protected AsyncTaskResult<ServerPacket> doInBackground(Void... params) {
                try {
                    ServerPacket serverPacket = quickPayApi.findOrder(user, orderNum);
                    return new AsyncTaskResult<ServerPacket>(serverPacket, null);
                } catch (QuickPayException ex) {
                    return new AsyncTaskResult<ServerPacket>(ex);
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
    public void getRefdAsync(final User user, final String orderNum, final QuickPayCallbackListener<ServerPacket> listener) {
        new AsyncTask<Void, Integer, AsyncTaskResult<ServerPacket>>() {

            @Override
            protected AsyncTaskResult<ServerPacket> doInBackground(Void... params) {
                try {
                    ServerPacket serverPacket = quickPayApi.getRefd(user, orderNum);//退款
                    return new AsyncTaskResult<ServerPacket>(serverPacket, null);
                } catch (QuickPayException ex) {
                    return new AsyncTaskResult<ServerPacket>(ex);
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
    public void forgetPasswordAsync(final String username, final QuickPayCallbackListener<ServerPacket> listener) {
        new AsyncTask<Void, Integer, AsyncTaskResult<ServerPacket>>() {

            @Override
            protected AsyncTaskResult<ServerPacket> doInBackground(Void... params) {
                try {
                    ServerPacket serverPacket = quickPayApi.forgetPassword(username);//忘记密码
                    return new AsyncTaskResult<ServerPacket>(serverPacket, null);
                } catch (QuickPayException ex) {
                    return new AsyncTaskResult<ServerPacket>(ex);
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
    public void getUploadTokenAsync(final User user, final QuickPayCallbackListener<String> listener) {
        new AsyncTask<Void, Integer, AsyncTaskResult<String>>() {

            @Override
            protected AsyncTaskResult<String> doInBackground(Void... params) {
                try {
                    String token = quickPayApi.getUploadToken(user);//获取token
                    return new AsyncTaskResult<String>(token);
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
    public void improveCertInfoAsync(final User user, final String certName, final String certAddr, final Map<String, String> imageMap, final QuickPayCallbackListener<Void> listener) {
        new AsyncTask<Void, Integer, AsyncTaskResult<Void>>() {

            @Override
            protected AsyncTaskResult<Void> doInBackground(Void... params) {
                try {
                    quickPayApi.improveCertInfo(user, certName, certAddr, imageMap);
                    return null;
                } catch (QuickPayException ex) {
                    return new AsyncTaskResult<Void>(ex);
                }
            }

            @Override
            protected void onPostExecute(AsyncTaskResult<Void> result) {
                if (result == null) {//这里为null表示成功
                    listener.onSuccess(null);
                } else {
                    listener.onFailure(result.getException());
                }
            }
        }.execute();
    }

    @Override
    public void pullinfoAsync(final String username, final String password, final String size, final String lasttime, final String maxtime, final QuickPayCallbackListener<ServerPacket> listener) {
        new AsyncTask<Void, Integer, AsyncTaskResult<ServerPacket>>() {

            @Override
            protected AsyncTaskResult<ServerPacket> doInBackground(Void... params) {
                try {
                    ServerPacket serverPacket = quickPayApi.pullinfo(username, password, size, lasttime, maxtime);
                    return new AsyncTaskResult<ServerPacket>(serverPacket);
                } catch (QuickPayException ex) {
                    return new AsyncTaskResult<ServerPacket>(ex);
                }
            }

            @Override
            protected void onPostExecute(AsyncTaskResult<ServerPacket> result) {
                if (result.getException() == null) {//这里为null表示成功
                    listener.onSuccess(result.getResult());
                } else {
                    listener.onFailure(result.getException());
                }
            }
        }.execute();
    }

    @Override
    public void updateMessageAsync(final String username, final String password, final Message[] messages, final String status, final QuickPayCallbackListener<ServerPacket> listener) {

        new AsyncTask<Void, Integer, AsyncTaskResult<ServerPacket>>() {

            @Override
            protected AsyncTaskResult<ServerPacket> doInBackground(Void... params) {
                try {
                    ServerPacket serverPacket = quickPayApi.updateMessage(username, password, status, messages);
                    return new AsyncTaskResult<ServerPacket>(serverPacket);
                } catch (QuickPayException ex) {
                    return new AsyncTaskResult<ServerPacket>(ex);
                }
            }

            @Override
            protected void onPostExecute(AsyncTaskResult<ServerPacket> result) {
                if (result.getException() == null) {//这里为null表示成功
                    listener.onSuccess(result.getResult());
                } else {
                    listener.onFailure(result.getException());
                }
            }
        }.execute();
    }
}
