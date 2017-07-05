package com.cardinfolink.yunshouyin.data;

import android.app.Activity;
import android.content.Context;
import android.content.SharedPreferences;

public class SaveData {

    private static User user = new User();


    public static User getUser(Context context) {
        SharedPreferences mySharedPreferences = context.getSharedPreferences("savedata",
                Activity.MODE_PRIVATE);
        user.setUsername(mySharedPreferences.getString("username", ""));
        user.setPassword(mySharedPreferences.getString("password", ""));
        user.setAutoLogin(mySharedPreferences.getBoolean("autologin", false));
        return user;
    }

    public static void setUser(Context context, User user) {
        SaveData.user = user;
        SharedPreferences mySharedPreferences = context.getSharedPreferences("savedata",
                Activity.MODE_PRIVATE);
        SharedPreferences.Editor editor = mySharedPreferences.edit();
        editor.putString("username", user.getUsername());
        editor.putString("password", user.getPassword());
        editor.putBoolean("autologin", user.isAutoLogin());
        editor.commit();
    }


    public static void setCameraFront(Context context, boolean isfront) {
        SharedPreferences mySharedPreferences = context.getSharedPreferences("savedata",
                Activity.MODE_PRIVATE);
        SharedPreferences.Editor editor = mySharedPreferences.edit();
        editor.putBoolean("isfront", isfront);
        editor.commit();
    }


    public static boolean isCameraFront(Context context) {
        SharedPreferences mySharedPreferences = context.getSharedPreferences("savedata",
                Activity.MODE_PRIVATE);
        return mySharedPreferences.getBoolean("isfront", false);
    }

    public static void setMessageTime(Context context, String lasttime) {
        SharedPreferences mySharedPreferences = context.getSharedPreferences("savedata",
                Activity.MODE_PRIVATE);
        SharedPreferences.Editor editor = mySharedPreferences.edit();
        editor.putString("lasttime", lasttime);
        editor.commit();
    }

    public static String  getMessageTime(Context context) {
        SharedPreferences mySharedPreferences = context.getSharedPreferences("savedata",
                Activity.MODE_PRIVATE);
        return mySharedPreferences.getString("lasttime", "");
    }



    public static void setMessageClick(Context context, boolean isClick) {
        SharedPreferences mySharedPreferences = context.getSharedPreferences("savedata",
                Activity.MODE_PRIVATE);
        SharedPreferences.Editor editor = mySharedPreferences.edit();
        editor.putBoolean("isClick", isClick);
        editor.commit();
    }


    public static boolean  getMessageClick(Context context) {
        SharedPreferences mySharedPreferences = context.getSharedPreferences("savedata",
                Activity.MODE_PRIVATE);
        return mySharedPreferences.getBoolean("isClick", false);
    }


}
