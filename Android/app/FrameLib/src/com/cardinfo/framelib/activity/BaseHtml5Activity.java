package com.cardinfo.framelib.activity;


import com.cardinfo.framelib.R;
import com.cardinfo.framelib.constant.Msg;
import com.cardinfo.framelib.model.JavaJSParam;
import com.cardinfo.framelib.version.VersionManage;
import com.cardinfo.framelib.view.Loading_Dialog;
import android.annotation.SuppressLint;
import android.app.Activity;
import android.content.Context;
import android.os.Bundle;
import android.os.Handler;
import android.os.Message;

import android.util.Log;
import android.view.Gravity;
import android.webkit.JavascriptInterface;
import android.webkit.WebChromeClient;
import android.webkit.WebSettings;
import android.webkit.WebView;
import android.webkit.WebViewClient;
import android.widget.Toast;


public class BaseHtml5Activity extends Activity {
	private static final String TAG="BaseHtml5Activity";
	public WebView mWebView;
	public Handler mHandler;
	public Context mContext;
    public Loading_Dialog dialog;
	public VersionManage versionManage;
	@Override
	protected void onCreate(Bundle savedInstanceState) {		
		super.onCreate(savedInstanceState);
		setContentView(R.layout.base_html5_activity);
		mContext=this;
		initLayout();
		initHandler();
	   
	}
	
	@SuppressLint("SetJavaScriptEnabled") private void initLayout(){
		dialog=new Loading_Dialog(mContext, findViewById(R.id.loading_dialog));
		 mWebView=(WebView) findViewById(R.id.base_webview);
	     mWebView.setWebViewClient(new WebViewClient(){
			  @Override
			public boolean shouldOverrideUrlLoading(WebView view, String url) {
				view.loadUrl(url);
				return true;
			}
			  
			@Override
			public void onPageFinished(WebView view, String url) {
				super.onPageFinished(view, url);
				onWebViewPageFinished(view,url);
			}
		  });
	     
	     
	     mWebView.setWebChromeClient(new WebChromeClient(){
			  @Override
			  public boolean onJsAlert(WebView view, String url, String message, android.webkit.JsResult result) {
				  Toast toast=Toast.makeText(getApplicationContext(), message, Toast.LENGTH_SHORT); 
				  toast.setGravity(Gravity.CENTER, 0, 250);
			      toast.show();			 
				  result.confirm();
				   return true;
				  
				  
			  };		  }
		  );
	     
	     
	      WebSettings webSettings= mWebView.getSettings();
		 
	      webSettings.setDefaultTextEncodingName("utf-8") ;
		  webSettings.setJavaScriptEnabled(true);		  
		  webSettings.setAllowFileAccess(true);// 设置允许访问文件数据		 
		  webSettings.setJavaScriptCanOpenWindowsAutomatically(true);
		  webSettings.setDomStorageEnabled(true);
		  webSettings.setDatabaseEnabled(true);
		  mWebView.addJavascriptInterface(new JsObject(), "couldCashier"); 
		  mWebView.requestFocus();
		  webSettings.setSupportZoom(true);
		  
		
	}
	
	
	@SuppressLint("HandlerLeak") private void initHandler(){
		mHandler=new Handler(){
			@Override
			public void handleMessage(Message msg) {
				switch(msg.what){
				case Msg.MSG_FROM_CLIENT_LOAD_START:{
					dialog.startLoading();
					break;
				}
				
				case Msg.MSG_FROM_CLIENT_LOAD_END:{
					dialog.endLoading();
					break;
				}
				
				case Msg.MSG_FROM_CLIENT_REMIND_UPDATE:{
					versionManage.showUpdataDialog();
					break;
				}
				default:ProcessingMessage(msg);
				}
				
				
				super.handleMessage(msg);
			}
		};
	}
	

	
	public void ProcessingMessage(Message msg){
		
	}
	
	public void onWebViewPageFinished(WebView view, String url){
		
	}
	
	public class JsObject{
		 @JavascriptInterface
		public void OnDo(String method,String parameter){
			
			Log.i(TAG, "JsObject.OnDO:method="+method+" parameter="+parameter);
			JavaJSParam jsParam=new JavaJSParam();
			jsParam.setMethod(method);
			jsParam.setParam(parameter);
			Message msg=mHandler.obtainMessage();
			msg.what=Msg.MSG_FROM_JS;
			msg.obj=jsParam;
			mHandler.sendMessageDelayed(msg, 0);
		 }
		
	 } 

}
