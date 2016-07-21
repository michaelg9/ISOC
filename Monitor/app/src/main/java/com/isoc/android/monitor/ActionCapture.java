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
        Cursor cursor = db.query(Database.DatabaseSchema.Actions.TABLE_NAME,null,null,null,null,null,null);
        StringBuilder sb=new StringBuilder();
        int dateIndex = cursor.getColumnIndex(Database.DatabaseSchema.Actions.COLUMN_NAME_DATE);
        int actionIndex = cursor.getColumnIndex(Database.DatabaseSchema.Actions.COLUMN_NAME_ACTION);
        while (cursor.moveToNext()) {
            sb.append("<action date=\"" + cursor.getString(dateIndex) +"\">" + cursor.getString(actionIndex) + "</action>\n");
        }
        cursor.close();
        return sb.toString();
    }
}
