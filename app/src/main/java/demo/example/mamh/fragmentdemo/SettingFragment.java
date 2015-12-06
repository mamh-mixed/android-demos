package demo.example.mamh.fragmentdemo;

import android.os.Bundle;
import android.app.Fragment;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;

/**
 * Created by mamh on 15-12-6.
 */
public class SettingFragment extends Fragment {

    @Override
    public View onCreateView(LayoutInflater inflater, ViewGroup container, Bundle savedInstanceState) {
        View settingLayout = inflater.inflate(R.layout.setting, container, false);
        return settingLayout;
    }

}
