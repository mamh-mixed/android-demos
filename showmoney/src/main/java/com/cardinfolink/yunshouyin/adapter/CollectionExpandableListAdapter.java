package com.cardinfolink.yunshouyin.adapter;

import android.content.Context;
import android.content.Intent;
import android.graphics.Bitmap;
import android.graphics.BitmapFactory;
import android.graphics.Color;
import android.os.AsyncTask;
import android.os.Bundle;
import android.text.TextUtils;
import android.view.View;
import android.view.ViewGroup;
import android.widget.BaseExpandableListAdapter;
import android.widget.ImageView;
import android.widget.TextView;

import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.activity.DetailActivity;
import com.cardinfolink.yunshouyin.data.MonthBill;
import com.cardinfolink.yunshouyin.data.TradeBill;
import com.cardinfolink.yunshouyin.util.EncoderUtil;

import java.io.File;
import java.io.FileOutputStream;
import java.io.IOException;
import java.io.InputStream;
import java.net.HttpURLConnection;
import java.net.URL;
import java.text.ParseException;
import java.text.SimpleDateFormat;
import java.util.Date;
import java.util.List;

/**
 * Created by mamh on 15-12-31.
 * 收款码账单的 adapter
 */
public class CollectionExpandableListAdapter extends BaseExpandableListAdapter {
    private String TAG = "CollectionExpandableListAdapter";

    private List<MonthBill> groupData;
    private List<List<TradeBill>> childrenData;
    private Context mContext;

    private File mExternalCacheDir;

    public CollectionExpandableListAdapter(Context context, List<MonthBill> groupData, List<List<TradeBill>> childrenData) {
        this.mContext = context;
        this.groupData = groupData;
        this.childrenData = childrenData;

        //创建缓存目录，系统一运行就得创建缓存目录的，
        initCacheDir();
    }

    private void initCacheDir() {
        //创建缓存目录，系统一运行就得创建缓存目录的，
        mExternalCacheDir = mContext.getExternalCacheDir();
        //如果外部的不能用，就调用内部的
        if (mExternalCacheDir == null) {
            mExternalCacheDir = mContext.getCacheDir();
        }
        if (!mExternalCacheDir.exists()) {
            mExternalCacheDir.mkdirs();
        }

    }

