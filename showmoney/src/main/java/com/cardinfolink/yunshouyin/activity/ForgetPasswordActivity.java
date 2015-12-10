package com.cardinfolink.yunshouyin.activity;

import android.app.Activity;
import android.os.Bundle;
import android.util.Log;
import android.view.View;
import android.widget.Button;

import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.api.QuickPayException;
import com.cardinfolink.yunshouyin.core.QuickPayCallbackListener;
import com.cardinfolink.yunshouyin.model.ServerPacket;
import com.cardinfolink.yunshouyin.ui.SettingActionBarItem;
import com.cardinfolink.yunshouyin.ui.SettingInputItem;
import com.cardinfolink.yunshouyin.view.HintDialog;

/**
 * 这个是忘记密码的界面
 */
public class ForgetPasswordActivity extends BaseActivity {
    private static final String TAG = "ForgetPasswordActivity";
    private SettingActionBarItem mActionBar;//标题栏，自定义的标题栏

    private SettingInputItem mEmail;

    private HintDialog mHintDialog;

    private Button mCnfirm;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_forget_password);

        mHintDialog = new HintDialog(this, findViewById(R.id.hint_dialog));

        mActionBar = (SettingActionBarItem) findViewById(R.id.action_bar);
        mActionBar.setLeftTextOnclickListner(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                finish();
            }
        });

        mEmail = (SettingInputItem) findViewById(R.id.email);
        mEmail.setImageViewOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                String title = getResources().getString(R.string.forget_account_message);
                String ok = getResources().getString(R.string.forget_i_known);
                String cancel = getResources().getString(R.string.forget_cancel);
                mHintDialog.show(title, ok, cancel);
            }
        });

        mCnfirm = (Button) findViewById(R.id.confirm);

        mCnfirm.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                forgetPassowd(); //忘记密码了
            }
        });
    }

    /**
     * 点击 按钮 调用这里， 提交忘记密码的请求
     */
    private void forgetPassowd() {
        Log.e(TAG, "     * 点击 按钮 调用这里， 提交忘记密码的请求 ");
        quickPayService.forgetPasswordAsync(mEmail.getText(), new QuickPayCallbackListener<ServerPacket>() {
            @Override
            public void onSuccess(ServerPacket data) {

            }

            @Override
            public void onFailure(QuickPayException ex) {

            }
        });
    }

}
