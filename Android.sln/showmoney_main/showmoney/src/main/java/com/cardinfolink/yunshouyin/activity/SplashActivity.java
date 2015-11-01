package com.cardinfolink.yunshouyin.activity;

import android.content.Intent;
import android.os.Bundle;
import android.os.Handler;

import com.cardinfolink.yunshouyin.R;
import com.umeng.analytics.MobclickAgent;

public class SplashActivity extends BaseActivity {
    private final int SPLASH_DISPLAY_LENGHT = 3000; //延迟三秒

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.splash_activity);
        MobclickAgent.updateOnlineConfig(mContext);


        new Handler().postDelayed(new Runnable() {

            @Override
            public void run() {
                Intent mainIntent = new Intent(SplashActivity.this, LoginActivity.class);
                SplashActivity.this.startActivity(mainIntent);
                SplashActivity.this.finish();
            }

        }, SPLASH_DISPLAY_LENGHT);

    }

}
