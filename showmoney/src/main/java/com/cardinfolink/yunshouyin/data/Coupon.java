package com.cardinfolink.yunshouyin.data;

/**
 * Created by charles on 2016/1/4.
 * 保存卡券优惠相关的信息， 这个放在SessonData里面，作为一个static的成员变量
 */
public class Coupon {
    private String respcd;//返回交易结果
    private String busicd;//交易类型
    private String chcd;//交易渠道
    private String payType;//支付方式
    private String availCount;//卡券剩余次数
    private String expDate;//卡券有效期
    private String voucherType;//卡券的类型
    private String saleMinAmount;//满足优惠条件的最小金额
    private String saleDiscount;//抵扣金额
    private String actualPayAmount;//实际支付金额
    private String maxDiscountAmt;//最大优惠金额
    private String mchntid;//商户号
    private String scanCodeId;//扫码号
    private String cardId;//卡券类型
    private String cardInfo;//卡券详情
    private String cardbin;//银行卡的bin 或 用户标识
    private String txamt;//交易原始金额

    public void setRespcd(String respcd) {
        this.respcd = respcd;
    }

    public String getBusicd() {
        return busicd;
    }

    public void setBusicd(String busicd) {
        this.busicd = busicd;
    }

    public String getChcd() {
        return chcd;
    }

    public void setChcd(String chcd) {
        this.chcd = chcd;
    }

    public String getPayType() {
        return payType;
    }

    public void setPayType(String payType) {
        this.payType = payType;
    }

    public String getAvailCount() {
        return availCount;
    }

    public void setAvailCount(String availCount) {
        this.availCount = availCount;
    }

    public String getExpDate() {
        return expDate;
    }

    public void setExpDate(String expDate) {
        this.expDate = expDate;
    }

    public String getVoucherType() {
        return voucherType;
    }

    public void setVoucherType(String voucherType) {
        this.voucherType = voucherType;
    }

    public String getSaleMinAmount() {
        return saleMinAmount;
    }

    public void setSaleMinAmount(String saleMinAmount) {
        this.saleMinAmount = saleMinAmount;
    }

    public String getSaleDiscount() {
        return saleDiscount;
    }

    public void setSaleDiscount(String saleDiscount) {
        this.saleDiscount = saleDiscount;
    }

    public String getActualPayAmount() {
        return actualPayAmount;
    }

    public void setActualPayAmount(String actualPayAmount) {
        this.actualPayAmount = actualPayAmount;
    }

    public String getMaxDiscountAmt() {
        return maxDiscountAmt;
    }

    public void setMaxDiscountAmt(String maxDiscountAmt) {
        this.maxDiscountAmt = maxDiscountAmt;
    }

    public String getMchntid() {
        return mchntid;
    }

    public void setMchntid(String mchntid) {
        this.mchntid = mchntid;
    }

    public String getScanCodeId() {
        return scanCodeId;
    }

    public void setScanCodeId(String scanCodeId) {
        this.scanCodeId = scanCodeId;
    }

    public String getCardId() {
        return cardId;
    }

    public void setCardId(String cardId) {
        this.cardId = cardId;
    }

    public String getCardInfo() {
        return cardInfo;
    }

    public void setCardInfo(String cardInfo) {
        this.cardInfo = cardInfo;
    }

    public String getCardbin() {
        return cardbin;
    }

    public void setCardbin(String cardbin) {
        this.cardbin = cardbin;
    }

    public String getTxamt() {
        return txamt;
    }

    public void setTxamt(String txamt) {
        this.txamt = txamt;
    }

    public String getRespcd() {
        return respcd;
    }
}
