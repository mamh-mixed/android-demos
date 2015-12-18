package com.cardinfolink.cashiersdk.util;

import com.cardinfolink.cashiersdk.model.InitData;
import com.cardinfolink.cashiersdk.model.OrderData;

import org.json.JSONException;
import org.json.JSONObject;


/**
 * 2. 参数描述
 * 参数                参数名称    最大长度      备注
 * txndir              交易方向    String(1)     Q:请求，A:应答
 * busicd              交易类型    String(4)     PURC：下单支付，PAUT：预下单，INQY：查询订单，VOID：撤销，REFD：退款，VERI：卡券核销
 * respcd              交易结果    String(2)     "应答码，00:交易成功，09：处理中（需要用户在手机客户端输入支付密码）"
 * inscd               机构号      String(8)     机构号，商户所属机构标识
 * chcd                渠道机构    String(5)     ALP：支付宝，WXP：微信
 * mchntid             商户号      String(15)    商户号
 * terminalid          终端号      String（8）   终端号
 * txamt               订单金额    String(12)    12定长字符，单位为分，左补0
 * goodsInfo           商品名称    String(32)    商品详情，上送格式如下：商品名称 1,单价,数量;商品名称 n,单价,数量;各要素之间逗号分隔，每条之间分号分隔。可显示在用户手机上的支付订单上
 * channelOrderNum     渠道交易号  String(64)    支付宝／微信返回的订单号
 * consumerAccount     渠道账号    String(64)    用户账号，支付宝返回支付宝账号，微信支付返回用户OpenID
 * consumerId          渠道账号ID  String(16)    用户id，支付宝返回的用户标识，微信不返回
 * errorDetail         错误信息    String(64)    交易状态/错误详情
 * orderNum            订单号      String(64)    支付方的订单号，同一个商户下的订单号不可重复
 * origOrderNum        原订单号    String(64)    支付方的原交易的订单号
 * qrcode              二维码信息  String(128)   预下单支付宝返回的二维码url串
 * scanCodeId          扫码号      String(32)    终端扫出来的字符串（终端主拍）
 * sign                签名        String(128)
 * chcdDiscount        渠道优惠    String(13)    如123.5，对于当前订单渠道优惠掉的金额，此部分不属于用户支付
 * merDiscount         商户优惠    String(13)    如123.5，对于当前订单商户优惠掉的金额，此部分不属于用户支付
 * cardId              卡券类型    String(40)    卡券类型
 * cardInfo            卡券详情    String(256)   卡券详情
 */
public class ParamsUtil {

    /**
     * 3.1. 下单支付
     * <p/>
     * 下单支付接口适用于获取支付宝客户端的“付款码”或微信“刷卡”的条码号，
     * 并通过该接口上送此条码号（scanCodeId字段）进行支付。
     * <p/>
     * 参数        参数名称            请求    应答    类型            备注
     * 交易方向    txndir              M       M       String(1)
     * 交易类型    busicd              M       M       String(4)     PURC 注意这里的值
     * 交易结果    respcd              null    M       String(2)
     * 机构号      inscd               M       M       String(8)
     * 渠道        chcd                C       C       String(5)    成功应答中必选
     * 商户号      mchntid             M       M       String(15)
     * 终端号      terminalid          M       null    String（8）   有终端号时要求填写
     * 订单金额    txamt               M       M       String(12)
     * 商品名称    goodsInfo           C       null    String(32)    商品详情，上送格式如下：商品名称 1,单价,数量; 商品名称 n,单价,数量；各要素之间逗号分隔，每条之间分号分隔
     * 渠道交易号  channelOrderNum     null    C       String(64)
     * 渠道账号    consumerAccount     null    C       String(64)
     * 渠道账号ID  consumerId          null    C       String(16)
     * 错误信息    errorDetail         null    C       String(64)
     * 订单号      orderNum            M       C       String(64)
     * 扫码号      scanCodeId          M       null    String(32)
     * 签名        sign                M       M       String(128)
     * 渠道优惠    chcdDiscount        null    C       String(13)
     * 商户优惠    merDiscount         null    C       String(13)
     * <p/>
     * 上面表格中null字段。传递的时候不要加上这些为null的字段
     *
     * @param initData
     * @param orderData
     * @return
     */
    public static JSONObject getPay(InitData initData, OrderData orderData) {
        JSONObject json = new JSONObject();
        try {
            json.put("txndir", "Q");
            json.put("busicd", "PURC");//下单支付  ( PURC：下单支付)，(PAUT：预下单)，(INQY：查询订单)，(VOID：撤销)，(CANC 取消订单)，(REFD：退款)，(VERI：卡券核销)
            json.put("inscd", initData.inscd);
            json.put("mchntid", initData.mchntid);
            json.put("txamt", orderData.txamt);
            json.put("orderNum", orderData.orderNum);
            json.put("scanCodeId", orderData.scanCodeId);
            json.put("terminalid", initData.terminalid);
            json.put("tradeFrom", "android");
            if (orderData.goodsInfo != null) {
                json.put("goodsInfo", orderData.goodsInfo);
            }
            json.put("sign", getSign(MapUtil.getSignString(MapUtil.getMapForJson(json.toString())), initData.signKey, "SHA-1"));
        } catch (JSONException e) {
            e.printStackTrace();
        }
        return json;
    }

