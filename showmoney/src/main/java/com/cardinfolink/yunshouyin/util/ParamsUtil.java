package com.cardinfolink.yunshouyin.util;

import android.annotation.SuppressLint;

import com.cardinfolink.yunshouyin.constant.SystemConfig;
import com.cardinfolink.yunshouyin.data.RequestParam;
import com.cardinfolink.yunshouyin.data.User;

import org.apache.http.NameValuePair;
import org.apache.http.message.BasicNameValuePair;

import java.text.SimpleDateFormat;
import java.util.Collections;
import java.util.Date;
import java.util.LinkedList;
import java.util.List;


public class ParamsUtil {


    @SuppressLint("SimpleDateFormat")
    public static RequestParam getLogin(String username, String password) {
        RequestParam requestParam = new RequestParam();
        //String url="http://211.147.72.70:10003/login";
        //String url="http://192.168.199.174:8081/login";
        String url = SystemConfig.Server + "/login";
        requestParam.setUrl(url);
        Date now = new Date();
        SimpleDateFormat spf = new SimpleDateFormat("yyyyMMddHHmmss");
        String transTime = spf.format(now);
        password = EncoderUtil.Encrypt(password, "MD5");
        List<NameValuePair> params = new LinkedList<NameValuePair>();
        params.add(new BasicNameValuePair("username", username));
        params.add(new BasicNameValuePair("password", password));
        params.add(new BasicNameValuePair("transtime", transTime));
        params.add(new BasicNameValuePair("sign", getSign(params, "SHA-1")));
        requestParam.setParams(params);
        return requestParam;
    }


    @SuppressLint("SimpleDateFormat")
    public static RequestParam getRequestActivate(String username, String password) {
        RequestParam requestParam = new RequestParam();
        //String url="http://211.147.72.70:10003/request_activate";
        // String url="http://192.168.199.174:8081/request_activate";
        String url = SystemConfig.Server + "/request_activate";
        requestParam.setUrl(url);
        Date now = new Date();
        SimpleDateFormat spf = new SimpleDateFormat("yyyyMMddHHmmss");
        String transTime = spf.format(now);
        password = EncoderUtil.Encrypt(password, "MD5");
        List<NameValuePair> params = new LinkedList<NameValuePair>();
        params.add(new BasicNameValuePair("username", username));
        params.add(new BasicNameValuePair("password", password));
        params.add(new BasicNameValuePair("transtime", transTime));
        params.add(new BasicNameValuePair("sign", getSign(params, "SHA-1")));
        requestParam.setParams(params);
        return requestParam;
    }


    public static RequestParam getRegister(String username, String password) {
        RequestParam requestParam = new RequestParam();
        //String url="http://211.147.72.70:10003/register";
        //String url="http://192.168.199.174:8081/register";
        String url = SystemConfig.Server + "/register";
        requestParam.setUrl(url);
        Date now = new Date();
        SimpleDateFormat spf = new SimpleDateFormat("yyyyMMddHHmmss");
        String transTime = spf.format(now);
        password = EncoderUtil.Encrypt(password, "MD5");
        List<NameValuePair> params = new LinkedList<NameValuePair>();
        params.add(new BasicNameValuePair("username", username));
        params.add(new BasicNameValuePair("password", password));
        params.add(new BasicNameValuePair("transtime", transTime));
        params.add(new BasicNameValuePair("sign", getSign(params, "SHA-1")));
        requestParam.setParams(params);
        return requestParam;

    }

    public static RequestParam getImproveInfo(User user) {
        RequestParam requestParam = new RequestParam();
        //String url="http://211.147.72.70:10003/improveinfo";
        // String url="http://192.168.199.174:8081/improveinfo";
        String url = SystemConfig.Server + "/improveinfo";
        requestParam.setUrl(url);
        Date now = new Date();
        SimpleDateFormat spf = new SimpleDateFormat("yyyyMMddHHmmss");
        @SuppressWarnings("unused")
        String transTime = spf.format(now);
        String password = EncoderUtil.Encrypt(user.getPassword(), "MD5");
        List<NameValuePair> params = new LinkedList<NameValuePair>();
        params.add(new BasicNameValuePair("username", user.getUsername()));
        params.add(new BasicNameValuePair("password", password));
        params.add(new BasicNameValuePair("province", user.getProvince()));
        params.add(new BasicNameValuePair("city", user.getCity()));
        params.add(new BasicNameValuePair("bank_open", user.getBankOpen()));
        params.add(new BasicNameValuePair("branch_bank", user.getBranchBank()));
        params.add(new BasicNameValuePair("bankNo", user.getBankNo()));
        params.add(new BasicNameValuePair("payee", user.getPayee()));
        params.add(new BasicNameValuePair("payee_card", user.getPayeeCard()));
        params.add(new BasicNameValuePair("phone_num", user.getPhoneNum()));
        params.add(new BasicNameValuePair("transtime", transTime));
        params.add(new BasicNameValuePair("sign", getSign(params, "SHA-1")));
        requestParam.setParams(params);
        return requestParam;

    }

