package com.cardinfolink.yunshouyin.api;


import android.support.annotation.NonNull;
import android.text.TextUtils;

import com.cardinfolink.yunshouyin.data.User;
import com.cardinfolink.yunshouyin.model.BankInfo;
import com.cardinfolink.yunshouyin.model.ServerPacket;
import com.cardinfolink.yunshouyin.model.ServerPacketOrder;
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
        // {"state":"fail","error":"old_password_error","count":0,"size":0,"refdcount":0}

        String url = quickPayConfigStorage.getUrl() + "/updatepassword"; //更新密码

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
     */
    @Override
    public User improveInfo(User user) {
        String url = quickPayConfigStorage.getUrl() + "/improveinfo";

        Map<String, String> params = new LinkedHashMap<>();

        params.put("username", user.getUsername());

        String password = EncoderUtil.Encrypt(user.getPassword(), "MD5");
        params.put("password", password);

        params.put("province", user.getProvince());
        params.put("city", user.getCity());
        params.put("bank_open", user.getBankOpen());
        params.put("branch_bank", user.getBranchBank());
        params.put("bankNo", user.getBankNo());
        params.put("payee", user.getPayee());
        params.put("payee_card", user.getPayeeCard());
        params.put("phone_num", user.getPhoneNum());
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
    public User updateInfo(User user) {
        //{"state":"fail","error":"请联系您的服务商为您修改清算信息。","count":0,"size":0,"refdcount":0}

        String url = quickPayConfigStorage.getUrl() + "/updateinfo";

        Map<String, String> params = new LinkedHashMap<>();

        params.put("username", user.getUsername());

        String password = EncoderUtil.Encrypt(user.getPassword(), "MD5");
        params.put("password", password);

        params.put("province", user.getProvince());
        params.put("city", user.getCity());
        params.put("bank_open", user.getBankOpen());
        params.put("branch_bank", user.getBranchBank());
        params.put("bankNo", user.getBankNo());
        params.put("payee", user.getPayee());
        params.put("payee_card", user.getPayeeCard());
        params.put("phone_num", user.getPhoneNum());
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

    /**
     * 提升限额
     */
    @Override
    public void increaseLimit(User user) {
        String url = quickPayConfigStorage.getUrl() + "/limitincrease";

        Map<String, String> params = new LinkedHashMap<>();
        params.put("username", user.getUsername());
        String password = EncoderUtil.Encrypt(user.getPassword(), "MD5");
        params.put("password", password);

        params.put("payee", user.getLimitName());
        params.put("email", user.getLimitEmail());
        params.put("phone_num", user.getLimitPhone());

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
    public BankInfo getBankInfo(User user) {
        String url = quickPayConfigStorage.getUrl() + "/getinfo";

        Map<String, String> params = new LinkedHashMap<>();
        params.put("username", user.getUsername());
        String password = EncoderUtil.Encrypt(user.getPassword(), "MD5");
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

    @Override
    public String getTotal(User user, String date) {
        //result =={"state":"success","total":"0.00","count":5,"size":0,"refdcount":0}

        String url = quickPayConfigStorage.getUrl() + "/getTotal";
        Map<String, String> params = new LinkedHashMap<>();
        params.put("username", user.getUsername());
        String password = EncoderUtil.Encrypt(user.getPassword(), "MD5");
        params.put("password", password);
        params.put("clientid", user.getClientid());
        params.put("date", date);
        params.put("transtime", getTransTime());
        params.put("sign", createSign(params, "SHA-1"));

        try {
            String response = postEngine.post(url, params);

            ServerPacket serverPacket = ServerPacket.getServerPacketFrom(response);
            if (serverPacket.getState().equals(QUICK_PAY_SUCCESS)) {
                return serverPacket.getTotal();
            } else {
                throw new QuickPayException(serverPacket.getError());
            }
        } catch (IOException e) {
            throw new QuickPayException();
        }
    }

    @Override
    public ServerPacketOrder getOrder(User user, String orderNum) {
        /**
         * {
         "state": "success",
         "count": 0,
         "size": 0,
         "refdcount": 0,
         "txn": {
         "response": "09",
         "system_date": "20151204112740",
         "transStatus": "10",
         "refundAmt": 0,
         "m_request": {
         "busicd": "PAUT",
         "inscd": "99911888",
         "txndir": "Q",
         "terminalid": "000000000000000",
         "orderNum": "15120322232663574",
         "mchntid": "999118880000017",
         "tradeFrom": "android",
         "txamt": "000000089500",
         "chcd": "ALP",
         "currency": "CNY"
         }
         }
         }
         */
        String url = quickPayConfigStorage.getUrl() + "/getOrder";

        Map<String, String> params = new LinkedHashMap<>();
        params.put("username", user.getUsername());
        String password = EncoderUtil.Encrypt(user.getPassword(), "MD5");
        params.put("password", password);
        params.put("clientid", user.getClientid());
        params.put("orderNum", orderNum);
        params.put("transtime", getTransTime());
        params.put("sign", createSign(params, "SHA-1"));

        try {
            String response = postEngine.post(url, params);
            //特别注意这里用的是一个新的ServerPacketOrder类来解析json。
            ServerPacketOrder serverPacket = ServerPacketOrder.getServerPacketOrder(response);
            if (serverPacket.getState().equals(QUICK_PAY_SUCCESS)) {
                return serverPacket;
            } else {
                throw new QuickPayException(serverPacket.getError());
            }
        } catch (IOException e) {
            throw new QuickPayException();
        }
    }

    @Override
    public ServerPacket getRefd(User user, String orderNum) {
        //{"state":"success","count":0,"size":0,"refdcount":0,"refdtotal":"9.66"}
        //退款
        String url = quickPayConfigStorage.getUrl() + "/getrefd";

        Map<String, String> params = new LinkedHashMap<>();
        params.put("username", user.getUsername());
        String password = EncoderUtil.Encrypt(user.getPassword(), "MD5");
        params.put("password", password);
        params.put("clientid", user.getClientid());
        params.put("orderNum", orderNum);
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
