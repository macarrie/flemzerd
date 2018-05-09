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
        $http.get("/api/v1/tvshows/removed").then(function(response) {
            if (response.status == 200) {
                $rootScope.tvshows.removed = response.data;
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

    var deleteShow = function(id) {
        if (id !== 0) {
            $http.delete("/api/v1/tvshows/details/" +id).then(function(response) {
                if (response.status == 204) {
                    loadShows();
                    return true;
                }
            });
        }
        return false;
    };

    var restoreShow = function(id) {
        if (id !== 0) {
            $http.post("/api/v1/tvshows/restore/" +id).then(function(response) {
                if (response.status == 200) {
                    loadShows();
                    return true;
                }
            });
        }
        return false;
    };

    $rootScope.tvshows = {}
    loadShows();
    $interval(loadShows, 30000);

    return {
        setStatusText: setStatusText,
        loadShows: loadShows,
        deleteShow: deleteShow,
        restoreShow: restoreShow
    };
}]);
