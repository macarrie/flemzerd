flemzerd.component("tvshowdetails", {
    templateUrl: "/static/templates/tvshowdetails.template.html",
    controller: function TvShowDetails($rootScope, $scope, $http, $state, $stateParams, tvshows, fanart) {
        $scope.utils = $rootScope.utils;

        var getShow = function() {
            $http.get("/api/v1/tvshows/details/" + $stateParams.id).then(function(response) {
                if (response.status == 200) {
                    $scope.tvshow = tvshows.setStatusText(response.data);
                    fanart.GetTvShowFanart($scope.tvshow).then(function(data) {
                        $scope.fanarts = data;
                        background_container = document.getElementById("full_background")
                        background_container.style.backgroundImage = "url('" + data.showbackground[0].url + "')";
                    });
                    return;
                }
                // TODO: Check 404
            });
        };

        getShow();

        $scope.deleteShow = function(id) {
            tvshows.deleteShow(id);
            $state.go("tvshows")
        };

        $scope.restoreShow = function(id) {
            tvshows.restoreShow(id);
            getShow();
        };

        return;
    }
});
