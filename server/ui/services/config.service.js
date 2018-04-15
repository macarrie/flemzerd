flemzerd.factory('config', ['$rootScope', '$interval', '$http', function($rootScope, $interval, $http) {
    var loadConfig = function() {
        $http.get("/api/v1/config").then(function(response) {
            if (response.status == 200) {
                $rootScope.config = response.data;
            }
        });
    };

    // start periodic checking
    loadConfig();
    $interval(loadConfig, 30000);

    return {
        "config": $rootScope.config
    }
}]);
