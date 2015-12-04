package com.cardinfolink.yunshouyin.activity;

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
import com.cardinfolink.yunshouyin.data.SaveData;
import com.cardinfolink.yunshouyin.data.SessonData;
import com.cardinfolink.yunshouyin.data.User;
import com.cardinfolink.yunshouyin.util.ShowMoneyApp;
import com.cardinfolink.yunshouyin.util.VerifyUtil;
import com.cardinfolink.yunshouyin.view.ActivateDialog;

import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class RegisterActivity extends BaseActivity {
    private EditText mEmailEdit;
    private EditText mPasswordEdit;
    private EditText mQrPasswordEdit;

    private Button mRegisterNext;

    /**
     * 验证邮箱
     *
     * @param email
     * @return
     */
    public static boolean checkEmail(String email) {
        boolean flag = false;
        try {
            String check = "^([a-z0-9A-Z]+[-|_|\\.]?)+[a-z0-9A-Z]@([a-z0-9A-Z]+(-[a-z0-9A-Z]+)?\\.)+[a-zA-Z]{2,}$";
            Pattern regex = Pattern.compile(check);
            Matcher matcher = regex.matcher(email);
            flag = matcher.matches();
        } catch (Exception e) {
            flag = false;
        }
        return flag;
    }

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.register_activity);
        mContext = this;
        initLayout();
    }

    private void initLayout() {
        mEmailEdit = (EditText) findViewById(R.id.register_email);
        VerifyUtil.addEmailLimit(mEmailEdit);

        mPasswordEdit = (EditText) findViewById(R.id.register_password);
        VerifyUtil.addEmailLimit(mPasswordEdit);

        mQrPasswordEdit = (EditText) findViewById(R.id.register_qr_password);
        VerifyUtil.addEmailLimit(mQrPasswordEdit);

        mRegisterNext = (Button) findViewById(R.id.btnregister);
        mRegisterNext.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                btnRegisterNextOnClick(v);
            }
        });
    }

    public void btnRegisterNextOnClick(View view) {
        if (!validate()) {
            return;
        }
        mLoadingDialog.startLoading();
        final String username = mEmailEdit.getText().toString(); //用户名
        final String password = mPasswordEdit.getText().toString(); //密码，第一次输入的
        final String qrPassword = mQrPasswordEdit.getText().toString(); //确认密码，第二次输入的

        quickPayService.registerAsync(username, password, qrPassword, new QuickPayCallbackListener<Void>() {
            @Override
            public void onSuccess(Void data) {
                //没有返回值的,那边返回的是null.
                User user = new User();
                user.setUsername(username);
                user.setPassword(password);
                SaveData.setUser(mContext, user);

                SessonData.loginUser.setUsername(username);
                SessonData.loginUser.setPassword(password);

                mLoadingDialog.endLoading();
                View activateView = findViewById(R.id.activate_dialog);
                ActivateDialog activateDialog = new ActivateDialog(mContext, activateView, username);
                activateDialog.show();
            }

            @Override
            public void onFailure(QuickPayException ex) {
                String errorMsg = ex.getErrorMsg();
                //更新UI
                mLoadingDialog.endLoading();
                Bitmap alertBitmap = BitmapFactory.decodeResource(mContext.getResources(), R.drawable.wrong);
                mAlertDialog.show(errorMsg, alertBitmap);
            }
        });
    }

    private boolean validate() {
        String email = mEmailEdit.getText().toString();
        String password = mPasswordEdit.getText().toString();
        String qrPassword = mQrPasswordEdit.getText().toString();

        String alertMsg = "";
        Bitmap alertBitmap = BitmapFactory.decodeResource(this.getResources(), R.drawable.wrong);
        if (TextUtils.isEmpty(email)) {
            alertMsg = ShowMoneyApp.getResString(R.string.alert_error_email_cannot_empty);
            mAlertDialog.show(alertMsg, alertBitmap);
            return false;
        }
        if (!checkEmail(email)) {
            alertMsg = ShowMoneyApp.getResString(R.string.alert_error_email_format_error);
            mAlertDialog.show(alertMsg, alertBitmap);
            return false;
        }
        if (TextUtils.isEmpty(password)) {
            alertMsg = ShowMoneyApp.getResString(R.string.alert_error_password_cannot_empty);
            mAlertDialog.show(alertMsg, alertBitmap);
            return false;
        }
        if (password.length() < 6) {
            alertMsg = ShowMoneyApp.getResString(R.string.alert_error_password_short_six);
            mAlertDialog.show(alertMsg, alertBitmap);
            return false;
        }
        if (!password.equals(qrPassword)) {
            alertMsg = ShowMoneyApp.getResString(R.string.alert_error_qrpassword_error);
            mAlertDialog.show(alertMsg, alertBitmap);
            return false;
        }

        return true;
    }

}
