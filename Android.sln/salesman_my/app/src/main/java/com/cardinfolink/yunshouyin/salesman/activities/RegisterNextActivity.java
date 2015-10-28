package com.cardinfolink.yunshouyin.salesman.activities;

import android.annotation.SuppressLint;
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

import com.cardinfolink.yunshouyin.salesman.R;
import com.cardinfolink.yunshouyin.salesman.models.SAServerPacket;
import com.cardinfolink.yunshouyin.salesman.models.SessonData;
import com.cardinfolink.yunshouyin.salesman.models.User;
import com.cardinfolink.yunshouyin.salesman.utils.ActivityCollector;
import com.cardinfolink.yunshouyin.salesman.utils.BankBaseUtil;
import com.cardinfolink.yunshouyin.salesman.utils.CommunicationListener;
import com.cardinfolink.yunshouyin.salesman.utils.ErrorUtil;
import com.cardinfolink.yunshouyin.salesman.utils.HttpCommunicationUtil;
import com.cardinfolink.yunshouyin.salesman.utils.JsonUtil;
import com.cardinfolink.yunshouyin.salesman.utils.ParamsUtil;
import com.cardinfolink.yunshouyin.salesman.utils.VerifyUtil;

import org.json.JSONArray;
import org.json.JSONException;
import org.json.JSONObject;

import java.util.ArrayList;
import java.util.Iterator;
import java.util.List;

public class RegisterNextActivity extends BaseActivity {

    private EditText mNameEdit;
    private EditText mBanknumEdit;
    private EditText mPhonenumEdit;
    private EditText mMerchantNameEdit;

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
        mMerchantNameEdit = (EditText) findViewById(R.id.info_merchantname);
        VerifyUtil.bankCardNumAddSpace(mBanknumEdit);

        /**
         * setup for Province
         */
        // view
        mProvinceEdit = (AutoCompleteTextView) findViewById(R.id.edit_province);
        mProvinceSpinner = (Spinner) findViewById(R.id.spinner_province);

        // 适配器data
        mProvinceList = new ArrayList<String>();
        mProvinceList.add("开户行所在省份");

        // spinner view
        mProvinceAdapter = new ArrayAdapter<>(mContext, R.layout.spinner_item, mProvinceList);
        // 设置样式
        mProvinceAdapter.setDropDownViewResource(R.layout.spinner_drop_item);
        // 加载适配器
        mProvinceSpinner.setAdapter(mProvinceAdapter);

        // autocomplete text view
        mProvinceSearchAdapter = new SearchAdapter(mContext, mProvinceList);
        mProvinceEdit.setAdapter(mProvinceSearchAdapter);
        mProvinceEdit.setThreshold(1);


        mCityEdit = (AutoCompleteTextView) findViewById(R.id.edit_city);
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


