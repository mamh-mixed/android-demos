package com.cardinfolink.yunshouyin.salesman.model;

import com.cardinfolink.yunshouyin.salesman.api.QuickPayException;
import com.google.gson.Gson;
import com.google.gson.GsonBuilder;

public class SAServerPacket {
    private String state;
    private String error;
    private User user;
    private User[] users;
    private String accessToken;
    private String uploadToken;
    private String downloadUrl;

    public static SAServerPacket getServerPacketFrom(String json) {
        try {
            Gson gson = new GsonBuilder().setDateFormat("yyyy-MM-dd HH:mm:ss").create();
            SAServerPacket packet = gson.fromJson(json, SAServerPacket.class);
            return packet;
        }catch (Exception ex){
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

    public User[] getUsers() {
        return users;
    }

    public void setUsers(User[] users) {
        this.users = users;
    }

    public String getAccessToken() {
        return accessToken;
    }

    public void setAccessToken(String accessToken) {
        this.accessToken = accessToken;
    }

    public String getUploadToken() {
        return uploadToken;
    }

    public void setUploadToken(String uploadToken) {
        this.uploadToken = uploadToken;
    }

    public String getDownloadUrl() {
        return downloadUrl;
    }

    public void setDownloadUrl(String downloadUrl) {
        this.downloadUrl = downloadUrl;
    }
}
