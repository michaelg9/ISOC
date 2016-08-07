package com.isoc.android.monitor;

import android.content.ContentValues;
import android.content.Context;
import android.content.SharedPreferences;
import android.database.Cursor;
import android.database.sqlite.SQLiteDatabase;
import android.net.ConnectivityManager;
import android.net.NetworkInfo;
import android.net.TrafficStats;
import android.net.wifi.ScanResult;
import android.net.wifi.WifiInfo;
import android.net.wifi.WifiManager;
import android.os.Build;
import android.util.Log;
import android.widget.Toast;

import java.io.BufferedReader;
import java.io.File;
import java.io.FileNotFoundException;
import java.io.FileReader;
import java.io.IOException;
import java.net.NetworkInterface;
import java.net.SocketException;
import java.util.Enumeration;
import java.util.List;

/**
 * Captures connectivity changes.
 *
 * Current mobile intf statistics can't reliably be captured using trafficstats, because it's buggy:
 * it may return -1 or 0 if the mobile intf is turned off, even if the actual value in the
 * respective file isn't 0.
 *
 * The way used to captured mobile stats is as following:
 * 0) Check if we have saved the mobile intf name. If yes, use it to read files directly
 * 1)Try trafficstats. If the value returned isn't what expected go to 2. Otherwise use it to find intf name and save it.
 * 2)Usually the mobile intf is named rmnet0. Try to read its files. If no luck:
 * 3)Try intf ppp0 which is another possible name.
 * If everything fails, we return 0. However, when the mobile intf is turned on, Trafficstats
 * usually is accurate so we use that chance to determine the name of the intf (by comparing the returned value to all the
 * statistics files of the intfs listed in /sys/class/net/)
 *
 *  Wireless intf statistics are read directly from their files --> We can find the intf name by comparing MAC addresses
 */
public class NetworkCapture {

    private static String getWifiIntfName(Context context) {
        //compares the wifi intf mac address returned by wifimanager to all intfs' mac addresses to find the name of the wifi intf
        //then saves it into sharedpreferences for faster access. Used to read tx/rx directly from the files
        SharedPreferences preferences = context.getSharedPreferences(context.getString(R.string.shared_values_filename),Context.MODE_PRIVATE);
        String result = preferences.getString("wifiInterfaceName", null);
        if (result != null) return result;

        WifiManager wifi = (WifiManager) context.getSystemService(Context.WIFI_SERVICE);
        WifiInfo wifiInfo = wifi.getConnectionInfo();
        try {
            String wifiMAC = wifiInfo.getMacAddress().replaceAll(":", "");
            for (Enumeration<NetworkInterface> list = NetworkInterface.getNetworkInterfaces(); list.hasMoreElements(); ) {
                NetworkInterface i = list.nextElement();
                byte[] intfMACBytes = i.getHardwareAddress();
                StringBuilder intfMAC = new StringBuilder();
                if (!i.isLoopback() && intfMACBytes != null) {
                    for (byte b : intfMACBytes) {
                        intfMAC.append(String.format("%02x", b));
                        if (!wifiMAC.startsWith(intfMAC.toString())) break;
                    }
                }
                if (wifiMAC.equals(intfMAC.toString()))
                    result = i.getDisplayName();
            }
        } catch (SocketException e) {
            Toast.makeText(context, e.toString(), Toast.LENGTH_LONG).show();
        }
        preferences.edit().putString("wifiInterfaceName", result).apply();
        return result;
    }

    public static String readStatsFromFile(String fileLocation) {
        //reads the contents of a file in fileLocation.
        FileReader file = null;
        BufferedReader in = null;
        StringBuilder result = new StringBuilder();
        try {
            file = new FileReader(fileLocation);
            in = new BufferedReader(file);
            String s = in.readLine();
            while (s != null) {
                result.append(s);
                s = in.readLine();
            }

        } catch (FileNotFoundException e) {
            Log.e("EXCEPTION", e.getMessage());
        } catch (IOException e) {
            Log.e("EXCEPTION", e.getMessage());
        } finally {
            if (file != null) try {
                file.close();
            } catch (IOException e) {
                Log.e("EXCEPTION", e.getMessage());
            }
            if (in != null) try {
                in.close();
            } catch (IOException e) {
                Log.e("EXCEPTION", e.getMessage());
            }
        }
        return result.toString();
    }

