package com.cardinfolink.yunshouyin.api;

public class QuickPayConfigStorage {
    // proxy setting
    private String proxy_url;
    private int proxy_port;

    // quick pay setting
    private String appKey;
    private String url;

    // bank data setting
    private String bankbaseKey;
    private String bankbaseUrl;


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
