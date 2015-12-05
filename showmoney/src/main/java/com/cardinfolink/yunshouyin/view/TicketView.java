package com.cardinfolink.yunshouyin.view;

import android.content.Context;
import android.view.LayoutInflater;
import android.view.View;
import android.widget.LinearLayout;

import com.cardinfolink.yunshouyin.R;

/**
 * 销券的界面
 * Created by mamh on 15-12-7.
 */
public class TicketView extends LinearLayout {
    private static final String TAG = "TicketView";
    private Context mContext;

    public TicketView(Context context) {
        super(context);
        mContext = context;

        View contentView = LayoutInflater.from(context).inflate(R.layout.ticket_view, null);
        LinearLayout.LayoutParams layoutParams = new LayoutParams(LayoutParams.MATCH_PARENT, LayoutParams.MATCH_PARENT);
        contentView.setLayoutParams(layoutParams);
        addView(contentView);

    }




}
