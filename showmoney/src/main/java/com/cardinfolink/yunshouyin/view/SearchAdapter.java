package com.cardinfolink.yunshouyin.view;

import android.content.Context;
import android.text.TextUtils;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.BaseAdapter;
import android.widget.Filter;
import android.widget.Filterable;
import android.widget.TextView;

import com.cardinfolink.yunshouyin.R;

import java.util.ArrayList;
import java.util.List;


public class SearchAdapter extends BaseAdapter implements Filterable {
    private List<String> mOriginData;
    private List<String> mNewData;
    private Context mContext;


    public SearchAdapter(Context context, List<String> list) {
        mContext = context;
        mOriginData = new ArrayList<String>();
        mOriginData.addAll(list);
        mNewData = new ArrayList<String>();
    }

    public void setData(List<String> list) {
        mOriginData = new ArrayList<String>();
        mOriginData.addAll(list);
    }

    @Override
    public int getCount() {
        if (mNewData != null) {
            return mNewData.size();
        }
        return 0;
    }

    @Override
    public Object getItem(int position) {
        if (mNewData != null) {
            return mNewData.get(position);
        }
        return "";
    }

    @Override
    public long getItemId(int position) {
        return position;
    }

    @Override
    public View getView(int position, View convertView, ViewGroup parent) {
        ViewHolder holder = null;
        if (convertView == null) {
            holder = new ViewHolder();
            convertView = LayoutInflater.from(mContext).inflate(R.layout.search_item, null);
            holder.text = (TextView) convertView.findViewById(R.id.text1);
            convertView.setTag(holder);
        } else {
            holder = (ViewHolder) convertView.getTag();
        }
        holder.text.setText(mNewData.get(position));

        return convertView;
    }

    @Override
    public Filter getFilter() {
        return new SearchFilter();
    }

    private final class ViewHolder {
        public TextView text;
    }

    private class SearchFilter extends Filter {

        @Override
        protected FilterResults performFiltering(CharSequence constraint) {
            FilterResults results = new FilterResults();
            List<String> data = new ArrayList<String>();
            if (!TextUtils.isEmpty(constraint)) {
                for (String name : mOriginData) {
                    if (name.contains(constraint)) {
                        data.add(name);
                    }
                }
            }
            results.values = data;
            results.count = data.size();
            return results;
        }

        @Override
        protected void publishResults(CharSequence constraint, FilterResults results) {
            mNewData = (List<String>) results.values;
            if (results.count > 0) {
                notifyDataSetChanged();
            } else {
                notifyDataSetInvalidated();
            }
        }
    }
}