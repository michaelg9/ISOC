package com.isoc.android.monitor;

import android.accounts.AccountManager;
import android.accounts.AuthenticatorException;
import android.app.Notification;
import android.app.NotificationManager;
import android.content.Context;
import android.preference.PreferenceManager;
import android.support.v4.app.NotificationCompat;
import android.util.Log;

import java.io.BufferedOutputStream;
import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStreamReader;
import java.io.OutputStream;
import java.net.HttpURLConnection;
import java.net.MalformedURLException;
import java.net.ProtocolException;
import java.net.URL;
import java.util.Arrays;
import java.util.HashMap;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

/**
 * All the communication with the server methods are here
 * Note that these are all networking methods so they should be run in a separate thread
 */
public class ServerCommunication {
    private String serverURL;
    private Context context;
    private static String loginURL = "/auth/0.1/login";
    private static String sendURL = "/app/0.1/upload";
    private static String registerURL="/signup";
    private static String accessTokenRequestURL= "/auth/0.1/token";
    private static String logOutURL="/auth/0.1/logout";

    public ServerCommunication(Context context) {
        this.context = context;
        serverURL = PreferenceManager.getDefaultSharedPreferences(context).getString(context.getString(R.string.server_url), null);
    }

    public String refreshRefreshToken(String oldRefreshToken){
        final String refreshURL="/auth/0.1/refresh";
        String[] s=requestPOST(refreshURL,"refreshToken="+oldRefreshToken,null);

        if (s[0].equals("200")) {
            return extractToken(context.getString(R.string.token_refresh), s[2]);
        }
        return null;
    }

    public String[] logOut(String refreshToken){
        return requestPOST(logOutURL,"refreshToken="+refreshToken,null);
    }

    public String getAccessToken(String refreshToken) {
        String[] s = requestPOST(accessTokenRequestURL, "refreshToken=" + refreshToken, null);
        String result;
        if (s[0].equals("200")) {
            result = extractToken(context.getString(R.string.token_access), s[2]);
        } else {
            result = null;
        }
        Log.e("access:","here: "+result);
        return result;
    }

    public String[] login(String email, String password) {
        String[] s = requestPOST(loginURL, "email=" + email + "&" + "password=" + password, null);
        String result[] = new String[2];
        if (s[0].equals("200")) {
            result[0] = AccountManager.KEY_AUTHTOKEN;
            result[1] = extractToken(context.getString(R.string.token_refresh), s[2]);
        } else {
            result[0] = AccountManager.KEY_ERROR_MESSAGE;
            result[1] = s[0];
        }
        Log.e("token", result[1]);
        return result;
    }

    public String[] sendData(String accessToken, String body) {
        HashMap<String, String> headers = new HashMap<>(1);
        headers.put("Authorization", "Bearer " + accessToken);
        headers.put("Content-Type", "application/xml");
        String[] result= requestPOST(sendURL, body, headers);
        showNotification(result);
        return result;
    }

    public String[] register(String email, String password) {
        HashMap<String,String> header=new HashMap<>();
        header.put("Content-Type", "application/xml");
        String[] s = requestPOST(registerURL+"?email=" + email + "&" + "password=" + password,
                "<xml>\n"+new XMLProduce(context,null).getMetaData(null)+"\n</xml>",header);
        String result[] = new String[2];
        if (s[0].equals("200")) {
            result[0] = LoginActivity.REGISTER_SUCCESS;
            result[1] = s[2];
        } else {
            result[0] = AccountManager.KEY_ERROR_MESSAGE;
            result[1] = s[0];
        }
        return result;
    }

    private String[] requestPOST(String path, String body, HashMap<String, String> headers) {
        String URL = serverURL + path;
        HttpURLConnection client = null;
        String[] result = new String[3];
        try {
            URL url = new URL(URL);
            client = (HttpURLConnection) url.openConnection();
            client.setConnectTimeout(4000);
            client.setRequestMethod("POST");
            if (headers != null) {
                for (String s : headers.keySet()) {
                    client.addRequestProperty(s, headers.get(s));
                }
            }
            if (body!=null) {
                client.setDoOutput(true);
                client.setFixedLengthStreamingMode(body.getBytes().length);
                OutputStream out = new BufferedOutputStream(client.getOutputStream());
                out.write(body.getBytes());
                out.flush();
                out.close();
            }
            if (client.getResponseCode() == 200) {
                BufferedReader r = new BufferedReader(new InputStreamReader(client.getInputStream()));
                result[0] = "200";
                result[1] = client.getResponseMessage();
                result[2] = r.readLine();
                Log.e("answer", result[2]);
                r.close();
            } else {
                BufferedReader r = new BufferedReader(new InputStreamReader(client.getErrorStream()));
                String error = client.getResponseMessage() + '(' + client.getResponseCode() + ')' + ": " + r.readLine();
                r.close();
                throw new AuthenticatorException(error);
            }
        } catch (AuthenticatorException e) {
            result[0] = e.getMessage();
        } catch (java.net.SocketTimeoutException e) {
            result[0] = "Error: TimeOut";
        } catch (MalformedURLException e) {
            result[0] = "Error: Malformed URL";
        } catch (ProtocolException e) {
            result[0] = "Error: Protocol Exception " + e.getMessage();
        } catch (IOException e) {
            result[0] = "Error: IOException " + e.getMessage();
        } finally {
            if (client != null)
                client.disconnect();
        }
        Log.e("answer:", Arrays.toString(result));
        return result;
    }

    private void showNotification(String[] result) {
        String content=(result[0].equals("200")) ? "Successfuly sent!" : result[0];
        Notification n = new NotificationCompat.Builder(context).setContentTitle("Monitor Data Sent").setSmallIcon(R.mipmap.ic_launcher)
                .setContentText(content).build();
        NotificationManager notificationManager = (NotificationManager) context.getSystemService(Context.NOTIFICATION_SERVICE);
        notificationManager.notify(0, n);
    }

    //extracts the target token from the server's answer
    private String extractToken(String target, String keyAuthtoken) {
        //match everything in quotes,excluding the quotes
        Matcher tokenMatcher = Pattern.compile("\"([^\"]*)\"").matcher(keyAuthtoken);
        String result = new String();
        while (tokenMatcher.find()) {
            if (!tokenMatcher.group(1).equals(target)) continue;
            if (tokenMatcher.find()) result = tokenMatcher.group(1);
            else return null;
        }
        return result;
    }
}