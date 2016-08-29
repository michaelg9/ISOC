$(document).ready(function(){
    // Listener for Login button
    $("#login-btn").on("click", function(){
        var email = $("#email").val();
        var password = $("#password").val();
        TokenAuth.login(email, password).done(function() {
            // If login was successfull redirect to dashboad
            window.location = "../dashboard";
        }).fail(function(data, textStatus, jqXHR) {
            // If login not successfull display warning
            $("#alert-wrong-password").html('<div class="alert alert-danger">Wrong password or email.</div>');
            return false;
        });
    });

    // On enter click Login button
    $(document).keypress(function(e) {
        var enter = 13;
        var key = e.which;
        if (key == enter) {
            $("#login-btn").click();
            return false;
        }
    });
});
