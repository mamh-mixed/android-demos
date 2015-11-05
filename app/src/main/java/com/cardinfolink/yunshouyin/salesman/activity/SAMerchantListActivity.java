package com.cardinfolink.yunshouyin.salesman.activity;

import android.content.Intent;
import android.content.SharedPreferences;
import android.os.Bundle;
import android.support.v4.widget.SwipeRefreshLayout;
import android.text.Editable;
import android.text.TextWatcher;
import android.util.Log;
import android.view.View;
import android.widget.Button;
import android.widget.EditText;
import android.widget.ListView;
import android.widget.TextView;
import android.widget.Toast;

import com.cardinfolink.yunshouyin.salesman.R;
import com.cardinfolink.yunshouyin.salesman.adapter.MerchantListAdapter;
import com.cardinfolink.yunshouyin.salesman.model.SAServerPacket;
import com.cardinfolink.yunshouyin.salesman.model.SessonData;
import com.cardinfolink.yunshouyin.salesman.model.User;
import com.cardinfolink.yunshouyin.salesman.utils.ActivityCollector;
import com.cardinfolink.yunshouyin.salesman.utils.CommunicationListenerV2;
import com.cardinfolink.yunshouyin.salesman.utils.ErrorUtil;
import com.cardinfolink.yunshouyin.salesman.utils.HttpCommunicationUtil;
import com.cardinfolink.yunshouyin.salesman.utils.ParamsUtil;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.Calendar;
import java.util.Collections;
import java.util.Comparator;
import java.util.Date;
import java.util.List;

//TODO: 加入分页下载,搜索API,上拉更多
public class SAMerchantListActivity extends BaseActivity {
    private final String TAG = "SAMerchantListActivity";

    //该地址会被ArrayAdapter所引用,作为数据源,对merchantInfos所做的修改会影响到arrayAdapter
    private List<User> users = new ArrayList<>();
    MerchantListAdapter adapter;
    private SwipeRefreshLayout swipeRefreshLayout;
    private Button btnAddNewMer;
    private EditText searchText;
    private TextView txtMerchantCountThisMonth;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_samerchant_list);
        initLayout();
        setupListView();
    }

    @Override
    protected void onDestroy() {
        super.onDestroy();
        Log.d(TAG, "onDestroy(): will delete cache files");
        SharedPreferences sp = getSharedPreferences("data", MODE_PRIVATE);
        sp.edit().clear().commit();
    }


    /**
     * setup listener
     */
    private void initLayout() {
        txtMerchantCountThisMonth = (TextView)findViewById(R.id.txt_merchantcountthismonth);
        swipeRefreshLayout = (SwipeRefreshLayout) findViewById(R.id.swipe_container);
        swipeRefreshLayout.setOnRefreshListener(new SwipeRefreshLayout.OnRefreshListener() {
            @Override
            public void onRefresh() {
                refreshData();
            }
        });

        btnAddNewMer = (Button) findViewById(R.id.btnAddNewMerchant);
        btnAddNewMer.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                Intent intent = new Intent(mContext, RegisterActivity.class);
                mContext.startActivity(intent);
            }
        });

        //输入关键字快速定位
        searchText = (EditText) findViewById(R.id.mItem_txtSearch);
        searchText.addTextChangedListener(new TextWatcher() {
            @Override
            public void beforeTextChanged(CharSequence s, int start, int count, int after) {

            }

            @Override
            public void onTextChanged(CharSequence s, int start, int before, int count) {
                adapter.getFilter().filter(s);
            }

            @Override
            public void afterTextChanged(Editable s) {

            }
        });
    }

    /**
     * called when refresh data, pull down to refresh
     */
    private void refreshData() {
        //startLoading();
        swipeRefreshLayout.setRefreshing(true);

        // async network call, callbacks
        HttpCommunicationUtil.sendDataToQuickIpayServer(ParamsUtil.getUsers_SA(SessonData.getAccessToken()), new CommunicationListenerV2() {
            @Override
            public void onResult(SAServerPacket serverPacket) {
                //必须保持users的引用,因为adapter使用了
                final List<User> tempUsers = new ArrayList<>();
                if (serverPacket.getUsers() != null){
                    tempUsers.addAll(Arrays.asList(serverPacket.getUsers()));
                }
                //order by create time
                Collections.sort(tempUsers, new Comparator<User>() {
                    @Override
                    public int compare(User lhs, User rhs) {
                        return lhs.getCreateTime().after(rhs.getCreateTime())==true?-1:1;
                    }
                });

                int num =0;
                Date today = new Date();
                Calendar calendar = Calendar.getInstance();
                calendar.setTime(today);
                calendar.set(Calendar.DAY_OF_MONTH, 1);
                calendar.set(Calendar.HOUR, 0);
                calendar.set(Calendar.MINUTE, 0);
                calendar.set(Calendar.SECOND, 0);
                calendar.set(Calendar.MILLISECOND, 0);
                Date firstDayOfMonth = calendar.getTime();

                for (User user: tempUsers){
                    if (user.getCreateTime().after(firstDayOfMonth)){
                        num++;
                    }
                }

                final int finalNum = num;
                runOnUiThread(new Runnable() {
                    @Override
                    public void run() {
                        txtMerchantCountThisMonth.setText(String.format("本月已经发展商户: %d 家", finalNum));
                        users.clear();
                        users.addAll(tempUsers);
                        adapter.refreshDataSource(users);
                        adapter.notifyDataSetChanged();
                        swipeRefreshLayout.setRefreshing(false);
                        Toast.makeText(SAMerchantListActivity.this, "刷新成功", Toast.LENGTH_SHORT).show();
                        //endLoading();
                    }
                });
            }

            @Override
            public void onError(final String error) {
                runOnUiThread(new Runnable() {
                    @Override
                    public void run() {
                        Log.i("opp", "error:" + error);
                        String errorStr = ErrorUtil.getErrorString(error);
                        swipeRefreshLayout.setRefreshing(false);
                        //endLoadingWithError(errorStr);
                        alertError(errorStr);
                        if (error.equals("accessToken_error")) {
                            //关闭所有activity,除了登录框
                            ActivityCollector.goLoginAndFinishRest();
                        }
                    }
                });
            }
        });
    }

    private void setupListView() {
        // currently no data
        adapter = new MerchantListAdapter(this, users);
        ListView listView = (ListView) findViewById(R.id.listViewMerchants);
        listView.setAdapter(adapter);
    }

    @Override
    protected void onResume() {
        super.onResume();
        //回到页面之后,从服务器刷新数据
        Log.d(TAG,"onResume() will refresh data");
        refreshData();
    }
}
