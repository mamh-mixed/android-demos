package com.cardinfolink.yunshouyin.salesman.api;

import com.cardinfolink.yunshouyin.salesman.utils.ErrorUtil;

public class QuickPayException extends RuntimeException {
    //token不对,重新登录
    public static final String ACCESSTOKEN_NOT_FOUND = "accessToken_error";
    //post error
    public static final String NETWORK_ERROR = "network_error";
    //可能的情况:URL配置错误,预期得到JSON格式但是返回HTML
    public static final String CONFIG_ERROR = "config_error";
    /**
     * error code from server
     */
    private String errorCode;
    /**
     * user friendly translation
     */
    private String errorMsg;

    public QuickPayException(String errorCode) {
        super(errorCode);
        this.errorCode = errorCode;
        this.errorMsg = ErrorUtil.getErrorString(errorCode);
    }

    /**
     * default is network exception
     */
    public QuickPayException() {
        this.errorCode = NETWORK_ERROR;
        this.errorMsg = "网络错误,请检查网络是否连接";
    }

    public String getErrorCode() {
        return errorCode;
    }

    public String getErrorMsg() {
        return errorMsg;
    }
}