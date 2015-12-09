package com.cardinfolink.yunshouyin.activity;

import android.os.Bundle;
import android.view.View;

import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.api.QuickPayException;
import com.cardinfolink.yunshouyin.core.QuickPayCallbackListener;
import com.cardinfolink.yunshouyin.data.SessonData;
import com.cardinfolink.yunshouyin.model.BankInfo;
import com.cardinfolink.yunshouyin.ui.SettingActionBarItem;
import com.cardinfolink.yunshouyin.ui.SettingDetailItem;

public class AccountInfoActivity extends BaseActivity {

    private SettingActionBarItem mAccountInfo;//账户信息 界面的标题栏

    private SettingDetailItem mName;//这里显示收款人
    private SettingDetailItem mLocalCity;//所在城市
    private SettingDetailItem mBankName;//清算银行名字
    private SettingDetailItem mMerchantName;//商户名
    private SettingDetailItem mPhoneNumber;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_account_info);

        mAccountInfo = (SettingActionBarItem) findViewById(R.id.sabi_account_info);
        mAccountInfo.setLeftTextOnclickListner(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                finish();
            }
        });

        mName = (SettingDetailItem) findViewById(R.id.name);//这里显示收款人 payee的值
        mLocalCity = (SettingDetailItem) findViewById(R.id.local_city);
        mBankName = (SettingDetailItem) findViewById(R.id.bank_name);
        mMerchantName = (SettingDetailItem) findViewById(R.id.merchant_name);
        mPhoneNumber = (SettingDetailItem) findViewById(R.id.phone_number);

        initUserInfo();//初始化用户的信息

    }

    /**
     * 初始化用户的信息，例如显示用户手机号，银行，所在城市等。。。。
     * <p/>
     * {
     * "state": "success",
     * "count": 0,
     * "size": 0,
     * "refdcount": 0,
     * "info": {
     * "bank_open": "中国工商银行",
     * "payee": "马明辉",//收款人
     * "payee_card": "6222021001114863340",
     * "phone_num": "13014625286",
     * "province": "上海市",
     * "city": "上海市",
     * "branch_bank": "中国工商银行股份有限公司上海市漕宝路支行",
     * "bankNo": "102290004911|102100099996"
     * }
     * }
     */

    private void initUserInfo() {
        quickPayService.getBankInfoAsync(SessonData.loginUser, new QuickPayCallbackListener<BankInfo>() {
            @Override
            public void onSuccess(BankInfo bankInfo) {
                mName.setDetail(bankInfo.getPayee());
                mLocalCity.setDetail(bankInfo.getCity());
                mBankName.setDetail(bankInfo.getBankOpen());
                mPhoneNumber.setDetail(bankInfo.getPhoneNum());
                //这里商户名没有获得到！！！！！！！！
            }

            @Override
            public void onFailure(QuickPayException ex) {

            }
        });
    }

}
