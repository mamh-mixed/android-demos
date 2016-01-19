package com.cardinfolink.yunshouyin.activity;

import android.content.Intent;
import android.content.res.AssetFileDescriptor;
import android.graphics.Bitmap;
import android.graphics.Color;
import android.media.AudioManager;
import android.media.MediaPlayer;
import android.media.MediaPlayer.OnCompletionListener;
import android.os.Bundle;
import android.os.Handler;
import android.os.Message;
import android.os.Vibrator;
import android.text.TextUtils;
import android.view.KeyEvent;
import android.view.SurfaceHolder;
import android.view.SurfaceHolder.Callback;
import android.view.SurfaceView;
import android.view.View;
import android.view.View.OnClickListener;
import android.widget.Toast;

import com.cardinfolink.cashiersdk.listener.CashierListener;
import com.cardinfolink.cashiersdk.model.OrderData;
import com.cardinfolink.cashiersdk.model.ResultData;
import com.cardinfolink.cashiersdk.sdk.CashierSdk;
import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.carmera.CameraManager;
import com.cardinfolink.yunshouyin.constant.Msg;
import com.cardinfolink.yunshouyin.data.Coupon;
import com.cardinfolink.yunshouyin.data.SaveData;
import com.cardinfolink.yunshouyin.data.TradeBill;
import com.cardinfolink.yunshouyin.decoding.CaptureActivityHandler;
import com.cardinfolink.yunshouyin.decoding.InactivityTimer;
import com.cardinfolink.yunshouyin.ui.SettingActionBarItem;
import com.cardinfolink.yunshouyin.util.Log;
import com.cardinfolink.yunshouyin.util.Utility;
import com.cardinfolink.yunshouyin.view.HintDialog;
import com.cardinfolink.yunshouyin.view.TradingLoadDialog;
import com.cardinfolink.yunshouyin.view.ViewfinderView;
import com.google.zxing.BarcodeFormat;
import com.google.zxing.Result;

import java.io.IOException;
import java.math.BigDecimal;
import java.text.SimpleDateFormat;
import java.util.Date;
import java.util.Vector;

public class CaptureActivity extends BaseActivity implements Callback {
    private static final String TAG = "CaptureActivity";

    private static final float BEEP_VOLUME = 0.10f;
    private static final long VIBRATE_DURATION = 200L;
    /**
     * When the beep has finished playing, rewind to queue up another one.
     */
    private final OnCompletionListener beepListener = new OnCompletionListener() {
        public void onCompletion(MediaPlayer mediaPlayer) {
            mediaPlayer.seekTo(0);
        }
    };

    private CaptureActivityHandler captureActivityHandler;
    private ViewfinderView viewfinderView;
    private boolean hasSurface;
    private Vector<BarcodeFormat> decodeFormats;
    private String characterSet;
    private InactivityTimer inactivityTimer;
    private MediaPlayer mediaPlayer;
    private boolean playBeep;
    private boolean vibrate;

    private Handler mHandler;

    private TradingLoadDialog mTradingLoadDialog;//交易的load的对话框
    private TradingLoadDialog mCouponLoadDialog;//卡券时候用的对话框
    private HintDialog mHintDialog;//显示一些提示信息 下面两个按钮的 对话框
    private HintDialog mHintErrorDialog;//这个也只有卡券时候会用的对话框

    private String total;//实际要支付的金额，如果有优惠这个就是优惠后的金额
    private String originaltotal;//原始金额

    private String mOrderNum;
    private ResultData mResultData;

    private SettingActionBarItem mActionBar;
    private String mCurrentTime;

    private boolean isPolling = false;
    private int pollingCount = 0;

    //从哪里启动的这个activity
    private String originalFromFlag;

    private String chcd;//渠道

