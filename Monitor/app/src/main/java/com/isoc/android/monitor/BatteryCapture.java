package com.isoc.android.monitor;

import android.content.ContentValues;
import android.content.Context;
import android.content.Intent;
import android.content.IntentFilter;
import android.content.SharedPreferences;
import android.database.Cursor;
import android.database.sqlite.SQLiteDatabase;
import android.os.BatteryManager;

/**
 * Created by maik on 1/7/2016.
 * Capture batttery statistics
 */
public class BatteryCapture {

    public static void getBatteryStats(Context context) {
        IntentFilter iFilter = new IntentFilter(Intent.ACTION_BATTERY_CHANGED);
        Intent battery = context.getApplicationContext().registerReceiver(null, iFilter);
        if (battery==null) return;
        ContentValues values=new ContentValues();
        values.put(Database.DatabaseSchema.Battery.COLUMN_NAME_CHARGING,battery.getIntExtra(BatteryManager.EXTRA_PLUGGED, -1));
        values.put(Database.DatabaseSchema.Battery.COLUMN_NAME_LEVEL,battery.getIntExtra(BatteryManager.EXTRA_LEVEL, -1));
        values.put(Database.DatabaseSchema.Battery.COLUMN_NAME_DATE,TimeCapture.getTime());
        SQLiteDatabase db= new Database(context).getWritableDatabase();
        db.insert(Database.DatabaseSchema.Battery.TABLE_NAME,null,values);
        db.close();
    }

    protected static String getBatteryXML(Context context, SharedPreferences prefs) {
        SQLiteDatabase db=new Database(context).getReadableDatabase();
        Cursor cursor = db.query(Database.DatabaseSchema.Battery.TABLE_NAME,null,null,null,null,null,null);
        StringBuilder sb=new StringBuilder();
        int dateIndex = cursor.getColumnIndex(Database.DatabaseSchema.Battery.COLUMN_NAME_DATE);
        int chargeIndex = cursor.getColumnIndex(Database.DatabaseSchema.Battery.COLUMN_NAME_CHARGING);
        int levelIndex = cursor.getColumnIndex(Database.DatabaseSchema.Battery.COLUMN_NAME_LEVEL);
        while (cursor.moveToNext()) {
            boolean charging =cursor.getInt(chargeIndex)!=0;
            sb.append("<battery time=\"" + cursor.getString(dateIndex) + "\" charging=\"" + charging +
                    "\">" + cursor.getInt(levelIndex) + "</battery>\n");
        }
        db.close();
        cursor.close();
        return sb.toString();
    }

}