package com.cardinfolink.yunshouyin.view;

import android.content.Context;
import android.content.Intent;
import android.view.MotionEvent;
import android.view.View;
import android.view.View.OnClickListener;
import android.view.View.OnTouchListener;

import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.activity.ResetPasswordActivity;

public class ResetDialog {
    private Context mContext;
    private View dialogView;

    public ResetDialog(Context context, View view) {
        mContext = context;
        dialogView = view;
    }

    public void show() {
        dialogView.setVisibility(View.VISIBLE);

        dialogView.setOnTouchListener(new OnTouchListener() {

            @Override
            public boolean onTouch(View v, MotionEvent event) {
                return true;
            }
        });

        dialogView.findViewById(R.id.reset_dialog_ok).setOnClickListener(new OnClickListener() {

            @Override
            public void onClick(View v) {
                Intent intent = new Intent(mContext, ResetPasswordActivity.class);
                mContext.startActivity(intent);
                dialogView.setVisibility(View.GONE);
            }
        });

        dialogView.findViewById(R.id.reset_dialog_cancel).setOnClickListener(new OnClickListener() {
            @Override
            public void onClick(View v) {
                dialogView.setVisibility(View.GONE);
            }
        });
    }

}
