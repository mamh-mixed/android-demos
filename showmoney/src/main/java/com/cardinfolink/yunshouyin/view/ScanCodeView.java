package com.cardinfolink.yunshouyin.view;

import android.app.Activity;
import android.content.Context;
import android.content.Intent;
import android.content.SharedPreferences;
import android.graphics.Bitmap;
import android.graphics.BitmapFactory;
import android.graphics.Color;
import android.os.Bundle;
import android.os.Handler;
import android.os.Message;
import android.text.TextUtils;
import android.view.LayoutInflater;
import android.view.MotionEvent;
import android.view.View;
import android.view.animation.Animation;
import android.view.animation.AnimationUtils;
import android.view.animation.TranslateAnimation;
import android.widget.ImageView;
import android.widget.LinearLayout;
import android.widget.TextView;
import android.widget.Toast;

import com.cardinfolink.cashiersdk.listener.CashierListener;
import com.cardinfolink.cashiersdk.model.OrderData;
import com.cardinfolink.cashiersdk.model.ResultData;
import com.cardinfolink.cashiersdk.sdk.CashierSdk;
import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.activity.CaptureActivity;
import com.cardinfolink.yunshouyin.api.QuickPayException;
import com.cardinfolink.yunshouyin.constant.Msg;
import com.cardinfolink.yunshouyin.core.QuickPayCallbackListener;
import com.cardinfolink.yunshouyin.core.QuickPayService;
import com.cardinfolink.yunshouyin.data.SessonData;
import com.cardinfolink.yunshouyin.data.User;
import com.cardinfolink.yunshouyin.util.ShowMoneyApp;
import com.cardinfolink.yunshouyin.util.Untilly;
import com.google.zxing.BarcodeFormat;
import com.google.zxing.EncodeHintType;
import com.google.zxing.MultiFormatWriter;
import com.google.zxing.WriterException;
import com.google.zxing.common.BitMatrix;
import com.google.zxing.qrcode.decoder.ErrorCorrectionLevel;

import java.text.SimpleDateFormat;
import java.util.Date;
import java.util.Hashtable;
import java.util.Random;

public class ScanCodeView extends LinearLayout implements View.OnClickListener, View.OnTouchListener {
    private static final String TAG = "ScanCodeView";
    private static final int MAX_MONEY = 99999999;//能够进行交易的最大金额数
    private static final int MAX_LIMIT_MONEY = 500;//单日限额的最大金额数

    private TranslateAnimation mShowAnimation;
    private TranslateAnimation mHideAnimation;

    private TextView mScanTitle;//二维码界面中间上边的 标题文本
    private TextView mAccount;//显示账号的

    private TextView btn0;
    private TextView btn1;
    private TextView btn2;
    private TextView btn3;
    private TextView btn4;
    private TextView btn5;
    private TextView btn6;
    private TextView btn7;
    private TextView btn8;
    private TextView btn9;
    private TextView btnadd;
    private TextView btnpoint;
    private TextView btnclear;

    private ImageView btndelete;

    private TextView input;
    private TextView output;

    private View scanCodeView;//二维码界面
    private View keyboardView;//键盘界面

    private View leftArrow;
    private View rightArrow;
    private View bottomArrow;

    private ImageView mScanCodePay;//键盘界面右下的图片,扫码收款按钮
    private ImageView mKeyBoard;//左下返回键盘界面的图片
    private ImageView mScanQR;//右下进入照相机扫码的界面
    private ImageView mQRImage;

    private ImageView mLeftImage;//切换按钮 支付宝还是微信的切换按钮
    private ImageView mRightImage;//切换按钮 支付宝还是微信的切换按钮
    private TextView mLeftText;
    //切换按钮 下面的文本
    private TextView mRightText;
    //切换按钮 下面的文本

    private Context mContext;

    private boolean clearFlag = true;
    private boolean pointFlag = true;
    private boolean addFlag = true;
    private boolean numFlag = true;
    private String[] s = new String[100];

    private String[] CHCD_TYPE = {
            "WXP", "ALP"
    };
    private String mCHCD = CHCD_TYPE[0];//默认是微信支付

