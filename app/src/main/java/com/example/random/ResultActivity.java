package com.example.random;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStreamReader;
import java.text.ParseException;
import java.text.SimpleDateFormat;
import java.util.ArrayList;
import java.util.Collections;
import java.util.Date;

import android.annotation.SuppressLint;
import android.app.Activity;
import android.app.AlertDialog;
import android.os.Bundle;
import android.os.PowerManager;
import android.os.PowerManager.WakeLock;
import android.util.Log;
import android.view.View;
import android.widget.AdapterView;
import android.widget.AdapterView.OnItemClickListener;
import android.widget.ListView;

import com.example.model.AppInfo;
import com.example.model.KillLaunchInfo;
import com.example.model.AppInfoComparators;
import com.example.adapter.AppInfoListAdapter;

public class ResultActivity extends Activity {
    private WakeLock mWakeLock;
    private ListView resultListView;
    private AppInfoListAdapter mylistAdapter;
    private ArrayList<AppInfo> firstcheckedAppInfoList;
    private ArrayList<AppInfo> secondcheckedAppInfoList;
    private ArrayList<AppInfo> thirdcheckedAppInfoList;
    private ArrayList<AppInfo> allcheckedAppInfoList;
    private static final String TAG = "random";

    @SuppressLint("SimpleDateFormat")
    @SuppressWarnings("unchecked")
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);

        setContentView(R.layout.activity_result);

        firstcheckedAppInfoList = (ArrayList<AppInfo>) getIntent().getSerializableExtra("firstcheckedAppInfoList");

        secondcheckedAppInfoList = (ArrayList<AppInfo>) getIntent().getSerializableExtra("secondcheckedAppInfoList");

        thirdcheckedAppInfoList = (ArrayList<AppInfo>) getIntent().getSerializableExtra("thirdcheckedAppInfoList");

        resultListView = (ListView) findViewById(R.id.activity_result_ListView);

        createMyListAdapter();

        try {
            Process process = Runtime.getRuntime().exec("logcat -d -v time -b events -s am_kill");
            BufferedReader bufferedReader = new BufferedReader(new InputStreamReader(process.getInputStream()));

            String line = "";
            while ((line = bufferedReader.readLine()) != null) {
                String[] linearray = line.split(" ");
                String killdatetime = linearray[0] + " " + linearray[1];
                String[] words = linearray[6].split(",");
                String packagename = words[2];
                addKilldatetime(packagename, killdatetime, allcheckedAppInfoList);
            }
        } catch (IOException e) {

        }

        // sort the arraylist.
        for (AppInfo ai : allcheckedAppInfoList) {
            Collections.sort(ai.getResultArrayList(), new AppInfoComparators());
        }

        for (AppInfo ai : allcheckedAppInfoList) {
            ArrayList<KillLaunchInfo> resultArrayList = ai.getResultArrayList();
            for (KillLaunchInfo kli : resultArrayList) {
                if (kli.getKillorlaunch().equals(KillLaunchInfo.killFlag)) {
                    int index = resultArrayList.indexOf(kli);
                    if (index >= 1 && index <= resultArrayList.size()) {
                        SimpleDateFormat df = new SimpleDateFormat("MM-dd HH:mm:ss.SSS");
                        try {
                            Date killDate = df.parse(kli.getKillorlaunchdatetime());
                            Date launchDate = df.parse(resultArrayList.get(index - 1).getKillorlaunchdatetime());
                            long diff = killDate.getTime() - launchDate.getTime();
                            kli.setDifftime(diff);
                        } catch (ParseException e) {
                            e.printStackTrace();
                        }
                    }
                }
            }
        }
        
        
 
        Collections.sort(allcheckedAppInfoList, new AppInfoComparators());

        resultListView.setAdapter(mylistAdapter);

        resultListView.setOnItemClickListener(new OnItemClickListener() {

            @Override
            public void onItemClick(AdapterView<?> parent, View view, int position, long id) {

                ArrayList<KillLaunchInfo> resultArrayList = allcheckedAppInfoList.get(position).getResultArrayList();
                ArrayList<String> resultStringArrayList = new ArrayList<String>();
                for (KillLaunchInfo klinfo : resultArrayList) {
                    resultStringArrayList.add(klinfo.getKillorlaunch() + ": " + klinfo.getKillorlaunchdatetime() + (klinfo.getKillorlaunch().equals(KillLaunchInfo.killFlag) ? " diff time: (" + ((float) klinfo.getDifftime()) / 1000 + " s)" : ""));
                }
                Object[] array = resultStringArrayList.toArray(new CharSequence[resultStringArrayList.size()]);
                CharSequence[] cs = (CharSequence[]) array;

                AlertDialog.Builder builder = new AlertDialog.Builder(ResultActivity.this);
                builder.setTitle(getString(R.string.dialog_result_title));
                builder.setPositiveButton(getString(R.string.dialog_ok), null);

                builder.setItems(cs, null);
                builder.show();

            }
        });
    }

    private void addKilldatetime(String packageName, String killdatetime, ArrayList<AppInfo> allList) {
        for (AppInfo ai : allList) {
            if (ai.getPackageName().equals(packageName)) {
                ai.getResultArrayList().add(new KillLaunchInfo(KillLaunchInfo.killFlag, killdatetime));
                ai.setKillTimes(ai.getKillTimes() + 1);
            }
        }

    }

    private void createMyListAdapter() {
        allcheckedAppInfoList = new ArrayList<AppInfo>();
        if (firstcheckedAppInfoList != null) {
            for (AppInfo firstcheckedAppInfo : firstcheckedAppInfoList) {
                int launchtimes = firstcheckedAppInfo.getLaunchTimes();
                if (launchtimes != 0) {
                    allcheckedAppInfoList.add(firstcheckedAppInfo);
                }
            }
        }
        if (secondcheckedAppInfoList != null) {
            for (AppInfo secondcheckedAppInfo : secondcheckedAppInfoList) {
                int launchtimes = secondcheckedAppInfo.getLaunchTimes();
                if (launchtimes != 0) {
                    allcheckedAppInfoList.add(secondcheckedAppInfo);
                }
            }
        }
        if (thirdcheckedAppInfoList != null) {
            for (AppInfo thirdcheckedAppInfo : thirdcheckedAppInfoList) {
                int launchtimes = thirdcheckedAppInfo.getLaunchTimes();
                if (launchtimes != 0) {
                    allcheckedAppInfoList.add(thirdcheckedAppInfo);
                }
            }
        }
        mylistAdapter = new AppInfoListAdapter(this, allcheckedAppInfoList);
    }

    @Override
    protected void onPause() {
        Log.d(TAG, "ResultActivity,onPause");
        super.onPause();
        releaseWakeLock();
    }

    @Override
    protected void onResume() {
        Log.d(TAG, "ResultActivity,onResume");
        super.onResume();
        acquireWakeLock();
    }

    private void acquireWakeLock() {
        if (mWakeLock == null) {
            PowerManager mPowerManager = ((PowerManager) getSystemService(POWER_SERVICE));
            mWakeLock = mPowerManager.newWakeLock(PowerManager.SCREEN_BRIGHT_WAKE_LOCK, this.getClass().getCanonicalName());
            mWakeLock.acquire();
        }
    }

    private void releaseWakeLock() {
        if (mWakeLock != null && mWakeLock.isHeld()) {
            mWakeLock.release();
            mWakeLock = null;
        }
    }

}
