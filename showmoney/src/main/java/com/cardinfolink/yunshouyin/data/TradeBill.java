package com.cardinfolink.yunshouyin.data;

import java.io.Serializable;

public class TradeBill implements Serializable {
    private static final long serialVersionUID = 1L;

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

    public String errorDetail;//错误信息
    public String total;//交易金额
    public String originalTotal;//消费金额
}
