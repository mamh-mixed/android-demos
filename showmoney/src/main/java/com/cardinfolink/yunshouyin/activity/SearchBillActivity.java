package com.cardinfolink.yunshouyin.activity;

import android.content.Context;
import android.os.Bundle;
import android.app.Activity;
import android.text.Editable;
import android.text.TextUtils;
import android.text.TextWatcher;
import android.util.Log;
import android.view.View;
import android.widget.CheckBox;
import android.widget.CompoundButton;
import android.widget.ExpandableListView;
import android.widget.LinearLayout;
import android.widget.TextView;
import android.widget.Toast;

import com.cardinfolink.cashiersdk.util.TxamtUtil;
import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.adapter.BillExpandableListAdapter;
import com.cardinfolink.yunshouyin.api.QuickPayException;
import com.cardinfolink.yunshouyin.core.QuickPayCallbackListener;
import com.cardinfolink.yunshouyin.core.QuickPayService;
import com.cardinfolink.yunshouyin.data.MonthBill;
import com.cardinfolink.yunshouyin.data.SessonData;
import com.cardinfolink.yunshouyin.data.TradeBill;
import com.cardinfolink.yunshouyin.model.CouponInfo;
import com.cardinfolink.yunshouyin.model.QRequest;
import com.cardinfolink.yunshouyin.model.ServerPacket;
import com.cardinfolink.yunshouyin.model.Txn;
import com.cardinfolink.yunshouyin.ui.EditTextClear;
import com.cardinfolink.yunshouyin.ui.SettingActionBarItem;
import com.cardinfolink.yunshouyin.util.ShowMoneyApp;
import com.cardinfolink.yunshouyin.view.LoadingDialog;
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

public class SearchBillActivity extends Activity {
    private static final String TAG = "SearchBillActivity";
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


    private SettingActionBarItem mActionBar;


    private TextView mSearch;//搜索的按钮
    private EditTextClear mSearchEditText;
    private LinearLayout mSearchLinearLayout;
    private LinearLayout mSearchConditionLinearLayout;

    //定义几个搜索条件 的checkbpx组件
    private CheckBox mPaySuccessCheckBox;//支付成功的 1

    private CheckBox mRecAppCheckBox;//app 收款的 1
    private CheckBox mRecPCCheckBox;//pc 收款的   2
    private CheckBox mRecWebCheckBox;//网页收款的  4
    private CheckBox mRecOpenCheckBox;//开放接口的 8
    //这里少折扣券的
    private CheckBox mPayAliCheckBox;//支付宝支付的 1
    private CheckBox mPayWxCheckBox;//微信支付的    2

