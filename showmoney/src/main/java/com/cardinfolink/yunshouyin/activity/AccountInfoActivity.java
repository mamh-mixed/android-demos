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

    private SettingActionBarItem mActionBar;//账户信息 界面的标题栏

    private SettingDetailItem mName;//这里显示收款人
    private SettingDetailItem mLocalCity;//所在城市
    private SettingDetailItem mBankName;//清算银行名字
    private SettingDetailItem mMerchantName;//商户名
    private SettingDetailItem mPhoneNumber;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_account_info);

        mActionBar = (SettingActionBarItem) findViewById(R.id.action_bar);
        mActionBar.setLeftTextOnclickListner(new View.OnClickListener() {
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

        mMerchantName.setDetail(SessonData.loginUser.getMerName());

        initUserInfo();//初始化用户的信息

    }


    private void initUserInfo() {
        quickPayService.getBankInfoAsync(SessonData.loginUser, new QuickPayCallbackListener<BankInfo>() {
            @Override
            public void onSuccess(BankInfo bankInfo) {
                mName.setDetail(bankInfo.getPayee());
                mLocalCity.setDetail(bankInfo.getCity());
                mBankName.setDetail(bankInfo.getBankOpen());
                mPhoneNumber.setDetail(bankInfo.getPhoneNum());
                //这里商户名没有获得到！！！！！！！！
                mMerchantName.setDetail(SessonData.loginUser.getMerName());
            }

            @Override
            public void onFailure(QuickPayException ex) {

            }
        });
    }

}
