package com.isoc.android.monitor;

import android.os.SystemClock;

import java.text.SimpleDateFormat;
import java.util.Calendar;
import java.util.Date;

/**
 * Created by maik on 1/7/2016.
 */
public class TimeCapture {
    private static final String defaultTimeFormat="yyyy-MM-dd HH:mm:ss";

    //given time in given format
    protected static String getTime(String format,long milliseconds) {
        SimpleDateFormat sdf = new SimpleDateFormat(format);
        return sdf.format(new Date(milliseconds));
    }

    //Current time in given format
    protected static String getTime(String format){ return getTime(format,Calendar.getInstance().getTimeInMillis());}

    //Current time in default format
    protected static String getTime(){ return getTime(defaultTimeFormat,Calendar.getInstance().getTimeInMillis());}

    //given time in default format
    protected static String getTime(long seconds) {
        return getTime(defaultTimeFormat,seconds);
    }

    //uptime in unix time format
    protected static long getUpTime(){return (System.currentTimeMillis()-SystemClock.elapsedRealtime());}

    //uptime as a date
    protected static String getUpDate(){return getTime((System.currentTimeMillis()-SystemClock.elapsedRealtime()));}
}
