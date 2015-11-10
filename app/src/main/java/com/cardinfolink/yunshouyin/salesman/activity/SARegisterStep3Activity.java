package com.cardinfolink.yunshouyin.salesman.activity;

import android.annotation.TargetApi;
import android.content.ContentUris;
import android.content.Intent;
import android.database.Cursor;
import android.net.Uri;
import android.os.Build;
import android.os.Bundle;
import android.os.Looper;
import android.provider.DocumentsContract;
import android.provider.MediaStore;
import android.support.v7.widget.RecyclerView;
import android.support.v7.widget.StaggeredGridLayoutManager;
import android.util.Log;
import android.view.View;
import android.widget.Button;
import android.widget.Toast;

import com.cardinfolink.yunshouyin.salesman.R;
import com.cardinfolink.yunshouyin.salesman.adapter.MerchantPhotoRecyclerViewAdapter;
import com.cardinfolink.yunshouyin.salesman.api.QuickPayException;
import com.cardinfolink.yunshouyin.salesman.core.QiniuCallbackListener;
import com.cardinfolink.yunshouyin.salesman.core.QuickPayCallbackListener;
import com.cardinfolink.yunshouyin.salesman.model.SAMerchantPhoto;
import com.cardinfolink.yunshouyin.salesman.model.SessonData;
import com.cardinfolink.yunshouyin.salesman.model.User;
import com.cardinfolink.yunshouyin.salesman.utils.ActivityCollector;

import java.text.SimpleDateFormat;
import java.util.ArrayList;
import java.util.Date;
import java.util.List;

public class SARegisterStep3Activity extends BaseActivity {
    private static final String LOG_TAG = "RegisterStep3";
    private static final int PICK_IMAGE_REQUEST = 1000;
    private StaggeredGridLayoutManager staggeredGridLayoutManager;
    private List<SAMerchantPhoto> imageList;
    private MerchantPhotoRecyclerViewAdapter rcAdapter;
    private RecyclerView recyclerView;

