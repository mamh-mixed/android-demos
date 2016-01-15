package com.cardinfolink.yunshouyin.view;

import android.content.Context;
import android.os.Handler;
import android.os.Message;
import android.view.View;
import android.widget.TextView;

import com.cardinfolink.yunshouyin.R;


/**
 * 这个使用自定义toast 不是更好？？？
 */
public class YellowTips {
    private Context mContext;
    private TextView mTextView;
    private View mTipsView;
    private Handler mHandler;
    private final int HIDE_TIPS = 100;

    public YellowTips(Context context, View view) {
        mContext = context;
        mTipsView = view;
        mTextView = (TextView) mTipsView.findViewById(R.id.tips_textview);

        mHandler = new Handler() {
            public void handleMessage(Message msg) {
                switch (msg.what) {
                    case HIDE_TIPS:
                        hide();
                        break;
                }

                super.handleMessage(msg);
            }
        };
    }

    public void show(String text) {
        mTextView.setText(text);
        mTipsView.setVisibility(View.VISIBLE);
        mHandler.sendEmptyMessageDelayed(HIDE_TIPS, 3000);
    }

    private void hide() {
        mTipsView.setVisibility(View.GONE);
        mTextView.setText("");
    }


}
