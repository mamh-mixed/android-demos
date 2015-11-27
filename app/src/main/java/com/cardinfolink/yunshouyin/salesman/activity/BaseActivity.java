package com.cardinfolink.yunshouyin.salesman.activity;

import android.app.Activity;
import android.content.Context;
import android.content.Intent;
import android.content.SharedPreferences;
import android.graphics.Bitmap;
import android.graphics.BitmapFactory;
import android.os.Bundle;
import android.support.v7.app.AppCompatActivity;
import android.view.View;

import com.cardinfolink.yunshouyin.salesman.R;
import com.cardinfolink.yunshouyin.salesman.api.QuickPayConfigStorage;
import com.cardinfolink.yunshouyin.salesman.core.BankDataService;
import com.cardinfolink.yunshouyin.salesman.core.BankDataServiceImpl;
import com.cardinfolink.yunshouyin.salesman.core.QiniuMultiUploadService;
import com.cardinfolink.yunshouyin.salesman.core.QuickPayService;
import com.cardinfolink.yunshouyin.salesman.db.SalesmanDB;
import com.cardinfolink.yunshouyin.salesman.model.User;
import com.cardinfolink.yunshouyin.salesman.utils.ActivityCollector;
import com.cardinfolink.yunshouyin.salesman.utils.SalesmanApplication;
import com.cardinfolink.yunshouyin.salesman.view.AlertDialog;
import com.cardinfolink.yunshouyin.salesman.view.LoadingDialog;
import com.cardinfolink.yunshouyin.salesman.view.WorkBeforeExitListener;
import com.umeng.analytics.MobclickAgent;


public class BaseActivity extends AppCompatActivity {
    private final String TAG = "BaseActivity";

    private LoadingDialog mLoadingDialog;    //显示loading
    private AlertDialog mAlertDialog;       // 提示消息对话框
    protected User loginUser = new User();
    protected SalesmanApplication application;
    protected Context mContext;
    protected SharedPreferences mSharedPreferences;//保存数据。
    protected SharedPreferences mRegisterSharedPreferences;//注册缓存

    protected QuickPayService quickPayService;
    protected BankDataService bankDataService;

    protected QiniuMultiUploadService qiniuMultiUploadService;

    protected SalesmanDB salesmanDB;//使用数据库存储省份城市银行信息

    //重载 setContentView 初始化 mLoadingDialog,mAlertDialog
    @Override
    public void setContentView(int layoutResID) {
        super.setContentView(layoutResID);

        mContext = this;
        mSharedPreferences = getSharedPreferences("savedata", Activity.MODE_PRIVATE);
        mRegisterSharedPreferences = getSharedPreferences("registerdata", Activity.MODE_PRIVATE);
        salesmanDB = SalesmanDB.getInstance(this);

        View loadingDialogView = findViewById(R.id.loading_dialog);
        mLoadingDialog = new LoadingDialog(this, loadingDialogView);

        Bitmap alertDialogBitmap = BitmapFactory.decodeResource(getResources(), R.drawable.wrong);
        String alertDialogMsg = getResources().getString(R.string.username_password_error);
        View alertDialogView = findViewById(R.id.alert_dialog);
        mAlertDialog = new AlertDialog(this, null, alertDialogView, alertDialogMsg, alertDialogBitmap);

        application = SalesmanApplication.getInstance();
        quickPayService = application.getQuickPayService();
        bankDataService = application.getBankDataService();
        qiniuMultiUploadService = application.getQiniuMultiUploadService();
    }

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        ActivityCollector.addActivity(this);
    }

    @Override
    protected void onDestroy() {
        super.onDestroy();
        ActivityCollector.removeActivity(this);
    }


    public void startLoading() {
        mLoadingDialog.startLoading();
    }

    public void endLoading() {
        mLoadingDialog.endLoading();
    }

    public void endLoadingWithError(String msg) {
        endLoading();
        alertError(msg);
    }

    public void alertError(String msg) {
        mAlertDialog.show(msg, BitmapFactory.decodeResource(this.getResources(), R.drawable.wrong));
    }

    public void alertInfo(String msg) {
        mAlertDialog.show(msg, BitmapFactory.decodeResource(this.getResources(), R.drawable.right));
    }

    public void alertInfo(String msg, WorkBeforeExitListener listener) {
        mAlertDialog.show(msg, BitmapFactory.decodeResource(this.getResources(), R.drawable.right), listener);
    }

    /**
     * umeng integration
     */
    protected void onResume() {
        super.onResume();
        MobclickAgent.onResume(this);
    }

    /**
     * umeng integration
     */
    protected void onPause() {
        super.onPause();
        MobclickAgent.onPause(this);
    }

    public User getLoginUser() {
        loginUser.setUsername(mSharedPreferences.getString("username", ""));
        loginUser.setPassword(mSharedPreferences.getString("password", ""));
        loginUser.setAutoLogin(mSharedPreferences.getBoolean("autologin", false));
        return loginUser;
    }

    public void setLoginUser(User user) {
        this.loginUser = user;
        SharedPreferences.Editor editor = mSharedPreferences.edit();
        editor.putString("username", user.getUsername());
        editor.putString("password", user.getPassword());
        editor.putBoolean("autologin", user.isAutoLogin());
        editor.commit();
    }


}

