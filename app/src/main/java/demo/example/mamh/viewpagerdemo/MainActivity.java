package demo.example.mamh.viewpagerdemo;

import android.app.Activity;
import android.os.Bundle;
import android.support.v4.view.PagerAdapter;
import android.support.v4.view.ViewPager;
import android.view.View;
import android.view.ViewGroup;
import android.widget.ImageView;
import android.widget.LinearLayout;
import android.widget.TextView;

import java.util.ArrayList;

public class MainActivity extends Activity {

    private ViewPager viewPager;

    private LinearLayout pointGroup;

    private TextView msg;

    private ArrayList<ImageView> imageList;
    private int[] imageIds = new int[]{
            R.drawable.a, R.drawable.b, R.drawable.c, R.drawable.d, R.drawable.e

    };

    //图片标题集合
    private final String[] imageDescriptions = {
            "巩俐不低俗，我就不能低俗",
            "扑树又回来啦！再唱经典老歌引万人大合唱",
            "揭秘北京电影如何升级",
            "乐视网TV版大派送",
            "热血屌丝的反杀"
    };
    private ArrayList<ImageView> pointList;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);
        viewPager = (ViewPager) findViewById(R.id.view_pager);
        pointGroup = (LinearLayout) findViewById(R.id.point_group);
        msg = (TextView) findViewById(R.id.msg);
        imageList = new ArrayList<ImageView>();
        pointList = new ArrayList<ImageView>();
        for (int i = 0; i < imageIds.length; i++) {
            ImageView iv = new ImageView(this);
            iv.setBackgroundResource(imageIds[i]);
            imageList.add(iv);
            ImageView point = new ImageView(this);
            point.setBackgroundResource(R.mipmap.ic_launcher);
            point.setVisibility(View.GONE);
            pointList.add(point);
            pointGroup.addView(point);

        }
        msg.setText(imageDescriptions[0]);
        pointList.get(0).setVisibility(View.VISIBLE);

        PagerAdapter adapter = new ViewPagerAdapter();
        viewPager.setAdapter(adapter);
        viewPager.addOnPageChangeListener(new ViewPager.OnPageChangeListener() {
            @Override
            public void onPageScrolled(int position, float positionOffset, int positionOffsetPixels) {

            }

            @Override
            public void onPageSelected(int position) {
                msg.setText(imageDescriptions[position]);
                pointList.get(position).setVisibility(View.VISIBLE);
            }

            @Override
            public void onPageScrollStateChanged(int state) {

            }
        });

    }


    private class ViewPagerAdapter extends PagerAdapter {
        @Override
        public Object instantiateItem(ViewGroup container, int position) {
            //实例化某个条目
            //获得相应位置上的item
            //container view的容器，viewpager自身
            container.addView(imageList.get(position));
            return imageList.get(position);
        }

        @Override
        public void destroyItem(ViewGroup container, int position, Object object) {
            //销毁对应位置上的obj
            container.removeView((View) object);
            object = null;
        }

        @Override
        public int getCount() {
            //告诉viewpager 有多少个条目。多少个页面
            return imageList.size();
        }

        @Override
        public boolean isViewFromObject(View view, Object object) {
            return view == object;
        }
    }
}
