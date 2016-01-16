package com.cardinfolink.yunshouyin.adapter;

import android.content.Context;
import android.content.Intent;
import android.graphics.Color;
import android.os.Bundle;
import android.text.TextUtils;
import android.view.View;
import android.view.ViewGroup;
import android.widget.BaseExpandableListAdapter;
import android.widget.ImageView;
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
public class TicketExpandableListAdapter extends BaseExpandableListAdapter {
    private String TAG = "TicketExpandableListAdapter";

    private List<MonthBill> groupData;
    private List<List<TradeBill>> childrenData;
    private Context mContext;

    private QuickPayService quickPayService;

    private HintBillDialog mHintDialog;

    public TicketExpandableListAdapter(Context context, List<MonthBill> groupData, List<List<TradeBill>> childrenData) {
        this.mContext = context;
        this.groupData = groupData;
        this.childrenData = childrenData;
    }

    public void setHintDialog(View view) {
        this.mHintDialog = new HintBillDialog(mContext, view);
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
            convertView = View.inflate(mContext, R.layout.ticket_expandablelistview_group, null);

            groupViewHolder.month = (TextView) convertView.findViewById(R.id.tv_month);
            groupViewHolder.year = (TextView) convertView.findViewById(R.id.tv_year);
            groupViewHolder.count = (TextView) convertView.findViewById(R.id.tv_count);
            groupViewHolder.total = (TextView) convertView.findViewById(R.id.tv_total);
            groupViewHolder.folder = (ImageView) convertView.findViewById(R.id.iv_fold);

            convertView.setTag(groupViewHolder);
        } else {
            groupViewHolder = (GroupViewHolder) convertView.getTag();
        }
        //设置一下月份
        groupViewHolder.month.setText(groupData.get(groupPosition).getCurrentMonth());
        groupViewHolder.year.setText(groupData.get(groupPosition).getCurrentYear());


        int count = 0;
        try {
            count = childrenData.get(groupPosition).size();
        } catch (Exception e) {
            count = 0;
        }
        groupViewHolder.count.setText(String.valueOf(count));



        String totalStr = groupData.get(groupPosition).getTotal();
        if (TextUtils.isEmpty(totalStr)) {
            totalStr = "0.00";
        } else {
            try {
                BigDecimal totalBg = new BigDecimal(totalStr);
                totalStr = totalBg.setScale(2, BigDecimal.ROUND_HALF_UP).toString();
            } catch (Exception e) {
                totalStr = "0.0";
            }
        }
        groupViewHolder.total.setText(totalStr);


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
            convertView = View.inflate(mContext, R.layout.ticket_expandablelistview_child, null);
            childViewHolder.linearLayoutDay = convertView.findViewById(R.id.ll_day);//子条目，左边显示日期，周几的那个线性布局
            childViewHolder.linearLayoutBillItem = convertView.findViewById(R.id.ll_coupon_item);

            childViewHolder.day = (TextView) convertView.findViewById(R.id.tv_day);
            childViewHolder.weekday = (TextView) convertView.findViewById(R.id.tv_weekday);

            childViewHolder.couponType = (ImageView) convertView.findViewById(R.id.coupon_type);
            childViewHolder.couponTradeDate = (TextView) convertView.findViewById(R.id.coupon_tradedate);
            childViewHolder.couponTradeFrom = (TextView) convertView.findViewById(R.id.coupon_tv_tradefrom);
            childViewHolder.couponTradeFromImage = (ImageView) convertView.findViewById(R.id.coupon_iv_tradefrom);
            childViewHolder.couponTradeStatus = (TextView) convertView.findViewById(R.id.coupon_tradestatus);

