package com.cardinfolink.yunshouyin.model;


import android.net.Uri;

/**
 * 商户上传的照片
 */
public class MerchantPhoto {

    private Uri imageUri;
    private String qiniuKey;
    private String fileType;
    private String filename;


    /**
     * merName(店铺名称),
     * merAddr（店铺地址）,
     * legalCertPos（法人证书正面）,
     * legalCertOpp（法人证书反面）,
     * businessLicense（营业执照）,
     * taxRegistCert（税务登记证）,
     * organizeCodeCert（组织机构代码证）
     */
    private String imageName;//标记区别 图片的名字

    public MerchantPhoto(Uri imageUri, String filename) {
        this.imageUri = imageUri;
        this.filename = filename;
    }

    public String getFilename() {
        return filename;
    }

    public void setFilename(String filename) {
        this.filename = filename;
    }

    public String getQiniuKey() {
        return qiniuKey;
    }

    public void setQiniuKey(String qiniuKey) {
        this.qiniuKey = qiniuKey;
    }

    public Uri getImageUri() {
        return imageUri;
    }

    public void setImageUri(Uri imageUri) {
        this.imageUri = imageUri;
    }

    public String getFileType() {
        return fileType;
    }

    public void setFileType(String fileType) {
        this.fileType = fileType;
    }

    public String getImageName() {
        return imageName;
    }

    public void setImageName(String imageName) {
        this.imageName = imageName;
    }
}

