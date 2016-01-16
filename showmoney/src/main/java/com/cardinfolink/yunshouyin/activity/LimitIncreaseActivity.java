package com.cardinfolink.yunshouyin.activity;

import android.annotation.TargetApi;
import android.content.ContentUris;
import android.content.Intent;
import android.database.Cursor;
import android.graphics.Bitmap;
import android.graphics.BitmapFactory;
import android.net.Uri;
import android.os.Build;
import android.os.Bundle;
import android.os.Environment;
import android.provider.DocumentsContract;
import android.provider.MediaStore;
import android.text.TextUtils;
import android.util.Log;
import android.view.View;
import android.widget.Button;
import android.widget.TextView;

import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.api.QuickPayException;
import com.cardinfolink.yunshouyin.core.QiniuCallbackListener;
import com.cardinfolink.yunshouyin.core.QuickPayCallbackListener;
import com.cardinfolink.yunshouyin.data.SessonData;
import com.cardinfolink.yunshouyin.data.User;
import com.cardinfolink.yunshouyin.model.MerchantPhoto;
import com.cardinfolink.yunshouyin.ui.SettingActionBarItem;
import com.cardinfolink.yunshouyin.ui.SettingClikcItem;
import com.cardinfolink.yunshouyin.ui.SettingInputItem;
import com.cardinfolink.yunshouyin.util.Utility;
import com.cardinfolink.yunshouyin.view.SelectPicDialog;
import com.cardinfolink.yunshouyin.view.YellowTips;
import com.qiniu.android.http.ResponseInfo;

import org.json.JSONObject;

import java.io.File;
import java.text.NumberFormat;
import java.text.SimpleDateFormat;
import java.util.Date;
import java.util.Hashtable;
import java.util.Map;


/**
 * 提升限额的界面,用于上传图片，填写店铺名称
 */
public class LimitIncreaseActivity extends BaseActivity implements View.OnClickListener {
    private static final String TAG = "LimitIncreaseActivity";
    private static final String QINIU_FORMAT = "app/%s/%s";
    private static final String TYPE = "type";
    private static final int PERSON = 0;
    private static final int COMPANY = 1;
    private int mType;

    private SettingClikcItem mTax;//税务
    private SettingClikcItem mOrganization;//组织结构照片

    private SettingInputItem mShopName;//商铺名称
    private SettingInputItem mShopAddress;//商铺地址
    private SettingClikcItem mCardPositive;//身份证 正面
    private SettingClikcItem mCardNegative;//身份证 反面
    private SettingClikcItem mBusiness;//营业执照

    private Button mFinish;//完成按钮

    private TextView mMessage;

    private SettingActionBarItem mActionBar;

    private YellowTips mYellowTips;

    private SelectPicDialog selectPic;
    private static final int PICK_ID_P_REQUEST = 1;//身份证 正面
    private static final int PICK_ID_N_REQUEST = 2;//身份证 反面
    private static final int PICK_B_REQUEST = 3;//营业执照
    private static final int PICK_TAX_REQUEST = 4;//税务
    private static final int PICK_O_REQUEST = 5;//组织机构

    //这里创建一个map。来保存商户要上传的照片，个体的是三张，企业是五张
    private Map<Integer, MerchantPhoto> imageMap = new Hashtable<Integer, MerchantPhoto>(5);
    private boolean isAllUploadSuccess = false;

    Bitmap wrongBitmap;

