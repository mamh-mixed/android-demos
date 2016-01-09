package com.cardinfolink.yunshouyin.view;

import android.content.Context;
import android.view.View;
import android.widget.TextView;

import com.cardinfolink.yunshouyin.R;

/**
 * Created by mamh on 16-1-9.
 * 账单界面要显示的一个对话框， 这个和 hintdialog很像，这里继承一下HintDialog
 * 多了3个自己使用的控件而已
 */
public class HintBillDialog extends HintDialog {

    private TextView mBillDate;
    private TextView mBillTotal;
    private TextView mBillCount;

    public HintBillDialog(Context context, View view) {
        super(context, view);
        mBillDate = (TextView) view.findViewById(R.id.hint_bill_datetime);
        mBillCount = (TextView) view.findViewById(R.id.hint_bill_count);
        mBillTotal = (TextView) view.findViewById(R.id.hint_bill_total);
    }

    /**
     * 设置 日期
     *
     * @param date
     */
    public void setBillDate(String date) {
        mBillDate.setText(date);
    }


    /**
     * 设置 当日收入总金额
     *
     * @param total
     */
    public void setBillTotal(String total) {
        mBillTotal.setText(total);
    }

    /**
     * 设置当日收入 笔r
     *
     * @param count
     */
    public void setBillCount(String count) {
        mBillCount.setText(count);
    }


}
