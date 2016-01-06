package com.cardinfolink.yunshouyin.model;

import com.google.gson.annotations.SerializedName;

/**
 * {
 * response='09',
 * systemDate='20160105154422',
 * transStatus='40',
 * refundAmt='0',
 * nickName='null',
 * avatarUrl='null',
 * checkCode='null',
 * couponName='null',
 * couponChannel='WXP',
 * couponOrderNo='null',
 * couponDiscountAmt='null',
 * consumerAccount='null',
 * mRequest=com.cardinfolink.yunshouyin.model.QRequest@42dfbda0
 * }
 */
public class Txn {
    private String response;

    @SerializedName("system_date")
    private String systemDate;

    private String transStatus;

    private String refundAmt;

    /**
     * 微信昵称，如果是使用微信扫描收款码支付需要有
     */
    private String nickName;

    /**
     * 微信头像，如果是使用微信扫描收款码支付需要有
     */
    private String avatarUrl;

    /**
     * 校验码，微信扫描固定码支付有这个数据
     */
    private String checkCode;

    /**
     * 卡券名称
     */
    private String couponName;


    /**
     * 卡券渠道
     */
    private String couponChannel;

    /**
     * 卡券核销订单号
     */
    private String couponOrderNo;

    /**
     * 卡券优惠金额
     */
    private String couponDiscountAmt;

    /**
     * 这个暂时没有了
     */
    private String consumerAccount;

    @SerializedName("m_request")
    private QRequest mRequest;


    public String getResponse() {
        return response;
    }

    public void setResponse(String response) {
        this.response = response;
    }

    public String getSystemDate() {
        return systemDate;
    }

    public void setSystemDate(String systemDate) {
        this.systemDate = systemDate;
    }

    public String getConsumerAccount() {
        return consumerAccount;
    }

    public void setConsumerAccount(String consumerAccount) {
        this.consumerAccount = consumerAccount;
    }

    public QRequest getmRequest() {
        return mRequest;
    }

    public void setmRequest(QRequest mRequest) {
        this.mRequest = mRequest;
    }

    public String getRefundAmt() {
        return refundAmt;
    }

    public void setRefundAmt(String refundAmt) {
        this.refundAmt = refundAmt;
    }

    public String getTransStatus() {
        return transStatus;
    }

    public void setTransStatus(String transStatus) {
        this.transStatus = transStatus;
    }

    public String getNickName() {
        return nickName;
    }

    public void setNickName(String nickName) {
        this.nickName = nickName;
    }

    public String getAvatarUrl() {
        return avatarUrl;
    }

    public void setAvatarUrl(String avatarUrl) {
        this.avatarUrl = avatarUrl;
    }

    public String getCheckCode() {
        return checkCode;
    }

    public void setCheckCode(String checkCode) {
        this.checkCode = checkCode;
    }

    public String getCouponName() {
        return couponName;
    }

    public void setCouponName(String couponName) {
        this.couponName = couponName;
    }

    public String getCouponChannel() {
        return couponChannel;
    }

    public void setCouponChannel(String couponChannel) {
        this.couponChannel = couponChannel;
    }

    public String getCouponOrderNo() {
        return couponOrderNo;
    }

    public void setCouponOrderNo(String couponOrderNo) {
        this.couponOrderNo = couponOrderNo;
    }

    public String getCouponDiscountAmt() {
        return couponDiscountAmt;
    }

    public void setCouponDiscountAmt(String couponDiscountAmt) {
        this.couponDiscountAmt = couponDiscountAmt;
    }

    @Override
    public String toString() {
        return "Txn{" +
                "response='" + response + '\'' +
                ", systemDate='" + systemDate + '\'' +
                ", transStatus='" + transStatus + '\'' +
                ", refundAmt='" + refundAmt + '\'' +
                ", nickName='" + nickName + '\'' +
                ", avatarUrl='" + avatarUrl + '\'' +
                ", checkCode='" + checkCode + '\'' +
                ", couponName='" + couponName + '\'' +
                ", couponChannel='" + couponChannel + '\'' +
                ", couponOrderNo='" + couponOrderNo + '\'' +
                ", couponDiscountAmt='" + couponDiscountAmt + '\'' +
                ", consumerAccount='" + consumerAccount + '\'' +
                ", mRequest=" + mRequest +
                '}';
    }
}
