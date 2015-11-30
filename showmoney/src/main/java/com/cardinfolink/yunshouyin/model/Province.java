package com.cardinfolink.yunshouyin.model;

/**
 * Created by mamh on 15-11-26.
 * 对应的province的实体类
 */
public class Province {
    private String provinceName;

    public Province() {

    }

    public Province(String provinceName) {
        this.provinceName = provinceName;
    }

    public String getProvinceName() {
        return provinceName;
    }

    public void setProvinceName(String provinceName) {
        this.provinceName = provinceName;
    }
}
