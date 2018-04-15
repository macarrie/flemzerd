flemzerd.component("home", {
    templateUrl: "/static/templates/home.template.html",
    controller: function HomeCtrl($rootScope, $scope, $http, $interval, config, movies, tvshows) {
        var self = this;
        $scope.config = {};
        $scope.movies = {};

        self.refresh = function() {
            $scope.config = $rootScope.config;
            $scope.movies = $rootScope.movies;
            $scope.tvshows = $rootScope.tvshows;
        };

        self.refresh();
        $interval(self.refresh, 1000);

        return;
    }
});
