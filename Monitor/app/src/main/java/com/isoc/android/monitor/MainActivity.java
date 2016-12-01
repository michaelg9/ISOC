package com.isoc.android.monitor;

import android.accounts.Account;
import android.accounts.AccountManager;
import android.app.Dialog;
import android.app.DialogFragment;
import android.content.Context;
import android.content.DialogInterface;
import android.os.Bundle;
import android.support.v7.app.AlertDialog;
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
            // Main fragment is shown on launch (after registering, if
            // necessary)
            getFragmentManager().beginTransaction()
                    .add(R.id.fragment_container, new MainFragment()).commit();
        }
    }

    @Override
    protected void onResume() {
        super.onResume();
        // Checking if there is a monitoring account registered first.
        // If not, prompt login screen.

        AccountManager am = AccountManager.get(this);
        if (am.getAccountsByType(getString(R.string.authenticator_account_type)).length == 0) {
            Log.e("no accounts", "adding");
            am.addAccount(getString(R.string.authenticator_account_type),
                    getString(R.string.token_refresh), null, null, this, null,
                    null);
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
        switch (item.getItemId()) {
        case R.id.action_settings:
            getFragmentManager().beginTransaction()
                    .replace(R.id.fragment_container, new SettingsFragment())
                    .addToBackStack(null).commit();
            return true;
        case android.R.id.home:
            if (getFragmentManager().getBackStackEntryCount() > 0)
                getFragmentManager().popBackStack();
            return true;
        case R.id.action_logout:
            new InvalidateTokenDialog().show(getFragmentManager(), "dialog");
            return true;
        default:
            break;
        }
        return super.onOptionsItemSelected(item);
    }

    //dialog used to verify that user wants to invalidate auth token
    public static class InvalidateTokenDialog extends DialogFragment {
        @Override
        public Dialog onCreateDialog(Bundle savedInstanceState) {
            AlertDialog.Builder builder = new AlertDialog.Builder(getActivity());
            builder.setTitle("Invalidate Auth Token?")
                    .setMessage(
                            "You will be asked to login again on the next sync.")
                    .setPositiveButton("OK",
                            new DialogInterface.OnClickListener() {
                                @Override
                                public void onClick(
                                        DialogInterface dialogInterface, int i) {
                                    logOut(getActivity());
                                }
                            }).setNegativeButton("Cancel", null);
            return builder.create();
        }

        // sends an asnychronous logout request to the server
        private void logOut(final Context context) {
            new Thread(new Runnable() {
                @Override
                public void run() {
                    AccountManager accountManager = AccountManager.get(context);
                    Account[] accounts = accountManager
                            .getAccountsByType(getString(R.string.authenticator_account_type));
                    for (Account a : accounts) {
                        accountManager.invalidateAuthToken(
                                getString(R.string.authenticator_account_type),
                                accountManager.peekAuthToken(a,
                                        getString(R.string.token_refresh)));
                    }
                }
            }).start();
        }
    }

}