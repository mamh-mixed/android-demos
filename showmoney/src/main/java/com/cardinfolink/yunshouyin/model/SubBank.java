package com.cardinfolink.yunshouyin.model;


import com.google.gson.annotations.SerializedName;

public class SubBank {
    @SerializedName("bank_name")
    private String bankName;

    @SerializedName("city_code")
    private String cityCode;

    @SerializedName("one_bank_no")
    private String oneBankNo;

    @SerializedName("two_bank_no")
    private String twoBankNo;


    public SubBank(String bankName, String cityCode, String oneBankNo, String twoBankNo) {
        this.bankName = bankName;
        this.cityCode = cityCode;
        this.oneBankNo = oneBankNo;
        this.twoBankNo = twoBankNo;
    }

    public String getBankName() {
        return bankName;
    }

    public void setBankName(String bankName) {
        this.bankName = bankName;
    }

    public String getCityCode() {
        return cityCode;
    }

    public void setCityCode(String cityCode) {
        this.cityCode = cityCode;
    }

    public String getOneBankNo() {
        return oneBankNo;
    }

    public void setOneBankNo(String oneBankNo) {
        this.oneBankNo = oneBankNo;
    }

    public String getTwoBankNo() {
        return twoBankNo;
    }

    public void setTwoBankNo(String twoBankNo) {
        this.twoBankNo = twoBankNo;
    }
}
