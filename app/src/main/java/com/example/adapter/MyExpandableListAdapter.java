package com.example.adapter;

import java.util.List;

import android.content.Context;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.BaseExpandableListAdapter;
import android.widget.TextView;

import com.example.eatwhatdemo.R;

public class MyExpandableListAdapter extends BaseExpandableListAdapter {
    private List<List<String>> groupData;
    private List<List<String>> childrenData;
    private LayoutInflater inflater;

    public MyExpandableListAdapter(Context context, List<List<String>> groupData, List<List<String>> childrenData) {
        super();
        this.groupData = groupData;
        this.childrenData = childrenData;
        this.inflater = LayoutInflater.from(context);
    }


    @Override
    public Object getChild(int groupPosition, int childPosition) {
        return childrenData.get(groupPosition).get(childPosition);
    }

    @Override
    public long getChildId(int groupPosition, int childPosition) {
        return 0;
    }

    @Override
    public View getChildView(int groupPosition, int childPosition, boolean isLastChild, View convertView, ViewGroup parent) {
        if (convertView == null) {
            convertView = inflater.inflate(R.layout.expandablelistview_children, null);
        }
        TextView textView = (TextView) convertView.findViewById(R.id.expandablelistview_child_restaurantid_TextView);
        textView.setText(getChild(groupPosition, childPosition).toString());
        return convertView;
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
    public int getGroupCount() {
        return groupData.size();
    }

    @Override
    public long getGroupId(int groupPosition) {
        return 0;
    }

    @Override
    public View getGroupView(int groupPosition, boolean isExpanded, View convertView, ViewGroup parent) {
        if (convertView == null) {
            convertView = inflater.inflate(R.layout.expandablelistview_groups, null);
        }

        @SuppressWarnings("unchecked")
        List<String> groupList = (List<String>) getGroup(groupPosition);

        TextView idtextView = (TextView) convertView.findViewById(R.id.expandablelistview_group_id_TextView);
        idtextView.setText(groupList.get(0).toString());

        TextView nametextView = (TextView) convertView.findViewById(R.id.expandablelistview_group_name_TextView);
        nametextView.setText(groupList.get(1).toString());
        return convertView;
    }

    @Override
    public boolean hasStableIds() {
        return false;
    }

    @Override
    public boolean isChildSelectable(int groupPosition, int childPosition) {
        return false;
    }

}// end class
