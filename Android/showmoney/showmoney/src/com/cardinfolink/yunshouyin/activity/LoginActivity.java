package com.cardinfolink.yunshouyin.activity;

import com.cardinfolink.cashiersdk.model.InitData;
import com.cardinfolink.cashiersdk.sdk.CashierSdk;
import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.constant.SystemConfig;
import com.cardinfolink.yunshouyin.data.SaveData;
import com.cardinfolink.yunshouyin.data.SessonData;
import com.cardinfolink.yunshouyin.data.User;
import com.cardinfolink.yunshouyin.util.CommunicationListener;
import com.cardinfolink.yunshouyin.util.ErrorUtil;
import com.cardinfolink.yunshouyin.util.HttpCommunicationUtil;
import com.cardinfolink.yunshouyin.util.JsonUtil;
import com.cardinfolink.yunshouyin.util.ParamsUtil;
import com.cardinfolink.yunshouyin.util.TelephonyManagerUtil;
import com.cardinfolink.yunshouyin.util.VerifyUtil;
import com.cardinfolink.yunshouyin.view.Activate_dialog;
import com.cardinfolink.yunshouyin.view.Alert_Dialog;
import com.cardinfolink.yunshouyin.view.Loading_Dialog;
import com.umeng.update.UmengUpdateAgent;

import android.annotation.SuppressLint;
import android.app.Activity;
import android.content.Context;
import android.content.Intent;
import android.graphics.BitmapFactory;
import android.os.Bundle;
import android.view.View;
import android.widget.CheckBox;
import android.widget.EditText;

public class LoginActivity extends BaseActivity {
	
	private EditText mUsernameEdit;
	private EditText mPasswordEdit;
	private CheckBox mAutoLoginCheckBox;
	@Override
	protected void onCreate(Bundle savedInstanceState) {
		super.onCreate(savedInstanceState);
		setContentView(R.layout.login_activity);
		UmengUpdateAgent.setUpdateOnlyWifi(true);
		UmengUpdateAgent.setUpdateCheckConfig(false);
		UmengUpdateAgent.update(this);	
		mUsernameEdit=(EditText) findViewById(R.id.login_username);
		VerifyUtil.addEmialLimit(mUsernameEdit);
		mPasswordEdit=(EditText) findViewById(R.id.login_password);
		VerifyUtil.addEmialLimit(mPasswordEdit);
		mAutoLoginCheckBox=(CheckBox) findViewById(R.id.checkbox_auto_login);
		User user=SaveData.getUser(mContext);
		mAutoLoginCheckBox.setChecked(user.isAutoLogin());
		mUsernameEdit.setText(user.getUsername());
		mPasswordEdit.setText(user.getPassword());
		if(user.isAutoLogin()){
			 login();
		}
		
	}

	
	
	
	
	
	
	 public void BtnLoginOnClick(View view){ 
		 
		 login();
	    }   
	 
