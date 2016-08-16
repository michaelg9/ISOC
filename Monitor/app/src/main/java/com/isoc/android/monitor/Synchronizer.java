package com.isoc.android.monitor;

import android.accounts.Account;
import android.content.AbstractThreadedSyncAdapter;
import android.content.ContentProviderClient;
import android.content.Context;
import android.content.SyncResult;
import android.os.Bundle;

/**
 * Created by me on 16/08/16.
 */
public class Synchronizer extends AbstractThreadedSyncAdapter {
    public Synchronizer(Context context, boolean autoInitialize) {
        super(context, autoInitialize);
    }

    public Synchronizer(Context context, boolean autoInitialize, boolean allowParallelSyncs) {
        super(context, autoInitialize, allowParallelSyncs);
    }

    @Override
    public void onPerformSync(Account account, Bundle bundle, String s, ContentProviderClient contentProviderClient, SyncResult syncResult) {
        
    }
}
