package com.cardinfolink.yunshouyin.activity;

import android.app.Activity;
import android.content.Context;
import android.graphics.Bitmap;
import android.graphics.BitmapFactory;

import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.core.QuickPayService;
import com.cardinfolink.yunshouyin.util.ShowMoneyApp;
import com.cardinfolink.yunshouyin.view.AlertDialog;
import com.cardinfolink.yunshouyin.view.LoadingDialog;
import com.umeng.analytics.MobclickAgent;

public class BaseActivity extends Activity {

    protected LoadingDialog mLoading_Dialog;    //显示loading
    protected AlertDialog mAlert_Dialog;       // 提示消息对话框
    protected Context mContext;
    protected QuickPayService quickPayService;
    protected ShowMoneyApp yunApplication;


    //重载 setContentView 初始化 mLoading_Dialog,mAlert_Dialog
    @Override
    public void setContentView(int layoutResID) {
        super.setContentView(layoutResID);
        mContext = this;
        mLoading_Dialog = new LoadingDialog(this, findViewById(R.id.loading_dialog));
        mAlert_Dialog = new AlertDialog(this, null, findViewById(R.id.alert_dialog),
                getResources().getString(R.string.username_password_error), BitmapFactory.decodeResource(this.getResources(), R.drawable.wrong));

        yunApplication = (ShowMoneyApp)getApplication();
        quickPayService = yunApplication.getQuickPayService();
    }


    public void startLoading() {
        mLoading_Dialog.startLoading();
    }

    public void endLoading() {
        mLoading_Dialog.endLoading();
    }

    public void alertShow(String msg, Bitmap bitmap) {
        mAlert_Dialog.show(msg, bitmap);
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
