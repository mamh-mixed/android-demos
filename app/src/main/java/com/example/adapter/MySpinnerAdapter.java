package com.example.adapter;

import java.util.ArrayList;
import java.util.List;

import android.content.Context;
import android.database.Cursor;
import android.database.DataSetObserver;
import android.database.sqlite.SQLiteDatabase;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.SpinnerAdapter;
import android.widget.TextView;

import com.example.eatwhatdemo.R;
import com.example.library.DatabaseHelper;

public class MySpinnerAdapter implements SpinnerAdapter {
    private List<List<String>> items;
    private LayoutInflater inflater;
    private DatabaseHelper dbhelper;

    public MySpinnerAdapter(Context context, List<List<String>> items) {
        super();
        this.items = items;
        this.inflater = LayoutInflater.from(context);
        this.dbhelper = new DatabaseHelper(context);

    }

    public MySpinnerAdapter() {

    }

    public MySpinnerAdapter(Context context) {
        this.inflater = LayoutInflater.from(context);
        this.dbhelper = new DatabaseHelper(context);

        loadData();
    }

    private void loadData() {
        items = new ArrayList<List<String>>();
        SQLiteDatabase db = dbhelper.getReadableDatabase();
        Cursor cursor = db.rawQuery("select * from ew_restaurant", null);
        while (cursor.moveToNext()) {
            int id = cursor.getInt(0);
            String name = cursor.getString(1);
            List<String> sublist = new ArrayList<String>();
            sublist.add(id + "");
            sublist.add(name);
            this.items.add(sublist);
        }
        cursor.close();
        db.close();
    }

    @Override
    public int getCount() {
        return items.size();
    }

    @Override
    public Object getItem(int position) {
        return items.get(position);
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
            convertView = inflater.inflate(R.layout.spinner_item, null);
        }
        TextView idtextView = (TextView) convertView.findViewById(R.id.spinner_item_id_textView);
        idtextView.setText(items.get(position).get(0));

        TextView restaurantnametextView = (TextView) convertView.findViewById(R.id.spinner_item_restaurantname_textView);
        restaurantnametextView.setText(items.get(position).get(1));
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
    public View getDropDownView(int position, View convertView, ViewGroup parent) {
        if (convertView == null) {
            convertView = inflater.inflate(R.layout.spinner_item, null);
        }
        TextView idtextView = (TextView) convertView.findViewById(R.id.spinner_item_id_textView);
        idtextView.setText(items.get(position).get(0));
        TextView restaurantnametextView = (TextView) convertView.findViewById(R.id.spinner_item_restaurantname_textView);
        restaurantnametextView.setText(items.get(position).get(1));
        return convertView;
    }

}