    //创建二维码需要的类成员变量
    private static final int IMAGE_HALFWIDTH = 40;
    private static final int FOREGROUND_COLOR = 0xff000000;
    private static final int BACKGROUND_COLOR = 0xffffffff;
    private ResultData mResultData;
    private Handler mHandler;
    private OrderData originOrder;//原始订单,取消订单时会用到

    private SharedPreferences sp;

    private int mLastDownY = 0;

    public ScanCodeView(Context context) {
        super(context);
        mContext = context;

        //初始化SharedPreferences sp
        sp = mContext.getSharedPreferences("savedata", Context.MODE_PRIVATE);
        mCHCD = sp.getString("CHCD", "WXP");//默认是微信支付。每次用户切换都记录一下。

        View contentView = LayoutInflater.from(context).inflate(R.layout.scancode_view, null);
        LinearLayout.LayoutParams layoutParams = new LayoutParams(LayoutParams.MATCH_PARENT, LayoutParams.MATCH_PARENT);
        contentView.setLayoutParams(layoutParams);
        addView(contentView);

        initAnimation();//初始化一个动画

        initView();//初始各个改组件
    }

    /**
     * 初始化各个组件
     */
    private void initView() {
        mAccount = (TextView) findViewById(R.id.tv_account);//显示账号的
        mAccount.setText(SessonData.loginUser.getUsername());

        scanCodeView = findViewById(R.id.ll_qrcode);
        keyboardView = findViewById(R.id.ll_keyboard);

        leftArrow = findViewById(R.id.left_around_arrow);//左边的布局
        rightArrow = findViewById(R.id.right_around_arrow);//右边的布局
        bottomArrow = findViewById(R.id.bottom_around_arrow);//下面的布局

        mKeyBoard = (ImageView) findViewById(R.id.iv_keyboard);//左下角返回或者进入键盘界面的按钮
        mScanQR = (ImageView) findViewById(R.id.scan_qr);//右下角进入照相机扫码的按钮
        mScanCodePay = (ImageView) findViewById(R.id.scancodepay);//右下角进入二维码界面的按钮

        mQRImage = (ImageView) findViewById(R.id.iv_center);

        mLeftImage = (ImageView) findViewById(R.id.iv_left);
        mRightImage = (ImageView) findViewById(R.id.iv_right);

        mLeftText = (TextView) findViewById(R.id.tv_left);
        mRightText = (TextView) findViewById(R.id.tv_right);

        mScanTitle = (TextView) findViewById(R.id.scan_title);

        //初始化键盘 0到9， 加，删除， 清空等 TextView。
        btn0 = (TextView) findViewById(R.id.tv_zero);
        btn1 = (TextView) findViewById(R.id.tv_one);
        btn2 = (TextView) findViewById(R.id.tv_two);
        btn3 = (TextView) findViewById(R.id.tv_three);
        btn4 = (TextView) findViewById(R.id.tv_four);
        btn5 = (TextView) findViewById(R.id.tv_five);
        btn6 = (TextView) findViewById(R.id.tv_six);
        btn7 = (TextView) findViewById(R.id.tv_seven);
        btn8 = (TextView) findViewById(R.id.tv_eight);
        btn9 = (TextView) findViewById(R.id.tv_nine);
        btnadd = (TextView) findViewById(R.id.tv_add);
        btnpoint = (TextView) findViewById(R.id.tv_point);
        btnclear = (TextView) findViewById(R.id.tv_clear);
        btndelete = (ImageView) findViewById(R.id.iv_del);

        input = (TextView) findViewById(R.id.input);
        output = (TextView) findViewById(R.id.output);

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

        mKeyBoard.setOnClickListener(this);
        mScanCodePay.setOnClickListener(this);
        mScanQR.setOnClickListener(this);

        mLeftImage.setOnClickListener(this);
        mRightImage.setOnClickListener(this);

        //这里添加上滑的事件
        btn0.setOnTouchListener(this);
        btn1.setOnTouchListener(this);
        btn2.setOnTouchListener(this);
        btn3.setOnTouchListener(this);
        btn4.setOnTouchListener(this);
        btn5.setOnTouchListener(this);
        btn6.setOnTouchListener(this);
        btn7.setOnTouchListener(this);
        btn8.setOnTouchListener(this);
        btn9.setOnTouchListener(this);
        btnadd.setOnTouchListener(this);
        btnpoint.setOnTouchListener(this);
        btnclear.setOnTouchListener(this);
        btndelete.setOnTouchListener(this);
    }

