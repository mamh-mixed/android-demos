package com.cardinfolink.yunshouyin.data;

import android.app.Activity;
import android.content.Context;
import android.content.SharedPreferences;

public class SaveData {

private static User user=new User();
private static boolean isLogin=false;
private static boolean isAutoLogin=false;


public static boolean isAutoLogin(Context context) {
	SharedPreferences mySharedPreferences = context.getSharedPreferences("savedata", 
			Activity.MODE_PRIVATE);
	isAutoLogin=mySharedPreferences.getBoolean("auto_login", false);
	return isAutoLogin;
}

public static void setAutoLogin(Context context,boolean isAutoLogin) {
	SaveData.isAutoLogin = isAutoLogin;
	SharedPreferences mySharedPreferences = context.getSharedPreferences("savedata", 
			Activity.MODE_PRIVATE); 
			SharedPreferences.Editor editor = mySharedPreferences.edit(); 
			editor.putBoolean("auto_login",isAutoLogin); 
			editor.commit(); 
}

public static boolean isLogin(Context context) {
	SharedPreferences mySharedPreferences = context.getSharedPreferences("savedata", 
			Activity.MODE_PRIVATE); 
	isLogin=mySharedPreferences.getBoolean("is_login", false);
	return isLogin;
}

public static void setLogin(Context context,boolean isLogin) {
	SaveData.isLogin = isLogin;
	SharedPreferences mySharedPreferences = context.getSharedPreferences("savedata", 
			Activity.MODE_PRIVATE); 
			SharedPreferences.Editor editor = mySharedPreferences.edit(); 
			editor.putBoolean("is_login",isLogin); 
			editor.commit(); 
	
}

public static User getUser(Context context) {
	SharedPreferences mySharedPreferences = context.getSharedPreferences("savedata", 
			Activity.MODE_PRIVATE); 
	user.setUsername(mySharedPreferences.getString("username", ""));
	user.setPassword(mySharedPreferences.getString("password", ""));
	user.setClientid(mySharedPreferences.getString("clientid", ""));
	user.setObject_id(mySharedPreferences.getString("objectid", ""));
	System.out.println(user.getUsername());
	return user;
}

public static void setUser(Context context,User user) {
	SaveData.user = user;
	SharedPreferences mySharedPreferences = context.getSharedPreferences("savedata", 
			Activity.MODE_PRIVATE); 
			SharedPreferences.Editor editor = mySharedPreferences.edit(); 
			editor.putString("username", user.getUsername());
			editor.putString("password", user.getPassword());
			editor.putString("clientid", user.getClientid());
			editor.putString("objectid", user.getObject_id());
			editor.commit(); 
}



}
