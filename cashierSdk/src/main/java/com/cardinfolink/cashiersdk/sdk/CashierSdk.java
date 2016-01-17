package com.cardinfolink.cashiersdk.sdk;

import android.util.Log;

import com.cardinfolink.cashiersdk.listener.CashierListener;
import com.cardinfolink.cashiersdk.listener.CommunicationListener;
import com.cardinfolink.cashiersdk.model.InitData;
import com.cardinfolink.cashiersdk.model.OrderData;
import com.cardinfolink.cashiersdk.model.ResultData;
import com.cardinfolink.cashiersdk.model.Server;
import com.cardinfolink.cashiersdk.util.CommunicationUtil;
import com.cardinfolink.cashiersdk.util.MapUtil;
import com.cardinfolink.cashiersdk.util.ParamsUtil;
import com.cardinfolink.cashiersdk.util.TxamtUtil;

import java.util.Map;

public class CashierSdk {
    private static final String TAG = "CashierSdk";

    private static final String mProduceHost = "121.40.167.112";
    private static final String mProducePort = "6001";
    private static final String mTestHost = "121.40.86.222";
    private static final String mTestPort = "6001";
    private static InitData mInitData;

    public static void init(InitData data) {
        mInitData = data;
        if (mInitData.isProduce) {
            Server server = new Server();
            server.setHost(mProduceHost);
            server.setPort(mProducePort);
            CommunicationUtil.setServer(server);
        } else {
            Server server = new Server();
            server.setHost(mTestHost);
            server.setPort(mTestPort);
            CommunicationUtil.setDEBUG(true);
            CommunicationUtil.setServer(server);
        }

    }

    /**
     * 成功	00	成功
     * 交易失败	01	交易失败
     * 商户号错误、未找到商户号	03	商户错误
     * 不支持该交易类型	05	不支持该交易类型
     * 处理中	09	处理中
     * 签名错误	12	签名错误
     * 条码错误或过期	14	条码错误或过期
     * 没有合适渠道（非支付宝和微信）	15	无此渠道
     * 订单号重复	19	订单号重复
     * 订单不存在	25	订单不存在
     * 没有配置路由策略	31	权限不足
     * 余额不足	51	余额不足
     * 退款失败，商家账户余额不足	61	退款失败，商家账户余额不足
     * 退款金额超过原订单金额	64	退款金额超过原订单金额
     * 外部系统错误	91	外部系统错误
     * 内部系统错误、连接支付宝、微信网关错误	96	系统错误
     * 超时、渠道规定时间内未响应	98	交易超时
     * 没有对应渠道商户	H1	商户无此交易渠道权限
     * 无此接口权限	H2	商户无此接口权限
     * 订单已关闭或取消	H3	订单已关闭或取消
     * 字段必填非空	H4	[字段名]不能为空
     * 字段格式错误	H5	[字段名]格式错误
     * 交易状态不合法	H6	交易状态不合法
     * 交易信息中包含违禁词汇	H7	交易信息中包含违禁词汇
     * 字段内容错误	H8	[字段名]字段错误
     * 原交易非支付交易	R1	原交易非支付交易
     * 原交易未成功支付	R2	原交易未成功支付
     * 存在部分退款	R3	原交易已退款
     * 原交易已退款	R3	原交易已退款
     * 只能隔天退款	R4	只能隔天退款
     * 不支持的卡类型，请换卡或绑新卡	R5	不支持的卡类型，请换卡或绑新卡
     * 只能撤销当天交易	R6	只能撤销当天交易
     * 买家账户状态异常	R7	买家账户状态异常
     * 卖家账户状态异常	R8	卖家账户状态异常
     * 付款码与交易渠道不符	R9	付款码与交易渠道不符
     * 分账信息不正确	S1	分账信息不正确
     * 软件版本过低，请升级	S2	软件版本过低，请升级
     * 付款金额小于最低限额	S3	付款金额小于最低限额
     * Openid错误	S4	Openid错误
     * 姓名校验错误	S5	姓名校验错误
     * 授权码code错误	J1	授权码code错误
     */

