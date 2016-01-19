package com.cardinfolink.yunshouyin.activity;

import android.content.Context;
import android.content.Intent;
import android.content.SharedPreferences;
import android.os.Bundle;
import android.os.Handler;

import com.cardinfolink.yunshouyin.R;
import com.umeng.analytics.MobclickAgent;

public class SplashActivity extends BaseActivity {
    private static final String TAG = "SplashActivity";
    private static final int SPLASH_DISPLAY_LENGHT = 3000; //延迟三秒

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_splash);
        MobclickAgent.updateOnlineConfig(mContext);

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


}
