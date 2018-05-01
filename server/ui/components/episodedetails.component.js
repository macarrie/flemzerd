flemzerd.component("episodedetails", {
    templateUrl: "/static/templates/moviedetails.template.html",
    controller: function MovieDetailsCtrl($rootScope, $scope, $http, $stateParams, fanart) {
        $scope.utils = $rootScope.utils;
        $scope.movie = {};

        $http.get("/api/v1/movies/details/" + $stateParams.id).then(function(response) {
            if (response.status == 200) {
                $scope.movie = response.data;
                fanart.GetMovieFanart($scope.movie).then(function(data) {
                    $scope.fanarts = data;
                    background_container = document.getElementById("movie_full_background")
                    background_container.style.backgroundImage = "url('" + data.moviebackground[0].url + "')";
                });
                return;
            }
            // TODO: Check 404
        });

        return;
    }
});