    public static RequestParam getHistory(User user, String month, long index, String status) {
        RequestParam requestParam = new RequestParam();
        //String url="http://211.147.72.70:10003/bill";
        //String url="http://192.168.199.174:8081/bill";
        String url = SystemConfig.Server + "/bill";
        requestParam.setUrl(url);
        Date now = new Date();
        SimpleDateFormat spf = new SimpleDateFormat("yyyyMMddHHmmss");
        @SuppressWarnings("unused")
        String transTime = spf.format(now);
        String password = EncoderUtil.Encrypt(user.getPassword(), "MD5");
        List<NameValuePair> params = new LinkedList<NameValuePair>();
        params.add(new BasicNameValuePair("username", user.getUsername()));
        params.add(new BasicNameValuePair("password", password));
        params.add(new BasicNameValuePair("clientid", user.getClientid()));
        params.add(new BasicNameValuePair("month", month));
        params.add(new BasicNameValuePair("index", "" + index));
        params.add(new BasicNameValuePair("status", "" + status));
        params.add(new BasicNameValuePair("transtime", transTime));
        params.add(new BasicNameValuePair("sign", getSign(params, "SHA-1")));
        requestParam.setParams(params);
        return requestParam;

    }


    public static RequestParam getTotal(User user, String date) {
        RequestParam requestParam = new RequestParam();
        //String url="http://211.147.72.70:10003/getTotal";
        //String url="http://192.168.199.174:8081/getTotal";
        String url = SystemConfig.Server + "/getTotal";
        requestParam.setUrl(url);
        Date now = new Date();
        SimpleDateFormat spf = new SimpleDateFormat("yyyyMMddHHmmss");
        @SuppressWarnings("unused")
        String transTime = spf.format(now);
        String password = EncoderUtil.Encrypt(user.getPassword(), "MD5");
        List<NameValuePair> params = new LinkedList<NameValuePair>();
        params.add(new BasicNameValuePair("username", user.getUsername()));
        params.add(new BasicNameValuePair("password", password));
        params.add(new BasicNameValuePair("clientid", user.getClientid()));
        params.add(new BasicNameValuePair("date", date));
        params.add(new BasicNameValuePair("transtime", transTime));
        params.add(new BasicNameValuePair("sign", getSign(params, "SHA-1")));
        requestParam.setParams(params);
        return requestParam;

    }


    public static RequestParam getOrder(User user, String orderNum) {
        RequestParam requestParam = new RequestParam();
        //String url="http://211.147.72.70:10003/getrefd";
        //String url="http://192.168.199.174:8081/getrefd";
        String url = SystemConfig.Server + "/getOrder";
        requestParam.setUrl(url);
        Date now = new Date();
        SimpleDateFormat spf = new SimpleDateFormat("yyyyMMddHHmmss");
        @SuppressWarnings("unused")
        String transTime = spf.format(now);
        String password = EncoderUtil.Encrypt(user.getPassword(), "MD5");
        List<NameValuePair> params = new LinkedList<NameValuePair>();
        params.add(new BasicNameValuePair("username", user.getUsername()));
        params.add(new BasicNameValuePair("password", password));
        params.add(new BasicNameValuePair("clientid", user.getClientid()));
        params.add(new BasicNameValuePair("orderNum", orderNum));
        params.add(new BasicNameValuePair("transtime", transTime));
        params.add(new BasicNameValuePair("sign", getSign(params, "SHA-1")));
        requestParam.setParams(params);
        return requestParam;

    }


