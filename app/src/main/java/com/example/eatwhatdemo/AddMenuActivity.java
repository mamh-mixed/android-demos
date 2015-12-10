package com.example.eatwhatdemo;

import android.app.Activity;
import android.content.ContentValues;
import android.database.sqlite.SQLiteDatabase;
import android.os.Bundle;
import android.view.Menu;
import android.view.View;
import android.view.View.OnClickListener;
import android.widget.Button;
import android.widget.EditText;
import android.widget.Spinner;
import android.widget.TextView;

import com.example.library.DatabaseHelper;
import com.example.adapter.MySpinnerAdapter;

public class AddMenuActivity extends Activity {

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_addmenu);

        dbhelper = new DatabaseHelper(this);

        commitButton = (Button) findViewById(R.id.add_menu_commit_button);
        backButton = (Button) findViewById(R.id.add_menu_back_button);
        nameeditText = (EditText) findViewById(R.id.add_menu_name_editText);
        descriptioneditText = (EditText) findViewById(R.id.add_menu_description_editText);
        restaurantSpinner = (Spinner) findViewById(R.id.add_restaurant_spinner);

        restaurantSpinner.setAdapter(new MySpinnerAdapter(this));

        commitButton.setOnClickListener(new OnClickListener() {
            @Override
            public void onClick(View v) {
                SQLiteDatabase db = dbhelper.getReadableDatabase();

                menuName = nameeditText.getText().toString();
                if (menuName.length() == 0) {
                    menuName = "menuname";
                }
                View view = restaurantSpinner.getSelectedView();
                TextView idtextView = (TextView) view.findViewById(R.id.spinner_item_id_textView);
                restaurantId = Integer.parseInt(idtextView.getText().toString());

                menuDescription = descriptioneditText.getText().toString();
                if (menuDescription.length() == 0) {
                    menuDescription = "menudescription";
                }

                ContentValues values = new ContentValues();
                values.put("menu_name", menuName);
                values.put("restaurant_id", restaurantId);
                values.put("menu_description", menuDescription);

                db.insert("ew_menu", null, values);
                db.close();
                finish();
            }
        });

        backButton.setOnClickListener(new OnClickListener() {

            @Override
            public void onClick(View v) {
                finish();
            }
        });

    }

    @Override
    public boolean onCreateOptionsMenu(Menu menu) {

        return true;
    }

    private DatabaseHelper dbhelper;

    private Spinner restaurantSpinner;

    private Button commitButton;
    private Button backButton;

    private EditText nameeditText;
    private EditText descriptioneditText;

    private String menuName;
    private int restaurantId;
    private String menuDescription;

}
