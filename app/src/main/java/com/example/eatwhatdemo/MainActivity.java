package com.example.eatwhatdemo;

import java.util.ArrayList;
import java.util.List;

import android.annotation.SuppressLint;
import android.app.Activity;
import android.app.AlertDialog.Builder;
import android.content.DialogInterface;
import android.content.DialogInterface.OnClickListener;
import android.content.Intent;
import android.database.Cursor;
import android.database.sqlite.SQLiteDatabase;
import android.hardware.Sensor;
import android.hardware.SensorEvent;
import android.hardware.SensorEventListener;
import android.hardware.SensorManager;
import android.media.MediaPlayer;
import android.os.Bundle;
import android.util.Log;
import android.view.LayoutInflater;
import android.view.View;
import android.widget.AdapterView;
import android.widget.AdapterView.OnItemClickListener;
import android.widget.ExpandableListView;
import android.widget.RadioButton;
import android.widget.RadioGroup;
import android.widget.TabHost;
import android.widget.TabHost.OnTabChangeListener;
import android.widget.TabHost.TabSpec;

import com.example.library.CornerListView;
import com.example.library.DatabaseHelper;
import com.example.adapter.MyExpandableListAdapter;
import com.example.adapter.MyListAdapter;

public class MainActivity extends Activity implements OnTabChangeListener,
        OnItemClickListener {

    private CornerListView add_del_update_serach_cornerListView = null;
    private ExpandableListView resultExpandableListView;
    private TabHost tabHost = null;
    private MediaPlayer mediaPlayer;
    private View dialogLayout;
    private boolean isRunAfterSensorChangedDo = false;

    // ��Ӧ������
    private SensorManager mSensorManager;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);

        // ���MediaPlayerʵ������create���������ģ���ô��һ����������ǰ����Ҫ�ٵ���prepare�����ˣ���Ϊcreate�������Ѿ����ù��ˡ�
        mediaPlayer = MediaPlayer.create(this, R.raw.shake);

        // tab host
        tabHost = (TabHost) findViewById(android.R.id.tabhost);
        tabHost.setup();

        // tab 1
        TabSpec tabSpec1 = tabHost.newTabSpec("todayTab");
        tabSpec1.setContent(R.id.LinearLayout001);
        tabSpec1.setIndicator("today", MainActivity.this.getResources()
                .getDrawable(R.drawable.today));
        tabHost.addTab(tabSpec1);

        // tab 2
        TabSpec tabSpec2 = tabHost.newTabSpec("historyTab");
        tabSpec2.setContent(R.id.LinearLayout002);
        tabSpec2.setIndicator("history", MainActivity.this.getResources()
                .getDrawable(R.drawable.history));
        tabHost.addTab(tabSpec2);

        // tab 3
        TabSpec tabSpec3 = tabHost.newTabSpec("settingTab");
        tabSpec3.setContent(R.id.LinearLayout003);
        tabSpec3.setIndicator("setting", MainActivity.this.getResources()
                .getDrawable(R.drawable.today));
        tabHost.addTab(tabSpec3);

        // ҡһҡ֮��� ��ʾ�Ľ��
        resultExpandableListView = (ExpandableListView) findViewById(R.id.result_expandablelistview);

        add_del_update_serach_cornerListView = (CornerListView) findViewById(R.id.add_del_update_serach_cornerListView);
        add_del_update_serach_cornerListView
                .setAdapter(new MyListAdapter(this));

        add_del_update_serach_cornerListView.setOnItemClickListener(this);

        // ��� tab�ı�ʱ���¼�����
        tabHost.setOnTabChangedListener(this);

        // ��ʼ��sensor��
        initSensor();

    }// end onCreate()

    @Override
    protected void onDestroy() {
        if (mediaPlayer != null) {
            mediaPlayer.release();
        }
        super.onDestroy();
    }

    @Override
    protected void onPause() {
        // TODO Auto-generated method stub
        super.onPause();
    }

    @Override
    protected void onRestart() {
        // TODO Auto-generated method stub
        super.onRestart();
    }

    @Override
    protected void onResume() {
        // TODO Auto-generated method stub
        super.onResume();
    }

    @Override
    protected void onStart() {
        // TODO Auto-generated method stub
        super.onStart();
    }

    @Override
    protected void onStop() {
        // TODO Auto-generated method stub
        super.onStop();
    }

    @SuppressLint("NewApi")
    private void initSensor() {
        mSensorManager = (SensorManager) getSystemService(SENSOR_SERVICE);

        List<Sensor> sensors = mSensorManager
                .getSensorList(Sensor.TYPE_ACCELEROMETER);
        if (sensors != null) {
            if (sensors.size() == 0) {
                return;
            }
        }// end if()

        // ��ɸ�Ӧ�����¼�
        SensorEventListener sensorelistener = new SensorEventListener() {
            @Override
            public void onAccuracyChanged(Sensor sensor, int accuracy) {
                // TODO Auto-generated method stub

            }

            // ��Ӧ������ı�
            @Override
            public void onSensorChanged(SensorEvent event) {
                // TODO Auto-generated method stub
                int sensorType = event.sensor.getType();

                // ��ȡҡһҡ����ֵ
                int shakeSenseValue = Integer.parseInt("14");
                // values[0]:X�ᣬvalues[1]��Y�ᣬvalues[2]��Z��
                float[] values = event.values;

                // �ڵ�һ��tab�²���Ӧҡһҡ�¼�
                if (tabHost.getCurrentTab() == 0
                        && sensorType == Sensor.TYPE_ACCELEROMETER) {
                    if ((Math.abs(values[0]) > shakeSenseValue)) {
                        // �����¼���ִ�д�Ӧ����Ϊ
                        afterSensorChangeDo();
                    }
                }// end if()

            }// end onSensorChanged()
        };

        // ע�Ქ��������ϵ������¼�
        mSensorManager.registerListener(sensorelistener,
                mSensorManager.getDefaultSensor(Sensor.TYPE_ACCELEROMETER),
                SensorManager.SENSOR_DELAY_NORMAL);

    }// end initSensor()

    /*
     * ҡ���ֻ�Ҫ����
     */
    private void afterSensorChangeDo() {
        if (isRunAfterSensorChangedDo) {
            return;
        }
        // �Ȱѱ�־��Ϊtrue
        isRunAfterSensorChangedDo = true;

        Log.d("eat", "after sensor change do!!!");

        List<List<String>> groupData;
        List<List<String>> childrenData;
        DatabaseHelper dbhelper;

        groupData = new ArrayList<List<String>>();
        childrenData = new ArrayList<List<String>>();

        dbhelper = new DatabaseHelper(this);
        SQLiteDatabase db = dbhelper.getReadableDatabase();

        Cursor cursor = db.rawQuery(
                "select * from ew_restaurant order by random() limit 1", null);
        while (cursor.moveToNext()) {
            int id = cursor.getInt(0);
            String name = cursor.getString(1);
            String address = cursor.getString(2);
            String phone = cursor.getString(3);
            String description = cursor.getString(4);

            List<String> subgrouplist = new ArrayList<String>();
            subgrouplist.add(id + "");
            subgrouplist.add(name);
            groupData.add(subgrouplist);

            List<String> subchildlist = new ArrayList<String>();

            // ��ַ�� + adress
            subchildlist.add(getResources().getString(R.string.address)
                    + address);
            // �绰��+ phone
            subchildlist.add(getResources().getString(R.string.phone) + phone);
            // ������+ description
            subchildlist.add(getResources().getString(R.string.description)
                    + description);

            childrenData.add(subchildlist);
        }// end while()

        cursor.close();
        db.close();

        mediaPlayer.start();

        resultExpandableListView.setAdapter(new MyExpandableListAdapter(this,
                groupData, childrenData));

        LayoutInflater inflater = getLayoutInflater();
        dialogLayout = inflater.inflate(R.layout.eatornot_dialog, null);
        Builder builder = new Builder(this);
        builder.setView(dialogLayout);
        builder.setPositiveButton(R.string.dialog_ok, new OnClickListener() {
            @Override
            public void onClick(DialogInterface dialog, int which) {
                isRunAfterSensorChangedDo = false;
                RadioGroup radiogroup = (RadioGroup) dialogLayout
                        .findViewById(R.id.dialog_radiogroup);
                RadioButton radiobutton = (RadioButton) findViewById(radiogroup
                        .getCheckedRadioButtonId());

                String meal = radiobutton.getText().toString();

            }
        });
        builder.setNegativeButton(R.string.dialog_shakeagain,
                new OnClickListener() {

                    @Override
                    public void onClick(DialogInterface dialog, int which) {
                        isRunAfterSensorChangedDo = false;
                        return;
                    }
                });

        builder.show();
    }// end afterSensorChangeDo()

    @Override
    public void onItemClick(AdapterView<?> parent, View view, int position,
                            long id) {
        Intent it = new Intent();
        switch (position) {
            case 0:
                it.setClass(MainActivity.this, AddRestaurantActivity.class);
                startActivity(it);
                break;
            case 1:
                it.setClass(MainActivity.this, AddMenuActivity.class);
                startActivity(it);
                break;
            case 2:
                it.setClass(MainActivity.this, DelUpdateRestaurantActivity.class);
                startActivity(it);
                break;
            case 3:
                it.setClass(MainActivity.this, DelUpdateMenuActivity.class);
                startActivity(it);
                break;
            case 4:
                it.setClass(MainActivity.this, SearchRestaurantActivity.class);
                startActivity(it);
                break;
            case 5:
                it.setClass(MainActivity.this, SearchMenuActivity.class);
                startActivity(it);
                break;
            default:

        }
    }

    @Override
    public void onTabChanged(String tabId) {
        Log.d("eat", "currentid = " + tabHost.getCurrentTab());
    }

}// end pulbic class MainActivity