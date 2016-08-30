package com.isoc.android.monitor;

import android.content.Context;
import android.content.Intent;
import android.content.pm.PackageManager;
import android.content.pm.ResolveInfo;
import android.database.Cursor;
import android.database.sqlite.SQLiteDatabase;
import android.net.Uri;
import android.os.Build;
import android.telephony.TelephonyManager;

import java.util.ArrayList;
import java.util.TimeZone;

/**
 * Methods to produce and send the XML string.
 * A db connection opens on creation of an XMLProduce object. Since the connection needs to close before the object is
 * garbage collected, I decided to close it upon retrieval of the results. For this reason, a new object needs to be created
 * each time. Otherwise there can be a finish method that closes the db connection.
 */
public class XMLProduce {
    private static SQLiteDatabase db;
    private StringBuilder xmlString;
    private static Context context;

    public XMLProduce(Context c, SQLiteDatabase database) {
        context = c;
        db = database;
        xmlString=new StringBuilder();
    }

    private String cursorToXML(Cursor c, String tag, String text) {
        //all the columns of the cursor will be parsed as attributes, except for the column named text. The name of the attribute is the name of the column
        if (c == null) return null;
        String[] attributes = c.getColumnNames();
        StringBuilder result = new StringBuilder();
        String ending = "/>\n";
        while (c.moveToNext()) {
            result.append("<" + tag);
            for (String attribute : attributes) {
                if (attribute.equals(text))
                    ending = ">" + c.getString(c.getColumnIndex(text)) + "</" + tag + ">\n";
                else
                    result.append(" " + attribute + "=\"" + c.getString(c.getColumnIndex(attribute)) + "\"");
            }
            result.append(ending);
        }
        c.close();
        return result.toString();
    }

    private void getActions() {
        String[] projection = new String[]{Database.DatabaseSchema.Actions.COLUMN_NAME_ACTION,
                Database.DatabaseSchema.Actions.COLUMN_NAME_DATE};
        Cursor cursor = db.query(Database.DatabaseSchema.Actions.TABLE_NAME, projection,
                Database.DatabaseSchema.GLOBAL_COLUMN_NAME_SENT + "='FALSE'", null, null, null, null);
        String result = cursorToXML(cursor, Database.DatabaseSchema.Actions.TAG, Database.DatabaseSchema.Actions.COLUMN_NAME_ACTION);
        xmlString.append(result);
    }

    private void getBattery() {
        String[] projection = new String[]{Database.DatabaseSchema.Battery.COLUMN_NAME_LEVEL,
                Database.DatabaseSchema.Battery.COLUMN_NAME_TIME,
                Database.DatabaseSchema.Battery.COLUMN_NAME_CHARGING,
                Database.DatabaseSchema.Battery.COLUMN_NAME_TEMP};
        Cursor cursor = db.query(Database.DatabaseSchema.Battery.TABLE_NAME,
                projection, Database.DatabaseSchema.GLOBAL_COLUMN_NAME_SENT + "='FALSE'", null, null, null, null);
        String result = cursorToXML(cursor, Database.DatabaseSchema.Battery.TAG, Database.DatabaseSchema.Battery.COLUMN_NAME_LEVEL);
        xmlString.append(result);
    }

