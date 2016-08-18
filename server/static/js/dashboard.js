var retrieveDataURL = "../data/0.1/user";
var batteryChart;

// TODO: Use jQuery instead of Angular
// TODO: Error handling

// Angular app
var app = angular.module("dashboardApp", []);
app.factory("userService", function($http) {
    return {
        getUser: function() {
            return $http.get("../data/" + sessionStorage.userID, {
                headers: {"Authorization": "Bearer " + sessionStorage.accessToken}
            });
        },
        saveUser: function($scope) {
            return function(response) {
                $scope.user = response.data;
                $scope.deviceInfo = $scope.user.devices[0].aboutDevice;
                $scope.userInfo = $scope.user.user;
                createBatteryGraph($scope.user.devices[0].data.battery);
            };
        }
    };
});

app.factory("downloadService", function($http) {
    return {
        getFeature: function(deviceID, feature, type) {
            return $http.get("../data/" + sessionStorage.userID + "/" + deviceID + "/" + feature + "?out=" + type, {
                headers: {"Authorization": "Bearer " + sessionStorage.accessToken},
                responseType: "arraybuffer"
            });
        },
        getAccessToken: function() {
            return $http.post("../auth/0.1/token");
        }
    };
});

app.controller("dashboardController", function($scope, userService, downloadService) {
    userService.getUser().then(userService.saveUser($scope)).catch(function(response) {
        if (response.status == 403) {
            downloadService.getAccessToken().then(function(response) {
                sessionStorage.accessToken = response.data.accessToken;
                userService.getUser().then(userService.saveUser($scope));
            });
        }
    });
    $scope.saveBatteryJSON = function() {
        var controllerElement = document.querySelector("[ng-controller=dashboardController]");
        var scope = angular.element(controllerElement).scope();
        var fileName = "battery.json";
        var a = document.createElement("a");
        document.body.appendChild(a);
        a.style = "display: none";
        downloadService.getFeature(scope.deviceInfo.id, "Battery", "json").then(function (result) {
            var file = new Blob([result.data], {type: "application/json"});
            var fileURL = window.URL.createObjectURL(file);
            a.href = fileURL;
            a.download = fileName;
            a.click();
        });
    };
});

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
    var info;
    var devices;

    var changeUserInfo = function (data) {
        var controllerElement = document.querySelector("[ng-controller=dashboardController]");
        var $scope = angular.element(controllerElement).scope();
        $scope.$apply(function() {
            $scope.userInfo = data;
        });
    };

    var updateUserInfo = function() {
        var userDataURL = "../data/" + sessionStorage.userID;
        tokenAuth.makeAuthRequest(userDataURL, "GET", {}).done(function(data) {
            changeUserInfo(data.user);
        }).fail(function (data, textStatus, jqXHR) {
            console.error(data);
        });
    };

    var updateEmail = function(newEmail) {
        var updateData = {email: newEmail};
        tokenAuth.makeAuthRequest(updateUserURL, "POST", updateData).done(function () {
            updateUserInfo();
        }).fail(function(data, textStatus, jqXHR) {
            console.error(data);
        });
    };

    var updateAPIKey = function() {
        var updateData = {apiKey: "1"}; // Use 1 for true
        tokenAuth.makeAuthRequest(updateUserURL, "POST", updateData).done(function () {
            updateUserInfo();
        }).fail(function(data, textStatus, jqXHR) {
            console.error(data);
        });
    };

    return {
        updateEmail: updateEmail,
        updateAPIKey: updateAPIKey
    };
})();

// JQuery listeners
$(document).ready(function() {
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
    $("#editEmail").on("click", function() {
        bootbox.prompt("Please enter your new email", function(result) {
            if (result !== "") {
                user.updateEmail(result);
            }
        });
    });

    $("#updateAPIKey").on("click", function () {
        user.updateAPIKey();
    });
});
