package com.cardinfolink.yunshouyin.adapter;

import android.annotation.SuppressLint;
import android.content.Context;
import android.content.Intent;
import android.graphics.Color;
import android.os.Bundle;
import android.text.TextUtils;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.BaseAdapter;
import android.widget.Button;
import android.widget.ImageView;
import android.widget.TextView;

import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.activity.DetailActivity;
import com.cardinfolink.yunshouyin.data.TradeBill;
import com.cardinfolink.yunshouyin.view.TransManageView;

import java.text.ParseException;
import java.text.SimpleDateFormat;
import java.util.ArrayList;
import java.util.Date;
import java.util.List;

/**
 * Created by mamh on 15-11-20.
 */
public class BillAdapter extends BaseAdapter {

    private final TransManageView transManageView;
    private List<TradeBill> tradeBillList;
    private Context mContext;

    public BillAdapter(TransManageView transManageView, List<TradeBill> tradeBillList) {
        this.transManageView = transManageView;
        this.tradeBillList = new ArrayList<TradeBill>();
        this.tradeBillList.addAll(tradeBillList);
        this.mContext = transManageView.getContext();
    }

    public void setData(List<TradeBill> list) {
        tradeBillList.clear();
        tradeBillList.addAll(list);
    }

    @Override
    public int getCount() {
        // TODO Auto-generated method stub
        return tradeBillList.size();
    }

    @Override
    public Object getItem(int position) {
        // TODO Auto-generated method stub
        return null;
    }

    @Override
    public long getItemId(int position) {
        // TODO Auto-generated method stub
        return 0;
    }

    @SuppressLint("NewApi")
    @Override
    public View getView(int position, View convertView, ViewGroup parent) {
        ViewHolder holder = null;
        if (convertView == null) {
            holder = new ViewHolder();
            convertView = LayoutInflater.from(mContext).inflate(
                    R.layout.bill_list_item, null);
            holder.paylogo = (ImageView) convertView
                    .findViewById(R.id.paylogo);
            holder.billTradeDate = (TextView) convertView
                    .findViewById(R.id.bill_tradedate);
            holder.billTradeFrom = (TextView) convertView
                    .findViewById(R.id.bill_tradefrom);
            holder.billTradeStatus = (TextView) convertView
                    .findViewById(R.id.bill_tradestatus);
            convertView.setTag(holder);

            holder.billTradeAmount = (TextView) convertView
                    .findViewById(R.id.bill_tradeamount);

            holder.billTradeDeatil = (Button) convertView
                    .findViewById(R.id.bill_tradedetail);

            convertView.setTag(holder);
        } else {
            holder = (ViewHolder) convertView.getTag();
        }

        if (position < 0 || position > tradeBillList.size()) {
            return convertView;
        }

        final TradeBill bill = tradeBillList.get(position);
        if (bill.chcd.equals("WXP")) {
            holder.paylogo.setImageResource(R.drawable.wpay);
        } else {
            holder.paylogo.setImageResource(R.drawable.apay);
        }
        SimpleDateFormat spf1 = new SimpleDateFormat("yyyyMMddHHmmss");
        SimpleDateFormat spf2 = new SimpleDateFormat("yyyy-MM-dd HH:mm:ss");
        try {
            Date tandeDate = spf1.parse(bill.tandeDate);
            holder.billTradeDate.setText(spf2.format(tandeDate));
        } catch (ParseException e) {
            // TODO Auto-generated catch block
            e.printStackTrace();
        }
        String tradeFrom = "PC";
        if (!TextUtils.isEmpty(bill.tradeFrom)) {
            tradeFrom = bill.tradeFrom;
        }

        String busicd = transManageView.getResources().getString(R.string.detail_activity_busicd_pay);
        if (bill.busicd.equals("REFD")) {
            busicd = transManageView.getResources().getString(R.string.detail_activity_busicd_refd);
        }

        holder.billTradeFrom.setText(tradeFrom + busicd);
        String tradeStatus;
        if (bill.response.equals("00")) {
            tradeStatus = transManageView.getResources().getString(R.string.detail_activity_trade_status_success);
            holder.billTradeStatus
                    .setTextColor(Color.parseColor("#888888"));
        } else if (bill.response.equals("09")) {
            tradeStatus = transManageView.getResources().getString(R.string.detail_activity_trade_status_nopay);
            holder.billTradeStatus.setTextColor(Color.RED);
        } else {
            tradeStatus = transManageView.getResources().getString(R.string.detail_activity_trade_status_fail);
            holder.billTradeStatus.setTextColor(Color.RED);
        }
        holder.billTradeStatus.setText(tradeStatus);
        holder.billTradeAmount.setText("￥" + bill.amount);
        holder.billTradeDeatil.setOnClickListener(new View.OnClickListener() {

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

    public final class ViewHolder {
        public ImageView paylogo;
        public TextView billTradeDate;
        public TextView billTradeFrom;
        public TextView billTradeStatus;
        public TextView billTradeAmount;
        public Button billTradeDeatil;
    }

}
