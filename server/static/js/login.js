// Listener for Login button
$(document).ready(function(){
    $("#login-btn").on('click', function(){
        var username = $("#username").val();
        var password = $("#password").val();
        var loginParams = {username: username, password: password};
        var loginURL = "../auth/0.1/login?";
        $.post({
            url: loginURL,
            data: loginParams
        }).done(function(data, textStatus, jqXHR) {
            if (data == "Success") {
                window.location = "../dashboard";
            }
        }).fail(function() {
            console.log("Wrong password!");
            $("#alert-wrong-password").html('<div class="alert alert-danger">Wrong password or username.</div>');
            return false;
        });
    });
});

// On enter hit click Login button
$(document).ready(function() {
    $(document).keypress(function(e) {
        var enter = 13;
        var key = e.which;
        if (key == enter) {
            $("#login-btn").click();
            return false;
        }
    });
});
