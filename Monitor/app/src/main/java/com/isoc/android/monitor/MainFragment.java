package com.isoc.android.monitor;

import android.app.Fragment;
import android.content.Context;
import android.os.Bundle;
import android.preference.PreferenceManager;
import android.support.v7.app.AppCompatActivity;
import android.support.v7.widget.Toolbar;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.Button;


/**
 * TO DO:
 * LOCATION
 * TIMEZONE
 * BROWSER HISTORY?
 * CELL TOWER CHANGE
 * String resources
 * capture accounts
 * ----------
 * BUGS:
 * Deprecated connectivity onReceive method, implement type?
 * Matching sockets to specific app.
 * SYSTEM APPS REPORTING OLD INSTALLED DATE
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
        View view = inflater.inflate(R.layout.fragment_main, container, false);


        final Button deleteDB = (Button) view.findViewById(R.id.delete_db);
        deleteDB.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View view) {
                getActivity().deleteDatabase(Database.DatabaseSchema.dbName);
                getActivity().getSharedPreferences(getActivity().getString(R.string.shared_values_filename), Context.MODE_PRIVATE).edit().clear().apply();
            }
        });

        final Button showResults = (Button) view.findViewById(R.id.buttonShow);
        showResults.setOnClickListener(new View.OnClickListener() {

            @Override
            public void onClick(View view) {
                showResults();
            }
        });

        Toolbar toolbar = (Toolbar) getActivity().findViewById(R.id.toolbar);
        toolbar.setTitle("Monitor");

        if (((AppCompatActivity) getActivity()).getSupportActionBar() != null)
            ((AppCompatActivity) getActivity()).getSupportActionBar().setDisplayHomeAsUpEnabled(false);

        return view;
    }

    @Override
    public void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        //sets the default preferences values if any preference is unset
        PreferenceManager.setDefaultValues(getActivity(), R.xml.preferences, false);
    }

    public String getResults(Context context) {
        XMLProduce XML = new XMLProduce(context);
        return XML.getXML();
    }

    public void showResults() {
        String results = getResults(getActivity());
        Bundle bundle = new Bundle();
        bundle.putString("results", results);
        ShowFragment showFragment = new ShowFragment();
        showFragment.setArguments(bundle);
        getFragmentManager().beginTransaction().replace(R.id.fragment_container, showFragment).addToBackStack(null).commit();
    }
}