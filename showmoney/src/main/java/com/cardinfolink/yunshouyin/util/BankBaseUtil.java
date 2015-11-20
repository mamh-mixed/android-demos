package com.cardinfolink.yunshouyin.util;

import android.util.Log;

import com.cardinfolink.yunshouyin.constant.SystemConfig;
import com.cardinfolink.yunshouyin.data.RequestParam;

import org.apache.commons.codec.binary.Hex;
import org.apache.http.NameValuePair;
import org.apache.http.message.BasicNameValuePair;

import java.util.Collections;
import java.util.LinkedList;
import java.util.List;

import javax.crypto.Mac;
import javax.crypto.spec.SecretKeySpec;

public class BankBaseUtil {
    private static final String TAG = "BankBaseUtil";

    public static RequestParam getProvince() {
        RequestParam requestParam = new RequestParam();
        String url = SystemConfig.bankbase_url + "/city/provinces/list.json";
        Log.i(TAG, "url = " + requestParam.getUrl());
        List<NameValuePair> params = new LinkedList<NameValuePair>();
        params.add(new BasicNameValuePair("appkey", SystemConfig.bankbase_key));
        params.add(new BasicNameValuePair("sig", getSign(params, SystemConfig.bankbase_key)));
        url = url + "?" + getValue(params);
        requestParam.setUrl(url);
        Log.i(TAG, "url = " + requestParam.getUrl());
        return requestParam;

    }


    public static RequestParam getCity(String province) {
        RequestParam requestParam = new RequestParam();
        String url = SystemConfig.bankbase_url + "/city/province/cities.json";
        Log.i(TAG, "url = " + requestParam.getUrl());
        List<NameValuePair> params = new LinkedList<NameValuePair>();
        params.add(new BasicNameValuePair("appkey", SystemConfig.bankbase_key));
        params.add(new BasicNameValuePair("province", province));
        params.add(new BasicNameValuePair("sig", getSign(params, SystemConfig.bankbase_key)));
        url = url + "?" + getValue(params);
        requestParam.setUrl(url);
        Log.i(TAG, "url = " + requestParam.getUrl());
        return requestParam;

    }


    public static RequestParam getBank() {
        RequestParam requestParam = new RequestParam();
        String url = SystemConfig.bankbase_url + "/bank/ids.json";
        Log.i(TAG, "url = " + requestParam.getUrl());
        List<NameValuePair> params = new LinkedList<NameValuePair>();
        params.add(new BasicNameValuePair("appkey", SystemConfig.bankbase_key));
        params.add(new BasicNameValuePair("sig", getSign(params, SystemConfig.bankbase_key)));
        url = url + "?" + getValue(params);
        requestParam.setUrl(url);
        Log.i(TAG, "url = " + requestParam.getUrl());
        return requestParam;

    }


    public static RequestParam getSerach(String city_code, String bank_id) {
        RequestParam requestParam = new RequestParam();
        String url = SystemConfig.bankbase_url + "/bank/search.json";
        Log.i(TAG, "bank_id=" + bank_id);
        List<NameValuePair> params = new LinkedList<NameValuePair>();
        params.add(new BasicNameValuePair("appkey", SystemConfig.bankbase_key));
        if (city_code != null && !city_code.isEmpty()) {
            params.add(new BasicNameValuePair("city_code", city_code));
        }

        if (bank_id != null && !bank_id.isEmpty()) {
            params.add(new BasicNameValuePair("bank_id", bank_id));
        }
        params.add(new BasicNameValuePair("size", "500"));
        params.add(new BasicNameValuePair("sig", getSign(params, SystemConfig.bankbase_key)));
        url = url + "?" + getValue(params);
        requestParam.setUrl(url);
        Log.i(TAG, "url = " + requestParam.getUrl());
        return requestParam;

    }

    public static String hmacSha1(String value, String key) {
        try {
            // Get an hmac_sha1 key from the raw key bytes
            byte[] keyBytes = key.getBytes();
            SecretKeySpec signingKey = new SecretKeySpec(keyBytes, "HmacSHA1");

            // Get an hmac_sha1 Mac instance and initialize with the signing key
            Mac mac = Mac.getInstance("HmacSHA1");
            mac.init(signingKey);

            // Compute the hmac on input data bytes
            byte[] rawHmac = mac.doFinal(value.getBytes());

            // Convert raw bytes to Hex
            byte[] hexBytes = new Hex().encode(rawHmac);

            // Covert array of Hex bytes to a String
            return new String(hexBytes, "UTF-8");
        } catch (Exception e) {

            throw new RuntimeException(e);
        }
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
        if (params.size() != 0) {
            sb.append(params.get(params.size() - 1).getName());
            sb.append('=');
            sb.append(params.get(params.size() - 1).getValue());
        } else {

        }
        String sign = "";
        Log.i(TAG, sb.toString());

        sign = hmacSha1(sb.toString(), SystemConfig.bankbase_key);
        Log.i(TAG, "sign=" + sign);
        return sign;
    }

    @SuppressWarnings("unchecked")
    public static String getValue(List<NameValuePair> params) {
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


        return sb.toString();
    }

}