    /**
     * 成功	00	成功
     */
    public static final String SDK_RESPONSE_CODE_SUCCESS = "00";

    /**
     * 交易失败	01	交易失败
     */
    public static final String SDK_RESPONSE_CODE_FAIL = "01";


    /**
     * 商户号错误、未找到商户号	03	商户错误
     */
    public static final String SDK_RESPONSE_CODE_MERCHANT_ERROR = "03";


    /**
     * 不支持该交易类型	05	不支持该交易类型
     */
    public static final String SDK_RESPONSE_CODE_NOT_SUPPORT_TRADE_TYPE = "05";


    /**
     * 处理中	09	处理中
     */
    public static final String SDK_RESPONSE_CODE_IN_PROGRESS = "09";

    /**
     * 名错误	12	签名错误
     */
    public static final String SDK_RESPONSE_CODE_SIGN_ERROR = "12";


    /**
     * 条码错误或过期	14	条码错误或过期
     */
    public static final String SDK_RESPONSE_CODE_QRCODE_ERROR = "14";

    /**
     * 没有合适渠道（非支付宝和微信）	15	无此渠道
     */
    public static final String SDK_RESPONSE_CODE_CHCD_ERROR = "15";

    /**
     * 订单号重复	19	订单号重复duplicate
     */
    public static final String SDK_RESPONSE_CODE_ORDER_NUMBER_DUPLICATE = "19";


    /**
     * 订单不存在	25	订单不存在
     */
    public static final String SDK_RESPONSE_CODE_ORDER_NUMBER_NOT_EXIST = "25";


    /**
     * 没有配置路由策略	31	权限不足
     */
    public static final String SDK_RESPONSE_CODE_NO_PERMISSION = "31";


    /**
     * 余额不足	51	余额不足
     */
    public static final String SDK_RESPONSE_CODE_CUSTOM_MONEY_LOW = "51";

    /**
     * 退款失败，商家账户余额不足	61	退款失败，商家账户余额不足
     */
    public static final String SDK_RESPONSE_CODE_MERCHANT_MONEY_LOW = "61";

    /**
     * 退款金额超过原订单金额	64	退款金额超过原订单金额
     */
    public static final String SDK_RESPONSE_CODE_REFD_MONEY_LARGE = "64";


    /**
     * * 外部系统错误	91	外部系统错误
     */
    public static final String SDK_RESPONSE_CODE_OUT_SYSTEM_ERROR = "91";


    /**
     * * 内部系统错误、连接支付宝、微信网关错误	96	系统错误
     */
    public static final String SDK_RESPONSE_CODE_IN_SYSTEM_ERROR = "96";

    /**
     * * 超时、渠道规定时间内未响应	98	交易超时
     */
    public static final String SDK_RESPONSE_CODE_TRADE_TIME_OUT = "98";

    /**
     * * 没有对应渠道商户	H1	商户无此交易渠道权限
     */
    public static final String SDK_RESPONSE_CODE_MERCHANT_NOT_SUCH_CHCD_PERMISSION = "H1";


    /**
     * * 无此接口权限	H2	商户无此接口权限
     */
    public static final String SDK_RESPONSE_CODE_MERCHANT_NOT_SUCH_PORT_PERMISSION = "H2";

    /**
     * * 订单已关闭或取消	H3	订单已关闭或取消
     */
    public static final String SDK_RESPONSE_CODE_ORDER_CLOSED = "H3";


    /**
     * 字段必填非空	H4	[字段名]不能为空
     */
    public static final String SDK_RESPONSE_CODE_PARAMS_EMPTY = "H4";

    /**
     * * 字段格式错误	H5	[字段名]格式错误
     */
    public static final String SDK_RESPONSE_CODE_PARAMS_FORMAT_ERROR = "H5";

    /**
     * * 字段内容错误	H8	[字段名]字段错误
     */
    public static final String SDK_RESPONSE_CODE_PARAMS_ERROR = "H8";


