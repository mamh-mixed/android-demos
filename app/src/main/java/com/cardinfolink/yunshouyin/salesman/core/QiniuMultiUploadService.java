package com.cardinfolink.yunshouyin.salesman.core;

import android.os.Looper;
import android.util.Log;

import com.cardinfolink.yunshouyin.salesman.api.QuickPayException;
import com.cardinfolink.yunshouyin.salesman.model.MerchantPhoto;
import com.qiniu.android.http.ResponseInfo;
import com.qiniu.android.storage.UpCompletionHandler;
import com.qiniu.android.storage.UploadManager;

import org.json.JSONObject;

import java.util.List;
import java.util.UUID;

public class QiniuMultiUploadService {
    private static final String TAG = "QiniuMultiUploadService";

    private QuickPayService quickPayService;
    private UploadManager uploadManager = new UploadManager();
    private List<MerchantPhoto> imageList;
    private QiniuCallbackListener listener;
    private String qiniuKeyPattern;

    public QiniuMultiUploadService(QuickPayService quickPayService) {
        this.quickPayService = quickPayService;
    }

    /**
     * 同时只能有一个upload使用, 加上同步关键字
     *
     * @param imageList
     * @param qiniuKeyPattern
     * @param listener
     */
    public synchronized void upload(final List<MerchantPhoto> imageList, String qiniuKeyPattern, final QiniuCallbackListener listener) {
        //UI线程内
        if (imageList == null || imageList.size() == 0) {
            //TODO: if no picture, show something to user

            return;
        }

        this.imageList = imageList;
        this.listener = listener;
        this.qiniuKeyPattern = qiniuKeyPattern;

        quickPayService.getUploadTokenAsync(new QuickPayCallbackListener<String>() {
            @Override
            public void onSuccess(final String data) {
                //进入后台线程
                upload(0, data);
            }

            @Override
            public void onFailure(QuickPayException ex) {
                //进入后台线程
                //TODO: 七牛的exception得包含在QuickPayException
                listener.onFailure(ex);
            }
        });
    }

    /**
     * Running in background thread
     *
     * @param index
     */
    public void upload(final int index, final String uploadToken) {
        if (Looper.myLooper() == Looper.getMainLooper()) {
            Log.d(TAG, "uploadImageToQiniu" + " is in main thread");
        } else {
            Log.d(TAG, "uploadImageToQiniu" + " is in background thread");
        }

        if (index == imageList.size()) {
            // oncomplete在后台线程
            listener.onComplete();
            return;
        }
        try {
            MerchantPhoto merchantPhoto = imageList.get(index);
            String filename = merchantPhoto.getFilename();
            //NOTE: 文件可能没有后缀名
            String fileType = filename.lastIndexOf(".") >= 0 ? filename.substring(filename.lastIndexOf(".") + 1) : "";
            String qiniuKey = String.format(qiniuKeyPattern, UUID.randomUUID().toString().replace("-", ""), fileType);

            Log.d(TAG, filename);
            Log.d(TAG, qiniuKey);

            merchantPhoto.setQiniuKey(qiniuKey);

            uploadManager.put(filename, qiniuKey, uploadToken,
                    new UpCompletionHandler() {
                        /**
                         * 此方法进入了主线程
                         * @param key
                         * @param info
                         * @param response
                         */
                        @Override
                        public void complete(String key, ResponseInfo info, JSONObject response) {
                            Log.i(TAG, key + ",\r\n " + info + ",\r\n " + response);
                            if (!info.isOK()) {
                                listener.onFailure(new QuickPayException());
                            }
                            if (Looper.myLooper() == Looper.getMainLooper()) {
                                Log.d(TAG, "7ncomplete" + " is in main thread");
                            } else {
                                Log.d(TAG, "7ncomplete" + " is in background thread");
                            }
                            new Thread(new Runnable() {
                                /**
                                 * 再次进入后台线程
                                 */
                                @Override
                                public void run() {
                                    upload(index + 1, uploadToken);
                                }
                            }).start();
                        }
                    }, null);
        } catch (Exception ex) {
            // onFailure 在后台线程
            listener.onFailure(ex);
        }
    }
}
