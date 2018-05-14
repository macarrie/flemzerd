flemzerd.component("tvshowdetails", {
    templateUrl: "/static/templates/tvshowdetails.template.html",
    controller: function TvShowDetails($rootScope, $scope, $http, $state, $stateParams, tvshows, fanart) {
        $scope.utils = $rootScope.utils;

        var getShow = function() {
            $http.get("/api/v1/tvshows/details/" + $stateParams.id).then(function(response) {
                if (response.status == 200) {
                    $scope.tvshow = tvshows.setStatusText(response.data);
                    $scope.seasons = [];

                    fanart.GetTvShowFanart($scope.tvshow).then(function(data) {
                        $scope.fanarts = data;
                        background_container = document.getElementById("full_background")
                        background_container.style.backgroundImage = "url('" + data.showbackground[0].url + "')";
                    });

                    for (var season in response.data.Seasons) {
                        s = response.data.Seasons[season];
                        getEpisodeList(s.SeasonNumber);
                    }
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

        var getEpisodeList = function(seasonNb) {
            $http.get("/api/v1/tvshows/details/" + $stateParams.id + "/seasons/" +seasonNb).then(function(response) {
                if (response.status == 200) {
                    $scope.seasons[seasonNb] = response.data;
                    return response.data;
                }
            });
        };

        $scope.downloadEpisode = function(episode) {
            $http.post("/api/v1/tvshows/episodes/" + episode.ID + "/download").then(function(response) {
                if (response.status == 200) {
                    getEpisodeList(episode.Season);
                    return;
                }
            });
        };

        return;
    }
});
