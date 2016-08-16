package com.isoc.android.monitor;

import android.accounts.AccountManager;
import android.os.Bundle;
import android.support.v7.app.AppCompatActivity;
import android.support.v7.widget.Toolbar;
import android.util.Log;
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
            //Main fragment is shown on launch (after registering, if necessary)
            getFragmentManager().beginTransaction().add(R.id.fragment_container,new MainFragment()).commit();
        }
    }

    @Override
    protected void onResume() {
        super.onResume();
        //Checking if there is a monitoring account registered first.
        //If not, prompt login screen.
        AccountManager am=AccountManager.get(this);
        if (am.getAccountsByType(getString(R.string.authenticator_account_type)).length==0){
            Log.e("no accounts","adding");
            am.addAccount(getString(R.string.authenticator_account_type),getString(R.string.token_refresh),null,null,this,null,null);
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