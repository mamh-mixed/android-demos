package com.cardinfolink.yunshouyin.constant;

import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.util.ContextUtil;

public class SystemConfig {

  
	  
  public static final String APP_KEY=ContextUtil.getInstance().getResources().getString(R.string.app_key);//app用户系统交互key
   
   //用户系统服务器地址  
   public static final String Server=ContextUtil.getInstance().getResources().getString(R.string.user_server);
   
  //扫固定码支付网页订单支付
   public static final String WEB_BILL_URL=ContextUtil.getInstance().getResources().getString(R.string.web_bill_url);;
 
  // SDK 环境
   public static final boolean IS_PRODUCE =true;
   
   //公共数据平台秘钥和key
   
   public static final String bankbase_key="20e786206dcf4aae8a63fe34553fd274";
   public static final String bankbase_url="http://211.144.213.120:443/bdp";
   
}
