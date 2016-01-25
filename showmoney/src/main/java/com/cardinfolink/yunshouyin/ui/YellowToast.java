package com.cardinfolink.yunshouyin.ui;

import android.content.Context;
import android.view.Gravity;
import android.view.LayoutInflater;
import android.view.View;
import android.widget.TextView;
import android.widget.Toast;

import com.cardinfolink.yunshouyin.R;

/**
 * Created by mamh on 16-1-25.
 */
public class YellowToast {
    private Context mContext;

    private View mLayout;
    private TextView mTipsTextView;
    private Toast mToast;


    public YellowToast(Context context) {
        this.mContext = context;

        LayoutInflater inflater = LayoutInflater.from(context);
        mLayout = inflater.inflate(R.layout.yellow_toast_layout, null);

        mTipsTextView = (TextView) mLayout.findViewById(R.id.tips_textview);

        mToast = new Toast(mContext);

        mToast.setView(mLayout);
        mToast.setGravity(Gravity.TOP | Gravity.FILL_HORIZONTAL, 0, 0);
        mToast.setDuration(Toast.LENGTH_SHORT);
    }

    public void show(String msg) {
        mTipsTextView.setText(msg);
        mToast.show();
    }

}
