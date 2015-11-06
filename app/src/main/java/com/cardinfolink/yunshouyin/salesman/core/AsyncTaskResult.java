package com.cardinfolink.yunshouyin.salesman.core;

import com.cardinfolink.yunshouyin.salesman.api.QuickPayException;

public class AsyncTaskResult<T> {
    private T result;
    private QuickPayException exception;

    public AsyncTaskResult(T result, QuickPayException exception) {
        this.result = result;
        this.exception = exception;
    }

    public T getResult() {
        return result;
    }

    public QuickPayException getException() {
        return exception;
    }
}
