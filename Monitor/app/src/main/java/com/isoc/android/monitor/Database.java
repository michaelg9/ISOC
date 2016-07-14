package com.isoc.android.monitor;

import android.content.Context;
import android.database.sqlite.SQLiteDatabase;
import android.database.sqlite.SQLiteOpenHelper;
import android.provider.BaseColumns;

/**
 * Created by me on 11/07/16.
 */
public class Database extends SQLiteOpenHelper {

    public Database(Context context) {
        super(context, DatabaseSchema.dbName, null, 1);
    }

    @Override
    public void onCreate(SQLiteDatabase db) {
        db.execSQL(DatabaseSchema.RunningServices.SQL_CREATE_TABLE);
        db.execSQL(DatabaseSchema.InstalledPackages.SQL_CREATE_TABLE);
        db.execSQL(DatabaseSchema.CallLog.SQL_CREATE_TABLE);
        db.execSQL(DatabaseSchema.Battery.SQL_CREATE_TABLE);
    }

    @Override
    public void onUpgrade(SQLiteDatabase db, int i, int i1) {
        db.execSQL(DatabaseSchema.SQL_DELETE_DB);
        onCreate(db);
    }

    public final static class DatabaseSchema {
        public final static String dbName = "Statistics.db";
        public final static String SQL_DELETE_DB="DROP DATABASE "+dbName;

        private DatabaseSchema() {
        }

        public static abstract class InstalledPackages implements BaseColumns{
            public final static String TABLE_NAME = "InstalledPackages";
            public final static String COLUMN_NAME_PACKAGE_NAME = "pckgName";
            public static final String COLUMN_NAME_INSTALLED_DATE = "InstDate";
            public static final String COLUMN_NAME_VERSION = "version";
            public static final String COLUMN_NAME_UID = "uid";
            public static final String COLUMN_NAME_LABEL = "label";

            public static final String SQL_CREATE_TABLE = "CREATE TABLE " + TABLE_NAME + "(" +
                    _ID + " INTEGER PRIMARY KEY," +
                    COLUMN_NAME_PACKAGE_NAME + " TEXT UNIQUE," +
                    COLUMN_NAME_INSTALLED_DATE + " INTEGER," +
                    COLUMN_NAME_VERSION + " TEXT," +
                    COLUMN_NAME_UID + " TEXT," +
                    COLUMN_NAME_LABEL + " TEXT);";
        }

        public static abstract class RunningServices implements BaseColumns {
            public final static String TABLE_NAME = "RunningServices";
            public final static String COLUMN_NAME_PROCESS_NAME = "pckgName";
            public static final String COLUMN_NAME_UP_TIME = "uptime";
            public static final String COLUMN_NAME_UID = "uid";
            public static final String COLUMN_NAME_RX = "rx";
            public static final String COLUMN_NAME_TX = "tx";

            public static final String SQL_CREATE_TABLE = "CREATE TABLE " + TABLE_NAME + "(" +
                    _ID + " INTEGER PRIMARY KEY," +
                    COLUMN_NAME_PROCESS_NAME + " TEXT UNIQUE," +
                    COLUMN_NAME_UP_TIME + " TEXT," +
                    COLUMN_NAME_UID + " TEXT," +
                    COLUMN_NAME_RX + " TEXT," +
                    COLUMN_NAME_TX + " TEXT);";
        }

        public static abstract class CallLog implements BaseColumns {
            public final static String TABLE_NAME = "CallLog";
            public final static String COLUMN_NAME_NUMBER = "number";
            public static final String COLUMN_NAME_TYPE = "type";
            public static final String COLUMN_NAME_DURATION = "duration";
            public static final String COLUMN_NAME_DATE = "date";
            public static final String COLUMN_NAME_NAME = "name";

            public static final String SQL_CREATE_TABLE = "CREATE TABLE " + TABLE_NAME + "(" +
                    _ID + " INTEGER PRIMARY KEY," +
                    COLUMN_NAME_NUMBER + " TEXT," +
                    COLUMN_NAME_TYPE + " TEXT," +
                    COLUMN_NAME_DURATION + " TEXT," +
                    COLUMN_NAME_DATE + " INTEGER," +
                    COLUMN_NAME_NAME + " TEXT);";
        }

        public static abstract class Battery implements BaseColumns{
            public final static String TABLE_NAME = "Battery";
            public final static String COLUMN_NAME_DATE = "date";
            public static final String COLUMN_NAME_CHARGING = "charging";
            public static final String COLUMN_NAME_LEVEL = "level";

            public static final String SQL_CREATE_TABLE = "CREATE TABLE " + TABLE_NAME + "(" +
                    _ID + " INTEGER PRIMARY KEY," +
                    COLUMN_NAME_DATE + " TEXT UNIQUE," +
                    COLUMN_NAME_CHARGING + " INTEGER," +
                    COLUMN_NAME_LEVEL + " INTEGER);";
        }
    }

    public static abstract class WifiNetworkInterface{
        public final static String PREFERENCES_FILENAME = "com.isoc.monitor.wifi";
        public final static String KEY_ACTIVE = "active";
        public static final String KEY_INTF_NAME = "intfName";
        public static final String KEY_SINCE = "since";
        public static final String KEY_CURRENT_RX = "crx";
        public static final String KEY_CURRENT_TX = "ctx";
        public static final String KEY_TOTAL_RX = "trx";
        public static final String KEY_TOTAL_TX = "ttx";
    }

    public static abstract class MobileNetworkInterface{
        public final static String PREFERENCES_FILENAME = "com.isoc.monitor.mobile";
        public final static String KEY_ACTIVE = "active";
        //public static final String KEY_INTF_NAME = "intfName";
        public static final String KEY_SINCE = "since";
        public static final String KEY_CURRENT_RX = "crx";
        public static final String KEY_CURRENT_TX = "ctx";
        public static final String KEY_TOTAL_RX = "trx";
        public static final String KEY_TOTAL_TX = "ttx";
    }
}