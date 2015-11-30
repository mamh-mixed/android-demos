package com.cardinfolink.yunshouyin.core;


import com.cardinfolink.yunshouyin.api.QuickPayException;

public class AsyncTaskResult<T> {
    private T result;
    private QuickPayException exception;

    public AsyncTaskResult(T result) {
        this.result = result;
        this.exception = null;
    }

    public AsyncTaskResult(QuickPayException exception) {
        this.result = null;
        this.exception = exception;
    }

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
