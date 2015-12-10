package com.example.adapter;

import java.util.ArrayList;
import java.util.List;

import com.example.eatwhatdemo.R;

import android.content.Context;
import android.database.DataSetObserver;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.ListAdapter;
import android.widget.TextView;

public class MyListAdapter implements ListAdapter {
    private List<String> items;
    private LayoutInflater inflater;

    public MyListAdapter() {

    }

    public MyListAdapter(Context context) {
        inflater = LayoutInflater.from(context);
        //����Դ�����л�ȡ
        String[] strArray = context.getResources().getStringArray(R.array.adddelupdatesearch);

        items = new ArrayList<String>();
        for (String str : strArray) {
            items.add(str);
        }

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
            convertView = inflater.inflate(R.layout.cornerlistview_item, null);
        }
        TextView textView = (TextView) convertView.findViewById(R.id.cornerlist_item_textView);
        textView.setText(items.get(position));
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
