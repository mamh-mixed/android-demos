package com.cardinfolink.yunshouyin.activity;

import java.text.ParseException;
import java.text.SimpleDateFormat;
import java.util.Date;
import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.data.SessonData;
import com.cardinfolink.yunshouyin.data.TradeBill;
import com.cardinfolink.yunshouyin.util.CommunicationListener;
import com.cardinfolink.yunshouyin.util.ErrorUtil;
import com.cardinfolink.yunshouyin.util.HttpCommunicationUtil;
import com.cardinfolink.yunshouyin.util.JsonUtil;
import com.cardinfolink.yunshouyin.util.ParamsUtil;
import com.cardinfolink.yunshouyin.view.Refd_Dialog;

import android.annotation.SuppressLint;
import android.app.Activity;
import android.content.Intent;
import android.graphics.BitmapFactory;
import android.graphics.Color;
import android.os.Bundle;
import android.view.View;
import android.widget.Button;
import android.widget.ImageView;
import android.widget.TextView;

public class DetailActivity extends BaseActivity {
	 private TradeBill mTradeBill;
	 private ImageView mPaylogoImage;
	 private TextView mTradeFromText;
	 private TextView mTradeDateText;
	 private TextView mTradeStatusText;
	 private TextView mConsumerAccount;
	 private TextView mTradeAmountText;
	 private TextView mOrderNumText;
	 private TextView mGoodInfoText;
	 private Button mRefdButton;
	@Override
	protected void onCreate(Bundle savedInstanceState) {
		super.onCreate(savedInstanceState);
		setContentView(R.layout.detail_activity);
		Intent intent=getIntent();
		Bundle billBundle=intent.getBundleExtra("BillBundle");
		mTradeBill=(TradeBill) billBundle.get("TradeBill");
		initLayout();
		initData();
	}

	private void initLayout(){
		mPaylogoImage=(ImageView) findViewById(R.id.paylogo);
		mTradeFromText=(TextView) findViewById(R.id.tradefrom);
		mTradeDateText=(TextView) findViewById(R.id.tradedate);
		mTradeStatusText=(TextView) findViewById(R.id.tradestatus);
		mConsumerAccount=(TextView) findViewById(R.id.consumer_account);
		mTradeAmountText=(TextView) findViewById(R.id.tradeamount);
		mOrderNumText=(TextView) findViewById(R.id.ordernum);
		mGoodInfoText=(TextView) findViewById(R.id.goodinfo);
		mRefdButton=(Button) findViewById(R.id.detail_btn_refd);
	}
	
	@SuppressLint("NewApi") private void initData(){
		if (mTradeBill.chcd.equals("WXP")) {
			mPaylogoImage.setImageResource(R.drawable.wpay);
		} else {
			mPaylogoImage.setImageResource(R.drawable.apay);
		}
		SimpleDateFormat spf1 = new SimpleDateFormat("yyyyMMddHHmmss");
		SimpleDateFormat spf2 = new SimpleDateFormat("yyyy-MM-dd HH:mm:ss");
		try {
			Date tandeDate = spf1.parse(mTradeBill.tandeDate);
			mTradeDateText.setText(spf2.format(tandeDate));
		} catch (ParseException e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
		}
		String tradeFrom = "PC";
		if (! mTradeBill.tradeFrom.isEmpty()) {
			tradeFrom = mTradeBill.tradeFrom;
		}
		String busicd = "支付";
		if (mTradeBill.busicd.equals("REFD")) {
			busicd = "退款";
		}

		mTradeFromText.setText(tradeFrom + busicd);
		String tradeStatus = "交易成功";
		if (mTradeBill.response.equals("00")) {
			tradeStatus = "交易成功";
			mTradeStatusText
					.setTextColor(Color.parseColor("#888888"));
		} else if (mTradeBill.response.equals("09")) {
			tradeStatus = "未支付";
			mTradeStatusText.setTextColor(Color.RED);
		} else {
			tradeStatus = "交易失败";
			mTradeStatusText.setTextColor(Color.RED);
		}
		mTradeStatusText.setText(tradeStatus);
		mConsumerAccount.setText(mTradeBill.consumerAccount);
		mTradeAmountText.setText("￥"+mTradeBill.amount);
		mOrderNumText.setText(mTradeBill.orderNum);
		mGoodInfoText.setText(mTradeBill.goodsInfo);
		if(mTradeBill.busicd.equals("REFD")||!mTradeBill.response.equals("00")){
			mRefdButton.setVisibility(View.INVISIBLE);
		}else{
			mRefdButton.setVisibility(View.VISIBLE);
		}
	}
	
	
	 public void BtnBackOnClick(View view){    
	        	    
		 DetailActivity.this.finish();  
          
     }  
	 
	 
	 public void BtnRefdOnClick(View view){ 
		 startLoading();
		 HttpCommunicationUtil.sendDataToServer(ParamsUtil.getRefd(SessonData.loginUser, mTradeBill.orderNum), new CommunicationListener() {
			
			@Override
			public void onResult(final String result) {
				String state=JsonUtil.getParam(result, "state");
				if(state.equals("success")){				
					final String refdtotal=JsonUtil.getParam(result, "refdtotal");
					runOnUiThread(new Runnable(){ 
						
						  
	                    @Override  
	                    public void run() {  
	                        //更新UI                      	
	                    	endLoading();
	                    	 Refd_Dialog refd_Dialog=new Refd_Dialog(DetailActivity.this, null, findViewById(R.id.refd_dialog), mTradeBill.orderNum, refdtotal,mTradeBill.amount);
	                  	    refd_Dialog.show();
	                    }  
	                      
	                }); 
					
				}else{
					runOnUiThread(new Runnable(){ 
						
						  
	                    @Override  
	                    public void run() {  
	                        //更新UI                      	
	                    	endLoading();
	                    	 mAlert_Dialog.show(ErrorUtil.getErrorString(JsonUtil.getParam(result, "error")),  BitmapFactory.decodeResource(mContext.getResources(), R.drawable.wrong));
	                    }  
	                      
	                }); 
				}
				
				
			}
			
			@Override
			public void onError(final String error) {
				runOnUiThread(new Runnable(){ 
					
					  
                    @Override  
                    public void run() {  
                        //更新UI                      	
                    	endLoading();
                    	 mAlert_Dialog.show(error,  BitmapFactory.decodeResource(mContext.getResources(), R.drawable.wrong));
                    }  
                      
                }); 
				
			}
		});
		 
 	   
          
     }  
	 
	
}