    @Override
    public void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);

        setContentView(R.layout.activity_capture);
        initCamera();
        initHandler();
        initLayout();
        initListener();


    }

    private void initLayout() {
        viewfinderView = (ViewfinderView) findViewById(R.id.viewfinder_view);

        //初始化对话框
        mTradingLoadDialog = new TradingLoadDialog(mContext, mHandler, findViewById(R.id.trading_load_dialog), mOrderNum);
        mHintDialog = new HintDialog(mContext, findViewById(R.id.hint_dialog));
        mCouponLoadDialog = new TradingLoadDialog(mContext, mHandler, findViewById(R.id.coupon_load_dialog), mOrderNum);
        mHintErrorDialog = new HintDialog(mContext, findViewById(R.id.hint_error_dialog));
    }

    private void initListener() {
        mActionBar = (SettingActionBarItem) findViewById(R.id.sabi_back);
        mActionBar.setLeftTextOnclickListner(new OnClickListener() {
            @Override
            public void onClick(View v) {

                if (Coupon.getInstance().getVoucherType() != null) {
                    if (Coupon.getInstance().getVoucherType().startsWith("4") || Coupon.getInstance().getVoucherType().startsWith("5")) {
                        showPayFailPref();
                    } else {
                        cleanAfterPay();
                    }
                } else {
                    Coupon.getInstance().clear();//清空卡券信息
                    finish();
                }
            }
        });
        //这里需要设置一下颜色
        mActionBar.setBackgroundColor(Color.BLACK);
        mActionBar.setLeftTextColor(Color.WHITE);
        mActionBar.setTitleColor(Color.WHITE);

        findViewById(R.id.flashlight).setOnClickListener(new OnClickListener() {

            @Override
            public void onClick(View v) {
                if (CameraManager.get().isFlashlight()) {
                    CameraManager.get().closeFlashlight();
                } else {
                    CameraManager.get().openFlashlight();
                }
            }
        });


        findViewById(R.id.reversal).setOnClickListener(new OnClickListener() {

            @Override
            public void onClick(View v) {
                CameraManager.isCameraFront = !CameraManager.isCameraFront;
                SaveData.setCameraFront(getApplication(), CameraManager.isCameraFront);
                if (captureActivityHandler != null) {
                    captureActivityHandler.quitSynchronously();
                    captureActivityHandler = null;
                }
                CameraManager.get().closeDriver();
                openCamera();


            }
        });


    }

    private void initCamera() {
        CameraManager.isCameraFront = SaveData.isCameraFront(this);
        CameraManager.init(getApplication());
        hasSurface = false;
        inactivityTimer = new InactivityTimer(this);
    }

    public void cleanAfterPay() {
        Coupon.getInstance().clear();//清空卡券信息
        if (ScanCodeActivity.getScanCodehandler() != null) {
            ScanCodeActivity.getScanCodehandler().sendEmptyMessage(Msg.MSG_FINISH_BIG_SCANCODEVIEW);
            finish();
        }
    }

    private void openCamera() {
        SurfaceView surfaceView = (SurfaceView) findViewById(R.id.preview_view);
        SurfaceHolder surfaceHolder = surfaceView.getHolder();
        if (hasSurface) {
            initCamera(surfaceHolder);
        } else {
            surfaceHolder.addCallback(this);
            surfaceHolder.setType(SurfaceHolder.SURFACE_TYPE_PUSH_BUFFERS);
        }
        decodeFormats = null;
        characterSet = null;

        playBeep = true;
        AudioManager audioService = (AudioManager) getSystemService(AUDIO_SERVICE);
        if (audioService.getRingerMode() != AudioManager.RINGER_MODE_NORMAL) {
            playBeep = false;
        }
        initBeepSound();
        vibrate = true;
    }

    public void initHandler() {
        mHandler = new Handler() {
            @Override
            public void handleMessage(Message msg) {
                switch (msg.what) {
                    case Msg.MSG_FROM_SCANCODE_SUCCESS: {
                        if ("scancodeview".equals(originalFromFlag)) {
                            //这边是扫码支付
                            mTradingLoadDialog.loading();
                            final OrderData orderData = new OrderData();
                            orderData.orderNum = mOrderNum;
                            orderData.txamt = total;
                            orderData.currency = CashierSdk.SDK_CURRENCY;
                            orderData.scanCodeId = (String) msg.obj;
                            //有优惠金额的时候
                            if (Coupon.getInstance().getSaleDiscount() != null && !"0".equals(Coupon.getInstance().getSaleDiscount())) {
                                orderData.discountMoney = String.valueOf(new BigDecimal(originaltotal).subtract(new BigDecimal(total)).doubleValue());
                                orderData.couponOrderNum = Coupon.getInstance().getOrderNum();
                            }
                            // /orderData.scanCodeId="13241252555";
                            CashierSdk.startPay(orderData, new CashierListener() {

                                @Override
                                public void onResult(ResultData resultData) {

                                    mResultData = resultData;
                                    if (mResultData.respcd.equals("00")) {
                                        mHandler.sendEmptyMessage(Msg.MSG_FROM_SERVER_TRADE_SUCCESS);
                                    } else if (mResultData.respcd.equals("09")) {
                                        mHandler.sendEmptyMessage(Msg.MSG_FROM_SERVER_TRADE_NOPAY);
                                    } else {
                                        //返回14 表示 条码错误或过期
                                        mHandler.sendEmptyMessage(Msg.MSG_FROM_SERVER_TRADE_FAIL);
                                    }
                                }

                                @Override
                                public void onError(int errorCode) {
                                    mHandler.sendEmptyMessage(Msg.MSG_FROM_SERVER_TIMEOUT);
                                }

                            });
                        } else if ("ticketview".equals(originalFromFlag)) {
                            //这里是卡券核销
                            String scancode = (String) msg.obj;
                            Intent intentForTicketView = new Intent();
                            intentForTicketView.putExtra("ticketcode", scancode);
                            setResult(0, intentForTicketView);

                            final OrderData orderData = new OrderData();
                            orderData.orderNum = Utility.geneOrderNumber();
                            orderData.scanCodeId = scancode;
                            mCouponLoadDialog.waiting();
                            CashierSdk.startVeri(orderData, new CashierListener() {
                                @Override
                                public void onResult(ResultData resultData) {
                                    mResultData = resultData;
                                    Coupon.getInstance().setAvailCount(resultData.availCount);
                                    Coupon.getInstance().setCardId(resultData.cardId);
                                    Coupon.getInstance().setVoucherType(resultData.voucherType);
                                    Coupon.getInstance().setSaleDiscount(resultData.saleDiscount);
                                    Coupon.getInstance().setMaxDiscountAmt(resultData.maxDiscountAmt);
                                    Coupon.getInstance().setExpDate(resultData.expDate);
                                    Coupon.getInstance().setSaleMinAmount(resultData.saleMinAmount);//保存卡券核销返回信息
                                    Coupon.getInstance().setOrderNum(resultData.orderNum);
                                    Coupon.getInstance().setScanCodeId(resultData.scanCodeId);
                                    if ("00".equals(mResultData.respcd)) {
                                        Intent intent = new Intent(mContext, CouponResultActivity.class);
                                        Bundle bundle = new Bundle();
                                        bundle.putBoolean("check_coupon_result_flag", true);
                                        intent.putExtras(bundle);
                                        mContext.startActivity(intent);
                                        runOnUiThread(new Runnable() {
                                            @Override
                                            public void run() {
                                                mCouponLoadDialog.hide();
                                            }
                                        });

                                        finish();
                                    } else {
                                        //核销失败
                                        runOnUiThread(new Runnable() {
                                            @Override
                                            public void run() {
                                                mCouponLoadDialog.hide();
                                                mHintErrorDialog.setText(getResources().getString(R.string.coupon_ver_fail), getResources().getString(R.string.coupon_ver_try_again), getResources().getString(R.string.coupon_ver_close));
                                                mHintErrorDialog.show();
                                                mHintErrorDialog.setCancelOnClickListener(new OnClickListener() {
                                                    @Override
                                                    public void onClick(View v) {
                                                        Coupon.getInstance().clear();
                                                        finish();
                                                    }
                                                });
                                                mHintErrorDialog.setOkOnClickListener(new OnClickListener() {
                                                    @Override
                                                    public void onClick(View v) {
                                                        Intent intent = new Intent(mContext, CaptureActivity.class);
                                                        Bundle bundle = new Bundle();
                                                        bundle.putString("original", "ticketview");
                                                        intent.putExtras(bundle);
                                                        mContext.startActivity(intent);
                                                        mHintErrorDialog.hide();
                                                    }
                                                });
                                            }
                                        });
                                    }//end if()
                                }

                                @Override
                                public void onError(int errorCode) {
                                    Log.e(TAG, " starVeri fail===" + errorCode);
                                    MainActivity.getHandler().sendEmptyMessage(Msg.MSG_FROM_SERVER_COUPON_FAIL);
                                }
                            });
                            break;

                        } else if ("searchbill".equals(originalFromFlag)) {
                            String qrCode = (String) msg.obj;
                            Intent mIntent = new Intent();
                            mIntent.putExtra("qrcode", qrCode);
                            setResult(RESULT_OK, mIntent);
                            finish();
                        }
                        break;
                    }
                    case Msg.MSG_FROM_SERVER_TRADE_SUCCESS: {
                        enterPaySuccessActivity();
                        break;
                    }
                    case Msg.MSG_FROM_SERVER_TRADE_FAIL: {
                        enterPayFailActivity();
                        break;
                    }
                    case Msg.MSG_FROM_SERVER_TRADE_NOPAY: {
                        showNopayDialog();
                        break;
                    }
                    case Msg.MSG_FROM_SEARCHING_POLLING: {
                        String title = String.format(getString(R.string.capture_activity_wait_user_input_password), pollingCount);
                        mHintDialog.setTitle(title);
                        searchBill();
                        break;
                    }
                    case Msg.MSG_FROM_SERVER_TIMEOUT: {
                        showPayTimeoutDialog();
                        break;
                    }
                    case Msg.MSG_FROM_SERVER_CLOSEBILL_SUCCESS: {
                        //关单成功
                        showCancelBillSuccess();
                        break;
                    }
                    case Msg.MSG_FROM_SERVER_CLOSEBILL_DOING: {
                        //关单返回09，
                        break;
                    }
                    case Msg.MSG_FROM_SERVER_CLOSEBILL_FAIL: {
                        //关单失败
                        showPayTimeoutDialog();
                        break;
                    }
                    case Msg.MSG_COUPON_CANCEL:
                        cancelCouponVerial();
                        break;
                }
                super.handleMessage(msg);
            }
        };
    }

    /**
     * 取消核销卡券
     */
    public void cancelCouponVerial() {
        finish();
    }


    // 未付款对话框，上面文本，下面一个按钮的对话框
    public void showNopayDialog() {
        //关闭计时器
        mTradingLoadDialog.hide();

        pollingCount = 0;
        String title = String.format(getString(R.string.capture_activity_wait_user_input_password), pollingCount);
        mHintDialog.setTitle(title);

        //左边的对话框
        mHintDialog.setCancelVisibility(View.GONE);

        //右边的对话框
        mHintDialog.setOkText(mContext.getString(R.string.capture_activity_trade));//取消交易
        mHintDialog.setOkOnClickListener(new OnClickListener() {

            @Override
            public void onClick(View v) {
                cancelBill();//取消交易
            }
        });

        mHintDialog.show();

        startPolling();

    }

    //结束轮询
    public void stopPolling() {
        isPolling = false;
        pollingCount = 0;
    }

    //开启轮询
    public void startPolling() {
        if (isPolling) {
            return;
        }

        pollingCount = 0;

        //开启一个线程轮询服务器5次
        isPolling = true;
        new Thread(new Runnable() {

            @Override
            public void run() {
                while (isPolling) {
                    try {
                        pollingCount++;
                        if (pollingCount >= 10) {
                            stopPolling();
                            cancelBill();
                        }
                        if (isPolling) {
                            mHandler.sendEmptyMessage(Msg.MSG_FROM_SEARCHING_POLLING);
                        }
                        Thread.sleep(5000);
                    } catch (InterruptedException e) {
                        e.printStackTrace();
                    }
                }
            }
        }).start();
    }


    // 取消订单
    private void cancelBill() {
        stopPolling();

        if (TextUtils.isEmpty(mOrderNum)) {
            return;
        }
        OrderData orderData = new OrderData();
        orderData.origOrderNum = mOrderNum;
        orderData.orderNum = Utility.geneOrderNumber();//新生成一个订单号
        CashierSdk.startCanc(orderData, new CashierListener() {

            @Override
            public void onResult(ResultData resultData) {
                if (resultData.respcd.equals("00")) {
                    mHandler.sendEmptyMessage(Msg.MSG_FROM_SERVER_CLOSEBILL_SUCCESS);
                } else if (resultData.respcd.equals("09")) {
                    mHandler.sendEmptyMessage(Msg.MSG_FROM_SERVER_CLOSEBILL_DOING);
                } else {
                    mHandler.sendEmptyMessage(Msg.MSG_FROM_SERVER_CLOSEBILL_FAIL);
                }
            }

            @Override
            public void onError(int errorCode) {

            }

        });
    }

    //取消账单，不去判断取消是否成功,在 onPause()方法里面调用
    private void cancelBillInPause() {
        mTradingLoadDialog.hide();
        mHintDialog.hide();

        stopPolling();

        if (TextUtils.isEmpty(mOrderNum)) {
            return;
        }

        final OrderData orderData = new OrderData();
        orderData.origOrderNum = mOrderNum;
        //这里先查询一下，如果还未支付再取消也不迟呀
        Log.e(TAG, "[onPause] cancel before search");
        CashierSdk.startQy(orderData, new CashierListener() {
            @Override
            public void onResult(ResultData resultData) {
                if (resultData.respcd.equals("09")) {
                    Log.e(TAG, "[onPause] not pay yet, will cancel");
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

    // 查询订单
    public void searchBill() {
        OrderData orderData = new OrderData();
        orderData.origOrderNum = mOrderNum;
        CashierSdk.startQy(orderData, new CashierListener() {

            @Override
            public void onResult(ResultData resultData) {
                mResultData = resultData;
                if (resultData.respcd.equals("00")) {
                    mHandler.sendEmptyMessage(Msg.MSG_FROM_SERVER_TRADE_SUCCESS);
                } else if (resultData.respcd.equals("09")) {
                    //09 状态
                } else {
                    mHandler.sendEmptyMessage(Msg.MSG_FROM_SERVER_TRADE_FAIL);
                }
            }

            @Override
            public void onError(int errorCode) {

            }
        });
    }


    //跳转到交易成功的界面
    public void enterPaySuccessActivity() {
        stopPolling();
        mTradingLoadDialog.hide();
        Intent intent = new Intent(CaptureActivity.this, PayResultActivity.class);
        Bundle bun = new Bundle();

        TradeBill tradeBill = new TradeBill();
        tradeBill.orderNum = mResultData.orderNum;
        tradeBill.chcd = mResultData.chcd;
        tradeBill.tandeDate = mCurrentTime;
        tradeBill.response = "success";
        tradeBill.total = total;//付款金额
        if (Coupon.getInstance().getSaleDiscount() != null &&
                !"0".equals(Coupon.getInstance().getSaleDiscount())) {
            //有优惠卡券支付
            tradeBill.originalTotal = originaltotal;//消费金额
        } else {
            //无优惠卡券支付

        }
        bun.putSerializable("TradeBill", tradeBill);
        intent.putExtra("BillBundle", bun);

        startActivity(intent);

        finish();
    }

    //跳转到交易失败的界面
    public void enterPayFailActivity() {
        stopPolling();
        mTradingLoadDialog.hide();
        Intent intent = new Intent(CaptureActivity.this, PayResultActivity.class);
        Bundle bun = new Bundle();

        TradeBill tradeBill = new TradeBill();
        tradeBill.orderNum = mResultData.orderNum;
        tradeBill.chcd = mResultData.chcd;
        tradeBill.tandeDate = mCurrentTime;
        tradeBill.errorDetail = mResultData.errorDetail;
        tradeBill.response = "fail";
        tradeBill.total = total;//付款金额

        //判断是否有优惠金额
        boolean hasDiscount = Coupon.getInstance().getSaleDiscount() != null && !"0".equals(Coupon.getInstance().getSaleDiscount());
        if (hasDiscount) {
            //有优惠卡券支付
            tradeBill.originalTotal = originaltotal;

        } else {

        }
        bun.putSerializable("TradeBill", tradeBill);
        intent.putExtra("BillBundle", bun);

        startActivity(intent);

        finish();
    }


    //服务器超时对话框，中间文本，下面两个按钮
    public void showPayTimeoutDialog() {
        //这里要包 loading 对话框关闭了。而且要结束loading对话框里面的一个线程。
        mTradingLoadDialog.hide();

        mHintDialog.setTitle(mContext.getString(R.string.capture_activity_trade_fail_timerout));
        //返回
        mHintDialog.setCancelText(mContext.getString(R.string.capture_activity_return));
        mHintDialog.setCancelOnClickListener(new OnClickListener() {
            @Override
            public void onClick(View v) {
                mHintDialog.hide();//关闭对话框
            }
        });

        //去账单
        mHintDialog.setOkText(mContext.getString(R.string.capture_activity_goto_bill));
        mHintDialog.setOkOnClickListener(new OnClickListener() {
            @Override
            public void onClick(View v) {
                //TODO 去账单
                mHintDialog.hide();//关闭对话框
                finish();
                MainActivity.getHandler().sendEmptyMessage(Msg.MSG_GO_BILL_VIEW);
            }
        });

        mHintDialog.show();
    }


    // 取消订单成功
    public void showCancelBillSuccess() {
        //这里要包 loading 对话框关闭了。而且要结束loading对话框里面的一个线程。
        mTradingLoadDialog.hide();

        mHintDialog.setTitle(getString(R.string.capture_activity_this_order_had_cancel));
        mHintDialog.setCancelVisibility(View.GONE);
        mHintDialog.setCancelOnClickListener(new OnClickListener() {
            @Override
            public void onClick(View v) {
                mHintDialog.hide();
            }
        });
        mHintDialog.setOkText(getString(R.string.capture_activity_confirm));
        mHintDialog.setOkOnClickListener(new OnClickListener() {
            @Override
            public void onClick(View v) {
                mHintDialog.hide();
                finish();
            }
        });
        mHintDialog.show();
    }


    @Override
    protected void onResume() {
        super.onResume();
        //把初始化的部分移到这里了
        Intent intent = getIntent();
        Bundle bundle = intent.getExtras();
        originalFromFlag = bundle.getString("original");
        if ("scancodeview".equals(originalFromFlag)) {
            total = bundle.getString("total");
            chcd = bundle.getString("chcd");
            originaltotal = bundle.getString("originaltotal");
            Date now = new Date();
            SimpleDateFormat mspf = new SimpleDateFormat("yyyy-MM-dd HH:mm:ss");
            mCurrentTime = mspf.format(now);

            mOrderNum = Utility.geneOrderNumber();
        } else if ("ticketview".equals(originalFromFlag)) {
            //扫卡券
            mActionBar.setTitle(getResources().getString(R.string.coupon_title_first));
        } else if ("searchbill".equals(originalFromFlag)) {
            mActionBar.setTitle(getString(R.string.bill_search));
        }

        SurfaceView surfaceView = (SurfaceView) findViewById(R.id.preview_view);
        SurfaceHolder surfaceHolder = surfaceView.getHolder();
        if (hasSurface) {
            initCamera(surfaceHolder);
        } else {
            surfaceHolder.addCallback(this);
            surfaceHolder.setType(SurfaceHolder.SURFACE_TYPE_PUSH_BUFFERS);
        }
        decodeFormats = null;
        characterSet = null;

        playBeep = true;
        AudioManager audioService = (AudioManager) getSystemService(AUDIO_SERVICE);
        if (audioService.getRingerMode() != AudioManager.RINGER_MODE_NORMAL) {
            playBeep = false;
        }
        initBeepSound();
        vibrate = true;
    }

    @Override
    protected void onPause() {
        cancelBillInPause();
        super.onPause();
        if (captureActivityHandler != null) {
            captureActivityHandler.quitSynchronously();
            captureActivityHandler = null;
        }
        CameraManager.get().closeDriver();
    }

    @Override
    protected void onDestroy() {
        inactivityTimer.shutdown();
        super.onDestroy();
    }

    private void initCamera(SurfaceHolder surfaceHolder) {
        try {
            CameraManager.get().openDriver(surfaceHolder);
        } catch (IOException ioe) {
            return;
        } catch (RuntimeException e) {
            return;
        }
        if (captureActivityHandler == null) {
            captureActivityHandler = new CaptureActivityHandler(this, decodeFormats, characterSet);
        }
    }

    @Override
    public void surfaceChanged(SurfaceHolder holder, int format, int width, int height) {

    }

    @Override
    public void surfaceCreated(SurfaceHolder holder) {
        if (!hasSurface) {
            hasSurface = true;
            initCamera(holder);
        }
    }

    @Override
    public void surfaceDestroyed(SurfaceHolder holder) {
        hasSurface = false;
    }

    public ViewfinderView getViewfinderView() {
        return viewfinderView;
    }

    public Handler getHandler() {
        return captureActivityHandler;
    }

    public void drawViewfinder() {
        viewfinderView.drawViewfinder();

    }

    public void handleDecode(final Result obj, Bitmap barcode) {
        inactivityTimer.onActivity();
        playBeepSoundAndVibrate();
        //从这里发送了扫二维码成功，之后就要调用sdk里面的付款了
        Message msg = mHandler.obtainMessage(Msg.MSG_FROM_SCANCODE_SUCCESS);

        msg.obj = obj.getText().toString();
        mHandler.sendMessageDelayed(msg, 0);

        CameraManager.get().stopPreview();//停止camera的preview
        CameraManager.get().closeDriver();
    }

    private void initBeepSound() {
        if (playBeep && mediaPlayer == null) {
            // The volume on STREAM_SYSTEM is not adjustable, and users found it
            // too loud,
            // so we now play on the music stream.
            setVolumeControlStream(AudioManager.STREAM_MUSIC);
            mediaPlayer = new MediaPlayer();
            mediaPlayer.setAudioStreamType(AudioManager.STREAM_MUSIC);
            mediaPlayer.setOnCompletionListener(beepListener);

            AssetFileDescriptor file = getResources().openRawResourceFd(R.raw.beep);
            try {
                mediaPlayer.setDataSource(file.getFileDescriptor(), file.getStartOffset(), file.getLength());
                file.close();
                mediaPlayer.setVolume(BEEP_VOLUME, BEEP_VOLUME);
                mediaPlayer.prepare();
            } catch (IOException e) {
                mediaPlayer = null;
            }
        }
    }

    private void playBeepSoundAndVibrate() {
        if (playBeep && mediaPlayer != null) {
            mediaPlayer.start();
        }
        if (vibrate) {
            Vibrator vibrator = (Vibrator) getSystemService(VIBRATOR_SERVICE);
            vibrator.vibrate(VIBRATE_DURATION);
        }
    }

    public void showPayFailPref() {
        mHintDialog.setText(getString(R.string.coupon_abandom_verial_or_not), getString(R.string.coupon_pay_again), getString(R.string.coupon_abandom));
        mHintDialog.setCancelOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                //卡券冲正
                OrderData orderData = new OrderData();
                orderData.orderNum = Utility.geneOrderNumber();//订单号
                orderData.origOrderNum = Coupon.getInstance().getOrderNum();//设置原始订单号
                CashierSdk.startReversal(orderData, new CashierListener() {
                    @Override
                    public void onResult(ResultData resultData) {
                        if ("00".equals(resultData.respcd)) {
                            //冲正成功
                            Coupon.getInstance().clear();
                            runOnUiThread(new Runnable() {
                                @Override
                                public void run() {
                                    Toast.makeText(CaptureActivity.this, getString(R.string.coupon_verial_success), Toast.LENGTH_SHORT).show();
                                    finish();
                                }
                            });
                        } else {
                            //冲正失败
                            runOnUiThread(new Runnable() {
                                @Override
                                public void run() {
                                    Toast.makeText(CaptureActivity.this, getString(R.string.coupon_verial_fail), Toast.LENGTH_SHORT).show();
                                }
                            });
                        }
                    }

                    @Override
                    public void onError(int errorCode) {

                    }
                });
                ScanCodeActivity.getScanCodehandler().sendEmptyMessage(Msg.MSG_FINISH_BIG_SCANCODEVIEW);
                mHintDialog.hide();
            }
        });
        mHintDialog.setOkOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                //重新支付
                MainActivity.getHandler().sendEmptyMessage(Msg.MSG_FROM_SERVER_COUPON_SUCCESS);
                mHintDialog.hide();
                finish();
            }

        });
        mHintDialog.show();
    }

    @Override
    public boolean onKeyDown(int keyCode, KeyEvent event) {
        if (keyCode == KeyEvent.KEYCODE_BACK) {
            if (keyCode == KeyEvent.KEYCODE_BACK) {
                if (Coupon.getInstance().getVoucherType() != null) {
                    if (Coupon.getInstance().getVoucherType().startsWith("4") || Coupon.getInstance().getVoucherType().startsWith("5")) {
                        showPayFailPref();
                    } else {
                        cleanAfterPay();
                    }
                } else {
                    Coupon.getInstance().clear();//清空卡券信息
                    finish();
                }
            }
        }

        return false;
    }
}