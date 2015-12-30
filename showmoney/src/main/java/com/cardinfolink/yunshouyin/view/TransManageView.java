package com.cardinfolink.yunshouyin.view;

import android.content.Context;
import android.os.Handler;
import android.text.TextUtils;
import android.util.Log;
import android.view.LayoutInflater;
import android.view.View;
import android.widget.ExpandableListView;
import android.widget.LinearLayout;
import android.widget.RadioButton;
import android.widget.RadioGroup;
import android.widget.TextView;

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
import com.cardinfolink.yunshouyin.ui.EditTextClear;
import com.cardinfolink.yunshouyin.util.ShowMoneyApp;
import com.handmark.pulltorefresh.library.PullToRefreshBase;
import com.handmark.pulltorefresh.library.PullToRefreshExpandableListView;

import java.text.SimpleDateFormat;
import java.util.ArrayList;
import java.util.Calendar;
import java.util.Collections;
import java.util.Comparator;
import java.util.Date;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Set;

public class TransManageView extends LinearLayout {
    private static final String TAG = "TransManageView";
    private Context mContext;

    //****************************************************************************************
    private PullToRefreshExpandableListView mBillPullRefreshListView;//第一个第一个账单的listview
    private BillExpandableListAdapter mBillAdapter;

    private Map<String, MonthBill> mMonthBillMap;
    private List<MonthBill> mMonthBilList;//月账单

    private List<List<TradeBill>> mTradeBillList;//日账单，这个条目会很多的
    private Map<String, List<TradeBill>> mTradeBillMap;

    private int billIndex;
    private int mMonthBillAgo;
    //****************************************************************************************

    //****************************************************************************************
    private PullToRefreshExpandableListView mTicketPullRefreshListView;//第2个第2个卡券账单的listview

    //****************************************************************************************
    private PullToRefreshExpandableListView mCollectionPullRefreshListView;//第3个第3个收款码账单的listview
    //****************************************************************************************

    private String mCurrentYearMonth;//当前年份+月份的一个字符串


    private TextView mTitle;
    private RadioButton mRaidoBill;//收款账单
    private RadioButton mRadioTicket;//卡券账单
    private RadioButton mRadioCollection;//收款码账单

    private RadioGroup mRadioGroup;

    private TextView mSearch;//搜索的按钮
    private EditTextClear mSearchEditText;
    private LinearLayout mSearchLinearLayout;
    private LinearLayout mSearchConditionLinearLayout;

    private Handler mMainactivityHandler;

    public TransManageView(Context context) {
        this(context, null);
    }

    public TransManageView(Context context, Handler handler) {
        super(context);
        mContext = context;
        mMainactivityHandler = handler;

        View contentView = LayoutInflater.from(context).inflate(R.layout.transmanage_view, null);
        LinearLayout.LayoutParams layoutParams = new LayoutParams(LayoutParams.MATCH_PARENT, LayoutParams.MATCH_PARENT);
        contentView.setLayoutParams(layoutParams);
        addView(contentView);

        initLayout();
    }

