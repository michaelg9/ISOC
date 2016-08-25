var tokenAuth = (function() {
    var tokenURL = "../auth/0.1/token";

    var refreshAccessToken = function() {
        $.post(tokenURL).done(function(data) {
            sessionStorage.accessToken = data.accessToken;
        });
    };

    var request = function(url, type, params) {
        return $.ajax(url, {
            type: type,
            data: params,
            beforeSend: function(jqXHR) {
                jqXHR.setRequestHeader("Authorization", "Bearer " + sessionStorage.accessToken);
            }
        });
    };

    var makeAuthRequest = function(url, type, params) {
        return request(url, type, params).done(function(data, textStatus, jqXHR) {
            return jqXHR;
        }).fail(function(data, statusText, jqXHR) {
            // FIXME
            var statusUnauthorized = 401;
            if (jqXHR.status === statusUnauthorized) {
                // If authentication failed refresh access token
                refreshAccessToken();
                return request(url, type, params);
            }
            return jqXHR;
        });
    };

    var login = function(email, password) {
        var loginParams = {email: email, password: password};
        var loginURL = "../login?";
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
        var logoutURL = "../logout";
        $.post({
            url: logoutURL
        }).done(function(data, textStatus, jqXHR) {
            if (data == "Success") {
                window.location = "../";
            }
        }).fail(function() {
            // This should never happen
            console.error("Failed logout!");
            return false;
        });
    };

    return {
        makeAuthRequest: makeAuthRequest,
        login: login,
        logout: logout
    };
})();
