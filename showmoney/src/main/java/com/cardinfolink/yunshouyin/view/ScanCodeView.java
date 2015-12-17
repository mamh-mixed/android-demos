package com.cardinfolink.yunshouyin.view;

import android.app.Activity;
import android.content.Context;
import android.content.Intent;
import android.graphics.Bitmap;
import android.graphics.BitmapFactory;
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
import com.cardinfolink.yunshouyin.activity.CaptureActivity;
import com.cardinfolink.yunshouyin.activity.CreateQRcodeActivity;
import com.cardinfolink.yunshouyin.api.QuickPayException;
import com.cardinfolink.yunshouyin.core.QuickPayCallbackListener;
import com.cardinfolink.yunshouyin.core.QuickPayService;
import com.cardinfolink.yunshouyin.data.SessonData;
import com.cardinfolink.yunshouyin.data.User;
import com.cardinfolink.yunshouyin.util.ShowMoneyApp;

import java.text.SimpleDateFormat;
import java.util.Date;

public class ScanCodeView extends LinearLayout implements OnClickListener {
    private Button btn0, btn1, btn2, btn3, btn4, btn5, btn6,
            btn7, btn8, btn9, btnadd, btnpoint, btnsm, btnclear, btndelete,
            swh;
    private RadioButton btnzhifubao, btnweixin;
    private EditText input;
    private TextView output;
    private boolean clearFlag = true;
    private boolean pointFlag = true;
    private boolean addFlag = true;
    private boolean switchFlag = true;
    private boolean numFlag = true;
    private String[] s = new String[100];
    private Context mContext;

    public ScanCodeView(Context context) {
        super(context);
        mContext = context;
        View contentView = LayoutInflater.from(context).inflate(R.layout.scancode_view, null);
        LinearLayout.LayoutParams layoutParams = new LayoutParams(LayoutParams.MATCH_PARENT, LayoutParams.MATCH_PARENT);
        contentView.setLayoutParams(layoutParams);
        addView(contentView);
        initLayout();
    }

