package cn.weipass.biz;

import android.app.Activity;
import android.app.AlertDialog;
import android.app.ProgressDialog;
import android.content.Context;
import android.content.DialogInterface;
import android.widget.Toast;

import com.cardinfolink.cashiersdk.model.ResultData;
import com.cardinfolink.yunshouyin.data.TradeBill;
import com.cardinfolink.yunshouyin.util.Log;

import org.json.JSONException;
import org.json.JSONObject;

import cn.weipass.biz.util.ToolsUtil;
import cn.weipass.pos.sdk.AuthorizationManager;
import cn.weipass.pos.sdk.BizServiceInvoker;
import cn.weipass.pos.sdk.IPrint;
import cn.weipass.pos.sdk.LatticePrinter;
import cn.weipass.pos.sdk.MagneticReader;
import cn.weipass.pos.sdk.Photograph;
import cn.weipass.pos.sdk.Printer;
import cn.weipass.pos.sdk.PsamManager;
import cn.weipass.pos.sdk.Scanner;
import cn.weipass.pos.sdk.ServiceManager;
import cn.weipass.pos.sdk.Sonar;
import cn.weipass.pos.sdk.Weipos;
import cn.weipass.pos.sdk.impl.WeiposImpl;

/**
 * Created by feng.chen@cardinfolink.com on 16/1/19.
 */
public class WeipassManager {

    private static final String TAG = "WeipassManager";

    private static WeipassManager instance;

    private Context context;
    private Activity activity;
//    private ProgressDialog pd = null;

    private Scanner sacner = null;
    private Printer printer = null;
    private Sonar sonar = null;
    private ServiceManager mServiceManager = null;
    private MagneticReader mMagneticReader;// 磁条卡管理
    private Photograph mPhotograph;
    private LatticePrinter latticePrinter;// 点阵打印

    private AuthorizationManager mAuthorizationManager;

    private BizServiceInvoker mBizServiceInvoker;

    private PsamManager psamManager;

    private boolean is2s = false;

    private WeipassManager() {

    }

    private WeipassManager(Context context) {
        this.context = context;
        this.activity = (Activity) context;
    }

    public static synchronized WeipassManager getInstance(Context context) {
        if (instance == null) {
            instance = new WeipassManager(context);
            Log.i(TAG, "new WeipassManager : " + instance.getClass().getSimpleName());
            return instance;
        } else {
            return instance;
        }
    }


    /**
     * WeiposImpl的初始化（init函数）和销毁（destroy函数），
     * 最好分别放在一级页面的onCreate和onDestroy中执行。 其他子页面不用再调用，可以直接获取能力对象并使用。
     */
    public void init() {

        WeiposImpl.as().init(context, new Weipos.OnInitListener() {

            @Override
            public void onInitOk() {
                /**
                 * onInitOk()方法中也不能直接操作ui
                 */
                // 获取设备基础信息
                String deviceInfo = WeiposImpl.as().getDeviceInfo();
                Log.i(TAG, "deviceInfo ------------------ " + deviceInfo);
                try {
                    if (deviceInfo != null) {
                        JSONObject deviceJson = new JSONObject(deviceInfo);
                        if (deviceJson.has("deviceType")) {
                            String deviceType = deviceJson.getString("deviceType");
                            if (deviceType.equals("2")) {
                                //旺POS2设备
                                is2s = false;
                            } else if (deviceType.equalsIgnoreCase("2s")) {
                                //旺POS2s设备
                                is2s = true;
                            }
                        } else {
                            //旺POS2设备
                            is2s = false;
                        }
                    }
                } catch (JSONException e1) {
                    // TODO Auto-generated catch block
                    e1.printStackTrace();
                }

                // TODO Auto-generated method stub
                sacner = WeiposImpl.as().openScanner();
                sonar = WeiposImpl.as().openSonar();

                mServiceManager = WeiposImpl.as().openServiceManager();

                try {
                    // 设备可能没有打印机，open会抛异常
                    printer = WeiposImpl.as().openPrinter();
                } catch (Exception e) {
                    // TODO: handle exception
                }
                // 回调函数中不能做UI操作，所以可以使用runOnUiThread函数来包装一下代码块
                activity.runOnUiThread(new Runnable() {
                    public void run() {
                        Toast.makeText(context, "微POS 设备初始化完成", Toast.LENGTH_SHORT).show();
                    }
                });
                try {
                    // 初始化磁条卡sdk对象
                    mMagneticReader = WeiposImpl.as().openMagneticReader();
                } catch (Exception e) {
                    // TODO: handle exception
                }

                try {
                    // 初始化相机
                    mPhotograph = WeiposImpl.as().openPhotograph();
                } catch (Exception e) {
                    // TODO: handle exception
                }

                try {
                    latticePrinter = WeiposImpl.as().openLatticePrinter();
                } catch (Exception e) {
                    // TODO: handle exception
                }

                try {
                    mAuthorizationManager = WeiposImpl.as()
                            .openAuthorizationManager();
                } catch (Exception e) {
                    // TODO: handle exception
                }
                try {
                    psamManager = WeiposImpl.as().openPsamManager();
                } catch (Exception e) {
                    // TODO: handle exception
                }

                try {
                    // 初始化服务调用
                    mBizServiceInvoker = WeiposImpl.as()
                            .openBizServiceInvoker();
                } catch (Exception e) {
                    // TODO: handle exception
                }
            }

            @Override
            public void onError(String message) {
                // TODO Auto-generated method stub
                final String msg = message;
                // 回调函数中不能做UI操作，所以可以使用runOnUiThread函数来包装一下代码块
                activity.runOnUiThread(new Runnable() {
                    public void run() {
                        Toast.makeText(context, msg, Toast.LENGTH_SHORT).show();
                    }
                });
            }
        });
    }


