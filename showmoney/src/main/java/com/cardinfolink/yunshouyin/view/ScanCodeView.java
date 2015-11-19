package com.cardinfolink.yunshouyin.view;

import android.content.Context;
import android.content.Intent;
import android.support.annotation.NonNull;
import android.util.Log;
import android.view.LayoutInflater;
import android.view.View;
import android.view.View.OnClickListener;
import android.widget.Button;
import android.widget.EditText;
import android.widget.LinearLayout;
import android.widget.Toast;

import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.activity.CaptureActivity;

public class ScanCodeView extends LinearLayout implements OnClickListener {
    public static final String CHANNEL_ALP = "ALP";
    public static final String CHANNEL_WXP = "WXP";
    public static final long MAX_PURCHASE_AMOUNT = 999999999;

    private long amount;

    Button btn0, btn1, btn2, btn3, btn4, btn5, btn6,
            btn7, btn8, btn9, btn00,
            btnsm, btnclear;

    int[] numberBtnArr;

    EditText edt_input;

    public ScanCodeView(Context context) {
        super(context);
        View contentView = LayoutInflater.from(context).inflate(
                R.layout.scancode_view, null);

        LinearLayout.LayoutParams layoutParams = new LayoutParams(
                LayoutParams.MATCH_PARENT, LayoutParams.MATCH_PARENT);
        contentView.setLayoutParams(layoutParams);
        addView(contentView);
        initLayout();
    }

    private void initLayout() {
        btn00 = (Button) findViewById(R.id.btn00);
        btn0 = (Button) findViewById(R.id.btn0);
        btn1 = (Button) findViewById(R.id.btn1);
        btn2 = (Button) findViewById(R.id.btn2);
        btn3 = (Button) findViewById(R.id.btn3);
        btn4 = (Button) findViewById(R.id.btn4);
        btn5 = (Button) findViewById(R.id.btn5);
        btn6 = (Button) findViewById(R.id.btn6);
        btn7 = (Button) findViewById(R.id.btn7);
        btn8 = (Button) findViewById(R.id.btn8);
        btn9 = (Button) findViewById(R.id.btn9);

        numberBtnArr = new int[10];
        numberBtnArr[0] = R.id.btn0;
        numberBtnArr[1] = R.id.btn1;
        numberBtnArr[2] = R.id.btn2;
        numberBtnArr[3] = R.id.btn3;
        numberBtnArr[4] = R.id.btn4;
        numberBtnArr[5] = R.id.btn5;
        numberBtnArr[6] = R.id.btn6;
        numberBtnArr[7] = R.id.btn7;
        numberBtnArr[8] = R.id.btn8;
        numberBtnArr[9] = R.id.btn9;

        btnsm = (Button) findViewById(R.id.btnsm);
        btnclear = (Button) findViewById(R.id.btnclear);

        edt_input = (EditText) findViewById(R.id.edt_input);

        btn00.setOnClickListener(this);
        btn0.setOnClickListener(this);
        btn1.setOnClickListener(this);
        btn2.setOnClickListener(this);
        btn3.setOnClickListener(this);
        btn4.setOnClickListener(this);
        btn5.setOnClickListener(this);
        btn6.setOnClickListener(this);
        btn7.setOnClickListener(this);
        btn8.setOnClickListener(this);
        btn9.setOnClickListener(this);

        btnclear.setOnClickListener(this);
        btnsm.setOnClickListener(this);
    }

    public void clearValue() {
        amount = 0;
        edt_input.setText(getDisplayAmount());
    }

    @Override
    public void onClick(View v) {
        long tempAmount = amount;
        for (int i = 0; i < 10; i++) {
            if (v.getId() == numberBtnArr[i]) {
                tempAmount = tempAmount * 10 + i;

                //TODO: 最大额度待确认,超过long类型最大值?
                if (tempAmount > MAX_PURCHASE_AMOUNT) {
                    Toast.makeText(getContext(), getResources().getString(R.string.qr_amount_exceed), Toast.LENGTH_SHORT).show();
                    return;
                }
                amount = tempAmount;
                edt_input.setText(getDisplayAmount());
                return;
            }
        }
        if (v.getId() == R.id.btn00) {
            tempAmount = tempAmount * 100;

            if (tempAmount > MAX_PURCHASE_AMOUNT) {
                Toast.makeText(getContext(), getResources().getString(R.string.qr_amount_exceed), Toast.LENGTH_SHORT).show();
                return;
            }
            amount = tempAmount;
            edt_input.setText(getDisplayAmount());
            return;
        }

        Log.d("TEST", "Amount " + amount);
        Log.d("TEST", "DisplayAmount: " + getDisplayAmount());

        switch (v.getId()) {
            case R.id.btnclear:
                clearValue();
                break;

            case R.id.btnsm:
                final double sum = amount;
                if (sum <= 0) {
                    Toast.makeText(getContext(), getResources().getString(R.string.qr_amount_nonzero), Toast.LENGTH_SHORT).show();
                    return;
                }

                String chcd = CHANNEL_ALP;

                Intent intent = new Intent(getContext(), CaptureActivity.class);
                intent.putExtra("chcd", chcd);
                intent.putExtra("total", "" + sum);
                getContext().startActivity(intent);
                break;
        }
    }

    @NonNull
    private String getDisplayAmount() {
        StringBuilder displayAmount = new StringBuilder();
        displayAmount.append(amount);

        // insert ',' every there number
        for (int i = displayAmount.length() - 3; i > 0; i = i - 3) {
            displayAmount.insert(i, ',');
        }

        return displayAmount.toString();
    }
}
