package com.cardinfolink.yunshouyin.activity;

import android.content.Intent;
import android.os.Bundle;
import android.view.View;
import android.widget.Button;
import android.widget.ImageView;

import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.ui.SettingActionBarItem;
import com.cardinfolink.yunshouyin.view.HintDialog;

/**
 * 提升限额 页面，也是 免费升级页面
 */
public class StartIncreaseActivity extends BaseActivity implements View.OnClickListener {
    private static final String TYPE = "type";
    private static final int PERSON = 0;
    private static final int COMPANY = 1;

    private SettingActionBarItem mActionBar;// action bar
    private ImageView mHelp;
    private HintDialog mHintDialog;
    private Button mPerson;
    private Button mCompany;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_start_increase);

        mActionBar = (SettingActionBarItem) findViewById(R.id.action_bar);
        mActionBar.setLeftTextOnclickListner(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                finish();
            }
        });

        mPerson = (Button) findViewById(R.id.btnperson);
        mCompany = (Button) findViewById(R.id.btncompany);
        mHelp = (ImageView) findViewById(R.id.iv_help);

        mPerson.setOnClickListener(this);
        mCompany.setOnClickListener(this);
        mHelp.setOnClickListener(this);


        mHintDialog = new HintDialog(this, findViewById(R.id.hint_dialog));

    }

    @Override
    public void onClick(View v) {
        Intent intent;

        switch (v.getId()) {
            case R.id.iv_help:
                //显示一个提示信息的对话框
                mHintDialog.show();
                break;
            case R.id.btnperson:
                //个体用户
                intent = new Intent(StartIncreaseActivity.this, LimitIncreaseActivity.class);
                intent.putExtra(TYPE, PERSON);
                startActivity(intent);
                finish();
                break;
            case R.id.btncompany:
                intent = new Intent(StartIncreaseActivity.this, LimitIncreaseActivity.class);
                intent.putExtra(TYPE, COMPANY);
                startActivity(intent);
                finish();
                //企业用户
                break;
        }


    }
}
