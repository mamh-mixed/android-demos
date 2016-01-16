package com.cardinfolink.yunshouyin.util;

import android.app.Application;
import android.app.Notification;
import android.content.Context;
import android.os.Handler;
import android.support.v4.app.NotificationCompat;
import android.widget.RemoteViews;
import android.widget.Toast;

import com.cardinfolink.yunshouyin.BuildConfig;
import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.api.QuickPayConfigStorage;
import com.cardinfolink.yunshouyin.constant.SystemConfig;
import com.cardinfolink.yunshouyin.core.BankDataService;
import com.cardinfolink.yunshouyin.core.BankDataServiceImpl;
import com.cardinfolink.yunshouyin.core.QiniuMultiUploadService;
import com.cardinfolink.yunshouyin.core.QuickPayService;
import com.cardinfolink.yunshouyin.core.QuickPayServiceImpl;
import com.umeng.message.PushAgent;
import com.umeng.message.UTrack;
import com.umeng.message.UmengMessageHandler;
import com.umeng.message.UmengNotificationClickHandler;
import com.umeng.message.entity.UMessage;

public class ShowMoneyApp extends Application {
    private static final String TAG = "ShowMoneyApp";

    /**
     * 环境配置，是生产环境，测试环境，开发环境
     */
    public static final String ENVIRONMENT = BuildConfig.ENVIRONMENT;

    /**
     * git log相关的信息在这里保存
     * 这个是不可以改变的，所以没必要放在SystemConfig里面了
     */
    public static final String GIT = BuildConfig.GIT;

    private static ShowMoneyApp instance;

    private QuickPayConfigStorage quickPayConfigStorage;
    private QuickPayService quickPayService;
    private BankDataService bankDataService;
    private QiniuMultiUploadService qiniuMultiUploadService;

    private PushAgent mPushAgent;


    public static ShowMoneyApp getInstance() {
        return instance;
    }

    public static String getResString(int id) {
        return instance.getResources().getString(id);
    }


    public QuickPayService getQuickPayService() {
        return quickPayService;
    }

    public BankDataService getBankDataService() {
        return bankDataService;
    }

    public QiniuMultiUploadService getQiniuMultiUploadService() {
        return qiniuMultiUploadService;
    }

    @Override
    public void onCreate() {
        super.onCreate();
        instance = this;
        initEnvironment();
        initPushAgent();
    }


    private void initEnvironment() {
        quickPayConfigStorage = new QuickPayConfigStorage();
        //dev, test, pro 是一样的
        quickPayConfigStorage.setAppKey(SystemConfig.APP_KEY);

        //default is pro
        SystemConfig.IS_PRODUCE = true;
        SystemConfig.SERVER = SystemConfig.PRO_SERVER;
        SystemConfig.DEBUG = BuildConfig.DEBUG;
        Log.LOG = BuildConfig.DEBUG;
        switch (ENVIRONMENT) {
            case "pro":
                SystemConfig.IS_PRODUCE = true;
                SystemConfig.SERVER = SystemConfig.PRO_SERVER;
                break;
            case "test":
                SystemConfig.IS_PRODUCE = false;
                SystemConfig.SERVER = SystemConfig.TEST_SERVER;
                break;
            case "dev":
                SystemConfig.IS_PRODUCE = false;
                SystemConfig.SERVER = SystemConfig.DEV_SERVER;
                break;
            default:
                break;
        }

        quickPayConfigStorage.setUrl(SystemConfig.SERVER);

        quickPayConfigStorage.setBankbaseKey(SystemConfig.BANKBASE_KEY);
        quickPayConfigStorage.setBankbaseUrl(SystemConfig.BANKBASE_URL);

        quickPayService = new QuickPayServiceImpl(quickPayConfigStorage);
        bankDataService = new BankDataServiceImpl(quickPayConfigStorage);
        qiniuMultiUploadService = new QiniuMultiUploadService(quickPayService);


    }


    private void initPushAgent() {

        mPushAgent = PushAgent.getInstance(this);
        mPushAgent.setDebugMode(false);

        UmengMessageHandler messageHandler = new UmengMessageHandler() {
            /**
             * 参考集成文档的1.6.3
             * http://dev.umeng.com/push/android/integration#1_6_3
             * */
            @Override
            public void dealWithCustomMessage(final Context context, final UMessage msg) {
                new Handler().post(new Runnable() {

                    @Override
                    public void run() {
                        // 对自定义消息的处理方式，点击或者忽略
                        boolean isClickOrDismissed = true;
                        if (isClickOrDismissed) {
                            //自定义消息的点击统计
                            UTrack.getInstance(getApplicationContext()).trackMsgClick(msg);
                        } else {
                            //自定义消息的忽略统计
                            UTrack.getInstance(getApplicationContext()).trackMsgDismissed(msg);
                        }
                        Toast.makeText(context, msg.custom, Toast.LENGTH_LONG).show();
                    }
                });
            }

            /**
             * 参考集成文档的1.6.4
             * http://dev.umeng.com/push/android/integration#1_6_4
             * */
            @Override
            public Notification getNotification(Context context, UMessage msg) {
                switch (msg.builder_id) {
                    case 1:
                        NotificationCompat.Builder builder = new NotificationCompat.Builder(context);
                        RemoteViews myNotificationView = new RemoteViews(context.getPackageName(), R.layout.notification_view);
                        myNotificationView.setTextViewText(R.id.notification_title, msg.title);
                        myNotificationView.setTextViewText(R.id.notification_text, msg.text);
                        myNotificationView.setImageViewBitmap(R.id.notification_large_icon, getLargeIcon(context, msg));
                        myNotificationView.setImageViewResource(R.id.notification_small_icon, getSmallIconId(context, msg));
                        builder.setContent(myNotificationView);
                        builder.setContentTitle(msg.title)
                                .setContentText(msg.text)
                                .setTicker(msg.ticker)
                                .setAutoCancel(true);
                        Notification mNotification = builder.build();
                        //由于Android v4包的bug，在2.3及以下系统，Builder创建出来的Notification，并没有设置RemoteView，故需要添加此代码
                        mNotification.contentView = myNotificationView;
                        return mNotification;
                    default:
                        //默认为0，若填写的builder_id并不存在，也使用默认。
                        return super.getNotification(context, msg);
                }
            }
        };
        mPushAgent.setMessageHandler(messageHandler);

        /**
         * 该Handler是在BroadcastReceiver中被调用，故
         * 如果需启动Activity，需添加Intent.FLAG_ACTIVITY_NEW_TASK
         * 参考集成文档的1.6.2
         * http://dev.umeng.com/push/android/integration#1_6_2
         * */
        UmengNotificationClickHandler notificationClickHandler = new UmengNotificationClickHandler() {
            @Override
            public void dealWithCustomAction(Context context, UMessage msg) {
                Log.e(TAG, " ====public void dealWithCustomAction(Context context, UMessage msg) ");
                Toast.makeText(context, msg.custom, Toast.LENGTH_LONG).show();
            }
        };
        mPushAgent.setNotificationClickHandler(notificationClickHandler);

    }
}
