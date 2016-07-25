package com.isoc.android.monitor;

import android.database.Cursor;

/**
 * Created by me on 25/07/16.
 */
public class XMLProduce {

    public static String tableToXML(Cursor c, String tag, String text){
        if (c==null) return null;
        String[] attributes=c.getColumnNames();
        StringBuilder result=new StringBuilder();
        String ending="//>\n";
        while (c.moveToNext()){
            result.append("<"+tag);
            for (String attribute : attributes){
                if (attribute.equals(text)) ending=">"+c.getString(c.getColumnIndex(text))+"</"+tag+">\n";
                else result.append(" "+attribute+"=\""+c.getString(c.getColumnIndex(attribute))+"\"");
            }
            result.append(ending);
        }
        c.close();
        return result.toString();
    }
}