    /**
     * 3.2. 预下单
     * 预下单接口为一笔指定金额的交易生成一个url（qrcode字段），将此url直接转换成二维码，
     * 使用支付宝或微信的扫一扫功能即可在手机端完成支付。可通过查询订单接口确定此交易的交易状态。
     * 参数        参数名称            请求    应答    类型             备注
     * 交易方向     txndir              M       M       String(1)
     * 交易类型     busicd              M       M       String(4)     PAUT 注意这里的值
     * 交易结果     respcd              null    M       String(2)
     * 机构号       inscd               M       M       String(8)
     * 渠道         chcd                M       M       String(5)
     * 商户号       mchntid             M       M       String(15)
     * 终端号       terminalid          M       null    String（8）   有终端号时要求填写
     * 订单金额     txamt               M       M       String(12)
     * 商品名称     goodsInfo           C       null    String(32)    商品详情，上送格式如下：商品名称 1,单价,数量; 商品名称 n,单价,数量；各要素之间逗号分隔，每条之间分号分隔
     * 渠道交易号   channelOrderNum     null    C       String(64)
     * 错误信息     errorDetail         null    C       String(64)
     * 订单号       orderNum            M       M       String(64)
     * 二维码信息   qrcode              null    C       String(128)
     * 签名        sign                M       M       String(128)
     *
     * @param initData
     * @param orderData
     * @return
     */
    public static JSONObject getPrePay(InitData initData, OrderData orderData) {
        JSONObject json = new JSONObject();
        try {
            json.put("txndir", "Q");
            json.put("busicd", "PAUT");//( PURC：下单支付)，(PAUT：预下单)，(INQY：查询订单)，(VOID：撤销)，(CANC 取消订单)，(REFD：退款)，(VERI：卡券核销)
            json.put("inscd", initData.inscd);
            json.put("mchntid", initData.mchntid);
            json.put("txamt", orderData.txamt);
            json.put("orderNum", orderData.orderNum);
            json.put("chcd", orderData.chcd);
            json.put("terminalid", initData.terminalid);
            json.put("tradeFrom", "android");
            if (orderData.goodsInfo != null) {
                json.put("goodsInfo", orderData.goodsInfo);
            }
            json.put("sign", getSign(MapUtil.getSignString(MapUtil.getMapForJson(json.toString())), initData.signKey, "SHA-1"));
        } catch (JSONException e) {
            e.printStackTrace();
        }
        return json;
    }


    /**
     * 3.3. 查询订单
     * 此接口用于查询下单支付和预下单交易的交易状态，当下单支付或预下单
     * 支付返回09：交易处理中或98：交易超时的应答码或者接入方没有收到应答时（网络原因等），推荐需要调用此接口，以明确订单状态。
     * 参数        参数名称            请求    应答    类型             备注
     * 交易方向    txndir              M       M       String(1)
     * 交易类型    busicd              M       M       String(4)      INQY 注意这里的值
     * 交易结果    respcd              null    M       String(2)
     * 机构号      inscd               M       M       String(8)
     * 渠道        chcd                C       C       String(5)      成功应答中必选
     * 商户号      mchntid             M       M       String(15)
     * 终端号      terminalid          M       null    String（8）     有终端号时要求填写
     * 渠道交易号  channelOrderNum     null    C       String(64)
     * 渠道账号    consumerAccount     null    C       String(64)
     * 渠道账号ID  consumerId          null    C       String(16)
     * 错误信息    errorDetail         null    C       String(64)
     * 原订单号    origOrderNum        M       M       String(64)
     * 签名        sign                M       M       String(128)
     * 渠道优惠    chcdDiscount        null    C       String(13)
     * 商户优惠    merDiscount         null    C       String(13)
     */
    public static JSONObject getQy(InitData initData, OrderData orderData) {
        JSONObject json = new JSONObject();
        try {
            json.put("txndir", "Q");
            json.put("busicd", "INQY");//( PURC：下单支付)，(PAUT：预下单)，(INQY：查询订单)，(VOID：撤销)，(CANC 取消订单)，(REFD：退款)，(VERI：卡券核销)
            json.put("inscd", initData.inscd);
            json.put("mchntid", initData.mchntid);
            json.put("origOrderNum", orderData.origOrderNum);
            json.put("terminalid", initData.terminalid);
            json.put("sign", getSign(MapUtil.getSignString(MapUtil.getMapForJson(json.toString())), initData.signKey, "SHA-1"));
        } catch (JSONException e) {
            e.printStackTrace();
        }
        return json;
    }


