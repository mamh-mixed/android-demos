package com.cardinfo.framelib.view;

import com.cardinfo.framelib.R;
import com.cardinfo.framelib.constant.Msg;
import android.content.Context;
import android.os.Message;
import android.view.View;
import android.view.animation.Animation;
import android.view.animation.AnimationUtils;
import android.widget.ImageView;
import android.widget.TextView;

public class Loading_Dialog {
	
	private View dialogView;
	private ImageView loadImg;
	private Context mContext;
	private boolean isLoading=false;
	public Loading_Dialog(Context context,View view){
		dialogView=view;
		mContext=context;
		loadImg=(ImageView) dialogView.findViewById(R.id.load_img);
	}
	
	
	public void startLoading(){
		   dialogView.setVisibility(View.VISIBLE);
		   Animation loadingAnimation = AnimationUtils.loadAnimation(  
				   mContext, R.anim.loading_animation);		  
		   loadImg.startAnimation(loadingAnimation);		  
	
		 
	}
	
	public void endLoading(){
		 dialogView.setVisibility(View.GONE);
		 loadImg.clearAnimation();
	}

}
