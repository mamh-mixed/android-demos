package com.cardinfolink.yunshouyin.model;

import java.io.Serializable;

/**
 * Created by Tommy on 2015/12/28.
 */
public class Message implements Serializable {

    private String msgId;
    private String username;
    private String title;
    private String message;
    private String pushtime;
    private String updateTime;
    private String status;

    public Message() {
    }

    public Message(String msgId, String username, String title, String message, String pushtime, String updateTime) {
        this.msgId = msgId;
        this.username = username;
        this.title = title;
        this.message = message;
        this.pushtime = pushtime;
        this.updateTime = updateTime;
    }

    public Message(String msgId, String username, String title, String message, String pushtime, String updateTime, String status) {
        this(msgId, username, title, message, pushtime, updateTime);
        this.status = status;
    }

    public String getStatus() {
        return status;
    }

    public void setStatus(String status) {
        this.status = status;
    }

    public String getMsgId() {
        return msgId;
    }

    public void setMsgId(String msgId) {
        this.msgId = msgId;
    }

    public String getUsername() {
        return username;
    }

    public void setUsername(String username) {
        this.username = username;
    }

    public String getTitle() {
        return title;
    }

    public void setTitle(String title) {
        this.title = title;
    }

    public String getMessage() {
        return message;
    }

    public void setMessage(String message) {
        this.message = message;
    }

    public String getPushtime() {
        return pushtime;
    }

    public void setPushtime(String pushtime) {
        this.pushtime = pushtime;
    }

    public String getUpdateTime() {
        return updateTime;
    }

    public void setUpdateTime(String updateTime) {
        this.updateTime = updateTime;
    }
}
