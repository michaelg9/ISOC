var apiKey = "37e72ff927f511e688adb827ebf7e157";
var requestURL = "../data/0.1/q?"
var requestParams = {appid: apiKey}
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
    console.log(batteryLevel);
    console.log(batteryTimes);
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

$(document).ready(function(){
    $('.input-group.date').datepicker({
        format: "dd/mm/yyyy"
    });
});
