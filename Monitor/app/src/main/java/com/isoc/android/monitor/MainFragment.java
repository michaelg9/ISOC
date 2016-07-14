package com.isoc.android.monitor;

import android.app.Fragment;
import android.content.ComponentName;
import android.content.Context;
import android.content.Intent;
import android.content.ServiceConnection;
import android.net.ConnectivityManager;
import android.net.NetworkInfo;
import android.os.AsyncTask;
import android.os.Bundle;
import android.os.IBinder;
import android.preference.PreferenceManager;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.Button;
import android.widget.Toast;

import java.io.BufferedOutputStream;
import java.io.IOException;
import java.io.OutputStream;
import java.net.HttpURLConnection;
import java.net.MalformedURLException;
import java.net.ProtocolException;
import java.net.URL;


/**
 * A simple {@link Fragment} subclass.
 */
public class MainFragment extends Fragment {
    private MyService mService;
    private boolean mBound = false;

    public MainFragment() {
        // Required empty public constructor
    }

    @Override
    public void onStart() {
        super.onStart();
        Intent intent = new Intent(getActivity(), MyService.class);
        getActivity().bindService(intent, mConnection, Context.BIND_AUTO_CREATE);

    }

    @Override
    public void onStop() {
        super.onStop();
        if (mBound) {
            getActivity().unbindService(mConnection);
            mBound = false;
        }
    }

    @Override
    public View onCreateView(LayoutInflater inflater, ViewGroup container,
                             Bundle savedInstanceState) {
        // Inflate the layout for this fragment
        View view= inflater.inflate(R.layout.fragment_main, container, false);


        final Button buttonStart = (Button) view.findViewById(R.id.buttonStart);
        buttonStart.setOnClickListener(new View.OnClickListener(){
            @Override
            public void onClick(View view) {
                startService();
            }
        });

        final Button buttonStop = (Button) view.findViewById(R.id.buttonStop);
        buttonStop.setOnClickListener(new View.OnClickListener(){

            @Override
            public void onClick(View view) {
                stopService();
            }
        });

        final Button showResults = (Button) view.findViewById(R.id.buttonShow);
        showResults.setOnClickListener(new View.OnClickListener(){

            @Override
            public void onClick(View view) {
                showResults();
            }
        });

        final Button sendResults = (Button) view.findViewById(R.id.buttonSend);
        sendResults.setOnClickListener(new View.OnClickListener(){

            @Override
            public void onClick(View view) {
                sendXML();
            }
        });

        return view;
    }

    @Override
    public void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
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

    public void startService() {
        getActivity().startService(new Intent(getActivity(), MyService.class));
    }

    public void stopService() {
        getActivity().stopService(new Intent(getActivity(), MyService.class));
    }

    public String getResults() {
        String result=new String();
        if (mService !=null)
            result = mService.generateXML();
        else
            Toast.makeText(getActivity(),"Service not binded",Toast.LENGTH_LONG).show();
        return result;
    }

    public void showResults() {
        Bundle bundle = new Bundle();
        bundle.putString("results",getResults());
        ShowFragment showFragment=new ShowFragment();
        showFragment.setArguments(bundle);
        getFragmentManager().beginTransaction().replace(R.id.fragment_container,showFragment).addToBackStack(null).commit();
    }

    public boolean checkNet(){
        ConnectivityManager connection = (ConnectivityManager) getActivity().getSystemService(Context.CONNECTIVITY_SERVICE);
        NetworkInfo netInfo = connection.getActiveNetworkInfo();
        if (!(netInfo != null && netInfo.isConnected())) {
            Toast.makeText(getActivity(), "No active Connection", Toast.LENGTH_LONG).show();
            return false;
        }
        return true;
    }

    public void sendXML(){
        if (!checkNet()) return;
        String ip =PreferenceManager.getDefaultSharedPreferences(getActivity()).getString("server_url",null);
        new Post().execute(ip,getResults());

    }


    private class Post extends AsyncTask<String,Void,String> {

        @Override
        protected String doInBackground(String... args) {
            String result;
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
                result="Done!";
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
            Toast.makeText(getActivity(),s,Toast.LENGTH_LONG).show();
        }
    }
}