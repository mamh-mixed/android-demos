package com.cardinfolink.yunshouyin.activity;
import java.text.SimpleDateFormat;
import java.util.Date;
import java.util.Map;

import org.json.JSONException;
import org.json.JSONObject;

import android.content.Intent;
import android.graphics.BitmapFactory;
import android.net.Uri;
import android.os.Bundle;
import android.os.Message;
import android.util.Log;
import android.view.Gravity;
import android.view.KeyEvent;


import android.widget.Toast;

import com.cardinfo.framelib.activity.BaseHtml5Activity;
import com.cardinfo.framelib.constant.Msg;
import com.cardinfo.framelib.listener.CommunicationListener;
import com.cardinfo.framelib.model.JavaJSParam;
import com.cardinfo.framelib.util.HttpCommunicationUtil;
import com.cardinfo.framelib.util.JsonUtil;
import com.cardinfo.framelib.util.TelephonyManagerUtil;
import com.cardinfo.framelib.version.VersionManage;
import com.cardinfolink.cashiersdk.model.InitData;
import com.cardinfolink.cashiersdk.sdk.CashierSdk;
import com.cardinfolink.cashiersdk.util.MapUtil;
import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.constant.SystemConfig;
import com.cardinfolink.yunshouyin.data.CurrentState;
import com.cardinfolink.yunshouyin.data.SaveData;
import com.cardinfolink.yunshouyin.data.User;
import com.cardinfolink.yunshouyin.util.ParamsUtil;
import com.cardinfolink.yunshouyin.view.Activate_dialog;
import com.cardinfolink.yunshouyin.view.Alert_Dialog;
import com.cardinfolink.yunshouyin.view.Refd_Dialog;

public class MainActivity extends BaseHtml5Activity {

	private static final String TAG = "MainActivity";

	private long exitTime = 0;

	private boolean isRegister = false;
	
	private String today;
	


	@Override
	protected void onCreate(Bundle savedInstanceState) {
		super.onCreate(savedInstanceState);
		 versionManage=new VersionManage(mContext, mHandler,SystemConfig.appserverurl);
		SimpleDateFormat spf = new SimpleDateFormat("yyyyMMdd");
		
		today = spf.format(new Date());
		initLayout();
	}

	private void initLayout() {
		
		
		mWebView.post(new Runnable() {

			@Override
			public void run() {
				mWebView.loadUrl(SystemConfig.HTML_HEAD_ADDRESS
						+ "/home-page.html");

			}
		});

		
		//进行更新检测，和更新相关
		versionManage.update();
	
		
		//Message msg = mHandler.obtainMessage(Msg.MSG_FROM_CLIENT_JUMP_LOGIN);
		//mHandler.sendMessageDelayed(msg, 3000);

	}

