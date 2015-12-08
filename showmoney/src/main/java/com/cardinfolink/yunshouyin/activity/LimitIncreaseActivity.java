package com.cardinfolink.yunshouyin.activity;

import android.content.Context;
import android.graphics.Bitmap;
import android.graphics.BitmapFactory;
import android.os.Bundle;
import android.text.TextUtils;
import android.view.View;
import android.widget.Button;
import android.widget.EditText;

import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.api.QuickPayException;
import com.cardinfolink.yunshouyin.core.QuickPayCallbackListener;
import com.cardinfolink.yunshouyin.data.SessonData;
import com.cardinfolink.yunshouyin.util.ShowMoneyApp;
import com.cardinfolink.yunshouyin.util.VerifyUtil;


/**
 * 提升限额的界面
 */
public class LimitIncreaseActivity extends BaseActivity {
    private EditText mNameEdit;
    private EditText mEmailEdit;
    private EditText mPhonenumEdit;
    private Button mCommit;

    private Context mContext;


    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_limit_increase);
        mContext = this;

        mNameEdit = (EditText) findViewById(R.id.limitincrease_name);
        mEmailEdit = (EditText) findViewById(R.id.limitincrease_email);
        mPhonenumEdit = (EditText) findViewById(R.id.limitincrease_phone);
        mCommit = (Button) findViewById(R.id.limitincrease_btnsubmit);
        mCommit.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                increaseLimitOnClick(v);
            }
        });
    }

    public void increaseLimitOnClick(View v) {
        if (!validate()) {
            return;
        }

        startLoading();
        String name = mNameEdit.getText().toString().replace(" ", "");
        String email = mEmailEdit.getText().toString().replace(" ", "");
        String phone = mPhonenumEdit.getText().toString().replace(" ", "");
        SessonData.loginUser.setLimitEmail(email);
        SessonData.loginUser.setLimitPhone(phone);
        SessonData.loginUser.setLimitName(name);

        quickPayService.increaseLimitAsync(SessonData.loginUser, new QuickPayCallbackListener<Void>() {
            @Override
            public void onSuccess(Void data) {
                endLoading();
                String alertMsg = ShowMoneyApp.getResString(R.string.alert_commit_success);
                Bitmap bitmap = BitmapFactory.decodeResource(mContext.getResources(), R.drawable.right);
                alertShow(alertMsg, bitmap);
            }

            @Override
            public void onFailure(QuickPayException ex) {
                endLoading();
                String alertMsg = ex.getErrorMsg();
                Bitmap bitmap = BitmapFactory.decodeResource(mContext.getResources(), R.drawable.wrong);
                alertShow(alertMsg, bitmap);
            }
        });
    }

    private boolean validate() {
        String name = mNameEdit.getText().toString().replace(" ", "");
        String email = mEmailEdit.getText().toString().replace(" ", "");
        String phone = mPhonenumEdit.getText().toString().replace(" ", "");

        String alertMsg = "";
        Bitmap bitmap = BitmapFactory.decodeResource(this.getResources(), R.drawable.wrong);
        if (TextUtils.isEmpty(phone)) {
            //这里调用了父类里的方法了。
            alertMsg = ShowMoneyApp.getResString(R.string.alert_error_phonenum_cannot_empty);
            alertShow(alertMsg, bitmap);
            return false;
        }

        if (!VerifyUtil.isMobileNO(phone)) {
            alertMsg = ShowMoneyApp.getResString(R.string.alert_error_phonenum_format_error);
            alertShow(alertMsg, bitmap);
            return false;
        }

        if (TextUtils.isEmpty(name)) {
            alertMsg = ShowMoneyApp.getResString(R.string.alert_error_name_cannot_empty);
            alertShow(alertMsg, bitmap);
            return false;
        }

        if (TextUtils.isEmpty(email)) {
            alertMsg = ShowMoneyApp.getResString(R.string.alert_error_email_cannot_empty);
            alertShow(alertMsg, bitmap);
            return false;
        }

        if (!VerifyUtil.checkEmail(email)) {
            alertMsg = ShowMoneyApp.getResString(R.string.alert_error_email_format_error);
            alertShow(alertMsg, bitmap);
            return false;
        }

        return true;
    }


}