    private void initAnimation() {
        mShowAnimation = new TranslateAnimation(
                Animation.RELATIVE_TO_SELF,
                0.0f,
                Animation.RELATIVE_TO_SELF,
                0.0f,
                Animation.RELATIVE_TO_SELF,
                1.0f,
                Animation.RELATIVE_TO_SELF,
                0.0f);
        mHideAnimation = new TranslateAnimation(
                Animation.RELATIVE_TO_SELF,
                0.0f,
                Animation.RELATIVE_TO_SELF,
                0.0f,
                Animation.RELATIVE_TO_SELF,
                0.0f,
                Animation.RELATIVE_TO_SELF,
                1.0f);
        mShowAnimation.setDuration(500);
        mHideAnimation.setDuration(500);
    }

    public void showKeyBoard() {
        //键盘不是显示的话就调用这个if里面的，让键盘界面显示
        if (keyboardView.getVisibility() != VISIBLE) {
            keyboardView.startAnimation(mShowAnimation);
            keyboardView.setVisibility(VISIBLE);
        }
        //同时要隐藏二维码界面
        if (scanCodeView.getVisibility() != GONE) {
            scanCodeView.startAnimation(mHideAnimation);
            scanCodeView.setVisibility(GONE);
        }
    }

    public void showScanCode() {
        //键盘 显示的话就调用这个if里面的，让键盘界面隐藏
        if (keyboardView.getVisibility() != GONE) {
            keyboardView.startAnimation(mHideAnimation);
            keyboardView.setVisibility(GONE);
        }
        if (scanCodeView.getVisibility() != VISIBLE) {
            scanCodeView.startAnimation(mShowAnimation);
            scanCodeView.setVisibility(VISIBLE);
        }
    }

    public void clearValue() {
        numFlag = true;
        input.setText("=0");
        output.setText("0");
        addFlag = true;
        pointFlag = true;
        clearFlag = true;
    }


    private interface CheckLimitInterface {
        void start();
    }

    /**
     * 检查限额，然后里面根据不同的 captureOrCreate来 创建二维码还是打开摄像头扫码。
     *
     * @param captureOrCreate
     */
    private void checkLimit(final CheckLimitInterface captureOrCreate) {
        QuickPayService quickPayService = ShowMoneyApp.getInstance().getQuickPayService();
        String date = (new SimpleDateFormat("yyyyMMdd")).format(new Date());
        User user = SessonData.loginUser;
        if (user.getLimit().equals("true")) {
            quickPayService.getTotalAsync(user, date, new QuickPayCallbackListener<String>() {
                @Override
                public void onSuccess(String data) {
                    double limitValue = Double.parseDouble(data);
                    if (limitValue >= MAX_LIMIT_MONEY) {
                        //"当日交易已超过限额,请申请提升限额!";
                        String alertMsg = ShowMoneyApp.getResString(R.string.alert_error_limit_error);
                        View alertView = ((Activity) mContext).findViewById(R.id.alert_dialog);
                        Bitmap alertBitmap = BitmapFactory.decodeResource(mContext.getResources(), R.drawable.wrong);
                        AlertDialog alertDialog = new AlertDialog(mContext, null, alertView, alertMsg, alertBitmap);
                        alertDialog.show();
                    } else {
                        captureOrCreate.start();//这里调用扫码还是生成二维码
                    }
                }

                @Override
                public void onFailure(QuickPayException ex) {
                    String errorMsg = ex.getErrorMsg();
                    View alertView = ((Activity) mContext).findViewById(R.id.alert_dialog);
                    Bitmap alertBitmap = BitmapFactory.decodeResource(mContext.getResources(), R.drawable.wrong);
                    AlertDialog alertDialog = new AlertDialog(mContext, null, alertView, errorMsg, alertBitmap);
                    alertDialog.show();
                    endLoading();
                    mScanTitle.setText(errorMsg);
                }
            });
        } else {
            captureOrCreate.start();//这里调用扫码还是生成二维码
        }
    }

