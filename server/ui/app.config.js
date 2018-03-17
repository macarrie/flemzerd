angular.
    module('flemzer').
    config(['$routeProvider',
        function config($routeProvider) {
            $routeProvider.when('/', {
                templateUrl: 'static/home/home.template.html',
                controller: "HomeCtrl"
            })
            .when('/tvshows', {
                templateUrl: 'static/tvshows/tvshows.template.html',
                controller: "TVShowsCtrl",
                controllerAs: "ctrl"
            })
            .when('/movies', {
                templateUrl: 'static/movies/movies.template.html',
                controller: "MoviesCtrl",
                controllerAs: "ctrl"
            })
            .when('/status', {
                templateUrl: 'static/status/status.template.html',
                controller: "StatusCtrl",
                controllerAs: "ctrl"
            })
            .when('/settings', {
                templateUrl: 'static/settings/settings.template.html',
                controller: "SettingsCtrl",
                controllerAs: "ctrl"
            })
            .otherwise('/');
        }
    ]);
