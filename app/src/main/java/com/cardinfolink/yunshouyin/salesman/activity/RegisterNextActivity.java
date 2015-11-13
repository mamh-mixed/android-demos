package com.cardinfolink.yunshouyin.salesman.activity;

import android.annotation.SuppressLint;
import android.content.SharedPreferences;
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
import com.cardinfolink.yunshouyin.salesman.adapter.SearchAdapter;
import com.cardinfolink.yunshouyin.salesman.api.QuickPayException;
import com.cardinfolink.yunshouyin.salesman.core.QuickPayCallbackListener;
import com.cardinfolink.yunshouyin.salesman.model.SessonData;
import com.cardinfolink.yunshouyin.salesman.model.User;
import com.cardinfolink.yunshouyin.salesman.utils.ActivityCollector;
import com.cardinfolink.yunshouyin.salesman.utils.BankBaseUtil;
import com.cardinfolink.yunshouyin.salesman.utils.CommunicationListener;
import com.cardinfolink.yunshouyin.salesman.utils.HttpCommunicationUtil;
import com.cardinfolink.yunshouyin.salesman.utils.JsonUtil;
import com.cardinfolink.yunshouyin.salesman.utils.RequestParam;
import com.cardinfolink.yunshouyin.salesman.utils.VerifyUtil;

import org.json.JSONArray;
import org.json.JSONException;
import org.json.JSONObject;

import java.util.ArrayList;
import java.util.Iterator;
import java.util.List;

public class RegisterNextActivity extends BaseActivity {
    private static final String TAG = "RegisterNextActivity";
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

    //分行，支行
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

        mOpenBankAdapter = new ArrayAdapter<String>(mContext, R.layout.spinner_item, mOpenBankList);
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

