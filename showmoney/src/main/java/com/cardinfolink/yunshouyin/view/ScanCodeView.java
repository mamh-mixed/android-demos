package com.cardinfolink.yunshouyin.view;

import android.content.Context;
import android.content.Intent;
import android.content.SharedPreferences;
import android.graphics.Bitmap;
import android.graphics.BitmapFactory;
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

import com.cardinfolink.cashiersdk.listener.CashierListener;
import com.cardinfolink.cashiersdk.model.OrderData;
import com.cardinfolink.cashiersdk.model.ResultData;
import com.cardinfolink.cashiersdk.sdk.CashierSdk;
import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.activity.CaptureActivity;
import com.cardinfolink.yunshouyin.activity.PayResultActivity;
import com.cardinfolink.yunshouyin.api.QuickPayException;
import com.cardinfolink.yunshouyin.constant.Msg;
import com.cardinfolink.yunshouyin.core.QuickPayCallbackListener;
import com.cardinfolink.yunshouyin.core.QuickPayService;
import com.cardinfolink.yunshouyin.data.Coupon;
import com.cardinfolink.yunshouyin.data.SessonData;
import com.cardinfolink.yunshouyin.data.TradeBill;
import com.cardinfolink.yunshouyin.data.User;
import com.cardinfolink.yunshouyin.util.ErrorUtil;
import com.cardinfolink.yunshouyin.util.Log;
import com.cardinfolink.yunshouyin.util.ShowMoneyApp;
import com.cardinfolink.yunshouyin.util.Utility;
import com.google.zxing.WriterException;

import java.math.BigDecimal;
import java.text.SimpleDateFormat;
import java.util.Date;

public class ScanCodeView extends LinearLayout implements View.OnClickListener, View.OnTouchListener {
    private static final String TAG = "ScanCodeView";
    private static final double MAX_MONEY = 99999999.99;//能够进行交易的最大金额数
    private TranslateAnimation mShowAnimation;
    private TranslateAnimation mHideAnimation;
    private static final int SPLASH_DISPLAY_LENGHT = 2000;

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
    private TextView output;//上边的文本框
    private TextView mHasDiscount;//提示有没有折扣价

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


    private ResultData mResultData;
    private Handler mMainActivityHandler;//来自MainActivity的handler，注意区分mHandler
    private Handler mHandler;
    private String mOrderNum;//原始订单,取消订单时会用到

    private SharedPreferences sp;

    private int mLastDownY = 0;


    private HintDialog mHintDialog;
    private YellowTips mYellowTips;


    private BigDecimal mOriginalTotal = new BigDecimal("0");//原始金额
    private double mTotal;//优惠后的金额,实际支付的金额，如果有优惠就是优惠后的金额。如果没有优惠就和原始金额是一样的
    private String mCurrentTime;


    public ScanCodeView(Context context) {
        this(context, null);
    }

    public ScanCodeView(Context context, Handler handler) {
        super(context);
        mContext = context;
        mMainActivityHandler = handler;
        //初始化SharedPreferences sp
        sp = mContext.getSharedPreferences("savedata", Context.MODE_PRIVATE);
        mCHCD = sp.getString("CHCD", "WXP");//默认是微信支付。每次用户切换都记录一下。

        View contentView = LayoutInflater.from(context).inflate(R.layout.scancode_view, null);
        LinearLayout.LayoutParams layoutParams = new LayoutParams(LayoutParams.MATCH_PARENT, LayoutParams.MATCH_PARENT);
        contentView.setLayoutParams(layoutParams);
        addView(contentView);

        initHandler();

        initAnimation();//初始化一个动画

        initLayout();//初始各个改组件
    }


