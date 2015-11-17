package com.cardinfolink.yunshouyin.salesman.activity;

import android.app.AlertDialog;
import android.content.DialogInterface;
import android.content.Intent;
import android.content.SharedPreferences;
import android.graphics.Color;
import android.graphics.drawable.ColorDrawable;
import android.os.Bundle;
import android.support.v4.widget.SwipeRefreshLayout;
import android.text.Editable;
import android.text.TextWatcher;
import android.util.Log;
import android.util.TypedValue;
import android.view.View;
import android.widget.Button;
import android.widget.EditText;
import android.widget.ListView;
import android.widget.TextView;
import android.widget.Toast;

import com.cardinfolink.yunshouyin.salesman.R;
import com.cardinfolink.yunshouyin.salesman.adapter.MerchantListAdapter;
import com.cardinfolink.yunshouyin.salesman.api.QuickPayException;
import com.cardinfolink.yunshouyin.salesman.core.QuickPayCallbackListener;
import com.cardinfolink.yunshouyin.salesman.model.User;
import com.cardinfolink.yunshouyin.salesman.swipeview.OnMenuItemClickListener;
import com.cardinfolink.yunshouyin.salesman.swipeview.SwipeMenu;
import com.cardinfolink.yunshouyin.salesman.swipeview.SwipeMenuCreator;
import com.cardinfolink.yunshouyin.salesman.swipeview.SwipeMenuItem;
import com.cardinfolink.yunshouyin.salesman.swipeview.SwipeMenuListView;
import com.cardinfolink.yunshouyin.salesman.utils.ActivityCollector;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.Calendar;
import java.util.Collections;
import java.util.Comparator;
import java.util.Date;
import java.util.List;

