package com.isoc.android.monitor;

import android.accounts.AbstractAccountAuthenticator;
import android.accounts.Account;
import android.accounts.AccountAuthenticatorResponse;
import android.accounts.AccountManager;
import android.accounts.NetworkErrorException;
import android.app.Service;
import android.content.Context;
import android.content.Intent;
import android.os.Bundle;
import android.os.IBinder;

public class AuthenticatorService extends Service {
    private Authenticator authenticator;
    public AuthenticatorService() {
    }

    @Override
    public void onCreate() {
        super.onCreate();
        authenticator=new Authenticator(this);
    }

    @Override
    public IBinder onBind(Intent intent) {
        return authenticator.getIBinder();
    }


    private static class Authenticator extends AbstractAccountAuthenticator {
        private Context context;

        public Authenticator(Context context) {
            super(context);
            this.context=context;
        }

        @Override
        public Bundle editProperties(AccountAuthenticatorResponse accountAuthenticatorResponse, String s) {
            return null;
        }

        @Override
        public Bundle addAccount(AccountAuthenticatorResponse accountAuthenticatorResponse, String accountType, String authTokenType, String[] requiredFeatures, Bundle options) throws NetworkErrorException {
            final Intent intent=new Intent(context,LoginActivity.class);
            intent.putExtra(LoginActivity.ARG_ACCOUNT_TYPE,accountType);
            intent.putExtra(LoginActivity.ARG_AUTH_TYPE,authTokenType);
            intent.putExtra(LoginActivity.ARG_IS_ADDING_NEW_ACCOUNT,true);
            intent.putExtra(AccountManager.KEY_ACCOUNT_AUTHENTICATOR_RESPONSE, accountAuthenticatorResponse);
            final Bundle bundle=new Bundle();
            bundle.putParcelable(AccountManager.KEY_INTENT,intent);
            return bundle;
        }

        @Override
        public Bundle confirmCredentials(AccountAuthenticatorResponse accountAuthenticatorResponse, Account account, Bundle bundle) throws NetworkErrorException {
            return null;
        }

        @Override
        public Bundle getAuthToken(AccountAuthenticatorResponse accountAuthenticatorResponse, Account account, String authTokenType, Bundle bundle) throws NetworkErrorException {
            final AccountManager am=AccountManager.get(context);
            String authToken=am.peekAuthToken(account,authTokenType);
            final Bundle result=new Bundle();
            if (authToken!=null){
                result.putString(AccountManager.KEY_ACCOUNT_NAME,account.name);
                result.putString(AccountManager.KEY_ACCOUNT_TYPE,account.type);
                result.putString(AccountManager.KEY_AUTHTOKEN,authToken);
            }else {
                Intent intent=new Intent(context, LoginActivity.class);
                intent.putExtra(AccountManager.KEY_ACCOUNT_AUTHENTICATOR_RESPONSE,accountAuthenticatorResponse);
                intent.putExtra(LoginActivity.ARG_ACCOUNT_TYPE,account.type);
                intent.putExtra(LoginActivity.ARG_AUTH_TYPE,authTokenType);
                result.putParcelable(AccountManager.KEY_INTENT, intent);
            }
            return result;
        }

        @Override
        public String getAuthTokenLabel(String s) {
            return null;
        }

        @Override
        public Bundle updateCredentials(AccountAuthenticatorResponse accountAuthenticatorResponse, Account account, String s, Bundle bundle) throws NetworkErrorException {
            return null;
        }

        @Override
        public Bundle hasFeatures(AccountAuthenticatorResponse accountAuthenticatorResponse, Account account, String[] strings) throws NetworkErrorException {
            return null;
        }
    }
}
