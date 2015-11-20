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

public class TradingCustomDialog {

    private Context mContext;
    private Handler mHandler;
    private Handler dialogHandler;
    private boolean isLoading = false;
    private int second = 0;
    private View mLoadView;
    private View mSuccessView;
    private View mFailView;
    private View mNoPayView;
    private View dialogView;
    private TextView mSecondText;
    private String mOrderNum;

    public TradingCustomDialog(Context context, Handler handler, View view, String orderNum) {
        mContext = context;
        mHandler = handler;
        dialogView = view;
        initLayout();
        initHandler();
        initListener();
        mOrderNum = orderNum;

    }


    private void initLayout() {
        mLoadView = dialogView.findViewById(R.id.trading_custom_dialog_loading);
        mSuccessView = dialogView.findViewById(R.id.trading_custom_dialog_success);
        mFailView = dialogView.findViewById(R.id.trading_custom_dialog_fail);
        mNoPayView = dialogView.findViewById(R.id.trading_custom_dialog_nopay);
    }

    private void initListener() {
        mLoadView.findViewById(R.id.close).setOnClickListener(new OnClickListener() {

            @Override
            public void onClick(View v) {

                mHandler.sendEmptyMessage(Msg.MSG_FROM_DIGLOG_CLOSE);
            }
        });

        mLoadView.findViewById(R.id.query).setOnClickListener(new OnClickListener() {

            @Override
            public void onClick(View v) {
                OrderData orderData = new OrderData();
                orderData.origOrderNum = mOrderNum;
                CashierSdk.startQy(orderData, new CashierListener() {

                    @Override
                    public void onResult(ResultData resultData) {
                        if (resultData.respcd.equals("00")) {
                            mHandler.sendEmptyMessage(Msg.MSG_FROM_SERVER_TRADE_SUCCESS);
                        } else if (resultData.respcd.equals("09")) {
                            mHandler.sendEmptyMessage(Msg.MSG_FROM_SERVER_TRADE_NOPAY);
                        } else {
                            mHandler.sendEmptyMessage(Msg.MSG_FROM_SERVER_TRADE_FAIL);
                        }
                    }

                    @Override
                    public void onError(int errorCode) {
                        mHandler.sendEmptyMessage(Msg.MSG_FROM_SERVER_TIMEOUT);
                    }

                });
            }
        });

        mNoPayView.findViewById(R.id.nopay_query).setOnClickListener(new OnClickListener() {

            @Override
            public void onClick(View v) {
                OrderData orderData = new OrderData();
                orderData.origOrderNum = mOrderNum;
                CashierSdk.startQy(orderData, new CashierListener() {

                    @Override
                    public void onResult(ResultData resultData) {

                        if (resultData.respcd.equals("00")) {
                            mHandler.sendEmptyMessage(Msg.MSG_FROM_SERVER_TRADE_SUCCESS);
                        } else if (resultData.respcd.equals("09")) {
                            mHandler.sendEmptyMessage(Msg.MSG_FROM_SERVER_TRADE_NOPAY);
                        } else {
                            mHandler.sendEmptyMessage(Msg.MSG_FROM_SERVER_TRADE_FAIL);
                        }
                    }

                    @Override
                    public void onError(int errorCode) {
                        mHandler.sendEmptyMessage(Msg.MSG_FROM_SERVER_TIMEOUT);
                    }

                });
            }
        });


        mNoPayView.findViewById(R.id.nopay_close).setOnClickListener(new OnClickListener() {

            @Override
            public void onClick(View v) {
                dialogView.setVisibility(View.GONE);
                //mHandler.sendEmptyMessage(Msg.MSG_FROM_DIGLOG_CLOSE);
            }
        });

        mSuccessView.findViewById(R.id.success_dialog_back).setOnClickListener(new OnClickListener() {

            @Override
            public void onClick(View v) {

                mHandler.sendEmptyMessage(Msg.MSG_FROM_DIGLOG_CLOSE);
            }
        });


        mSuccessView.findViewById(R.id.histroy).setOnClickListener(new OnClickListener() {

            @Override
            public void onClick(View v) {
                mHandler.sendEmptyMessage(Msg.MSG_FROM_SUCCESS_DIGLOG_HISTORY);
            }
        });


        mFailView.findViewById(R.id.fail_dialog_back).setOnClickListener(new OnClickListener() {

            @Override
            public void onClick(View v) {
                mHandler.sendEmptyMessage(Msg.MSG_FROM_DIGLOG_CLOSE);
            }
        });

        mFailView.findViewById(R.id.fail_query).setOnClickListener(new OnClickListener() {

            @Override
            public void onClick(View v) {
                mHandler.sendEmptyMessage(Msg.MSG_FROM_SUCCESS_DIGLOG_HISTORY);
            }
        });
    }

    public void nopay() {
        isLoading = false;
        dialogView.setVisibility(View.VISIBLE);
        mLoadView.setVisibility(View.GONE);
        mNoPayView.setVisibility(View.VISIBLE);
    }


    public void loading() {
        dialogView.setVisibility(View.VISIBLE);
        mLoadView.setVisibility(View.VISIBLE);
        Animation loadingAnimation = AnimationUtils.loadAnimation(
                mContext, R.anim.loading_animation);
        ImageView loadView = (ImageView) mLoadView.findViewById(R.id.trading_custom_dialog_load_img);
        mSecondText = (TextView) mLoadView.findViewById(R.id.second);
        loadView.startAnimation(loadingAnimation);

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


    private void initHandler() {
        dialogHandler = new Handler() {
            @Override
            public void handleMessage(Message msg) {
                switch (msg.what) {
                    case Msg.MSG_FROM_DIGLOG_SECOND: {
                        mSecondText.setText(second + "S");

                    }
                }
                super.handleMessage(msg);
            }
        };
    }

    public void success() {
        isLoading = false;
        dialogView.setVisibility(View.VISIBLE);
        mLoadView.setVisibility(View.GONE);
        mSuccessView.setVisibility(View.VISIBLE);
    }

    public void fail() {
        isLoading = false;
        dialogView.setVisibility(View.VISIBLE);
        mLoadView.setVisibility(View.GONE);
        mFailView.setVisibility(View.VISIBLE);
    }
}
