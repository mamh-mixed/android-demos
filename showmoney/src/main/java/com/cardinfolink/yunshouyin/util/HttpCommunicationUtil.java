package com.cardinfolink.yunshouyin.util;

import android.util.Log;

import com.cardinfolink.yunshouyin.data.RequestParam;

import org.apache.http.HttpResponse;
import org.apache.http.client.HttpClient;
import org.apache.http.client.entity.UrlEncodedFormEntity;
import org.apache.http.client.methods.HttpGet;
import org.apache.http.client.methods.HttpPost;
import org.apache.http.impl.client.DefaultHttpClient;
import org.apache.http.util.EntityUtils;

import java.net.URI;


public class HttpCommunicationUtil {
    private static final String TAG = "CommunicationUtil";

    public static void sendDataToServer(
            final RequestParam requestParam, final CommunicationListener listener) {

        new Thread(new Runnable() {

            @Override
            public void run() {
                try {
                    URI baseUrl = new URI(requestParam.getUrl());
                    Log.i(TAG, "url = " + requestParam.getUrl());
                    final HttpPost postMethod = new HttpPost(baseUrl);
                    final HttpClient httpClient = new DefaultHttpClient();
                    ;
                    postMethod.setEntity(new UrlEncodedFormEntity(requestParam.getParams(),
                            "utf-8")); // 将参数填入POST Entity中
                    HttpResponse response = httpClient.execute(postMethod); // 执行POST方法
                    Log.i(TAG, "resCode = "
                            + response.getStatusLine().getStatusCode()); // 获取响应码
                    String result = EntityUtils.toString(response.getEntity(),
                            "utf-8");
                    Log.i(TAG, "result = " + result); // 获取响应内容
                    listener.onResult(result);

                } catch (Exception e) {
                    Log.i(TAG, "error = " + e.getMessage());
                    listener.onError("网络错误");
                    e.printStackTrace();
                }
            }
        }).start();
    }


    public static void sendGetDataToServer(
            final RequestParam requestParam, final CommunicationListener listener) {

        new Thread(new Runnable() {

            @Override
            public void run() {
                try {
                    URI baseUrl = new URI(requestParam.getUrl());
                    Log.i(TAG, "url = " + requestParam.getUrl());
                    final HttpGet getMethod = new HttpGet(baseUrl);
                    final HttpClient httpClient = new DefaultHttpClient();
                    HttpResponse response = httpClient.execute(getMethod); // 执行POST方法
                    Log.i(TAG, "resCode = "
                            + response.getStatusLine().getStatusCode()); // 获取响应码
                    String result = EntityUtils.toString(response.getEntity(),
                            "utf-8");
                    Log.i(TAG, "result = " + result); // 获取响应内容
                    listener.onResult(result);

                } catch (Exception e) {
                    Log.i(TAG, "error = " + e.getMessage());
                    listener.onError("网络错误");
                    e.printStackTrace();
                }
            }
        }).start();
    }

}
