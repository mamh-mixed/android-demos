package com.cardinfolink.yunshouyin.activity;

import android.os.Bundle;
import android.view.View;
import android.widget.ImageView;

import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.ui.SettingActionBarItem;
import com.cardinfolink.yunshouyin.view.HintDialog;

/**
 * 提升限额 页面，也是 免费升级页面
 */
public class StartIncreaseActivity extends BaseActivity implements View.OnClickListener {
    private SettingActionBarItem mActionBar;// action bar
    private ImageView mHelp;
    private HintDialog mHintDialog;


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

        mHelp = (ImageView) findViewById(R.id.iv_help);
        mHelp.setOnClickListener(this);

        mHintDialog = new HintDialog(this, findViewById(R.id.hint_dialog));

    }

    @Override
    public void onClick(View v) {
        switch (v.getId()) {
            case R.id.iv_help:
                //显示一个提示信息的对话框
                mHintDialog.show();
                mHintDialog.setTitle("个体工商户 所需材料\n法人身份证\n营业执照\n\n企业商户 所需材料\n法人身份证\n企业营业执照\n税务登记证\n组织机构代码证");
                mHintDialog.setOkText("我知道了");
                break;
        }


    }
}
