package com.cardinfolink.yunshouyin.view;

import android.content.Context;
import android.content.Intent;
import android.graphics.Bitmap;
import android.graphics.BitmapFactory;
import android.text.TextUtils;
import android.view.LayoutInflater;
import android.view.View;
import android.widget.Button;
import android.widget.EditText;
import android.widget.LinearLayout;

import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.activity.BaseActivity;
import com.cardinfolink.yunshouyin.activity.SplashActivity;
import com.cardinfolink.yunshouyin.api.QuickPayException;
import com.cardinfolink.yunshouyin.core.QuickPayCallbackListener;
import com.cardinfolink.yunshouyin.core.QuickPayService;
import com.cardinfolink.yunshouyin.data.SessonData;
import com.cardinfolink.yunshouyin.util.ShowMoneyApp;
import com.cardinfolink.yunshouyin.util.VerifyUtil;

public class PasswordUpdateView extends LinearLayout {
    private Context mContext;
    private EditText mOldPwdEdit;
    private EditText mNewPwdEdit;
    private EditText mQrPwdEdit;
    private BaseActivity mBaseActivity;
    private Button updatePasswordButton;

    private QuickPayService quickPayService;

    public PasswordUpdateView(Context context) {
        super(context);
        mContext = context;
        mBaseActivity = (BaseActivity) mContext;

        quickPayService = ShowMoneyApp.getInstance().getQuickPayService();

        View contentView = LayoutInflater.from(context).inflate(R.layout.password_update_view, null);
        LinearLayout.LayoutParams layoutParams = new LayoutParams(LayoutParams.MATCH_PARENT, LayoutParams.MATCH_PARENT);
        contentView.setLayoutParams(layoutParams);
        addView(contentView);

        mOldPwdEdit = (EditText) contentView.findViewById(R.id.update_password_oldpwd);
        VerifyUtil.addEmailLimit(mOldPwdEdit);

        mNewPwdEdit = (EditText) contentView.findViewById(R.id.update_password_newpwd);
        VerifyUtil.addEmailLimit(mNewPwdEdit);

        mQrPwdEdit = (EditText) contentView.findViewById(R.id.update_password_qr_newpwd);
        VerifyUtil.addEmailLimit(mQrPwdEdit);


        updatePasswordButton = (Button) contentView.findViewById(R.id.btn_update_password);
        updatePasswordButton.setOnClickListener(new OnClickListener() {

            @Override
            public void onClick(View v) {
                updatePasswordOnClick(v);
            }
        });

    }

    public void updatePasswordOnClick(View v) {
        if (!validate()) {
            return;
        }

        mBaseActivity.startLoading();
        final String oldpwd = mOldPwdEdit.getText().toString().replace(" ", "");
        final String newpwd = mNewPwdEdit.getText().toString().replace(" ", "");

        quickPayService.updatePasswordAsync(SessonData.loginUser.getUsername(), oldpwd, newpwd, new QuickPayCallbackListener<Void>() {
            @Override
            public void onSuccess(Void data) {
                //更新一下UI
                mBaseActivity.endLoading();
                View alertView = mBaseActivity.findViewById(R.id.alert_dialog);
                String alertMsg = getResources().getString(R.string.alert_update_success);
                Bitmap bitmap = BitmapFactory.decodeResource(mContext.getResources(), R.drawable.right);
                AlertDialog alertDialog = new AlertDialog(mContext, null, alertView, alertMsg, bitmap);
                mOldPwdEdit.setText("");
                mNewPwdEdit.setText("");
                mQrPwdEdit.setText("");
                //更新密码之后就退出登录。
                alertDialog.show(new OnClickListener() {
                    @Override
                    public void onClick(View v) {
                        mBaseActivity.startActivity(new Intent(mBaseActivity, SplashActivity.class));
                        mBaseActivity.finish();
                    }
                });
            }

            @Override
            public void onFailure(QuickPayException ex) {
                mBaseActivity.endLoading();
                String error = ex.getErrorMsg();
                View alertView = mBaseActivity.findViewById(R.id.alert_dialog);
                Bitmap bitmap = BitmapFactory.decodeResource(mContext.getResources(), R.drawable.wrong);
                AlertDialog alertDialog = new AlertDialog(mContext, null, alertView, error, bitmap);
                alertDialog.show();
            }
        });
    }

    private boolean validate() {
        String oldpwd = mOldPwdEdit.getText().toString().replace(" ", "");//注意这里把所有的空格都删除了
        String newpwd = mNewPwdEdit.getText().toString().replace(" ", "");
        String qrpwd = mQrPwdEdit.getText().toString().replace(" ", "");

        Bitmap bitmap = BitmapFactory.decodeResource(this.getResources(), R.drawable.wrong);
        String alertMsg = "";
        if (TextUtils.isEmpty(oldpwd)) {
            alertMsg = ShowMoneyApp.getResString(R.string.alert_error_old_password_cannot_empty);
            mBaseActivity.alertShow(alertMsg, bitmap);
            return false;
        }

        if (oldpwd.length() < 6) {
            alertMsg = ShowMoneyApp.getResString(R.string.alert_error_old_password_short_six);
            mBaseActivity.alertShow(alertMsg, bitmap);
            return false;
        }


        if (TextUtils.isEmpty(newpwd)) {
            alertMsg = ShowMoneyApp.getResString(R.string.alert_error_new_password_cannot_empty);
            mBaseActivity.alertShow(alertMsg, bitmap);
            return false;
        }

        if (newpwd.length() < 6) {
            alertMsg = ShowMoneyApp.getResString(R.string.alert_error_new_password_short_six);
            mBaseActivity.alertShow(alertMsg, bitmap);
            return false;
        }

        if (!qrpwd.equals(newpwd)) {
            alertMsg = ShowMoneyApp.getResString(R.string.alert_error_qrpassword_error);
            mBaseActivity.alertShow(alertMsg, bitmap);
            return false;
        }

        return true;
    }

}
