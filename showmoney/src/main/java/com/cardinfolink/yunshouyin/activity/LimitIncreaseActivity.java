package com.cardinfolink.yunshouyin.activity;

import android.os.Bundle;
import android.view.View;
import android.widget.Button;
import android.widget.TextView;

import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.ui.SettingActionBarItem;
import com.cardinfolink.yunshouyin.ui.SettingClikcItem;
import com.cardinfolink.yunshouyin.ui.SettingInputItem;


/**
 * 提升限额的界面,用于上传图片，填写店铺名称
 */
public class LimitIncreaseActivity extends BaseActivity implements View.OnClickListener {
    private static final String TYPE = "type";
    private static final int PERSON = 0;
    private static final int COMPANY = 1;
    private int mType;

    private SettingClikcItem mTax;//税务
    private SettingClikcItem mOrganization;//组织结构照片

    private SettingInputItem mMerchant;//商铺名称
    private SettingInputItem mMerchantAddress;//商铺地址
    private SettingClikcItem mCardPositive;//身份证 正面
    private SettingClikcItem mCardNegative;//身份证 反面
    private SettingClikcItem mBusiness;//营业执照

    private Button mFinish;//完成按钮

    private TextView mMessage;

    private SettingActionBarItem mActionBar;

    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_limit_increase);

        mActionBar= (SettingActionBarItem) findViewById(R.id.action_bar);
        mActionBar.setLeftTextOnclickListner(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                finish();
            }
        });

        //需要输入内容的
        mMerchant = (SettingInputItem) findViewById(R.id.merchant_name);//商铺名称
        mMerchantAddress = (SettingInputItem) findViewById(R.id.merchant_address);//商铺地址

        //上传图片
        mCardPositive = (SettingClikcItem) findViewById(R.id.id_card_positive);//身份证 正面
        mCardNegative = (SettingClikcItem) findViewById(R.id.id_card_negaitive);//身份证 反面
        mBusiness = (SettingClikcItem) findViewById(R.id.business);//营业执照

        //上传图片，只有企业商户才有的
        mTax = (SettingClikcItem) findViewById(R.id.tax);//税务
        mOrganization = (SettingClikcItem) findViewById(R.id.organization);//组织机构

        mMessage = (TextView) findViewById(R.id.increase_message);

        mFinish = (Button) findViewById(R.id.btnfinish);//完成按钮

        mType = getIntent().getIntExtra(TYPE, PERSON);
        if (PERSON == mType) {
            mMessage.setText(getString(R.string.limit_increase_message));
            mTax.setVisibility(View.GONE);
            mOrganization.setVisibility(View.GONE);
        } else if (COMPANY == mType) {
            mMessage.setText(getString(R.string.limit_increase_message1));
            mTax.setVisibility(View.VISIBLE);
            mOrganization.setVisibility(View.VISIBLE);
        }

        mCardPositive.setOnClickListener(this);//身份证 正面
        mCardNegative.setOnClickListener(this);//身份证 反面
        mBusiness.setOnClickListener(this);//营业执照
        mTax.setOnClickListener(this);//税务
        mOrganization.setOnClickListener(this);//组织机构
        mFinish.setOnClickListener(this);//完成按钮
    }


    @Override
    public void onClick(View v) {
        switch (v.getId()) {
            case R.id.id_card_positive:
                //身份证 正面
                break;
            case R.id.id_card_negaitive:
                //身份证 反面
                break;
            case R.id.business:
                //营业执照
                break;
            case R.id.tax:
                //税务
                break;
            case R.id.organization:
                //组织机构
                break;
            case R.id.btnfinish:
                break;
        }

    }
}
