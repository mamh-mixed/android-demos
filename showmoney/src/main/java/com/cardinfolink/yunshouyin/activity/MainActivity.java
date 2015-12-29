package com.cardinfolink.yunshouyin.activity;

import android.content.Intent;
import android.os.Bundle;
import android.os.Handler;
import android.os.Message;
import android.support.v4.view.PagerAdapter;
import android.support.v4.view.ViewPager;
import android.util.Log;
import android.view.Display;
import android.view.Gravity;
import android.view.KeyEvent;
import android.view.View;
import android.view.WindowManager;
import android.view.animation.Animation;
import android.view.animation.TranslateAnimation;
import android.widget.ImageView;
import android.widget.LinearLayout;
import android.widget.LinearLayout.LayoutParams;
import android.widget.Toast;

import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.constant.Msg;
import com.cardinfolink.yunshouyin.data.SessonData;
import com.cardinfolink.yunshouyin.view.MySettingView;
import com.cardinfolink.yunshouyin.view.ScanCodeView;
import com.cardinfolink.yunshouyin.view.TicketView;
import com.cardinfolink.yunshouyin.view.TransManageView;
import com.umeng.common.message.UmengMessageDeviceConfig;
import com.umeng.message.IUmengRegisterCallback;
import com.umeng.message.MsgConstant;
import com.umeng.message.PushAgent;
import com.umeng.update.UmengUpdateAgent;

import java.util.ArrayList;

public class MainActivity extends BaseActivity {
    private static final String TAG = "MainActivity";

    private ScanCodeView mScanCodeView;
    private TransManageView mTransManageView;
    private TicketView mTicketView;
    private MySettingView mMySettingView;

    private long exitTime = 0;

    private Handler handler = new Handler();
    private PushAgent mPushAgent;
    //此处是注册的回调处理
    //参考集成文档的1.7.10
    //http://dev.umeng.com/push/android/integration#1_7_10
    private IUmengRegisterCallback mRegisterCallback = new UmengPushAgengRegisterCallback();

    private ViewPager mTabPager;//声明对象
    private ImageView mTabImg;// 动画图片
    private ImageView mTab1, mTab2, mTab3, mTab4;
    private int zero = 0;// 动画图片偏移量
    private int currIndex = 0;// 当前页卡编号
    private int one;// 单个水平动画位移
    private int two;
    private int three;

