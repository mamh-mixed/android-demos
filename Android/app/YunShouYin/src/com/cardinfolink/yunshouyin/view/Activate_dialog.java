package com.cardinfolink.yunshouyin.view;

import com.cardinfo.framelib.constant.Msg;
import com.cardinfolink.yunshouyin.R;

import android.content.Context;
import android.os.Handler;
import android.view.MotionEvent;
import android.view.View;
import android.view.View.OnClickListener;
import android.view.View.OnTouchListener;
import android.widget.TextView;

public class Activate_dialog {
	private Context mContext;
	private Handler mHandler;
	private View dialogView;
	private String mEmali;
	
	public Activate_dialog(Context context,Handler handler,View view,String email) {
		mContext=context;
		mHandler=handler;
		dialogView=view;
		mEmali=email;
	}
	
	
	public void show(){
		TextView textView=(TextView) dialogView.findViewById(R.id.email);
		textView.setText("激活链接将发送到该邮箱:\n\n"+mEmali+"");
		dialogView.setVisibility(View.VISIBLE);
		dialogView.setOnTouchListener(new OnTouchListener() {
			
			@Override
			public boolean onTouch(View v, MotionEvent event) {
				// TODO Auto-generated method stub
				return true;
			}
		});
		dialogView.findViewById(R.id.activate_dialog_cancel).setOnClickListener(new OnClickListener() {
			
			@Override
			public void onClick(View v) {
				dialogView.setVisibility(View.GONE);
				mHandler.sendEmptyMessage(Msg.MSG_FROM_ACTIVATE_DIGLOG_CANCEL);
				
			}
		});
		
		dialogView.findViewById(R.id.activate_dialog_ok).setOnClickListener(new OnClickListener() {
			
			@Override
			public void onClick(View v) {
				dialogView.setVisibility(View.GONE);
				mHandler.sendEmptyMessage(Msg.MSG_FROM_ACTIVATE_DIGLOG_OK);
			}
		});
	}
}