    public static RequestParam getRefd(User user, String orderNum) {
        RequestParam requestParam = new RequestParam();
        //String url="http://211.147.72.70:10003/getrefd";
        //String url="http://192.168.199.174:8081/getrefd";
        String url = SystemConfig.Server + "/getrefd";
        requestParam.setUrl(url);
        Date now = new Date();
        SimpleDateFormat spf = new SimpleDateFormat("yyyyMMddHHmmss");
        @SuppressWarnings("unused")
        String transTime = spf.format(now);
        String password = EncoderUtil.Encrypt(user.getPassword(), "MD5");
        List<NameValuePair> params = new LinkedList<NameValuePair>();
        params.add(new BasicNameValuePair("username", user.getUsername()));
        params.add(new BasicNameValuePair("password", password));
        params.add(new BasicNameValuePair("clientid", user.getClientid()));
        params.add(new BasicNameValuePair("orderNum", orderNum));
        params.add(new BasicNameValuePair("transtime", transTime));
        params.add(new BasicNameValuePair("sign", getSign(params, "SHA-1")));
        requestParam.setParams(params);
        return requestParam;

    }


    public static RequestParam getInfo(User user) {
        RequestParam requestParam = new RequestParam();
        //String url="http://211.147.72.70:10003/getinfo";
        //String url="http://192.168.199.174:8081/getinfo";
        String url = SystemConfig.Server + "/getinfo";
        requestParam.setUrl(url);
        Date now = new Date();
        SimpleDateFormat spf = new SimpleDateFormat("yyyyMMddHHmmss");
        @SuppressWarnings("unused")
        String transTime = spf.format(now);
        String password = EncoderUtil.Encrypt(user.getPassword(), "MD5");
        List<NameValuePair> params = new LinkedList<NameValuePair>();
        params.add(new BasicNameValuePair("username", user.getUsername()));
        params.add(new BasicNameValuePair("password", password));
        params.add(new BasicNameValuePair("transtime", transTime));
        params.add(new BasicNameValuePair("sign", getSign(params, "SHA-1")));
        requestParam.setParams(params);
        return requestParam;

    }


    public static RequestParam forgetPassword(String username) {
        RequestParam requestParam = new RequestParam();
        //String url="http://211.147.72.70:10003/forgetpassword";
        //String url="http://192.168.199.174:8081/forgetpassword";
        String url = SystemConfig.Server + "/forgetpassword";
        requestParam.setUrl(url);
        Date now = new Date();
        SimpleDateFormat spf = new SimpleDateFormat("yyyyMMddHHmmss");
        @SuppressWarnings("unused")
        String transTime = spf.format(now);
        List<NameValuePair> params = new LinkedList<NameValuePair>();
        params.add(new BasicNameValuePair("username", username));
        params.add(new BasicNameValuePair("transtime", transTime));
        params.add(new BasicNameValuePair("sign", getSign(params, "SHA-1")));
        requestParam.setParams(params);
        return requestParam;

    }


    @SuppressLint("SimpleDateFormat")
    public static RequestParam getUpdate(String username, String oldpassword, String newpassword) {
        RequestParam requestParam = new RequestParam();
        //String url="http://211.147.72.70:10003/updatepassword";
        //String url="http://192.168.199.174:8081/updatepassword";
        String url = SystemConfig.Server + "/updatepassword";
        requestParam.setUrl(url);
        Date now = new Date();
        SimpleDateFormat spf = new SimpleDateFormat("yyyyMMddHHmmss");
        @SuppressWarnings("unused")
        String transTime = spf.format(now);
        oldpassword = EncoderUtil.Encrypt(oldpassword, "MD5");
        newpassword = EncoderUtil.Encrypt(newpassword, "MD5");
        List<NameValuePair> params = new LinkedList<NameValuePair>();
        params.add(new BasicNameValuePair("username", username));
        params.add(new BasicNameValuePair("oldpassword", oldpassword));
        params.add(new BasicNameValuePair("newpassword", newpassword));
        params.add(new BasicNameValuePair("transtime", transTime));
        params.add(new BasicNameValuePair("sign", getSign(params, "SHA-1")));
        requestParam.setParams(params);
        return requestParam;
    }


