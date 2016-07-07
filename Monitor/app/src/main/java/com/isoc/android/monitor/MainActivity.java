package com.isoc.android.monitor;

import android.content.ComponentName;
import android.content.Context;
import android.content.Intent;
import android.content.ServiceConnection;
import android.net.ConnectivityManager;
import android.net.NetworkInfo;
import android.os.AsyncTask;
import android.os.Bundle;
import android.os.IBinder;
import android.support.v7.app.AppCompatActivity;
import android.support.v7.widget.Toolbar;
import android.text.method.ScrollingMovementMethod;
import android.view.Menu;
import android.view.MenuItem;
import android.view.View;
import android.widget.EditText;
import android.widget.TextView;
import android.widget.Toast;

import java.io.BufferedOutputStream;
import java.io.IOException;
import java.io.OutputStream;
import java.net.HttpURLConnection;
import java.net.MalformedURLException;
import java.net.ProtocolException;
import java.net.URL;

public class MainActivity extends AppCompatActivity {
    MyService mService;
    boolean mBound = false;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);
        Toolbar toolbar = (Toolbar) findViewById(R.id.toolbar);
        setSupportActionBar(toolbar);
    }

    @Override
    public boolean onCreateOptionsMenu(Menu menu) {
        // Inflate the menu; this adds items to the action bar if it is present.
        getMenuInflater().inflate(R.menu.menu_main, menu);
        return true;
    }

    @Override
    public boolean onOptionsItemSelected(MenuItem item) {
        // Handle action bar item clicks here. The action bar will
        // automatically handle clicks on the Home/Up button, so long
        // as you specify a parent activity in AndroidManifest.xml.
        int id = item.getItemId();

        //noinspection SimplifiableIfStatement
        if (id == R.id.action_settings) {
            return true;
        }

        return super.onOptionsItemSelected(item);
    }

    @Override
    protected void onStop() {
        super.onStop();
        if (mBound) {
            unbindService(mConnection);
            mBound = false;
        }
    }

    @Override
    protected void onStart(){
        super.onStart();
        Intent intent = new Intent(this, MyService.class);
        bindService(intent, mConnection, Context.BIND_AUTO_CREATE);
    }

    private ServiceConnection mConnection = new ServiceConnection() {
        @Override
        public void onServiceConnected(ComponentName componentName, IBinder service) {
            mBound = true;
            MyService.LocalBinder binder = (MyService.LocalBinder) service;
            mService = binder.getService();

        }

        @Override
        public void onServiceDisconnected(ComponentName componentName) {
            mBound = false;
        }
    };

    public void startService(View view) {
        startService(new Intent(getBaseContext(), MyService.class));
    }

    public void stopService(View view) {
        stopService(new Intent(getBaseContext(), MyService.class));
    }

    public String getResults() {
        String result=new String();
        if (mService !=null)
            result = mService.generateXML();
        else
            Toast.makeText(this,"Service not binded",Toast.LENGTH_LONG).show();
        return result;
    }

    public void showResults(View view) {
        TextView text = (TextView) findViewById(R.id.textResults);
        text.setMovementMethod(new ScrollingMovementMethod());
        text.setText(getResults());
    }

    public boolean checkNet(){
        ConnectivityManager connection = (ConnectivityManager) getSystemService(Context.CONNECTIVITY_SERVICE);
        NetworkInfo netInfo = connection.getActiveNetworkInfo();
        if (!(netInfo != null && netInfo.isConnected())) {
            Toast.makeText(this, "No active Connection", Toast.LENGTH_LONG).show();
            return false;
        }
        return true;
    }

    public void sendXML(View view){
        if (!checkNet()) return;
        EditText tIP = (EditText) findViewById(R.id.textIP);
        String ip = tIP.getText().toString();
        new Post().execute(ip,getResults());

    }


    private class Post extends AsyncTask<String,Void,String> {

        @Override
        protected String doInBackground(String... args) {
            String result=new String();
            URL url;
            HttpURLConnection client = null;
            String xml= args[1];
            try {
                url = new URL(args[0]);
                client = (HttpURLConnection) url.openConnection();
                client.setConnectTimeout(4000);
                client.setFixedLengthStreamingMode(xml.getBytes().length);
                client.setRequestMethod("POST");
                client.setDoOutput(true);
                OutputStream out = new BufferedOutputStream(client.getOutputStream());
                out.write(xml.getBytes());
                out.flush();
                out.close();
                result="done!";
            }catch (java.net.SocketTimeoutException e){
                result="TimeOut";
            }catch (MalformedURLException e) {
                result="Malformed URL";
            } catch (ProtocolException e) {
                result="Protocol Exception: "+e.getMessage();
            } catch (IOException e) {
                result="IOException: "+e.getMessage();
            }finally{
                if (client != null)
                    client.disconnect();
            }
            return result;
        }

        @Override
        protected void onPostExecute(String s) {
            Toast.makeText(getBaseContext(),s,Toast.LENGTH_LONG).show();
        }
    }
}
