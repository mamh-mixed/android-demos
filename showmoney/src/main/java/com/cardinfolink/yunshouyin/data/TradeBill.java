package com.cardinfolink.yunshouyin.data;

import java.io.Serializable;

public class TradeBill implements Serializable {
    private static final long serialVersionUID = 1L;

    public static final String BILL_TYPE = "bill";

    public static final String COUPON_TYPE = "coupon";

    /**
     * 默认的构造方法是new出来的是 普通的账单
     */
    public TradeBill() {
        this.billType = BILL_TYPE;
    }

    public TradeBill(String billType) {
        this.billType = billType;
    }

    /**
     * 这个字段用来区分是普通的账单 还是 卡券账单
     */
    public String billType;

    public String chcd;//交易渠道
    public String tradeFrom;//交易来自哪里
    public String busicd;
    public String orderNum;//订单号
    public String tandeDate;//订单日期
    public String response;//订单状态
    public String amount;
    public String consumerAccount;
    public String goodsInfo;
    public String refundAmt;
    public String transStatus;//这个是新增的一个表示订单状态的字段

    public String couponType;
    public String couponName;
    public String couponChannel;//卡券渠道
    public String couponDiscountAmt;//卡券优惠金额


    public String errorDetail;//错误信息
    public String total;//交易金额
    public String originalTotal;//消费金额


    @Override
    public String toString() {
        return "TradeBill{" +
                "billType='" + billType + '\'' +
                ", chcd='" + chcd + '\'' +
                ", tradeFrom='" + tradeFrom + '\'' +
                ", busicd='" + busicd + '\'' +
                ", orderNum='" + orderNum + '\'' +
                ", tandeDate='" + tandeDate + '\'' +
                ", response='" + response + '\'' +
                ", amount='" + amount + '\'' +
                ", consumerAccount='" + consumerAccount + '\'' +
                ", goodsInfo='" + goodsInfo + '\'' +
                ", refundAmt='" + refundAmt + '\'' +
                ", transStatus='" + transStatus + '\'' +
                ", couponType='" + couponType + '\'' +
                ", couponName='" + couponName + '\'' +
                ", couponChannel='" + couponChannel + '\'' +
                ", couponDiscountAmt='" + couponDiscountAmt + '\'' +
                ", errorDetail='" + errorDetail + '\'' +
                ", total='" + total + '\'' +
                ", originalTotal='" + originalTotal + '\'' +
                '}';
    }
}