    // 初始化各个组件
    private void initLayout() {
        mHintDialog = new HintDialog(mContext, findViewById(R.id.hint_dialog));
        mYellowTips = new YellowTips(mContext, findViewById(R.id.yellow_tips));

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
        mHasDiscount = (TextView) findViewById(R.id.tv_hasdiscount);

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
        output = (TextView) findViewById(R.id.output);//上边的文本框

        btn0.setOnClickListener(new NumberOnClickListener("0"));
        btn1.setOnClickListener(new NumberOnClickListener("1"));
        btn2.setOnClickListener(new NumberOnClickListener("2"));
        btn3.setOnClickListener(new NumberOnClickListener("3"));
        btn4.setOnClickListener(new NumberOnClickListener("4"));
        btn5.setOnClickListener(new NumberOnClickListener("5"));
        btn6.setOnClickListener(new NumberOnClickListener("6"));
        btn7.setOnClickListener(new NumberOnClickListener("7"));
        btn8.setOnClickListener(new NumberOnClickListener("8"));
        btn9.setOnClickListener(new NumberOnClickListener("9"));
        btnadd.setOnClickListener(new AddOnClickListener());
        btnpoint.setOnClickListener(new PointOnClickListener());
        btnclear.setOnClickListener(new ClearOnClickListener());
        btndelete.setOnClickListener(new DelOnClickListener());

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

    //初始化动画
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

    //获取 mhintdialog 实例，为了ScanCodeActivity里面用的，真他妈的瞎搞！
    public HintDialog getmHintDialog() {
        return mHintDialog;
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
        input.setText("0");
        output.setText("0");
        addFlag = true;
        pointFlag = true;
        clearFlag = true;
        mHasDiscount.setVisibility(View.INVISIBLE);
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
        if ("true".equals(user.getLimit())) {
            quickPayService.getTotalAsync(user, date, new QuickPayCallbackListener<String>() {
                @Override
                public void onSuccess(String data) {
                    if (TextUtils.isEmpty(data)) {
                        data = "0";
                    }
                    try {
                        BigDecimal maxLimitBg = new BigDecimal(SessonData.loginUser.getLimitAmt());
                        BigDecimal limitValueBg = new BigDecimal(data);
                        if (limitValueBg.compareTo(maxLimitBg) >= 0) {
                            //"当日交易已超过限额,请申请提升限额!";
                            String alertMsg = ShowMoneyApp.getResString(R.string.alert_error_limit_error);
                            mYellowTips.show(alertMsg);
                        } else {
                            captureOrCreate.start();//这


                        }
                    } catch (Exception e) {
                    }
                }

                @Override
                public void onFailure(QuickPayException ex) {
                    endLoading();
                    String errorMsg = ex.getErrorMsg();
//                    mScanTitle.setText(errorMsg);
                    mYellowTips.show(ErrorUtil.getErrorString(errorMsg));
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
        int color = mContext.getResources().getColor(R.color.background_scan_qrcode_layout_wexin);
        scanCodeView.setBackgroundColor(color);//设置背景颜色
        mLeftText.setVisibility(INVISIBLE);
        mLeftImage.setImageResource(R.drawable.scan_left_disable);
        mRightText.setVisibility(VISIBLE);
        mRightImage.setImageResource(R.drawable.scan_right_able);
        mScanTitle.setText(getResources().getString(R.string.create_qrcode_activity_open_wx));
    }

    private void setRight() {
        SharedPreferences.Editor editor = sp.edit();
        editor.putString("CHCD", mCHCD);
        editor.commit();
        int color = mContext.getResources().getColor(R.color.background_scan_qrcode_layout_ali);
        scanCodeView.setBackgroundColor(color);//设置背景颜色
        mLeftText.setVisibility(VISIBLE);
        mLeftImage.setImageResource(R.drawable.scan_left_able);
        mRightText.setVisibility(INVISIBLE);
        mRightImage.setImageResource(R.drawable.scan_right_disable);
        mScanTitle.setText(getResources().getString(R.string.create_qrcode_activity_open_ali));
    }


    public void startLoading() {
        //二维码的loading的动画
        mQRImage.setImageResource(R.drawable.loading);
        startLoading(mQRImage);
    }

    public void startLoading(View view) {
        Animation loadingAnimation = AnimationUtils.loadAnimation(mContext, R.anim.loading_animation);
        view.startAnimation(loadingAnimation);
    }

    public void endLoading(View view) {
        view.clearAnimation();
    }

    public void endLoading() {
        //二维码的loading的动画
        endLoading(mQRImage);
    }

    private boolean validate(double sum) {
        if (sum < 0.01) {
            //"金额不能少于 0.01元!"
            String toastMsg = ShowMoneyApp.getResString(R.string.toast_money_cannot_zero);
            mYellowTips.show(toastMsg);
            return false;
        }
        if (sum > MAX_MONEY) {
            //金额太大了
            String toastMsg = ShowMoneyApp.getResString(R.string.toast_money_too_large);
            mYellowTips.show(toastMsg);
            return false;
        }
        return true;
    }

    private void startQRPay(final double total, final double originaiTotal) {

        boolean hasDiscount = (Coupon.getInstance().getSaleDiscount() != null)
                && (!"0".equals(Coupon.getInstance().getSaleDiscount()));
        if (!validate(total)) {
            return;
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
                createOrder(String.valueOf(total), String.valueOf(originaiTotal), mCHCD);
            }
        });

        showScanCode();//显示二维码界面
    }

    @Override
    public boolean onTouch(View v, MotionEvent event) {
        mTotal = Double.parseDouble(input.getText().toString());

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
                        startQRPay(mTotal, mOriginalTotal.doubleValue());
                    }
                }
                break;
        }
        return false;
    }

