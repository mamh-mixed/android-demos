package com.cardinfolink.yunshouyin.salesman.adapter;

import android.graphics.Bitmap;
import android.os.Handler;
import android.os.Looper;
import android.util.Log;
import android.view.View;
import android.view.ViewGroup;
import android.widget.ArrayAdapter;
import android.widget.Filter;
import android.widget.TableLayout;
import android.widget.TextView;
import android.widget.Toast;

import com.cardinfolink.yunshouyin.salesman.R;
import com.cardinfolink.yunshouyin.salesman.activity.MerchantListActivity;
import com.cardinfolink.yunshouyin.salesman.api.QuickPayException;
import com.cardinfolink.yunshouyin.salesman.core.QuickPayCallbackListener;
import com.cardinfolink.yunshouyin.salesman.core.QuickPayService;
import com.cardinfolink.yunshouyin.salesman.model.User;
import com.cardinfolink.yunshouyin.salesman.utils.SalesmanApplication;
import com.cardinfolink.yunshouyin.salesman.utils.Downloader;
import com.cardinfolink.yunshouyin.salesman.utils.ImageUtil;

import java.util.ArrayList;
import java.util.List;

public class MerchantListAdapter extends ArrayAdapter<User> {
    private static final String TAG = "MerchantListAdapter";

    private Filter myFilter;
    private MerchantListActivity merchantListActivity;
    private List<User> usersOrigin = new ArrayList<>();
    private List<User> users;

    public MerchantListAdapter(MerchantListActivity merchantListActivity, final List<User> users) {
        super(merchantListActivity, R.layout.merchant_item_view, users);
        this.merchantListActivity = merchantListActivity;
        this.users = users;
        this.usersOrigin.addAll(users);

        myFilter = new MerchantFilter();

    }

    private class MerchantFilter extends Filter {
        @Override
        protected FilterResults performFiltering(CharSequence constraint) {
            FilterResults filterResults = new FilterResults();
            ArrayList<User> tmpUsers = new ArrayList<>();
            // 没有关键字,数据内容使用原始的数据拷贝引用
            if (constraint == null || constraint.length() == 0) {
                filterResults.values = usersOrigin;
                filterResults.count = usersOrigin.size();
            } else if (usersOrigin != null) {
                for (User user : usersOrigin) {
                    if (user.getMerName() != null && user.getMerName().contains(constraint)) {
                        tmpUsers.add(user);
                    }
                }
                filterResults.values = tmpUsers;
                filterResults.count = tmpUsers.size();
            }
            return filterResults;
        }

        @Override
        protected void publishResults(CharSequence constraint, FilterResults results) {
            ArrayList<User> objects = (ArrayList<User>) results.values;
            //这里并没有引用arrayList的地址,而是对list内的item逐个加入adapter
            users.clear();
            if (objects != null && objects.size() > 0) {
                //MerchantListActivity.adapter.clear();
                users.addAll(objects);
            }
            notifyDataSetChanged();
        }
    }

