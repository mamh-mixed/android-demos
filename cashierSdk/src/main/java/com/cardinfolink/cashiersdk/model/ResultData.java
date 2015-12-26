package com.cardinfolink.cashiersdk.model;

public class ResultData {
    public String respcd;
    public String busicd;
    public String chcd;
    public String txamt;
    public String channelOrderNum;
    public String consumerAccount;
    public String consumerId;
    public String errorDetail;
    public String orderNum;
    public String chcdDiscount;
    public String merDiscount;
    public String qrcode;
    public String origOrderNum;
    public String scanCodeId;
    public String cardId;
    public String cardInfo;
    public String voucherType;
    public String saleMinAmount;
    public String saleDiscount;
    public String actualPayAmount;
    public String maxDiscountAmt;

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
                ", voucherType='" + voucherType + '\'' +
                ", saleMinAmount='" + saleMinAmount + '\'' +
                ", saleDiscount='" + saleDiscount + '\'' +
                ", actualPayAmount='" + actualPayAmount + '\'' +
                ", maxDiscountAmt='" + maxDiscountAmt + '\'' +
                '}';
    }
}
