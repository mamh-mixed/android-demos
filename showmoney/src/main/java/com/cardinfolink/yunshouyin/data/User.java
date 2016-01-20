package com.cardinfolink.yunshouyin.data;

public class User {
    private String username;
    private String password;
    private String activate;
    private String clientid;
    private String limit = "true";
    private String createTime;
    private String signKey;
    private String inscd;
    private String objectId;

    private String province;
    private String city;
    private String branchBank;
    private String bankNo;
    private String bankOpen;
    private String payee;
    private String payeeCard;
    private String phoneNum;

    private String limitEmail;
    private boolean isAutoLogin;
    private String limitName;
    private String limitPhone;

    private String payUrl;

    private String limitAmt;

    private String merName;

    public String getMerName() {
        return merName;
    }

    public void setMerName(String merName) {
        this.merName = merName;
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

    public String getLimitEmail() {
        return limitEmail;
    }

    public void setLimitEmail(String limitEmail) {
        this.limitEmail = limitEmail;
    }

    public String getLimitName() {
        return limitName;
    }

    public void setLimitName(String limitName) {
        this.limitName = limitName;
    }

    public String getLimitPhone() {
        return limitPhone;
    }

    public void setLimitPhone(String limitPhone) {
        this.limitPhone = limitPhone;
    }

    public boolean isAutoLogin() {
        return isAutoLogin;
    }

    public void setAutoLogin(boolean isAutoLogin) {
        this.isAutoLogin = isAutoLogin;
    }

    public String getObjectId() {
        return objectId;
    }

    public void setObjectId(String objectId) {
        this.objectId = objectId;
    }

    public String getLimit() {
        return limit;
    }


    public void setLimit(String limit) {
        this.limit = limit;
    }


    public String getClientid() {
        return clientid;
    }


    public void setClientid(String clientid) {
        this.clientid = clientid;
    }


    public String getUsername() {
        return username;
    }

    public void setUsername(String username) {
        this.username = username;
    }

    public String getPassword() {
        return password;
    }

    public void setPassword(String password) {
        this.password = password;
    }

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


    public String getActivate() {
        return activate;
    }

    public void setActivate(String activate) {
        this.activate = activate;
    }

    public String getCreateTime() {
        return createTime;
    }

    public void setCreateTime(String createTime) {
        this.createTime = createTime;
    }

    public String getSignKey() {
        return signKey;
    }

    public void setSignKey(String signKey) {
        this.signKey = signKey;
    }

    public String getInscd() {
        return inscd;
    }

    public void setInscd(String inscd) {
        this.inscd = inscd;
    }

    public void setIsAutoLogin(boolean isAutoLogin) {
        this.isAutoLogin = isAutoLogin;
    }

    public String getPayUrl() {
        return payUrl;
    }

    public void setPayUrl(String payUrl) {
        this.payUrl = payUrl;
    }

    public String getLimitAmt() {
        return limitAmt;
    }

    public void setLimitAmt(String limitAmt) {
        this.limitAmt = limitAmt;
    }

    @Override
    public String toString() {
        return "User{" +
                "username='" + username + '\'' +
                ", password='" + password + '\'' +
                ", activate='" + activate + '\'' +
                ", clientid='" + clientid + '\'' +
                ", limit='" + limit + '\'' +
                ", createTime='" + createTime + '\'' +
                ", signKey='" + signKey + '\'' +
                ", inscd='" + inscd + '\'' +
                ", objectId='" + objectId + '\'' +
                ", province='" + province + '\'' +
                ", city='" + city + '\'' +
                ", branchBank='" + branchBank + '\'' +
                ", bankNo='" + bankNo + '\'' +
                ", bankOpen='" + bankOpen + '\'' +
                ", payee='" + payee + '\'' +
                ", payeeCard='" + payeeCard + '\'' +
                ", phoneNum='" + phoneNum + '\'' +
                ", limitEmail='" + limitEmail + '\'' +
                ", isAutoLogin=" + isAutoLogin +
                ", limitName='" + limitName + '\'' +
                ", limitPhone='" + limitPhone + '\'' +
                ", payUrl='" + payUrl + '\'' +
                ", limitAmt='" + limitAmt + '\'' +
                ", merName='" + merName + '\'' +
                '}';
    }
}
