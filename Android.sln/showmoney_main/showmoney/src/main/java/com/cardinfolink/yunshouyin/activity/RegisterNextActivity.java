package com.cardinfolink.yunshouyin.activity;

import java.util.ArrayList;
import java.util.Iterator;
import java.util.List;

import org.json.JSONArray;
import org.json.JSONException;
import org.json.JSONObject;

import com.cardinfolink.cashiersdk.model.InitData;
import com.cardinfolink.cashiersdk.sdk.CashierSdk;
import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.constant.SystemConfig;
import com.cardinfolink.yunshouyin.data.SessonData;
import com.cardinfolink.yunshouyin.data.User;
import com.cardinfolink.yunshouyin.util.BankBaseUtil;
import com.cardinfolink.yunshouyin.util.CommunicationListener;
import com.cardinfolink.yunshouyin.util.HttpCommunicationUtil;
import com.cardinfolink.yunshouyin.util.JsonUtil;
import com.cardinfolink.yunshouyin.util.ParamsUtil;
import com.cardinfolink.yunshouyin.util.TelephonyManagerUtil;
import com.cardinfolink.yunshouyin.util.VerifyUtil;
import com.cardinfolink.yunshouyin.view.SearchAdapter;

import android.annotation.SuppressLint;
import android.app.Activity;
import android.content.Context;
import android.content.Intent;
import android.graphics.BitmapFactory;
import android.os.Bundle;
import android.text.Editable;
import android.text.TextWatcher;
import android.util.Log;
import android.view.View;
import android.view.View.OnFocusChangeListener;
import android.widget.AdapterView;
import android.widget.AdapterView.OnItemSelectedListener;
import android.widget.ArrayAdapter;
import android.widget.AutoCompleteTextView;
import android.widget.EditText;
import android.widget.Spinner;

public class RegisterNextActivity extends BaseActivity {

	private EditText mNameEdit;
	private EditText mBanknumEdit;
	private EditText mPhonenumEdit;

	private AutoCompleteTextView mProvinceEdit;
	private Spinner mProvinceSpinner;
	private List<String> mProvinceList;
	private ArrayAdapter mProvinceAdapter;
	private SearchAdapter mProvinceSearchAdapter;

	private AutoCompleteTextView mCityEdit;
	private Spinner mCitySpinner;
	private List<String> mCityList;
	private List<String> mBankIdList;
	private ArrayAdapter mCityAdapter;
	private SearchAdapter mCitySearchAdapter;

	private AutoCompleteTextView mOpenBankEdit;
	private Spinner mOpenBankSpinner;
	private List<String> mOpenBankList;
	private List<String> mCityCodeList;
	private ArrayAdapter mOpenBankAdapter;
	private SearchAdapter mOpenBankSearchAdapter;

	private AutoCompleteTextView mBranchBankEdit;
	private Spinner mBranchBankSpinner;
	private List<String> mBranchBankList;
	private List<String> mBankNoList;
	private ArrayAdapter mBranchBankAdapter;
	private SearchAdapter mBranchBankSearchAdapter;

	private String info_province;
	private String info_city;
	private String info_openbank;
	private String info_branch_bank;
	


	@Override
	protected void onCreate(Bundle savedInstanceState) {
		super.onCreate(savedInstanceState);
		setContentView(R.layout.register_next_activity);
		initLayout();
		initListener();
		initData();
	}

