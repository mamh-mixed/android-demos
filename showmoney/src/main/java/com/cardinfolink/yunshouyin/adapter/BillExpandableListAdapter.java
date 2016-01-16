package com.cardinfolink.yunshouyin.adapter;

import android.content.Context;
import android.content.Intent;
import android.graphics.Color;
import android.graphics.Paint;
import android.os.Bundle;
import android.text.TextUtils;
import android.view.View;
import android.view.ViewGroup;
import android.widget.BaseExpandableListAdapter;
import android.widget.ImageView;
import android.widget.LinearLayout;
import android.widget.TextView;

import com.cardinfolink.cashiersdk.util.TxamtUtil;
import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.activity.DetailActivity;
import com.cardinfolink.yunshouyin.api.QuickPayException;
import com.cardinfolink.yunshouyin.core.QuickPayCallbackListener;
import com.cardinfolink.yunshouyin.core.QuickPayService;
import com.cardinfolink.yunshouyin.data.MonthBill;
import com.cardinfolink.yunshouyin.data.SessonData;
import com.cardinfolink.yunshouyin.data.TradeBill;
import com.cardinfolink.yunshouyin.model.ServerPacket;
import com.cardinfolink.yunshouyin.util.ShowMoneyApp;
import com.cardinfolink.yunshouyin.view.HintBillDialog;

import java.math.BigDecimal;
import java.text.ParseException;
import java.text.SimpleDateFormat;
import java.util.Date;
import java.util.List;

/**
 * Created by mamh on 15-12-26.
 */
public class BillExpandableListAdapter extends BaseExpandableListAdapter {
    private String TAG = "BillExpandableListAdapter";

    private List<MonthBill> groupData;
    private List<List<TradeBill>> childrenData;
    private Context mContext;

    private QuickPayService quickPayService;

    private HintBillDialog mHintDialog;

    private boolean isSearch = false;

    public BillExpandableListAdapter(Context context, List<MonthBill> groupData, List<List<TradeBill>> childrenData) {
        this.mContext = context;
        this.groupData = groupData;
        this.childrenData = childrenData;
    }

    public void setHintDialog(View view) {
        this.mHintDialog = new HintBillDialog(mContext, view);
    }

    public void setIsSearch(boolean isSearch) {
        this.isSearch = isSearch;
    }

    @Override
    public int getGroupCount() {
        if (groupData != null) {
            return groupData.size();
        } else {
            return 0;
        }
    }

    @Override
    public int getChildrenCount(int groupPosition) {
        return childrenData.get(groupPosition).size();
    }

    @Override
    public Object getGroup(int groupPosition) {
        return groupData.get(groupPosition);
    }

    @Override
    public Object getChild(int groupPosition, int childPosition) {
        return childrenData.get(groupPosition).get(childPosition);
    }

    @Override
    public long getGroupId(int groupPosition) {
        return 0;
    }

    @Override
    public long getChildId(int groupPosition, int childPosition) {
        return 0;
    }

    @Override
    public boolean hasStableIds() {
        return false;
    }

    @Override
    public View getGroupView(int groupPosition, boolean isExpanded, View convertView, ViewGroup parent) {
        GroupViewHolder groupViewHolder = null;

        if (convertView == null) {
            groupViewHolder = new GroupViewHolder();
            convertView = View.inflate(mContext, R.layout.bill_expandablelistview_group, null);

            groupViewHolder.month = (TextView) convertView.findViewById(R.id.tv_month);
            groupViewHolder.year = (TextView) convertView.findViewById(R.id.tv_year);
            groupViewHolder.total = (TextView) convertView.findViewById(R.id.tv_total);
            groupViewHolder.count = (TextView) convertView.findViewById(R.id.tv_count);
            groupViewHolder.folder = (ImageView) convertView.findViewById(R.id.iv_fold);
            groupViewHolder.leftLinearLayout = (LinearLayout) convertView.findViewById(R.id.ll_left);

            convertView.setTag(groupViewHolder);
        } else {
            groupViewHolder = (GroupViewHolder) convertView.getTag();
        }

        //如果是搜索 出来的账单 不显示 总收入
        if (isSearch) {
            groupViewHolder.leftLinearLayout.setVisibility(View.INVISIBLE);
        } else {
            groupViewHolder.leftLinearLayout.setVisibility(View.VISIBLE);
        }

        //设置一下月份
        groupViewHolder.month.setText(groupData.get(groupPosition).getCurrentMonth());
        groupViewHolder.year.setText(groupData.get(groupPosition).getCurrentYear());


        String totalStr = groupData.get(groupPosition).getTotal();
        if (TextUtils.isEmpty(totalStr)) {
            totalStr = "0.00";
        }
        groupViewHolder.total.setText(totalStr);

        //如果是搜索 出来的账单 使用childrenData list的size来
        if (isSearch) {
            String countStr = String.valueOf(childrenData.get(groupPosition).size());
            groupViewHolder.count.setText(countStr);
        } else {
            String countStr = String.valueOf(groupData.get(groupPosition).getTotalRecord());
            groupViewHolder.count.setText(countStr);
        }


        if (isExpanded) {
            groupViewHolder.folder.setBackgroundResource(R.drawable.bill_pack);
        } else {
            groupViewHolder.folder.setBackgroundResource(R.drawable.bill_unfold);
        }

        return convertView;
    }

