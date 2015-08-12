package com.cardinfolink.yunshouyin.view;

import java.math.BigDecimal;
import java.text.SimpleDateFormat;
import java.util.Date;
import java.util.Random;

import com.cardinfo.framelib.constant.Msg;
import com.cardinfo.framelib.util.DeviceManageUtil;
import com.cardinfolink.cashiersdk.listener.CashierListener;
import com.cardinfolink.cashiersdk.model.OrderData;
import com.cardinfolink.cashiersdk.model.ResultData;
import com.cardinfolink.cashiersdk.sdk.CashierSdk;
import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.data.SaveData;
import com.cardinfolink.yunshouyin.data.User;

import android.content.Context;
import android.graphics.Bitmap;
import android.os.Handler;
import android.os.Message;
import android.util.Log;
import android.view.Gravity;
import android.view.MotionEvent;
import android.view.View;
import android.view.View.OnClickListener;
import android.view.View.OnTouchListener;
import android.view.inputmethod.InputMethodManager;
import android.widget.EditText;
import android.widget.TextView;
import android.widget.Toast;

public class Refd_Dialog {
	private Context mContext;
	private Handler mHandler;
	private View dialogView;
	private double maxRefd=0;
	private String mOrderNum;
	private Handler dialogHandler;
	public Refd_Dialog(Context context, Handler handler, View view,String orderNum,String refdTotal,String total){
		mContext=context;
		mHandler=handler;
	    dialogView=view;
	    mOrderNum=orderNum;
	    initHandler();
	    maxRefd=Double.parseDouble(total)-Double.parseDouble(refdTotal);
	    BigDecimal   b   =   new   BigDecimal(maxRefd);  
	    maxRefd   =   b.setScale(2,   BigDecimal.ROUND_HALF_UP).doubleValue();  
	}
	
	public void show(){
		
		TextView textView=(TextView) dialogView.findViewById(R.id.refd_title);
		textView.setText("本次可退款额度：¥"+maxRefd);
		final EditText refdValue=(EditText) dialogView.findViewById(R.id.refd_value_edit);
		final EditText refdPassword=(EditText) dialogView.findViewById(R.id.refd_password_edit);
		
		
		dialogView.setVisibility(View.VISIBLE);
		dialogView.setOnTouchListener(new OnTouchListener() {

			@Override
			public boolean onTouch(View v, MotionEvent event) {
				// TODO Auto-generated method stub
				return true;
			}
		});
		
		dialogView.findViewById(R.id.refd_dialog_cancel).setOnClickListener(new OnClickListener() {
			
			@Override
			public void onClick(View v) {
			    DeviceManageUtil.hideInput(mContext);
				dialogView.setVisibility(View.GONE);
				
			}
		});
		
		
   dialogView.findViewById(R.id.refd_dialog_ok).setOnClickListener(new OnClickListener() {
			
			@Override
			public void onClick(View v) {
				 DeviceManageUtil.hideInput(mContext);
				dialogView.setVisibility(View.GONE);
				
				String value=refdValue.getText().toString();
				String password=refdPassword.getText().toString();
				if(value.length()==0){
					Toast toast = Toast.makeText(mContext,
							"金额不能为空",
							Toast.LENGTH_SHORT);
					toast.setGravity(Gravity.CENTER, 0, 250);
					toast.show();
					return;
				}
				
				if(Double.parseDouble(value)>maxRefd){
					Toast toast = Toast.makeText(mContext,
							"余额不足",
							Toast.LENGTH_SHORT);
					toast.setGravity(Gravity.CENTER, 0, 250);
					toast.show();
					return;
				}
				User user=SaveData.getUser(mContext);
				
				if(!password.equals(user.getPassword())){
					Toast toast = Toast.makeText(mContext,
							"密码错误",
							Toast.LENGTH_SHORT);
					toast.setGravity(Gravity.CENTER, 0, 250);
					toast.show();
					return;
				}
				   OrderData orderData=new OrderData();
				    orderData.origOrderNum=mOrderNum;
				    Log.i("opp",  orderData.origOrderNum);
				    Date now = new Date();
					SimpleDateFormat spf = new SimpleDateFormat("yyMMddHHmmss");
					String orderNmuber =spf.format(now);
					Random random = new Random();
					for (int i = 0; i < 5; i++) {
						orderNmuber = orderNmuber + random.nextInt(10);
					};	
				    orderData.orderNum=orderNmuber;
				    orderData.currency="156";
				    orderData.txamt=value;
				    
				CashierSdk.startRefd(orderData, new CashierListener() {
					
					@Override
					public void onResult(ResultData resultData) {
						dialogHandler.sendEmptyMessage(100);
					  if(resultData.respcd.equals("00")){
						  mHandler.sendEmptyMessage(Msg.MSG_FROM_SERVER_REFD_SUCCESS);
					  }else{
						  mHandler.sendEmptyMessage( Msg.MSG_FROM_SERVER_REFD_FAIL);
					  }
						
					}
					
					@Override
					public void onError(int errorCode) {
						dialogHandler.sendEmptyMessage(100);
						mHandler.sendEmptyMessage(Msg.MSG_FROM_SERVER_TIMEOUT);
						
					}
				});
			}
		});
	}
	
	
	 private void initHandler(){
		   dialogHandler=new Handler(){
			   @Override
			public void handleMessage(Message msg) {
				   switch(msg.what){
				   case 100:{
					    final EditText refdValue=(EditText) dialogView.findViewById(R.id.refd_value_edit);
						final EditText refdPassword=(EditText) dialogView.findViewById(R.id.refd_password_edit);
					    refdPassword.setText("");
					    refdValue.setText("");
				   }
				   }
				super.handleMessage(msg);
			}
		   };
	   }
}
