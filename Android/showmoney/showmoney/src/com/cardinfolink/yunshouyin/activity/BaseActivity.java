package com.cardinfolink.yunshouyin.activity;

import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.view.Alert_Dialog;
import com.cardinfolink.yunshouyin.view.Loading_Dialog;
import com.umeng.analytics.MobclickAgent;

import android.app.Activity;
import android.content.Context;
import android.graphics.Bitmap;
import android.graphics.BitmapFactory;
import android.os.Bundle;


public class BaseActivity extends Activity {
	
	protected Loading_Dialog mLoading_Dialog;	
	protected Alert_Dialog mAlert_Dialog;
	protected Context mContext;
	
	@Override
	protected void onCreate(Bundle savedInstanceState) {
		super.onCreate(savedInstanceState);
		
	}
	
	@Override
	public void setContentView(int layoutResID) {		
		super.setContentView(layoutResID);
		mContext=this;
		mLoading_Dialog=new Loading_Dialog(this, findViewById(R.id.loading_dialog));
		mAlert_Dialog=new Alert_Dialog(this, null, findViewById(R.id.alert_dialog), 
				   getResources().getString(R.string.username_password_error), BitmapFactory.decodeResource(this.getResources(), R.drawable.wrong));
	}

	
	public void startLoading(){
		mLoading_Dialog.startLoading();
	}
	
	public void endLoading(){
		mLoading_Dialog.endLoading();
	}
	
	public void alertShow(String msg,Bitmap bitmap){
		mAlert_Dialog.show(msg, bitmap);
	}
	
	@Override
	protected void onResume() {
		super.onResume();
		MobclickAgent.onResume(this);
	}
	
	@Override
	protected void onPause() {
		super.onPause();
		MobclickAgent.onPause(this);
	}
}
