var requestURL = "../data/0.1/user";
var batteryChart;

// TODO: Commenting
// TODO: Look into global variables in JS and use JS linter

// Angular app
var app = angular.module("dashboardApp", []);
app.controller("deviceController", function($scope) {
    $scope.deviceInfo = {};
});
app.controller("userController", function($scope) {
    $scope.userInfo = {}
});

function changeDeviceInfo(deviceInfo) {
    var controllerElement = document.querySelector("[ng-controller=deviceController]");
    var $scope = angular.element(controllerElement).scope();
    $scope.$apply(function() {
        $scope.deviceInfo = deviceInfo;
    });
}

function changeUserInfo(userInfo) {
    var controllerElement = document.querySelector("[ng-controller=userController]");
    var $scope = angular.element(controllerElement).scope();
    $scope.$apply(function() {
        $scope.userInfo = userInfo;
    });
}


// AJAX call to server
var batteryData = $.get({
    url: requestURL
}).done(function(data, textStatus, jqXHR) {
    var ctx = $("#batteryChart");
    var userData = JSON.parse(data);
    changeDeviceInfo(userData.devices[0].deviceInfo);
    changeUserInfo(userData.user);
    // TODO: Check if that works
    createBatteryGraph(userData.devices[0].data.battery)
});

function createBatteryGraph(batteryData) {
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
    })
}

// TODO: Check if you can put all that into one

// Listener for daterangepicker
$(document).ready(function() {
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
});

// Logout the user on logout link
$(document).ready(function() {
    $(".logout").on("click", function() {
        var logoutURL = "../logout";
        $.post({
            url: logoutURL
        }).done(function(data, textStatus, jqXHR) {
            if (data == "Success") {
                window.location = "../";
            }
        }).fail(function() {
            // This should never happen
            console.log("Failed logout!");
            return false;
        });
    });
});

// Add the modal prompt for new email
$(document).ready(function() {
    $("#editEmail").on("click", function() {
        bootbox.prompt("Please enter your new email", function(result) {
            console.log(result);
        });
    });
});
