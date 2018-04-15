flemzerd.component("tvshows", {
    templateUrl: "/static/templates/tvshows.template.html",
    controller: function TVShowsCtrl($rootScope, $scope, $interval, $http, config, tvshows) {
        var self = this;
        $scope.config = {};
        $scope.tvshows = {};

        self.refresh = function() {
            $scope.config = $rootScope.config;
            $scope.tvshows = $rootScope.tvshows;
        };

        self.refresh();
        $interval(self.refresh, 1000);

        return;
    }
});
