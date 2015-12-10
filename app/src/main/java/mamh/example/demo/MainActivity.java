package mamh.example.demo;

import android.app.Activity;
import android.support.v7.app.AppCompatActivity;
import android.os.Bundle;
import android.view.KeyEvent;
import android.view.View;
import android.view.animation.RotateAnimation;
import android.widget.ImageView;
import android.widget.RelativeLayout;

import static android.view.KeyEvent.KEYCODE_MENU;

public class MainActivity extends Activity implements View.OnClickListener {

    private ImageView iconHome;
    private ImageView iconMenu;

    private RelativeLayout level1;
    private RelativeLayout level2;
    private RelativeLayout level3;
    private boolean isLevel3Show = true;//用于标记三级菜单是否显示
    private boolean isLevel2Show = true;
    private boolean isLevel1Show = true;


    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);

        iconHome = (ImageView) findViewById(R.id.icon_home);
        iconMenu = (ImageView) findViewById(R.id.icon_menu);

        level1 = (RelativeLayout) findViewById(R.id.level1);
        level2 = (RelativeLayout) findViewById(R.id.level2);
        level3 = (RelativeLayout) findViewById(R.id.level3);

        iconHome.setOnClickListener(this);
        iconMenu.setOnClickListener(this);

    }

    @Override
    public void onClick(View v) {
        switch (v.getId()) {
            case R.id.icon_menu:
                //显示与隐藏第三集菜单
                if (isLevel3Show) {
                    //hide
                    startAnimationOut(level3, 0);
                } else {
                    //show
                    startAnimationIn(level3);
                }
                isLevel3Show = !isLevel3Show;
                break;
            case R.id.icon_home:
                if (isLevel2Show) {
                    startAnimationOut(level2, 500);
                    isLevel2Show = false;
                    if (isLevel3Show) {
                        //hide
                        startAnimationOut(level3, 0);
                        isLevel3Show = false;
                    }

                } else {
                    startAnimationIn(level2);
                    isLevel2Show = true;

                }

                break;
        }

    }

    @Override
    public boolean onKeyDown(int keyCode, KeyEvent event) {
        if (keyCode == KeyEvent.KEYCODE_MENU
                && event.getAction() == KeyEvent.ACTION_DOWN) {
            changeLevel1();
        }
        return super.onKeyDown(keyCode, event);
    }

    private void changeLevel1() {
        //改变第一季菜单状态
        if (isLevel1Show) {
            startAnimationOut(level1, 0);
            isLevel1Show = false;
            if (isLevel2Show) {
                startAnimationOut(level2, 300);
                isLevel2Show = false;
                if (isLevel3Show) {
                    startAnimationOut(level3, 500);
                    isLevel3Show = false;
                }
            }
        } else {
            startAnimationIn(level1);
            startAnimationIn(level2, 200);
            startAnimationIn(level3, 500);
        }
    }

    /**
     * 旋转进入的动画
     *
     * @param v
     */
    private void startAnimationIn(View v) {
        startAnimationIn(v, 0);
    }

    private void startAnimationIn(View v, int i) {
        //顺时针进入 是从 180  度 ---  360 度
        float x = v.getWidth() / 2;
        float y = v.getHeight();
        float from = 180;
        float to = 360;
        RotateAnimation rotateAnimation = new RotateAnimation(from, to, x, y);
        rotateAnimation.setDuration(500);//设置动画时间
        rotateAnimation.setFillAfter(true);//动画执行晚后保持最后的状态
        rotateAnimation.setStartOffset(i);
        v.startAnimation(rotateAnimation);
    }

    /**
     * 让指定的view执行旋转离开的动画
     *
     * @param v
     */
    private void startAnimationOut(View v, int i) {
        float x = v.getWidth() / 2;
        float y = v.getHeight();
        float from = 0;
        float to = 180;
        RotateAnimation rotateAnimation = new RotateAnimation(from, to, x, y);
        //顺时针离开 是从 0  度 ---  180 度
        rotateAnimation.setDuration(500);//设置动画时间
        rotateAnimation.setFillAfter(true);//动画执行晚后保持最后的状态
        rotateAnimation.setStartOffset(i);
        v.startAnimation(rotateAnimation);
    }

}
