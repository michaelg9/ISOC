package com.isoc.android.monitor;

import android.content.ContentValues;
import android.content.Context;
import android.content.SharedPreferences;
import android.database.Cursor;
import android.database.sqlite.SQLiteDatabase;
import android.provider.CallLog;
import android.telephony.TelephonyManager;

import com.google.i18n.phonenumbers.NumberParseException;
import com.google.i18n.phonenumbers.PhoneNumberUtil;
import com.google.i18n.phonenumbers.Phonenumber;

 /**
 * Use of external library libphonenumber
 * (https://github.com/googlei18n/libphonenumber/) to normalize
 * numbers. 
 * Call log monitoring is triggered by MyService
 * The time of the call field in our DB is unique so there are no duplicate
 * entries from re-runs. Phone numbers are not sent for privacy. Instead, the
 * table CallLogNumberReplacements matches each unique number (after it's been
 * normalized using libphonenumber) to an auto-increment integer, and this
 * integer masks the real number when we send the data to the server.
 **/

public class ContactsCapture {

    public static void getCallLog(Context context, SQLiteDatabase db) {
        // we need to get the last captured call's date from preferences file
        // "values"
        SharedPreferences preferences = context.getSharedPreferences(
                context.getString(R.string.shared_values_filename),
                Context.MODE_PRIVATE);
        // if there's a long different than 0, then we restore progress.
        long lastCallDate = preferences.getLong("lastCallDate", 0);
        String[] projection = new String[] { CallLog.Calls.NUMBER,
                CallLog.Calls.TYPE, CallLog.Calls.DURATION, CallLog.Calls.DATE,
                CallLog.Calls.CACHED_NAME };
        // querying the call log for calls made after lastCallDate
        Cursor cursor = context.getContentResolver().query(
                CallLog.Calls.CONTENT_URI, projection,
                projection[3] + ">" + lastCallDate, null,
                CallLog.Calls.DATE + " DESC");
        if (cursor == null) {
            return;
        }
        int numberIndex = cursor.getColumnIndex(CallLog.Calls.NUMBER);
        int typeIndex = cursor.getColumnIndex(CallLog.Calls.TYPE);
        int durationIndex = cursor.getColumnIndex(CallLog.Calls.DURATION);
        int dateIndex = cursor.getColumnIndex(CallLog.Calls.DATE);
        int nameIndex = cursor.getColumnIndex(CallLog.Calls.CACHED_NAME);
        // we need the sim operator county code to normalize numbers
        // appropriately
        String countryISO = ContactsCapture.getCountryISO(preferences, context);

        while (cursor.moveToNext()) {
            // log content values will be inserted into call log table
            ContentValues log = new ContentValues();
            // replacement content values will be inserted into the number
            // replacements table
            ContentValues replacement = new ContentValues();
            String number = formatNumber(cursor.getString(numberIndex),
                    countryISO);

            log.put(Database.DatabaseSchema.CallLog.COLUMN_NAME_NUMBER, number);
            replacement
                    .put(Database.DatabaseSchema.CallLogNumberReplacements.COLUMN_NAME_NUMBER,
                            number);
            log.put(Database.DatabaseSchema.CallLog.COLUMN_NAME_DURATION,
                    cursor.getString(durationIndex));
            log.put(Database.DatabaseSchema.CallLog.COLUMN_NAME_DATE,
                    TimeCapture.getGivenStringTime(cursor.getLong(dateIndex)));
            boolean saved = cursor.getString(nameIndex) != null;
            log.put(Database.DatabaseSchema.CallLog.COLUMN_NAME_SAVED,
                    Boolean.toString(saved));
            String callType = resolveType(cursor.getString(typeIndex));
            log.put(Database.DatabaseSchema.CallLog.COLUMN_NAME_TYPE, callType);
            // conflict ignore makes sure that duplicate rows (date not unique)
            // are not re-inserted and the "error is ignored"
            db.insertWithOnConflict(
                    Database.DatabaseSchema.CallLogNumberReplacements.TABLE_NAME,
                    null, replacement, SQLiteDatabase.CONFLICT_IGNORE);
            db.insertWithOnConflict(Database.DatabaseSchema.CallLog.TABLE_NAME,
                    null, log, SQLiteDatabase.CONFLICT_IGNORE);
        }
        // updating lastCallDate into the preferences file
        if (cursor.moveToFirst())
            preferences.edit()
                    .putLong("lastCallDate", cursor.getLong(dateIndex)).apply();

        cursor.close();
    }

    // retrieves the sim operator's country iso code from telephony manager and
    // then saves into the preferences file
    // for faster access
    public static String getCountryISO(SharedPreferences preferences,
            Context context) {
        String countryISO = preferences.getString("countryISO", null);
        if (countryISO == null) {
            TelephonyManager tm = (TelephonyManager) context
                    .getSystemService(Context.TELEPHONY_SERVICE);
            countryISO = tm.getSimCountryIso();
        }
        return countryISO;
    }

    // normalizes a number using libphonenumber. if the number can't be
    // formatted then it's used as-is
    // the reason for normalizing is that we match a number replacement integer
    // to a call log number. To do this
    // efficiently, each number must always be in the same format (for example,
    // country code is always present)
    public static String formatNumber(String number, String countryISO) {
        String result;
        PhoneNumberUtil phoneUtil = com.google.i18n.phonenumbers.PhoneNumberUtil
                .getInstance();
        Phonenumber.PhoneNumber myNum;
        try {
            myNum = phoneUtil.parse(number, countryISO);
            result = phoneUtil.format(myNum,
                    PhoneNumberUtil.PhoneNumberFormat.E164);
        } catch (NumberParseException e) {
            result = number;
        }
        return result;
    }

    //converts type String to a more readable form
    private static String resolveType(String type) {
        String callType;
        switch (Integer.parseInt(type)) {
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
        return callType;
    }

}