    private void startCapturePay(final double total, final double originaltotal) {
        boolean hasDiscount = (Coupon.getInstance().getSaleDiscount() != null)
                && (!"0".equals(Coupon.getInstance().getSaleDiscount()));
        if (!validate(total)) {
            return;
        }

        checkLimit(new CheckLimitInterface() {
            @Override
            public void start() {
                //进入照相机扫码界面
                Intent intent = new Intent(mContext, CaptureActivity.class);
                //这里要传人 支付类型，是微信还是支付宝,这里不需要传人支付类型了，服务器判断。
                Bundle bundle = new Bundle();
                bundle.putString("chcd", mCHCD); //这里要传人 支付类型，是微信还是支付宝
                bundle.putString("total", "" + total);//实际支付
                bundle.putString("originaltotal", "" + originaltotal);
                bundle.putString("original", "scancodeview");
                intent.putExtras(bundle);
                mContext.startActivity(intent);

                Message msg = Message.obtain();
                msg.what = Msg.MSG_SCAN_CODE_VIEW_CLEAR_INPUT_OUTPUT;
                mHandler.sendMessageDelayed(msg, SPLASH_DISPLAY_LENGHT);
            }
        });
    }

    private class NumberOnClickListener implements OnClickListener {
        private String number = "";

        public NumberOnClickListener(String number) {
            this.number = number;
        }

        @Override
        public void onClick(View v) {
            if (numFlag) {
                clearZero();
                output.append(this.number);
                getResult();
                addFlag = true;
            }
        }
    }

    private class PointOnClickListener implements OnClickListener {

        @Override
        public void onClick(View v) {
            String outputText = output.getText().toString();
            String s1 = outputText.substring(outputText.lastIndexOf("+") + 1);
            if (s1.contains(".")) {
                return;
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
        }
    }

    private class AddOnClickListener implements OnClickListener {

        @Override
        public void onClick(View v) {
            String outputText = output.getText().toString();
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
        }
    }

    private class DelOnClickListener implements OnClickListener {

        @Override
        public void onClick(View v) {
            String outputText = output.getText().toString();
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
                if ("0".equals(outputText.toString())) {
                    mHasDiscount.setVisibility(View.INVISIBLE);
                }
            } else {
                mHasDiscount.setVisibility(View.INVISIBLE);
            }
            addCheck();
            numFlag = true;
        }
    }

    private class ClearOnClickListener implements OnClickListener {

        @Override
        public void onClick(View v) {
            clearValue();
        }
    }

