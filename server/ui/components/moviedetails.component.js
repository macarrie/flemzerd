flemzerd.component("moviedetails", {
    templateUrl: "/static/templates/moviedetails.template.html",
    controller: function MovieDetailsCtrl($rootScope, $scope, $http, $stateParams) {
        $scope.movie = {};

        $http.get("/api/v1/movies/details/" + $stateParams.id).then(function(response) {
            if (response.status == 200) {
                $scope.movie = response.data;
                return;
            }
            // TODO: Check 404
        });

        return;
    }
});
