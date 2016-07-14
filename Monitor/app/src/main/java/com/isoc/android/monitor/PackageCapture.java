package com.isoc.android.monitor;

import android.app.ActivityManager;
import android.content.ContentValues;
import android.content.Context;
import android.content.SharedPreferences;
import android.content.pm.ApplicationInfo;
import android.content.pm.PackageInfo;
import android.content.pm.PackageManager;
import android.database.Cursor;
import android.database.sqlite.SQLiteDatabase;
import android.net.TrafficStats;
import android.os.SystemClock;

import java.util.List;


/**
 * Created by maik on 1/7/2016.
 * CPU time and memory of running apps?
 */
public class PackageCapture {

    protected static void getInstalledPackages(Context context){
        PackageManager packageManager= context.getPackageManager();
        int flags = PackageManager.GET_META_DATA;
        List<PackageInfo> packages = packageManager.getInstalledPackages(flags);
        SQLiteDatabase db=new Database(context).getWritableDatabase();

        for (int i =0; i<packages.size();i++){
            if ((packages.get(i).applicationInfo.flags & ApplicationInfo.FLAG_SYSTEM) !=1) continue;
            ContentValues values=new ContentValues();
            values.put(Database.DatabaseSchema.InstalledPackages.COLUMN_NAME_PACKAGE_NAME,packages.get(i).packageName);
            values.put(Database.DatabaseSchema.InstalledPackages.COLUMN_NAME_INSTALLED_DATE,packages.get(i).firstInstallTime);
            values.put(Database.DatabaseSchema.InstalledPackages.COLUMN_NAME_VERSION,packages.get(i).versionName);
            values.put(Database.DatabaseSchema.InstalledPackages.COLUMN_NAME_UID,Integer.toString(packages.get(i).applicationInfo.uid));
            values.put(Database.DatabaseSchema.InstalledPackages.COLUMN_NAME_LABEL,packageManager.getApplicationLabel(packages.get(i).applicationInfo).toString());
            db.insertWithOnConflict(Database.DatabaseSchema.InstalledPackages.TABLE_NAME,null,values,SQLiteDatabase.CONFLICT_IGNORE);

        }
        db.close();
    }

    protected static void getRunningServices(Context context){
        ActivityManager am=(ActivityManager) context.getSystemService(Context.ACTIVITY_SERVICE);
        List<ActivityManager.RunningServiceInfo> runningApps = am.getRunningServices(Integer.MAX_VALUE);
        SQLiteDatabase db=new Database(context).getWritableDatabase();

        for (int i =0; i<runningApps.size();i++){
            ContentValues values=new ContentValues();
            values.put(Database.DatabaseSchema.RunningServices.COLUMN_NAME_PROCESS_NAME,runningApps.get(i).process);
            values.put(Database.DatabaseSchema.RunningServices.COLUMN_NAME_UP_TIME,(SystemClock.elapsedRealtime()-runningApps.get(i).activeSince)/1000);
            int uid=runningApps.get(i).uid;
            values.put(Database.DatabaseSchema.RunningServices.COLUMN_NAME_UID,Integer.toString(uid));
            values.put(Database.DatabaseSchema.RunningServices.COLUMN_NAME_RX,Long.toString(TrafficStats.getUidRxBytes(uid)));
            values.put(Database.DatabaseSchema.RunningServices.COLUMN_NAME_TX,Long.toString(TrafficStats.getUidTxBytes(uid)));
            db.insertWithOnConflict(Database.DatabaseSchema.RunningServices.TABLE_NAME,null,values,SQLiteDatabase.CONFLICT_IGNORE);
        }
        db.close();
    }

    protected static String getRunningServicesXML(Context context, SharedPreferences prefs) {
        SQLiteDatabase db=new Database(context).getReadableDatabase();
        Cursor cursor = db.query(Database.DatabaseSchema.RunningServices.TABLE_NAME,null,null,null,null,null,null);
        StringBuilder result=new StringBuilder();
        int uid = cursor.getColumnIndex(Database.DatabaseSchema.RunningServices.COLUMN_NAME_UID);
        int uptime = cursor.getColumnIndex(Database.DatabaseSchema.RunningServices.COLUMN_NAME_UP_TIME);
        int rx = cursor.getColumnIndex(Database.DatabaseSchema.RunningServices.COLUMN_NAME_RX);
        int tx = cursor.getColumnIndex(Database.DatabaseSchema.RunningServices.COLUMN_NAME_TX);
        int name = cursor.getColumnIndex(Database.DatabaseSchema.RunningServices.COLUMN_NAME_PROCESS_NAME);

        while (cursor.moveToNext()){
            result.append("<runservice uid=\""+cursor.getString(uid)+"\" uptime=\"" + cursor.getString(uptime) + "\" rx=\"" +
                    cursor.getString(rx) + "\" tx=\"" +
                    cursor.getString(tx) +"\">" +  cursor.getString(name) + "</runservice>\n");
        }
        cursor.close();
        db.close();
        return result.toString();
    }


    protected static String getInstalledPackagesXML(Context context) {
        StringBuilder result=new StringBuilder();
        SQLiteDatabase db=new Database(context).getReadableDatabase();
        Cursor cursor = db.query(Database.DatabaseSchema.InstalledPackages.TABLE_NAME,null,null,null,null,null,null);
        int uid = cursor.getColumnIndex(Database.DatabaseSchema.InstalledPackages.COLUMN_NAME_UID);
        int label = cursor.getColumnIndex(Database.DatabaseSchema.InstalledPackages.COLUMN_NAME_LABEL);
        int version = cursor.getColumnIndex(Database.DatabaseSchema.InstalledPackages.COLUMN_NAME_VERSION);
        int date = cursor.getColumnIndex(Database.DatabaseSchema.InstalledPackages.COLUMN_NAME_INSTALLED_DATE);
        int name = cursor.getColumnIndex(Database.DatabaseSchema.InstalledPackages.COLUMN_NAME_PACKAGE_NAME);

        while (cursor.moveToNext()){
            result.append("<installedapp name=\"" + cursor.getString(name) + "\" installed=\"" +
                    TimeCapture.getTime(cursor.getLong(date)) + "\" version=\"" + cursor.getString(version) + "\" uid=\"" +
                    cursor.getString(uid) +"\">" + cursor.getString(label) + "</installedapp>\n");
        }
        db.close();
        cursor.close();
        return result.toString();
    }
}