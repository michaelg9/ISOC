var graphs = (function() {
    var batteryGraph;

    var createBatteryGraph = function (context, batteryData) {
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
        batteryGraph = new Chart(context, {
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
    };

    var updateBatteryGraph = function(startDate, endDate) {
        batteryGraph.options.scales.xAxes[0].time.min = startDate;
        batteryGraph.options.scales.xAxes[0].time.max = endDate;
        batteryGraph.update();
    };

    return {
        createBatteryGraph: createBatteryGraph,
        updateBatteryGraph: updateBatteryGraph
    };
})();
