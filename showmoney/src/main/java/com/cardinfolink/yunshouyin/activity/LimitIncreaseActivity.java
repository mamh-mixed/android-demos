package com.cardinfolink.yunshouyin.activity;

import android.content.ContentUris;
import android.content.Intent;
import android.database.Cursor;
import android.net.Uri;
import android.os.Build;
import android.os.Bundle;
import android.provider.DocumentsContract;
import android.provider.MediaStore;
import android.util.Log;
import android.view.View;
import android.widget.Button;
import android.widget.TextView;

import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.model.MerchantPhoto;
import com.cardinfolink.yunshouyin.ui.SettingActionBarItem;
import com.cardinfolink.yunshouyin.ui.SettingClikcItem;
import com.cardinfolink.yunshouyin.ui.SettingInputItem;
import com.cardinfolink.yunshouyin.view.SelectPicDialog;


/**
 * 提升限额的界面,用于上传图片，填写店铺名称
 */
public class LimitIncreaseActivity extends BaseActivity implements View.OnClickListener {
    private static final String TAG = "LimitIncreaseActivity";

    private static final String TYPE = "type";
    private static final int PERSON = 0;
    private static final int COMPANY = 1;
    private int mType;

    private SettingClikcItem mTax;//税务
    private SettingClikcItem mOrganization;//组织结构照片

    private SettingInputItem mMerchant;//商铺名称
    private SettingInputItem mMerchantAddress;//商铺地址
    private SettingClikcItem mCardPositive;//身份证 正面
    private SettingClikcItem mCardNegative;//身份证 反面
    private SettingClikcItem mBusiness;//营业执照

    private Button mFinish;//完成按钮

    private TextView mMessage;

    private SettingActionBarItem mActionBar;

    private SelectPicDialog selectPic;
    private static final int PICK_ID_P_REQUEST = 1;//身份证 正面
    private static final int PICK_ID_N_REQUEST = 2;//身份证 反面
    private static final int PICK_B_REQUEST = 3;//营业执照
    private static final int PICK_TAX_REQUEST = 4;//税务
    private static final int PICK_O_REQUEST = 5;//组织机构