    private void getCall() {
        //query to join the call table with the number replacements table, matching equality on the number field
        String query = String.format("SELECT %s,%s,%s,%s,R._id FROM %s AS C JOIN %s AS R USING (%s) WHERE %s='FALSE'",
                Database.DatabaseSchema.CallLog.COLUMN_NAME_DATE, Database.DatabaseSchema.CallLog.COLUMN_NAME_TYPE,
                Database.DatabaseSchema.CallLog.COLUMN_NAME_DURATION, Database.DatabaseSchema.CallLog.COLUMN_NAME_SAVED,
                Database.DatabaseSchema.CallLog.TABLE_NAME, Database.DatabaseSchema.CallLogNumberReplacements.TABLE_NAME,
                Database.DatabaseSchema.CallLog.COLUMN_NAME_NUMBER, Database.DatabaseSchema.GLOBAL_COLUMN_NAME_SENT);
        Cursor cursor = db.rawQuery(query, null);
        StringBuilder result = new StringBuilder();
        int dateIndex = cursor.getColumnIndex(Database.DatabaseSchema.CallLog.COLUMN_NAME_DATE);
        int typeIndex = cursor.getColumnIndex(Database.DatabaseSchema.CallLog.COLUMN_NAME_TYPE);
        int durationIndex = cursor.getColumnIndex(Database.DatabaseSchema.CallLog.COLUMN_NAME_DURATION);
        int savedIndex = cursor.getColumnIndex(Database.DatabaseSchema.CallLog.COLUMN_NAME_SAVED);
        int iIndex = cursor.getColumnIndex(Database.DatabaseSchema.CallLogNumberReplacements._ID);

        while (cursor.moveToNext()) {
            result.append("<call time=\"" + TimeCapture.getGivenStringTime(cursor.getLong(dateIndex)) + "\" type=\"" + cursor.getString(typeIndex) +
                    "\" duration=\"" + cursor.getString(durationIndex) + "\" saved=\"" + cursor.getString(savedIndex) + "\">" +
                    cursor.getString(iIndex) + "</call>\n");
        }
        cursor.close();
        xmlString.append(result);
    }

    private void getCall2() {
        String query = String.format("SELECT %s,%s,%s,%s,R._id FROM %s AS C JOIN %s AS R USING (%s)",
                Database.DatabaseSchema.CallLog.COLUMN_NAME_DATE, Database.DatabaseSchema.CallLog.COLUMN_NAME_TYPE,
                Database.DatabaseSchema.CallLog.COLUMN_NAME_DURATION, Database.DatabaseSchema.CallLog.COLUMN_NAME_SAVED,
                Database.DatabaseSchema.CallLog.TABLE_NAME, Database.DatabaseSchema.CallLogNumberReplacements.TABLE_NAME,
                Database.DatabaseSchema.CallLog.COLUMN_NAME_NUMBER);
        Cursor cursor = db.rawQuery(query, null);
        String result = cursorToXML(cursor, Database.DatabaseSchema.CallLog.TAG, Database.DatabaseSchema.CallLogNumberReplacements._ID);
        xmlString.append(result);
    }

    private void getNumberReplacments() {
        Cursor cursor = db.query(Database.DatabaseSchema.CallLogNumberReplacements.TABLE_NAME, null, null, null, null, null, null);
        String result = cursorToXML(cursor, "rep",null);
        xmlString.append(result);
    }

