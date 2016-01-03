package com.cardinfolink.cashiersdk.util;

import android.text.TextUtils;

import com.cardinfolink.cashiersdk.listener.CommunicationListener;
import com.cardinfolink.cashiersdk.model.Server;
import com.cardinfolink.cashiersdk.sdk.CashierSdk;

import org.json.JSONObject;

public class CommunicationUtil {
    private static final String TAG = "CommunicationUtil";
    private static String mHost = "211.147.72.70";
    private static String mPort = "10008";

    public static void setServer(Server server) {
        mHost = server.getHost();
        mPort = server.getPort();
    }

    public static void sendDataToServer(final JSONObject json, final CommunicationListener listener) {
        new Thread(new Runnable() {

            @Override
            public void run() {
                SocketClient socketClient = new SocketClient(mHost, mPort, 15000);
                String result = socketClient.reqToServer(json.toString());
                if (!TextUtils.isEmpty(result)) {
                    if (result.contains("}")) {
                        result = result.substring(4, result.lastIndexOf("}") + 1);
                        listener.onResult(result);
                    } else {
                        //返回的 格式有误
                        listener.onError(CashierSdk.SDK_ERROR_RESULT_FORMAT);
                    }
                } else {
                    //返回的结构为空
                    listener.onError(CashierSdk.SDK_ERROR_RESULT_NULL);
                }
            }
        }).start();
    }

}
