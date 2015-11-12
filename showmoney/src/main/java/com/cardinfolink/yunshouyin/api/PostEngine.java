package com.cardinfolink.yunshouyin.api;

import org.apache.commons.io.IOUtils;

import java.io.BufferedInputStream;
import java.io.DataOutputStream;
import java.io.IOException;
import java.io.InputStream;
import java.net.HttpURLConnection;
import java.net.Proxy;
import java.net.URL;
import java.util.Map;

public class PostEngine {
    private Proxy httpProxy;

    public PostEngine() {
    }

    public PostEngine(Proxy httpProxy) {
        this.httpProxy = httpProxy;
    }

    /**
     * IOException 就是网络无法连接的问题了
     *
     * @param _url
     * @param params
     * @return
     * @throws IOException
     */
    public String post(String _url, Map<String, String> params) throws IOException {
        StringBuilder sb = new StringBuilder();
        boolean isFirst = true;
        for (String s : params.keySet()) {
            if (isFirst) {
                sb.append(s + "=");
                isFirst = false;
            } else {
                sb.append("&" + s + "=");
            }
            sb.append(params.get(s));
        }

        byte[] postData = sb.toString().getBytes("UTF-8");
        int postDataLength = postData.length;

        URL url = new URL(_url);
        HttpURLConnection conn = null;
        if (httpProxy != null) {
            conn = (HttpURLConnection) url.openConnection(httpProxy);
        } else {
            conn = (HttpURLConnection) url.openConnection();
        }
        try {
            conn.setDoOutput(true);
            conn.setRequestProperty("Content-Type", "application/x-www-form-urlencoded");
            conn.setRequestProperty("charset", "utf-8");
            conn.setRequestProperty("Content-Length", Integer.toString(postDataLength));
            conn.setRequestMethod("POST");

            DataOutputStream wr = new DataOutputStream(conn.getOutputStream());
            wr.write(postData);
            if (conn.getResponseCode() == 200) {
                InputStream in = new BufferedInputStream(conn.getInputStream());
                String response = IOUtils.toString(in, "UTF-8");
                return response;
            } else {
                return null;
            }

        } finally {
            conn.disconnect();
        }
    }


    public String get(String _url, Map<String, String> params) throws IOException {
        StringBuilder sb = new StringBuilder();
        boolean isFirst = true;
        for (String s : params.keySet()) {
            if (isFirst) {
                sb.append(s + "=");
                isFirst = false;
            } else {
                sb.append("&" + s + "=");
            }
            sb.append(params.get(s));
        }
        _url += "?" + sb.toString();

        URL url = new URL(_url);
        HttpURLConnection conn = null;
        if (httpProxy != null) {
            conn = (HttpURLConnection) url.openConnection(httpProxy);
        } else {
            conn = (HttpURLConnection) url.openConnection();
        }
        try {
            conn.setRequestProperty("charset", "utf-8");
            conn.setRequestMethod("GET");

            if (conn.getResponseCode() == 200) {
                InputStream in = new BufferedInputStream(conn.getInputStream());
                String response = IOUtils.toString(in, "UTF-8");
                return response;
            } else {
                //TODO: handle other http code
                return null;
            }

        } finally {
            conn.disconnect();
        }
    }
}
