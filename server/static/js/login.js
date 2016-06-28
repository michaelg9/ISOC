// Listener for Login button
$(document).ready(function(){
    $("#login-btn").on('click', function(){
        var email = $("#email").val();
        var password = $("#password").val();
        var loginParams = {email: email, password: password};
        var loginURL = "../auth/0.1/login?";
        $.post({
            url: loginURL,
            data: loginParams
        }).done(function(data, textStatus, jqXHR) {
            if (data == "Success") {
                window.location = "../dashboard";
            }
        }).fail(function(data, textStatus, jqXHR) {
            $("#alert-wrong-password").html('<div class="alert alert-danger">Wrong password or email.</div>');
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
