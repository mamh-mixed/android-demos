package com.cardinfolink.yunshouyin.adapter;

import android.content.Context;
import android.content.Intent;
import android.graphics.Color;
import android.os.Bundle;
import android.text.TextUtils;
import android.util.Log;
import android.view.View;
import android.view.ViewGroup;
import android.widget.BaseExpandableListAdapter;
import android.widget.ImageView;
import android.widget.TextView;

import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.activity.DetailActivity;
import com.cardinfolink.yunshouyin.data.MonthBill;
import com.cardinfolink.yunshouyin.data.TradeBill;

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


    public BillExpandableListAdapter(Context context, List<MonthBill> groupData, List<List<TradeBill>> childrenData) {
        this.mContext = context;
        this.groupData = groupData;
        this.childrenData = childrenData;
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

            convertView.setTag(groupViewHolder);
        } else {
            groupViewHolder = (GroupViewHolder) convertView.getTag();
        }
        //设置一下月份
        groupViewHolder.month.setText(groupData.get(groupPosition).getCurrentMonth());
        groupViewHolder.year.setText(groupData.get(groupPosition).getCurrentYear());
        groupViewHolder.total.setText(groupData.get(groupPosition).getTotal());
        groupViewHolder.count.setText("" + groupData.get(groupPosition).getTotalRecord());

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
        SimpleDateFormat spf4 = new SimpleDateFormat("EEEE");
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
        if ("10".equals(bill.transStatus)) {
            //处理中
            tradeStatus = mContext.getString(R.string.expandable_listview_trade_status_nopay);
            childViewHolder.billTradeStatus.setTextColor(Color.parseColor("#888888"));
        } else if ("30".equals(bill.transStatus)) {
            double amt = Double.parseDouble(bill.refundAmt);
            if (amt == 0) {
                //成功的
                tradeStatus = mContext.getString(R.string.expandable_listview_trade_status_success);
                childViewHolder.billTradeStatus.setTextColor(Color.parseColor("#888888"));
            } else {
                //部分退款的
                tradeStatus = mContext.getString(R.string.expandable_listview_trade_status_partrefd);
                childViewHolder.billTradeStatus.setTextColor(Color.parseColor("#888888"));
            }
        } else if ("40".equals(bill.transStatus)) {
            if ("09".equals(bill.response)) {
                //已关闭
                tradeStatus = mContext.getString(R.string.expandable_listview_trade_status_closed);
                childViewHolder.billTradeStatus.setTextColor(Color.parseColor("#888888"));
            } else {
                //全额退款
                tradeStatus = mContext.getString(R.string.expandable_listview_trade_status_partrefd);
                childViewHolder.billTradeStatus.setTextColor(Color.parseColor("#888888"));
            }
        } else {
            //失败的
            tradeStatus = mContext.getString(R.string.expandable_listview_trade_status_fail);
            childViewHolder.billTradeStatus.setTextColor(Color.RED);
        }
        childViewHolder.billTradeStatus.setText(tradeStatus);
        childViewHolder.billTradeAmount.setText("￥" + bill.amount);


        childViewHolder.linearLayoutDay.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
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
    }

    public final class ChildViewHolder {
        public TextView day;
        public TextView weekday;
        public ImageView paylogo;
        public TextView billTradeDate;
        public TextView billTradeFrom;
        public ImageView billTradeFromImage;
        public TextView billTradeStatus;
        public TextView billTradeAmount;
        public View linearLayoutDay;//左边显示日期，周几的一个线性布局
        public View linearLayoutBillItem;//右边显示详情账单信息的一个线性布局
    }
}
