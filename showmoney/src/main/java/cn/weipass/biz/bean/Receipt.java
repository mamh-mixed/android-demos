package cn.weipass.biz.bean;

import android.os.Parcel;
import android.os.Parcelable;

/**
 * 打印出来凭条的内容
 * <p/>
 * <p/>
 * <p/>
 * <p/>
 *         签购单（顾客存根）
 * --------------------------------
 * 商户名称：上海讯联测试商户
 * 商户编号：100000000010001
 * 收银员：收银员1
 * 交易类型：支付宝 扫码付
 * 日期时间：2015-11-30 14:32:50
 * 交易账号：414***@qq.com
 * 渠道订单号：2015113021001004970243401199
 * 商家订单号：20151130143248075112509
 * 金额：RMB  0.01
 * <p/>
 * <p/>
 * 退款：
 *        签购单（顾客存根）
 * --------------------------------
 * 商户名称：上海讯联测试商户
 * 商户编号：100000000010001
 * 收银员：收银员1
 * 交易类型：支付宝 退款
 * 日期时间：2015-11-30 14:32:50
 * 交易账号：414***@qq.com
 * 渠道订单号：2015113021001004970243401199
 * 商家订单号：20151130143248075112509
 * 原交易订单号：20151130143248075112509
 * 金额：RMB  -0.01
 * <p/>
 * 核销：
 *        签购单（顾客存根）
 * --------------------------------
 * 商户名称：上海讯联测试商户
 * 商户编号：100000000010001
 * 收银员：收银员1
 * 交易类型：卡券核销
 * 日期时间：2015-11-30 14:32:50
 * 商家订单号：20151130143248075112509
 * 卡券号：30143248075112
 * 详情：
 * 冰火名家烤全鱼一份
 */
public class Receipt extends BaseBean {


    private String title; //凭条抬头
    private String merchantName; //商户名称
    private String merchantNum; //商户编号
    private String cashier; //收银员
    private String tradeType; //交易类型
    private String tradeDate; //日期时间
    private String orderChannel; //渠道订单号(3.0 目前没有返回)
    private String orderMerchant; //商家订单号
    private String orderOriginal; //原交易订单号(退款凭条所需字段)
    private String coupon; //卡券号(核销凭条所需字段)
    private String couponDes; //详情(核销凭条所需字段)
    private String amount; //金额

    private Receipt(Parcel source) {

        title = source.readString();
        merchantName = source.readString();
        merchantNum = source.readString();
        cashier = source.readString();
        tradeType = source.readString();
        tradeDate = source.readString();
        orderChannel = source.readString();
        orderMerchant = source.readString();
        orderOriginal = source.readString();
        coupon = source.readString();
        couponDes = source.readString();
        amount = source.readString();

    }

    public static final Parcelable.Creator<Receipt> CREATOR = new Parcelable.Creator<Receipt>() {

        @Override
        public Receipt createFromParcel(Parcel source) {
            return new Receipt(source);
        }

        @Override
        public Receipt[] newArray(int size) {
            return new Receipt[size];
        }
    };


    @Override
    public int describeContents() {
        return 0;
    }

    @Override
    public void writeToParcel(Parcel dest, int flags) {

        dest.writeString(title);
        dest.writeString(merchantName);
        dest.writeString(merchantNum);
        dest.writeString(cashier);
        dest.writeString(tradeType);
        dest.writeString(tradeDate);
        dest.writeString(orderChannel);
        dest.writeString(orderMerchant);
        dest.writeString(orderOriginal);
        dest.writeString(coupon);
        dest.writeString(couponDes);
        dest.writeString(amount);
    }

}
