// TODO: Comment

var TokenAuth = (function() {
    var tokenURL = "../auth/0.1/token";

    var refreshAccessToken = function() {
        return $.post(tokenURL).done(function(data, textStatus, jqXHR) {
            sessionStorage.accessToken = data.accessToken;
            return jqXHR;
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
            var statusForbidden = 403;
            if (data.status === statusForbidden) {
                // If authentication failed refresh access token
                return refreshAccessToken().done(function() {
                    return request(url, type, params);
                });
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
