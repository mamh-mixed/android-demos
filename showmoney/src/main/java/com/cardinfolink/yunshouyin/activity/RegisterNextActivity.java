package com.cardinfolink.yunshouyin.activity;

import android.os.Bundle;
import android.util.Log;
import android.view.MotionEvent;
import android.view.View;
import android.widget.Button;

import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.ui.SettingActionBarItem;
import com.cardinfolink.yunshouyin.ui.SettingClikcItem;
import com.cardinfolink.yunshouyin.view.SelectDialog;

public class RegisterNextActivity extends BaseActivity implements View.OnClickListener {
    private static final String TAG = "RegisterNextActivity";

    private SettingActionBarItem mActionBar;

    private SettingClikcItem mSetProvinceCity;//点击去设置省市信息
    private SettingClikcItem mSetBank;//点击设置银行信息，主行，或分行

    private Button mRegisterFinished;

    private SelectDialog selectDialog;


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

        mSetProvinceCity.setOnClickListener(this);
        mSetBank.setOnClickListener(this);
        mRegisterFinished.setOnClickListener(this);

        selectDialog = new SelectDialog(this, findViewById(R.id.select_dialog));
    }


    @Override
    public void onClick(View v) {
        switch (v.getId()) {
            case R.id.btnregister:
                Log.e(TAG, " onclick register");
                break;
            case R.id.province_city:
                Log.e(TAG, " onclick province_city");
                selectDialog.show();
                break;
            case R.id.bank:

                Log.e(TAG, " onclick bank");
                break;
        }
    }


}
