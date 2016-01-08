package com.cardinfolink.yunshouyin.api;

import org.apache.commons.io.IOUtils;

import java.io.BufferedInputStream;
import java.io.DataOutputStream;
import java.io.IOException;
import java.io.InputStream;
import java.net.HttpURLConnection;
import java.net.Proxy;
import java.net.URL;
import java.security.cert.CertificateException;
import java.security.cert.X509Certificate;
import java.util.Map;

import javax.net.ssl.HostnameVerifier;
import javax.net.ssl.HttpsURLConnection;
import javax.net.ssl.SSLContext;
import javax.net.ssl.SSLSession;
import javax.net.ssl.TrustManager;
import javax.net.ssl.X509TrustManager;

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

        if ("https".equals(url.getProtocol().toLowerCase())) {
            HttpsURLConnection https = (HttpsURLConnection)url.openConnection();
            trustAllHosts();
            https.setHostnameVerifier(DO_NOT_VERIFY);
            conn = https;
        } else {
            conn = (HttpURLConnection)url.openConnection();
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
                throw new IOException();
            }

        } finally {
            conn.disconnect();
        }
    }

    /**
     *
     */
    final static HostnameVerifier DO_NOT_VERIFY = new HostnameVerifier() {

        public boolean verify(String hostname, SSLSession session) {
            return true;
        }
    };
    /**
     * Trust every server - dont check for any certificate
     */
    private static void trustAllHosts() {
        final String TAG = "trustAllHosts";
        // Create a trust manager that does not validate certificate chains
        TrustManager[] trustAllCerts = new TrustManager[]{new X509TrustManager() {

            public java.security.cert.X509Certificate[] getAcceptedIssuers() {
                return new java.security.cert.X509Certificate[]{};
            }

            public void checkClientTrusted(X509Certificate[] chain, String authType) throws CertificateException {

            }

            public void checkServerTrusted(X509Certificate[] chain, String authType) throws CertificateException {

            }
        }};

        // Install the all-trusting trust manager
        try {
            SSLContext sc = SSLContext.getInstance("TLS");
            sc.init(null, trustAllCerts, new java.security.SecureRandom());
            HttpsURLConnection.setDefaultSSLSocketFactory(sc.getSocketFactory());
        } catch (Exception e) {
            e.printStackTrace();
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
