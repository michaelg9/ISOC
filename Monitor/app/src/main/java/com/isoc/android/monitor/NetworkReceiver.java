package com.isoc.android.monitor;

import android.content.BroadcastReceiver;
import android.content.Context;
import android.content.Intent;
import android.net.ConnectivityManager;
import android.net.NetworkInfo;
import android.util.Log;

public class NetworkReceiver extends BroadcastReceiver {
    public NetworkReceiver() {
    }


    //Deprecated method, implement type?
    @Override
    public void onReceive(Context context, Intent intent) {
        NetworkCapture.getTrafficStats(context,(NetworkInfo)intent.getParcelableExtra(ConnectivityManager.EXTRA_NETWORK_INFO));
   }
}

