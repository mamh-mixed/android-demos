package com.cardinfolink.yunshouyin.salesman.activities;

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
import com.cardinfolink.yunshouyin.salesman.models.SAServerPacket;
import com.cardinfolink.yunshouyin.salesman.models.SessonData;
import com.cardinfolink.yunshouyin.salesman.models.User;
import com.cardinfolink.yunshouyin.salesman.utils.CommunicationListenerV2;
import com.cardinfolink.yunshouyin.salesman.utils.HttpCommunicationUtil;
import com.cardinfolink.yunshouyin.salesman.utils.ParamsUtil;
import com.cardinfolink.yunshouyin.salesman.utils.SAApplication;
import com.cardinfolink.yunshouyin.salesman.utils.SADownloader;
import com.cardinfolink.yunshouyin.salesman.utils.SAImageUtil;
import com.cardinfolink.yunshouyin.salesman.utils.SAQrCodeUtil;

import java.util.ArrayList;
import java.util.List;

class MerchantListAdapter extends ArrayAdapter<User> {
    private SAMerchantListActivity saMerchantListActivity;
    Filter myFilter;
    private List<User> users_origin = new ArrayList<>();
    private List<User> users;

    public MerchantListAdapter(SAMerchantListActivity saMerchantListActivity, final List<User> users) {
        super(saMerchantListActivity, R.layout.merchant_item_view, users);
        this.saMerchantListActivity = saMerchantListActivity;
        this.users = users;
        this.users_origin.addAll(users);

        myFilter = new Filter() {
            @Override
            protected FilterResults performFiltering(CharSequence constraint) {
                FilterResults filterResults = new FilterResults();
                ArrayList<User> tmpUsers = new ArrayList<>();
                // 没有关键字,数据内容使用原始的数据拷贝引用
                if (constraint == null || constraint.length() == 0) {
                    filterResults.values = users_origin;
                    filterResults.count = users_origin.size();
                } else if (users_origin != null) {
                    for (User user : users_origin) {
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
//                    saMerchantListActivity.adapter.clear();
                    users.addAll(objects);
                }
                notifyDataSetChanged();
//                else{
//                        //TODO:失效有何用???
//                        notifyDataSetInvalidated();
//                    }
            }
        };
    }

    @Override
    public View getView(final int position, View convertView, ViewGroup parent) {
        View itemView = convertView;
        if (itemView == null) {
            itemView = saMerchantListActivity.getLayoutInflater().inflate(R.layout.merchant_item_view, parent, false);
        }

        final User merchant = users.get(position);

        final TableLayout detailViewGroup = (TableLayout) itemView.findViewById(R.id.mItem_detailViewGroup);
        TextView merchantNameText = (TextView) itemView.findViewById(R.id.mItem_txtMerchantName);
        merchantNameText.setText(merchant.getMerName());

        merchantNameText.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                Toast.makeText(saMerchantListActivity, "点击商户名收起或展开商户项", Toast.LENGTH_SHORT).show();
                int status = detailViewGroup.getVisibility();
                if (status == View.GONE) {
                    detailViewGroup.setVisibility(View.VISIBLE);
                } else if (status == View.VISIBLE) {
                    detailViewGroup.setVisibility(View.GONE);
                }
            }
        });

        TextView emailText = (TextView) itemView.findViewById(R.id.mItem_txtEmail);
        emailText.setText(merchant.getUsername());
        TextView midText = (TextView) itemView.findViewById(R.id.mItem_txtMerchantId);
        midText.setText(merchant.getClientid());
        TextView secretText = (TextView) itemView.findViewById(R.id.mItem_txtSecretKey);
        secretText.setText(merchant.getSignKey());

        TextView downloadQRText = (TextView) itemView.findViewById(R.id.merchantlist_download_qrcode);
        downloadQRText.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                Log.d("jiahua:", "download qrcode");
                Toast.makeText(SAApplication.getContext(), "二维码生成中...", Toast.LENGTH_LONG).show();

                String merchantId = merchant.getClientid();
                HttpCommunicationUtil.sendDataToQuickIpayServer(ParamsUtil.getDownload(SessonData.getAccessToken(), merchantId, "bill"), new CommunicationListenerV2() {
                    @Override
                    public void onResult(SAServerPacket serverPacket) {
                        String imageUrl = serverPacket.getDownloadUrl();
                        try {
                            Bitmap bitmap = SADownloader.downloadBitmap(imageUrl);
                            new SAImageUtil().saveImageToExternalStorage(bitmap);
                        }catch (final Exception ex){
                            Log.d("jiahua", ex.getMessage());
                            new Handler(Looper.getMainLooper()).post(new Runnable() {
                                @Override
                                public void run() {
                                    Toast.makeText(SAApplication.getContext(), "下载错误:"+ex.getMessage(), Toast.LENGTH_LONG).show();
                                }
                            });

                            return;
                        }

                        new Handler(Looper.getMainLooper()).post(new Runnable() {
                            @Override
                            public void run() {
                                Toast.makeText(SAApplication.getContext(), "账单二维码已经下载到相册", Toast.LENGTH_LONG).show();
                            }
                        });
                    }

                    @Override
                    public void onError(final String error) {
                        new Handler(Looper.getMainLooper()).post(new Runnable() {
                            @Override
                            public void run() {
                                Toast.makeText(SAApplication.getContext(), "下载失败:"+error, Toast.LENGTH_LONG).show();
                            }
                        });
                    }
                });

                HttpCommunicationUtil.sendDataToQuickIpayServer(ParamsUtil.getDownload(SessonData.getAccessToken(), merchantId, "pay"), new CommunicationListenerV2() {
                    @Override
                    public void onResult(SAServerPacket serverPacket) {
                        String imageUrl = serverPacket.getDownloadUrl();
                        try {
                            Bitmap bitmap = SADownloader.downloadBitmap(imageUrl);
                            new SAImageUtil().saveImageToExternalStorage(bitmap);
                        } catch (final Exception ex) {
                            Log.d("jiahua", ex.getMessage());
                            new Handler(Looper.getMainLooper()).post(new Runnable() {
                                @Override
                                public void run() {
                                    Toast.makeText(SAApplication.getContext(), "下载错误:" + ex.getMessage(), Toast.LENGTH_LONG).show();
                                }
                            });

                            return;
                        }

                        new Handler(Looper.getMainLooper()).post(new Runnable() {
                            @Override
                            public void run() {
                                Toast.makeText(SAApplication.getContext(), "支付二维码已经下载到相册", Toast.LENGTH_LONG).show();
                            }
                        });
                    }

                    @Override
                    public void onError(final String error) {
                        new Handler(Looper.getMainLooper()).post(new Runnable() {
                            @Override
                            public void run() {
                                Toast.makeText(SAApplication.getContext(), "下载失败:" + error, Toast.LENGTH_LONG).show();
                            }
                        });
                    }
                });
            }
        });

        return itemView;
    }

    @Override
    public Filter getFilter() {
        return myFilter;
    }

    public void refreshDataSource(List<User> users) {
        users_origin.clear();
        users_origin.addAll(users);
    }
}
