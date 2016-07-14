package com.isoc.android.monitor;

import android.app.Service;
import android.content.Intent;
import android.content.SharedPreferences;
import android.os.Binder;
import android.os.IBinder;
import android.preference.PreferenceManager;
import android.widget.Toast;

public class MyService extends Service {
    private SharedPreferences prefs;

    private final IBinder mBinder = new LocalBinder();

    public class LocalBinder extends Binder {
        MyService getService() {
            return MyService.this;
        }
    }

    @Override
    public void onCreate() {
        super.onCreate();
        prefs= PreferenceManager.getDefaultSharedPreferences(getApplicationContext());
    }

    @Override
    public IBinder onBind(Intent intent) {
        return mBinder;
    }

    @Override
    public int onStartCommand(Intent intent, int flags, int startId) {
        Toast.makeText(this, "Service started", Toast.LENGTH_SHORT).show();
        ContactsCapture.getCallLog(this);
        NetworkCapture.getTrafficStats(this);
        PackageCapture.getInstalledPackages(this);
        PackageCapture.getRunningServices(this);
        return super.onStartCommand(intent, flags, startId);
    }

    @Override
    public boolean onUnbind(Intent intent) {
        return true;
    }

    @Override
    public void onRebind(Intent intent) {
        super.onRebind(intent);
    }

    @Override
    public void onDestroy() {
        super.onDestroy();
        Toast.makeText(this, "Service destroyed", Toast.LENGTH_SHORT).show();
    }


    public String generateXML() {

        StringBuilder result = new StringBuilder("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" +
                "<data>\n"+
                "<metadata>\n" +
                MetaDataCapture.getMetaDataXML(this)+
                "</metadata>\n" +
                "<device-data>\n");
        result.append(BatteryCapture.getBatteryXML(this,prefs));
        result.append(NetworkCapture.getTrafficXML(this));
//        result.append(ContactsCapture.getCallXML(this));
        //result.append(PackageCapture.getRunningServicesXML(this));
        result.append(PackageCapture.getInstalledPackagesXML(this));
        result.append("</device-data>\n</data>");
        return result.toString();

    }

}