    private File mExternalCacheDir;

    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_limit_increase);
        mYellowTips = new YellowTips(this, findViewById(R.id.yellow_tips));
        mActionBar = (SettingActionBarItem) findViewById(R.id.action_bar);
        mActionBar.setLeftTextOnclickListner(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                finish();
            }
        });

        //需要输入内容的
        mShopName = (SettingInputItem) findViewById(R.id.merchant_name);//商铺名称
        mShopName.setImageViewDrawable(null);//去掉组件右边的问号
        mShopAddress = (SettingInputItem) findViewById(R.id.merchant_address);//商铺地址
        mShopAddress.setImageViewDrawable(null);

        //上传图片
        mCardPositive = (SettingClikcItem) findViewById(R.id.id_card_positive);//身份证 正面
        mCardNegative = (SettingClikcItem) findViewById(R.id.id_card_negaitive);//身份证 反面
        mBusiness = (SettingClikcItem) findViewById(R.id.business);//营业执照

        //上传图片，只有企业商户才有的
        mTax = (SettingClikcItem) findViewById(R.id.tax);//税务
        mOrganization = (SettingClikcItem) findViewById(R.id.organization);//组织机构

        mMessage = (TextView) findViewById(R.id.increase_message);

        mFinish = (Button) findViewById(R.id.btnfinish);//完成按钮

        mType = getIntent().getIntExtra(TYPE, PERSON);
        if (PERSON == mType) {
            mMessage.setText(getString(R.string.limit_increase_message));
            mTax.setVisibility(View.GONE);
            mOrganization.setVisibility(View.GONE);
        } else if (COMPANY == mType) {
            mMessage.setText(getString(R.string.limit_increase_message1));
            mTax.setVisibility(View.VISIBLE);
            mOrganization.setVisibility(View.VISIBLE);
        }

        mCardPositive.setOnClickListener(this);//身份证 正面
        mCardNegative.setOnClickListener(this);//身份证 反面
        mBusiness.setOnClickListener(this);//营业执照
        mTax.setOnClickListener(this);//税务
        mOrganization.setOnClickListener(this);//组织机构
        mFinish.setOnClickListener(this);//完成按钮

        selectPic = new SelectPicDialog(this, findViewById(R.id.select_pic_dialog));

        wrongBitmap = BitmapFactory.decodeResource(this.getResources(), R.drawable.wrong);

        //初始化一个存放图片的目录
        initCacheDir();
    }

    private void initCacheDir() {
        //创建缓存目录，系统一运行就得创建缓存目录的，
        mExternalCacheDir = mContext.getExternalCacheDir();
        //如果外部的不能用，就调用内部的
        if (mExternalCacheDir == null) {
            mExternalCacheDir = mContext.getCacheDir();
        }
        if (!mExternalCacheDir.exists()) {
            mExternalCacheDir.mkdirs();
        }

    }

    @Override
    public void onClick(View v) {
        Utility.hideInput(mContext, v);
        switch (v.getId()) {
            case R.id.id_card_positive:
                //身份证 正面
                setTakePhotoOnClickListener(PICK_ID_P_REQUEST);
                setPickPhotoOnClickListener(mCardPositive.getTitle(), PICK_ID_P_REQUEST);
                selectPic.show();
                break;
            case R.id.id_card_negaitive:
                //身份证 反面
                setTakePhotoOnClickListener(PICK_ID_N_REQUEST);
                setPickPhotoOnClickListener(mCardNegative.getTitle(), PICK_ID_N_REQUEST);
                selectPic.show();
                break;
            case R.id.business:
                //营业执照
                setTakePhotoOnClickListener(PICK_B_REQUEST);
                setPickPhotoOnClickListener(mBusiness.getTitle(), PICK_B_REQUEST);
                selectPic.show();
                break;
            case R.id.tax:
                //税务
                setTakePhotoOnClickListener(PICK_TAX_REQUEST);
                setPickPhotoOnClickListener(mTax.getTitle(), PICK_TAX_REQUEST);
                selectPic.show();
                break;
            case R.id.organization:
                //组织机构
                setTakePhotoOnClickListener(PICK_O_REQUEST);
                setPickPhotoOnClickListener(mOrganization.getTitle(), PICK_O_REQUEST);
                selectPic.show();
                break;
            case R.id.btnfinish:
                uploadPhoto();
                break;
        }

    }


    private void setTakePhotoOnClickListener(final int requestCode) {
        final Intent intent = new Intent(MediaStore.ACTION_IMAGE_CAPTURE);
        File imageFile = new File(mExternalCacheDir, "yunshouyin_" + requestCode + "_.jpg");
        if (imageFile.exists()) {
            imageFile.delete();
        }
        Uri uri = Uri.fromFile(imageFile);
        intent.putExtra(MediaStore.EXTRA_OUTPUT, uri);

        selectPic.setTakePhotoOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                startActivityForResult(intent, requestCode);
                selectPic.hide();
            }
        });

    }

    //设置选择照片 按钮的监听事件，不同的 title，不同的rquestCode
    private void setPickPhotoOnClickListener(final String title, final int requestCode) {
        final Intent intent = new Intent();
        // 开启Pictures画面Type设定为image
        intent.setType("image/*");
        // 使用Intent.ACTION_GET_CONTENT这个Action
        intent.setAction(Intent.ACTION_GET_CONTENT);

        //组织机构
        selectPic.setPickPhotoOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                //取得相片后返回本画面
                startActivityForResult(Intent.createChooser(intent, title), requestCode);
                selectPic.hide();
            }
        });
    }

    private boolean validate(String shopName, String shopAddress) {
        //先检查是否都填写了

        //先检验 是否都选择好照片了
        String unselected = getResources().getString(R.string.limit_increase_unselected);
        String alertMsg = "";

        if (TextUtils.isEmpty(shopName)) {
            alertMsg = getResources().getString(R.string.alert_error__shop_name_cannot_empty);//店铺名称不能为空
            mYellowTips.show(alertMsg);
            return false;
        }

        if (TextUtils.isEmpty(shopAddress)) {
            alertMsg = getResources().getString(R.string.alert_error__shop_address_cannot_empty);//店铺地址不能为空
            mYellowTips.show(alertMsg);
            return false;
        }

        if (imageMap.get(PICK_ID_P_REQUEST) == null) {
            // 身份证 正面
            alertMsg = mCardPositive.getTitle() + unselected;
            mYellowTips.show(alertMsg);
            return false;
        }

        if (imageMap.get(PICK_ID_N_REQUEST) == null) {
            // 身份证 反面
            alertMsg = mCardNegative.getTitle() + unselected;
            mYellowTips.show(alertMsg);
            return false;
        }

        if (imageMap.get(PICK_B_REQUEST) == null) {
            // 营业执照
            alertMsg = mBusiness.getTitle() + unselected;
            mYellowTips.show(alertMsg);
            return false;
        }

        if (mType == COMPANY) {//当是企业用户的时候要多判断这两个
            if (imageMap.get(PICK_TAX_REQUEST) == null) {
                ////税务
                alertMsg = mTax.getTitle() + unselected;
                mYellowTips.show(alertMsg);
                return false;
            }
            if (imageMap.get(PICK_O_REQUEST) == null) {
                //组织机构
                alertMsg = mOrganization.getTitle() + unselected;
                mYellowTips.show(alertMsg);
                return false;
            }
        }

        return true;
    }

    private void uploadPhoto() {
        String shopName = mShopName.getText();//商户名称
        String shopAddress = mShopAddress.getText();//商户地址

        if (!validate(shopName, shopAddress)) {
            //先检查是否都填写了
            return;
        }

        if (isAllUploadSuccess) {
            //都上传成功了 updateUser
            return;
        }

        Date now = new Date();
        SimpleDateFormat yyMMdd = new SimpleDateFormat("yyyyMMdd");
        String clientId = SessonData.loginUser.getClientid();
        final String qiniuKeyPattern = String.format(QINIU_FORMAT, yyMMdd.format(now), clientId) + "/%s.%s";

        startLoading();

        //如果map size很大这里就不太好了 后台线程,这里会根据map的大小创建多少个子线程用来上传
        qiniuMultiUploadService.upload(imageMap, qiniuKeyPattern, new QiniuCallbackListener() {
            @Override
            public void onComplete(String key, ResponseInfo info, JSONObject response, int photoKey) {
                updateUploadSuccess(photoKey);
            }

            @Override
            public void onFailure(QuickPayException ex, int photoKey) {
                //失败时根据photoKey来更新不同的ui组件
                updateUploadFailure(ex, photoKey);
            }

            @Override
            public void onProgress(String key, double percent, int photoKey) {
                updateUploadProgress(key, percent, photoKey);
            }
        });

    }


    private boolean checkAllUploadSuccess() {
        //这里遍历map，判断是否都上传成功了,这里只是遍历一下value
        for (MerchantPhoto merchantPhoto : imageMap.values()) {
            //如果发现qiniuKey是null就表明还没有上传成功
            if (merchantPhoto.getQiniuKey() == null) {
                return false;
            }
        }
        return true;
    }

    private void updateUploadSuccess(int photoKey) {
        String successStr = getString(R.string.limit_increase_upload_success);
        switch (photoKey) {//根据不同的photoKey来更新不同的ui组件
            case PICK_ID_P_REQUEST: // 身份证 正面
                mCardPositive.setRightText(successStr);
                mCardPositive.setImageResource(R.drawable.login_added);
                Log.e(TAG, " upload success == 身份证 正面 ");
                break;
            case PICK_ID_N_REQUEST://身份证 反面
                mCardNegative.setRightText(successStr);
                mCardNegative.setImageResource(R.drawable.login_added);
                Log.e(TAG, " upload success == 身份证 反面 ");
                break;
            case PICK_B_REQUEST://营业执照
                mBusiness.setRightText(successStr);
                mBusiness.setImageResource(R.drawable.login_added);
                Log.e(TAG, " upload success == 营业执照 ");
                break;
            case PICK_TAX_REQUEST://税务
                mTax.setRightText(successStr);
                mTax.setImageResource(R.drawable.login_added);
                Log.e(TAG, " upload success == 税务 ");
                break;
            case PICK_O_REQUEST://组织机构
                mOrganization.setRightText(successStr);
                mOrganization.setImageResource(R.drawable.login_added);
                Log.e(TAG, " upload success == 组织机构 ");
                break;
        }

        //因为是有五个线程同时在上传，这时候走到这里表明有一个上传成功了。
        //判断一下是否都上传成功了，都上传成功了 才能进行下一步。1.关闭load。
        if (checkAllUploadSuccess()) {
            Log.e(TAG, " all upload success");
            endLoading();//这个要判断一下是否map中图片都上传失败，或者成功，然后才能结束loading的动画。
            isAllUploadSuccess = true;

            improveCertInfo();
        }
    }

    private void improveCertInfo() {
        /**
         * merName(店铺名称),
         * merAddr（店铺地址）,
         * legalCertPos（法人证书正面）,
         * legalCertOpp（法人证书反面）,
         * businessLicense（营业执照）,
         * taxRegistCert（税务登记证）,
         * organizeCodeCert（组织机构代码证）
         */
        Map<String, String> map = new Hashtable<>();
        for (MerchantPhoto merchantPhoto : imageMap.values()) {
            //这里要不要判断一下是否为空呢？？？
            map.put(merchantPhoto.getImageName(), merchantPhoto.getQiniuKey());
        }

        User user = SessonData.loginUser;
        String shopName = mShopName.getText();//商户名称
        String shopAddress = mShopAddress.getText();//商户地址
        quickPayService.improveCertInfoAsync(user, shopName, shopAddress, map, new QuickPayCallbackListener<Void>() {
            @Override
            public void onSuccess(Void data) {
                //improve成功，接下来 activateUser激活
                activateUser();
                Log.e(TAG, " improveCertInfoAsync  onSuccess");
            }

            @Override
            public void onFailure(QuickPayException ex) {
                Log.e(TAG, "improveCertInfoAsync onFailure");
                alertShow(ex.getErrorMsg(), wrongBitmap);
            }
        });

    }

    private void activateUser() {
        String username = SessonData.loginUser.getUsername();
        String password = SessonData.loginUser.getPassword();
        // 3.激活
        quickPayService.activateAsync(username, password, new QuickPayCallbackListener<Void>() {
            @Override
            public void onSuccess(Void data) {
                Log.e(TAG, " activateAsync  onSuccess");
                Intent intent = new Intent(LimitIncreaseActivity.this, FinalIncreaseActivity.class);
                startActivity(intent);
                finish();
            }

            @Override
            public void onFailure(QuickPayException ex) {
                Log.e(TAG, " activateAsync  onFailure");
                alertShow(ex.getErrorMsg(), wrongBitmap);
            }
        });
    }

    private void updateUploadFailure(QuickPayException ex, int photoKey) {
        endLoading();//这个要判断一下是否map中图片都上传失败，或者成功，然后才能结束loading的动画。
        String failStr = getString(R.string.limit_increase_upload_fail);

        switch (photoKey) {//根据不同的photoKey来更新不同的ui组件
            case PICK_ID_P_REQUEST: // 身份证 正面
                mCardPositive.setRightText(failStr);
                Log.e(TAG, " upload Failure == 身份证 正面 ");
                break;
            case PICK_ID_N_REQUEST://身份证 反面
                mCardNegative.setRightText(failStr);
                Log.e(TAG, " upload Failure == 身份证 反面 ");
                break;
            case PICK_B_REQUEST://营业执照
                mBusiness.setRightText(failStr);
                Log.e(TAG, " upload Failure == 营业执照 ");
                break;
            case PICK_TAX_REQUEST://税务
                mTax.setRightText(failStr);
                Log.e(TAG, " upload Failure == 税务 ");
                break;
            case PICK_O_REQUEST://组织机构
                mOrganization.setRightText(failStr);
                Log.e(TAG, " upload Failure == 组织机构 ");
                break;
            default:
                Log.e(TAG, " upload Failure ==  default ");

                ex.getErrorMsg();
                String alertMsg = ex.getErrorMsg();
                mAlertDialog.show(alertMsg, wrongBitmap);
                break;
        }
    }

    //更新上传进度，根据不同的photokey来更新不同的ui组件
    private void updateUploadProgress(String key, double percent, int photoKey) {
        NumberFormat fmt = NumberFormat.getPercentInstance();
        fmt.setMaximumFractionDigits(2);//最多两位百分小数，如25.23%
        String percentString = fmt.format(percent);
        switch (photoKey) {//根据不同的photoKey来更新不同的ui组件
            case PICK_ID_P_REQUEST: // 身份证 正面
                mCardPositive.setRightText(percentString);
                break;
            case PICK_ID_N_REQUEST://身份证 反面
                mCardNegative.setRightText(percentString);
                break;
            case PICK_B_REQUEST://营业执照
                mBusiness.setRightText(percentString);
                break;
            case PICK_TAX_REQUEST://税务
                mTax.setRightText(percentString);
                break;
            case PICK_O_REQUEST://组织机构
                mOrganization.setRightText(percentString);
                break;
        }
    }

    @Override
    protected void onActivityResult(int requestCode, int resultCode, Intent data) {
        String selectedStr = getString(R.string.limit_increase_selected);//已选择
        String unselectedStr = getString(R.string.limit_increase_update_licence);//上传证件,相当于提示 用户 未选择图片
        switch (requestCode) {
            case PICK_ID_P_REQUEST:            // 身份证 正面
                if (resultCode == RESULT_OK) {
                    mCardPositive.setRightText(selectedStr);
                    imageMap.put(PICK_ID_P_REQUEST, getMerchantPhoto(requestCode, data));
                    imageMap.get(PICK_ID_P_REQUEST).setImageName("legalCertPos");//标记一下这个图片的名称
                } else {
                    mCardPositive.setRightText(unselectedStr);
                    //这里取消选择的照片，这里设置为null，当在使用的时候要检查是否为null
                    imageMap.remove(PICK_ID_P_REQUEST);//取消选择的照片，也就是删除map中对应key的value
                }
                break;
            case PICK_ID_N_REQUEST://身份证 反面
                if (resultCode == RESULT_OK) {
                    mCardNegative.setRightText(selectedStr);
                    imageMap.put(PICK_ID_N_REQUEST, getMerchantPhoto(requestCode, data));
                    imageMap.get(PICK_ID_N_REQUEST).setImageName("legalCertOpp");//标记一下这个图片的名称
                } else {
                    mCardNegative.setRightText(unselectedStr);
                    imageMap.remove(PICK_ID_N_REQUEST);//取消选择的照片也就是删除map中对应key的value
                }
                break;
            case PICK_B_REQUEST://营业执照
                if (resultCode == RESULT_OK) {
                    mBusiness.setRightText(selectedStr);
                    imageMap.put(PICK_B_REQUEST, getMerchantPhoto(requestCode, data));
                    imageMap.get(PICK_B_REQUEST).setImageName("businessLicense");//标记一下这个图片的名称
                } else {
                    mBusiness.setRightText(unselectedStr);
                    imageMap.remove(PICK_B_REQUEST);//取消选择的照片也就是删除map中对应key的value
                }
                break;
            case PICK_TAX_REQUEST://税务
                if (resultCode == RESULT_OK) {
                    mTax.setRightText(selectedStr);
                    imageMap.put(PICK_TAX_REQUEST, getMerchantPhoto(requestCode, data));
                    imageMap.get(PICK_TAX_REQUEST).setImageName("taxRegistCert");//标记一下这个图片的名称
                } else {
                    mTax.setRightText(unselectedStr);
                    imageMap.remove(PICK_TAX_REQUEST);//取消选择的照片也就是删除map中对应key的value
                }
                break;
            case PICK_O_REQUEST://组织机构
                if (resultCode == RESULT_OK) {
                    mOrganization.setRightText(selectedStr);
                    imageMap.put(PICK_O_REQUEST, getMerchantPhoto(requestCode, data));
                    imageMap.get(PICK_O_REQUEST).setImageName("organizeCodeCert");//标记一下这个图片的名称
                } else {
                    mOrganization.setRightText(unselectedStr);
                    imageMap.remove(PICK_O_REQUEST);//取消选择的照片也就是删除map中对应key的value
                }
                break;
        }
        super.onActivityResult(requestCode, resultCode, data);
    }

    private MerchantPhoto getMerchantPhoto(int requestCode, Intent data) {
        if (data == null) {//等于null说明是拍照得来的
            File imageFile = new File(mExternalCacheDir, "yunshouyin_" + requestCode + "_.jpg");
            Uri uri = Uri.fromFile(imageFile);
            return new MerchantPhoto(uri, imageFile.getAbsolutePath());
        } else {
            return getMerchantPhoto(data);
        }
    }

    private MerchantPhoto getMerchantPhoto(Intent data) {
        if (Build.VERSION.SDK_INT >= 19) {
            // 4.4以上使用
            return handleImageOnKitKat(data);
        } else {
            return handleImageBeforeKitKat(data);
        }
    }

    public MerchantPhoto handleImageBeforeKitKat(Intent data) {
        if (data != null) {
            Uri uri = data.getData();
            String imagePath = getImagePath(uri, null);
            return new MerchantPhoto(uri, imagePath);
        }
        return null;
    }

    @TargetApi(Build.VERSION_CODES.KITKAT)
    public MerchantPhoto handleImageOnKitKat(Intent data) {
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

            return new MerchantPhoto(uri, imagePath);
        }
        return null;
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

}
