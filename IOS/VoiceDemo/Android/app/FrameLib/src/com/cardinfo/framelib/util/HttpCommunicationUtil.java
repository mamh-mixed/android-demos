package com.cardinfo.framelib.util;

import java.net.URI;


import org.apache.http.HttpResponse;
import org.apache.http.client.HttpClient;
import org.apache.http.client.entity.UrlEncodedFormEntity;
import org.apache.http.client.methods.HttpPost;
import org.apache.http.impl.client.DefaultHttpClient;
import org.apache.http.util.EntityUtils;

import com.cardinfo.framelib.listener.CommunicationListener;
import com.cardinfo.framelib.model.RequestParam;

import android.util.Log;



public class HttpCommunicationUtil {
	private static final String TAG = "CommunicationUtil";
	public static void sendDataToServer(
			final RequestParam requestParam,final CommunicationListener listener) {

		new Thread(new Runnable() {

			@Override
			public void run() {
				try {
					URI baseUrl = new URI(requestParam.getUrl());
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
					listener.onError("网络错误");
					e.printStackTrace();
				}
			}
		}).start();
	}
	
}