    private void initLayout() {
        SimpleDateFormat spf = new SimpleDateFormat("yyyyMM");
        mCurrentYearMonth = spf.format(new Date());

        mTitle = (TextView) findViewById(R.id.tv_title);

        //******************************************************************************************
        mBillPullRefreshListView = (PullToRefreshExpandableListView) findViewById(R.id.bill_list_view);
        mBillPullRefreshListView.setMode(PullToRefreshBase.Mode.BOTH);
        mBillPullRefreshListView.setOnRefreshListener(new PullToRefreshBase.OnRefreshListener<ExpandableListView>() {
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
                    mCurrentYearMonth = spf.format(new Date());
                    getBill();
                } else {
                    getBill();
                }

            }
        });

        mMonthBilList = new ArrayList<>();
        mTradeBillList = new ArrayList<>();
        mMonthBillMap = new HashMap<>();
        mTradeBillMap = new HashMap<>();

        mBillAdapter = new BillExpandableListAdapter(mContext, mMonthBilList, mTradeBillList);

        ExpandableListView ActualView = mBillPullRefreshListView.getRefreshableView();
        ActualView.setAdapter(mBillAdapter);
        ActualView.setGroupIndicator(null);

        //******************************************************************************************
        mTicketPullRefreshListView = (PullToRefreshExpandableListView) findViewById(R.id.ticket_list_view);
        mTicketPullRefreshListView.setVisibility(GONE);


        //******************************************************************************************
        mCollectionPullRefreshListView = (PullToRefreshExpandableListView) findViewById(R.id.colloction_list_view);
        mCollectionPullRefreshListView.setVisibility(GONE);


        //******************************************************************************************

        mRadioGroup = (RadioGroup) findViewById(R.id.redio_group);

        mSearchLinearLayout = (LinearLayout) findViewById(R.id.ll_search);
        mSearchConditionLinearLayout = (LinearLayout) findViewById(R.id.ll_search_condition);
        mSearch = (TextView) findViewById(R.id.tv_search);
        mSearchEditText = (EditTextClear) findViewById(R.id.et_search);

        mSearch.setOnClickListener(new OnClickListener() {
            @Override
            public void onClick(View v) {
                if (mRadioGroup.getVisibility() == VISIBLE) {
                    mRadioGroup.setVisibility(GONE);
                    mSearchLinearLayout.setVisibility(VISIBLE);
                    mSearchConditionLinearLayout.setVisibility(VISIBLE);
                } else if (mRadioGroup.getVisibility() == GONE) {
                    mRadioGroup.setVisibility(VISIBLE);
                    mSearchLinearLayout.setVisibility(GONE);
                    mSearchConditionLinearLayout.setVisibility(GONE);
                }
            }
        });

        //账单
        mRaidoBill = (RadioButton) findViewById(R.id.radio_bill);
        mRaidoBill.setOnClickListener(new OnClickListener() {

            @Override
            public void onClick(View v) {
                mBillPullRefreshListView.setVisibility(VISIBLE);
                mTicketPullRefreshListView.setVisibility(GONE);
                mCollectionPullRefreshListView.setVisibility(GONE);

                mTitle.setText(mRaidoBill.getText());//设置标题
            }
        });

        //卡券
        mRadioTicket = (RadioButton) findViewById(R.id.radio_ticket);
        mRadioTicket.setOnClickListener(new OnClickListener() {

            @Override
            public void onClick(View v) {
                mBillPullRefreshListView.setVisibility(GONE);
                mTicketPullRefreshListView.setVisibility(VISIBLE);
                mCollectionPullRefreshListView.setVisibility(GONE);
                mTitle.setText(mRadioTicket.getText());
            }
        });

        //收款码
        mRadioCollection = (RadioButton) findViewById(R.id.radio_collection);
        mRadioCollection.setOnClickListener(new OnClickListener() {

            @Override
            public void onClick(View v) {
                mBillPullRefreshListView.setVisibility(GONE);
                mTicketPullRefreshListView.setVisibility(GONE);
                mCollectionPullRefreshListView.setVisibility(VISIBLE);
                mTitle.setText(mRadioCollection.getText());
            }
        });
    }


    public void refresh() {
        getBill();
    }


    //获取收款的账单账单
    private void getBill() {
        QuickPayService quickPayService = ShowMoneyApp.getInstance().getQuickPayService();
        quickPayService.getHistoryBillsAsync(SessonData.loginUser, mCurrentYearMonth, String.valueOf(billIndex), "100", "all", new QuickPayCallbackListener<ServerPacket>() {
            @Override
            public void onSuccess(ServerPacket data) {
                //这里可以在ui线程里执行的
                final int totalRecord = data.getTotalRecord();//这个字段表示当月的总条数
                final int count = data.getCount();//返回的条数
                final String total = data.getTotal();
                final int refdcount = data.getRefdcount();
                final String refdtotal = data.getRefdtotal();
                final int size = data.getSize();

                //这里开始遍历这个账单的数组************************************************************
                for (Txn txn : data.getTxn()) {
                    TradeBill tradeBill = new TradeBill();
                    tradeBill.response = txn.getResponse();
                    tradeBill.tandeDate = txn.getSystemDate();
                    tradeBill.consumerAccount = txn.getConsumerAccount();
                    tradeBill.refundAmt = TxamtUtil.getNormal(txn.getRefundAmt());

                    QRequest req = txn.getmRequest();
                    if (req != null) {
                        tradeBill.orderNum = req.getOrderNum();
                        tradeBill.amount = TxamtUtil.getNormal(req.getTxamt());
                        tradeBill.busicd = req.getBusicd();

                        //使用/v3/bill接口 退款的好像也没有拉取到
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

                    //渠道为空的 不列入统计，这样totalRecord和实际的list的size可能不一样
                    if (TextUtils.isEmpty(tradeBill.chcd)) {
                        continue;
                    }

                    //添加到相应的map中，最后再转换到list中，按照月份的先后排序转换到list中
                    if (mMonthBillMap.containsKey(currentYearMonth)) {
                        mMonthBillMap.get(currentYearMonth).setCount(count);
                        mMonthBillMap.get(currentYearMonth).setTotal(total);
                        mMonthBillMap.get(currentYearMonth).setRefdcount(refdcount);
                        mMonthBillMap.get(currentYearMonth).setRefdtotal(refdtotal);
                        mMonthBillMap.get(currentYearMonth).setSize(size);
                        mMonthBillMap.get(currentYearMonth).setTotalRecord(totalRecord);
                    } else {
                        MonthBill monthBill = new MonthBill(currentYear, currentMonth);
                        monthBill.setCount(count);
                        monthBill.setTotal(total);
                        monthBill.setRefdcount(refdcount);
                        monthBill.setRefdtotal(refdtotal);
                        monthBill.setSize(size);
                        monthBill.setTotalRecord(totalRecord);
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
                //**********************************************************************************

                mapToMonthBillList(mMonthBillMap);
                mapToTradeBillList(mTradeBillMap);

                mBillAdapter.notifyDataSetChanged();
                mBillPullRefreshListView.onRefreshComplete();

                billIndex += size;
                if (billIndex == totalRecord) {
                    //之前用的是size来判断的。size等于零 表示 加载到这个月的全部的了，这时候就要加载前一个月的数据了
                    //现在用totalRecord来判断，相等表明这个月的数据加载完了，这个时候就要加载前一个月的数据了
                    billIndex = 0;
                    mMonthBillAgo += 1;
                    Calendar calendar = Calendar.getInstance();
                    calendar.add(Calendar.MONTH, 0 - mMonthBillAgo);    //得到前一个月
                    String year = String.valueOf(calendar.get(Calendar.YEAR));
                    String month = String.valueOf(calendar.get(Calendar.MONTH) + 1);
                    mCurrentYearMonth = year + month;
                }
            }

            @Override
            public void onFailure(QuickPayException ex) {
                Log.e(TAG, " get history bill fail" + ex.getErrorMsg());
                mBillPullRefreshListView.onRefreshComplete();
            }
        });
    }

    //获取卡券账单
    public void getTicketBill() {

    }

    //获取 收款码 账单
    public void getCollectionBill() {

    }

    private void mapToMonthBillList(Map<String, MonthBill> map) {
        Set<String> keyset = map.keySet();
        ArrayList<String> list = new ArrayList<>();
        list.addAll(keyset);
        Comparator<String> com = new Comparator<String>() {

            @Override
            public int compare(String lhs, String rhs) {
                if (Integer.valueOf(lhs) > Integer.valueOf(rhs)) {
                    return -1;
                } else if (Integer.valueOf(lhs) < Integer.valueOf(rhs)) {
                    return 1;
                } else {
                    return 0;
                }
            }
        };

        Collections.sort(list, com);

        mMonthBilList.clear();
        for (String key : list) {
            mMonthBilList.add(map.get(key));
        }
    }

    private void mapToTradeBillList(Map<String, List<TradeBill>> map) {
        Set<String> keyset = map.keySet();
        ArrayList<String> list = new ArrayList<>();
        list.addAll(keyset);
        Comparator<String> com = new Comparator<String>() {

            @Override
            public int compare(String lhs, String rhs) {
                if (Integer.valueOf(lhs) > Integer.valueOf(rhs)) {
                    return -1;
                } else if (Integer.valueOf(lhs) < Integer.valueOf(rhs)) {
                    return 1;
                } else {
                    return 0;
                }
            }
        };
        Collections.sort(list, com);//把 key 排序一下

        mTradeBillList.clear();
        for (String key : list) {
            mTradeBillList.add(map.get(key));
        }
    }
}
