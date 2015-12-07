package com.cardinfolink.yunshouyin.view;

import android.content.Context;
import android.graphics.Bitmap;
import android.graphics.BitmapFactory;
import android.text.Editable;
import android.text.TextUtils;
import android.text.TextWatcher;
import android.view.LayoutInflater;
import android.view.View;
import android.widget.AdapterView;
import android.widget.AdapterView.OnItemSelectedListener;
import android.widget.ArrayAdapter;
import android.widget.AutoCompleteTextView;
import android.widget.Button;
import android.widget.EditText;
import android.widget.LinearLayout;
import android.widget.Spinner;

import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.activity.BaseActivity;
import com.cardinfolink.yunshouyin.adapter.SearchAdapter;
import com.cardinfolink.yunshouyin.api.QuickPayException;
import com.cardinfolink.yunshouyin.core.BankDataService;
import com.cardinfolink.yunshouyin.core.QuickPayCallbackListener;
import com.cardinfolink.yunshouyin.core.QuickPayService;
import com.cardinfolink.yunshouyin.data.SessonData;
import com.cardinfolink.yunshouyin.data.User;
import com.cardinfolink.yunshouyin.model.Bank;
import com.cardinfolink.yunshouyin.model.BankInfo;
import com.cardinfolink.yunshouyin.model.City;
import com.cardinfolink.yunshouyin.model.Province;
import com.cardinfolink.yunshouyin.model.SubBank;
import com.cardinfolink.yunshouyin.util.ShowMoneyApp;
import com.cardinfolink.yunshouyin.util.VerifyUtil;

import java.util.ArrayList;
import java.util.Iterator;
import java.util.List;

public class AccountUpdateView extends LinearLayout {
    private static final String TAG = "AccountUpdateView";
    private Context mContext;
    private BaseActivity mBaseActivity;

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

    private Button mSubmitButton;

    private BankDataService bankDataService;
    private QuickPayService quickPayService;

