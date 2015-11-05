package com.cardinfolink.yunshouyin.salesman.activity;

import android.annotation.SuppressLint;
import android.content.Intent;
import android.content.pm.ApplicationInfo;
import android.content.pm.PackageManager;
import android.os.Bundle;
import android.util.Log;
import android.view.View;
import android.widget.CheckBox;
import android.widget.EditText;

import com.cardinfolink.yunshouyin.salesman.R;
import com.cardinfolink.yunshouyin.salesman.model.SAServerPacket;
import com.cardinfolink.yunshouyin.salesman.model.SaveData;
import com.cardinfolink.yunshouyin.salesman.model.SessonData;
import com.cardinfolink.yunshouyin.salesman.model.SystemConfig;
import com.cardinfolink.yunshouyin.salesman.model.User;
import com.cardinfolink.yunshouyin.salesman.utils.CommunicationListenerV2;
import com.cardinfolink.yunshouyin.salesman.utils.ErrorUtil;
import com.cardinfolink.yunshouyin.salesman.utils.HttpCommunicationUtil;
import com.cardinfolink.yunshouyin.salesman.utils.ParamsUtil;
import com.cardinfolink.yunshouyin.salesman.utils.VerifyUtil;
import com.umeng.update.UmengUpdateAgent;


public class LoginActivity extends BaseActivity {
    private final String TAG = "LoginActivity";
    private EditText mUsernameEdit;
    private EditText mPasswordEdit;
    private CheckBox mAutoLoginCheckBox;

    @SuppressLint("NewApi")
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        // check update
        UmengUpdateAgent.update(this);
        initEnvironment();
        setContentView(R.layout.login_activity);

        mUsernameEdit = (EditText) findViewById(R.id.login_username);
        VerifyUtil.addEmailLimit(mUsernameEdit);
        mPasswordEdit = (EditText) findViewById(R.id.login_password);
        VerifyUtil.addEmailLimit(mPasswordEdit);
        mAutoLoginCheckBox = (CheckBox) findViewById(R.id.checkbox_auto_login);
        User user = SaveData.getUser(mContext);
        mAutoLoginCheckBox.setChecked(user.isAutoLogin());
        mUsernameEdit.setText(user.getUsername());
        mPasswordEdit.setText(user.getPassword());

        if (user.isAutoLogin()) {
            login();
        }
    }

    /**
     *
     */
    private void initEnvironment(){
        try {
            ApplicationInfo ai = getPackageManager().getApplicationInfo(
                    getPackageName(), PackageManager.GET_META_DATA);
            Bundle bundle = ai.metaData;
            String environment = bundle.getString("ENVIRONMENT");
            Log.d(TAG, "ENVIRONMENT is " + environment);
            SystemConfig.initEnvironment(environment);
        } catch (PackageManager.NameNotFoundException e) {
            Log.e(TAG, "Failed to load meta-data, NameNotFound: " + e.getMessage());
        } catch (Exception e) {
            Log.e(TAG, "Failed to load meta-data: " + e.getMessage());
        }
    }

    public void BtnLoginOnClick(View view) {
        login();
    }

    /**
     * 验证用户名密码是否有值
     *
     * @return
     */
    @SuppressLint("NewApi")
    private boolean validate() {
        String username, password;
        username = mUsernameEdit.getText().toString();
        password = mPasswordEdit.getText().toString();
        if (username.isEmpty()) {
            alertError("用户名不能为空!");
            return false;
        }

        if (password.isEmpty()) {
            alertError("密码不能为空!");
            return false;
        }
        return true;
    }


    private void login() {
        Log.d(TAG, "======================login========================");
        // Test only
        //mUsernameEdit.setText("toolstest");
        //mPasswordEdit.setText("Yun#1016");


        if (validate()) {
            startLoading();

            final String username = mUsernameEdit.getText().toString();
            final String password = mPasswordEdit.getText().toString();

            /**
             * save to share preference
             */
            User user = new User();
            user.setUsername(username);
            user.setPassword(password);

            // 自动登录checkbox 保存密码
            if (mAutoLoginCheckBox.isChecked()) {
                user.setAutoLogin(true);
                user.setPassword(password);
            }

            SaveData.setUser(mContext, user);


            /**
             * save to session
             */
            SessonData.loginUser.setUsername(username);
            SessonData.loginUser.setPassword(password);


            /**
             * async network call and callbacks
             */
            HttpCommunicationUtil.sendDataToQuickIpayServer(ParamsUtil.getLogin_SA(username, password), new CommunicationListenerV2() {
                @Override
                public void onResult(SAServerPacket serverPacket) {
                    SessonData.loginUser.setAccessToken(serverPacket.getAccessToken());
                    runOnUiThread(new Runnable() {
                        @Override
                        public void run() {
                            endLoading();
                            Intent intent = new Intent(mContext, SAMerchantListActivity.class);
                            mContext.startActivity(intent);
                        }
                    });
                }

                @Override
                public void onError(final String error) {
                    runOnUiThread(new Runnable() {
                        @Override
                        public void run() {
                            if (error.equals("username_password_error")) {
                                mPasswordEdit.setText("");
                            }
                            // convert to user friendly message
                            String errorStr = ErrorUtil.getErrorString(error);
                            Log.i(TAG, "error:" + error);
                            endLoadingWithError(errorStr);
                        }
                    });
                }
            });
        }
    }



    @Override
    protected void onResume() {
        super.onResume();
        User user = SaveData.getUser(mContext);
        mAutoLoginCheckBox.setChecked(user.isAutoLogin());
        mUsernameEdit.setText(user.getUsername());
        mPasswordEdit.setText(user.getPassword());
    }
}
