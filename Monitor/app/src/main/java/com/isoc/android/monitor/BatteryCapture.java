package com.isoc.android.monitor;

import android.content.Context;
import android.content.Intent;
import android.content.IntentFilter;
import android.os.BatteryManager;

/**
 * Created by maik on 1/7/2016.
 */
public class BatteryCapture {

    private static int[] getBatteryStats(Context context) {
        IntentFilter iFilter = new IntentFilter(Intent.ACTION_BATTERY_CHANGED);
        Intent battery = context.registerReceiver(null, iFilter);
        int plugged = battery.getIntExtra(BatteryManager.EXTRA_PLUGGED, -1);
        int level = battery.getIntExtra(BatteryManager.EXTRA_LEVEL, -1);
        return new int[]{level, plugged};
    }

    protected static String getBatteryXML(Context context,String timeFormat) {
        int[] batteryStats = getBatteryStats(context);
        boolean charging = batteryStats[1]!=0;
        return("<battery time=\""+TimeCapture.getTime(timeFormat)+"\" charging=\""+charging+"\">"+batteryStats[0]+"</battery>\n");
    }

}
