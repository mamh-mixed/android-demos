package com.cardinfolink.yunshouyin.view;

import android.app.Activity;
import android.content.Context;
import android.content.Intent;
import android.graphics.BitmapFactory;
import android.util.Log;
import android.view.LayoutInflater;
import android.view.View;
import android.view.View.OnClickListener;
import android.widget.Button;
import android.widget.EditText;
import android.widget.LinearLayout;
import android.widget.RadioButton;
import android.widget.TextView;
import android.widget.Toast;

import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.activity.BaseActivity;
import com.cardinfolink.yunshouyin.activity.CaptureActivity;
import com.cardinfolink.yunshouyin.activity.CreateQRcodeActivity;
import com.cardinfolink.yunshouyin.data.SessonData;
import com.cardinfolink.yunshouyin.util.CommunicationListener;
import com.cardinfolink.yunshouyin.util.ContextUtil;
import com.cardinfolink.yunshouyin.util.ErrorUtil;
import com.cardinfolink.yunshouyin.util.HttpCommunicationUtil;
import com.cardinfolink.yunshouyin.util.JsonUtil;
import com.cardinfolink.yunshouyin.util.ParamsUtil;

import java.text.SimpleDateFormat;
import java.util.ArrayList;
import java.util.Date;
import java.util.List;

public class ScanCodeView extends LinearLayout implements OnClickListener {
    Button btn0, btn1, btn2, btn3, btn4, btn5, btn6,
            btn7, btn8, btn9, btnadd, btnpoint, btnsm, btnclear, btndelete,
            swh;
    RadioButton btnzhifubao, btnweixin;
    EditText edt_input;
    TextView txt_output;
    boolean clear_flag = true;
    boolean point_flag = true;
    boolean add_flag = true;
    boolean switch_flag = true;
    boolean num_flag = true;
    double result = 0;
    String[] s = new String[100];
    private Context mContext;
    private List<Item> items = new ArrayList<Item>();
    private BaseActivity mBaseActivity;

    public ScanCodeView(Context context) {
        super(context);
        mContext = context;
        mBaseActivity = (BaseActivity) mContext;
        View contentView = LayoutInflater.from(context).inflate(
                R.layout.scancode_view, null);
        LinearLayout.LayoutParams layoutParams = new LayoutParams(
                LayoutParams.MATCH_PARENT, LayoutParams.MATCH_PARENT);
        contentView.setLayoutParams(layoutParams);
        addView(contentView);
        initLayout();

    }

    public void clearValue() {
        num_flag = true;
        edt_input.setText("=0");
        txt_output.setText("0");
        add_flag = true;
        point_flag = true;
        clear_flag = true;
    }

    private void initLayout() {
        swh = (Button) findViewById(R.id.swh);
        btnzhifubao = (RadioButton) findViewById(R.id.btnzhifubao);
        btnweixin = (RadioButton) findViewById(R.id.btnweixin);
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
        btnadd = (Button) findViewById(R.id.btnadd);
        btnpoint = (Button) findViewById(R.id.btnpoint);
        btnsm = (Button) findViewById(R.id.btnsm);
        btnclear = (Button) findViewById(R.id.btnclear);
        btndelete = (Button) findViewById(R.id.btndelete);

        edt_input = (EditText) findViewById(R.id.edt_input);
        txt_output = (TextView) findViewById(R.id.txt_output);

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
        btnadd.setOnClickListener(this);
        btnpoint.setOnClickListener(this);
        btnclear.setOnClickListener(this);
        btndelete.setOnClickListener(this);
        btnzhifubao.setOnClickListener(this);
        btnweixin.setOnClickListener(this);
        btnsm.setOnClickListener(this);

        swh.setOnClickListener(new OnClickListener() {

            @Override
            public void onClick(View v) {
                if (switch_flag) {
                    btnsm.setText(ContextUtil.getResString(R.string.scancode_view_create_code));
                    switch_flag = false;
                } else {
                    btnsm.setText(ContextUtil.getResString(R.string.scancode_view_scaning_code));
                    switch_flag = true;
                }

            }
        });

    }

