package com.cardinfolink.yunshouyin.salesman.utils;


import com.cardinfolink.yunshouyin.salesman.models.SystemConfig;
import com.cardinfolink.yunshouyin.salesman.models.User;

import org.apache.http.NameValuePair;
import org.apache.http.message.BasicNameValuePair;

import java.text.SimpleDateFormat;
import java.util.Collections;
import java.util.Date;
import java.util.LinkedList;
import java.util.List;


public class ParamsUtil {

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

    public static RequestParam getLogin_SA(String username, String password) {
        RequestParam requestParam = new RequestParam();
        String url = SystemConfig.getServer_Tool() + "/login";
        requestParam.setUrl(url);
        Date now = new Date();
        SimpleDateFormat spf = new SimpleDateFormat("yyyyMMddHHmmss");
        String transTime = spf.format(now);
        List<NameValuePair> params = new LinkedList();
        params.add(new BasicNameValuePair("username", username));
        params.add(new BasicNameValuePair("password", password));
        params.add(new BasicNameValuePair("transtime", transTime));
        params.add(new BasicNameValuePair("sign", getSign(params, "SHA-1")));
        requestParam.setParams(params);
        return requestParam;
    }

    public static RequestParam getUsers_SA(String accessToken) {
        RequestParam requestParam = new RequestParam();
        String url = SystemConfig.getServer_Tool() + "/users";
        requestParam.setUrl(url);
        List<NameValuePair> params = new LinkedList();
        params.add(new BasicNameValuePair("accessToken", accessToken));
        requestParam.setParams(params);
        return requestParam;
    }

    /**
     * username_exist用户存在的错误
     *
     * @param accessToken
     * @param
     * @return
     */
    public static RequestParam getRegister_SA(String accessToken, String username, String password) {
        RequestParam requestParam = new RequestParam();
        String url = SystemConfig.getServer_Tool() + "/register";
        requestParam.setUrl(url);
        List<NameValuePair> params = new LinkedList();
        params.add(new BasicNameValuePair("accessToken", accessToken));
        params.add(new BasicNameValuePair("username", username));
        password = EncoderUtil.Encrypt(password, "MD5");
        params.add(new BasicNameValuePair("password", password));
        requestParam.setParams(params);
        return requestParam;
    }

    public static RequestParam getUploadToken_SA(String accessToken) {
        RequestParam requestParam = new RequestParam();
        String url = SystemConfig.getServer_Tool() + "/uploadToken";
        requestParam.setUrl(url);
        List<NameValuePair> params = new LinkedList();
        params.add(new BasicNameValuePair("accessToken", accessToken));
        requestParam.setParams(params);
        return requestParam;
    }

    /**
     * @param accessToken
     * @param user
     * @return
     */
    public static RequestParam getUpdate_SA(String accessToken, User user) {
        RequestParam requestParam = new RequestParam();
        String url = SystemConfig.getServer_Tool() + "/update";
        requestParam.setUrl(url);
        List<NameValuePair> params = new LinkedList<>();
        params.add(new BasicNameValuePair("accessToken", accessToken));

        // copy from 云收银
        if (user.getUsername() != null && !user.getUsername().equals("")) {
            params.add(new BasicNameValuePair("username", user.getUsername()));
        }
        if (user.getProvince() != null && !user.getProvince().equals("")) {
            params.add(new BasicNameValuePair("province", user.getProvince()));
        }
        if (user.getCity() != null && !user.getCity().equals("")) {
            params.add(new BasicNameValuePair("city", user.getCity()));
        }
        if (user.getBank_open() != null && !user.getBank_open().equals("")) {
            params.add(new BasicNameValuePair("bank_open", user.getBank_open()));
        }
        if (user.getBranch_bank() != null && !user.getBranch_bank().equals("")) {
            params.add(new BasicNameValuePair("branch_bank", user.getBranch_bank()));
        }
        if (user.getBankNo() != null && !user.getBankNo().equals("")) {
            params.add(new BasicNameValuePair("bankNo", user.getBankNo()));
        }
        if (user.getPayee() != null && !user.getPayee().equals("")) {
            params.add(new BasicNameValuePair("payee", user.getPayee()));
        }
        if (user.getPayee_card() != null && !user.getPayee_card().equals("")) {
            params.add(new BasicNameValuePair("payee_card", user.getPayee_card()));
        }
        if (user.getPhone_num() != null && !user.getPhone_num().equals("")) {
            params.add(new BasicNameValuePair("phone_num", user.getPhone_num()));
        }
        if (user.getMerName() != null && !user.getMerName().equals("")) {
            params.add(new BasicNameValuePair("merName", user.getMerName()));
        }
        if (user.getImages() != null) {
            for (String uri : user.getImages()) {
                params.add(new BasicNameValuePair("image", uri));
            }
        }
        requestParam.setParams(params);
        return requestParam;
    }

    /**
     * params_empty的错误,应该是注册信息不完备
     *
     * @param accessToken
     * @param username
     * @return
     */
    public static RequestParam getActivate_SA(String accessToken, String username, String qrCodeURI) {
        RequestParam requestParam = new RequestParam();
        String url = SystemConfig.getServer_Tool() + "/activate";
        requestParam.setUrl(url);
        List<NameValuePair> params = new LinkedList<>();
        params.add(new BasicNameValuePair("accessToken", accessToken));
        params.add(new BasicNameValuePair("username", username));
        //二维码生成暂时由服务端生成
//        params.add(new BasicNameValuePair("imageUrl", qrCodeURI));
        requestParam.setParams(params);
        return requestParam;
    }

    public static RequestParam getDownload(String accessToken, String merchantId, String imageType) {
        RequestParam requestParam = new RequestParam();
        String url = SystemConfig.getServer_Tool() + "/download";
        requestParam.setUrl(url);
        List<NameValuePair> params = new LinkedList<>();
        params.add(new BasicNameValuePair("accessToken", accessToken));
        params.add(new BasicNameValuePair("merId", merchantId));
        params.add(new BasicNameValuePair("imageType", imageType));
        requestParam.setParams(params);
        return requestParam;
    }
}
