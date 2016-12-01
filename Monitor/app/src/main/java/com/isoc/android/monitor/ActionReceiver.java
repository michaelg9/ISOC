package com.isoc.android.monitor;

import android.content.BroadcastReceiver;
import android.content.ContentValues;
import android.content.Context;
import android.content.Intent;
import android.content.SharedPreferences;
import android.database.sqlite.SQLiteDatabase;
import android.net.ConnectivityManager;
import android.preference.PreferenceManager;

/* 
 * ActionReceiver is triggered when an action like shutdown / boot / airplane mode is performed
 * This receiver isn't disabled when the preference is disabled, because it's also used
 * to persist interface data and start the service on reboot and to save interface data on shutdown
 */
public class ActionReceiver extends BroadcastReceiver {
    public ActionReceiver() {}

    @Override
    public void onReceive(Context context, Intent intent) {
        String action = intent.getAction();
        SharedPreferences prefs = PreferenceManager
                .getDefaultSharedPreferences(context);
        //checking user preferences
        boolean actionsEnabled = (prefs.getBoolean(
                context.getString(R.string.actions_key), false));
        boolean connectivityEnabled = prefs.getBoolean(
                context.getString(R.string.connectivity_key), false);

        if (action.equals("android.intent.action.BOOT_COMPLETED")) {
            if (connectivityEnabled){
                // Even if actions checkbox is disabled, we need to save the current
                // interface rx/tx stats and restart the alarm (if monitoring is
                // enabled)
                NetworkCapture.saveCurrentStats(context);
            
            }
            //if monitoring is enabled, start the service after boot
            if (prefs.getBoolean("monitoring", false)) {
                MyService.ServiceControls.start(context);
            }
            //if actions monitoring enabled, record the event
            if (actionsEnabled)
                ActionCapture.getAction(context,
                        Database.DatabaseSchema.Actions.ACTION_BOOT);
        } else if (action.equals("android.intent.action.ACTION_SHUTDOWN")
                || action.equals("android.intent.action.QUICKBOOT_POWEROFF")
                || action.equals("android.intent.action.ACTION_REBOOT")) {
            // saving current tx and rx of last active interface (or none if no
            // connectivity exists while shutting down)
            if (connectivityEnabled)
                NetworkCapture
                        .getTrafficStats(
                                context,
                                ((ConnectivityManager) context
                                        .getSystemService(Context.CONNECTIVITY_SERVICE))
                                        .getActiveNetworkInfo());
            if (actionsEnabled)
                ActionCapture.getAction(context,
                        Database.DatabaseSchema.Actions.ACTION_SHUTDOWN);
        } else if (action.equals("android.intent.action.AIRPLANE_MODE")
                && actionsEnabled) {
            if (intent.getBooleanExtra("state", true))
                ActionCapture.getAction(context,
                        Database.DatabaseSchema.Actions.ACTION_AIRPLANE_ON);
            else
                ActionCapture.getAction(context,
                        Database.DatabaseSchema.Actions.ACTION_AIRPLANE_OFF);
        }
    }

    /**
     * ActionReceiver calls ActionCapture to record actions
     * Actions are captured only when they are broadcasted, not on intervals.
     * Look at the database schema for the captured actions.
     */
    private static class ActionCapture {

        public static void getAction(Context context, String action) {
            String date = TimeCapture.getCurrentStringTime();
            ContentValues values = new ContentValues();
            values.put(Database.DatabaseSchema.Actions.COLUMN_NAME_ACTION,
                    action);
            values.put(Database.DatabaseSchema.Actions.COLUMN_NAME_DATE, date);
            SQLiteDatabase db = new Database(context).getWritableDatabase();
            db.insert(Database.DatabaseSchema.Actions.TABLE_NAME, null, values);
            db.close();
        }
    }
}