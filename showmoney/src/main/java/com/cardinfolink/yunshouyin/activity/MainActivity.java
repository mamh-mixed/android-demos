package com.cardinfolink.yunshouyin.activity;

import android.content.Context;
import android.content.Intent;
import android.content.SharedPreferences;
import android.os.Bundle;
import android.os.Handler;
import android.os.Message;
import android.support.v4.view.PagerAdapter;
import android.support.v4.view.ViewPager;
import android.view.Gravity;
import android.view.KeyEvent;
import android.view.View;
import android.view.WindowManager;
import android.widget.ImageView;
import android.widget.LinearLayout;
import android.widget.LinearLayout.LayoutParams;
import android.widget.Toast;

import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.constant.Msg;
import com.cardinfolink.yunshouyin.view.MySettingView;
import com.cardinfolink.yunshouyin.view.ScanCodeView;
import com.cardinfolink.yunshouyin.view.TicketView;
import com.cardinfolink.yunshouyin.view.TransManageView;
import com.umeng.message.IUmengRegisterCallback;
import com.umeng.message.PushAgent;
import com.umeng.update.UmengUpdateAgent;

import java.util.ArrayList;

public class MainActivity extends BaseActivity {
    private static final String TAG = "MainActivity";
    public boolean mFirsInTicketView = true;
    private ScanCodeView mScanCodeView;
    private TransManageView mTransManageView;
    private TicketView mTicketView;
    private MySettingView mMySettingView;

    private long exitTime = 0;

    private Handler handler = new Handler();
    //此处是注册的回调处理
    //参考集成文档的1.7.10
    //http://dev.umeng.com/push/android/integration#1_7_10
    private IUmengRegisterCallback mRegisterCallback = new UmengPushAgengRegisterCallback();

    private ViewPager mTabPager;//声明对象
    private ImageView mTab1, mTab2, mTab3, mTab4;
    private int currIndex = 0;// 当前页卡编号

    // 每个页面的view数据,存放4个界面
    private ArrayList<View> mViews;

    private static Handler mMainActivityHandler;//main activity里面的handler，用来切换界面的
    private SharedPreferences sp;

