package com.isoc.android.monitor;

import android.content.Context;
import android.os.Build;
import android.os.SystemClock;
import android.telephony.TelephonyManager;

/**
 * Created by maik on 4/7/2016.
 */
public class MetaDataCapture {

    private static String[][] getTelephonyDetails(Context context){
        String[] datatype=new String[]{"unknown","gprs","edge","umts","cdma","evdo0","evdoA","1xrtt","hsdpa","hsupa","hspa","iden","evdoB","lte","ehrpd","hspap"};
        String[] datastate=new String[]{"disconnected","connecting","connected","suspended"};
        TelephonyManager tm=(TelephonyManager) context.getSystemService(context.TELEPHONY_SERVICE);
        String result[][]=new String[6][2];
        result[0]=new String[]{"imei",tm.getDeviceId()};
        result[1]=new String[]{"datanettype",datatype[tm.getNetworkType()]};
        result[2]=new String[]{"datastate",datastate[tm.getDataState()]};
        result[3]=new String[]{"country",tm.getNetworkCountryIso()};
        result[4]=new String[]{"network",tm.getNetworkOperatorName()};
        result[5]=new String[]{"carrier",tm.getSimOperatorName()};
        return result;
    }

    private static String[][] getPhoneDetails(){
        String result[][] = new String[4][2];
        result[0]=new String[]{"manufacturer",Build.MANUFACTURER};
        result[1]=new String[]{"model",Build.MODEL};
        result[2]=new String[]{"androidver",Build.VERSION.RELEASE};
        result[3]=new String[]{"uptime",Long.toString(SystemClock.elapsedRealtime()/1000)};
        return result;
    }

    protected static String getMetaDataXML(Context context){
        StringBuilder result=new StringBuilder();
        for (String[] d : getTelephonyDetails(context)) {
            result.append("<"+d[0]+">" + d[1] + "</"+d[0]+">\n");
    }
        for (String[] d : getPhoneDetails())
            result.append("<"+d[0]+">" + d[1] + "</"+d[0]+">\n");
        return result.toString();
    }
}
