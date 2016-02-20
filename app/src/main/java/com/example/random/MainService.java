package com.example.random;

import java.io.IOException;
import java.text.SimpleDateFormat;
import java.util.ArrayList;
import java.util.Date;
import java.util.Random;

import android.annotation.SuppressLint;
import android.app.ActivityManager;
import android.app.ActivityManager.MemoryInfo;
import android.app.Notification;
import android.app.NotificationManager;
import android.app.PendingIntent;
import android.app.Service;
import android.content.Context;
import android.content.Intent;
import android.content.SharedPreferences;
import android.content.SharedPreferences.Editor;
import android.content.pm.PackageManager;
import android.content.res.Configuration;
import android.graphics.PixelFormat;
import android.os.Bundle;
import android.os.Handler;
import android.os.IBinder;
import android.os.Message;
import android.os.PowerManager;
import android.os.PowerManager.WakeLock;
import android.text.format.Formatter;
import android.util.Log;
import android.view.LayoutInflater;
import android.view.MotionEvent;
import android.view.View;
import android.view.View.MeasureSpec;
import android.view.View.OnClickListener;
import android.view.View.OnTouchListener;
import android.view.WindowManager;
import android.view.WindowManager.LayoutParams;
import android.widget.Button;
import android.widget.LinearLayout;
import android.widget.TextView;
import android.widget.Toast;

import com.example.model.AppInfo;
import com.example.model.KillLaunchInfo;

@SuppressLint({"WorldReadableFiles", "WorldWriteableFiles", "HandlerLeak"})
public class MainService extends Service {
    private WakeLock mWakeLock;
    private ActivityManager mActivityManager = null;
    /**
     * For showing and hiding our notification.
     */
    private NotificationManager mNotificationManager;
    // 更新悬浮窗的handler
    private Handler mHandler = new Handler() {
        public void handleMessage(Message msg) {
            updateView(msg);
            super.handleMessage(msg);
        }
    };

    private WindowManager windowManager;
    private LayoutParams params;
    private LinearLayout floatLayout;

    private Button stopButton;
    private TextView packagenameTextView;
    private TextView launchtimesTextView;
    private TextView meminfoTextView;

    private LaunchAppThread launchAppThread;

    private ArrayList<AppInfo> firstcheckedAppInfoList;
    private ArrayList<AppInfo> secondcheckedAppInfoList;
    private ArrayList<AppInfo> thirdcheckedAppInfoList;

    // 共享的偏好
    private SharedPreferences mPreferences;

    private int firstRatio;

    private int secondRatio;

    private int thirdRatio;

    private int total;

    private float mTouchX = -1;
    private float mTouchY = -1;

    private float mOrgTouchX = -1;
    private float mOrgTouchY = -1;

    private static final String START_ENABLE = "START_ENABLE";

    // logcat的标签
    private static final String TAG = "random";

    @Override
    public IBinder onBind(Intent intent) {
        return null;
    }

    @Override
    public void onCreate() {
        Log.d(TAG, "mainservice,onCreate");
        mActivityManager = (ActivityManager) getSystemService(Context.ACTIVITY_SERVICE);
        mNotificationManager = (NotificationManager) getSystemService(NOTIFICATION_SERVICE);
        // 创建悬浮窗
        createFloatView();
        showNotification();
        super.onCreate();
    }

    @SuppressWarnings("unchecked")
    public int onStartCommand(Intent intent, int flags, int startId) {
        Log.d(TAG, "mainservice,onStartCommand");
        acquireWakeLock();

        if (intent == null) {
            return START_NOT_STICKY;
        }

        firstcheckedAppInfoList = (ArrayList<AppInfo>) intent.getSerializableExtra("firstcheckedAppInfoList");
        secondcheckedAppInfoList = (ArrayList<AppInfo>) intent.getSerializableExtra("secondcheckedAppInfoList");
        thirdcheckedAppInfoList = (ArrayList<AppInfo>) intent.getSerializableExtra("thirdcheckedAppInfoList");

        firstRatio = intent.getIntExtra("firstRatio", 5);
        secondRatio = intent.getIntExtra("secondRatio", 3);
        thirdRatio = intent.getIntExtra("thirdRatio", 2);
        total = intent.getIntExtra("total", Integer.MAX_VALUE);

        try {
            Runtime.getRuntime().exec("logcat -b events -c");
        } catch (IOException e) {
            e.printStackTrace();
        }

        launchAppThread = new LaunchAppThread();
        launchAppThread.start();

        return START_NOT_STICKY;
    }

    @Override
    public void onDestroy() {
        Log.d(TAG, "mainservice,onDestroy");
        // 调用保存偏好设置的函数
        saveToPreference();

        // 释放wakelock的锁
        releaseWakeLock();

        // 移除悬浮窗
        windowManager.removeView(floatLayout);

        // Cancel the persistent notification.
        mNotificationManager.cancel(R.string.service_started);

        // Tell the user we stopped.
        Toast.makeText(this, R.string.service_stopped, Toast.LENGTH_SHORT).show();
        super.onDestroy();
    }

