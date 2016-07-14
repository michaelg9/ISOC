package com.isoc.android.monitor;

import android.content.Context;
import android.content.SharedPreferences;
import android.net.ConnectivityManager;
import android.net.TrafficStats;
import android.net.wifi.WifiInfo;
import android.net.wifi.WifiManager;
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
 * !!!!!!!!!!!!!!!!!!!!!!!!!!!when mobile off, counters =0
 */
public class NetworkCapture {

    private static String getWifiIntfName(Context context){
        WifiManager wifi = (WifiManager) context.getSystemService(Context.WIFI_SERVICE);
        WifiInfo wifiInfo=wifi.getConnectionInfo();
        String result=new String();
        try {
            String wifiMAC=wifiInfo.getMacAddress().replaceAll(":","");
            for (Enumeration<NetworkInterface> list = NetworkInterface.getNetworkInterfaces(); list.hasMoreElements(); ) {
                NetworkInterface i = list.nextElement();
                byte[] intfMACBytes=i.getHardwareAddress();
                StringBuilder intfMAC = new StringBuilder();
                if (!i.isLoopback() && intfMACBytes!=null){
                    for (byte b : intfMACBytes) {
                        intfMAC.append(String.format("%02x", b));
                        if (!wifiMAC.startsWith(intfMAC.toString())) break;
                    }
                }
                if (wifiMAC.equals(intfMAC.toString()))
                    result= i.getDisplayName();
            }
        } catch (SocketException e) {
            Toast.makeText(context,e.toString(),Toast.LENGTH_LONG).show();
        }
        SharedPreferences.Editor wifiEdit = context.getSharedPreferences(Database.WifiNetworkInterface.PREFERENCES_FILENAME,Context.MODE_PRIVATE).edit();
        wifiEdit.putString(Database.WifiNetworkInterface.KEY_INTF_NAME,result);
        wifiEdit.apply();
        return result;
    }

    private static String readStatsFromFile(Context context, String fileLocation) {
        FileReader file = null;
        BufferedReader in = null;
        String result = new String();
        try {
            file = new FileReader(fileLocation);
            in = new BufferedReader(file);
            result = in.readLine();
        } catch (FileNotFoundException e) {
            Toast.makeText(context, e.getMessage(), Toast.LENGTH_SHORT).show();
        } catch (IOException e) {
            Toast.makeText(context, e.getMessage(), Toast.LENGTH_SHORT).show();
        } finally {
            if (file != null) try {
                file.close();
            } catch (IOException e) {
                Toast.makeText(context, e.getMessage(), Toast.LENGTH_SHORT).show();
            }
            if (in != null) try {
                in.close();
            } catch (IOException e) {
                Toast.makeText(context, e.getMessage(), Toast.LENGTH_SHORT).show();
            }
        }
        return result;
    }

