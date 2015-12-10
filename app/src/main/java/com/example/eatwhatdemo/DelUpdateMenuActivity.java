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
import android.util.Log;
import android.view.Menu;
import android.view.View;
import android.widget.AdapterView;
import android.widget.AdapterView.OnItemLongClickListener;
import android.widget.ExpandableListView;
import android.widget.ExpandableListView.OnGroupClickListener;
import android.widget.TextView;

import com.example.library.DatabaseHelper;
import com.example.adapter.MyExpandableListAdapter;

public class DelUpdateMenuActivity extends Activity {
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_delupdatemenu);

        // ��ʼ�� groupData��childrenData���顣
        loadData();

        expandableListView = (ExpandableListView) findViewById(R.id.del_update_menu_expandablelistview);
        Log.d("eat", "ex = " + expandableListView);

        myExpandableListAdapter = new MyExpandableListAdapter(this.getApplicationContext(), groupData, childrenData);
        expandableListView.setAdapter(myExpandableListAdapter);

        expandableListView.setOnGroupClickListener(new OnGroupClickListener() {
            @Override
            public boolean onGroupClick(ExpandableListView parent, View v, int groupPosition, long id) {

                return false;
            }
        });

        expandableListView.setOnItemLongClickListener(new OnItemLongClickListener() {
            private View convertView;
            private int posi;

            public boolean onItemLongClick(AdapterView<?> parent, View view, int position, long id) {
                convertView = view;
                posi = position;

                AlertDialog.Builder builder = new Builder(DelUpdateMenuActivity.this);
                builder.setTitle(R.string.delete_or_edit_dialog_title);
                builder.setItems(R.array.delete_or_edit_menu, new DialogInterface.OnClickListener() {
                    public void onClick(DialogInterface dialog, int which) {
                        switch (which) {
                            // delete����
                            case 0:
                                TextView idtextView = (TextView) convertView.findViewById(R.id.expandablelistview_group_id_TextView);

                                DatabaseHelper dbhelper = new DatabaseHelper(DelUpdateMenuActivity.this);
                                SQLiteDatabase db = dbhelper.getReadableDatabase();

                                db.beginTransaction();
                                try {
                                    db.execSQL("delete from ew_menu where menu_id = '" + idtextView.getText().toString() + "'");
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
                            case 1:
                                break;
                            default:
                                break;
                        }
                    }
                });
                builder.setNegativeButton("Cancel", new OnClickListener() {
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
     * ��ʼ��groupData��childData���顣 ��ʼ����dbhelper
     */
    private void loadData() {
        groupData = new ArrayList<List<String>>();
        childrenData = new ArrayList<List<String>>();

        dbhelper = new DatabaseHelper(this);
        SQLiteDatabase db = dbhelper.getReadableDatabase();

        Cursor cursor = db.rawQuery("select menu_id,menu_name,restaurant_name,menu_description from ew_menu,ew_restaurant where ew_menu.restaurant_id=ew_restaurant.restaurant_id", null);
        while (cursor.moveToNext()) {
            int id = cursor.getInt(0);
            String menuname = cursor.getString(1);
            String restaurantname = cursor.getString(2);
            String description = cursor.getString(3);

            List<String> subgrouplist = new ArrayList<String>();
            subgrouplist.add("" + id);
            subgrouplist.add(menuname);
            groupData.add(subgrouplist);

            List<String> subchildlist = new ArrayList<String>();
            subchildlist.add(getResources().getString(R.string.restaurant) + restaurantname);
            subchildlist.add(getResources().getString(R.string.description) + description);
            childrenData.add(subchildlist);
        }
        cursor.close();
        db.close();

    }

    private List<List<String>> groupData;
    private List<List<String>> childrenData;
    private DatabaseHelper dbhelper;
    private MyExpandableListAdapter myExpandableListAdapter;
    private ExpandableListView expandableListView;
}
