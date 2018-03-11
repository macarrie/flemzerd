angular.module('flemzer')
    .controller('HomeCtrl', HomeCtrl);

function HomeCtrl($routeParams, $http) {
    $http.get("https://api.trakt.tv/oauth/device/code").then(function(response) {
        console.log(response);
    });

    return;
}
