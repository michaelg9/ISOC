// TODO: Get User data dynamically
// Hardcoded API key for the user
var apiKey = "37e72ff927f511e688adb827ebf7e157";

// Request variables for API call to get userdata
var requestURL = "../data/0.1/q?"
var requestParams = {appid: apiKey}

// AJAX call to server
var batteryData = $.get({
    url: requestURL,
    data: requestParams
}).done(function(data, textStatus, jqXHR) {
    var ctx = $("#batteryChart");
    var dataJSON = JSON.parse(data);
    var batteryData = dataJSON.devices[0].data.battery;
    var batteryLevel = batteryData.map(function(battery) {
        return battery.value;
    });
    var batteryTimes = batteryData.map(function(battery) {
        return battery.time;
    });
    var batteryChart = new Chart(ctx, {
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
                    gridLines: {
                        display: false
                    },
                    display: false
                }]
            }
        }
    })
});

// Listeners for datepickers
// TODO: Make one for "From" and "To"
// TODO: Find way to generalise to more datepickers
$(document).ready(function(){
    $('.input-group.date').datepicker({
        format: "dd/mm/yyyy"
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
