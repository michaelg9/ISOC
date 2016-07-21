package com.isoc.android.monitor;

import android.content.Context;
import android.content.Intent;
import android.support.v4.content.WakefulBroadcastReceiver;

public class ServiceReceiver extends WakefulBroadcastReceiver {
    public ServiceReceiver() {
    }

    @Override
    public void onReceive(Context context, Intent intent) {
        startWakefulService(context,new Intent(context,MyService.class));
    }
}