    // 开户省份和总行列表先读取
    public void initData() {

        HttpCommunicationUtil.sendGetDataToServer(BankBaseUtil.getProvince(),
                new CommunicationListener() {
                    @Override
                    public void onResult(String result) {
                        final List<String> tempProvinceList = new ArrayList<>();
                        tempProvinceList.add("开户行所在省份");
                        try {
                            JSONArray jsonArray = new JSONArray(result);
                            for (int i = 0; i < jsonArray.length(); i++) {
                                tempProvinceList.add(jsonArray.getString(i));
                            }

                        } catch (JSONException e) {
                            // TODO Auto-generated catch block
                            e.printStackTrace();
                        }

                        runOnUiThread(new Runnable() {

                            @Override
                            public void run() {
                                mProvinceList.clear();
                                mProvinceList.addAll(tempProvinceList);
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
                        final List<String> tempOpenBankList = new ArrayList<String>();
                        final List<String> tempBankIdList = new ArrayList<String>();
                        tempOpenBankList.add("请选择开户银行");
                        tempBankIdList.add("");//为了使index对应起来

                        try {
                            JSONObject jsonObj = new JSONObject(result);
                            Iterator it = jsonObj.keys();

                            while (it.hasNext()) {
                                String key = it.next().toString();
                                tempOpenBankList.add(JsonUtil.getParam(
                                        JsonUtil.getParam(result, key),
                                        "bank_name"));
                                tempBankIdList.add(JsonUtil.getParam(
                                        JsonUtil.getParam(result, key), "id"));
                            }

                        } catch (JSONException e) {
                            // TODO Auto-generated catch block
                            e.printStackTrace();
                        }

                        runOnUiThread(new Runnable() {

                            @Override
                            public void run() {
                                mOpenBankList.clear();
                                mBankIdList.clear();
                                mOpenBankList.addAll(tempOpenBankList);
                                mBankIdList.addAll(tempBankIdList);
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
                                    final List<String> tempCityList = new ArrayList<>();
                                    final List<String> tempCityCodeList = new ArrayList<>();

                                    tempCityList.add("开户行所在城市");
                                    tempCityCodeList.add("");

                                    try {
                                        JSONArray jsonArray = new JSONArray(
                                                result);

                                        for (int i = 0; i < jsonArray.length(); i++) {

                                            tempCityList.add(JsonUtil.getParam(jsonArray.getString(i), "city_name"));
                                            tempCityCodeList.add(JsonUtil.getParam(jsonArray.getString(i), "city_code"));
                                        }

                                    } catch (JSONException e) {
                                        // TODO Auto-generated catch block
                                        e.printStackTrace();
                                    }

                                    runOnUiThread(new Runnable() {

                                        @Override
                                        public void run() {
                                            mCityList.clear();
                                            mCityCodeList.clear();

                                            mCityList.addAll(tempCityList);
                                            mCityCodeList.addAll(tempCityCodeList);
                                            mCitySpinner.setSelection(0);
                                            mCityEdit.setText("");
                                            mCityAdapter.notifyDataSetChanged();
                                            mCitySearchAdapter.setData(mCityList);
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
                mOpenBankEdit.setText("");
                if (mCityList.indexOf(mCityEdit.getText().toString()) > 0) {
                    HttpCommunicationUtil.sendGetDataToServer(
                            BankBaseUtil.getBank(),
                            new CommunicationListener() {

                                @Override
                                public void onResult(String result) {
                                    Log.i("opp", "result:" + result);
                                    final List<String> tempOpenBankList = new ArrayList<>();
                                    final List<String> tempBankIdList = new ArrayList<>();
                                    tempOpenBankList.add("请选择开户银行");
                                    tempBankIdList.add("");

                                    try {
                                        JSONObject jsonObj = new JSONObject(
                                                result);
                                        Iterator it = jsonObj.keys();

                                        while (it.hasNext()) {
                                            String key = it.next().toString();
                                            tempOpenBankList.add(JsonUtil
                                                    .getParam(JsonUtil
                                                                    .getParam(result,
                                                                            key),
                                                            "bank_name"));
                                            tempBankIdList.add(JsonUtil.getParam(
                                                    JsonUtil.getParam(result,
                                                            key), "id"));
                                        }

                                    } catch (JSONException e) {
                                        // TODO Auto-generated catch block
                                        e.printStackTrace();
                                    }

                                    runOnUiThread(new Runnable() {

                                        @Override
                                        public void run() {
                                            mOpenBankList.clear();
                                            mBankIdList.clear();
                                            mOpenBankList.addAll(tempOpenBankList);
                                            mBankIdList.addAll(tempBankIdList);

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
                                    final List<String> tempBranchBankList = new ArrayList<>();
                                    final List<String> tempBankNoList = new ArrayList<>();
                                    tempBranchBankList.add("请选择开户支行");
                                    tempBankNoList.add("行号");

                                    try {
                                        JSONArray jsonArray = new JSONArray(
                                                result);

                                        for (int i = 0; i < jsonArray.length(); i++) {

                                            tempBranchBankList.add(JsonUtil
                                                    .getParam(jsonArray
                                                                    .getString(i),
                                                            "bank_name"));
                                            tempBankNoList.add(JsonUtil.getParam(
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


                                    runOnUiThread(new Runnable() {

                                        @Override
                                        public void run() {
                                            mBranchBankList.clear();
                                            mBankNoList.clear();
                                            mBranchBankList.addAll(tempBranchBankList);
                                            mBankNoList.addAll(tempBankNoList);

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

        mProvinceSpinner.setOnItemSelectedListener(new OnItemSelectedListener() {

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
                }

            }
        });


    }

    public void btnRegisterFinishedOnClick(View view) {
        // Test only
//        mBanknumEdit.setText("6228482371938777011");
//        mPhonenumEdit.setText("18964153831");
//        mMerchantNameEdit.setText("香辣鸡翅");
//        mNameEdit.setText("鸡哥");

        if (validate()) {
            startLoading();

            final User user = SessonData.registerUser;

            user.setProvince(mProvinceEdit.getText().toString());
            user.setCity(mCityEdit.getText().toString());
            user.setBank_open(mOpenBankEdit.getText().toString());
            user.setBranch_bank(mBranchBankEdit.getText().toString());

            user.setBankNo(mBankNoList.get(mBranchBankList.indexOf(mBranchBankEdit.getText().toString())));
            user.setPayee(mNameEdit.getText().toString());
            user.setPayee_card(mBanknumEdit.getText().toString().replace(" ", ""));
            user.setPhone_num(mPhonenumEdit.getText().toString());
            user.setMerName(mMerchantNameEdit.getText().toString());

            Log.d("register user", SessonData.registerUser.getJsonString());

            HttpCommunicationUtil.sendDataToServer(
                    ParamsUtil.getUpdate_SA(SessonData.getAccessToken(), user),
                    new CommunicationListener() {

                        @Override
                        public void onResult(String result) {
                            final SAServerPacket serverPacket = SAServerPacket.getServerPacketFrom(result);
                            String state = JsonUtil.getParam(result, "state");
                            if (serverPacket.getState().equals("success")) {
                                //NOTE:clientID也是merchantId,用于在七牛那边创建唯一id
                                SessonData.registerUser.setClientid(serverPacket.getUser().getClientid());
//                                SessonData.registerUser.setObject_id(JsonUtil
//                                        .getParam(user_json, "objectId"));

                                intentToActivity(SARegisterStep3Activity.class);

                            } else {
                                runOnUiThread(new Runnable() {

                                    @Override
                                    public void run() {
                                        final String error = serverPacket.getError();
                                        String errorStr = ErrorUtil.getErrorString(error);
                                        endLoadingWithError(errorStr);
                                        if (error.equals("accessToken_error")) {
                                            //关闭所有activity,除了登录框
                                            ActivityCollector.goLoginAndFinishRest();
                                        }
                                    }
                                });
                            }
                        }

                        @Override
                        public void onError(final String error) {
                            runOnUiThread(new Runnable() {

                                @Override
                                public void run() {
                                    endLoadingWithError(error);
                                }
                            });
                        }
                    });
        }
    }

    @SuppressLint("NewApi")
    private boolean validate() {
        String openbank = "";
        String name = mNameEdit.getText().toString().replace(" ", "");
        String banknum = mBanknumEdit.getText().toString().replace(" ", "");
        String phonenum = mPhonenumEdit.getText().toString().replace(" ", "");
        String merchantname = mMerchantNameEdit.getText().toString().replace(" ", "");

        if (mProvinceEdit.getText().toString().isEmpty()) {
            alertError("开户行所在省份不能为空!");
            return false;
        }

        if (mCityEdit.getText().toString().isEmpty()) {
            alertError("开户行所在城市不能为空!");
            return false;
        }

        if (mOpenBankEdit.getText().toString().isEmpty()) {
            alertError("开户行不能为空!");
            return false;
        }

        if (mBranchBankEdit.getText().toString().isEmpty()) {
            alertError("开户支行不能为空!");
            return false;
        }

        if (name.isEmpty()) {
            alertError("姓名不能为空!");
            return false;
        }

        if (banknum.isEmpty()) {
            alertError("银行卡号不能为空!");
            return false;
        }

        if (!VerifyUtil.checkBankCard(banknum)) {
            alertError("请输入正确的银行卡号!");
            return false;
        }

        if (phonenum.isEmpty()) {
            alertError("手机号不能为空!");
            return false;
        }

        if (!VerifyUtil.isMobileNO(phonenum)) {
            alertError("请输入正确的手机号!");
            return false;
        }

        if (merchantname.isEmpty()) {
            alertError("请输入商店名称");
        }

        return true;
    }
}