    /**
     * * 交易状态不合法	H6	交易状态不合法illegal
     */
    public static final String SDK_RESPONSE_CODE_TRADE_STATUS_ILLEGAL = "H6";

    /**
     * * 交易信息中包含违禁词汇	H7	交易信息中包含违禁词汇
     */
    public static final String SDK_RESPONSE_CODE_TRADE_INFO_ERROR = "H7";


    /**
     * * 原交易非支付交易	R1	原交易非支付交易
     */
    public static final String SDK_RESPONSE_CODE_ORIGIN_NOT_PAY_TRADE = "R1";


    /**
     * * 原交易未成功支付	R2	原交易未成功支付
     */
    public static final String SDK_RESPONSE_CODE_ORIGIN_NOT_PAY_SUCCESS = "R2";


    /**
     * * 存在部分退款	R3	原交易已退款
     */
    public static final String SDK_RESPONSE_CODE_ORIGIN_PART_REFD = "R3";


    /**
     * * 原交易已退款	R3	原交易已退款
     */
    public static final String SDK_RESPONSE_CODE_ORIGIN_ALL_REFD = "R3";


    /**
     * * 只能隔天退款	R4	只能隔天退款
     */
    public static final String SDK_RESPONSE_CODE_YESTERDAY_CANNOT_REFD = "R4";


    /**
     * * 不支持的卡类型，请换卡或绑新卡	R5	不支持的卡类型，请换卡或绑新卡
     */
    public static final String SDK_RESPONSE_CODE_NOT_SUPPORT_SUCH_BANK_CARD = "R5";


    /**
     * * 只能撤销当天交易	R6	只能撤销当天交易
     */
    public static final String SDK_RESPONSE_CODE_ONLY_CANCEL_TODAY_TRADE = "R6";


    /**
     * * 买家账户状态异常	R7	买家账户状态异常
     */
    public static final String SDK_RESPONSE_CODE_CUSTOM_ACCOUNT_ERROR = "R7";


    /**
     * * 卖家账户状态异常	R8	卖家账户状态异常
     */
    public static final String SDK_RESPONSE_CODE_MERCHANT_ACCOUNT_ERROR = "R8";

    /**
     * * 付款码与交易渠道不符	R9	付款码与交易渠道不符
     */
    public static final String SDK_RESPONSE_CODE_CHCD_NOT_MATCH = "R9";


    /**
     * * 分账信息不正确	S1	分账信息不正确Ledger information
     */
    public static final String SDK_RESPONSE_CODE_LEDGER_INFO_ERROR = "S1";


    /**
     * * 软件版本过低，请升级	S2	软件版本过低，请升级
     */
    public static final String SDK_RESPONSE_CODE_VERSION_LOW = "S2";


    /**
     * * 付款金额小于最低限额	S3	付款金额小于最低限额
     */
    public static final String SDK_RESPONSE_CODE_PAY_MONEY_LESS_THAN_LIMIT = "S3";


    /**
     * * Openid错误	S4	Openid错误
     */
    public static final String SDK_RESPONSE_CODE_OPEN_ID_ERROR = "S4";

    /**
     * * 姓名校验错误	S5	姓名校验错误
     */
    public static final String SDK_RESPONSE_CODE_NAME_CHECK_ERROR = "S5";

    /**
     * * 授权码code错误	J1	授权码code错误Authorization code
     */
    public static final String SDK_RESPONSE_CODE_AUTHOR_CODE_CHECK_ERROR = "J1";


    public static final String CHARSET_GBK = "gbk";
    public static final String CHARSET = CHARSET_GBK;

    public static final String SIGN_TYPE_SHA_1 = "SHA-1";
    public static final String SIGN_TYPE_SHA_256 = "SHA-256";
    public static final String SIGN_TYPE = SIGN_TYPE_SHA_1;

    public static final String SDK_TRADE_FROM = "android";