    @SuppressLint("SimpleDateFormat")
    class LaunchAppThread extends Thread {

        public void run() {
            int firstTotal = 0;
            int secondTotal = 0;
            int thirdTotal = 0;
            if (firstRatio + secondRatio + thirdRatio != 0) {
                firstTotal = (firstRatio) * total / (firstRatio + secondRatio + thirdRatio);
                secondTotal = (secondRatio) * total / (firstRatio + secondRatio + thirdRatio);
                thirdTotal = total - firstTotal - secondTotal;
            }

            for (int i = 0; getRunningFlag() && i < firstTotal; i++) {
                runLaunch(firstcheckedAppInfoList);
            }// end for()
            for (int i = 0; getRunningFlag() && i < secondTotal; i++) {
                runLaunch(secondcheckedAppInfoList);
            }// end for()
            for (int i = 0; getRunningFlag() && i < thirdTotal; i++) {
                runLaunch(thirdcheckedAppInfoList);
            }// end for()

            Intent intent = new Intent();
            intent.setClass(MainService.this, ResultActivity.class);
            intent.addFlags(Intent.FLAG_ACTIVITY_NEW_TASK);

            intent.putExtra("firstcheckedAppInfoList", firstcheckedAppInfoList);
            intent.putExtra("secondcheckedAppInfoList", secondcheckedAppInfoList);
            intent.putExtra("thirdcheckedAppInfoList", thirdcheckedAppInfoList);

            // 启动显示结果的activity
            startActivity(intent);

            // 停止service
            stopSelf();

        }// edn run()

        private void runLaunch(ArrayList<AppInfo> appList) {
            Random random = new Random();
            int index = 0;

            if (appList.size() > 0) {
                index = Math.abs(random.nextInt() % appList.size());
            } else {
                return;
            }

            String packageName = appList.get(index).getPackageName();

            PackageManager packageManager = getPackageManager();
            Intent it = packageManager.getLaunchIntentForPackage(packageName);
            if (it != null) {
                int launchtimes = appList.get(index).getLaunchTimes();
                it.addFlags(Intent.FLAG_ACTIVITY_NEW_TASK);
                startActivity(it);
                appList.get(index).setLaunchTimes(launchtimes + 1);
                SimpleDateFormat formatter = new SimpleDateFormat("MM-dd HH:mm:ss.SSS");
                String currentDate = formatter.format(new Date(System.currentTimeMillis()));
                appList.get(index).setLaunchTimes(launchtimes + 1);
                appList.get(index).getResultArrayList().add(new KillLaunchInfo(KillLaunchInfo.launchFlag, currentDate));

                try {
                    String packagenamestr = "package name: " + packageName;
                    String launchtimesstr = "launch times: " + appList.get(index).getLaunchTimes();
                    String meminfostr = "meminfo:" + getSystemAvaialbeMemorySize();

                    Message msg = new Message();

                    Bundle bundle = new Bundle();
                    bundle.putString("packagename", packagenamestr);
                    bundle.putString("launchtimes", launchtimesstr);
                    bundle.putString("meminfo", meminfostr);

                    msg.setData(bundle);
                    mHandler.sendMessage(msg);

                    Thread.sleep(1000);
                } catch (InterruptedException e) {
                    e.printStackTrace();
                }
            }
        }// end public void runLaunch()

        public void setRunningFlag(boolean run) {
            this.running = run;
        }

        public boolean getRunningFlag() {
            return this.running;
        }

        private boolean running = true;

    }// end class LaunchFirstAppThread

    private String getSystemAvaialbeMemorySize() {
        MemoryInfo memoryInfo = new MemoryInfo();
        mActivityManager.getMemoryInfo(memoryInfo);
        long memSize = memoryInfo.availMem;

        String availMemStr = Formatter.formatFileSize(MainService.this, memSize);
        return availMemStr;
    }

    private void showNotification() {
        // In this sample, we'll use the same text for the ticker and the
        // expanded notification
        CharSequence text = getText(R.string.service_started);

        // Set the icon, scrolling text and timestamp
        Notification.Builder builder = new Notification.Builder(this);
        builder.setSmallIcon(R.drawable.ic_launcher);
        builder.setTicker("");
        builder.setContentTitle("title");
        builder.setContentText("nei rong");
        Notification notification = builder.build();
        notification.flags |= Notification.FLAG_NO_CLEAR;

        // The PendingIntent to launch our activity if the user selects this
        // notification
        PendingIntent contentIntent = PendingIntent.getActivity(this, 0, new Intent(this, MainActivity.class), 0);

        // Set the info for the views that show in the notification panel.
        //notification.setLatestEventInfo(this, getText(R.string.service_label), text, contentIntent);

        // Send the notification.
        // We use a string id because it is a unique number. We use it later to
        // cancel.
        mNotificationManager.notify(R.string.service_started, notification);

    }

