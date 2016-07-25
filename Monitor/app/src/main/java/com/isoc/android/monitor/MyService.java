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

public class MyService extends IntentService {


    public MyService() {
        super("RecordService");
    }

    @Override
    protected void onHandleIntent(Intent intent) {
        SQLiteDatabase db = new Database(getApplicationContext()).getWritableDatabase();

        Log.e("RecordService","started: "+TimeCapture.getTime());
        SharedPreferences preferences = PreferenceManager.getDefaultSharedPreferences(this);
        if (preferences.getBoolean("calls",false)) ContactsCapture.getCallLog(this,db);

        String packagesPreference=preferences.getString("inst_pack","all");
        if (!packagesPreference.equals("none")) PackageCapture.getInstalledPackages(this,packagesPreference,db);

        if (preferences.getBoolean("run_service",false)) PackageCapture.getRunningServices(this,db);

        if (preferences.getBoolean("sockets",false)) SocketsCapture.getSockets(this,db);

        if (preferences.getBoolean("wifiScan",false)) NetworkCapture.getWifiAPs(this,db);

        Log.e("RecordService","ended: "+TimeCapture.getTime());

        db.close();
        ServiceReceiver.completeWakefulIntent(intent);
    }


    public static class ServiceControls{
        private ServiceControls(){}

        public static void startRepeated(Context context){
            Intent intent=new Intent(context,MyService.class);
            context.startService(intent);
            intent = new Intent(context,ServiceReceiver.class);
            final PendingIntent pendingIntent= PendingIntent.getBroadcast(context,1,intent,PendingIntent.FLAG_UPDATE_CURRENT);
            AlarmManager alarmManager=(AlarmManager) context.getSystemService(Context.ALARM_SERVICE);
            alarmManager.setRepeating(AlarmManager.ELAPSED_REALTIME_WAKEUP, SystemClock.elapsedRealtime(),300000,pendingIntent);
        }

        public static void stopRepeated(Context context){
            Intent intent=new Intent(context,ServiceReceiver.class);
            final PendingIntent pendingIntent= PendingIntent.getBroadcast(context,1,intent,PendingIntent.FLAG_UPDATE_CURRENT);
            AlarmManager alarmManager=(AlarmManager) context.getSystemService(Context.ALARM_SERVICE);
            alarmManager.cancel(pendingIntent);
            pendingIntent.cancel();
        }

        public static boolean checkExistence(Context context){
            Intent intent = new Intent(context,ServiceReceiver.class);
            PendingIntent pendingIntent= PendingIntent.getBroadcast(context,1,intent,PendingIntent.FLAG_NO_CREATE);
            return pendingIntent!=null;

        }
    }
}