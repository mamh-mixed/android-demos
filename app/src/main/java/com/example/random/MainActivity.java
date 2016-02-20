package com.example.random;

import java.util.ArrayList;
import java.util.List;

import android.annotation.SuppressLint;
import android.app.AlertDialog;
import android.content.Context;
import android.content.DialogInterface;
import android.content.DialogInterface.OnClickListener;
import android.content.Intent;
import android.content.SharedPreferences;
import android.content.SharedPreferences.Editor;
import android.content.pm.ApplicationInfo;
import android.content.pm.PackageInfo;
import android.os.Bundle;
import android.os.PowerManager;
import android.os.PowerManager.WakeLock;
import android.preference.CheckBoxPreference;
import android.preference.Preference;
import android.preference.Preference.OnPreferenceChangeListener;
import android.preference.Preference.OnPreferenceClickListener;
import android.preference.PreferenceActivity;
import android.preference.PreferenceScreen;
import android.preference.SwitchPreference;
import android.util.Log;
import android.view.View;
import android.widget.EditText;

import com.example.model.AppInfo;
import com.example.wheellibrary.NumericWheelAdapter;
import com.example.wheellibrary.WheelView;

@SuppressLint("NewApi")
public class MainActivity extends PreferenceActivity implements
        OnPreferenceClickListener, OnPreferenceChangeListener {
    private PreferenceScreen mFirstPreferenceScreen;
    private PreferenceScreen mSecondPreferenceScreen;
    private PreferenceScreen mThirdPreferenceScreen;
    private PreferenceScreen mRatioPreferenceScreen;
    private PreferenceScreen mTotalPreferenceScreen;
    private SwitchPreference mStartSwitchPreference;
    private CheckBoxPreference mSystemappCheckBoxPreference;

    private SharedPreferences mPreferences;

    private boolean startEnable;
    private boolean systemApp;

    private WakeLock mWakeLock;

    private ArrayList<AppInfo> firstcheckedAppInfoList;
    private ArrayList<AppInfo> secondcheckedAppInfoList;
    private ArrayList<AppInfo> thirdcheckedAppInfoList;
    private ArrayList<AppInfo> allAppInfoList;

    private static final int firstCurrentItemIndex = 5;
    private static final int secondCurrentItemIndex = 3;
    private static final int thirdCurrentItemIndex = 2;

    private int[] ratioArray = {0, 1, 2, 3, 4, 5, 6, 7, 8, 9};

    private int firstRatio;
    private int secondRatio;
    private int thirdRatio;
    private int total;

    private WheelView firstratioWheelView;
    private WheelView secondratioWheelView;
    private WheelView thirdratioWheelView;

    private EditText totalEditText;

    private static final String TAG = "random";

    public static final int FIRST_RESULT = 1;
    public static final int SECOND_RESULT = 2;
    public static final int THIRD_RESULT = 3;

    private static final String KEY_main_preferencescreen_first = "main_preferencescreen_first";
    private static final String KEY_main_preferencescreen_second = "main_preferencescreen_second";
    private static final String KEY_main_preferencescreen_third = "main_preferencescreen_third";
    private static final String KEY_main_preferencescreen_ratio = "main_preferencescreen_ratio";
    private static final String KEY_main_preferencescreen_total = "main_preferencescreen_total";
    private static final String KEY_main_switchpreference_start = "main_switchpreference_start";
    private static final String KEY_main_checkboxpreference_systemapp = "main_checkboxpreference_systemapp";

    private static final String START_ENABLE = "START_ENABLE";
    private static final String FIRST_RATIO = "FIRST_RATIO";
    private static final String SECOND_RATIO = "SECOND_RATIO";
    private static final String THIRD_RATIO = "THIRD_RATIO";
    private static final String TOTAL = "TOTAL";
    private static final String SYSTEM_APP = "SYSTEM_APP";

    @SuppressWarnings("deprecation")
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        addPreferencesFromResource(R.xml.main);
        readFromPreference(getBaseContext());

        mFirstPreferenceScreen = (PreferenceScreen) findPreference(KEY_main_preferencescreen_first);
        mSecondPreferenceScreen = (PreferenceScreen) findPreference(KEY_main_preferencescreen_second);
        mThirdPreferenceScreen = (PreferenceScreen) findPreference(KEY_main_preferencescreen_third);
        mRatioPreferenceScreen = (PreferenceScreen) findPreference(KEY_main_preferencescreen_ratio);
        mTotalPreferenceScreen = (PreferenceScreen) findPreference(KEY_main_preferencescreen_total);
        mStartSwitchPreference = (SwitchPreference) findPreference(KEY_main_switchpreference_start);
        mSystemappCheckBoxPreference = (CheckBoxPreference) findPreference(KEY_main_checkboxpreference_systemapp);

        allAppInfoList = new ArrayList<AppInfo>();

        firstcheckedAppInfoList = new ArrayList<AppInfo>();
        secondcheckedAppInfoList = new ArrayList<AppInfo>();
        thirdcheckedAppInfoList = new ArrayList<AppInfo>();

        if (systemApp) {
            allAppInfoList.clear();
            List<PackageInfo> packages = getPackageManager().getInstalledPackages(0);
            for (int i = 0; i < packages.size(); i++) {
                PackageInfo packageInfo = packages.get(i);
                AppInfo tmpInfo = new AppInfo();
                tmpInfo.setAppName(packageInfo.applicationInfo.loadLabel(getPackageManager()).toString());
                tmpInfo.setPackageName(packageInfo.packageName);
                tmpInfo.setVersionName(packageInfo.versionName);
                tmpInfo.setVersionCode(packageInfo.versionCode);
                // Only display the system app info
                if ((packageInfo.applicationInfo.flags & ApplicationInfo.FLAG_SYSTEM) != 0) {
                    if (tmpInfo.getAppName().equals(getResources().getString(R.string.app_name))) {
                        continue;
                    }
                    allAppInfoList.add(tmpInfo);//
                }
            }// end for()
        } else {
            allAppInfoList.clear();
            List<PackageInfo> packages = getPackageManager().getInstalledPackages(0);
            for (int i = 0; i < packages.size(); i++) {
                PackageInfo packageInfo = packages.get(i);
                AppInfo tmpInfo = new AppInfo();
                tmpInfo.setAppName(packageInfo.applicationInfo.loadLabel(getPackageManager()).toString());
                tmpInfo.setPackageName(packageInfo.packageName);
                tmpInfo.setVersionName(packageInfo.versionName);
                tmpInfo.setVersionCode(packageInfo.versionCode);
                // Only display the non-system app info
                if ((packageInfo.applicationInfo.flags & ApplicationInfo.FLAG_SYSTEM) == 0) {
                    if (tmpInfo.getAppName().equals(getResources().getString(R.string.app_name))) {
                        continue;
                    }
                    allAppInfoList.add(tmpInfo);//
                }// end if()
            }// end for()
        }

        mFirstPreferenceScreen.setOnPreferenceClickListener(this);
        mSecondPreferenceScreen.setOnPreferenceClickListener(this);
        mThirdPreferenceScreen.setOnPreferenceClickListener(this);
        mRatioPreferenceScreen.setOnPreferenceClickListener(this);
        mTotalPreferenceScreen.setOnPreferenceClickListener(this);
        mStartSwitchPreference.setOnPreferenceChangeListener(this);
        mSystemappCheckBoxPreference.setOnPreferenceChangeListener(this);

        setSavedPreference();
    }

    protected void onActivityResult(int requestCode, int resultCode, Intent data) {
        Log.d(TAG, "MainActivity,onActivityResult");
        switch (resultCode) {
            case FIRST_RESULT:
                firstcheckedAppInfoList = (ArrayList<AppInfo>) data.getSerializableExtra("firstcheckedAppInfoList");
                break;
            case SECOND_RESULT:
                secondcheckedAppInfoList = (ArrayList<AppInfo>) data.getSerializableExtra("secondcheckedAppInfoList");
                break;
            case THIRD_RESULT:
                thirdcheckedAppInfoList = (ArrayList<AppInfo>) data.getSerializableExtra("thirdcheckedAppInfoList");
                break;

            default:
        }// end switch()
    }

    @Override
    protected void onDestroy() {
        Log.d(TAG, "MainActivity,onDestroy");
        Intent intent = new Intent(MainActivity.this, MainService.class);
        stopService(intent);
        super.onDestroy();
    }

    @Override
    protected void onStart() {
        Log.d(TAG, "MainActivity,onStart");
        super.onStart();
    }

    @Override
    protected void onStop() {
        Log.d(TAG, "MainActivity,onStop");
        super.onStop();
    }

    @Override
    protected void onPause() {
        Log.d(TAG, "MainActivity,onPause");
        super.onPause();
        releaseWakeLock();
    }

    @Override
    protected void onRestart() {
        Log.d(TAG, "MainActivity,onRestart");
        super.onRestart();
    }

    @Override
    protected void onResume() {
        Log.d(TAG, "MainActivity,onResume");
        super.onResume();
        // 获取锁，保持屏幕亮度
        readFromPreference(this);
        setSavedPreference();
        acquireWakeLock();
    }

    @SuppressWarnings("deprecation")
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

    @Override
    public boolean onPreferenceClick(Preference preference) {
        Log.d(TAG, "click" + preference.getKey());

        if (KEY_main_preferencescreen_first.equals(preference.getKey())) {
            Intent intent = new Intent();
            intent.setClass(MainActivity.this, FirstActivity.class);
            intent.putExtra("allAppInfoList", allAppInfoList);
            intent.putExtra("firstcheckedAppInfoList", firstcheckedAppInfoList);
            intent.putExtra("secondcheckedAppInfoList", secondcheckedAppInfoList);
            intent.putExtra("thirdcheckedAppInfoList", thirdcheckedAppInfoList);
            startActivityForResult(intent, FIRST_RESULT);
        } else if (KEY_main_preferencescreen_second.equals(preference.getKey())) {
            Intent intent = new Intent();
            intent.setClass(MainActivity.this, SecondActivity.class);
            intent.putExtra("allAppInfoList", allAppInfoList);
            intent.putExtra("firstcheckedAppInfoList", firstcheckedAppInfoList);
            intent.putExtra("secondcheckedAppInfoList", secondcheckedAppInfoList);
            intent.putExtra("thirdcheckedAppInfoList", thirdcheckedAppInfoList);
            startActivityForResult(intent, SECOND_RESULT);
        } else if (KEY_main_preferencescreen_third.equals(preference.getKey())) {
            Intent intent = new Intent();
            intent.setClass(MainActivity.this, ThirdActivity.class);
            intent.putExtra("allAppInfoList", allAppInfoList);
            intent.putExtra("firstcheckedAppInfoList", firstcheckedAppInfoList);
            intent.putExtra("secondcheckedAppInfoList", secondcheckedAppInfoList);
            intent.putExtra("thirdcheckedAppInfoList", thirdcheckedAppInfoList);
            startActivityForResult(intent, THIRD_RESULT);
        } else if (KEY_main_preferencescreen_ratio.equals(preference.getKey())) {
            View dialogLayout = getLayoutInflater().inflate(R.layout.ratio, null);
            firstratioWheelView = (WheelView) dialogLayout.findViewById(R.id.activity_main_WheelView1);
            secondratioWheelView = (WheelView) dialogLayout.findViewById(R.id.activity_main_WheelView2);
            thirdratioWheelView = (WheelView) dialogLayout.findViewById(R.id.activity_main_WheelView3);

            firstratioWheelView.setAdapter(new NumericWheelAdapter(ratioArray));
            secondratioWheelView.setAdapter(new NumericWheelAdapter(ratioArray));
            thirdratioWheelView.setAdapter(new NumericWheelAdapter(ratioArray));

            firstratioWheelView.setCurrentItem(firstRatio);
            secondratioWheelView.setCurrentItem(secondRatio);
            thirdratioWheelView.setCurrentItem(thirdRatio);

            AlertDialog.Builder builder = new AlertDialog.Builder(this);
            builder.setTitle(getString(R.string.dialog_ratio_title));
            builder.setView(dialogLayout);
            builder.setPositiveButton(getString(R.string.dialog_ok), new OnClickListener() {
                @Override
                public void onClick(DialogInterface arg0, int arg1) {
                    firstRatio = firstratioWheelView.getCurrentItem();
                    secondRatio = secondratioWheelView.getCurrentItem();
                    thirdRatio = thirdratioWheelView.getCurrentItem();
                    mRatioPreferenceScreen.setSummary(getString(R.string.main_preferencescreen_ratio_summary)
                            + ", "
                            + firstRatio
                            + " : "
                            + secondRatio + " : " + thirdRatio);
                }
            });
            builder.setNegativeButton(getString(R.string.dialog_cancel), null);
            builder.show();
        } else if (KEY_main_preferencescreen_total.equals(preference.getKey())) {
            View dialogLayout = getLayoutInflater().inflate(R.layout.total, null);
            totalEditText = (EditText) dialogLayout.findViewById(R.id.total_edittext_total);
            totalEditText.setText(total + "");
            AlertDialog.Builder builder = new AlertDialog.Builder(this);
            builder.setTitle(getString(R.string.dialog_total_title));
            builder.setView(dialogLayout);
            builder.setPositiveButton(getString(R.string.dialog_ok), new OnClickListener() {
                @Override
                public void onClick(DialogInterface arg0, int arg1) {
                    if (totalEditText.getText().toString().length() == 0) {
                        total = (Integer.MAX_VALUE);
                    } else {
                        try {
                            total = (Integer.parseInt(totalEditText.getText().toString()));
                        } catch (NumberFormatException e) {
                            total = Integer.MAX_VALUE;
                        }
                    }
                    // 设置一下totalPreferenceScreen的副标题
                    mTotalPreferenceScreen.setSummary(getString(R.string.main_preferencescreen_total_summary)
                            + ", "
                            + ((total == Integer.MAX_VALUE) ? "MAX"
                            : ("" + total)));
                }
            });
            builder.setNegativeButton(getString(R.string.dialog_cancel), null);
            builder.show();
        }
        saveToPreference();
        return true;
    }

    @Override
    public boolean onPreferenceChange(Preference preference, Object newValue) {
        Log.d(TAG, "change" + preference.getKey());
        if (KEY_main_switchpreference_start.equals(preference.getKey())) {
            boolean enable = (Boolean) newValue;

            Intent intent = new Intent(MainActivity.this, MainService.class);

            intent.putExtra("firstRatio", firstRatio);
            intent.putExtra("secondRatio", secondRatio);
            intent.putExtra("thirdRatio", thirdRatio);
            intent.putExtra("total", total);
            intent.putExtra("allAppInfoList", allAppInfoList);
            intent.putExtra("firstcheckedAppInfoList", firstcheckedAppInfoList);
            intent.putExtra("secondcheckedAppInfoList", secondcheckedAppInfoList);
            intent.putExtra("thirdcheckedAppInfoList", thirdcheckedAppInfoList);

            if (enable) {
                startService(intent);
            } else {
                stopService(intent);
            }
            startEnable = enable;
        } else if (KEY_main_checkboxpreference_systemapp.equals(preference.getKey())) {
            boolean enable = (Boolean) newValue;
            firstcheckedAppInfoList.clear();
            secondcheckedAppInfoList.clear();
            thirdcheckedAppInfoList.clear();
            if (enable) {
                allAppInfoList.clear();
                List<PackageInfo> packages = getPackageManager().getInstalledPackages(0);
                for (int i = 0; i < packages.size(); i++) {
                    PackageInfo packageInfo = packages.get(i);
                    AppInfo tmpInfo = new AppInfo();
                    tmpInfo.setAppName(packageInfo.applicationInfo.loadLabel(getPackageManager()).toString());
                    tmpInfo.setPackageName(packageInfo.packageName);
                    tmpInfo.setVersionName(packageInfo.versionName);
                    tmpInfo.setVersionCode(packageInfo.versionCode);
                    // Only display the system app info
                    if ((packageInfo.applicationInfo.flags & ApplicationInfo.FLAG_SYSTEM) != 0) {
                        if (tmpInfo.getAppName().equals(getString(R.string.app_name))) {
                            continue;
                        }
                        allAppInfoList.add(tmpInfo);//
                    }
                }// end for()
            } else {
                allAppInfoList.clear();
                List<PackageInfo> packages = getPackageManager().getInstalledPackages(0);
                for (int i = 0; i < packages.size(); i++) {
                    PackageInfo packageInfo = packages.get(i);
                    AppInfo tmpInfo = new AppInfo();
                    tmpInfo.setAppName(packageInfo.applicationInfo.loadLabel(getPackageManager()).toString());
                    tmpInfo.setPackageName(packageInfo.packageName);
                    tmpInfo.setVersionName(packageInfo.versionName);
                    tmpInfo.setVersionCode(packageInfo.versionCode);
                    // Only display the non-system app info
                    if ((packageInfo.applicationInfo.flags & ApplicationInfo.FLAG_SYSTEM) == 0) {
                        if (tmpInfo.getAppName().equals(getString(R.string.app_name))) {
                            continue;
                        }
                        allAppInfoList.add(tmpInfo);//
                    }// end if()
                }// end for()
            }
            systemApp = enable;
        }
        saveToPreference();
        return true;
    }

    // 设置startSwitchPreference的值
    private void setSavedPreference() {
        mStartSwitchPreference.setChecked(startEnable);
        mSystemappCheckBoxPreference.setChecked(systemApp);

        String setRatioStr = getString(R.string.main_preferencescreen_ratio_summary);
        // 设置一下ratioPreferenceScreen的副标题
        mRatioPreferenceScreen.setSummary(setRatioStr + ", " + firstRatio + " : " + secondRatio + " : " + thirdRatio);
        // 设置一下totalPreferenceScreen的副标题
        mTotalPreferenceScreen.setSummary(setRatioStr + ", " + ((total == Integer.MAX_VALUE) ? "MAX" : ("" + total)));
    }

    private void saveToPreference() {
        mPreferences = getSharedPreferences("MainActivity", Context.MODE_PRIVATE);
        Editor preferenceEditor = mPreferences.edit();
        preferenceEditor.putBoolean(START_ENABLE, startEnable);
        preferenceEditor.putBoolean(SYSTEM_APP, systemApp);
        preferenceEditor.putInt(FIRST_RATIO, firstRatio);
        preferenceEditor.putInt(SECOND_RATIO, secondRatio);
        preferenceEditor.putInt(THIRD_RATIO, thirdRatio);
        preferenceEditor.putInt(TOTAL, total);
        preferenceEditor.apply();
    }

    private void readFromPreference(Context context) {
        mPreferences = context.getSharedPreferences("MainActivity", Context.MODE_PRIVATE);
        if (mPreferences == null) {
            return;
        }
        // 是否开始了。默认是false。
        startEnable = mPreferences.getBoolean(START_ENABLE, false);
        // 是否是系统的应用？ 默认true。
        systemApp = mPreferences.getBoolean(SYSTEM_APP, true);
        firstRatio = mPreferences.getInt(FIRST_RATIO, ratioArray[firstCurrentItemIndex]);
        secondRatio = mPreferences.getInt(SECOND_RATIO, ratioArray[secondCurrentItemIndex]);
        thirdRatio = mPreferences.getInt(THIRD_RATIO, ratioArray[thirdCurrentItemIndex]);
        total = mPreferences.getInt(TOTAL, Integer.MAX_VALUE);
    }
}