	@Override
	public void ProcessingMessage(Message msg) {
		super.ProcessingMessage(msg);
		switch (msg.what) {
		case com.cardinfo.framelib.constant.Msg.MSG_FROM_JS:
			JavaJSParam jsParam = (JavaJSParam) msg.obj;
			String method = jsParam.getMethod();
			if (method.equals("login")) {
				String json = jsParam.getParam();
				String username = JsonUtil.getParam(json, "username");
				String password = JsonUtil.getParam(json, "password");
				boolean isAutoLogin = JsonUtil.getParam(json, "autologin")
						.equals("true");
				login(username, password, isAutoLogin);
			} else if (method.equals("safeexit")) {
				 final User user = SaveData.getUser(mContext);			    
				user.setPassword("");
				SaveData.setUser(mContext, user);
			    final boolean	isAuotLogin=false;
				SaveData.setAutoLogin(mContext, false);
				mWebView.post(new Runnable() {

					@Override
					public void run() {
						mWebView.loadUrl(SystemConfig.HTML_HEAD_ADDRESS
								+ "/login.html?device=android&username="
								+ user.getUsername() + "&password="
								+ user.getPassword() + "&autologin="
								+ isAuotLogin);

					}
				});

			} else if (method.equals("jump_register")) {
				isRegister = true;
				mWebView.post(new Runnable() {

					@Override
					public void run() {
						mWebView.loadUrl(SystemConfig.HTML_HEAD_ADDRESS
								+ "/register.html?device=android");

					}
				});

			} else if (method.equals("improveinfo")) {
				String json = jsParam.getParam();
				Log.i(TAG, "improveinfo");				
				improveinfo(json);
			}else if (method.equals("register")) {
				String json = jsParam.getParam();
				Log.i(TAG, "register");
				final String username = JsonUtil.getParam(json, "username");
				String password = JsonUtil.getParam(json, "password");
				register(username,password);
			}else if (method.equals("scancode")) {
				final String json = jsParam.getParam();
				Log.i(TAG, "scancode");
				User user=SaveData.getUser(mContext);
				if (!user.getLimit().equals("false")) {
				HttpCommunicationUtil.sendDataToServer(ParamsUtil.getTotal(user, today), new CommunicationListener() {
					
					@Override
					public void onResult(String result) {
						String state = JsonUtil.getParam(result,
								"state");

						if (state.equals("success")) {

							double limitValue = Double
									.parseDouble(JsonUtil.getParam(
											result, "total"));
							if (limitValue >= 500) {
								Message msg = mHandler
										.obtainMessage(Msg.MSG_FROM_SERVER_LIMIT_BEYOND);
								mHandler.sendMessageDelayed(msg, 0);
							} else {
								Double total=Double.parseDouble(JsonUtil.getParam(json, "sum"));
								if (total > 0) {
									
									String chcd = JsonUtil.getParam(json, "chcd");
									String busicd= JsonUtil.getParam(json, "busicd");
									if(busicd.equals("PURC")){
										Intent intent = new Intent(mContext,
												CaptureActivity.class);
										intent.putExtra("chcd", chcd);
										intent.putExtra("total", "" + total);
										startActivityForResult(intent, 100);
									}else{
										Log.i(TAG, "PAUT");
										Intent intent = new Intent(mContext,
												CreateQRcodeActivity.class);
										intent.putExtra("chcd", chcd);
										intent.putExtra("total", "" + total);
										startActivityForResult(intent, 100);
									}
									
								} else {
									Toast toast = Toast.makeText(getApplicationContext(),
											"金额不能为0", Toast.LENGTH_SHORT);
									toast.show();
								}

							}

						}
						
					}
					
					@Override
					public void onError(String error) {
						mHandler.sendEmptyMessage(Msg.MSG_FROM_SERVER_TIMEOUT);
					}
				});
				}else{
					Double total=Double.parseDouble(JsonUtil.getParam(json, "sum"));
					if (total > 0) {
						
						String chcd = JsonUtil.getParam(json, "chcd");
						String busicd= JsonUtil.getParam(json, "busicd");
						if(busicd.equals("PURC")){
							Intent intent = new Intent(mContext,
									CaptureActivity.class);
							intent.putExtra("chcd", chcd);
							intent.putExtra("total", "" + total);
							startActivityForResult(intent, 100);
						}else{
							Log.i(TAG, "PAUT");
							Intent intent = new Intent(mContext,
									CreateQRcodeActivity.class);
							intent.putExtra("chcd", chcd);
							intent.putExtra("total", "" + total);
							startActivityForResult(intent, 100);
						}
						
					} else {
						Toast toast = Toast.makeText(getApplicationContext(),
								"金额不能为0", Toast.LENGTH_SHORT);
						toast.show();
					}
				}
					
			
			}else if (method.equals("updatepassword")) {
				String json = jsParam.getParam();
				Log.i(TAG, "register");
				final String oldpwd = JsonUtil.getParam(json, "oldpwd");
				final String newpwd = JsonUtil.getParam(json, "newpwd");
				final User user=SaveData.getUser(mContext);
				HttpCommunicationUtil.sendDataToServer(ParamsUtil.getUpdate(user.getUsername(),oldpwd, newpwd), new CommunicationListener() {
					
					@Override
					public void onResult(String result) {
						Log.i(TAG, result);
						String state=JsonUtil.getParam(result, "state");
						if(state.equals("success")){
							user.setPassword(newpwd);
							SaveData.setUser(mContext, user);
							mHandler.sendEmptyMessage(Msg.MSG_FROM_SERVER_UPDATEPASSWORD_SUCCESS);
							
						}else{
						String error = JsonUtil.getParam(
								result, "error");
						Message msg = mHandler
								.obtainMessage(Msg.MSG_FROM_SERVER_ERROR);
						msg.obj = error;
						mHandler.sendMessageDelayed(msg, 0);
						}
					}
					
					@Override
					public void onError(String error) {
						mHandler.sendEmptyMessage(Msg.MSG_FROM_SERVER_TIMEOUT);
					}
				});
				
			}else if (method.equals("updateaccount")) {
				String json = jsParam.getParam();
				Log.i(TAG, "updateaccount");
				User user=SaveData.getUser(mContext);
				user.setBank_open(JsonUtil.getParam(json, "bank_open"));
				user.setPayee(JsonUtil.getParam(json, "payee"));
				user.setPayee_card(JsonUtil.getParam(json, "payee_card"));
				user.setPhone_num(JsonUtil.getParam(json, "phone_num"));
				HttpCommunicationUtil.sendDataToServer(ParamsUtil.getUpdateInfo(user), new CommunicationListener() {
					
					@Override
					public void onResult(String result) {
						Log.i(TAG, result);
						String state=JsonUtil.getParam(result, "state");
						if(state.equals("success")){
							mHandler.sendEmptyMessage(Msg.MSG_FROM_SERVER_UPDATEPASSWORD_SUCCESS);
						}else{
						String error = JsonUtil.getParam(
								result, "error");
						Message msg = mHandler
								.obtainMessage(Msg.MSG_FROM_SERVER_ERROR);
						msg.obj = error;
						mHandler.sendMessageDelayed(msg, 0);
						}
					}
					
					@Override
					public void onError(String error) {
						mHandler.sendEmptyMessage(Msg.MSG_FROM_SERVER_TIMEOUT);
					}
				});
				
			}else if (method.equals("limitincrease")) {
				String json = jsParam.getParam();
				Log.i(TAG, "limitincrease");
				User user=SaveData.getUser(mContext);
				user.setPayee(JsonUtil.getParam(json, "payee"));
				 user.setEmail(JsonUtil.getParam(json, "email"));
				user.setPhone_num(JsonUtil.getParam(json, "phone_num"));
				HttpCommunicationUtil.sendDataToServer(ParamsUtil.getLimitincrease(user), new CommunicationListener() {
					
					@Override
					public void onResult(String result) {
						Log.i(TAG, result);
						String state=JsonUtil.getParam(result, "state");
						if(state.equals("success")){
							mHandler.sendEmptyMessage(Msg.MSG_FROM_SERVER_LIMITINCREASE_SUCCESS);
						}else{
						String error = JsonUtil.getParam(
								result, "error");
						Message msg = mHandler
								.obtainMessage(Msg.MSG_FROM_SERVER_ERROR);
						msg.obj = error;
						mHandler.sendMessageDelayed(msg, 0);
						}
					}
					
					@Override
					public void onError(String error) {
						mHandler.sendEmptyMessage(Msg.MSG_FROM_SERVER_TIMEOUT);
					}
				});
				
			}else if (method.equals("refd")) {
				String json = jsParam.getParam();
				Log.i(TAG, "refd");
				User user=SaveData.getUser(mContext);
				final String orderNum=JsonUtil.getParam(json, "orderNum");
				final String total=JsonUtil.getParam(json, "total");
				HttpCommunicationUtil.sendDataToServer(ParamsUtil.getRefd(user, orderNum), new CommunicationListener() {
					
					@Override
					public void onResult(String result) {
						Log.i(TAG, result);
						String state=JsonUtil.getParam(result, "state");
						if(state.equals("success")){
							Message msg=mHandler.obtainMessage(Msg.MSG_FROM_SERVER_QYREFD_SUCCESS);
							String refdtotal=JsonUtil.getParam(result, "refdtotal");
							JSONObject json=new JSONObject();
							try {
								json.put("refdtotal", refdtotal);
								json.put("total", total);
								json.put("orderNum", orderNum);
								msg.obj=json.toString();
								mHandler.sendMessageDelayed(msg, 0);
							} catch (JSONException e) {
								// TODO Auto-generated catch block
								e.printStackTrace();
							}
							
							
						}else{
							String error = JsonUtil.getParam(
									result, "error");
							Message msg = mHandler
									.obtainMessage(Msg.MSG_FROM_SERVER_ERROR);
							msg.obj = error;
							mHandler.sendMessageDelayed(msg, 0);
							}
						
						
					}
					
					@Override
					public void onError(String error) {
						mHandler.sendEmptyMessage(Msg.MSG_FROM_SERVER_TIMEOUT);
						
					}
				});
			}else if (method.equals("openwapbill")) {
				
				   User user=SaveData.getUser(mContext);
				   Uri uri;
				 uri = Uri.parse(SystemConfig.WEB_BILL_URL+user.getObject_id());
				 Intent  intent = new  Intent(Intent.ACTION_VIEW, uri);
				  startActivity(intent);
				   
				
				break;
			}


			break;

		case Msg.MSG_FROM_SERVER_TIMEOUT: {
			Toast toast = Toast.makeText(getApplicationContext(),
					getResources().getString(R.string.server_timeout),
					Toast.LENGTH_SHORT);
			toast.setGravity(Gravity.CENTER, 0, 250);
			toast.show();
			final User user = SaveData.getUser(mContext);
			final boolean isAuotLogin = SaveData.isAutoLogin(mContext);
			mWebView.post(new Runnable() {

				@Override
				public void run() {
					mWebView.loadUrl(SystemConfig.HTML_HEAD_ADDRESS
							+ "/login.html?device=android&username="
							+ user.getUsername() + "&password="
							+ user.getPassword() + "&autologin=" + isAuotLogin);

				}
			});
			break;
		}
		case Msg.MSG_FROM_SERVER_ERROR: {
			String error = (String) msg.obj;
			if (error.equals("username_password_error")) {
//				Toast toast = Toast.makeText(
//						getApplicationContext(),
//						getResources().getString(R.string.username_password_error),
//						Toast.LENGTH_SHORT);
//				toast.setGravity(Gravity.CENTER, 0, 250);
//				toast.show();
			   Alert_Dialog alert_Dialog=new Alert_Dialog(mContext, mHandler, findViewById(R.id.alert_dialog), 
					   getResources().getString(R.string.username_password_error), BitmapFactory.decodeResource(this.getResources(), R.drawable.wrong));
			   alert_Dialog.show();
			   final User user = SaveData.getUser(mContext);
			    user.setPassword("");
			    SaveData.setUser(mContext, user);
			    SaveData.setAutoLogin(mContext, false);
				final boolean isAuotLogin = SaveData.isAutoLogin(mContext);
				mWebView.post(new Runnable() {

					@Override
					public void run() {
						mWebView.loadUrl(SystemConfig.HTML_HEAD_ADDRESS
								+ "/login.html?device=android&username="
								+ user.getUsername() + "&password="
								+ user.getPassword() + "&autologin="
								+ isAuotLogin);

					}
				});
				
				
				
				
			} else if (error.equals("username_exist")) {
				Toast toast = Toast.makeText(getApplicationContext(),
						getResources().getString(R.string.username_exist),
						Toast.LENGTH_SHORT);
				toast.setGravity(Gravity.CENTER, 0, 250);
				toast.show();
			}else if(error.equals("user_no_activate")){
				User user = SaveData.getUser(mContext);
				Activate_dialog activate_dialog=new Activate_dialog(mContext,mHandler,findViewById(R.id.activate_dialog),user.getUsername());
				activate_dialog.show();
				
			}else if(error.equals("username_no_exist")){
			
				Alert_Dialog alert_Dialog=new Alert_Dialog(mContext, mHandler, findViewById(R.id.alert_dialog), 
						  "用户名不存在!", BitmapFactory.decodeResource(this.getResources(), R.drawable.wrong));
				   alert_Dialog.show();
				
			}else if(error.equals("old_password_error")){
			
				Alert_Dialog alert_Dialog=new Alert_Dialog(mContext, mHandler, findViewById(R.id.alert_dialog), 
						  "原密码错误!", BitmapFactory.decodeResource(this.getResources(), R.drawable.wrong));
				   alert_Dialog.show();
				
			}else if(error.equals("username_no_exist")){
			
				Alert_Dialog alert_Dialog=new Alert_Dialog(mContext, mHandler, findViewById(R.id.alert_dialog), 
						  "用户名不存在!", BitmapFactory.decodeResource(this.getResources(), R.drawable.wrong));
				   alert_Dialog.show();
				
			}
			break;
		}
		case Msg.MSG_FROM_CLIENT_JUMP_LOGIN: {
			final User user = SaveData.getUser(mContext);
			final boolean isAuotLogin = SaveData.isAutoLogin(mContext);
			if (isAuotLogin && !CurrentState.isSafeExit) {
				login(user.getUsername(), user.getPassword(), isAuotLogin);
				return;
			}

			mWebView.post(new Runnable() {

				@Override
				public void run() {
					mWebView.loadUrl(SystemConfig.HTML_HEAD_ADDRESS
							+ "/login.html?device=android&username="
							+ user.getUsername() + "&password="
							+ user.getPassword() + "&autologin=" + isAuotLogin);

				}
			});
			break;
		}

		case Msg.MSG_FROM_CLIENT_REGISTER_SUCCESS: {
			
			User user=SaveData.getUser(mContext);
			login(user.getUsername(), user.getPassword(), false);
			break;
		}

		case Msg.MSG_FROM_ACTIVATE_DIGLOG_OK: {
			User user=SaveData.getUser(mContext);
			requestActivate(user.getUsername(), user.getPassword());
			
			break;
		}
		
		case Msg.MSG_FROM_CLIENT_JUMP_NEXT_REGISTER: {
			mWebView.post(new Runnable() {

				@Override
				public void run() {
					mWebView.loadUrl(SystemConfig.HTML_HEAD_ADDRESS
							+ "/register-next.html?device=android");

				}
			});
			
			break;
		}
		
		
		case Msg.MSG_FROM_CLIENT_JUMP_SCANCODE: {
			final User user=SaveData.getUser(mContext);
			mWebView.post(new Runnable() {

				@Override
				public void run() {
					mWebView.loadUrl(SystemConfig.HTML_HEAD_ADDRESS
							+ "/index.html?device=android&username="
							+ user.getUsername() + "&password="
							+ user.getPassword() +"&clientid="+user.getClientid()+"&key="+SystemConfig.APP_KEY+"#/scanPage");

				}
			});
			
			break;
		}
		
		case Msg.MSG_FROM_SERVER_UPDATEPASSWORD_SUCCESS: {
			User user=SaveData.getUser(mContext);
			JSONObject jsonObject=new JSONObject();
			try {
				jsonObject.put("username", user.getUsername());
				jsonObject.put("password", user.getPassword());
				jsonObject.put("device", "android");
				jsonObject.put("key", SystemConfig.APP_KEY);
				jsonObject.put("clientid", user.getClientid());
				
				mWebView.loadUrl("javascript:CloudCashierBridge.saveUserData("+jsonObject+")");  
				
				 Alert_Dialog alert_Dialog=new Alert_Dialog(mContext, mHandler, findViewById(R.id.alert_dialog), 
						  "修改成功!", BitmapFactory.decodeResource(this.getResources(), R.drawable.right));
				   alert_Dialog.show();
			} catch (JSONException e) {
				// TODO Auto-generated catch block
				e.printStackTrace();
			}
			
			
			break;
		}
		
		case Msg.MSG_FROM_SERVER_QYREFD_SUCCESS: {
			 String json=(String) msg.obj;
			 Log.i(TAG, json);
			 Refd_Dialog refd_Dialog=new Refd_Dialog(mContext, mHandler, findViewById(R.id.refd_dialog),JsonUtil.getParam(json, "orderNum"), JsonUtil.getParam(json, "refdtotal"), JsonUtil.getParam(json, "total"));
			 refd_Dialog.show();
			break;
		}
		
		case Msg.MSG_FROM_SERVER_LIMITINCREASE_SUCCESS: {
			 Alert_Dialog alert_Dialog=new Alert_Dialog(mContext, mHandler, findViewById(R.id.alert_dialog), 
					  "提交成功!", BitmapFactory.decodeResource(this.getResources(), R.drawable.right));
			   alert_Dialog.show();
			
			break;
		}case Msg.MSG_FROM_SERVER_REFD_SUCCESS: {
			 Alert_Dialog alert_Dialog=new Alert_Dialog(mContext, mHandler, findViewById(R.id.alert_dialog), 
					  "退款成功!", BitmapFactory.decodeResource(this.getResources(), R.drawable.right));
			   alert_Dialog.show();
			
			break;
		}case Msg.MSG_FROM_SERVER_REFD_FAIL: {
			 Alert_Dialog alert_Dialog=new Alert_Dialog(mContext, mHandler, findViewById(R.id.alert_dialog), 
					  "退款失败!", BitmapFactory.decodeResource(this.getResources(), R.drawable.wrong));
			   alert_Dialog.show();
			
			break;
		}
		default:
			break;
		}
	}

	
	private void register(final String username,final String password){
		if (isRegister) {
			HttpCommunicationUtil.sendDataToServer(
					ParamsUtil.getRegister(username, password),
					new CommunicationListener() {

						@Override
						public void onResult(String result) {
							Log.i(TAG, result);
							String state = JsonUtil.getParam(result,
									"state");
							if (state.equals("success")) {
								isRegister=false;
								User user = new User();
								user.setUsername(username);
								user.setPassword(password);
								SaveData.setUser(mContext, user);
								SaveData.setAutoLogin(mContext, false);
								mHandler.sendEmptyMessage(Msg.MSG_FROM_CLIENT_REGISTER_SUCCESS);
							} else {
								String error = JsonUtil.getParam(
										result, "error");
								Message msg = mHandler
										.obtainMessage(Msg.MSG_FROM_SERVER_ERROR);
								msg.obj = error;
								mHandler.sendMessageDelayed(msg, 0);
							}

						}

						@Override
						public void onError(String error) {
							mHandler.sendEmptyMessage(Msg.MSG_FROM_SERVER_TIMEOUT);
						}
					});

		} else {
			User user=SaveData.getUser(mContext);
			login(user.getUsername(), user.getPassword(), false);
		}
	}
	
	
	private void login(final String username, final String password,
			final boolean isAutoLogin) {
		
		
		   mHandler.sendEmptyMessage(Msg.MSG_FROM_CLIENT_LOAD_START);
		HttpCommunicationUtil.sendDataToServer(
				ParamsUtil.getLogin(username, password),
				new CommunicationListener() {

					@Override
					public void onResult(String result) {
						Log.i(TAG, result);
						 mHandler.sendEmptyMessage(Msg.MSG_FROM_CLIENT_LOAD_END);
						Map<String, Object> map = MapUtil.getMapForJson(result);
						String state = (String) map.get("state");
						if (state.equals("success")) {
							CurrentState.isLogin = true;
							Log.i(TAG, "" + isAutoLogin);
							SaveData.setAutoLogin(mContext, isAutoLogin);
							Map<String, Object> usermap = MapUtil
									.getMapForJson(map.get("user").toString());
							final User user = new User();
							user.setUsername(username);
							user.setPassword(password);					
							String clientid = (String) usermap.get("clientid");
							user.setClientid(clientid);
							user.setObject_id((String)usermap.get("objectId"));
							String limit = (String) usermap.get("limit");// 是否限制日交易金额
							if (limit != null) {
								user.setLimit(limit);
							}

							SaveData.setUser(mContext, user);
							if (clientid == null) {
								// clientid为空,跳转到完善信息页面
								// Intent intent=new Intent(mContext,
								// ImproveInfoActivity.class);
								// mContext.startActivity(intent);
								mHandler.sendEmptyMessage(Msg.MSG_FROM_CLIENT_JUMP_NEXT_REGISTER);
							} else {
								InitData data = new InitData();
								data.mchntid = clientid;// 商户号
								data.inscd = (String) usermap.get("inscd");// 机构号
								data.signKey = (String) usermap.get("signKey");// 秘钥
								// Log.e("opp",
								// ""+TelephonyManagerUtil.getDeviceId(mContext));
								data.terminalid = TelephonyManagerUtil
										.getDeviceId(mContext);// 设备号
								data.isProduce = SystemConfig.IS_PRODUCE;// 是否生产环境
								CashierSdk.init(data);
								CurrentState.isSafeExit = false;// 是否安全退出

								// 启动扫码界面
								mWebView.post(new Runnable() {

									@Override
									public void run() {
										mWebView.loadUrl(SystemConfig.HTML_HEAD_ADDRESS
												+ "/index.html?device=android&username="
												+ user.getUsername() + "&password="
												+ user.getPassword() +"&clientid="+user.getClientid()+"&key="+SystemConfig.APP_KEY+"#/scanPage");

									}
								});
							}
						} else {
							
							User user = new User();
							user.setUsername(username);
							user.setPassword(password);
							SaveData.setUser(mContext, user);
							SaveData.setAutoLogin(mContext, false);
							String error = (String) map.get("error");
							Message msg = mHandler
									.obtainMessage(Msg.MSG_FROM_SERVER_ERROR);
							msg.obj = error;
							mHandler.sendMessageDelayed(msg, 0);

						}

					}

					@Override
					public void onError(String error) {
					    mHandler.sendEmptyMessage(Msg.MSG_FROM_CLIENT_LOAD_END);
						mHandler.sendEmptyMessage(Msg.MSG_FROM_SERVER_TIMEOUT);

					}
				});
	}

	
	private void requestActivate(final String username, final String password){
		Log.i(TAG, "requestActivate  username="+username+" password="+password);
		
		HttpCommunicationUtil.sendDataToServer(ParamsUtil.getRequestActivate(username, password), new CommunicationListener() {
			
			@Override
			public void onResult(String result) {
				Log.i(TAG, result);
				String state=JsonUtil.getParam(result, "state");
				if (state.equals("success")) {
					
				}else{
					String error =  JsonUtil.getParam(result, "error");
					Message msg = mHandler
							.obtainMessage(Msg.MSG_FROM_SERVER_ERROR);
					msg.obj = error;
					mHandler.sendMessageDelayed(msg, 0);
				}
			}
			
			@Override
			public void onError(String error) {
				mHandler.sendEmptyMessage(Msg.MSG_FROM_SERVER_TIMEOUT);
				
			}
		});
	}
	
