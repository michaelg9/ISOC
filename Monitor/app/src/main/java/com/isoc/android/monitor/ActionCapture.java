package com.isoc.android.monitor;

import android.content.ContentValues;
import android.content.Context;
import android.database.sqlite.SQLiteDatabase;

/**
 * Actions are captured only when they are broadcasted, not on intervals.
 * Look at the database schema for the captured actions.
 */
public class ActionCapture {

    public static void getAction(Context context,String action){
        String date = TimeCapture.getTime();
        ContentValues values=new ContentValues();
        values.put(Database.DatabaseSchema.Actions.COLUMN_NAME_ACTION,action);
        values.put(Database.DatabaseSchema.Actions.COLUMN_NAME_DATE,date);
        SQLiteDatabase db =new Database(context).getWritableDatabase();
        db.insert(Database.DatabaseSchema.Actions.TABLE_NAME,null,values);
        db.close();
    }
}
