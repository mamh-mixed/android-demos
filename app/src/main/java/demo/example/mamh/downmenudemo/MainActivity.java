package demo.example.mamh.downmenudemo;

import android.support.v7.app.AppCompatActivity;
import android.os.Bundle;
import android.util.Log;
import android.view.View;
import android.view.ViewGroup;
import android.widget.BaseAdapter;
import android.widget.EditText;
import android.widget.ImageView;
import android.widget.ListAdapter;
import android.widget.ListView;
import android.widget.PopupWindow;
import android.widget.TextView;

import java.util.ArrayList;
import java.util.List;

public class MainActivity extends AppCompatActivity {

    private EditText input;
    private ImageView downarrow;


    private List<String> msgList = new ArrayList<String>();

    private PopupWindow popupWindow;
    private ListView listView;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);

        input = (EditText) findViewById(R.id.input);
        downarrow = (ImageView) findViewById(R.id.down_arrow);


        for (int i = 0; i < 20; i++) {
            msgList.add("sdkfjalfjasdlf" + i);
        }

        initListView();



        downarrow.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                Log.e("xxx", "click");
                popupWindow = new PopupWindow(MainActivity.this);
                popupWindow.setWidth(input.getWidth());
                popupWindow.setHeight(200);
                popupWindow.setContentView(listView);
                popupWindow.setOutsideTouchable(true);
                popupWindow.showAsDropDown(input, 0, 0);
            }
        });


    }

    private void initListView() {
        listView = new ListView(this);
        ListAdapter adapter = new ListAdapter();
        listView.setAdapter(adapter);
        listView.setBackgroundResource(R.drawable.listview_background);
    }

    private class ListAdapter extends BaseAdapter {

        @Override
        public int getCount() {
            return msgList.size();
        }

        @Override
        public Object getItem(int position) {
            return null;
        }

        @Override
        public long getItemId(int position) {
            return 0;
        }

        @Override
        public View getView(final int position, View convertView, ViewGroup parent) {
            ViewHolder holder = null;
            if (convertView == null) {
                convertView = View.inflate(getApplicationContext(), R.layout.list_item, null);
                holder = new ViewHolder();
                holder.delete = (ImageView) convertView.findViewById(R.id.iv_list_delete);
                holder.msg = (TextView) convertView.findViewById(R.id.tv_list_item);
                convertView.setTag(holder);

            } else {
                holder = (ViewHolder) convertView.getTag();
            }
            holder.msg.setText(msgList.get(position));
            holder.delete.setOnClickListener(new View.OnClickListener() {
                @Override
                public void onClick(View v) {
                    //delete item
                    msgList.remove(position);
                    notifyDataSetChanged();
                }
            });

            convertView.setOnClickListener(new View.OnClickListener() {
                @Override
                public void onClick(View v) {
                    input.setText(msgList.get(position));
                    popupWindow.dismiss();
                }
            });
            return convertView;

        }
    }


    private class ViewHolder {
        TextView msg;
        ImageView delete;
    }
}
