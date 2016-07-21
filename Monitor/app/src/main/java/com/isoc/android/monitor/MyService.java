package com.isoc.android.monitor;

import android.app.AlarmManager;
import android.app.IntentService;
import android.app.PendingIntent;
import android.content.Context;
import android.content.Intent;
import android.content.SharedPreferences;
import android.os.SystemClock;
import android.preference.PreferenceManager;
import android.util.Log;

public class MyService extends IntentService {


    public MyService() {
        super("RecordService");
    }

    @Override
    protected void onHandleIntent(Intent intent) {
        Log.e("RecordService","started: "+TimeCapture.getTime());
        SharedPreferences preferences = PreferenceManager.getDefaultSharedPreferences(this);
        if (preferences.getBoolean("calls",true)) ContactsCapture.getCallLog(this);

        String packagesPreference=preferences.getString("inst_pack","all");
        if (!packagesPreference.equals("none")) PackageCapture.getInstalledPackages(this,packagesPreference);

        String servicesPreference=preferences.getString("run_service","all");
        if (!servicesPreference.equals("none")) PackageCapture.getRunningServices(this,servicesPreference);

        if (preferences.getBoolean("sockets",true)) SocketsCapture.getSockets(this);

        Log.e("RecordService","ended: "+TimeCapture.getTime());

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