    public static final String SDK_CURRENCY_RMB = "156";//币种类型，这里表示的人民币
    public static final String SDK_CURRENCY = SDK_CURRENCY_RMB;//币种类型，这里表示的人民币

    public static final int SDK_ERROR_RESULT_FORMAT = 0;//socket 返回的结果格式不对
    public static final int SDK_ERROR_RESULT_NULL = 5;//socket 返回的结果是空

    public static final int SDK_ERROR_TXAMT_NULL = 1;//交易金额为空
    public static final int SDK_ERROR_CURRENCY_NOT_RMB = 2;//币种不是rmb
    public static final int SDK_ERROR_SIGN_NOT_MATCH = 3;//签名不匹配

    /**
     * 3.1. 下单支付
     * 下单支付接口适用于获取支付宝客户端的“付款码”或微信“刷卡”的条码号，
     * 并通过该接口上送此条码号（scanCodeId字段）进行支付。
     *
     * @param orderData
     * @param listener
     */
    public static void startPay(OrderData orderData, final CashierListener listener) {

        String str = orderData.txamt;
        String discount = orderData.discountMoney;
        if (str != null && orderData.currency != null) {
            if (SDK_CURRENCY.equals(orderData.currency)) {
                orderData.txamt = TxamtUtil.getTxamtUtil(str);
                orderData.discountMoney = TxamtUtil.getTxamtUtil(discount);
                if (orderData.txamt == null) {
                    listener.onError(SDK_ERROR_TXAMT_NULL);
                    return;
                }
            } else {
                listener.onError(SDK_ERROR_CURRENCY_NOT_RMB);
                return;
            }
        }


        CommunicationUtil.sendDataToServer(ParamsUtil.getPay(mInitData, orderData), new CommunicationListener() {

            @Override
            public void onResult(String result) {
                Map<String, Object> map = MapUtil.getMapForJson(result);
                String sign = (String) map.get("sign");
                if (sign != null) {
                    map.remove("sign");
                    String veriSign = ParamsUtil.getSign(MapUtil.getSignString(map), mInitData.signKey);
                    if (sign.equals(veriSign)) {
                        ResultData resultData = MapUtil.getResultData(map);
                        listener.onResult(resultData);
                    } else {
                        listener.onError(SDK_ERROR_SIGN_NOT_MATCH);
                    }
                }
            }

            @Override
            public void onError(int error) {
                listener.onError(error);
            }
        });
    }


    /**
     * 3.2. 预下单
     * 预下单接口为一笔指定金额的交易生成一个url（qrcode字段），将此url直接转换成二维码，
     * 使用支付宝或微信的扫一扫功能即可在手机端完成支付。可通过查询订单接口确定此交易的交易状态。
     *
     * @param orderData
     * @param listener
     */
    public static void startPrePay(OrderData orderData, final CashierListener listener) {
        String str = orderData.txamt;
        String discount = orderData.discountMoney;
        if (str != null && orderData.currency != null) {
            if (SDK_CURRENCY.equals(orderData.currency)) {
                orderData.txamt = TxamtUtil.getTxamtUtil(str);
                orderData.discountMoney = TxamtUtil.getTxamtUtil(discount);
                if (orderData.txamt == null) {
                    listener.onError(SDK_ERROR_TXAMT_NULL);
                    return;
                }
            } else {
                listener.onError(SDK_ERROR_CURRENCY_NOT_RMB);
                return;
            }
        }

        CommunicationUtil.sendDataToServer(ParamsUtil.getPrePay(mInitData, orderData), new CommunicationListener() {

            @Override
            public void onResult(String result) {
                Map<String, Object> map = MapUtil.getMapForJson(result);
                String sign = (String) map.get("sign");
                if (sign != null) {
                    map.remove("sign");
                    String veriSign = ParamsUtil.getSign(MapUtil.getSignString(map), mInitData.signKey);
                    if (sign.equals(veriSign)) {
                        ResultData resultData = MapUtil.getResultData(map);
                        listener.onResult(resultData);
                    } else {
                        listener.onError(SDK_ERROR_SIGN_NOT_MATCH);
                    }
                }
            }

            @Override
            public void onError(int error) {
                listener.onError(error);

            }
        });
    }


