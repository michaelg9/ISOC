package com.isoc.android.monitor;

import android.accounts.Account;
import android.accounts.AccountManager;
import android.accounts.AuthenticatorException;
import android.accounts.OperationCanceledException;
import android.content.AbstractThreadedSyncAdapter;
import android.content.ContentProviderClient;
import android.content.Context;
import android.content.SyncResult;
import android.database.sqlite.SQLiteDatabase;
import android.os.Bundle;
import android.util.Log;

import java.io.IOException;

/**
 * Synchronizer module used by Android to automatically send our data 
 * Run by default once every 24h
 */
public class Synchronizer extends AbstractThreadedSyncAdapter {
    private Context context;

    public Synchronizer(Context context, boolean autoInitialize) {
        super(context, autoInitialize);
        this.context = context;
    }

    public Synchronizer(Context context, boolean autoInitialize,
            boolean allowParallelSyncs) {
        super(context, autoInitialize, allowParallelSyncs);
        this.context = context;
    }

    //called when asked to sync
    @Override
    public void onPerformSync(Account account, Bundle bundle, String s,
            ContentProviderClient contentProviderClient, SyncResult syncResult) {
        AccountManager am = AccountManager.get(context);
        String refreshTokenType = context.getString(R.string.token_refresh);
        Database dbHelper = new Database(context);
        SQLiteDatabase db = dbHelper.getWritableDatabase();
        try {
            // we're already on a separate thread than gui, so synchronous token
            // request
            String refreshToken = am.blockingGetAuthToken(account,
                    refreshTokenType, true);
            /*
             * if there's no refresh token saved, don't attempt to send, Account
             * manager is going to notify the user //that re-authentication is
             * needed
             */
            if (refreshToken == null)
                return;
            ServerCommunication serverCommunication = new ServerCommunication(
                    context);
            // request new access token
            String accessToken = serverCommunication
                    .getAccessToken(refreshToken);
            String deviceID = am.getUserData(account,
                    context.getString(R.string.am_deviceID));
            if (deviceID == null) {
                Log.e("ERROR", "deviceID null, return...");
                return;
            }
            Log.e("devIDtoSEND", deviceID);
            String[] sendResponse = serverCommunication.sendData(accessToken,
                    new XMLProduce(context, db).getXML(Integer
                            .parseInt(deviceID)));

            /*
             * We request a new refresh token if: 
             * a)the request was successful and it's been more than 
             * 5 days(432000000 mills) since the last refresh token renewal
             * b)accessToken is null (the request for a new access token failed)
             * c)the response code of the request to send the data is 401 (Unauthorized)
             */
            if ((TimeCapture.getCurrentLongTime()
                    - Long.parseLong(am.getUserData(account,
                            context.getString(R.string.am_refreshDateKey))) > 432000000)
                    || accessToken == null || sendResponse[0].equals("401")) {
                String newRefreshToken = serverCommunication
                        .refreshRefreshToken(refreshToken);
                if (newRefreshToken != null) {
                    am.invalidateAuthToken(refreshTokenType, refreshToken);
                    am.setAuthToken(account, refreshTokenType, newRefreshToken);
                    am.setUserData(account,
                            context.getString(R.string.am_refreshDateKey),
                            Long.toString(TimeCapture.getCurrentLongTime()));
                }
            }

            // if request was successful, we need to mark as sent every database
            // field that was marked as unsent
            if (sendResponse[0].equals("200")) {
                dbHelper.markSend(db);
            }

        } catch (OperationCanceledException e) {
            e.printStackTrace();
        } catch (IOException e) {
            e.printStackTrace();
        } catch (AuthenticatorException e) {
            e.printStackTrace();
        } finally {
            db.close();
        }
    }
}