package com.cardinfolink.cashiersdk.util;

import android.annotation.SuppressLint;

import com.cardinfolink.cashiersdk.model.InitData;
import com.cardinfolink.cashiersdk.model.OrderData;
import com.cardinfolink.cashiersdk.sdk.CashierSdk;

import org.json.JSONException;
import org.json.JSONObject;


public class ParamsUtil {
    public static final String SIGN_TYPE = CashierSdk.SIGN_TYPE;

    @SuppressLint("SimpleDateFormat")
    public static JSONObject getPay(InitData initData, OrderData orderData) {
        JSONObject json = new JSONObject();
        try {
            json.put("txndir", "Q");
            json.put("busicd", "PURC");
            json.put("inscd", initData.inscd);
            json.put("mchntid", initData.mchntid);
            json.put("txamt", orderData.txamt);
            json.put("orderNum", orderData.orderNum);
            json.put("scanCodeId", orderData.scanCodeId);
            json.put("terminalid", initData.terminalid);
            json.put("currency", orderData.currency);
            json.put("tradeFrom", "android");
            if (orderData.goodsInfo != null) {
                json.put("goodsInfo", orderData.goodsInfo);
            }

            json.put("sign", getSign(MapUtil.getSignString(MapUtil.getMapForJson(json.toString())), initData.signKey, SIGN_TYPE));

        } catch (JSONException e) {
            // TODO Auto-generated catch block
            e.printStackTrace();
        }


        return json;
    }

    public static JSONObject getPrePay(InitData initData, OrderData orderData) {
        JSONObject json = new JSONObject();
        try {
            json.put("txndir", "Q");
            json.put("busicd", "PAUT");
            json.put("inscd", initData.inscd);
            json.put("mchntid", initData.mchntid);
            json.put("txamt", orderData.txamt);
            json.put("orderNum", orderData.orderNum);
            json.put("chcd", orderData.chcd);
            json.put("terminalid", initData.terminalid);
            json.put("currency", orderData.currency);
            json.put("tradeFrom", "android");
            if (orderData.goodsInfo != null) {
                json.put("goodsInfo", orderData.goodsInfo);
            }

            json.put("sign", getSign(MapUtil.getSignString(MapUtil.getMapForJson(json.toString())), initData.signKey, SIGN_TYPE));

        } catch (JSONException e) {
            // TODO Auto-generated catch block
            e.printStackTrace();
        }


        return json;
    }


    public static JSONObject getQy(InitData initData, OrderData orderData) {
        JSONObject json = new JSONObject();
        try {
            json.put("txndir", "Q");
            json.put("busicd", "INQY");
            json.put("inscd", initData.inscd);
            json.put("mchntid", initData.mchntid);
            json.put("origOrderNum", orderData.origOrderNum);
            json.put("terminalid", initData.terminalid);
            json.put("sign", getSign(MapUtil.getSignString(MapUtil.getMapForJson(json.toString())), initData.signKey, SIGN_TYPE));

        } catch (JSONException e) {
            // TODO Auto-generated catch block
            e.printStackTrace();
        }


        return json;
    }

    public static JSONObject getVoid(InitData initData, OrderData orderData) {
        JSONObject json = new JSONObject();
        try {
            json.put("txndir", "Q");
            json.put("busicd", "VOID");
            json.put("inscd", initData.inscd);
            json.put("mchntid", initData.mchntid);
            json.put("origOrderNum", orderData.origOrderNum);
            json.put("orderNum", orderData.orderNum);
            json.put("terminalid", initData.terminalid);
            json.put("tradeFrom", "android");
            json.put("sign", getSign(MapUtil.getSignString(MapUtil.getMapForJson(json.toString())), initData.signKey, SIGN_TYPE));

        } catch (JSONException e) {
            // TODO Auto-generated catch block
            e.printStackTrace();
        }


        return json;
    }


    public static JSONObject getRefd(InitData initData, OrderData orderData) {
        JSONObject json = new JSONObject();
        try {
            json.put("txndir", "Q");
            json.put("busicd", "REFD");
            json.put("inscd", initData.inscd);
            json.put("mchntid", initData.mchntid);
            json.put("origOrderNum", orderData.origOrderNum);
            json.put("orderNum", orderData.orderNum);
            json.put("txamt", orderData.txamt);
            json.put("terminalid", initData.terminalid);
            json.put("currency", orderData.currency);
            json.put("tradeFrom", "android");
            json.put("sign", getSign(MapUtil.getSignString(MapUtil.getMapForJson(json.toString())), initData.signKey, SIGN_TYPE));

        } catch (JSONException e) {
            // TODO Auto-generated catch block
            e.printStackTrace();
        }


        return json;
    }


    public static JSONObject getVeri(InitData initData, OrderData orderData) {
        JSONObject json = new JSONObject();
        try {
            json.put("txndir", "Q");
            json.put("busicd", "VERI");
            json.put("inscd", initData.inscd);
            json.put("mchntid", initData.mchntid);
            json.put("scanCodeId", orderData.scanCodeId);
            json.put("orderNum", orderData.orderNum);
            json.put("terminalid", initData.terminalid);
            json.put("tradeFrom", "android");
            json.put("sign", getSign(MapUtil.getSignString(MapUtil.getMapForJson(json.toString())), initData.signKey, SIGN_TYPE));

        } catch (JSONException e) {
            // TODO Auto-generated catch block
            e.printStackTrace();
        }


        return json;
    }


    @SuppressWarnings("unchecked")
    public static String getSign(String str, String key, String signType) {
        str = str + key;
        String sign = EncoderUtil.Encrypt(str, signType);
        return sign;
    }
}