    protected LoadingDialog mLoadingDialog;    //显示loading


    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_search_bill);
        mContext = this;
        quickPayService = ShowMoneyApp.getInstance().getQuickPayService();
        initLayout();
    }

    private void initLayout() {

        SimpleDateFormat spf = new SimpleDateFormat("yyyyMM");
        mCurrentYearMonth = spf.format(new Date());
        billIndex = 0;

        mActionBar = (SettingActionBarItem) findViewById(R.id.action_bar);
        mActionBar.setLeftTextOnclickListner(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                finish();
            }
        });

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

                //通过判断 radio group是否是隐藏的状态，隐藏的状态就是按条件查找的情况
                searchBill();

            }

            @Override
            public void onPullUpToRefresh(PullToRefreshBase<ExpandableListView> refreshView) {
                searchBill();
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


        mSearch = (TextView) findViewById(R.id.tv_search);
        mSearchEditText = (EditTextClear) findViewById(R.id.et_search);

        mSearch.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                //订单号大于等于17位才触发查询操作
                String orderNum = mSearchEditText.getText().toString();
                if (!TextUtils.isEmpty(orderNum)) {
                    //这里精确查找
                    billIndex = 0;
                    mMonthBillAgo = 0;//注意这里要清零
                    mTradeBillList.clear();
                    mMonthBilList.clear();
                    mMonthBillMap.clear();
                    mTradeBillMap.clear();
                    SimpleDateFormat spf = new SimpleDateFormat("yyyyMM");
                    mCurrentYearMonth = spf.format(new Date());

                    findBill(orderNum);
                }
            }
        });

        mSearchEditText.addTextChangedListener(new TextWatcher() {
            @Override
            public void beforeTextChanged(CharSequence s, int start, int count, int after) {

            }

            @Override
            public void onTextChanged(CharSequence s, int start, int before, int count) {

            }


            @Override
            public void afterTextChanged(Editable s) {
                //订单号大于等于17位才触发查询操作
                String orderNum = mSearchEditText.getText().toString();
                if (!TextUtils.isEmpty(orderNum) && orderNum.length() == 17) {
                    Log.e(TAG, " search order: " + orderNum);
                    //这里精确查找
                    billIndex = 0;
                    mMonthBillAgo = 0;//注意这里要清零
                    mTradeBillList.clear();
                    mMonthBilList.clear();
                    mMonthBillMap.clear();
                    mTradeBillMap.clear();
                    SimpleDateFormat spf = new SimpleDateFormat("yyyyMM");
                    mCurrentYearMonth = spf.format(new Date());

                    findBill(orderNum);
                }
            }
        });


        //定义几个搜索条件 的checkbpx组件
        mPaySuccessCheckBox = (CheckBox) findViewById(R.id.cb_success);//支付成功的 1
        mPaySuccessCheckBox.setOnCheckedChangeListener(new SearchCheckBoxOnCheckedChangeListener());

        mRecAppCheckBox = (CheckBox) findViewById(R.id.cb_rec_type1);//app 收款的 1
        mRecAppCheckBox.setOnCheckedChangeListener(new SearchCheckBoxOnCheckedChangeListener());
        mRecPCCheckBox = (CheckBox) findViewById(R.id.cb_rec_type2);//pc 收款的   2
        mRecPCCheckBox.setOnCheckedChangeListener(new SearchCheckBoxOnCheckedChangeListener());
        mRecWebCheckBox = (CheckBox) findViewById(R.id.cb_rec_type3);//网页收款的  4
        mRecWebCheckBox.setOnCheckedChangeListener(new SearchCheckBoxOnCheckedChangeListener());
        mRecOpenCheckBox = (CheckBox) findViewById(R.id.cb_rec_type4);//其他收款的开放接口的 8
        mRecOpenCheckBox.setOnCheckedChangeListener(new SearchCheckBoxOnCheckedChangeListener());
        //这里少折扣券的
        mPayAliCheckBox = (CheckBox) findViewById(R.id.cb_pay_type1);//支付宝支付的 1
        mPayAliCheckBox.setOnCheckedChangeListener(new SearchCheckBoxOnCheckedChangeListener());
        mPayWxCheckBox = (CheckBox) findViewById(R.id.cb_pay_type2);//微信支付的    2
        mPayWxCheckBox.setOnCheckedChangeListener(new SearchCheckBoxOnCheckedChangeListener());

    }


    private class SearchCheckBoxOnCheckedChangeListener implements CompoundButton.OnCheckedChangeListener {

        @Override
        public void onCheckedChanged(CompoundButton buttonView, boolean isChecked) {
            if (!TextUtils.isEmpty(mSearchEditText.getText())) {
                return;
            }
            //按条件查找之前先清空一下数据
            billIndex = 0;
            mMonthBillAgo = 0;//注意这里要清零
            mTradeBillList.clear();
            mMonthBilList.clear();
            mMonthBillMap.clear();
            mTradeBillMap.clear();
            mBillAdapter.notifyDataSetChanged();
            SimpleDateFormat spf = new SimpleDateFormat("yyyyMM");
            mCurrentYearMonth = spf.format(new Date());
            searchBill();
        }
    }

    //精确查找某个账单
    private void findBill(String orderNum) {
        mLoadingDialog.startLoading();
        quickPayService.getOrderAsync(SessonData.loginUser, orderNum, new QuickPayCallbackListener<ServerPacket>() {
            @Override
            public void onSuccess(ServerPacket data) {
                parseServerPacket(data, mMonthBillMap, mTradeBillMap, mMonthBilList, mTradeBillList);

                mBillAdapter.notifyDataSetChanged();
                mBillPullRefreshListView.onRefreshComplete();

                if (mMonthBilList.size() <= 0) {
                    String msg = mContext.getString(R.string.bill_search_result_message1);
                    Toast.makeText(mContext, msg, Toast.LENGTH_SHORT).show();
                }

                mLoadingDialog.endLoading();
            }

            @Override
            public void onFailure(QuickPayException ex) {
                mBillAdapter.notifyDataSetChanged();
                mBillPullRefreshListView.onRefreshComplete();
                String msg = mContext.getString(R.string.bill_search_result_message2) + ex.getErrorMsg();
                Toast.makeText(mContext, msg, Toast.LENGTH_SHORT).show();
                mLoadingDialog.endLoading();
            }
        });

    }

    private void searchBill() {
        int recValue = 0;//收款方式
        if (mRecAppCheckBox.isChecked()) {
            recValue += 1;
        }
        if (mRecPCCheckBox.isChecked()) {
            recValue += 2;
        }
        if (mRecWebCheckBox.isChecked()) {
            recValue += 4;
        }
        if (mRecOpenCheckBox.isChecked()) {
            recValue += 8;
        }
        if (recValue == 0) {
            recValue = 15;
        }


        int payValue = 0;//支付方式
        if (mPayAliCheckBox.isChecked()) {
            payValue += 1;
        }
        if (mPayWxCheckBox.isChecked()) {
            payValue += 2;
        }
        if (payValue == 0) {
            payValue = 3;
        }

        int txnStatus = 0;//支付状态
        if (mPaySuccessCheckBox.isChecked()) {
            txnStatus += 1;
        }
        if (txnStatus == 0) {
            txnStatus = 7;
        }

        getBill(String.valueOf(recValue), String.valueOf(payValue), String.valueOf(txnStatus));
    }

    private void getBill(String recType, String payType, String txnStatus) {
        mLoadingDialog.startLoading();
        String sizeStr = "100";
        String index = String.valueOf(billIndex);
        quickPayService.findOrderAsync(SessonData.loginUser, index, sizeStr, recType, payType, txnStatus, new QuickPayCallbackListener<ServerPacket>() {
            @Override
            public void onSuccess(ServerPacket data) {
                //这里特殊一些，需要用的size。
                final int size = data.getSize();

                parseServerPacket(data, mMonthBillMap, mTradeBillMap, mMonthBilList, mTradeBillList);

                mBillAdapter.notifyDataSetChanged();
                mBillPullRefreshListView.onRefreshComplete();
                if (mBillAdapter.getGroupCount() >= 1) {
                    mBillPullRefreshListView.getRefreshableView().expandGroup(0);
                }
                billIndex += size;
                if (mMonthBilList.size() <= 0) {
                    String msg = mContext.getString(R.string.bill_search_result_message3);
                    Toast.makeText(mContext, msg, Toast.LENGTH_SHORT).show();
                }
                mLoadingDialog.endLoading();
            }

            @Override
            public void onFailure(QuickPayException ex) {
                mBillAdapter.notifyDataSetChanged();
                mBillPullRefreshListView.onRefreshComplete();
                String msg = mContext.getString(R.string.bill_search_result_message2) + ex.getErrorMsg();
                Toast.makeText(mContext, msg, Toast.LENGTH_SHORT).show();
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
