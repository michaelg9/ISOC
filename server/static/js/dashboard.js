var requestURL = "../data/0.1/user";

// Angular app
var app = angular.module("deviceApp", []);
app.controller("deviceController", function($scope) {
    $scope.deviceInfo = {};
});

function changeDeviceInfo(deviceInfo) {
    var appElement = document.querySelector("[ng-app=deviceApp]");
    var $scope = angular.element(appElement).scope();
    $scope.$apply(function() {
        $scope.deviceInfo = deviceInfo;
    });
}

// AJAX call to server
var batteryChart;
var batteryData = $.get({
    url: requestURL
}).done(function(data, textStatus, jqXHR) {
    var ctx = $("#batteryChart");
    var userData = JSON.parse(data);
    changeDeviceInfo(userData.devices[0].deviceInfo);
    var batteryData = userData.devices[0].data.battery;
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
        type: 'line',
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
});

$(document).ready(function() {
    $('#daterangepicker').daterangepicker({
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

// Listeners for datepickers
/*
$(document).ready(function(){
    $('.input-group.date.startdate').datepicker({
        format: "dd/mm/yyyy"
    }).on('changeDate', function(e) {
        $('.input-group.date.enddate').datepicker('setStartDate', e.date);
        batteryChart.options.scales.xAxes[0].time.min = e.date;
        batteryChart.update();
    });
});

$(document).ready(function(){
    $('.input-group.date.enddate').datepicker({
        format: "dd/mm/yyyy",
        endDate: "0d"
    }).on('changeDate', function(e) {
        $('.input-group.date.startdate').datepicker('setEndDate', e.date);
        batteryChart.options.scales.xAxes[0].time.max = e.date;
        batteryChart.update();
    });
});*/

// Logout the user on logout link
$(document).ready(function(){
    $('.logout').on('click', function() {
        var logoutURL = "../auth/0.1/logout";
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
