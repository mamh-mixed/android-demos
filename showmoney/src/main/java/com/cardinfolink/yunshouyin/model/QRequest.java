package com.cardinfolink.yunshouyin.model;

public class QRequest {

    private String busicd;

    private String inscd;

    private String txndir;

    private String terminalid;

    private String origOrderNum;

    private String orderNum;

    private String mchntid;

    /**
     * 交易来源，例如android，pc，ios等
     */
    private String tradeFrom;

    /**
     * 交易总金额
     */
    private String txamt;

    /**
     * 总金额
     */
    private String totalFee;

    /**
     * 交易渠道
     */
    private String chcd;

    /**
     * 商品信息，现在估计用不到了
     */
    private String goodsInfo;

    /**
     * 交易币种类型
     */
    private String currency;

    /**
     * 卡券优惠金额
     */
    private String couponDiscountAmt;

    public String getBusicd() {
        return busicd;
    }

    public void setBusicd(String busicd) {
        this.busicd = busicd;
    }

    public String getInscd() {
        return inscd;
    }

    public void setInscd(String inscd) {
        this.inscd = inscd;
    }

    public String getTxndir() {
        return txndir;
    }

    public void setTxndir(String txndir) {
        this.txndir = txndir;
    }

    public String getTerminalid() {
        return terminalid;
    }

    public void setTerminalid(String terminalid) {
        this.terminalid = terminalid;
    }

    public String getOrigOrderNum() {
        return origOrderNum;
    }

    public void setOrigOrderNum(String origOrderNum) {
        this.origOrderNum = origOrderNum;
    }

    public String getOrderNum() {
        return orderNum;
    }

    public void setOrderNum(String orderNum) {
        this.orderNum = orderNum;
    }

    public String getMchntid() {
        return mchntid;
    }

    public void setMchntid(String mchntid) {
        this.mchntid = mchntid;
    }

    public String getTradeFrom() {
        return tradeFrom;
    }

    public void setTradeFrom(String tradeFrom) {
        this.tradeFrom = tradeFrom;
    }

    public String getTxamt() {
        return txamt;
    }

    public void setTxamt(String txamt) {
        this.txamt = txamt;
    }

    public String getChcd() {
        return chcd;
    }

    public void setChcd(String chcd) {
        this.chcd = chcd;
    }

    public String getGoodsInfo() {
        return goodsInfo;
    }

    public void setGoodsInfo(String goodsInfo) {
        this.goodsInfo = goodsInfo;
    }

    public String getTotalFee() {
        return totalFee;
    }

    public void setTotalFee(String totalFee) {
        this.totalFee = totalFee;
    }

    public String getCurrency() {
        return currency;
    }

    public void setCurrency(String currency) {
        this.currency = currency;
    }

    public String getCouponDiscountAmt() {
        return couponDiscountAmt;
    }

    public void setCouponDiscountAmt(String couponDiscountAmt) {
        this.couponDiscountAmt = couponDiscountAmt;
    }

    @Override
    public String toString() {
        return "QRequest{" +
                "busicd='" + busicd + '\'' +
                ", inscd='" + inscd + '\'' +
                ", txndir='" + txndir + '\'' +
                ", terminalid='" + terminalid + '\'' +
                ", origOrderNum='" + origOrderNum + '\'' +
                ", orderNum='" + orderNum + '\'' +
                ", mchntid='" + mchntid + '\'' +
                ", tradeFrom='" + tradeFrom + '\'' +
                ", txamt='" + txamt + '\'' +
                ", totalFee='" + totalFee + '\'' +
                ", chcd='" + chcd + '\'' +
                ", goodsInfo='" + goodsInfo + '\'' +
                ", currency='" + currency + '\'' +
                ", couponDiscountAmt='" + couponDiscountAmt + '\'' +
                '}';
    }
}
