var flemzerd = angular.module('flemzer', ['ui.router'])

flemzerd.config(function config($stateProvider) {
    $stateProvider.state({
        name: "home",
        url: "/",
        component: "home"
    });
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
        name: "moviedetails",
        url: "/movies/{id}",
        component: "moviedetails"
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
