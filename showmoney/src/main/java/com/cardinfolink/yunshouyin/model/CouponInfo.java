package com.cardinfolink.yunshouyin.model;

/**
 * Created by mamh on 16-1-6.
 * 卡券数组。每个卡券需包含：卡券类型：（减、兑、折）、卡券名称、卡券渠道等等
 * json对应的实体类
 */
public class CouponInfo {

    /**
     * 卡券类型
     */
    private String type;

    /**
     * 卡券名称
     */
    private String name;

    /**
     * 卡券渠道
     */
    private String channel;


    public String getType() {
        return type;
    }

    public void setType(String type) {
        this.type = type;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public String getChannel() {
        return channel;
    }

    public void setChannel(String channel) {
        this.channel = channel;
    }

    @Override
    public String toString() {
        return "CouponInfo{" +
                "type='" + type + '\'' +
                ", name='" + name + '\'' +
                ", channel='" + channel + '\'' +
                '}';
    }
}
