package com.cardinfolink.yunshouyin.view;

import android.content.Context;
import android.os.Handler;
import android.os.Message;
import android.view.View;
import android.view.View.OnClickListener;
import android.view.animation.Animation;
import android.view.animation.AnimationUtils;
import android.widget.ImageView;
import android.widget.TextView;

import com.cardinfolink.cashiersdk.listener.CashierListener;
import com.cardinfolink.cashiersdk.model.OrderData;
import com.cardinfolink.cashiersdk.model.ResultData;
import com.cardinfolink.cashiersdk.sdk.CashierSdk;
import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.constant.Msg;

/**
 * 和交易相关的几个对话框。
 * load对话框界面,上面一个转圈的load的图片，中间一段文本，下面一个按钮
 */
public class TradingLoadDialog {
    private Context mContext;
    private Handler mHandler;
    private Handler dialogHandler;
    private boolean isLoading = false;
    private int second = 0;

    private String mOrderNum;

    private View dialogView;//父布局view
    private TextView mCancel;//取消订单的按钮
    private ImageView mLoadImage;//显示load图片
    private TextView mLoadMessage;//显示文本消息的，位于图片下面
    private TextView mLoadSecond;//显示计数秒数的文本


    public TradingLoadDialog(Context context, Handler handler, View view, String orderNum) {
        mContext = context;
        mHandler = handler;
        dialogView = view;//父布局view

        //loading  对话框里面的 取消订单 的 按钮
        mCancel = (TextView) dialogView.findViewById(R.id.trading_load_cancel);
        mLoadImage = (ImageView) dialogView.findViewById(R.id.trading_load_img);
        mLoadSecond = (TextView) dialogView.findViewById(R.id.trading_load_second);
        //这里默认显示：正在处理您的交易，请稍后
        mLoadMessage = (TextView) dialogView.findViewById(R.id.trading_load_message);

        initHandler();

        mOrderNum = orderNum;
    }


    /**
     * 交易时候用的 loading的对话框
     */
    public void loading() {
        dialogView.setVisibility(View.VISIBLE);//先让 对话框显示

        Animation loadingAnimation = AnimationUtils.loadAnimation(mContext, R.anim.loading_animation);
        mLoadImage.startAnimation(loadingAnimation);//给图片加上动画 转圈的动画

        //取消订单
        mCancel.setOnClickListener(new OnClickListener() {
            @Override
            public void onClick(View v) {
                OrderData orderData = new OrderData();
                orderData.origOrderNum = mOrderNum;
                CashierSdk.startCanc(orderData, new CashierListener() {
                    @Override
                    public void onResult(ResultData resultData) {
                        if (resultData.respcd.equals("00")) {
                            mHandler.sendEmptyMessage(Msg.MSG_FROM_SERVER_CLOSEBILL_SUCCESS);
                        } else if (resultData.respcd.equals("09")) {
                            mHandler.sendEmptyMessage(Msg.MSG_FROM_SERVER_CLOSEBILL_DOING);
                        } else {
                            mHandler.sendEmptyMessage(Msg.MSG_FROM_SERVER_CLOSEBILL_FAIL);
                        }
                    }

                    @Override
                    public void onError(int errorCode) {
                        mHandler.sendEmptyMessage(Msg.MSG_FROM_SERVER_TIMEOUT);
                    }
                });
            }
        });

        isLoading = true;
        second = 0;
        new Thread(new Runnable() {

            @Override
            public void run() {
                while (isLoading) {
                    try {
                        Thread.sleep(1000);
                        second++;
                        Message msg = dialogHandler.obtainMessage(Msg.MSG_FROM_DIGLOG_SECOND);
                        dialogHandler.sendMessageDelayed(msg, 0);
                    } catch (InterruptedException e) {
                        // TODO Auto-generated catch block
                        e.printStackTrace();
                    }

                }
            }
        }).start();

    }


    public void hide() {
        isLoading = false;
        dialogView.setVisibility(View.GONE);
    }

    public void stopLoading() {
        isLoading = false;
    }

    private void initHandler() {
        dialogHandler = new Handler() {
            @Override
            public void handleMessage(Message msg) {
                switch (msg.what) {
                    case Msg.MSG_FROM_DIGLOG_SECOND: {
                        mLoadSecond.setText(second + "S");
                        if (second > 45) {
                            mHandler.sendEmptyMessage(Msg.MSG_FROM_SERVER_TIMEOUT);
                        }
                    }
                }
                super.handleMessage(msg);
            }
        };
    }

}
