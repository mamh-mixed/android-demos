package com.cardinfolink.yunshouyin.activity;

import android.app.Activity;
import android.content.Context;
import android.content.Intent;
import android.content.SharedPreferences;
import android.os.Bundle;
import android.os.Handler;
import android.support.v4.view.PagerAdapter;
import android.support.v4.view.ViewPager;
import android.util.Log;
import android.view.MotionEvent;
import android.view.View;
import android.view.ViewGroup;
import android.view.WindowManager;
import android.widget.Button;
import android.widget.ImageView;

import com.cardinfolink.yunshouyin.R;
import com.jeremyfeinstein.slidingmenu.lib.CustomViewAbove;

import java.util.ArrayList;

/**
 * Created by charles on 2016/1/5.
 */
public class GuideActivity extends Activity {

    private static final int[] mImageIds = new int[]{R.drawable.guide1,
            R.drawable.guide2, R.drawable.guide3};

    private ViewPager mViewPager;
    private ArrayList<ImageView> imageViewList;
    private Button mStart;
    private SharedPreferences sp;
    private int mFirsDownX = 0;
    private int mCurrentUpX = 0;
    private Context mContext;
    private String original;
    private static final int GUIDE_DISPLAY_LENGHT = 1000;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        sp = getSharedPreferences("savedata", Context.MODE_PRIVATE);
        setContentView(R.layout.activity_guide);
        mContext = this;
        mViewPager = (ViewPager) findViewById(R.id.vp_guide);
        mViewPager.setOnTouchListener(new View.OnTouchListener() {
            @Override
            public boolean onTouch(View v, MotionEvent event) {
                return false;
            }
        });
        mStart = (Button) findViewById(R.id.btn_start);
        initView();

        mViewPager.setAdapter(new GuideAdapter());
        mViewPager.setOnPageChangeListener(new GuidePageListener());
        Intent intent=getIntent();
        Bundle bundle=intent.getExtras();
        original = bundle.getString("original");

    }

    private void initView() {
        imageViewList = new ArrayList<ImageView>();

        for (int i = 0; i < mImageIds.length; i++) {
            ImageView imageView = new ImageView(this);
            imageView.setBackgroundResource(mImageIds[i]);
            imageViewList.add(imageView);
        }


    }

    /***
     * Viewpager的适配器
     */
    class GuideAdapter extends PagerAdapter {

        @Override
        public int getCount() {
            return mImageIds.length;
        }

        @Override
        public boolean isViewFromObject(View view, Object object) {
            return view == object;
        }

        @Override
        public Object instantiateItem(ViewGroup container, int position) {
            container.addView(imageViewList.get(position));
            return imageViewList.get(position);
        }

        @Override
        public void destroyItem(ViewGroup container, int position, Object object) {
            container.removeView((View) object);
        }
    }

    /**
     * ViewPager的滑动监听
     */
    class GuidePageListener implements ViewPager.OnPageChangeListener {

        @Override
        public void onPageScrolled(int position, float positionOffset, int positionOffsetPixels) {
            if (position == mImageIds.length - 1) {
                    if("AboutActivity".equals(original)) {
                        new Handler().postDelayed(new Runnable() {
                            @Override
                            public void run() {
                                finish();
                            }
                        },GUIDE_DISPLAY_LENGHT);

                    }else {
                        SharedPreferences.Editor editor = sp.edit();
                        editor.putBoolean("is_user_guide_show", false).commit();
                        Intent intent = new Intent(GuideActivity.this, LoginActivity.class);
                        startActivity(intent);
                        finish();
                    }
            }
        }

        @Override
        public void onPageSelected(int position) {

        }

        @Override
        public void onPageScrollStateChanged(int state) {

        }
    }
}