    @SuppressLint("SimpleDateFormat")
    public static RequestParam getReset(String username, String code, String newpassword) {
        RequestParam requestParam = new RequestParam();
        //String url="http://211.147.72.70:10003/resetpassword";
        //String url="http://192.168.199.174:8081/resetpassword";
        String url = SystemConfig.Server + "/resetpassword";
        requestParam.setUrl(url);
        Date now = new Date();
        SimpleDateFormat spf = new SimpleDateFormat("yyyyMMddHHmmss");
        @SuppressWarnings("unused")
        String transTime = spf.format(now);
        newpassword = EncoderUtil.Encrypt(newpassword, "MD5");
        List<NameValuePair> params = new LinkedList<NameValuePair>();
        params.add(new BasicNameValuePair("username", username));
        params.add(new BasicNameValuePair("code", code));
        params.add(new BasicNameValuePair("newpassword", newpassword));
        params.add(new BasicNameValuePair("transtime", transTime));
        params.add(new BasicNameValuePair("sign", getSign(params, "SHA-1")));
        requestParam.setParams(params);
        return requestParam;
    }


    public static RequestParam getUpdateInfo(User user) {
        RequestParam requestParam = new RequestParam();
        //String url="http://211.147.72.70:10003/updateinfo";
        // String url="http://192.168.199.174:8081/updateinfo";
        String url = SystemConfig.Server + "/updateinfo";
        requestParam.setUrl(url);
        Date now = new Date();
        SimpleDateFormat spf = new SimpleDateFormat("yyyyMMddHHmmss");
        @SuppressWarnings("unused")
        String transTime = spf.format(now);
        String password = EncoderUtil.Encrypt(user.getPassword(), "MD5");
        List<NameValuePair> params = new LinkedList<NameValuePair>();
        params.add(new BasicNameValuePair("username", user.getUsername()));
        params.add(new BasicNameValuePair("password", password));
        params.add(new BasicNameValuePair("province", user.getProvince()));
        params.add(new BasicNameValuePair("city", user.getCity()));
        params.add(new BasicNameValuePair("bank_open", user.getBankOpen()));
        params.add(new BasicNameValuePair("branch_bank", user.getBranchBank()));
        params.add(new BasicNameValuePair("bankNo", user.getBankNo()));
        params.add(new BasicNameValuePair("payee", user.getPayee()));
        params.add(new BasicNameValuePair("payee_card", user.getPayeeCard()));
        params.add(new BasicNameValuePair("phone_num", user.getPhoneNum()));
        params.add(new BasicNameValuePair("transtime", transTime));
        params.add(new BasicNameValuePair("sign", getSign(params, "SHA-1")));
        requestParam.setParams(params);
        return requestParam;

    }

    public static RequestParam getLimitincrease(User user) {
        RequestParam requestParam = new RequestParam();
        //String url="http://211.147.72.70:10003/limitincrease";
        // String url="http://192.168.199.174:8081/limitincrease";
        String url = SystemConfig.Server + "/limitincrease";
        requestParam.setUrl(url);
        Date now = new Date();
        SimpleDateFormat spf = new SimpleDateFormat("yyyyMMddHHmmss");
        @SuppressWarnings("unused")
        String transTime = spf.format(now);
        String password = EncoderUtil.Encrypt(user.getPassword(), "MD5");
        List<NameValuePair> params = new LinkedList<NameValuePair>();
        params.add(new BasicNameValuePair("username", user.getUsername()));
        params.add(new BasicNameValuePair("password", password));
        params.add(new BasicNameValuePair("payee", user.getLimitName()));
        params.add(new BasicNameValuePair("email", user.getLimitEmail()));
        params.add(new BasicNameValuePair("phone_num", user.getLimitPhone()));
        params.add(new BasicNameValuePair("transtime", transTime));
        params.add(new BasicNameValuePair("sign", getSign(params, "SHA-1")));
        requestParam.setParams(params);
        return requestParam;

    }


    @SuppressWarnings("unchecked")
    public static String getSign(List<NameValuePair> params, String signType) {
        ComparatorNameValuePair comparatorNameValuePair = new ComparatorNameValuePair();
        Collections.sort(params, comparatorNameValuePair);
        StringBuilder sb = new StringBuilder();
        for (int i = 0; i < params.size() - 1; i++) {
            sb.append(params.get(i).getName());
            sb.append('=');
            sb.append(params.get(i).getValue());
            sb.append('&');
        }
        sb.append(params.get(params.size() - 1).getName());
        sb.append('=');
        sb.append(params.get(params.size() - 1).getValue());
        sb.append(SystemConfig.APP_KEY);
        String sign = EncoderUtil.Encrypt(sb.toString(), signType);
        return sign;
    }
}
