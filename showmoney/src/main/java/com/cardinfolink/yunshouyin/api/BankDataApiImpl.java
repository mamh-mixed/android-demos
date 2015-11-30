package com.cardinfolink.yunshouyin.api;

import android.text.TextUtils;

import com.cardinfolink.yunshouyin.model.Bank;
import com.cardinfolink.yunshouyin.model.City;
import com.cardinfolink.yunshouyin.model.Province;
import com.cardinfolink.yunshouyin.model.SubBank;
import com.google.gson.Gson;
import com.google.gson.reflect.TypeToken;

import org.apache.commons.codec.binary.Hex;

import java.net.InetSocketAddress;
import java.net.Proxy;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collections;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import java.util.SortedMap;
import java.util.TreeMap;

import javax.crypto.Mac;
import javax.crypto.spec.SecretKeySpec;

/**
 * 要么取到数据,要么抛出网络异常的exception,需改进
 */
public class BankDataApiImpl implements BankDataApi {
    protected QuickPayConfigStorage quickPayConfigStorage;
    protected PostEngine postEngine;

    public BankDataApiImpl(QuickPayConfigStorage quickPayConfigStorage) {
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

    @Override
    public List<Province> getProvince() {
        String url = quickPayConfigStorage.getBankbaseUrl() + "/city/provinces/list.json";

        Map<String, String> params = new LinkedHashMap<>();
        params.put("appkey", quickPayConfigStorage.getBankbaseKey());
        params.put("sig", createSign(params));

        try {
            String response = postEngine.get(url, params);
            Gson gson = new Gson();
            String[] arr = gson.fromJson(response, String[].class);
            List<Province> list = new ArrayList<Province>();
            for (String province : arr) {
                list.add(new Province(province));
            }
            return list;
        } catch (Exception e) {
            throw new QuickPayException();
        }
    }

    @Override
    public List<City> getCity(String province) {
        String url = quickPayConfigStorage.getBankbaseUrl() + "/city/province/cities.json";

        Map<String, String> params = new LinkedHashMap<>();
        params.put("appkey", quickPayConfigStorage.getBankbaseKey());
        params.put("province", province);
        params.put("sig", createSign(params));

        try {
            String response = postEngine.get(url, params);
            Gson gson = new Gson();
            City[] arr = gson.fromJson(response, City[].class);
            return Arrays.asList(arr);
        } catch (Exception e) {
            throw new QuickPayException();
        }
    }

    @Override
    public List<SubBank> getBranchBank(String city_code, String bank_id) {
        String url = quickPayConfigStorage.getBankbaseUrl() + "/bank/search.json";

        Map<String, String> params = new LinkedHashMap<>();
        params.put("appkey", quickPayConfigStorage.getBankbaseKey());
        params.put("city_code", city_code);
        params.put("bank_id", bank_id);
        params.put("size", "500");
        params.put("sig", createSign(params));

        try {
            String response = postEngine.get(url, params);
            Gson gson = new Gson();
            SubBank[] arr = gson.fromJson(response, SubBank[].class);
            return Arrays.asList(arr);
        } catch (Exception e) {
            throw new QuickPayException();
        }
    }

    @Override
    public Map<String, Bank> getBank() {
        String url = quickPayConfigStorage.getBankbaseUrl() + "/bank/ids.json";

        Map<String, String> params = new LinkedHashMap<>();
        params.put("appkey", quickPayConfigStorage.getBankbaseKey());
        params.put("sig", createSign(params));

        try {
            String response = postEngine.get(url, params);
            Gson gson = new Gson();
            Map<String, Bank> decoded = gson.fromJson(response, new TypeToken<Map<String, Bank>>() {
            }.getType());
            return decoded;
        } catch (Exception e) {
            throw new QuickPayException();
        }
    }

    private String createSign(Map<String, String> params) {
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
        String sign = hmacSha1(toSign.toString(), quickPayConfigStorage.getBankbaseKey());
        return sign;
    }

    private String hmacSha1(String value, String key) {
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
}
