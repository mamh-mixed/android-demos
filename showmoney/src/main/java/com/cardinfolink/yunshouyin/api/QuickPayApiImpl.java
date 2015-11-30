package com.cardinfolink.yunshouyin.api;


import android.support.annotation.NonNull;
import android.text.TextUtils;

import com.cardinfolink.yunshouyin.data.User;
import com.cardinfolink.yunshouyin.model.BankInfo;
import com.cardinfolink.yunshouyin.model.ServerPacket;
import com.cardinfolink.yunshouyin.util.EncoderUtil;

import java.io.IOException;
import java.net.InetSocketAddress;
import java.net.Proxy;
import java.text.SimpleDateFormat;
import java.util.ArrayList;
import java.util.Collections;
import java.util.Date;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import java.util.SortedMap;
import java.util.TreeMap;

public class QuickPayApiImpl implements QuickPayApi {
    private static final String QUICK_PAY_SUCCESS = "success";

    private static final String TAG = "QuickPayApiImpl";
    protected QuickPayConfigStorage quickPayConfigStorage;
    protected PostEngine postEngine;

    public QuickPayApiImpl(QuickPayConfigStorage quickPayConfigStorage) {
        this.quickPayConfigStorage = quickPayConfigStorage;

        String proxyUrl = quickPayConfigStorage.getProxyUrl();
        int proxyPort = quickPayConfigStorage.getProxyPort();
        if (!TextUtils.isEmpty(proxyUrl)) {
            InetSocketAddress inetSocketAddress = new InetSocketAddress(proxyUrl, proxyPort);
            Proxy httpProxy = new Proxy(Proxy.Type.HTTP, inetSocketAddress);
            postEngine = new PostEngine(httpProxy);
        } else {
            postEngine = new PostEngine();
        }
    }

    @NonNull
    private String getTransTime() {
        Date now = new Date();
        SimpleDateFormat spf = new SimpleDateFormat("yyyyMMddHHmmss");
        return spf.format(now);
    }

    /**
     * 1. Sort by key name
     * 2. Prepare string, append app key
     * 3. Sign string
     *
     * @param params
     * @param signType
     * @return
     */
    private String createSign(Map<String, String> params, String signType) {
        SortedMap<String, String> sortedMap = new TreeMap<>();
        sortedMap.putAll(params);

        List<String> keys = new ArrayList<>(params.keySet());
        Collections.sort(keys);

        StringBuffer toSign = new StringBuffer();
        for (int i = 0; i < keys.size(); i++) {
            String key = keys.get(i);
            String value = params.get(key);
            if (null != value && !"".equals(value)) {
                if (i == keys.size() - 1) {
                    toSign.append(key + "=" + value);
                } else {
                    toSign.append(key + "=" + value + "&");
                }
            }
        }
        toSign.append(quickPayConfigStorage.getAppKey());
//        Log.d(TAG, "Raw string: " + toSign.toString());
        String sign = EncoderUtil.Encrypt(toSign.toString(), signType);
//        Log.d(TAG, "Signed string: " + sign);
        return sign;
    }

    /**
     * errors:
     * username_exist
     *
     * @param username
     * @param password
     */
    @Override
    public void register(String username, String password) {
        String url = quickPayConfigStorage.getUrl() + "/register";

        Map<String, String> params = new LinkedHashMap<>();
        params.put("username", username);
        password = EncoderUtil.Encrypt(password, "MD5");
        params.put("password", password);
        params.put("transtime", getTransTime());
        params.put("sign", createSign(params, "SHA-1"));

        try {
            String response = postEngine.post(url, params);
            ServerPacket serverPacket = ServerPacket.getServerPacketFrom(response);
            if (serverPacket.getState().equals(QUICK_PAY_SUCCESS)) {
                return;
            } else {
                throw new QuickPayException(serverPacket.getError());
            }
        } catch (IOException e) {
            throw new QuickPayException();
        }
    }

