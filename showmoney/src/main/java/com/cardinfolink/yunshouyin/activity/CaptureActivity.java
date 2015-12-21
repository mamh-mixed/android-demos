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
import android.view.SurfaceHolder;
import android.view.SurfaceHolder.Callback;
import android.view.SurfaceView;
import android.view.View;
import android.view.View.OnClickListener;

import com.cardinfolink.cashiersdk.listener.CashierListener;
import com.cardinfolink.cashiersdk.model.OrderData;
import com.cardinfolink.cashiersdk.model.ResultData;
import com.cardinfolink.cashiersdk.sdk.CashierSdk;
import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.carmera.CameraManager;
import com.cardinfolink.yunshouyin.constant.Msg;
import com.cardinfolink.yunshouyin.data.SessonData;
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

    /**
     * Called when the activity is first created.
     */
    @Override
    public void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.scancode_activity);
        Intent intent = getIntent();
        total = intent.getStringExtra("total");
        //这里不需要传人支付类型了，服务器判断。
        Date now = new Date();
        SimpleDateFormat spf = new SimpleDateFormat("yyMMddHHmmss");
        mOrderNum = spf.format(now);
        Random random = new Random();
        for (int i = 0; i < 5; i++) {
            mOrderNum = mOrderNum + random.nextInt(10);
        }

        CameraManager.init(getApplication());

        initHandler();
        initLayout();
        initListener();

        hasSurface = false;
        inactivityTimer = new InactivityTimer(this);
    }

    private void initLayout() {
        viewfinderView = (ViewfinderView) findViewById(R.id.viewfinder_view);
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
        mActionBar.setBackgroundColor(Color.BLACK);
        mActionBar.setLeftTextColor(Color.WHITE);

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
                                    mHandler.sendEmptyMessage(Msg.MSG_FROM_SERVER_TRADE_FAIL);
                                }
                            }

                            @Override
                            public void onError(int errorCode) {
                                mHandler.sendEmptyMessage(Msg.MSG_FROM_SERVER_TIMEOUT);
                            }

                        });

                        break;
                    }

                    case Msg.MSG_FROM_DIGLOG_CLOSE: {
                        CaptureActivity.this.finish();
                        break;
                    }
                    case Msg.MSG_FROM_SERVER_TRADE_SUCCESS: {
                        showPaySuccessDialog();
                        break;
                    }
                    case Msg.MSG_FROM_SERVER_TRADE_FAIL: {
                        showPayFailDialog();
                        break;
                    }
                    case Msg.MSG_FROM_SERVER_TRADE_NOPAY: {
                        showNopayDialog();
                        break;
                    }
                    case Msg.MSG_FROM_SUCCESS_DIGLOG_HISTORY: {
                        SessonData.positionView = 1;
                        setResult(101);
                        finish();
                        break;
                    }
                }
                super.handleMessage(msg);
            }
        };
    }


    /**
     * 未付款的对话框
     * <p/>
     * 这两个完全一样的对话框可以复用一样的layout文件。
     * 显示交易成功的对话框，上边一个图片，中间显示文本，下边两个按钮 对话框
     * 显示本次交易出错的对话框 上边一个图片，中间显示文本，下边两个按钮对话框
     * <p/>
     * 未付款对话框，上面文本，下面一个按钮的对话框
     */
    public void showNopayDialog() {
        //左边的对话框
        mHintDialog.setCancelText(mContext.getResources().getString(R.string.txt_query_result));//查询结果
        mHintDialog.setCancelOnClickListener(new OnClickListener() {

            @Override
            public void onClick(View v) {
                OrderData orderData = new OrderData();
                orderData.origOrderNum = mOrderNum;
                CashierSdk.startQy(orderData, new CashierListener() {

                    @Override
                    public void onResult(ResultData resultData) {

                        if (resultData.respcd.equals("00")) {
                            mHandler.sendEmptyMessage(Msg.MSG_FROM_SERVER_TRADE_SUCCESS);
                        } else if (resultData.respcd.equals("09")) {
                            mHandler.sendEmptyMessage(Msg.MSG_FROM_SERVER_TRADE_NOPAY);
                        } else {
                            mHandler.sendEmptyMessage(Msg.MSG_FROM_SERVER_TRADE_FAIL);
                        }
                    }

                    @Override
                    public void onError(int errorCode) {
                        mHandler.sendEmptyMessage(Msg.MSG_FROM_SERVER_TIMEOUT);
                    }

                });
            }
        });
        //右边的对话框
        mHintDialog.setOkText(mContext.getString(R.string.txt_close));//关闭
        mHintDialog.setOkOnClickListener(new OnClickListener() {

            @Override
            public void onClick(View v) {
                mHintDialog.hide();
            }
        });

        //这里要包 loading 对话框关闭了。而且要结束loading对话框里面的一个线程。
        //isLoading = false;
        //loadDialogView.setVisibility(View.GONE);//先要把loading的对话框隐藏了

        mHintDialog.show();
    }


    /**
     * 显示交易成功的对话框，上边一个图片，中间显示文本，下边两个按钮 对话框
     */
    public void showPaySuccessDialog() {
        //左边按钮
        mHintDialog.setCancelText(mContext.getString(R.string.txt_history_txns));//历史交易
        mHintDialog.setCancelOnClickListener(new OnClickListener() {

            @Override
            public void onClick(View v) {
                mHandler.sendEmptyMessage(Msg.MSG_FROM_DIGLOG_CLOSE);
            }
        });


        //右边按钮
        mHintDialog.setOkText(mContext.getString(R.string.txt_return));//返回
        mHintDialog.setOkOnClickListener(new OnClickListener() {

            @Override
            public void onClick(View v) {
                mHandler.sendEmptyMessage(Msg.MSG_FROM_SUCCESS_DIGLOG_HISTORY);
            }
        });

        //这里要包 loading 对话框关闭了。而且要结束loading对话框里面的一个线程。
        //isLoading = false;
        //loadDialogView.setVisibility(View.GONE);

        mHintDialog.show();
    }

    /**
     * 显示本次交易出错的对话框 上边一个图片，中间显示文本，下边两个按钮对话框
     */
    public void showPayFailDialog() {
        //左边按钮
        mHintDialog.setCancelText(mContext.getString(R.string.txt_query_result));//查询结果
        mHintDialog.setCancelOnClickListener(new OnClickListener() {

            @Override
            public void onClick(View v) {
                mHandler.sendEmptyMessage(Msg.MSG_FROM_DIGLOG_CLOSE);
            }
        });

        //右边按钮
        mHintDialog.setOkText(mContext.getString(R.string.txt_return));//返回
        mHintDialog.setOkOnClickListener(new OnClickListener() {

            @Override
            public void onClick(View v) {
                mHandler.sendEmptyMessage(Msg.MSG_FROM_SUCCESS_DIGLOG_HISTORY);
            }
        });
        //这里要包 loading 对话框关闭了。而且要结束loading对话框里面的一个线程。
        //isLoading = false;
        //loadDialogView.setVisibility(View.GONE);

        mHintDialog.show();
    }

    /**
     * 服务器超时对话框，中间文本，下面两个按钮
     */
    public void showPayTimeoutDialog() {
        //返回
        mHintDialog.setCancelText(mContext.getString(R.string.txt_return));
        mHintDialog.setCancelOnClickListener(new OnClickListener() {
            @Override
            public void onClick(View v) {
                mHandler.sendEmptyMessage(Msg.MSG_FROM_DIGLOG_CLOSE);
            }
        });

        //去账单
        mHintDialog.setOkText(mContext.getString(R.string.txt_goto_bill));
        mHintDialog.setOkOnClickListener(new OnClickListener() {
            @Override
            public void onClick(View v) {
                //TODO 去账单
            }
        });
        //这里要包 loading 对话框关闭了。而且要结束loading对话框里面的一个线程。
        //isLoading = false;
        //loadDialogView.setVisibility(View.GONE);

        mHintDialog.show();
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
            handler = new CaptureActivityHandler(this, decodeFormats,
                    characterSet);
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
        Message msg = mHandler.obtainMessage(Msg.MSG_FROM_SCANCODE_SUCCESS);
        msg.obj = obj.getText().toString();
        mHandler.sendMessageDelayed(msg, 0);

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