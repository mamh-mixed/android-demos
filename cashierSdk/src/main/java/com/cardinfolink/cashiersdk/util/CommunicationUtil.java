package com.cardinfolink.cashiersdk.util;

import android.util.Log;

import com.cardinfolink.cashiersdk.listener.CommunicationListener;
import com.cardinfolink.cashiersdk.model.Server;

import org.json.JSONObject;

public class CommunicationUtil {
    private static final String TAG = "CommunicationUtil";
    private static String mHost = "211.147.72.70";
    private static String mPort = "10008";

    //private static String mHost="192.168.199.174:";
    //private static String mPort="3001";

    public static void setServer(Server server) {
        mHost = server.getHost();
        mPort = server.getPort();
    }

    public static void sendDataToServer(
            final JSONObject json, final CommunicationListener listener) {

        new Thread(new Runnable() {

            @Override
            public void run() {
                Log.i(TAG, "mHost=" + mHost + " mPort" + mPort);
                SocketClient socketClient = new SocketClient(mHost, mPort, 15000);
                String result = socketClient.reqToServer(json.toString());
                Log.e(TAG, "result" + result);
                if (result != null && result.length() > 0) {
                    if (result.contains("}")) {
                        result = result.substring(4, result.lastIndexOf("}") + 1);
                        Log.e(TAG, "result" + result);
                        listener.onResult(result);
                    } else {
                        listener.onError(0);
                    }
                } else {
                    listener.onError(4);
                }


            }
        }).start();
    }


}
