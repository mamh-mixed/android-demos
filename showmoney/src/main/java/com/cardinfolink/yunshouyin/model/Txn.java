package com.cardinfolink.yunshouyin.model;

public class Txn {
    private String response;
    private String system_date;
    private String consumerAccount;
    private QRequest mRequest;

    public String getResponse() {
        return response;
    }

    public void setResponse(String response) {
        this.response = response;
    }

    public String getSystem_date() {
        return system_date;
    }

    public void setSystem_date(String system_date) {
        this.system_date = system_date;
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
