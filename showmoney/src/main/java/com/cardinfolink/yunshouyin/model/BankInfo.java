package com.cardinfolink.yunshouyin.model;

import com.google.gson.annotations.SerializedName;

public class BankInfo {
    /**
     * {
     * "state": "success",
     * "count": 0,
     * "size": 0,
     * "refdcount": 0,
     * "info": {这里对应BankInfo类
     * "bank_open": "中国工商银行",
     * "payee": "司瑞华",
     * "payee_card": "6222021001090585214",
     * "phone_num": "13821936284",
     * "province": "天津市",
     * "city": "天津市",
     * "branch_bank": "中国工商银行天津市金街支行",
     * "bankNo": "102110000076|102100099996"
     * }
     * }
     */


    @SerializedName("bank_open")
    private String bankOpen;

    private String payee;

    @SerializedName("payee_card")
    private String payeeCard;

    @SerializedName("phone_num")
    private String phoneNum;

    @SerializedName("province")
    private String province;

    @SerializedName("city")
    private String city;

    @SerializedName("branch_bank")
    private String branchBank;

    @SerializedName("bankNo")
    private String bankNo;

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

    public String getProvince() {
        return province;
    }

    public void setProvince(String province) {
        this.province = province;
    }

    public String getCity() {
        return city;
    }

    public void setCity(String city) {
        this.city = city;
    }

    public String getBranchBank() {
        return branchBank;
    }

    public void setBranchBank(String branchBank) {
        this.branchBank = branchBank;
    }

    public String getBankNo() {
        return bankNo;
    }

    public void setBankNo(String bankNo) {
        this.bankNo = bankNo;
    }
}
