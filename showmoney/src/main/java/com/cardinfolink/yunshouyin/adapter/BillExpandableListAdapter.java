package com.cardinfolink.yunshouyin.adapter;

import android.content.Context;
import android.util.Log;
import android.view.View;
import android.view.ViewGroup;
import android.widget.ArrayAdapter;
import android.widget.BaseExpandableListAdapter;
import android.widget.ListAdapter;
import android.widget.ListView;
import android.widget.TextView;

import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.data.MonthBill;
import com.cardinfolink.yunshouyin.data.TradeBill;

import java.util.List;
import java.util.Map;
import java.util.Objects;
import java.util.Set;

/**
 * Created by mamh on 15-12-26.
 */
public class BillExpandableListAdapter extends BaseExpandableListAdapter {
    private String TAG = "BillExpandableListAdapter";

    private List<MonthBill> groupData;
    private List<Map<String, List<TradeBill>>> childrenData;
    private Context mContext;


    public BillExpandableListAdapter(Context context, List<MonthBill> groupData, List<Map<String, List<TradeBill>>> childrenData) {
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
        return childrenData.get(groupPosition).get(String.valueOf(childPosition + 1));
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
        if (convertView == null) {
            convertView = View.inflate(mContext, R.layout.expandablelistview_group, null);
        }

        //设置一下月份
        TextView month = (TextView) convertView.findViewById(R.id.tv_month);
        month.setText(groupData.get(groupPosition).getCurrentMonth());

        TextView year = (TextView) convertView.findViewById(R.id.tv_year);
        year.setText(groupData.get(groupPosition).getCurrentYear());

        TextView total = (TextView) convertView.findViewById(R.id.tv_total);
        total.setText(groupData.get(groupPosition).getTotal());

        TextView count = (TextView) convertView.findViewById(R.id.tv_count);
        count.setText(groupData.get(groupPosition).getCount() + "");

        return convertView;
    }

    @Override
    public View getChildView(int groupPosition, int childPosition, boolean isLastChild, View convertView, ViewGroup parent) {
        if (convertView == null) {
            convertView = View.inflate(mContext, R.layout.expandablelistview_child, null);
        }

        Map<String, List<TradeBill>> tradeBillMap = childrenData.get(groupPosition);
        Object[] keyArray = tradeBillMap.keySet().toArray();
        String dayStr = keyArray[childPosition].toString();

        TextView day = (TextView) convertView.findViewById(R.id.tv_day);
        day.setText(dayStr);

        ListView listView = (ListView) convertView.findViewById(R.id.child_list_view);
        listView.setAdapter(new BillAdapter(mContext, tradeBillMap.get(dayStr)));

        return convertView;
    }

    @Override
    public boolean isChildSelectable(int groupPosition, int childPosition) {
        return true;
    }



}
