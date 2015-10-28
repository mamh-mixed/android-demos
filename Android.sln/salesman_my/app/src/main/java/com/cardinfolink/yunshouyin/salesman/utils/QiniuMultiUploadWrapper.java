package com.cardinfolink.yunshouyin.salesman.utils;

import android.content.Context;
import android.os.Looper;
import android.util.Log;

import com.cardinfolink.yunshouyin.salesman.models.SAMerchantPhoto;
import com.cardinfolink.yunshouyin.salesman.models.SessonData;
import com.qiniu.android.http.ResponseInfo;
import com.qiniu.android.storage.UpCompletionHandler;
import com.qiniu.android.storage.UploadManager;

import org.json.JSONObject;

import java.util.List;
import java.util.UUID;

public class QiniuMultiUploadWrapper {
    private UploadManager uploadManager = new UploadManager();

    private List<SAMerchantPhoto> imageList;
    private QiniuTaskListener listener;
    private Context context;
    private String qiniuKeyPattern;
    private String uploadToken;

    public void uploadImageToQiniu(Context context, final List<SAMerchantPhoto> imageList, String qiniuKeyPattern, final QiniuTaskListener listener) {
        if (imageList == null || imageList.size() == 0) {
            return;
        }

        this.imageList = imageList;
        this.context = context;
        this.listener = listener;
        this.qiniuKeyPattern = qiniuKeyPattern;

        new Thread(new Runnable() {
            @Override
            public void run() {
                uploadToken = SessonData.getUploadToken();
                uploadImageToQiniu(0);
            }
        }).start();
    }

    /**
     * Running in background thread
     *
     * @param index
     */
    public void uploadImageToQiniu(final int index) {
        if (Looper.myLooper() == Looper.getMainLooper()) {
            Log.d("jiahua:", "uploadImageToQiniu" + " is in main thread");
        } else {
            Log.d("jiahua:", "uploadImageToQiniu" + " is in background thread");
        }
        if (index == imageList.size()) {
            // oncomplete在后台线程
            listener.onComplete();
            return;
        }
        try {
            SAMerchantPhoto saMerchantPhoto = imageList.get(index);
            String filename = saMerchantPhoto.getFilename();
            //NOTE: 文件可能没有后缀名
            String fileType = filename.lastIndexOf(".") >= 0 ? filename.substring(filename.lastIndexOf(".") + 1) : "";
            String qiniuKey = String.format(qiniuKeyPattern, UUID.randomUUID().toString().replace("-", ""), fileType);

            Log.d("jiahua:", filename);
            Log.d("jiahua:", qiniuKey);

            saMerchantPhoto.setQiniuKey(qiniuKey);

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
                            Log.i("qiniu", key + ",\r\n " + info + ",\r\n " + response);
                            if (Looper.myLooper() == Looper.getMainLooper()) {
                                Log.d("jiahua:", "7ncomplete" + " is in main thread");
                            } else {
                                Log.d("jiahua:", "7ncomplete" + " is in background thread");
                            }
                            new Thread(new Runnable() {
                                /**
                                 * 再次进入后台线程
                                 */
                                @Override
                                public void run() {
                                    uploadImageToQiniu(index + 1);
                                }
                            }).start();
                        }
                    }, null);
        } catch (Exception ex) {
            // onError 在后台线程
            listener.onError(ex);
        }
    }
}