    /**
     * errors:
     * username_no_exist
     *
     * @param username
     * @param password
     * @return
     */
    @Override
    public User login(String username, String password) {
        String url = quickPayConfigStorage.getUrl() + "/login";

        Map<String, String> params = new LinkedHashMap<>();
        params.put("username", username);
        password = EncoderUtil.Encrypt(password, "MD5");
        params.put("password", password);
        params.put("transtime", getTransTime());
        params.put("sign", createSign(params, "SHA-1"));

        try {
            String response = postEngine.post(url, params);
            ServerPacket serverPacket = ServerPacket.getServerPacketFrom(response);
            if (serverPacket.getState().equals(QUICK_PAY_SUCCESS)) {
                return serverPacket.getUser();
            } else {
                throw new QuickPayException(serverPacket.getError());
            }
        } catch (IOException e) {
            throw new QuickPayException();
        }
    }

    @Override
    public void updatePassword(String username, String oldPassword, String newPassword) {
        String url = quickPayConfigStorage.getUrl() + "/updatepassword";

        Map<String, String> params = new LinkedHashMap<>();
        params.put("username", username);
        oldPassword = EncoderUtil.Encrypt(oldPassword, "MD5");
        newPassword = EncoderUtil.Encrypt(newPassword, "MD5");
        params.put("oldpassword", oldPassword);
        params.put("newpassword", newPassword);
        params.put("transtime", getTransTime());
        params.put("sign", createSign(params, "SHA-1"));

        try {
            String response = postEngine.post(url, params);
            ServerPacket serverPacket = ServerPacket.getServerPacketFrom(response);
            if (serverPacket.getState().equals(QUICK_PAY_SUCCESS)) {
                return;
            } else {
                throw new QuickPayException(serverPacket.getError());
            }
        } catch (IOException e) {
            throw new QuickPayException();
        }
    }

    /**
     * errors:
     * username_password_error
     *
     * @param username
     * @param password
     */
    @Override
    public void activate(String username, String password) {
        String url = quickPayConfigStorage.getUrl() + "/request_activate";

        Map<String, String> params = new LinkedHashMap<>();
        params.put("username", username);
        password = EncoderUtil.Encrypt(password, "MD5");
        params.put("password", password);
        params.put("transtime", getTransTime());
        params.put("sign", createSign(params, "SHA-1"));

        try {
            String response = postEngine.post(url, params);
            ServerPacket serverPacket = ServerPacket.getServerPacketFrom(response);
            if (serverPacket.getState().equals(QUICK_PAY_SUCCESS)) {
                return;
            } else {
                throw new QuickPayException(serverPacket.getError());
            }
        } catch (IOException e) {
            throw new QuickPayException();
        }
    }

    /**
     * errors:
     * user_already_improved
     *
     * @param username
     * @param password
     * @param province
     * @param city
     * @param bank_open
     * @param branch_bank
     * @param bankNo
     * @param payee
     * @param payee_card
     * @param phone_num
     */
    @Override
    public void updateInfo(String username, String password, String province, String city, String bank_open, String branch_bank, String bankNo, String payee, String payee_card, String phone_num) {
        String url = quickPayConfigStorage.getUrl() + "/improveinfo";

        Map<String, String> params = new LinkedHashMap<>();

        params.put("username", username);
        password = EncoderUtil.Encrypt(password, "MD5");
        params.put("password", password);

        params.put("province", province);
        params.put("city", city);
        params.put("bank_open", bank_open);
        params.put("branch_bank", branch_bank);
        params.put("bankNo", bankNo);
        params.put("payee", payee);
        params.put("payee_card", payee_card);
        params.put("phone_num", phone_num);
        params.put("transtime", getTransTime());
        params.put("sign", createSign(params, "SHA-1"));


        try {
            String response = postEngine.post(url, params);
            ServerPacket serverPacket = ServerPacket.getServerPacketFrom(response);
            if (serverPacket.getState().equals(QUICK_PAY_SUCCESS)) {
                return;
            } else {
                throw new QuickPayException(serverPacket.getError());
            }
        } catch (IOException e) {
            throw new QuickPayException();
        }
    }

