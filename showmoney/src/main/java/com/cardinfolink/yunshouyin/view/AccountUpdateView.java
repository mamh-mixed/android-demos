package com.cardinfolink.yunshouyin.view;

import android.annotation.SuppressLint;
import android.app.Activity;
import android.content.Context;
import android.graphics.BitmapFactory;
import android.text.Editable;
import android.text.TextWatcher;
import android.util.Log;
import android.view.LayoutInflater;
import android.view.View;
import android.widget.AdapterView;
import android.widget.AdapterView.OnItemSelectedListener;
import android.widget.ArrayAdapter;
import android.widget.AutoCompleteTextView;
import android.widget.EditText;
import android.widget.LinearLayout;
import android.widget.Spinner;

import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.activity.BaseActivity;
import com.cardinfolink.yunshouyin.data.SessonData;
import com.cardinfolink.yunshouyin.data.User;
import com.cardinfolink.yunshouyin.util.BankBaseUtil;
import com.cardinfolink.yunshouyin.util.CommunicationListener;
import com.cardinfolink.yunshouyin.util.ErrorUtil;
import com.cardinfolink.yunshouyin.util.HttpCommunicationUtil;
import com.cardinfolink.yunshouyin.util.JsonUtil;
import com.cardinfolink.yunshouyin.util.ParamsUtil;
import com.cardinfolink.yunshouyin.util.VerifyUtil;

import org.json.JSONArray;
import org.json.JSONException;
import org.json.JSONObject;

import java.util.ArrayList;
import java.util.Iterator;
import java.util.List;

public class AccountUpdateView extends LinearLayout {

    private EditText mNameEdit;
    private EditText mBanknumEdit;
    private EditText mPhonenumEdit;
    private Context mContext;
    private BaseActivity mBaseActivity;

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

    private boolean isInit = false;

    public AccountUpdateView(Context context) {
        super(context);
        mContext = context;
        mBaseActivity = (BaseActivity) mContext;
        View contentView = LayoutInflater.from(context).inflate(
                R.layout.account_update_view, null);
        LinearLayout.LayoutParams layoutParams = new LayoutParams(
                LayoutParams.MATCH_PARENT, LayoutParams.MATCH_PARENT);
        contentView.setLayoutParams(layoutParams);
        addView(contentView);
        initLayout(contentView);
        initData();
        //initListener();
        //	initListener();
        // getInfo();

        contentView.findViewById(R.id.btn_submit).setOnClickListener(
                new OnClickListener() {

                    @Override
                    public void onClick(View v) {
                        finishedOnClick();
                    }
                });

    }

