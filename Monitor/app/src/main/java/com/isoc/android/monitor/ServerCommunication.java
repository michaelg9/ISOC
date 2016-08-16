package com.isoc.android.monitor;

import android.accounts.AccountManager;
import android.accounts.AuthenticatorException;
import android.content.Context;
import android.preference.PreferenceManager;
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
import java.util.regex.Matcher;
import java.util.regex.Pattern;

/**
 *
 */
public class ServerCommunication {
    private String serverURL;
    private Context context;

    public ServerCommunication(Context context) {
        this.context=context;
        serverURL= PreferenceManager.getDefaultSharedPreferences(context).getString(context.getString(R.string.server_url), null);
    }

    public String[] login(String email,String password){
        String loginURL = serverURL +"/auth/0.1/login";
        HttpURLConnection client = null;
        String result[]=new String[2];
        result[0]=AccountManager.KEY_ERROR_MESSAGE;
        try {
            URL url = new URL(loginURL);
            client = (HttpURLConnection) url.openConnection();
            client.setConnectTimeout(4000);
            client.setRequestMethod("POST");
            client.setDoOutput(true);
            String payload="email="+email+"&"+"password="+password;
            OutputStream out = new BufferedOutputStream(client.getOutputStream());
            out.write(payload.getBytes());
            out.flush();
            out.close();
            if(client.getResponseCode()==200) {
                BufferedReader r = new BufferedReader(new InputStreamReader(client.getInputStream()));
                String answer=r.readLine();
                Log.e("answer",answer);
                result[1]=extractToken(context.getString(R.string.token_refresh),answer);
                result[0]= AccountManager.KEY_AUTHTOKEN;
                r.close();
            }else{
                BufferedReader r = new BufferedReader(new InputStreamReader(client.getErrorStream()));
                String error=client.getResponseMessage()+'('+client.getResponseCode()+')'+": "+r.readLine();
                r.close();
                throw new AuthenticatorException(error);
            }
        }catch (AuthenticatorException e){
            result[1]=e.getMessage();
        }catch (java.net.SocketTimeoutException e) {
            result[1]="Error: TimeOut";
        } catch (MalformedURLException e) {
            result[1]="Error: Malformed URL";
        } catch (ProtocolException e) {
            result[1]="Error: Protocol Exception " + e.getMessage();
        } catch (IOException e) {
            result[1]="Error: IOException " + e.getMessage();
        } finally {
            if (client != null)
                client.disconnect();
        }
        Log.e("token",result[1]);
        return result;
    }

    //extracts the target token from the server's answer
    private String extractToken(String target,String keyAuthtoken) {
        //match everything in quotes,excluding the quotes
        Matcher tokenMatcher=Pattern.compile("\"([^\"]*)\"").matcher(keyAuthtoken);
        String result=new String();
        while (tokenMatcher.find()){
            if (!tokenMatcher.group(1).equals(target)) continue;
            if (tokenMatcher.find()) result=tokenMatcher.group(1);
            else return null;
        }
        return result;
    }
}
