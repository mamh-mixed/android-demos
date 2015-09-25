package com.cardinfolink.yunshouyin.activity;

import com.cardinfolink.cashiersdk.model.InitData;
import com.cardinfolink.cashiersdk.sdk.CashierSdk;
import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.constant.SystemConfig;
import com.cardinfolink.yunshouyin.data.SessonData;
import com.cardinfolink.yunshouyin.data.User;
import com.cardinfolink.yunshouyin.util.CommunicationListener;
import com.cardinfolink.yunshouyin.util.HttpCommunicationUtil;
import com.cardinfolink.yunshouyin.util.JsonUtil;
import com.cardinfolink.yunshouyin.util.ParamsUtil;
import com.cardinfolink.yunshouyin.util.TelephonyManagerUtil;
import com.cardinfolink.yunshouyin.util.VerifyUtil;
import com.cardinfolink.yunshouyin.view.Alert_Dialog;

import android.annotation.SuppressLint;
import android.app.Activity;
import android.content.Context;
import android.content.Intent;
import android.graphics.BitmapFactory;
import android.os.Bundle;
import android.view.View;
import android.widget.EditText;

public class RegisterNextActivity extends BaseActivity {
	private EditText mOpenBankEdit;
	private EditText mNameEdit;
	private EditText mBanknumEdit;
	private EditText mPhonenumEdit;

	@Override
	protected void onCreate(Bundle savedInstanceState) {
		super.onCreate(savedInstanceState);
		setContentView(R.layout.register_next_activity);
		mOpenBankEdit = (EditText) findViewById(R.id.info_openbank);
		mNameEdit = (EditText) findViewById(R.id.info_name);
		mBanknumEdit = (EditText) findViewById(R.id.info_banknum);
		mPhonenumEdit = (EditText) findViewById(R.id.info_phonenum);		
		VerifyUtil.bankCardNumAddSpace(mBanknumEdit);

	}

	public void BtnRegisterFinishedOnClick(View view) {
		if (validate()) {
			mLoading_Dialog.startLoading();
			User user = new User();
			user.setUsername(SessonData.loginUser.getUsername());
			user.setPassword(SessonData.loginUser.getPassword());
			user.setBank_open(mOpenBankEdit.getText().toString());
			user.setPayee(mNameEdit.getText().toString());
			user.setPayee_card(mBanknumEdit.getText().toString()
					.replace(" ", ""));
			user.setPhone_num(mPhonenumEdit.getText().toString());
			HttpCommunicationUtil.sendDataToServer(
					ParamsUtil.getImproveInfo(user),
					new CommunicationListener() {

						@Override
						public void onResult(String result) {
							String state = JsonUtil.getParam(result, "state");
							if (state.equals("success")) {
								String user_json=JsonUtil.getParam(result,"user");
								SessonData.loginUser.setClientid(JsonUtil.getParam(user_json,"clientid"));
								SessonData.loginUser.setObject_id(JsonUtil.getParam(user_json,"objectId"));
								SessonData.loginUser.setLimit(JsonUtil.getParam(user_json,"limit"));
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
								 Intent intent = new Intent(RegisterNextActivity.this,MainActivity.class);    
		                    	 RegisterNextActivity.this.startActivity(intent); 
		                    	 RegisterNextActivity.this.finish();
								

							} else {
								runOnUiThread(new Runnable() {

									@Override
									public void run() {
										// 更新UI
										mLoading_Dialog.endLoading();
										mAlert_Dialog.show(
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
							runOnUiThread(new Runnable() {

								@Override
								public void run() {
									// 更新UI
									mLoading_Dialog.endLoading();
									mAlert_Dialog.show(error, BitmapFactory
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
			mAlert_Dialog.show("开户行不能为空!", BitmapFactory.decodeResource(
					this.getResources(), R.drawable.wrong));
			return false;
		}

		if (name.isEmpty()) {
			mAlert_Dialog.show("姓名不能为空!", BitmapFactory.decodeResource(
					this.getResources(), R.drawable.wrong));
			return false;
		}

		if (banknum.isEmpty()) {
			mAlert_Dialog.show("银行卡号不能为空!", BitmapFactory.decodeResource(
					this.getResources(), R.drawable.wrong));
			return false;
		}

		if (!VerifyUtil.checkBankCard(banknum)) {
			mAlert_Dialog.show("请输入正确的银行卡号!", BitmapFactory.decodeResource(
					this.getResources(), R.drawable.wrong));
			return false;
		}

		if (phonenum.isEmpty()) {
			mAlert_Dialog.show("手机号不能为空!", BitmapFactory.decodeResource(
					this.getResources(), R.drawable.wrong));
			return false;
		}

		if (!VerifyUtil.isMobileNO(phonenum)) {
			mAlert_Dialog.show("请输入正确的手机号!", BitmapFactory.decodeResource(
					this.getResources(), R.drawable.wrong));
			return false;
		}

		return true;
	}

}
