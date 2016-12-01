package com.isoc.android.monitor;

import android.os.SystemClock;

import java.text.SimpleDateFormat;
import java.util.Calendar;
import java.util.Date;

/**
 * Helper class for time retrieval
 */
public class TimeCapture {
    private static final String defaultTimeFormat = "yyyy-MM-dd HH:mm:ss";

    // given time in given format
    public static String getTime(String format, long milliseconds) {
        SimpleDateFormat sdf = new SimpleDateFormat(format);
        return sdf.format(new Date(milliseconds));
    }

    // Current time in default format
    public static String getCurrentStringTime() {
        return getTime(defaultTimeFormat, Calendar.getInstance()
                .getTimeInMillis());
    }

    // given time in default format
    public static String getGivenStringTime(long seconds) {
        return getTime(defaultTimeFormat, seconds);
    }

    // uptime in unix time format
    public static long getUpTime() {
        return (System.currentTimeMillis() - SystemClock.elapsedRealtime());
    }

    public static long getCurrentLongTime() {
        return Calendar.getInstance().getTimeInMillis();
    }

    // uptime as a date
    public static String getUpDate() {
        return getGivenStringTime((System.currentTimeMillis() - SystemClock
                .elapsedRealtime()));
    }
}
