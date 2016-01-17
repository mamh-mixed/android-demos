package com.cardinfolink.yunshouyin.activity;

import android.app.Activity;
import android.os.Bundle;
import android.widget.LinearLayout;

import com.cardinfolink.yunshouyin.R;
import com.cardinfolink.yunshouyin.view.ScanCodeView;

public class ScanCodeActivity extends Activity {

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_scan_code);


        LinearLayout.LayoutParams layoutParams = new LinearLayout.LayoutParams(LinearLayout.LayoutParams.MATCH_PARENT, LinearLayout.LayoutParams.MATCH_PARENT);
        ScanCodeView scanCodeView = new ScanCodeView(this, MainActivity.getHandler());
        scanCodeView.setLayoutParams(layoutParams);

        addContentView(scanCodeView, layoutParams);
    }

}