	private void improveinfo(String json){
		final User user=SaveData.getUser(mContext);
		user.setBank_open(JsonUtil.getParam(json, "bank_open"));
		user.setPayee(JsonUtil.getParam(json, "payee"));
		user.setPayee_card(JsonUtil.getParam(json, "payee_card"));
		user.setPhone_num(JsonUtil.getParam(json, "phone_num"));
		HttpCommunicationUtil.sendDataToServer(ParamsUtil.getImproveInfo(user), new CommunicationListener() {
			
			@Override
			public void onResult(String result) {
				Log.i(TAG, result);
				Map<String, Object>map= MapUtil.getMapForJson(result);
				String state=(String) map.get("state");
				if(state.equals("success")){
					Map<String, Object>usermap=MapUtil.getMapForJson(map.get("user").toString());
					 InitData data=new InitData();
					 data.mchntid=(String) usermap.get("clientid");
					 data.inscd=(String) usermap.get("inscd");
					 data.signKey=(String) usermap.get("signKey"); 
					 data.terminalid = TelephonyManagerUtil
								.getDeviceId(mContext);// 设备号
					 data.isProduce = SystemConfig.IS_PRODUCE;// 是否生产环境
					 CashierSdk.init(data);
					 user.setClientid( data.mchntid);
					 SaveData.setUser(mContext, user);
					 Message msg = mHandler
								.obtainMessage(Msg.MSG_FROM_CLIENT_JUMP_SCANCODE);
						mHandler.sendMessageDelayed(msg, 0);
				}else{
					String error =  JsonUtil.getParam(result, "error");
					Message msg = mHandler
							.obtainMessage(Msg.MSG_FROM_SERVER_ERROR);
					msg.obj = error;
					mHandler.sendMessageDelayed(msg, 0);
				}
			}
			
			@Override
			public void onError(String error) {
				mHandler.sendEmptyMessage(Msg.MSG_FROM_SERVER_TIMEOUT);
			}
		});
	}
	
	
	
