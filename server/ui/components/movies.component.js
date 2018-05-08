flemzerd.component("movies", {
    templateUrl: "/static/templates/movies.template.html",
    controller: function MoviesCtrl($rootScope, $scope, $http, $interval, config, movies) {
        var self = this;
        $scope.utils = $rootScope.utils;
        $scope.config = {};
        $scope.movies = {};

        self.refresh = function() {
            $scope.config = $rootScope.config;
            $scope.movies = $rootScope.movies;
        };

        $scope.refreshMovies = function() {
            $scope.refreshing = true;
            $http.post("/api/v1/modules/watchlists/refresh").then(function(response) {
                if (response.status == 200) {
                    movies.loadMovies();
                    $scope.refreshing = false;
                    return;
                }
            });
        };

        self.refresh();
        $interval(self.refresh, 1000);

        return;
    }
});