    @Override
    public View getView(final int position, View convertView, ViewGroup parent) {
        ViewHolder holder = null;
        if (convertView == null) {
            convertView = merchantListActivity.getLayoutInflater().inflate(R.layout.merchant_item_view, parent, false);
            holder = new ViewHolder();
            holder.detailViewGroup = (TableLayout) convertView.findViewById(R.id.mItem_detailViewGroup);
            holder.merchantName = (TextView) convertView.findViewById(R.id.mItem_txtMerchantName);
            holder.merchantEmail = (TextView) convertView.findViewById(R.id.mItem_txtEmail);
            holder.merchantId = (TextView) convertView.findViewById(R.id.mItem_txtMerchantId);
            holder.merchantSecret = (TextView) convertView.findViewById(R.id.mItem_txtSecretKey);
            holder.downloadQR = (TextView) convertView.findViewById(R.id.merchantlist_download_qrcode);

            convertView.setTag(holder);

        }else{
            holder = (ViewHolder)convertView.getTag();
        }
        final User merchant = users.get(position);
        final TableLayout detailViewGroup = holder.detailViewGroup;

        holder.merchantName.setText(merchant.getMerName());
        holder.merchantName.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                int status = detailViewGroup.getVisibility();
                if (status == View.GONE) {
                    detailViewGroup.setVisibility(View.VISIBLE);
                } else if (status == View.VISIBLE) {
                    detailViewGroup.setVisibility(View.GONE);
                }
            }
        });

        holder.merchantEmail.setText(merchant.getUsername());

        holder.merchantId.setText(merchant.getClientid());

        holder.merchantSecret.setText(merchant.getSignKey());

        holder.downloadQR.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                Log.d(TAG, "download qrcode");
                Toast.makeText(merchantListActivity, "二维码生成中...", Toast.LENGTH_LONG).show();

                String merchantId = merchant.getClientid();

                QuickPayService quickPayService = SalesmanApplication.getInstance().getQuickPayService();
                quickPayService.getQrPostUrlAsync(merchantId, "bill", new QuickPayCallbackListener<String>() {
                    @Override
                    public void onSuccess(String data) {
                        String imageUrl = data;
                        try {
                            Bitmap bitmap = Downloader.downloadBitmap(imageUrl);
                            new ImageUtil().saveImageToExternalStorage(bitmap);
                        } catch (final Exception ex) {
                            Log.d(TAG, ex.getMessage());
                            new Handler(Looper.getMainLooper()).post(new Runnable() {
                                @Override
                                public void run() {
                                    Toast.makeText(merchantListActivity, "下载错误:" + ex.getMessage(), Toast.LENGTH_LONG).show();
                                }
                            });

                            return;
                        }

                        new Handler(Looper.getMainLooper()).post(new Runnable() {
                            @Override
                            public void run() {
                                Toast.makeText(merchantListActivity, "账单二维码已经下载到相册", Toast.LENGTH_LONG).show();
                            }
                        });
                    }

                    @Override
                    public void onFailure(final QuickPayException ex) {
                        new Handler(Looper.getMainLooper()).post(new Runnable() {
                            @Override
                            public void run() {
                                Toast.makeText(merchantListActivity, "下载失败:" + ex.getErrorMsg(), Toast.LENGTH_LONG).show();
                            }
                        });
                    }
                });

                quickPayService.getQrPostUrlAsync(merchantId, "pay", new QuickPayCallbackListener<String>() {
                    @Override
                    public void onSuccess(String data) {
                        String imageUrl = data;
                        try {
                            Bitmap bitmap = Downloader.downloadBitmap(imageUrl);
                            new ImageUtil().saveImageToExternalStorage(bitmap);
                        } catch (final Exception ex) {
                            Log.d(TAG, ex.getMessage());
                            new Handler(Looper.getMainLooper()).post(new Runnable() {
                                @Override
                                public void run() {
                                    Toast.makeText(merchantListActivity, "下载错误:" + ex.getMessage(), Toast.LENGTH_LONG).show();
                                }
                            });

                            return;
                        }

                        new Handler(Looper.getMainLooper()).post(new Runnable() {
                            @Override
                            public void run() {
                                Toast.makeText(merchantListActivity, "支付二维码已经下载到相册", Toast.LENGTH_LONG).show();
                            }
                        });
                    }

                    @Override
                    public void onFailure(final QuickPayException ex) {
                        new Handler(Looper.getMainLooper()).post(new Runnable() {
                            @Override
                            public void run() {
                                Toast.makeText(merchantListActivity, "下载失败:" + ex.getErrorMsg(), Toast.LENGTH_LONG).show();
                            }
                        });
                    }
                });
            }
        });

        return convertView;
    }

    @Override
    public Filter getFilter() {
        return myFilter;
    }

    public void refreshDataSource(List<User> users) {
        usersOrigin.clear();
        usersOrigin.addAll(users);
    }

    private static class ViewHolder {
        public TableLayout detailViewGroup;
        public TextView merchantName;
        public TextView merchantEmail;
        public TextView merchantId;
        public TextView merchantSecret;
        public TextView downloadQR;
    }
}
