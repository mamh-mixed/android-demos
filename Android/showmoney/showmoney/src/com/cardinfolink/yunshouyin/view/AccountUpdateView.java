package com.cardinfolink.yunshouyin.view;


import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.activity.BaseActivity;
import com.cardinfolink.yunshouyin.data.SessonData;
import com.cardinfolink.yunshouyin.data.User;
import com.cardinfolink.yunshouyin.util.CommunicationListener;
import com.cardinfolink.yunshouyin.util.HttpCommunicationUtil;
import com.cardinfolink.yunshouyin.util.JsonUtil;
import com.cardinfolink.yunshouyin.util.ParamsUtil;
import com.cardinfolink.yunshouyin.util.VerifyUtil;

import android.annotation.SuppressLint;
import android.app.Activity;
import android.content.Context;
import android.graphics.BitmapFactory;
import android.view.LayoutInflater;
import android.view.View;
import android.widget.EditText;
import android.widget.LinearLayout;

public class AccountUpdateView extends LinearLayout {
	private EditText mOpenBankEdit;
	private EditText mNameEdit;
	private EditText mBanknumEdit;
	private EditText mPhonenumEdit;
	private Context mContext;
	private BaseActivity mBaseActivity;

	public AccountUpdateView(Context context) {
		super(context);
		mContext=context;
		mBaseActivity=(BaseActivity) mContext;
		View contentView = LayoutInflater.from(context).inflate(
				R.layout.account_update_view, null);
		LinearLayout.LayoutParams layoutParams = new LayoutParams(
				LayoutParams.MATCH_PARENT, LayoutParams.MATCH_PARENT);
		contentView.setLayoutParams(layoutParams);
		addView(contentView);
		mOpenBankEdit = (EditText) contentView.findViewById(R.id.info_openbank);
		mNameEdit = (EditText) contentView.findViewById(R.id.info_name);
		mBanknumEdit = (EditText)contentView.findViewById(R.id.info_banknum);
		mPhonenumEdit = (EditText) contentView.findViewById(R.id.info_phonenum);
		
	
		VerifyUtil.bankCardNumAddSpace(mBanknumEdit);
		
		
		contentView.findViewById(R.id.btn_submit).setOnClickListener(new OnClickListener() {
			
			@Override
			public void onClick(View v) {
				finishedOnClick() ;
			}
		});
		
	}
	
	public void initData(){
		HttpCommunicationUtil.sendDataToServer(ParamsUtil.getInfo(SessonData.loginUser), new CommunicationListener() {
			
			@Override
			public void onResult(final String result) {
				if(JsonUtil.getParam(result, "state").equals("success")){
					((Activity) mContext).runOnUiThread(new Runnable(){  
						  
	                    @Override  
	                    public void run() {  
	                    	 String info=JsonUtil.getParam(result, "info");
	                    	mOpenBankEdit.setText(JsonUtil.getParam(info, "bank_open"));
	    					mNameEdit.setText(JsonUtil.getParam(info, "payee"));
	    					mBanknumEdit.setText(JsonUtil.getParam(info, "payee_card"));
	    					mPhonenumEdit.setText(JsonUtil.getParam(info, "phone_num"));
							 
	                    }  
	                      
	                }); 
					
				}
				
			}
			
			@Override
			public void onError(String error) {
				// TODO Auto-generated method stub
				
			}
		});
	}
	
	
	public void finishedOnClick() {
		if (validate()) {
			mBaseActivity.startLoading();
			User user = new User();
			user.setUsername(SessonData.loginUser.getUsername());
			user.setPassword(SessonData.loginUser.getPassword());
			user.setBank_open(mOpenBankEdit.getText().toString());
			user.setPayee(mNameEdit.getText().toString());
			user.setPayee_card(mBanknumEdit.getText().toString()
					.replace(" ", ""));
			user.setPhone_num(mPhonenumEdit.getText().toString());
			HttpCommunicationUtil.sendDataToServer(
					ParamsUtil.getUpdateInfo(user),
					new CommunicationListener() {

						@Override
						public void onResult(String result) {
							String state = JsonUtil.getParam(result, "state");
							if (state.equals("success")) {
								((Activity) mContext).runOnUiThread(new Runnable() {

									@Override
									public void run() {
										// 更新UI
										mBaseActivity.endLoading();
										Alert_Dialog alert_Dialog=new Alert_Dialog(mContext, null, ((Activity)mContext).findViewById(R.id.alert_dialog), "修改成功",  BitmapFactory.decodeResource(mContext.getResources(), R.drawable.right));
										  alert_Dialog.show();
									}

								});
								
							} else {
								((Activity) mContext).runOnUiThread(new Runnable() {

									@Override
									public void run() {
										// 更新UI
										mBaseActivity.endLoading();
										mBaseActivity.alertShow(
												"提交失败!",
												BitmapFactory.decodeResource(
														mContext.getResources(),
														R.drawable.wrong));
									}

								});
							}
						}

						@Override
						public void onError(final String error) {
							((Activity) mContext).runOnUiThread(new Runnable() {

								@Override
								public void run() {
									// 更新UI
									mBaseActivity.endLoading();
									mBaseActivity.alertShow(error, BitmapFactory
											.decodeResource(
													mContext.getResources(),
													R.drawable.wrong));
								}

							});
						}
					});

			// Intent intent = new
			// Intent(RegisterNextActivity.this,MainActivity.class);
			// RegisterNextActivity.this.startActivity(intent);
			// RegisterNextActivity.this.finish();
		}

		// Intent intent = new
		// Intent(RegisterNextActivity.this,MainActivity.class);
		// RegisterNextActivity.this.startActivity(intent);
		// RegisterNextActivity.this.finish();
	}

	@SuppressLint("NewApi") private boolean validate() {
		String openbank = mOpenBankEdit.getText().toString().replace(" ", "");
		String name = mNameEdit.getText().toString().replace(" ", "");
		String banknum = mBanknumEdit.getText().toString().replace(" ", "");
		String phonenum = mPhonenumEdit.getText().toString().replace(" ", "");
		if (openbank.isEmpty()) {
			mBaseActivity.alertShow("开户行不能为空!", BitmapFactory.decodeResource(
					this.getResources(), R.drawable.wrong));
			return false;
		}

		if (name.isEmpty()) {
			mBaseActivity.alertShow("姓名不能为空!", BitmapFactory.decodeResource(
					this.getResources(), R.drawable.wrong));
			return false;
		}

		if (banknum.isEmpty()) {
			mBaseActivity.alertShow("银行卡号不能为空!", BitmapFactory.decodeResource(
					this.getResources(), R.drawable.wrong));
			return false;
		}

		if (!VerifyUtil.checkBankCard(banknum)) {
			mBaseActivity.alertShow("请输入正确的银行卡号!", BitmapFactory.decodeResource(
					this.getResources(), R.drawable.wrong));
			return false;
		}

		if (phonenum.isEmpty()) {
			mBaseActivity.alertShow("手机号不能为空!", BitmapFactory.decodeResource(
					this.getResources(), R.drawable.wrong));
			return false;
		}

		if (!VerifyUtil.isMobileNO(phonenum)) {
			mBaseActivity.alertShow("请输入正确的手机号!", BitmapFactory.decodeResource(
					this.getResources(), R.drawable.wrong));
			return false;
		}

		return true;
	}

}
