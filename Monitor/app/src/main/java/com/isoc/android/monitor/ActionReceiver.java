package com.isoc.android.monitor;

import android.content.BroadcastReceiver;
import android.content.Context;
import android.content.Intent;
import android.content.SharedPreferences;
import android.net.ConnectivityManager;
import android.preference.PreferenceManager;

public class ActionReceiver extends BroadcastReceiver {
    public ActionReceiver() {
    }

    @Override
    public void onReceive(Context context, Intent intent) {
        String action = intent.getAction();
        SharedPreferences prefs = PreferenceManager.getDefaultSharedPreferences(context);
        boolean actionsEnabled=(prefs.getBoolean(context.getString(R.string.actions_key), false));
        boolean connectivityEnabled=prefs.getBoolean(context.getString(R.string.connectivity_key),false);
        //Even if actions checkbox is disabled, we need to save the current interface rx/tx stats and restart the alarm (if monitoring is enabled)
        if (action.equals("android.intent.action.BOOT_COMPLETED")) {
            if (connectivityEnabled) NetworkCapture.saveCurrentStats(context);
            if (prefs.getBoolean("monitoring", false)) {
                MyService.ServiceControls.startRepeated(context);
            }
            if (actionsEnabled)
                ActionCapture.getAction(context, Database.DatabaseSchema.Actions.ACTION_BOOT);
        }else if (action.equals("android.intent.action.ACTION_SHUTDOWN") || action.equals("android.intent.action.QUICKBOOT_POWEROFF") || action.equals("android.intent.action.ACTION_REBOOT")) {
            //saving current tx and rx of last active interface (or none if no connectivity exists while shutting down)
            if (connectivityEnabled)
                NetworkCapture.getTrafficStats(context, ((ConnectivityManager) context.getSystemService(Context.CONNECTIVITY_SERVICE)).getActiveNetworkInfo());
            if (actionsEnabled)
                ActionCapture.getAction(context, Database.DatabaseSchema.Actions.ACTION_SHUTDOWN);
        }else if (action.equals("android.intent.action.AIRPLANE_MODE") && actionsEnabled) {
            if (intent.getBooleanExtra("state", true))
                ActionCapture.getAction(context, Database.DatabaseSchema.Actions.ACTION_AIRPLANE_ON);
            else
                ActionCapture.getAction(context, Database.DatabaseSchema.Actions.ACTION_AIRPLANE_OFF);
        }
    }
}