angular.module('flemzer')
    .controller('MoviesCtrl', MoviesCtrl);

function MoviesCtrl($routeParams, $http) {
    var self = this;
    self.config = {};

    self.LoadConfig = function() {
        $http.get("/api/v1/config").then(function(response) {
            if (response.status == 200) {
                self.config = response.data;
            }
        });
    };

    self.LoadConfig();

    $http.get("/api/v1/movies").then(function(response) {
        if (response.status == 200) {
            self.movies = response.data;
            return;
        }
    });
}
