var TokenAuth = (function() {

    var refreshAccessToken = function() {
        var tokenURL = "/auth/0.1/token";
        // Request a new access token, save it in the sessionStorage and
        // return the request, so that we can handle the error.
        return $.post(tokenURL).done(function(data, textStatus, jqXHR) {
            sessionStorage.accessToken = data.accessToken;
            return jqXHR;
        });
    };



    var makeAuthRequest = function(url, type, params) {
        // "request" makes an AJAX request to the server authorised
        // with the access token saved in the sessions storage
        var request = function(url, type, params) {
            return $.ajax(url, {
                type: type,
                data: params,
                beforeSend: function(jqXHR) {
                    jqXHR.setRequestHeader("Authorization", "Bearer " + sessionStorage.accessToken);
                }
            });
        };

        return request(url, type, params).done(function(data, textStatus, jqXHR) {
            // If the request was successfull simply return it
            return jqXHR;
        }).fail(function(data, statusText, jqXHR) {
            var statusForbidden = 403;
            if (data.status === statusForbidden) {
                // If authentication failed refresh access token
                return refreshAccessToken().done(function() {
                    // If we succesfully got a new access token
                    // make a new request
                    return request(url, type, params);
                });
            }
            // If the request did not fail because of the access token
            // return the request for custom error handling
            return jqXHR;
        });
    };

    // "login" gets the user ID and the access token and stores them in
    // the session storage. The refresh token is stored in an encrypted cookie
    // and automatically refreshed by the server on each request (hence the user
    // gets logged out automatically after a week of inactivity).
    // TODO: Protect against CSRF.
    var login = function(email, password) {
        var loginParams = {email: email, password: password};
        var loginURL = "/login?";
        return $.post({
            url: loginURL,
            data: loginParams
        }).done(function(data, textStatus, jqXHR) {
            sessionStorage.userID = data.id;
            sessionStorage.accessToken = data.accessToken;
            return jqXHR;
        });
    };

    var logout = function() {
        var logoutURL = "/logout";
        return $.post({
            url: logoutURL
        }).done(function(data, statusText, jqXHR) {
            // If logout was successfull delete the accessToken and user ID
            sessionStorage.clear();
            return jqXHR;
        });
    };

    return {
        makeAuthRequest: makeAuthRequest,
        login: login,
        logout: logout
    };
})();
