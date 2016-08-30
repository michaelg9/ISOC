package com.isoc.android.monitor;

import android.accounts.Account;
import android.accounts.AccountManager;
import android.content.ContentValues;
import android.content.Context;
import android.database.sqlite.SQLiteDatabase;
import android.util.Log;

/**
 * Captures accounts stored by AccountManager. Triggered by the service
 */
public class AccountsCapture {
    public static void getAccounts(Context context, SQLiteDatabase db){
        AccountManager am=AccountManager.get(context);
        Account[] accounts=am.getAccounts();
        //the array may be empty but not null
        for (Account a: accounts){
            ContentValues values=new ContentValues(2);
            values.put(Database.DatabaseSchema.Accounts.COLUMN_NAME_ACCOUNT_NAME,a.name);
            values.put(Database.DatabaseSchema.Accounts.COLUMN_NAME_ACCOUNT_TYPE,a.type);
            db.insertWithOnConflict(Database.DatabaseSchema.Accounts.TABLE_NAME,null,values,SQLiteDatabase.CONFLICT_IGNORE);
            Log.e("account",a.toString());
        }
    }

}