	@Override
	public boolean onKeyDown(int keyCode, KeyEvent event) {
		if (keyCode == KeyEvent.KEYCODE_BACK
				&& event.getAction() == KeyEvent.ACTION_DOWN) {
			if ((System.currentTimeMillis() - exitTime) > 2000) {
				Toast toast = Toast.makeText(getApplicationContext(),
						getResources().getString(R.string.press_again_exit),
						Toast.LENGTH_SHORT);
				toast.setGravity(Gravity.CENTER, 0, 250);
				toast.show();
				exitTime = System.currentTimeMillis();
			} else {
				finish();
				System.exit(0);
			}
			return true;
		}
		return super.onKeyDown(keyCode, event);
	}

	@Override
	protected void onActivityResult(int requestCode, int resultCode, Intent data) {
		super.onActivityResult(requestCode, resultCode, data);
		
		if(requestCode==100){
			if(resultCode==101){
			final User user=SaveData.getUser(mContext);
			mWebView.post(new Runnable() {

				@Override
				public void run() {
					mWebView.loadUrl(SystemConfig.HTML_HEAD_ADDRESS
							+ "/index.html?device=android&username="
							+ user.getUsername() + "&password="
							+ user.getPassword() +"&clientid="+user.getClientid()+"&key="+SystemConfig.APP_KEY+"#/transManage");

				}
			});
		}
			
		}else if(requestCode==103){
			 mHandler.sendEmptyMessage(Msg.MSG_FROM_CLIENT_JUMP_LOGIN);
		}
	}
	
	
}
