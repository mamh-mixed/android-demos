package com.cardinfolink.yunshouyin.activity;

import java.util.ArrayList;

import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.constant.SystemConfig;
import com.cardinfolink.yunshouyin.data.SessonData;
import com.cardinfolink.yunshouyin.view.AccountUpdateView;
import com.cardinfolink.yunshouyin.view.LimitIncreaseView;
import com.cardinfolink.yunshouyin.view.PasswordUpdateView;
import com.cardinfolink.yunshouyin.view.ScanCodeView;
import com.cardinfolink.yunshouyin.view.TransManageView;
import com.jeremyfeinstein.slidingmenu.lib.SlidingMenu;
import com.umeng.update.UmengUpdateAgent;
import android.content.Intent;
import android.net.Uri;
import android.os.Bundle;
import android.util.Log;
import android.view.View;
import android.widget.AdapterView;
import android.widget.AdapterView.OnItemClickListener;
import android.widget.ArrayAdapter;
import android.widget.LinearLayout;
import android.widget.ListView;
import android.widget.LinearLayout.LayoutParams;

public class MainActivity extends BaseActivity {
	SlidingMenu mLeftMenu;
	private ScanCodeView mScanCodeView;
	private TransManageView mTransManageView;
	private PasswordUpdateView mPasswordUpdateView;
	private AccountUpdateView mAccountUpdateView;
	private LimitIncreaseView mLimitIncreaseView;
	private LinearLayout mMainContent;
	private ListView mDrawerList;
	private ArrayList<String> menuLists;
	private ArrayAdapter<String> adapter;

	@Override
	protected void onCreate(Bundle savedInstanceState) {
		super.onCreate(savedInstanceState);
		setContentView(R.layout.main_activity);
		initLayout();
		UmengUpdateAgent.setUpdateOnlyWifi(true);
		UmengUpdateAgent.setUpdateCheckConfig(false);
		UmengUpdateAgent.update(this);	
	}

	private void initLayout() {
		mMainContent = (LinearLayout) findViewById(R.id.main_content);
		
		LinearLayout.LayoutParams layoutParams = new LayoutParams(
				LayoutParams.MATCH_PARENT, LayoutParams.MATCH_PARENT);
		mScanCodeView = new ScanCodeView(mContext);
		mScanCodeView.setLayoutParams(layoutParams);
		
		mTransManageView=new TransManageView(mContext);
		mTransManageView.setLayoutParams(new LayoutParams(
				LayoutParams.MATCH_PARENT, LayoutParams.MATCH_PARENT));
		
		mPasswordUpdateView=new PasswordUpdateView(mContext);
		mPasswordUpdateView.setLayoutParams(new LayoutParams(
				LayoutParams.MATCH_PARENT, LayoutParams.MATCH_PARENT));
		
		mAccountUpdateView=new AccountUpdateView(mContext);
		mAccountUpdateView.setLayoutParams(new LayoutParams(
				LayoutParams.MATCH_PARENT, LayoutParams.MATCH_PARENT));
		
		mLimitIncreaseView=new LimitIncreaseView(mContext);
		mLimitIncreaseView.setLayoutParams(new LayoutParams(
				LayoutParams.MATCH_PARENT, LayoutParams.MATCH_PARENT));
		
		mMainContent.addView(mScanCodeView);

		mLeftMenu = new SlidingMenu(this);
		mLeftMenu.setMode(SlidingMenu.LEFT);
		mLeftMenu.setTouchModeAbove(SlidingMenu.TOUCHMODE_FULLSCREEN);
		mLeftMenu.setShadowWidthRes(R.dimen.shadow_width);
		mLeftMenu.setShadowDrawable(R.drawable.shadow);

		// 设置滑动菜单视图的宽度
		mLeftMenu.setBehindOffsetRes(R.dimen.slidingmenu_offset);
		// 设置渐入渐出效果的值
		mLeftMenu.setFadeDegree(0.35f);
		/**
		 * SLIDING_WINDOW will include the Title/ActionBar in the content
		 * section of the SlidingMenu, while SLIDING_CONTENT does not.
		 */
		mLeftMenu.attachToActivity(this, SlidingMenu.SLIDING_CONTENT);
		// 为侧滑菜单设置布局
		mLeftMenu.setMenu(R.layout.leftmenu);
		
		 mDrawerList=(ListView) mLeftMenu.findViewById(R.id.left_drawer);
		 menuLists= new ArrayList<String>();
	        
     	menuLists.add("扫码支付");
     	menuLists.add("交易管理");
     	menuLists.add("密码修改");
 	    menuLists.add("账户修改");
     	menuLists.add("限额提升");
     	menuLists.add("网页账单");
 	    menuLists.add("安全退出");

     
     adapter=new ArrayAdapter<String>(this,R.layout.menu_list_item, menuLists);
     mDrawerList.setAdapter(adapter);
     mDrawerList.setOnItemClickListener(new MenuOnItemClick());
		
	}
	
	
	
	

	public void BtnMenuOnClick(View view) {

		if (mLeftMenu.isMenuShowing()) {
			mLeftMenu.toggle();
		} else {
			mLeftMenu.showMenu();
		}

	}
	
	 	
	private class MenuOnItemClick implements OnItemClickListener{

		@Override
		public void onItemClick(AdapterView<?> parent, View view, int position,
				long id) {
			SessonData.position_view=position;
			switch(position)
			{
				case 0:
					mMainContent.removeAllViews();
					mMainContent.addView(mScanCodeView);
					mLeftMenu.toggle();
					break;
				case 1:				
					mMainContent.removeAllViews();
					mMainContent.addView(mTransManageView);
					mTransManageView.initData();
					mLeftMenu.toggle();
					break;
				case 2:
					mMainContent.removeAllViews();
					mMainContent.addView(mPasswordUpdateView);
					mLeftMenu.toggle();
					break;
				case 3:
					mMainContent.removeAllViews();
					mMainContent.addView(mAccountUpdateView);
					mAccountUpdateView.initData();
					mLeftMenu.toggle();
					break;
				case 4:
					mMainContent.removeAllViews();
					mMainContent.addView(mLimitIncreaseView);
					mLeftMenu.toggle();
					break;
				case 5:
					 mLeftMenu.toggle();
					 Uri uri;
					 uri = Uri.parse(SystemConfig.WEB_BILL_URL+"?merchantCode="+SessonData.loginUser.getObject_id());
					 Intent  intent = new  Intent(Intent.ACTION_VIEW, uri);
					 startActivity(intent);	
					 
					 break;
				case 6:
					finish();
					break;
				
				
				
			}
			
		}
		
	}
	
	private void openView(int position){
		SessonData.position_view=position;
		switch(position)
		{
			case 0:
				mMainContent.removeAllViews();
				mMainContent.addView(mScanCodeView);
				
				break;
			case 1:				
				mMainContent.removeAllViews();
				mMainContent.addView(mTransManageView);
				mTransManageView.initData();
				break;
			case 2:
				mMainContent.removeAllViews();
				mMainContent.addView(mPasswordUpdateView);
				
				break;
			case 3:
				mMainContent.removeAllViews();
				mMainContent.addView(mAccountUpdateView);
				mAccountUpdateView.initData();
				
				break;
			case 4:
				mMainContent.removeAllViews();
				mMainContent.addView(mLimitIncreaseView);
				
				break;
	
		
			
		}
	}
	
	@Override
	protected void onResume() {		
		super.onResume();

		Log.i("opp", "position="+SessonData.position_view);
		openView(SessonData.position_view);
	}
	
}
