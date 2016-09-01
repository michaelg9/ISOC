var Downloads = (function() {
    // Get the given feature for the specified device in the given format.
    var getFeature = function(deviceID, feature, format) {
        var url = "/data/" + sessionStorage.userID + "/" + deviceID + "/" + feature;
        return TokenAuth.makeAuthRequest(url, "GET", {out: format});
    };

    // Get all the data of the feature not just from the current user.
    var getAllOfFeature = function(feature, format) {
        var url = "/data/all/features/" + feature;
        return TokenAuth.makeAuthRequest(url, "GET", {out: format});
    };

    return {
        getFeature: getFeature,
        getAllOfFeature: getAllOfFeature
    };
})();
