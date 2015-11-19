package com.cardinfolink.yunshouyin.view;


import android.content.Context;
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
        Animation loadingAnimation = AnimationUtils.loadAnimation(
                mContext, R.anim.loading_animation);
        loadImg.startAnimation(loadingAnimation);


    }

    public void endLoading() {
        dialogView.setVisibility(View.GONE);
        loadImg.clearAnimation();
    }

}
