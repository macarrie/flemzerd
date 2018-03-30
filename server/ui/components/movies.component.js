flemzerd.component("movies", {
    templateUrl: "/static/templates/movies.template.html",
    controller: function MoviesCtrl($http) {
        var self = this;
        self.config = {};
        self.movies = {};

        self.LoadConfig = function() {
            $http.get("/api/v1/config").then(function(response) {
                if (response.status == 200) {
                    self.config = response.data;
                }
            });
        };
        self.LoadConfig();

        self.GetMovieCount = function() {
            return Object.keys(self.movies).length;
        }

        $http.get("/api/v1/movies").then(function(response) {
            if (response.status == 200) {
                self.movies = response.data;
                return;
            }
        });
    }
});
