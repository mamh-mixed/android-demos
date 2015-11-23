package com.cardinfolink.yunshouyin.model;

import com.google.gson.annotations.SerializedName;

public class BankInfo {
    @SerializedName("bank_open")
    private String bankOpen;
    private String payee;
    @SerializedName("payee_card")
    private String payeeCard;
    @SerializedName("phone_num")
    private String phoneNum;

    public String getBankOpen() {
        return bankOpen;
    }

    public void setBankOpen(String bankOpen) {
        this.bankOpen = bankOpen;
    }

    public String getPayee() {
        return payee;
    }

    public void setPayee(String payee) {
        this.payee = payee;
    }

    public String getPayeeCard() {
        return payeeCard;
    }

    public void setPayeeCard(String payeeCard) {
        this.payeeCard = payeeCard;
    }

    public String getPhoneNum() {
        return phoneNum;
    }

    public void setPhoneNum(String phoneNum) {
        this.phoneNum = phoneNum;
    }
}
