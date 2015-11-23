package com.cardinfolink.yunshouyin.salesman.model;

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
