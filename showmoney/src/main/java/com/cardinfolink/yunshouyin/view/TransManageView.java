package com.cardinfolink.yunshouyin.view;

import android.content.Context;
import android.os.AsyncTask;
import android.util.Log;
import android.view.LayoutInflater;
import android.view.View;
import android.widget.ExpandableListView;
import android.widget.LinearLayout;
import android.widget.SimpleExpandableListAdapter;

import com.cardinfolink.yunshouyin.R;
import com.handmark.pulltorefresh.library.PullToRefreshBase;
import com.handmark.pulltorefresh.library.PullToRefreshExpandableListView;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class TransManageView extends LinearLayout {
    private static final String TAG = "TransManageView";
    private Context mContext;
    private PullToRefreshExpandableListView mPullRefreshListView;
    private SimpleExpandableListAdapter mAdapter;

    private String[] mChildStrings = {"Child One", "Child Two", "Child Three", "Child Four", "Child Five", "Child Six"};

    private String[] mGroupStrings = {"Group One", "Group Two", "Group Three"};
    private static final String KEY = "key";
    private List<Map<String, String>> groupData = new ArrayList<Map<String, String>>();
    private List<List<Map<String, String>>> childData = new ArrayList<List<Map<String, String>>>();


    public TransManageView(Context context) {
        super(context);
        mContext = context;
        View contentView = LayoutInflater.from(context).inflate(R.layout.transmanage_view, null);
        LinearLayout.LayoutParams layoutParams = new LayoutParams(LayoutParams.MATCH_PARENT, LayoutParams.MATCH_PARENT);
        contentView.setLayoutParams(layoutParams);
        addView(contentView);

        initLayout();
    }

    private void initLayout() {
        mPullRefreshListView = (PullToRefreshExpandableListView) findViewById(R.id.bill_list_view);

        mPullRefreshListView.setMode(PullToRefreshBase.Mode.BOTH);


        mPullRefreshListView.setOnRefreshListener(new PullToRefreshBase.OnRefreshListener<ExpandableListView>() {
            @Override
            public void onRefresh(PullToRefreshBase<ExpandableListView> refreshView) {
                // Do work to refresh the list here.
                Log.e(TAG, " pull to refersh");
                new GetDataTask().execute();
            }
        });
        for (String group : mGroupStrings) {
            Map<String, String> groupMap1 = new HashMap<String, String>();
            groupData.add(groupMap1);
            groupMap1.put(KEY, group);

            List<Map<String, String>> childList = new ArrayList<Map<String, String>>();
            for (String string : mChildStrings) {
                Map<String, String> childMap = new HashMap<String, String>();
                childList.add(childMap);
                childMap.put(KEY, string);
            }
            childData.add(childList);
        }

        mAdapter = new SimpleExpandableListAdapter(mContext,
                groupData,
                android.R.layout.simple_expandable_list_item_1,
                new String[]{KEY}, new int[]{android.R.id.text1},
                childData,
                android.R.layout.simple_expandable_list_item_2, new String[]{KEY}, new int[]{android.R.id.text1});

        ExpandableListView ActualView = mPullRefreshListView.getRefreshableView();
        ActualView.setAdapter(mAdapter);
    }


    public void refresh() {


    }
    private class GetDataTask extends AsyncTask<Void, Void, String[]> {

        @Override
        protected String[] doInBackground(Void... params) {
            // Simulates a background job.
            try {
                Thread.sleep(2000);
            } catch (InterruptedException e) {
            }
            return mChildStrings;
        }

        @Override
        protected void onPostExecute(String[] result) {
            Map<String, String> newMap = new HashMap<String, String>();
            newMap.put(KEY, "Added after refresh...");
            groupData.add(newMap);

            List<Map<String, String>> childList = new ArrayList<Map<String, String>>();
            for (String string : mChildStrings) {
                Map<String, String> childMap = new HashMap<String, String>();
                childMap.put(KEY, string);
                childList.add(childMap);
            }
            childData.add(childList);

            mAdapter.notifyDataSetChanged();

            // Call onRefreshComplete when the list has been refreshed.
            mPullRefreshListView.onRefreshComplete();

            super.onPostExecute(result);
        }
    }

}
