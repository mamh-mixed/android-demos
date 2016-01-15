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
    public String terminalid;
    public String refundAmt;
    public String transStatus;//这个是新增的一个表示订单状态的字段

    public String couponType;//卡券类型
    public String couponName;//卡券描述，名称
    public String couponChannel;//卡券渠道
    public String couponDiscountAmt;//卡券优惠金额
    public String couponOrderNum;//卡券核销的订单号


    public String errorDetail;//错误信息
    public String total;//交易金额
    public String originalTotal;//消费金额

    /**
     * 微信昵称名字
     */
    public String nickName;

    /**
     * 微信头像地址
     */
    public String avatarUrl;

    /**
     * 检验码
     */
    public String checkCode;

    /**
     * 小票号
     */
    public String smallTicketNumber;

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
                ", terminalid='" + terminalid + '\'' +
                ", refundAmt='" + refundAmt + '\'' +
                ", transStatus='" + transStatus + '\'' +
                ", couponType='" + couponType + '\'' +
                ", couponName='" + couponName + '\'' +
                ", couponChannel='" + couponChannel + '\'' +
                ", couponDiscountAmt='" + couponDiscountAmt + '\'' +
                ", couponOrderNum='" + couponOrderNum + '\'' +
                ", errorDetail='" + errorDetail + '\'' +
                ", total='" + total + '\'' +
                ", originalTotal='" + originalTotal + '\'' +
                ", nickName='" + nickName + '\'' +
                ", avatarUrl='" + avatarUrl + '\'' +
                ", checkCode='" + checkCode + '\'' +
                ", smallTicketNumber='" + smallTicketNumber + '\'' +
                '}';
    }
}
