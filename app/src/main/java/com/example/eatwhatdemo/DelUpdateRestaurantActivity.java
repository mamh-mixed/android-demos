package com.example.eatwhatdemo;

import java.util.ArrayList;
import java.util.List;

import android.app.Activity;
import android.app.AlertDialog;
import android.app.AlertDialog.Builder;
import android.content.DialogInterface;
import android.content.DialogInterface.OnClickListener;
import android.database.Cursor;
import android.database.sqlite.SQLiteDatabase;
import android.os.Bundle;
import android.view.LayoutInflater;
import android.view.Menu;
import android.view.View;
import android.widget.AdapterView;
import android.widget.AdapterView.OnItemLongClickListener;
import android.widget.EditText;
import android.widget.ExpandableListView;
import android.widget.ExpandableListView.OnGroupClickListener;
import android.widget.TextView;

import com.example.library.DatabaseHelper;
import com.example.adapter.MyExpandableListAdapter;

public class DelUpdateRestaurantActivity extends Activity {
    private List<List<String>> groupData;
    private List<List<String>> childrenData;
    private DatabaseHelper dbhelper;
    private MyExpandableListAdapter myExpandableListAdapter;
    private ExpandableListView expandableListView;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);

        setContentView(R.layout.activity_delupdaterestaurant);

        // ��ʼ�� groupData��childrenData���顣
        loadData();

        expandableListView = (ExpandableListView) findViewById(R.id.del_update_restaurant_expandablelistview);

        myExpandableListAdapter = new MyExpandableListAdapter(this.getApplicationContext(), groupData, childrenData);

        expandableListView.setAdapter(myExpandableListAdapter);

		/*
         * ��expandablelistview�����鱻����ļ�����
		 */
        expandableListView.setOnGroupClickListener(new OnGroupClickListener() {
            @Override
            public boolean onGroupClick(ExpandableListView parent, View v, int groupPosition, long id) {

                return false;
            }
        });

		/*
		 * add lintener for expandablelistview
		 */
        expandableListView.setOnItemLongClickListener(new OnItemLongClickListener() {
            private View convertView;
            private int posi;

            public boolean onItemLongClick(AdapterView<?> parent, View view, int position, long id) {
                convertView = view;
                posi = position;

                AlertDialog.Builder builder = new Builder(DelUpdateRestaurantActivity.this);
                builder.setTitle(R.string.delete_or_edit_dialog_title);

                builder.setItems(R.array.delete_or_edit_menu, new DialogInterface.OnClickListener() {
                    public void onClick(DialogInterface dialog, int which) {

                        switch (which) {
						/*
						 * delete operation.
						 */
                            case 0: {
                                TextView idtextView = (TextView) convertView.findViewById(R.id.expandablelistview_group_id_TextView);

                                DatabaseHelper dbhelper = new DatabaseHelper(DelUpdateRestaurantActivity.this);
                                SQLiteDatabase db = dbhelper.getReadableDatabase();

                                db.beginTransaction();
                                try {
                                    db.execSQL("delete from ew_menu where restaurant_id = \"" + idtextView.getText().toString() + "\"");
                                    db.execSQL("delete from ew_restaurant where restaurant_id = \"" + idtextView.getText().toString() + "\"");
                                    db.setTransactionSuccessful();
                                } catch (Exception e) {
                                    e.printStackTrace();
                                } finally {
                                    db.endTransaction();
                                }

                                db.close();

                                groupData.remove(posi);
                                childrenData.remove(posi);
                                myExpandableListAdapter.notifyDataSetChanged();
                                expandableListView.invalidate();

                                break;
                            }
						/*
						 * edit operation.
						 */
                            case 1: {
                                LayoutInflater infalter = getLayoutInflater();
                                View layout = infalter.inflate(R.layout.edit_restaurant_dialog, null);
                                final TextView idtextView = (TextView) convertView.findViewById(R.id.expandablelistview_group_id_TextView);

                                final EditText nameEditText = (EditText) layout.findViewById(R.id.edit_restaurant_name_editText);
                                final EditText addressEditText = (EditText) layout.findViewById(R.id.edit_restaurant_addr_editText);
                                final EditText phoneEditText = (EditText) layout.findViewById(R.id.edit_restaurant_phone_editText);
                                final EditText descriptionEditText = (EditText) layout.findViewById(R.id.edit_restaurant_description_editText);

                                DatabaseHelper dbhelper = new DatabaseHelper(DelUpdateRestaurantActivity.this);
                                SQLiteDatabase db = dbhelper.getReadableDatabase();

                                String name = "";
                                String address = "";
                                String phone = "";
                                String description = "";
                                Cursor cursor = db.rawQuery("select * from ew_restaurant where restaurant_id='" + idtextView.getText().toString() + "'", null);
                                while (cursor.moveToNext()) {
                                    int id = cursor.getInt(0);
                                    name = cursor.getString(1);
                                    address = cursor.getString(2);
                                    phone = cursor.getString(3);
                                    description = cursor.getString(4);
                                }
                                db.close();

                                nameEditText.setText(name);
                                addressEditText.setText(address);
                                phoneEditText.setText(phone);
                                descriptionEditText.setText(description);

                                AlertDialog.Builder builder = new Builder(DelUpdateRestaurantActivity.this);
                                builder.setTitle(R.string.edit_restaurant_dialog_title);

                                builder.setView(layout);
                                builder.setPositiveButton("Ok", new OnClickListener() {
                                    @Override
                                    public void onClick(DialogInterface dialog, int which) {
                                        DatabaseHelper dbhelper = new DatabaseHelper(DelUpdateRestaurantActivity.this);
                                        SQLiteDatabase db = dbhelper.getReadableDatabase();
                                        db.beginTransaction();
                                        try {
                                            db.execSQL("update ew_restaurant set restaurant_name='" + nameEditText.getText().toString() + "' where restaurant_id = \"" + idtextView.getText().toString() + "\"");
                                            db.execSQL("update ew_restaurant set restaurant_address='" + addressEditText.getText().toString() + "' where restaurant_id = \"" + idtextView.getText().toString() + "\"");
                                            db.execSQL("update ew_restaurant set restaurant_phone='" + phoneEditText.getText().toString() + "' where restaurant_id = \"" + idtextView.getText().toString() + "\"");
                                            db.execSQL("update ew_restaurant set restaurant_description='" + descriptionEditText.getText().toString() + "' where restaurant_id = \"" + idtextView.getText().toString() + "\"");
                                            db.setTransactionSuccessful();
                                        } catch (Exception e) {
                                            e.printStackTrace();
                                        } finally {
                                            db.endTransaction();
                                        }

                                        db.close();

                                        myExpandableListAdapter.notifyDataSetChanged();
                                        expandableListView.invalidate();
                                    }
                                });
                                builder.setNegativeButton("Cancel", null);
                                builder.show();
                                break;
                            }
                            default:
                                break;
                        }// end switch()
                    }
                });

                builder.setNegativeButton("ȡ��", new OnClickListener() {
                    public void onClick(DialogInterface dialog, int which) {

                    }
                });

                builder.show();

                return true;
            }
        });

    }

    @Override
    public boolean onCreateOptionsMenu(Menu menu) {
        return true;
    }

    /*
     * init groupData��childrenData the datas from the database by select
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
            subchildlist.add(getResources().getString(R.string.address) + address);
            subchildlist.add(getResources().getString(R.string.phone) + phone);
            subchildlist.add(getResources().getString(R.string.description) + description);
            childrenData.add(subchildlist);
        }
        cursor.close();
        db.close();

    }

}
