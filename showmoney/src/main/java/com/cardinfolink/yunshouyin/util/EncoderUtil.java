package com.cardinfolink.yunshouyin.util;

import android.text.TextUtils;

import java.security.MessageDigest;

public class EncoderUtil {
    public static String Encrypt(String strSrc, String encName) {
        //这里加入了 为null 的判断
        if (TextUtils.isEmpty(strSrc)) {
            return null;
        }
        //parameter strSrc is a string will be encrypted,
        //parameter encName is the algorithm name will be used.
        //encName dafault to "MD5"
        MessageDigest md = null;
        String strDes = null;

        byte[] bt = strSrc.getBytes();
        try {
            if (encName == null || encName.equals("")) {
                encName = "MD5";
            }
            md = MessageDigest.getInstance(encName);
            md.update(bt);
            strDes = bytes2Hex(md.digest()); //to HexString
        } catch (Exception e) {
            return null;
        }
        return strDes;
    }

    public static String Encrypt(byte[] bt, String encName) {
        //parameter strSrc is a string will be encrypted,
        //parameter encName is the algorithm name will be used.
        //encName dafault to "MD5"
        MessageDigest md = null;
        String strDes = null;
        try {
            if (encName == null || encName.equals("")) {
                encName = "MD5";
            }
            md = MessageDigest.getInstance(encName);
            md.update(bt);
            strDes = bytes2Hex(md.digest()); //to HexString
        } catch (Exception e) {
            return null;
        }
        return strDes;
    }


    private static String bytes2Hex(byte[] bts) {
        String des = "";
        String tmp = null;
        for (int i = 0; i < bts.length; i++) {
            tmp = (Integer.toHexString(bts[i] & 0xFF));
            if (tmp.length() == 1) {
                des += "0";
            }
            des += tmp;
        }
        return des;
    }


}