    @Override
    public int getGroupCount() {
        if (groupData != null) {
            return groupData.size();
        } else {
            return 0;
        }
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
    public Object getChild(int groupPosition, int childPosition) {
        return childrenData.get(groupPosition).get(childPosition);
    }

    @Override
    public long getGroupId(int groupPosition) {
        return 0;
    }

    @Override
    public long getChildId(int groupPosition, int childPosition) {
        return 0;
    }

    @Override
    public boolean hasStableIds() {
        return false;
    }

    @Override
    public View getGroupView(int groupPosition, boolean isExpanded, View convertView, ViewGroup parent) {
        GroupViewHolder groupViewHolder = null;

        if (convertView == null) {
            groupViewHolder = new GroupViewHolder();
            convertView = View.inflate(mContext, R.layout.collection_expandablelistview_group, null);

            groupViewHolder.month = (TextView) convertView.findViewById(R.id.tv_month);
            groupViewHolder.year = (TextView) convertView.findViewById(R.id.tv_year);
            groupViewHolder.count = (TextView) convertView.findViewById(R.id.tv_count);
            groupViewHolder.folder = (ImageView) convertView.findViewById(R.id.iv_fold);

            convertView.setTag(groupViewHolder);
        } else {
            groupViewHolder = (GroupViewHolder) convertView.getTag();
        }
        //设置一下月份
        groupViewHolder.month.setText(groupData.get(groupPosition).getCurrentMonth());
        groupViewHolder.year.setText(groupData.get(groupPosition).getCurrentYear());
        int count = 0;
        try {
            count = childrenData.get(groupPosition).size();
        } catch (Exception e) {
            count = 0;
        }
        groupViewHolder.count.setText(String.valueOf(count));

        if (isExpanded) {
            groupViewHolder.folder.setBackgroundResource(R.drawable.bill_pack);
        } else {
            groupViewHolder.folder.setBackgroundResource(R.drawable.bill_unfold);
        }

        return convertView;
    }

    @Override
    public View getChildView(int groupPosition, int childPosition, boolean isLastChild, View convertView, ViewGroup parent) {
        ChildViewHolder childViewHolder = null;

        if (convertView == null) {
            childViewHolder = new ChildViewHolder();
            convertView = View.inflate(mContext, R.layout.collection_expandablelistview_child, null);
            childViewHolder.linearLayoutDay = convertView.findViewById(R.id.ll_day);
            childViewHolder.linearLayoutBillItem = convertView.findViewById(R.id.ll_bill_item);

            childViewHolder.day = (TextView) convertView.findViewById(R.id.tv_day);
            childViewHolder.weekday = (TextView) convertView.findViewById(R.id.tv_weekday);

            childViewHolder.paylogo = (ImageView) convertView.findViewById(R.id.paylogo);
            childViewHolder.billTradeDate = (TextView) convertView.findViewById(R.id.bill_tradedate);
            childViewHolder.billTradeStatus = (TextView) convertView.findViewById(R.id.bill_tradestatus);
            childViewHolder.billTradeAmount = (TextView) convertView.findViewById(R.id.bill_tradeamount);

            childViewHolder.billCheckCode = (TextView) convertView.findViewById(R.id.bill_checkcode);
            childViewHolder.billNickName = (TextView) convertView.findViewById(R.id.bill_nickname);

            convertView.setTag(childViewHolder);
        } else {
            childViewHolder = (ChildViewHolder) convertView.getTag();
        }

        //从list中根据位置获取到相应的bill项
        final TradeBill bill = childrenData.get(groupPosition).get(childPosition);

        getUrlBitmapAsync(bill.avatarUrl, childViewHolder.paylogo);

        SimpleDateFormat spf1 = new SimpleDateFormat("yyyyMMddHHmmss");
        SimpleDateFormat spf2 = new SimpleDateFormat("HH:mm:ss");
        SimpleDateFormat spf3 = new SimpleDateFormat("dd");
        SimpleDateFormat spf4 = new SimpleDateFormat("E");
        try {
            Date tandeDate = spf1.parse(bill.tandeDate);
            childViewHolder.billTradeDate.setText(spf2.format(tandeDate));
            childViewHolder.day.setText(spf3.format(tandeDate));
            childViewHolder.weekday.setText(spf4.format(tandeDate));
        } catch (ParseException e) {
            e.printStackTrace();
        }


        String tradeStatus;
        if ("10".equals(bill.transStatus)) {
            //处理中
            tradeStatus = mContext.getString(R.string.expandable_listview_trade_status_nopay);
            childViewHolder.billTradeStatus.setTextColor(Color.parseColor("#888888"));
        } else if ("30".equals(bill.transStatus)) {
            double amt = Double.parseDouble(bill.refundAmt);
            if (amt == 0) {
                //成功的
                tradeStatus = mContext.getString(R.string.expandable_listview_trade_status_success);
                childViewHolder.billTradeStatus.setTextColor(Color.parseColor("#888888"));
            } else {
                //部分退款的
                tradeStatus = mContext.getString(R.string.expandable_listview_trade_status_partrefd);
                childViewHolder.billTradeStatus.setTextColor(Color.parseColor("#888888"));
            }
        } else if ("40".equals(bill.transStatus)) {
            if ("09".equals(bill.response)) {
                //已关闭
                tradeStatus = mContext.getString(R.string.expandable_listview_trade_status_closed);
                childViewHolder.billTradeStatus.setTextColor(Color.parseColor("#888888"));
            } else {
                //全额退款
                tradeStatus = mContext.getString(R.string.expandable_listview_trade_status_partrefd);
                childViewHolder.billTradeStatus.setTextColor(Color.parseColor("#888888"));
            }
        } else {
            //失败的
            tradeStatus = mContext.getString(R.string.expandable_listview_trade_status_fail);
            childViewHolder.billTradeStatus.setTextColor(Color.RED);
        }
        childViewHolder.billTradeStatus.setText(tradeStatus);
        childViewHolder.billTradeAmount.setText("￥" + bill.amount);

        childViewHolder.billNickName.setText(bill.nickName);
        childViewHolder.billCheckCode.setText(bill.checkCode);

        childViewHolder.linearLayoutDay.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
            }
        });

        childViewHolder.linearLayoutBillItem.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                Intent intent = new Intent(mContext, DetailActivity.class);
                Bundle bundle = new Bundle();
                bundle.putSerializable("TradeBill", bill);
                intent.putExtra("BillBundle", bundle);
                mContext.startActivity(intent);
            }
        });

        return convertView;
    }

    @Override
    public boolean isChildSelectable(int groupPosition, int childPosition) {
        return true;
    }

    public final class GroupViewHolder {
        public ImageView folder;
        public TextView month;
        public TextView year;
        public TextView count;
    }

    public final class ChildViewHolder {
        public TextView day;
        public TextView weekday;
        public ImageView paylogo;
        public TextView billTradeDate;
        public TextView billTradeStatus;
        public TextView billTradeAmount;

        public TextView billNickName;
        public TextView billCheckCode;

        public View linearLayoutDay;//左边显示日期，周几的一个线性布局
        public View linearLayoutBillItem;//右边显示详情账单信息的一个线性布局
    }


    private void getUrlBitmapAsync(final String url, final ImageView imageView) {
        new AsyncTask<Void, Integer, Bitmap>() {
            @Override
            protected Bitmap doInBackground(Void... params) {
                return getUrlBitmap(url);
            }

            @Override
            protected void onPostExecute(Bitmap bitmap) {
                if (bitmap == null) {
                    imageView.setImageResource(R.drawable.wpay);
                } else {
                    imageView.setImageBitmap(bitmap);
                }
            }
        }.execute();
    }

    /**
     * 获取网络上的图片
     *
     * @param url
     * @return
     */
    private Bitmap getUrlBitmap(String url) {
        URL myFileUrl = null;
        Bitmap bitmap = null;

        if (TextUtils.isEmpty(url)) {
            return bitmap;
        }

        String name = EncoderUtil.Encrypt(url, "MD5");
        File imageFile = new File(mExternalCacheDir, name);

        if (imageFile.exists()) {
            //如果头像图片文件存在就直接使用
            bitmap = BitmapFactory.decodeFile(imageFile.getPath());
            return bitmap;
        }

        //走到这里表明头像图片不存在，这里就要下载了
        try {
            myFileUrl = new URL(url);
        } catch (Exception e) {
            e.printStackTrace();
        }
        try {
            HttpURLConnection conn = (HttpURLConnection) myFileUrl.openConnection();
            conn.setDoInput(true);
            conn.connect();
            InputStream is = conn.getInputStream();

            FileOutputStream fos = new FileOutputStream(imageFile);
            byte[] buffer = new byte[1024];
            int len = 0;
            while ((len = is.read(buffer)) != -1) {
                fos.write(buffer, 0, len);
            }

            bitmap = BitmapFactory.decodeFile(imageFile.getPath());

            is.close();
            fos.close();
        } catch (IOException e) {
            e.printStackTrace();
        }
        return bitmap;
    }


}
