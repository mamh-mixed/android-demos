package com.cardinfolink.yunshouyin.salesman.activity;

import android.annotation.SuppressLint;
import android.content.Intent;
import android.content.SharedPreferences;
import android.os.Bundle;
import android.text.Editable;
import android.text.TextUtils;
import android.text.TextWatcher;
import android.util.Log;
import android.view.View;
import android.view.View.OnFocusChangeListener;
import android.widget.AdapterView;
import android.widget.AdapterView.OnItemSelectedListener;
import android.widget.ArrayAdapter;
import android.widget.AutoCompleteTextView;
import android.widget.Button;
import android.widget.EditText;
import android.widget.Spinner;

import com.cardinfolink.yunshouyin.salesman.R;
import com.cardinfolink.yunshouyin.salesman.adapter.SearchAdapter;
import com.cardinfolink.yunshouyin.salesman.api.QuickPayException;
import com.cardinfolink.yunshouyin.salesman.core.QuickPayCallbackListener;
import com.cardinfolink.yunshouyin.salesman.model.City;
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

    private Button mRegisterNext;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.register_next_activity);
        initLayout();
        initListener();
    }

    private void initSpinner() {
        //初始化四个spinner
        mProvinceSpinner = (Spinner) findViewById(R.id.spinner_province);
        mCitySpinner = (Spinner) findViewById(R.id.spinner_city);
        mOpenBankSpinner = (Spinner) findViewById(R.id.spinner_openbank);
        mBranchBankSpinner = (Spinner) findViewById(R.id.spinner_branchbank);
    }//end initSpinner()

    private void initEditText() {
        mNameEdit = (EditText) findViewById(R.id.info_name);
        mBanknumEdit = (EditText) findViewById(R.id.info_banknum);
        mPhonenumEdit = (EditText) findViewById(R.id.info_phonenum);
        mMerchantNameEdit = (EditText) findViewById(R.id.info_merchantname);

        mProvinceEdit = (AutoCompleteTextView) findViewById(R.id.edit_province);
        mCityEdit = (AutoCompleteTextView) findViewById(R.id.edit_city);
        mOpenBankEdit = (AutoCompleteTextView) findViewById(R.id.edit_openbank);
        mBranchBankEdit = (AutoCompleteTextView) findViewById(R.id.edit_branchbank);
    }//end initEditText()

    private void initButton() {
        mRegisterNext = (Button) findViewById(R.id.bt_register_next);
        mRegisterNext.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                registerNext(v);
            }
        });
    }

    private void initArrayList() {
        mProvinceList = new ArrayList<String>();
        mProvinceList.add("开户行所在省份");

        mCityList = new ArrayList<String>();
        mCityList.add("开户行所在城市");

        mCityCodeList = new ArrayList<String>();
        mCityCodeList.add("");

        mOpenBankList = new ArrayList<String>();
        mOpenBankList.add("请选择开户银行");

        mBankIdList = new ArrayList<String>();
        mBankIdList.add("");

        mBranchBankList = new ArrayList<String>();
        mBranchBankList.add("请选择开户支行");

        mBankNoList = new ArrayList<String>();
        mBankNoList.add("行号");
    }//end initArrayList()

    private void initAdapter() {
        mProvinceAdapter = new ArrayAdapter<>(mContext, R.layout.spinner_item, mProvinceList);
        mProvinceAdapter.setDropDownViewResource(R.layout.spinner_drop_item);// 设置样式
        mProvinceSpinner.setAdapter(mProvinceAdapter);// 加载适配器

        mProvinceSearchAdapter = new SearchAdapter(mContext, mProvinceList);
        mProvinceEdit.setAdapter(mProvinceSearchAdapter);
        mProvinceEdit.setThreshold(1);

        mCityAdapter = new ArrayAdapter<String>(mContext, R.layout.spinner_item, mCityList);
        mCityAdapter.setDropDownViewResource(R.layout.spinner_drop_item); // 设置样式
        mCitySpinner.setAdapter(mCityAdapter);

        mCitySearchAdapter = new SearchAdapter(mContext, mCityList);
        mCityEdit.setAdapter(mCitySearchAdapter);
        mCityEdit.setThreshold(1);


        mOpenBankAdapter = new ArrayAdapter<String>(mContext, R.layout.spinner_item, mOpenBankList);
        mOpenBankAdapter.setDropDownViewResource(R.layout.spinner_drop_item);
        mOpenBankSpinner.setAdapter(mOpenBankAdapter);

        mOpenBankSearchAdapter = new SearchAdapter(mContext, mOpenBankList);
        mOpenBankEdit.setAdapter(mOpenBankSearchAdapter);
        mOpenBankEdit.setThreshold(1);


        mBranchBankAdapter = new ArrayAdapter<String>(mContext, R.layout.spinner_item, mBranchBankList);
        mBranchBankAdapter.setDropDownViewResource(R.layout.spinner_drop_item);// 设置样式
        mBranchBankSpinner.setAdapter(mBranchBankAdapter);

        mBranchBankSearchAdapter = new SearchAdapter(mContext, mBranchBankList);
        mBranchBankEdit.setAdapter(mBranchBankSearchAdapter);
        mBranchBankEdit.setThreshold(1);
    }//end initAdapter()

    private void initLayout() {
        initSpinner();
        initEditText();
        initButton();
        initArrayList();//一定要注意初始化的顺序
        initAdapter();

        initProvinceData();
        String province = mRegisterSharedPreferences.getString("register_province", "");
        if (!TextUtils.isEmpty(province)) {
            mProvinceEdit.setText(province);
            initCityData(province);
        }

        String city = mRegisterSharedPreferences.getString("register_city", "");
        if (!TextUtils.isEmpty(city)) {
            mCityEdit.setText(city);
        }
        initBankData();

        String bankopen = mRegisterSharedPreferences.getString("register_bankopen", "");
        if (!TextUtils.isEmpty(bankopen)) {
            mOpenBankEdit.setText(bankopen);
            int indexBank = mOpenBankList.indexOf(bankopen);
            int indexCity = mCityList.indexOf(city);

            if (indexBank > 0 && indexCity > 0) {
                String cityCode = mCityCodeList.get(indexCity);
                String bankId = mBankIdList.get(indexBank);
                initBranchBankData(cityCode, bankId);
            }
        }

        String branchBank = mRegisterSharedPreferences.getString("register_branchbank", "");
        if (!TextUtils.isEmpty(branchBank)) {
            int index = mBranchBankList.indexOf(branchBank);
            if (index > 0) {
                mBranchBankEdit.setText(branchBank);
            }
        }
        //收款人
        String payee = mRegisterSharedPreferences.getString("register_payee", "");
        if (!TextUtils.isEmpty(payee)) {
            mNameEdit.setText(payee);
        }

        //收款的银行账号
        String payeeCard = mRegisterSharedPreferences.getString("register_payeecard", "");
        if (!TextUtils.isEmpty(payeeCard)) {
            mBanknumEdit.setText(payeeCard);
        }
        VerifyUtil.bankCardNumAddSpace(mBanknumEdit);

        //电话号
        String phoneNum = mRegisterSharedPreferences.getString("register_phonenum", "");
        if (!TextUtils.isEmpty(phoneNum)) {
            mPhonenumEdit.setText(phoneNum);
        }
        String merName = mRegisterSharedPreferences.getString("register_mername", "");
        if (!TextUtils.isEmpty(merName)) {
            mMerchantNameEdit.setText(merName);
        }

    }

    private void initListener() {
        //添加province EditText框变化事件
        mProvinceEdit.addTextChangedListener(new RegisterTextWatcher(mProvinceEdit));

        //添加City EditText 框变化事件
        mCityEdit.addTextChangedListener(new RegisterTextWatcher(mCityEdit));

        //添加bank EditText 事件
        mOpenBankEdit.addTextChangedListener(new RegisterTextWatcher(mOpenBankEdit));

        mBranchBankEdit.addTextChangedListener(new RegisterTextWatcher(mBranchBankEdit));

        mBanknumEdit.addTextChangedListener(new RegisterTextWatcher(mBanknumEdit));

        mNameEdit.addTextChangedListener(new RegisterTextWatcher(mNameEdit));

        mPhonenumEdit.addTextChangedListener(new RegisterTextWatcher(mPhonenumEdit));

        mMerchantNameEdit.addTextChangedListener(new RegisterTextWatcher(mMerchantNameEdit));

        mProvinceSpinner.setOnItemSelectedListener(new RegisterOnItemSelectedListener());

        mCitySpinner.setOnItemSelectedListener(new RegisterOnItemSelectedListener());

        mOpenBankSpinner.setOnItemSelectedListener(new RegisterOnItemSelectedListener());

        mBranchBankSpinner.setOnItemSelectedListener(new RegisterOnItemSelectedListener());

        mCityEdit.setOnFocusChangeListener(new RegisterOnFocusChangeListener());

        mOpenBankEdit.setOnFocusChangeListener(new RegisterOnFocusChangeListener());

        mBranchBankEdit.setOnFocusChangeListener(new RegisterOnFocusChangeListener());

    }


    public void registerNext(View view) {
        if (!validate()) {
            return;
        }

        startLoading();

        if (SessonData.registerUser == null) {
            SessonData.registerUser = new User();
            SessonData.registerUser.setUsername(mRegisterSharedPreferences.getString("register_username", ""));
            SessonData.registerUser.setPassword(mRegisterSharedPreferences.getString("register_password", ""));
        }
        final User user = SessonData.registerUser;


        final String province = mProvinceEdit.getText().toString();
        user.setProvince(province);
        final String city = mCityEdit.getText().toString();
        user.setCity(city);

        final String bankopen = mOpenBankEdit.getText().toString();
        user.setBankOpen(bankopen);
        final String branchBank = mBranchBankEdit.getText().toString();
        user.setBranchBank(branchBank);

        //有些地方没有支行，get()会抛出outofindex异常
        int index = mBranchBankList.indexOf(branchBank);
        final String bankNo = (index != -1) ? mBankNoList.get(index) : "";
        user.setBankNo(bankNo);

        final String payee = mNameEdit.getText().toString();
        user.setPayee(payee);
        final String payeeCard = mBanknumEdit.getText().toString().replace(" ", "");
        user.setPayeeCard(payeeCard);
        final String phoneNum = mPhonenumEdit.getText().toString();
        user.setPhoneNum(phoneNum);
        final String merName = mMerchantNameEdit.getText().toString();
        user.setMerName(merName);

        quickPayService.updateUserAsync(user, new QuickPayCallbackListener<User>() {
            @Override
            public void onSuccess(final User data) {
                runOnUiThread(new Runnable() {
                    @Override
                    public void run() {
                        //NOTE:clientID也是merchantId,用于在七牛那边创建唯一id
                        String clientId = data.getClientid();
                        SessonData.registerUser.setClientid(clientId);

                        saveToRegisterSharedPreferences("register_clientid", clientId);
                        saveToRegisterSharedPreferences("register_province", province);
                        saveToRegisterSharedPreferences("register_city", city);
                        saveToRegisterSharedPreferences("register_bankopen", bankopen);
                        saveToRegisterSharedPreferences("register_branchbank", branchBank);
                        saveToRegisterSharedPreferences("register_bankno", bankNo);
                        saveToRegisterSharedPreferences("register_payee", payee);
                        saveToRegisterSharedPreferences("register_payeecard", payeeCard);
                        saveToRegisterSharedPreferences("register_phonenum", phoneNum);
                        saveToRegisterSharedPreferences("register_mername", merName);
                        saveToRegisterSharedPreferences("register_step_finish", 2);

                        endLoading();
                        Intent intent = new Intent(RegisterNextActivity.this, RegisterStep3Activity.class);
                        startActivity(intent);
                        finish();
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


    private void initProvinceData() {

        List<String> data = null;
        // TODO: 15-11-24 这里需要读取缓存,或者在 getProvince（）里的一个异步任务里做会更好呢？？！！
        if (data != null && data.size() != 0) {
            Log.d(TAG, "will use cache data to get province");
            updateProvinceAdapter(data);
        } else {
            bankDataService.getProvince(new ProvinceQuickPayCallbackListener());
        }
    }

    private void initBankData() {
        //获取bank的数据
        String data = readFromSharePreference("bank");
        if (data != null && data.length() != 0) {
            Log.d(TAG, "will use cache data to get bank");
            updateBankAdapter(data);
        } else {
            Log.d(TAG, "will do post to get bank data");
            RequestParam bankParam = BankBaseUtil.getBank();
            BankCommunicationListener bankCommunicationListener = new BankCommunicationListener();
            HttpCommunicationUtil.sendGetDataToServer(bankParam, bankCommunicationListener);
        }
    }

    private void initCityData(String province) {
        List<City> data = null;
        if (data != null && data.size() != 0) {
            Log.d(TAG, "will use cache data to get City");
            //updateCityAdapter(data);
        } else {
            Log.d(TAG, "will post to get City: " + data);
            bankDataService.getCity(province, new CityQuickPayCallbackListener());
        }
    }

    private void initBranchBankData(String cityCode, String bankId) {
        String key = cityCode + "_" + bankId;
        String data = readFromSharePreference(key);
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

    private void updateProvinceAdapter(List<String> data) {
        //这里直接得到的就是一个省份的list，不需要再去用json去解析了。
        mProvinceList.clear();
        mProvinceList.add(0, "开户行所在省份");
        mProvinceList.addAll(data);
        mProvinceAdapter.notifyDataSetChanged();
        mProvinceSearchAdapter.setData(mProvinceList);
        mProvinceSearchAdapter.notifyDataSetChanged();
    }

    private void updateCityAdapter(final List<City> data) {
        ArrayList<String> tempCityList = new ArrayList<String>();
        ArrayList<String> tempCityCodeList =  new ArrayList<String>();
        tempCityList.add(0, "开户行所在城市");
        tempCityCodeList.add(0, "");
        Iterator<City> it = data.iterator();
        while(it.hasNext()){
            City c = it.next();
            tempCityList.add(c.getCityName());//"city_name"这个要注意别弄成getCity（）了。
            tempCityCodeList.add(c.getCityCode());//"city_code"
        }
        mCityList.clear();
        mCityList.addAll(tempCityList);
        mCityCodeList.clear();
        mCityCodeList.addAll(tempCityCodeList);

        mCitySpinner.setSelection(0);
        mCityAdapter.notifyDataSetChanged();
        mCitySearchAdapter.setData(mCityList);
    }

    private void updateBankAdapter(final String data) {
        new Thread() {
            @Override
            public void run() {
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
        }.start();


    }

    private void updateBranchBankAdapter(final String data) {
        new Thread() {
            @Override
            public void run() {
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
        }.start();
    }

    private String readFromSharePreference(String key) {
        return mDataSharedPreferences.getString(key, "");
    }

    private void saveToSharePreferences(String result, String key) {
        SharedPreferences.Editor editor = mDataSharedPreferences.edit();
        editor.putString(key, result);
        editor.commit();
    }

    private void saveToRegisterSharedPreferences(String key, String value) {
        SharedPreferences.Editor editor = mRegisterSharedPreferences.edit();
        editor.putString(key, value);
        editor.commit();
    }

    private void saveToRegisterSharedPreferences(String key, int value) {
        SharedPreferences.Editor editor = mRegisterSharedPreferences.edit();
        editor.putInt(key, value);
        editor.commit();
    }

    @SuppressLint("NewApi")
    private boolean validate() {
        String province = mProvinceEdit.getText().toString();
        if (TextUtils.isEmpty(province)) {
            alertError("开户行所在省份不能为空!");
            return false;
        }

        String city = mCityEdit.getText().toString();
        if (TextUtils.isEmpty(city)) {
            alertError("开户行所在城市不能为空!");
            return false;
        }

        String openbank = mOpenBankEdit.getText().toString();
        if (TextUtils.isEmpty(openbank)) {
            alertError("开户行不能为空!");
            return false;
        }
        String branchbank = mBranchBankEdit.getText().toString();
        if (TextUtils.isEmpty(branchbank)) {
            if (mBranchBankList.size() == 1 && mBranchBankList.get(0).equals("请选择开户支行")) {
                //有些地方没有支行，这里不填写就不能下一步
            } else {
                alertError("开户支行不能为空!");
                return false;
            }
        }

        String name = mNameEdit.getText().toString().replace(" ", "");
        if (TextUtils.isEmpty(name)) {
            alertError("姓名不能为空!");
            return false;
        }

        String banknum = mBanknumEdit.getText().toString().replace(" ", "");
        if (TextUtils.isEmpty(banknum)) {
            alertError("银行卡号不能为空!");
            return false;
        }

        if (!VerifyUtil.checkBankCard(banknum)) {
            alertError("请输入正确的银行卡号!");
            return false;
        }

        String phonenum = mPhonenumEdit.getText().toString().replace(" ", "");
        if (TextUtils.isEmpty(phonenum)) {
            alertError("手机号不能为空!");
            return false;
        }
        if (!VerifyUtil.isMobileNO(phonenum)) {
            alertError("请输入正确的手机号!");
            return false;
        }

        String merchantname = mMerchantNameEdit.getText().toString().replace(" ", "");
        if (TextUtils.isEmpty(merchantname)) {
            alertError("请输入商店名称");
            return false;
        }
        return true;
    }

    //内部类，实现QuickPayCallbackListener接口
    private class ProvinceQuickPayCallbackListener implements QuickPayCallbackListener<List<String>> {

        @Override
        public void onSuccess(List<String> data) {
            // TODO: 15-11-24 save data to file or sqlite??!!
            updateProvinceAdapter(data);
        }

        @Override
        public void onFailure(QuickPayException ex) {

        }
    }

    //内部类，实现CommunicationListener接口,用来获取bank信息
    private class BankCommunicationListener implements CommunicationListener {
        @Override
        public void onResult(String result) {
            saveToSharePreferences(result, "bank");
            updateBankAdapter(result);
        }

        @Override
        public void onError(String error) {
            Log.i(TAG, "get bank data error:" + error);
        }
    }

    private class CityQuickPayCallbackListener implements QuickPayCallbackListener<List<City>> {

        @Override
        public void onSuccess(List<City> data) {
            updateCityAdapter(data);
        }

        @Override
        public void onFailure(QuickPayException ex) {

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
            saveToSharePreferences(result, key);
            updateBranchBankAdapter(result);
        }

        @Override
        public void onError(String error) {
            Log.i(TAG, "get branch bank error:" + error);
        }
    }


    private class RegisterTextWatcher implements TextWatcher {
        private View view;

        public RegisterTextWatcher(View view) {
            this.view = view;
        }

        @Override
        public void beforeTextChanged(CharSequence s, int start, int count, int after) {

        }

        @Override
        public void onTextChanged(CharSequence s, int start, int before, int count) {

        }

        @Override
        public void afterTextChanged(Editable s) {

            switch (view.getId()) {
                case R.id.edit_province:
                    //province
                    mCityEdit.setText("");//先把city的清空
                    saveToRegisterSharedPreferences("register_city", "");
                    String province = mProvinceEdit.getText().toString();
                    saveToRegisterSharedPreferences("register_province", province);
                    if (mProvinceList.indexOf(province) > 0) {
                        initCityData(province);
                    }//end if()
                    break;
                case R.id.edit_city:
                    //city
                    mOpenBankEdit.setText("");
                    saveToRegisterSharedPreferences("register_bankopen", "");
                    String city = mCityEdit.getText().toString();
                    saveToRegisterSharedPreferences("register_city", city);
                    if (mCityList.indexOf(city) > 0) {
                        initBankData();
                    }
                    break;
                case R.id.edit_openbank:
                    mBranchBankEdit.setText("");
                    saveToRegisterSharedPreferences("register_branchbank", "");
                    String openBank = mOpenBankEdit.getText().toString();
                    city = mCityEdit.getText().toString();
                    if (mOpenBankList.indexOf(openBank) > 0) {

                        int indexOfCity = mCityList.indexOf(city);
                        int indexOfBank = mOpenBankList.indexOf(openBank);
                        if (indexOfCity > 0 && indexOfBank > 0) {
                            String cityCode = mCityCodeList.get(indexOfCity);
                            String bankId = mBankIdList.get(indexOfBank);

                            saveToRegisterSharedPreferences("register_bankopen", openBank);
                            initBranchBankData(cityCode, bankId);
                        }
                    }
                    break;
                case R.id.edit_branchbank:
                    String branchBank = mBranchBankEdit.getText().toString();
                    saveToRegisterSharedPreferences("register_branchbank", branchBank);
                    break;
                case R.id.info_name:
                    String payee = mNameEdit.getText().toString();
                    saveToRegisterSharedPreferences("register_payee", payee);
                    break;
                case R.id.info_banknum:
                    String payeeCard = mBanknumEdit.getText().toString();
                    saveToRegisterSharedPreferences("register_payeecard", payeeCard);
                    break;
                case R.id.info_phonenum:
                    String phoneNum = mPhonenumEdit.getText().toString();
                    saveToRegisterSharedPreferences("register_phonenum", phoneNum);
                    break;
                case R.id.info_merchantname:
                    String merName = mMerchantNameEdit.getText().toString();
                    saveToRegisterSharedPreferences("register_mername", merName);
                    break;
            }

        }
    }

    private class RegisterOnItemSelectedListener implements OnItemSelectedListener {

        @Override
        public void onItemSelected(AdapterView<?> parent, View view, int position, long id) {
            switch (parent.getId()) {
                case R.id.spinner_province:
                    if (position > 0) {
                        mProvinceEdit.setText(mProvinceList.get(position));
                    }
                    break;
                case R.id.spinner_city:
                    if (position > 0) {
                        mCityEdit.setText(mCityList.get(position));
                    }
                    break;
                case R.id.spinner_openbank:
                    if (position > 0) {
                        mOpenBankEdit.setText(mOpenBankList.get(position));
                    }
                    break;
                case R.id.spinner_branchbank:
                    if (position > 0) {
                        mBranchBankEdit.setText(mBranchBankList.get(position));
                    }
                    break;
            }
        }

        @Override
        public void onNothingSelected(AdapterView<?> parent) {

        }
    }

    private class RegisterOnFocusChangeListener implements OnFocusChangeListener {

        @Override
        public void onFocusChange(View view, boolean hasFocus) {
            if (!hasFocus) {
                return;
            }
            switch (view.getId()) {
                case R.id.edit_city:
                    if (mProvinceList.indexOf(mProvinceEdit.getText().toString()) < 0) {
                        mProvinceEdit.setText("");
                    }
                    break;
                case R.id.edit_openbank:
                    if (mCityList.indexOf(mCityEdit.getText().toString()) < 0) {
                        mCityEdit.setText("");
                    }
                    break;
                case R.id.edit_branchbank:
                    if (mOpenBankList.indexOf(mOpenBankEdit.getText().toString()) < 0) {
                        mOpenBankEdit.setText("");
                    }
                    break;
                case R.id.edit_province:
                    break;
            }
        }
    }
}
