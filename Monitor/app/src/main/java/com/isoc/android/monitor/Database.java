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
        db.execSQL(DatabaseSchema.NetworkInterface.SQL_CREATE_TABLE);
        db.execSQL(DatabaseSchema.Actions.SQL_CREATE_TABLE);
        db.execSQL(DatabaseSchema.Sockets.SQL_CREATE_TABLE);
    }

    @Override
    public void onUpgrade(SQLiteDatabase db, int i, int i1) {
        db.execSQL(DatabaseSchema.SQL_DELETE_DB);
        onCreate(db);
    }

    public final static class DatabaseSchema {
        public final static String dbName = "Statistics.db";
        public final static String SQL_DELETE_DB = "DROP DATABASE " + dbName;

        private DatabaseSchema() {
        }

        public static abstract class InstalledPackages implements BaseColumns {
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
                    COLUMN_NAME_LABEL + " TEXT);";}

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
                    COLUMN_NAME_DATE + " INTEGER UNIQUE," +
                    COLUMN_NAME_NAME + " TEXT);";
        }

        public static abstract class Battery implements BaseColumns {
            public final static String TABLE_NAME = "Battery";
            public final static String COLUMN_NAME_DATE = "date";
            public static final String COLUMN_NAME_CHARGING = "charging";
            public static final String COLUMN_NAME_LEVEL = "level";
            public static final String COLUMN_NAME_TEMP = "temp";

            public static final String SQL_CREATE_TABLE = "CREATE TABLE " + TABLE_NAME + "(" +
                    _ID + " INTEGER PRIMARY KEY," +
                    COLUMN_NAME_DATE + " TEXT UNIQUE," +
                    COLUMN_NAME_CHARGING + " INTEGER," +
                    COLUMN_NAME_TEMP + " REAL," +
                    COLUMN_NAME_LEVEL + " INTEGER);";
        }

        public static abstract class NetworkInterface implements BaseColumns {
            public final static String TABLE_NAME = "NetIntf";
            public final static String COLUMN_NAME_ACTIVE = "active";
            public static final String COLUMN_NAME_TYPE = "type";
            public static final String COLUMN_NAME_SINCE = "since";
            public static final String COLUMN_NAME_TIME = "time";
            public static final String COLUMN_NAME_CURRENT_RX = "crx";
            public static final String COLUMN_NAME_CURRENT_TX = "ctx";
            public static final String COLUMN_NAME_TOTAL_RX = "trx";
            public static final String COLUMN_NAME_TOTAL_TX = "ttx";

            public static final String SQL_CREATE_TABLE = "CREATE TABLE " + TABLE_NAME + "(" +
                    _ID + " INTEGER PRIMARY KEY," +
                    COLUMN_NAME_ACTIVE + " TEXT," +
                    COLUMN_NAME_TYPE + " TEXT NOT NULL," +
                    COLUMN_NAME_TIME + " TEXT UNIQUE," +
                    COLUMN_NAME_SINCE + " TEXT DEFAULT \"" + TimeCapture.getUpDate() + "\"," +
                    COLUMN_NAME_CURRENT_RX + " INTEGER DEFAULT 0," +
                    COLUMN_NAME_CURRENT_TX + " INTEGER DEFAULT 0," +
                    COLUMN_NAME_TOTAL_RX + " INTEGER DEFAULT 0," +
                    COLUMN_NAME_TOTAL_TX + " INTEGER DEFAULT 0);";
        }

        public static abstract class Actions implements BaseColumns {
            public final static String TABLE_NAME = "Actions";
            public final static String COLUMN_NAME_ACTION = "action";
            public final static String COLUMN_NAME_DATE = "date";
            public final static String ACTION_BOOT="boot";
            public final static String ACTION_SHUTDOWN="shutdown";
            public final static String ACTION_REBOOT="reboot";
            public final static String ACTION_AIRPLANE_ON="airplaneOn";
            public final static String ACTION_AIRPLANE_OFF="airplaneOff";

            public static final String SQL_CREATE_TABLE = "CREATE TABLE " + TABLE_NAME + "(" +
                    _ID + " INTEGER PRIMARY KEY," +
                    COLUMN_NAME_ACTION + " TEXT," +
                    COLUMN_NAME_DATE + " TEXT);";
        }


        public static abstract class Sockets implements BaseColumns {
            public final static String TABLE_NAME = "Sockets";
            public final static String COLUMN_NAME_TYPE = "type";
            public final static String COLUMN_NAME_DATE = "date";
            public final static String COLUMN_NAME_UID = "uid";
            public final static String COLUMN_NAME_STATUS = "status";
            public final static String COLUMN_NAME_LIP = "lip";
            public final static String COLUMN_NAME_LPORT = "lport";
            public final static String COLUMN_NAME_RIP = "rip";
            public final static String COLUMN_NAME_RPORT = "rport";

            public static final String SQL_CREATE_TABLE = "CREATE TABLE " + TABLE_NAME + "(" +
                    _ID + " INTEGER PRIMARY KEY," +
                    COLUMN_NAME_TYPE + " TEXT," +
                    COLUMN_NAME_UID + " TEXT," +
                    COLUMN_NAME_STATUS + " TEXT," +
                    COLUMN_NAME_LIP + " TEXT," +
                    COLUMN_NAME_LPORT + " TEXT," +
                    COLUMN_NAME_RIP + " TEXT," +
                    COLUMN_NAME_RPORT + " TEXT," +
                    COLUMN_NAME_DATE + " TEXT);";
        }
    }
}