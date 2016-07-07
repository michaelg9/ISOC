package com.isoc.android.monitor;

import java.text.SimpleDateFormat;
import java.util.Calendar;
import java.util.Date;

/**
 * Created by maik on 1/7/2016.
 */
public class TimeCapture {

    protected static String getTime(String format,long seconds) {
        SimpleDateFormat sdf = new SimpleDateFormat(format);
        return sdf.format(new Date(seconds));
    }

    protected static String getTime(String format){
        long seconds = Calendar.getInstance().getTimeInMillis();
        return getTime(format,seconds);
    }
}
