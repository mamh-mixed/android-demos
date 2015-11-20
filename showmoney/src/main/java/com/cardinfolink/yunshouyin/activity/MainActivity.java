package com.cardinfolink.yunshouyin.activity;

import android.content.Intent;
import android.os.Bundle;
import android.view.Gravity;
import android.view.KeyEvent;
import android.view.View;
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
import com.cardinfolink.yunshouyin.view.LimitIncreaseView;
import com.cardinfolink.yunshouyin.view.PasswordUpdateView;
import com.cardinfolink.yunshouyin.view.ScanCodeView;
import com.cardinfolink.yunshouyin.view.TransManageView;
import com.cardinfolink.yunshouyin.view.WapView;
import com.jeremyfeinstein.slidingmenu.lib.SlidingMenu;
import com.umeng.update.UmengUpdateAgent;

import java.util.ArrayList;

public class MainActivity extends BaseActivity {
    SlidingMenu mLeftMenu;
    private ScanCodeView mScanCodeView;
    private TransManageView mTransManageView;
    private PasswordUpdateView mPasswordUpdateView;
    private AccountUpdateView mAccountUpdateView;
    private LimitIncreaseView mLimitIncreaseView;
    private WapView mWapBillView;

    private LinearLayout mMainContent;
    private ListView mDrawerList;
    private ArrayList<String> menuLists;
    private ArrayAdapter<String> adapter;
    private long exitTime = 0;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.main_activity);
        initLayout();
        UmengUpdateAgent.setUpdateOnlyWifi(false);
        UmengUpdateAgent.setUpdateCheckConfig(false);
        UmengUpdateAgent.update(this);
    }

    private void initLayout() {
        mMainContent = (LinearLayout) findViewById(R.id.main_content);

        LinearLayout.LayoutParams layoutParams = new LayoutParams(
                LayoutParams.MATCH_PARENT, LayoutParams.MATCH_PARENT);
        mScanCodeView = new ScanCodeView(mContext);
        mScanCodeView.setLayoutParams(layoutParams);

        mTransManageView = new TransManageView(mContext);
        mTransManageView.setLayoutParams(new LayoutParams(
                LayoutParams.MATCH_PARENT, LayoutParams.MATCH_PARENT));

        mPasswordUpdateView = new PasswordUpdateView(mContext);
        mPasswordUpdateView.setLayoutParams(new LayoutParams(
                LayoutParams.MATCH_PARENT, LayoutParams.MATCH_PARENT));

        mAccountUpdateView = new AccountUpdateView(mContext);
        mAccountUpdateView.setLayoutParams(new LayoutParams(
                LayoutParams.MATCH_PARENT, LayoutParams.MATCH_PARENT));

        mLimitIncreaseView = new LimitIncreaseView(mContext);
        mLimitIncreaseView.setLayoutParams(new LayoutParams(
                LayoutParams.MATCH_PARENT, LayoutParams.MATCH_PARENT));

        mWapBillView = new WapView(mContext);
        mWapBillView.setLayoutParams(new LayoutParams(
                LayoutParams.MATCH_PARENT, LayoutParams.MATCH_PARENT));


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

        mDrawerList = (ListView) mLeftMenu.findViewById(R.id.left_drawer);
        menuLists = new ArrayList<String>();

        menuLists.add(ShowMoneyApp.getResString(R.string.main_activity_menu_scancode));
        menuLists.add(ShowMoneyApp.getResString(R.string.main_activity_menu_transmange));
        menuLists.add(ShowMoneyApp.getResString(R.string.main_activity_menu_passwordupdate));
        menuLists.add(ShowMoneyApp.getResString(R.string.main_activity_menu_accountupdate));
        //menuLists.add("限额提升");
        menuLists.add(ShowMoneyApp.getResString(R.string.main_activity_menu_webbill));
        menuLists.add(ShowMoneyApp.getResString(R.string.main_activity_menu_safeexit));


        adapter = new ArrayAdapter<String>(this, R.layout.menu_list_item, menuLists);
        mDrawerList.setAdapter(adapter);
        mDrawerList.setOnItemClickListener(new MenuOnItemClick());

    }


    public void BtnMenuOnClick(View view) {

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
                mAccountUpdateView.setLayoutParams(new LayoutParams(
                        LayoutParams.MATCH_PARENT, LayoutParams.MATCH_PARENT));
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
        if (keyCode == KeyEvent.KEYCODE_BACK
                && event.getAction() == KeyEvent.ACTION_DOWN) {
            if ((System.currentTimeMillis() - exitTime) > 2000) {
                Toast toast = Toast.makeText(getApplicationContext(),
                        getResources().getString(R.string.press_again_exit),
                        Toast.LENGTH_SHORT);
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
        public void onItemClick(AdapterView<?> parent, View view, int position,
                                long id) {
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
//					 mLeftMenu.toggle();
//					 Uri uri;
//					 uri = Uri.parse(SystemConfig.WEB_BILL_URL+"?merchantCode="+SessonData.loginUser.getObjectId());
//					 Intent  intent = new  Intent(Intent.ACTION_VIEW, uri);
//					 startActivity(intent);

                    break;
                case 5:
                    finish();
                    break;


            }

        }

    }

}
