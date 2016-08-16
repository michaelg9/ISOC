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
Intent service launched on a specific interval from the wakeful broadcast receiver ServiceReceiver
Captures everything that isn't triggered by broadcasts.
Alarms are used to re-launch the service.
alarms are also used to send new data to the server
*/

public class MyService extends IntentService {

    public MyService() {
        super("RecordService");
    }

    @Override
    protected void onHandleIntent(Intent intent) {
        Log.e("RecordService","started: "+TimeCapture.getTime());
        if (intent.hasExtra(getString(R.string.sent_action))){
            XMLProduce p=new XMLProduce(this);
            p.new XMLSend(p.getXML()).sendXML();
            ServiceControls.scheduleSend(this);
        }else {
            SQLiteDatabase db = new Database(getApplicationContext()).getWritableDatabase();
            SharedPreferences preferences = PreferenceManager.getDefaultSharedPreferences(this);
            if (preferences.getBoolean(getString(R.string.calls_key), false))
                ContactsCapture.getCallLog(this, db);
            String packagesPreference = preferences.getString(getString(R.string.installed_packages_key), "none");
            if (!packagesPreference.equals("none"))
                PackageCapture.getInstalledPackages(this, packagesPreference, db);
            if (preferences.getBoolean(getString(R.string.running_services_key), false))
                PackageCapture.getRunningServices(this, db);
            if (preferences.getBoolean(getString(R.string.sockets_key), false))
                SocketsCapture.getSockets(db);
            if (preferences.getBoolean(getString(R.string.wifi_APs_key), false))
                NetworkCapture.getWifiAPs(this, db);
            if (preferences.getBoolean(getString(R.string.sms_key), false))
                SMSCapture.getSMS(this, db);
            db.close();
            ServiceControls.scheduleLaunch(this);
        }
        Log.e("RecordService","ended: "+TimeCapture.getTime());
        ServiceReceiver.completeWakefulIntent(intent);
    }


    public static class ServiceControls{
        private ServiceControls(){}
        private static int pendingLaunchCode=1;
        private static int pendingSendCode=2;

        public static void startRepeated(Context context){
            Intent i = new Intent(context,MyService.class);
            context.startService(i);
            scheduleLaunch(context);
            scheduleSend(context);
        }


        private static void scheduleLaunch(Context context){
            //sets an alarm to trigger according to the interval preference
            //Each alarm is set explicitly, there are no repeated alarms.
            //The reason is that I had problems with repeated alarms on android m.
            Intent i = new Intent(context.getApplicationContext(),ServiceReceiver.class);
            int time=PreferenceManager.getDefaultSharedPreferences(context).getInt(context.getString(R.string.capture_interval),context.getResources().getInteger(R.integer.timer_def));
            PendingIntent pendingIntent= PendingIntent.getBroadcast(context.getApplicationContext(),pendingLaunchCode,i,PendingIntent.FLAG_UPDATE_CURRENT);
            AlarmManager alarmManager=(AlarmManager) context.getSystemService(Context.ALARM_SERVICE);
            alarmManager.set(AlarmManager.ELAPSED_REALTIME_WAKEUP, SystemClock.elapsedRealtime()+time*60000,pendingIntent);
        }

        private static void scheduleSend(Context context){
            Intent i=new Intent(context.getApplicationContext(),MyService.class);
            i.putExtra(context.getString(R.string.sent_action),true);
            int time=PreferenceManager.getDefaultSharedPreferences(context).getInt(context.getString(R.string.sent_interval),context.getResources().getInteger(R.integer.sent_def));
            PendingIntent pendingIntent= PendingIntent.getService(context.getApplicationContext(),pendingSendCode,i,PendingIntent.FLAG_UPDATE_CURRENT);
            AlarmManager alarmManager=(AlarmManager) context.getSystemService(Context.ALARM_SERVICE);
            alarmManager.set(AlarmManager.ELAPSED_REALTIME, SystemClock.elapsedRealtime()+time*3600000,pendingIntent);
        }

        public static void stop(Context context){
            AlarmManager alarmManager=(AlarmManager) context.getSystemService(Context.ALARM_SERVICE);
            Intent intent=new Intent(context,ServiceReceiver.class);
            PendingIntent pendingIntent= PendingIntent.getBroadcast(context.getApplicationContext(),pendingLaunchCode,intent,PendingIntent.FLAG_UPDATE_CURRENT);
            alarmManager.cancel(pendingIntent);
            pendingIntent.cancel();
            pendingIntent=PendingIntent.getBroadcast(context.getApplicationContext(),pendingSendCode,intent,PendingIntent.FLAG_UPDATE_CURRENT);
            alarmManager.cancel(pendingIntent);
            pendingIntent.cancel();
        }
    }
}