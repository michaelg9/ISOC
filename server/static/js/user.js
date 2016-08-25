var user = (function() {
    var updateUserURL = "../update/user?";

    // Used to store the info about the current user
    var info = {};
    var devices = [];
    var currentDevice = {};

    var getCurrentDevice = function() {
        return currentDevice;
    };

    var initUser = function() {
        var userURL = "../data/" + sessionStorage.userID;
        tokenAuth.makeAuthRequest(userURL, "GET", {}).done(function(result) {
            setUserInfo(result.user);
            devices = result.devices;
            currentDevice = devices[0];
            rivets.bind($("#userInfo"), {userInfo: info});
            rivets.bind($("#deviceInfo"), {deviceInfo: currentDevice.aboutDevice});
            graphs.createBatteryGraph($("#batteryGraph"), currentDevice.data.battery);
        }).fail(function(result) {
            console.error(result);
        });
    };

    var setUserInfo = function (data) {
        // We have to update each attribute seperately because otherwise
        // rivetjs does not update the view
        info.email = data.email;
        info.apiKey = data.apiKey;
    };

    var updateUserInfo = function() {
        var userDataURL = "../data/" + sessionStorage.userID;
        tokenAuth.makeAuthRequest(userDataURL, "GET", {}).done(function(data) {
            setUserInfo(data.user);
        }).fail(function (result) {
            console.error(result);
        });
    };

    var updateEmail = function(newEmail) {
        var updateData = {email: newEmail};
        tokenAuth.makeAuthRequest(updateUserURL, "POST", updateData).done(function () {
            updateUserInfo();
        }).fail(function(result) {
            console.error(result);
        });
    };

    var updateAPIKey = function() {
        var updateData = {apiKey: "1"}; // Use 1 for true
        tokenAuth.makeAuthRequest(updateUserURL, "POST", updateData).done(function () {
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
