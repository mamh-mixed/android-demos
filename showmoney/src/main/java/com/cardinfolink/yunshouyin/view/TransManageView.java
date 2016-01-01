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
import com.cardinfolink.yunshouyin.adapter.CollectionExpandableListAdapter;
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

    private QuickPayService quickPayService;

    //***普通的收款账单*************************************************************************************
    private PullToRefreshExpandableListView mBillPullRefreshListView;//第一个第一个账单的listview
    private BillExpandableListAdapter mBillAdapter;

    private Map<String, MonthBill> mMonthBillMap;
    private List<MonthBill> mMonthBilList;//月账单

    private List<List<TradeBill>> mTradeBillList;//日账单，这个条目会很多的
    private Map<String, List<TradeBill>> mTradeBillMap;

    private int billIndex;
    private int mMonthBillAgo;

    //***卡券账单*************************************************************************************
    private PullToRefreshExpandableListView mTicketPullRefreshListView;//第2个第2个卡券账单的listview


    //**收款码账单**************************************************************************************
    private PullToRefreshExpandableListView mCollectionPullRefreshListView;//第3个第3个收款码账单的listview
    private CollectionExpandableListAdapter mCollectionAdapter;

    private Map<String, MonthBill> mMonthCollectionBillMap;
    private List<MonthBill> mMonthCollectionBillList;//收款码的月账单

    private List<List<TradeBill>> mCollectionBillList;//收款码的日账单，这个条目会很多的
    private Map<String, List<TradeBill>> mCollectionBillMap;

    private int collectionIndex;//收款码账单 使用到的 index索引值
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

        quickPayService = ShowMoneyApp.getInstance().getQuickPayService();

        initLayout();
    }

    private void initLayout() {
        SimpleDateFormat spf = new SimpleDateFormat("yyyyMM");
        mCurrentYearMonth = spf.format(new Date());

        mTitle = (TextView) findViewById(R.id.tv_title);

        //***普通的收款账单***************************************************************************************
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

        //***卡券账单***************************************************************************************
        mTicketPullRefreshListView = (PullToRefreshExpandableListView) findViewById(R.id.ticket_list_view);
        mTicketPullRefreshListView.setVisibility(GONE);


        //***收款码账单***************************************************************************************
        mCollectionPullRefreshListView = (PullToRefreshExpandableListView) findViewById(R.id.colloction_list_view);
        mCollectionPullRefreshListView.setMode(PullToRefreshBase.Mode.BOTH);
        mCollectionPullRefreshListView.setVisibility(GONE);
        mCollectionPullRefreshListView.setOnRefreshListener(new PullToRefreshBase.OnRefreshListener<ExpandableListView>() {
            @Override
            public void onRefresh(PullToRefreshBase<ExpandableListView> refreshView) {
                getCollectionBill();
            }
        });

        mMonthCollectionBillList = new ArrayList<>();
        mCollectionBillList = new ArrayList<>();
        mMonthCollectionBillMap = new HashMap<>();
        mCollectionBillMap = new HashMap<>();

        mCollectionAdapter = new CollectionExpandableListAdapter(mContext, mMonthCollectionBillList, mCollectionBillList);
        ExpandableListView ActualView1 = mCollectionPullRefreshListView.getRefreshableView();
        ActualView1.setAdapter(mCollectionAdapter);
        ActualView1.setGroupIndicator(null);

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
        getCollectionBill();
    }


    //获取收款的账单账单
    private void getBill() {
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
                    tradeBill.transStatus = txn.getTransStatus();
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

                mapToMonthBillList(mMonthBillMap, mMonthBilList);
                mapToBillList(mTradeBillMap, mTradeBillList);

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
        String size = "100";
        String index = String.valueOf(collectionIndex);
        /**
         * recType
         * "收款方式：移动版 桌面版 收款码 开放接口
         *移动版 1
         *桌面版 2
         *收款码 4
         *开放接口 8
         *移动版｜桌面版：1 | 2 = 3
         *移动版 | 收款码:1 | 4 = 5
         *全部：1 | 2 | 4 | 8 = 15 "
         */
        String recType = "4";

        /**
         * payType
         *  "支付方式：支付宝 微信
         *  支付宝 1
         *  微信 2
         *  全部：1 | 2 = 3"
         */
        String payType = "3";

        /**
         *txnStatus
         *"交易状态：交易成功 部分退款 全额退款
         *  交易成功 1
         *  部分退款 2
         *  全额退款 4
         *  部分退款 ｜ 全额退款：2 | 4 = 6
         *  全部：1 | 2 | 4 = 7"
         */
        String txnStatus = "7";

        quickPayService.findOrderAsync(SessonData.loginUser, index, size, recType, payType, txnStatus, new QuickPayCallbackListener<ServerPacket>() {
            @Override
            public void onSuccess(ServerPacket data) {
                final int totalRecord = data.getTotalRecord();//这个字段表示当月的总条数
                final int count = data.getCount();//返回的条数
                final String total = data.getTotal();
                final int refdcount = data.getRefdcount();
                final String refdtotal = data.getRefdtotal();
                final int size = data.getSize();

                //**********************************************************************************
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

                    //添加到收款码账单相应的map中，最后再转换到list中，按照月份的先后排序转换到list中
                    if (mMonthCollectionBillMap.containsKey(currentYearMonth)) {
                        mMonthCollectionBillMap.get(currentYearMonth).setCount(count);
                        mMonthCollectionBillMap.get(currentYearMonth).setTotal(total);
                        mMonthCollectionBillMap.get(currentYearMonth).setRefdcount(refdcount);
                        mMonthCollectionBillMap.get(currentYearMonth).setRefdtotal(refdtotal);
                        mMonthCollectionBillMap.get(currentYearMonth).setSize(size);
                        mMonthCollectionBillMap.get(currentYearMonth).setTotalRecord(totalRecord);
                    } else {
                        MonthBill monthBill = new MonthBill(currentYear, currentMonth);
                        monthBill.setCount(count);
                        monthBill.setTotal(total);
                        monthBill.setRefdcount(refdcount);
                        monthBill.setRefdtotal(refdtotal);
                        monthBill.setSize(size);
                        monthBill.setTotalRecord(totalRecord);
                        mMonthCollectionBillMap.put(currentYearMonth, monthBill);
                    }

                    //收款码账单对应的map
                    if (mCollectionBillMap.containsKey(currentYearMonth)) {
                        mCollectionBillMap.get(currentYearMonth).add(tradeBill);
                    } else {
                        List<TradeBill> list = new ArrayList<TradeBill>();
                        list.add(tradeBill);
                        mCollectionBillMap.put(currentYearMonth, list);
                    }
                }//end for()
                //**********************************************************************************

                //这里通过这个两个方法，把map类型转换为list类型
                mapToMonthBillList(mMonthCollectionBillMap, mMonthCollectionBillList);
                mapToBillList(mCollectionBillMap, mCollectionBillList);

                mCollectionAdapter.notifyDataSetChanged();//这一句很重要的
                mCollectionPullRefreshListView.onRefreshComplete();

                collectionIndex += size;
            }

            @Override
            public void onFailure(QuickPayException ex) {
                mCollectionPullRefreshListView.onRefreshComplete();
            }
        });


    }

    private void mapToMonthBillList(Map<String, MonthBill> map, List<MonthBill> monthList) {
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

        monthList.clear();
        for (String key : list) {
            monthList.add(map.get(key));
        }
    }

    private void mapToBillList(Map<String, List<TradeBill>> map, List<List<TradeBill>> monthList) {
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

        monthList.clear();
        for (String key : list) {
            monthList.add(map.get(key));
        }
    }
}