    protected static void getTrafficStats(Context context) {
        ConnectivityManager cm = (ConnectivityManager) context.getSystemService(Context.CONNECTIVITY_SERVICE);
        boolean wifiOn=false;
        boolean mobileOn=false;
        if (cm.getActiveNetworkInfo()!=null){
            switch (cm.getActiveNetworkInfo().getType()) {
                case ConnectivityManager.TYPE_WIFI:
                    wifiOn = true;
                    break;
                case (ConnectivityManager.TYPE_MOBILE):
                    mobileOn = true;
                    break;
                default:
                    break;
            }
        }

        SharedPreferences mobilePref = context.getSharedPreferences(Database.MobileNetworkInterface.PREFERENCES_FILENAME,Context.MODE_PRIVATE);
        SharedPreferences.Editor mobileEdit = mobilePref.edit();
        long mobileRxSaved = mobilePref.getLong(Database.MobileNetworkInterface.KEY_CURRENT_RX,0);
        long mobileTxSaved = mobilePref.getLong(Database.MobileNetworkInterface.KEY_CURRENT_TX,0);
        long mobileRx = (mobileRxSaved > TrafficStats.getMobileRxBytes()) ? mobileRxSaved : TrafficStats.getMobileRxBytes();
        long mobileTx = (mobileTxSaved > TrafficStats.getMobileTxBytes()) ? mobileTxSaved : TrafficStats.getMobileTxBytes();
        mobileEdit.putString(Database.MobileNetworkInterface.KEY_ACTIVE,Boolean.toString(mobileOn));
        mobileEdit.putLong(Database.MobileNetworkInterface.KEY_CURRENT_RX,mobileRx);
        mobileEdit.putLong(Database.MobileNetworkInterface.KEY_CURRENT_TX,mobileTx);
        mobileEdit.apply();

        SharedPreferences wifiPref = context.getSharedPreferences(Database.WifiNetworkInterface.PREFERENCES_FILENAME,Context.MODE_PRIVATE);
        SharedPreferences.Editor wifiEdit = wifiPref.edit();
        long wifiRxSaved = wifiPref.getLong(Database.WifiNetworkInterface.KEY_CURRENT_RX,0);
        long wifiTxSaved = wifiPref.getLong(Database.WifiNetworkInterface.KEY_CURRENT_TX,0);
        long wifiRxCurrent = Long.parseLong(readStatsFromFile(context, "/sys/class/net/" + wifiPref.getString(Database.WifiNetworkInterface.KEY_INTF_NAME,getWifiIntfName(context)) + "/statistics/rx_bytes"));
        long wifiTxCurrent = Long.parseLong(readStatsFromFile(context, "/sys/class/net/" + wifiPref.getString(Database.WifiNetworkInterface.KEY_INTF_NAME,getWifiIntfName(context)) + "/statistics/tx_bytes"));
        long wifiRx = (wifiRxSaved > wifiRxCurrent) ? wifiRxSaved :wifiRxCurrent;
        long wifiTx = (wifiTxSaved > wifiTxCurrent) ? wifiTxSaved : wifiTxCurrent;
        wifiEdit.putString(Database.WifiNetworkInterface.KEY_ACTIVE,Boolean.toString(wifiOn));
        wifiEdit.putLong(Database.WifiNetworkInterface.KEY_CURRENT_RX,wifiRx);
        wifiEdit.putLong(Database.WifiNetworkInterface.KEY_CURRENT_TX,wifiTx);
        wifiEdit.apply();
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

    protected static String getTrafficXML(Context context) {
        SharedPreferences mobilePref = context.getSharedPreferences(Database.MobileNetworkInterface.PREFERENCES_FILENAME,Context.MODE_PRIVATE);
        SharedPreferences wifiPref = context.getSharedPreferences(Database.WifiNetworkInterface.PREFERENCES_FILENAME,Context.MODE_PRIVATE);

        return ("<data active=\""+mobilePref.getString(Database.MobileNetworkInterface.KEY_ACTIVE,"false")+"\" time=\"" +
                TimeCapture.getTime() + "\" since=\"" +
                TimeCapture.getTime(mobilePref.getLong(Database.MobileNetworkInterface.KEY_SINCE,0)) +"\" rx=\""+
                mobilePref.getLong(Database.MobileNetworkInterface.KEY_CURRENT_RX,0)+
                mobilePref.getLong(Database.MobileNetworkInterface.KEY_TOTAL_RX,0)+"\" tx=\""+
                mobilePref.getLong(Database.MobileNetworkInterface.KEY_CURRENT_TX,0)+
                mobilePref.getLong(Database.MobileNetworkInterface.KEY_TOTAL_TX,0) +"\">mobile</data>\n" +
                "<data active=\""+wifiPref.getString(Database.WifiNetworkInterface.KEY_ACTIVE,"false")+"\" time=\"" +
                TimeCapture.getTime() + "\" since=\"" +
                TimeCapture.getTime(wifiPref.getLong(Database.WifiNetworkInterface.KEY_SINCE,0)) + "\" rx=\"" +
                wifiPref.getLong(Database.WifiNetworkInterface.KEY_CURRENT_RX,0)+
                wifiPref.getLong(Database.WifiNetworkInterface.KEY_TOTAL_RX,0)+"\" tx=\""+
                wifiPref.getLong(Database.WifiNetworkInterface.KEY_CURRENT_TX,0)+
                wifiPref.getLong(Database.WifiNetworkInterface.KEY_TOTAL_TX,0)+"\">wifi</data>\n");
    }
}