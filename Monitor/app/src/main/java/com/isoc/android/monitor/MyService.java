package com.isoc.android.monitor;

import android.app.AlarmManager;
import android.app.IntentService;
import android.app.PendingIntent;
import android.content.Context;
import android.content.Intent;
import android.content.SharedPreferences;
import android.database.sqlite.SQLiteDatabase;
import android.os.SystemClock;
import android.preference.PreferenceManager;
import android.util.Log;

/*
 * Intent service launched on periodic intervals using 
 * the wakeful broadcast receiver ServiceReceiver
 * Captures everything that isn't triggered by broadcasts.
 * Alarms are used to re-launch the service. 
 * The alarms are explicitly set each time (rather than repeating alarms)
 * because I noticed that repeating alarms were removed after sleep mode
 */

public class MyService extends IntentService {

    public MyService() {
        super("RecordService");
    }

    @Override
    protected void onHandleIntent(Intent intent) {
        Log.e("RecordService", "started: " + TimeCapture.getCurrentStringTime());
        SQLiteDatabase db = new Database(getApplicationContext())
                .getWritableDatabase();
        SharedPreferences preferences = PreferenceManager
                .getDefaultSharedPreferences(this);
        // trigger all enabled parts
        if (preferences.getBoolean(getString(R.string.calls_key), false))
            ContactsCapture.getCallLog(this, db);
        String packagesPreference = preferences.getString(
                getString(R.string.installed_packages_key), "none");
        if (!packagesPreference.equals("none"))
            PackageCapture.getInstalledPackages(this, packagesPreference, db);
        if (preferences.getBoolean(getString(R.string.running_services_key),
                false))
            PackageCapture.getRunningServices(this, db);
        if (preferences.getBoolean(getString(R.string.sockets_key), false))
            SocketsCapture.getSockets(db);
        if (preferences.getBoolean(getString(R.string.wifi_APs_key), false))
            NetworkCapture.getWifiAPs(this, db);
        if (preferences.getBoolean(getString(R.string.sms_key), false))
            SMSCapture.getSMS(this, db);
        if (preferences.getBoolean(getString(R.string.accounts_key), false)) {
            AccountsCapture.getAccounts(this, db);
        }
        db.close();
        ServiceControls.scheduleLaunch(this);
        Log.e("RecordService", "ended: " + TimeCapture.getCurrentStringTime());
        ServiceReceiver.completeWakefulIntent(intent);
    }

    //Service controls are used to turn on / off the service
    //and schedule relaunches
    public static class ServiceControls {
        private ServiceControls() {
        }

        private static int pendingLaunchCode = 1;

        public static void start(Context context) {
            Intent i = new Intent(context, MyService.class);
            context.startService(i);
        }

        private static void scheduleLaunch(Context context) {
            // sets an alarm to trigger according to the interval preference
            // Each alarm is set explicitly, there are no repeated alarms.
            // The reason is that I had problems with repeated alarms on android
            // m.
            Intent i = new Intent(context.getApplicationContext(),
                    ServiceReceiver.class);
            int time = PreferenceManager.getDefaultSharedPreferences(context)
                    .getInt(context.getString(R.string.capture_interval),
                            context.getResources().getInteger(
                                    R.integer.timer_def));
            PendingIntent pendingIntent = PendingIntent.getBroadcast(
                    context.getApplicationContext(), pendingLaunchCode, i,
                    PendingIntent.FLAG_UPDATE_CURRENT);
            AlarmManager alarmManager = (AlarmManager) context
                    .getSystemService(Context.ALARM_SERVICE);
            alarmManager
                    .set(AlarmManager.ELAPSED_REALTIME_WAKEUP,
                            SystemClock.elapsedRealtime() + time * 60000,
                            pendingIntent);
        }

        public static void stop(Context context) {
            AlarmManager alarmManager = (AlarmManager) context
                    .getSystemService(Context.ALARM_SERVICE);
            Intent intent = new Intent(context, ServiceReceiver.class);
            PendingIntent pendingIntent = PendingIntent.getBroadcast(
                    context.getApplicationContext(), pendingLaunchCode, intent,
                    PendingIntent.FLAG_UPDATE_CURRENT);
            alarmManager.cancel(pendingIntent);
            pendingIntent.cancel();
        }
    }
}