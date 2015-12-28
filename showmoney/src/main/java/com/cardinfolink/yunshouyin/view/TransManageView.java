package com.cardinfolink.yunshouyin.view;

import android.content.Context;
import android.text.TextUtils;
import android.util.Log;
import android.view.LayoutInflater;
import android.view.View;
import android.widget.ExpandableListView;
import android.widget.LinearLayout;

import com.cardinfolink.cashiersdk.util.TxamtUtil;
import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.adapter.BillExpandableListAdapter;
import com.cardinfolink.yunshouyin.api.QuickPayException;
import com.cardinfolink.yunshouyin.core.QuickPayCallbackListener;
import com.cardinfolink.yunshouyin.core.QuickPayService;
import com.cardinfolink.yunshouyin.data.MonthBill;
import com.cardinfolink.yunshouyin.data.SessonData;
import com.cardinfolink.yunshouyin.data.TradeBill;
import com.cardinfolink.yunshouyin.model.QRequest;
import com.cardinfolink.yunshouyin.model.ServerPacket;
import com.cardinfolink.yunshouyin.model.Txn;
import com.cardinfolink.yunshouyin.util.ShowMoneyApp;
import com.handmark.pulltorefresh.library.PullToRefreshBase;
import com.handmark.pulltorefresh.library.PullToRefreshExpandableListView;

import java.text.SimpleDateFormat;
import java.util.ArrayList;
import java.util.Calendar;
import java.util.Date;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class TransManageView extends LinearLayout {
    private static final String TAG = "TransManageView";
    private Context mContext;
    private PullToRefreshExpandableListView mPullRefreshListView;
    private BillExpandableListAdapter mAdapter;

    private List<Map<String, List<TradeBill>>> mTradeBillList;//日账单，这个条目会很多的
    private List<MonthBill> mMonthList;//月账单
    private int monthAgo = 0;


    public TransManageView(Context context) {
        super(context);
        mContext = context;
        View contentView = LayoutInflater.from(context).inflate(R.layout.transmanage_view, null);
        LinearLayout.LayoutParams layoutParams = new LayoutParams(LayoutParams.MATCH_PARENT, LayoutParams.MATCH_PARENT);
        contentView.setLayoutParams(layoutParams);
        addView(contentView);

        initLayout();

    }

    private void initLayout() {
        mPullRefreshListView = (PullToRefreshExpandableListView) findViewById(R.id.bill_list_view);

        mPullRefreshListView.setMode(PullToRefreshBase.Mode.BOTH);


        mPullRefreshListView.setOnRefreshListener(new PullToRefreshBase.OnRefreshListener<ExpandableListView>() {
            @Override
            public void onRefresh(PullToRefreshBase<ExpandableListView> refreshView) {
                // Do work to refresh the list here.
                Log.e(TAG, " pull to refersh");
                mPullRefreshListView.onRefreshComplete();

            }
        });

        mMonthList = new ArrayList<>();

        mTradeBillList = new ArrayList<>();

        mAdapter = new BillExpandableListAdapter(mContext, mMonthList, mTradeBillList);

        final ExpandableListView ActualView = mPullRefreshListView.getRefreshableView();
        ActualView.setAdapter(mAdapter);
        ActualView.setGroupIndicator(null);
    }


    public void refresh() {
        getTradeBill();
    }


    private void getTradeBill() {
        Calendar calendar = Calendar.getInstance();
        calendar.add(Calendar.MONTH, 0 - monthAgo);    //得到前一个月
        String year = String.valueOf(calendar.get(Calendar.YEAR));
        String month = String.valueOf(calendar.get(Calendar.MONTH) + 1);


        final MonthBill monthBill = new MonthBill(year, month);

        QuickPayService quickPayService = ShowMoneyApp.getInstance().getQuickPayService();
        quickPayService.getHistoryBillsAsync(SessonData.loginUser, year + month, "0", "all", new QuickPayCallbackListener<ServerPacket>() {
            @Override
            public void onSuccess(ServerPacket data) {
                //这里可以在ui线程里执行的

                final int count = data.getCount();//返回的条数
                final String total = data.getTotal();
                final int refdcount = data.getRefdcount();
                final String refdtotal = data.getRefdtotal();
                final int size = data.getSize();
                Map<String, List<TradeBill>> billMap = new HashMap<String, List<TradeBill>>();
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

                    //获取这个账单里面的日期,年月日 的 日
                    String currentDay = tradeBill.tandeDate.substring(6, 8);
                    Log.e(TAG, tradeBill.tandeDate + " === " + currentDay);

                    if (!TextUtils.isEmpty(tradeBill.chcd)) {
                        if (billMap.containsKey(currentDay)) {
                            billMap.get(currentDay).add(tradeBill);
                        } else {
                            ArrayList<TradeBill> list = new ArrayList<TradeBill>();
                            list.add(tradeBill);
                            billMap.put(currentDay, list);//把这个list添加到map里面
                        }
                    }
                }
                // 更新UI
                monthBill.setCount(count);
                monthBill.setTotal(total);
                monthBill.setRefdcount(refdcount);
                monthBill.setRefdtotal(refdtotal);
                monthBill.setSize(size);

                mMonthList.add(monthBill);
                mTradeBillList.add(billMap);
                mAdapter.notifyDataSetChanged();
                mPullRefreshListView.onRefreshComplete();

                monthAgo++;//月份往前计数一个月
            }

            @Override
            public void onFailure(QuickPayException ex) {
                Log.e(TAG, " get history bill fail"+ex.getErrorMsg());
                mPullRefreshListView.onRefreshComplete();
            }
        });
    }

}
