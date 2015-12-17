package com.cardinfolink.yunshouyin.activity;

import android.content.Context;
import android.content.Intent;
import android.os.Bundle;
import android.os.Handler;
import android.util.Log;
import android.view.Gravity;
import android.view.KeyEvent;
import android.view.View;
import android.view.inputmethod.InputMethodManager;
import android.widget.AdapterView;
import android.widget.AdapterView.OnItemClickListener;
import android.widget.ArrayAdapter;
import android.widget.LinearLayout;
import android.widget.LinearLayout.LayoutParams;
import android.widget.ListView;
import android.widget.Toast;

import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.data.SessonData;
import com.cardinfolink.yunshouyin.util.ShowMoneyApp;
import com.cardinfolink.yunshouyin.view.AccountUpdateView;
import com.cardinfolink.yunshouyin.view.PasswordUpdateView;
import com.cardinfolink.yunshouyin.view.ScanCodeView;
import com.cardinfolink.yunshouyin.view.TransManageView;
import com.cardinfolink.yunshouyin.view.WapView;
import com.jeremyfeinstein.slidingmenu.lib.SlidingMenu;
import com.umeng.common.message.UmengMessageDeviceConfig;
import com.umeng.message.IUmengRegisterCallback;
import com.umeng.message.MsgConstant;
import com.umeng.message.PushAgent;
import com.umeng.update.UmengUpdateAgent;

import java.util.ArrayList;

public class MainActivity extends BaseActivity {
    private static final String TAG = "MainActivity";

    private SlidingMenu mLeftMenu;
    private ScanCodeView mScanCodeView;
    private TransManageView mTransManageView;
    private PasswordUpdateView mPasswordUpdateView;
    private AccountUpdateView mAccountUpdateView;
    private WapView mWapBillView;

    private LinearLayout mMainContent;
    private ListView mDrawerList;
    private ArrayList<String> menuLists;
    private ArrayAdapter<String> adapter;
    private long exitTime = 0;