    /**
     * 3.4. 撤销
     * 撤销下单支付或预下单交易。
     * 撤销下单支付或预下单交易。
     * 参数        参数名称            请求    应答    类型           备注
     * 交易方向    txndir              M       M       String(1)
     * 交易类型    busicd              M       M       String(4)    VOID 注意这里的值
     * 交易结果    respcd              null    M       String(2)
     * 机构号      inscd               M       M       String(8)
     * 渠道        chcd                C       C       String(5)    成功应答中必选
     * 商户号      mchntid             M       M       String(15)
     * 终端号      terminalid          M       null    String（8）   有终端号时要求填写
     * 渠道交易号  channelOrderNum     null    C       String(64)
     * 渠道账号    consumerAccount     null    C       String(64)
     * 渠道账号ID  consumerId          null    C       String(16)
     * 错误信息    errorDetail         null    C       String(64)
     * 订单号      orderNum            M       M       String(64)
     * 原订单号    origOrderNum        M       M       String(64)
     * 签名        sign                M       M       String(128)
     * 渠道优惠    chcdDiscount        null    C       String(13)
     * 商户优惠    merDiscount         null    C       String(13)
     *
     * @param initData
     * @param orderData
     * @return
     */
    public static JSONObject getVoid(InitData initData, OrderData orderData) {
        JSONObject json = new JSONObject();
        try {
            json.put("txndir", "Q");
            json.put("busicd", "VOID");//( PURC：下单支付)，(PAUT：预下单)，(INQY：查询订单)，(VOID：撤销)，(CANC 取消订单)，(REFD：退款)，(VERI：卡券核销)
            json.put("inscd", initData.inscd);
            json.put("mchntid", initData.mchntid);
            json.put("origOrderNum", orderData.origOrderNum);
            json.put("orderNum", orderData.orderNum);
            json.put("terminalid", initData.terminalid);
            json.put("tradeFrom", "android");
            json.put("sign", getSign(MapUtil.getSignString(MapUtil.getMapForJson(json.toString())), initData.signKey, "SHA-1"));
        } catch (JSONException e) {
            e.printStackTrace();
        }
        return json;
    }

    /**
     * 3.5. 退款
     * 撤销下单支付或预下单交易。
     * 参数        参数名称            请求    应答    类型        备注
     * 交易方向    txndir              M       M       String(1)
     * 交易类型    busicd              M       M       String(4)    REFD 注意这里的值
     * 交易结果    respcd              null    M       String(2)
     * 机构号      inscd               M       M       String(8)
     * 渠道        chcd                C       C       String(5)    成功应答中必选
     * 商户号      mchntid             M       M       String(15)
     * 终端号      terminalid          M       null    String（8）    有终端号时要求填写
     * 订单金额    txamt               M       null    String(12)
     * 渠道交易号  channelOrderNum     null    C       String(64)
     * 渠道账号    consumerAccount     null    C       String(64)
     * 渠道账号ID  consumerId          null    C       String(16)
     * 错误信息    errorDetail         null    C       String(64)
     * 订单号      orderNum            M       M       String(64)
     * 原订单号    origOrderNum        M       M       String(64)
     * 签名        sign                M       M       String(128)
     * 渠道优惠    chcdDiscount        null    C       String(13)
     * 商户优惠    merDiscount         null    C       String(13)
     */
    public static JSONObject getRefd(InitData initData, OrderData orderData) {
        JSONObject json = new JSONObject();
        try {
            json.put("txndir", "Q");
            json.put("busicd", "REFD");//( PURC：下单支付)，(PAUT：预下单)，(INQY：查询订单)，(VOID：撤销)，(CANC 取消订单)，(REFD：退款)，(VERI：卡券核销)
            json.put("inscd", initData.inscd);
            json.put("mchntid", initData.mchntid);
            json.put("origOrderNum", orderData.origOrderNum);
            json.put("orderNum", orderData.orderNum);
            json.put("txamt", orderData.txamt);
            json.put("terminalid", initData.terminalid);
            json.put("tradeFrom", "android");
            json.put("sign", getSign(MapUtil.getSignString(MapUtil.getMapForJson(json.toString())), initData.signKey, "SHA-1"));
        } catch (JSONException e) {
            e.printStackTrace();
        }
        return json;
    }

