var User = (function() {
    var updateUserURL = "../update/user?";

    // Used to store the info about the current user
    var info = {};
    var devices = [];
    var currentDevice = {};

    var getCurrentDevice = function() {
        return currentDevice;
    };

    // Get the info for the current user from the server, bind info
    // to rivetjs and create graphs.
    var initUser = function() {
        var userURL = "../data/" + sessionStorage.userID;
        TokenAuth.makeAuthRequest(userURL, "GET", {}).done(function(result) {
            setUserInfo(result.user);
            devices = result.devices;
            currentDevice = devices[0];
            rivets.bind($("#userInfo"), {userInfo: info});
            rivets.bind($("#deviceInfo"), {deviceInfo: currentDevice.aboutDevice});
            Graphs.createBatteryGraph($("#batteryGraph"), currentDevice.data.battery);
        }).fail(function(result) {
            console.error(result);
        });
    };

    // Update the rivetjs binding with the given data.
    var setUserInfo = function (data) {
        // We have to update each attribute seperately because otherwise
        // rivetjs does not update the view
        info.email = data.email;
        info.apiKey = data.apiKey;
    };

    // Get new userdata from the server
    var updateUserInfo = function() {
        var userDataURL = "../data/" + sessionStorage.userID;
        TokenAuth.makeAuthRequest(userDataURL, "GET", {}).done(function(data) {
            setUserInfo(data.user);
        }).fail(function (result) {
            console.error(result);
        });
    };

    var updateEmail = function(newEmail) {
        var updateData = {email: newEmail};
        TokenAuth.makeAuthRequest(updateUserURL, "POST", updateData).done(function () {
            updateUserInfo();
        }).fail(function(result) {
            console.error(result);
        });
    };

    var updateAPIKey = function() {
        var updateData = {apiKey: "1"}; // Use 1 for true
        TokenAuth.makeAuthRequest(updateUserURL, "POST", updateData).done(function () {
            updateUserInfo();
        }).fail(function(result) {
            console.error(result);
        });
    };

    return {
        getCurrentDevice: getCurrentDevice,
        initUser: initUser,
        updateEmail: updateEmail,
        updateAPIKey: updateAPIKey
    };
})();