    private Handler handler = new Handler();
    private PushAgent mPushAgent;
    //此处是注册的回调处理
    //参考集成文档的1.7.10
    //http://dev.umeng.com/push/android/integration#1_7_10
    private IUmengRegisterCallback mRegisterCallback = new UmengPushAgengRegisterCallback();


    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.main_activity);
        initLayout();
        initUmeng();
    }


    private void initUmeng() {
        UmengUpdateAgent.setUpdateOnlyWifi(false);
        UmengUpdateAgent.setUpdateCheckConfig(false);
        UmengUpdateAgent.update(this);

        mPushAgent = PushAgent.getInstance(this);
        //mPushAgent.setPushCheck(true);    //默认不检查集成配置文件
        //mPushAgent.setLocalNotificationIntervalLimit(false);  //默认本地通知间隔最少是10分钟

        //应用程序启动统计
        //参考集成文档的1.5.1.2
        //http://dev.umeng.com/push/android/integration#1_5_1
        mPushAgent.setResourcePackageName("com.cardinfolink.yunshouyin");
        mPushAgent.onAppStart();

        //开启推送并设置注册的回调处理
        if (!mPushAgent.isEnabled()) {
            mPushAgent.enable(mRegisterCallback);
        }
        String info = String.format("enabled:%s \n isRegistered:%s \n DeviceToken:%s " +
                        "SdkVersion:%s  \n AppVersionCode:%s \n AppVersionName:%s",
                mPushAgent.isEnabled(), mPushAgent.isRegistered(),
                mPushAgent.getRegistrationId(), MsgConstant.SDK_VERSION,
                UmengMessageDeviceConfig.getAppVersionCode(this), UmengMessageDeviceConfig.getAppVersionName(this));
        Log.e(TAG, "pushAgeng info :" + info);
    }

    private void initLayout() {
        mMainContent = (LinearLayout) findViewById(R.id.main_content);

        LinearLayout.LayoutParams layoutParams = new LayoutParams(LayoutParams.MATCH_PARENT, LayoutParams.MATCH_PARENT);
        mScanCodeView = new ScanCodeView(mContext);
        mScanCodeView.setLayoutParams(layoutParams);

        mTransManageView = new TransManageView(mContext);
        mTransManageView.setLayoutParams(new LayoutParams(LayoutParams.MATCH_PARENT, LayoutParams.MATCH_PARENT));

        mPasswordUpdateView = new PasswordUpdateView(mContext);
        mPasswordUpdateView.setLayoutParams(new LayoutParams(LayoutParams.MATCH_PARENT, LayoutParams.MATCH_PARENT));

        mAccountUpdateView = new AccountUpdateView(mContext);
        mAccountUpdateView.setLayoutParams(new LayoutParams(LayoutParams.MATCH_PARENT, LayoutParams.MATCH_PARENT));

        mWapBillView = new WapView(mContext);
        mWapBillView.setLayoutParams(new LayoutParams(LayoutParams.MATCH_PARENT, LayoutParams.MATCH_PARENT));


        mMainContent.addView(mScanCodeView);

        mLeftMenu = new SlidingMenu(this);
        mLeftMenu.setMode(SlidingMenu.LEFT);
        mLeftMenu.setTouchModeAbove(SlidingMenu.TOUCHMODE_FULLSCREEN);
        mLeftMenu.setShadowWidthRes(R.dimen.shadow_width);
        mLeftMenu.setShadowDrawable(R.drawable.shadow);

        // 设置滑动菜单视图的宽度
        mLeftMenu.setBehindOffsetRes(R.dimen.slidingmenu_offset);
        // 设置渐入渐出效果的值
        mLeftMenu.setFadeDegree(0.35f);
        /**
         * SLIDING_WINDOW will include the Title/ActionBar in the content
         * section of the SlidingMenu, while SLIDING_CONTENT does not.
         */
        mLeftMenu.attachToActivity(this, SlidingMenu.SLIDING_CONTENT);
        // 为侧滑菜单设置布局
        mLeftMenu.setMenu(R.layout.leftmenu);
        mLeftMenu.setOnOpenListener(new SlidingMenuOnOpenListener());

        mDrawerList = (ListView) mLeftMenu.findViewById(R.id.left_drawer);
        menuLists = new ArrayList<String>();

        menuLists.add(ShowMoneyApp.getResString(R.string.main_activity_menu_scancode));
        menuLists.add(ShowMoneyApp.getResString(R.string.main_activity_menu_transmange));
        menuLists.add(ShowMoneyApp.getResString(R.string.main_activity_menu_passwordupdate));
        menuLists.add(ShowMoneyApp.getResString(R.string.main_activity_menu_accountupdate));
        menuLists.add(ShowMoneyApp.getResString(R.string.main_activity_menu_webbill));
        menuLists.add(ShowMoneyApp.getResString(R.string.main_activity_menu_safeexit));


        adapter = new ArrayAdapter<String>(this, R.layout.menu_list_item, menuLists);
        mDrawerList.setAdapter(adapter);
        mDrawerList.setOnItemClickListener(new MenuOnItemClick());
    }


    public void BtnMenuOnClick(View view) {
        hideInput();

        if (mLeftMenu.isMenuShowing()) {
            mLeftMenu.toggle();
        } else {
            mLeftMenu.showMenu();
        }
    }

    private void openView(int position) {
        SessonData.positionView = position;
        switch (position) {
            case 0:
                mMainContent.removeAllViews();
                mMainContent.addView(mScanCodeView);
                mScanCodeView.clearValue();
                break;
            case 1:
                mMainContent.removeAllViews();
                mMainContent.addView(mTransManageView);
                mTransManageView.refresh();
                break;
            case 2:
                mMainContent.removeAllViews();
                mMainContent.addView(mPasswordUpdateView);
                break;
            case 3:
                mMainContent.removeAllViews();
                mAccountUpdateView = new AccountUpdateView(mContext);
                mAccountUpdateView.setLayoutParams(new LayoutParams(LayoutParams.MATCH_PARENT, LayoutParams.MATCH_PARENT));
                mMainContent.addView(mAccountUpdateView);
                mAccountUpdateView.getInfo();
                break;
            case 4:
                mMainContent.removeAllViews();
                mMainContent.addView(mWapBillView);
                mWapBillView.initData();
                break;
        }
    }

    @Override
    public boolean onKeyDown(int keyCode, KeyEvent event) {
        if (keyCode == KeyEvent.KEYCODE_BACK && event.getAction() == KeyEvent.ACTION_DOWN) {
            if ((System.currentTimeMillis() - exitTime) > 2000) {
                String pressText = getResources().getString(R.string.press_again_exit);
                Toast toast = Toast.makeText(getApplicationContext(), pressText, Toast.LENGTH_SHORT);
                toast.setGravity(Gravity.CENTER, 0, 250);
                toast.show();
                exitTime = System.currentTimeMillis();
            } else {
                //finish();
                Intent intent = new Intent(Intent.ACTION_MAIN);
                intent.addCategory(Intent.CATEGORY_HOME);
                intent.setFlags(Intent.FLAG_ACTIVITY_CLEAR_TOP);
                startActivity(intent);
                android.os.Process.killProcess(android.os.Process.myPid());
            }
            return true;
        }
        return super.onKeyDown(keyCode, event);
    }

    @Override
    protected void onResume() {
        super.onResume();
        openView(SessonData.positionView);
    }

    private class MenuOnItemClick implements OnItemClickListener {

        @Override
        public void onItemClick(AdapterView<?> parent, View view, int position, long id) {
            hideInput();

            SessonData.positionView = position;
            switch (position) {
                case 0:
                    mMainContent.removeAllViews();
                    mMainContent.addView(mScanCodeView);
                    mScanCodeView.clearValue();
                    mLeftMenu.toggle();
                    break;
                case 1:
                    mMainContent.removeAllViews();
                    mMainContent.addView(mTransManageView);
                    mTransManageView.initData();
                    mLeftMenu.toggle();
                    break;
                case 2:
                    mMainContent.removeAllViews();
                    mMainContent.addView(mPasswordUpdateView);
                    mLeftMenu.toggle();
                    break;
                case 3:
                    mMainContent.removeAllViews();
                    mAccountUpdateView = new AccountUpdateView(mContext);
                    mAccountUpdateView.setLayoutParams(new LayoutParams(
                            LayoutParams.MATCH_PARENT, LayoutParams.MATCH_PARENT));
                    mMainContent.addView(mAccountUpdateView);
                    mAccountUpdateView.getInfo();
                    mLeftMenu.toggle();
                    break;
                case 4:
                    mMainContent.removeAllViews();
                    mMainContent.addView(mWapBillView);
                    mWapBillView.initData();
                    mLeftMenu.toggle();
                    break;
                case 5:
                    finish();
                    break;
            }

        }

    }


    private class UmengPushAgengRegisterCallback implements IUmengRegisterCallback {

        @Override
        public void onRegistered(String s) {
            handler.post(new Runnable() {

                @Override
                public void run() {
                    Log.e(TAG, " =========onRegistered(String s) =============run()");
                }
            });
        }
    }

    private class SlidingMenuOnOpenListener implements SlidingMenu.OnOpenListener {

        @Override
        public void onOpen() {
            hideInput();
        }
    }

    private void hideInput() {
        try {
            InputMethodManager inputMethodManager = (InputMethodManager) getSystemService(Context.INPUT_METHOD_SERVICE);
            inputMethodManager.hideSoftInputFromWindow(getCurrentFocus().getWindowToken(), InputMethodManager.HIDE_NOT_ALWAYS);
        } catch (Exception e) {
            e.printStackTrace();
        }

    }
}
