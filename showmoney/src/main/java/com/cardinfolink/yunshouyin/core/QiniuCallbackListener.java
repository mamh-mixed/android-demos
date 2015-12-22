package com.cardinfolink.yunshouyin.core;

public interface QiniuCallbackListener {
    void onComplete();

    void onFailure(Exception ex);
}
