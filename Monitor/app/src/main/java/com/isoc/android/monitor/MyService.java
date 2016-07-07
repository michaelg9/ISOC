package com.isoc.android.monitor;

import android.app.Service;
import android.content.Intent;
import android.os.Binder;
import android.os.IBinder;
import android.telephony.TelephonyManager;
import android.util.Log;
import android.widget.Toast;

public class MyService extends Service {
    private final IBinder mBinder = new LocalBinder();
    private final String timeFormat="yyyy-MM-dd HH:mm:ss";

    public class LocalBinder extends Binder {
        MyService getService() {
            return MyService.this;
        }
    }

    @Override
    public IBinder onBind(Intent intent) {
        return mBinder;
    }

    @Override
    public int onStartCommand(Intent intent, int flags, int startId) {
        Toast.makeText(this, "Service started", Toast.LENGTH_SHORT).show();
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
        result.append(BatteryCapture.getBatteryXML(this,timeFormat));
        result.append(NetworkCapture.getTrafficXML(this,timeFormat));
        result.append(ContactsCapture.getCallXML(this,timeFormat));
        result.append(PackageCapture.getRunningServicesXML(this,timeFormat));
        result.append(PackageCapture.getInstalledPackagesXML(this,timeFormat));
        result.append("</device-data>\n</data>");
        return result.toString();

    }

}