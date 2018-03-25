angular.module('flemzer')
    .controller('TVShowsCtrl', TVShowsCtrl);

function TVShowsCtrl($routeParams, $http) {
    var self = this;

    $http.get("/api/v1/tvshows").then(function(response) {
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
                console.log(shows[show]);
            }
            self.tvshows = shows;
            return;
        }
    });
}
