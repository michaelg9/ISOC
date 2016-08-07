package com.isoc.android.monitor;

import android.content.ContentValues;
import android.content.Context;
import android.content.Intent;
import android.content.IntentFilter;
import android.database.sqlite.SQLiteDatabase;
import android.os.BatteryManager;

/**
 * Captures battery details when any of these actions are broadcasted: plugged, unplugged,
 * battery low and battery recovers from low. Not captured on intervals
 */
public class BatteryCapture {

    public static void getBatteryStats(Context context) {
        IntentFilter iFilter = new IntentFilter(Intent.ACTION_BATTERY_CHANGED);
        Intent battery = context.getApplicationContext().registerReceiver(null, iFilter); //capturing a sticky broadcast doesn't need a receiver
        if (battery==null) return;
        ContentValues values=new ContentValues();
        String charging;
        switch (battery.getIntExtra(BatteryManager.EXTRA_PLUGGED, -1)){
            case 0:
                charging="no";
                break;
            case BatteryManager.BATTERY_PLUGGED_AC:
                charging="ac";
                break;
            case BatteryManager.BATTERY_PLUGGED_USB:
                charging="usb";
                break;
            case BatteryManager.BATTERY_PLUGGED_WIRELESS:
                charging="wireless";
                break;
            default:
                charging="unknown";
                break;
        }
        values.put(Database.DatabaseSchema.Battery.COLUMN_NAME_TEMP,((float) battery.getIntExtra(BatteryManager.EXTRA_TEMPERATURE, -1))/10);
        values.put(Database.DatabaseSchema.Battery.COLUMN_NAME_CHARGING,charging);
        values.put(Database.DatabaseSchema.Battery.COLUMN_NAME_LEVEL,battery.getIntExtra(BatteryManager.EXTRA_LEVEL, -1));
        values.put(Database.DatabaseSchema.Battery.COLUMN_NAME_TIME,TimeCapture.getTime());
        SQLiteDatabase db= new Database(context).getWritableDatabase();
        db.insert(Database.DatabaseSchema.Battery.TABLE_NAME,null,values);
        db.close();
    }

}