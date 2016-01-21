package com.lvren.demoqq;

import android.app.Activity;
import android.os.Bundle;

public class MainActivity extends Activity {

    private TestAdapter adapter;
    private PinnedHeaderListView listView;

    @Override
    public void onCreate(final Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.main);
        adapter = new TestAdapter(getLayoutInflater());

        listView = (PinnedHeaderListView) findViewById(R.id.section_list_view);
        listView.setAdapter(adapter);
        listView.setOnScrollListener(adapter);
        listView.setPinnedHeaderView(getLayoutInflater().inflate(
                R.layout.list_section, listView, false));
    }
}
