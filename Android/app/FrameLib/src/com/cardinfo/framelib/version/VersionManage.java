package com.cardinfo.framelib.version;

import java.io.BufferedInputStream;
import java.io.File;
import java.io.FileOutputStream;
import java.io.InputStream;
import java.net.HttpURLConnection;
import java.net.URL;

import org.xmlpull.v1.XmlPullParser;

import com.cardinfo.framelib.R;
import com.cardinfo.framelib.constant.Msg;

import android.app.Activity;
import android.app.AlertDialog;
import android.app.AlertDialog.Builder;
import android.app.ProgressDialog;
import android.content.Context;
import android.content.DialogInterface;
import android.content.DialogInterface.OnClickListener;
import android.content.Intent;
import android.content.pm.PackageInfo;
import android.content.pm.PackageManager;
import android.net.Uri;
import android.os.Environment;
import android.os.Handler;
import android.os.Message;
import android.util.Log;
import android.util.Xml;

public class VersionManage {
	private static String TAG="VersionManage";
	private Context mContext;
	private Handler mhandler;
	private UpdateInfo info;
	private String mVersionServer;
	public VersionManage(Context context,Handler handler,String versionServer) {
		mContext=context;
		mhandler=handler;
		mVersionServer=versionServer;
	}

	/*
	 * 获取当前程序的版本号
	 */
	private String getVersionName() throws Exception {
		// 获取packagemanager的实例
		PackageManager packageManager = mContext.getPackageManager();
		// getPackageName()是你当前类的包名，0代表是获取版本信息
		PackageInfo packInfo = packageManager.getPackageInfo( mContext.getPackageName(),
				0);
		return packInfo.versionName;
	}
	
	
	/* 
	 * 用pull解析器解析服务器返回的xml文件 (xml封装了版本号) 
	 */  
	public static UpdateInfo getUpdateInfo(InputStream is) throws Exception{  
	    XmlPullParser  parser = Xml.newPullParser();    
	    parser.setInput(is, "utf-8");//设置解析的数据源    
	   int type = parser.getEventType();  
	    UpdateInfo info = new UpdateInfo();//实体   
	    while(type != XmlPullParser.END_DOCUMENT ){  
	        switch (type) {  
	        case XmlPullParser.START_TAG:  
	           if("version".equals(parser.getName())){  
	                info.setVersion(parser.nextText()); //获取版本号   
	            }else if ("url".equals(parser.getName())){  
	                info.setUrl(parser.nextText()); //获取要升级的APK文件   
	           }else if ("description".equals(parser.getName())){  
	                info.setDescription(parser.nextText()); //获取该文件的信息   
	            }  
	            break;  
	        }  
	        type = parser.next();  
	    }  
	    return info;  
	}
	
	
	
	
	public  File getFileFromServer(String path, ProgressDialog pd) throws Exception{  
		    //如果相等的话表示当前的sdcard挂载在手机上并且是可用的   
		   if(Environment.getExternalStorageState().equals(Environment.MEDIA_MOUNTED)){  
		        URL url = new URL(path);  
		       HttpURLConnection conn =  (HttpURLConnection) url.openConnection();  
		       conn.setConnectTimeout(5000);  
		        //获取到文件的大小    
		        pd.setMax(conn.getContentLength());  
		       InputStream is = conn.getInputStream();  
		        File file = new File(Environment.getExternalStorageDirectory(), "updata.apk");  
		        FileOutputStream fos = new FileOutputStream(file);  
		        BufferedInputStream bis = new BufferedInputStream(is);  
		        byte[] buffer = new byte[1024];  
		        int len ;  
		        int total=0;  
		        while((len =bis.read(buffer))!=-1){  
		            fos.write(buffer, 0, len);  
		           total+= len;  
		           //获取当前下载量   
		            pd.setProgress(total);  
		        }  
		        fos.close();  
		        bis.close();  
		        is.close();  
		        return file;  
		    }  
		    else{  
		        return null;  
		    }  
		}
	
	
	public void update(){
		new Thread(new CheckVersionTask()).start();
	}
	
	
	
