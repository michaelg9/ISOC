package com.isoc.android.monitor;

import android.os.Bundle;
import android.support.v7.app.AppCompatActivity;
import android.support.v7.widget.Toolbar;
import android.view.Menu;
import android.view.MenuItem;
//empty activity that holds the active fragment

public class MainActivity extends AppCompatActivity {


    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);
        Toolbar toolbar = (Toolbar) findViewById(R.id.toolbar);
        setSupportActionBar(toolbar);
        if (findViewById(R.id.fragment_container) != null) {
            if (savedInstanceState != null) {
                return;
            }
            //Main fragment is shown on launch
            getFragmentManager().beginTransaction().add(R.id.fragment_container,new MainFragment()).commit();
        }
    }

    @Override
    public boolean onCreateOptionsMenu(Menu menu) {
        // Inflate the menu; this adds items to the action bar if it is present.
        getMenuInflater().inflate(R.menu.menu_main, menu);
        return true;
    }


    @Override
    public boolean onOptionsItemSelected(MenuItem item) {
        // Handle action bar item clicks here. The action bar will
        // automatically handle clicks on the Home/Up button, so long
        // as you specify a parent activity in AndroidManifest.xml.
        switch (item.getItemId()){
            case R.id.action_settings:
                getFragmentManager().beginTransaction().replace(R.id.fragment_container,new SettingsFragment()).addToBackStack(null).commit();
                return true;
            case android.R.id.home:
                if (getFragmentManager().getBackStackEntryCount()>0) getFragmentManager().popBackStack();
                break;
            default:
                break;
        }
        return super.onOptionsItemSelected(item);
    }
}