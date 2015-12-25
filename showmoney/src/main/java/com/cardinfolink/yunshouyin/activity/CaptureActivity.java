package com.cardinfolink.yunshouyin.activity;

import android.content.Intent;
import android.content.res.AssetFileDescriptor;
import android.graphics.Bitmap;
import android.graphics.Color;
import android.graphics.drawable.GradientDrawable;
import android.media.AudioManager;
import android.media.MediaPlayer;
import android.media.MediaPlayer.OnCompletionListener;
import android.os.Bundle;
import android.os.Handler;
import android.os.Message;
import android.os.Vibrator;
import android.text.TextUtils;
import android.util.Log;
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
import com.cardinfolink.yunshouyin.decoding.CaptureActivityHandler;
import com.cardinfolink.yunshouyin.decoding.InactivityTimer;
import com.cardinfolink.yunshouyin.ui.SettingActionBarItem;
import com.cardinfolink.yunshouyin.view.HintDialog;
import com.cardinfolink.yunshouyin.view.TradingLoadDialog;
import com.cardinfolink.yunshouyin.view.ViewfinderView;
import com.google.zxing.BarcodeFormat;
import com.google.zxing.Result;

import java.io.IOException;
import java.text.SimpleDateFormat;
import java.util.Date;
import java.util.Random;
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

    private CaptureActivityHandler handler;
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
    private HintDialog mHintDialog;//显示一些提示信息 下面两个按钮的 对话框

    private String total;

    private String mOrderNum;
    private ResultData mResultData;

    private SettingActionBarItem mActionBar;
    private String mCurrentTime;

    private boolean isPolling = false;
    private int pollingCount = 0;

    //从哪里启动的这个activity
    private String originalFromFlag;

    @Override
    public void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_capture);
        Intent intent = getIntent();
        Bundle bundle = intent.getExtras();

        CameraManager.init(getApplication());

        initHandler();
        initLayout();
        initListener();

        originalFromFlag = bundle.getString("original");
        if ("scancodeview".equals(originalFromFlag)) {
            total = bundle.getString("total");
            //这里不需要传人支付类型了，服务器判断。
            Date now = new Date();

            SimpleDateFormat mspf = new SimpleDateFormat("yyyy/MM/dd HH:mm:ss");
            mCurrentTime = mspf.format(now);

            mOrderNum = geneOrderNumber();

        } else if ("ticketview".equals(originalFromFlag)) {
            //扫卡券
            mActionBar.setTitle(getResources().getString(R.string.coupon_title_first));
        }

        hasSurface = false;
        inactivityTimer = new InactivityTimer(this);
    }

    private void initLayout() {
        viewfinderView = (ViewfinderView) findViewById(R.id.viewfinder_view);

        //初始化对话框
        mTradingLoadDialog = new TradingLoadDialog(mContext, mHandler, findViewById(R.id.trading_load_dialog), mOrderNum);
        mHintDialog = new HintDialog(mContext, findViewById(R.id.hint_dialog));
    }

    private void initListener() {
        mActionBar = (SettingActionBarItem) findViewById(R.id.sabi_back);
        mActionBar.setLeftTextOnclickListner(new OnClickListener() {

            @Override
            public void onClick(View v) {
                finish();
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
                    GradientDrawable myGrad = (GradientDrawable) v.getBackground();
                    myGrad.setColor(Color.parseColor("#222222"));
                    CameraManager.get().closeFlashlight();
                } else {
                    GradientDrawable myGrad = (GradientDrawable) v.getBackground();
                    myGrad.setColor(Color.parseColor("#444444"));
                    CameraManager.get().openFlashlight();
                }
            }
        });
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
                            orderData.currency = "156";
                            orderData.scanCodeId = (String) msg.obj;
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
                            Toast.makeText(mContext, "ticketview", Toast.LENGTH_SHORT).show();
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
                }
                super.handleMessage(msg);
            }
        };
    }


    // 未付款对话框，上面文本，下面一个按钮的对话框
    public void showNopayDialog() {
        //关闭计时器
        mTradingLoadDialog.hide();

        pollingCount = 0;
        String title = String.format(getString(R.string.capture_activity_wait_user_input_password), pollingCount);
        mHintDialog.setTitle(title);

        //左边的对话框
        mHintDialog.setCancelText(mContext.getResources().getString(R.string.capture_activity_query_manual));//手动查询
        mHintDialog.setCancelOnClickListener(new OnClickListener() {
            private int pressCount = 0;
            private long lastClickTime;

            public synchronized boolean isFastClick() {
                long time = System.currentTimeMillis();
                if (time - lastClickTime < 5000) {
                    return true;
                }
                lastClickTime = time;
                return false;
            }

            @Override
            public void onClick(View v) {
                if (isFastClick()) {
                    return;
                }
                mHintDialog.setTitle(String.format("手动查询：%s 次", pressCount));
                pressCount++;
                stopPolling();//结束轮询
                searchBill();//手动查询,手动查询把 轮询关闭然后每次按一下按钮查询一下
            }
        });
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
                            Log.e(TAG, "[Thread] is polling = mHandler.sendEmptyMessage(Msg.MSG_FROM_SEARCHING_POLLING) = " + pollingCount);
                            mHandler.sendEmptyMessage(Msg.MSG_FROM_SEARCHING_POLLING);
                        }
                        Thread.sleep(5000);
                    } catch (InterruptedException e) {
                        e.printStackTrace();
                    }
                }
                Log.e(TAG, "[Thread] end while() = " + pollingCount);
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
        orderData.orderNum = geneOrderNumber();//新生成一个订单号
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
        bun.putString("txamt", mResultData.txamt);
        bun.putString("orderNum", mResultData.orderNum);
        bun.putString("chcd", mResultData.chcd);
        bun.putString("mCurrentTime", mCurrentTime);
        bun.putBoolean("result", true);

        intent.putExtras(bun);

        intent.setClass(CaptureActivity.this, PayResultActivity.class);

        startActivity(intent);

        finish();
    }

    //跳转到交易失败的界面
    public void enterPayFailActivity() {
        stopPolling();
        mTradingLoadDialog.hide();
        Intent intent = new Intent(CaptureActivity.this, PayResultActivity.class);
        Bundle bun = new Bundle();

        bun.putString("txamt", mResultData.txamt);
        bun.putString("orderNum", mResultData.orderNum);
        bun.putString("chcd", mResultData.chcd);
        bun.putString("errorDetail", mResultData.errorDetail);
        bun.putString("mCurrentTime", mCurrentTime);
        bun.putBoolean("result", false);

        intent.putExtras(bun);

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
            }
        });

        mHintDialog.show();
    }


    // 取消订单成功
    public void showCancelBillSuccess() {
        //这里要包 loading 对话框关闭了。而且要结束loading对话框里面的一个线程。
        mTradingLoadDialog.hide();

        mHintDialog.setTitle(getString(R.string.capture_activity_this_order_had_cancel));
        mHintDialog.setCancelText(getString(R.string.capture_activity_i_know));
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

    // 时间加上一个随机数 生成账单号
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

    @Override
    protected void onResume() {
        super.onResume();
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
        super.onPause();
        if (handler != null) {
            handler.quitSynchronously();
            handler = null;
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
        if (handler == null) {
            handler = new CaptureActivityHandler(this, decodeFormats, characterSet);
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
        return handler;
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

}