    @Override
    public View getChildView(int groupPosition, int childPosition, boolean isLastChild, View convertView, ViewGroup parent) {
        ChildViewHolder childViewHolder = null;

        if (convertView == null) {
            childViewHolder = new ChildViewHolder();
            convertView = View.inflate(mContext, R.layout.bill_expandablelistview_child, null);
            childViewHolder.linearLayoutDay = convertView.findViewById(R.id.ll_day);
            childViewHolder.linearLayoutBillItem = convertView.findViewById(R.id.ll_bill_item);

            childViewHolder.day = (TextView) convertView.findViewById(R.id.tv_day);
            childViewHolder.weekday = (TextView) convertView.findViewById(R.id.tv_weekday);

            childViewHolder.paylogo = (ImageView) convertView.findViewById(R.id.paylogo);
            childViewHolder.billTradeDate = (TextView) convertView.findViewById(R.id.bill_tradedate);
            childViewHolder.billTradeFrom = (TextView) convertView.findViewById(R.id.bill_tv_tradefrom);
            childViewHolder.billTradeFromImage = (ImageView) convertView.findViewById(R.id.bill_iv_tradefrom);
            childViewHolder.billTradeStatus = (TextView) convertView.findViewById(R.id.bill_tradestatus);
            childViewHolder.billTradeAmount = (TextView) convertView.findViewById(R.id.bill_tradeamount);
            childViewHolder.billOriginTradeAmount = (TextView) convertView.findViewById(R.id.bill_origin_trademount);
            childViewHolder.billDiscount = (ImageView) convertView.findViewById(R.id.bill_descount);

            convertView.setTag(childViewHolder);
        } else {
            childViewHolder = (ChildViewHolder) convertView.getTag();
        }

        //从list中根据位置获取到相应的bill项
        final TradeBill bill = childrenData.get(groupPosition).get(childPosition);

        if (!TextUtils.isEmpty(bill.chcd)) {
            //有chcd渠道的话,这里设置不同渠道的图片
            if ("WXP".equals(bill.chcd)) {
                childViewHolder.paylogo.setImageResource(R.drawable.wpay);
            } else if ("ALP".equals(bill.chcd)) {
                childViewHolder.paylogo.setImageResource(R.drawable.apay);
            } else {
                childViewHolder.paylogo.setImageDrawable(null);
            }
        } else {
            childViewHolder.paylogo.setImageDrawable(null);
        }
        SimpleDateFormat spf1 = new SimpleDateFormat("yyyyMMddHHmmss");
        SimpleDateFormat spf2 = new SimpleDateFormat("HH:mm:ss");
        SimpleDateFormat spf3 = new SimpleDateFormat("dd");
        SimpleDateFormat spf4 = new SimpleDateFormat("E");
        try {
            Date tandeDate = spf1.parse(bill.tandeDate);
            childViewHolder.billTradeDate.setText(spf2.format(tandeDate));
            childViewHolder.day.setText(spf3.format(tandeDate));
            childViewHolder.weekday.setText(spf4.format(tandeDate));
        } catch (ParseException e) {
            e.printStackTrace();
        }

        if ("android".equals(bill.tradeFrom) || "ios".equals(bill.tradeFrom)) {
            childViewHolder.billTradeFromImage.setImageResource(R.drawable.bill_phone);
            childViewHolder.billTradeFrom.setText(mContext.getString(R.string.expandable_listview_pay_type1));
        } else if ("wap".equals(bill.tradeFrom)) {
            childViewHolder.billTradeFromImage.setImageResource(R.drawable.bill_web);
            childViewHolder.billTradeFrom.setText(mContext.getString(R.string.expandable_listview_pay_type2));
        } else if ("PC".equals(bill.tradeFrom)) {
            childViewHolder.billTradeFromImage.setImageResource(R.drawable.bill_pc);
            childViewHolder.billTradeFrom.setText(mContext.getString(R.string.expandable_listview_pay_type3));
        } else {
            childViewHolder.billTradeFromImage.setImageResource(R.drawable.bill_else);
            childViewHolder.billTradeFrom.setText(mContext.getString(R.string.expandable_listview_pay_type4));
        }

        String tradeStatus;
        int colorStatus = mContext.getResources().getColor(R.color.textview_textcolor_bill_status1);
        if ("10".equals(bill.transStatus)) {
            //处理中
            tradeStatus = mContext.getString(R.string.expandable_listview_trade_status_nopay);
        } else if ("30".equals(bill.transStatus)) {
            double amt = Double.parseDouble(bill.refundAmt);
            if (amt == 0) {
                //成功的
                colorStatus = mContext.getResources().getColor(R.color.textview_textcolor_bill_status);
                tradeStatus = mContext.getString(R.string.expandable_listview_trade_status_success);
            } else {
                //部分退款的
                tradeStatus = mContext.getString(R.string.expandable_listview_trade_status_partrefd);
            }
        } else if ("40".equals(bill.transStatus)) {
            if ("09".equals(bill.response)) {
                //已关闭
                tradeStatus = mContext.getString(R.string.expandable_listview_trade_status_closed);
            } else {
                //全额退款
                tradeStatus = mContext.getString(R.string.expandable_listview_trade_status_allrefd);
            }
        } else {
            //失败的
            tradeStatus = mContext.getString(R.string.expandable_listview_trade_status_fail);
        }
        childViewHolder.billTradeStatus.setText(tradeStatus);
        childViewHolder.billTradeStatus.setTextColor(colorStatus);

        try {
            BigDecimal txAmt = new BigDecimal(bill.amount);
            String txAmtStr = txAmt.setScale(2, BigDecimal.ROUND_HALF_UP).toString();
            childViewHolder.billTradeAmount.setText("￥" + txAmtStr);//实际支付金额
        } catch (Exception e) {

        }

        try {
            //如果有优惠
            BigDecimal discountAmt = new BigDecimal(bill.couponDiscountAmt);//优惠的金额
            BigDecimal txAmt = new BigDecimal(bill.amount);
            BigDecimal b0 = new BigDecimal("0");
            if (discountAmt.compareTo(b0) > 0) {
                //大于0 说明有优惠金额
                childViewHolder.billOriginTradeAmount.setVisibility(View.VISIBLE);
                childViewHolder.billDiscount.setVisibility(View.VISIBLE);
                String origin = discountAmt.add(txAmt).setScale(2, BigDecimal.ROUND_HALF_UP).toString();
                childViewHolder.billOriginTradeAmount.setText("￥" + origin);
                childViewHolder.billOriginTradeAmount.getPaint().setFlags(Paint.STRIKE_THRU_TEXT_FLAG);
            } else {
                childViewHolder.billOriginTradeAmount.setVisibility(View.INVISIBLE);
                childViewHolder.billDiscount.setVisibility(View.INVISIBLE);
            }
        } catch (Exception e) {
            childViewHolder.billOriginTradeAmount.setVisibility(View.INVISIBLE);
            childViewHolder.billDiscount.setVisibility(View.INVISIBLE);
        }

        childViewHolder.linearLayoutDay.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                quickPayService = ShowMoneyApp.getInstance().getQuickPayService();
                SimpleDateFormat spf = new SimpleDateFormat("yyyyMMddHHmmss");
                SimpleDateFormat spf1 = new SimpleDateFormat("yyyyMMdd");
                SimpleDateFormat spf2 = new SimpleDateFormat("yyyy/MM/dd");
                String date;
                String currentdate;
                try {
                    Date tandeDate = spf.parse(bill.tandeDate);
                    date = spf1.format(tandeDate);
                    currentdate = spf2.format(tandeDate);
                } catch (ParseException e) {
                    date = spf1.format(new Date());
                    currentdate = spf2.format(new Date());
                }
                final String finalCurrentdate = currentdate;
                quickPayService.getSummaryDayAsync(SessonData.loginUser, date, "1", new QuickPayCallbackListener<ServerPacket>() {
                    @Override
                    public void onSuccess(ServerPacket data) {
                        //笔数
                        String countStr = String.valueOf(data.getCount());
                        mHintDialog.setBillCount(countStr);

                        //收入
                        String totalStr = TxamtUtil.getNormal(data.getTotal());
                        if (TextUtils.isEmpty(totalStr)) {
                            totalStr = "0.00";
                        }
                        mHintDialog.setBillTotal(totalStr);

                        mHintDialog.setBillDate(finalCurrentdate);
                        mHintDialog.show();
                    }

                    @Override
                    public void onFailure(QuickPayException ex) {
                        mHintDialog.setTitle(ex.getErrorMsg());
                        mHintDialog.show();
                    }
                });
            }
        });

        childViewHolder.linearLayoutBillItem.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                Intent intent = new Intent(mContext, DetailActivity.class);
                Bundle bundle = new Bundle();
                bundle.putSerializable("TradeBill", bill);
                intent.putExtra("BillBundle", bundle);
                mContext.startActivity(intent);
            }
        });


        return convertView;
    }

    @Override
    public boolean isChildSelectable(int groupPosition, int childPosition) {
        return true;
    }

    public final class GroupViewHolder {
        public ImageView folder;
        public TextView month;
        public TextView year;
        public TextView total;
        public TextView count;
        public LinearLayout leftLinearLayout;
    }

    public final class ChildViewHolder {
        public TextView day;
        public TextView weekday;
        public ImageView paylogo;
        public TextView billTradeDate;
        public TextView billTradeFrom;
        public ImageView billTradeFromImage;
        public TextView billTradeStatus;

        public TextView billTradeAmount;//优惠后的实际支付的金额
        public TextView billOriginTradeAmount;//优惠之前的金额

        public ImageView billDiscount;//显示折扣的一个图片

        public View linearLayoutDay;//左边显示日期，周几的一个线性布局
        public View linearLayoutBillItem;//右边显示详情账单信息的一个线性布局
    }
}
