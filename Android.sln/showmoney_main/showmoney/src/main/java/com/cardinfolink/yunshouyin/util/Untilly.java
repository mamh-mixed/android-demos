package com.cardinfolink.yunshouyin.util;

import android.graphics.Bitmap;
import android.graphics.Matrix;

import java.io.ByteArrayOutputStream;
import java.io.File;
import java.io.FileInputStream;
import java.io.FileOutputStream;
import java.io.IOException;


/**
 * ������
 *
 * @author LJC
 */
public class Untilly {
    /**
     * ����ͼƬ
     *
     * @param icon
     * @param h
     * @return
     */
    public static Bitmap zoomBitmap(Bitmap icon, int h) {
        // ����ͼƬ
        Matrix m = new Matrix();
        float sx = (float) 2 * h / icon.getWidth();
        float sy = (float) 2 * h / icon.getHeight();
        m.setScale(sx, sy);
        // ���¹���һ��2h*2h��ͼƬ
        return Bitmap.createBitmap(icon, 0, 0, icon.getWidth(), icon.getHeight(), m, false);
    }


    /**
     * @param fileName
     * @return byte[]
     */
    public static byte[] readFile(String fileName) {
        FileInputStream fis = null;
        ByteArrayOutputStream baos = null;
        byte[] data = null;
        try {
            fis = new FileInputStream(fileName);
            byte[] buffer = new byte[8 * 1024];
            int readSize = -1;
            baos = new ByteArrayOutputStream();
            while ((readSize = fis.read(buffer)) != -1) {
                baos.write(buffer, 0, readSize);
            }
            return baos.toByteArray();
        } catch (Exception e) {
            // TODO: handle exception
            e.printStackTrace();
        } finally {
            try {
                if (fis != null) {
                    fis.close();
                }
                if (baos != null)
                    baos.close();
            } catch (IOException e) {
                // TODO Auto-generated catch block
                e.printStackTrace();
            }
        }
        return data;

    }

    /**
     * @param data     ����
     * @param path     ·��
     * @param fileName �ļ���
     * @return true�ɹ� falseʧ��
     */
    public static boolean writeToSdcard(byte[] data, String path, String fileName) {
        FileOutputStream fos = null;
        try {
            //�ж���û���ļ���
            File filePath = new File(path);
            if (!filePath.exists()) {
                //�����ļ���
                filePath.mkdirs();
            }

            //�ж���û��ͬ�����ļ�
            File file = new File(path + fileName);
            //�еĻ���ɾ��
            if (file.exists()) {
                file.delete();
            }
            //д�ļ�
            fos = new FileOutputStream(file);
            fos.write(data);
            fos.flush();
            return true;
            //		}

        } catch (Exception e) {
            return false;
            // TODO: handle exception
        } finally {
            try {
                if (fos != null)
                    fos.close();
            } catch (IOException e) {
                // TODO Auto-generated catch block
                e.printStackTrace();
            }
        }
    }


}
