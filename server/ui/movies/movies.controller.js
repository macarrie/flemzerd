angular.module('flemzer')
    .controller('MoviesCtrl', MoviesCtrl);

function MoviesCtrl($routeParams, $http) {
    var self = this;

    $http.get("/api/v1/movies").then(function(response) {
        if (response.status == 200) {
            self.movies = response.data;
            return;
        }
    });
}
