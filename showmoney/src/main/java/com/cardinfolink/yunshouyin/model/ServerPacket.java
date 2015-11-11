package com.cardinfolink.yunshouyin.model;

import com.cardinfolink.yunshouyin.api.QuickPayException;
import com.cardinfolink.yunshouyin.data.User;
import com.google.gson.Gson;
import com.google.gson.GsonBuilder;

public class ServerPacket {
    private String state;
    private String error;
    private User user;
    private int count;
    private double total;
    private int refdcount;
    private double refdtotal;
    private int size;
    private BankInfo info;

    public static ServerPacket getServerPacketFrom(String json) {
        try {
            Gson gson = new GsonBuilder().setDateFormat("yyyy-MM-dd HH:mm:ss").create();
            ServerPacket packet = gson.fromJson(json, ServerPacket.class);
            return packet;
        } catch (Exception ex) {
            throw new QuickPayException(QuickPayException.CONFIG_ERROR);
        }
    }

    public String getState() {
        return state;
    }

    public void setState(String state) {
        this.state = state;
    }

    public String getError() {
        return error;
    }

    public void setError(String error) {
        this.error = error;
    }

    public User getUser() {
        return user;
    }

    public void setUser(User user) {
        this.user = user;
    }

    public int getCount() {
        return count;
    }

    public void setCount(int count) {
        this.count = count;
    }

    public double getTotal() {
        return total;
    }

    public void setTotal(double total) {
        this.total = total;
    }

    public int getRefdcount() {
        return refdcount;
    }

    public void setRefdcount(int refdcount) {
        this.refdcount = refdcount;
    }

    public double getRefdtotal() {
        return refdtotal;
    }

    public void setRefdtotal(double refdtotal) {
        this.refdtotal = refdtotal;
    }

    public int getSize() {
        return size;
    }

    public void setSize(int size) {
        this.size = size;
    }

    public BankInfo getInfo() {
        return info;
    }

    public void setInfo(BankInfo info) {
        this.info = info;
    }
}