    public static void saveCurrentStats(Context context) {
        //adds the latest current rx/tx to the totals and resets the currents
        // (for both wifi & mobile). Executed when boot action is broadcasted.
        SQLiteDatabase db = new Database(context).getWritableDatabase();
        //we can't use update method to directly add one column's value to another using ContentValues, so execute directly
        String query = String.format("UPDATE %s SET %s=%s+%s,%s=%s+%s,%s=0,%s=0 " +
                        "WHERE %s=(SELECT MAX(%s) FROM %s WHERE %s LIKE 'wifi') OR " +
                        "%s=(SELECT MAX(%s) FROM %s WHERE %s LIKE 'mobile')",
                Database.DatabaseSchema.NetworkInterface.TABLE_NAME,
                Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_TOTAL_TX,
                Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_TOTAL_TX,
                Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_CURRENT_TX,
                Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_TOTAL_RX,
                Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_TOTAL_RX,
                Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_CURRENT_RX,
                Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_CURRENT_TX,
                Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_CURRENT_RX,
                Database.DatabaseSchema.NetworkInterface._ID,
                Database.DatabaseSchema.NetworkInterface._ID,
                Database.DatabaseSchema.NetworkInterface.TABLE_NAME,
                Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_TYPE,
                Database.DatabaseSchema.NetworkInterface._ID,
                Database.DatabaseSchema.NetworkInterface._ID,
                Database.DatabaseSchema.NetworkInterface.TABLE_NAME,
                Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_TYPE);
        db.execSQL(query);

        /*String query2=String.format("SELECT * FROM %s " +
                        "WHERE %s=(SELECT MAX(%s) FROM %s WHERE %s LIKE 'wifi') OR " +
                        "%s=(SELECT MAX(%s) FROM %s WHERE %s LIKE 'mobile')",
                Database.DatabaseSchema.NetworkInterface.TABLE_NAME,
                Database.DatabaseSchema.NetworkInterface._ID,
                Database.DatabaseSchema.NetworkInterface._ID,
                Database.DatabaseSchema.NetworkInterface.TABLE_NAME,
                Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_TYPE,
                Database.DatabaseSchema.NetworkInterface._ID,
                Database.DatabaseSchema.NetworkInterface._ID,
                Database.DatabaseSchema.NetworkInterface.TABLE_NAME,
                Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_TYPE);
        Cursor c = db.rawQuery(query2,null);
        while (c.moveToNext()){
            StringBuilder sb=new StringBuilder();
            for (String s : c.getColumnNames()){
                sb.append(s+':'+c.getString(c.getColumnIndex(s))+" ");
            }
            Log.e("ROW",sb.toString());
        }c.close();*/


        db.close();
    }

    private static long getMobileRx(Context context) {
        //gets the value of mobile rx bytes.
        // If it's not working, fall back to reading directly common intf names

        //First we check if the mobile intf name is saved.
        SharedPreferences preferences = context.getSharedPreferences(context.getString(R.string.shared_values_filename),Context.MODE_PRIVATE);
        String mobileIntf = preferences.getString("mobileInterfaceName", null);
        if (mobileIntf != null) return Long.parseLong(readStatsFromFile("/sys/class/net/"+mobileIntf+"/statistics/rx_bytes"));

        // We try Trafficstats. If it's working, we try to find intf name using the returned value and save it.
        long result = TrafficStats.getMobileRxBytes();
        if (result > 0){
            String[] interfaces=(new File("/sys/class/net/")).list();
            for (String intf : interfaces){
                String rx=readStatsFromFile("/sys/class/net/"+intf+"/statistics/rx_bytes");
                if (rx!=null && Long.parseLong(rx)==TrafficStats.getMobileRxBytes()){
                    //we found the intf!! Save it for future reference
                    preferences.edit().putString("mobileInterfaceName", intf).apply();
                }
            }
            return result;
        }

        //Otherwise try most common name for the mobile intf: rmnet0
        String s = readStatsFromFile("/sys/class/net/rmnet0/statistics/rx_bytes");
        if (s != null) result = Long.parseLong(s);
        if (result > 0) return result;
        //another less common name is ppp0
        s = readStatsFromFile("/sys/class/net/ppp0/statistics/rx_bytes");
        if (s != null) result = Long.parseLong(s);
        if (result > 0) return result;
        return 0;
    }

    private static long getMobileTx(Context context) {
        //gets the value of mobile tx bytes.

        SharedPreferences preferences = context.getSharedPreferences(context.getString(R.string.shared_values_filename),Context.MODE_PRIVATE);
        String mobileIntf = preferences.getString("mobileInterfaceName", null);
        if (mobileIntf != null) return Long.parseLong(readStatsFromFile("/sys/class/net/"+mobileIntf+"/statistics/tx_bytes"));

        long result = TrafficStats.getMobileTxBytes();
        if (result > 0) return result;

        String s = readStatsFromFile("/sys/class/net/rmnet0/statistics/tx_bytes");
        if (s != null) result = Long.parseLong(s);
        if (result > 0) return result;

        s = readStatsFromFile("/sys/class/net/ppp0/statistics/tx_bytes");
        if (s != null) result = Long.parseLong(s);
        if (result > 0) return result;
        return 0;
    }

