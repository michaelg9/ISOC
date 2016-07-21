package com.isoc.android.monitor;

import android.content.Context;
import android.content.Intent;
import android.content.pm.PackageManager;
import android.content.pm.ResolveInfo;
import android.net.Uri;
import android.os.Build;
import android.telephony.TelephonyManager;

import java.util.ArrayList;

/**
 * Created by maik on 4/7/2016.
 */
public class MetaDataCapture {

    private static ArrayList<String[]> getMetaData(Context context) {
        ArrayList<String[]> result = new ArrayList<String[]>();
        getTelephonyDetails(context, result);
        getPhoneDetails(result);
        getDefaultBrowser(context,result);
        return result;
    }

    private static void getTelephonyDetails(Context context, ArrayList<String[]> result) {
        String[] datatype = new String[]{"unknown", "gprs", "edge", "umts", "cdma", "evdo0", "evdoA", "1xrtt", "hsdpa", "hsupa", "hspa", "iden", "evdoB", "lte", "ehrpd", "hspap"};
        TelephonyManager tm = (TelephonyManager) context.getSystemService(Context.TELEPHONY_SERVICE);
        result.add(new String[]{"imei", tm.getDeviceId()});
        result.add(new String[]{"datanettype", datatype[tm.getNetworkType()]});
        result.add(new String[]{"country", tm.getNetworkCountryIso()});
        result.add(new String[]{"network", tm.getNetworkOperatorName()});
        result.add(new String[]{"carrier", tm.getSimOperatorName()});
    }

    private static void getPhoneDetails(ArrayList<String[]> result) {
        result.add(new String[]{"manufacturer", Build.MANUFACTURER});
        result.add(new String[]{"model", Build.MODEL});
        result.add(new String[]{"androidver", Build.VERSION.RELEASE});
        result.add(new String[]{"uptime", Long.toString(TimeCapture.getUpTime())});
    }

    private static void getDefaultBrowser(Context context,ArrayList<String[]> result){
        Intent browseIntent =new Intent("android.intent.action.VIEW", Uri.parse("http://"));
        ResolveInfo defaultBrowse=context.getPackageManager().resolveActivity(browseIntent, PackageManager.MATCH_DEFAULT_ONLY);
        result.add(new String[]{"defaultBrowser",defaultBrowse.activityInfo.packageName});
    }

    protected static String getMetaDataXML(Context context) {
        StringBuilder result = new StringBuilder();
        for (String[] d : getMetaData(context)) {
            result.append("<" + d[0] + ">" + d[1] + "</" + d[0] + ">\n");
        }
        return result.toString();
    }
}