    private void createFloatView() {
        windowManager = (WindowManager) getApplicationContext().getSystemService(WINDOW_SERVICE);
        params = new WindowManager.LayoutParams();
        params.type = WindowManager.LayoutParams.TYPE_PHONE;
        params.width = windowManager.getDefaultDisplay().getWidth();
        params.height = WindowManager.LayoutParams.WRAP_CONTENT;
        params.flags = LayoutParams.FLAG_NOT_FOCUSABLE;
        params.format = PixelFormat.TRANSPARENT;

        LayoutInflater inflater = LayoutInflater.from(getApplication());
        floatLayout = (LinearLayout) inflater.inflate(R.layout.activity_floatview, null);

        stopButton = (Button) floatLayout.findViewById(R.id.activity_floatview_stopbutton);

        packagenameTextView = (TextView) floatLayout.findViewById(R.id.activity_floatview_packagename_textview);
        launchtimesTextView = (TextView) floatLayout.findViewById(R.id.activity_floatview_launchtimes_textview);
        meminfoTextView = (TextView) floatLayout.findViewById(R.id.activity_floatview_meminfo_textview);

        floatLayout.measure(MeasureSpec.EXACTLY + floatLayout.getWidth(), MeasureSpec.EXACTLY);

        floatLayout.setOnTouchListener(new OnTouchListener() {

            @Override
            public boolean onTouch(View v, MotionEvent event) {
                switch (event.getAction()) {
                    case MotionEvent.ACTION_DOWN:
                        mOrgTouchX = event.getRawX();
                        mOrgTouchY = event.getRawY();

                        break;
                    case MotionEvent.ACTION_MOVE:

                        mTouchX = event.getRawX();
                        mTouchY = event.getRawY();
                        updateViewPosition();

                        break;
                    case MotionEvent.ACTION_UP:
                        mTouchX = event.getRawX();
                        mTouchY = event.getRawY();
                        updateViewPosition();

                        break;
                }
                return false;
            }
        });

        stopButton.setOnClickListener(new OnClickListener() {

            @Override
            public void onClick(View v) {
                // 停止3个线程的循环
                if (launchAppThread != null) {
                    launchAppThread.setRunningFlag(false);
                }

            }
        });

        windowManager.addView(floatLayout, params);
    }

    @Override
    public void onConfigurationChanged(Configuration newConfig) {
        super.onConfigurationChanged(newConfig);
        // 切换为竖屏
        if (newConfig.orientation == Configuration.ORIENTATION_PORTRAIT) {
            updateViewPosition();
        } else if (newConfig.orientation == Configuration.ORIENTATION_LANDSCAPE) {
            updateViewPosition();
        }
    }

    private void updateViewPosition() {
        params.y = params.y + (int) (mTouchY - mOrgTouchY);
        params.x = params.x + (int) (-mTouchX + mOrgTouchX);
        params.width = windowManager.getDefaultDisplay().getWidth();

        mOrgTouchX = mTouchX;
        mOrgTouchY = mTouchY;

        if (windowManager != null) {
            windowManager.updateViewLayout(floatLayout, params);
        }
    }

    private void updateView(Message msg) {
        packagenameTextView.setText("");
        Bundle bundle = msg.getData();
        if (bundle != null) {
            packagenameTextView.setText("");
            launchtimesTextView.setText("");
            meminfoTextView.setText("");

            String packagename = bundle.getString("packagename");
            String launchtimes = bundle.getString("launchtimes");
            String meminfo = bundle.getString("meminfo");

            packagenameTextView.setText(packagename);
            launchtimesTextView.setText(launchtimes);
            meminfoTextView.setText(meminfo);
        }
    }

    private void acquireWakeLock() {
        if (mWakeLock == null) {
            PowerManager mPowerManager = ((PowerManager) getSystemService(POWER_SERVICE));
            mWakeLock = mPowerManager.newWakeLock(PowerManager.SCREEN_BRIGHT_WAKE_LOCK, this.getClass().getCanonicalName());
            mWakeLock.acquire();
        }
    }

    private void releaseWakeLock() {
        if (mWakeLock != null && mWakeLock.isHeld()) {
            mWakeLock.release();
            mWakeLock = null;
        }
    }

    private void saveToPreference() {
        mPreferences = getSharedPreferences("MainActivity", Context.MODE_WORLD_READABLE | Context.MODE_WORLD_WRITEABLE);
        Editor preferenceEditor = mPreferences.edit();

        // 保存 start的默认值
        preferenceEditor.putBoolean(START_ENABLE, false);
        preferenceEditor.commit();
    }

}
