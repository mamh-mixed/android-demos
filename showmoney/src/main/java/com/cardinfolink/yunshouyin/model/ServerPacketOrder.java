package com.cardinfolink.yunshouyin.model;

import com.cardinfolink.yunshouyin.api.QuickPayException;
import com.cardinfolink.yunshouyin.data.User;
import com.google.gson.Gson;
import com.google.gson.GsonBuilder;

/**
 * Created by mamh on 15-12-4.
 * 这个和ServerPacket唯一的区别就是里面的Txn是一个Txn对象，还是一个Txn数组。
 */
public class ServerPacketOrder {
    /** 刷新是返回的 serverpacketOder 。json对应的实体类。
     * 成功返回的json字符串：
     * {
     *      "state": "success",
     *      "count": 0,
     *      "size": 0,
     *      "refdcount": 0,
     *      "txn": {
     *             "response": "09",
     *             "system_date": "20151204112740",
     *             "transStatus": "10",
     *             "refundAmt": 0,
     *             "m_request": {
     *                   "busicd": "PAUT",
     *                   "inscd": "99911888",
     *                   "txndir": "Q",
     *                   "terminalid": "000000000000000",
     *                   "orderNum": "15120322232663574",
     *                   "mchntid": "999118880000017",
     *                   "tradeFrom": "android",
     *                   "txamt": "000000089500",
     *                   "chcd": "ALP",
     *                   "currency": "CNY"
     *              }
     *      }
     * }
     *
     *
     * {"state":"fail","error":"params_empty","count":0,"size":0,"refdcount":0}
     * {"state":"fail","error":"sign_fail","count":0,"size":0,"refdcount":0}
     */
    private String state;
    private String error;
    private User user;
    private int count;
    private String total;
    private int refdcount;
    private String refdtotal;
    private int size;
    private BankInfo info;
    private Txn txn;

    public static ServerPacketOrder getServerPacketOrder(String json) {
        try {
            Gson gson = new GsonBuilder().setDateFormat("yyyy-MM-dd HH:mm:ss").create();
            ServerPacketOrder packet = gson.fromJson(json, ServerPacketOrder.class);
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

    public Txn getTxn() {
        return txn;
    }

    public void setTxn(Txn txn) {
        this.txn = txn;
    }
}
