flemzerd.factory('tvshows', ['$rootScope', '$interval', '$http', function($rootScope, $interval, $http) {
    var setStatusText = function(show) {
        switch (show["Status"]) {
            case 1:
                show["StatusText"] = "Continuing";
                break;
            case 2:
                show["StatusText"] = "Planned";
                break;
            case 3:
                show["StatusText"] = "Ended";
                break;
            default:
                show["StatusText"] = "Unknown";
                break;
        }

        return show;
    };

    var loadShows = function() {
        $http.get("/api/v1/tvshows/tracked").then(function(response) {
            if (response.status == 200) {
                shows = response.data;
                for (var show in shows) {
                    shows[show] = setStatusText(shows[show]);
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
    };

    $rootScope.tvshows = {}
    loadShows();
    $interval(loadShows, 30000);

    return {
        "tvshows": $rootScope.tvshows,
        setStatusText: setStatusText,
        loadShows: loadShows
    };
}]);
