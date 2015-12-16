package com.cardinfolink.yunshouyin.api;


import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.util.ErrorUtil;
import com.cardinfolink.yunshouyin.util.ShowMoneyApp;

public class QuickPayException extends RuntimeException {
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
        //这里把error code转换为文本可读的错误消息提示
        this.errorMsg = ErrorUtil.getErrorString(errorCode);
    }

    /**
     * 这个暂时代码里没用到这样的构造函数
     *
     * @param errorCode
     * @param errorMsg
     */
    public QuickPayException(String errorCode, String errorMsg) {
        super(errorCode);
        this.errorCode = errorCode;
        this.errorMsg = errorMsg;
    }

    /**
     * default is network exception
     * 默认显示的出错信息。
     * <string name="alert_error_network">网络错误,请检查网络是否连接</string>
     * 这个版本加入的多语言的支持，所以把这个字段移到了strings.xml文件里了
     */
    public QuickPayException() {
        super(NETWORK_ERROR);
        this.errorCode = NETWORK_ERROR;
        this.errorMsg = ShowMoneyApp.getResString(R.string.alert_error_network);
    }

    public String getErrorCode() {
        return errorCode;
    }

    public String getErrorMsg() {
        return errorMsg;
    }
}