    private void getWifiAPs() {
        String[] projection;
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.JELLY_BEAN_MR1) {
            projection = new String[]{Database.DatabaseSchema.WifiAP.COLUMN_NAME_BSSID,
                    Database.DatabaseSchema.WifiAP.COLUMN_NAME_SSID,
                    Database.DatabaseSchema.WifiAP.COLUMN_NAME_CAPABILITIES,
                    Database.DatabaseSchema.WifiAP.COLUMN_NAME_SIGNAL,
                    Database.DatabaseSchema.WifiAP.COLUMN_NAME_FREQ,
                    Database.DatabaseSchema.WifiAP.COLUMN_NAME_SEEN};
        } else {
            projection = new String[]{Database.DatabaseSchema.WifiAP.COLUMN_NAME_BSSID,
                    Database.DatabaseSchema.WifiAP.COLUMN_NAME_SSID,
                    Database.DatabaseSchema.WifiAP.COLUMN_NAME_CAPABILITIES,
                    Database.DatabaseSchema.WifiAP.COLUMN_NAME_SIGNAL,
                    Database.DatabaseSchema.WifiAP.COLUMN_NAME_FREQ};
        }
        Cursor cursor = db.query(Database.DatabaseSchema.WifiAP.TABLE_NAME, projection, null, null, null, null, null);
        String result = cursorToXML(cursor, Database.DatabaseSchema.WifiAP.TAG, Database.DatabaseSchema.WifiAP.COLUMN_NAME_SSID);
        xmlString.append(result);
    }

    private void getConnectivity() {
        //query to sum the totals with the current columns first
        String query = String.format("SELECT %s,%s,%s,(%s+%s) AS rx,(%s+%s) AS tx,%s FROM %s WHERE %s='FALSE'",
                Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_ACTIVE,
                Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_TIME,
                Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_SINCE,
                Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_TOTAL_RX,
                Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_CURRENT_RX,
                Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_TOTAL_TX,
                Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_CURRENT_TX,
                Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_TYPE,
                Database.DatabaseSchema.NetworkInterface.TABLE_NAME,
                Database.DatabaseSchema.GLOBAL_COLUMN_NAME_SENT);
        Cursor cursor = db.rawQuery(query, null);
        String result = cursorToXML(cursor, Database.DatabaseSchema.NetworkInterface.TAG, Database.DatabaseSchema.NetworkInterface.COLUMN_NAME_TYPE);
        xmlString.append(result);
    }

    private void getRunningServices() {
        String[] projection = new String[]{Database.DatabaseSchema.RunningServices.COLUMN_NAME_SINCE,
                Database.DatabaseSchema.RunningServices.COLUMN_NAME_TIME,
                Database.DatabaseSchema.RunningServices.COLUMN_NAME_RX,
                Database.DatabaseSchema.RunningServices.COLUMN_NAME_TX,
                Database.DatabaseSchema.RunningServices.COLUMN_NAME_PACKAGE_NAME};
        Cursor cursor = db.query(Database.DatabaseSchema.RunningServices.TABLE_NAME, projection,
                Database.DatabaseSchema.GLOBAL_COLUMN_NAME_SENT + "='FALSE'", null, null, null, null);
        String result = cursorToXML(cursor, Database.DatabaseSchema.RunningServices.TAG,
                Database.DatabaseSchema.RunningServices.COLUMN_NAME_PACKAGE_NAME);
        xmlString.append(result);
    }

    private void getInstalledPackages2() {
        String[] projection = new String[]{Database.DatabaseSchema.InstalledPackages.COLUMN_NAME_PACKAGE_NAME,
                Database.DatabaseSchema.InstalledPackages.COLUMN_NAME_INSTALLED_DATE,
                Database.DatabaseSchema.InstalledPackages.COLUMN_NAME_VERSION,
                Database.DatabaseSchema.InstalledPackages.COLUMN_NAME_LABEL};
        Cursor cursor = db.query(Database.DatabaseSchema.InstalledPackages.TABLE_NAME, projection,
                Database.DatabaseSchema.GLOBAL_COLUMN_NAME_SENT + "='FALSE'", null, null, null,
                Database.DatabaseSchema.InstalledPackages.COLUMN_NAME_INSTALLED_DATE + " DESC");
        String result = cursorToXML(cursor, Database.DatabaseSchema.InstalledPackages.TAG, Database.DatabaseSchema.InstalledPackages.COLUMN_NAME_LABEL);
        xmlString.append(result);
    }

    private void getInstalledPackages() {
        StringBuilder result = new StringBuilder();
        Cursor cursor = db.query(Database.DatabaseSchema.InstalledPackages.TABLE_NAME, null, null, null, null, null,
                Database.DatabaseSchema.InstalledPackages.COLUMN_NAME_INSTALLED_DATE + " DESC");
        int uid = cursor.getColumnIndex(Database.DatabaseSchema.InstalledPackages.COLUMN_NAME_UID);
        int label = cursor.getColumnIndex(Database.DatabaseSchema.InstalledPackages.COLUMN_NAME_LABEL);
        int version = cursor.getColumnIndex(Database.DatabaseSchema.InstalledPackages.COLUMN_NAME_VERSION);
        int date = cursor.getColumnIndex(Database.DatabaseSchema.InstalledPackages.COLUMN_NAME_INSTALLED_DATE);
        int name = cursor.getColumnIndex(Database.DatabaseSchema.InstalledPackages.COLUMN_NAME_PACKAGE_NAME);

        while (cursor.moveToNext()) {
            result.append("<installedapp name=\"" + cursor.getString(name) + "\" installed=\"" +
                    TimeCapture.getGivenStringTime(cursor.getLong(date)) + "\" version=\"" + cursor.getString(version) + "\" uid=\"" +
                    cursor.getString(uid) + "\">" + cursor.getString(label) + "</installedapp>\n");
        }
        cursor.close();
        xmlString.append(result);
    }

    private void getSMS() {
        String query = String.format("SELECT %s,%s,%s,R._id FROM %s AS C JOIN %s AS R USING (%s) WHERE %s='FALSE'",
                Database.DatabaseSchema.SMSLog.COLUMN_NAME_DATE, Database.DatabaseSchema.SMSLog.COLUMN_NAME_TYPE,
                Database.DatabaseSchema.SMSLog.COLUMN_NAME_READ,
                Database.DatabaseSchema.SMSLog.TABLE_NAME, Database.DatabaseSchema.CallLogNumberReplacements.TABLE_NAME,
                Database.DatabaseSchema.CallLogNumberReplacements.COLUMN_NAME_NUMBER,
                Database.DatabaseSchema.GLOBAL_COLUMN_NAME_SENT);
        Cursor cursor = db.rawQuery(query, null);
        String result = cursorToXML(cursor, Database.DatabaseSchema.SMSLog.TAG, Database.DatabaseSchema.CallLogNumberReplacements._ID);
        xmlString.append(result);
    }

    private void getSockets() {
        Cursor connections = db.query(Database.DatabaseSchema.Sockets.TABLE_NAME, null, null, null, null, null, null);
        StringBuilder result = new StringBuilder();
        int dateIndex = connections.getColumnIndex(Database.DatabaseSchema.Sockets.COLUMN_NAME_DATE);
        int lipIndex = connections.getColumnIndex(Database.DatabaseSchema.Sockets.COLUMN_NAME_LIP);
        int lportIndex = connections.getColumnIndex(Database.DatabaseSchema.Sockets.COLUMN_NAME_LPORT);
        int ripIndex = connections.getColumnIndex(Database.DatabaseSchema.Sockets.COLUMN_NAME_RIP);
        int rportIndex = connections.getColumnIndex(Database.DatabaseSchema.Sockets.COLUMN_NAME_RPORT);
        int typeIndex = connections.getColumnIndex(Database.DatabaseSchema.Sockets.COLUMN_NAME_TYPE);
        int statusIndex = connections.getColumnIndex(Database.DatabaseSchema.Sockets.COLUMN_NAME_STATUS);
        int uidIndex = connections.getColumnIndex(Database.DatabaseSchema.Sockets.COLUMN_NAME_UID);

        while (connections.moveToNext()) {
            Cursor listeners=db.query(Database.DatabaseSchema.InstalledPackages.TABLE_NAME,
                    new String[]{Database.DatabaseSchema.InstalledPackages.COLUMN_NAME_LABEL},
                    Database.DatabaseSchema.InstalledPackages.COLUMN_NAME_UID+"="+connections.getString(uidIndex), null, null, null, null);
            StringBuilder proc=new StringBuilder();
            if (listeners.moveToFirst()){
                int nameIndex=listeners.getColumnIndex(Database.DatabaseSchema.InstalledPackages.COLUMN_NAME_LABEL);
                proc.append(listeners.getString(nameIndex));
                while (listeners.moveToNext()) proc.append('/'+listeners.getString(nameIndex));
            }else proc.append("unknown");
            result.append("<connection time=\"" + connections.getString(dateIndex) + "\" lip=\"" +
                    connections.getString(lipIndex) + "\" lport=\"" + connections.getString(lportIndex) + "\" rip=\"" +
                    connections.getString(ripIndex) + "\" rport=\"" +connections.getString(rportIndex) + "\" type=\"" +
                    connections.getString(typeIndex) + "\" status=\"" +connections.getString(statusIndex) +  "\">" +
                    proc.toString() + "</connection>\n");
            listeners.close();
        }
        connections.close();

        xmlString.append(result);
    }

    public String getMetaData(Integer deviceId){
        return new MetaDataCapture(deviceId).getMetaDataXML();
    }

    public String getXML(Integer deviceID) {
        if (db==null) return null;
        xmlString.append("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" + "<data>\n"
                + getMetaData(deviceID) + "<device-data>\n");
        //we call all methods even if some functionality is disabled, because it may have been disabled after collecting some new data
        getActions();
        getBattery();
        getCall2();
        getWifiAPs();
        getConnectivity();
        getRunningServices();
        getSMS();
        getSockets();
        getInstalledPackages2();
        xmlString.append("</device-data>\n</data>");
        return xmlString.toString();
    }

    private class MetaDataCapture {
        //device id can be null(there will be no <device> tag in this case)
        private Integer deviceID;
        private ArrayList<String[]> data;

        public MetaDataCapture(Integer deviceID) {
            data=new ArrayList<String[]>();
            this.deviceID=deviceID;
            getMetaData();
        }

        private void getMetaData() {
            getTelephonyDetails();
            getPhoneDetails();
            getDefaultBrowser();
        }

        private void getTelephonyDetails() {
            String[] datatype = new String[]{"unknown", "gprs", "edge", "umts", "cdma", "evdo0", "evdoA", "1xrtt", "hsdpa", "hsupa", "hspa", "iden", "evdoB", "lte", "ehrpd", "hspap"};
            TelephonyManager tm = (TelephonyManager) context.getSystemService(Context.TELEPHONY_SERVICE);
            data.add(new String[]{"imei", tm.getDeviceId()});
            if (deviceID!=null) data.add(new String[]{"device",deviceID.toString()});
            data.add(new String[]{"dataNetType", datatype[tm.getNetworkType()]});
            data.add(new String[]{"country", tm.getNetworkCountryIso()});
            data.add(new String[]{"network", tm.getNetworkOperatorName()});
            data.add(new String[]{"carrier", tm.getSimOperatorName()});
        }

        private void getPhoneDetails() {
            data.add(new String[]{"manufacturer", Build.MANUFACTURER});
            data.add(new String[]{"model", Build.MODEL});
            data.add(new String[]{"androidVer", Build.VERSION.RELEASE});
            data.add(new String[]{"lastReboot", TimeCapture.getUpDate()});
            data.add(new String[]{"timeZone",  TimeZone.getDefault().getDisplayName(false,TimeZone.SHORT)});
        }

        private void getDefaultBrowser(){
            //resolve the default application used to open http
            Intent browseIntent =new Intent("android.intent.action.VIEW", Uri.parse("http://"));
            ResolveInfo defaultBrowse=context.getPackageManager().resolveActivity(browseIntent, PackageManager.MATCH_DEFAULT_ONLY);
            data.add(new String[]{"defaultBrowser",defaultBrowse.activityInfo.packageName});
        }

        public String getMetaDataXML() {
            StringBuilder result = new StringBuilder("<metadata>\n");
            for (String[] d : data) {
                result.append("<" + d[0] + ">" + d[1] + "</" + d[0] + ">\n");
            }
            result.append("</metadata>\n");
            return result.toString();
        }
    }
}