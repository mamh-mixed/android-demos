package com.cardinfolink.cashiersdk.model;

public class ResultData {
    public String respcd;//返回交易结果
    public String busicd;//交易类型
    public String chcd;//交易渠道
    public String txamt;//交易原始金额
    public String channelOrderNum;//
    public String consumerAccount;
    public String consumerId;
    public String errorDetail;//错误详情
    public String orderNum;//订单号
    public String chcdDiscount;
    public String merDiscount;
    public String qrcode;
    public String origOrderNum;//原始订单号
    public String scanCodeId;//扫码号
    public String cardId;//卡券类型
    public String cardInfo;//卡券详情
    public String cardbin;//银行卡的bin 或 用户标识
    public String mchntid;//商户号
    public String veriTime;//核销次数
    public String payType;//支付方式
    public String availCount;//卡券剩余次数
    public String expDate;//卡券有效期
    public String voucherType;//卡券的类型
    public String saleMinAmount;//满足优惠条件的最小金额
    public String saleDiscount;//抵扣金额
    public String actualPayAmount;//实际支付金额
    public String maxDiscountAmt;//最大优惠金额


    @Override
    public String toString() {
        return "ResultData{" +
                "respcd='" + respcd + '\'' +
                ", busicd='" + busicd + '\'' +
                ", chcd='" + chcd + '\'' +
                ", txamt='" + txamt + '\'' +
                ", channelOrderNum='" + channelOrderNum + '\'' +
                ", consumerAccount='" + consumerAccount + '\'' +
                ", consumerId='" + consumerId + '\'' +
                ", errorDetail='" + errorDetail + '\'' +
                ", orderNum='" + orderNum + '\'' +
                ", chcdDiscount='" + chcdDiscount + '\'' +
                ", merDiscount='" + merDiscount + '\'' +
                ", qrcode='" + qrcode + '\'' +
                ", origOrderNum='" + origOrderNum + '\'' +
                ", scanCodeId='" + scanCodeId + '\'' +
                ", cardId='" + cardId + '\'' +
                ", cardInfo='" + cardInfo + '\'' +
                ", cardbin='" + cardbin + '\'' +
                ", mchntid='" + mchntid + '\'' +
                ", veriTime='" + veriTime + '\'' +
                ", payType='" + payType + '\'' +
                ", availCount='" + availCount + '\'' +
                ", expDate='" + expDate + '\'' +
                ", voucherType='" + voucherType + '\'' +
                ", saleMinAmount='" + saleMinAmount + '\'' +
                ", saleDiscount='" + saleDiscount + '\'' +
                ", actualPayAmount='" + actualPayAmount + '\'' +
                ", maxDiscountAmt='" + maxDiscountAmt + '\'' +
                '}';
    }
}