    @Override
    public void onClick(View v) {
        mTotal = Double.parseDouble(input.getText().toString());
        switch (v.getId()) {
            case R.id.scancodepay:
                startCapturePay(mTotal, mOriginalTotal.doubleValue());
                break;
            case R.id.iv_left:
                //切换了支付方式
                if (!mCHCD.equals(CHCD_TYPE[0])) {
                    cancelBill(mOrderNum);
                    mCHCD = CHCD_TYPE[0];
                    startQRPay(mTotal, mOriginalTotal.doubleValue());
                }
                break;
            case R.id.iv_right:
                //切换了支付方式
                if (!mCHCD.equals(CHCD_TYPE[1])) {
                    cancelBill(mOrderNum);
                    mCHCD = CHCD_TYPE[1];
                    startQRPay(mTotal, mOriginalTotal.doubleValue());
                }
                break;
            case R.id.scan_qr:
                cancelBill(mOrderNum);
                showKeyBoard();
                startCapturePay(mTotal, mOriginalTotal.doubleValue());
                break;
            case R.id.iv_keyboard:
                cancelBill(mOrderNum);
                showKeyBoard();//显示键盘界面
                break;
        }
    }


    /**
     * 生成二维码入口，从这里进的
     * 生成二维码，也就是 调用了 预下单 的接口
     *
     * @param total
     * @param chcd
     */
    private void createOrder(String total, String originalTotal, String chcd) {
        final OrderData orderData = new OrderData();
        if (!total.equals(originalTotal)) {
            orderData.payType = Coupon.getInstance().getPayType();
            orderData.discountMoney = new BigDecimal(originalTotal).subtract(new BigDecimal(total)).toString();
            orderData.couponOrderNum = Coupon.getInstance().getScanCodeId();
        }

        SimpleDateFormat mspf = new SimpleDateFormat("yyyy-MM-dd HH:mm:ss");
        mCurrentTime = mspf.format(new Date());

        orderData.orderNum = Utility.geneOrderNumber();
        orderData.txamt = total;
        orderData.currency = CashierSdk.SDK_CURRENCY;
        orderData.chcd = chcd;

        mOrderNum = orderData.orderNum;//保存为原始订单


        if (Coupon.getInstance().getSaleDiscount() != null && !"0".equals(Coupon.getInstance().getSaleDiscount())) {
            orderData.discountMoney = String.valueOf(new BigDecimal(originalTotal).subtract(new BigDecimal(total)).doubleValue());
            orderData.couponOrderNum = Coupon.getInstance().getOrderNum();
        }
        startLoading();

        CashierSdk.startPrePay(orderData, new CashierListener() {
            @Override
            public void onResult(ResultData resultData) {
                mResultData = resultData;
                Message msg = new Message();
                msg.what = Msg.MSG_CREATE_QR_SUCCESS;
                mHandler.sendMessageDelayed(msg, 0);
            }

            @Override
            public void onError(int errorCode) {
                Message msg = new Message();
                msg.what = Msg.MSG_CREATE_QR_FAIL;
                mHandler.sendMessageDelayed(msg, 0);
            }
        });
    }

    //取消账单，不去判断取消是否成功
    private void cancelBill(final String mOrderNum) {
        if (TextUtils.isEmpty(mOrderNum)) {
            return;
        }

        final OrderData orderData = new OrderData();
        orderData.origOrderNum = mOrderNum;
        //这里先查询一下，如果还未支付再取消也不迟呀
        CashierSdk.startQy(orderData, new CashierListener() {
            @Override
            public void onResult(ResultData resultData) {
                if (resultData.respcd.equals("09")) {
                    //如果是还未支付 这时候再取消
                    orderData.origOrderNum = mOrderNum;
                    orderData.orderNum = Utility.geneOrderNumber();//新生成一个订单号
                    CashierSdk.startCanc(orderData, new CashierListener() {
                        @Override
                        public void onResult(ResultData resultData) {
                        }

                        @Override
                        public void onError(int errorCode) {
                        }
                    });
                }//end if()
            }

            @Override
            public void onError(int errorCode) {

            }
        });

    }


