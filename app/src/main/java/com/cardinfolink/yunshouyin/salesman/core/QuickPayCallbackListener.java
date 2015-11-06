package com.cardinfolink.yunshouyin.salesman.core;

import com.cardinfolink.yunshouyin.salesman.api.QuickPayException;

public interface QuickPayCallbackListener<T> {
    void onSuccess(T data);

    void onFailure(QuickPayException ex);
}
