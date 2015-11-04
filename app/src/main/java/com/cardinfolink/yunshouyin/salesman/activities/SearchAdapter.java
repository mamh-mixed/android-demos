package com.cardinfolink.yunshouyin.salesman.activities;

import android.content.Context;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.BaseAdapter;
import android.widget.Filter;
import android.widget.Filterable;
import android.widget.TextView;

import com.cardinfolink.yunshouyin.salesman.R;

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
        // TODO Auto-generated method stub
        return mNewData.size();
    }

    @Override
    public Object getItem(int position) {
        // TODO Auto-generated method stub
        return mNewData.get(position);
    }

    @Override
    public long getItemId(int position) {
        // TODO Auto-generated method stub
        return position;
    }

    @Override
    public View getView(int position, View convertView, ViewGroup parent) {
        ViewHolder holder = null;
        if (convertView == null) {
            holder = new ViewHolder();
            convertView = LayoutInflater.from(mContext).inflate(
                    R.layout.search_item, null);
            holder.text = (TextView) convertView.findViewById(R.id.text1);
            convertView.setTag(holder);

        } else {
            holder = (ViewHolder) convertView.getTag();
        }
        holder.text.setText(mNewData.get(position));


        return convertView;
    }


    public final class ViewHolder {

        public TextView text;

    }


    @Override
    public Filter getFilter() {
        // TODO Auto-generated method stub
        return new SearchFilter();
    }


    class SearchFilter extends Filter {

        @Override
        protected FilterResults performFiltering(CharSequence constraint) {
            FilterResults results = new FilterResults();
            List<String> data = new ArrayList<>();
            if (constraint!=null && !constraint.equals("")) {
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
        protected void publishResults(CharSequence constraint,
                                      FilterResults results) {
            mNewData = (List<String>) results.values;
            if (results.count > 0) {
                notifyDataSetChanged();
            } else {
                notifyDataSetInvalidated();
            }

        }

    }
}