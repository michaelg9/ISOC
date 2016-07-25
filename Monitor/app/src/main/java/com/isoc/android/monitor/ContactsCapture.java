package com.isoc.android.monitor;

import android.content.ContentValues;
import android.content.Context;
import android.database.Cursor;
import android.database.sqlite.SQLiteDatabase;
import android.provider.CallLog;

/**

 */

public class ContactsCapture {

    //TO DO: get only calls earlier than last captured call
    protected static void getCallLog(Context context,SQLiteDatabase db) {
        String[] projection=new String[]{CallLog.Calls.NUMBER,CallLog.Calls.TYPE,CallLog.Calls.DURATION,CallLog.Calls.DATE,CallLog.Calls.CACHED_NAME};
        Cursor cursor = context.getContentResolver().query(CallLog.Calls.CONTENT_URI, projection, null, null, CallLog.Calls.DATE+ " DESC");
        if (cursor==null) return;
        int number = cursor.getColumnIndex(CallLog.Calls.NUMBER);
        int type = cursor.getColumnIndex(CallLog.Calls.TYPE);
        int duration = cursor.getColumnIndex(CallLog.Calls.DURATION);
        int date = cursor.getColumnIndex(CallLog.Calls.DATE);
        int name = cursor.getColumnIndex(CallLog.Calls.CACHED_NAME);


        while (cursor.moveToNext()) {
            ContentValues log = new ContentValues();
            ContentValues replacement=new ContentValues();
            String formattedNumber = cursor.getString(number).replace("+","00");        //still, the same number without country code is taken as a different number...
            log.put(Database.DatabaseSchema.CallLog.COLUMN_NAME_NUMBER,formattedNumber);
            replacement.put(Database.DatabaseSchema.CallLogNumberReplacements.COLUMN_NAME_NUMBER,formattedNumber);
            log.put(Database.DatabaseSchema.CallLog.COLUMN_NAME_DURATION,cursor.getString(duration));
            log.put(Database.DatabaseSchema.CallLog.COLUMN_NAME_DATE,cursor.getString(date));
            String contactName=cursor.getString(name);
            boolean saved=contactName!=null;
            log.put(Database.DatabaseSchema.CallLog.COLUMN_NAME_SAVED,Boolean.toString(saved));

            String callType;
            switch (Integer.parseInt(cursor.getString(type))) {
                case CallLog.Calls.OUTGOING_TYPE:
                    callType = "Outgoing";
                    break;
                case CallLog.Calls.INCOMING_TYPE:
                    callType = "Incoming";
                    break;
                case CallLog.Calls.MISSED_TYPE:
                    callType = "Missed";
                    break;
                default:
                    callType = "Unknown";
                    break;
            }
            log.put(Database.DatabaseSchema.CallLog.COLUMN_NAME_TYPE,callType);
            db.insertWithOnConflict(Database.DatabaseSchema.CallLog.TABLE_NAME,null,log,SQLiteDatabase.CONFLICT_IGNORE);
            db.insertWithOnConflict(Database.DatabaseSchema.CallLogNumberReplacements.TABLE_NAME,null,replacement,SQLiteDatabase.CONFLICT_IGNORE);
        }
        cursor.close();
    }

    protected static String getCallXML(SQLiteDatabase db) {
        String query=String.format("SELECT %s,%s,%s,%s,R._id FROM %s AS C JOIN %s AS R USING (%s)",
                Database.DatabaseSchema.CallLog.COLUMN_NAME_DATE,Database.DatabaseSchema.CallLog.COLUMN_NAME_TYPE,
                Database.DatabaseSchema.CallLog.COLUMN_NAME_DURATION,Database.DatabaseSchema.CallLog.COLUMN_NAME_SAVED,
                Database.DatabaseSchema.CallLog.TABLE_NAME,Database.DatabaseSchema.CallLogNumberReplacements.TABLE_NAME,
                Database.DatabaseSchema.CallLog.COLUMN_NAME_NUMBER);
        Cursor cursor = db.rawQuery(query,null);
        StringBuilder result=new StringBuilder();
        int dateIndex= cursor.getColumnIndex(Database.DatabaseSchema.CallLog.COLUMN_NAME_DATE);
        int typeIndex=cursor.getColumnIndex(Database.DatabaseSchema.CallLog.COLUMN_NAME_TYPE);
        int durationIndex=cursor.getColumnIndex(Database.DatabaseSchema.CallLog.COLUMN_NAME_DURATION);
        int savedIndex =cursor.getColumnIndex(Database.DatabaseSchema.CallLog.COLUMN_NAME_SAVED);
        int iIndex=cursor.getColumnIndex(Database.DatabaseSchema.CallLogNumberReplacements._ID);

        while (cursor.moveToNext()) {
                result.append("<call time=\"" + TimeCapture.getTime(cursor.getLong(dateIndex)) +"\" type=\"" + cursor.getString(typeIndex) +
                        "\" duration=\"" + cursor.getString(durationIndex) + "\" saved=\"" + cursor.getString(savedIndex) + "\">" +
                        cursor.getString(iIndex) + "</call>\n");
        }
        cursor.close();
        return result.toString();
    }

    protected static String getCallXML2(SQLiteDatabase db) {
        String query=String.format("SELECT %s,%s,%s,%s,R._id FROM %s AS C JOIN %s AS R USING (%s)",
                Database.DatabaseSchema.CallLog.COLUMN_NAME_DATE,Database.DatabaseSchema.CallLog.COLUMN_NAME_TYPE,
                Database.DatabaseSchema.CallLog.COLUMN_NAME_DURATION,Database.DatabaseSchema.CallLog.COLUMN_NAME_SAVED,
                Database.DatabaseSchema.CallLog.TABLE_NAME,Database.DatabaseSchema.CallLogNumberReplacements.TABLE_NAME,
                Database.DatabaseSchema.CallLog.COLUMN_NAME_NUMBER);
        Cursor cursor = db.rawQuery(query,null);
        String result = XMLProduce.tableToXML(cursor, Database.DatabaseSchema.CallLog.TAG,Database.DatabaseSchema.CallLogNumberReplacements._ID);
        cursor.close();
        return result;
    }
}
