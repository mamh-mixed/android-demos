package com.cardinfolink.yunshouyin.salesman.adapter;


import android.content.Context;
import android.support.v7.widget.RecyclerView;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;

import com.cardinfolink.yunshouyin.salesman.R;
import com.cardinfolink.yunshouyin.salesman.model.SAMerchantPhoto;
import com.cardinfolink.yunshouyin.salesman.view.MerchantImageViewHolder;

import java.util.List;

public class MerchantPhotoRecyclerViewAdapter extends RecyclerView.Adapter<MerchantImageViewHolder> {

    private List<SAMerchantPhoto> itemList;
    private Context context;

    public MerchantPhotoRecyclerViewAdapter(Context context, List<SAMerchantPhoto> itemList) {
        this.itemList = itemList;
        this.context = context;
    }

    /**
     * 创建view
     *
     * @param parent
     * @param viewType
     * @return
     */
    @Override
    public MerchantImageViewHolder onCreateViewHolder(ViewGroup parent, int viewType) {
        View layoutView = LayoutInflater.from(parent.getContext()).inflate(R.layout.solvent_list, null);
        MerchantImageViewHolder miv = new MerchantImageViewHolder(layoutView);
        return miv;
    }

    /**
     * 绑定view到数据
     *
     * @param holder
     * @param position
     */
    @Override
    public void onBindViewHolder(MerchantImageViewHolder holder, int position) {
        SAMerchantPhoto saMerchantPhoto = itemList.get(position);
//        String imagePath = saMerchantPhoto.getFilename();
//        if (imagePath != null){//这样特别耗内存
//            Bitmap bitmap = BitmapFactory.decodeFile(imagePath);
//            holder.merchantPhoto.setImageBitmap(bitmap);
//        }
        holder.merchantPhoto.setImageURI(saMerchantPhoto.getImageUri());
    }

    @Override
    public int getItemCount() {
        return this.itemList.size();
    }

}