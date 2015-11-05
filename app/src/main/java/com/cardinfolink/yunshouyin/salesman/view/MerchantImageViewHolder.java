package com.cardinfolink.yunshouyin.salesman.view;

import android.support.v7.app.AppCompatActivity;
import android.support.v7.widget.RecyclerView;
import android.util.Log;
import android.view.View;
import android.widget.ImageView;

import com.baoyz.actionsheet.ActionSheet;
import com.cardinfolink.yunshouyin.salesman.R;
import com.cardinfolink.yunshouyin.salesman.activity.SARegisterStep3Activity;

public class MerchantImageViewHolder extends RecyclerView.ViewHolder implements View.OnClickListener, View.OnLongClickListener, ActionSheet.ActionSheetListener {
    private static final String TAG = "MerchantImageViewHolder";
    public ImageView merchantPhoto;

    private int itemIndex;
    private AppCompatActivity activity;

    public MerchantImageViewHolder(View itemView) {
        super(itemView);
        itemView.setOnClickListener(this);
        itemView.setOnLongClickListener(this);
        this.merchantPhoto = (ImageView) itemView.findViewById(R.id.merchant_photo);
    }

    @Override
    public void onClick(View v) {
        Log.d(TAG, "clicked position = " + getAdapterPosition());
    }

    @Override
    public boolean onLongClick(View v) {
        Log.d(TAG, "long clicked position = " + getLayoutPosition());

        activity = (AppCompatActivity) v.getContext();
        itemIndex = getAdapterPosition();
        Log.d("testpos", "long clicked AdapterPosition: "+itemIndex);
        Log.d("testpos", "long clicked LayoutPosition: "+getLayoutPosition());

        ActionSheet.createBuilder(activity, activity.getSupportFragmentManager())
                .setOtherButtonTitles("删除")
                .setCancelButtonTitle("取消")
                .setCancelableOnTouchOutside(true)
                .setListener(this).show();

        return true;
    }

    @Override
    public void onDismiss(ActionSheet actionSheet, boolean isCancel) {
        if (isCancel) {
            Log.d(TAG, "cancel clicked");
        }
    }

    @Override
    public void onOtherButtonClick(ActionSheet actionSheet, int index) {
        switch (index) {
            case 0:
                if (activity instanceof SARegisterStep3Activity) {
                    Log.d(TAG, "delete item at position: "+itemIndex);
                    SARegisterStep3Activity SARegisterStep3Activity = (SARegisterStep3Activity) activity;
                    SARegisterStep3Activity.removeItemAt(itemIndex);
                }
                break;
            default:
                break;
        }
    }
}
