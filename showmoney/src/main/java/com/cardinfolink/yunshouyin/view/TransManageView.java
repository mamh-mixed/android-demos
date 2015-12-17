package com.cardinfolink.yunshouyin.view;

import android.content.Context;
import android.text.TextUtils;
import android.view.LayoutInflater;
import android.view.MotionEvent;
import android.view.View;
import android.widget.LinearLayout;
import android.widget.ListView;
import android.widget.TextView;

import com.cardinfolink.cashiersdk.util.TxamtUtil;
import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.adapter.BillAdapter;
import com.cardinfolink.yunshouyin.api.QuickPayException;
import com.cardinfolink.yunshouyin.core.QuickPayCallbackListener;
import com.cardinfolink.yunshouyin.core.QuickPayService;
import com.cardinfolink.yunshouyin.data.SessonData;
import com.cardinfolink.yunshouyin.data.TradeBill;
import com.cardinfolink.yunshouyin.model.QRequest;
import com.cardinfolink.yunshouyin.model.ServerPacket;
import com.cardinfolink.yunshouyin.model.Txn;
import com.cardinfolink.yunshouyin.util.ShowMoneyApp;
import com.handmark.pulltorefresh.library.PullToRefreshBase;
import com.handmark.pulltorefresh.library.PullToRefreshBase.Mode;
import com.handmark.pulltorefresh.library.PullToRefreshBase.OnRefreshListener;
import com.handmark.pulltorefresh.library.PullToRefreshListView;

import java.text.SimpleDateFormat;
import java.util.ArrayList;
import java.util.Date;
import java.util.List;

public class TransManageView extends LinearLayout {
    private static final String TAG = "TransManageView";
    private Context mContext;

    private PullToRefreshListView mPullToRefreshListView;
    private List<TradeBill> mTradeBillList;
    private TextView mBillTipsText;

    private String mMonth;
    private String tipsYearMonth;
    private int billIndex;
    private String mBillStatus;
    private BillAdapter mBillAdapter;


    public TransManageView(Context context) {
        super(context);
        mContext = context;
        View contentView = LayoutInflater.from(context).inflate(R.layout.transmanage_view, null);
        LinearLayout.LayoutParams layoutParams = new LayoutParams(LayoutParams.MATCH_PARENT, LayoutParams.MATCH_PARENT);
        contentView.setLayoutParams(layoutParams);
        addView(contentView);
        SimpleDateFormat spf = new SimpleDateFormat("yyyyMM");
        mMonth = spf.format(new Date());
        tipsYearMonth = (new SimpleDateFormat("yyyy" + getResources().getString(R.string.year) + "MM" + getResources().getString(R.string.month))).format(new Date());
        mTradeBillList = new ArrayList<TradeBill>();
        initLayout();
        initListener();
        billIndex = 0;
        mBillStatus = "all";
    }

    public void initData() {
        billIndex = 0;
        mTradeBillList.clear();
        getTradeBill();
    }

    public void refresh() {
        billIndex = 0;
        mTradeBillList.clear();
        new Thread(new Runnable() {

            @Override
            public void run() {
                try {
                    Thread.sleep(500);
                } catch (InterruptedException e) {
                    e.printStackTrace();
                }
                getTradeBill();
            }
        }).start();


    }


    private void initLayout() {
        mBillTipsText = (TextView) findViewById(R.id.bill_tips);
        mPullToRefreshListView = (PullToRefreshListView) findViewById(R.id.pull_refresh_list);
        //mListView = mPullToRefreshListView.getRefreshableView();
        mBillAdapter = new BillAdapter(this, mTradeBillList);
        mPullToRefreshListView.setAdapter(mBillAdapter);
        // 设置pull-to-refresh模式为Mode.Both
        mPullToRefreshListView.setMode(Mode.BOTH);
        mPullToRefreshListView.setOnTouchListener(new View.OnTouchListener() {

            @Override
            public boolean onTouch(View v, MotionEvent event) {
                return false;
            }
        });

    }

    private void initListener() {
        // 设置上拉下拉事件
        mPullToRefreshListView.setOnRefreshListener(new OnRefreshListener<ListView>() {

            @Override
            public void onRefresh(PullToRefreshBase<ListView> refreshView) {
                if (refreshView.isHeaderShown()) {
                    billIndex = 0;
                    mTradeBillList.clear();
                    getTradeBill();
                } else {
                    getTradeBill();
                }
            }
        });

        findViewById(R.id.radio_all).setOnClickListener(new OnClickListener() {

            @Override
            public void onClick(View v) {
                billIndex = 0;
                mBillStatus = "all";
                mTradeBillList.clear();
                mPullToRefreshListView.setRefreshing();
            }
        });

        findViewById(R.id.radio_success).setOnClickListener(new OnClickListener() {

            @Override
            public void onClick(View v) {
                billIndex = 0;
                mBillStatus = "success";
                mTradeBillList.clear();
                mPullToRefreshListView.setRefreshing();
            }
        });

        findViewById(R.id.radio_fail).setOnClickListener(new OnClickListener() {

            @Override
            public void onClick(View v) {
                billIndex = 0;
                mBillStatus = "fail";
                mTradeBillList.clear();
                mPullToRefreshListView.setRefreshing();
            }
        });
    }

    private void getTradeBill() {
        QuickPayService quickPayService = ShowMoneyApp.getInstance().getQuickPayService();
        quickPayService.getHistoryBillsAsync(SessonData.loginUser, mMonth, billIndex + "", mBillStatus, new QuickPayCallbackListener<ServerPacket>() {
            @Override
            public void onSuccess(ServerPacket data) {
                final int count = data.getCount();
                final String total = data.getTotal();
                final int refdcount = data.getRefdcount();
                final String refdtotal = data.getRefdtotal();
                final int size = data.getSize();

                for (Txn txn : data.getTxn()) {
                    TradeBill tradeBill = new TradeBill();
                    tradeBill.response = txn.getResponse();
                    tradeBill.tandeDate = txn.getSystemDate();
                    tradeBill.consumerAccount = txn.getConsumerAccount();

                    QRequest req = txn.getmRequest();
                    if (req != null) {
                        tradeBill.orderNum = req.getOrderNum();
                        tradeBill.amount = TxamtUtil.getNormal(req.getTxamt());
                        tradeBill.busicd = req.getBusicd();
                        if (tradeBill.busicd.equals("REFD")) {
                            tradeBill.amount = "-" + tradeBill.amount;
                        }
                        tradeBill.chcd = req.getChcd();
                        tradeBill.tradeFrom = req.getTradeFrom();
                        tradeBill.goodsInfo = req.getGoodsInfo();
                    }
                    if (!TextUtils.isEmpty(tradeBill.chcd)) {
                        mTradeBillList.add(tradeBill);
                    }
                }
                // 更新UI
                mBillAdapter.setData(mTradeBillList);
                mBillAdapter.notifyDataSetChanged();
                mPullToRefreshListView.onRefreshComplete();


                mBillTipsText.setText(
                        tipsYearMonth +
                                "  " +
                                getResources().getString(R.string.txn_total_times) + count +
                                " " +
                                getResources().getString(R.string.txn_total_amount) + total + getResources().getString(R.string.txn_currency) +
                                " " +
                                getResources().getString(R.string.txn_refund) + refdcount
                                +
                                getResources().getString(R.string.txn_unit) + "(" + refdtotal + getResources().getString(R.string.txn_currency) + ")");
                billIndex += size;
            }

            @Override
            public void onFailure(QuickPayException ex) {
                mBillAdapter.notifyDataSetChanged();
                mPullToRefreshListView.onRefreshComplete();
            }
        });
    }


}
