flemzerd.component("tvshows", {
    templateUrl: "/static/templates/tvshows.template.html",
    controller: function TVShowsCtrl($http) {
        var self = this;
        self.tvshows = {};

        self.LoadConfig = function() {
            $http.get("/api/v1/config").then(function(response) {
                if (response.status == 200) {
                    self.config = response.data;
                }
            });
        };
        self.LoadConfig();

        self.GetShowCount = function() {
            return Object.keys(self.tvshows).length;
        }

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
                    console.log(shows[show]);
                }
                self.tvshows = shows;
                return;
            }
        });
    }
});
