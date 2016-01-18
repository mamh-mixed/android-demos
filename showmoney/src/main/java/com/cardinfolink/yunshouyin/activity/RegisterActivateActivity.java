package com.cardinfolink.yunshouyin.activity;

import android.content.Intent;
import android.os.Bundle;
import android.view.View;
import android.widget.Button;
import android.widget.TextView;

import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.api.QuickPayException;
import com.cardinfolink.yunshouyin.core.QuickPayCallbackListener;
import com.cardinfolink.yunshouyin.data.SessonData;
import com.cardinfolink.yunshouyin.data.User;
import com.cardinfolink.yunshouyin.ui.SettingActionBarItem;
import com.cardinfolink.yunshouyin.view.YellowTips;


public class RegisterActivateActivity extends BaseActivity {
    private SettingActionBarItem mActionBar;
    private Button mNext;
    private TextView mActivateMessage;//显示账号是否激活的一些提示信息的
    private TextView mAccount;//显示账号
    private YellowTips mYellowTips;

    private boolean isActivate = false;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_register_activate);

        mActionBar = (SettingActionBarItem) findViewById(R.id.action_bar);
        mActionBar.setLeftTextOnclickListner(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                finish();
            }
        });
        mYellowTips = new YellowTips(this, findViewById(R.id.yellow_tips));

        mAccount = (TextView) findViewById(R.id.account);
        //这里得到上一个activity中（register acitvity）设置的账号
        final String username = SessonData.loginUser.getUsername();
        final String password = SessonData.loginUser.getPassword();
        mAccount.setText(username);

        mNext = (Button) findViewById(R.id.btnnext);
        mNext.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                if (isActivate) {
                    Intent intent = new Intent(RegisterActivateActivity.this, RegisterNextActivity.class);
                    startActivity(intent);
                    finish();
                } else {
                    checkActivate(username, password);
                }
            }
        });

        mActivateMessage = (TextView) findViewById(R.id.activate_message);


        checkActivate(username, password);
    }

    /**
     * 检查用户是否是激活的
     *
     * @param username
     * @param password
     */
    private void checkActivate(String username, String password) {
        if (isActivate) {
            //这里先判断一下是不是激活的，是话就不执行了
            return;
        }

        quickPayService.loginAsync(username, password, new QuickPayCallbackListener<User>() {
            @Override
            public void onSuccess(User data) {
                isActivate = true;

                // clientid为空,跳转到完善信息页面
                mNext.setText(getResources().getString(R.string.activate_i_had_activated));
                mActivateMessage.setText(getResources().getString(R.string.activate_your_had_activate));

            }

            @Override
            public void onFailure(QuickPayException ex) {
                String errorCode = ex.getErrorCode();
                String errorMsg = ex.getErrorMsg();
                mYellowTips.show(errorMsg);
                //其他出错的状态怎么设置这个？？？
                isActivate = false;
            }
        });
    }

}
