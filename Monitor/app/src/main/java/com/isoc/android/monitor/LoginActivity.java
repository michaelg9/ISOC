package com.isoc.android.monitor;

import android.accounts.Account;
import android.accounts.AccountAuthenticatorActivity;
import android.accounts.AccountManager;
import android.animation.Animator;
import android.animation.AnimatorListenerAdapter;
import android.app.LoaderManager.LoaderCallbacks;
import android.content.ContentResolver;
import android.content.CursorLoader;
import android.content.Intent;
import android.content.Loader;
import android.database.Cursor;
import android.net.Uri;
import android.os.AsyncTask;
import android.os.Bundle;
import android.provider.ContactsContract;
import android.text.TextUtils;
import android.util.Log;
import android.view.KeyEvent;
import android.view.View;
import android.view.View.OnClickListener;
import android.view.inputmethod.EditorInfo;
import android.widget.ArrayAdapter;
import android.widget.AutoCompleteTextView;
import android.widget.Button;
import android.widget.EditText;
import android.widget.TextView;

import java.util.ArrayList;
import java.util.List;

/**
 * A login screen that offers login via email and password.
 * Triggered when there's no account stored in the Account manager
 * or if it the refresh token is expired.
 * Only the email and the refreshToken (not the password) are saved
 * in the AccountManager for security purposes.
 * An AccessToken is requested each time we're about to upload xml data
 */
public class LoginActivity extends AccountAuthenticatorActivity implements LoaderCallbacks<Cursor> {
    public static final String ARG_ACCOUNT_TYPE = "account_type";
    public static final String ARG_AUTH_TYPE = "authentication_type";
    public static final String ARG_IS_ADDING_NEW_ACCOUNT = "is_new_account";
    public static final String REGISTER_SUCCESS="registration_success";

    //Keep track of the login task to ensure we can cancel it if requested.
    private UserLoginTask mAuthTask = null;

    // UI references.
    private AutoCompleteTextView mEmailView;
    private EditText mPasswordView;
    private View mProgressView;
    private View mLoginFormView;
    private TextView mErrorView;
    private TextView mSignUpLink;

    //interface-contract for user-profile on device
    private interface ProfileQuery {
        String[] PROJECTION = {
                ContactsContract.CommonDataKinds.Email.ADDRESS,
                ContactsContract.CommonDataKinds.Email.IS_PRIMARY,
        };

