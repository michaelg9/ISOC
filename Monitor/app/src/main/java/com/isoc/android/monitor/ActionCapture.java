package com.isoc.android.monitor;

import android.content.ContentValues;
import android.content.Context;
import android.database.Cursor;
import android.database.sqlite.SQLiteDatabase;

/**
 * Created by me on 17/07/16.
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

    public static String getActionsXML(SQLiteDatabase db){
        String[] projection = new String[]{Database.DatabaseSchema.Actions.COLUMN_NAME_ACTION,
                Database.DatabaseSchema.Actions.COLUMN_NAME_DATE};
        Cursor cursor = db.query(Database.DatabaseSchema.Actions.TABLE_NAME,projection,null,null,null,null,null);
        String result= XMLProduce.tableToXML(cursor,Database.DatabaseSchema.Actions.TAG,Database.DatabaseSchema.Actions.COLUMN_NAME_ACTION);
        cursor.close();
        return result;
    }
}