    @Override
    public void increaseLimit(String username, String password, String payee, String phone_num, String email) {
        String url = quickPayConfigStorage.getUrl() + "/limitincrease";

        Map<String, String> params = new LinkedHashMap<>();
        params.put("username", username);
        password = EncoderUtil.Encrypt(password, "MD5");
        params.put("password", password);
        params.put("payee", payee);
        params.put("email", email);
        params.put("phone_num", phone_num);
        params.put("transtime", getTransTime());
        params.put("sign", createSign(params, "SHA-1"));

        try {
            String response = postEngine.post(url, params);
            ServerPacket serverPacket = ServerPacket.getServerPacketFrom(response);
            if (serverPacket.getState().equals(QUICK_PAY_SUCCESS)) {
                return;
            } else {
                throw new QuickPayException(serverPacket.getError());
            }
        } catch (IOException e) {
            throw new QuickPayException();
        }
    }

    @Override
    public BankInfo getBankInfo(String username, String password) {
        String url = quickPayConfigStorage.getUrl() + "/getinfo";

        Map<String, String> params = new LinkedHashMap<>();
        params.put("username", username);
        password = EncoderUtil.Encrypt(password, "MD5");
        params.put("password", password);
        params.put("transtime", getTransTime());
        params.put("sign", createSign(params, "SHA-1"));

        try {
            String response = postEngine.post(url, params);
            ServerPacket serverPacket = ServerPacket.getServerPacketFrom(response);
            if (serverPacket.getState().equals(QUICK_PAY_SUCCESS)) {
                return serverPacket.getInfo();
            } else {
                throw new QuickPayException(serverPacket.getError());
            }
        } catch (IOException e) {
            throw new QuickPayException();
        }
    }


    /**
     * server not support
     *
     * @param username
     */
    @Override
    public void forgetPassword(String username) {
        String url = quickPayConfigStorage.getUrl() + "/forgetpassword";

        Map<String, String> params = new LinkedHashMap<>();
        params.put("username", username);
        params.put("transtime", getTransTime());
        params.put("sign", createSign(params, "SHA-1"));

        try {
            String response = postEngine.post(url, params);
            ServerPacket serverPacket = ServerPacket.getServerPacketFrom(response);
            //TODO: issue, what if serverPacket has not state?
            if (serverPacket.getState().equals(QUICK_PAY_SUCCESS)) {
                return;
            } else {
                throw new QuickPayException(serverPacket.getError());
            }
        } catch (IOException e) {
            throw new QuickPayException();
        }
    }

    /**
     * Not tested, no one use
     *
     * @param username
     * @param code
     * @param newPassword
     */
    @Override
    public void resetPassword(String username, String code, String newPassword) {
        String url = quickPayConfigStorage.getUrl() + "/resetpassword";

        Map<String, String> params = new LinkedHashMap<>();
        params.put("username", username);
        newPassword = EncoderUtil.Encrypt(newPassword, "MD5");
        params.put("code", code);
        params.put("newpassword", newPassword);
        params.put("transtime", getTransTime());
        params.put("sign", createSign(params, "SHA-1"));

        try {
            String response = postEngine.post(url, params);
            ServerPacket serverPacket = ServerPacket.getServerPacketFrom(response);
            if (serverPacket.getState().equals(QUICK_PAY_SUCCESS)) {
                return;
            } else {
                throw new QuickPayException(serverPacket.getError());
            }
        } catch (IOException e) {
            throw new QuickPayException();
        }
    }

    @Override
    public ServerPacket getHistoryBills(String username, String password, String clientid, String month, long index, String status) {
        String url = quickPayConfigStorage.getUrl() + "/bill";

        Map<String, String> params = new LinkedHashMap<>();
        params.put("username", username);
        password = EncoderUtil.Encrypt(password, "MD5");
        params.put("password", password);
        params.put("clientid", clientid);
        params.put("month", month);
        params.put("index", "" + index);
        params.put("status", status);
        params.put("transtime", getTransTime());
        params.put("sign", createSign(params, "SHA-1"));

        try {
            String response = postEngine.post(url, params);
            ServerPacket serverPacket = ServerPacket.getServerPacketFrom(response);
            if (serverPacket.getState().equals(QUICK_PAY_SUCCESS)) {
                return serverPacket;
            } else {
                throw new QuickPayException(serverPacket.getError());
            }
        } catch (IOException e) {
            throw new QuickPayException();
        }
    }
}
