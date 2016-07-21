package com.isoc.android.monitor;

import android.content.ContentValues;
import android.content.Context;
import android.content.SharedPreferences;
import android.database.Cursor;
import android.database.sqlite.SQLiteDatabase;
import android.net.ConnectivityManager;
import android.net.NetworkInfo;
import android.net.TrafficStats;
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

    public static String readStatsFromFile(Context context, String fileLocation) {
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
                rx = Long.parseLong(readStatsFromFile(context, "/sys/class/net/" + wifiIntfName + "/statistics/rx_bytes"));
                tx = Long.parseLong(readStatsFromFile(context, "/sys/class/net/" + wifiIntfName + "/statistics/tx_bytes"));
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
                new String[]{Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_TOTAL_RX, Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_TOTAL_TX},
                null,null,null,null, Database.DatabaseSchema.NetworkInterface._ID+" DESC","1");


        long totalrx=0;
        long totaltx=0;
        if (cursor.moveToFirst()){  //if there is at least one entry, get the totals
            totalrx=cursor.getLong(cursor.getColumnIndex(Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_TOTAL_RX));
            totaltx=cursor.getLong(cursor.getColumnIndex(Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_TOTAL_TX));
        }
        cursor.close();
        ContentValues values = new ContentValues();
        values.put(Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_TYPE, type);
        values.put(Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_ACTIVE, Boolean.toString(affectedNet.isConnected()));
        values.put(Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_CURRENT_RX, rx);
        values.put(Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_CURRENT_TX, tx);
        values.put(Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_TOTAL_RX,totalrx);
        values.put(Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_TOTAL_TX,totaltx);
        values.put(Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_TIME, time);

        db.insert(Database.DatabaseSchema.NetworkInterface.TABLE_NAME, null, values);
        db.close();
    }

    /*
    //retrieve active interface(s) on API>=21
    @TargetApi(21)
    private static String[][] getActiveInterface21(Context context) {
        ConnectivityManager connectivityManager = (ConnectivityManager) context.getSystemService(Context.CONNECTIVITY_SERVICE);
        Network[] networks = connectivityManager.getAllNetworks();
        String[][] result = new String[networks.length][4];
        for (int i = 0; i < networks.length; i++) {
            if (connectivityManager.getNetworkInfo(networks[i]).isConnected()) {
                result[i][0] = connectivityManager.getLinkProperties(networks[i]).getInterfaceName();
                result[i][1] = connectivityManager.getNetworkInfo(networks[i]).getTypeName();
                result[i][2] = readStatsFromFile(context, "/sys/class/net/" + result[i][0] + "/statistics/rx_bytes");
                result[i][3] = readStatsFromFile(context, "/sys/class/net/" + result[i][0] + "/statistics/tx_bytes");
                //result[i][2] = Integer.toString(connectivityManager.getNetworkCapabilities(networks[i]).getLinkDownstreamBandwidthKbps());
                //result[i][3] = Integer.toString(connectivityManager.getNetworkCapabilities(networks[i]).getLinkUpstreamBandwidthKbps());

            }
        }
        return result;
    }

    //retrieve active interface on API<21
    private static String[] getActiveInterface15(Context context) {
        ConnectivityManager connectivityManager = (ConnectivityManager) context.getSystemService(Context.CONNECTIVITY_SERVICE);
        NetworkInfo network = connectivityManager.getActiveNetworkInfo();

        ArrayList<NetworkInterface> active = new ArrayList<NetworkInterface>();
        try {
            for (Enumeration<NetworkInterface> list = NetworkInterface.getNetworkInterfaces(); list.hasMoreElements(); ) {
                NetworkInterface i = list.nextElement();
                if (i.isUp() && !i.isLoopback() && i.getInterfaceAddresses().size() > 0)
                    active.add(i);
            }
        } catch (SocketException e) {
            Toast.makeText(context,e.toString(),Toast.LENGTH_LONG).show();
        }

        if (network == null || active.size() != 1) return null;

        String[] result = new String[4];
        result[0] = active.get(0).getDisplayName();
        result[1] = network.getTypeName();
        //result[2] =network.getSubtypeName();// Integer.toString(connectivityManager.getNetworkCapabilities(networks[i]).getLinkDownstreamBandwidthKbps());
        //result[3] =null;// Integer.toString(connectivityManager.getNetworkCapabilities(networks[i]).getLinkUpstreamBandwidthKbps());
        result[2] = readStatsFromFile(context, "/sys/class/net/" + result[0] + "/statistics/rx_bytes");
        result[3] = readStatsFromFile(context, "/sys/class/net/" + result[0] + "/statistics/tx_bytes");
        return result;
    }

    protected static String getActiveInterfaceXML(Context context) {
        StringBuilder result = new StringBuilder();
        if (Build.VERSION.SDK_INT < Build.VERSION_CODES.LOLLIPOP) {
            String[] i =getActiveInterface15(context);
            result.append("<intf name=\""+i[0]+"\" rx=\"" + i[2] + "\" tx=\"" + i[3] + "\">" + i[1] + "</intf>\n");
        }
        else{
            for (String[] i : getActiveInterface21(context))
                result.append("<activeintf name=\""+i[0]+"\" rx=\"" + i[2] + "\" tx=\"" + i[3] + "\">" + i[1] + "</activeintf>\n");
        }

        return result.toString();
    }*/

    protected static String getTrafficXML(SQLiteDatabase db) {
        Cursor cursor = db.query(Database.DatabaseSchema.NetworkInterface.TABLE_NAME, null, null, null, null, null, null);
        int sinceIndex = cursor.getColumnIndex(Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_SINCE);
        int activeIndex = cursor.getColumnIndex(Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_ACTIVE);
        int curRxIndex = cursor.getColumnIndex(Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_CURRENT_RX);
        int curTxIndex = cursor.getColumnIndex(Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_CURRENT_TX);
        int totRxIndex = cursor.getColumnIndex(Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_TOTAL_RX);
        int totTxIndex = cursor.getColumnIndex(Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_TOTAL_TX);
        int typeIndex = cursor.getColumnIndex(Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_TYPE);
        int timeIndex = cursor.getColumnIndex(Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_TIME);

        StringBuilder result = new StringBuilder();
        while (cursor.moveToNext()) {
            Log.e("tx-tot/cur",cursor.getLong(totTxIndex)+"/"+cursor.getLong(curTxIndex));
            Log.e("rx-tot/cur",cursor.getLong(totRxIndex)+"/"+cursor.getLong(curRxIndex));
            result.append("<data active=\"" + cursor.getString(activeIndex) + "\" time=\"" +
                    cursor.getString(timeIndex) + "\" since=\"" + cursor.getString(sinceIndex) + "\" rx=\"" +
                    Long.toString(cursor.getLong(curRxIndex) + cursor.getLong(totRxIndex)) + "\" tx=\"" + (cursor.getLong(curTxIndex) +
                    cursor.getLong(totTxIndex)) + "\">" + cursor.getString(typeIndex) + "</data>\n");
        }
        cursor.close();
        return result.toString();
    }
}