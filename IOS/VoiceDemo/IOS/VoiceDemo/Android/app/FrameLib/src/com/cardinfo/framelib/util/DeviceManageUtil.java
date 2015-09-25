package com.cardinfo.framelib.util;

import java.lang.reflect.Method;

import android.app.Activity;
import android.content.Context;
import android.content.pm.FeatureInfo;
import android.content.pm.PackageManager;
import android.hardware.Camera;
import android.os.IBinder;
import android.os.IHardwareService;
import android.view.inputmethod.InputMethodManager;


public class DeviceManageUtil {
	 /** 
     * 设置闪光灯的开启和关闭 
     * @param isEnable 
     * @author linc 
     * @date 2012-3-18 
     */  
    @SuppressWarnings("unused")
	public static void setFlashlightEnabled(boolean isEnable)  
    {  
        try  
        {  
            Method method = Class.forName("android.os.ServiceManager").getMethod("getService", String.class);  
            IBinder binder = (IBinder) method.invoke(null, new Object[] { "hardware" });  
              
            IHardwareService localhardwareservice = IHardwareService.Stub.asInterface(binder);  
            localhardwareservice.setFlashlightEnabled(isEnable);  
        }  
        catch (Exception e)  
        {  
            e.printStackTrace();  
        }  
          
    }  
    
    
   static public void hideInput(Context context){
    	((InputMethodManager)context.getSystemService((context).INPUT_METHOD_SERVICE)).hideSoftInputFromWindow(((Activity)context).getCurrentFocus().getWindowToken(), InputMethodManager.HIDE_NOT_ALWAYS); 
    }
    
}
