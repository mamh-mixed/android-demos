package com.cardinfolink.cashiersdk.util;

import com.cardinfolink.cashiersdk.model.ResultData;

import org.json.JSONArray;
import org.json.JSONObject;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.HashMap;
import java.util.Iterator;
import java.util.List;
import java.util.Map;

public class MapUtil {
    /**
     * Json 转成 Map<>
     *
     * @param jsonStr
     * @return
     */
    public static Map<String, Object> getMapForJson(String jsonStr) {
        JSONObject jsonObject;
        try {
            jsonObject = new JSONObject(jsonStr);

            Iterator<String> keyIter = jsonObject.keys();
            String key;
            Object value;
            Map<String, Object> valueMap = new HashMap<String, Object>();
            while (keyIter.hasNext()) {
                key = keyIter.next();
                value = jsonObject.get(key);
                valueMap.put(key, value);
            }
            return valueMap;
        } catch (Exception e) {
            // TODO: handle exception
            e.printStackTrace();

        }
        return null;
    }

    /**
     * Json 转成 List<Map<>>
     *
     * @param jsonStr
     * @return
     */
    public static List<Map<String, Object>> getlistForJson(String jsonStr) {
        List<Map<String, Object>> list = null;
        try {
            JSONArray jsonArray = new JSONArray(jsonStr);
            JSONObject jsonObj;
            list = new ArrayList<Map<String, Object>>();
            for (int i = 0; i < jsonArray.length(); i++) {
                jsonObj = (JSONObject) jsonArray.get(i);
                list.add(getMapForJson(jsonObj.toString()));
            }
        } catch (Exception e) {
            // TODO: handle exception
            e.printStackTrace();
        }
        return list;
    }

    public static String getSignString(Map<String, Object> map) {
        StringBuilder sb = new StringBuilder();
        Object[] unsort_key = map.keySet().toArray();
        Arrays.sort(unsort_key);
        for (int i = 0; i < unsort_key.length - 1; i++) {
            sb.append(unsort_key[i].toString());
            sb.append('=');
            sb.append(map.get(unsort_key[i]).toString());
            sb.append('&');
        }
        sb.append(unsort_key[unsort_key.length - 1].toString());
        sb.append('=');
        sb.append(map.get(unsort_key[unsort_key.length - 1]).toString());
        return sb.toString();
    }

    public static ResultData getResultData(Map<String, Object> map) {

        ResultData resultData = new ResultData();
        resultData.busicd = (String) map.get("busicd");
        resultData.respcd = (String) map.get("respcd");
        resultData.chcd = (String) map.get("chcd");
        resultData.txamt = (String) map.get("txamt");
        resultData.channelOrderNum = (String) map.get("channelOrderNum");
        resultData.consumerAccount = (String) map.get("consumerAccount");
        resultData.consumerId = (String) map.get("consumerId");
        resultData.errorDetail = (String) map.get("errorDetail");
        resultData.orderNum = (String) map.get("orderNum");
        resultData.chcdDiscount = (String) map.get("chcdDiscount");
        resultData.merDiscount = (String) map.get("merDiscount");
        resultData.qrcode = (String) map.get("qrcode");
        resultData.origOrderNum = (String) map.get("origOrderNum");
        resultData.cardId = (String) map.get("cardId");
        resultData.cardInfo = (String) map.get("cardInfo");
        resultData.scanCodeId = (String) map.get("scanCodeId");

        resultData.txamt = TxamtUtil.getNormal(resultData.txamt);
        return resultData;
    }
}
