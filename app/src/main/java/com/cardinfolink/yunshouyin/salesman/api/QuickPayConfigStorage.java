package com.cardinfolink.yunshouyin.salesman.api;

public class QuickPayConfigStorage {
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
    private String proxy_url;
    private int proxy_port;

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

    public String getProxy_url() {
        return proxy_url;
    }

    public void setProxy_url(String proxy_url) {
        this.proxy_url = proxy_url;
    }

    public int getProxy_port() {
        return proxy_port;
    }

    public void setProxy_port(int proxy_port) {
        this.proxy_port = proxy_port;
    }
}
