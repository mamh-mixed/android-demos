package com.cardinfolink.yunshouyin.salesman.api;

public class QuickPayConfigStorage {
    private String proxyUrl;
    private int proxyPort;


    /**
     * access to quick pay
     */
    private String accessToken;
    /**
     * access to qiniu cloud
     */
    private String uploadToken;
    private String appKey;
    private String url;

    // bank data setting
    private String bankbaseKey;
    private String bankbaseUrl;

    public String getAccessToken() {
        return accessToken;
    }

    public void setAccessToken(String accessToken) {
        this.accessToken = accessToken;
    }

    public String getUploadToken() {
        return uploadToken;
    }

    public void setUploadToken(String uploadToken) {
        this.uploadToken = uploadToken;
    }

    public String getAppKey() {
        return appKey;
    }

    public void setAppKey(String appKey) {
        this.appKey = appKey;
    }

    public String getUrl() {
        return url;
    }

    public void setUrl(String url) {
        this.url = url;
    }

    public String getProxyUrl() {
        return proxyUrl;
    }

    public void setProxyUrl(String proxyUrl) {
        this.proxyUrl = proxyUrl;
    }

    public int getProxyPort() {
        return proxyPort;
    }

    public void setProxyPort(int proxyPort) {
        this.proxyPort = proxyPort;
    }

    public String getBankbaseKey() {
        return bankbaseKey;
    }

    public void setBankbaseKey(String bankbaseKey) {
        this.bankbaseKey = bankbaseKey;
    }

    public String getBankbaseUrl() {
        return bankbaseUrl;
    }

    public void setBankbaseUrl(String bankbaseUrl) {
        this.bankbaseUrl = bankbaseUrl;
    }
}
