package com.isoc.android.monitor;

import android.app.Fragment;
import android.content.Context;
import android.database.sqlite.SQLiteDatabase;
import android.net.ConnectivityManager;
import android.net.NetworkInfo;
import android.os.AsyncTask;
import android.os.Bundle;
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
 * TO DO:
 * Formatting call log numbers: a number without country code is taken as a different number than the same number with country code
 * WIFI / BLUETOOTH SCAN
 * SMS
 * CONNECTIONS
 * LOCATION
 * TIMEZONE
 * ACTIONS PREFRENCE DETAILED NEW SCREENPREFERENCE
 * MOBILE INTF READ DIRECTLY
 * BROWSER HISTORY
 * CELL TOWER CHANGE
 * SYSTEM APPS REPORTING OLD INSTALLED DATE
 * ----------
 * NETWORK:
 *  * when mobile off, counters =0. OK we get that directly from the source too...
 * http://randomizedsort.blogspot.co.uk/2010/10/where-does-android-gets-its-traffic.html
 *
 * default since uptime, what if no capture and restart?
 *
 * Deprecated onreceive method, implement type?
 * ------
 */
public class MainFragment extends Fragment {

    public MainFragment() {
        // Required empty public constructor
    }


    @Override
    public View onCreateView(LayoutInflater inflater, ViewGroup container,
                             Bundle savedInstanceState) {
        // Inflate the layout for this fragment
        View view= inflater.inflate(R.layout.fragment_main, container, false);


        final Button deleteDB = (Button) view.findViewById(R.id.delete_db);
        deleteDB.setOnClickListener(new View.OnClickListener(){
            @Override
            public void onClick(View view) {
                getActivity().deleteDatabase(Database.DatabaseSchema.dbName);
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


        final Button checkAlarm = (Button) view.findViewById(R.id.checkAlarm);
        checkAlarm.setOnClickListener(new View.OnClickListener(){

            @Override
            public void onClick(View view) {
                String s=(MyService.ServiceControls.checkExistence(getActivity())) ? "Exists!!" : "Doesn't exist :(";
                Toast.makeText(getActivity(),s,Toast.LENGTH_LONG).show();
            }
        });

        return view;
    }

    @Override
    public void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        PreferenceManager.setDefaultValues(getActivity(), R.xml.preferences, false);
    }

    public String getResults(Context context){
        SQLiteDatabase db = new Database(context).getReadableDatabase();
        String result= "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" + "<data>\n"+"<metadata>\n" +
                MetaDataCapture.getMetaDataXML(context)+
                "</metadata>\n" +"<device-data>\n"+
                SocketsCapture.getSocketsXML(db)+
                BatteryCapture.getBatteryXML(db)+
                ActionCapture.getActionsXML(db)+
                NetworkCapture.getTrafficXML(db)+
                NetworkCapture.getWifiAPResultsXML(db)+
                ContactsCapture.getCallXML2(db)+
                PackageCapture.getRunningServicesXML(db)+
                PackageCapture.getInstalledPackagesXML2(db)+
                "</device-data>\n</data>";
        db.close();

        return result;

    }

    public void showResults() {

        String results=getResults(getActivity());
        Bundle bundle = new Bundle();
        bundle.putString("results",results);
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
        new Post().execute(ip,getResults(getActivity()));

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
                result="Done! Response code: "+client.getResponseCode();
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