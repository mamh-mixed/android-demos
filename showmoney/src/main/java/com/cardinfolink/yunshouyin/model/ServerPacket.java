package com.cardinfolink.yunshouyin.model;

import com.cardinfolink.yunshouyin.api.QuickPayException;
import com.cardinfolink.yunshouyin.data.User;
import com.google.gson.Gson;
import com.google.gson.GsonBuilder;
import com.google.gson.annotations.SerializedName;

public class ServerPacket {

    /**
     * succes 或者 fail
     */
    private String state;

    /**
     * 如果state是fail，一定返回，建议统一使用英文返回码，
     * 多语言化在客户端完成，接口返回错误也需要放入文档
     */
    private String error;

    /**
     * 该月支付金额，扣去了退款的金额。如果传入month字段，则必须返回。
     * 根据币种不同单位不一样。如果币种是CNY，则212表示2.12元，单位是分；如果是JPY，则212表示212元，单位是元。
     * json里面会对应为totalFee
     */
    @SerializedName("totalFee")
    private String total;

    /**
     * 该月支付笔数。如果传入month字段，则必须返回。
     */
    private int count;


    /**
     * 该月退款金额。如果传入month字段，则必须返回。根据币种不同单位不一样。
     * 如果币种是CNY，则212表示2.12元，单位是分；如果是JPY，则212表示212元，单位是元。
     */
    @SerializedName("refdTotalFee")
    private String refdtotal;

    /**
     * 该月退款笔数。如果传入month字段，则必须返回。
     */
    private int refdcount;


    /**
     * 其实是txn数组的长度
     */
    private int size;

    /**
     * 总纪录数,当月的总笔数
     */
    private int totalRecord;

    private BankInfo info;

    /**
     * 订单数组，成功返回
     */
    private Txn[] txn;

    private User user;

    /**
     * 七牛上传token，成功返回
     */
    private String uploadToken;

    private Message[] message;


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

    public String getTotal() {
        return total;
    }

    public void setTotal(String total) {
        this.total = total;
    }

    public int getRefdcount() {
        return refdcount;
    }

    public void setRefdcount(int refdcount) {
        this.refdcount = refdcount;
    }

    public String getRefdtotal() {
        return refdtotal;
    }

    public void setRefdtotal(String refdtotal) {
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

    public Txn[] getTxn() {
        return txn;
    }

    public void setTxn(Txn[] txn) {
        this.txn = txn;
    }


    public String getUploadToken() {
        return uploadToken;
    }

    public void setUploadToken(String uploadToken) {
        this.uploadToken = uploadToken;
    }

    public int getTotalRecord() {
        return totalRecord;
    }

    public void setTotalRecord(int totalRecord) {
        this.totalRecord = totalRecord;
    }

    public Message[] getMessage() {
        return message;
    }

    public void setMessage(Message[] message) {
        this.message = message;
    }
}
