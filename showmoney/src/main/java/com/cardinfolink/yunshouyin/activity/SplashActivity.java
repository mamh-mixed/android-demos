package com.cardinfolink.yunshouyin.activity;

import android.content.Context;
import android.content.Intent;
import android.content.SharedPreferences;
import android.os.Bundle;
import android.os.Handler;

import com.cardinfolink.yunshouyin.R;
import com.umeng.analytics.MobclickAgent;
import com.umeng.message.IUmengRegisterCallback;
import com.umeng.message.PushAgent;

public class SplashActivity extends BaseActivity {
    private static final String TAG = "SplashActivity";
    private static final int SPLASH_DISPLAY_LENGHT = 3000; //延迟三秒


    private Handler handler = new Handler();
    //此处是注册的回调处理
    //参考集成文档的1.7.10
    //http://dev.umeng.com/push/android/integration#1_7_10
    private IUmengRegisterCallback mRegisterCallback = new UmengPushAgengRegisterCallback();


    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_splash);
        MobclickAgent.updateOnlineConfig(mContext);

        initUmeng();

        SharedPreferences sp = getSharedPreferences("savedata", Context.MODE_PRIVATE);
        Boolean isFirst = sp.getBoolean("is_user_guide_show", true);

        //默认是登录，除非从其他地方传人了original的值
        String originalFromFlag = "login";
        try {
            Intent intent = getIntent();
            Bundle bundle = intent.getExtras();
            originalFromFlag = bundle.getString("original");
        } catch (Exception e) {

        }

        if ("login".equals(originalFromFlag)) {
            if (isFirst) {
                new Handler().postDelayed(new Runnable() {
                    @Override
                    public void run() {
                        Intent mainIntent = new Intent(SplashActivity.this, GuideActivity.class);
                        startActivity(mainIntent);
                        finish();
                    }
                }, SPLASH_DISPLAY_LENGHT);
            } else {
                new Handler().postDelayed(new Runnable() {
                    @Override
                    public void run() {
                        Intent mainIntent = new Intent(SplashActivity.this, LoginActivity.class);
                        startActivity(mainIntent);
                        finish();
                    }
                }, SPLASH_DISPLAY_LENGHT);
            }
        } else {
            new Handler().postDelayed(new Runnable() {

                @Override
                public void run() {
                    finish();
                }

            }, SPLASH_DISPLAY_LENGHT);
        }

    }

    private void initUmeng() {
        PushAgent mPushAgent = PushAgent.getInstance(this);

        //应用程序启动统计
        //参考集成文档的1.5.1.2
        //http://dev.umeng.com/push/android/integration#1_5_1
        mPushAgent.setResourcePackageName("com.cardinfolink.yunshouyin");
        mPushAgent.onAppStart();

        //开启推送并设置注册的回调处理
        mPushAgent.enable(mRegisterCallback);
        mPushAgent.setMergeNotificaiton(false);//不合并消息 通知，这样通知栏会有多条消息显示

        /**
         *    4.2.3  获取设备的device_token（可选）
         * 可以在Debug模式下输出的logcat中看到device_token，也可以使用下面的方法来获取device_token。
         * String device_token = UmengRegistrar.getRegistrationId(context)
         *
         *说明
         *
         * device_token为友盟生成的用于标识设备的id，长度为44位，不能定制和修改。同一台设备上每个应用对应的device_token不一样。
         * 获取device_token的代码需要放在mPushAgent.enable();后面，注册成功以后调用才能获得device_token。
         * 如果返回值为空，说明设备还没有注册成功， 需要等待几秒钟，同时请确保测试手机网络畅通。
         * 如果一直都获取不到device_token，请参考： 安卓获取不到device_token
         * 在回调函数中获取测试设备的device_token可参考下文5.1部分进行获取。
         *
         */
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
}