    private void initLayout(View contentView) {

        mNameEdit = (EditText) contentView.findViewById(R.id.info_name);
        mBanknumEdit = (EditText) contentView.findViewById(R.id.info_banknum);
        mPhonenumEdit = (EditText) contentView.findViewById(R.id.info_phonenum);
        VerifyUtil.bankCardNumAddSpace(mBanknumEdit);

        VerifyUtil.bankCardNumAddSpace(mBanknumEdit);

        mProvinceEdit = (AutoCompleteTextView) contentView
                .findViewById(R.id.edit_province);
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

        mCityEdit = (AutoCompleteTextView) contentView
                .findViewById(R.id.edit_city);
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

        mOpenBankEdit = (AutoCompleteTextView) contentView
                .findViewById(R.id.edit_openbank);
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

        mBranchBankEdit = (AutoCompleteTextView) contentView
                .findViewById(R.id.edit_branchbank);
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

    public void getInfo() {

        HttpCommunicationUtil.sendDataToServer(ParamsUtil.getInfo(SessonData.loginUser),
                new CommunicationListener() {

                    @Override
                    public void onResult(final String result) {
                        if (JsonUtil.getParam(result, "state").equals("success")) {
                            ((Activity) mContext).runOnUiThread(new Runnable() {

                                @Override
                                public void run() {
                                    String info = JsonUtil.getParam(result, "info");

                                    info_province = JsonUtil.getParam(info, "province");
                                    info_city = JsonUtil.getParam(info, "city");
                                    info_openbank = JsonUtil.getParam(info, "bank_open");
                                    info_branch_bank = JsonUtil.getParam(info, "branch_bank");

                                    mProvinceEdit.setText(info_province);
                                    mCityEdit.setText(info_city);
                                    mOpenBankEdit.setText(info_openbank);
                                    mBranchBankEdit.setText(info_branch_bank);

                                    mNameEdit.setText(JsonUtil.getParam(info, "bank_open"));
                                    mNameEdit.setText(JsonUtil.getParam(info, "payee"));
                                    mBanknumEdit.setText(JsonUtil.getParam(info, "payee_card"));
                                    mPhonenumEdit.setText(JsonUtil.getParam(info, "phone_num"));


                                    mCityList.clear();
                                    mCityCodeList.clear();
                                    mCityList.add("开户行所在城市");
                                    mCityCodeList.add("");
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

                                                                        mCityAdapter
                                                                                .notifyDataSetChanged();
                                                                        mCitySearchAdapter
                                                                                .setData(mCityList);


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


                                                                                                mBranchBankList.clear();
                                                                                                mBranchBankList.add("请选择开户支行");
                                                                                                mBankNoList.clear();
                                                                                                mBankNoList.add("行号");
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

                                                                                                                                    mBranchBankAdapter
                                                                                                                                            .notifyDataSetChanged();
                                                                                                                                    mBranchBankSearchAdapter
                                                                                                                                            .setData(mBranchBankList);

                                                                                                                                    initListener();
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

                                                                                    }

                                                                                    @Override
                                                                                    public void onError(String error) {
                                                                                        Log.i("opp", "error:" + error);

                                                                                    }
                                                                                });


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

                        }

                    }

                    @Override
                    public void onError(String error) {
                        // TODO Auto-generated method stub

                    }
                });


        new Thread(new Runnable() {

            @Override
            public void run() {
                try {
                    Thread.sleep(4000);
                } catch (InterruptedException e) {
                    // TODO Auto-generated catch block
                    e.printStackTrace();
                }
                initListener();

            }
        }).start();


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


    }

    private void initListener() {

        if (isInit) {


            return;
        }
        isInit = true;
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
                } else {
                    if (mBranchBankList.indexOf(mBranchBankEdit.getText().toString()) < 0) {
                        mBranchBankEdit.setText("");
                    }
                }

            }
        });


    }

    public void finishedOnClick() {
        if (validate()) {
            mBaseActivity.startLoading();
            User user = new User();
            user.setUsername(SessonData.loginUser.getUsername());
            user.setPassword(SessonData.loginUser.getPassword());
            // user.setBankOpen(mOpenBankEdit.getText().toString());
            user.setProvince(mProvinceEdit.getText().toString());
            user.setBankOpen(mOpenBankEdit.getText().toString());
            user.setCity(mCityEdit.getText().toString());
            user.setBranchBank(mBranchBankEdit.getText().toString());
            int index = mBranchBankList.indexOf(mBranchBankEdit.getText().toString());
            if (index >= 0) {

                user.setBankNo(mBankNoList.get(index));
            } else {
                user.setBankNo("");
            }


            user.setPayee(mNameEdit.getText().toString());
            user.setPayeeCard(mBanknumEdit.getText().toString()
                    .replace(" ", ""));
            user.setPhoneNum(mPhonenumEdit.getText().toString());
            HttpCommunicationUtil.sendDataToServer(
                    ParamsUtil.getUpdateInfo(user),
                    new CommunicationListener() {

                        @Override
                        public void onResult(final String result) {
                            String state = JsonUtil.getParam(result, "state");
                            if (state.equals("success")) {
                                ((Activity) mContext)
                                        .runOnUiThread(new Runnable() {

                                            @Override
                                            public void run() {
                                                // 更新UI
                                                mBaseActivity.endLoading();
                                                AlertDialog alert_Dialog = new AlertDialog(
                                                        mContext,
                                                        null,
                                                        ((Activity) mContext)
                                                                .findViewById(R.id.alert_dialog),
                                                        getResources().getString(R.string.alert_update_success),
                                                        BitmapFactory
                                                                .decodeResource(
                                                                        mContext.getResources(),
                                                                        R.drawable.right));
                                                alert_Dialog.show();
                                            }

                                        });

                            } else {
                                ((Activity) mContext)
                                        .runOnUiThread(new Runnable() {

                                            @Override
                                            public void run() {
                                                // 更新UI
                                                mBaseActivity.endLoading();
                                                mBaseActivity.alertShow(
                                                        ErrorUtil.getErrorString(JsonUtil.getParam(result, "error")),
                                                        BitmapFactory
                                                                .decodeResource(
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
                                    mBaseActivity.alertShow(error,
                                            BitmapFactory.decodeResource(
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
            mBaseActivity.alertShow("开户行所在省份不能为空!", BitmapFactory
                    .decodeResource(this.getResources(), R.drawable.wrong));
            return false;
        }

        if (mCityEdit.getText().toString().isEmpty()) {
            mBaseActivity.alertShow("开户行所在城市不能为空!", BitmapFactory
                    .decodeResource(this.getResources(), R.drawable.wrong));
            return false;
        }

        if (mOpenBankEdit.getText().toString().isEmpty()) {
            mBaseActivity.alertShow("开户行不能为空!", BitmapFactory.decodeResource(
                    this.getResources(), R.drawable.wrong));
            return false;
        }

        if (mBranchBankEdit.getText().toString().isEmpty()) {
            mBaseActivity.alertShow("开户支行不能为空!", BitmapFactory.decodeResource(
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
            mBaseActivity.alertShow("请输入正确的银行卡号!", BitmapFactory
                    .decodeResource(this.getResources(), R.drawable.wrong));
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
