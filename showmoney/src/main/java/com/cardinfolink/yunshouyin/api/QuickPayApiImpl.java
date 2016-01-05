package com.cardinfolink.yunshouyin.api;


import android.text.TextUtils;

import com.cardinfolink.yunshouyin.data.User;
import com.cardinfolink.yunshouyin.model.BankInfo;
import com.cardinfolink.yunshouyin.model.Message;
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
    private static final String TAG = "QuickPayApiImpl";

    private static final String QUICK_PAY_SUCCESS = "success";

    private static final String DEVICE_TYPE = "Android";//这里就是android

    private static final String SIGN_TYPE_MD5 = "MD5";//密码加密使用这个

    private static final String SIGN_TYPE_SHA_1 = "SHA-1";
    private static final String SIGN_TYPE_SHA_256 = "SHA-256";
    private static final String SIGN_TYPE = SIGN_TYPE_SHA_256;//报文加密使用这个

    private static final String URL_PATH_REGISTER = "/v3/register";
    private static final String URL_PATH_LOGIN = "/v3/login";
    private static final String URL_PATH_PASSWORD_FORGET = "/v3/password/forget";
    private static final String URL_PATH_PASSWORD_UPDATE = "/v3/password/update";
    private static final String URL_PATH_ACCOUNT_ACTIVATE = "/v3/account/activate";
    private static final String URL_PATH_ACCOUNT_IMPROVE = "/v3/account/improve";
    private static final String URL_PATH_ACCOUNT_CERTIFICATE = "/v3/account/certificate";
    private static final String URL_PATH_ACCOUNT_INFO_SETTLE = "/v3/account/info/settle";
    private static final String URL_PATH_TOKEN_QINIU = "/v3/token/qiniu";
    private static final String URL_PATH_BILLS = "/v3/bills";
    private static final String URL_PATH_ORDERS = "/v3/orders";
    private static final String URL_PATH_COUPONS = "/v3/coupons";
    private static final String URL_PATH_SUMMARY_DAY = "/v3/summary/day";
    private static final String URL_PATH_MESSAGE_PULL = "/v3/message/pull";
    private static final String URL_PATH_MESSAGE_UPDATE = "/v3/message/update";

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

    private String getTransTime() {
        Date now = new Date();
        SimpleDateFormat spf = new SimpleDateFormat("yyyyMMddHHmmss");
        return spf.format(now);
    }

    private String createSign(Map<String, String> params) {
        return createSign(params, SIGN_TYPE);
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
        String sign = EncoderUtil.Encrypt(toSign.toString(), signType);
        return sign;
    }


    /**
     * 注册
     * /register:
     *
     * @param username
     * @param password
     * @param invite
     */
    @Override
    public void register(String username, String password, String invite) {
        String url = quickPayConfigStorage.getUrl() + URL_PATH_REGISTER;

        Map<String, String> params = new LinkedHashMap<>();
        params.put("username", username);
        password = EncoderUtil.Encrypt(password, SIGN_TYPE_MD5);
        params.put("password", password);
        if (!TextUtils.isEmpty(invite)) {
            params.put("invitationCode", invite);//邀请码,不是必须的
        }
        params.put("transtime", getTransTime());
        params.put("sign", createSign(params));

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
     * /login:登录
     *
     * @param username
     * @param password
     * @param deviceToken
     * @return
     */
    @Override
    public User login(String username, String password, String deviceToken) {
        String url = quickPayConfigStorage.getUrl() + URL_PATH_LOGIN;

        Map<String, String> params = new LinkedHashMap<>();
        params.put("username", username);
        password = EncoderUtil.Encrypt(password, SIGN_TYPE_MD5);
        params.put("password", password);
        params.put("transtime", getTransTime());
        params.put("deviceType", DEVICE_TYPE);//可选项：iOS, Android 建议：客户端上传大小写可能有变，建议后端统一成大写或是小写。如果传送了数据，需要更新数据库对应字段。
        if (!TextUtils.isEmpty(deviceToken)) {
            params.put("deviceToken", deviceToken);//消息推送的唯一标识。如果传送了数据，需要更新数据库对应字段。
        }
        params.put("sign", createSign(params));

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
     * 修改密码
     *
     * @param username
     * @param oldPassword
     * @param newPassword
     */
    @Override
    public void updatePassword(String username, String oldPassword, String newPassword) {
        // {"state":"fail","error":"old_password_error","count":0,"size":0,"refdcount":0}

        String url = quickPayConfigStorage.getUrl() + URL_PATH_PASSWORD_UPDATE; //更新密码

        Map<String, String> params = new LinkedHashMap<>();
        params.put("username", username);
        oldPassword = EncoderUtil.Encrypt(oldPassword, SIGN_TYPE_MD5);
        newPassword = EncoderUtil.Encrypt(newPassword, SIGN_TYPE_MD5);
        params.put("oldpassword", oldPassword);
        params.put("newpassword", newPassword);
        params.put("transtime", getTransTime());
        params.put("sign", createSign(params));

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
     * 账户激活
     *
     * @param username
     * @param password
     */
    @Override
    public void activate(String username, String password) {
        String url = quickPayConfigStorage.getUrl() + URL_PATH_ACCOUNT_ACTIVATE;

        Map<String, String> params = new LinkedHashMap<>();
        params.put("username", username);
        password = EncoderUtil.Encrypt(password, SIGN_TYPE_MD5);
        params.put("password", password);
        params.put("transtime", getTransTime());
        params.put("sign", createSign(params));

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
     * 清算银行卡信息完善
     *
     * @param user
     */
    @Override
    public User improveInfo(User user) {
        String url = quickPayConfigStorage.getUrl() + URL_PATH_ACCOUNT_IMPROVE;

        Map<String, String> params = new LinkedHashMap<>();

        params.put("username", user.getUsername());
        String password = EncoderUtil.Encrypt(user.getPassword(), SIGN_TYPE_MD5);
        params.put("password", password);

        params.put("province", user.getProvince());
        params.put("city", user.getCity());
        params.put("bankOpen", user.getBankOpen());
        params.put("branchBank", user.getBranchBank());
        params.put("bankNo", user.getBankNo());
        params.put("payee", user.getPayee());
        params.put("payeeCard", user.getPayeeCard());
        params.put("phoneNum", user.getPhoneNum());

        params.put("transtime", getTransTime());
        params.put("sign", createSign(params));

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
    public BankInfo getBankInfo(User user) {
        String url = quickPayConfigStorage.getUrl() + URL_PATH_ACCOUNT_INFO_SETTLE;

        Map<String, String> params = new LinkedHashMap<>();
        params.put("username", user.getUsername());
        String password = EncoderUtil.Encrypt(user.getPassword(), SIGN_TYPE_MD5);
        params.put("password", password);
        params.put("transtime", getTransTime());
        params.put("sign", createSign(params, SIGN_TYPE_SHA_1));

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
     * 忘记密码
     *
     * @param username
     */
    @Override
    public ServerPacket forgetPassword(String username) {
        String url = quickPayConfigStorage.getUrl() + URL_PATH_PASSWORD_FORGET;

        Map<String, String> params = new LinkedHashMap<>();
        params.put("username", username);
        params.put("transtime", getTransTime());
        params.put("sign", createSign(params));

        try {
            String response = postEngine.post(url, params);
            ServerPacket serverPacket = ServerPacket.getServerPacketFrom(response);
            //TODO: issue, what if serverPacket has not state?
            if (serverPacket.getState().equals(QUICK_PAY_SUCCESS)) {
                return serverPacket;
            } else {
                throw new QuickPayException(serverPacket.getError());
            }
        } catch (IOException e) {
            throw new QuickPayException();
        }
    }


    /**
     * 收款账单列表
     * 提供给APP获取收款账单的接口，只返回支付订单，没有退款订单
     *
     * @param user
     * @param month
     * @param index
     * @param size
     * @param status
     * @return
     */
    @Override
    public ServerPacket getHistoryBills(User user, String month, String index, String size, String status) {
        String url = quickPayConfigStorage.getUrl() + URL_PATH_BILLS;

        Map<String, String> params = new LinkedHashMap<>();
        params.put("username", user.getUsername());
        String password = EncoderUtil.Encrypt(user.getPassword(), SIGN_TYPE_MD5);
        params.put("password", password);

        params.put("clientId", user.getClientid());
        params.put("month", month);
        params.put("index", index);
        params.put("size", size);
        params.put("status", status);

        params.put("transtime", getTransTime());
        params.put("sign", createSign(params));

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

    //两个参数的用来精确查找某个订单的
    @Override
    public ServerPacket findOrder(User user, String orderNum) {
        return findOrder(user, orderNum, null, null, null, null, null);
    }

    //这是多个参数的用来根据条件查找的
    @Override
    public ServerPacket findOrder(User user, String index, String size, String recType, String payType, String txnStatus) {
        return findOrder(user, null, index, size, recType, payType, txnStatus);
    }


    /**
     * 查找订单
     * 查找订单，返回收款账单列表，本接口提供按照订单号搜索，以及状态位搜索，recType，payType，txnStatus。本接口包含了(/getOrde
     * 取单个订单的功能)
     *
     * @param user
     * @param orderNum
     * @param index
     * @param size
     * @param recType
     * @param payType
     * @param txnStatus
     * @return
     */
    @Override
    public ServerPacket findOrder(User user, String orderNum, String index, String size, String recType, String payType, String txnStatus) {
        String url = quickPayConfigStorage.getUrl() + URL_PATH_ORDERS;

        Map<String, String> params = new LinkedHashMap<>();
        params.put("username", user.getUsername());
        String password = EncoderUtil.Encrypt(user.getPassword(), SIGN_TYPE_MD5);
        params.put("password", password);

        if (TextUtils.isEmpty(orderNum)) {
            params.put("index", index);//分页起始位置
            params.put("size", size);//每页条数，默认15条
            params.put("recType", recType);//'收款方式：移动版(1) 桌面版(2) 收款码(4) 开放接口(8)。移动版｜桌面版：1 | 2 = 3移动版 | 收款码: 1 | 4 = 5'
            params.put("payType", payType);//'支付方式：支付宝 微信。支付宝 1，微信 2，全部：1 | 2 = 3'
            params.put("txnStatus", txnStatus);//'交易状态：交易成功 部分退款 全额退款。交易成功 1，部分退款 2，全额退款 4。部分退款 ｜ 全额退款：2 | 4 = 6，全部：1 | 2 | 4 = 7'
        } else {
            //如果提供了账单号，这里是精确查询，其他的参数不要传人
            params.put("orderNum", orderNum);//订单号，用于订单精准搜索
        }

        params.put("transtime", getTransTime());
        params.put("sign", createSign(params));

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

    /**
     * 这个有替代的方案   /summary/day:
     *
     * @param user
     * @param date
     * @return
     */
    @Override
    public String getTotal(User user, String date) {
        //result =={"state":"success","total":"0.00","count":5,"size":0,"refdcount":0}

        String url = quickPayConfigStorage.getUrl() + "/getTotal";
        Map<String, String> params = new LinkedHashMap<>();
        params.put("username", user.getUsername());
        String password = EncoderUtil.Encrypt(user.getPassword(), SIGN_TYPE_MD5);
        params.put("password", password);
        params.put("clientid", user.getClientid());
        params.put("date", date);
        params.put("transtime", getTransTime());
        params.put("sign", createSign(params, SIGN_TYPE_SHA_1));

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


    /**
     * 这个暂时还不确定 新的接口 v3里面要不要？？？？！！！！
     *
     * @param user
     * @param orderNum
     * @return
     */
    @Override
    public ServerPacket getRefd(User user, String orderNum) {
        //{"state":"success","count":0,"size":0,"refdcount":0,"refdtotal":"9.66"}
        //退款
        String url = quickPayConfigStorage.getUrl() + "/getrefd";

        Map<String, String> params = new LinkedHashMap<>();
        params.put("username", user.getUsername());
        String password = EncoderUtil.Encrypt(user.getPassword(), SIGN_TYPE_MD5);
        params.put("password", password);
        params.put("clientid", user.getClientid());
        params.put("orderNum", orderNum);
        params.put("transtime", getTransTime());
        params.put("sign", createSign(params, SIGN_TYPE_SHA_1));
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

    /**
     * 获取七牛上传token
     *
     * @param user
     * @return
     */
    @Override
    public String getUploadToken(User user) {
        /**
         * {
         "state": "success",
         "count": 0,
         "totalRecord": 0,
         "size": 0,
         "refdcount": 0,
         "uploadToken": "-OOrgfZJbxz29kiW6HQsJ_OQJcjX6gaPRDf6xOcc:8qhOkFfrq8whZ9QeekRCAh0gIPI=:eyJzY29wZSI6InRlc3QiLCJkZWFkbGluZSI6MTQ1MDc4NDExOSwiZW5kVXNlciI6InVzZXJJZCJ9"
         }
         */
        //送username，password
        String url = quickPayConfigStorage.getUrl() + URL_PATH_TOKEN_QINIU;

        Map<String, String> params = new LinkedHashMap<>();
        params.put("username", user.getUsername());
        String password = EncoderUtil.Encrypt(user.getPassword(), SIGN_TYPE_MD5);
        params.put("password", password);
        params.put("transtime", getTransTime());
        params.put("sign", createSign(params));
        try {
            String response = postEngine.post(url, params);
            ServerPacket serverPacket = ServerPacket.getServerPacketFrom(response);
            if (serverPacket.getState().equals(QUICK_PAY_SUCCESS)) {
                return serverPacket.getUploadToken();
            } else {
                throw new QuickPayException(serverPacket.getError());
            }
        } catch (IOException e) {
            throw new QuickPayException();
        }
    }

    /**
     * 账户认证
     *
     * @param user
     * @param imageMap
     */
    @Override
    public void improveCertInfo(User user, String certName, String certAddr, Map<String, String> imageMap) {
        /**
         * 送username，password
         * merName(店铺名称),
         * merAddr（店铺地址）,
         * legalCertPos（法人证书正面）,
         * legalCertOpp（法人证书反面）,
         * businessLicense（营业执照）,
         * taxRegistCert（税务登记证）,
         * organizeCodeCert（组织机构代码证）
         */
        String url = quickPayConfigStorage.getUrl() + URL_PATH_ACCOUNT_CERTIFICATE;
        Map<String, String> params = new LinkedHashMap<>();
        params.put("username", user.getUsername());
        String password = EncoderUtil.Encrypt(user.getPassword(), SIGN_TYPE_MD5);
        params.put("password", password);

        params.put("transtime", getTransTime());

        params.put("certName", certName);
        params.put("certAddr", certAddr);

        //把名字对应值都放入到params里面
        for (Map.Entry<String, String> map : imageMap.entrySet()) {
            params.put(map.getKey(), map.getValue());
        }
        params.put("sign", createSign(params));
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
     * 消息接口
     * 提供给APP拉取消息的接口
     *
     * @param username
     * @param password
     * @param size
     * @param lasttime
     * @param maxtime
     * @return
     */
    @Override
    public ServerPacket pullinfo(String username, String password, String size, String lasttime, String maxtime) {
        String url = quickPayConfigStorage.getUrl() + URL_PATH_MESSAGE_PULL;

        Map<String, String> params = new LinkedHashMap<>();
        params.put("username", username);
        password = EncoderUtil.Encrypt(password, SIGN_TYPE_MD5);
        params.put("password", password);
        if (!TextUtils.isEmpty(lasttime)) {
            params.put("lasttime", lasttime);
        }
        if (!TextUtils.isEmpty(maxtime)) {
            params.put("maxtime", maxtime);
        }
        params.put("size", size);
        params.put("transtime", getTransTime());
        params.put("sign", createSign(params));
        try {
            //// TODO: mamh  这里没有判断 serverPacket.getState()的状态？？？？
            String response = postEngine.post(url, params);
            ServerPacket serverPacket = ServerPacket.getServerPacketFrom(response);
            return serverPacket;
        } catch (IOException e) {
            throw new QuickPayException();
        }
    }

    /**
     * 消息更新
     * 提供给APP消息更新的接口,已读、已删
     *
     * @param username
     * @param password
     * @param status
     * @param messages
     * @return
     */
    @Override
    public ServerPacket updateMessage(String username, String password, String status, Message[] messages) {
        String url = quickPayConfigStorage.getUrl() + URL_PATH_MESSAGE_UPDATE;
        Map<String, String> params = new LinkedHashMap<>();
        params.put("username", username);
        password = EncoderUtil.Encrypt(password, SIGN_TYPE_MD5);
        params.put("password", password);
        StringBuilder sb = new StringBuilder("[");
        for (Message message : messages) {
            message.setStatus(status);
            sb.append("{").append("\"msgId\":").append("\"").append(message.getMsgId()).append("\"").append(",");
            sb.append("\"status\":").append(message.getStatus()).append("}").append(",");
        }
        sb.append("]");
        String messageStr = sb.toString().replace(",]", "]");
        params.put("message", messageStr);
        params.put("transtime", getTransTime());
        params.put("sign", createSign(params));
        try {
            String response = postEngine.post(url, params);
            ServerPacket serverPacket = ServerPacket.getServerPacketFrom(response);
            return serverPacket;
        } catch (IOException e) {
            throw new QuickPayException();
        }
    }
}