    public static void getTrafficStats(Context context, NetworkInfo affectedNet) {
        if (affectedNet == null) return;
        String time = TimeCapture.getTime();
        String type = new String();
        Log.e("NETCHANGE", affectedNet.toString());
        long tx;
        long rx;
        switch (affectedNet.getType()) {
            case ConnectivityManager.TYPE_WIFI:
                type = "wifi";
                String wifiIntfName = getWifiIntfName(context);
                rx = Long.parseLong(readStatsFromFile("/sys/class/net/" + wifiIntfName + "/statistics/rx_bytes"));
                tx = Long.parseLong(readStatsFromFile("/sys/class/net/" + wifiIntfName + "/statistics/tx_bytes"));
                break;
            case (ConnectivityManager.TYPE_MOBILE):
                type = "mobile";
                rx = getMobileRx(context);
                tx = getMobileTx(context);
                break;
            default:
                return;
        }
        SQLiteDatabase db = new Database(context).getWritableDatabase();
        //get the last record of the connectivity DB that matches affected type. Used to retrieve totals and since fields.
        Cursor cursor = db.query(Database.DatabaseSchema.NetworkInterface.TABLE_NAME,
                new String[]{Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_TOTAL_RX,
                        Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_TOTAL_TX,
                        Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_SINCE},
                Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_TYPE + "='" + type+'\'',
                null, null, null, Database.DatabaseSchema.NetworkInterface._ID + " DESC", "1");


        long totalrx = 0;
        long totaltx = 0;
        String since = TimeCapture.getUpDate();  // if there's no entry in the table, 'since' will be last reboot date
        if (cursor.moveToFirst()) {  //if there is at least one entry, get the totals and since
            totalrx = cursor.getLong(cursor.getColumnIndex(Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_TOTAL_RX));
            totaltx = cursor.getLong(cursor.getColumnIndex(Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_TOTAL_TX));
            since = cursor.getString(cursor.getColumnIndex(Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_SINCE));
        }
        cursor.close();
        ContentValues values = new ContentValues();
        values.put(Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_TYPE, type);
        values.put(Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_SINCE, since);
        values.put(Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_ACTIVE, Boolean.toString(affectedNet.isConnected()));
        values.put(Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_CURRENT_RX, rx);
        values.put(Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_CURRENT_TX, tx);
        values.put(Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_TOTAL_RX, totalrx);
        values.put(Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_TOTAL_TX, totaltx);
        values.put(Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_TIME, time);

        db.insertWithOnConflict(Database.DatabaseSchema.NetworkInterface.TABLE_NAME, null, values, SQLiteDatabase.CONFLICT_IGNORE);
        db.close();
    }

    public static void getWifiAPs(Context context, SQLiteDatabase db) {
        //retrieves last scan's results, so we avoid draining the battery

        WifiManager wifiManager = (WifiManager) context.getSystemService(Context.WIFI_SERVICE);
        List<ScanResult> scanResults = wifiManager.getScanResults();
        if (scanResults==null) return;
        for (ScanResult sr : scanResults) {
            ContentValues values = new ContentValues();
            values.put(Database.DatabaseSchema.WifiAP.COLUMN_NAME_BSSID, sr.BSSID);
            values.put(Database.DatabaseSchema.WifiAP.COLUMN_NAME_SSID, sr.SSID);
            values.put(Database.DatabaseSchema.WifiAP.COLUMN_NAME_CAPABILITIES, sr.capabilities);

            //timestamp only available on api 17+
            if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.JELLY_BEAN_MR1)
                values.put(Database.DatabaseSchema.WifiAP.COLUMN_NAME_SEEN, TimeCapture.getTime(sr.timestamp));

            //signal quality from 0 to 10
            values.put(Database.DatabaseSchema.WifiAP.COLUMN_NAME_SIGNAL, WifiManager.calculateSignalLevel(sr.level, 11));
            values.put(Database.DatabaseSchema.WifiAP.COLUMN_NAME_FREQ, sr.frequency);
            db.insertWithOnConflict(Database.DatabaseSchema.WifiAP.TABLE_NAME, null, values, SQLiteDatabase.CONFLICT_REPLACE);
        }
    }
}