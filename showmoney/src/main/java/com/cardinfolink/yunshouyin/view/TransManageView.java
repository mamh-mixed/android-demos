package com.cardinfolink.yunshouyin.view;

import android.content.Context;
import android.content.Intent;
import android.os.Handler;
import android.text.TextUtils;
import android.view.LayoutInflater;
import android.view.View;
import android.widget.ExpandableListView;
import android.widget.LinearLayout;
import android.widget.RadioButton;
import android.widget.TextView;

import com.cardinfolink.cashiersdk.util.TxamtUtil;
import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.activity.SearchBillActivity;
import com.cardinfolink.yunshouyin.adapter.BillExpandableListAdapter;
import com.cardinfolink.yunshouyin.adapter.CollectionExpandableListAdapter;
import com.cardinfolink.yunshouyin.adapter.TicketExpandableListAdapter;
import com.cardinfolink.yunshouyin.api.QuickPayException;
import com.cardinfolink.yunshouyin.core.QuickPayCallbackListener;
import com.cardinfolink.yunshouyin.core.QuickPayService;
import com.cardinfolink.yunshouyin.data.MonthBill;
import com.cardinfolink.yunshouyin.data.SessionData;
import com.cardinfolink.yunshouyin.data.TradeBill;
import com.cardinfolink.yunshouyin.model.CouponInfo;
import com.cardinfolink.yunshouyin.model.QRequest;
import com.cardinfolink.yunshouyin.model.ServerPacket;
import com.cardinfolink.yunshouyin.model.Txn;
import com.cardinfolink.yunshouyin.util.ShowMoneyApp;
import com.handmark.pulltorefresh.library.PullToRefreshBase;
import com.handmark.pulltorefresh.library.PullToRefreshExpandableListView;

