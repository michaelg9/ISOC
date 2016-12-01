package com.isoc.android.monitor;

import android.content.BroadcastReceiver;
import android.content.Context;
import android.content.Intent;
import android.net.ConnectivityManager;
import android.net.NetworkInfo;

/*
 * Triggered by connectivity changes
 * BUG: some devices echo the same broadcast intent for a few seconds after the first trigger.
 * Some are caught by the unique date field but not always.... 
 * Possible solution: Throttle them:
 * https://plus.google.com/107087775691692585833/posts/FCmYXsWptiQ
 * BUG: Use of deprecated extra_network_info. implement type?
 */
public class NetworkReceiver extends BroadcastReceiver {
    public NetworkReceiver() {
    }

    @Override
    public void onReceive(Context context, Intent intent) {
        NetworkCapture.getTrafficStats(context, (NetworkInfo) intent
                .getParcelableExtra(ConnectivityManager.EXTRA_NETWORK_INFO));
    }
}
