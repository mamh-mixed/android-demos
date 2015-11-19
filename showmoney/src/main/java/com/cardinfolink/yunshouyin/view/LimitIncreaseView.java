package com.cardinfolink.yunshouyin.view;

import android.annotation.SuppressLint;
import android.app.Activity;
import android.content.Context;
import android.graphics.BitmapFactory;
import android.view.LayoutInflater;
import android.view.View;
import android.widget.EditText;
import android.widget.LinearLayout;

import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.activity.BaseActivity;
import com.cardinfolink.yunshouyin.data.SessonData;
import com.cardinfolink.yunshouyin.util.CommunicationListener;
import com.cardinfolink.yunshouyin.util.HttpCommunicationUtil;
import com.cardinfolink.yunshouyin.util.JsonUtil;
import com.cardinfolink.yunshouyin.util.ParamsUtil;
import com.cardinfolink.yunshouyin.util.VerifyUtil;

public class LimitIncreaseView extends LinearLayout {
    private EditText mNameEdit;
    private EditText mEmailEdit;
    private EditText mPhonenumEdit;
    private Context mContext;
    private BaseActivity mBaseActivity;

    public LimitIncreaseView(Context context) {
        super(context);
        mContext = context;
        mBaseActivity = (BaseActivity) mContext;
        View contentView = LayoutInflater.from(context).inflate(
                R.layout.limit_increase, null);
        LinearLayout.LayoutParams layoutParams = new LayoutParams(
                LayoutParams.MATCH_PARENT, LayoutParams.MATCH_PARENT);
        contentView.setLayoutParams(layoutParams);
        addView(contentView);
        mNameEdit = (EditText) contentView.findViewById(R.id.limitincrease_name);
        mEmailEdit = (EditText) contentView.findViewById(R.id.limitincrease_email);
        mPhonenumEdit = (EditText) contentView.findViewById(R.id.limitincrease_phone);
        contentView.findViewById(R.id.limitincrease_btnsubmit).setOnClickListener(new OnClickListener() {

            @Override
            public void onClick(View v) {

                if (check()) {
                    mBaseActivity.startLoading();
                    String name = mNameEdit.getText().toString().replace(" ", "");
                    ;
                    String email = mEmailEdit.getText().toString().replace(" ", "");
                    String phone = mPhonenumEdit.getText().toString().replace(" ", "");
                    SessonData.loginUser.setLimitEmail(email);
                    SessonData.loginUser.setLimitPhone(phone);
                    SessonData.loginUser.setLimitName(name);
                    HttpCommunicationUtil.sendDataToServer(ParamsUtil.getLimitincrease(SessonData.loginUser), new CommunicationListener() {

                        @Override
                        public void onResult(String result) {
                            String state = JsonUtil.getParam(result, "state");
                            if (state.equals("success")) {
                                ((Activity) mContext).runOnUiThread(new Runnable() {

                                    @Override
                                    public void run() {
                                        // 更新UI
                                        mBaseActivity.endLoading();
                                        mBaseActivity.alertShow("提交完成", BitmapFactory.decodeResource(mContext.getResources(), R.drawable.right));
                                    }

                                });

                            } else {
                                ((Activity) mContext).runOnUiThread(new Runnable() {

                                    @Override
                                    public void run() {
                                        // 更新UI
                                        mBaseActivity.endLoading();
                                        mBaseActivity.alertShow(
                                                "提交失败!",
                                                BitmapFactory.decodeResource(
                                                        mContext.getResources(),
                                                        R.drawable.wrong));
                                    }

                                });
                            }
                        }

                        @Override
                        public void onError(final String error) {
                            ((Activity) mContext).runOnUiThread(new Runnable() {

                                @Override
                                public void run() {
                                    // 更新UI
                                    mBaseActivity.endLoading();
                                    mBaseActivity.alertShow(error, BitmapFactory
                                            .decodeResource(
                                                    mContext.getResources(),
                                                    R.drawable.wrong));
                                }

                            });
                        }
                    });
                }
            }
        });
    }

    @SuppressLint("NewApi")
    private boolean check() {
        String name = mNameEdit.getText().toString().replace(" ", "");
        ;
        String email = mEmailEdit.getText().toString().replace(" ", "");
        String phone = mPhonenumEdit.getText().toString().replace(" ", "");
        if (phone.isEmpty()) {
            mBaseActivity.alertShow("手机号不能为空!", BitmapFactory.decodeResource(this.getResources(), R.drawable.wrong));
            return false;

        }

        if (!VerifyUtil.isMobileNO(phone)) {
            mBaseActivity.alertShow("手机号格式不正确!", BitmapFactory.decodeResource(this.getResources(), R.drawable.wrong));
            return false;

        }
        if (name.isEmpty()) {
            mBaseActivity.alertShow("姓名不能为空!", BitmapFactory.decodeResource(this.getResources(), R.drawable.wrong));
            return false;

        }

        if (email.isEmpty()) {
            mBaseActivity.alertShow("邮箱不能为空!", BitmapFactory.decodeResource(this.getResources(), R.drawable.wrong));
            return false;

        }

        if (!VerifyUtil.checkEmail(email)) {
            mBaseActivity.alertShow("邮箱格式不正确!", BitmapFactory.decodeResource(this.getResources(), R.drawable.wrong));
            return false;

        }

        return true;
    }


}
