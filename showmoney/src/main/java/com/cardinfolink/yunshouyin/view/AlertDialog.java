package com.cardinfolink.yunshouyin.view;

import android.content.Context;
import android.graphics.Bitmap;
import android.graphics.drawable.Drawable;
import android.os.Handler;
import android.text.TextUtils;
import android.view.MotionEvent;
import android.view.View;
import android.view.View.OnClickListener;
import android.view.View.OnTouchListener;
import android.widget.ImageView;
import android.widget.TextView;

import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.constant.Msg;
import com.cardinfolink.yunshouyin.util.ShowMoneyApp;

/**
 * 警告对话框，上面一个图片，中间一行文本，下面一个按钮
 */
public class AlertDialog {
    private Context mContext;
    private Handler mHandler;
    private View dialogView;
    private String mMessage;
    private Bitmap mBitmap;

    private TextView mTitle;
    private ImageView mImageView;
    private TextView mOk;

    private OnClickListener mOkOnClickListener;

    public AlertDialog(Context context, Handler handler, View view, String message, Bitmap bitmap) {
        mContext = context;
        mHandler = handler;
        dialogView = view;
        mMessage = message;
        mBitmap = bitmap;

        mTitle = (TextView) dialogView.findViewById(R.id.alert_message);
        mOk = (TextView) dialogView.findViewById(R.id.alert_ok);
        mImageView = (ImageView) dialogView.findViewById(R.id.alert_img);

        setTitle(mMessage);

        if (mBitmap != null) {
            setImageViewBitmap(mBitmap);
        }

        //初始化就设置这个事件
        dialogView.setOnTouchListener(new OnTouchListener() {

            @Override
            public boolean onTouch(View v, MotionEvent event) {
                return true;
            }
        });

        //ok按钮的默认行为
        mOkOnClickListener = new OnClickListener() {

            @Override
            public void onClick(View v) {
                hide();
                if (mHandler != null) {
                    mHandler.sendEmptyMessage(Msg.MSG_FROM_CLIENT_ALERT_OK);
                }

            }
        };
        mOk.setOnClickListener(mOkOnClickListener);
    }

    /**
     * alert对话框的默认行为。
     * 默认按钮的行为是 关闭然后发送了handler消息
     */
    public void show() {
        dialogView.setVisibility(View.VISIBLE);
    }


    /**
     * 这是另外一种常用的方式，自定义显示的title，自定义显示的图片
     * 按钮的行为默认是关闭对话框。
     *
     * @param message
     * @param bitmap
     */
    public void show(String message, Bitmap bitmap) {
        mMessage = message;
        if (TextUtils.isEmpty(message)) {
            mMessage = ShowMoneyApp.getResString(R.string.server_timeout);
        }
        setTitle(mMessage);

        if (bitmap != null) {
            setImageViewBitmap(bitmap);
        }

        //也可以这样设置号按钮的行为
        setOkOnClickListener(new OnClickListener() {
            @Override
            public void onClick(View v) {
                hide();
            }
        });

        show();//这里显示对话框
    }

    /**
     * 这是一个种完全自定义的对话框
     *
     * @param message
     * @param ok
     * @param bitmap
     * @param listener
     */
    public void show(String message, String ok, Bitmap bitmap, OnClickListener listener) {
        setTitle(message);
        setOkText(ok);
        setImageViewBitmap(bitmap);

        setOkOnClickListener(listener);

        show();
    }

    public void hide() {
        dialogView.setVisibility(View.GONE);
    }

    public void setOkOnClickListener(OnClickListener l) {
        mOkOnClickListener = l;
        mOk.setOnClickListener(mOkOnClickListener);
    }

    public void setTitle(String title) {
        mTitle.setText(title);
    }

    public void setOkText(String ok) {
        mOk.setText(ok);
    }

    /**
     * 设置图片的方法
     *
     * @param d
     */
    public void setImageViewDrawable(Drawable d) {
        mImageView.setImageDrawable(d);
    }

    // 设置图片的方法
    public void setImageViewResource(int id) {
        mImageView.setImageResource(id);
    }

    // 设置图片的方法
    public void setImageViewBitmap(Bitmap bitmap) {
        mImageView.setImageBitmap(bitmap);
    }
}
