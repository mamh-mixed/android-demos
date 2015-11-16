package com.cardinfolink.yunshouyin.salesman.model;

import com.google.gson.Gson;

import java.util.Arrays;
import java.util.Date;

public class User {
    private String username;
    private String password;
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

    private String province;
    private String city;
    private String branchBank;
    private String bankNo;

    /**
     * NEW fields for Salesman
     */
    private String accessToken;
    private String imageUrl;
    private String[] images;
    private Date createTime;
    private String merName;
    private String signKey;


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

    public String getAccessToken() {
        return accessToken;
    }

    public void setAccessToken(String accessToken) {
        this.accessToken = accessToken;
    }

    public String getImageUrl() {
        return imageUrl;
    }

    public void setImageUrl(String imageUrl) {
        this.imageUrl = imageUrl;
    }

    public String getSignKey() {
        return signKey;
    }

    public void setSignKey(String signKey) {
        this.signKey = signKey;
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

    public void setIsAutoLogin(boolean isAutoLogin) {
        this.isAutoLogin = isAutoLogin;
    }

    public String[] getImages() {
        return images;
    }

    public void setImages(String[] images) {
        this.images = images;
    }

    public Date getCreateTime() {
        return createTime;
    }

    public void setCreateTime(Date createTime) {
        this.createTime = createTime;
    }

    public String getMerName() {
        return merName;
    }

    public void setMerName(String merName) {
        this.merName = merName;
    }

    public String getJsonString() {
        Gson gson = new Gson();
        return gson.toJson(this);
    }

    @Override
    public String toString() {
        return "User{" +
                "username='" + username + '\'' +
                ", password='" + password + '\'' +
                ", bankOpen='" + bankOpen + '\'' +
                ", payee='" + payee + '\'' +
                ", payeeCard='" + payeeCard + '\'' +
                ", phoneNum='" + phoneNum + '\'' +
                ", clientid='" + clientid + '\'' +
                ", limitEmail='" + limitEmail + '\'' +
                ", objectId='" + objectId + '\'' +
                ", isAutoLogin=" + isAutoLogin +
                ", limitName='" + limitName + '\'' +
                ", limitPhone='" + limitPhone + '\'' +
                ", limit='" + limit + '\'' +
                ", province='" + province + '\'' +
                ", city='" + city + '\'' +
                ", branchBank='" + branchBank + '\'' +
                ", bankNo='" + bankNo + '\'' +
                ", accessToken='" + accessToken + '\'' +
                ", imageUrl='" + imageUrl + '\'' +
                ", images=" + Arrays.toString(images) +
                ", createTime=" + createTime +
                ", merName='" + merName + '\'' +
                ", signKey='" + signKey + '\'' +
                '}';
    }
}
