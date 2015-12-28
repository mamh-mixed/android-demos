package com.cardinfolink.yunshouyin.adapter;

import android.content.Context;
import android.graphics.Color;
import android.text.TextUtils;
import android.view.View;
import android.view.ViewGroup;
import android.widget.BaseExpandableListAdapter;
import android.widget.ImageView;
import android.widget.ListView;
import android.widget.TextView;

import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.data.MonthBill;
import com.cardinfolink.yunshouyin.data.TradeBill;

import java.text.ParseException;
import java.text.SimpleDateFormat;
import java.util.Date;
import java.util.List;
import java.util.Map;

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
            convertView = View.inflate(mContext, R.layout.expandablelistview_group, null);

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
        int size = childrenData.get(groupPosition).size();
        groupViewHolder.count.setText(size + "");

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
            convertView = View.inflate(mContext, R.layout.expandablelistview_child, null);
            childViewHolder.day = (TextView) convertView.findViewById(R.id.tv_day);
            childViewHolder.weekday = (TextView) convertView.findViewById(R.id.tv_weekday);

            childViewHolder.paylogo = (ImageView) convertView.findViewById(R.id.paylogo);
            childViewHolder.billTradeDate = (TextView) convertView.findViewById(R.id.bill_tradedate);
            childViewHolder.billTradeFrom = (TextView) convertView.findViewById(R.id.bill_tradefrom);
            childViewHolder.billTradeStatus = (TextView) convertView.findViewById(R.id.bill_tradestatus);
            childViewHolder.billTradeAmount = (TextView) convertView.findViewById(R.id.bill_tradeamount);

            convertView.setTag(childViewHolder);
        } else {
            childViewHolder = (ChildViewHolder) convertView.getTag();
        }

        TradeBill bill = childrenData.get(groupPosition).get(childPosition);

        if (!TextUtils.isEmpty(bill.chcd)) {
            //有chcd渠道的话,这里设置不同渠道的图片
            if (bill.chcd.equals("WXP")) {
                childViewHolder.paylogo.setImageResource(R.drawable.wpay);
            } else {
                childViewHolder.paylogo.setImageResource(R.drawable.apay);
            }
        } else {
            childViewHolder.paylogo.setImageDrawable(null);
        }
        SimpleDateFormat spf1 = new SimpleDateFormat("yyyyMMddHHmmss");
        SimpleDateFormat spf2 = new SimpleDateFormat("MM-dd HH:mm:ss");
        SimpleDateFormat spf3 = new SimpleDateFormat("dd");
        try {
            Date tandeDate = spf1.parse(bill.tandeDate);
            childViewHolder.billTradeDate.setText(spf2.format(tandeDate));
            childViewHolder.day.setText(spf3.format(tandeDate));
        } catch (ParseException e) {
            e.printStackTrace();
        }
        String tradeFrom = "PC";
        if (!TextUtils.isEmpty(bill.tradeFrom)) {
            tradeFrom = bill.tradeFrom;
        }

        String busicd = mContext.getResources().getString(R.string.detail_activity_busicd_pay);
        if (bill.busicd.equals("REFD")) {
            //退款的
            busicd = mContext.getResources().getString(R.string.detail_activity_busicd_refd);
        } else if ("CANC".equals(bill.busicd)) {
            //取消订单
            busicd = mContext.getResources().getString(R.string.detail_activity_busicd_canc);
        }

        childViewHolder.billTradeFrom.setText(tradeFrom + busicd);
        String tradeStatus;
        if (bill.response.equals("00")) {
            tradeStatus = mContext.getResources().getString(R.string.detail_activity_trade_status_success);
            childViewHolder.billTradeStatus.setTextColor(Color.parseColor("#888888"));
        } else if (bill.response.equals("09")) {
            tradeStatus = mContext.getResources().getString(R.string.detail_activity_trade_status_nopay);
            childViewHolder.billTradeStatus.setTextColor(Color.RED);
        } else {
            tradeStatus = mContext.getResources().getString(R.string.detail_activity_trade_status_fail);
            childViewHolder.billTradeStatus.setTextColor(Color.RED);
        }
        childViewHolder.billTradeStatus.setText(tradeStatus);
        childViewHolder.billTradeAmount.setText("￥" + bill.amount);

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
        public TextView billTradeStatus;
        public TextView billTradeAmount;
    }
}