    /**
     * 3.6  取消订单
     * 对于未成功付款的订单进行取消，则关闭交易，使用户后期不能支付成功；
     * 对于成功付款的订单进行取消，系统将订单金额返还给用户，相当于对此交易做撤销。
     * 参数        参数名称            请求    应答    类型            备注
     * 交易方向    txndir              M       M       String(1)    Q:请求，A:应答
     * 交易类型    busicd              M       M       String(4)    CANC
     * 交易结果    respcd              null    M       String(2)
     * 机构号      inscd               M       M       String(8)
     * 渠道        chcd                C       C       String(5)    成功应答中必选
     * 商户号      mchntid             M       M       String(15)
     * 终端号      terminalid          M       null    String（8）    有终端号时要求填写
     * 渠道交易号  channelOrderNum     null    C       String(64)
     * 渠道账号    consumerAccount     null    C       String(64)
     * 渠道账号ID  consumerId          null    C       String(16)
     * 错误信息    errorDetail         null    C       String(64)
     * 订单号      orderNum            M       M       String(64)
     * 原订单号    origOrderNum        M       M       String(64)
     * 签名        sign                M       M       String(128)
     * 渠道优惠    chcdDiscount        null    C       String(13)
     * 商户优惠    merDiscount         null    C       String(13)
     */
    public static JSONObject getCanc(InitData initData, OrderData orderData) {
        JSONObject json = new JSONObject();
        try {
            json.put("txndir", "Q");//Q:请求，A:应答
            json.put("busicd", "CANC");//CANC取消订单 ( PURC：下单支付)，(PAUT：预下单)，(INQY：查询订单)，(VOID：撤销)，(CANC 取消订单)，(REFD：退款)，(VERI：卡券核销)
            json.put("inscd", initData.inscd);//机构号
            json.put("mchntid", initData.mchntid);//商户号
            json.put("terminalid", initData.terminalid);//终端号
            json.put("orderNum", orderData.orderNum);//订单号
            json.put("origOrderNum", orderData.origOrderNum);//原订单号
            json.put("tradeFrom", "android");//这个要加上的
            json.put("sign", getSign(MapUtil.getSignString(MapUtil.getMapForJson(json.toString())), initData.signKey, "SHA-1"));
        } catch (Exception e) {
            e.printStackTrace();
        }
        return json;
    }

    /**
     * 3.7卡券核销
     * 对上送的卡券号进行核销。
     * 参数        参数名称        请求    应答    类型           备注
     * 交易方向    txndir          M       M       String(1)
     * 交易类型    busicd          M       M       String(4)    VERI   注意这里的值
     * 交易结果    respcd          null    M       String(2)
     * 机构号      inscd           M       M       String(8)
     * 渠道        chcd            C       C       String(5)    成功应答中必选
     * 商户号      mchntid         M       M       String(15)
     * 终端号      terminalid      M       null    String（8）   有终端号时要求填写
     * 错误信息    errorDetail     null    C       String(64)
     * 订单号      orderNum        M       M       String(64)
     * 扫码号      scanCodeId      M       M       String(32)
     * 签名        sign            M       M       String(128)
     * 卡券类型    cardId          null    C       String(40)
     * 卡券详情    cardInfo        null    C       String(256)
     *
     * @param initData
     * @param orderData
     * @return
     */
    public static JSONObject getVeri(InitData initData, OrderData orderData) {
        JSONObject json = new JSONObject();
        try {
            json.put("txndir", "Q");
            json.put("busicd", "VERI");// ( PURC：下单支付)，(PAUT：预下单)，(INQY：查询订单)，(VOID：撤销)，(CANC 取消订单)，(REFD：退款)，(VERI：卡券核销)
            json.put("inscd", initData.inscd);
            json.put("mchntid", initData.mchntid);
            json.put("scanCodeId", orderData.scanCodeId);
            json.put("orderNum", orderData.orderNum);
            json.put("terminalid", initData.terminalid);
            json.put("tradeFrom", "android");
            json.put("sign", getSign(MapUtil.getSignString(MapUtil.getMapForJson(json.toString())), initData.signKey, "SHA-1"));
        } catch (JSONException e) {
            e.printStackTrace();
        }
        return json;
    }


    @SuppressWarnings("unchecked")
    public static String getSign(String str, String key, String signType) {
        str = str + key;
        String sign = EncoderUtil.Encrypt(str, signType);
        return sign;
    }
}
