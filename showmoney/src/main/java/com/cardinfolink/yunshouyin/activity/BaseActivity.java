package com.cardinfolink.yunshouyin.activity;

import android.app.Activity;
import android.content.Context;
import android.graphics.Bitmap;
import android.graphics.BitmapFactory;
import android.os.Bundle;
import android.view.View;

import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.core.BankDataService;
import com.cardinfolink.yunshouyin.core.QiniuMultiUploadService;
import com.cardinfolink.yunshouyin.core.QuickPayService;
import com.cardinfolink.yunshouyin.util.ActivityCollector;
import com.cardinfolink.yunshouyin.util.ShowMoneyApp;
import com.cardinfolink.yunshouyin.view.AlertDialog;
import com.cardinfolink.yunshouyin.view.LoadingDialog;
import com.umeng.analytics.MobclickAgent;

public class BaseActivity extends Activity {
    private static final String TAG = "BaseActivity";

    protected LoadingDialog mLoadingDialog;    //显示loading
    protected AlertDialog mAlertDialog;       // 提示消息对话框
    protected Context mContext;
    protected QuickPayService quickPayService;
    protected BankDataService bankDataService;
    protected ShowMoneyApp yunApplication;

    protected QiniuMultiUploadService qiniuMultiUploadService;

    //重载 setContentView 初始化 mLoadingDialog,mAlertDialog
    @Override
    public void setContentView(int layoutResID) {
        super.setContentView(layoutResID);
        mContext = this;

        View LoadView = findViewById(R.id.loading_dialog);
        mLoadingDialog = new LoadingDialog(this, LoadView);

        View alertView = findViewById(R.id.alert_dialog);
        String alertMsg = getResources().getString(R.string.username_password_error);
        Bitmap alertBitmap = BitmapFactory.decodeResource(this.getResources(), R.drawable.wrong);
        mAlertDialog = new AlertDialog(this, null, alertView, alertMsg, alertBitmap);

        yunApplication = (ShowMoneyApp) getApplication();
        quickPayService = yunApplication.getQuickPayService();
        bankDataService = yunApplication.getBankDataService();
        qiniuMultiUploadService = yunApplication.getQiniuMultiUploadService();

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

    public void alertShow(String msg, Bitmap bitmap) {
        mAlertDialog.show(msg, bitmap);
    }

    @Override
    protected void onResume() {
        super.onResume();
        //友盟统计
        MobclickAgent.onResume(this);
    }

    @Override
    protected void onPause() {
        super.onPause();
        //友盟统计
        MobclickAgent.onPause(this);
    }
}
