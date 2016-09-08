package com.isoc.android.monitor;

import android.content.ContentValues;
import android.database.sqlite.SQLiteDatabase;

import java.util.regex.Matcher;
import java.util.regex.Pattern;

/**
 * Captures active sockets by reading files tcp tcp6 udp udp6 raw raw6 stored at /net/proc/
 * BUG: Shared uids produce uncertainty over which app initiated the connection... It's only vendor / system apps though that share uids...
 */
public class SocketsCapture {
    //pattern for an IPv4 and v6 address in hexadecimal
    private static String socket4 = "\\p{XDigit}{8}:\\p{XDigit}{4}";
    private static String socket6 = "\\p{XDigit}{32}:\\p{XDigit}{4}";
    //pattern to match syntax in these files
    //                                    sl :      local          remote       status        tx_queue:rx_queue           tr:tm->when                 retrnsmt           uid
    private static String Connection4 = "\\d+: " + socket4 + " " + socket4 + " \\p{XDigit}{2} \\p{XDigit}*:\\p{XDigit}* \\p{XDigit}{2}:\\p{XDigit}{8} \\p{XDigit}{8}\\s+\\d+";
    private static String Connection6 = "\\d+: " + socket6 + " " + socket6 + " \\p{XDigit}{2} \\p{XDigit}*:\\p{XDigit}* \\p{XDigit}{2}:\\p{XDigit}{8} \\p{XDigit}{8}\\s+\\d+";

    public static void getSockets(SQLiteDatabase db) {
        String time = TimeCapture.getCurrentStringTime();
        String[] types={"tcp","udp","raw","tcp6","udp6","raw6"}; //the names of the interesting files
        for (String type :types) {
            String content = NetworkCapture.readStatsFromFile("/proc/net/"+type);
            Matcher m;
            boolean is6 = type.contains("6");
            //if it's an ipv6 file, use socket6 pattern
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

    //class representing an IPv6/v4 socket. Must be extended to implement formatIP
    private abstract static class Socket {
        private String localIP;
        private String localPort;
        private String remoteIP;
        private String remotePort;
        private String uid;
        private String status;

        //all the possible socket states,taken from linux source code
        private enum states {
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
            localIP = formatIP(convertFromEndian(lSocket[0])); //must reverse the bytes first, then format appropriately
            localPort = Integer.toString(Integer.parseInt(lSocket[1], 16)); //convert from hexadecimal
        }

        public abstract String formatIP(String[] localIP);

        private void setRemoteSocket(String[] rSocket) {
            remoteIP = formatIP(convertFromEndian(rSocket[0]));
            remotePort = Integer.toString(Integer.parseInt(rSocket[1], 16));
        }


        private String[] convertFromEndian(String s) {
            //reverts the bytes into their appropriate position
            String[] result = new String[s.length() / 2];
            //take substrings of length 8
            for (int i = 0; i < s.length(); i+=8) {
                String part=s.substring(i,i+8);
                int j=0;
                for (int k = 6; k >= 0; k-=2) {
                    //append 2 character strings in reverse order
                    result[i/2+j] = part.substring(k,k+2);
                    j++;
                }
            }
            return result;
        }

        //convert status code into the appropriate enum
        private void setStatus(String status) {
            this.status = states.values()[Integer.parseInt(status, 16)].name().toLowerCase();
        }
    }

    private static class Socket4 extends Socket {

        public Socket4(String entry) {
            super(entry);
        }

        //must convert to decimal and append dots inbetween
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

        // just append a dots inbeetween. IPv6 addresses are in hexadecimal by nature
        public String formatIP(String[] localIP) {
            StringBuilder result = new StringBuilder(localIP[0] + localIP[1]);
            for (int i = 2; i < localIP.length; i += 2) {
                result.append('.' + localIP[i] + localIP[i + 1]);
            }
            return result.toString();
        }
    }
}
