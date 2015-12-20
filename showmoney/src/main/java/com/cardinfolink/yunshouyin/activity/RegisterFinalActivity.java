package com.cardinfolink.yunshouyin.activity;

import android.content.Intent;
import android.graphics.BitmapFactory;
import android.os.Bundle;
import android.text.TextUtils;
import android.util.Log;
import android.view.View;
import android.widget.Button;

import com.cardinfolink.cashiersdk.model.InitData;
import com.cardinfolink.cashiersdk.sdk.CashierSdk;
import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.api.QuickPayException;
import com.cardinfolink.yunshouyin.constant.SystemConfig;
import com.cardinfolink.yunshouyin.core.QuickPayCallbackListener;
import com.cardinfolink.yunshouyin.data.SaveData;
import com.cardinfolink.yunshouyin.data.SessonData;
import com.cardinfolink.yunshouyin.data.User;
import com.cardinfolink.yunshouyin.ui.SettingActionBarItem;
import com.cardinfolink.yunshouyin.util.TelephonyManagerUtil;
import com.cardinfolink.yunshouyin.view.ActivateDialog;

public class RegisterFinalActivity extends BaseActivity implements View.OnClickListener {
    private SettingActionBarItem mActionBar;
    private Button mUseNow;//立即使用
    private Button mIncreaseLimit;//提升限额

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_register_final);
        mActionBar = (SettingActionBarItem) findViewById(R.id.action_bar);
        mActionBar.setLeftTextOnclickListner(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                finish();
            }
        });

        mUseNow = (Button) findViewById(R.id.btnnow);
        mIncreaseLimit = (Button) findViewById(R.id.btnlimit);

        mUseNow.setOnClickListener(this);
        mIncreaseLimit.setOnClickListener(this);
    }

    @Override
    public void onClick(View v) {
        Intent intent;
        switch (v.getId()) {
            case R.id.btnnow:
                login();//注意 这里和loginActivity里面的login（）方法不太一样的

                //立即使用
                break;
            case R.id.btnlimit:
                SessonData.loginUser.setAutoLogin(true);
                SaveData.setUser(RegisterFinalActivity.this, SessonData.loginUser);
                //提升限额,进入到 提升限额的界面，提升用户 选择商户类型
                intent = new Intent(RegisterFinalActivity.this, StartIncreaseActivity.class);
                startActivity(intent);
                finish();
                break;
        }


    }


    private void login() {
        final String username = SessonData.loginUser.getUsername();
        final String password = SessonData.loginUser.getPassword();

        User user = new User();
        user.setUsername(username);
        user.setPassword(password);
        SaveData.setUser(mContext, user);//保存密码到文件

        quickPayService.loginAsync(username, password, new QuickPayCallbackListener<User>() {
            @Override
            public void onSuccess(User data) {
                SessonData.loginUser.setClientid(data.getClientid());
                SessonData.loginUser.setObjectId(data.getObjectId());
                SessonData.loginUser.setLimit(data.getLimit());

                InitData initData = new InitData();
                initData.setMchntid(data.getClientid());
                initData.setInscd(data.getInscd());
                initData.setSignKey(data.getSignKey());
                initData.setTerminalid(TelephonyManagerUtil.getDeviceId(mContext));
                initData.setIsProduce(SystemConfig.IS_PRODUCE);
                CashierSdk.init(initData);//初始化sdk

                //更新UI,这里直接跳转到mainActivity界面了
                Intent intent = new Intent(mContext, MainActivity.class);
                intent.setFlags(Intent.FLAG_ACTIVITY_CLEAR_TOP);
                mContext.startActivity(intent);
                finish();
            }

            @Override
            public void onFailure(QuickPayException ex) {
                String errorCode = ex.getErrorCode();
                String errorMsg = ex.getErrorMsg();
                if (errorCode.equals("user_no_activate")) {
                    //更新UI,这里不太可能是 没激活状态吧
                    View view = findViewById(R.id.activate_dialog);
                    String eMail = SessonData.loginUser.getUsername();
                    ActivateDialog activateDialog = new ActivateDialog(mContext, view, eMail);
                    activateDialog.show();
                } else {
                    mAlertDialog.show(errorMsg, BitmapFactory.decodeResource(mContext.getResources(), R.drawable.wrong));
                }
            }
        });


    }


}
