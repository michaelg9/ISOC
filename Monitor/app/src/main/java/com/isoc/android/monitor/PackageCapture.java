package com.isoc.android.monitor;

import android.app.ActivityManager;
import android.content.ContentValues;
import android.content.Context;
import android.content.pm.ApplicationInfo;
import android.content.pm.PackageInfo;
import android.content.pm.PackageManager;
import android.database.sqlite.SQLiteDatabase;
import android.net.TrafficStats;

import java.util.List;


/**
 * TO DO:CPU time and memory of running services?
 * TO DO: indicate which apps provide public providers?
 * BUG: Some vendor / system apps report installed date close to 0...
 */
public class PackageCapture {

    protected static void getInstalledPackages(Context context, String pref, SQLiteDatabase db) {
        PackageManager packageManager = context.getPackageManager();
        int flags = PackageManager.GET_META_DATA;
        List<PackageInfo> packages = packageManager.getInstalledPackages(flags);

        for (int i = 0; i < packages.size(); i++) {
            if ((pref.equals("sys")) && (packages.get(i).applicationInfo.flags & ApplicationInfo.FLAG_SYSTEM) != ApplicationInfo.FLAG_SYSTEM)
                continue;
            else if ((pref.equals("usr")) && (packages.get(i).applicationInfo.flags & ApplicationInfo.FLAG_SYSTEM) == ApplicationInfo.FLAG_SYSTEM)
                continue;
            ContentValues values = new ContentValues();

            values.put(Database.DatabaseSchema.InstalledPackages.COLUMN_NAME_PACKAGE_NAME, packages.get(i).packageName);
            values.put(Database.DatabaseSchema.InstalledPackages.COLUMN_NAME_INSTALLED_DATE, TimeCapture.getTime(packages.get(i).firstInstallTime));
            values.put(Database.DatabaseSchema.InstalledPackages.COLUMN_NAME_VERSION, packages.get(i).versionName);
            values.put(Database.DatabaseSchema.InstalledPackages.COLUMN_NAME_UID, Integer.toString(packages.get(i).applicationInfo.uid));
            values.put(Database.DatabaseSchema.InstalledPackages.COLUMN_NAME_LABEL, packageManager.getApplicationLabel(packages.get(i).applicationInfo).toString());
            db.insertWithOnConflict(Database.DatabaseSchema.InstalledPackages.TABLE_NAME, null, values, SQLiteDatabase.CONFLICT_IGNORE);
        }
    }


    protected static void getRunningServices(Context context, SQLiteDatabase db) {
        ActivityManager am = (ActivityManager) context.getSystemService(Context.ACTIVITY_SERVICE);
        List<ActivityManager.RunningServiceInfo> runningServices = am.getRunningServices(Integer.MAX_VALUE);
        String time = TimeCapture.getTime();

        for (int i = 0; i < runningServices.size(); i++) {
            ContentValues values = new ContentValues();
            values.put(Database.DatabaseSchema.RunningServices.COLUMN_NAME_PACKAGE_NAME, runningServices.get(i).process);
            values.put(Database.DatabaseSchema.RunningServices.COLUMN_NAME_SINCE, TimeCapture.getTime(TimeCapture.getUpTime() + runningServices.get(i).activeSince));
            values.put(Database.DatabaseSchema.RunningServices.COLUMN_NAME_TIME, time);
            int uid = runningServices.get(i).uid;
            values.put(Database.DatabaseSchema.RunningServices.COLUMN_NAME_UID, Integer.toString(uid));
            long bytes = (TrafficStats.getUidRxBytes(uid) == -1) ? 0 : TrafficStats.getUidRxBytes(uid);
            values.put(Database.DatabaseSchema.RunningServices.COLUMN_NAME_RX, Long.toString(bytes));
            bytes = (TrafficStats.getUidTxBytes(uid) == -1) ? 0 : TrafficStats.getUidTxBytes(uid);
            values.put(Database.DatabaseSchema.RunningServices.COLUMN_NAME_TX, Long.toString(bytes));
            db.insertWithOnConflict(Database.DatabaseSchema.RunningServices.TABLE_NAME, null, values, SQLiteDatabase.CONFLICT_IGNORE);
        }
    }

}