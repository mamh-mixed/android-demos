package com.cardinfolink.yunshouyin.view;

import android.content.Context;
import android.view.MotionEvent;
import android.view.View;
import android.widget.Button;

import com.cardinfolink.yunshouyin.R;

/**
 * Created by mamh on 15-12-22.
 */
public class SelectPicDialog {
    private Context mContext;
    private View dialogView;

    private Button mTakePhoto, mPickPhoto, mCancel;

    public SelectPicDialog(Context context, View view) {
        mContext = context;
        dialogView = view;

        mTakePhoto = (Button) dialogView.findViewById(R.id.select_pic_take_photo);
        mPickPhoto = (Button) dialogView.findViewById(R.id.select_pic_pick_photo);
        mCancel = (Button) dialogView.findViewById(R.id.select_pic_cancel);


        dialogView.setOnTouchListener(new View.OnTouchListener() {

            @Override
            public boolean onTouch(View v, MotionEvent event) {
                return true;
            }
        });


        mCancel.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                hide();
            }
        });

    }

    public void show() {
        dialogView.setVisibility(View.VISIBLE);
    }

    public void hide() {
        dialogView.setVisibility(View.GONE);
    }

    public void setTakePhotoOnClickListener(View.OnClickListener l) {
        mTakePhoto.setOnClickListener(l);
    }

    public void setPickPhotoOnClickListener(View.OnClickListener l) {
        mPickPhoto.setOnClickListener(l);
    }

    public void setCancelnClickListener(View.OnClickListener l) {
        mCancel.setOnClickListener(l);
    }
}