        int ADDRESS = 0;
    }

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_login);

        // Set up the login form.
        mEmailView = (AutoCompleteTextView) findViewById(R.id.email);
        populateAutoComplete();
        mErrorView=(TextView) findViewById(R.id.text_error_login);

        mPasswordView = (EditText) findViewById(R.id.password);
        mPasswordView.setOnEditorActionListener(new TextView.OnEditorActionListener() {
            @Override
            public boolean onEditorAction(TextView textView, int id, KeyEvent keyEvent) {
                if (id == R.id.login || id == EditorInfo.IME_NULL) {
                    attemptTask(true);
                    return true;
                }
                return false;
            }
        });

        Button mEmailSignInButton = (Button) findViewById(R.id.sign_in_button);
        mEmailSignInButton.setOnClickListener(new OnClickListener() {
            @Override
            public void onClick(View view) {
                attemptTask(true);
            }
        });

        mLoginFormView = findViewById(R.id.login_form);
        mProgressView = findViewById(R.id.login_progress);

        mSignUpLink= (TextView) findViewById(R.id.link_signup_register);
        mSignUpLink.setOnClickListener(new OnClickListener() {
            @Override
            public void onClick(View view) {
                mSignUpLink.setVisibility(View.GONE);
                Button signUpButton=(Button)findViewById(R.id.sign_up_button);
                signUpButton.setVisibility(View.VISIBLE);
                signUpButton.setOnClickListener(new OnClickListener() {
                    @Override
                    public void onClick(View view) {
                        attemptTask(false);
                    }
                });
            }
        });

    }

    private void populateAutoComplete() {
        getLoaderManager().initLoader(0, null, this);
    }

    private boolean isEmailValid(String email) {
        return email.contains("@");
    }

    private boolean isPasswordValid(String password) {
        return password!=null;
    }

      //Shows the progress UI and hides the login form.
    private void showProgress(final boolean show) {
        int shortAnimTime = getResources().getInteger(android.R.integer.config_shortAnimTime);

        mLoginFormView.animate().setDuration(shortAnimTime).alpha(
                show ? 0 : 1).setListener(new AnimatorListenerAdapter() {
            @Override
            public void onAnimationEnd(Animator animation) {
                mLoginFormView.setVisibility(show ? View.GONE : View.VISIBLE);
            }
        });
        mLoginFormView.setVisibility(show ? View.GONE : View.VISIBLE);

        mProgressView.setVisibility(show ? View.VISIBLE : View.GONE);
        mProgressView.animate().setDuration(shortAnimTime).alpha(
                show ? 1 : 0).setListener(new AnimatorListenerAdapter() {
            @Override
            public void onAnimationEnd(Animator animation) {
                mProgressView.setVisibility(show ? View.VISIBLE : View.GONE);
            }
        });
    }

    @Override
    public Loader<Cursor> onCreateLoader(int i, Bundle bundle) {
        return new CursorLoader(this,
                // Retrieve data rows for the device user's 'profile' contact.
                Uri.withAppendedPath(ContactsContract.Profile.CONTENT_URI,
                        ContactsContract.Contacts.Data.CONTENT_DIRECTORY), ProfileQuery.PROJECTION,

                // Select only email addresses.
                ContactsContract.Contacts.Data.MIMETYPE +
                        " = ?", new String[]{ContactsContract.CommonDataKinds.Email
                .CONTENT_ITEM_TYPE},

                // Show primary email addresses first. Note that there won't be
                // a primary email address if the user hasn't specified one.
                ContactsContract.Contacts.Data.IS_PRIMARY + " DESC");
    }

    @Override
    public void onLoadFinished(Loader<Cursor> cursorLoader, Cursor cursor) {
        List<String> emails = new ArrayList<>();
        cursor.moveToFirst();
        while (!cursor.isAfterLast()) {
            emails.add(cursor.getString(ProfileQuery.ADDRESS));
            cursor.moveToNext();
        }
        //Adding monitor account emails too
        AccountManager am = AccountManager.get(this);
        Account[] accounts=am.getAccountsByType(getString(R.string.authenticator_account_type));
        //returned account[] may be empty but never null
        for (Account a : accounts) {
            emails.add(a.name);
        }

        addEmailsToAutoComplete(emails);
    }

    @Override
    public void onLoaderReset(Loader<Cursor> cursorLoader) {}

    private void addEmailsToAutoComplete(List<String> emailAddressCollection) {
        //Create adapter to tell the AutoCompleteTextView what to show in its dropdown list.
        ArrayAdapter<String> adapter =
                new ArrayAdapter<>(LoginActivity.this,
                        android.R.layout.simple_dropdown_item_1line, emailAddressCollection);

        mEmailView.setAdapter(adapter);
    }

    /**
     * Attempts to sign in or register the account specified by the login form.
     * If there are form errors (invalid email, missing fields, etc.), the
     * errors are presented and no actual login attempt is made.
     */
    private void attemptTask(boolean isLoggingIn) {
        if (mAuthTask != null) {
            return;
        }
        // Reset errors.
        mEmailView.setError(null);
        mPasswordView.setError(null);

        // Store values at the time of the login attempt.
        String email = mEmailView.getText().toString();
        String password = mPasswordView.getText().toString();

        boolean cancel = false;
        View focusView = null;

        // Check for a valid password, if the user entered one.
        if (TextUtils.isEmpty(password) || !isPasswordValid(password)) {
            mPasswordView.setError(getString(R.string.error_invalid_password));
            focusView = mPasswordView;
            cancel = true;
        }

        // Check for a valid email address.
        if (TextUtils.isEmpty(email)) {
            mEmailView.setError(getString(R.string.error_field_required));
            focusView = mEmailView;
            cancel = true;
        } else if (!isEmailValid(email)) {
            mEmailView.setError(getString(R.string.error_invalid_email));
            focusView = mEmailView;
            cancel = true;
        }

        if (cancel) {
            // There was an error; don't attempt login and focus the first
            // form field with an error.
            focusView.requestFocus();
        } else {
            // Show a progress spinner, and kick off a background task to
            // perform the user login attempt.
            showProgress(true);
            mAuthTask = new UserLoginTask(email, password,isLoggingIn);
            mAuthTask.execute((Void) null);
        }
    }

    @Override
    public void onBackPressed() {
        moveTaskToBack(true);
    }

    private void finishLogin(Intent accountDetails) {
        String user = accountDetails.getStringExtra(AccountManager.KEY_ACCOUNT_NAME);
        if (getIntent().getBooleanExtra(ARG_IS_ADDING_NEW_ACCOUNT, false)) {
            final Account account = new Account(user, accountDetails.getStringExtra(AccountManager.KEY_ACCOUNT_TYPE));
            String refreshToken = accountDetails.getStringExtra(AccountManager.KEY_AUTHTOKEN);
            AccountManager accountManager = AccountManager.get(getApplicationContext());
            accountManager.addAccountExplicitly(account,null,null);
            accountManager.setAuthToken(account, getString(R.string.token_refresh), refreshToken);
            accountManager.setUserData(account,getString(R.string.am_refreshDateKey),Long.toString(TimeCapture.getCurrentLongTime()));
            ContentResolver.setSyncAutomatically(account,getString(R.string.provider_authority),true);
            //saving device id
            if (accountDetails.hasExtra(getString(R.string.am_deviceID))){
                int dev=accountDetails.getIntExtra(getString(R.string.am_deviceID),-1);
                Log.e("devID",Integer.toString(dev));
                accountManager.setUserData(account,getString(R.string.am_deviceID),Integer.toString(dev));

            }
        }
        setAccountAuthenticatorResult(accountDetails.getExtras());
        setResult(RESULT_OK, accountDetails);
    }

     //Represents an asynchronous login task used to authenticate the user.

    public class UserLoginTask extends AsyncTask<Void, Void, Intent> {
        private final String mEmail;
        private final String mPassword;
        private final boolean mIsLoggingIn;

        UserLoginTask(String email, String password,boolean isLoggingIn) {
            mEmail = email;
            mPassword = password;
            mIsLoggingIn=isLoggingIn;
        }

        @Override
        protected Intent doInBackground(Void... params) {
            Intent accountDetails=new Intent();
            accountDetails.putExtra(AccountManager.KEY_ACCOUNT_NAME,mEmail);
            accountDetails.putExtra(AccountManager.KEY_ACCOUNT_TYPE,getString(R.string.authenticator_account_type));
            if (!mIsLoggingIn){
                //if the user is registering, send a register request first
                String[] registerResponse=new ServerCommunication(getApplicationContext()).register(mEmail,mPassword);
                if (!registerResponse[0].equals(REGISTER_SUCCESS)){
                    //if the request failed do not attempt to login
                    accountDetails.putExtra(registerResponse[0],registerResponse[1]);
                    return accountDetails;
                }else{
                    //save device id
                    Log.e("device-id",registerResponse[1]);
                    accountDetails.putExtra(getString(R.string.am_deviceID),Integer.parseInt(registerResponse[1]));
                }
            }
            //if the register request was successful, save the device id and login
            String[] loginResponse=new ServerCommunication(getApplicationContext()).login(mEmail,mPassword);
            accountDetails.putExtra(loginResponse[0],loginResponse[1]);
            return accountDetails;
        }

        @Override
        protected void onPostExecute(final Intent accountDetails) {
            mAuthTask = null;
            showProgress(false);
            if (!accountDetails.hasExtra(AccountManager.KEY_ERROR_MESSAGE)) {
                finishLogin(accountDetails);
                finish();
            } else {
                mErrorView.setVisibility(View.VISIBLE);
                mErrorView.setText(accountDetails.getStringExtra(AccountManager.KEY_ERROR_MESSAGE));
                mEmailView.setError(getString(R.string.login_try_again));
                mPasswordView.setError(getString(R.string.login_try_again));
                mEmailView.requestFocus();
            }
        }

        @Override
        protected void onCancelled() {
            mAuthTask = null;
            showProgress(false);
        }
    }
}