    private Button btnActivate;
    private Button btnPickImage;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_saregister_step3);

        btnActivate = (Button) findViewById(R.id.btnActivate);
        btnActivate.setEnabled(false);

        btnPickImage = (Button) findViewById(R.id.btnPickImage);
        btnPickImage.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                Intent intent = new Intent(Intent.ACTION_GET_CONTENT);
                //TODO: 多张选择,原生没有,需要第三方库
                intent.putExtra(Intent.EXTRA_ALLOW_MULTIPLE, true);
                intent.setType("image/*");
                startActivityForResult(Intent.createChooser(intent, "选择图片"), PICK_IMAGE_REQUEST);
            }
        });

        // set up recyclerView
        recyclerView = (RecyclerView) findViewById(R.id.recycler_view);
        recyclerView.setHasFixedSize(true);

        // layout manager, 2 columns and VERTICAL
        staggeredGridLayoutManager = new StaggeredGridLayoutManager(2, StaggeredGridLayoutManager.VERTICAL);
        recyclerView.setLayoutManager(staggeredGridLayoutManager);

        // fill data adapter
        imageList = getListItemData();
        rcAdapter = new MerchantPhotoRecyclerViewAdapter(SARegisterStep3Activity.this, imageList);
        recyclerView.setAdapter(rcAdapter);
    }

    private List<SAMerchantPhoto> getListItemData() {
        List<SAMerchantPhoto> listViewItems = new ArrayList<SAMerchantPhoto>();
        return listViewItems;
    }


    public void confirmUpload(View view) {
        startLoading();

        // tools/20150202/clientId/uuid.png
        Date now = new Date();
        SimpleDateFormat yyMMdd = new SimpleDateFormat("yyyyMMdd");
        String clientId = SessonData.registerUser.getClientid();
        final String qiniuKeyPattern = String.format("tools/%s/%s", yyMMdd.format(now), clientId) + "/%s.%s";

        // imageList会生成出qiniukey出来
        // 1. upload images to qiniu server
        application.getQiniuMultiUploadService().upload(imageList, qiniuKeyPattern, new QiniuCallbackListener() {
            @Override
            public void onComplete() {
                if (Looper.myLooper() == Looper.getMainLooper()) {
                    Log.d(LOG_TAG, "onComplete" + " is in main thread");
                } else {
                    Log.d(LOG_TAG, "onComplete" + " is in background thread");
                }

                // 7n上传成功
                runOnUiThread(new Runnable() {
                    @Override
                    public void run() {
                        Toast.makeText(SARegisterStep3Activity.this, "上传图片成功中,请等待", Toast.LENGTH_SHORT);
                    }
                });

                // 2. update to cardinfolink server
                final User user = SessonData.registerUser;
                List<String> images = new ArrayList<String>();
                for (SAMerchantPhoto SAMerchantPhoto : imageList) {
                    images.add(SAMerchantPhoto.getQiniuKey());
                }
                user.setImages((images.toArray(new String[images.size()])));


                application.getQuickPayService().updateUser(user, new QuickPayCallbackListener<User>() {
                    @Override
                    public void onSuccess(User data) {
                        // TODO: no check user result
                        runOnUiThread(new Runnable() {
                            @Override
                            public void run() {
                                Toast.makeText(SARegisterStep3Activity.this, "更新到服务器,激活中", Toast.LENGTH_SHORT);
                            }
                        });

                        // 3.激活
                        application.getQuickPayService().activateUser(user.getUsername(), new QuickPayCallbackListener<User>() {
                            @Override
                            public void onSuccess(User data) {
                                runOnUiThread(new Runnable() {
                                    @Override
                                    public void run() {
                                        endLoading();
                                        Toast.makeText(SARegisterStep3Activity.this, "成功新增商户，参数已经发送到您的邮箱和商户邮箱，请查收。", Toast.LENGTH_LONG);
                                        ActivityCollector.goHomeAndFinishRest();
//                                        alertInfo("成功新增商户，参数已经发送到您的邮箱和商户邮箱，请查收。", new WorkBeforeExitListener() {
//                                            @Override
//                                            public void complete() {
//                                                ActivityCollector.goHomeAndFinishRest();
//                                            }
//                                        });
                                    }
                                });
                            }

                            @Override
                            public void onFailure(final QuickPayException ex) {
                                runOnUiThread(new Runnable() {
                                    @Override
                                    public void run() {
                                        String errorStr = ex.getErrorMsg();
                                        endLoadingWithError(errorStr);
                                    }
                                });
                            }
                        });
                    }

                    @Override
                    public void onFailure(final QuickPayException ex) {
                        runOnUiThread(new Runnable() {
                            @Override
                            public void run() {
                                endLoadingWithError(ex.getErrorMsg());
                            }
                        });
                    }
                });
            }

            @Override
            public void onFailure(Exception ex) {
                Log.d(LOG_TAG, ex.getMessage());
                if (Looper.myLooper() == Looper.getMainLooper()) {
                    Log.d(LOG_TAG, "onFailure" + " is in main thread");
                } else {
                    Log.d(LOG_TAG, "onFailure" + " is in background thread");
                }
                runOnUiThread(new Runnable() {
                    @Override
                    public void run() {
                        endLoadingWithError("上传图片失败");
                    }
                });
            }
        });
    }

    @Override
    protected void onActivityResult(int requestCode, int resultCode, Intent data) {
        Log.d(LOG_TAG, String.format("requestCode: %d, resultCode: %d", requestCode, resultCode));
        switch (requestCode) {
            case PICK_IMAGE_REQUEST:
                //至少有一张图之后确认激活才能打开
                btnActivate.setEnabled(true);

                if (resultCode == RESULT_OK) {
                    if (Build.VERSION.SDK_INT >= 19) {
                        // 4.4以上使用
                        handleImageOnKitKat(data);
                    } else {
                        handleImageBeforeKitKat(data);
                    }
                }
                break;
            default:
                break;
        }
    }

    private void handleImageBeforeKitKat(Intent data) {
        if (data != null) {
            Uri uri = data.getData();
            String imagePath = getImagePath(uri, null);
            Log.d(LOG_TAG, imagePath);
            imageList.add(new SAMerchantPhoto(uri, imagePath));
            rcAdapter.notifyDataSetChanged();
//            staggeredGridLayoutManager.onItemsAdded(recyclerView, 0, 1);
//            staggeredGridLayoutManager.onItemsChanged(recyclerView);
        }
    }

    @TargetApi(Build.VERSION_CODES.KITKAT)
    private void handleImageOnKitKat(Intent data) {
        if (data != null) {
            String imagePath = null;
            Uri uri = data.getData();
            if (DocumentsContract.isDocumentUri(this, uri)) {
                String docId = DocumentsContract.getDocumentId(uri);
                if ("com.android.providers.media.documents".equals(uri.getAuthority())) {
                    String id = docId.split(":")[1];
                    String selection = MediaStore.Images.Media._ID + "=" + id;
                    imagePath = getImagePath(MediaStore.Images.Media.EXTERNAL_CONTENT_URI, selection);
                } else if ("com.android.providers.downloads.documents".equals(uri.getAuthority())) {
                    Uri contentUri = ContentUris.withAppendedId(Uri.parse("content://downloads/public_downloads"), Long.valueOf(docId));
                    imagePath = getImagePath(contentUri, null);
                }
            } else if ("content".equalsIgnoreCase(uri.getScheme())) {
                imagePath = getImagePath(uri, null);
            }

            Log.d(LOG_TAG, imagePath);
            imageList.add(new SAMerchantPhoto(uri, imagePath));
            rcAdapter.notifyDataSetChanged();
        }
    }

    private String getImagePath(Uri uri, String selection) {
        String path = null;
        Cursor cursor = getContentResolver().query(uri, null, selection, null, null);
        if (cursor != null) {
            if (cursor.moveToFirst()) {
                path = cursor.getString(cursor.getColumnIndex(MediaStore.Images.Media.DATA));
            }
            cursor.close();
        }
        return path;
    }

    /**
     * callback, ViewHolder调用
     *
     * @param index
     */
    public void removeItemAt(int index) {
        imageList.remove(index);
        rcAdapter.notifyDataSetChanged();
    }
}
