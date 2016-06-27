package com.isoc.android.monitor;

import android.Manifest;
import android.app.Service;
import android.content.Intent;
import android.content.IntentFilter;
import android.content.pm.PackageManager;
import android.database.Cursor;
import android.os.BatteryManager;
import android.os.Binder;
import android.os.IBinder;
import android.provider.CallLog;
import android.support.annotation.Nullable;
import android.support.v4.app.ActivityCompat;
import android.widget.Toast;

import java.text.SimpleDateFormat;
import java.util.Calendar;
import java.util.Date;

public class MyService extends Service {
    private final IBinder mBinder = new LocalBinder();

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

    private int getBatteryLevel() {
        IntentFilter iFilter = new IntentFilter(Intent.ACTION_BATTERY_CHANGED);
        Intent battery = this.registerReceiver(null, iFilter);
        return battery.getIntExtra(BatteryManager.EXTRA_LEVEL, -1);
    }

    private String getTime(String format) {
        Calendar c = Calendar.getInstance();
        SimpleDateFormat sdf = new SimpleDateFormat(format);
        return sdf.format(c.getTime());
    }

    @Nullable
    private String[][] getCallLog() {
        if (ActivityCompat.checkSelfPermission(this, Manifest.permission.READ_CALL_LOG) != PackageManager.PERMISSION_GRANTED) {
            return null;
        }
        Cursor cursor = getApplicationContext().getContentResolver().query(CallLog.Calls.CONTENT_URI, null, null, null, CallLog.Calls.DATE+ " DESC");
        int number = cursor.getColumnIndex(CallLog.Calls.NUMBER);
        int type = cursor.getColumnIndex(CallLog.Calls.TYPE);
        int duration = cursor.getColumnIndex(CallLog.Calls.DURATION);
        int date = cursor.getColumnIndex(CallLog.Calls.DATE);

        String[][] result = new String[cursor.getCount()][4];
        int i = 0;
        while (cursor.moveToNext()) {
            result[i][0] = cursor.getString(number);
            result[i][1] = cursor.getString(duration);
            result[i][3] = cursor.getString(date);

            String callType = new String();
            switch (Integer.parseInt(cursor.getString(type))) {
                case CallLog.Calls.OUTGOING_TYPE:
                    callType = "Outgoing";
                    break;
                case CallLog.Calls.INCOMING_TYPE:
                    callType = "Incoming";
                    break;
                case CallLog.Calls.MISSED_TYPE:
                    callType = "Missed";
                    break;
                default:
                    callType = "Unknown";
                    break;
            }
            result[i][2] = callType;
            i++;
        }
        cursor.close();
        return result;
    }

    public String generateXML() {

        StringBuilder result = new StringBuilder("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" +
                "<data>\n"+
                "<metadata>\n" +
                "<device>1</device>\n" +
                "</metadata>\n" +
                "<device-data>\n");
        result.append("<battery time=\"" + getTime("yyyy-MM-dd HH:mm:ss") + "\">" + getBatteryLevel() + "</battery>\n");

        String[][] callLog = getCallLog();
        if (callLog != null)
            for (String[] call : callLog) {
                long seconds=Long.parseLong(call[3]);
                SimpleDateFormat formatter = new SimpleDateFormat("yyyy-MM-dd HH:mm:ss");
                String time=formatter.format(new Date(seconds));
                result.append("<call time=\""+time +"\" type=\"" + call[2] + "\" duration=\"" + call[1]+"\">"+ call[0]+"</call>\n");
            }

        result.append("</device-data>\n</data>");
        return result.toString();
    }
}
