angular.module('tvshows')
    .component('tvshows', {
        templateUrl: 'static/tvshows/tvshows.template.html',
        controller: ['$http', function TVShowsController($http) {
            var self = this;
            $http.get("/api/v1/tvshows").then(function(response) {
                console.log(response)
                if (response.status == 200) {
                    self.tvshows = response.data;
                    return;
                }
            });
        }]
    })
