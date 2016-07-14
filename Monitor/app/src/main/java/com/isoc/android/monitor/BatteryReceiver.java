package com.isoc.android.monitor;

import android.content.BroadcastReceiver;
import android.content.Context;
import android.content.Intent;

public class BatteryReceiver extends BroadcastReceiver {
    public BatteryReceiver() {
    }

    @Override
    public void onReceive(Context context, Intent intent) {
        BatteryCapture.getBatteryStats(context);
    }
}
