package com.cardinfolink.yunshouyin.view;

import android.content.Context;
import android.view.MotionEvent;
import android.view.View;
import android.view.View.OnClickListener;
import android.view.View.OnTouchListener;
import android.widget.TextView;

import com.cardinfolink.yunshouyin.R;

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
