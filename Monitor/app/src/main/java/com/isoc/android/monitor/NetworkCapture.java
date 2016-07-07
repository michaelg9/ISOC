package com.isoc.android.monitor;

import android.annotation.TargetApi;
import android.content.Context;
import android.net.ConnectivityManager;
import android.net.Network;
import android.net.NetworkInfo;
import android.net.TrafficStats;
import android.net.wifi.WifiInfo;
import android.net.wifi.WifiManager;
import android.os.Build;
import android.os.SystemClock;
import android.text.format.Formatter;
import android.util.Log;
import android.widget.Toast;

import java.io.BufferedReader;
import java.io.File;
import java.io.FileNotFoundException;
import java.io.FileReader;
import java.io.IOException;
import java.math.BigInteger;
import java.net.InetAddress;
import java.net.NetworkInterface;
import java.net.SocketException;
import java.net.UnknownHostException;
import java.nio.ByteBuffer;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Enumeration;

/**
 * Created by maik on 1/7/2016.
 * !!!!!!!!!!!!!!!!!!!!!!!!!!!when mobile off, counters =0
 */
public class NetworkCapture {
    public static String wifiIntfName=new String();
    private static String mobileIntf=new String();

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
        wifiIntfName=result;
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

    private static String[] getTrafficStats(Context context) {
        String wifiIntf=getWifiIntfName(context);
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
        String[] result = new String[7];
        result[0]=Boolean.toString(mobileOn);
        result[1] = Long.toString(TrafficStats.getMobileRxBytes());
        result[2] = Long.toString(TrafficStats.getMobileTxBytes());
        result[3]=Boolean.toString(wifiOn);
        result[4] = readStatsFromFile(context, "/sys/class/net/" + wifiIntf + "/statistics/rx_bytes");
        result[5] = readStatsFromFile(context, "/sys/class/net/" + wifiIntf + "/statistics/tx_bytes");
        result[6] = Long.toString(System.currentTimeMillis() - SystemClock.elapsedRealtime());
        return result;

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

    protected static String getTrafficXML(Context context, String timeFormat) {
        String[] traffic = getTrafficStats(context);
        return ("<data active=\""+traffic[0]+"\" time=\"" + TimeCapture.getTime(timeFormat) + "\" since=\"" + TimeCapture.getTime(timeFormat, Long.parseLong(traffic[6])) + "\" rx=\""+traffic[1]+"\" tx=\""+traffic[2]+"\">mobile</data>\n") +
                ("<data active=\""+traffic[3]+"\" time=\"" + TimeCapture.getTime(timeFormat) + "\" since=\"" + TimeCapture.getTime(timeFormat, Long.parseLong(traffic[6])) + "\" rx=\""+traffic[4]+"\" tx=\""+traffic[5]+"\">wifi</data>\n");
    }

}