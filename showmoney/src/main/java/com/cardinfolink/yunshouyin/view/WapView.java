package com.cardinfolink.yunshouyin.view;

import android.content.Context;
import android.view.LayoutInflater;
import android.view.View;
import android.webkit.WebChromeClient;
import android.webkit.WebSettings;
import android.webkit.WebView;
import android.webkit.WebViewClient;
import android.widget.LinearLayout;

import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.constant.SystemConfig;
import com.cardinfolink.yunshouyin.data.SessonData;


public class WapView extends LinearLayout {
    private Context mContext;
    private WebView mWebView;

    public WapView(Context context) {
        super(context);
        mContext = context;

        View contentView = LayoutInflater.from(context).inflate(
                R.layout.wap_bill_view, null);
        LinearLayout.LayoutParams layoutParams = new LayoutParams(
                LayoutParams.MATCH_PARENT, LayoutParams.MATCH_PARENT);
        contentView.setLayoutParams(layoutParams);
        addView(contentView);
        mWebView = (WebView) contentView.findViewById(R.id.base_webview);
        mWebView.setWebViewClient(new WebViewClient() {
            @Override
            public boolean shouldOverrideUrlLoading(WebView view, String url) {
                view.loadUrl(url);
                return true;
            }

            @Override
            public void onPageFinished(WebView view, String url) {
                super.onPageFinished(view, url);

            }
        });


        mWebView.setWebChromeClient(new WebChromeClient() {
                                        @Override
                                        public boolean onJsAlert(WebView view, String url, String message, android.webkit.JsResult result) {

                                            return false;


                                        }

                                        ;
                                    }
        );


        WebSettings webSettings = mWebView.getSettings();

        webSettings.setDefaultTextEncodingName("utf-8");
        webSettings.setJavaScriptEnabled(true);
        webSettings.setAllowFileAccess(true);// 设置允许访问文件数据
        webSettings.setJavaScriptCanOpenWindowsAutomatically(true);
        webSettings.setDomStorageEnabled(true);
        webSettings.setDatabaseEnabled(true);
        mWebView.requestFocus();
        webSettings.setSupportZoom(true);

        initData();

    }

    public void initData() {
        if (SessonData.loginUser.getObjectId() != null && SessonData.loginUser.getObjectId().length() > 0) {

            mWebView.loadUrl(SystemConfig.WEB_BILL_URL + "?merchantCode=" + SessonData.loginUser.getObjectId());
        }

    }

}
