// TODO: Comment
// TODO: RequireJS and capitalize modules

// JQuery listeners
$(document).ready(function() {
    // Get the userdata from the server
    user.initUser();

    // Listener for daterangepicker
    $("#daterangepicker").daterangepicker({
        startDate: moment().subtract(7, "days"),
        endDate: moment(),
        maxDate: moment(),
        locale: {
            format: "MMMM D, YYYY"
        }
    }, function(startDate, endDate, label) {
        // When user selects daterange update the graph
        graphs.updateBatteryGraph(startDate, endDate);
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
        // Get the battery data from the current device in JSON format
        downloads
        .getFeature(user.getCurrentDevice().aboutDevice.id, "Battery", "json")
        .done(function(data) {
            var blob = new Blob([JSON.stringify(data)], {type: "application/json"});
            saveAs(blob, "Battery.json");
        }).fail(function(data, textStatus, jqXHR) {
            console.error(textStatus);
        });
    });
});
