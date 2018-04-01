flemzerd.component("tvshows", {
    templateUrl: "/static/templates/tvshows.template.html",
    controller: function TVShowsCtrl($scope, $http) {
        var self = this;
        $scope.tvshows = {};
        $scope.config = {};
        $scope.tvshows_found = false;

        self.LoadConfig = function() {
            $http.get("/api/v1/config").then(function(response) {
                if (response.status == 200) {
                    $scope.config = response.data;
                }
            });
        };
        self.LoadConfig();

        $http.get("/api/v1/tvshows").then(function(response) {
            if (response.status == 200) {
                shows = response.data;
                for (var show in shows) {
                    switch (shows[show]["Status"]) {
                        case 1:
                            shows[show]["StatusText"] = "Continuing";
                            break;
                        case 2:
                            shows[show]["StatusText"] = "Planned";
                            break;
                        case 3:
                            shows[show]["StatusText"] = "Ended";
                            break;
                        default:
                            shows[show]["StatusText"] = "Unknown";
                            break;
                    }
                }

                if (shows != null) {
                     $scope.tvshows_found = true;
                }

                $scope.tvshows = shows;
                return;
            }
        });
    }
});
