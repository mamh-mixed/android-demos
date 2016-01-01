package com.cardinfolink.yunshouyin.model;

import com.google.gson.annotations.SerializedName;

public class Txn {
    private String response;

    @SerializedName("system_date")
    private String systemDate;

    private String consumerAccount;
    @SerializedName("m_request")
    private QRequest mRequest;

    private String refundAmt;

    private String transStatus;


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

    public String getRefundAmt() {
        return refundAmt;
    }

    public void setRefundAmt(String refundAmt) {
        this.refundAmt = refundAmt;
    }

    public String getTransStatus() {
        return transStatus;
    }

    public void setTransStatus(String transStatus) {
        this.transStatus = transStatus;
    }

    @Override
    public String toString() {
        return "Txn{" +
                "response='" + response + '\'' +
                ", systemDate='" + systemDate + '\'' +
                ", consumerAccount='" + consumerAccount + '\'' +
                ", mRequest=" + mRequest +
                ", refundAmt='" + refundAmt + '\'' +
                ", transStatus='" + transStatus + '\'' +
                '}';
    }
}
