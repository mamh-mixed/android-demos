package com.cardinfolink.yunshouyin.activity;

import android.app.Activity;
import android.content.Intent;
import android.net.Uri;
import android.os.Bundle;
import android.os.SystemClock;
import android.text.TextUtils;
import android.view.View;
import android.widget.TextView;
import android.widget.Toast;

import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.ui.SettingActionBarItem;
import com.cardinfolink.yunshouyin.ui.SettingClikcItem;
import com.cardinfolink.yunshouyin.ui.SettingDetailItem;
import com.umeng.update.UmengUpdateAgent;
import com.umeng.update.UmengUpdateListener;
import com.umeng.update.UpdateResponse;
import com.umeng.update.UpdateStatus;

public class AboutActivity extends Activity implements View.OnClickListener {
    private static final String TAG = "AboutActivity";

    private SettingActionBarItem mActionBar;//关于云收银的 action bar

    private SettingDetailItem mVersion;//显示版本信息的 item

    private SettingClikcItem mWebsite;//产品网站
    private SettingClikcItem mWelcome;//显示欢迎页面
    private SettingClikcItem mUpdate;//检测更新
    private TextView mAgreement;

    //存储时间的数组
    long[] mHits = new long[6];

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_about);

        mActionBar = (SettingActionBarItem) findViewById(R.id.action_bar);
        mActionBar.setLeftTextOnclickListner(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                finish();
            }
        });

        mVersion = (SettingDetailItem) findViewById(R.id.version);

        mWebsite = (SettingClikcItem) findViewById(R.id.website);
        mWelcome = (SettingClikcItem) findViewById(R.id.welcome);
        mUpdate = (SettingClikcItem) findViewById(R.id.update);
        mAgreement = (TextView) findViewById(R.id.tv_agreement);

        mWebsite.setOnClickListener(this);
        mWelcome.setOnClickListener(this);
        mUpdate.setOnClickListener(this);
        mAgreement.setOnClickListener(this);
        mVersion.setOnClickListener(this);

        setVersionName();//获取versionName的值并设置到mVersion里面
    }

    private void setVersionName() {
        try {
            String pkName = getPackageName();
            String versionName = getPackageManager().getPackageInfo(pkName, 0).versionName;
            mVersion.setDetail(versionName);
        } catch (Exception e) {
            mVersion.setDetail("");
        }
    }

    @Override
    public void onClick(View v) {
        Intent intent = null;
        switch (v.getId()) {
            case R.id.website:
                String urlStr = mWebsite.getRightText();
                if (!TextUtils.isEmpty(urlStr)) {
                    //为什么要这样做？？？？
                    try {
                        Uri uri = Uri.parse(urlStr);
                        intent = new Intent(Intent.ACTION_VIEW, uri);
                        startActivity(intent);
                    } catch (Exception e) {

                    }
                }
                break;
            case R.id.welcome:
                intent = new Intent(AboutActivity.this, GuideActivity.class);
                Bundle bundle = new Bundle();
                bundle.putString("original", "AboutActivity");
                intent.putExtras(bundle);
                startActivity(intent);
                break;
            case R.id.update:
                //检查更新
                checkUpdate();
                break;
            case R.id.tv_agreement:
                intent = new Intent(AboutActivity.this, AgreementActivity.class);
                startActivity(intent);
                break;
            case R.id.version:
                versionClick();
                break;
        }
    }


    private void versionClick() {
        //实现数组的移位操作，点击一次，左移一位，末尾补上当前开机时间（cpu的时间）
        System.arraycopy(mHits, 1, mHits, 0, mHits.length - 1);
        mHits[mHits.length - 1] = SystemClock.uptimeMillis();
        if (mHits[0] >= (SystemClock.uptimeMillis() - 500)) {
            Intent intent = new Intent(AboutActivity.this, ConfigActivity.class);
            startActivity(intent);
        }
    }


    private void checkUpdate() {
        UmengUpdateAgent.setUpdateAutoPopup(false);
        UmengUpdateAgent.setUpdateListener(new UmengUpdateListener() {
            @Override
            public void onUpdateReturned(int updateStatus, UpdateResponse updateInfo) {
                String toastMsg;
                switch (updateStatus) {
                    case UpdateStatus.Yes: // has update
                        UmengUpdateAgent.showUpdateDialog(AboutActivity.this, updateInfo);
                        break;
                    case UpdateStatus.No: // has no update
                        toastMsg = getResources().getString(R.string.setting_no_update);
                        Toast.makeText(AboutActivity.this, toastMsg, Toast.LENGTH_SHORT).show();
                        break;
                    case UpdateStatus.NoneWifi: // none wifi
                        toastMsg = getResources().getString(R.string.setting_no_wifi_no_update);
                        Toast.makeText(AboutActivity.this, toastMsg, Toast.LENGTH_SHORT).show();
                        break;
                    case UpdateStatus.Timeout: // time out
                        toastMsg = getResources().getString(R.string.setting_update_timeout);
                        Toast.makeText(AboutActivity.this, toastMsg, Toast.LENGTH_SHORT).show();
                        break;
                }
            }
        });
        UmengUpdateAgent.update(this);
    }
}