    @Override
    protected void onResume() {
        super.onResume();
        getWindow().setSoftInputMode(WindowManager.LayoutParams.SOFT_INPUT_ADJUST_PAN);


        if (mMySettingView != null && mTabPager.getCurrentItem() == 3) {
            mMySettingView.checkMessageCount();
        }
    }

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.main_activity);
        sp = getSharedPreferences("savedata", Context.MODE_PRIVATE);
        mFirsInTicketView = sp.getBoolean("mFirsInTicketView", true);
        initHandler();
        initLayout();
        initUmeng();
    }


    public static Handler getHandler() {
        return mMainActivityHandler;
    }

    private void initUmeng() {
        UmengUpdateAgent.setDefault();
        UmengUpdateAgent.setUpdateOnlyWifi(false);
        UmengUpdateAgent.setUpdateCheckConfig(false);
        UmengUpdateAgent.update(this);

        PushAgent mPushAgent = PushAgent.getInstance(this);
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
    }

    public void initHandler() {
        //mainactivity里面的handler主要是用来切换主界面上的四个界面的。从其他地方发个消息过来
        //在这里进行切换界面的操作。例如在扫码界面 要切换到账单界面，就
        mMainActivityHandler = new Handler() {

            @Override
            public void handleMessage(Message msg) {
                super.handleMessage(msg);
                switch (msg.what) {
                    case Msg.MSG_FROM_SERVER_COUPON_SUCCESS:
                        mTabPager.setCurrentItem(0);
                        break;
                    case Msg.MSG_GO_SCAN_CODE_VIEW:
                        mTabPager.setCurrentItem(0);
                        break;
                    case Msg.MSG_GO_TICKET_VIEW:
                        mTabPager.setCurrentItem(1);
                        break;
                    case Msg.MSG_GO_BILL_VIEW:
                        //去账单界面
                        mTabPager.setCurrentItem(2);
                        break;
                    case Msg.MSG_GO_MY_SETTING_VIEW:
                        mTabPager.setCurrentItem(3);
                        break;
                    case Msg.MSG_REFRESH_BILL_LIST_VIEW:
                        if (mTransManageView != null && mTabPager.getCurrentItem() == 2) {
                            mTransManageView.refresh();
                        }

                }
            }
        };
    }

    private void initLayout() {
        //扫码的界面，第一个界面
        LinearLayout.LayoutParams layoutParams = new LayoutParams(LayoutParams.MATCH_PARENT, LayoutParams.MATCH_PARENT);
        mScanCodeView = new ScanCodeView(mContext, mMainActivityHandler);
        mScanCodeView.setLayoutParams(layoutParams);

        //销券的界面，第二个界面
        mTicketView = new TicketView(mContext, mMainActivityHandler);
        mTicketView.setLayoutParams(layoutParams);

        //账单界面，第三个界面
        mTransManageView = new TransManageView(mContext);
        mTransManageView.setLayoutParams(layoutParams);

        //我的设置界面，第四个界面
        mMySettingView = new MySettingView(mContext);
        mMySettingView.setLayoutParams(layoutParams);

        // 启动activity时不自动弹出软键盘
        getWindow().setSoftInputMode(WindowManager.LayoutParams.SOFT_INPUT_STATE_ALWAYS_HIDDEN);

        mTabPager = (ViewPager) findViewById(R.id.tabpager);
        mTabPager.addOnPageChangeListener(new MainPagerOnPageChangeListener());

        mTab1 = (ImageView) findViewById(R.id.img_gathering);
        mTab2 = (ImageView) findViewById(R.id.img_ticket);
        mTab3 = (ImageView) findViewById(R.id.img_bill);
        mTab4 = (ImageView) findViewById(R.id.img_my);

        mTab1.setOnClickListener(new MainPagerItemOnClickListener(0));
        mTab2.setOnClickListener(new MainPagerItemOnClickListener(1));
        mTab3.setOnClickListener(new MainPagerItemOnClickListener(2));
        mTab4.setOnClickListener(new MainPagerItemOnClickListener(3));


        // 每个页面的view数据
        mViews = new ArrayList<View>();
        mViews.add(mScanCodeView);
        mViews.add(mTicketView);
        mViews.add(mTransManageView);
        mViews.add(mMySettingView);

        // 填充ViewPager的数据适配器
        PagerAdapter mPagerAdapter = new MainPagerAdapter();
        mTabPager.setAdapter(mPagerAdapter);
        mTab1.setImageResource(R.drawable.gathering_selected);

        mTransManageView.refresh();
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


    private class MainPagerAdapter extends android.support.v4.view.PagerAdapter {
        @Override
        public boolean isViewFromObject(View arg0, Object arg1) {
            return arg0 == arg1;
        }

        @Override
        public int getCount() {
            return mViews.size();
        }

        @Override
        public void destroyItem(View container, int position, Object object) {
            ((ViewPager) container).removeView(mViews.get(position));
        }


        @Override
        public Object instantiateItem(View container, int position) {
            ((ViewPager) container).addView(mViews.get(position));
            return mViews.get(position);
        }
    }

    /**
     * 头标点击监听
     */
    private class MainPagerItemOnClickListener implements View.OnClickListener {
        private int index = 0;

        public MainPagerItemOnClickListener(int i) {
            index = i;
        }

        public void onClick(View v) {
            mTabPager.setCurrentItem(index);
        }
    }


    /*
     * 页卡切换监听
     */
    private class MainPagerOnPageChangeListener implements ViewPager.OnPageChangeListener {
        public void onPageSelected(int position) {
            switch (position) {
                case 0:
                    mTab1.setImageResource(R.drawable.gathering_selected);

                    if (currIndex == 1) {
                        mTab2.setImageResource(R.drawable.ticket_not_selected);
                    } else if (currIndex == 2) {
                        mTab3.setImageResource(R.drawable.bill_not_selected);
                    } else if (currIndex == 3) {
                        mTab4.setImageResource(R.drawable.my_not_selected);
                    }
                    break;
                case 1:
                    mTab2.setImageResource(R.drawable.ticket_selected);
                    if (currIndex == 0) {
                        mTab1.setImageResource(R.drawable.gathering_not_selected);
                    } else if (currIndex == 2) {
                        mTab3.setImageResource(R.drawable.bill_not_selected);
                    } else if (currIndex == 3) {
                        mTab4.setImageResource(R.drawable.my_not_selected);
                    }
                    if (mFirsInTicketView) {
                        mTicketView.showCouponHintDialog();
                    }
                    SharedPreferences.Editor editor = sp.edit();
                    editor.putBoolean("mFirsInTicketView", false).commit();
                    mFirsInTicketView = false;
                    break;
                case 2:
                    mTab3.setImageResource(R.drawable.bill_selected);
                    if (currIndex == 0) {
                        mTab1.setImageResource(R.drawable.gathering_not_selected);
                    } else if (currIndex == 1) {
                        mTab2.setImageResource(R.drawable.ticket_not_selected);
                    } else if (currIndex == 3) {
                        mTab4.setImageResource(R.drawable.my_not_selected);
                    }
                    break;
                case 3:
                    mTab4.setImageResource(R.drawable.my_selected);
                    if (currIndex == 0) {
                        mTab1.setImageResource(R.drawable.gathering_not_selected);
                    } else if (currIndex == 1) {
                        mTab2.setImageResource(R.drawable.ticket_not_selected);
                    } else if (currIndex == 2) {
                        mTab3.setImageResource(R.drawable.bill_not_selected);
                    }
                    break;
            }
            currIndex = position;
        }

        @Override
        public void onPageScrollStateChanged(int state) {

        }

        @Override
        public void onPageScrolled(int position, float positionOffset, int positionOffsetPixels) {

        }

    }


    private class UmengPushAgengRegisterCallback implements IUmengRegisterCallback {

        @Override
        public void onRegistered(String s) {
            handler.post(new Runnable() {

                @Override
                public void run() {
                }
            });
        }
    }

    @Override
    protected void onActivityResult(int requestCode, int resultCode, Intent data) {
        if (data != null) {
            String ticketcode = data.getStringExtra("ticketcode");
            mTicketView.setTicketCode(ticketcode);
        }
    }
}
