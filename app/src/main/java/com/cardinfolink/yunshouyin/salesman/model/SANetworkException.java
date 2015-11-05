package com.cardinfolink.yunshouyin.salesman.model;

public class SANetworkException extends RuntimeException {
    private String errorCode;
    public SANetworkException(String errorCode){
        super(errorCode);
        this.errorCode = errorCode;
    }
}