    private MerchantPhoto[] imageList = new MerchantPhoto[5];


    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_limit_increase);

        mActionBar = (SettingActionBarItem) findViewById(R.id.action_bar);
        mActionBar.setLeftTextOnclickListner(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                finish();
            }
        });

        //需要输入内容的
        mMerchant = (SettingInputItem) findViewById(R.id.merchant_name);//商铺名称
        mMerchantAddress = (SettingInputItem) findViewById(R.id.merchant_address);//商铺地址

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
    }


    @Override
    public void onClick(View v) {
        final Intent intent = new Intent();
        // 开启Pictures画面Type设定为image
        intent.setType("image/*");
        // 使用Intent.ACTION_GET_CONTENT这个Action
        intent.setAction(Intent.ACTION_GET_CONTENT);

        switch (v.getId()) {
            case R.id.id_card_positive:
                //身份证 正面
                selectPic.setPickPhotoOnClickListener(new View.OnClickListener() {
                    @Override
                    public void onClick(View v) {
                        //取得相片后返回本画面
                        String title = getResources().getString(R.string.limit_increase_card_positive);
                        startActivityForResult(Intent.createChooser(intent, title), PICK_ID_P_REQUEST);
                        selectPic.hide();
                    }
                });
                selectPic.show();
                break;
            case R.id.id_card_negaitive:
                //身份证 反面
                selectPic.setPickPhotoOnClickListener(new View.OnClickListener() {
                    @Override
                    public void onClick(View v) {
                        //取得相片后返回本画面
                        String title = getResources().getString(R.string.limit_increase_card_negaitive);
                        startActivityForResult(Intent.createChooser(intent, title), PICK_ID_N_REQUEST);
                        selectPic.hide();
                    }
                });
                selectPic.show();
                break;
            case R.id.business:
                //营业执照
                selectPic.setPickPhotoOnClickListener(new View.OnClickListener() {
                    @Override
                    public void onClick(View v) {
                        //取得相片后返回本画面
                        String title = getResources().getString(R.string.limit_increase_business_licence);
                        startActivityForResult(Intent.createChooser(intent, title), PICK_B_REQUEST);
                        selectPic.hide();
                    }
                });
                selectPic.show();
                break;
            case R.id.tax:
                //税务
                selectPic.setPickPhotoOnClickListener(new View.OnClickListener() {
                    @Override
                    public void onClick(View v) {
                        //取得相片后返回本画面
                        String title = getResources().getString(R.string.limit_increase_tax);
                        startActivityForResult(Intent.createChooser(intent, title), PICK_TAX_REQUEST);
                        selectPic.hide();
                    }
                });
                selectPic.show();
                break;
            case R.id.organization:
                //组织机构
                selectPic.setPickPhotoOnClickListener(new View.OnClickListener() {
                    @Override
                    public void onClick(View v) {
                        //取得相片后返回本画面
                        String title = getResources().getString(R.string.limit_increase_organization);
                        startActivityForResult(Intent.createChooser(intent, title), PICK_O_REQUEST);
                        selectPic.hide();
                    }
                });
                selectPic.show();
                break;
            case R.id.btnfinish:
                break;
        }

    }


    @Override
    protected void onActivityResult(int requestCode, int resultCode, Intent data) {
        Log.e(TAG, "===========================onActivityResult+++++++++++++++++++++++++++" + resultCode);
        String selectedStr = getString(R.string.limit_increase_selected);//已选择
        String unselectedStr = getString(R.string.limit_increase_update_licence);//上传证件,相当于提示 用户 未选择图片
        switch (requestCode) {
            case PICK_ID_P_REQUEST:
                if (resultCode == RESULT_OK) {
                    mCardPositive.setRightText(selectedStr);
                    imageList[0] = getMerchantPhoto(data);
                } else {
                    mCardPositive.setRightText(unselectedStr);
                    imageList[0] = null;
                }
                break;
            case PICK_ID_N_REQUEST:
                if (resultCode == RESULT_OK) {
                    mCardNegative.setRightText(selectedStr);
                    imageList[1] = getMerchantPhoto(data);
                } else {
                    mCardNegative.setRightText(unselectedStr);
                    imageList[1] = null;
                }
                break;
            case PICK_B_REQUEST:
                if (resultCode == RESULT_OK) {
                    mBusiness.setRightText(selectedStr);
                    imageList[2] = getMerchantPhoto(data);
                } else {
                    mBusiness.setRightText(unselectedStr);
                    imageList[2] = null;
                }
                break;
            case PICK_TAX_REQUEST:
                if (resultCode == RESULT_OK) {
                    mTax.setRightText(selectedStr);
                    imageList[3] = getMerchantPhoto(data);
                } else {
                    mTax.setRightText(unselectedStr);
                    imageList[3] = null;
                }
                break;
            case PICK_O_REQUEST:
                if (resultCode == RESULT_OK) {
                    mOrganization.setRightText(selectedStr);
                    imageList[4] = getMerchantPhoto(data);
                } else {
                    mOrganization.setRightText(unselectedStr);
                    imageList[4] = null;
                }
                break;
        }
        super.onActivityResult(requestCode, resultCode, data);
    }

    private MerchantPhoto getMerchantPhoto(Intent data) {
        if (Build.VERSION.SDK_INT >= 19) {
            // 4.4以上使用
            return handleImageOnKitKat(data);
        } else {
            return handleImageBeforeKitKat(data);
        }
    }

    private MerchantPhoto handleImageBeforeKitKat(Intent data) {
        if (data != null) {
            Uri uri = data.getData();
            String imagePath = getImagePath(uri, null);
            return new MerchantPhoto(uri, imagePath);
        }
        return null;
    }

    private MerchantPhoto handleImageOnKitKat(Intent data) {
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
