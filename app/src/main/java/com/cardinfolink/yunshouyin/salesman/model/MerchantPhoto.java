package com.cardinfolink.yunshouyin.salesman.model;


import android.net.Uri;

public class MerchantPhoto {
    private Uri imageUri;
    private String qiniuKey;
    private String fileType;
    private String filename;

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
}

