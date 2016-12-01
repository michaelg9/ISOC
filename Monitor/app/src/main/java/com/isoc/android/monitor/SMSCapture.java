package com.isoc.android.monitor;

import android.content.ContentValues;
import android.content.Context;
import android.content.SharedPreferences;
import android.database.Cursor;
import android.database.sqlite.SQLiteDatabase;
import android.net.Uri;
import android.os.Build;
import android.provider.Telephony;

/*****
 * Use of external library libphonenumber
 * (https://github.com/googlei18n/libphonenumber/) to normalize
 * numbers.
 * Captures sms messages, triggered by MyService. 
 * Normalizes the number first to match to the 
 * replacement integer.
 * Only new sms are captured each time 
 * (with date after the saved captured date)
 */
public class SMSCapture {
    //Called by MyService
    //Depending on OS version, uses Telephony.SMS API or not
    protected static void getSMS(Context context, SQLiteDatabase db) {
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.KITKAT) {
            // we can use Telephony.SMS API
            String[] projection = new String[] { Telephony.Sms.ADDRESS,
                    Telephony.Sms.DATE, Telephony.Sms.READ, Telephony.Sms.TYPE };
            getSMSLog(context, db, projection, Telephony.Sms.CONTENT_URI);
        } else {
            // manually query content://sms/ ------>What if not existent?
            String[] projection = new String[] { "address", "date", "read",
                    "type" };
            getSMSLog(context, db, projection, Uri.parse("content://sms/"));
        }
    }
    
    //queries the content provider
    private static void getSMSLog(Context context, SQLiteDatabase db,
            String[] projection, Uri uri) {
        // restore last captured sms's date to restore progress
        SharedPreferences preferences = context.getSharedPreferences(
                context.getString(R.string.shared_values_filename),
                Context.MODE_PRIVATE);
        long lastSMSDate = preferences.getLong("lastSMSDate", 0);
        // get all sms with date after lastSMSDate
        Cursor cursor = context.getContentResolver().query(uri, projection,
                projection[1] + ">" + lastSMSDate, null, null);
        if (cursor == null)
            return;

        int[] indexes = new int[projection.length];
        for (int i = 0; i < projection.length; i++)
            indexes[i] = cursor.getColumnIndex(projection[i]);
        // we need country iso to normilize numbers, just as we do in Contact
        // capturing
        String countryISO = ContactsCapture.getCountryISO(preferences, context);
        while (cursor.moveToNext()) {
            ContentValues log = new ContentValues();
            ContentValues replacement = new ContentValues();
            String number = ContactsCapture.formatNumber(
                    cursor.getString(indexes[0]), countryISO);
            log.put(Database.DatabaseSchema.SMSLog.COLUMN_NAME_NUMBER, number);
            replacement
                    .put(Database.DatabaseSchema.CallLogNumberReplacements.COLUMN_NAME_NUMBER,
                            number);
            log.put(Database.DatabaseSchema.SMSLog.COLUMN_NAME_DATE,
                    TimeCapture.getGivenStringTime(cursor.getLong(indexes[1])));
            log.put(Database.DatabaseSchema.SMSLog.COLUMN_NAME_READ,
                    isRead(cursor.getString(indexes[2])));
            log.put(Database.DatabaseSchema.SMSLog.COLUMN_NAME_TYPE,
                    resolveFolder(cursor.getString(indexes[3])));
            // conflict ignore makes sure that duplicate rows (date not unique)
            // are not re-inserted
            db.insertWithOnConflict(
                    Database.DatabaseSchema.CallLogNumberReplacements.TABLE_NAME,
                    null, replacement, SQLiteDatabase.CONFLICT_IGNORE);
            db.insertWithOnConflict(Database.DatabaseSchema.SMSLog.TABLE_NAME,
                    null, log, SQLiteDatabase.CONFLICT_IGNORE);
        }
        // saving progress
        if (cursor.moveToFirst())
            preferences.edit()
                    .putLong("lastSMSDate", cursor.getLong(indexes[1])).apply();
        cursor.close();
    }
    
    //converts read status string to a more readble one 
    private static String isRead(String flag) {
        String read;
        switch (Integer.parseInt(flag)) {
        case 0: // not read
            read = "false";
            break;
        case 1: // read
            read = "true";
            break;
        default:
            read = "unknown";
        }
        return read;
    }
    
    //converts folder string to a more readble one
    private static String resolveFolder(String f) {
        String folder;
        switch (Integer.parseInt(f)) {
        case 0: // Telephony.Sms.MESSAGE_TYPE_ALL
            folder = "all";
            break;
        case 1: // Telephony.Sms.MESSAGE_TYPE_INBOX
            folder = "inbox";
            break;
        case 2: // Telephony.Sms.MESSAGE_TYPE_SENT
            folder = "sent";
            break;
        case 3: // Telephony.Sms.MESSAGE_TYPE_DRAFT
            folder = "draft";
            break;
        case 4: // Telephony.Sms.MESSAGE_TYPE_OUTBOX
            folder = "outbox";
            break;
        case 5: // Telephony.Sms.MESSAGE_TYPE_FAILED
            folder = "failed";
            break;
        case 6: // Telephony.Sms.MESSAGE_TYPE_QUEUED
            folder = "queued";
            break;
        default:
            folder = "unknown";
            break;
        }
        return folder;
    }

}