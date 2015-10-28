package com.cardinfolink.yunshouyin.salesman.activities;

import android.content.Context;
import android.content.Intent;
import android.graphics.BitmapFactory;
import android.os.Bundle;
import android.os.PersistableBundle;
import android.support.v7.app.AppCompatActivity;
import android.util.Log;

import com.cardinfolink.yunshouyin.salesman.R;
import com.cardinfolink.yunshouyin.salesman.utils.ActivityCollector;
import com.cardinfolink.yunshouyin.salesman.views.Alert_Dialog;
import com.cardinfolink.yunshouyin.salesman.views.Loading_Dialog;
import com.umeng.analytics.MobclickAgent;


public class BaseActivity extends AppCompatActivity {

    private Loading_Dialog mLoading_Dialog;    //显示loading
    private Alert_Dialog mAlert_Dialog;       // 提示消息对话框
    protected Context mContext;


    //重载 setContentView 初始化 mLoading_Dialog,mAlert_Dialog
    @Override
    public void setContentView(int layoutResID) {
        super.setContentView(layoutResID);
        mContext = this;
        mLoading_Dialog = new Loading_Dialog(this, findViewById(R.id.loading_dialog));
        mAlert_Dialog = new Alert_Dialog(this, null, findViewById(R.id.alert_dialog),
                getResources().getString(R.string.username_password_error), BitmapFactory.decodeResource(this.getResources(), R.drawable.wrong));
    }

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        Log.d("BaseActivity", getClass().getSimpleName());
        ActivityCollector.addActivity(this);
    }

    @Override
    protected void onDestroy() {
        super.onDestroy();
        ActivityCollector.removeActivity(this);
    }

    public void intentToActivity(Class cls) {
        Intent intent = new Intent(
                mContext,
                cls);
        mContext.startActivity(intent);
    }

    public void startLoading() {
        mLoading_Dialog.startLoading();
    }

    public void endLoading() {
        mLoading_Dialog.endLoading();
    }

    public void endLoadingWithError(String msg) {
        endLoading();
        alertError(msg);
    }

    public void alertError(String msg) {
        mAlert_Dialog.show(msg, BitmapFactory.decodeResource(this.getResources(), R.drawable.wrong));
    }

    public void alertInfo(String msg) {
        mAlert_Dialog.show(msg, BitmapFactory.decodeResource(this.getResources(), R.drawable.right));
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
}

