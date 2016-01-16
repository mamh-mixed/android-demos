package com.cardinfolink.yunshouyin.activity;

import android.os.Bundle;
import android.app.Activity;
import android.view.View;
import android.widget.CheckBox;
import android.widget.CompoundButton;
import android.widget.TextView;

import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.constant.SystemConfig;
import com.cardinfolink.yunshouyin.ui.ResultInfoItem;
import com.cardinfolink.yunshouyin.ui.SettingActionBarItem;
import com.cardinfolink.yunshouyin.util.Log;
import com.cardinfolink.yunshouyin.util.ShowMoneyApp;

public class ConfigActivity extends Activity {

    private SettingActionBarItem mActionBar;

    private ResultInfoItem mGitInfo;
    private ResultInfoItem mServer;
    private ResultInfoItem mDebug;
    private TextView mGitCommit;
    private CheckBox mDebugCheckBox;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_config);

        mActionBar = (SettingActionBarItem) findViewById(R.id.action_bar);
        mActionBar.setLeftTextOnclickListner(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                finish();
            }
        });

        mGitInfo = (ResultInfoItem) findViewById(R.id.git);
        mGitCommit = (TextView) findViewById(R.id.git_commit);
        mServer = (ResultInfoItem) findViewById(R.id.server);
        mDebug = (ResultInfoItem) findViewById(R.id.debug);

        String git = ShowMoneyApp.GIT;
        mGitCommit.setText(git);
        mServer.setRightText(SystemConfig.SERVER);

        mDebug.setRightText(String.valueOf(SystemConfig.DEBUG));

        mDebugCheckBox = (CheckBox) findViewById(R.id.cb_debug);
        mDebugCheckBox.setChecked(SystemConfig.DEBUG);
        mDebugCheckBox.setOnCheckedChangeListener(new CompoundButton.OnCheckedChangeListener() {
            @Override
            public void onCheckedChanged(CompoundButton buttonView, boolean isChecked) {
                SystemConfig.DEBUG = isChecked;
                Log.LOG = isChecked;
                mDebug.setRightText(String.valueOf(SystemConfig.DEBUG));
            }
        });


    }

}
