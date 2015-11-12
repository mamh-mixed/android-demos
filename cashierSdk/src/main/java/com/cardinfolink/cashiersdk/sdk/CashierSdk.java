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


    public static final String CURRENCY_CNY = "156";
    public static final String CURRENCY_JPY = "JPY";

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

    //NOTE: 服务器端存到分这个程度,补全成12位,位数不够再补零
    //TODO: 需要单元测试
    private static String formatAmount(String source, String currency) {
        double parsed;
        try {
            parsed = Double.parseDouble(source);
        } catch (Exception ex) {
            //非真实数字
            return null;
        }

        switch (currency) {
            case CURRENCY_JPY:
                break;
            case CURRENCY_CNY:
            default:
                //人民币小数点后有两位
                parsed = parsed * 100;
                break;
        }

        long sum = (long) parsed;

        String result = String.format("%012d", sum);
        return result;
    }


    public static void startPay(OrderData orderData, final CashierListener listener) {
        String str = orderData.txamt;
        // TODO: validation
        if (str != null) {
            orderData.txamt = formatAmount(str, orderData.currency);
            if (orderData.txamt == null) {
                listener.onError(1);
                return;
            }
        } else {
            listener.onError(1);
            return;
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


    public static void startPrePay(OrderData orderData, final CashierListener listener) {
        String str = orderData.txamt;
        // TODO: validation
        if (str != null) {
            orderData.txamt = formatAmount(str, orderData.currency);
            if (orderData.txamt == null) {
                listener.onError(1);
                return;
            }
        } else {
            listener.onError(1);
            return;
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

    public static void startRefd(OrderData orderData, final CashierListener listener) {
        String str = orderData.txamt;
        // TODO: validation
        if (str != null) {
            orderData.txamt = formatAmount(str, orderData.currency);
            if (orderData.txamt == null) {
                listener.onError(1);
                return;
            }
        } else {
            listener.onError(1);
            return;
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

    public static void startVeri(OrderData orderData, final CashierListener listener) {
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
