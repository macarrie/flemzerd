flemzerd.component("moviedetails", {
    templateUrl: "/static/templates/moviedetails.template.html",
    controller: function MovieDetailsCtrl($rootScope, $scope, $http, $state, $stateParams, movies, fanart) {
        $scope.utils = $rootScope.utils;
        $scope.movie = {};

        var getMovie = function() {
            $http.get("/api/v1/movies/details/" + $stateParams.id).then(function(response) {
                if (response.status == 200) {
                    $scope.movie = response.data;
                    fanart.GetMovieFanart($scope.movie).then(function(data) {
                        $scope.fanarts = data;
                        background_container = document.getElementById("full_background")
                        background_container.style.backgroundImage = "url('" + data.moviebackground[0].url + "')";
                    });
                    return;
                }
                // TODO: Check 404
            });
        };

        getMovie();

        $scope.deleteMovie = function(id) {
            movies.deleteMovie(id);
            $state.go("movies")
        };

        $scope.restoreMovie = function(id) {
            movies.restoreMovie(id);
            getMovie();
        };

        $scope.changeDownloadedState = function(m, downloaded) {
            ret = movies.changeDownloadedState(m, downloaded);
            if (ret != false) {
                ret.then(function(response) {
                    if (response.status == 200) {
                        $scope.movie = response.data;
                    }
                });
            }
        };

        return;
    }
});
