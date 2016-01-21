package com.cardinfolink.yunshouyin.activity;

import android.content.Context;
import android.content.Intent;
import android.os.Bundle;
import android.text.Editable;
import android.text.InputFilter;
import android.text.InputType;
import android.text.TextUtils;
import android.text.TextWatcher;
import android.view.View;
import android.widget.Button;
import android.widget.TextView;

import com.cardinfolink.cashiersdk.model.InitData;
import com.cardinfolink.cashiersdk.sdk.CashierSdk;
import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.api.QuickPayException;
import com.cardinfolink.yunshouyin.constant.SystemConfig;
import com.cardinfolink.yunshouyin.core.QuickPayCallbackListener;
import com.cardinfolink.yunshouyin.data.SessionData;
import com.cardinfolink.yunshouyin.data.User;
import com.cardinfolink.yunshouyin.model.Bank;
import com.cardinfolink.yunshouyin.model.City;
import com.cardinfolink.yunshouyin.model.Province;
import com.cardinfolink.yunshouyin.model.SubBank;
import com.cardinfolink.yunshouyin.ui.SettingActionBarItem;
import com.cardinfolink.yunshouyin.ui.SettingClikcItem;
import com.cardinfolink.yunshouyin.ui.SettingInputItem;
import com.cardinfolink.yunshouyin.util.Log;
import com.cardinfolink.yunshouyin.util.ShowMoneyApp;
import com.cardinfolink.yunshouyin.util.VerifyUtil;
import com.cardinfolink.yunshouyin.view.SelectDialog;
import com.cardinfolink.yunshouyin.view.YellowTips;

import java.util.ArrayList;
import java.util.Hashtable;
import java.util.List;
import java.util.Map;

import kankan.wheel.widget.OnWheelScrollListener;
import kankan.wheel.widget.WheelView;
import kankan.wheel.widget.adapters.AbstractWheelTextAdapter;


public class RegisterNextActivity extends BaseActivity implements View.OnClickListener {
    private static final String TAG = "RegisterNextActivity";
    private static final int MAX_BANK_NUMBER_LENGTH = 23;//银行卡号的最大允许的长度
    private static final int MAX_PHONE_NUMBER_LENTH = 11;//大陆手机号的最大长度

    private SettingActionBarItem mActionBar;

    private SettingClikcItem mSetProvinceCity;//点击去设置省市信息
    private SettingClikcItem mSetBank;//点击设置银行信息，主行，或分行
    private SettingInputItem mName;//姓名
    private SettingInputItem mBankNumber;//银行卡号
    private SettingInputItem mPhone;//手机号
    private TextView mAgreement;
    private Button mRegisterFinished;//注册按钮，

    private SelectDialog selectDialog;

    private List<Province> provinceList = new ArrayList<Province>();//省份的list
    private List<City> cityList = new ArrayList<City>();//城市的list
    private Map<Province, List<City>> provinceCityMap = new Hashtable<>();

    private String mProvinceName;//保存设置的省份，调用银行时候会检查这个是否有值
    private String mCityCode;
    private String mCityName;

    private List<Bank> bankList = new ArrayList<Bank>();//Bank的list
    private List<SubBank> subbankList = new ArrayList<SubBank>();//SubBank 的list
    private String mBankName;
    private String mSubBankName;
    private String mBankNo;//sb.getOneBankNo() + "|" + sb.getTwoBankNo()
    private static final String SEPARATOR = "__";//分隔符

    //这个maps的可以是 mCityCode + SEPARATOR + currentBank.getBankName()
    private Map<String, List<SubBank>> bankSubBankMap = new Hashtable<>();

