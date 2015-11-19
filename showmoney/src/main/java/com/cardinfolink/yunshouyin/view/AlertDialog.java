package com.cardinfolink.yunshouyin.view;

import android.content.Context;
import android.graphics.Bitmap;
import android.os.Handler;
import android.view.MotionEvent;
import android.view.View;
import android.view.View.OnClickListener;
import android.view.View.OnTouchListener;
import android.widget.ImageView;
import android.widget.TextView;

import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.constant.Msg;
import com.cardinfolink.yunshouyin.util.ShowMoneyApp;

public class AlertDialog {
    private Context mContext;
    private Handler mHandler;
    private View dialogView;
    private String mMessage;
    private Bitmap mBitmap;

    public AlertDialog(Context context, Handler handler, View view,
                       String message, Bitmap bitmap) {
        mContext = context;
        mHandler = handler;
        dialogView = view;
        mMessage = message;
        mBitmap = bitmap;
    }

    public void show() {
        TextView textView = (TextView) dialogView.findViewById(R.id.alert_message);
        textView.setText(mMessage);
        ImageView imageView = (ImageView) dialogView.findViewById(R.id.alert_img);
        if (mBitmap != null) {
            imageView.setImageBitmap(mBitmap);
        }
        dialogView.setVisibility(View.VISIBLE);

        dialogView.setOnTouchListener(new OnTouchListener() {

            @Override
            public boolean onTouch(View v, MotionEvent event) {
                // TODO Auto-generated method stub
                return true;
            }
        });

        dialogView.findViewById(R.id.alert_ok).setOnClickListener(
                new OnClickListener() {

                    @Override
                    public void onClick(View v) {
                        dialogView.setVisibility(View.GONE);
                        if (mHandler != null) {
                            mHandler.sendEmptyMessage(Msg.MSG_FROM_CLIENT_ALERT_OK);
                        }

                    }
                });
    }


    public void show(String message, Bitmap bitmap) {

        TextView textView = (TextView) dialogView.findViewById(R.id.alert_message);
        if (message.length() == 0) {
            message = ShowMoneyApp.getResString(R.string.server_timeout);
        }
        textView.setText(message);
        ImageView imageView = (ImageView) dialogView.findViewById(R.id.alert_img);
        if (bitmap != null) {
            imageView.setImageBitmap(bitmap);
        }
        dialogView.setVisibility(View.VISIBLE);

        dialogView.setOnTouchListener(new OnTouchListener() {

            @Override
            public boolean onTouch(View v, MotionEvent event) {
                // TODO Auto-generated method stub
                return true;
            }
        });

        dialogView.findViewById(R.id.alert_ok).setOnClickListener(
                new OnClickListener() {

                    @Override
                    public void onClick(View v) {
                        dialogView.setVisibility(View.GONE);

                    }
                });
    }
}
