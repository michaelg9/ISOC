package com.isoc.android.monitor;

import android.app.ActivityManager;
import android.content.Context;
import android.content.pm.PackageInfo;
import android.content.pm.PackageManager;
import android.net.TrafficStats;
import android.os.SystemClock;

import java.util.List;


/**
 * Created by maik on 1/7/2016.
 * CPU time and memory of running apps?
 */
public class PackageCapture {

    private static String[][] getInstalledPackages(Context context, String timeFormat){
        PackageManager packageManager= context.getPackageManager();
        List<PackageInfo> packages = packageManager.getInstalledPackages(PackageManager.GET_META_DATA);
        String[][] result= new String[packages.size()][5];
        for (int i =0; i<packages.size();i++){
            result[i][0]=packages.get(i).packageName;
            result[i][1]=TimeCapture.getTime(timeFormat,packages.get(i).firstInstallTime);
            result[i][2]=packages.get(i).versionName;
            result[i][3]=Integer.toString(packages.get(i).applicationInfo.uid);
            result[i][4]=packageManager.getApplicationLabel(packages.get(i).applicationInfo).toString();
        }
        return result;
    }

    private static String[][] getRunningServices(Context context){
        ActivityManager am=(ActivityManager) context.getSystemService(context.ACTIVITY_SERVICE);
        List<ActivityManager.RunningServiceInfo> runningApps = am.getRunningServices(Integer.MAX_VALUE);
        String[][] result= new String[runningApps.size()][5];
        for (int i =0; i<runningApps.size();i++){
            result[i][0]=runningApps.get(i).process;
            result[i][1]=Long.toString((SystemClock.elapsedRealtime()-runningApps.get(i).activeSince)/1000);
            int uid=runningApps.get(i).uid;
            result[i][2]=Integer.toString(uid);
            result[i][3]= Long.toString(TrafficStats.getUidRxBytes(uid));
            result[i][4]= Long.toString(TrafficStats.getUidTxBytes(uid));
        }
        return result;
    }

    protected static String getRunningServicesXML(Context context,String format) {
        StringBuilder result=new StringBuilder();
        for (String[] p : getRunningServices(context)) {
            result.append("<runservice uid=\""+p[2]+"\" uptime=\"" + p[1] + "\" rx=\"" + p[3] + "\" tx=\"" + p[4] +"\">" + p[0] + "</runservice>\n");
        }
        return result.toString();
    }


    protected static String getInstalledPackagesXML(Context context, String timeFormat) {
        StringBuilder result=new StringBuilder();
        for (String[] p : getInstalledPackages(context,timeFormat)) {
            result.append("<installedapp name=\"" + p[0] + "\" installed=\"" + p[1] + "\" version=\"" + p[2] + "\" uid=\"" + p[3] +"\">" + p[4] + "</installedapp>\n");
        }
        return result.toString();
    }
}