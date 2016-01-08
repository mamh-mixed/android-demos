package com.cardinfolink.yunshouyin.model;

import com.google.gson.annotations.SerializedName;

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

    /**
     * 请求来源，例如android
     */
    private String tradeFrom;

    /**
     * 订单号
     */
    private String orderNum;

    /**
     * 响应码，例如00
     */
    private String response;

    /**
     * 交易时间，格式yyyyMMddHHmmss，例如20151229190742
     */
    @SerializedName("system_date")
    private String systemDate;

    /**
     * 终端号
     */
    private String terminalid;

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

    public String getTradeFrom() {
        return tradeFrom;
    }

    public void setTradeFrom(String tradeFrom) {
        this.tradeFrom = tradeFrom;
    }

    public String getResponse() {
        return response;
    }

    public void setResponse(String response) {
        this.response = response;
    }

    public String getSystemDate() {
        return systemDate;
    }

    public void setSystemDate(String systemDate) {
        this.systemDate = systemDate;
    }


    public String getOrderNum() {
        return orderNum;
    }

    public void setOrderNum(String orderNum) {
        this.orderNum = orderNum;
    }

    public String getTerminalid() {
        return terminalid;
    }

    public void setTerminalid(String terminalid) {
        this.terminalid = terminalid;
    }

    @Override
    public String toString() {
        return "CouponInfo{" +
                "type='" + type + '\'' +
                ", name='" + name + '\'' +
                ", channel='" + channel + '\'' +
                ", tradeFrom='" + tradeFrom + '\'' +
                ", orderNum='" + orderNum + '\'' +
                ", response='" + response + '\'' +
                ", systemDate='" + systemDate + '\'' +
                ", terminalid='" + terminalid + '\'' +
                '}';
    }
}
