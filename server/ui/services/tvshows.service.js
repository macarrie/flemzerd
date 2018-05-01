flemzerd.factory('tvshows', ['$rootScope', '$interval', '$http', function($rootScope, $interval, $http) {
    var loadShows = function() {
        $http.get("/api/v1/tvshows/tracked").then(function(response) {
            if (response.status == 200) {
                shows = response.data;
                for (var show in shows) {
                    switch (shows[show]["Status"]) {
                        case 1:
                            shows[show]["StatusText"] = "Continuing";
                            break;
                        case 2:
                            shows[show]["StatusText"] = "Planned";
                            break;
                        case 3:
                            shows[show]["StatusText"] = "Ended";
                            break;
                        default:
                            shows[show]["StatusText"] = "Unknown";
                            break;
                    }
                }

                $rootScope.tvshows.tracked = shows;
                return;
            }
        });
        $http.get("/api/v1/tvshows/downloading").then(function(response) {
            if (response.status == 200) {
                $rootScope.tvshows.downloading = response.data;
                return;
            }
        });
        $http.get("/api/v1/tvshows/downloaded").then(function(response) {
            if (response.status == 200) {
                $rootScope.tvshows.downloaded = response.data;
                return;
            }
        });
    };

    $rootScope.tvshows = {}
    loadShows();
    $interval(loadShows, 30000);

    return {
        "tvshows": $rootScope.tvshows,
    }
}]);
