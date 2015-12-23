package com.cardinfolink.yunshouyin.core;

import com.cardinfolink.yunshouyin.api.QuickPayException;
import com.qiniu.android.http.ResponseInfo;

import org.json.JSONObject;

public interface QiniuCallbackListener {
    void onComplete(String key, ResponseInfo info, JSONObject response, int photoKey);

    void onFailure(QuickPayException ex, int photoKey);

    void onProgress(String key, double percent, int photoKey);
}
