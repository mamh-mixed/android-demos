package com.cardinfolink.cashiersdk.data;

import android.app.Activity;
import android.content.Context;
import android.content.SharedPreferences;

import com.cardinfolink.cashiersdk.model.Server;

public class SaveData {
    private static Server mServer = new Server();

    public static Server getmServer(Context context) {
        SharedPreferences mySharedPreferences = context.getSharedPreferences("savedata_sdk",
                Activity.MODE_PRIVATE);
        mServer.setHost(mySharedPreferences.getString("host", "211.147.72.70"));
        mServer.setPort(mySharedPreferences.getString("port", "10008"));
        return mServer;
    }

    public static void setmServer(Context context, Server mServer) {
        SaveData.mServer = mServer;
        SharedPreferences mySharedPreferences = context.getSharedPreferences("savedata_sdk",
                Activity.MODE_PRIVATE);
        SharedPreferences.Editor editor = mySharedPreferences.edit();
        editor.putString("host", mServer.getHost());
        editor.putString("port", mServer.getPort());
        editor.commit();
    }


}