    /**
     * 跳转到交易成功的界面
     */
    public void enterPaySuccessActivity() {
        Intent intent = new Intent(mContext, PayResultActivity.class);
        Bundle bun = new Bundle();

        TradeBill tradeBill = new TradeBill();
        //这里是改了 还是怎么了 resultdata没有返回orderNum了?????
        tradeBill.orderNum = mResultData.origOrderNum;
        tradeBill.chcd = mResultData.chcd;
        tradeBill.tandeDate = mCurrentTime;
        tradeBill.response = "success";
        tradeBill.total = String.valueOf(mTotal);//付款金额
        boolean hasCoupon = Coupon.getInstance().getSaleDiscount() != null &&
                !"0".equals(Coupon.getInstance().getSaleDiscount());
        if (hasCoupon) {
            //有优惠卡券支付
            tradeBill.originalTotal = String.valueOf(mOriginalTotal);//消费金额
        } else {
            //无优惠卡券支付

        }
        bun.putSerializable("TradeBill", tradeBill);
        intent.putExtra("BillBundle", bun);

        mContext.startActivity(intent);
    }

    /**
     * 跳转到交易失败的界面
     */
    public void enterPayFailActivity() {
        Intent intent = new Intent(mContext, PayResultActivity.class);

        Bundle bun = new Bundle();
        TradeBill tradeBill = new TradeBill();
        //这里是改了 还是怎么了 resultdata没有返回orderNum了?????
        tradeBill.orderNum = mResultData.origOrderNum;
        tradeBill.chcd = mResultData.chcd;
        tradeBill.tandeDate = mCurrentTime;
        tradeBill.errorDetail = mResultData.errorDetail;
        tradeBill.response = "fail";
        tradeBill.total = String.valueOf(mTotal);//付款金额
        boolean hasCoupon = Coupon.getInstance().getSaleDiscount() != null &&
                !"0".equals(Coupon.getInstance().getSaleDiscount());//判断是否有优惠金额
        if (hasCoupon) {
            //有优惠卡券支付
            tradeBill.originalTotal = String.valueOf(mOriginalTotal);//消费金额
        } else {

        }
        bun.putSerializable("TradeBill", tradeBill);
        intent.putExtra("BillBundle", bun);


        mContext.startActivity(intent);
    }


    public void startPolling() {
        //开启一个线程轮询服务器5次
        PollingThread pollingThread = new PollingThread(mOrderNum);
        pollingThread.start();
    }


    private class PollingThread extends Thread {
        private int pollingCount = 0;
        private String orderNumInThread;
        private OrderData orderDataInThread = new OrderData();


        public PollingThread(String orderNum) {
            this.orderNumInThread = orderNum;
            this.orderDataInThread.origOrderNum = orderNum;
        }