    /**
     * 3.3. 查询订单
     * 此接口用于查询下单支付和预下单交易的交易状态，当下单支付或预下单
     * 支付返回09：交易处理中或98：交易超时的应答码或者接入方没有收到应答时
     * （网络原因等），推荐需要调用此接口，以明确订单状态。
     *
     * @param orderData
     * @param listener
     */
    public static void startQy(OrderData orderData, final CashierListener listener) {
        CommunicationUtil.sendDataToServer(ParamsUtil.getQy(mInitData, orderData), new CommunicationListener() {

            @Override
            public void onResult(String result) {
                Map<String, Object> map = MapUtil.getMapForJson(result);
                String sign = (String) map.get("sign");
                if (sign != null) {
                    map.remove("sign");
                    String veriSign = ParamsUtil.getSign(MapUtil.getSignString(map), mInitData.signKey, "SHA-1");
                    Log.i(TAG, "verisign: " + veriSign);
                    if (sign.equals(veriSign)) {
                        ResultData resultData = MapUtil.getResultData(map);
                        listener.onResult(resultData);
                    } else {
                        listener.onError(SDK_ERROR_SIGN_NOT_MATCH);
                    }
                }
            }

            @Override
            public void onError(int error) {
                listener.onError(error);
            }
        });
    }

    public static ResultData startQy(OrderData orderData) {
        String result = CommunicationUtil.sendDataToServer(ParamsUtil.getQy(mInitData, orderData));
        if (result == null || result.length() == 0) {
            return null;
        }
        Map<String, Object> map = MapUtil.getMapForJson(result);
        String sign = (String) map.get("sign");
        if (sign != null) {
            map.remove("sign");
            String veriSign = ParamsUtil.getSign(MapUtil.getSignString(map), mInitData.signKey, "SHA-1");
            Log.i(TAG, "verisign: " + veriSign);
            if (sign.equals(veriSign)) {
                ResultData resultData = MapUtil.getResultData(map);
                return resultData;
            } else {
                return null;
            }
        }
        return null;
    }


    /**
     * * 3.4. 撤销
     * 撤销下单支付或预下单交易
     *
     * @param orderData
     * @param listener
     */
    public static void startVoid(OrderData orderData, final CashierListener listener) {
        CommunicationUtil.sendDataToServer(ParamsUtil.getVoid(mInitData, orderData), new CommunicationListener() {

            @Override
            public void onResult(String result) {
                Map<String, Object> map = MapUtil.getMapForJson(result);
                String sign = (String) map.get("sign");
                if (sign != null) {
                    map.remove("sign");
                    String veriSign = ParamsUtil.getSign(MapUtil.getSignString(map), mInitData.signKey, "SHA-1");
                    if (sign.equals(veriSign)) {
                        ResultData resultData = MapUtil.getResultData(map);
                        listener.onResult(resultData);
                    } else {
                        listener.onError(SDK_ERROR_SIGN_NOT_MATCH);
                    }
                }
            }

            @Override
            public void onError(int error) {
                listener.onError(error);

            }
        });
    }

