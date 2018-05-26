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

        $http.get("/api/v1/movies/removed").then(function(response) {
            if (response.status == 200) {
                $rootScope.movies.removed = response.data;
                return;
            }
        });
    };

    var deleteMovie = function(id) {
        if (id !== 0) {
            $http.delete("/api/v1/movies/details/" +id).then(function(response) {
                if (response.status == 204) {
                    loadMovies();
                    return true;
                }
            });
        }
        return false;
    };

    var restoreMovie = function(id) {
        if (id !== 0) {
            $http.post("/api/v1/movies/restore/" +id).then(function(response) {
                if (response.status == 200) {
                    loadMovies();
                    return true;
                }
            });
        }
        return false;
    };

    var changeDownloadedState = function(movie, downloaded) {
        id = movie.ID;
        if (id !== 0) {
            return $http.put("/api/v1/movies/details/" +id, {
                "Downloaded": downloaded
            }).then(function(response) {
                return response;
            });
        }
        return false;
    };

    $rootScope.movies = {}
    loadMovies();
    $interval(loadMovies, 30000);

    return {
        loadMovies: loadMovies,
        deleteMovie: deleteMovie,
        restoreMovie: restoreMovie,
        changeDownloadedState: changeDownloadedState
    }
}]);
