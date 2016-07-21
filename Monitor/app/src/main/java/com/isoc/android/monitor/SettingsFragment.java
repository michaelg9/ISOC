package com.isoc.android.monitor;


import android.content.ComponentName;
import android.content.SharedPreferences;
import android.content.pm.PackageManager;
import android.os.Bundle;
import android.preference.CheckBoxPreference;
import android.preference.EditTextPreference;
import android.preference.Preference;
import android.preference.PreferenceFragment;
import android.preference.PreferenceScreen;
import android.preference.SwitchPreference;
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
    public boolean onPreferenceTreeClick(PreferenceScreen preferenceScreen, Preference preference) {
        super.onPreferenceTreeClick(preferenceScreen, preference);
        if (preference != null && preference instanceof PreferenceScreen) {
            if (((PreferenceScreen) preference).getDialog() != null) {
                //        preferenceScreen.getDialog().getActionBar().hide();
                //getWindow().getDecorView().setBackground( getActivity().getWindow().getDecorView().getBackground().getConstantState().newDrawable());
            }
        }
        return false;
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
        View view = super.onCreateView(inflater, container, savedInstanceState);
        addPreferencesFromResource(R.xml.preferences);
        EditTextPreference url = (EditTextPreference) findPreference("server_url");
        url.setSummary(url.getText());
        Toolbar toolbar = (Toolbar) getActivity().findViewById(R.id.toolbar);
        toolbar.setTitle("Settings");
        return view;
    }

    @Override
    public void onPrepareOptionsMenu(Menu menu) {
        MenuItem item = menu.findItem(R.id.action_settings);
        item.setVisible(false);
        super.onPrepareOptionsMenu(menu);
    }


    @Override
    public void onSharedPreferenceChanged(SharedPreferences sharedPreferences, String s) {
        Preference preference = findPreference(s);
        if (preference instanceof SwitchPreference){
            if (((SwitchPreference)preference).isChecked()) MyService.ServiceControls.startRepeated(getActivity());
            else MyService.ServiceControls.stopRepeated(getActivity());
        }else if (preference instanceof EditTextPreference) {
            EditTextPreference editPreference = (EditTextPreference) preference;
            editPreference.setSummary(editPreference.getText());
        } else if (preference instanceof CheckBoxPreference) {
            switch (preference.getKey()) {
                case "battery":
                    ComponentName batteryReceiver = new ComponentName(getActivity(), BatteryReceiver.class);
                    if (((CheckBoxPreference) preference).isChecked())
                        setComponentState(true, batteryReceiver);
                    else setComponentState(false, batteryReceiver);
                    break;
                case "connectivity":
                    ComponentName networkReceiver = new ComponentName(getActivity(), NetworkReceiver.class);
                    if (((CheckBoxPreference) preference).isChecked())
                        setComponentState(true, networkReceiver);
                    else setComponentState(false, networkReceiver);
                    break;
                default:
                    break;
            }
        }
    }


    private void setComponentState(boolean state, ComponentName name) {
        int newState = (state) ? PackageManager.COMPONENT_ENABLED_STATE_ENABLED : PackageManager.COMPONENT_ENABLED_STATE_DISABLED;
        getActivity().getPackageManager().setComponentEnabledSetting(name, newState, PackageManager.DONT_KILL_APP);
    }

}