    /**
     * * 3.5. 退款
     * 撤销下单支付或预下单交易
     *
     * @param orderData
     * @param listener
     */
    public static void startRefd(OrderData orderData, final CashierListener listener) {
        String str = orderData.txamt;
        if (str != null && orderData.currency != null) {
            if (SDK_CURRENCY.equals(orderData.currency)) {
                orderData.txamt = TxamtUtil.getTxamtUtil(str);
                if (orderData.txamt == null) {
                    listener.onError(SDK_ERROR_TXAMT_NULL);
                    return;
                }
            } else {
                listener.onError(SDK_ERROR_CURRENCY_NOT_RMB);
                return;
            }
        }

        CommunicationUtil.sendDataToServer(ParamsUtil.getRefd(mInitData, orderData), new CommunicationListener() {

            @Override
            public void onResult(String result) {
                Map<String, Object> map = MapUtil.getMapForJson(result);
                String sign = (String) map.get("sign");
                if (sign != null) {
                    map.remove("sign");
                    String veriSign = ParamsUtil.getSign(MapUtil.getSignString(map), mInitData.signKey);
                    Log.i(TAG, "veriSign: " + veriSign);
                    if (sign.equals(veriSign)) {
                        ResultData resultData = MapUtil.getResultData(map);
                        listener.onResult(resultData);
                    } else {
                        listener.onError(SDK_ERROR_SIGN_NOT_MATCH);
                    }
                }
            }

            @Override
            public void onError(int error) {
                listener.onError(error);
            }
        });
    }


    /**
     * 3.6  取消订单
     * 对于未成功付款的订单进行取消，则关闭交易，使用户后期不能支付成功；
     * 对于成功付款的订单进行取消，系统将订单金额返还给用户，相当于对此交易做撤销。
     *
     * @param orderData
     * @param listener
     */
    public static void startCanc(OrderData orderData, final CashierListener listener) {
        CommunicationUtil.sendDataToServer(ParamsUtil.getCanc(mInitData, orderData), new CommunicationListener() {

            @Override
            public void onResult(String result) {
                Map<String, Object> map = MapUtil.getMapForJson(result);
                String sign = (String) map.get("sign");
                if (sign != null) {
                    map.remove("sign");
                    String veriSign = ParamsUtil.getSign(MapUtil.getSignString(map), mInitData.signKey);
                    Log.i(TAG, "veriSign: " + veriSign);
                    if (sign.equals(veriSign)) {
                        ResultData resultData = MapUtil.getResultData(map);
                        listener.onResult(resultData);
                    } else {
                        listener.onError(SDK_ERROR_SIGN_NOT_MATCH);
                    }
                }
            }

            @Override
            public void onError(int error) {
                listener.onError(error);
            }
        });
    }

    /**
     * 3.7卡券核销
     * 对上送的卡券号进行核销
     *
     * @param orderData
     * @param listener
     */
    public static void startVeri(OrderData orderData, final CashierListener listener) {
        CommunicationUtil.sendDataToServer(ParamsUtil.getVeri(mInitData, orderData), new CommunicationListener() {

            @Override
            public void onResult(String result) {
                Map<String, Object> map = MapUtil.getMapForJson(result);
                String sign = (String) map.get("sign");
                if (sign != null) {
                    map.remove("sign");
                    String veriSign = ParamsUtil.getSign(MapUtil.getSignString(map), mInitData.signKey);
                    if (sign.equals(veriSign)) {
                        ResultData resultData = MapUtil.getResultData(map);
                        listener.onResult(resultData);
                    } else {
                        listener.onError(SDK_ERROR_SIGN_NOT_MATCH);
                    }
                }
            }

            @Override
            public void onError(int error) {
                listener.onError(error);
            }
        });
    }


    /**
     * 卡券冲正
     */
    public static void startReversal(OrderData orderData, final CashierListener listener) {
        CommunicationUtil.sendDataToServer(ParamsUtil.getReveral(mInitData, orderData), new CommunicationListener() {

            @Override
            public void onResult(String result) {
                Map<String, Object> map = MapUtil.getMapForJson(result);
                String sign = (String) map.get("sign");
                if (sign != null) {
                    map.remove("sign");
                    String veriSign = ParamsUtil.getSign(MapUtil.getSignString(map), mInitData.signKey);
                    if (sign.equals(veriSign)) {
                        ResultData resultData = MapUtil.getResultData(map);
                        listener.onResult(resultData);
                    } else {
                        listener.onError(SDK_ERROR_SIGN_NOT_MATCH);
                    }
                }
            }

            @Override
            public void onError(int error) {
                listener.onError(error);
            }
        });

    }


}