    /**
     * 注意：destroy函数在一级根页面的onDestroy调用，以防止在二级页面或者返回到一级页面中
     * 使用weipos能力对象（例如：Printer）抛出服务未初始化的异常.
     */
    public void destory() {

        WeiposImpl.as().destroy();
    }


    /**
     * 打印凭条
     *
     * @param bill
     */
    public void print(ResultData resultData) {

        if (printer == null) {
            Toast.makeText(context, "尚未初始化打印sdk，请稍后再试", Toast.LENGTH_SHORT).show();
            return;
        }
        Log.i(TAG, "print start :正在打印小票...");
//        if (pd == null) {
//            pd = new ProgressDialog(context);
//        }
//        pd.setMessage("正在打印小票...");
//        pd.show();
        printer.setOnEventListener(new IPrint.OnEventListener() {

            @Override
            public void onEvent(final int what, String in) {
                // TODO Auto-generated method stub
                Log.i(TAG, "print result : " + in);
                final String info = in;
                // 回调函数中不能做UI操作，所以可以使用runOnUiThread函数来包装一下代码块
                activity.runOnUiThread(new Runnable() {
                    public void run() {
//                        if (pd != null) {
//                            pd.hide();
//                        }
                        final String message = ToolsUtil.getPrintErrorInfo(what, info);
                        if (message == null || message.length() < 1) {
                            return;
                        }
                        if (message.equals("EVENT_OK")) {
                            printer.cutting();
                        }
//                        showResultInfo(context, "打印", "打印结果信息", message);
                    }
                });
            }
        });
        ToolsUtil.printNormal(context, printer, resultData);

    }


    /**
     * 弹出对话框
     *
     * @param context
     * @param operInfo
     * @param titleHeader
     * @param info
     */
    private void showResultInfo(Context context, String operInfo, String titleHeader, String info) {
        AlertDialog.Builder builder = new AlertDialog.Builder(context);

        builder.setMessage(titleHeader + ":" + info);
        builder.setTitle(operInfo);
        builder.setPositiveButton("确认",
                new android.content.DialogInterface.OnClickListener() {
                    @Override
                    public void onClick(DialogInterface dialog, int which) {

                        dialog.dismiss();
                    }
                });
        builder.setNegativeButton("取消",
                new android.content.DialogInterface.OnClickListener() {
                    @Override
                    public void onClick(DialogInterface dialog, int which) {
                        dialog.dismiss();
                    }
                });
        builder.create().show();
    }
}