    private void setLeft() {
        SharedPreferences.Editor editor = sp.edit();
        editor.putString("CHCD", mCHCD);
        editor.commit();
        scanCodeView.setBackgroundColor(Color.parseColor("#339933"));//设置背景颜色
        mLeftText.setVisibility(VISIBLE);//左边显示
        mLeftImage.setImageResource(R.drawable.scan_left_disable);
        mRightText.setVisibility(INVISIBLE);//右边不显示
        mRightImage.setImageResource(R.drawable.scan_right_able);
        mScanTitle.setText(getResources().getString(R.string.create_qrcode_activity_open_wx));
    }

    private void setRight() {
        SharedPreferences.Editor editor = sp.edit();
        editor.putString("CHCD", mCHCD);
        editor.commit();
        scanCodeView.setBackgroundColor(Color.parseColor("#0099ff"));//设置背景颜色
        mLeftText.setVisibility(INVISIBLE);//左边显示
        mLeftImage.setImageResource(R.drawable.scan_left_able);
        mRightText.setVisibility(VISIBLE);//右边不显示
        mRightImage.setImageResource(R.drawable.scan_right_disable);
        mScanTitle.setText(getResources().getString(R.string.create_qrcode_activity_open_ali));
    }


    public void startLoading() {
        //二维码的loading的动画
        mQRImage.setImageResource(R.drawable.loading);
        Animation loadingAnimation = AnimationUtils.loadAnimation(mContext, R.anim.loading_animation);
        mQRImage.startAnimation(loadingAnimation);
    }

    public void endLoading() {
        //二维码的loading的动画
        mQRImage.clearAnimation();
    }


    @Override
    public boolean onTouch(View v, MotionEvent event) {
        final double sum = Double.parseDouble(input.getText().toString().substring(1));

        int currentY = 0;
        switch (event.getAction()) {
            case MotionEvent.ACTION_DOWN:
                mLastDownY = (int) event.getY();
                return false;
            case MotionEvent.ACTION_UP:
                currentY = (int) event.getY();
                int dy = currentY - mLastDownY;
                if (dy < 0) {
                    if (Math.abs(dy) > keyboardView.getHeight() / 2) {
                        if (sum <= 0) {
                            //"金额不能为零!"
                            String toastMsg = ShowMoneyApp.getResString(R.string.toast_money_cannot_zero);
                            Toast.makeText(mContext, toastMsg, Toast.LENGTH_SHORT).show();
                            break;
                        }
                        if (sum > MAX_MONEY) {
                            //金额太大了
                            String toastMsg = ShowMoneyApp.getResString(R.string.toast_money_too_large);
                            Toast.makeText(mContext, toastMsg, Toast.LENGTH_SHORT).show();
                            break;
                        }

                        if (mCHCD.equals(CHCD_TYPE[0])) {
                            setLeft();//微信
                        } else {
                            setRight();//支付宝
                        }

                        startLoading();//load 二维码是的动画

                        //生成二维码比较慢？？？！！！
                        checkLimit(new CheckLimitInterface() {
                            @Override
                            public void start() {
                                createQRcode(String.valueOf(sum), mCHCD);
                            }
                        });

                        showScanCode();//显示二维码界面
                    }
                }
                break;
        }
        return false;
    }

