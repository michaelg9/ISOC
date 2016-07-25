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
import android.preference.PreferenceManager;
import android.util.Log;
import android.widget.Toast;

import java.io.BufferedReader;
import java.io.FileNotFoundException;
import java.io.FileReader;
import java.io.IOException;
import java.net.NetworkInterface;
import java.net.SocketException;
import java.util.Enumeration;
import java.util.List;

/**
 * Created by maik on 1/7/2016.

 */
public class NetworkCapture {

    private static String getWifiIntfName(Context context) {
        SharedPreferences prefs = PreferenceManager.getDefaultSharedPreferences(context);
        String result = prefs.getString("wifiInterfaceName", null);                        //in resources...
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
        prefs.edit().putString("wifiInterfaceName", result).apply();                 //in resources....
        return result;
    }

    public static String readStatsFromFile(String fileLocation) {
        FileReader file = null;
        BufferedReader in = null;
        StringBuilder result = new StringBuilder();
        try {
            file = new FileReader(fileLocation);
            in = new BufferedReader(file);
            String s = in.readLine();
            while (s!=null){
                result.append(s);
                s = in.readLine();
            }

        } catch (FileNotFoundException e) {
            Log.e("EXCEPTION",e.getMessage());
        } catch (IOException e) {
            Log.e("EXCEPTION",e.getMessage());
        } finally {
            if (file != null) try {
                file.close();
            } catch (IOException e) {
                Log.e("EXCEPTION",e.getMessage());
            }
            if (in != null) try {
                in.close();
            } catch (IOException e) {
                Log.e("EXCEPTION",e.getMessage());
            }
        }
        return result.toString();
    }

    protected static void saveCurrentStats(Context context){
        SQLiteDatabase db=new Database(context).getWritableDatabase();
        db.execSQL("UPDATE "+ Database.DatabaseSchema.NetworkInterface.TABLE_NAME+" SET "+
                Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_TOTAL_TX+" ="+
                Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_CURRENT_TX+ ","+
                Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_TOTAL_RX+" ="+
                Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_CURRENT_RX+ ","+
                Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_CURRENT_RX+"=0,"+
                Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_CURRENT_TX+"=0 WHERE "+
                Database.DatabaseSchema.NetworkInterface._ID+"= (SELECT MAX("+ Database.DatabaseSchema.NetworkInterface._ID+") FROM "+
                Database.DatabaseSchema.NetworkInterface.TABLE_NAME+")");
        db.close();
    }

    protected static void getTrafficStats(Context context, NetworkInfo affectedNet) {
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
                rx = TrafficStats.getMobileRxBytes();   //stupid method
                tx = TrafficStats.getMobileTxBytes();   //another stupid method
                break;
            default:
                return;
        }
        SQLiteDatabase db = new Database(context).getWritableDatabase();
        Cursor cursor=db.query(Database.DatabaseSchema.NetworkInterface.TABLE_NAME,
                new String[]{Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_TOTAL_RX,
                        Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_TOTAL_TX,
                        Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_SINCE},null,null,null,null,
                Database.DatabaseSchema.NetworkInterface._ID+" DESC","1");


        long totalrx=0;
        long totaltx=0;
        String since = TimeCapture.getUpDate();  // if there's no entry in the table, 'since' will be last reboot date
        if (cursor.moveToFirst()){  //if there is at least one entry, get the totals
            totalrx=cursor.getLong(cursor.getColumnIndex(Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_TOTAL_RX));
            totaltx=cursor.getLong(cursor.getColumnIndex(Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_TOTAL_TX));
            since = cursor.getString(cursor.getColumnIndex(Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_SINCE));
        }
        cursor.close();
        ContentValues values = new ContentValues();
        values.put(Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_TYPE, type);
        values.put(Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_SINCE, since);
        values.put(Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_ACTIVE, Boolean.toString(affectedNet.isConnected()));
        values.put(Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_CURRENT_RX, rx);
        values.put(Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_CURRENT_TX, tx);
        values.put(Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_TOTAL_RX,totalrx);
        values.put(Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_TOTAL_TX,totaltx);
        values.put(Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_TIME, time);

        db.insertWithOnConflict(Database.DatabaseSchema.NetworkInterface.TABLE_NAME, null, values,SQLiteDatabase.CONFLICT_IGNORE);
        db.close();
    }

    public static void getWifiAPs(Context context,SQLiteDatabase db){
        WifiManager wifiManager=(WifiManager) context.getSystemService(Context.WIFI_SERVICE);
        List<ScanResult> scanResults=wifiManager.getScanResults();
        for (ScanResult sr : scanResults){
            ContentValues values=new ContentValues();
            values.put(Database.DatabaseSchema.WifiAP.COLUMN_NAME_BSSID,sr.BSSID);
            values.put(Database.DatabaseSchema.WifiAP.COLUMN_NAME_SSID,sr.SSID);
            values.put(Database.DatabaseSchema.WifiAP.COLUMN_NAME_CAPABILITIES,sr.capabilities);
            //values.put(Database.DatabaseSchema.WifiAP.COLUMN_NAME_SEEN,sr.timestamp);
            values.put(Database.DatabaseSchema.WifiAP.COLUMN_NAME_SIGNAL,WifiManager.calculateSignalLevel(sr.level,11));
            values.put(Database.DatabaseSchema.WifiAP.COLUMN_NAME_FREQ,sr.frequency);
            db.insertWithOnConflict(Database.DatabaseSchema.WifiAP.TABLE_NAME,null,values,SQLiteDatabase.CONFLICT_IGNORE);
        }
    }

    public static String getWifiAPResultsXML(SQLiteDatabase db){
        String[] projection = new String[]{Database.DatabaseSchema.WifiAP.COLUMN_NAME_BSSID,
                Database.DatabaseSchema.WifiAP.COLUMN_NAME_SSID,
                Database.DatabaseSchema.WifiAP.COLUMN_NAME_CAPABILITIES,
                Database.DatabaseSchema.WifiAP.COLUMN_NAME_SIGNAL,
                Database.DatabaseSchema.WifiAP.COLUMN_NAME_FREQ};
        Cursor cursor=db.query(Database.DatabaseSchema.WifiAP.TABLE_NAME,projection,null,null,null,null,null);
        String result=XMLProduce.tableToXML(cursor,Database.DatabaseSchema.WifiAP.TAG,Database.DatabaseSchema.WifiAP.COLUMN_NAME_SSID);
        cursor.close();
        return result;
    }

    protected static String getTrafficXML(SQLiteDatabase db) {
        String query=String.format("SELECT %s,%s,%s,(%s+%s) AS rx,(%s+%s) AS tx,%s FROM %s",
                Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_ACTIVE,
                Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_TIME,
                Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_SINCE,
                Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_TOTAL_RX,
                Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_CURRENT_RX,
                Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_TOTAL_TX,
                Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_CURRENT_TX,
                Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_TYPE,
                Database.DatabaseSchema.NetworkInterface.TABLE_NAME);
        Cursor cursor = db.rawQuery(query, null);
        String result = XMLProduce.tableToXML(cursor,Database.DatabaseSchema.NetworkInterface.TAG,Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_TYPE);
        cursor.close();
        return result;
    }



}