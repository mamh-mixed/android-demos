package com.example.random;

import java.util.ArrayList;

import android.app.Activity;
import android.content.Intent;
import android.os.Bundle;
import android.view.View;
import android.view.View.OnClickListener;
import android.widget.Button;
import android.widget.CheckBox;
import android.widget.CompoundButton;
import android.widget.CompoundButton.OnCheckedChangeListener;
import android.widget.LinearLayout;
import android.widget.ScrollView;

import com.example.model.AppInfo;

public class ThirdActivity extends Activity {

    @SuppressWarnings("unchecked")
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);

        LinearLayout linearlayout = new LinearLayout(this);
        linearlayout.setOrientation(LinearLayout.VERTICAL);

        allAppInfoList = new ArrayList<AppInfo>();
        allAppInfoList = ((ArrayList<AppInfo>) getIntent().getSerializableExtra("allAppInfoList"));

        firstcheckedAppInfoList = new ArrayList<AppInfo>();
        firstcheckedAppInfoList = (ArrayList<AppInfo>) getIntent().getSerializableExtra("firstcheckedAppInfoList");
        secondcheckedAppInfoList = new ArrayList<AppInfo>();
        secondcheckedAppInfoList = (ArrayList<AppInfo>) getIntent().getSerializableExtra("secondcheckedAppInfoList");
        thirdcheckedAppInfoList = new ArrayList<AppInfo>();
        thirdcheckedAppInfoList = (ArrayList<AppInfo>) getIntent().getSerializableExtra("thirdcheckedAppInfoList");

        CheckBox allCheckBox = new CheckBox(this);
        allCheckBox.setText(R.string.allcheckbox_text);
        linearlayout.addView(allCheckBox);

        for (int i = 0; i < allAppInfoList.size(); i++) {
            AppInfo appinfo = allAppInfoList.get(i);
            CheckBox cb = new CheckBox(this);
            cb.setText(appinfo.getAppName());
            if (contains(appinfo, thirdcheckedAppInfoList)) {
                cb.setChecked(true);
            } else {
                cb.setChecked(false);
            }
            if (contains(appinfo, firstcheckedAppInfoList)
                    || contains(appinfo, secondcheckedAppInfoList)) {
                cb.setEnabled(false);
            }
            cb.setId(i);
            linearlayout.addView(cb);
        }

        Button confirmButton = new Button(this);
        confirmButton.setText(R.string.confirmbutton_text);
        linearlayout.addView(confirmButton);

        ScrollView sv = (ScrollView) this.getLayoutInflater().inflate(R.layout.activity_first, null);
        sv.addView(linearlayout);
        setContentView(sv);

        confirmButton.setOnClickListener(new OnClickListener() {
            @Override
            public void onClick(View arg0) {
                ArrayList<AppInfo> checkedAppInfoList = new ArrayList<AppInfo>();

                for (int i = 0; i < allAppInfoList.size(); i++) {
                    CheckBox cb = (CheckBox) findViewById(i);
                    if (cb.isChecked()) {
                        checkedAppInfoList.add(allAppInfoList.get(i));
                    }
                }

                Intent intent = new Intent();
                intent.setClass(ThirdActivity.this, MainActivity.class);
                intent.putExtra("thirdcheckedAppInfoList", checkedAppInfoList);
                setResult(THIRD_RESULT, intent);
                finish();
            }

        });

        allCheckBox.setOnCheckedChangeListener(new OnCheckedChangeListener() {
            @Override
            public void onCheckedChanged(CompoundButton buttonView,
                                         boolean isChecked) {
                if (isChecked) {
                    for (int i = 0; i < allAppInfoList.size(); i++) {
                        CheckBox cb = (CheckBox) findViewById(i);
                        if (cb.isEnabled()) {
                            cb.setChecked(true);
                        }
                    }
                } else {
                    for (int i = 0; i < allAppInfoList.size(); i++) {
                        CheckBox cb = (CheckBox) findViewById(i);
                        if (cb.isEnabled()) {
                            cb.setChecked(false);
                        }
                    }
                }
            }

        });
    }

    public boolean contains(AppInfo appinfo, ArrayList<AppInfo> appinfolist) {
        for (int i = 0; i < appinfolist.size(); i++) {
            if (appinfolist.get(i).getAppName().equals(appinfo.getAppName())) {
                return true;
            }
        }
        return false;
    }

    private ArrayList<AppInfo> firstcheckedAppInfoList;
    private ArrayList<AppInfo> secondcheckedAppInfoList;
    private ArrayList<AppInfo> thirdcheckedAppInfoList;

    private ArrayList<AppInfo> allAppInfoList;
    public static final int FIRST_RESULT = 1;
    public static final int SECOND_RESULT = 2;
    public static final int THIRD_RESULT = 3;

}
