package com.cardinfolink.yunshouyin.data;


import org.apache.http.NameValuePair;

import java.util.List;

public class RequestParam {
    private String url;
    private List<NameValuePair> params;

    public String getUrl() {
        return url;
    }

    public void setUrl(String url) {
        this.url = url;
    }

    public List<NameValuePair> getParams() {
        return params;
    }

    public void setParams(List<NameValuePair> params) {
        this.params = params;
    }


}
