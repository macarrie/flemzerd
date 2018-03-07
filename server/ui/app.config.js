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
            .when('/status', {
                templateUrl: 'static/status/status.template.html',
                controller: "StatusCtrl",
                controllerAs: "ctrl"
            })
            .otherwise('/');
        }
    ]);
