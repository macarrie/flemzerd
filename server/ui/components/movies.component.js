flemzerd.component("movies", {
    templateUrl: "/static/templates/movies.template.html",
    controller: function MoviesCtrl($scope, $http, $timeout) {
        var self = this;
        $scope.config = {};
        $scope.movies = {};
        $scope.movies_found = false;

        self.LoadConfig = function() {
            $http.get("/api/v1/config").then(function(response) {
                if (response.status == 200) {
                    $scope.config = response.data;
                }
            });
        };
        self.LoadConfig();

        $http.get("/api/v1/movies").then(function(response) {
            if (response.status == 200) {
                $scope.movies = response.data;

                if ($scope.movies != null) {
                    $scope.movies_found = true;
                }

                return;
            }
        });
    }
});
