var batteryChart;

// TODO: Restructure in seperate files and comment

function createBatteryGraph(batteryData) {
    // Sort data according to time so it gets displayed properly
    batteryData.sort(function(a,b){
        var dateA = new Date(a.time);
        var dateB = new Date(b.time);
        return dateB - dateA;
    });

    var batteryLevel = batteryData.map(function(battery) {
        return battery.value;
    });
    var batteryTimes = batteryData.map(function(battery) {
        return battery.time;
    });

    // Get context
    var ctx = $("#batteryChart");
    batteryChart = new Chart(ctx, {
        type: "line",
        data: {
            labels: batteryTimes,
            datasets: [{
                label: "Moto X",
                borderColor: "rgba(75,192,192,1)",
                pointBorderColor: "rgba(75,192,192,1)",
                pointBackgroundColor: "#fff",
                pointBorderWidth: 1,
                pointHoverRadius: 5,
                pointHoverBackgroundColor: "rgba(75,192,192,1)",
                pointHoverBorderColor: "rgba(220,220,220,1)",
                pointHoverBorderWidth: 2,
                pointRadius: 1,
                pointHitRadius: 10,
                data: batteryLevel
            }]
        },
        options: {
            scales: {
                yAxes: [{
                    gridLines: {
                        display: false
                    },
                    display: true,
                    ticks: {
                        beginAtZero: true,
                        max: 100
                    }
                }],
                xAxes: [{
                    type: 'time',
                    time: {
                        max: null,
                        min: null
                    },
                    gridLines: {
                        display: false
                    },
                    display: true
                }]
            }
        }
    });
}

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
            var statusUnauthorized = 401;
            if (jqXHR.status === statusUnauthorized) {
                // If authentication failed refresh access token
                refreshAccessToken();
                return request(url, type, params);
            }
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
        logout: logout
    };
})();

var user = (function() {
    var updateUserURL = "../update/user?";

    // Used to store the info about the current user
    var info = {};
    var devices = [];
    var currentDevice = {};

    var getCurrentDevice = function() {
        return currentDevice;
    };

    var getUser = function() {
        var userURL = "../data/" + sessionStorage.userID;
        tokenAuth.makeAuthRequest(userURL, "GET", {}).done(function(result) {
            setUserInfo(result.user);
            devices = result.devices;
            currentDevice = devices[0];
            rivets.bind($("#userInfo"), {userInfo: info});
            rivets.bind($("#deviceInfo"), {deviceInfo: currentDevice.aboutDevice});
            createBatteryGraph(currentDevice.data.battery);
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
        getUser: getUser,
        updateEmail: updateEmail,
        updateAPIKey: updateAPIKey
    };
})();

var downloads = (function() {
    var getFeature = function(deviceID, feature, format) {
        var featureURL = "../data/" + sessionStorage.userID + "/" + deviceID + "/" + feature;
        return tokenAuth.makeAuthRequest(featureURL, "GET", {out: format});
    };

    return {
        getFeature: getFeature
    };
})();

// JQuery listeners
$(document).ready(function() {
    user.getUser();

    // Listener for daterangepicker
    $("#daterangepicker").daterangepicker({
        startDate: moment().subtract(7, "days"),
        endDate: moment(),
        maxDate: moment(),
        locale: {
            format: "MMMM D, YYYY"
        }
    }, function(startDate, endDate, label) {
        batteryChart.options.scales.xAxes[0].time.min = startDate;
        batteryChart.options.scales.xAxes[0].time.max = endDate;
        batteryChart.update();
    });

    // Logout the user on logout link
    $(".logout").on("click", function() {
        tokenAuth.logout();
    });

    // Add the modal prompt for new email
    $("#editEmail").on("click", function(e) {
        e.preventDefault();
        bootbox.prompt("Please enter your new email", function(result) {
            if (result !== "") {
                user.updateEmail(result);
            }
        });
    });

    $("#updateAPIKey").on("click", function (e) {
        e.preventDefault();
        user.updateAPIKey();
    });

    $("#batteryJSON").on("click", function() {
        downloads.getFeature(user.getCurrentDevice().aboutDevice.id, "Battery", "json").done(function(data) {
            var blob = new Blob([JSON.stringify(data)], {type: "application/json"});
            saveAs(blob, "Battery.json");
        }).fail(function(data, textStatus, jqXHR) {
            console.error(textStatus);
        });
    });
});
