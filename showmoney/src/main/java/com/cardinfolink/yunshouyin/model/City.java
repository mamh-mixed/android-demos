package com.cardinfolink.yunshouyin.model;

import com.google.gson.annotations.SerializedName;

public class City {
    private String id;

    @SerializedName("city_code")
    private String cityCode;

    @SerializedName("province_code")
    private String provinceCode;

    @SerializedName("city_name")
    private String cityName;

    @SerializedName("city_jb")
    private String cityJb;

    private String city;
    private String province;


    public City(String id, String cityCode, String provinceCode, String cityName, String cityJb, String city, String province) {
        this.id = id;
        this.cityCode = cityCode;
        this.provinceCode = provinceCode;
        this.cityName = cityName;
        this.cityJb = cityJb;
        this.city = city;
        this.province = province;
    }

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getCityCode() {
        return cityCode;
    }

    public void setCityCode(String cityCode) {
        this.cityCode = cityCode;
    }

    public String getProvinceCode() {
        return provinceCode;
    }

    public void setProvinceCode(String provinceCode) {
        this.provinceCode = provinceCode;
    }

    public String getCityName() {
        return cityName;
    }

    public void setCityName(String cityName) {
        this.cityName = cityName;
    }

    public String getCityJb() {
        return cityJb;
    }

    public void setCityJb(String cityJb) {
        this.cityJb = cityJb;
    }

    public String getCity() {
        return city;
    }

    public void setCity(String city) {
        this.city = city;
    }

    public String getProvince() {
        return province;
    }

    public void setProvince(String province) {
        this.province = province;
    }
}