    public void clearValue() {
        numFlag = true;
        input.setText("=0");
        output.setText("0");
        addFlag = true;
        pointFlag = true;
        clearFlag = true;
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

        input = (EditText) findViewById(R.id.edt_input);
        output = (TextView) findViewById(R.id.txt_output);

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
                if (switchFlag) {
                    btnsm.setText(ShowMoneyApp.getResString(R.string.scancode_view_create_code));
                    switchFlag = false;
                } else {
                    btnsm.setText(ShowMoneyApp.getResString(R.string.scancode_view_scaning_code));
                    switchFlag = true;
                }

            }
        });

    }

    @Override
    public void onClick(View v) {
        String outputStr = output.getText().toString();
        switch (v.getId()) {
            case R.id.btn0:
                if (numFlag) {
                    clearzero();
                    output.append("0");
                    addFlag = true;
                    getResult();
                }
                break;
            case R.id.btn1:
                if (numFlag) {
                    clearzero();
                    output.append("1");
                    getResult();
                    addFlag = true;
                }
                break;
            case R.id.btn2:
                if (numFlag) {
                    clearzero();
                    output.append("2");
                    getResult();
                    addFlag = true;
                }
                break;
            case R.id.btn3:
                if (numFlag) {
                    clearzero();
                    output.append("3");
                    getResult();
                    addFlag = true;
                }
                break;
            case R.id.btn4:
                if (numFlag) {
                    clearzero();
                    output.append("4");
                    getResult();
                    addFlag = true;
                }
                break;
            case R.id.btn5:
                if (numFlag) {
                    clearzero();
                    output.append("5");
                    getResult();
                    addFlag = true;
                }
                break;
            case R.id.btn6:
                if (numFlag) {
                    clearzero();
                    output.append("6");
                    getResult();
                    addFlag = true;
                }
                break;
            case R.id.btn7:
                if (numFlag) {
                    clearzero();
                    output.append("7");
                    getResult();
                    addFlag = true;
                }
                break;
            case R.id.btn8:
                if (numFlag) {
                    clearzero();
                    output.append("8");
                    addFlag = true;
                    getResult();
                }
                break;
            case R.id.btn9:
                if (numFlag) {
                    clearzero();
                    output.append("9");
                    getResult();
                    addFlag = true;
                }
                break;
            case R.id.btnpoint:
                String s1 = outputStr.substring(outputStr.lastIndexOf("+") + 1);
                if (s1.contains(".")) {
                    break;
                }

                if (outputStr.contains(".")) {
                    String k = outputStr.substring(outputStr.lastIndexOf("."));
                    if (k.equals(".")) {
                        return;
                    } else {
                        clearzero(".");
                        pointFlag(".");
                    }
                } else {
                    clearzero(".");
                    pointFlag(".");
                }
                break;
            case R.id.btnadd:
                if (outputStr.contains("+")) {
                    String k = outputStr.substring(outputStr.lastIndexOf("+"));
                    if (k.equals("+")) {
                        return;
                    } else {
                        clearzero("+");
                        addFlag("+");
                    }
                } else {
                    clearzero("+");
                    addFlag("+");
                }
                break;
            case R.id.btnclear:
                numFlag = true;
                input.setText("=0");
                output.setText("0");
                addFlag = true;
                pointFlag = true;
                clearFlag = true;
                break;
            case R.id.btndelete:
                String r = input.getText().toString();
                addFlag = true;
                if (outputStr.contains(".")) {
                    String k = outputStr.substring(outputStr.lastIndexOf("."));
                    if (k.equals(".")) {
                        pointFlag = true;
                    }
                }
                if (outputStr != null && !outputStr.equals("")) {
                    String k = outputStr.substring(outputStr.lastIndexOf("+") + 1);
                    output.setText(outputStr.substring(0, outputStr.length() - 1));
                    if (outputStr.contains("+")) {
                        if (k.equals("+")) {
                            addFlag = false;
                        } else {
                            addFlag = true;
                        }
                    } else {
                        addFlag = true;
                    }
                }
                addCheck();
                break;
            case R.id.btnzhifubao:
                break;
            case R.id.btnweixin:
                break;
            case R.id.btnsm:
                final String sumString = input.getText().toString().substring(1);//传递时候用string类型
                final double sum = Double.parseDouble(sumString);//比较大小用double类型
                if (sum <= 0) {
                    //"金额不能为零!"
                    String toastMsg = ShowMoneyApp.getResString(R.string.toast_money_cannot_zero);
                    Toast.makeText(mContext, toastMsg, Toast.LENGTH_SHORT).show();
                    return;
                }
                if (sum > 99999999) {
                    // "金额过大!"
                    String toastMsg = ShowMoneyApp.getResString(R.string.toast_money_too_large);
                    Toast.makeText(mContext, toastMsg, Toast.LENGTH_SHORT).show();
                    return;
                }
                QuickPayService quickPayService = ShowMoneyApp.getInstance().getQuickPayService();
                String date = (new SimpleDateFormat("yyyyMMdd")).format(new Date());
                User user = SessonData.loginUser;
                if (user.getLimit().equals("true")) {
                    quickPayService.getTotalAsync(user, date, new QuickPayCallbackListener<String>() {
                        @Override
                        public void onSuccess(String data) {
                            double limitValue = Double.parseDouble(data);
                            if (limitValue >= 500) {
                                //"当日交易已超过限额,请申请提升限额!";
                                String alertMsg = ShowMoneyApp.getResString(R.string.alert_error_limit_error);
                                View alertView = ((Activity) mContext).findViewById(R.id.alert_dialog);
                                Bitmap alertBitmap = BitmapFactory.decodeResource(mContext.getResources(), R.drawable.wrong);
                                AlertDialog alertDialog = new AlertDialog(mContext, null, alertView, alertMsg, alertBitmap);
                                alertDialog.show();
                            } else {
                                String chcd = "ALP";
                                if (btnweixin.isChecked()) {
                                    chcd = "WXP";
                                } else {
                                    chcd = "ALP";
                                }
                                if (switchFlag) {
                                    Intent intent = new Intent(mContext, CaptureActivity.class);
                                    intent.putExtra("chcd", chcd);
                                    intent.putExtra("total", sumString);
                                    mContext.startActivity(intent);
                                } else {
                                    Intent intent = new Intent(mContext, CreateQRcodeActivity.class);
                                    intent.putExtra("chcd", chcd);
                                    intent.putExtra("total",sumString);
                                    mContext.startActivity(intent);
                                }
                            }
                        }

                        @Override
                        public void onFailure(QuickPayException ex) {
                            String errorMsg = ex.getErrorMsg();
                            View alertView = ((Activity) mContext).findViewById(R.id.alert_dialog);
                            Bitmap alertBitmap = BitmapFactory.decodeResource(mContext.getResources(), R.drawable.wrong);
                            AlertDialog alertDialog = new AlertDialog(mContext, null, alertView, errorMsg, alertBitmap);
                            alertDialog.show();
                        }
                    });
                } else {
                    String chcd = "ALP";
                    if (btnweixin.isChecked()) {
                        chcd = "WXP";
                    } else {
                        chcd = "ALP";
                    }

                    if (switchFlag) {
                        Intent intent = new Intent(mContext, CaptureActivity.class);
                        intent.putExtra("chcd", chcd);
                        intent.putExtra("total", sumString);
                        mContext.startActivity(intent);
                    } else {
                        Intent intent = new Intent(mContext, CreateQRcodeActivity.class);
                        intent.putExtra("chcd", chcd);
                        intent.putExtra("total", sumString);
                        mContext.startActivity(intent);
                    }
                }

                break;
        }

    }


    public void getResult() {
        double result = 0;
        String x = output.getText().toString();
        String t = "";
        int i = 0;

        if (x.indexOf("+") == -1) {
            result = Double.parseDouble(x);
            input.setText("=" + String.format("%.2f", result));
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
            input.setText("=" + String.format("%.2f", result));
        }


        if (result > 99999999) {
            // "金额过大!"
            String toastMsg = ShowMoneyApp.getResString(R.string.toast_money_too_large);
            Toast.makeText(mContext, toastMsg, Toast.LENGTH_SHORT).show();
            numFlag = false;
        } else {
            numFlag = true;
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
        input.setText("=" + String.format("%.2f", result));

        if (result > 99999999) {
            Toast.makeText(mContext, "金额过大!", Toast.LENGTH_SHORT).show();
            numFlag = false;
        } else {
            numFlag = true;
        }
    }

    public void clearzero() {
        if (clearFlag) {
            output.setText("");
            clearFlag = false;
        }
    }

    public void clearzero(String z) {
        clearFlag = false;
    }

    public void addFlag(String q) {
        if (addFlag) {
            String x = output.getText().toString();
            if (x.contains(".")) {
                String k = x.substring(x.lastIndexOf("."));
                if (k.equals(".")) {
                    output.setText(x.substring(0, x.length() - 1));
                    output.append(q);
                    addFlag = false;
                    pointFlag = true;
                } else {
                    output.append(q);
                    addFlag = false;
                    pointFlag = true;
                }
            } else {
                output.append(q);
                addFlag = false;
                pointFlag = true;
            }
        } else {
            return;
        }

    }

    public void pointFlag(String q) {
        if (pointFlag) {
            String x = output.getText().toString();
            if (x.contains("+")) {
                String k = x.substring(x.lastIndexOf("+"));
                if (k.equals("+")) {
                    output.append("0" + q);
                    pointFlag = false;
                    addFlag = true;
                } else {
                    output.append(q);
                    pointFlag = false;
                    addFlag = true;
                }
            } else {
                output.append(q);
                pointFlag = false;
                addFlag = true;
            }
        } else {
            return;
        }

    }

    public void addCheck() {
        String x = output.getText().toString();
        if (x.contains("+")) {
            String k = x.substring(x.lastIndexOf("+"));
            if (k.equals("+")) {
                x = x.substring(0, x.lastIndexOf("+"));
                getResult(x);
                return;
            }
            getResult();
        } else if (x.length() == 0) {
            output.setText(0 + "");
            input.setText("=0");
            clearFlag = true;
        } else {
            getResult(x);
        }

    }

}