//TODO: 加入分页下载,搜索API,上拉更多
public class MerchantListActivity extends BaseActivity {
    private final String TAG = "MerchantListActivity";
    MerchantListAdapter adapter;
    //该地址会被ArrayAdapter所引用,作为数据源,对merchantInfos所做的修改会影响到arrayAdapter
    private List<User> users = new ArrayList<>();
    private SwipeRefreshLayout swipeRefreshLayout;
    private Button btnAddNewMer;
    private EditText searchText;
    private TextView txtMerchantCountThisMonth;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_merchant_list);
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
        txtMerchantCountThisMonth = (TextView) findViewById(R.id.txt_merchantcountthismonth);
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
                int step = mSharedPreferences.getInt("register_step_finish", 0);
                final AlertDialog.Builder builder = new AlertDialog.Builder(MerchantListActivity.this);
                builder.setNegativeButton("否", new DialogInterface.OnClickListener() {
                    @Override
                    public void onClick(DialogInterface dialog, int which) {
                        SharedPreferences.Editor editor = mSharedPreferences.edit();
                        editor.putInt("register_step_finish", 0);
                        editor.commit();
                        intentToActivity(RegisterActivity.class);
                    }
                });
                switch (step) {
                    case 1:
                        //有未完成的注册步骤,这里是第一步完成了
                        builder.setMessage("有未完成的注册是否继续？");
                        builder.setPositiveButton("继续", new DialogInterface.OnClickListener() {
                            @Override
                            public void onClick(DialogInterface dialog, int which) {
                                Intent intent = new Intent(MerchantListActivity.this, RegisterNextActivity.class);
                                startActivity(intent);
                            }
                        });
                        builder.show();
                        break;
                    case 2:
                        //有未完成的注册步骤,这里是第2步完成了
                        builder.setMessage("有未完成的注册是否继续？");
                        builder.setPositiveButton("继续", new DialogInterface.OnClickListener() {
                            @Override
                            public void onClick(DialogInterface dialog, int which) {
                                Intent intent = new Intent(MerchantListActivity.this, RegisterStep3Activity.class);
                                startActivity(intent);
                            }
                        });
                        builder.show();
                        break;
                    default:
                        intentToActivity(RegisterActivity.class);
                        break;
                }
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
        quickPayService.getUsersAsync(new QuickPayCallbackListener<User[]>() {
            @Override
            public void onSuccess(User[] data) {
                //必须保持users的引用,因为adapter使用了
                final List<User> tempUsers = new ArrayList<>();
                if (data != null) {
                    tempUsers.addAll(Arrays.asList(data));
                }
                //order by create time
                Collections.sort(tempUsers, new Comparator<User>() {
                    @Override
                    public int compare(User lhs, User rhs) {
                        return lhs.getCreateTime().after(rhs.getCreateTime()) == true ? -1 : 1;
                    }
                });

                int num = 0;
                Date today = new Date();
                Calendar calendar = Calendar.getInstance();
                calendar.setTime(today);
                calendar.set(Calendar.DAY_OF_MONTH, 1);
                calendar.set(Calendar.HOUR, 0);
                calendar.set(Calendar.MINUTE, 0);
                calendar.set(Calendar.SECOND, 0);
                calendar.set(Calendar.MILLISECOND, 0);
                Date firstDayOfMonth = calendar.getTime();

                for (User user : tempUsers) {
                    if (user.getCreateTime().after(firstDayOfMonth)) {
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
                        Toast.makeText(MerchantListActivity.this, "刷新成功", Toast.LENGTH_SHORT).show();
                        //endLoading();
                    }
                });
            }

            @Override
            public void onFailure(final QuickPayException ex) {
                runOnUiThread(new Runnable() {
                    @Override
                    public void run() {
                        String errorStr = ex.getErrorMsg();
                        swipeRefreshLayout.setRefreshing(false);
                        alertError(errorStr);
                        if (ex.getErrorCode().equals(QuickPayException.ACCESSTOKEN_NOT_FOUND)) {
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
        SwipeMenuListView listView = (SwipeMenuListView) findViewById(R.id.listViewMerchants);
        listView.setAdapter(adapter);

        // step 1. create a MenuCreator
        SwipeMenuCreator creator = new MerchantSwipeMenuCreator();
        // set creator
        listView.setMenuCreator(creator);

        listView.setOnMenuItemClickListener(new OnMenuItemClickListener() {
            @Override
            public void onMenuItemClick(int position, SwipeMenu menu, int index) {
                switch (index) {
                    case 0:
                        User user = users.get(position);
                        Log.d(TAG,"will update user = " + user);
                        break;
                    case 1:
                        Log.d(TAG, "will delete " + position);
                        users.remove(position);
                        adapter.notifyDataSetChanged();
                        break;
                }
            }
        });
    }


    @Override
    protected void onResume() {
        super.onResume();
        //回到页面之后,从服务器刷新数据
        Log.d(TAG, "onResume() will refresh data");
        refreshData();
    }


    private class MerchantSwipeMenuCreator implements SwipeMenuCreator {
        private int dp2px(int dp) {
            return (int) TypedValue.applyDimension(TypedValue.COMPLEX_UNIT_DIP, dp,
                    getResources().getDisplayMetrics());
        }

        @Override
        public void create(SwipeMenu menu) {
            // create "open" item
            SwipeMenuItem openItem = new SwipeMenuItem(getApplicationContext());
            // set item background
            openItem.setBackground(new ColorDrawable(Color.rgb(0xC9, 0xC9, 0xCE)));
            // set item width
            openItem.setWidth(dp2px(90));
            // set item title
            openItem.setTitle("Update");
            // set item title fontsize
            openItem.setTitleSize(18);
            // set item title font color
            openItem.setTitleColor(Color.WHITE);
            // add to menu
            menu.addMenuItem(openItem);

            // create "delete" item
            SwipeMenuItem deleteItem = new SwipeMenuItem(getApplicationContext());
            // set item background
            deleteItem.setBackground(new ColorDrawable(Color.rgb(0xF9, 0x3F, 0x25)));
            // set item width
            deleteItem.setWidth(dp2px(90));
            // set a icon
            deleteItem.setIcon(R.drawable.delete);
            // add to menu
            menu.addMenuItem(deleteItem);
        }
    }
}