            childViewHolder.couponeName = (TextView) convertView.findViewById(R.id.coupon_name);
            childViewHolder.couponeChannel = (TextView) convertView.findViewById(R.id.coupon_channel);
            convertView.setTag(childViewHolder);
        } else {
            childViewHolder = (ChildViewHolder) convertView.getTag();
        }

        //从list中根据位置获取到相应的bill项
        final TradeBill bill = childrenData.get(groupPosition).get(childPosition);

        if (!TextUtils.isEmpty(bill.couponType)) {
            //有chcd渠道的话,这里设置不同渠道的图片
            if ("1".equals(bill.couponType)) {
                childViewHolder.couponType.setImageResource(R.drawable.bill_subtract);
            } else if ("2".equals(bill.couponType)) {
                childViewHolder.couponType.setImageResource(R.drawable.bill_exchang);
            } else if ("3".equals(bill.couponType)) {
                childViewHolder.couponType.setImageResource(R.drawable.bill_discount1);
            } else {
                childViewHolder.couponType.setImageDrawable(null);
            }
        } else {
            childViewHolder.couponType.setImageDrawable(null);
        }

        SimpleDateFormat spf1 = new SimpleDateFormat("yyyyMMddHHmmss");
        SimpleDateFormat spf2 = new SimpleDateFormat("HH:mm:ss");
        SimpleDateFormat spf3 = new SimpleDateFormat("dd");
        SimpleDateFormat spf4 = new SimpleDateFormat("E");
        try {
            Date tandeDate = spf1.parse(bill.tandeDate);
            childViewHolder.couponTradeDate.setText(spf2.format(tandeDate));
            childViewHolder.day.setText(spf3.format(tandeDate));
            childViewHolder.weekday.setText(spf4.format(tandeDate));
        } catch (ParseException e) {
            e.printStackTrace();
        }

        if ("android".equals(bill.tradeFrom) || "ios".equals(bill.tradeFrom)) {
            childViewHolder.couponTradeFromImage.setImageResource(R.drawable.bill_phone);
            childViewHolder.couponTradeFrom.setText(mContext.getString(R.string.expandable_listview_pay_type1));
        } else if ("wap".equals(bill.tradeFrom)) {
            childViewHolder.couponTradeFromImage.setImageResource(R.drawable.bill_web);
            childViewHolder.couponTradeFrom.setText(mContext.getString(R.string.expandable_listview_pay_type2));
        } else if ("PC".equals(bill.tradeFrom)) {
            childViewHolder.couponTradeFromImage.setImageResource(R.drawable.bill_pc);
            childViewHolder.couponTradeFrom.setText(mContext.getString(R.string.expandable_listview_pay_type3));
        } else {
            childViewHolder.couponTradeFromImage.setImageResource(R.drawable.bill_else);
            childViewHolder.couponTradeFrom.setText(mContext.getString(R.string.expandable_listview_pay_type4));
        }

        String tradeStatus;
        if ("00".equals(bill.response)) {
            tradeStatus = mContext.getString(R.string.expandable_listview_trade_status_success);
        } else {
            //失败的,这里估计不会有失败的吧
            tradeStatus = mContext.getString(R.string.expandable_listview_trade_status_fail);
        }
        childViewHolder.couponTradeStatus.setText(tradeStatus);

        childViewHolder.couponeName.setText(bill.couponName);
        childViewHolder.couponeChannel.setText(bill.couponChannel);

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
                quickPayService.getSummaryDayAsync(SessonData.loginUser, date, "2", new QuickPayCallbackListener<ServerPacket>() {
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
                bill.billType = TradeBill.COUPON_TYPE;
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

        public ImageView couponType;
        public TextView couponTradeDate;//日期
        public TextView couponTradeFrom;//交易来源
        public ImageView couponTradeFromImage;//交易来源显示的图片
        public TextView couponTradeStatus;//状态

        public TextView couponeName;
        public TextView couponeChannel;

        public View linearLayoutDay;//左边显示日期，周几的一个线性布局
        public View linearLayoutBillItem;//右边显示详情账单信息的一个线性布局
    }
}
