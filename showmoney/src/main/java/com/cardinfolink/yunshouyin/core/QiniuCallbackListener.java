package com.cardinfolink.yunshouyin.core;

import com.qiniu.android.http.ResponseInfo;

import org.json.JSONObject;

public interface QiniuCallbackListener {
    void onComplete(String key, ResponseInfo info, JSONObject response);

    void onFailure(Exception ex);

    void onProgress(String key, double percent);
}
