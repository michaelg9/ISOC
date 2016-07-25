package com.isoc.android.monitor;

import android.content.ContentValues;
import android.content.Context;
import android.database.Cursor;
import android.database.sqlite.SQLiteDatabase;

import java.util.regex.Matcher;
import java.util.regex.Pattern;

/**
 * Created by me on 19/07/16.
 */
public class SocketsCapture {
    private static String socket4 = "\\p{XDigit}{8}:\\p{XDigit}{4}";
    private static String socket6 = "\\p{XDigit}{32}:\\p{XDigit}{4}";
    //                                    sl :      local          remote       status        tx_queue:rx_queue           tr:tm->when                 retrnsmt           uid
    private static String Connection4 = "\\d+: " + socket4 + " " + socket4 + " \\p{XDigit}{2} \\p{XDigit}*:\\p{XDigit}* \\p{XDigit}{2}:\\p{XDigit}{8} \\p{XDigit}{8}\\s+\\d+";
    private static String Connection6 = "\\d+: " + socket6 + " " + socket6 + " \\p{XDigit}{2} \\p{XDigit}*:\\p{XDigit}* \\p{XDigit}{2}:\\p{XDigit}{8} \\p{XDigit}{8}\\s+\\d+";

    public static void getSockets(Context context,SQLiteDatabase db) {
        String time = TimeCapture.getTime();
        String[] types={"tcp","udp","raw","tcp6","udp6","raw6"};
        for (String type :types) {
            String content = NetworkCapture.readStatsFromFile("/proc/net/"+type);
            Matcher m;
            boolean is6 = type.contains("6");
            if (is6) m=Pattern.compile(Connection6).matcher(content);
                else m=Pattern.compile(Connection4).matcher(content);
            while (m.find()) {
                Socket socket = (is6) ? new Socket6(m.group()) : new Socket4(m.group());
                ContentValues values = new ContentValues();
                values.put(Database.DatabaseSchema.Sockets.COLUMN_NAME_DATE, time);
                values.put(Database.DatabaseSchema.Sockets.COLUMN_NAME_LIP, socket.getLocalIP());
                values.put(Database.DatabaseSchema.Sockets.COLUMN_NAME_LPORT, socket.getLocalPort());
                values.put(Database.DatabaseSchema.Sockets.COLUMN_NAME_RIP, socket.getRemoteIP());
                values.put(Database.DatabaseSchema.Sockets.COLUMN_NAME_RPORT, socket.getRemotePort());
                values.put(Database.DatabaseSchema.Sockets.COLUMN_NAME_STATUS, socket.getStatus());
                values.put(Database.DatabaseSchema.Sockets.COLUMN_NAME_UID, socket.getUid());
                values.put(Database.DatabaseSchema.Sockets.COLUMN_NAME_TYPE, type);
                db.insert(Database.DatabaseSchema.Sockets.TABLE_NAME,null,values);
            }
        }
    }

    public static String getSocketsXML(SQLiteDatabase db){
        Cursor cursor = db.query(Database.DatabaseSchema.Sockets.TABLE_NAME,null,null,null,null,null,null);
        StringBuilder sb=new StringBuilder();
        int dateIndex = cursor.getColumnIndex(Database.DatabaseSchema.Sockets.COLUMN_NAME_DATE);
        int lipIndex = cursor.getColumnIndex(Database.DatabaseSchema.Sockets.COLUMN_NAME_LIP);
        int lportIndex = cursor.getColumnIndex(Database.DatabaseSchema.Sockets.COLUMN_NAME_LPORT);
        int ripIndex = cursor.getColumnIndex(Database.DatabaseSchema.Sockets.COLUMN_NAME_RIP);
        int rportIndex = cursor.getColumnIndex(Database.DatabaseSchema.Sockets.COLUMN_NAME_RPORT);
        int typeIndex = cursor.getColumnIndex(Database.DatabaseSchema.Sockets.COLUMN_NAME_TYPE);
        int statusIndex = cursor.getColumnIndex(Database.DatabaseSchema.Sockets.COLUMN_NAME_STATUS);
        int uidIndex = cursor.getColumnIndex(Database.DatabaseSchema.Sockets.COLUMN_NAME_UID);

        while (cursor.moveToNext()) {
            sb.append("<connection time=\"" + cursor.getString(dateIndex) + "\" local=\"" + cursor.getString(lipIndex)+':'+
                    cursor.getString(lportIndex)+"\" remote=\""+cursor.getString(ripIndex)+':'+
                            cursor.getString(rportIndex)+"\" type=\""+cursor.getString(typeIndex)+"\" status=\""+
                    cursor.getString(statusIndex)+"\" uid=\""+cursor.getString(uidIndex)+"\">" +
                    "proc" + "</connection>\n");
        }
        cursor.close();
        return sb.toString();
    }