	/*
	 * 从服务器获取xml解析并进行比对版本号 
	 */ 
	public class CheckVersionTask implements Runnable{ 
	   
	    public void run() { 
	        try { 
	            //从资源文件获取服务器 地址   
	            String path = mVersionServer;	            
	            //包装成url的对象   
	            URL url = new URL(path); 
	            HttpURLConnection conn =  (HttpURLConnection) url.openConnection();  
	            conn.setConnectTimeout(5000); 
	            InputStream is =conn.getInputStream();  
	            info = getUpdateInfo(is); 
	               
	            if(info.getVersion().equals(getVersionName())){ 
	                Log.i(TAG,"版本号相同无需升级");
	                mhandler.sendEmptyMessageDelayed(Msg.MSG_FROM_CLIENT_JUMP_LOGIN,1000);
	               
	            }else{ 
	                Log.i(TAG,"版本号不同 ,提示用户升级 "); 
	                mhandler.sendEmptyMessage(Msg.MSG_FROM_CLIENT_REMIND_UPDATE);
	            } 
	        } catch (Exception e) { 
	            // 待处理   
	        	mhandler.sendEmptyMessageDelayed(Msg.MSG_FROM_CLIENT_JUMP_LOGIN,1000);
//	            msg.what = GET_UNDATAINFO_ERROR; 
//	            handler.sendMessage(msg); 
//	            e.printStackTrace(); 
	        }  
	    } 
	} 
	
	
	/*
	 * 
	 * 弹出对话框通知用户更新程序 
	 * 
	 * 弹出对话框的步骤：
	 *  1.创建alertDialog的builder.  
	 *  2.要给builder设置属性, 对话框的内容,样式,按钮
	 *  3.通过builder 创建一个对话框
	 *  4.对话框show()出来  
	 */ 
	public void showUpdataDialog() { 
	    AlertDialog.Builder builer = new Builder(mContext) ;  
	    builer.setTitle("版本升级"); 
	    builer.setMessage(info.getDescription()); 
	    //当点确定按钮时从服务器上下载 新的apk 然后安装   
	    builer.setPositiveButton("确定", new OnClickListener() { 
	    public void onClick(DialogInterface dialog, int which) { 
	            Log.i(TAG,"下载apk,更新"); 
	            downLoadApk(); 
	        }    
	    }); 
	    //当点取消按钮时进行登录  
	    builer.setNegativeButton("取消", new OnClickListener() { 
	        public void onClick(DialogInterface dialog, int which) { 
	        	 mhandler.sendEmptyMessage(Msg.MSG_FROM_CLIENT_JUMP_LOGIN);
	        } 
	    }); 
	    AlertDialog dialog = builer.create(); 
	    dialog.setCancelable(false);
	    dialog.show(); 
	}
	
	
	/*
	 * 从服务器中下载APK
	 */ 
	public void downLoadApk() { 
	    final ProgressDialog pd;    //进度条对话框  
	    pd = new  ProgressDialog(mContext); 
	    pd.setProgressStyle(ProgressDialog.STYLE_HORIZONTAL); 
	    pd.setMessage("正在下载更新"); 
	    pd.show(); 
	    new Thread(){ 
	        @Override 
	        public void run() { 
	            try { 
	                File file = getFileFromServer(info.getUrl(), pd); 
	                sleep(3000); 
	                installApk(file); 
	                pd.dismiss(); //结束掉进度条对话框  
	            } catch (Exception e) { 
	               
	                e.printStackTrace(); 
	            } 
	        }}.start(); 
	}
	
	//安装apk   
	protected void installApk(File file) { 
		Intent intent = new Intent(Intent.ACTION_VIEW); intent.setFlags(Intent.FLAG_ACTIVITY_NEW_TASK); 
		intent.setDataAndType(Uri.fromFile(file),"application/vnd.android.package-archive");  
		mContext.startActivity(intent);  
		android.os.Process.killProcess(android.os.Process.myPid());
	}
	
}
