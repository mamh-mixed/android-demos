package com.cardinfolink.yunshouyin.core;

import com.cardinfolink.yunshouyin.api.QuickPayException;

public interface QiniuCallbackListener {
    void onComplete(String pattern, int photoKey);

    void onFailure(QuickPayException ex, int photoKey);

    void onProgress(String key, double percent, int photoKey);
}
