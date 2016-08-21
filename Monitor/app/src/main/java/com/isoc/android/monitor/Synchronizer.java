package com.isoc.android.monitor;

import android.accounts.Account;
import android.accounts.AccountManager;
import android.accounts.AuthenticatorException;
import android.accounts.OperationCanceledException;
import android.content.AbstractThreadedSyncAdapter;
import android.content.ContentProviderClient;
import android.content.Context;
import android.content.SyncResult;
import android.os.Bundle;

import java.io.IOException;

/**
 * Synchronizer module used by Android to automatically send our data
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
        AccountManager am=AccountManager.get(getContext());
        try {
            String refreshToken=am.blockingGetAuthToken(account,getContext().getString(R.string.token_refresh),true);
            ServerCommunication serverCommunication= new ServerCommunication(getContext());
            String accessToken=serverCommunication.getAccessToken(refreshToken);
            serverCommunication.sendData(accessToken,new XMLProduce(getContext()).getXML());
        } catch (OperationCanceledException e) {
            e.printStackTrace();
        } catch (IOException e) {
            e.printStackTrace();
        } catch (AuthenticatorException e) {
            e.printStackTrace();
        }
    }

}
