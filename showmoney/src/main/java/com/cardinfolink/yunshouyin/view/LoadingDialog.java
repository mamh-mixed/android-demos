package com.cardinfolink.yunshouyin.view;


import android.content.Context;
import android.view.MotionEvent;
import android.view.View;
import android.view.animation.Animation;
import android.view.animation.AnimationUtils;
import android.widget.ImageView;

import com.cardinfolink.yunshouyin.R;

public class LoadingDialog {

    private View dialogView;
    private ImageView loadImg;
    private Context mContext;

    public LoadingDialog(Context context, View view) {
        dialogView = view;
        mContext = context;
        loadImg = (ImageView) dialogView.findViewById(R.id.load_img);
    }


    public void startLoading() {
        dialogView.setVisibility(View.VISIBLE);
        Animation loadingAnimation = AnimationUtils.loadAnimation(mContext, R.anim.loading_animation);
        loadImg.startAnimation(loadingAnimation);
        dialogView.setOnTouchListener(new View.OnTouchListener() {

            @Override
            public boolean onTouch(View v, MotionEvent event) {
                //加入这个 在loading的时候，点击其他任何地方都会没有反应的。
                return true;
            }
        });
    }

    public void endLoading() {
        dialogView.setVisibility(View.GONE);
        loadImg.clearAnimation();
    }

}
