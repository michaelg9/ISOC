package com.isoc.android.monitor;

import android.app.ActivityManager;
import android.content.ContentValues;
import android.content.Context;
import android.content.pm.ApplicationInfo;
import android.content.pm.PackageInfo;
import android.content.pm.PackageManager;
import android.content.pm.ProviderInfo;
import android.database.Cursor;
import android.database.sqlite.SQLiteDatabase;
import android.net.TrafficStats;
import android.util.Log;

import java.util.List;


/**
 * Created by maik on 1/7/2016.
 * CPU time and memory of running apps?
 */
public class PackageCapture {

    protected static void getInstalledPackages(Context context, String pref, SQLiteDatabase db) {
        PackageManager packageManager = context.getPackageManager();
        int flags = PackageManager.GET_META_DATA | PackageManager.GET_PROVIDERS;
        List<PackageInfo> packages = packageManager.getInstalledPackages(flags);

        for (int i = 0; i < packages.size(); i++) {
            if ((pref.equals("sys")) && (packages.get(i).applicationInfo.flags & ApplicationInfo.FLAG_SYSTEM) != ApplicationInfo.FLAG_SYSTEM)
                continue;
            else if ((pref.equals("usr")) && (packages.get(i).applicationInfo.flags & ApplicationInfo.FLAG_SYSTEM) == ApplicationInfo.FLAG_SYSTEM)
                continue;
            ContentValues values = new ContentValues();


            if (packages.get(i).providers != null) {
                Log.e("App:", packages.get(i).packageName);
                for (ProviderInfo pi : packages.get(i).providers) {
                    if (pi.exported) {
                        Log.e("Permission"," "+pi.readPermission);
                        Log.e("Provider", pi.toString());
                    }
                }
            }

            values.put(Database.DatabaseSchema.InstalledPackages.COLUMN_NAME_PACKAGE_NAME, packages.get(i).packageName);
            values.put(Database.DatabaseSchema.InstalledPackages.COLUMN_NAME_INSTALLED_DATE, packages.get(i).firstInstallTime);
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

    protected static String getRunningServicesXML(SQLiteDatabase db) {
        String[] projection = new String[]{Database.DatabaseSchema.RunningServices.COLUMN_NAME_UID,
                Database.DatabaseSchema.RunningServices.COLUMN_NAME_SINCE,
                Database.DatabaseSchema.RunningServices.COLUMN_NAME_TIME,
                Database.DatabaseSchema.RunningServices.COLUMN_NAME_RX,
                Database.DatabaseSchema.RunningServices.COLUMN_NAME_TX,
                Database.DatabaseSchema.RunningServices.COLUMN_NAME_PACKAGE_NAME};
        Cursor cursor = db.query(Database.DatabaseSchema.RunningServices.TABLE_NAME, projection, null, null, null, null, null);
        String result = XMLProduce.tableToXML(cursor, Database.DatabaseSchema.RunningServices.TAG,
                Database.DatabaseSchema.RunningServices.COLUMN_NAME_PACKAGE_NAME);
        cursor.close();
        return result;
    }

    protected static String getInstalledPackagesXML2(SQLiteDatabase db) {
        String[] projection = new String[]{Database.DatabaseSchema.InstalledPackages.COLUMN_NAME_PACKAGE_NAME,
                Database.DatabaseSchema.InstalledPackages.COLUMN_NAME_INSTALLED_DATE,
                Database.DatabaseSchema.InstalledPackages.COLUMN_NAME_VERSION,
                Database.DatabaseSchema.InstalledPackages.COLUMN_NAME_UID,
                Database.DatabaseSchema.InstalledPackages.COLUMN_NAME_LABEL};
        Cursor cursor = db.query(Database.DatabaseSchema.InstalledPackages.TABLE_NAME, projection, null, null, null, null,
                Database.DatabaseSchema.InstalledPackages.COLUMN_NAME_INSTALLED_DATE + " DESC");
        String result=XMLProduce.tableToXML(cursor,Database.DatabaseSchema.InstalledPackages.TAG,Database.DatabaseSchema.InstalledPackages.COLUMN_NAME_LABEL);
        cursor.close();
        return result;
    }

    protected static String getInstalledPackagesXML(SQLiteDatabase db) {
        StringBuilder result = new StringBuilder();
        Cursor cursor = db.query(Database.DatabaseSchema.InstalledPackages.TABLE_NAME, null, null, null, null, null,
                Database.DatabaseSchema.InstalledPackages.COLUMN_NAME_INSTALLED_DATE + " DESC");
        int uid = cursor.getColumnIndex(Database.DatabaseSchema.InstalledPackages.COLUMN_NAME_UID);
        int label = cursor.getColumnIndex(Database.DatabaseSchema.InstalledPackages.COLUMN_NAME_LABEL);
        int version = cursor.getColumnIndex(Database.DatabaseSchema.InstalledPackages.COLUMN_NAME_VERSION);
        int date = cursor.getColumnIndex(Database.DatabaseSchema.InstalledPackages.COLUMN_NAME_INSTALLED_DATE);
        int name = cursor.getColumnIndex(Database.DatabaseSchema.InstalledPackages.COLUMN_NAME_PACKAGE_NAME);

        while (cursor.moveToNext()) {
            result.append("<installedapp name=\"" + cursor.getString(name) + "\" installed=\"" +
                    TimeCapture.getTime(cursor.getLong(date)) + "\" version=\"" + cursor.getString(version) + "\" uid=\"" +
                    cursor.getString(uid) + "\">" + cursor.getString(label) + "</installedapp>\n");
        }
        cursor.close();
        return result.toString();
    }
}