        mBranchBankAdapter = new ArrayAdapter<String>(mContext, R.layout.spinner_item, mBranchBankList);
        // 设置样式
        mBranchBankAdapter.setDropDownViewResource(R.layout.spinner_drop_item);
        // 加载适配器
        mBranchBankSpinner.setAdapter(mBranchBankAdapter);
        mBranchBankSearchAdapter = new SearchAdapter(mContext, mBranchBankList);
        mBranchBankEdit.setAdapter(mBranchBankSearchAdapter);
        mBranchBankEdit.setThreshold(1);
    }

    private void initProvinceData() {
        String data = readFromSharePreference("data", "province");
        if (data != null && data.length() != 0) {
            Log.d(TAG, "will use cache data to get province: \n" + data);
            updateProvinceAdapter(data);
        } else {
            Log.d(TAG, "will do post to get province data");
            RequestParam provinceParam = BankBaseUtil.getProvince();
            ProvinceCommunicationListener provinceCommunicationListener = new ProvinceCommunicationListener();
            HttpCommunicationUtil.sendGetDataToServer(provinceParam, provinceCommunicationListener);
        }
    }

    private void initBankData() {
        //获取bank的数据
        String data = readFromSharePreference("data", "bank");
        if (data != null && data.length() != 0) {
            Log.d(TAG, "will use cache data to get bank: \n" + data);
            updateBankAdapter(data);
        } else {
            Log.d(TAG, "will do post to get bank data");
            RequestParam bankParam = BankBaseUtil.getBank();
            BankCommunicationListener bankCommunicationListener = new BankCommunicationListener();
            HttpCommunicationUtil.sendGetDataToServer(bankParam, bankCommunicationListener);
        }
    }

    private void initCityData(String province) {
        String data = readFromSharePreference("data", province);
        if (data != null && data.length() != 0) {
            Log.d(TAG, "will use cache data to get City: " + data);
            updateCityAdapter(data);
        } else {
            Log.d(TAG, "will post to get City: " + data);
            RequestParam cityParam = BankBaseUtil.getCity(province);
            CityCommunicationListener cityCommunicationListener = new CityCommunicationListener(province);
            HttpCommunicationUtil.sendGetDataToServer(cityParam, cityCommunicationListener);
        }
    }

    private void initBranchBankData(String cityCode, String bankId) {
        String key = cityCode + "_" + bankId;
        String data = readFromSharePreference("data", key);
        if (data != null && data.length() != 0) {
            Log.d(TAG, "will use cache data to get branch bank");
            updateBranchBankAdapter(data);
        } else {
            Log.d(TAG, "will use post to get branch bank");
            RequestParam bbParam = BankBaseUtil.getSerach(cityCode, bankId);
            BranchBankCommunicationListener bbCL = new BranchBankCommunicationListener(cityCode, bankId);
            HttpCommunicationUtil.sendGetDataToServer(bbParam, bbCL);
        }
    }

    // 开户省份和总行列表先读取
    public void initData() {
        Log.d(TAG, "initData start");
        initProvinceData();
        initBankData();
        Log.d(TAG, "initData end");
    }//end public void initData()


    //字符串  -->  jsonArray --> List
    private List<String> jsonArrayToList(String result) {
        final List<String> list = new ArrayList<>();
        try {
            JSONArray jsonArray = new JSONArray(result);
            for (int i = 0; i < jsonArray.length(); i++) {
                list.add(jsonArray.getString(i));
            }
        } catch (JSONException e) {
            e.printStackTrace();
        }
        return list;
    }

    //字符串  -->  jsonArray(每个数组元素又是一个jsonObject) ----> List
    private List<String> jsonArrayToList(String result, String key) {
        final List<String> list = new ArrayList<>();
        try {
            JSONArray jsonArray = new JSONArray(result);
            for (int i = 0; i < jsonArray.length(); i++) {
                String tempjson = jsonArray.getString(i);
                list.add(JsonUtil.getParam(tempjson, key));
            }
        } catch (JSONException e) {
            // TODO Auto-generated catch block
            e.printStackTrace();
        }
        return list;
    }

    private List<String> jsonArrayToList(String result, String key1, String key2, String sep) {
        final List<String> list = new ArrayList<>();
        try {
            JSONArray jsonArray = new JSONArray(result);
            for (int i = 0; i < jsonArray.length(); i++) {
                String tempjson = jsonArray.getString(i);
                String value1 = JsonUtil.getParam(tempjson, key1);
                String value2 = JsonUtil.getParam(tempjson, key2);
                list.add(value1 + sep + value2);
            }
        } catch (JSONException e) {
            e.printStackTrace();
        }
        return list;
    }

    //字符串 --> jsonObject --> List
    private List<String> jsonObjectToList(String result, String key) {
        final List<String> list = new ArrayList<String>();
        try {
            //string 转换为  jsonObj
            JSONObject jsonObj = new JSONObject(result);
            Iterator it = jsonObj.keys();
            while (it.hasNext()) {
                String tempkey = it.next().toString();
                String subJson = JsonUtil.getParam(result, tempkey);
                list.add(JsonUtil.getParam(subJson, key));
            }
        } catch (JSONException e) {
            e.printStackTrace();
        }
        return list;
    }

    private void updateProvinceAdapter(String data) {
        final List<String> tempProvinceList = jsonArrayToList(data);
        tempProvinceList.add(0, "开户行所在省份");
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

    private void updateCityAdapter(String data) {
        final List<String> tempCityList = jsonArrayToList(data, "city_name");
        final List<String> tempCityCodeList = jsonArrayToList(data, "city_code");
        tempCityList.add(0, "开户行所在城市");
        tempCityCodeList.add(0, "");
        runOnUiThread(new Runnable() {
            @Override
            public void run() {
                mCityList.clear();
                mCityList.addAll(tempCityList);
                mCityCodeList.clear();
                mCityCodeList.addAll(tempCityCodeList);
                mCitySpinner.setSelection(0);
                mCityAdapter.notifyDataSetChanged();
                mCitySearchAdapter.setData(mCityList);
            }
        });
    }


    private void updateBankAdapter(String data) {
        final List<String> tempOpenBankList = jsonObjectToList(data, "bank_name");
        final List<String> tempBankIdList = jsonObjectToList(data, "id");
        tempOpenBankList.add(0, "请选择开户银行");
        tempBankIdList.add(0, "");//为了使index对应起来
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

    private void updateBranchBankAdapter(String data) {
        final List<String> tempBranchBankList = jsonArrayToList(data, "bank_name");
        final List<String> tempBankNoList = jsonArrayToList(data, "one_bank_no", "two_bank_no", "|");
        tempBranchBankList.add(0, "请选择开户支行");
        tempBankNoList.add(0, "行号");
        runOnUiThread(new Runnable() {
            @Override
            public void run() {
                mBranchBankList.clear();
                mBranchBankList.addAll(tempBranchBankList);

                mBankNoList.clear();
                mBankNoList.addAll(tempBankNoList);

                mBranchBankSpinner.setSelection(0);

                mBranchBankAdapter.notifyDataSetChanged();
                mBranchBankSearchAdapter.setData(mBranchBankList);
            }

        });
    }

    private String readFromSharePreference(String name, String key) {
        SharedPreferences sp = getSharedPreferences(name, MODE_PRIVATE);
        return sp.getString(key, "");
    }

    private void saveToSharePreferences(String result, String name, String key) {
        SharedPreferences sp = getSharedPreferences(name, MODE_PRIVATE);
        SharedPreferences.Editor editor = sp.edit();
        editor.putString(key, result);
        editor.commit();
    }

    private void initListener() {
        //添加province EditText框变化事件
        mProvinceEdit.addTextChangedListener(new TextWatcher() {
            @Override
            public void onTextChanged(CharSequence s, int start, int before, int count) {
                // TODO Auto-generated method stub
            }

            @Override
            public void beforeTextChanged(CharSequence s, int start, int count, int after) {
                // TODO Auto-generated method stub
            }

            @Override
            public void afterTextChanged(Editable s) {
                mCityEdit.setText("");
                if (mProvinceList.indexOf(mProvinceEdit.getText().toString()) > 0) {
                    String province = mProvinceEdit.getText().toString();
                    initCityData(province);
                }//end if()
            }//end public void afterTextChange(Editable s)
        });

        //添加City EditText 框变化事件
        mCityEdit.addTextChangedListener(new TextWatcher() {
            @Override
            public void onTextChanged(CharSequence s, int start, int before, int count) {
                // TODO Auto-generated method stub
            }

            @Override
            public void beforeTextChanged(CharSequence s, int start, int count, int after) {
                // TODO Auto-generated method stub
            }

            @Override
            public void afterTextChanged(Editable s) {
                mOpenBankEdit.setText("");
                if (mCityList.indexOf(mCityEdit.getText().toString()) > 0) {
                    initBankData();
                }
            }
        });

        //添加bank EditText 事件
        mOpenBankEdit.addTextChangedListener(new TextWatcher() {
            @Override
            public void onTextChanged(CharSequence s, int start, int before, int count) {
                // TODO Auto-generated method stub
            }

            @Override
            public void beforeTextChanged(CharSequence s, int start, int count, int after) {
                // TODO Auto-generated method stub
            }

            @Override
            public void afterTextChanged(Editable s) {
                mBranchBankEdit.setText("");
                if (mOpenBankList.indexOf(mOpenBankEdit.getText().toString()) > 0) {
                    int indexOfCity = mCityList.indexOf(mCityEdit.getText().toString());
                    int indexOfBank = mOpenBankList.indexOf(mOpenBankEdit.getText().toString());
                    String cityCode = mCityCodeList.get(indexOfCity);
                    String bankId = mBankIdList.get(indexOfBank);
                    initBranchBankData(cityCode, bankId);
                }
            }
        });

        mProvinceSpinner.setOnItemSelectedListener(new OnItemSelectedListener() {
            @Override
            public void onItemSelected(AdapterView<?> parent, View view, int position, long id) {
                if (position > 0) {
                    mProvinceEdit.setText(mProvinceList.get(position));
                }
            }

            @Override
            public void onNothingSelected(AdapterView<?> parent) {
            }
        });

        mCitySpinner.setOnItemSelectedListener(new OnItemSelectedListener() {
            @Override
            public void onItemSelected(AdapterView<?> parent, View view, int position, long id) {
                if (position > 0) {
                    mCityEdit.setText(mCityList.get(position));
                }
            }

            @Override
            public void onNothingSelected(AdapterView<?> parent) {
            }
        });

        mOpenBankSpinner.setOnItemSelectedListener(new OnItemSelectedListener() {
            @Override
            public void onItemSelected(AdapterView<?> parent, View view, int position, long id) {
                if (position > 0) {
                    mOpenBankEdit.setText(mOpenBankList.get(position));
                }
            }

            @Override
            public void onNothingSelected(AdapterView<?> parent) {
            }
        });

        mBranchBankSpinner.setOnItemSelectedListener(new OnItemSelectedListener() {
            @Override
            public void onItemSelected(AdapterView<?> parent, View view, int position, long id) {
                if (position > 0) {
                    mBranchBankEdit.setText(mBranchBankList.get(position));
                }
            }

            @Override
            public void onNothingSelected(AdapterView<?> parent) {
            }
        });

        mCityEdit.setOnFocusChangeListener(new OnFocusChangeListener() {
            @Override
            public void onFocusChange(View v, boolean hasFocus) {
                if (hasFocus) {
                    if (mProvinceList.indexOf(mProvinceEdit.getText().toString()) < 0) {
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
                    if (mOpenBankList.indexOf(mOpenBankEdit.getText().toString()) < 0) {
                        mOpenBankEdit.setText("");
                    }
                }
            }
        });

    }

    public void btnRegisterFinishedOnClick(View view) {
        if (validate()) {
            startLoading();
            final User user = SessonData.registerUser;
            user.setProvince(mProvinceEdit.getText().toString());
            user.setCity(mCityEdit.getText().toString());
            user.setBankOpen(mOpenBankEdit.getText().toString());
            user.setBranchBank(mBranchBankEdit.getText().toString());

            //有些地方没有支行，get()会抛出outofindex异常
            String branchBank = mBranchBankEdit.getText().toString();
            int index = mBranchBankList.indexOf(branchBank);
            user.setBankNo((index != -1) ? mBankNoList.get(index) : "");

            user.setPayee(mNameEdit.getText().toString());
            user.setPayeeCard(mBanknumEdit.getText().toString().replace(" ", ""));
            user.setPhoneNum(mPhonenumEdit.getText().toString());
            user.setMerName(mMerchantNameEdit.getText().toString());


            application.getQuickPayService().updateUserAsync(user, new QuickPayCallbackListener<User>() {
                @Override
                public void onSuccess(final User data) {
                    runOnUiThread(new Runnable() {
                        @Override
                        public void run() {
                            //NOTE:clientID也是merchantId,用于在七牛那边创建唯一id
                            SessonData.registerUser.setClientid(data.getClientid());
                            endLoading();
                            intentToActivity(RegisterStep3Activity.class);
                        }
                    });

                }

                @Override
                public void onFailure(final QuickPayException ex) {
                    runOnUiThread(new Runnable() {
                        @Override
                        public void run() {
                            String error = ex.getErrorCode();
                            String errorStr = ex.getErrorMsg();
                            endLoadingWithError(errorStr);
                            if (error.equals(QuickPayException.ACCESSTOKEN_NOT_FOUND)) {
                                //关闭所有activity,除了登录框
                                ActivityCollector.goLoginAndFinishRest();
                            }
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
            if (mBranchBankList.size() == 1 && mBranchBankList.get(0).equals("请选择开户支行")) {
                //有些地方没有支行，这里不填写就不能下一步
            }else {
                alertError("开户支行不能为空!");
                return false;
            }
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
            return false;
        }

        return true;
    }

    //内部类，实现CommunicationListener接口
    private class ProvinceCommunicationListener implements CommunicationListener {
        @Override
        public void onResult(String result) {
            saveToSharePreferences(result, "data", "province");
            updateProvinceAdapter(result);
        }

        @Override
        public void onError(String error) {
            Log.i(TAG, "get province data error:" + error);
        }
    }

    //内部类，实现CommunicationListener接口,用来获取bank信息
    private class BankCommunicationListener implements CommunicationListener {
        @Override
        public void onResult(String result) {
            saveToSharePreferences(result, "data", "bank");
            updateBankAdapter(result);
        }

        @Override
        public void onError(String error) {
            Log.i(TAG, "get bank data error:" + error);
        }
    }

    private class CityCommunicationListener implements CommunicationListener {
        private String province;

        public CityCommunicationListener(String province) {
            this.province = province;
        }

        @Override
        public void onResult(String result) {
            saveToSharePreferences(result, "data", province);
            updateCityAdapter(result);
        }

        @Override
        public void onError(String error) {
            Log.d(TAG, "get city data error");
        }
    }

    private class BranchBankCommunicationListener implements CommunicationListener {
        private String cityCode;
        private String bankId;

        public BranchBankCommunicationListener(String cityCode, String bankId) {
            this.cityCode = cityCode;
            this.bankId = bankId;
        }

        @Override
        public void onResult(String result) {
            String key = cityCode + "_" + bankId;
            saveToSharePreferences(result, "data", key);
            updateBranchBankAdapter(result);
        }

        @Override
        public void onError(String error) {
            Log.i(TAG, "get branch bank error:" + error);
        }
    }
}
