package com.cardinfolink.yunshouyin.model;

import com.google.gson.annotations.SerializedName;

public class Bank {
    private String id;

    @SerializedName("bank_name")
    private String bankName;

    public Bank(String id, String bankName) {
        this.id = id;
        this.bankName = bankName;
    }

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getBankName() {
        return bankName;
    }

    public void setBankName(String bankName) {
        this.bankName = bankName;
    }
}
