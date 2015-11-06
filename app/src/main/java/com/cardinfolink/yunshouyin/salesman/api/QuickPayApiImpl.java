package com.cardinfolink.yunshouyin.salesman.api;

import com.cardinfolink.yunshouyin.salesman.model.SAServerPacket;
import com.cardinfolink.yunshouyin.salesman.model.User;
import com.cardinfolink.yunshouyin.salesman.utils.EncoderUtil;

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

        if (this.quickPayConfigStorage.getProxy_url() != null && !"".equals(this.quickPayConfigStorage.getProxy_url())) {
            Proxy httpProxy = new Proxy(Proxy.Type.HTTP, new InetSocketAddress(this.quickPayConfigStorage.getProxy_url(), this.quickPayConfigStorage.getProxy_port()));
            postEngine = new PostEngine(httpProxy);
        } else {
            postEngine = new PostEngine();
        }
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
     * 1. prepare request url and post data
     * 2. get response from server
     * 3. throw QuickPayException if response state code is fail
     * 4. throw QuickPayException if network error
     *
     * @param username
     * @param password
     * @return
     */
    @Override
    public String login(String username, String password) {
        String url = quickPayConfigStorage.getUrl() + "/login";

        Date now = new Date();
        SimpleDateFormat spf = new SimpleDateFormat("yyyyMMddHHmmss");
        String transTime = spf.format(now);
        Map<String, String> params = new LinkedHashMap<>();
        params.put("username", username);
        params.put("password", password);
        params.put("transtime", transTime);
        params.put("sign", createSign(params, "SHA-1"));

        try {
            String response = postEngine.post(url, params);
            SAServerPacket serverPacket = SAServerPacket.getServerPacketFrom(response);
            if (serverPacket.getState().equals(QUICK_PAY_SUCCESS)) {
                // cache the accessToken
                quickPayConfigStorage.setAccessToken(serverPacket.getAccessToken());
                return quickPayConfigStorage.getAccessToken();
            } else {
                throw new QuickPayException(serverPacket.getError());
            }
        } catch (IOException e) {
            throw new QuickPayException();
        }
    }

    @Override
    public String getUploadToken() {
        checkAccessToken();

        String url = quickPayConfigStorage.getUrl() + "/uploadToken";
        Map<String, String> params = new LinkedHashMap<>();
        params.put("accessToken", quickPayConfigStorage.getAccessToken());

        try {
            String response = postEngine.post(url, params);
            SAServerPacket serverPacket = SAServerPacket.getServerPacketFrom(response);
            if (serverPacket.getState().equals(QUICK_PAY_SUCCESS)) {
                return serverPacket.getUploadToken();
            } else {
                throw new QuickPayException(serverPacket.getError());
            }
        } catch (IOException e) {
            throw new QuickPayException();
        }
    }

    @Override
    public User[] getUsers() {
        checkAccessToken();

        String url = quickPayConfigStorage.getUrl() + "/users";
        Map<String, String> params = new LinkedHashMap<>();
        params.put("accessToken", quickPayConfigStorage.getAccessToken());

        try {
            String response = postEngine.post(url, params);
            SAServerPacket serverPacket = SAServerPacket.getServerPacketFrom(response);
            if (serverPacket.getState().equals(QUICK_PAY_SUCCESS)) {
                return serverPacket.getUsers();
            } else {
                throw new QuickPayException(serverPacket.getError());
            }
        } catch (IOException e) {
            throw new QuickPayException();
        }
    }

    @Override
    public User registerUser(String username, String password) {
        checkAccessToken();

        String url = quickPayConfigStorage.getUrl() + "/register";
        Map<String, String> params = new LinkedHashMap<>();
        params.put("accessToken", quickPayConfigStorage.getAccessToken());
        params.put("username", username);
        password = EncoderUtil.Encrypt(password, "MD5");
        params.put("password", password);

        try {
            String response = postEngine.post(url, params);
            SAServerPacket serverPacket = SAServerPacket.getServerPacketFrom(response);
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
    public User updateUser(User user) {
        checkAccessToken();

        String url = quickPayConfigStorage.getUrl() + "/update";
        Map<String, String> params = new LinkedHashMap<>();
        params.put("accessToken", quickPayConfigStorage.getAccessToken());
        if (user.getUsername() != null && !user.getUsername().equals("")) {
            params.put("username", user.getUsername());
        }
        if (user.getProvince() != null && !user.getProvince().equals("")) {
            params.put("province", user.getProvince());
        }
        if (user.getCity() != null && !user.getCity().equals("")) {
            params.put("city", user.getCity());
        }
        if (user.getBank_open() != null && !user.getBank_open().equals("")) {
            params.put("bank_open", user.getBank_open());
        }
        if (user.getBranch_bank() != null && !user.getBranch_bank().equals("")) {
            params.put("branch_bank", user.getBranch_bank());
        }
        if (user.getBankNo() != null && !user.getBankNo().equals("")) {
            params.put("bankNo", user.getBankNo());
        }
        if (user.getPayee() != null && !user.getPayee().equals("")) {
            params.put("payee", user.getPayee());
        }
        if (user.getPayee_card() != null && !user.getPayee_card().equals("")) {
            params.put("payee_card", user.getPayee_card());
        }
        if (user.getPhone_num() != null && !user.getPhone_num().equals("")) {
            params.put("phone_num", user.getPhone_num());
        }
        if (user.getMerName() != null && !user.getMerName().equals("")) {
            params.put("merName", user.getMerName());
        }
        if (user.getImages() != null) {
            for (String uri : user.getImages()) {
                params.put("image", uri);
            }
        }

        try {
            String response = postEngine.post(url, params);
            SAServerPacket serverPacket = SAServerPacket.getServerPacketFrom(response);
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
    public User activateUser(String username) {
        checkAccessToken();

        String url = quickPayConfigStorage.getUrl() + "/activate";
        Map<String, String> params = new LinkedHashMap<>();
        params.put("accessToken", quickPayConfigStorage.getAccessToken());
        params.put("username", username);

        try {
            String response = postEngine.post(url, params);
            SAServerPacket serverPacket = SAServerPacket.getServerPacketFrom(response);
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
    public String getQrPostUrl(String merchantId, String imageType) {
        checkAccessToken();

        String url = quickPayConfigStorage.getUrl() + "/download";
        Map<String, String> params = new LinkedHashMap<>();
        params.put("accessToken", quickPayConfigStorage.getAccessToken());
        params.put("merId", merchantId);
        params.put("imageType", imageType);

        try {
            String response = postEngine.post(url, params);
            SAServerPacket serverPacket = SAServerPacket.getServerPacketFrom(response);
            if (serverPacket.getState().equals(QUICK_PAY_SUCCESS)) {
                return serverPacket.getDownloadUrl();
            } else {
                throw new QuickPayException(serverPacket.getError());
            }
        } catch (IOException e) {
            throw new QuickPayException();
        }
    }

    private void checkAccessToken() {
        if (quickPayConfigStorage.getAccessToken() == null || "".equals(quickPayConfigStorage.getAccessToken())) {
            throw new QuickPayException(QuickPayException.ACCESSTOKEN_NOT_FOUND);
        }
    }
}