    private YellowTips mYellowTips;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_register_next);

        mActionBar = (SettingActionBarItem) findViewById(R.id.action_bar);
        mActionBar.setLeftTextOnclickListner(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                finish();
            }
        });
        mRegisterFinished = (Button) findViewById(R.id.btnregister);

        mSetProvinceCity = (SettingClikcItem) findViewById(R.id.province_city);
        mSetBank = (SettingClikcItem) findViewById(R.id.bank);

        mName = (SettingInputItem) findViewById(R.id.name);//姓名
        mName.setImageViewDrawable(null);
        mBankNumber = (SettingInputItem) findViewById(R.id.bank_number);//银行卡号
        mBankNumber.setInputType(InputType.TYPE_CLASS_NUMBER);//限制银行卡号输入法只能是数字
        mBankNumber.setTextFilters(new InputFilter[]{new InputFilter.LengthFilter(MAX_BANK_NUMBER_LENGTH)});//限制输入长度
        VerifyUtil.bankCardNumAddSpace(mBankNumber.getmText());
        mBankNumber.setImageViewDrawable(null);

        mPhone = (SettingInputItem) findViewById(R.id.phone_number);//手机号
        mPhone.setImageViewDrawable(null);
        mPhone.setInputType(InputType.TYPE_CLASS_PHONE);
        mPhone.setTextFilters(new InputFilter[]{new InputFilter.LengthFilter(MAX_PHONE_NUMBER_LENTH)});//限制输入长度

        mAgreement = (TextView) findViewById(R.id.tv_agreement);

        mSetProvinceCity.setOnClickListener(this);
        mSetBank.setOnClickListener(this);
        mRegisterFinished.setOnClickListener(this);
        mAgreement.setOnClickListener(this);

        selectDialog = new SelectDialog(this, findViewById(R.id.select_dialog));

        mYellowTips = new YellowTips(this, findViewById(R.id.yellow_tips));
    }


    @Override
    public void onClick(View v) {
        switch (v.getId()) {
            case R.id.btnregister:
                Log.e(TAG, " onclick register");
                register();//调用注册的方法
                break;
            case R.id.province_city:
                Log.e(TAG, " onclick province_city");
                showProvinceCity();
                break;
            case R.id.bank:
                Log.e(TAG, " onclick bank");
                showBankSubBank();
                break;
            case R.id.tv_agreement:
                Intent intent = new Intent(RegisterNextActivity.this, AgreementActivity.class);
                startActivity(intent);
                break;
        }
    }

    private void register() {
        String name = mName.getText().replace(" ", ""); //姓名
        String bankCardNum = mBankNumber.getText().replace(" ", ""); //银行卡号
        String phonenum = mPhone.getText().replace(" ", "");

        String province = mProvinceName;
        String city = mCityName;

        String bank = mBankName;
        String subbank = mSubBankName;
        String bankNo = mBankNo;

        if (!validate(name, bankCardNum, phonenum, province, city, bank, subbank)) {
            return;
        }

        mLoadingDialog.startLoading();
        User user = new User();
        user.setUsername(SessionData.loginUser.getUsername());
        user.setPassword(SessionData.loginUser.getPassword());
        user.setProvince(province);
        user.setCity(city);
        user.setBankOpen(bank);
        if (TextUtils.isEmpty(subbank)) {
            subbank = " ";
        }
        user.setBranchBank(subbank);
        if (TextUtils.isEmpty(bankNo)) {
            bankNo = " ";
        }
        user.setBankNo(bankNo);//注意这里 是两个 bankNo 拼接的
        user.setPayee(name);//姓名
        user.setPayeeCard(bankCardNum);//银行卡号
        user.setPhoneNum(phonenum);
        quickPayService.improveInfoAsync(user, new QuickPayCallbackListener<User>() {
            @Override
            public void onSuccess(User data) {
                SessionData.loginUser.setClientid(data.getClientid());
                SessionData.loginUser.setObjectId(data.getObjectId());
                SessionData.loginUser.setLimit(data.getLimit());

                InitData initData = new InitData();
                initData.setMchntid(data.getClientid());    // 商户号
                initData.setInscd(data.getInscd());         // 机构号
                initData.setSignKey(data.getSignKey());     // 秘钥
                initData.setTerminalid(data.getUsername());// 设备号
                initData.setIsProduce(SystemConfig.IS_PRODUCE);// 是否生产环境

                CashierSdk.init(initData);
                Intent intent = new Intent(RegisterNextActivity.this, RegisterFinalActivity.class);
                startActivity(intent);
                finish();
            }

            @Override
            public void onFailure(QuickPayException ex) {
                String errorMsg = ex.getErrorMsg();
                mLoadingDialog.endLoading();
                mYellowTips.show(errorMsg);
            }
        });

    }

    private boolean validate(String name, String banknum, String phonenum,
                             String province, String city, String bank, String subbank) {
        //注意这里 subbank没有检查，因为支行可能是空的

        String alertMsg = "";
        if (TextUtils.isEmpty(province)) {
            alertMsg = ShowMoneyApp.getResString(R.string.alert_error_province_cannot_empty);
            mYellowTips.show(alertMsg);
            return false;
        }

        if (TextUtils.isEmpty(city)) {
            alertMsg = ShowMoneyApp.getResString(R.string.alert_error_city_cannot_empty);
            mYellowTips.show(alertMsg);
            return false;
        }

        if (TextUtils.isEmpty(bank)) {
            alertMsg = ShowMoneyApp.getResString(R.string.alert_error_bank_cannot_empty);
            mYellowTips.show(alertMsg);
            return false;
        }

        if (TextUtils.isEmpty(name)) {
            alertMsg = ShowMoneyApp.getResString(R.string.alert_error_name_cannot_empty);
            mYellowTips.show(alertMsg);
            return false;
        }

        if (TextUtils.isEmpty(banknum)) {
            alertMsg = ShowMoneyApp.getResString(R.string.alert_error_banknum_cannot_empty);
            mYellowTips.show(alertMsg);
            return false;
        }

        if (!VerifyUtil.checkBankCard(banknum)) {
            alertMsg = ShowMoneyApp.getResString(R.string.alert_error_banknum_format_error);
            mYellowTips.show(alertMsg);
            return false;
        }

        if (TextUtils.isEmpty(phonenum)) {
            alertMsg = ShowMoneyApp.getResString(R.string.alert_error_phonenum_cannot_empty);
            mYellowTips.show(alertMsg);
            return false;
        }

        if (!VerifyUtil.isMobileNO(phonenum)) {
            alertMsg = ShowMoneyApp.getResString(R.string.alert_error_phonenum_format_error);
            mYellowTips.show(alertMsg);
            return false;
        }

        return true;
    }


    /**
     * 显示银行 银行支行 滚轮的界面
     */
    private void showBankSubBank() {
        if (TextUtils.isEmpty(mProvinceName)) {
            String msg = getResources().getString(R.string.alert_error_province_cannot_empty);
            mYellowTips.show(msg);
            return;
        }
        if (TextUtils.isEmpty(mCityCode)) {
            String msg = getResources().getString(R.string.alert_error_city_cannot_empty);
            mYellowTips.show(msg);
            return;
        }


        updateBankData();//去获取银行信息

        selectDialog.setSearchText("");

        selectDialog.addLeftScrollingListener(new BankOnWheelScrollListener());

        selectDialog.setOkOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {

                int bankIndex = selectDialog.getWheelLeftCurrentItem();
                int subbankIndex = selectDialog.getWheelRightCurrentItem();

                String bankName = bankList.get(bankIndex).getBankName();
                mSetBank.setTitle(bankName);
                mBankName = bankName;


                List<SubBank> list = bankSubBankMap.get(mCityCode + SEPARATOR + bankName);
                if (list != null && subbankIndex >= 0 && subbankIndex < list.size()) {
                    SubBank subBank = list.get(subbankIndex);
                    String subbankName = subBank.getBankName();
                    mSetBank.setRightText(subbankName);
                    mSubBankName = subbankName;
                    mBankNo = subBank.getOneBankNo() + "|" + subBank.getTwoBankNo();
                } else {
                    mSetBank.setRightText("");
                    mSubBankName = null;
                    mBankNo = null;
                }

                selectDialog.hide();
            }
        });

        selectDialog.setCancelOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                selectDialog.hide();
            }
        });

        selectDialog.addSearchTextChangedListener(new TextWatcher() {
            @Override
            public void beforeTextChanged(CharSequence s, int start, int count, int after) {

            }

            @Override
            public void onTextChanged(CharSequence s, int start, int before, int count) {

            }

            @Override
            public void afterTextChanged(Editable s) {
                searchBank();
            }
        });

        selectDialog.show();
    }

    /**
     * 模糊搜索银行
     */
    private void searchBank() {
        try {
            String bank = selectDialog.getSearchText();
            if (TextUtils.isEmpty(bank)) {
                return;
            }
            try {
                for (int i = 0; i <= bankList.size(); i++) {
                    String bankName = bankList.get(i).getBankName();
                    if (bankName.contains(bank)) {
                        selectDialog.setWheelLeftCurrentItem(i);
                        updateSubBankData(bankList.get(i));
                        bank = bankName;
                        break;
                    }
                }
            } catch (Exception e) {

            }

            try {
                for (int i = 0; i <= subbankList.size(); i++) {
                    String bankName = subbankList.get(i).getBankName();
                    if (bankName.contains(bank)) {
                        selectDialog.setWheelRightCurrentItem(i);
                        break;
                    }
                }
            } catch (Exception e) {

            }
        } catch (Exception e) {
            e.printStackTrace();
        }
    }

    /**
     * 模糊搜索省份
     */
    private void searchProvince() {
        try {
            String provice = selectDialog.getSearchText();
            if (TextUtils.isEmpty(provice)) {
                return;
            }
            try {
                for (int i = 0; i <= provinceList.size(); i++) {
                    String proviceName = provinceList.get(i).getProvinceName();
                    if (proviceName.contains(provice)) {
                        selectDialog.setWheelLeftCurrentItem(i);
                        updateCityData(provinceList.get(i));
                        provice = proviceName;
                        break;
                    }
                }
            } catch (Exception e) {
                e.printStackTrace();
            }

            try {
                for (int i = 0; i <= cityList.size(); i++) {
                    String cityName = cityList.get(i).getCityName();
                    if (cityName.contains(provice)) {
                        selectDialog.setWheelRightCurrentItem(i);
                        break;
                    }
                }
            } catch (Exception e) {
                e.printStackTrace();
            }
        } catch (Exception e) {
            e.printStackTrace();
        }
    }

    /**
     * 这里调用显示选择 省份城市 滚轮 的界面
     */
    private void showProvinceCity() {

        updateProvinceData();

        selectDialog.addLeftScrollingListener(new ProvinceOnWheelScrollListener());

        //点击确定按钮 就把滚轮当前的 值设置到相应的位置
        selectDialog.setOkOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {

                int provinceIndex = selectDialog.getWheelLeftCurrentItem();
                int cityIndex = selectDialog.getWheelRightCurrentItem();
                String provinceName = provinceList.get(provinceIndex).getProvinceName();

                //省份是肯定有的
                mSetProvinceCity.setTitle(provinceName);
                mProvinceName = provinceName;//再点击ok按钮的时候把省份名字保存一下

                //城市码不一定刷新出来，然后用户就按确定了，这时候城市就是空的。存到一个hashtable里
                List<City> list = provinceCityMap.get(provinceList.get(provinceIndex));
                if (list != null && cityIndex >= 0 && cityIndex < list.size()) {
                    String cityName = list.get(cityIndex).getCityName();
                    mCityCode = list.get(cityIndex).getCityCode();//保存一下citycode，之后获取银行信息的时候会用到
                    mCityName = cityName;
                    mSetProvinceCity.setRightText(cityName);
                } else {
                    mSetProvinceCity.setRightText("");
                    mCityCode = null;//这里设置为空
                    mCityName = null;
                }

                mSetBank.setTitle(getResources().getString(R.string.register_bank_branch_bank));
                mSetBank.setRightText("");
                mBankName = "";
                mSubBankName = "";
                selectDialog.hide();
            }
        });

        selectDialog.setCancelOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                selectDialog.hide();
            }
        });

        selectDialog.addSearchTextChangedListener(new TextWatcher() {
            @Override
            public void beforeTextChanged(CharSequence s, int start, int count, int after) {

            }

            @Override
            public void onTextChanged(CharSequence s, int start, int before, int count) {
            }

            @Override
            public void afterTextChanged(Editable s) {
                searchProvince();
            }
        });

        //调用显示对话框
        selectDialog.show();
    }


    private void updateProvinceData() {
        //这里调用 获取省份的信息。成功之后调用获取城市的信息。
        bankDataService.getProvince(new QuickPayCallbackListener<List<Province>>() {
            @Override
            public void onSuccess(List<Province> data) {
                provinceList.clear();
                provinceList.addAll(data);

                selectDialog.setWheelLeftAdapter(new RegisterWheelAdapter<Province>(mContext, provinceList));

                selectDialog.setWheelLeftCurrentItem(0);
                Province currentProvince = data.get(0);

                //调用获取城市信息，并或刷新相应的data和view
                updateCityData(currentProvince);
            }

            @Override
            public void onFailure(QuickPayException ex) {

            }
        });
    }


    /**
     * 调用更新城市的方法，传人一个省份，得到这个省份下面所有的城市
     */
    private void updateCityData(final Province currentProvince) {
        bankDataService.getCity(currentProvince.getProvinceName(), new QuickPayCallbackListener<List<City>>() {
            @Override
            public void onSuccess(List<City> data) {
                cityList.clear();
                cityList.addAll(data);
                provinceCityMap.put(currentProvince, data);
                selectDialog.setWheelRigthAdapter(new RegisterWheelAdapter<City>(mContext, cityList));
                selectDialog.setWheelRightCurrentItem(0);
            }

            @Override
            public void onFailure(QuickPayException ex) {

            }
        });
    }


    private void updateBankData() {
        bankDataService.getBank(new QuickPayCallbackListener<List<Bank>>() {
            @Override
            public void onSuccess(List<Bank> data) {
                bankList.clear();
                bankList.addAll(data);
                selectDialog.setWheelLeftAdapter(new RegisterWheelAdapter<Bank>(mContext, bankList));

                selectDialog.setWheelLeftCurrentItem(0);
                Bank currentBank = data.get(0);

                //调用update subbank的方法
                updateSubBankData(currentBank);
            }

            @Override
            public void onFailure(QuickPayException ex) {

            }
        });
    }

    private void updateSubBankData(final Bank currentBank) {
        bankDataService.getBranchBank(mCityCode, currentBank.getId(), new QuickPayCallbackListener<List<SubBank>>() {
            @Override
            public void onSuccess(List<SubBank> data) {
                subbankList.clear();
                subbankList.addAll(data);
                bankSubBankMap.put(mCityCode + SEPARATOR + currentBank.getBankName(), data);//存入map字典中
                selectDialog.setWheelRigthAdapter(new RegisterWheelAdapter<SubBank>(mContext, subbankList));
                selectDialog.setWheelRightCurrentItem(0);
            }

            @Override
            public void onFailure(QuickPayException ex) {

            }
        });
    }

    /**
     * 左边Province 滚轮 滑动 事件 的监听类
     */
    private class ProvinceOnWheelScrollListener implements OnWheelScrollListener {
        @Override
        public void onScrollingStarted(WheelView wheel) {
            cityList.clear();//开始滚动就清除城市列表
        }

        @Override
        public void onScrollingFinished(WheelView wheel) {
            //滚动结束调用跟新城市的方法
            Province currentProvince = provinceList.get(wheel.getCurrentItem());
            if (provinceCityMap != null && provinceCityMap.get(currentProvince) != null) {
                selectDialog.setWheelRigthAdapter(new RegisterWheelAdapter<City>(mContext, provinceCityMap.get(currentProvince)));
            } else {
                updateCityData(currentProvince);
            }
        }
    }

    /**
     * 当左边显示的是银行信息的时候会添加这个监听事件
     * 左边bak 滚轮 滑动 事件 的监听类
     */
    private class BankOnWheelScrollListener implements OnWheelScrollListener {
        @Override
        public void onScrollingStarted(WheelView wheel) {
            subbankList.clear();
        }

        @Override
        public void onScrollingFinished(WheelView wheel) {
            Bank currentBank = bankList.get(wheel.getCurrentItem());
            String mapKey = mCityCode + SEPARATOR + currentBank.getBankName();
            if (bankSubBankMap != null && bankSubBankMap.get(mapKey) != null) {
                selectDialog.setWheelRigthAdapter(new RegisterWheelAdapter<SubBank>(mContext, bankSubBankMap.get(mapKey)));
            } else {
                updateSubBankData(currentBank);
            }
        }
    }


    /**
     * 滚轮组件的适配器，泛型的，里面存了一个 list，构建方法的时候传人这个list
     *
     * @param <T>
     */
    private class RegisterWheelAdapter<T> extends AbstractWheelTextAdapter {

        // items
        private List<T> items;

        /**
         * Constructor
         *
         * @param context the current context
         * @param items   the items
         */
        public RegisterWheelAdapter(Context context, List<T> items) {
            super(context);
            this.items = items;
            setTextSize(8);
        }

        @Override
        public CharSequence getItemText(int index) {
            if (index >= 0 && index < items.size()) {
                T item = items.get(index);
                if (item instanceof CharSequence) {
                    return (CharSequence) item;
                }
                if (item instanceof Province) {
                    return ((Province) item).getProvinceName();
                }
                if (item instanceof City) {
                    return ((City) item).getCityName();
                }
                if (item instanceof Bank) {
                    return ((Bank) item).getBankName();
                }
                if (item instanceof SubBank) {
                    return ((SubBank) item).getBankName();
                }
                return item.toString();
            }
            return null;
        }

        @Override
        public int getItemsCount() {
            if (items != null) {
                return items.size();
            } else {
                return 0;
            }
        }
    }
}
