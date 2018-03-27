var flemzerd = angular.module('flemzer', ['ui.router'])

flemzerd.config(function config($stateProvider) {
    $stateProvider.state({
        name: "tvshows",
        url: "/tvshows",
        component: "tvshows"
    });
    $stateProvider.state({
        name: "movies",
        url: "/movies",
        component: "movies"
    });
    $stateProvider.state({
        name: "status",
        url: "/status",
        component: "status"
    });
    $stateProvider.state({
        name: "settings",
        url: "/settings",
        component: "settings"
    });
});
