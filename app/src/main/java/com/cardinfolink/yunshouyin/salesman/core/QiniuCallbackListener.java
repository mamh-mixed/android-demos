package com.cardinfolink.yunshouyin.salesman.core;

public interface QiniuCallbackListener {
    void onComplete();

    void onFailure(Exception ex);
}
