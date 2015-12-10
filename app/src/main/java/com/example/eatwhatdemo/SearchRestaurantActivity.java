package com.example.eatwhatdemo;

import android.os.Bundle;
import android.app.Activity;
import android.view.Menu;

public class SearchRestaurantActivity extends Activity {

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_searchrestaurant);
    }

    @Override
    public boolean onCreateOptionsMenu(Menu menu) {

        return true;
    }

}
