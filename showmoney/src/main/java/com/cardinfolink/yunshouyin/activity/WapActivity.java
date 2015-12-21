package com.cardinfolink.yunshouyin.activity;

import android.os.Bundle;
import android.view.View;
import android.webkit.WebChromeClient;
import android.webkit.WebSettings;
import android.webkit.WebView;
import android.webkit.WebViewClient;

import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.constant.SystemConfig;
import com.cardinfolink.yunshouyin.data.SessonData;
import com.cardinfolink.yunshouyin.ui.SettingActionBarItem;


public class WapActivity extends BaseActivity {
    private static final String TAG = "WapActivity";
    private SettingActionBarItem mActionBar;
    private WebView mWebView;

    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.wap_bill_view);

        mActionBar = (SettingActionBarItem) findViewById(R.id.action_bar);
        mActionBar.setLeftTextOnclickListner(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                finish();
            }
        });

        mWebView = (WebView) findViewById(R.id.base_webview);
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
        });


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

    @Override
    public void onBackPressed() {
        super.onBackPressed();
        finish();
    }
}
