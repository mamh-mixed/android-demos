package com.cardinfolink.yunshouyin.core;

import com.cardinfolink.yunshouyin.api.QuickPayException;
import com.cardinfolink.yunshouyin.data.SessionData;
import com.cardinfolink.yunshouyin.model.MerchantPhoto;
import com.qiniu.android.http.ResponseInfo;
import com.qiniu.android.storage.UpCompletionHandler;
import com.qiniu.android.storage.UpProgressHandler;
import com.qiniu.android.storage.UploadManager;
import com.qiniu.android.storage.UploadOptions;

import org.json.JSONObject;

import java.util.Map;
import java.util.UUID;

public class QiniuMultiUploadService {
    private static final String TAG = "QiniuMultiUploadService";

    private QuickPayService quickPayService;
    private UploadManager uploadManager = new UploadManager();
    private Map<Integer, MerchantPhoto> photoMap;
    private QiniuCallbackListener mListener;
    private String mQiniuKeyPattern;

    public QiniuMultiUploadService(QuickPayService quickPayService) {
        this.quickPayService = quickPayService;
    }

    /**
     * 同时只能有一个upload使用, 加上同步关键字
     *
     * @param photoMap
     * @param qiniuKeyPattern
     * @param listener
     */
    public synchronized void upload(final Map<Integer, MerchantPhoto> photoMap, String qiniuKeyPattern, final QiniuCallbackListener listener) {
        //UI线程内
        if (photoMap == null || photoMap.size() == 0) {
            //map是空的，这个不太可能
            return;
        }

        this.photoMap = photoMap;
        this.mListener = listener;
        this.mQiniuKeyPattern = qiniuKeyPattern;

        quickPayService.getUploadTokenAsync(SessionData.loginUser, new QuickPayCallbackListener<String>() {
            @Override
            public void onSuccess(final String uploadToken) {
                //进入后台线程,这里会根据map的大小创建多少个子线程用来上传
                //如果map size很大这里就不太好了
                upload(uploadToken);
            }

            @Override
            public void onFailure(QuickPayException ex) {
                //这里失败表面获取toke就失败了
                mListener.onFailure(ex, -1);
            }
        });
    }

    public void upload(final String uploadToken) {
        for (Map.Entry<Integer, MerchantPhoto> map : photoMap.entrySet()) {
            final int photoKey = map.getKey();
            final MerchantPhoto merchantPhoto = map.getValue();
            if (merchantPhoto == null) {
                //这里判断一下是不是null,hashTable 的value不能是null
                continue;
            }

            String filename = merchantPhoto.getFilename();
            //NOTE: 文件可能没有后缀名
            final String fileType = filename.lastIndexOf(".") >= 0 ? filename.substring(filename.lastIndexOf(".") + 1) : "";
            final String qiniuKey = String.format(mQiniuKeyPattern, UUID.randomUUID().toString().replace("-", ""), fileType);


            //上传完成的 回调
            UpCompletionHandler completionHandler = new UpCompletionHandler() {
                @Override
                public void complete(String key, ResponseInfo info, JSONObject response) {
                    if (info.isOK()) {
                        merchantPhoto.setQiniuKey(qiniuKey);
                        merchantPhoto.setFileType(fileType);
                        mListener.onComplete(key, info, response, photoKey);
                    } else {
                        merchantPhoto.setQiniuKey(null);
                        merchantPhoto.setFileType(null);
                        mListener.onFailure(new QuickPayException(), photoKey);
                    }
                }
            };

            //上传正在进行中的回调
            UpProgressHandler upProgressHandler = new UpProgressHandler() {
                public void progress(String key, double percent) {
                    mListener.onProgress(key, percent, photoKey);//photoKey标记是正在上传的哪个照片
                }
            };

            //上传的 选项
            UploadOptions uploadOptions = new UploadOptions(null, null, false, upProgressHandler, null);

            //这里开始真正的上传
            uploadManager.put(filename, qiniuKey, uploadToken, completionHandler, uploadOptions);
        }
    }
}
