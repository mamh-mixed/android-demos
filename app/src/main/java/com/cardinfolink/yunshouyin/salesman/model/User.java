package com.cardinfolink.yunshouyin.salesman.model;

import com.google.gson.Gson;

import java.util.Date;

public class User {
    private String username;
    private String password;
    private String bank_open;
    private String payee;
    private String payee_card;
    private String phone_num;
    private String clientid;
    private String limit_email;
    private String object_id;
    private boolean isAutoLogin;
    private String limit_name;

    private String limit_phone;
    private String limit = "true";

    private String province;
    private String city;
    private String branch_bank;
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

    public String getBranch_bank() {
        return branch_bank;
    }

    public void setBranch_bank(String branch_bank) {
        this.branch_bank = branch_bank;
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

    public String getLimit_email() {
        return limit_email;
    }

    public void setLimit_email(String limit_email) {
        this.limit_email = limit_email;
    }

    public String getLimit_name() {
        return limit_name;
    }

    public void setLimit_name(String limit_name) {
        this.limit_name = limit_name;
    }

    public String getLimit_phone() {
        return limit_phone;
    }

    public void setLimit_phone(String limit_phone) {
        this.limit_phone = limit_phone;
    }

    public boolean isAutoLogin() {
        return isAutoLogin;
    }

    public void setAutoLogin(boolean isAutoLogin) {
        this.isAutoLogin = isAutoLogin;
    }

    public String getObject_id() {
        return object_id;
    }

    public void setObject_id(String object_id) {
        this.object_id = object_id;
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

    public String getBank_open() {
        return bank_open;
    }

    public void setBank_open(String bank_open) {
        this.bank_open = bank_open;
    }

    public String getPayee() {
        return payee;
    }

    public void setPayee(String payee) {
        this.payee = payee;
    }

    public String getPayee_card() {
        return payee_card;
    }

    public void setPayee_card(String payee_card) {
        this.payee_card = payee_card;
    }

    public String getPhone_num() {
        return phone_num;
    }

    public void setPhone_num(String phone_num) {
        this.phone_num = phone_num;
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

    public String getJsonString(){
        Gson gson = new Gson();
        return gson.toJson(this);
    }
}
