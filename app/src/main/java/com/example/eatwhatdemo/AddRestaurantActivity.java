package com.example.eatwhatdemo;

import java.util.ArrayList;
import java.util.List;

import android.app.Activity;
import android.content.ContentValues;
import android.database.Cursor;
import android.database.sqlite.SQLiteDatabase;
import android.os.Bundle;
import android.view.Menu;
import android.view.View;
import android.view.View.OnClickListener;
import android.widget.Button;
import android.widget.EditText;
import android.widget.ExpandableListView;
import android.widget.ExpandableListView.OnChildClickListener;
import android.widget.ExpandableListView.OnGroupClickListener;

import com.example.library.DatabaseHelper;
import com.example.adapter.MyExpandableListAdapter;

public class AddRestaurantActivity extends Activity implements OnClickListener,
		OnChildClickListener, OnGroupClickListener {
	private List<List<String>> groupData;
	private List<List<String>> childrenData;
	private MyExpandableListAdapter myExpandableListAdapter;
	private ExpandableListView expandableListView;
	private DatabaseHelper dbhelper;

	private Button commitButton;
	private Button backButton;

	private EditText nameeditText;
	private EditText addreditText;
	private EditText phoneeditText;
	private EditText descriptioneditText;

	private String restaurantName;
	private String restaurantAddress;
	private String restaurantPhone;
	private String restaurantDescription;

	@Override
	protected void onCreate(Bundle savedInstanceState) {
		super.onCreate(savedInstanceState);
		setContentView(R.layout.activity_addrestaurant);

		// init the adapter data, groupData and childrenData.
		loadData();

		commitButton = (Button) findViewById(R.id.addrestaurantcommitbutton);
		backButton = (Button) findViewById(R.id.addrestaurantbackbutton);
		nameeditText = (EditText) findViewById(R.id.add_restaurant_name_editText);
		addreditText = (EditText) findViewById(R.id.add_restaurant_addr_editText);
		phoneeditText = (EditText) findViewById(R.id.add_restaurant_phone_editText);
		descriptioneditText = (EditText) findViewById(R.id.add_restaurant_description_editText);

		expandableListView = (ExpandableListView) findViewById(R.id.add_restaurant_expandablelistview);
		myExpandableListAdapter = new MyExpandableListAdapter(
				this.getApplicationContext(), groupData, childrenData);
		expandableListView.setAdapter(myExpandableListAdapter);

		commitButton.setOnClickListener(this);
		backButton.setOnClickListener(this);

		expandableListView.setOnChildClickListener(this);

		expandableListView.setOnGroupClickListener(this);
	}

	@Override
	public boolean onCreateOptionsMenu(Menu menu) {
		return true;
	}

	/*
	 * init grouData, childrenData. the datas from db.
	 */
	private void loadData() {
		groupData = new ArrayList<List<String>>();
		childrenData = new ArrayList<List<String>>();

		dbhelper = new DatabaseHelper(this);
		SQLiteDatabase db = dbhelper.getReadableDatabase();

		Cursor cursor = db.rawQuery("select * from ew_restaurant", null);
		while (cursor.moveToNext()) {
			int id = cursor.getInt(0);
			String name = cursor.getString(1);
			String address = cursor.getString(2);
			String phone = cursor.getString(3);
			String description = cursor.getString(4);

			List<String> subgrouplist = new ArrayList<String>();
			subgrouplist.add("" + id);
			subgrouplist.add(name);
			groupData.add(subgrouplist);

			List<String> subchildlist = new ArrayList<String>();
			subchildlist.add(address);
			subchildlist.add(phone);
			subchildlist.add(description);
			childrenData.add(subchildlist);
		}
		cursor.close();
		db.close();

	}

	@Override
	public void onClick(View v) {
		// TODO Auto-generated method stub
		if (commitButton == (Button) v) {
			SQLiteDatabase db = dbhelper.getReadableDatabase();

			restaurantName = nameeditText.getText().toString();
			if (restaurantName.length() == 0) {
				restaurantName = "name";
			}
			restaurantAddress = addreditText.getText().toString();
			if (restaurantAddress.length() == 0) {
				restaurantAddress = "address";
			}
			restaurantPhone = phoneeditText.getText().toString();
			if (restaurantPhone.length() == 0) {
				restaurantPhone = "13636443123";
			}
			restaurantDescription = descriptioneditText.getText().toString();
			if (restaurantDescription.length() == 0) {
				restaurantDescription = "miaoshu";
			}

			ContentValues values = new ContentValues();
			values.put("restaurant_name", restaurantName);
			values.put("restaurant_address", restaurantAddress);
			values.put("restaurant_phone", restaurantPhone);
			values.put("restaurant_description", restaurantDescription);

			db.insert("ew_restaurant", null, values);
			db.close();
			finish();
		} else if (backButton == (Button) v) {
			finish();
		}

	}

	@Override
	public boolean onGroupClick(ExpandableListView parent, View v,
			int groupPosition, long id) {
		// TODO Auto-generated method stub
		return false;
	}

	@Override
	public boolean onChildClick(ExpandableListView parent, View v,
			int groupPosition, int childPosition, long id) {
		// TODO Auto-generated method stub
		return false;
	}
}
