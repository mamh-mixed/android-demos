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
            CommunicationUtil.setServer(server);
        }

    }


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
        if (str != null && orderData.currency != null) {
            if (orderData.currency.equals("156")) {
                orderData.txamt = TxamtUtil.getTxamtUtil(str);
                if (orderData.txamt == null) {
                    listener.onError(1);
                    return;
                }
            } else {
                listener.onError(2);
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

                    String veriSign = ParamsUtil.getSign(MapUtil.getSignString(map), mInitData.signKey, "SHA-1");
                    Log.i(TAG, "veriSign: " + veriSign);
                    if (sign.equals(veriSign)) {
                        ResultData resultData = MapUtil.getResultData(map);
                        listener.onResult(resultData);

                    } else {
                        Log.i(TAG, "签名不一致");
                        listener.onError(3);
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
        if (str != null && orderData.currency != null) {
            if (orderData.currency.equals("156")) {
                orderData.txamt = TxamtUtil.getTxamtUtil(str);
                if (orderData.txamt == null) {
                    listener.onError(1);
                    return;
                }
            } else {
                listener.onError(2);
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

                    String veriSign = ParamsUtil.getSign(MapUtil.getSignString(map), mInitData.signKey, "SHA-1");
                    if (sign.equals(veriSign)) {
                        ResultData resultData = MapUtil.getResultData(map);
                        listener.onResult(resultData);

                    } else {
                        listener.onError(3);
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
        String str = orderData.txamt;
        if (str != null && orderData.currency != null) {
            if (orderData.currency.equals("156")) {
                orderData.txamt = TxamtUtil.getTxamtUtil(str);
                if (orderData.txamt == null) {
                    listener.onError(1);
                    return;
                }
            } else {
                listener.onError(2);
                return;
            }
        }

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
                        listener.onError(3);
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
     * * 3.4. 撤销
     * 撤销下单支付或预下单交易
     *
     * @param orderData
     * @param listener
     */
    public static void startVoid(OrderData orderData, final CashierListener listener) {

        String str = orderData.txamt;
        if (str != null && orderData.currency != null) {
            if (orderData.currency.equals("156")) {
                orderData.txamt = TxamtUtil.getTxamtUtil(str);
                if (orderData.txamt == null) {
                    listener.onError(1);
                    return;
                }
            } else {
                listener.onError(2);
                return;
            }
        }

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
                        listener.onError(3);
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
            if (orderData.currency.equals("156")) {
                orderData.txamt = TxamtUtil.getTxamtUtil(str);
                if (orderData.txamt == null) {
                    listener.onError(1);
                    return;
                }
            } else {
                listener.onError(2);
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

                    String veriSign = ParamsUtil.getSign(MapUtil.getSignString(map), mInitData.signKey, "SHA-1");
                    Log.i(TAG, "veriSign: " + veriSign);
                    if (sign.equals(veriSign)) {
                        ResultData resultData = MapUtil.getResultData(map);
                        listener.onResult(resultData);

                    } else {
                        listener.onError(3);
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
                    String veriSign = ParamsUtil.getSign(MapUtil.getSignString(map), mInitData.signKey, "SHA-1");
                    Log.i(TAG, "veriSign: " + veriSign);
                    if (sign.equals(veriSign)) {
                        ResultData resultData = MapUtil.getResultData(map);
                        listener.onResult(resultData);
                    } else {
                        listener.onError(3);
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

        String str = orderData.txamt;
        if (str != null && orderData.currency != null) {
            if (orderData.currency.equals("156")) {
                orderData.txamt = TxamtUtil.getTxamtUtil(str);
                if (orderData.txamt == null) {
                    listener.onError(1);
                    return;
                }
            } else {
                listener.onError(2);
                return;
            }
        }

        CommunicationUtil.sendDataToServer(ParamsUtil.getVeri(mInitData, orderData), new CommunicationListener() {

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
                        listener.onError(3);
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
