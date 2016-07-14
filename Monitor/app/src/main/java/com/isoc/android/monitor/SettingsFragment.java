package com.isoc.android.monitor;


import android.content.ComponentName;
import android.content.SharedPreferences;
import android.content.pm.PackageManager;
import android.os.Bundle;
import android.preference.CheckBoxPreference;
import android.preference.EditTextPreference;
import android.preference.Preference;
import android.preference.PreferenceFragment;
import android.support.v4.app.Fragment;
import android.support.v7.widget.Toolbar;
import android.view.LayoutInflater;
import android.view.Menu;
import android.view.MenuItem;
import android.view.View;
import android.view.ViewGroup;

/**
 * A simple {@link Fragment} subclass.
 */
public class SettingsFragment extends PreferenceFragment implements SharedPreferences.OnSharedPreferenceChangeListener {
    @Override
    public void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setHasOptionsMenu(true);

    }

    @Override
    public void onResume() {
        super.onResume();
        getPreferenceScreen().getSharedPreferences().registerOnSharedPreferenceChangeListener(this);

    }

    @Override
    public void onPause() {
        super.onPause();
        getPreferenceScreen().getSharedPreferences().unregisterOnSharedPreferenceChangeListener(this);
    }

    @Override
    public View onCreateView(LayoutInflater inflater, ViewGroup container, Bundle savedInstanceState) {
        View view=super.onCreateView(inflater, container, savedInstanceState);
        addPreferencesFromResource(R.xml.preferences);
        EditTextPreference url = (EditTextPreference) findPreference("server_url");
        url.setSummary(url.getText());
        Toolbar toolbar= (Toolbar) getActivity().findViewById(R.id.toolbar);
        toolbar.setTitle("Settings");
        return view;
    }

    @Override
    public void onPrepareOptionsMenu(Menu menu) {
        MenuItem item=menu.findItem(R.id.action_settings);
        item.setVisible(false);
        super.onPrepareOptionsMenu(menu);
    }


    @Override
    public void onSharedPreferenceChanged(SharedPreferences sharedPreferences, String s) {
        Preference preference=findPreference(s);
        if (preference instanceof EditTextPreference){
            EditTextPreference editPreference= (EditTextPreference) preference;
            editPreference.setSummary(editPreference.getText());
        }
        if (preference instanceof CheckBoxPreference) {
            if (preference.getKey().equals("battery")) {
                ComponentName batteryReceiver = new ComponentName(getActivity(),BatteryReceiver.class);
                if (((CheckBoxPreference) preference).isChecked()) setComponentState(true,batteryReceiver);
                else setComponentState(false,batteryReceiver);
            }
        }
    }

    private void setComponentState(boolean state,ComponentName name){
        int newState = (state) ? PackageManager.COMPONENT_ENABLED_STATE_ENABLED : PackageManager.COMPONENT_ENABLED_STATE_DISABLED;
        getActivity().getPackageManager().setComponentEnabledSetting(name,newState,PackageManager.DONT_KILL_APP);
    }

}
