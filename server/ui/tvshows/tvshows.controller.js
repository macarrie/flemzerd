angular.module('flemzer')
    .controller('TVShowsCtrl', TVShowsCtrl);

function TVShowsCtrl($routeParams, $http) {
    var self = this;

    $http.get("/api/v1/tvshows").then(function(response) {
        if (response.status == 200) {
            self.tvshows = response.data;
            return;
        }
    });
}
