package com.isoc.android.monitor;

import android.content.BroadcastReceiver;
import android.content.Context;
import android.content.Intent;
import android.preference.PreferenceManager;
import android.util.Log;

public class ActionReceiver extends BroadcastReceiver {
    public ActionReceiver() {
    }

    @Override
    public void onReceive(Context context, Intent intent) {
        String action = intent.getAction();
        if (!PreferenceManager.getDefaultSharedPreferences(context).getBoolean("actions",true)){
            if (action.equals("android.intent.action.BOOT_COMPLETED")) NetworkCapture.saveCurrentStats(context);
            return;
        }
        Log.e("CAPTURE-ACTION",intent.toString());

        switch (action) {
            case "android.intent.action.BOOT_COMPLETED":
                NetworkCapture.saveCurrentStats(context);
                ActionCapture.getAction(context, Database.DatabaseSchema.Actions.ACTION_BOOT);
                break;
            case "android.intent.action.ACTION_SHUTDOWN":
                ActionCapture.getAction(context, Database.DatabaseSchema.Actions.ACTION_SHUTDOWN);
                break;
            case "android.intent.action.ACTION_REBOOT":
                ActionCapture.getAction(context, Database.DatabaseSchema.Actions.ACTION_REBOOT);
                break;
            case "android.intent.action.AIRPLANE_MODE":
                if (intent.getBooleanExtra("state",true)) ActionCapture.getAction(context, Database.DatabaseSchema.Actions.ACTION_AIRPLANE_ON);
                else ActionCapture.getAction(context, Database.DatabaseSchema.Actions.ACTION_AIRPLANE_OFF);
                break;
            default:
                return;
        }
    }
}
