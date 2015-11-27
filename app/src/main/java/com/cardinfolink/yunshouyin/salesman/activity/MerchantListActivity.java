package com.cardinfolink.yunshouyin.salesman.activity;

import android.app.AlertDialog;
import android.content.DialogInterface;
import android.content.Intent;
import android.content.SharedPreferences;
import android.os.Bundle;
import android.support.v4.widget.SwipeRefreshLayout;
import android.text.Editable;
import android.text.TextWatcher;
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
import com.cardinfolink.yunshouyin.salesman.utils.ActivityCollector;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.Calendar;
import java.util.Collections;
import java.util.Comparator;
import java.util.Date;
import java.util.List;

public class MerchantListActivity extends BaseActivity {
    private final String TAG = "MerchantListActivity";
    MerchantListAdapter merchantListAdapter;
    //该地址会被ArrayAdapter所引用,作为数据源,对merchantInfos所做的修改会影响到arrayAdapter
    private List<User> users = new ArrayList<>();

    private SwipeRefreshLayout swipeRefreshLayout;
    private Button mAddNewMerchant;
    private EditText mSearchMerchant;
    private TextView mMerchantCount;//this month mer count

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_merchant_list);
        initLayout();
    }

    @Override
    protected void onDestroy() {
        super.onDestroy();
    }


    /**
     * setup listener
     */
    private void initLayout() {
        merchantListAdapter = new MerchantListAdapter(this, users);
        ListView listView = (ListView) findViewById(R.id.listViewMerchants);
        listView.setAdapter(merchantListAdapter);

        mMerchantCount = (TextView) findViewById(R.id.txt_merchantcountthismonth);

        swipeRefreshLayout = (SwipeRefreshLayout) findViewById(R.id.swipe_container);
        swipeRefreshLayout.setOnRefreshListener(new MerchantOnRefreshListener());

        mAddNewMerchant = (Button) findViewById(R.id.btnAddNewMerchant);
        mAddNewMerchant.setOnClickListener(new AddNewMerchantOnClickListener());

        //输入关键字快速定位
        mSearchMerchant = (EditText) findViewById(R.id.mItem_txtSearch);
        mSearchMerchant.addTextChangedListener(new SearchMerchantTextChangedListener());
    }


    /**
     * called when refresh data, pull down to refresh
     */
    private void refreshData() {
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

                Date today = new Date();
                Calendar calendar = Calendar.getInstance();
                calendar.setTime(today);
                calendar.set(Calendar.DAY_OF_MONTH, 1);
                calendar.set(Calendar.HOUR, 0);
                calendar.set(Calendar.MINUTE, 0);
                calendar.set(Calendar.SECOND, 0);
                calendar.set(Calendar.MILLISECOND, 0);
                Date firstDayOfMonth = calendar.getTime();
                int num = 0;
                for (User user : tempUsers) {
                    if (user.getCreateTime().after(firstDayOfMonth)) {
                        num++;
                    }
                }

                mMerchantCount.setText(String.format("本月已经发展商户: %d 家", num));
                users.clear();
                users.addAll(tempUsers);
                merchantListAdapter.refreshDataSource(users);
                merchantListAdapter.notifyDataSetChanged();
                swipeRefreshLayout.setRefreshing(false);
                Toast.makeText(MerchantListActivity.this, "刷新成功", Toast.LENGTH_SHORT).show();
            }

            @Override
            public void onFailure(final QuickPayException ex) {
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


    @Override
    protected void onResume() {
        super.onResume();
        refreshData();
    }

    private class MerchantOnRefreshListener implements SwipeRefreshLayout.OnRefreshListener {

        @Override
        public void onRefresh() {
            refreshData();
        }
    }

    private class AddNewMerchantOnClickListener implements View.OnClickListener {

        @Override
        public void onClick(View v) {
            int step = mRegisterSharedPreferences.getInt("register_step_finish", 0);
            final AlertDialog.Builder builder = new AlertDialog.Builder(MerchantListActivity.this);
            builder.setNegativeButton("否", new DialogInterface.OnClickListener() {
                @Override
                public void onClick(DialogInterface dialog, int which) {
                    mRegisterSharedPreferences.edit().clear().commit();
                    Intent intent = new Intent(MerchantListActivity.this, RegisterActivity.class);
                    startActivity(intent);
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
                    Intent intent = new Intent(MerchantListActivity.this, RegisterActivity.class);
                    startActivity(intent);
                    break;
            }
        }
    }

    private class SearchMerchantTextChangedListener implements TextWatcher {

        @Override
        public void beforeTextChanged(CharSequence s, int start, int count, int after) {

        }

        @Override
        public void onTextChanged(CharSequence s, int start, int before, int count) {
            merchantListAdapter.getFilter().filter(s);
        }

        @Override
        public void afterTextChanged(Editable s) {

        }
    }
}
