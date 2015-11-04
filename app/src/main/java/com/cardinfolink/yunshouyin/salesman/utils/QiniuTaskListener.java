package com.cardinfolink.yunshouyin.salesman.utils;

public interface QiniuTaskListener {
    void onComplete();

    void onError(Exception ex);
}
