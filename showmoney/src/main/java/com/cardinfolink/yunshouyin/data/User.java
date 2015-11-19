package com.cardinfolink.yunshouyin.data;

public class User {
    private String username;
    private String password;

    private String province;
    private String city;
    private String branchBank;
    private String bankNo;

    private String bankOpen;
    private String payee;
    private String payeeCard;
    private String phoneNum;
    private String clientid;
    private String limitEmail;
    private String objectId;
    private boolean isAutoLogin;
    private String limitName;

    private String limitPhone;
    private String limit = "true";

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


}
