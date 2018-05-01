flemzerd.component("tvshowdetails", {
    templateUrl: "/static/templates/tvshowdetails.template.html",
    controller: function TvShowDetails($rootScope, $scope, $http, $stateParams, fanart) {
        $scope.utils = $rootScope.utils;
        $scope.movie = {};

        $http.get("/api/v1/tvshows/details/" + $stateParams.id).then(function(response) {
            if (response.status == 200) {
                $scope.tvshow = response.data;
                fanart.GetTvShowFanart($scope.tvshow).then(function(data) {
                    $scope.fanarts = data;
                    console.log("Fanarts: ", data);
                    background_container = document.getElementById("full_background")
                    background_container.style.backgroundImage = "url('" + data.showbackground[0].url + "')";
                });
                return;
            }
            // TODO: Check 404
        });

        return;
    }
});
