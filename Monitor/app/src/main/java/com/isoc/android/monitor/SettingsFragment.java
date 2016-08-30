package com.isoc.android.monitor;


import android.content.ComponentName;
import android.content.SharedPreferences;
import android.content.pm.PackageManager;
import android.os.Bundle;
import android.preference.CheckBoxPreference;
import android.preference.EditTextPreference;
import android.preference.Preference;
import android.preference.PreferenceFragment;
import android.preference.SwitchPreference;
import android.support.v7.app.AppCompatActivity;
import android.support.v7.widget.Toolbar;
import android.view.LayoutInflater;
import android.view.Menu;
import android.view.MenuItem;
import android.view.View;
import android.view.ViewGroup;

/**
 * Settings fragment.
 */
public class SettingsFragment extends PreferenceFragment implements SharedPreferences.OnSharedPreferenceChangeListener {
    @Override
    public void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        //we have to change the options menu to hide the settings option
        setHasOptionsMenu(true);
    }

    @Override
    public void onResume() {
        super.onResume();
        //listener for preference changes
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
        //show the current value of the server url preference in the summary
        EditTextPreference url = (EditTextPreference) findPreference(getString(R.string.server_url));
        url.setSummary(url.getText());
        //change the title of the toolbar to settings
        Toolbar toolbar = (Toolbar) getActivity().findViewById(R.id.toolbar);
        toolbar.setTitle("Settings");
        //enabling back button
        if (((AppCompatActivity) getActivity()).getSupportActionBar() != null)
            ((AppCompatActivity) getActivity()).getSupportActionBar().setDisplayHomeAsUpEnabled(true);

        return view;
    }

    @Override
    public void onPrepareOptionsMenu(Menu menu) {
        MenuItem item = menu.findItem(R.id.action_settings);
        //hide the settings option since we're already here
        item.setVisible(false);
        super.onPrepareOptionsMenu(menu);
    }


    @Override
    public void onSharedPreferenceChanged(SharedPreferences sharedPreferences, String s) {
        Preference preference = findPreference(s);
        if (preference instanceof SwitchPreference) {
            //if the switch is changed, we disable / enable everything
            if (((SwitchPreference) preference).isChecked()) {
                MyService.ServiceControls.start(getActivity());
            } else {
                MyService.ServiceControls.stop(getActivity());
                //we need to disable receivers too
                ((CheckBoxPreference) findPreference(getString(R.string.battery_key))).setChecked(false);
                ((CheckBoxPreference) findPreference(getString(R.string.connectivity_key))).setChecked(false);
            }
        } else if (preference instanceof EditTextPreference) {
            //if an edit text preference changes, update the summary to the new value
            EditTextPreference editPreference = (EditTextPreference) preference;
            editPreference.setSummary(editPreference.getText());
        } else if (preference instanceof CheckBoxPreference) {
            //disable appropriate receivers
            switch (preference.getKey()) {
                case "battery":
                    ComponentName[] batteryReceiver = new ComponentName[]{new ComponentName(getActivity(), BatteryReceiver.class)};
                    if (((CheckBoxPreference) preference).isChecked())
                        setComponentState(true, batteryReceiver);
                    else setComponentState(false, batteryReceiver);
                    break;
                case "connectivity":
                    ComponentName[] networkReceiver = new ComponentName[]{new ComponentName(getActivity(), NetworkReceiver.class)};
                    if (((CheckBoxPreference) preference).isChecked())
                        setComponentState(true, networkReceiver);
                    else setComponentState(false, networkReceiver);
                    break;
                default:
                    break;
            }
        }
    }


    private void setComponentState(boolean state, ComponentName[] names) {
        int newState = (state) ? PackageManager.COMPONENT_ENABLED_STATE_ENABLED : PackageManager.COMPONENT_ENABLED_STATE_DISABLED;
        for (ComponentName name : names) {
            getActivity().getPackageManager().setComponentEnabledSetting(name, newState, PackageManager.DONT_KILL_APP);
        }
    }

}
