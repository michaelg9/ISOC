package com.isoc.android.monitor;


import android.os.Bundle;
import android.support.v4.app.Fragment;
import android.text.method.ScrollingMovementMethod;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.TextView;


/**
 * A simple {@link Fragment} subclass.
 */
public class ShowFragment extends android.app.Fragment {


    public ShowFragment() {
        // Required empty public constructor
    }


    @Override
    public View onCreateView(LayoutInflater inflater, ViewGroup container,
                             Bundle savedInstanceState) {
        // Inflate the layout for this fragment
        View view= inflater.inflate(R.layout.fragment_show, container, false);
        TextView textView=(TextView) view.findViewById(R.id.frag_show_text);
        textView.setText(getArguments().getString("results"));
        textView.setMovementMethod(new ScrollingMovementMethod());

        return view;
    }

}
