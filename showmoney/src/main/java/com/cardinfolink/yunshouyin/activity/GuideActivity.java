package com.cardinfolink.yunshouyin.activity;

import android.app.Activity;
import android.content.Context;
import android.content.Intent;
import android.content.SharedPreferences;
import android.os.Bundle;
import android.support.v4.view.PagerAdapter;
import android.support.v4.view.ViewPager;
import android.view.View;
import android.view.ViewGroup;
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

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        sp = getSharedPreferences("savedata", Context.MODE_PRIVATE);
        setContentView(R.layout.activity_guide);
        mViewPager = (ViewPager) findViewById(R.id.vp_guide);
        mStart = (Button) findViewById(R.id.btn_start);
        initView();

        mViewPager.setAdapter(new GuideAdapter());
        mViewPager.setOnPageChangeListener(new GuidePageListener());

    }

    private void initView() {
        imageViewList = new ArrayList<ImageView>();

        for (int i = 0; i < mImageIds.length; i++) {
            ImageView imageView = new ImageView(this);
            imageView.setBackgroundResource(mImageIds[i]);
            imageViewList.add(imageView);
        }

        mStart.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                SharedPreferences.Editor editor = sp.edit();
                editor.putBoolean("is_user_guide_show", false).commit();
                Intent intent = new Intent(GuideActivity.this, LoginActivity.class);
                startActivity(intent);
                finish();
            }
        });

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

        }

        @Override
        public void onPageSelected(int position) {
            if (position == mImageIds.length - 1) {//最后一个界面
                mStart.setVisibility(View.VISIBLE);
            } else {
                mStart.setVisibility(View.INVISIBLE);
            }
        }

        @Override
        public void onPageScrollStateChanged(int state) {

        }
    }
}
