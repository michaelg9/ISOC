package com.isoc.android.monitor;

import android.content.BroadcastReceiver;
import android.content.ContentValues;
import android.content.Context;
import android.content.Intent;
import android.content.IntentFilter;
import android.database.sqlite.SQLiteDatabase;
import android.os.BatteryManager;
/*Broadcast receiver for battery events.
* Triggered when power is plugged or unplugged, battery is low or battery recovers from low.
* This receiver is disabled if the user disables battery capturing from preferences
 */
public class BatteryReceiver extends BroadcastReceiver {
    public BatteryReceiver() {
    }

    @Override
    public void onReceive(Context context, Intent intent) {
        new BatteryCapture().getBatteryStats(context);
    }

    /**
     * Captures battery details, after being triggered by BatteryReceiver. Not captured on intervals
     */
    private class BatteryCapture {
        public void getBatteryStats(Context context) {
            IntentFilter iFilter = new IntentFilter(Intent.ACTION_BATTERY_CHANGED);
            Intent battery = context.getApplicationContext().registerReceiver(null, iFilter); //capturing a sticky broadcast doesn't need a receiver
            if (battery==null) return;
            ContentValues values=new ContentValues();
            String charging=resolveChargingMode(battery.getIntExtra(BatteryManager.EXTRA_PLUGGED, -1));
            values.put(Database.DatabaseSchema.Battery.COLUMN_NAME_TEMP,((float) battery.getIntExtra(BatteryManager.EXTRA_TEMPERATURE, -1))/10);
            values.put(Database.DatabaseSchema.Battery.COLUMN_NAME_CHARGING,charging);
            values.put(Database.DatabaseSchema.Battery.COLUMN_NAME_LEVEL,battery.getIntExtra(BatteryManager.EXTRA_LEVEL, -1));
            values.put(Database.DatabaseSchema.Battery.COLUMN_NAME_TIME,TimeCapture.getCurrentStringTime());
            SQLiteDatabase db= new Database(context).getWritableDatabase();
            db.insert(Database.DatabaseSchema.Battery.TABLE_NAME,null,values);
            db.close();
        }

        private String resolveChargingMode(int mode){
            String charging;
            switch (mode){
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
            return charging;
        }

    }
}
