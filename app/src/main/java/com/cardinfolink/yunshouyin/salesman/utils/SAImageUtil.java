package com.cardinfolink.yunshouyin.salesman.utils;

import android.content.ContentValues;
import android.content.Context;
import android.graphics.Bitmap;
import android.os.Environment;
import android.provider.MediaStore;
import android.util.Log;

import java.io.File;
import java.io.FileOutputStream;
import java.util.UUID;

public class SAImageUtil {
    public static final String LOG_TAG = "picture";

    // Ref: http://stackoverflow.com/questions/3401579/get-filename-and-path-from-uri-from-mediastore
//    public static String getRealPathFromURI(Context context, Uri contentUri) {
//        Cursor cursor = null;
//        try {
//            String[] proj = {MediaStore.Images.Media.DATA};
//            cursor = context.getContentResolver().query(contentUri, proj, null, null, null);
//            int column_index = cursor.getColumnIndexOrThrow(MediaStore.Images.Media.DATA);
//            cursor.moveToFirst();
//            return cursor.getString(column_index);
//        } finally {
//            if (cursor != null) {
//                cursor.close();
//            }
//        }
//    }

    /**
     * 如果没有外部存储卡,就提示用户安装后再调用
     *
     * @return
     */
    private boolean isExternalStorageWritable() {
        String state = Environment.getExternalStorageState();
        if (Environment.MEDIA_MOUNTED.equals(state)) {
            Log.d(LOG_TAG, "external storage is writable");
            return true;
        }
        return false;
    }

    private boolean isExternalStorageReadable() {
        String state = Environment.getExternalStorageState();
        if (Environment.MEDIA_MOUNTED.equals(state) || Environment.MEDIA_MOUNTED_READ_ONLY.equals(state)) {
            Log.d(LOG_TAG, "external storage is readable");
            return true;
        }
        return false;
    }

    public void readSample() {
        String filename = "myfile";
        String string = "Hello World!";
        FileOutputStream outputStream;

        try {
            outputStream = SAApplication.getInstance().getContext().openFileOutput(filename, Context.MODE_PRIVATE);
            outputStream.write(string.getBytes());
            outputStream.close();
        } catch (Exception ex) {
            ex.printStackTrace();
        }
    }

    /**
     * create a public album
     *
     * @param albumName
     * @return
     */
    private File getAlbumStorageDir(String albumName) {
        File file = new File(Environment.getExternalStoragePublicDirectory(Environment.DIRECTORY_PICTURES), albumName);
        if (!file.mkdirs()) {
            Log.e(LOG_TAG, "Directory not created");
        }
        return file;
    }

    public File getAlbumStorageDir(Context context, String albumName) {
        // Get the directory for the app's private pictures directory.
        File file = new File(context.getExternalFilesDir(
                Environment.DIRECTORY_PICTURES), albumName);
        if (!file.mkdirs()) {
            Log.e(LOG_TAG, "Directory not created");
        }
        return file;
    }

    /**
     * 存在外部存储里,图库里是不显示的
     */
    public void saveImageToExternalStorage(Bitmap bitmap) {
        //check if SD card available
        if (!isExternalStorageWritable()) {
            Log.e(LOG_TAG, "请配置外部存储");
        }

        //get app's album folder
        File albumDir = getAlbumStorageDir("云收银销售");

        //generate file name
        String filename = UUID.randomUUID().toString().replace("-", "") + ".jpg";
        File file = new File(albumDir, filename);

        //bitmap to png
        try {
            FileOutputStream outputStream = new FileOutputStream(file);
            bitmap.compress(Bitmap.CompressFormat.JPEG, 100, outputStream);
            outputStream.flush();
            outputStream.close();
            Log.d(LOG_TAG, file.toString());
        } catch (Exception ex) {
            Log.e(LOG_TAG, ex.getMessage());
        }

        //show in gallery
        saveImageToSystemGallery(file.getAbsolutePath());
    }

    /**
     * 所有应用都是直接取gallery里的照片的
     * meta data setup for gallery
     * ref: http://stackoverflow.com/questions/20859584/how-save-image-in-android-gallery
     */
    private void saveImageToSystemGallery(final String filePath) {
        ContentValues values = new ContentValues();

        values.put(MediaStore.Images.Media.DATE_TAKEN, System.currentTimeMillis());
        values.put(MediaStore.Images.Media.MIME_TYPE, "image/jpeg");
        values.put(MediaStore.MediaColumns.DATA, filePath);

        SAApplication.getInstance().getContext().getContentResolver().insert(MediaStore.Images.Media.EXTERNAL_CONTENT_URI, values);
    }
}
