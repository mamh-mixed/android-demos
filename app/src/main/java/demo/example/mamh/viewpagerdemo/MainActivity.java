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

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);
        viewPager = (ViewPager) findViewById(R.id.view_pager);
        pointGroup = (LinearLayout) findViewById(R.id.point_group);
        msg = (TextView) findViewById(R.id.msg);
        imageList = new ArrayList<ImageView>();
        for (int i = 0; i < imageIds.length; i++) {
            ImageView iv = new ImageView(this);
            iv.setBackgroundResource(imageIds[i]);
            imageList.add(iv);
        }
        PagerAdapter adapter = new ViewPagerAdapter();
        viewPager.setAdapter(adapter);


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
