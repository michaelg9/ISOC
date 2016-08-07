package com.isoc.android.monitor;

import android.content.Context;
import android.content.res.TypedArray;
import android.os.Parcel;
import android.os.Parcelable;
import android.preference.DialogPreference;
import android.util.AttributeSet;
import android.view.View;
import android.view.ViewGroup;
import android.widget.NumberPicker;

/**
 * Custom preference for choosing an integer. Used to select the update interval of the service
 * and the interval to sent data to the server
 */
public class NumberPreference extends DialogPreference {
    private static final int DEFAULT_VALUE=5;
    private static final int DEFAULT_MIN_VALUE=1;
    private static final int DEFAULT_MAX_VALUE=60;

    private int min;
    private int max;
    private String time;

    private int timer;
    private NumberPicker numberPicker;

    public NumberPreference(Context context, AttributeSet attrs) {
        super(context, attrs);

        setDialogLayoutResource(R.layout.number_preference);
        setNegativeButtonText("Cancel");
        setPositiveButtonText("OK");
        TypedArray a = context.getTheme().obtainStyledAttributes(attrs, R.styleable.number_preference, 0, 0);
        try{
            min=a.getInteger(R.styleable.number_preference_min, DEFAULT_MIN_VALUE);
            max=a.getInteger(R.styleable.number_preference_max, DEFAULT_MAX_VALUE);
            time=a.getString(R.styleable.number_preference_time);
        }finally{
            a.recycle();
        }
        setDialogIcon(null);
    }


    public void setSummary() {
        super.setSummary("Every "+getTimer()+' '+time);
    }

    @Override
    protected View onCreateView(ViewGroup parent) {
        View result=super.onCreateView(parent);
        setSummary();
        return result;
    }

    @Override
    protected void onDialogClosed(boolean positiveResult) {
        if (positiveResult) {
            int number = numberPicker.getValue();
            if (callChangeListener(number)){
                timer=number;
                persistInt(timer);
                setSummary();
            }
        }
    }

    @Override
    protected Object onGetDefaultValue(TypedArray a, int index) {
        return a.getInt(index,DEFAULT_VALUE);
    }

    @Override
    protected void onSetInitialValue(boolean restorePersistedValue, Object defaultValue) {
        if (restorePersistedValue) {
            timer = getPersistedInt(DEFAULT_VALUE);
        }
        else{
            timer =(Integer) defaultValue;
            persistInt(timer);
        }
    }

    @Override
    protected void onBindDialogView(View view) {
        super.onBindDialogView(view);
        numberPicker=(NumberPicker) view.findViewById(R.id.numpref_picker);
        numberPicker.setMinValue(min);
        numberPicker.setMaxValue(max);
        numberPicker.setValue(timer);
    }

    public int getTimer() {
        return getPersistedInt(DEFAULT_VALUE);
    }

    @Override
    protected Parcelable onSaveInstanceState() {
        final Parcelable superState = super.onSaveInstanceState();
        final SavedState myState=new SavedState(superState);
        if (numberPicker!=null)myState.value=numberPicker.getValue();
        return myState;
    }

    @Override
    protected void onRestoreInstanceState(Parcelable state) {
        if (state==null || !state.getClass().equals(SavedState.class)){
            super.onRestoreInstanceState(state);
            return;
        }
        SavedState myState=(SavedState)state;
        super.onRestoreInstanceState(myState.getSuperState());
        if (numberPicker!=null) numberPicker.setValue(myState.value);
    }

    private static class SavedState extends BaseSavedState {
        // field that holds the setting's value
        int value;


        public SavedState(Parcelable superState) {
            super(superState);
        }

        public SavedState(Parcel source) {
            super(source);
            // Get the current preference's value
            value = source.readInt();
        }

        @Override
        public void writeToParcel(Parcel dest, int flags) {
            super.writeToParcel(dest, flags);
            // Write the preference's value
            dest.writeInt(value);
        }

        // Standard creator object using an instance of this class
        public static final Parcelable.Creator<SavedState> CREATOR =
                new Parcelable.Creator<SavedState>() {

                    public SavedState createFromParcel(Parcel in) {
                        return new SavedState(in);
                    }

                    public SavedState[] newArray(int size) {
                        return new SavedState[size];
                    }
                };
    }

}
