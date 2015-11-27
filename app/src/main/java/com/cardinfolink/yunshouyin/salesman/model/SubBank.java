package com.cardinfolink.yunshouyin.salesman.model;


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

    private String bankId;  //对应的大银行号，json里不会返回这个

    public SubBank(String bankName, String cityCode, String oneBankNo, String twoBankNo, String bankId) {
        this.bankName = bankName;
        this.cityCode = cityCode;
        this.oneBankNo = oneBankNo;
        this.twoBankNo = twoBankNo;
        this.bankId = bankId;
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

    public String getBankId() {
        return bankId;
    }

    public void setBankId(String bankId) {
        this.bankId = bankId;
    }
}
