flemzerd.component("movies", {
    templateUrl: "/static/templates/movies.template.html",
    controller: function MoviesCtrl($rootScope, $scope, $http, $interval, config, movies) {
        var self = this;
        $scope.config = {};
        $scope.movies = {};

        $scope.getFailedTorrentsNb = function(dl) {
            return Object.keys(dl.FailedTorrents).length;
        };

        self.refresh = function() {
            $scope.config = $rootScope.config;
            $scope.movies = $rootScope.movies;
        };

        self.refresh();
        $interval(self.refresh, 1000);

        return;
    }
});
