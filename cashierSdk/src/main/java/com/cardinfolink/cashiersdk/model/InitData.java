package com.cardinfolink.cashiersdk.model;

public class InitData {
    public String inscd;                // 机构号
    public String mchntid;              // 商户号
    public String signKey;              // 秘钥
    public String terminalid;           // 设备号
    public boolean isProduce = false;   //是否是生产环境

    public InitData() {

    }

    public InitData(String inscd, String mchntid, String signKey, String terminalid, boolean isProduce) {
        this.inscd = inscd;
        this.mchntid = mchntid;
        this.signKey = signKey;
        this.terminalid = terminalid;
        this.isProduce = isProduce;
    }


    public String getInscd() {
        return inscd;
    }

    public void setInscd(String inscd) {
        this.inscd = inscd;
    }

    public String getMchntid() {
        return mchntid;
    }

    public void setMchntid(String mchntid) {
        this.mchntid = mchntid;
    }

    public String getSignKey() {
        return signKey;
    }

    public void setSignKey(String signKey) {
        this.signKey = signKey;
    }

    public String getTerminalid() {
        return terminalid;
    }

    public void setTerminalid(String terminalid) {
        this.terminalid = terminalid;
    }

    public boolean isProduce() {
        return isProduce;
    }

    public void setIsProduce(boolean isProduce) {
        this.isProduce = isProduce;
    }
}