	 @SuppressLint("NewApi") private boolean validate(){
		 String username,password;
		 username=mUsernameEdit.getText().toString();
		 password=mPasswordEdit.getText().toString();
		 if(username.isEmpty()){
			  mAlert_Dialog.show("用户名不能为空!",  BitmapFactory.decodeResource(this.getResources(), R.drawable.wrong));
			  return false;
		 }
		 
		 if(password.isEmpty()){
			  mAlert_Dialog.show("密码不能为空!",  BitmapFactory.decodeResource(this.getResources(), R.drawable.wrong));
			  return false;
		 }
		 return true;
	 }
	
	 
	 private void login(){
		 if(validate()){
			 
			 mLoading_Dialog.startLoading();
			 
			final String username=mUsernameEdit.getText().toString();
			final String password=mPasswordEdit.getText().toString();			
			 HttpCommunicationUtil.sendDataToServer(ParamsUtil.getLogin(username, password), new CommunicationListener() {
				
				@SuppressLint("NewApi") @Override
				public void onResult(final String result) {
					String state = JsonUtil.getParam(result,"state");
					if(state.equals("success")){
						User  user =new User();
						user.setUsername(username);
						
						if(mAutoLoginCheckBox.isChecked()){
							user.setAutoLogin(true);
							user.setPassword(password);
						}
					SaveData.setUser(mContext, user);
				    SessonData.loginUser.setUsername(username);
				    SessonData.loginUser.setPassword(password);
					String user_json=JsonUtil.getParam(result,"user");
					SessonData.loginUser.setClientid(JsonUtil.getParam(user_json,"clientid"));
					SessonData.loginUser.setObject_id(JsonUtil.getParam(user_json,"objectId"));
					SessonData.loginUser.setLimit(JsonUtil.getParam(user_json,"limit"));
					
					if (SessonData.loginUser.getClientid() == null||SessonData.loginUser.getClientid().isEmpty()) {
						// clientid为空,跳转到完善信息页面
						
						runOnUiThread(new Runnable(){  
							  
		                    @Override  
		                    public void run() {  
		                        //更新UI  
		                    	 mLoading_Dialog.endLoading();	
								 Intent intent = new Intent(mContext,RegisterNextActivity.class);    
								 mContext.startActivity(intent);  
		                    }  
		                      
		                }); 
						   
						
						
					} else {
						InitData data = new InitData();
						data.mchntid = SessonData.loginUser.getClientid();// 商户号
						data.inscd =JsonUtil.getParam(user_json,"inscd");// 机构号
						data.signKey = JsonUtil.getParam(user_json,"signKey");// 秘钥
						// Log.e("opp",
						// ""+TelephonyManagerUtil.getDeviceId(mContext));
						data.terminalid = TelephonyManagerUtil
								.getDeviceId(mContext);// 设备号
						data.isProduce = SystemConfig.IS_PRODUCE;// 是否生产环境
						CashierSdk.init(data);
						runOnUiThread(new Runnable(){  
							  
		                    @Override  
		                    public void run() {  
		                        //更新UI  
		                    	 mLoading_Dialog.endLoading();
		                    	 SessonData.position_view=0;
		                    	 Intent intent = new Intent(LoginActivity.this,MainActivity.class);    
		                		 LoginActivity.this.startActivity(intent); 
		                    }  
		                      
		                }); 
						
					}
					
						
						
					}else{
						User user =new User();
						user.setUsername(username);
						user.setPassword(password);
						SaveData.setUser(mContext, user);
						SessonData.loginUser.setUsername(username);
					    SessonData.loginUser.setPassword(password);
						final String error=JsonUtil.getParam(result,"error");
						if(error.equals("user_no_activate")){
							runOnUiThread(new Runnable(){  
								  
			                    @Override  
			                    public void run() {  
			                        //更新UI  
			                    	 mLoading_Dialog.endLoading();	
			                    	 Activate_dialog activate_dialog=new Activate_dialog(mContext, LoginActivity.this.findViewById(R.id.activate_dialog), SessonData.loginUser.getUsername());
				                     activate_dialog.show();
			                    }  
			                      
			                }); 
						}else{
							runOnUiThread(new Runnable(){  
								  
			                    @Override  
			                    public void run() {  
			                        //更新UI  
			                    	String  errorStr=ErrorUtil.getErrorString(error);
			                    	 mLoading_Dialog.endLoading();		                    	 
			                    	 mAlert_Dialog.show(errorStr,  BitmapFactory.decodeResource(mContext.getResources(), R.drawable.wrong));
			                    }  
			                      
			                }); 
						}
										
					}
					
					
					   
					
				}
				
				@Override
				public void onError(final String error) {
					  runOnUiThread(new Runnable(){  
						  
		                    @Override  
		                    public void run() {  
		                        //更新UI  
		                    	 mLoading_Dialog.endLoading();	
		                    	 mAlert_Dialog.show(error,  BitmapFactory.decodeResource(mContext.getResources(), R.drawable.wrong));
		                    }  
		                      
		                });  
					
				}
			});
		 }
	 }
	 
	 public void BtnRegisterOnClick(View view){    
	        
		 Intent intent = new Intent(LoginActivity.this,RegisterActivity.class);    
		 LoginActivity.this.startActivity(intent);    
          
     }   
	 
	 @Override
	protected void onResume() {
		super.onResume();
		User user=SaveData.getUser(mContext);
		mAutoLoginCheckBox.setChecked(user.isAutoLogin());
		mUsernameEdit.setText(user.getUsername());
		mPasswordEdit.setText(user.getPassword());
		
	}
	 
	    
}