    private abstract static class Socket {
        private String localIP;
        private String localPort;
        private String remoteIP;
        private String remotePort;
        private String uid;
        private String status;

        private enum states {
            //taken from linux source code
            UNKNOWN(0),
            ESTABLISHED(1),
            SYN_SENT(2),
            SYN_RECV(3),
            FIN_WAIT1(4),
            FIN_WAIT2(5),
            TIME_WAIT(6),
            CLOSE(7),
            CLOSE_WAIT(8),
            LAST_ACK(9),
            LISTEN(10),
            CLOSING(11),
            NEW_SYN_RECV(12),
            MAX_STATES(13);

            private final int n;

            states(int n) {
                this.n = n;
            }
        }

        private Socket(String entry) {
            String[] tokens = entry.trim().split(" ");
            setLocalSocket(tokens[1].split(":"));
            setRemoteSocket(tokens[2].split(":"));
            uid = tokens[tokens.length - 1];
            setStatus(tokens[3]);
        }

        public String getLocalIP() {
            return localIP;
        }

        public String getLocalPort() {
            return localPort;
        }

        public String getRemoteIP() {
            return remoteIP;
        }

        public String getRemotePort() {
            return remotePort;
        }

        public String getUid() {
            return uid;
        }

        public String getStatus() {
            return status;
        }

        private void setLocalSocket(String[] lSocket) {
            localIP = formatIP(convertFromEndian(lSocket[0]));
            localPort = Integer.toString(Integer.parseInt(lSocket[1], 16));
        }

        public abstract String formatIP(String[] localIP);

        private void setRemoteSocket(String[] rSocket) {
            remoteIP = formatIP(convertFromEndian(rSocket[0]));
            remotePort = Integer.toString(Integer.parseInt(rSocket[1], 16));
        }

        private String[] convertFromEndian(String s) {
            String[] result = new String[s.length() / 2];
            String[] tokens = splitEqually(s, 8);
            for (int i = 0; i < tokens.length; i++) {
                String[] parts = splitEqually(tokens[i], 2);
                int j = 0;
                for (int k = 3; k >= 0; k--) {
                    result[i * 4 + j] = parts[k];
                    j++;
                }
            }
            return result;
        }

        private String[] splitEqually(String s, int num) {
            int length = s.length();
            if ((length % num) != 0) return null;
            String[] result = new String[length / num];
            for (int k = 0; k <= length - num; k += num) {
                result[k / num] = s.substring(k, k + num);
            }
            return result;
        }

        private void setStatus(String status) {
            this.status = states.values()[Integer.parseInt(status, 16)].name().toLowerCase();
        }
    }

    private static class Socket4 extends Socket {

        public Socket4(String entry) {
            super(entry);
        }

        public String formatIP(String[] localIP) {
            StringBuilder result = new StringBuilder().append(Integer.decode("0x" + localIP[0]));
            for (int i = 1; i < localIP.length; i++) {
                result.append('.' + Integer.decode("0x" + localIP[i]).toString());
            }
            return result.toString();
        }
    }

    private static class Socket6 extends Socket {
        public Socket6(String entry) {
            super(entry);
        }

        public String formatIP(String[] localIP) {
            StringBuilder result = new StringBuilder(localIP[0] + localIP[1]);
            for (int i = 2; i < localIP.length; i += 2) {
                result.append(':' + localIP[i] + localIP[i + 1]);
            }
            return result.toString();
        }
    }
}
