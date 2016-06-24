// TODO: Get User data dynamically
// Hardcoded API key for the user
var apiKey = "37e72ff927f511e688adb827ebf7e157";

// Request variables for API call to get userdata
var requestURL = "../data/0.1/q?"
var requestParams = {appid: apiKey}

// AJAX call to server
var batteryChart;
var batteryData = $.get({
    url: requestURL,
    data: requestParams
}).done(function(data, textStatus, jqXHR) {
    var ctx = $("#batteryChart");
    var dataJSON = JSON.parse(data);
    var batteryData = dataJSON.devices[0].data.battery;
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

// Listeners for datepickers
// TODO: Make one for "From" and "To"
// TODO: Find way to generalise to more datepickers
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
});

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
