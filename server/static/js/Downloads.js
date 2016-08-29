var Downloads = (function() {
    // Get the given feature for the specified device in the given format.
    var getFeature = function(deviceID, feature, format) {
        var featureURL = "../data/" + sessionStorage.userID + "/" + deviceID + "/" + feature;
        return TokenAuth.makeAuthRequest(featureURL, "GET", {out: format});
    };

    return {
        getFeature: getFeature
    };
})();
