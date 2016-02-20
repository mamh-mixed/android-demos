package com.example.adapter;

import java.util.ArrayList;

import android.content.Context;
import android.database.DataSetObserver;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.ListAdapter;
import android.widget.TextView;

import com.example.model.AppInfo;
import com.example.random.R;

public class AppInfoListAdapter implements ListAdapter {
    private ArrayList<AppInfo> items;
    private LayoutInflater inflater;
    private Context context;

    public AppInfoListAdapter() {
    }

    public AppInfoListAdapter(Context context) {
        this.context = context;
        this.inflater = LayoutInflater.from(context);
        this.items = new ArrayList<AppInfo>();
    }

    public AppInfoListAdapter(Context context, ArrayList<AppInfo> items) {
        this.context = context;
        this.inflater = LayoutInflater.from(context);
        this.items = items;
    }

    @Override
    public int getCount() {
        return this.items.size();
    }

    @Override
    public Object getItem(int position) {
        return this.items.get(position);
    }

    @Override
    public long getItemId(int position) {
        return position;
    }

    @Override
    public int getItemViewType(int position) {
        return 0;
    }

    @Override
    public View getView(int position, View convertView, ViewGroup parent) {
        if (convertView == null) {
            convertView = inflater.inflate(R.layout.list_item, null);
        }

        TextView appNameTextView = (TextView) convertView.findViewById(R.id.list_item_appname_textView);
        TextView launchTextView = (TextView) convertView.findViewById(R.id.list_item_launch_textView);
        TextView killTextView = (TextView) convertView.findViewById(R.id.list_item_kill_textView);
        appNameTextView.setText(items.get(position).getAppName());
        launchTextView.setText(context.getString(R.string.list_item_launch) + items.get(position).getLaunchTimes());
        killTextView.setText(context.getString(R.string.list_item_kill) + items.get(position).getKillTimes());
        return convertView;
    }

    @Override
    public int getViewTypeCount() {
        return 1;
    }

    @Override
    public boolean hasStableIds() {
        return false;
    }

    @Override
    public boolean isEmpty() {
        return false;
    }

    @Override
    public void registerDataSetObserver(DataSetObserver arg0) {

    }

    @Override
    public void unregisterDataSetObserver(DataSetObserver arg0) {

    }

    @Override
    public boolean areAllItemsEnabled() {
        return true;
    }

    @Override
    public boolean isEnabled(int position) {
        return true;
    }

}