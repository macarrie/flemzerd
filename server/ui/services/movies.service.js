flemzerd.factory('movies', ['$rootScope', '$interval', '$http', function($rootScope, $interval, $http) {
    var loadMovies = function() {
        $http.get("/api/v1/movies/tracked").then(function(response) {
            if (response.status == 200) {
                $rootScope.movies.tracked = response.data;
                return;
            }
        });

        $http.get("/api/v1/movies/downloading").then(function(response) {
            if (response.status == 200) {
                $rootScope.movies.downloading = response.data;
                return;
            }
        });

        $http.get("/api/v1/movies/downloaded").then(function(response) {
            if (response.status == 200) {
                $rootScope.movies.downloaded = response.data;
                return;
            }
        });
    };

    $rootScope.movies = {}
    loadMovies();
    $interval(loadMovies, 30000);

    return {
        "movies": $rootScope.movies,
    }
}]);