    @Override
    public void onClick(View v) {

        String x = txt_output.getText().toString();

        switch (v.getId()) {

            case R.id.btn0:

                if (num_flag) {

                    clearzero();
                    txt_output.append("0");
                    add_flag = true;
                    getResult();

                }

                break;
            case R.id.btn1:

                if (num_flag) {
                    clearzero();
                    txt_output.append("1");
                    getResult();
                    add_flag = true;

                }

                break;
            case R.id.btn2:
                if (num_flag) {
                    clearzero();
                    txt_output.append("2");
                    getResult();
                    add_flag = true;
                }
                break;
            case R.id.btn3:
                if (num_flag) {
                    clearzero();
                    txt_output.append("3");
                    getResult();
                }
                break;
            case R.id.btn4:
                if (num_flag) {
                    clearzero();
                    txt_output.append("4");
                    getResult();
                    add_flag = true;
                }
                break;
            case R.id.btn5:
                if (num_flag) {
                    clearzero();
                    txt_output.append("5");
                    getResult();
                    add_flag = true;
                }
                break;
            case R.id.btn6:
                if (num_flag) {
                    clearzero();
                    txt_output.append("6");
                    getResult();
                    add_flag = true;
                }
                break;
            case R.id.btn7:
                if (num_flag) {
                    clearzero();
                    txt_output.append("7");
                    getResult();
                    add_flag = true;
                }
                break;
            case R.id.btn8:
                if (num_flag) {
                    clearzero();
                    txt_output.append("8");
                    add_flag = true;
                    getResult();
                }
                break;
            case R.id.btn9:
                if (num_flag) {
                    clearzero();
                    txt_output.append("9");
                    getResult();
                    add_flag = true;
                }
                break;
            case R.id.btnpoint:

                String s1 = x.substring(x.lastIndexOf("+") + 1);
                if (s1.contains(".")) {
                    break;
                }

                if (x.contains(".")) {
                    String k = x.substring(x.lastIndexOf("."));
                    if (k.equals("."))

                    {
                        return;

                    } else {
                        clearzero(".");
                        point_falg(".");
                    }
                } else {

                    clearzero(".");
                    point_falg(".");
                }
                break;
            case R.id.btnadd:

                if (x.contains("+")) {
                    String k = x.substring(x.lastIndexOf("+"));

                    if (k.equals("+"))

                    {
                        return;
                    } else {
                        clearzero("+");
                        add_falg("+");

                    }

                } else {

                    clearzero("+");
                    add_falg("+");

                }
                break;
            case R.id.btnclear:
                num_flag = true;
                edt_input.setText("=0");
                txt_output.setText("0");
                add_flag = true;
                point_flag = true;
                clear_flag = true;
                break;
            case R.id.btndelete:

                String r = edt_input.getText().toString();

                add_flag = true;
                if (x.contains(".")) {
                    String k = x.substring(x.lastIndexOf("."));
                    if (k.equals("."))

                    {
                        point_flag = true;

                    }
                }
                if (x != null && !x.equals("")) {
                    String k = x.substring(x.lastIndexOf("+") + 1);
                    // if (x.contains("+")&&k.equals("+")) {
                    //
                    // add_flag=true;
                    //
                    // }

                    txt_output.setText(x.substring(0, x.length() - 1));
                    if (x.contains("+")) {
                        // String k=x.substring(x.lastIndexOf("+"));
                        if (k.equals("+"))

                        {

                            add_flag = false;
                        } else {
                            add_flag = true;
                        }
                    } else {
                        add_flag = true;
                    }

                }
                addcheck();

                break;
            case R.id.btnzhifubao:

                break;
            case R.id.btnweixin:

                break;
            case R.id.btnsm:
                final double sum = Double.parseDouble(edt_input.getText().toString().substring(1));
                if (sum <= 0) {
                    Toast.makeText(mContext, "金额不能为零!", Toast.LENGTH_SHORT).show();
                    return;
                }

                if (SessonData.loginUser.getLimit().equals("true")) {
                    HttpCommunicationUtil.sendDataToServer(ParamsUtil.getTotal(SessonData.loginUser, (new SimpleDateFormat("yyyyMMdd")).format(new Date())), new CommunicationListener() {

                        @Override
                        public void onResult(String result) {
                            String state = JsonUtil.getParam(result,
                                    "state");

                            if (state.equals("success")) {

                                double limitValue = Double
                                        .parseDouble(JsonUtil.getParam(
                                                result, "total"));
                                if (limitValue >= 500) {
                                    ((Activity) mContext).runOnUiThread(new Runnable() {

                                        @Override
                                        public void run() {

                                            Alert_Dialog alert_Dialog = new Alert_Dialog(mContext, null, ((Activity) mContext).findViewById(R.id.alert_dialog), "当日交易已超过限额,请申请提升限额!", BitmapFactory.decodeResource(mContext.getResources(), R.drawable.wrong));
                                            alert_Dialog.show();

                                        }

                                    });
                                } else {
                                    String chcd = "ALP";
                                    if (btnweixin.isChecked()) {
                                        chcd = "WXP";
                                    } else {
                                        chcd = "ALP";
                                    }
                                    Log.i("opp", "" + chcd);

                                    if (switch_flag) {
                                        Intent intent = new Intent(mContext, CaptureActivity.class);
                                        intent.putExtra("chcd", chcd);
                                        intent.putExtra("total", "" + sum);
                                        mContext.startActivity(intent);

                                    } else {
                                        Intent intent = new Intent(mContext, CreateQRcodeActivity.class);
                                        intent.putExtra("chcd", chcd);
                                        intent.putExtra("total", "" + sum);
                                        mContext.startActivity(intent);

                                    }

                                }

                            } else {
                                final String error = JsonUtil.getParam(result, "error");
                                ((Activity) mContext).runOnUiThread(new Runnable() {

                                    @Override
                                    public void run() {

                                        Alert_Dialog alert_Dialog = new Alert_Dialog(mContext, null, ((Activity) mContext).findViewById(R.id.alert_dialog), ErrorUtil.getErrorString(error), BitmapFactory.decodeResource(mContext.getResources(), R.drawable.wrong));
                                        alert_Dialog.show();

                                    }

                                });
                            }

                        }

                        @Override
                        public void onError(final String error) {

                            ((Activity) mContext).runOnUiThread(new Runnable() {

                                @Override
                                public void run() {

                                    Alert_Dialog alert_Dialog = new Alert_Dialog(mContext, null, ((Activity) mContext).findViewById(R.id.alert_dialog), ErrorUtil.getErrorString(error), BitmapFactory.decodeResource(mContext.getResources(), R.drawable.wrong));
                                    alert_Dialog.show();

                                }

                            });
                        }
                    });
                } else {

                    String chcd = "ALP";
                    if (btnweixin.isChecked()) {
                        chcd = "WXP";
                    } else {
                        chcd = "ALP";
                    }
                    Log.i("opp", "" + chcd);

                    if (switch_flag) {
                        Intent intent = new Intent(mContext, CaptureActivity.class);
                        intent.putExtra("chcd", chcd);
                        intent.putExtra("total", "" + sum);
                        mContext.startActivity(intent);

                    } else {
                        Intent intent = new Intent(mContext, CreateQRcodeActivity.class);
                        intent.putExtra("chcd", chcd);
                        intent.putExtra("total", "" + sum);
                        mContext.startActivity(intent);

                    }

                }


                break;
        }

    }