	private void initLayout() {
		
		mNameEdit = (EditText) findViewById(R.id.info_name);
		mBanknumEdit = (EditText) findViewById(R.id.info_banknum);
		mPhonenumEdit = (EditText) findViewById(R.id.info_phonenum);
		VerifyUtil.bankCardNumAddSpace(mBanknumEdit);

		VerifyUtil.bankCardNumAddSpace(mBanknumEdit);

		mProvinceEdit = (AutoCompleteTextView) findViewById(R.id.edit_province);
		mProvinceSpinner = (Spinner) findViewById(R.id.spinner_province);
		// 适配器
		mProvinceList = new ArrayList<String>();
		mProvinceList.add("开户行所在省份");
		mProvinceAdapter = new ArrayAdapter<String>(mContext,
				R.layout.spinner_item, mProvinceList);

		// 设置样式
		mProvinceAdapter.setDropDownViewResource(R.layout.spinner_drop_item);
		// 加载适配器
		mProvinceSpinner.setAdapter(mProvinceAdapter);

		mProvinceSearchAdapter = new SearchAdapter(mContext, mProvinceList);

		mProvinceEdit.setAdapter(mProvinceSearchAdapter);
		mProvinceEdit.setThreshold(1);
		// mProvinceEdit.setf

		mCityEdit = (AutoCompleteTextView)findViewById(R.id.edit_city);
		mCitySpinner = (Spinner) findViewById(R.id.spinner_city);
		// 适配器
		mCityList = new ArrayList<String>();
		mCityCodeList = new ArrayList<String>();
		mCityList.add("开户行所在城市");
		mCityCodeList.add("");
		mCityAdapter = new ArrayAdapter<String>(mContext,
				R.layout.spinner_item, mCityList);
		// 设置样式
		mCityAdapter.setDropDownViewResource(R.layout.spinner_drop_item);
		// 加载适配器
		mCitySpinner.setAdapter(mCityAdapter);

		mCitySearchAdapter = new SearchAdapter(mContext, mCityList);

		mCityEdit.setAdapter(mCitySearchAdapter);
		mCityEdit.setThreshold(1);

		mOpenBankEdit = (AutoCompleteTextView) findViewById(R.id.edit_openbank);
		mOpenBankSpinner = (Spinner) findViewById(R.id.spinner_openbank);
		// 适配器
		mOpenBankList = new ArrayList<String>();
		mOpenBankList.add("请选择开户银行");
		mBankIdList = new ArrayList<String>();
		mBankIdList.add("");

		mOpenBankAdapter = new ArrayAdapter<String>(mContext,
				R.layout.spinner_item, mOpenBankList);
		// 设置样式
		mOpenBankAdapter.setDropDownViewResource(R.layout.spinner_drop_item);
		// 加载适配器
		mOpenBankSpinner.setAdapter(mOpenBankAdapter);

		mOpenBankSearchAdapter = new SearchAdapter(mContext, mOpenBankList);
		mOpenBankEdit.setAdapter(mOpenBankSearchAdapter);
		mOpenBankEdit.setThreshold(1);

		mBranchBankEdit = (AutoCompleteTextView) findViewById(R.id.edit_branchbank);
		mBranchBankSpinner = (Spinner) findViewById(R.id.spinner_branchbank);
		// 适配器
		mBranchBankList = new ArrayList<String>();
		mBranchBankList.add("请选择开户支行");
		mBankNoList = new ArrayList<String>();
		mBankNoList.add("行号");

		mBranchBankAdapter = new ArrayAdapter<String>(mContext,
				R.layout.spinner_item, mBranchBankList);
		// 设置样式
		mBranchBankAdapter.setDropDownViewResource(R.layout.spinner_drop_item);
		// 加载适配器
		mBranchBankSpinner.setAdapter(mBranchBankAdapter);

		mBranchBankSearchAdapter = new SearchAdapter(mContext, mBranchBankList);
		mBranchBankEdit.setAdapter(mBranchBankSearchAdapter);
		mBranchBankEdit.setThreshold(1);

	}

	
	public void initData() {

		HttpCommunicationUtil.sendGetDataToServer(BankBaseUtil.getProvince(),
				new CommunicationListener() {

					@Override
					public void onResult(String result) {

						try {
							JSONArray jsonArray = new JSONArray(result);
							mProvinceList.clear();
							mProvinceList.add("开户行所在省份");
							for (int i = 0; i < jsonArray.length(); i++) {
								mProvinceList.add(jsonArray.getString(i));
							}

						} catch (JSONException e) {
							// TODO Auto-generated catch block
							e.printStackTrace();
						}

						((Activity) mContext).runOnUiThread(new Runnable() {

							@Override
							public void run() {
								// 更新UI
								mProvinceAdapter.notifyDataSetChanged();
								mProvinceSearchAdapter.setData(mProvinceList);
								mProvinceSearchAdapter.notifyDataSetChanged();

							}

						});

					}

					@Override
					public void onError(String error) {
						Log.i("opp", "error:" + error);

					}
				});

		HttpCommunicationUtil.sendGetDataToServer(BankBaseUtil.getBank(),
				new CommunicationListener() {

					@Override
					public void onResult(String result) {
						Log.i("opp", "result:" + result);
						try {
							JSONObject jsonObj = new JSONObject(result);
							Iterator it = jsonObj.keys();
							mOpenBankList.clear();
							mOpenBankList.add("请选择开户银行");

							mBankIdList.clear();
							mBankIdList.add("");

							while (it.hasNext()) {
								String key = it.next().toString();
								mOpenBankList.add(JsonUtil.getParam(
										JsonUtil.getParam(result, key),
										"bank_name"));
								mBankIdList.add(JsonUtil.getParam(
										JsonUtil.getParam(result, key), "id"));
							}

						} catch (JSONException e) {
							// TODO Auto-generated catch block
							e.printStackTrace();
						}

						((Activity) mContext).runOnUiThread(new Runnable() {

							@Override
							public void run() {
								// 更新UI
								mOpenBankSpinner.setSelection(0);
								mOpenBankAdapter.notifyDataSetChanged();
								mOpenBankSearchAdapter.notifyDataSetChanged();
							}

						});

					}

					@Override
					public void onError(String error) {
						Log.i("opp", "error:" + error);

					}
				});

	}