import java.text.SimpleDateFormat;
import java.util.ArrayList;
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

    private String mCurrentYearMonth;//当前年份+月份的一个字符串

    private View mEmptyViewBill;
    private TextView mEmptyTextviewBill;

    //***卡券账单*************************************************************************************
    private PullToRefreshExpandableListView mTicketPullRefreshListView;//第2个第2个卡券账单的listview
    private TicketExpandableListAdapter mTicketAdapter;

    private Map<String, MonthBill> mMonthTicketBillMap;
    private List<MonthBill> mMonthTicketBilltList;

    private List<List<TradeBill>> mTicketBillList;
    private Map<String, List<TradeBill>> mTicketBillMap;

    private int ticketIndex;
    private int mMonthTicketAgo;
    private String mTicketCurrentYearMonth;

    private View mEmptyViewTicket;
    private TextView mEmptyTextviewTicket;

    //**收款码账单**************************************************************************************
    private PullToRefreshExpandableListView mCollectionPullRefreshListView;//第3个第3个收款码账单的listview
    private CollectionExpandableListAdapter mCollectionAdapter;

    private Map<String, MonthBill> mMonthCollectionBillMap;
    private List<MonthBill> mMonthCollectionBillList;//收款码的月账单

    private List<List<TradeBill>> mCollectionBillList;//收款码的日账单，这个条目会很多的
    private Map<String, List<TradeBill>> mCollectionBillMap;

    private int collectionIndex;//收款码账单 使用到的 index索引值

    private View mEmptyViewCollection;
    private TextView mEmptyTextviewCollection;

    //****************************************************************************************


    private RadioButton mRaidoBill;//收款账单
    private RadioButton mRadioTicket;//卡券账单
    private RadioButton mRadioCollection;//收款码账单


    private TextView mSearch;//搜索的按钮


    private TextView mTitle;

    protected LoadingDialog mLoadingDialog;    //显示loading

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
        mTicketCurrentYearMonth = mCurrentYearMonth;
        billIndex = ticketIndex = 0;

        mTitle = (TextView) findViewById(R.id.tv_title);

        mLoadingDialog = new LoadingDialog(mContext, findViewById(R.id.loading_dialog));

        //***普通的收款账单***************************************************************************************
        mBillPullRefreshListView = (PullToRefreshExpandableListView) findViewById(R.id.bill_list_view);
        mBillPullRefreshListView.setMode(PullToRefreshBase.Mode.BOTH);
        mBillPullRefreshListView.setOnRefreshListener(new PullToRefreshBase.OnRefreshListener2<ExpandableListView>() {
            @Override
            public void onPullDownToRefresh(PullToRefreshBase<ExpandableListView> refreshView) {
                billIndex = 0;
                mMonthBillAgo = 0;//注意这里要清零
                mTradeBillList.clear();
                mMonthBilList.clear();
                mMonthBillMap.clear();
                mTradeBillMap.clear();
                mBillAdapter.notifyDataSetChanged();
                SimpleDateFormat spf = new SimpleDateFormat("yyyyMM");
                mCurrentYearMonth = spf.format(new Date());

                getBill();
            }

            @Override
            public void onPullUpToRefresh(PullToRefreshBase<ExpandableListView> refreshView) {
                getBill();
            }
        });

        mMonthBilList = new ArrayList<>();
        mTradeBillList = new ArrayList<>();
        mMonthBillMap = new HashMap<>();
        mTradeBillMap = new HashMap<>();

        mBillAdapter = new BillExpandableListAdapter(mContext, mMonthBilList, mTradeBillList);
        mBillAdapter.setHintDialog(findViewById(R.id.hint_dialog));

        ExpandableListView billActualView = mBillPullRefreshListView.getRefreshableView();
        billActualView.setAdapter(mBillAdapter);
        billActualView.setGroupIndicator(null);

        //***卡券账单***************************************************************************************
        mTicketPullRefreshListView = (PullToRefreshExpandableListView) findViewById(R.id.ticket_list_view);
        mTicketPullRefreshListView.setMode(PullToRefreshBase.Mode.BOTH);
        mTicketPullRefreshListView.setVisibility(GONE);
        mTicketPullRefreshListView.setOnRefreshListener(new PullToRefreshBase.OnRefreshListener2<ExpandableListView>() {
            @Override
            public void onPullDownToRefresh(PullToRefreshBase<ExpandableListView> refreshView) {
                ticketIndex = 0;
                mMonthTicketAgo = 0;
                mMonthTicketBillMap.clear();
                mMonthTicketBilltList.clear();
                mTicketBillList.clear();
                mTicketBillMap.clear();

                mTicketAdapter.notifyDataSetChanged();
                SimpleDateFormat spf = new SimpleDateFormat("yyyyMM");
                mTicketCurrentYearMonth = spf.format(new Date());
                getTicketBill();
            }

            @Override
            public void onPullUpToRefresh(PullToRefreshBase<ExpandableListView> refreshView) {
                getTicketBill();
            }
        });

        mMonthTicketBillMap = new HashMap<>();
        mMonthTicketBilltList = new ArrayList<>();
        mTicketBillList = new ArrayList<>();
        mTicketBillMap = new HashMap<>();

        mTicketAdapter = new TicketExpandableListAdapter(mContext, mMonthTicketBilltList, mTicketBillList);
        mTicketAdapter.setHintDialog(findViewById(R.id.hint_dialog));

        ExpandableListView ticketActualView = mTicketPullRefreshListView.getRefreshableView();
        ticketActualView.setAdapter(mTicketAdapter);
        ticketActualView.setGroupIndicator(null);

        //***收款码账单***************************************************************************************
        mCollectionPullRefreshListView = (PullToRefreshExpandableListView) findViewById(R.id.colloction_list_view);
        mCollectionPullRefreshListView.setMode(PullToRefreshBase.Mode.BOTH);
        mCollectionPullRefreshListView.setVisibility(GONE);
        mCollectionPullRefreshListView.setOnRefreshListener(new PullToRefreshBase.OnRefreshListener2<ExpandableListView>() {
            @Override
            public void onPullDownToRefresh(PullToRefreshBase<ExpandableListView> refreshView) {
                collectionIndex = 0;
                mMonthCollectionBillList.clear();
                mCollectionBillList.clear();
                mMonthCollectionBillMap.clear();
                mCollectionBillMap.clear();
                mCollectionAdapter.notifyDataSetChanged();
                getCollectionBill();
            }

            @Override
            public void onPullUpToRefresh(PullToRefreshBase<ExpandableListView> refreshView) {
                getCollectionBill();
            }
        });

        mMonthCollectionBillList = new ArrayList<>();
        mCollectionBillList = new ArrayList<>();
        mMonthCollectionBillMap = new HashMap<>();
        mCollectionBillMap = new HashMap<>();

        mCollectionAdapter = new CollectionExpandableListAdapter(mContext, mMonthCollectionBillList, mCollectionBillList);
        ExpandableListView collectionActualView = mCollectionPullRefreshListView.getRefreshableView();
        collectionActualView.setAdapter(mCollectionAdapter);
        collectionActualView.setGroupIndicator(null);

        //******************************************************************************************


        mSearch = (TextView) findViewById(R.id.tv_search);
        mSearch.setOnClickListener(new OnClickListener() {
            @Override
            public void onClick(View v) {
                Intent intent = new Intent(mContext, SearchBillActivity.class);
                mContext.startActivity(intent);
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
                mTitle.setText(mRaidoBill.getText().toString());//设置标题
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
                mTitle.setText(mRadioTicket.getText().toString());
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
                mTitle.setText(mRadioCollection.getText().toString());
            }
        });


        //设置一个空的view，当listview为空的时候
        mEmptyViewBill = View.inflate(mContext, R.layout.expandable_list_view_empty, null);
        mEmptyTextviewBill = (TextView) mEmptyViewBill.findViewById(R.id.tv_message);
        mEmptyTextviewBill.setText(mRaidoBill.getText());
        mBillPullRefreshListView.setEmptyView(mEmptyViewBill);

        mEmptyViewTicket = View.inflate(mContext, R.layout.expandable_list_view_empty, null);
        mEmptyTextviewTicket = (TextView) mEmptyViewTicket.findViewById(R.id.tv_message);
        mEmptyTextviewTicket.setText(mRadioTicket.getText());
        mTicketPullRefreshListView.setEmptyView(mEmptyViewTicket);

        mEmptyViewCollection = View.inflate(mContext, R.layout.expandable_list_view_empty, null);
        mEmptyTextviewCollection = (TextView) mEmptyViewCollection.findViewById(R.id.tv_message);
        mEmptyTextviewCollection.setText(mRadioCollection.getText());
        mCollectionPullRefreshListView.setEmptyView(mEmptyViewCollection);
    }


    public void refresh() {
        billIndex = 0;
        mMonthBillAgo = 0;//注意这里要清零
        mTradeBillList.clear();
        mMonthBilList.clear();
        mMonthBillMap.clear();
        mTradeBillMap.clear();
        mBillAdapter.notifyDataSetChanged();
        SimpleDateFormat spf = new SimpleDateFormat("yyyyMM");
        mCurrentYearMonth = spf.format(new Date());

        getBill();

        getCollectionBill();

        ticketIndex = 0;
        mMonthTicketAgo = 0;
        mMonthTicketBillMap.clear();
        mMonthTicketBilltList.clear();
        mTicketBillList.clear();
        mTicketBillMap.clear();
        mTicketAdapter.notifyDataSetChanged();
        mTicketCurrentYearMonth = spf.format(new Date());

        getTicketBill();
    }

    //局部账单跟新的方法，传人groupdata位置，和child的位置
    public void refresh(int groupPosition, int childPosition) {
        String orderNum;
        final TradeBill bill;
        try {
            bill = mTradeBillList.get(groupPosition).get(childPosition);
            orderNum = bill.orderNum;//账单号
        } catch (Exception e) {
            return;
        }

        quickPayService.getOrderAsync(SessionData.loginUser, orderNum, new QuickPayCallbackListener<ServerPacket>() {
            @Override
            public void onSuccess(ServerPacket data) {
                try {
                    Txn[] txn = data.getTxn();
                    bill.refundAmt = TxamtUtil.getNormal(txn[0].getRefundAmt());//对于人民币的金额都需要除以100
                    bill.transStatus = txn[0].getTransStatus();
                    mBillAdapter.notifyDataSetChanged();
                } catch (Exception e) {

                }
            }

            @Override
            public void onFailure(QuickPayException ex) {
                //出错不处理
            }
        });


    }

    //获取收款的账单账单
    private void getBill() {
        mLoadingDialog.startLoading();

        quickPayService.getHistoryBillsAsync(SessionData.loginUser, mCurrentYearMonth, String.valueOf(billIndex), "100", "success", new QuickPayCallbackListener<ServerPacket>() {
            @Override
            public void onSuccess(ServerPacket data) {
                //这里可以在ui线程里执行的
                //这里特殊一些，需要用的size

                parseServerPacket(data, mMonthBillMap, mTradeBillMap, mMonthBilList, mTradeBillList);

                mBillAdapter.notifyDataSetChanged();
                mBillPullRefreshListView.onRefreshComplete();
                if (mBillAdapter.getGroupCount() >= 1) {
                    mBillPullRefreshListView.getRefreshableView().expandGroup(0);
                }

                mCurrentYearMonth = data.getNextMonth();

                mLoadingDialog.endLoading();
            }

            @Override
            public void onFailure(QuickPayException ex) {
                mBillAdapter.notifyDataSetChanged();
                mBillPullRefreshListView.onRefreshComplete();
                mLoadingDialog.endLoading();
            }
        });
    }

    //获取卡券账单
    public void getTicketBill() {
        mLoadingDialog.startLoading();

        quickPayService.getHistoryCouponsAsync(SessionData.loginUser, mTicketCurrentYearMonth, String.valueOf(ticketIndex), "100", new QuickPayCallbackListener<ServerPacket>() {
            @Override
            public void onSuccess(ServerPacket data) {
                parseServerPacket(data, mMonthTicketBillMap, mTicketBillMap, mMonthTicketBilltList, mTicketBillList);

                mTicketAdapter.notifyDataSetChanged();
                mTicketPullRefreshListView.onRefreshComplete();

                if (mTicketAdapter.getGroupCount() >= 1) {
                    mTicketPullRefreshListView.getRefreshableView().expandGroup(0);
                }

                mTicketCurrentYearMonth = data.getNextMonth();

                mLoadingDialog.endLoading();
            }

            @Override
            public void onFailure(QuickPayException ex) {
                mTicketAdapter.notifyDataSetChanged();
                mTicketPullRefreshListView.onRefreshComplete();
                mLoadingDialog.endLoading();
            }
        });

    }

    //获取 收款码 账单
    public void getCollectionBill() {
        mLoadingDialog.startLoading();
        String sizeStr = "100";
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

        quickPayService.findOrderAsync(SessionData.loginUser, index, sizeStr, recType, payType, txnStatus, new QuickPayCallbackListener<ServerPacket>() {
            @Override
            public void onSuccess(ServerPacket data) {
                //这里特殊一些，需要用的size。
                final int size = data.getSize();

                parseServerPacket(data, mMonthCollectionBillMap, mCollectionBillMap, mMonthCollectionBillList, mCollectionBillList);

                mCollectionAdapter.notifyDataSetChanged();//这一句很重要的
                mCollectionPullRefreshListView.onRefreshComplete();
                if (mCollectionAdapter.getGroupCount() >= 1) {
                    mCollectionPullRefreshListView.getRefreshableView().expandGroup(0);
                }
                collectionIndex += size;

                mLoadingDialog.endLoading();
            }

            @Override
            public void onFailure(QuickPayException ex) {
                mCollectionAdapter.notifyDataSetChanged();//这一句很重要的
                mCollectionPullRefreshListView.onRefreshComplete();
                mLoadingDialog.endLoading();
            }
        });


    }


    private void parseServerPacket(ServerPacket data, Map<String, MonthBill> monthMap, Map<String, List<TradeBill>> tradeBillMap, List<MonthBill> monthList, List<List<TradeBill>> tradeBillList) {
        //处理服务器返回的ServerPacket data的数据，把他们都放在相应的list和map中去
        //这个方法是把serverpacket 数据加工成我们需要的list类型的数据
        // ServerPacket data, 这个就是服务器返回的data数据，
        // Map<String,MonthBill> monthMap, 这个是传人的月的账单，保存的是和月份相关的收入总和，
        //      账单总数的统计信息 key,value 键-值对的形象保存在map中。 例如：201601 --》 new MonthBill();
        // Map<String, List<TradeBill>> tradeBillMap, 这个是日账单详情，保存的是每天每天的账单的
        //      详情 key,value 键-值对的形象保存在map中。 例如：201601 --》对应一个 list，list中保存了这个月的日账单。
        // List<MonthBill> monthList, 这个是map转换来的
        // List<List<TradeBill>> tradeBillList) { 这个是map转换来的 }
        final int totalRecord = data.getTotalRecord();//这个字段表示当月的总条数
        final int count = data.getCount();//返回的条数,这个基本上没有用到
        final String total = TxamtUtil.getNormal(data.getTotal());
        final int refdcount = data.getRefdcount();
        final String refdtotal = data.getRefdtotal();
        final int size = data.getSize();

        //这里开始遍历这个账单的数组************************************************************
        if (data.getTxn() != null) {
            for (Txn txn : data.getTxn()) {
                TradeBill tradeBill = new TradeBill();
                tradeBill.response = txn.getResponse();
                tradeBill.tandeDate = txn.getSystemDate();
                tradeBill.consumerAccount = txn.getConsumerAccount();
                tradeBill.transStatus = txn.getTransStatus();
                tradeBill.refundAmt = TxamtUtil.getNormal(txn.getRefundAmt());//对于人民币的金额都需要除以100
                tradeBill.couponDiscountAmt = TxamtUtil.getNormal(txn.getCouponDiscountAmt());//卡券优惠金额，人民币需要除以100

                //收款码账单需要的三个数据，
                tradeBill.nickName = txn.getNickName();//微信账号的昵称
                tradeBill.avatarUrl = txn.getAvatarUrl();//微信头像地址
                tradeBill.checkCode = txn.getCheckCode();//检验码

                QRequest req = txn.getmRequest();
                if (req != null) {
                    tradeBill.orderNum = req.getOrderNum();
                    tradeBill.amount = TxamtUtil.getNormal(req.getTxamt());//对于人民币的金额都需要除以100
                    tradeBill.busicd = req.getBusicd();
                    tradeBill.terminalid = req.getTerminalid();

                    //使用/v3/bill接口 退款的好像也没有拉取到
                    if (tradeBill.busicd.equals("REFD")) {
                        tradeBill.amount = "-" + tradeBill.amount;
                    }
                    tradeBill.chcd = req.getChcd();
                    tradeBill.tradeFrom = req.getTradeFrom();
                    tradeBill.goodsInfo = req.getGoodsInfo();
                }

                if (TextUtils.isEmpty(tradeBill.tandeDate)) {
                    continue;
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
                if (monthMap.containsKey(currentYearMonth)) {
                    monthMap.get(currentYearMonth).setCount(count);
                    monthMap.get(currentYearMonth).setTotal(total);
                    monthMap.get(currentYearMonth).setRefdcount(refdcount);
                    monthMap.get(currentYearMonth).setRefdtotal(refdtotal);
                    monthMap.get(currentYearMonth).setSize(size);
                    monthMap.get(currentYearMonth).setTotalRecord(totalRecord);
                } else {
                    MonthBill monthBill = new MonthBill(currentYear, currentMonth);
                    monthBill.setCount(count);
                    monthBill.setTotal(total);
                    monthBill.setRefdcount(refdcount);
                    monthBill.setRefdtotal(refdtotal);
                    monthBill.setSize(size);
                    monthBill.setTotalRecord(totalRecord);
                    monthMap.put(currentYearMonth, monthBill);
                }

                if (tradeBillMap.containsKey(currentYearMonth)) {
                    tradeBillMap.get(currentYearMonth).add(tradeBill);
                } else {
                    List<TradeBill> list = new ArrayList<TradeBill>();
                    list.add(tradeBill);
                    tradeBillMap.put(currentYearMonth, list);
                }
            }//end for()
        }
        //**********************************************************************************

        //***这里是获取卡券账单数组**这个和上面的不可能同时有数据的*********************************
        if (data.getCoupons() != null) {
            for (CouponInfo couponInfo : data.getCoupons()) {
                TradeBill tradeBill = new TradeBill();
                tradeBill.response = couponInfo.getResponse();
                tradeBill.tandeDate = couponInfo.getSystemDate();
                tradeBill.tradeFrom = couponInfo.getTradeFrom();
                tradeBill.couponType = couponInfo.getType();
                tradeBill.couponName = couponInfo.getName();
                tradeBill.couponChannel = couponInfo.getChannel();
                tradeBill.terminalid = couponInfo.getTerminalid();
                tradeBill.couponOrderNum = couponInfo.getOrderNum();

                //Qrequest里面的都是交易相关的信息，外面的是和卡券相关的信息
                QRequest req = couponInfo.getmRequest();
                if (req != null) {
                    tradeBill.orderNum = req.getOrderNum();
                    tradeBill.amount = TxamtUtil.getNormal(req.getTxamt());//对于人民币的金额都需要除以100
                    tradeBill.busicd = req.getBusicd();
                    tradeBill.couponDiscountAmt = TxamtUtil.getNormal(req.getCouponDiscountAmt());

                    tradeBill.chcd = req.getChcd();
                }
                //获取这个账单里面的日期,年月日 的 日
                if (TextUtils.isEmpty(tradeBill.tandeDate)) {
                    continue;
                }
                String currentDay = tradeBill.tandeDate.substring(6, 8);
                String currentYear = tradeBill.tandeDate.substring(0, 4);
                String currentMonth = tradeBill.tandeDate.substring(4, 6);
                String currentYearMonth = tradeBill.tandeDate.substring(0, 6);


                //添加到相应的map中，最后再转换到list中，按照月份的先后排序转换到list中
                if (monthMap.containsKey(currentYearMonth)) {
                    monthMap.get(currentYearMonth).setCount(count);
                    monthMap.get(currentYearMonth).setTotal(total);
                    monthMap.get(currentYearMonth).setRefdcount(refdcount);
                    monthMap.get(currentYearMonth).setRefdtotal(refdtotal);
                    monthMap.get(currentYearMonth).setSize(size);
                    monthMap.get(currentYearMonth).setTotalRecord(totalRecord);
                } else {
                    MonthBill monthBill = new MonthBill(currentYear, currentMonth);
                    monthBill.setCount(count);
                    monthBill.setTotal(total);
                    monthBill.setRefdcount(refdcount);
                    monthBill.setRefdtotal(refdtotal);
                    monthBill.setSize(size);
                    monthBill.setTotalRecord(totalRecord);
                    monthMap.put(currentYearMonth, monthBill);
                }

                if (tradeBillMap.containsKey(currentYearMonth)) {
                    tradeBillMap.get(currentYearMonth).add(tradeBill);
                } else {
                    List<TradeBill> list = new ArrayList<TradeBill>();
                    list.add(tradeBill);
                    tradeBillMap.put(currentYearMonth, list);
                }

            }
        }


        //这是里把相应的map类型转换成list类型，因为expandablelistview里面group和child的数据需要是list类型的
        //也不一定没要是list类型，不过list类型在expandablelistview里面用起来方便。
        mapToMonthBillList(monthMap, monthList);
        mapToBillList(tradeBillMap, tradeBillList);
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
