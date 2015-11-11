package com.cardinfolink.yunshouyin.core;


import com.cardinfolink.yunshouyin.api.QuickPayException;

public interface QuickPayCallbackListener<T> {
    void onSuccess(T data);

    void onFailure(QuickPayException ex);
}
