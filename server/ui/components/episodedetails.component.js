flemzerd.component("episodedetails", {
    templateUrl: "/static/templates/episodedetails.template.html",
    controller: function EpisodeDetails($rootScope, $scope, $http, $stateParams, tvshows, fanart) {
        $scope.utils = $rootScope.utils;

        var getEpisode = function() {
        $http.get("/api/v1/tvshows/episodes/" + $stateParams.id).then(function(response) {
            if (response.status == 200) {
                $scope.episode = response.data;
                $scope.episode.TvShow = tvshows.setStatusText($scope.episode.TvShow);
                fanart.GetTvShowFanart($scope.episode.TvShow).then(function(data) {
                    $scope.fanarts = data;
                    background_container = document.getElementById("full_background")
                    background_container.style.backgroundImage = "url('" + data.showbackground[0].url + "')";
                });
                return;
            }
            // TODO: Check 404
        });
        };

        getEpisode();

        $scope.deleteEpisode = function(id) {
            tvshows.deleteEpisode(id);
        };

        return;
    }
});