	private void initListener() {

		

		mProvinceEdit.addTextChangedListener(new TextWatcher() {

			@Override
			public void onTextChanged(CharSequence s, int start, int before,
					int count) {
				// TODO Auto-generated method stub

			}

			@Override
			public void beforeTextChanged(CharSequence s, int start, int count,
					int after) {
				// TODO Auto-generated method stub

			}

			@Override
			public void afterTextChanged(Editable s) {
				mCityList.clear();
				mCityCodeList.clear();
				mCityList.add("开户行所在城市");
				mCityCodeList.add("");
				mCityEdit.setText("");
				if (mProvinceList.indexOf(mProvinceEdit.getText().toString()) > 0) {
					String province = mProvinceEdit.getText().toString();
					Log.i("xxx", "province" + province);

					HttpCommunicationUtil.sendGetDataToServer(
							BankBaseUtil.getCity(province),
							new CommunicationListener() {

								@Override
								public void onResult(String result) {
									Log.i("opp", "result:" + result);
									try {
										JSONArray jsonArray = new JSONArray(
												result);

										for (int i = 0; i < jsonArray.length(); i++) {

											mCityList.add(JsonUtil.getParam(
													jsonArray.getString(i),
													"city_name"));
											mCityCodeList.add(JsonUtil
													.getParam(jsonArray
															.getString(i),
															"city_code"));
										}

									} catch (JSONException e) {
										// TODO Auto-generated catch block
										e.printStackTrace();
									}

									((Activity) mContext)
											.runOnUiThread(new Runnable() {

												@Override
												public void run() {
													// 更新UI
													mCitySpinner
															.setSelection(0);
													mCityEdit.setText("");
													mCityAdapter
															.notifyDataSetChanged();
													mCitySearchAdapter
															.setData(mCityList);

												}

											});

								}

								@Override
								public void onError(String error) {
									Log.i("opp", "error:" + error);

								}
							});
				}

			}
		});

		mCityEdit.addTextChangedListener(new TextWatcher() {

			@Override
			public void onTextChanged(CharSequence s, int start, int before,
					int count) {
				// TODO Auto-generated method stub

			}

			@Override
			public void beforeTextChanged(CharSequence s, int start, int count,
					int after) {
				// TODO Auto-generated method stub

			}

			@Override
			public void afterTextChanged(Editable s) {
				mOpenBankList.clear();
				mOpenBankList.add("请选择开户银行");
				mBankIdList.clear();
				mBankIdList.add("");
				mOpenBankEdit.setText("");
				if (mCityList.indexOf(mCityEdit.getText().toString()) > 0) {
					HttpCommunicationUtil.sendGetDataToServer(
							BankBaseUtil.getBank(),
							new CommunicationListener() {

								@Override
								public void onResult(String result) {
									Log.i("opp", "result:" + result);
									try {
										JSONObject jsonObj = new JSONObject(
												result);
										Iterator it = jsonObj.keys();
										mOpenBankList.clear();
										mOpenBankList.add("请选择开户银行");

										mBankIdList.clear();
										mBankIdList.add("");

										while (it.hasNext()) {
											String key = it.next().toString();
											mOpenBankList.add(JsonUtil
													.getParam(JsonUtil
															.getParam(result,
																	key),
															"bank_name"));
											mBankIdList.add(JsonUtil.getParam(
													JsonUtil.getParam(result,
															key), "id"));
										}

									} catch (JSONException e) {
										// TODO Auto-generated catch block
										e.printStackTrace();
									}

									((Activity) mContext)
											.runOnUiThread(new Runnable() {

												@Override
												public void run() {
													// 更新UI

													mOpenBankSpinner
															.setSelection(0);
													mOpenBankEdit.setText("");
													mOpenBankAdapter
															.notifyDataSetChanged();
													mOpenBankSearchAdapter
															.setData(mOpenBankList);
												}

											});

								}

								@Override
								public void onError(String error) {
									Log.i("opp", "error:" + error);

								}
							});
				}

			}
		});

		mOpenBankEdit.addTextChangedListener(new TextWatcher() {

			@Override
			public void onTextChanged(CharSequence s, int start, int before,
					int count) {
				// TODO Auto-generated method stub

			}

			@Override
			public void beforeTextChanged(CharSequence s, int start, int count,
					int after) {
				// TODO Auto-generated method stub

			}

			@Override
			public void afterTextChanged(Editable s) {
				mBranchBankList.clear();
				mBranchBankList.add("请选择开户支行");
				mBankNoList.clear();
				mBankNoList.add("行号");
				mBranchBankEdit.setText("");
				if (mOpenBankList.indexOf(mOpenBankEdit.getText().toString()) > 0) {
					HttpCommunicationUtil.sendGetDataToServer(BankBaseUtil
							.getSerach(mCityCodeList.get(mCityList
									.indexOf(mCityEdit.getText().toString())),
									mBankIdList.get(mOpenBankList
											.indexOf(mOpenBankEdit.getText()
													.toString()))),
							new CommunicationListener() {

								@Override
								public void onResult(String result) {
									try {
										JSONArray jsonArray = new JSONArray(
												result);
										mBranchBankList.clear();
										mBranchBankList.add("请选择开户支行");
										mBankNoList.clear();
										mBankNoList.add("行号");

										for (int i = 0; i < jsonArray.length(); i++) {

											mBranchBankList.add(JsonUtil
													.getParam(jsonArray
															.getString(i),
															"bank_name"));
											mBankNoList.add(JsonUtil.getParam(
													jsonArray.getString(i),
													"one_bank_no")
													+ "|"
													+ JsonUtil.getParam(
															jsonArray
																	.getString(i),
															"two_bank_no"));
										}

									} catch (JSONException e) {
										// TODO Auto-generated catch block
										e.printStackTrace();
									}

									((Activity) mContext)
											.runOnUiThread(new Runnable() {

												@Override
												public void run() {
													// 更新UI

													mBranchBankSpinner
															.setSelection(0);
													mBranchBankEdit.setText("");
													mBranchBankAdapter
															.notifyDataSetChanged();
													mBranchBankSearchAdapter
															.setData(mBranchBankList);
												}

											});

								}

								@Override
								public void onError(String error) {
									Log.i("opp", "error:" + error);

								}
							});
				}

			}
		});

		mProvinceSpinner
				.setOnItemSelectedListener(new OnItemSelectedListener() {

					@Override
					public void onItemSelected(AdapterView<?> parent,
							View view, int position, long id) {
						if (position > 0) {
							mProvinceEdit.setText(mProvinceList.get(position));

						}

					}

					@Override
					public void onNothingSelected(AdapterView<?> parent) {
						// TODO Auto-generated method stub

					}
				});

		mCitySpinner.setOnItemSelectedListener(new OnItemSelectedListener() {

			@Override
			public void onItemSelected(AdapterView<?> parent, View view,
					int position, long id) {
				if (position > 0) {
					mCityEdit.setText(mCityList.get(position));
				}

			}

			@Override
			public void onNothingSelected(AdapterView<?> parent) {
				// TODO Auto-generated method stub

			}
		});

		mOpenBankSpinner
				.setOnItemSelectedListener(new OnItemSelectedListener() {

					@Override
					public void onItemSelected(AdapterView<?> parent,
							View view, int position, long id) {
						if (position > 0) {
							mOpenBankEdit.setText(mOpenBankList.get(position));
						}

					}

					@Override
					public void onNothingSelected(AdapterView<?> parent) {
						// TODO Auto-generated method stub

					}
				});

		mBranchBankSpinner
				.setOnItemSelectedListener(new OnItemSelectedListener() {

					@Override
					public void onItemSelected(AdapterView<?> parent,
							View view, int position, long id) {
						if (position > 0) {
							mBranchBankEdit.setText(mBranchBankList
									.get(position));
						}

					}

					@Override
					public void onNothingSelected(AdapterView<?> parent) {
						// TODO Auto-generated method stub

					}
				});

		mCityEdit.setOnFocusChangeListener(new OnFocusChangeListener() {

			@Override
			public void onFocusChange(View v, boolean hasFocus) {
				if (hasFocus) {
					if (mProvinceList.indexOf(mProvinceEdit.getText()
							.toString()) < 0) {
						mProvinceEdit.setText("");
					}
				}

			}
		});

		mOpenBankEdit.setOnFocusChangeListener(new OnFocusChangeListener() {

			@Override
			public void onFocusChange(View v, boolean hasFocus) {
				if (hasFocus) {
					if (mCityList.indexOf(mCityEdit.getText().toString()) < 0) {
						mCityEdit.setText("");
					}
				}

			}
		});

		mBranchBankEdit.setOnFocusChangeListener(new OnFocusChangeListener() {

			@Override
			public void onFocusChange(View v, boolean hasFocus) {
				if (hasFocus) {
					if (mOpenBankList.indexOf(mOpenBankEdit.getText()
							.toString()) < 0) {
						mOpenBankEdit.setText("");
					}
				}else{
					if (mBranchBankList.indexOf(mBranchBankEdit.getText().toString()) < 0) {
						mBranchBankEdit.setText("");
					}
				}

			}
		});

		
	}
	public void BtnRegisterFinishedOnClick(View view) {
		if (validate()) {
			mLoading_Dialog.startLoading();
			User user = new User();
			user.setUsername(SessonData.loginUser.getUsername());
			user.setPassword(SessonData.loginUser.getPassword());
			// user.setBank_open(mOpenBankEdit.getText().toString());
			user.setProvince(mProvinceEdit.getText().toString());
			user.setBank_open(mOpenBankEdit.getText().toString());
			user.setCity(mCityEdit.getText().toString());
			user.setBranch_bank(mBranchBankEdit.getText().toString());
			user.setBankNo(mBankNoList.get(mBranchBankList.indexOf(mBranchBankEdit.getText().toString())));
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
								String user_json = JsonUtil.getParam(result,
										"user");
								SessonData.loginUser.setClientid(JsonUtil
										.getParam(user_json, "clientid"));
								SessonData.loginUser.setObject_id(JsonUtil
										.getParam(user_json, "objectId"));
								SessonData.loginUser.setLimit(JsonUtil
										.getParam(user_json, "limit"));
								InitData data = new InitData();
								data.mchntid = SessonData.loginUser
										.getClientid();// 商户号
								data.inscd = JsonUtil.getParam(user_json,
										"inscd");// 机构号
								data.signKey = JsonUtil.getParam(user_json,
										"signKey");// 秘钥
								// Log.e("opp",
								// ""+TelephonyManagerUtil.getDeviceId(mContext));
								data.terminalid = TelephonyManagerUtil
										.getDeviceId(mContext);// 设备号
								data.isProduce = SystemConfig.IS_PRODUCE;// 是否生产环境
								CashierSdk.init(data);
								Intent intent = new Intent(
										RegisterNextActivity.this,
										MainActivity.class);
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

	@SuppressLint("NewApi")
	private boolean validate() {
		if (mBranchBankList.indexOf(mBranchBankEdit.getText().toString()) < 0) {
			mBranchBankEdit.setText("");
		}
		String openbank = "";
		String name = mNameEdit.getText().toString().replace(" ", "");
		String banknum = mBanknumEdit.getText().toString().replace(" ", "");
		String phonenum = mPhonenumEdit.getText().toString().replace(" ", "");

		if (mProvinceEdit.getText().toString().isEmpty()) {
			alertShow("开户行所在省份不能为空!", BitmapFactory
					.decodeResource(this.getResources(), R.drawable.wrong));
			return false;
		}

		if (mCityEdit.getText().toString().isEmpty()) {
			alertShow("开户行所在城市不能为空!", BitmapFactory
					.decodeResource(this.getResources(), R.drawable.wrong));
			return false;
		}

		if (mOpenBankEdit.getText().toString().isEmpty()) {
			alertShow("开户行不能为空!", BitmapFactory.decodeResource(
					this.getResources(), R.drawable.wrong));
			return false;
		}

		if (mBranchBankEdit.getText().toString().isEmpty()) {
			alertShow("开户支行不能为空!", BitmapFactory.decodeResource(
					this.getResources(), R.drawable.wrong));
			return false;
		}

		if (name.isEmpty()) {
			alertShow("姓名不能为空!", BitmapFactory.decodeResource(
					this.getResources(), R.drawable.wrong));
			return false;
		}

		if (banknum.isEmpty()) {
			alertShow("银行卡号不能为空!", BitmapFactory.decodeResource(
					this.getResources(), R.drawable.wrong));
			return false;
		}

		if (!VerifyUtil.checkBankCard(banknum)) {
			alertShow("请输入正确的银行卡号!", BitmapFactory
					.decodeResource(this.getResources(), R.drawable.wrong));
			return false;
		}

		if (phonenum.isEmpty()) {
			alertShow("手机号不能为空!", BitmapFactory.decodeResource(
					this.getResources(), R.drawable.wrong));
			return false;
		}

		if (!VerifyUtil.isMobileNO(phonenum)) {
			alertShow("请输入正确的手机号!", BitmapFactory.decodeResource(
					this.getResources(), R.drawable.wrong));
			return false;
		}

		return true;
	}
}
