package com.cardinfolink.yunshouyin.activity;

import android.content.Intent;
import android.net.Uri;
import android.os.Bundle;
import android.app.Activity;
import android.text.TextUtils;
import android.view.View;
import android.widget.Toast;

import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.ui.SettingActionBarItem;
import com.cardinfolink.yunshouyin.ui.SettingClikcItem;
import com.cardinfolink.yunshouyin.ui.SettingDetailItem;

public class AboutActivity extends Activity implements View.OnClickListener {

    private SettingActionBarItem mAbount;//关于云收银的 action bar

    private SettingDetailItem mVersion;//显示版本信息的 item

    private SettingClikcItem mWebsite;//产品网站
    private SettingClikcItem mWelcome;//显示欢迎页面
    private SettingClikcItem mUpdate;//检测更新

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_about);

        mAbount = (SettingActionBarItem) findViewById(R.id.sabi_about);
        mAbount.setLeftTextOnclickListner(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                finish();
            }
        });

        mVersion = (SettingDetailItem) findViewById(R.id.version);

        mWebsite = (SettingClikcItem) findViewById(R.id.website);
        mWelcome = (SettingClikcItem) findViewById(R.id.welcome);
        mUpdate = (SettingClikcItem) findViewById(R.id.update);

        mWebsite.setOnClickListener(this);
        mWelcome.setOnClickListener(this);
        mUpdate.setOnClickListener(this);
    }

    @Override
    public void onClick(View v) {
        Intent intent = null;
        switch (v.getId()) {
            case R.id.website:
                String urlStr = mWebsite.getRightText();
                if (!TextUtils.isEmpty(urlStr)) {
                    Uri uri = Uri.parse(urlStr);
                    intent = new Intent(Intent.ACTION_VIEW, uri);
                    startActivity(intent);
                }
                break;
            case R.id.welcome:
                intent = new Intent(AboutActivity.this, SplashActivity.class);
                startActivity(intent);
                finish();
                break;
            case R.id.update:
                //检查更新
                Toast.makeText(this, "检测更新", Toast.LENGTH_SHORT).show();
                break;
        }
    }
}
