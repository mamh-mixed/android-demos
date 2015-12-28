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
import java.util.Collections;
import java.util.Date;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Set;

public class TransManageView extends LinearLayout {
    private static final String TAG = "TransManageView";
    private Context mContext;
    private PullToRefreshExpandableListView mPullRefreshListView;
    private BillExpandableListAdapter mAdapter;

    private Map<String, MonthBill> mMonthBillMap;
    private List<MonthBill> mMonthBilList;//月账单

    private List<List<TradeBill>> mTradeBillList;//日账单，这个条目会很多的
    private Map<String, List<TradeBill>> mTradeBillMap;

    private int billIndex;
    private String mMonth;
    private int mMonthAgo;

    public TransManageView(Context context) {
        super(context);
        mContext = context;
        View contentView = LayoutInflater.from(context).inflate(R.layout.transmanage_view, null);
        LinearLayout.LayoutParams layoutParams = new LayoutParams(LayoutParams.MATCH_PARENT, LayoutParams.MATCH_PARENT);
        contentView.setLayoutParams(layoutParams);
        addView(contentView);
        initLayout();


        SimpleDateFormat spf = new SimpleDateFormat("yyyyMM");
        mMonth = spf.format(new Date());
    }

    private void initLayout() {
        mPullRefreshListView = (PullToRefreshExpandableListView) findViewById(R.id.bill_list_view);

        mPullRefreshListView.setMode(PullToRefreshBase.Mode.BOTH);


        mPullRefreshListView.setOnRefreshListener(new PullToRefreshBase.OnRefreshListener<ExpandableListView>() {
            @Override
            public void onRefresh(PullToRefreshBase<ExpandableListView> refreshView) {
                // Do work to refresh the list here.
                if (refreshView.isHeaderShown()) {
                    billIndex = 0;
                    mTradeBillList.clear();
                    mMonthBilList.clear();
                    mMonthBillMap.clear();
                    mTradeBillMap.clear();
                    SimpleDateFormat spf = new SimpleDateFormat("yyyyMM");
                    mMonth = spf.format(new Date());
                    getTradeBill();
                    Log.e(TAG, " billIndex == " + billIndex);
                } else {
                    getTradeBill();
                }

            }
        });

        mMonthBilList = new ArrayList<>();
        mTradeBillList = new ArrayList<>();
        mMonthBillMap = new HashMap<>();
        mTradeBillMap = new HashMap<>();


        mAdapter = new BillExpandableListAdapter(mContext, mMonthBilList, mTradeBillList);

        final ExpandableListView ActualView = mPullRefreshListView.getRefreshableView();
        ActualView.setAdapter(mAdapter);
        ActualView.setGroupIndicator(null);
    }


    public void refresh() {
        getTradeBill();
    }


    private void getTradeBill() {
        QuickPayService quickPayService = ShowMoneyApp.getInstance().getQuickPayService();
        quickPayService.getHistoryBillsAsync(SessonData.loginUser, mMonth, String.valueOf(billIndex), "100", "all", new QuickPayCallbackListener<ServerPacket>() {
            @Override
            public void onSuccess(ServerPacket data) {
                //这里可以在ui线程里执行的

                final int count = data.getCount();//返回的条数
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

                    //获取这个账单里面的日期,年月日 的 日
                    String currentDay = tradeBill.tandeDate.substring(6, 8);
                    String currentYear = tradeBill.tandeDate.substring(0, 4);
                    String currentMonth = tradeBill.tandeDate.substring(4, 6);
                    String currentYearMonth = tradeBill.tandeDate.substring(0, 6);

                    if (TextUtils.isEmpty(tradeBill.chcd)) {
                        continue;
                    }

                    //添加到相应的map中，最后在转换到list中，排序转换到list中
                    if (mMonthBillMap.containsKey(currentYearMonth)) {
                        mMonthBillMap.get(currentYearMonth).setCount(count);
                        mMonthBillMap.get(currentYearMonth).setTotal(total);
                        mMonthBillMap.get(currentYearMonth).setRefdcount(refdcount);
                        mMonthBillMap.get(currentYearMonth).setRefdtotal(refdtotal);
                        mMonthBillMap.get(currentYearMonth).setSize(billIndex);
                    } else {
                        MonthBill monthBill = new MonthBill(currentYear, currentMonth);
                        monthBill.setCount(count);
                        monthBill.setTotal(total);
                        monthBill.setRefdcount(refdcount);
                        monthBill.setRefdtotal(refdtotal);
                        monthBill.setSize(billIndex);
                        mMonthBillMap.put(currentYearMonth, monthBill);
                    }
                    if (mTradeBillMap.containsKey(currentYearMonth)) {
                        mTradeBillMap.get(currentYearMonth).add(tradeBill);
                    } else {
                        List<TradeBill> list = new ArrayList<TradeBill>();
                        list.add(tradeBill);
                        mTradeBillMap.put(currentYearMonth, list);
                    }
                }

                mapToMonthBillList(mMonthBillMap);
                mapToTradeBillList(mTradeBillMap);

                mAdapter.notifyDataSetChanged();
                mPullRefreshListView.onRefreshComplete();

                billIndex += size;
                if (size == 0) {
                    //size等于零 表示 加载到这个月的全部的了，这时候就要加载前一个月的数据了
                    billIndex = 0;
                    mMonthAgo += 1;
                    Calendar calendar = Calendar.getInstance();
                    calendar.add(Calendar.MONTH, 0 - mMonthAgo);    //得到前一个月
                    String year = String.valueOf(calendar.get(Calendar.YEAR));
                    String month = String.valueOf(calendar.get(Calendar.MONTH) + 1);
                    mMonth = year + month;
                }
            }

            @Override
            public void onFailure(QuickPayException ex) {
                mPullRefreshListView.onRefreshComplete();
            }
        });
    }

    private void mapToMonthBillList(Map<String, MonthBill> map) {
        Set<String> keyset = map.keySet();
        ArrayList<String> list = new ArrayList<>();
        list.addAll(keyset);
        Collections.reverse(list);

        mMonthBilList.clear();
        for (String key : list) {
            mMonthBilList.add(map.get(key));
        }
    }

    private void mapToTradeBillList(Map<String, List<TradeBill>> map) {
        Set<String> keyset = map.keySet();
        ArrayList<String> list = new ArrayList<>();
        list.addAll(keyset);
        Collections.reverse(list);//把 key 排序一下

        mTradeBillList.clear();
        for (String key : list) {
            mTradeBillList.add(map.get(key));
        }
    }
}
