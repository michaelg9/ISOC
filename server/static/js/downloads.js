var downloads = (function() {
    var getFeature = function(deviceID, feature, format) {
        var featureURL = "../data/" + sessionStorage.userID + "/" + deviceID + "/" + feature;
        return tokenAuth.makeAuthRequest(featureURL, "GET", {out: format});
    };

    return {
        getFeature: getFeature
    };
})();
