package com.isoc.android.monitor;

import android.content.ContentValues;
import android.database.Cursor;
import android.net.Uri;

/*
 * Stub content provider, needed for the sync adapter implementation.
 * TODO: implement it!
 */

public class DataProvider extends android.content.ContentProvider {
    private static final int accounts = 0;
    private static final int installedPackages = 5;
    private static final int runningServices = 10;
    private static final int smsLog = 15;
    private static final int callLog = 20;
    private static final int battery = 25;
    private static final int networkInterface = 30;
    private static final int wifiAp = 35;
    private static final int actions = 40;
    private static final int sockets = 45;

    public DataProvider() {
    }

    @Override
    public int delete(Uri uri, String selection, String[] selectionArgs) {
        return 0;
    }

    @Override
    public String getType(Uri uri) {
        return null;
    }

    @Override
    public Uri insert(Uri uri, ContentValues values) {
        return null;
    }

    @Override
    public boolean onCreate() {
        return false;
    }

    @Override
    public Cursor query(Uri uri, String[] projection, String selection,
            String[] selectionArgs, String sortOrder) {
        return null;
    }

    @Override
    public int update(Uri uri, ContentValues values, String selection,
            String[] selectionArgs) {
        return 0;
    }
}