    public AccountUpdateView(Context context) {
        super(context);
        mContext = context;
        //此类不是个Activity的子类，所有弄个这个来调用activity里的一些个方法
        mBaseActivity = (BaseActivity) mContext;

        bankDataService = ShowMoneyApp.getInstance().getBankDataService();
        quickPayService = ShowMoneyApp.getInstance().getQuickPayService();

        View contentView = LayoutInflater.from(context).inflate(R.layout.account_update_view, null);
        LinearLayout.LayoutParams layoutParams = new LayoutParams(LayoutParams.MATCH_PARENT, LayoutParams.MATCH_PARENT);
        contentView.setLayoutParams(layoutParams);
        addView(contentView);
        initLayout(contentView);
        initData();
        initListener();

        mSubmitButton = (Button) contentView.findViewById(R.id.btn_submit);
        mSubmitButton.setOnClickListener(new OnClickListener() {
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

        mProvinceEdit = (AutoCompleteTextView) contentView.findViewById(R.id.edit_province);
        mProvinceSpinner = (Spinner) findViewById(R.id.spinner_province);

        mProvinceList = new ArrayList<String>();
        mProvinceList.add("开户行所在省份");
        mProvinceAdapter = new ArrayAdapter<String>(mContext, R.layout.spinner_item, mProvinceList);

        // 设置样式
        mProvinceAdapter.setDropDownViewResource(R.layout.spinner_drop_item);
        // 加载适配器
        mProvinceSpinner.setAdapter(mProvinceAdapter);

        mProvinceSearchAdapter = new SearchAdapter(mContext, mProvinceList);

        mProvinceEdit.setAdapter(mProvinceSearchAdapter);
        mProvinceEdit.setThreshold(1);

        mCityEdit = (AutoCompleteTextView) contentView.findViewById(R.id.edit_city);
        mCitySpinner = (Spinner) findViewById(R.id.spinner_city);
        // 适配器
        mCityList = new ArrayList<String>();
        mCityCodeList = new ArrayList<String>();
        mCityList.add("开户行所在城市");
        mCityCodeList.add("");
        mCityAdapter = new ArrayAdapter<String>(mContext, R.layout.spinner_item, mCityList);
        // 设置样式
        mCityAdapter.setDropDownViewResource(R.layout.spinner_drop_item);
        // 加载适配器
        mCitySpinner.setAdapter(mCityAdapter);

        mCitySearchAdapter = new SearchAdapter(mContext, mCityList);

        mCityEdit.setAdapter(mCitySearchAdapter);
        mCityEdit.setThreshold(1);

        mOpenBankEdit = (AutoCompleteTextView) contentView.findViewById(R.id.edit_openbank);
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

        mBranchBankEdit = (AutoCompleteTextView) contentView.findViewById(R.id.edit_branchbank);
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

    public void getInfo() {
        quickPayService.getBankInfoAsync(SessonData.loginUser, new QuickPayCallbackListener<BankInfo>() {
            @Override
            public void onSuccess(BankInfo data) {
                mProvinceEdit.setText(data.getProvince());
                mCityEdit.setText(data.getCity());
                mOpenBankEdit.setText(data.getBankOpen());
                mBranchBankEdit.setText(data.getBranchBank());

                mNameEdit.setText(data.getPayee());
                mBanknumEdit.setText(data.getPayeeCard());
                mPhonenumEdit.setText(data.getPhoneNum());

                mCityList.clear();
                mCityCodeList.clear();
                mCityList.add("开户行所在城市");
                mCityCodeList.add("");
                if (mProvinceList.indexOf(mProvinceEdit.getText().toString()) > 0) {
                    String province = mProvinceEdit.getText().toString();
                    bankDataService.getCity(province, new QuickPayCallbackListener<List<City>>() {
                        @Override
                        public void onSuccess(List<City> dataCity) {
                            updateCityAdapter(dataCity);

                            bankDataService.getBank(new QuickPayCallbackListener<List<Bank>>() {
                                @Override
                                public void onSuccess(List<Bank> dataBank) {
                                    updateBankAdapter(dataBank);
                                    String _city = mCityEdit.getText().toString();
                                    String _bank = mOpenBankEdit.getText().toString();
                                    int _cityIndex = mCityList.indexOf(_city);
                                    int _bank_index = mOpenBankList.indexOf(_bank);

                                    if (_cityIndex > 0 && _bank_index > 0) {
                                        String _cityCode = mCityCodeList.get(_cityIndex);
                                        String _bankId = mBankIdList.get(_bank_index);
                                        bankDataService.getBranchBank(_cityCode, _bankId, new BranchBankQuickPayCallbackListener());
                                    }
                                }

                                @Override
                                public void onFailure(QuickPayException ex) {

                                }
                            });
                        }

                        @Override
                        public void onFailure(QuickPayException ex) {

                        }
                    });
                }
            }

            @Override
            public void onFailure(QuickPayException ex) {

            }
        });
    }

    public void initData() {
        bankDataService.getProvince(new ProvinceQuickPayCallbackListener());
    }

    private void initListener() {

        mProvinceEdit.addTextChangedListener(new AccountTextWatcher(mProvinceEdit));

        mCityEdit.addTextChangedListener(new AccountTextWatcher(mCityEdit));

        mOpenBankEdit.addTextChangedListener(new AccountTextWatcher(mOpenBankEdit));

        mProvinceSpinner.setOnItemSelectedListener(new AccountOnItemSelectedListener());

        mCitySpinner.setOnItemSelectedListener(new AccountOnItemSelectedListener());

        mOpenBankSpinner.setOnItemSelectedListener(new AccountOnItemSelectedListener());

        mBranchBankSpinner.setOnItemSelectedListener(new AccountOnItemSelectedListener());

        mCityEdit.setOnFocusChangeListener(new AccountOnFocusChangeListener());

        mOpenBankEdit.setOnFocusChangeListener(new AccountOnFocusChangeListener());

        mBranchBankEdit.setOnFocusChangeListener(new AccountOnFocusChangeListener());


    }

    public void finishedOnClick() {
        if (!validate()) {
            return;
        }
        mBaseActivity.startLoading();
        User user = new User();
        user.setUsername(SessonData.loginUser.getUsername());
        user.setPassword(SessonData.loginUser.getPassword());
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
        user.setPayeeCard(mBanknumEdit.getText().toString().replace(" ", ""));
        user.setPhoneNum(mPhonenumEdit.getText().toString());

        quickPayService.updateInfoAsync(user, new QuickPayCallbackListener<User>() {
            @Override
            public void onSuccess(User data) {
                // 更新UI
                mBaseActivity.endLoading();
                View alertView = mBaseActivity.findViewById(R.id.alert_dialog);
                String alertMsg = getResources().getString(R.string.alert_update_success);
                Bitmap bitmap = BitmapFactory.decodeResource(mContext.getResources(), R.drawable.right);
                AlertDialog alertDialog = new AlertDialog(mContext, null, alertView, alertMsg, bitmap);
                alertDialog.show();
            }

            @Override
            public void onFailure(QuickPayException ex) {
                String errorMsg = ex.getErrorMsg();
                mBaseActivity.endLoading();
                Bitmap bitmap = BitmapFactory.decodeResource(mContext.getResources(), R.drawable.wrong);
                mBaseActivity.alertShow(errorMsg, bitmap);

            }
        });

    }

    private boolean validate() {
        String name = mNameEdit.getText().toString().replace(" ", ""); //姓名,这里把空格统统干掉
        String banknum = mBanknumEdit.getText().toString().replace(" ", ""); //银行卡号,这里把空格统统干掉
        String phonenum = mPhonenumEdit.getText().toString().replace(" ", "");//手机号,这里把空格统统干掉

        String province = mProvinceEdit.getText().toString();
        String city = mCityEdit.getText().toString();
        String openbank = mOpenBankEdit.getText().toString();
        String branchbank = mBranchBankEdit.getText().toString();

        if (mBranchBankList.indexOf(branchbank) < 0) {
            mBranchBankEdit.setText("");
        }

        String alertMsg = "";
        Bitmap alertBitmap = BitmapFactory.decodeResource(this.getResources(), R.drawable.wrong);
        if (TextUtils.isEmpty(province)) {
            alertMsg = ShowMoneyApp.getResString(R.string.alert_error_province_cannot_empty);
            mBaseActivity.alertShow(alertMsg, alertBitmap);
            return false;
        }

        if (TextUtils.isEmpty(city)) {
            alertMsg = ShowMoneyApp.getResString(R.string.alert_error_city_cannot_empty);
            mBaseActivity.alertShow(alertMsg, alertBitmap);
            return false;
        }

        if (TextUtils.isEmpty(openbank)) {
            alertMsg = ShowMoneyApp.getResString(R.string.alert_error_bank_cannot_empty);
            mBaseActivity.alertShow(alertMsg, alertBitmap);
            return false;
        }

        if (TextUtils.isEmpty(branchbank)) {
            alertMsg = ShowMoneyApp.getResString(R.string.alert_error_bankbranch_cannot_empty);
            mBaseActivity.alertShow(alertMsg, alertBitmap);
            return false;
        }

        if (TextUtils.isEmpty(name)) {
            alertMsg = ShowMoneyApp.getResString(R.string.alert_error_name_cannot_empty);
            mBaseActivity.alertShow(alertMsg, alertBitmap);
            return false;
        }

        if (TextUtils.isEmpty(banknum)) {
            alertMsg = ShowMoneyApp.getResString(R.string.alert_error_banknum_cannot_empty);
            mBaseActivity.alertShow(alertMsg, alertBitmap);
            return false;
        }

        if (!VerifyUtil.checkBankCard(banknum)) {
            alertMsg = ShowMoneyApp.getResString(R.string.alert_error_banknum_format_error);
            mBaseActivity.alertShow(alertMsg, alertBitmap);
            return false;
        }

        if (TextUtils.isEmpty(phonenum)) {
            alertMsg = ShowMoneyApp.getResString(R.string.alert_error_phonenum_cannot_empty);
            mBaseActivity.alertShow(alertMsg, alertBitmap);
            return false;
        }

        if (!VerifyUtil.isMobileNO(phonenum)) {
            alertMsg = ShowMoneyApp.getResString(R.string.alert_error_phonenum_format_error);
            mBaseActivity.alertShow(alertMsg, alertBitmap);
            return false;
        }

        return true;
    }


    private void updateProvinceAdapter(List<Province> data) {
        //这里直接得到的就是一个省份的list，不需要再去用json去解析了。
        List<String> tempProvinceList = new ArrayList<>();
        tempProvinceList.add(0, "开户行所在省份");

        Iterator<Province> iterator = data.iterator();
        while (iterator.hasNext()) {
            Province p = iterator.next();
            tempProvinceList.add(p.getProvinceName());
        }
        mProvinceList.clear();
        mProvinceList.addAll(tempProvinceList);
        mProvinceAdapter.notifyDataSetChanged();
        mProvinceSearchAdapter.setData(mProvinceList);
        mProvinceSearchAdapter.notifyDataSetChanged();
    }

    private void updateBankAdapter(List<Bank> data) {
        List<String> tempOpenBankList = new ArrayList<String>();
        List<String> tempBankIdList = new ArrayList<String>();

        Iterator<Bank> it = data.iterator();
        while (it.hasNext()) {
            Bank b = it.next();
            tempOpenBankList.add(b.getBankName());
            tempBankIdList.add(b.getId());
        }
        tempOpenBankList.add(0, "请选择开户银行");
        tempBankIdList.add(0, "");//为了使index对应起来

        mOpenBankList.clear();
        mBankIdList.clear();
        mOpenBankList.addAll(tempOpenBankList);
        mBankIdList.addAll(tempBankIdList);
        mOpenBankSpinner.setSelection(0);
        mOpenBankAdapter.notifyDataSetChanged();
        mOpenBankSearchAdapter.notifyDataSetChanged();
    }

    private void updateCityAdapter(final List<City> data) {
        ArrayList<String> tempCityList = new ArrayList<String>();
        ArrayList<String> tempCityCodeList = new ArrayList<String>();
        tempCityList.add(0, "开户行所在城市");
        tempCityCodeList.add(0, "");
        Iterator<City> it = data.iterator();
        while (it.hasNext()) {
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


    private void updateBranchBankAdapter(List<SubBank> data) {
        final List<String> tempBranchBankList = new ArrayList<String>();
        final List<String> tempBankNoList = new ArrayList<String>();
        Iterator<SubBank> it = data.iterator();

        while (it.hasNext()) {
            SubBank sb = it.next();
            tempBranchBankList.add(sb.getBankName());
            tempBankNoList.add(sb.getOneBankNo() + "|" + sb.getTwoBankNo());
        }
        tempBranchBankList.add(0, "请选择开户支行");
        tempBankNoList.add(0, "行号");

        mBranchBankList.clear();
        mBranchBankList.addAll(tempBranchBankList);
        mBankNoList.clear();
        mBankNoList.addAll(tempBankNoList);
        mBranchBankSpinner.setSelection(0);
        mBranchBankAdapter.notifyDataSetChanged();
        mBranchBankSearchAdapter.setData(mBranchBankList);
    }

    //内部类，实现QuickPayCallbackListener接口,用来获取bank信息
    private class BankQuickPayCallbackListener implements QuickPayCallbackListener<List<Bank>> {

        @Override
        public void onSuccess(List<Bank> data) {
            updateBankAdapter(data);
        }

        @Override
        public void onFailure(QuickPayException ex) {

        }
    }

    //内部类，实现QuickPayCallbackListener接口
    private class ProvinceQuickPayCallbackListener implements QuickPayCallbackListener<List<Province>> {

        @Override
        public void onSuccess(List<Province> data) {
            updateProvinceAdapter(data);
        }

        @Override
        public void onFailure(QuickPayException ex) {

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

    private class BranchBankQuickPayCallbackListener implements QuickPayCallbackListener<List<SubBank>> {

        @Override
        public void onSuccess(List<SubBank> data) {
            updateBranchBankAdapter(data);
        }

        @Override
        public void onFailure(QuickPayException ex) {

        }
    }


    private class AccountTextWatcher implements TextWatcher {
        private View view;

        public AccountTextWatcher(View view) {
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
            String province, city, openBank;

            switch (view.getId()) {
                case R.id.edit_province:
                    //province
                    mCityEdit.setText("");//先把city的清空
                    province = mProvinceEdit.getText().toString();
                    if (mProvinceList.indexOf(province) > 0) {
                        bankDataService.getCity(province, new CityQuickPayCallbackListener());
                    }//end if()
                    break;
                case R.id.edit_city:
                    //city
                    mOpenBankEdit.setText("");
                    city = mCityEdit.getText().toString();
                    if (mCityList.indexOf(city) > 0) {
                        bankDataService.getBank(new BankQuickPayCallbackListener());
                    }
                    break;
                case R.id.edit_openbank:
                    //bank
                    mBranchBankEdit.setText("");
                    openBank = mOpenBankEdit.getText().toString();
                    city = mCityEdit.getText().toString();
                    province = mProvinceEdit.getText().toString();
                    if (mOpenBankList.indexOf(openBank) > 0) {
                        int indexCity = mCityList.indexOf(city);
                        int indexBank = mOpenBankList.indexOf(openBank);
                        if (indexBank > 0 && indexCity > 0) {
                            String cityCode = mCityCodeList.get(indexCity);
                            String bankId = mBankIdList.get(indexBank);
                            bankDataService.getBranchBank(cityCode, bankId, new BranchBankQuickPayCallbackListener());
                        }
                    }
                    break;
                case R.id.edit_branchbank:
                    //branch bank
                    break;
                case R.id.info_name:
                    //name
                    break;
                case R.id.info_banknum:
                    //bank number
                    break;
                case R.id.info_phonenum:
                    //phone number
                    break;
            }
        }
    }

    private class AccountOnItemSelectedListener implements OnItemSelectedListener {

        @Override
        public void onItemSelected(AdapterView<?> parent, View view, int position, long id) {
            if (position > 0) {
                switch (parent.getId()) {
                    case R.id.spinner_province:
                        mProvinceEdit.setText(mProvinceList.get(position));
                        break;
                    case R.id.spinner_city:
                        mCityEdit.setText(mCityList.get(position));
                        break;
                    case R.id.spinner_openbank:
                        mOpenBankEdit.setText(mOpenBankList.get(position));
                        break;
                    case R.id.spinner_branchbank:
                        mBranchBankEdit.setText(mBranchBankList.get(position));
                        break;
                }
            }
        }

        @Override
        public void onNothingSelected(AdapterView<?> parent) {

        }
    }

    private class AccountOnFocusChangeListener implements OnFocusChangeListener {

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
