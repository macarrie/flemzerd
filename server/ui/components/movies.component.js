flemzerd.component("movies", {
    templateUrl: "/static/templates/movies.template.html",
    controller: function MoviesCtrl($rootScope, $scope, $http, $interval, config, movies) {
        var self = this;
        $scope.config = {};
        $scope.movies = {};

        $scope.objLength = function(obj) {
            if (typeof obj == undefined) {
                return 0;
            }
            return Object.keys(obj).length;
        };

        $scope.getYear = function(movie) {
            var d = new Date(movie.Date);
            return d.getFullYear();
        };

        $scope.displayDate = function(dateString) {
            var d = new Date(dateString);
            return d.toUTCString();
        }

        self.refresh = function() {
            $scope.config = $rootScope.config;
            $scope.movies = $rootScope.movies;
        };

        self.refresh();
        $interval(self.refresh, 1000);

        return;
    }
});