    @Override
    public void onClick(View v) {
        String outputText = output.getText().toString();
        final double sum = Double.parseDouble(input.getText().toString().substring(1));
        switch (v.getId()) {
            case R.id.scancodepay:
                if (sum <= 0) {
                    //"金额不能为零!"
                    String toastMsg = ShowMoneyApp.getResString(R.string.toast_money_cannot_zero);
                    Toast.makeText(mContext, toastMsg, Toast.LENGTH_SHORT).show();
                    return;
                }
                if (sum > MAX_MONEY) {
                    //金额太大了
                    String toastMsg = ShowMoneyApp.getResString(R.string.toast_money_too_large);
                    Toast.makeText(mContext, toastMsg, Toast.LENGTH_SHORT).show();
                    return;
                }

                checkLimit(new CheckLimitInterface() {
                    @Override
                    public void start() {
                        //进入照相机扫码界面
                        Intent intent = new Intent(mContext, CaptureActivity.class);
                        //这里要传人 支付类型，是微信还是支付宝,这里不需要传人支付类型了，服务器判断。
                        Bundle bundle=new Bundle();
                        bundle.putString("chcd", mCHCD); //这里要传人 支付类型，是微信还是支付宝
                        bundle.putString("total", "" + sum);
                        bundle.putString("original","scancodeview");
                        intent.putExtras(bundle);
                        mContext.startActivity(intent);
                    }
                });

                break;
            case R.id.iv_left:
                //切换了支付方式
                if (!mCHCD.equals(CHCD_TYPE[0])) {
                    mCHCD = CHCD_TYPE[0];//切换了支付方式 ，就赋值给mCHCD
                    setLeft();
                    startLoading();
                    cancelOrder();//取消订单
                    checkLimit(new CheckLimitInterface() {
                        @Override
                        public void start() {
                            //这里 生成二维码
                            createQRcode(String.valueOf(sum), mCHCD);
                        }
                    });
                }
                break;
            case R.id.iv_right:
                //切换了支付方式
                if (!mCHCD.equals(CHCD_TYPE[1])) {
                    mCHCD = CHCD_TYPE[1];//切换了支付方式 ，就赋值给mCHCD
                    setRight();
                    startLoading();
                    cancelOrder();//取消订单
                    checkLimit(new CheckLimitInterface() {
                        @Override
                        public void start() {
                            //这里 生成二维码
                            createQRcode(String.valueOf(sum), mCHCD);
                        }
                    });
                }
                break;
            case R.id.scan_qr:
                cancelOrder();//取消订单
                showKeyBoard();

                checkLimit(new CheckLimitInterface() {
                    @Override
                    public void start() {
                        //进入照相机扫码界面
                        Intent intent = new Intent(mContext, CaptureActivity.class);
                        Bundle bundle=new Bundle();
                        bundle.putString("chcd", mCHCD); //这里要传人 支付类型，是微信还是支付宝
                        bundle.putString("total", "" + sum);
                        bundle.putString("original","scancodeview");
                        intent.putExtras(bundle);
                        mContext.startActivity(intent);
                    }
                });
                break;
            case R.id.iv_keyboard:
                cancelOrder();//取消订单
                showKeyBoard();//显示键盘界面
                break;
            case R.id.tv_zero:
                if (numFlag) {
                    clearZero();
                    output.append("0");
                    addFlag = true;
                    getResult();
                }
                break;
            case R.id.tv_one:
                if (numFlag) {
                    clearZero();
                    output.append("1");
                    getResult();
                    addFlag = true;
                }
                break;
            case R.id.tv_two:
                if (numFlag) {
                    clearZero();
                    output.append("2");
                    getResult();
                    addFlag = true;
                }
                break;
            case R.id.tv_three:
                if (numFlag) {
                    clearZero();
                    output.append("3");
                    getResult();
                    addFlag = true;
                }
                break;
            case R.id.tv_four:
                if (numFlag) {
                    clearZero();
                    output.append("4");
                    getResult();
                    addFlag = true;
                }
                break;
            case R.id.tv_five:
                if (numFlag) {
                    clearZero();
                    output.append("5");
                    getResult();
                    addFlag = true;
                }
                break;
            case R.id.tv_six:
                if (numFlag) {
                    clearZero();
                    output.append("6");
                    getResult();
                    addFlag = true;
                }
                break;
            case R.id.tv_seven:
                if (numFlag) {
                    clearZero();
                    output.append("7");
                    getResult();
                    addFlag = true;
                }
                break;
            case R.id.tv_eight:
                if (numFlag) {
                    clearZero();
                    output.append("8");
                    addFlag = true;
                    getResult();
                }
                break;
            case R.id.tv_nine:
                if (numFlag) {
                    clearZero();
                    output.append("9");
                    getResult();
                    addFlag = true;
                }
                break;
            case R.id.tv_point:
                String s1 = outputText.substring(outputText.lastIndexOf("+") + 1);
                if (s1.contains(".")) {
                    break;
                }

                if (outputText.contains(".")) {
                    String k = outputText.substring(outputText.lastIndexOf("."));
                    if (k.equals(".")) {
                        return;
                    } else {
                        clearZero(".");
                        pointFlag(".");
                    }
                } else {
                    clearZero(".");
                    pointFlag(".");
                }
                break;
            case R.id.tv_add:
                if (outputText.contains("+")) {
                    String k = outputText.substring(outputText.lastIndexOf("+"));
                    if (k.equals("+")) {
                        return;
                    } else {
                        clearZero("+");
                        addFlag("+");
                    }
                } else {
                    clearZero("+");
                    addFlag("+");
                }
                break;
            case R.id.tv_clear:
                numFlag = true;
                input.setText("=0");
                output.setText("0");
                addFlag = true;
                pointFlag = true;
                clearFlag = true;
                break;
            case R.id.iv_del:
                String r = input.getText().toString();
                addFlag = true;
                if (outputText.contains(".")) {
                    String k = outputText.substring(outputText.lastIndexOf("."));
                    if (k.equals(".")) {
                        pointFlag = true;
                    }
                }
                if (!TextUtils.isEmpty(outputText)) {
                    String k = outputText.substring(outputText.lastIndexOf("+") + 1);
                    output.setText(outputText.substring(0, outputText.length() - 1));
                    if (outputText.contains("+")) {
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


        if (result > MAX_MONEY) {
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

        if (result > MAX_MONEY) {
            String toastMsg = ShowMoneyApp.getResString(R.string.toast_money_too_large);
            Toast.makeText(mContext, toastMsg, Toast.LENGTH_SHORT).show();
            numFlag = false;
        } else {
            numFlag = true;
        }
    }

    public void clearZero() {
        if (clearFlag) {
            output.setText("");
            clearFlag = false;
        }
    }

    public void clearZero(String z) {
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

    /**
     * 生成账单号
     * 时间加上一个随机数
     *
     * @return
     */
    private String geneOrderNumber() {
        String mOrderNum;

        Date now = new Date();
        SimpleDateFormat spf = new SimpleDateFormat("yyMMddHHmmss");
        mOrderNum = spf.format(now);
        Random random = new Random();//订单号末尾随机的生成一个数
        for (int i = 0; i < 5; i++) {
            mOrderNum = mOrderNum + random.nextInt(10);
        }
        return mOrderNum;
    }

    /**
     * 生成二维码入口，从这里进的
     * 生成二维码，也就是 调用了 预下单 的接口
     *
     * @param total
     * @param chcd
     */
    private void createQRcode(String total, String chcd) {
        final OrderData orderData = new OrderData();
        orderData.orderNum = geneOrderNumber();
        orderData.txamt = total;
        orderData.currency = "156";
        orderData.chcd = chcd;

        originOrder = new OrderData();
        originOrder.origOrderNum = orderData.orderNum;//保存为原始订单
        originOrder.txamt = total;
        originOrder.currency = "156";
        originOrder.chcd = chcd;

        initHandler();//初始化handler

        CashierSdk.startPrePay(orderData, new CashierListener() {
            @Override
            public void onResult(ResultData resultData) {
                mResultData = resultData;
                Message msg = new Message();
                msg.what = 1;
                mHandler.sendMessageDelayed(msg, 0);
            }

            @Override
            public void onError(int errorCode) {
                endLoading();//出错结束loading
                mQRImage.setImageResource(R.drawable.wrong);
                mScanTitle.setText("Error: " + errorCode);
            }
        });
    }

    /**
     * 取消订单
     */
    public void cancelOrder() {
        if (originOrder == null) {
            return;
        }
        originOrder.orderNum = geneOrderNumber();//新生成一个订单号
        CashierSdk.startCanc(originOrder, new CashierListener() {
            @Override
            public void onResult(ResultData resultData) {
                originOrder = null;
            }

            @Override
            public void onError(int errorCode) {
                originOrder = null;
            }
        });

    }

    /**
     * 更新二维码图片
     * 在做微信支付 和 支付宝 切换会调用这个
     */
    private void updateQR() {
        endLoading();//停止load动画

        Bitmap icon = null;
        if (mResultData.chcd.equals(CHCD_TYPE[0])) {//微信支付
            icon = BitmapFactory.decodeResource(getResources(), R.drawable.scan_wechat);
            mScanTitle.setText(getResources().getString(R.string.create_qrcode_activity_open_wx));
        } else {
            icon = BitmapFactory.decodeResource(getResources(), R.drawable.scan_alipay);
            mScanTitle.setText(getResources().getString(R.string.create_qrcode_activity_open_ali));
        }
        Bitmap bitmap;
        //算出中间二维码图片最大的宽高。
        int dy = (int) (bottomArrow.getY() - mScanTitle.getY());
        int dx = (int) (rightArrow.getX() - leftArrow.getX());
        int min = Math.min(dx, dy);//求出最小的
        min = Math.abs(min);
        min = Math.abs(min - 2 * leftArrow.getWidth() - 10);
        try {
            //创建二维码图片
            if (!TextUtils.isEmpty(mResultData.qrcode)) {
                bitmap = cretaeBitmap(mResultData.qrcode, icon, min, min);
                mQRImage.setImageBitmap(bitmap);
            } else {
                mQRImage.setImageResource(R.drawable.wrong);
            }
        } catch (WriterException e) {
            e.printStackTrace();
        }

    }

    private void initHandler() {
        mHandler = new Handler() {
            @Override
            public void handleMessage(Message msg) {
                switch (msg.what) {
                    case 1: {
                        if (mResultData != null) {
                            if (mResultData.respcd.equals("00") || mResultData.respcd.equals("09")) {
                                updateQR();
                            }
                        }
                        break;
                    }

                    case 2: {
                        if (mResultData != null) {
                            if (mResultData.respcd.equals("00")) {

                            } else {

                            }
                        }
                        break;
                    }
                    case 3: {
                        Toast toast = Toast.makeText(mContext, getResources().getString(R.string.server_timeout), Toast.LENGTH_SHORT);
                        toast.show();
                    }
                    case Msg.MSG_FROM_DIGLOG_CLOSE: {
                        break;
                    }
                    case Msg.MSG_FROM_SERVER_TRADE_SUCCESS: {
                        break;
                    }
                    case Msg.MSG_FROM_SERVER_TRADE_FAIL: {
                        break;
                    }
                    case Msg.MSG_FROM_SERVER_TRADE_NOPAY: {

                        break;
                    }
                    case Msg.MSG_FROM_SUCCESS_DIGLOG_HISTORY: {

                        break;
                    }

                }
                super.handleMessage(msg);
            }
        };
    }


    /**
     * //生成bitmap 二维码图片,生成一个固定长宽都是width和height的二维码图片
     *
     * @param str
     * @param icon
     * @param widthx
     * @param heighty
     * @return
     * @throws WriterException
     */
    private Bitmap cretaeBitmap(String str, Bitmap icon, int widthx, int heighty) throws WriterException {
        icon = Untilly.zoomBitmap(icon, IMAGE_HALFWIDTH);
        Hashtable<EncodeHintType, Object> hints = new Hashtable<EncodeHintType, Object>();
        hints.put(EncodeHintType.ERROR_CORRECTION, ErrorCorrectionLevel.H);
        hints.put(EncodeHintType.CHARACTER_SET, "utf-8");
        hints.put(EncodeHintType.MARGIN, 1);
        //调用com.google.zxing里面的生成二维码的方法
        BitMatrix matrix = new MultiFormatWriter().encode(str, BarcodeFormat.QR_CODE, widthx, heighty, hints);

        int width = matrix.getWidth();
        int height = matrix.getHeight();

        int halfW = width / 2;
        int halfH = height / 2;
        int[] pixels = new int[width * height];
        for (int y = 0; y < height; y++) {
            for (int x = 0; x < width; x++) {
                if (x > halfW - IMAGE_HALFWIDTH && x < halfW + IMAGE_HALFWIDTH
                        && y > halfH - IMAGE_HALFWIDTH
                        && y < halfH + IMAGE_HALFWIDTH) {
                    pixels[y * width + x] = icon.getPixel(x - halfW + IMAGE_HALFWIDTH, y - halfH + IMAGE_HALFWIDTH);
                } else {
                    if (matrix.get(x, y)) {
                        pixels[y * width + x] = FOREGROUND_COLOR;
                    } else {
                        pixels[y * width + x] = BACKGROUND_COLOR;
                    }
                }

            }
        }
        Bitmap bitmap = Bitmap.createBitmap(width, height, Bitmap.Config.ARGB_8888);
        bitmap.setPixels(pixels, 0, width, 0, 0, width, height);

        return bitmap;
    }

    //生成bitmap图片,生成一个固定长宽都是300的二维码图片
    private Bitmap cretaeBitmap(String str, Bitmap icon) throws WriterException {
        return cretaeBitmap(str, icon, 300, 300);
    }


}
