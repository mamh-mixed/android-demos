/*
 * To change this template, choose Tools | Templates
 * and open the template in the editor.
 */
package com.cardinfolink.cashiersdk.util;

import android.util.Log;

import java.io.BufferedReader;
import java.io.InputStreamReader;
import java.io.OutputStreamWriter;
import java.io.PrintWriter;
import java.net.Socket;


/**
 * @author webapp
 */
public class SocketClient {
    private static final String TAG = "SocketClient";
    private Socket sk;
    private String host = null;
    private String port = null;
    private int timeout = 15000;

    public SocketClient(
            String host,
            String port,
            int timeout
    ) {
        this.host = host;
        this.port = port;
        this.timeout = timeout;
    }


    public String reqToServer(String msg) {
        int length = msg.length();
        String lengthStr = "" + length;
        while (lengthStr.length() < 4) {
            lengthStr = "0" + lengthStr;
        }
        msg = lengthStr + msg;
        Log.e(TAG, "json :" + msg);
        PrintWriter out = null;
        BufferedReader in = null;


        try {
            this.sk = new Socket(this.host, Integer.parseInt(this.port));
            this.sk.setSoTimeout(timeout);

            out = new PrintWriter(new OutputStreamWriter(this.sk.getOutputStream(), "gbk"));
            in = new BufferedReader(new InputStreamReader(this.sk.getInputStream(), "gbk"));

            out.print(msg);
            out.flush();

            //get response
            char[] cbuf = new char[4096];
            int ret = in.read(cbuf);

            String retstr = String.copyValueOf(cbuf);

            return retstr;

        } catch (Exception e) {
            e.printStackTrace();
        } finally {
            if (this.sk != null)
                try {
                    this.sk.close();
                    out.close();
                    in.close();


                } catch (Exception e) {

                }

        }
        return "error";

    }

}