    // 每个页面的view数据,存放4个界面
    private ArrayList<View> mViews;
    private Handler mHandler;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.main_activity);
        initHandler();
        initLayout();
        initUmeng();
    }

    public Handler getHandler() {
        return mHandler;
    }

    public void setHandler(Handler handler) {
        this.mHandler = handler;
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

    public void initHandler() {
        //mainactivity里面的handler主要是用来切换主界面上的四个界面的。从其他地方发个消息过来
        //在这里进行切换界面的操作。例如在扫码界面 要切换到账单界面，就
        mHandler = new Handler() {

            @Override
            public void handleMessage(Message msg) {
                super.handleMessage(msg);
                //// TODO: mamh  这里要进行切换界面的操作
            }
        };
    }

    private void initLayout() {
        //扫码的界面，第一个界面
        LinearLayout.LayoutParams layoutParams = new LayoutParams(LayoutParams.MATCH_PARENT, LayoutParams.MATCH_PARENT);
        mScanCodeView = new ScanCodeView(mContext, mHandler);
        mScanCodeView.setLayoutParams(layoutParams);

        //销券的界面，第二个界面
        mTicketView = new TicketView(mContext, mHandler);
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

        mTabImg = (ImageView) findViewById(R.id.img_tab_now);//动画图片

        mTab1.setOnClickListener(new MainPagerItemOnClickListener(0));
        mTab2.setOnClickListener(new MainPagerItemOnClickListener(1));
        mTab3.setOnClickListener(new MainPagerItemOnClickListener(2));
        mTab4.setOnClickListener(new MainPagerItemOnClickListener(3));

        Display currDisplay = getWindowManager().getDefaultDisplay();// 获取屏幕当前分辨率
        int displayWidth = currDisplay.getWidth();
        one = displayWidth / 4; // 设置水平动画平移大小
        two = one * 2;
        three = one * 3;

        // 每个页面的view数据
        mViews = new ArrayList<View>();
        mViews.add(mScanCodeView);
        mViews.add(mTicketView);
        mViews.add(mTransManageView);
        mViews.add(mMySettingView);

        // 填充ViewPager的数据适配器
        PagerAdapter mPagerAdapter = new MainPagerAdapter();
        mTabPager.setAdapter(mPagerAdapter);
        mTab1.setImageDrawable(getResources().getDrawable(R.drawable.gathering_selected));
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
        @Override
        public void onPageScrolled(int position, float positionOffset, int positionOffsetPixels) {

        }

        public void onPageSelected(int arg0) {
            Animation animation = null;
            switch (arg0) {
                case 0:
                    mTab1.setImageResource(R.drawable.gathering_selected);

                    if (currIndex == 1) {
                        animation = new TranslateAnimation(one, 0, 0, 0);
                        mTab2.setImageResource(R.drawable.ticket_not_selected);
                    } else if (currIndex == 2) {
                        animation = new TranslateAnimation(two, 0, 0, 0);
                        mTab3.setImageResource(R.drawable.bill_not_selected);
                    } else if (currIndex == 3) {
                        animation = new TranslateAnimation(three, 0, 0, 0);
                        mTab4.setImageResource(R.drawable.my_not_selected);
                    }
                    break;
                case 1:
                    mTab2.setImageResource(R.drawable.ticket_selected);
                    if (currIndex == 0) {
                        animation = new TranslateAnimation(zero, one, 0, 0);
                        mTab1.setImageResource(R.drawable.gathering_not_selected);
                    } else if (currIndex == 2) {
                        animation = new TranslateAnimation(two, one, 0, 0);
                        mTab3.setImageResource(R.drawable.bill_not_selected);
                    } else if (currIndex == 3) {
                        animation = new TranslateAnimation(three, one, 0, 0);
                        mTab4.setImageResource(R.drawable.my_not_selected);
                    }
                    break;
                case 2:
                    mTab3.setImageResource(R.drawable.bill_selected);
                    //刷新一下账单列表
                    mTransManageView.refresh();
                    if (currIndex == 0) {
                        animation = new TranslateAnimation(zero, two, 0, 0);
                        mTab1.setImageResource(R.drawable.gathering_not_selected);
                    } else if (currIndex == 1) {
                        animation = new TranslateAnimation(one, two, 0, 0);
                        mTab2.setImageResource(R.drawable.ticket_not_selected);
                    } else if (currIndex == 3) {
                        animation = new TranslateAnimation(three, two, 0, 0);
                        mTab4.setImageResource(R.drawable.my_not_selected);
                    }
                    break;
                case 3:
                    mTab4.setImageResource(R.drawable.my_selected);
                    if (currIndex == 0) {
                        animation = new TranslateAnimation(zero, three, 0, 0);
                        mTab1.setImageResource(R.drawable.gathering_not_selected);
                    } else if (currIndex == 1) {
                        animation = new TranslateAnimation(one, three, 0, 0);
                        mTab2.setImageResource(R.drawable.ticket_not_selected);
                    } else if (currIndex == 2) {
                        animation = new TranslateAnimation(two, three, 0, 0);
                        mTab3.setImageResource(R.drawable.bill_not_selected);
                    }
                    break;
            }
            currIndex = arg0;
            animation.setFillAfter(true);// True:图片停在动画结束位置
            animation.setDuration(150);// 动画持续时间
            mTabImg.startAnimation(animation);// 开始动画
        }

        @Override
        public void onPageScrollStateChanged(int state) {

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

}
