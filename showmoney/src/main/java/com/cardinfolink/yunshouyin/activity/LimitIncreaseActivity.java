package com.cardinfolink.yunshouyin.activity;

import android.os.Bundle;
import android.view.View;
import android.widget.TextView;

import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.ui.SettingClikcItem;


/**
 * 提升限额的界面,用于上传图片，填写店铺名称
 */
public class LimitIncreaseActivity extends BaseActivity {
    private static final String TYPE = "type";
    private static final int PERSON = 0;
    private static final int COMPANY = 1;
    private SettingClikcItem mTax;
    private SettingClikcItem mOrganization;
    private TextView mMessage;


    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_limit_increase);

        mTax = (SettingClikcItem) findViewById(R.id.tax);
        mOrganization = (SettingClikcItem) findViewById(R.id.organization);
        mMessage = (TextView) findViewById(R.id.increase_message);

        int type = getIntent().getIntExtra(TYPE, PERSON);
        if (PERSON == type) {
            mMessage.setText(getString(R.string.limit_increase_message));
            mTax.setVisibility(View.GONE);
            mOrganization.setVisibility(View.GONE);
        } else if (COMPANY == type) {
            mMessage.setText(getString(R.string.limit_increase_message1));
            mTax.setVisibility(View.VISIBLE);
            mOrganization.setVisibility(View.VISIBLE);
        }
    }


}
