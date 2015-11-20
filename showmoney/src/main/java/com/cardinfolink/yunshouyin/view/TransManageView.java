package com.cardinfolink.yunshouyin.view;

import android.annotation.SuppressLint;
import android.app.Activity;
import android.content.Context;
import android.view.LayoutInflater;
import android.view.MotionEvent;
import android.view.View;
import android.widget.LinearLayout;
import android.widget.ListView;
import android.widget.TextView;

import com.cardinfolink.cashiersdk.util.TxamtUtil;
import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.adapter.BillAdapter;
import com.cardinfolink.yunshouyin.data.SessonData;
import com.cardinfolink.yunshouyin.data.TradeBill;
import com.cardinfolink.yunshouyin.util.CommunicationListener;
import com.cardinfolink.yunshouyin.util.HttpCommunicationUtil;
import com.cardinfolink.yunshouyin.util.JsonUtil;
import com.cardinfolink.yunshouyin.util.ParamsUtil;
import com.handmark.pulltorefresh.library.PullToRefreshBase;
import com.handmark.pulltorefresh.library.PullToRefreshBase.Mode;
import com.handmark.pulltorefresh.library.PullToRefreshBase.OnRefreshListener;
import com.handmark.pulltorefresh.library.PullToRefreshListView;

import org.json.JSONArray;
import org.json.JSONException;

import java.text.SimpleDateFormat;
import java.util.ArrayList;
import java.util.Date;
import java.util.List;

public class TransManageView extends LinearLayout {
    private Context mContext;

    private PullToRefreshListView mPullToRefreshListView;
    //private ListView mListView;
    private List<TradeBill> mTradeBillList;
    private TextView mBillTipsText;

    private String mMonth;
    private String tips_year_month;
    private int bill_index;
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
        tips_year_month = (new SimpleDateFormat("yyyy" + getResources().getString(R.string.year) + "MM" + getResources().getString(R.string.month))).format(new Date());
        mTradeBillList = new ArrayList<TradeBill>();
        initLayout();
        initListener();
        bill_index = 0;
        mBillStatus = "all";
    }

    public void initData() {
        bill_index = 0;
        mTradeBillList.clear();
        getTradeBill();
    }

    public void refresh() {
        bill_index = 0;
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
                    bill_index = 0;
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
                bill_index = 0;
                mBillStatus = "all";
                mTradeBillList.clear();
                mPullToRefreshListView.setRefreshing();
            }
        });

        findViewById(R.id.radio_success).setOnClickListener(new OnClickListener() {

            @Override
            public void onClick(View v) {
                bill_index = 0;
                mBillStatus = "success";
                mTradeBillList.clear();
                mPullToRefreshListView.setRefreshing();
            }
        });

        findViewById(R.id.radio_fail).setOnClickListener(new OnClickListener() {

            @Override
            public void onClick(View v) {
                bill_index = 0;
                mBillStatus = "fail";
                mTradeBillList.clear();
                mPullToRefreshListView.setRefreshing();
            }
        });
    }

    private void getTradeBill() {

        HttpCommunicationUtil.sendDataToServer(ParamsUtil.getHistory(SessonData.loginUser, mMonth, bill_index, mBillStatus), new CommunicationListener() {

            @Override
            public void onResult(String result) {
                String state = JsonUtil.getParam(result, "state");
                if (state.equals("success")) {
                    final String count = JsonUtil.getParam(result, "count");
                    final String total = JsonUtil.getParam(result, "total");
                    final String refdcount = JsonUtil.getParam(result, "refdcount");
                    final String refdtotal = JsonUtil.getParam(result, "refdtotal");
                    String txn = JsonUtil.getParam(result, "txn");
                    final int size = Integer.parseInt(JsonUtil.getParam(result, "size"));
                    try {
                        JSONArray txnAJsonArray = new JSONArray(txn);
                        for (int i = 0; i < txnAJsonArray.length(); i++) {
                            String bill = txnAJsonArray.getString(i);
                            String m_request = JsonUtil.getParam(bill, "m_request");
                            TradeBill tradeBill = new TradeBill();
                            tradeBill.orderNum = JsonUtil.getParam(m_request, "orderNum");
                            tradeBill.amount = TxamtUtil.getNormal(JsonUtil.getParam(m_request, "txamt"));
                            tradeBill.busicd = JsonUtil.getParam(m_request, "busicd");
                            tradeBill.chcd = JsonUtil.getParam(m_request, "chcd");
                            tradeBill.response = JsonUtil.getParam(bill, "response");
                            tradeBill.tandeDate = JsonUtil.getParam(bill, "system_date");
                            tradeBill.consumerAccount = JsonUtil.getParam(bill, "consumerAccount");
                            tradeBill.tradeFrom = JsonUtil.getParam(m_request, "tradeFrom");
                            tradeBill.goodsInfo = JsonUtil.getParam(m_request, "goodsInfo");
                            mTradeBillList.add(tradeBill);
                        }

                        ((Activity) mContext).runOnUiThread(new Runnable() {

                            @SuppressLint("NewApi")
                            @Override
                            public void run() {
                                // 更新UI
                                mBillAdapter.setData(mTradeBillList);
                                mBillAdapter.notifyDataSetChanged();
                                mPullToRefreshListView.onRefreshComplete();


                                mBillTipsText.setText(
                                        tips_year_month +
                                                "  " +
                                                getResources().getString(R.string.txn_total_times) + count +
                                                " " +
                                                getResources().getString(R.string.txn_total_amount) + total + getResources().getString(R.string.txn_currency) +
                                                " " +
                                                getResources().getString(R.string.txn_refund) + refdcount
                                                +
                                                getResources().getString(R.string.txn_unit) + "(" + refdtotal + getResources().getString(R.string.txn_currency) + ")");
                                bill_index += size;
                            }

                        });

                    } catch (JSONException e) {
                        ((Activity) mContext).runOnUiThread(new Runnable() {

                            @SuppressLint("NewApi")
                            @Override
                            public void run() {
                                // 更新UI
                                mBillAdapter.notifyDataSetChanged();
                                mPullToRefreshListView.onRefreshComplete();

                            }

                        });
                        e.printStackTrace();
                    }
                    // bill_tips
                }
            }

            @Override
            public void onError(String error) {

            }
        });
    }


}
