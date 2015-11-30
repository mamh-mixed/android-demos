package com.cardinfolink.yunshouyin.activity;

import android.app.Activity;
import android.content.Context;
import android.graphics.Bitmap;
import android.graphics.BitmapFactory;

import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.core.BankDataService;
import com.cardinfolink.yunshouyin.core.QuickPayService;
import com.cardinfolink.yunshouyin.util.ShowMoneyApp;
import com.cardinfolink.yunshouyin.view.AlertDialog;
import com.cardinfolink.yunshouyin.view.LoadingDialog;
import com.umeng.analytics.MobclickAgent;

public class BaseActivity extends Activity {

    protected LoadingDialog mLoadingDialog;    //显示loading
    protected AlertDialog mAlertDialog;       // 提示消息对话框
    protected Context mContext;
    protected QuickPayService quickPayService;
    protected BankDataService bankDataService;
    protected ShowMoneyApp yunApplication;


    //重载 setContentView 初始化 mLoadingDialog,mAlertDialog
    @Override
    public void setContentView(int layoutResID) {
        super.setContentView(layoutResID);
        mContext = this;
        mLoadingDialog = new LoadingDialog(this, findViewById(R.id.loading_dialog));
        mAlertDialog = new AlertDialog(this, null, findViewById(R.id.alert_dialog),
                getResources().getString(R.string.username_password_error), BitmapFactory.decodeResource(this.getResources(), R.drawable.wrong));

        yunApplication = (ShowMoneyApp) getApplication();
        quickPayService = yunApplication.getQuickPayService();
        bankDataService = yunApplication.getBankDataService();
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