        @Override
        public void run() {
            while (true) {
                try {
                    pollingCount++;
                    if (pollingCount >= 10) {
                        //取消订单
                        cancelBill(orderNumInThread);
                        break;
                    }
                    //子线程内部的 订单号 和外面的 订单号,这里我们自己创建线程，自己来控制查询结果
                    if (orderNumInThread.equals(mOrderNum)) {
                        ResultData resultData = CashierSdk.startQy(this.orderDataInThread);
                        if (resultData != null) {
                            mResultData = resultData;
                            if ("00".equals(resultData.respcd)) {
                                mHandler.sendEmptyMessage(Msg.MSG_FROM_SERVER_TRADE_SUCCESS);
                                break;
                            } else if ("09".equals(resultData.respcd)) {
                                //09 状态
                            } else if ("H3".equals(resultData.respcd)) {
                                //订单已经关闭
                                break;
                            } else {
                                mHandler.sendEmptyMessage(Msg.MSG_FROM_SERVER_TRADE_FAIL);
                                break;
                            }
                        } else {
                            break;
                        }
                    } else {
                        cancelBill(orderNumInThread);
                        break;
                    }
                    Thread.sleep(5000);
                } catch (InterruptedException e) {
                    e.printStackTrace();
                }
            }
        }

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
        min = Math.abs(min - 2 * leftArrow.getWidth());
        try {
            //创建二维码图片
            if (!TextUtils.isEmpty(mResultData.qrcode)) {
                bitmap = Utility.cretaeBitmap(mResultData.qrcode, icon, min, min);
                mQRImage.setImageBitmap(bitmap);

                startPolling();
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
                    case Msg.MSG_CREATE_QR_SUCCESS: {//这里是创建二维码成功
                        if (mResultData != null) {
                            if (mResultData.respcd.equals("00") || mResultData.respcd.equals("09")) {
                                updateQR();//更新二维码
                            } else {
                                //这里返回其他，说明出错了
                                mScanTitle.setText(ErrorUtil.getErrorString(mResultData.errorDetail));
                                endLoading();
                            }
                        }
                        break;
                    }
                    case Msg.MSG_CREATE_QR_FAIL: {
                        endLoading();//出错结束loading
                        mQRImage.setImageResource(R.drawable.wrong);
                        break;
                    }
                    case Msg.MSG_FROM_SERVER_TRADE_SUCCESS: {
                        showKeyBoard();//显示键盘界面
                        enterPaySuccessActivity();
                        break;
                    }
                    case Msg.MSG_FROM_SERVER_TRADE_FAIL: {
                        showKeyBoard();//显示键盘界面
                        enterPayFailActivity();
                        break;
                    }
                    case Msg.MSG_FROM_SERVER_CLOSEBILL_SUCCESS: {
                        //关单成功//"关闭订单成功"
                        String title = mContext.getString(R.string.scancode_view_cancel_order_success);
                        mHintDialog.setTitle(title);
                        mHintDialog.setOkText(mContext.getString(R.string.scancode_view_had_cancel_order));//"已经关单了"
                        mHintDialog.setOkOnClickListener(new OnClickListener() {
                            @Override
                            public void onClick(View v) {
                                mHintDialog.hide();
                            }
                        });
                        mQRImage.setImageDrawable(null);
                        break;
                    }
                    case Msg.MSG_FROM_SERVER_CLOSEBILL_DOING: {
                        //关单返回09，表明进行中吧
                        break;
                    }
                    case Msg.MSG_FROM_SERVER_CLOSEBILL_FAIL: {
                        //关单失败
                        String title = mContext.getString(R.string.scancode_view_cancel_fail);
                        mHintDialog.setTitle(title);
                        mHintDialog.setOkText(title);//"关单失败"
                        mHintDialog.setOkOnClickListener(new OnClickListener() {
                            @Override
                            public void onClick(View v) {
                                mHintDialog.hide();
                            }
                        });
                        mQRImage.setImageDrawable(null);
                        break;
                    }
                    case Msg.MSG_SCAN_CODE_VIEW_CLEAR_INPUT_OUTPUT: {
                        clearValue();
                    }
                }
                super.handleMessage(msg);
            }
        };
    }


    public void getResult() {
        mOriginalTotal = new BigDecimal("0");
        String x = output.getText().toString();//上边的文本框


        String t = "";
        int i = 0;
        String tempInputResult;//优惠后的金额，
        if (x.indexOf("+") == -1) {
            mOriginalTotal = new BigDecimal(x);
            tempInputResult = discountMoneyResult(mOriginalTotal);
            input.setText(String.valueOf(tempInputResult));//下面的文本框
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
                mOriginalTotal = mOriginalTotal.add(new BigDecimal(s[c]));
            }
            //优惠后的金额，
            tempInputResult = discountMoneyResult(mOriginalTotal);
            input.setText(String.valueOf(tempInputResult));//下面的文本框
        }

        //添加浏览判断小数点后最多输入两位小数
        String[] nums = x.split("\\+");
        if (nums.length > 0) {
            String lastNum = nums[nums.length - 1];//按照加号分割后，取出最后一个数字
            if (lastNum != null && lastNum.contains(".")) {
                String[] floatsNum = lastNum.split("\\.");//分割字符串
                if (floatsNum.length > 0) {
                    String last = floatsNum[floatsNum.length - 1];
                    if (last != null && last.length() >= 2) {
                        numFlag = false;
                        return;
                    }
                }
            }
        }

        if (mOriginalTotal.compareTo(new BigDecimal(MAX_MONEY)) > 0) {
            // "金额过大!"
            String toastMsg = mContext.getString(R.string.toast_money_too_large);
            mYellowTips.show(toastMsg);
            numFlag = false;
        } else {
            numFlag = true;
        }

    }

    public void getResult(String w) {
        mOriginalTotal = new BigDecimal("0");

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
            mOriginalTotal = mOriginalTotal.add(new BigDecimal(s[c]));
        }
        //优惠后的金额，
        String tempInputResult = discountMoneyResult(mOriginalTotal);
        input.setText(String.valueOf(tempInputResult));//下面的文本框

        if (mOriginalTotal.compareTo(new BigDecimal(MAX_MONEY)) > 0) {
            String toastMsg = ShowMoneyApp.getResString(R.string.toast_money_too_large);
            mYellowTips.show(toastMsg);
            numFlag = false;
        } else {
            numFlag = true;
        }
    }

    //获取打折后的金额
    public String discountMoneyResult(BigDecimal result) {
        BigDecimal tempResult = result;
        boolean hasCoupon = Coupon.getInstance().getSaleDiscount() != null &&
                !"0".equals(Coupon.getInstance().getSaleDiscount());
        if (!hasCoupon) {
            mHasDiscount.setVisibility(View.INVISIBLE);
            return tempResult.toString();
        }
        //打折门限值
        BigDecimal limit = new BigDecimal("0");
        if (Coupon.getInstance().getSaleMinAmount() != null && !"0".equals(Coupon.getInstance().getSaleMinAmount())) {
            BigDecimal saleMinAmount = new BigDecimal(Double.valueOf(Coupon.getInstance().getSaleMinAmount()));
            limit = saleMinAmount.divide(new BigDecimal(100));
        }
        //折扣值
        BigDecimal saleDiscount = new BigDecimal(Double.valueOf(Coupon.getInstance().getSaleDiscount()));
        BigDecimal discount = saleDiscount.divide(new BigDecimal(100));

        //最高优惠金额
        boolean hasMaxDiscount = Coupon.getInstance().getMaxDiscountAmt() != null && !"0".equals(Coupon.getInstance().getMaxDiscountAmt());
        BigDecimal maxDiscount = new BigDecimal("0");
        if (hasMaxDiscount) {
            maxDiscount = new BigDecimal(Coupon.getInstance().getMaxDiscountAmt()).divide(new BigDecimal(100));
        }


        //满减券
        if (Coupon.getInstance().getVoucherType().endsWith("1")) {
            if (limit.compareTo(new BigDecimal(0)) > 0 && result.compareTo(limit) >= 0) {
                mHasDiscount.setVisibility(View.VISIBLE);
                tempResult = tempResult.subtract(discount);
                Log.e(TAG, tempResult + "满减");
            }
        } else if (Coupon.getInstance().getVoucherType().endsWith("2")) {
            //固定金额券
            mHasDiscount.setVisibility(View.VISIBLE);
            if (tempResult.compareTo(discount) <= 0) {
                tempResult = new BigDecimal(0);
            } else {
                tempResult.subtract(discount);
            }
            Log.e(TAG, tempResult + "固定金额");
        } else if (Coupon.getInstance().getVoucherType().endsWith("3")) {
            //满折券
            if (limit.compareTo(new BigDecimal(0)) > 0 && result.compareTo(limit) >= 0) {
                mHasDiscount.setVisibility(View.VISIBLE);
                tempResult = tempResult.multiply(discount).setScale(2, BigDecimal.ROUND_FLOOR);
                //判断优惠金额是否大于最大优惠金额
                // TODO: 2016/1/4  可能存在精度的问题
                if (hasMaxDiscount) {
                    if ((result.subtract(tempResult)).compareTo(maxDiscount) > 0) {
                        //tempResult = new BigDecimal(result).subtract(new BigDecimal(maxDiscount)).doubleValue();
                        tempResult = result.subtract(maxDiscount).setScale(2, BigDecimal.ROUND_FLOOR);
                    }
                }
            }
            Log.e(TAG, tempResult + "满折");
        }
        return tempResult.setScale(2, BigDecimal.ROUND_FLOOR).toString();
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
                    numFlag = true;
                } else {
                    output.append(q);
                    addFlag = false;
                    pointFlag = true;
                    numFlag = true;
                }
            } else {
                output.append(q);
                addFlag = false;
                pointFlag = true;
                numFlag = true;
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
            output.setText("0");
            input.setText("0");
            clearFlag = true;
        } else {
            getResult(x);
        }

    }


}
