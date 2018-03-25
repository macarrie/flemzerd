angular.module('flemzer')
    .controller('SettingsCtrl', SettingsCtrl);

function SettingsCtrl($routeParams, $http) {
    var self = this;
    self.trakt_auth = false;
    self.trakt_device_code = {};
    self.trakt_token = {};
    self.trakt_auth_errors = [];
    self.config = {};

    self.LoadConfig = function() {
        $http.get("/api/v1/config").then(function(response) {
            if (response.status == 200) {
                self.config = response.data;
            }
        });
    };

    self.startTraktAuth = function() {
        $http.get("/api/v1/modules/watchlists/trakt/auth").then(function(response) {
            if (response.status == 200) {
                var waitForDeviceCode = setInterval(function() {
                    $http.get("/api/v1/modules/watchlists/trakt/devicecode").then(function(response) {
                        if (response.data.device_code != "") {
                            self.trakt_device_code = response.data;
                            clearInterval(waitForDeviceCode);
                            var waitForTraktToken = setInterval(function() {
                                self.GetTraktToken();
                                self.checkTraktAuthErrors();

                                if (self.trakt_auth_errors.length != 0) {
                                    self.trakt_auth = false;
                                    clearInterval(waitForTraktToken);
                                    return;
                                }

                                if (Object.keys(self.trakt_token).length !== 0) {
                                    if (self.trakt_token.access_token != "") {
                                        clearInterval(waitForTraktToken);
                                    }
                                }
                            }, 3000);
                            setTimeout(function() {
                                clearInterval(waitForTraktToken);
                            }, self.trakt_device_code.expires_in * 1000);
                        }
                    });
                }, 2000);
                setTimeout(function() {
                    clearInterval(waitForDeviceCode);
                }, 10000);
                return;
            } else if (response.status == 204) {
                console.log("Trakt auth already done");
            } else {
                console.log("Trakt auth error");
            }
            return;
        });
    };

    self.checkTraktAuthErrors = function() {
        $http.get("/api/v1/modules/watchlists/trakt/auth_errors").then(function(response) {
            if (response.status == 200) {
                self.trakt_auth_errors = response.data;
            }
            return;
        });
    };

    self.GetTraktToken = function() {
        $http.get("/api/v1/modules/watchlists/trakt/token").then(function(response) {
            if (response.status == 200) {
                if (response.data.access_token != "") {
                    self.trakt_auth = true;
                    self.trakt_token = response.data
                }
            }
            return;
        });
    };

    self.LoadConfig();
    self.GetTraktToken();
}