    public void getResult() {

        double result = 0;
        String x = txt_output.getText().toString();
        String t = "";
        int i = 0;

        if (x.indexOf("+") == -1) {
            result = Double.parseDouble(x);
            edt_input.setText("=" + String.format("%.2f", result));
        } else {

            while (x.contains("+")) {

                t = x.substring(0, x.indexOf("+"));
                x = x.substring(x.indexOf("+") + 1);
                s[i] = t;
                i++;

            }
            s[i] = x;
            i++;
            for (int c = 0; c < i; c++) {

                result += Double.parseDouble(s[c]);
            }
            edt_input.setText("=" + String.format("%.2f", result));
        }

        // String r=edt_input.getText().toString();

        if (result > 99999999) {
            Toast.makeText(mContext, "金额过大!", Toast.LENGTH_SHORT).show();
            num_flag = false;

        } else {
            num_flag = true;
        }

    }

    public void getResult(String w) {

        double result = 0;
        String x = w;
        String t = "";
        int i = 0;

        while (x.contains("+")) {

            t = x.substring(0, x.indexOf("+"));
            x = x.substring(x.indexOf("+") + 1);
            s[i] = t;
            i++;

        }
        s[i] = x;
        i++;
        for (int c = 0; c < i; c++) {

            result += Double.parseDouble(s[c]);
        }
        edt_input.setText("=" + String.format("%.2f", result));

        // String r=edt_input.getText().toString();
        if (result > 99999999) {
            Toast.makeText(mContext, "金额过大!", Toast.LENGTH_SHORT).show();
            num_flag = false;

        } else {
            num_flag = true;
        }

    }

    public void clearzero() {

        if (clear_flag) {
            txt_output.setText("");
            clear_flag = false;
        }

    }

    public void clearzero(String z) {

        clear_flag = false;

    }

    public void add_falg(String q) {
        if (add_flag) {
            String x = txt_output.getText().toString();
            if (x.contains(".")) {
                String k = x.substring(x.lastIndexOf("."));
                if (k.equals("."))

                {
                    txt_output.setText(x.substring(0, x.length() - 1));
                    txt_output.append(q);
                    add_flag = false;
                    point_flag = true;

                } else {
                    txt_output.append(q);
                    add_flag = false;
                    point_flag = true;
                }
            } else {

                txt_output.append(q);
                add_flag = false;
                point_flag = true;

            }
        } else {
            return;
        }

    }

    public void point_falg(String q) {
        if (point_flag) {
            String x = txt_output.getText().toString();
            if (x.contains("+")) {
                String k = x.substring(x.lastIndexOf("+"));
                if (k.equals("+"))

                {
                    txt_output.append("0" + q);
                    point_flag = false;
                    add_flag = true;
                } else {

                    txt_output.append(q);
                    point_flag = false;
                    add_flag = true;
                }
            } else {

                txt_output.append(q);
                point_flag = false;
                add_flag = true;
            }
        } else {
            return;
        }

    }

    public void addcheck() {
        String x = txt_output.getText().toString();
        if (x.contains("+")) {
            String k = x.substring(x.lastIndexOf("+"));
            if (k.equals("+"))

            {
                x = x.substring(0, x.lastIndexOf("+"));
                getResult(x);
                return;
            }
            getResult();

        } else if (x.length() == 0) {

            txt_output.setText(0 + "");
            edt_input.setText("=0");
            clear_flag = true;
        } else {
            getResult(x);
        }

    }

    public class Item {

        public double value = 0;
        public int type = 0;

        public Item(double value, int type) {
            this.value = value;
            this.type = type;

        }

    }

}
