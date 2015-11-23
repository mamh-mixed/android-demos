package com.cardinfolink.yunshouyin.model;

import com.google.gson.annotations.SerializedName;

public class Txn {
    private String response;

    @SerializedName("system_date")
    private String systemDate;

    private String consumerAccount;

    @SerializedName("m_request")
    private QRequest mRequest;

    public String getResponse() {
        return response;
    }

    public void setResponse(String response) {
        this.response = response;
    }

    public String getSystemDate() {
        return systemDate;
    }

    public void setSystemDate(String systemDate) {
        this.systemDate = systemDate;
    }

    public String getConsumerAccount() {
        return consumerAccount;
    }

    public void setConsumerAccount(String consumerAccount) {
        this.consumerAccount = consumerAccount;
    }

    public QRequest getmRequest() {
        return mRequest;
    }

    public void setmRequest(QRequest mRequest) {
        this.mRequest = mRequest;
    }
}
