flemzerd.component("tvshows", {
    templateUrl: "/static/templates/tvshows.template.html",
    controller: function TVShowsCtrl($rootScope, $scope, $interval, $http, config, tvshows) {
        var self = this;
        $scope.utils = $rootScope.utils;
        $scope.config = {};
        $scope.tvshows = {};

        self.refresh = function() {
            $scope.config = $rootScope.config;
            $scope.tvshows = $rootScope.tvshows;
        };

        $scope.refreshTvShows = function() {
            $scope.refreshing = true;
            $http.post("/api/v1/modules/watchlists/refresh").then(function(response) {
                if (response.status == 200) {
                    tvshows.loadShows();
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
